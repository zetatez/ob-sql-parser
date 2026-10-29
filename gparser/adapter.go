package gparser

import (
	"strings"

	"github.com/zetatez/ob-sql-parser/ast"
	"github.com/zetatez/ob-sql-parser/feature"
	"github.com/zetatez/ob-sql-parser/token"
	"github.com/zetatez/ob-sql-parser/version"
)

// Options configures the goyacc parser.
type Options struct {
	// Version is the target OceanBase version (zero means latest / all
	// features enabled).
	Version version.Version
	// Features is the feature registry used for version gating.
	Features *feature.Registry
}

// obLexer adapts a pre-lexed token slice to the goyacc yyLexer interface. It
// stops (returns EOF) at a top-level semicolon so the driver can parse one
// statement at a time.
type obLexer struct {
	toks        []token.Token
	i           int
	version     version.Version
	features    *feature.Registry
	err         error
	last        token.Token
	result      ast.Statement
	param       int
	inGrant     bool
	atStart     bool
	pendingUser bool
	justEnded   bool // a statement just ended at ';' ; snapshot next token start
	stmtStart   token.Position
	results     []ast.Statement
}

func newLexer(toks []token.Token, opts Options) *obLexer {
	f := opts.Features
	if f == nil {
		f = feature.DefaultRegistry()
	}
	clean := make([]token.Token, 0, len(toks))
	// A `/*+ ... */` comment is a hint. Only LOAD DATA statements use them;
	// elsewhere hints are ordinary comments and are dropped. To decide, split
	// into statement segments at SEMI and check each segment's first keyword.
	atSegStart := true
	inLoad := false
	for _, t := range toks {
		switch t.Type {
		case token.EOF:
			continue
		case token.SEMI:
			inLoad = false
			atSegStart = true
			clean = append(clean, t)
			continue
		case token.COMMENT:
			if inLoad && strings.HasPrefix(t.Literal, "/*+") {
				clean = append(clean, t)
			}
			continue
		}
		if t.Type == token.IDENT && atSegStart {
			atSegStart = false
			if strings.EqualFold(t.Value, "load") {
				inLoad = true
			}
		}
		clean = append(clean, t)
	}
	return &obLexer{toks: clean, version: opts.Version, features: f, justEnded: true, atStart: true}
}

// Error implements yyLexer. The generated parser calls it on a syntax error.
func (l *obLexer) Error(s string) {
	if l.err == nil {
		l.err = &Error{Line: l.last.Pos.Line, Col: l.last.Pos.Col, Msg: s}
	}
}

// Lex implements yyLexer, returning the next terminal code (0 == EOF).
func (l *obLexer) Lex(lval *yySymType) int {
	if l.i >= len(l.toks) {
		return 0
	}
	t := l.toks[l.i]
	l.last = t
	if t.Type == token.SEMI {
		l.i++
		l.inGrant = false
		l.pendingUser = false
		l.justEnded = true
		l.atStart = true
		return ';'
	}
	if l.justEnded {
		l.stmtStart = t.Pos
		l.justEnded = false
	}
	first := l.atStart
	l.atStart = false
	l.i++
	term := l.mapToken(t, lval)
	// "analyze" is non-reserved: only a statement-initial ANALYZE is the table
	// maintenance keyword; elsewhere (EXPLAIN / identifiers) it stays NAME_OB.
	if first && t.Type == token.IDENT && strings.EqualFold(t.Value, "analyze") {
		term = ANALYZE
	}
	if first && (term == GRANT || term == REVOKE || term == AUDIT || term == NOAUDIT) {
		l.inGrant = true
	}
	// USER/ROLE are context-sensitive: only the word right after a leading
	// CREATE/DROP is the keyword, so "user" stays usable as an identifier.
	if first && (term == CREATE || term == DROP || term == SET || term == LOAD || term == ALTER ||
		term == RENAME || term == LOCK || term == UNLOCK) {
		l.pendingUser = true
	} else if l.pendingUser {
		keep := false
		if t.Type == token.IDENT {
			switch strings.ToLower(t.Value) {
			case "user":
				term = USER
			case "role":
				term = ROLE
			case "password":
				term = PASSWORD
			case "data":
				term = DATA
			case "system":
				term = SYSTEM
			case "log":
				term = LOG
			case "materialized", "view":
				keep = true
			}
		}
		if !keep {
			l.pendingUser = false
		}
	}
	return term
}

