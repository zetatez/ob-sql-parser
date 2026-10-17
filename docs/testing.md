# 测试指南

> 本文回答"**怎么保证这个解析器是对的、怎么为它加测试**"。覆盖测试体系地图、
> 新增能力的标准流程、单点调试技巧与质量基线。
>
> 读者：维护者 / 想给本项目贡献的开发者。相关：[coverage](coverage.md) 的能力表即"活文档"。

## 1. 快速运行

```bash
go test ./...                  # 全部
go test -race ./...            # 竞态（解析器实例无共享可变状态）
go test ./... -cover           # 各包语句覆盖
go test ./gparser/ -run TestOfficialCorpus -v   # 官方语料（跑单文件）
go test ./gparser/ -run 'TestExtraRoundTrip/`select`' -v  # 单个子用例
```

## 2. 测试体系地图

| 测试文件                         | 主要用例                                                                                                                         | 职责                                            |
| ---                              | ---                                                                                                                              | ---                                             |
| `parse_test.go`                  | `TestRoundTrip`(35)、`ParseBasics`、`ParseMulti`、`ParseError`、`Comments`                                                       | 顶层 API 解析/幂等基础                          |
| `integration_test.go`            | `TestVersionGating`、`TestVisitorCollect`、`TestJSONSerialization`                                                               | 顶层集成（门控/遍历/JSON）                      |
| `ast/ast_test.go`                | `TestAllNodeKindsShape`                                                                                                          | 全 ~140 节点类型 Finalize/Children/Marshal 冒烟 |
| `formatter/formatter_test.go`    | 72 个表达式/语句精确输出断言 + 全部语句类型排版冒烟                                                                              | Formatter 输出正确性                            |
| `token/token_test.go`            | `LookupKeyword`/保留字/`Type.String`/Span/Position                                                                               | 关键字表与记号                                  |
| `visitor/visitor_test.go`        | 收集/剪枝/typed-nil 守卫/`WalkStatements`                                                                                        | 遍历框架                                        |
| `gparser/gparser_test.go`        | `TestRoundTrip`(35)、`TestExtraRoundTrip`(~150)、`TestFeatureGating`(~40 行边界)、`TestMulti`、`TestDoStmt`、`TestStatementSpan` | 主引擎 round-trip 与门控边界                    |
| `gparser/corpus_test.go`         | `TestOfficialCorpus`（官方 4.2.5 82 分支 × parse+round-trip+JSON 三重）                                                          | 语句面全量对照                                  |
| `gparser/coverage_test.go`       | `TestCoverageSupported`(~50)/`Unsupported`                                                                                       | 覆盖矩阵"活文档"断言                            |
| `gparser/version_corpus_test.go` | `TestVersionCorpus`（4.2.5 严格）                                                                                                | 版本门控回归                                    |
| `gparser/span_test.go`           | `TestExpressionSpans`                                                                                                            | 表达式级 Span 精确列号                          |
| `gparser/pretty_test.go`         | `TestPrettyFormat`                                                                                                               | 多行输出可再解析且幂等                          |
| `json/json_test.go`              | 45+ 语句 JSON 双向、`Dump`、`UnknownKind`、`MarshalIndent`、`Write`、错误路径                                                    | AST 与 JSON                                        |
| `lexer/lexer_test.go`            | 操作符/字面量/位置/引号标识符/注释/选项(ANSI_QUOTES)/错误/变量/引入符                                                            | 词法                                            |
| `feature/feature_test.go`        | `TestDefaultGating`、`TestCustomRegistry`                                                                                        | 门控矩阵/自定义表                               |
| `version/version_test.go`        | `ParseAndString`、`Compare`、`ParsePatchOrdering`                                                                                | 版本模型                                        |

## 3. 新增能力的标准流程

按顺序做，保证"能力 -> 测试 -> 文档 -> 门控"四同步：

1. **round-trip 先行**：在 `gparser.TestExtraRoundTrip` 或 `parse_test.TestRoundTrip`
   加一条 SQL。框架自动断言 `Parse -> Format -> Parse -> Format` 幂等。
