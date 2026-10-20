package yql

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/mia0x75/yql/grammar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleSyntaxErr(t *testing.T) {
	var testData = []string{
		"asdasdasd",
		"((a=1",
		"((a='foo'",
		"((a='foo')",
		"a=foo",
		"a in (foo,bar)",
		"a in (foo)",
		"a intersect (foo)",
		"a ∩ (1,2,3)",
		"a !∩ (1,2,3)",
		"a !=   foo",
	}
	data := map[string]any{
		"a": "foo",
	}
	ass := assert.New(t)
	for _, tc := range testData {
		ok, err := Match(tc, data)
		ass.False(ok)
		ass.Error(err)
	}
}

func TestMatch_Int(t *testing.T) {
	var testData = []struct {
		rawYql string
		data   map[string]any
		out    bool
	}{
		{
			rawYql: `a=10`,
			data: map[string]any{
				"a": 9,
			},
			out: false,
		},
		{
			rawYql: `a=10`,
			data: map[string]any{
				"a": 10,
			},
			out: true,
		},
		{
			rawYql: `a>10`,
			data: map[string]any{
				"a": 10,
			},
			out: false,
		},
		{
			rawYql: `a>10`,
			data: map[string]any{
				"a": 11,
			},
			out: true,
		},
		{
			rawYql: `a>=10`,
			data: map[string]any{
				"a": 10,
			},
			out: true,
		},
		{
			rawYql: `a>=10`,
			data: map[string]any{
				"a": 11,
			},
			out: true,
		},
		{
			rawYql: `a>=10`,
			data: map[string]any{
				"a": 1,
			},
			out: false,
		},
		{
			rawYql: `a<10`,
			data: map[string]any{
				"a": 1,
			},
			out: true,
		},
		{
			rawYql: `a<10`,
			data: map[string]any{
				"a": 10,
			},
			out: false,
		},
		{
			rawYql: `a<10`,
			data: map[string]any{
				"a": 11,
			},
			out: false,
		},
	}
	ass := assert.New(t)
	for _, tc := range testData {
		ok, err := Match(tc.rawYql, tc.data)
		ass.NoError(err)
		ass.Equal(tc.out, ok, "rawYql=%s||data=%+v", tc.rawYql, tc.data)
	}
}

func TestMatch_Int64(t *testing.T) {
	var testData = []struct {
		rawYql string
		data   map[string]any
		out    bool
	}{
		{
			rawYql: `a=10`,
			data: map[string]any{
				"a": int64(9),
			},
			out: false,
		},
		{
			rawYql: `a=10`,
			data: map[string]any{
				"a": int64(10),
			},
			out: true,
		},
		{
			rawYql: `a>10`,
			data: map[string]any{
				"a": int64(10),
			},
			out: false,
		},
		{
			rawYql: `a>10`,
			data: map[string]any{
				"a": int64(11),
			},
			out: true,
		},
		{
			rawYql: `a>=10`,
			data: map[string]any{
				"a": int64(10),
			},
			out: true,
		},
		{
			rawYql: `a>=10`,
			data: map[string]any{
				"a": int64(11),
			},
			out: true,
		},
		{
			rawYql: `a>=10`,
			data: map[string]any{
				"a": int64(1),
			},
			out: false,
		},
		{
			rawYql: `a<10`,
			data: map[string]any{
				"a": int64(1),
			},
			out: true,
		},
		{
			rawYql: `a<10`,
			data: map[string]any{
				"a": int64(10),
			},
			out: false,
		},
		{
			rawYql: `a<10`,
			data: map[string]any{
				"a": int64(11),
			},
			out: false,
		},
	}
	ass := assert.New(t)
	for _, tc := range testData {
		ok, err := Match(tc.rawYql, tc.data)
		ass.NoError(err)
		ass.Equal(tc.out, ok, "rawYql=%s||data=%+v", tc.rawYql, tc.data)
	}
}