func (l *obLexer) mapToken(t token.Token, lval *yySymType) int {
	switch t.Type {
	case token.IDENT:
		word := t.Value
		low := strings.ToLower(word)
		// Inside a GRANT/REVOKE privilege list, most words are privilege names
		// (SELECT, INSERT, CREATE, ...) and must be treated as identifiers until
		// the ON keyword ends the list.
		if l.inGrant {
			switch low {
			case "all":
				return ALL
			case "privileges":
				return PRIVILEGES
			case "on":
				l.inGrant = false
				return ON
			default:
				lval.node = &ast.Identifier{Parts: []string{word}}
				return NAME_OB
			}
		}
		switch low {
		case "true", "false":
			lval.node = &ast.Literal{LitType: ast.LitBool, Value: low, Text: word}
			return BOOL_VALUE
		case "null":
			lval.node = &ast.Literal{LitType: ast.LitNull, Value: "null", Text: word}
			return NULLX
		}
		if term, ok := keywordTerminal[low]; ok {
			return term
		}
		lval.node = &ast.Identifier{Parts: []string{word}}
		return NAME_OB
	case token.STRING:
		if i := strings.IndexByte(t.Value, 0); i >= 0 {
			lval.node = &ast.Literal{LitType: ast.LitString, Value: t.Value[i+1:], Text: t.Literal, Charset: t.Value[:i]}
		} else {
			lval.node = &ast.Literal{LitType: ast.LitString, Value: t.Value, Text: t.Literal}
		}
		return STRING_VALUE
	case token.INT:
		lval.node = &ast.Literal{LitType: ast.LitInt, Value: t.Value, Text: t.Literal}
		return INTNUM
	case token.DECIMAL:
		lval.node = &ast.Literal{LitType: ast.LitDecimal, Value: t.Value, Text: t.Literal}
		return DECIMAL_VAL
	case token.FLOAT:
		lval.node = &ast.Literal{LitType: ast.LitFloat, Value: t.Value, Text: t.Literal}
		return APPROXNUM
	case token.HEXNUM:
		lval.node = &ast.Literal{LitType: ast.LitHex, Value: t.Value, Text: t.Literal}
		return HEXNUM
	case token.BITNUM:
		lval.node = &ast.Literal{LitType: ast.LitBit, Value: t.Value, Text: t.Literal}
		return BITNUM
	case token.HEXSTRING:
		lval.node = &ast.Literal{LitType: ast.LitHex, Value: t.Value, Text: t.Literal}
		return HEX_STRING_VALUE
	case token.BITSTRING:
		lval.node = &ast.Literal{LitType: ast.LitBit, Value: t.Value, Text: t.Literal}
		return BIN_STRING_VALUE
	case token.QUESTION:
		l.param++
		lval.node = &ast.ParamMarker{Index: l.param}
		return QUESTIONMARK
	case token.VARIABLE, token.SYSVARIABLE:
		v := variableFromToken(t)
		lval.node = v
		if t.Type == token.SYSVARIABLE {
			return SYSTEM_VARIABLE
		}
		// A bare '@' with no variable name is the user@host separator used by
		// GRANT/REVOKE and similar statements.
		if v.Name == "" && v.Scope == "" {
			return AT
		}
		return USER_VARIABLE
	case token.COMMA:
		return ','
	case token.DOT:
		return '.'
	case token.LPAREN:
		return '('
	case token.RPAREN:
		return ')'
	case token.PLUS:
		return '+'
	case token.MINUS:
		return '-'
	case token.STAR:
		return '*'
	case token.SLASH:
		return '/'
	case token.PERCENT:
		return '%'
	case token.POW:
		return '^'
	case token.AMP:
		return '&'
	case token.BITOR:
		return '|'
	case token.TILDE:
		return '~'
	case token.EQ, token.EQEQ:
		return COMP_EQ
	case token.NE:
		return COMP_NE
	case token.LT:
		return COMP_LT
	case token.GT:
		return COMP_GT
	case token.LE:
		return COMP_LE
	case token.GE:
		return COMP_GE
	case token.NULLSAFE_EQ:
		return NSEQ
	case token.ASSIGN:
		// SET @x := 1 (MySQL legacy) is canonicalized to '='
		return COMP_EQ
	case token.SHL:
		return SHIFT_LEFT
	case token.SHR:
		return SHIFT_RIGHT
	case token.OTHER:
		switch t.Literal {
		case "||":
			return OROR
		case "&&":
			return ANDAND
		case ":":
			return ':'
		}
	case token.COMMENT:
		// Only /*+ ... */ hints survive newLexer, and only inside LOAD DATA.
		lval.node = &ast.Literal{LitType: ast.LitString, Text: t.Literal, Value: t.Literal}
		return DIRECT_HINT
	}
	l.err = &Error{Line: t.Pos.Line, Col: t.Pos.Col, Msg: "unexpected token " + t.Text()}
	return ILLEGAL_TOK
}

