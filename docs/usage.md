# 使用参考（API）

> 本文回答"**怎么用这个库**"。从零开始到能跑的最小程序、每个 API 的语义、
> 错误模型与常见问题。
>
> 读者：集成方。相关文档：[coverage](coverage.md) 能力清单、[version-coverage](version-coverage.md)
> 版本语义、[testing](testing.md) 如何自测。

## 0. 一个能直接跑的最小程序

```go
package main

import (
	"fmt"

	obparser "github.com/zetatez/ob-sql-parser"
)

func main() {
	st, err := obparser.Parse("SELECT a, b FROM t WHERE a > 1")
	if err != nil {
		panic(err)
	}
	fmt.Println(obparser.Format(st)) // SELECT a, b FROM t WHERE a > 1
}
```

```bash
go get github.com/zetatez/ob-sql-parser
go run .
```

## 1. 入口 API（根包 `obparser`）

```go
Parse(sql string, opts ...Option) (ast.Statement, error)     // 单语句
ParseMulti(sql string, opts ...Option) ([]ast.Statement, error) // 多语句（按 ; 分隔）
MustParse(sql string, opts ...Option) ast.Statement          // 出错即 panic（测试/脚本用）
Format(n ast.Node) string                                     // AST -> 规范 SQL
```

约定：
- 单语句也**允许**末尾 `;`（文法消费；多语句按 `;` 切分）。
- `ParseMulti("")` / 全空输入 -> 空列表，**不是错误**。
- 返回的 `ast.Statement` 已 `Finalize`（每个节点带 `Kind` 判别符）并注入语句级与
  表达式级 `Span`。

## 2. 版本门控（本库的核心差异化）

```go
obparser.Parse(sql)                                              // 默认：所有特性放行（按最新）
obparser.Parse(sql, obparser.WithVersion(version.MustParse("4.2.5"))) // 严格限定到 4.2.5
obparser.Parse(sql, obparser.WithFeatures(customRegistry))      // 自定义特性表
```

- `WithVersion`：目标 OB 版本。**严格**语义——4.2.5 之后引入的特性会被拒绝。
- `WithFeatures`：替换默认 `feature.DefaultRegistry`（用于裁剪/扩展门控）。
- 二者可同时传；**后传者生效**（Options 是列表）。

```go
// 4.3+ 的语法在 4.2.5 下必须报错
_, err := obparser.Parse("CREATE MATERIALIZED VIEW mv AS SELECT 1",
	obparser.WithVersion(version.MustParse("4.2.5")))
// err != nil，且错误信息带 line:col 与特性说明
```

版本语义与边界表见 [version-coverage](version-coverage.md)；机制见
[design](design.md) 第 7 节。

## 3. 错误模型

```go
type err struct{ ... } // 内部实现
```

- 所有解析/门控失败统一返回带 **`行:列`** 的错误，例如 `1:18: syntax error` 或
  `1:48: GROUPING SETS requires OceanBase 4.4+`（位置随输入而定，格式如 `line:col: message`）。
- **门控失败也是错误**（不是警告）；这使 `WithVersion` 可被测试直接断言"被拒绝"。
- `MustParse` 在任何错误下 panic（避免 `if err != nil` 噪音的快速路径）。

## 4. 源码位置（Span）

```go
st, _ := obparser.Parse("SELECT a + b FROM t WHERE a > 1")
fmt.Println(st.Span())                                 // 1:1 - 1:34
bodies := st.(*ast.SelectStmt).Bodies[0]
fmt.Println(bodies.Where.Span())                       // 表达式级 1:29 - 1:34
```

- `token.Span{Start, End token.Position}`，`Position{Offset, Line, Col}` 全 1-based。
- 半开区间 `[Start, End)`。
- 由 `gparser/genpos.go` 在生成期注入位置栈实现，**无需任何运行期配置**。
- 手工构造的 AST（未解析）位置为零值，`Position.IsValid()` 为 false——判断
  "有没有真实源码位置"用它。

## 5. 格式化

```go
obparser.Format(st)                                     // 紧凑规范输出
formatter.FormatWith(st, formatter.Options{Pretty: true}) // 多行缩进
```

- **round-trip 保证**：`Format` 输出可再解析，且 `Parse(Format(x))` 再 `Format`
  结果**不变**（幂等）。这是本项目"解析正确"的硬门槛，corpus 测试逐条断言。