func TestMatch_String(t *testing.T) {
	var testData = []struct {
		rawYql string
		data   map[string]any
		out    bool
	}{
		{
			rawYql: `a=10`,
			data: map[string]any{
				"a": "10",
			},
			out: true,
		},
		{
			rawYql: `a=10`,
			data: map[string]any{
				"a": "010",
			},
			out: false,
		},
		{
			rawYql: `a=10`,
			data: map[string]any{
				"a": "",
			},
			out: false,
		},
		{
			rawYql: `a>1`,
			data: map[string]any{
				"a": "1",
			},
			out: false,
		},
		{
			rawYql: `a>1`,
			data: map[string]any{
				"a": "2",
			},
			out: true,
		},
		{
			rawYql: `a>=1`,
			data: map[string]any{
				"a": "1",
			},
			out: true,
		},
		{
			rawYql: `a>=1`,
			data: map[string]any{
				"a": "2",
			},
			out: true,
		},
		{
			rawYql: `a>=1`,
			data: map[string]any{
				"a": "",
			},
			out: false,
		},
		{
			rawYql: `a>20`,
			data: map[string]any{
				"a": "21",
			},
			out: true,
		},
		{
			rawYql: `a>20`,
			data: map[string]any{
				"a": "3",
			},
			out: true,
		},
		{
			rawYql: `a>20`,
			data: map[string]any{
				"a": "1",
			},
			out: false,
		},
		{
			rawYql: `a!=20`,
			data: map[string]any{
				"a": "3",
			},
			out: true,
		},
		{
			rawYql: `a!=20`,
			data: map[string]any{
				"a": "20",
			},
			out: false,
		},
	}
	ass := assert.New(t)
	for _, tc := range testData {
		ok, err := Match(tc.rawYql, tc.data)
		ass.NoError(err)
		ass.Equal(tc.out, ok, "rawYql=%s||data=%+v", tc.rawYql, tc.data)
	}
}

func TestMatch__Float(t *testing.T) {
	var testData = []struct {
		rawYql string
		data   map[string]any
		out    bool
	}{
		{
			rawYql: `a=10`,
			data: map[string]any{
				"a": float64(10),
			},
			out: true,
		},
		{
			rawYql: `a=10`,
			data: map[string]any{
				"a": float64(10 - epsilon*0.1),
			},
			out: true,
		},
		{
			rawYql: `a!=10`,
			data: map[string]any{
				"a": float64(10),
			},
			out: false,
		},
		{
			rawYql: `a>10`,
			data: map[string]any{
				"a": float64(10),
			},
			out: false,
		},
		{
			rawYql: `a>10`,
			data: map[string]any{
				"a": float64(10.000000001),
			},
			out: true,
		},
		{
			rawYql: `a<10`,
			data: map[string]any{
				"a": float64(10 - epsilon*0.1),
			},
			out: true,
		},
		{
			rawYql: `a>=10`,
			data: map[string]any{
				"a": float64(10.0),
			},
			out: true,
		},
		{
			rawYql: `a>=10`,
			data: map[string]any{
				"a": float64(9.0),
			},
			out: false,
		},
		{
			rawYql: `a<=10`,
			data: map[string]any{
				"a": float64(10.0),
			},
			out: true,
		},
		{
			rawYql: `a<=10`,
			data: map[string]any{
				"a": float64(9.0),
			},
			out: true,
		},
		{
			rawYql: `a<=10`,
			data: map[string]any{
				"a": float64(11.0),
			},
			out: false,
		},
	}
	ass := assert.New(t)
	for _, tc := range testData {
		ok, _ := Match(tc.rawYql, tc.data)
		ass.Equal(tc.out, ok, "rawYql=%s||data=%+v", tc.rawYql, tc.data)
	}
}

