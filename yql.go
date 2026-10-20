package yql

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/mia0x75/yql/grammar"
)

func match(rawYQL string, data map[string]any) (ok bool, err error) {
	if 0 == len(data) {
		return false, nil
	}
	query, err := grammar.Parse(rawYQL)
	if err != nil {
		return false, err
	}
	return eval(query.Expr, data)
}

func eval(expr grammar.Expr, data map[string]any) (bool, error) {
	switch e := expr.(type) {
	case *grammar.BinaryExpr:
		left, err := eval(e.Left, data)
		if err != nil {
			return false, err
		}
		right, err := eval(e.Right, data)
		if err != nil {
			return false, err
		}
		switch e.Op {
		case "AND":
			return left && right, nil
		case "OR":
			return left || right, nil
		default:
			return false, fmt.Errorf("unknown operator %s", e.Op)
		}
	case *grammar.CompareExpr:
		actualValue, ok := data[e.Field]
		if !ok {
			return false, fmt.Errorf("%s not provided", e.Field)
		}
		var err error
		actualValue, err = funcRuner(actualValue, e.Funcs)
		if err != nil {
			return false, err
		}
		return compare(actualValue, e.Values, e.Op), nil
	default:
		return false, fmt.Errorf("unknown expression")
	}
}

// Ruler caches an AST
type cachedAST struct {
	query *grammar.Query
}

// Ruler represents an AST structure and could do the match according to input data
type Ruler interface {
	Match(map[string]any) (bool, error)
}

// Match do the comparison according to the input data
func (ast cachedAST) Match(data map[string]any) (bool, error) {
	if 0 == len(data) {
		return false, nil
	}
	if ast.query == nil {
		return false, fmt.Errorf("invalid yql")
	}
	return eval(ast.query.Expr, data)
}

// Rule parses a rule and stores the AST.
// It returns an error when the YQL expression is invalid.
// Reusing a Rule is faster than calling Match repeatedly.
func Rule(rawYQL string) (ruler Ruler, err error) {
	query, err := grammar.Parse(rawYQL)
	if err != nil {
		return cachedAST{}, err
	}
	return cachedAST{query: query}, nil
}

// Match parses rawYQL and evaluates it with the provided data.
// An error is returned if rawYQL contains invalid syntax.
func Match(rawYQL string, data map[string]any) (bool, error) {
	return match(rawYQL, data)
}

func compare(actualValue any, expectValue []string, op string) bool {
	if len(expectValue) > 1 {
		return compareSet(actualValue, expectValue, op)
	}
	if actualValue != nil && reflect.TypeOf(actualValue).Kind() == reflect.Slice {
		return compareSet(actualValue, expectValue, op)
	}
	e := removeQuote(expectValue[0])
	switch actual := actualValue.(type) {
	case int:
		expect, err := strconv.ParseInt(e, 10, 64)
		if nil != err {
			return false
		}
		return cmpInt(int64(actual), expect, op)
	case int64:
		expect, err := strconv.ParseInt(e, 10, 64)
		if nil != err {
			return false
		}
		return cmpInt(actual, expect, op)
	case json.Number:
		actualValue, err := strconv.ParseInt(actualValue.(json.Number).String(), 10, 64)
		if nil != err {
			return false
		}
		expect, err := strconv.ParseInt(e, 10, 64)
		if nil != err {
			return false
		}
		return cmpInt(actualValue, expect, op)
	case float64:
		expect, err := strconv.ParseFloat(e, 64)
		if nil != err {
			return false
		}
		return cmpFloat(actual, expect, op)
	case string:
		return cmpStr(actual, e, op)
	case bool:
		expect, err := strconv.ParseBool(e)
		if nil != err {
			return false
		}
		return cmpBool(actual, expect, op)
	default:
		return false
	}
}

func removeQuote(str string) string {
	l := 0
	r := len(str)
	if str[0] == '\'' {
		l++
	}
	if str[r-1] == '\'' {
		r--
	}
	return str[l:r]
}
