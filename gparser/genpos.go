//go:build ignore

// Command genpos regenerates gparser/sql.go from gparser/sql.y via goyacc and
// then injects a parallel position stack into the generated LALR parse loop,
// so every reduced node receives a source span (goyacc has no %locations).
//
// Run: go generate ./gparser/
package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if err := runGoyacc(); err != nil {
		panic(err)
	}
	b, err := os.ReadFile("sql.go")
	if err != nil {
		panic(err)
	}
	s := string(b)
	if err := applyEdits(&s, []edit{{
		anchor: "\tyyS := yyrcvr.stack[:]\n",
		insert: "\tyyPosTmp := make([]int, len(yyS))\n" +
			"\tyyPos := yyPosTmp\n" +
			"\tyyTok := -1\n" +
			"\tyyreducedStart := 0\n" +
			"\tyyreduced := false\n" +
			"\tyyFirst := true\n",
	}, {
		anchor: "\tyyS[yyp].yys = yystate\n",
		insert: "\tif yyp >= len(yyPos) {\n" +
			"\t\tnp := make([]int, len(yyPos)*2)\n" +
			"\t\tcopy(np, yyPos)\n" +
			"\t\tyyPos = np\n" +
			"\t}\n" +
			"\tif yyreduced {\n" +
			"\t\tyyreduced = false\n" +
			"\t\tyyPos[yyp] = yyreducedStart\n" +
			"\t} else if yyFirst {\n" +
			"\t\tyyFirst = false\n" +
			"\t\tyyPos[yyp] = -1 // initial dummy frame\n" +
			"\t} else {\n" +
			"\t\tyyTok++\n" +
			"\t\tyyPos[yyp] = yyTok\n" +
			"\t}\n",
	}, {
		anchor: "\tgoto yystack /* stack new state and value */\n",
		before: true,
		insert: "\tyyR := int(yyR2[yynt])\n" +
			"\tyyreducedStart = reduceSpan(yylex, yyVAL.node, yyPos, yypt, yyR, yyTok)\n" +
			"\tyyreduced = true\n" +
			"\n",
	}}); err != nil {
		panic(err)
	}
	if err := os.WriteFile("sql.go", []byte(s), 0o644); err != nil {
		panic(err)
	}
}

type edit struct {
	anchor string
	insert string
	before bool // insert before the anchor instead of after
}

func applyEdits(s *string, edits []edit) error {
	for _, e := range edits {
		i := strings.Index(*s, e.anchor)
		if i < 0 {
			return &errNotFound{anchor: e.anchor}
		}
		if e.before {
			*s = (*s)[:i] + e.insert + (*s)[i:]
		} else {
			*s = (*s)[:i+len(e.anchor)] + e.insert + (*s)[i+len(e.anchor):]
		}
	}
	return nil
}

type errNotFound struct{ anchor string }

func (e *errNotFound) Error() string { return "anchor not found: " + e.anchor }

func runGoyacc() error {
	cmd := exec.Command("goyacc", "-o", "sql.go", "sql.y")
	buf := &bytes.Buffer{}
	cmd.Stderr = buf
	cmd.Stdout = buf
	if err := cmd.Run(); err != nil {
		return err
	}
	if buf.Len() > 0 {
		// conflicts are warnings, not errors
		os.Stderr.WriteString(buf.String())
	}
	return nil
}