func TestMatch_Boolean(t *testing.T) {
	// comparison between booleans only support =,!=
	var testData = []struct {
		rawYql string
		data   map[string]any
		out    bool
	}{
		{
			rawYql: `a=true`,
			data: map[string]any{
				"a": true,
			},
			out: true,
		},
		{
			rawYql: `a=true`,
			data: map[string]any{
				"a": false,
			},
			out: false,
		},
		{
			rawYql: `a=false`,
			data: map[string]any{
				"a": true,
			},
			out: false,
		},
		{
			rawYql: `a=false`,
			data: map[string]any{
				"a": false,
			},
			out: true,
		},
		{
			rawYql: `a!=false`,
			data: map[string]any{
				"a": true,
			},
			out: true,
		},
		{
			rawYql: `a>false`,
			data: map[string]any{
				"a": true,
			},
			out: false,
		},
		{
			rawYql: `a<=false`,
			data: map[string]any{
				"a": true,
			},
			out: false,
		},
	}
	ass := assert.New(t)
	for _, tc := range testData {
		ok, err := Match(tc.rawYql, tc.data)
		ass.NoError(err)
		ass.Equal(tc.out, ok, "rawYql=%s||data=%+v", tc.rawYql, tc.data)
	}
}

func TestMatch_In(t *testing.T) {
	var testData = []struct {
		rawYql string
		data   map[string]any
		out    bool
	}{
		{
			rawYql: `a in ('swim')`,
			data: map[string]any{
				"a": "swim",
			},
			out: true,
		},
		{
			rawYql: `a in ('soccer')`,
			data: map[string]any{
				"a": "swim",
			},
			out: false,
		},
		{
			rawYql: `a !in ('swim')`,
			data: map[string]any{
				"a": "swim",
			},
			out: false,
		},
		{
			rawYql: `a !in ('soccer')`,
			data: map[string]any{
				"a": "swim",
			},
			out: true,
		},
		{
			rawYql: `a !in (1,2,3)`,
			data: map[string]any{
				"a": 1,
			},
			out: false,
		},
		{
			rawYql: `a in (1,2,3)`,
			data: map[string]any{
				"a": 3,
			},
			out: true,
		},
		{
			rawYql: `a in (1,5,9)`,
			data: map[string]any{
				"a": 5,
			},
			out: true,
		},
		{
			rawYql: `a !in (1,2, 10,   5)`,
			data: map[string]any{
				"a": int64(9),
			},
			out: true,
		},
		{
			rawYql: `a !in (1,2, 10,   5)`,
			data: map[string]any{
				"a": 9,
			},
			out: true,
		},
		{
			rawYql: `a in (1,2, 10,   5)`,
			data: map[string]any{
				"a": 10,
			},
			out: true,
		},
		{
			rawYql: `a in (1,2, 10,   5)`,
			data: map[string]any{
				"a": []int{10, 2},
			},
			out: true,
		},
		{
			rawYql: `a in (1,2, 10,   5)`,
			data: map[string]any{
				"a": []string{"1", "5"},
			},
			out: true,
		},
		{
			rawYql: `a !in (1,2, 10,   5)`,
			data: map[string]any{
				"a": []string{"1", "5"},
			},
			out: false,
		},
		{
			rawYql: `a !in (1,2, 10,   5)`,
			data: map[string]any{
				"a": "1",
			},
			out: false,
		},
		{
			rawYql: `a !in (1,2, 10,   5)`,
			data: map[string]any{
				"a": []string{"1", "5", "3"},
			},
			out: true,
		},
		{
			rawYql: `a in (1,2, 10,   5)`,
			data: map[string]any{
				"a": []int64{1, 5},
			},
			out: true,
		},
		{
			rawYql: `a in (1,2, 10.000,   5.000)`,
			data: map[string]any{
				"a": []float64{2.000000000000001, 5.000000000000000002},
			},
			out: true,
		},
		{
			rawYql: `a in (1,2,3, 10.00000000001)`,
			data: map[string]any{
				"a": float64(10.0),
			},
			out: true,
		},
		{
			rawYql: `a !in (1,2,3, 10.00001)`,
			data: map[string]any{
				"a": float64(10.0),
			},
			out: true,
		},
		{
			rawYql: `a in (1,2, 10)`,
			data: map[string]any{
				"a": []int64{2, 3},
			},
			out: false,
		},
		{
			rawYql: `a in (1,2,3, 10)`,
			data: map[string]any{
				"a": int64(10),
			},
			out: true,
		},
		{
			rawYql: `a in (1,2,3, 10)`,
			data: map[string]any{
				"a": []int64{2, 3},
			},
			out: true,
		},
		{
			rawYql: `a in (1,2, 10)`,
			data: map[string]any{
				"a": []int64{2, 3},
			},
			out: false,
		},
		{
			rawYql: `a in (1,2, 10,   5)`,
			data: map[string]any{
				"a": []float64{2, 5},
			},
			out: true,
		},
		{
			rawYql: `a in (1,2, 10,   5)`,
			data: map[string]any{
				"a": []float64{1, 10, 2, 5},
			},
			out: true,
		},
		{
			rawYql: `a in (1,2, 10,   5)`,
			data: map[string]any{
				"a": []float64{1, 10, 2, 5, 3},
			},
			out: false,
		},
		{
			rawYql: `a in (1,2, 10,   5)`,
			data: map[string]any{
				"a": []string{"1", "5"},
			},
			out: true,
		},
	}
	ass := assert.New(t)
	for _, tc := range testData {
		ok, _ := Match(tc.rawYql, tc.data)
		ass.Equal(tc.out, ok, "rawYql=%s||data=%+v", tc.rawYql, tc.data)
	}
}