2. **语义保留**：涉及子句顺序/别名时，先 `Parse -> Format` 人工确认输出与输入语义等价。
3. **版本门控**：有版本边界的语法，在 `TestFeatureGating` 加
   `(sql, 版本, 期望接受/拒绝)`；并在 `TestVersionCorpus` 归入 legal（4.2.5 合法）
   或 post425（4.2.5 拒绝）列表。`Since` 依据见 version-coverage 第 3 节。
4. **覆盖矩阵**：同步更新 `docs/coverage.md` 对应行（图例含门控标注），并迁移
   `TestCoverageSupported/Unsupported`。
5. **JSON**：新节点类型在 `json/registry.go`（kind 注册表）登记后自动生效；在
   `json.TestRoundTripJSON` 加一条代表 SQL。
6. **Span**：表达式/语句节点会被 `TestExpressionSpans`/`TestStatementSpan` 隐式覆盖；
   需精确列号时单独加断言。

## 4. 机制速览（写测试时有用）

- **整体一次解析**：每语句由 goyacc 一次 parse；`;` 是文法终结符（顶层分隔 +
  PL 体内分号），因此 PL 块可含任意多条语句。
- **上下文关键字**：`analyze`/`user`/`role`/`system` 等仅特定位置生效，仍可作
  标识符——有回归用例 `SELECT analyze FROM t`、`FROM user`。
- **门控=错误**：`gateFeat` 抛 `bailout`，驱动 `recover` 成带 `line:col` 的错误，
  因此 `WithVersion` 下可用 `err != nil` 断言"被拒绝"。
- **位置栈注入**：表达式级 Span 由 `genpos.go` 注入；**重新生成必须用
  `go generate`**，裸 `goyacc -o sql.go` 会丢注入（`grep reduceSpan gparser/sql.go`
  可快速校验，1 行即有、0 行即丢）。

## 5. 调试单个解析失败

```bash
# 快速跑单条：写进 tmp 目录最小程序
mkdir -p /tmp/opencode/check && cat > /tmp/opencode/check/main.go <<'EOF'
package main
import ("fmt"; ob "github.com/zetatez/ob-sql-parser")
func main(){
  st, e := ob.Parse("SELECT a FROM t WHERE a > 1")
  fmt.Printf("err=%v\n", e)
  if e == nil { fmt.Printf("-> %s\n", ob.Format(st)) }
}
EOF
cd /home/shiyi/git/ob-sql-parser && go run /tmp/opencode/check
```

排查顺序：① 词法（`lexer.LexAll` 看 token 流）-> ② 语法（`Parse` 裸跑，无版本）->
③ 门控（加 `WithVersion`）-> ④ round-trip（`Format` 再解析）。

## 6. 质量基线

```
go test -count=1 ./...   # 全绿
go test -race ./...      # 全绿
gofmt -l .               # 无输出
go vet ./...             # 无输出
staticcheck ./...        # 无告警（U1000 等已清零）
cd gparser && go generate # 可复现；grep -c reduceSpan sql.go == 1

冲突：24 shift/reduce、0 reduce/reduce（官方文法固有；默认移进优先 = 期望语义，
已由 round-trip/语料覆盖兜底）

语句覆盖（go test ./... -cover，各包自带测试口径）：
  token 100% · version 95.2% · feature 86% · gparser 84.1% · ast 79.3%
  visitor 73.7% · lexer 73.1% · 根包 65.4% · json 59.8% · formatter 48.4%
语句面（官方 4.2.5 的 91 个 stmt 分支 + 5.x 26 个）：100% 有用例。
说明：gparser 生成文件带 //line yaccpar 引用，无法用 -coverpkg 跨包合并；
上表为"各包自带测试"口径。
```

## 7. 建议的 CI

```yaml
# .github/workflows/ci.yml（参考）
steps:
  - uses: actions/checkout@v4
    with: { submodules: false }
  - uses: actions/setup-go@v5
    with: { go-version: "1.21" }
  - run: go install golang.org/x/tools/cmd/goyacc@latest
  - run: go generate ./gparser/
  - run: test -n "$(grep -c reduceSpan gparser/sql.go | grep -v '^0$')"
  - run: go test -count=1 ./...
  - run: go test -race ./...
  - run: go vet ./...
  - run: gofmt -l . | test -z "$(cat)"
```