- `Pretty` 只插入空白/换行，输出**仍是合法 SQL**，可再解析。
- 对"语义等价"（如别名、子句顺序）以规范化为准，不做保真切片。

## 6. 遍历（visitor）

```go
visitor.Walk(st, visitor.FuncVisitor(func(n ast.Node) visitor.Visitor {
	switch n.(type) {
	case *ast.TableRef:
		// 收集表引用
	}
	return visitor.NoOpVisitor{} // 继续下钻；返回 nil 则剪枝该子树
}))
```

- 回调返回新的 `visitor.Visitor` 可携带状态（下钻专用 visitor）。
- 空指针/typed-nil 子节点安全（手写 AST 也不会 panic）。
- `visitor.WalkStatements(stmts, v)` 遍历 `[]ast.Statement`。

## 7. JSON 双向

```go
data, _ := jsonpkg.Marshal(st)      // 紧凑；节点带 "kind" 判别符
data, _ = jsonpkg.MarshalIndent(st) // 缩进输出
back, _ := jsonpkg.Unmarshal(data)  // 重建 AST
out, _ := jsonpkg.Dump(stmts)       // 多语句列表（外层 StatementList）
nodes := back.(interface{ Nodes() []ast.Node }).Nodes()
```

- `kind` 判别符 = 节点类型名（`SelectStmt`/`BinaryExpr`…），由 `Finalize` 注入。
- 解码按 `json/registry.go` 的 kind 注册表反射重建；未知 kind 报错。
- 序列化结果与 `Format` 一样可用于持久化/传输。

## 8. 完整示例（版本 + 遍历 + JSON 组合）

```go
st, err := obparser.Parse(
	"SELECT t.a, COUNT(*) FROM t GROUP BY t.a HAVING COUNT(*) > 1",
	obparser.WithVersion(version.MustParse("4.2.5")))
if err != nil { panic(err) }

// 1) 规范化重建
println(obparser.Format(st))

// 2) 收集投影列名
var cols []string
visitor.Walk(st, visitor.FuncVisitor(func(n ast.Node) visitor.Visitor {
	if id, ok := n.(*ast.Identifier); ok {
		cols = append(cols, id.Parts[len(id.Parts)-1])
	}
	return visitor.NoOpVisitor{}
}))

// 3) JSON 持久化
data, _ := jsonpkg.Marshal(st)
back, _ := jsonpkg.Unmarshal(data)
if obparser.Format(back) != obparser.Format(st) {
	panic("JSON round-trip broken") // 序列化往返后规范化结果必须一致
}
```

## 9. FAQ / 常见坑

| 现象 | 原因与解法 |
| --- | --- |
| `Parse` 报 `syntax error at 1:x` | 语法不支持；先在 [coverage](coverage.md) 查是否支持，再查版本门控（`WithVersion` 过旧会拒绝新语法） |
| 裸写 `u@localhost` 报错 | 账号需带引号：`'u'@'%'`（裸写与 `@` 变量词法冲突，见 coverage 第 11 节） |
| `CAST(x AS BINARY)` 与官方行为不同 | `binary` 未映射为终结符（避免占用表选项名），但 `CAST(... AS BINARY)` 仍可解析且 round-trip；如与期望输出有差异，可用 `CHAR` 转换类型 |
| 想自定义门控 | 拷贝 `feature.DefaultRegistry` 后修改/追加 `feature.Rule`，用 `WithFeatures` 注入 |
| `Span()` 全零 | 手工构造的节点；解析产物才带位置 |
| 需要 Oracle 模式 | 本库仅 MySQL 模式（OB `compatibility mode = mysql`） |

## 10. 内部约定（扩展解析器时）

- 解析唯一引擎：`gparser`（goyacc 移植官方文法）。
- 上下文关键字：`user`/`role`/`analyze` 等仅特定语句位置生效，仍可作标识符。
- 只把文法用到的关键字映射为终结符，其余词一律 `NAME_OB`。
- `;` 由文法消费（顶层分隔 + PL 体内分号）。
- 版本门控在语义动作 `gateFeat` 抛 `bailout`，驱动 `recover` 为带 `line:col` 错误。
- **重新生成解析器必须走 `go generate`**（`gparser/`），裸 `goyacc` 会丢位置栈注入。