func TestMatch_And(t *testing.T) {
	var testData = []struct {
		rawYql string
		data   map[string]any
		out    bool
	}{
		{
			rawYql: `a=10 and b>'2' and c<9 and d!=2`,
			data: map[string]any{
				"a": int64(10),
				"b": int64(3),
				"c": int64(-1),
				"d": int64(2),
			},
			out: false,
		},
		{
			rawYql: `a=10 and b>'2' and c<9 and d!=2`,
			data: map[string]any{
				"a": int64(10),
				"b": int64(3),
				"c": int64(-1),
			},
			out: false,
		},
		{
			rawYql: `a=10 and b>'2'`,
			data: map[string]any{
				"a": int64(10),
				"b": int64(3),
			},
			out: true,
		},
		{
			rawYql: `a=10 and b>'2'`,
			data: map[string]any{
				"a": int64(10),
				"b": int64(2),
			},
			out: false,
		},
		{
			rawYql: `a=10 and b>'2'`,
			data: map[string]any{
				"a": int64(10),
				"b": int64(3),
			},
			out: true,
		},
		{
			rawYql: `a=10 and b>'2' and c<9`,
			data: map[string]any{
				"a": int64(10),
				"b": int64(3),
				"c": int64(-1),
			},
			out: true,
		},
		{
			rawYql: `a=10 and b>'2' and c<9 and d!=2`,
			data: map[string]any{
				"a": int64(10),
				"b": int64(3),
				"c": int64(-1),
				"d": int64(0),
			},
			out: true,
		},
	}
	ass := assert.New(t)
	for _, tc := range testData {
		ok, _ := Match(tc.rawYql, tc.data)
		ass.Equal(tc.out, ok, "rawYql=%s||data=%+v", tc.rawYql, tc.data)
	}
}

func TestMatch_Or(t *testing.T) {
	var testData = []struct {
		rawYql string
		data   map[string]any
		out    bool
	}{
		{
			rawYql: `a=10 or b>'2'`,
			data: map[string]any{
				"a": int64(10),
				"b": int64(1),
			},
			out: true,
		},
		{
			rawYql: `a=10 or b>'2'`,
			data: map[string]any{
				"a": int64(9),
				"b": int64(2),
			},
			out: false,
		},
		{
			rawYql: `a=10 or b>'2'`,
			data: map[string]any{
				"a": int64(10),
				"b": int64(3),
			},
			out: true,
		},
		{
			rawYql: `a=10 or b>'2' or c<9`,
			data: map[string]any{
				"a": int64(1),
				"b": int64(3),
				"c": int64(100),
			},
			out: true,
		},
		{
			rawYql: `a=10 or b>'2' or c<9 or d!=2`,
			data: map[string]any{
				"a": int64(1),
				"b": int64(2),
				"c": int64(10),
				"d": int64(0),
			},
			out: true,
		},
		{
			rawYql: `a=10 or b>'2' or c<9 or d!=2`,
			data: map[string]any{
				"a": int64(1),
				"b": int64(1),
				"c": int64(10),
				"d": int64(2),
			},
			out: false,
		},
	}
	ass := assert.New(t)
	for _, tc := range testData {
		ok, err := Match(tc.rawYql, tc.data)
		ass.NoError(err)
		ass.Equal(tc.out, ok, "rawYql=%s||data=%+v", tc.rawYql, tc.data)
	}
}