// keywordTerminal maps a normalized keyword to its grammar terminal. Words not
// present here are treated as ordinary identifiers (NAME_OB), which lets
// non-reserved keywords (INT, VARCHAR, ...) serve as names.
var keywordTerminal = map[string]int{
	"select": SELECT, "insert": INSERT, "replace": REPLACE, "update": UPDATE,
	"delete": DELETE, "create": CREATE, "drop": DROP, "truncate": TRUNCATE,
	"show": SHOW, "set": SET, "use": USE, "explain": EXPLAIN,
	"desc": DESC, "describe": DESCRIBE, "begin": BEGIN, "start": START,
	"commit": COMMIT, "rollback": ROLLBACK, "with": WITH, "values": VALUES,
	"value": VALUE, "into": INTO, "from": FROM, "where": WHERE,
	"group": GROUP, "by": BY, "having": HAVING, "order": ORDER,
	"limit": LIMIT, "offset": OFFSET, "union": UNION, "except": EXCEPT,
	"intersect": INTERSECT, "minus": MINUS, "all": ALL, "distinct": DISTINCT,
	"as": AS, "and": AND, "or": OR, "xor": XOR, "not": NOT,
	"is": IS, "like": LIKE, "rlike": RLIKE, "regexp": REGEXP, "in": IN,
	"between": BETWEEN, "case": CASE, "when": WHEN, "then": THEN,
	"else": ELSE, "end": END, "exists": EXISTS, "interval": INTERVAL,
	"asc": ASC, "nulls": NULLS, "first": FIRST, "last": LAST,
	"over": OVER, "partition": PARTITION, "rows": ROWS, "range": RANGE,
	"current": CURRENT, "row": ROW, "unbounded": UNBOUNDED,
	"preceding": PRECEDING, "following": FOLLOWING, "join": JOIN,
	"inner": INNER, "outer": OUTER, "left": LEFT, "right": RIGHT,
	"full": FULL, "cross": CROSS, "natural": NATURAL,
	"straight_join": STRAIGHT_JOIN, "on": ON, "using": USING, "dual": DUAL,
	"table": TABLE, "database": DATABASE, "schema": SCHEMA, "index": INDEX,
	"key": KEY, "primary": PRIMARY, "unique": UNIQUE, "foreign": FOREIGN,
	"references": REFERENCES, "constraint": CONSTRAINT, "check": CHECK,
	"if": IF, "temporary": TEMPORARY, "ignore": IGNORE,
	"low_priority": LOW_PRIORITY, "high_priority": HIGH_PRIORITY,
	"delayed": DELAYED, "duplicate": DUPLICATE, "returning": RETURNING,
	"global": GLOBAL, "session": SESSION, "local": LOCAL,
	"databases": DATABASES, "schemas": SCHEMAS, "tables": TABLES,
	"variables": VARIABLES, "status": STATUS, "warnings": WARNINGS,
	"errors": ERRORS, "grants": GRANTS, "processlist": PROCESSLIST,
	"engines": ENGINES, "storage": STORAGE, "collation": COLLATION,
	"columns": COLUMNS, "fields": FIELDS, "indexes": INDEXES, "keys": KEYS,
	"view": VIEW, "triggers": TRIGGERS, "procedure": PROCEDURE,
	"function": FUNCTION, "names": NAMES, "character": CHARACTER,
	"transaction": TRANSACTION, "work": WORK, "read": READ, "write": WRITE,
	"only": ONLY, "chain": CHAIN, "extended": EXTENDED, "pretty": PRETTY,
	"format": FORMAT, "cascade": CASCADE, "restrict": RESTRICT, "quick": QUICK,
	"unsigned": UNSIGNED, "default": DEFAULT, "comment": COMMENT,
	"auto_increment": AUTO_INCREMENT, "charset": CHARSET, "collate": COLLATE,
	"invisible": INVISIBLE, "cast": CAST, "for": FOR, "share": SHARE,
	"key_block_size": KEY_BLOCK_SIZE, "block_size": BLOCK_SIZE, "action": ACTION, "no": NO,
	"enable": ENABLE, "disable": DISABLE, "zerofill": ZEROFILL, "visible": VISIBLE,
	"window": WINDOW, "recursive": RECURSIVE, "escape": ESCAPE,
	"div": DIV, "mod": MOD, "sounds": SOUNDS,
	"alter": ALTER, "add": ADD, "modify": MODIFY, "change": CHANGE,
	"rename": RENAME, "column": COLUMN, "to": TO, "force": FORCE,
	"rollup": ROLLUP, "list": LIST, "hash": HASH, "maxvalue": MAXVALUE,
	"less": LESS, "than": THAN, "partitions": PARTITIONS,
	"generated": GENERATED, "always": ALWAYS, "virtual": VIRTUAL, "stored": STORED,
	"sequence": SEQUENCE, "increment": INCREMENT, "cache": CACHE,
	"minvalue": MINVALUE, "cycle": CYCLE, "nocycle": NOCYCLE, "noorder": NOORDER,
	"isolation": ISOLATION, "level": LEVEL, "uncommitted": UNCOMMITTED,
	"committed": COMMITTED, "repeatable": REPEATABLE, "serializable": SERIALIZABLE,
	"savepoint": SAVEPOINT, "release": RELEASE, "lock": LOCK, "unlock": UNLOCK,
	"plugins": PLUGINS, "profiles": PROFILES, "profile": PROFILE, "privileges": PRIVILEGES,
	"trigger": TRIGGER, "open": OPEN,
	"grant": GRANT, "revoke": REVOKE, "option": OPTION,
	"identified": IDENTIFIED, "prepare": PREPARE, "execute": EXECUTE,
	"deallocate": DEALLOCATE, "kill": KILL, "help": HELP,
	"optimize": OPTIMIZE, "repair": REPAIR, "connection": CONNECTION,
	"query": QUERY, "none": NONE, "load": LOAD, "infile": INFILE,
	"terminated": TERMINATED, "lines": LINES, "flush": FLUSH, "export": EXPORT,
	"xa": XA, "resume": RESUME, "suspend": SUSPEND, "migrate": MIGRATE,
	"one": ONE, "phase": PHASE, "recover": RECOVER, "convert": CONVERT,
	"purge": PURGE, "recyclebin": RECYCLEBIN, "xid": XID,
	"member": MEMBER, "against": AGAINST, "match": MATCH, "grouping": GROUPING,
	"sets": SETS, "cube": CUBE, "boolean": BOOLEAN, "mode": MODE, "expansion": EXPANSION,
	"of": OF, "merge": MERGE, "matched": MATCHED, "flashback": FLASHBACK,
	"before": BEFORE, "snapshot": SNAPSHOT, "materialized": MATERIALIZED,
	"refresh": REFRESH, "fast": FAST, "complete": COMPLETE, "demand": DEMAND,
	"subpartition": SUBPARTITION, "subpartitions": SUBPARTITIONS,
	"signal": SIGNAL, "resignal": RESIGNAL, "sqlstate": SQLSTATE, "get": GET,
	"diagnostics": DIAGNOSTICS, "stacked": STACKED, "tenant": TENANT,
	"resource": RESOURCE, "unit": UNIT, "pool": POOL, "tablegroup": TABLEGROUP,
	"tablespace": TABLESPACE, "outline": OUTLINE, "master": MASTER,
	"rotate": ROTATE, "including": INCLUDING, "new": NEW,
	"standby": STANDBY, "activate": ACTIVATE, "switchover": SWITCHOVER,
	"until": UNTIL, "unlimited": UNLIMITED, "cancel": CANCEL, "scn": SCN,
	"restore": RESTORE, "point": POINT, "ls": LS, "time": TIME,
	"audit": AUDIT, "noaudit": NOAUDIT, "dump": DUMP, "entity": ENTITY,
	"chunk": CHUNK, "service": SERVICE, "remove": REMOVE, "files": FILES,
	"location": LOCATION, "url": URL, "routine": ROUTINE, "stream": STREAM,
	"pause": PAUSE, "stop": STOP, "concurrent_limiting_rule": CONCURRENT_LIMITING_RULE,
	"sensitive": SENSITIVE, "rule": RULE, "link": LINK,
	"connect": CONNECT, "host": HOST, "port": PORT, "cluster": CLUSTER,
	"lateral": LATERAL, "qualify": QUALIFY, "auto": AUTO, "size": SIZE,
	"plugin": PLUGIN, "soname": SONAME, "install": INSTALL, "uninstall": UNINSTALL,
	"handler": HANDLER, "method_opt": METHOD_OPT, "analyse": ANALYSE,
	"next": NEXT, "prev": PREV, "close": CLOSE, "reorganize": REORGANIZE,
	"after": AFTER, "coalesce": COALESCE, "exchange": EXCHANGE,
	"events": EVENTS, "mutex": MUTEX,
	"dumpfile": DUMPFILE, "logs": LOGS, "hosts": HOSTS, "outfile": OUTFILE,
	"condition": CONDITION, "external": EXTERNAL, "catalog": CATALOG,
	"properties": PROPERTIES, "returns": RETURNS, "declare": DECLARE,
	"cursor": CURSOR, "fetch": FETCH, "loop": LOOP, "while": WHILE,
	"repeat": REPEAT, "return": RETURN, "inout": INOUT,
	"elseif": ELSEIF, "call": CALL, "out": OUT, "do": DO, "transfer": TRANSFER,
	"continue": CONTINUE, "exit": EXIT, "undo": UNDO, "leave": LEAVE,
	"iterate": ITERATE, "found": FOUND, "without": WITHOUT, "locked": LOCKED,
	"validation": VALIDATION, "nowait": NOWAIT, "skip": SKIP,
}
