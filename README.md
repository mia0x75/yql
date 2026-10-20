# YQL (Yet another Query Language)

[![Go Reference](https://pkg.go.dev/badge/github.com/mia0x75/yql.svg)](https://pkg.go.dev/github.com/mia0x75/yql)

YQL is a small, SQL-like expression language for Go rule engines. It covers the
`where` portion of SQL, adds set-oriented comparisons, and keeps the API to two
functions: `Match` and `Rule`.

```go
ok, err := yql.Match(`name='deen' and age>=23`, map[string]any{
    "name": "deen",
    "age":  23,
})
```

## Features

- Boolean expressions with `and`, `or`, and parentheses
- Type-aware comparison for numbers, strings, and booleans
- Set membership and intersection operators
- Cached rules through `Rule`
- Aggregate helpers such as `sum`, `avg`, `max`, `min`, and `count`

## Installation

```bash
go get github.com/mia0x75/yql
```

YQL is a library, so there is no binary to install. The module targets Go 1.26.8
or newer.

## Quick start

```go
package main

import (
    "fmt"

    "github.com/mia0x75/yql"
)

func main() {
    query := `name='deen' and age>=23 and (hobby in ('soccer', 'swim') or score>90)`
    ok, err := yql.Match(query, map[string]any{
        "name":  "deen",
        "age":   23,
        "hobby": "basketball",
        "score": int64(100),
    })
    fmt.Println(ok, err) // true <nil>
}
```

A runnable version of these examples is available in
[`example/main.go`](./example/main.go):

```bash
go run ./example
```

`Match` parses and evaluates the expression on every call. When the same rule is
evaluated repeatedly, use `Rule` to cache the parsed expression:

```go
ruler, err := yql.Rule(`name='deen' and age>=23`)
if err != nil {
    // invalid YQL
}

ok, err := ruler.Match(map[string]any{
    "name": "deen",
    "age":  23,
})
fmt.Println(ok, err) // true <nil>
```

`Rule` returns an error for invalid YQL. If the data map is empty, `Match`
returns `false, nil` without parsing the expression.

## Data types

Supported field types:

| Type | Examples |
| --- | --- |
| `int`, `int64`, `float64` | `age=23`, `score>=99.5` |
| `string` | `name='deen'` |
| `bool` | `enabled=true` |
| `json.Number` | integer JSON numbers |
| `[]int`, `[]int64`, `[]float64` | set comparisons and aggregates |
| `[]string` | set comparisons |

String values must be quoted, for example `name='deen'`. Numbers and booleans
are written directly, for example `age>=23` and `enabled=false`.

Comparison is type-aware: the value on the right is parsed according to the
dynamic type of the field on the left. Float comparisons use an epsilon of
`1e-10`.

## Operators

| Operator | Meaning |
| --- | --- |
| `=` | equal |
| `!=`, `<>` | not equal |
| `>` | greater than |
| `>=` | greater than or equal |
| `<` | less than |
| `<=` | less than or equal |
| `in` | value, or every element of a slice, is in the expected set |
| `!in` | negation of `in` |
| `intersect` | at least one element of a slice is in the expected set |
| `!intersect` | negation of `intersect` |

Set operators apply to slices. For example:

```go
yql.Match(`score in (7, 1, 9, 5, 3)`, map[string]any{
    "score": []int64{3, 5, 2},
}) // false: 2 is not in the expected set

yql.Match(`score intersect (7, 1, 9, 5, 3)`, map[string]any{
    "score": []int64{3, 100, 200},
}) // true: 3 is in both sets
```

Use parentheses to group expressions:

```text
a=1 and (b=2 or c=3)
```

`and` binds more tightly than `or`.

## Helpers

Helpers use `field.func()` syntax and can be chained. Supported helpers:

| Helper | Description |
| --- | --- |
| `sum` | total of a numeric slice |
| `avg` | average of a numeric slice |
| `max` | maximum value in a numeric slice |
| `min` | minimum value in a numeric slice |
| `count` | length of any slice; `1` for a scalar value |

`sum` accepts `[]int`, `[]int64`, and `[]float64`. `avg`, `max`, and `min`
require a non-empty numeric slice. `count` also accepts `[]string`.

```go
yql.Match(`score.sum() > 10`, map[string]any{
    "score": []int{1, 2, 3, 4, 5},
}) // true
```

## Syntax

- Field names start with a letter or underscore and contain letters or
  underscores.
- String values use single quotes: `'soccer'`.
- Numbers accept an optional sign and a decimal point: `-1`, `+1`, `1.5`.
- Boolean values are `true` and `false`.
- Use lowercase keywords, helper names, and set operators.
- `and` has higher precedence than `or`; use parentheses for clarity.

The grammar is defined in
[`grammar/parser.y`](./grammar/parser.y). The parser is generated with
[`goyacc`](https://pkg.go.dev/golang.org/x/tools/cmd/goyacc).

## Error handling

`Match` and `Rule` return an error when:

- the expression contains invalid syntax;
- a referenced field is missing from the data map;
- an unsupported helper is used;
- a value cannot be parsed for the field's type.

A failed comparison (for example, comparing incompatible types) returns `false`
without an error.

## Development

```bash
make help      # list available targets
make test      # go test ./...
make cover     # coverage summary
make vet       # go vet ./...
make generate  # regenerate grammar/parser.go (requires goyacc)
```

Do not edit `grammar/parser.go` by hand; update `grammar/parser.y` and run
`make generate`.

## License

This project is available under the [MIT License](./LICENSE) and the
[Anti 996 License](./LICENSE-ANTI-996).