func TestMatch_Or_And(t *testing.T) {
	var testData = []struct {
		rawYql string
		data   map[string]any
		out    bool
	}{
		{
			rawYql: `a=9 or c=1 and b!='1'`,
			data: map[string]any{
				"a": int64(10),
				"b": int64(1),
				"c": int64(1),
			},
			out: false,
		},
		{
			rawYql: `a=10 and b>'2' or c=1`,
			data: map[string]any{
				"a": int64(10),
				"b": int64(1),
				"c": int64(1),
			},
			out: true,
		},
		{
			rawYql: `a=10 or c=1 and b!='1'`,
			data: map[string]any{
				"a": int64(10),
				"b": int64(1),
				"c": int64(1),
			},
			out: true,
		},
		{
			rawYql: `a=10 and (c=1 or b!='1')`,
			data: map[string]any{
				"a": int64(10),
				"b": int64(1),
				"c": int64(1),
			},
			out: true,
		},
		{
			rawYql: `a=10 and (c=1 or b!='1') and d='123'`,
			data: map[string]any{
				"a": int64(10),
				"b": int64(1),
				"c": int64(1),
				"d": "123",
			},
			out: true,
		},
	}
	ass := assert.New(t)
	for _, tc := range testData {
		ok, err := Match(tc.rawYql, tc.data)
		ass.NoError(err)
		ass.Equal(tc.out, ok, "rawYql=%s||data=%+v", tc.rawYql, tc.data)
	}
}

func TestMatch_Inter(t *testing.T) {
	var testData = []struct {
		rawYql string
		data   map[string]any
		out    bool
	}{
		{
			rawYql: `letter intersect ('a','b','c','d','e')`,
			data: map[string]any{
				"letter": []string{"a", "e"},
			},
			out: true,
		},
		{
			rawYql: `letter !intersect (1,2,3)`,
			data: map[string]any{
				"letter": []float64{0.5, 3.01},
			},
			out: true,
		},
		{
			rawYql: `letter intersect ('a', 'b','c','d',  'e')`,
			data: map[string]any{
				"letter": []string{"a", "e", "f"},
			},
			out: true,
		},
		{
			rawYql: `letter intersect ('a','b','c', 'd',  'e')`,
			data: map[string]any{
				"letter": []string{"f"},
			},
			out: false,
		},
		{
			rawYql: `letter intersect ('a','b','c','d')`,
			data: map[string]any{
				"letter": "c",
			},
			out: true,
		},
		{
			rawYql: `letter intersect (1,2,3)`,
			data: map[string]any{
				"letter": []float64{2.0, 3.0},
			},
			out: true,
		},
		{
			rawYql: `letter intersect (1,2,3)`,
			data: map[string]any{
				"letter": []int64{2, 5},
			},
			out: true,
		},
		{
			rawYql: `letter !intersect (1,2,3)`,
			data: map[string]any{
				"letter": []int64{4, 5},
			},
			out: true,
		},
	}
	ass := assert.New(t)
	for _, tc := range testData {
		ok, err := Match(tc.rawYql, tc.data)
		ass.NoError(err)
		ass.Equal(tc.out, ok, "rawYql=%s||data=%+v", tc.rawYql, tc.data)
	}
}

