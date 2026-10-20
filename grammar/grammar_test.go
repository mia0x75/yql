package grammar

import (
	"strings"
	"testing"
)

func TestParse_CompareExpr(t *testing.T) {
	var testData = []struct {
		input string
		op    string
		value string
	}{
		{`a=1`, `=`, "1"},
		{`a!=1`, `!=`, "1"},
		{`a<>1`, `<>`, "1"},
		{`a>1`, `>`, "1"},
		{`a>=1`, `>=`, "1"},
		{`a<1`, `<`, "1"},
		{`a<=1`, `<=`, "1"},
		{`a=1.5`, `=`, "1.5"},
		{`a=true`, `=`, "true"},
		{`a=FALSE`, `=`, "FALSE"},
		{`a=-1`, `=`, "-1"},
		{`a=+1`, `=`, "+1"},
		{`a='hello world'`, `=`, "'hello world'"},
	}
	for _, tc := range testData {
		query, err := Parse(tc.input)
		if err != nil {
			t.Fatalf("parse %q: %v", tc.input, err)
		}
		expr, ok := query.Expr.(*CompareExpr)
		if !ok {
			t.Fatalf("parse %q: got %T, want *CompareExpr", tc.input, query.Expr)
		}
		if expr.Field != "a" || expr.Op != tc.op || len(expr.Values) != 1 || expr.Values[0] != tc.value {
			t.Fatalf("parse %q: got %+v", tc.input, expr)
		}
	}
}

func TestParse_SetExpr(t *testing.T) {
	var testData = []struct {
		input  string
		op     string
		values []string
	}{
		{`a in (1,2,3)`, `in`, []string{"1", "2", "3"}},
		{`a !in (1)`, `!in`, []string{"1"}},
		{`a intersect ('x','y')`, `intersect`, []string{"'x'", "'y'"}},
		{`a !intersect (1.5)`, `!intersect`, []string{"1.5"}},
	}
	for _, tc := range testData {
		query, err := Parse(tc.input)
		if err != nil {
			t.Fatalf("parse %q: %v", tc.input, err)
		}
		expr, ok := query.Expr.(*CompareExpr)
		if !ok {
			t.Fatalf("parse %q: got %T, want *CompareExpr", tc.input, query.Expr)
		}
		if expr.Field != "a" || expr.Op != tc.op {
			t.Fatalf("parse %q: got %+v", tc.input, expr)
		}
		if len(expr.Values) != len(tc.values) {
			t.Fatalf("parse %q: got values %v, want %v", tc.input, expr.Values, tc.values)
		}
		for i := range tc.values {
			if expr.Values[i] != tc.values[i] {
				t.Fatalf("parse %q: got values %v, want %v", tc.input, expr.Values, tc.values)
			}
		}
	}
}

func TestParse_FuncChain(t *testing.T) {
	query, err := Parse(`a.sum().count()=1`)
	if err != nil {
		t.Fatal(err)
	}
	expr, ok := query.Expr.(*CompareExpr)
	if !ok {
		t.Fatalf("got %T, want *CompareExpr", query.Expr)
	}
	if expr.Field != "a" || len(expr.Funcs) != 2 || expr.Funcs[0] != "sum" || expr.Funcs[1] != "count" {
		t.Fatalf("got %+v", expr)
	}
}

func TestParse_BinaryExpr(t *testing.T) {
	query, err := Parse(`a=1 and (b=2 or c=3)`)
	if err != nil {
		t.Fatal(err)
	}
	root, ok := query.Expr.(*BinaryExpr)
	if !ok || root.Op != "AND" {
		t.Fatalf("got %#v, want AND BinaryExpr", query.Expr)
	}
	right, ok := root.Right.(*BinaryExpr)
	if !ok || right.Op != "OR" {
		t.Fatalf("got %#v, want OR BinaryExpr", root.Right)
	}
}

func TestParse_FieldName(t *testing.T) {
	query, err := Parse(`_a=1`)
	if err != nil {
		t.Fatal(err)
	}
	expr, ok := query.Expr.(*CompareExpr)
	if !ok || expr.Field != "_a" {
		t.Fatalf("got %#v, want field _a", query.Expr)
	}
}

func TestParse_DeepExpression(t *testing.T) {
	input := strings.Repeat(`a=1 and `, 300) + `a=1`
	if _, err := Parse(input); err != nil {
		t.Fatal(err)
	}

	nested := strings.Repeat("(", 300) + `a=1` + strings.Repeat(")", 300)
	if _, err := Parse(nested); err != nil {
		t.Fatal(err)
	}
}

func TestParse_Errors(t *testing.T) {
	var testData = []string{
		"",
		" ",
		"a",
		"a=",
		"a=foo",
		"a='foo",
		"a=+x",
		"a=1 @",
		"a in (1,)",
		"a in (1 2)",
		"a in ()",
		"a.sum",
		"a.sum(",
		"a.sum()",
		"a ∩ (1)",
	}
	for _, input := range testData {
		if query, err := Parse(input); err == nil {
			t.Fatalf("parse %q: expected error, got %+v", input, query)
		}
	}
}

func TestExprMarkers(t *testing.T) {
	var expr Expr = &BinaryExpr{}
	expr.isExpr()
	expr = &CompareExpr{}
	expr.isExpr()
}

func TestYySymName(t *testing.T) {
	if got := yySymName('x'); got != "'x'" {
		t.Fatalf("yySymName('x')=%q", got)
	}
	if got := yySymName(0x100); got != "256" {
		t.Fatalf("yySymName(0x100)=%q", got)
	}
}

func TestParserDebug(t *testing.T) {
	old := yyDebug
	yyDebug = 4
	defer func() { yyDebug = old }()

	if _, err := Parse(`a=1 and (b=2 or c=3)`); err != nil {
		t.Fatal(err)
	}
	if _, err := Parse(`a=1 foo`); err == nil {
		t.Fatal("expected error")
	}
}

type unknownTokenLexer struct {
	called bool
}

func (l *unknownTokenLexer) Lex(lval *yySymType) int {
	if l.called {
		return 0
	}
	l.called = true
	return 0x100
}

func (l *unknownTokenLexer) Error(string) {}

func TestUnknownToken(t *testing.T) {
	if got := yyParse(&unknownTokenLexer{}); got == 0 {
		t.Fatal("expected parse failure for unknown token")
	}
}