func TestMatch(t *testing.T) {
	var testData = []struct {
		rawYql string
		data   map[string]any
		out    bool
	}{
		{
			rawYql: `age>23 and (sex in ('boy','girl') or sex='other') and score>=95 and rank !in ('b','c','d')`,
			data: map[string]any{
				"age":   int64(24),
				"sex":   "boy",
				"score": int64(95),
				"rank":  "s",
			},
			out: true,
		},
		{
			rawYql: `age>23 and (sex in ('boy','girl') or sex='other')`,
			data: map[string]any{
				"age": int64(24),
				"sex": "other",
			},
			out: true,
		},
		{
			rawYql: `age>23 and (sex in ('boy','girl') or sex='other')`,
			data: map[string]any{
				"age": int64(24),
				"sex": "boy",
			},
			out: true,
		},
		{
			rawYql: `age>23 and (sex in ('boy','girl') or some!=5) and words='hello world'`,
			data: map[string]any{
				"age":   int64(211),
				"sex":   "boy",
				"some":  int64(6),
				"words": "hello world",
			},
			out: true,
		},
		{
			rawYql: `age>23 and (sex in ('boy','girl') or some!=5) and words='hello world'`,
			data: map[string]any{
				"age":   int64(21),
				"sex":   "boy",
				"some":  int64(6),
				"words": "hello world",
			},
			out: false,
		},
		{
			rawYql: `tag in (1,3,5) and status!=0`,
			data: map[string]any{
				"tag":    []int64{1, 5},
				"status": int64(3),
			},
			out: true,
		},
	}
	ass := assert.New(t)
	for _, tc := range testData {
		ok, err := Match(tc.rawYql, tc.data)
		ass.NoError(err)
		ass.Equal(tc.out, ok, "rawYql=%s||data=%+v", tc.rawYql, tc.data)
	}
}

func TestRule_Match(t *testing.T) {
	var testData = []struct {
		rawYql string
		data   map[string]any
		out    bool
	}{
		{
			rawYql: `age>23 and (sex in ('boy','girl') or sex='other') and score>=95 and rank !in ('b','c','d')`,
			data: map[string]any{
				"age":   int64(24),
				"sex":   "boy",
				"score": int64(95),
				"rank":  "s",
			},
			out: true,
		},
		{
			rawYql: `age>23 and (sex in ('boy','girl') or sex='other')`,
			data: map[string]any{
				"age": int64(24),
				"sex": "other",
			},
			out: true,
		},
		{
			rawYql: `age>23 and (sex in ('boy','girl') or sex='other')`,
			data: map[string]any{
				"age": int64(24),
				"sex": "boy",
			},
			out: true,
		},
		{
			rawYql: `age>23 and (sex in ('boy','girl') or some!=5) and words='hello world'`,
			data: map[string]any{
				"age":   int64(211),
				"sex":   "boy",
				"some":  int64(6),
				"words": "hello world",
			},
			out: true,
		},
		{
			rawYql: `age>23 and (sex in ('boy','girl') or some!=5) and words='hello world'`,
			data: map[string]any{
				"age":   int64(21),
				"sex":   "boy",
				"some":  int64(6),
				"words": "hello world",
			},
			out: false,
		},
		{
			rawYql: `tag in (1,3,5) and status!=0`,
			data: map[string]any{
				"tag":    []int64{1, 5},
				"status": int64(3),
			},
			out: true,
		},
	}
	ass := assert.New(t)
	for _, tc := range testData {
		ruler, err := Rule(tc.rawYql)
		ass.NoError(err)
		ok, err := ruler.Match(tc.data)
		ass.NoError(err)
		ass.Equal(tc.out, ok, "rawYql=%s||data=%+v", tc.rawYql, tc.data)
	}
}

func TestRule_Match_Multi(t *testing.T) {
	var testData = []struct {
		data map[string]any
		out  bool
	}{
		{
			data: map[string]any{
				"age":   int64(24),
				"sex":   "boy",
				"score": int64(95),
				"rank":  "s",
			},
			out: true,
		},
		{
			data: map[string]any{
				"age":   int64(23),
				"sex":   "boy",
				"score": int64(95),
				"rank":  "s",
			},
			out: false,
		},
		{
			data: map[string]any{
				"age":   int64(24),
				"sex":   "boy",
				"score": int64(95),
				"rank":  "",
			},
			out: true,
		},
		{
			data: map[string]any{
				"age":   int64(24),
				"sex":   "boy",
				"score": int64(95),
				"rank":  "a",
			},
			out: true,
		},
		{
			data: map[string]any{
				"age":   int64(23),
				"sex":   "boy",
				"score": int64(92),
				"rank":  "s",
			},
			out: false,
		},
		{
			data: map[string]any{
				"age":   int64(23),
				"sex":   "other",
				"score": int64(97),
				"rank":  "s",
			},
			out: false,
		},
	}
	ass := assert.New(t)
	ruler, err := Rule(`age>23 and (sex in ('boy','girl') or sex='other') and score>=95 and rank !in ('b','c','d')`)
	if !ass.NoError(err) {
		t.FailNow()
	}
	for _, tc := range testData {
		ok, err := ruler.Match(tc.data)
		ass.NoError(err)
		ass.Equal(tc.out, ok, "data=%+v", tc.data)
	}
}

func Test_Compare_Slice_And_One_Element(t *testing.T) {
	should := require.New(t)
	var testData = []struct {
		rawYql string
		data   map[string]any
		out    bool
	}{
		{
			rawYql: `letter intersect ('a')`,
			data: map[string]any{
				"letter": []string{"a", "b", "c"},
			},
			out: true,
		},
		{
			rawYql: `letter intersect ('a', 'b')`,
			data: map[string]any{
				"letter": []string{"a", "b", "c"},
			},
			out: true,
		},
		{
			rawYql: `letter intersect ('d')`,
			data: map[string]any{
				"letter": []string{"a", "b", "c"},
			},
			out: false,
		},
		{
			rawYql: `letter in ('a')`,
			data: map[string]any{
				"letter": []string{"a", "b", "c"},
			},
			out: false,
		},
	}
	for _, tc := range testData {
		actual, err := Match(tc.rawYql, tc.data)
		should.NoError(err)
		should.Equal(tc.out, actual)
	}
}

func ExampleMatch() {
	rawYQL := `name='deen' and age>=23 and (hobby in ('soccer', 'swim') or score>90)`
	result, _ := Match(rawYQL, map[string]interface{}{
		"name":  "deen",
		"age":   23,
		"hobby": "basketball",
		"score": int64(100),
	})
	fmt.Println(result)
	rawYQL = `score intersect (7,1,9,5,3)`
	result, _ = Match(rawYQL, map[string]interface{}{
		"score": []int64{3, 100, 200},
	})
	fmt.Println(result)
	rawYQL = `score in (7,1,9,5,3)`
	result, _ = Match(rawYQL, map[string]interface{}{
		"score": []int64{3, 5, 2},
	})
	fmt.Println(result)
	rawYQL = `score.sum() > 10`
	result, _ = Match(rawYQL, map[string]interface{}{
		"score": []int{1, 2, 3, 4, 5},
	})
	fmt.Println(result)
	//Output:
	//true
	//true
	//false
	//true
}

func ExampleRule() {
	rawYQL := `name='deen' and age>=23 and (hobby in ('soccer', 'swim') or score>90)`
	ruler, _ := Rule(rawYQL)
	result, _ := ruler.Match(map[string]interface{}{
		"name":  "deen",
		"age":   23,
		"hobby": "basketball",
		"score": int64(100),
	})
	fmt.Println(result)
	result, _ = ruler.Match(map[string]interface{}{
		"name":  "deen",
		"age":   23,
		"hobby": "basketball",
		"score": int64(90),
	})
	fmt.Println(result)
	//Output:
	//true
	//false
}

func TestMatch_EmptyData(t *testing.T) {
	ass := assert.New(t)
	ok, err := Match("a=1", nil)
	ass.False(ok)
	ass.NoError(err)

	// invalid yql must not be parsed when data is empty
	ok, err = Match("invalid", nil)
	ass.False(ok)
	ass.NoError(err)
}

func TestRule_EmptyData(t *testing.T) {
	ass := assert.New(t)
	ruler, err := Rule("a=1")
	ass.NoError(err)

	ok, err := ruler.Match(nil)
	ass.False(ok)
	ass.NoError(err)
}

func TestCachedAST_Invalid(t *testing.T) {
	ok, err := cachedAST{}.Match(map[string]any{"a": 1})
	assert.False(t, ok)
	assert.Error(t, err)
}

func TestRule_SyntaxError(t *testing.T) {
	ruler, err := Rule("a=foo")
	assert.Error(t, err)
	assert.NotNil(t, ruler)
}

func TestEval_Errors(t *testing.T) {
	ass := assert.New(t)
	data := map[string]any{"a": 1}
	valid := &grammar.CompareExpr{Field: "a", Op: opEqual, Values: []string{"1"}}
	missing := &grammar.CompareExpr{Field: "missing", Op: opEqual, Values: []string{"1"}}

	_, err := eval(nil, data)
	ass.Error(err)

	_, err = eval(&grammar.BinaryExpr{Op: "unknown", Left: valid, Right: valid}, data)
	ass.Error(err)

	_, err = eval(&grammar.BinaryExpr{Op: "AND", Left: missing, Right: valid}, data)
	ass.Error(err)

	_, err = eval(&grammar.BinaryExpr{Op: "AND", Left: valid, Right: missing}, data)
	ass.Error(err)

	_, err = eval(&grammar.CompareExpr{Field: "a", Funcs: []string{"unknown"}, Op: opEqual, Values: []string{"1"}}, data)
	ass.Error(err)
}

func TestCompare_TypeParseErrors(t *testing.T) {
	var testData = []struct {
		rawYql string
		data   map[string]any
		out    bool
	}{
		{`a='x'`, map[string]any{"a": 1}, false},
		{`a='x'`, map[string]any{"a": int64(1)}, false},
		{`a='x'`, map[string]any{"a": float64(1)}, false},
		{`a='x'`, map[string]any{"a": true}, false},
		{`a=10`, map[string]any{"a": json.Number("10")}, true},
		{`a='x'`, map[string]any{"a": json.Number("10")}, false},
		{`a=10`, map[string]any{"a": json.Number("x")}, false},
		{`a=1`, map[string]any{"a": struct{}{}}, false},
	}
	ass := assert.New(t)
	for _, tc := range testData {
		ok, err := Match(tc.rawYql, tc.data)
		ass.NoError(err)
		ass.Equal(tc.out, ok, "rawYql=%s||data=%+v", tc.rawYql, tc.data)
	}
}

func TestMatch_SignedNumbers(t *testing.T) {
	ass := assert.New(t)
	ok, err := Match(`a=-1`, map[string]any{"a": -1})
	ass.NoError(err)
	ass.True(ok)

	ok, err = Match(`a=+1`, map[string]any{"a": 1})
	ass.NoError(err)
	ass.True(ok)
}

func TestMatch_NotEqualSynonym(t *testing.T) {
	ass := assert.New(t)
	ok, err := Match(`a<>1`, map[string]any{"a": 2})
	ass.NoError(err)
	ass.True(ok)

	ok, err = Match(`a<>1`, map[string]any{"a": 1})
	ass.NoError(err)
	ass.False(ok)
}

func TestMatch_ChainedFuncs(t *testing.T) {
	ass := assert.New(t)
	ok, err := Match(`a.sum().count()=1`, map[string]any{"a": []int{1, 2}})
	ass.NoError(err)
	ass.True(ok)
}

func TestMatch_UnsupportedFunc(t *testing.T) {
	ok, err := Match(`a.SUM()>1`, map[string]any{"a": []int{1, 2}})
	assert.False(t, ok)
	assert.Error(t, err)
}

func TestMatch_InvalidSyntaxVariants(t *testing.T) {
	var testData = []string{
		"",
		" ",
		"a='foo",
		"a=+x",
		"a=1 @",
	}
	data := map[string]any{"a": 1}
	ass := assert.New(t)
	for _, rawYql := range testData {
		ok, err := Match(rawYql, data)
		ass.False(ok, "rawYql=%q", rawYql)
		ass.Error(err, "rawYql=%q", rawYql)
	}
}
