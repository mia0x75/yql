package yql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntSetBelong(t *testing.T) {
	var testData = []struct {
		ina, inb []int64
		out      bool
	}{
		{
			ina: []int64{3, 1, 2},
			inb: []int64{5, 1, 4, 3, 2},
			out: true,
		},
		{
			ina: []int64{3, 1, 2},
			inb: []int64{5, 1, 4, 3, 0},
			out: false,
		},
		{
			ina: []int64{3, 1, 2, 6, 8, 10, 11},
			inb: []int64{5, 1, 4, 3, 0},
			out: false,
		},
		{
			ina: []int64{},
			inb: []int64{1},
			out: false,
		},
		{
			ina: []int64{1},
			inb: []int64{1},
			out: true,
		},
		{
			ina: []int64{},
			inb: []int64{},
			out: false,
		},
	}
	ass := assert.New(t)
	for _, tc := range testData {
		ass.Equal(tc.out, intSetBelong(tc.ina, tc.inb), "ina=%+v||inb=%+v", tc.ina, tc.inb)
	}
}

func TestIntSetInter(t *testing.T) {
	var testData = []struct {
		ina, inb []int64
		out      bool
	}{
		{
			ina: []int64{3, 1, 2},
			inb: []int64{5, 1, 4, 3, 2},
			out: true,
		},
		{
			ina: []int64{3, 1, 2},
			inb: []int64{5, 1, 4, 3, 0},
			out: true,
		},
		{
			ina: []int64{3, 1, 2, 6, 8, 10, 11},
			inb: []int64{5, 1, 4, 3, 0},
			out: true,
		},
		{
			ina: []int64{},
			inb: []int64{1},
			out: false,
		},
		{
			ina: []int64{1},
			inb: []int64{1},
			out: true,
		},
		{
			ina: []int64{},
			inb: []int64{},
			out: false,
		},
		{
			ina: []int64{1, 3, 5, 7, 9},
			inb: []int64{2, 4, 6, 8, 10, 5},
			out: true,
		},
		{
			ina: []int64{1, 3, 5, 7, 9},
			inb: []int64{2, 4, 6, 8, 10, 12},
			out: false,
		},
	}
	ass := assert.New(t)
	for _, tc := range testData {
		ass.Equal(tc.out, intSetInter(tc.ina, tc.inb), "ina=%+v||inb=%+v", tc.ina, tc.inb)
	}
}

func TestCmpInt_AllOps(t *testing.T) {
	var testData = []struct {
		actual, expect int64
		op             string
		out            bool
	}{
		{5, 5, opEqual, true},
		{5, 5, opIn, true},
		{5, 4, opNotEqual, true},
		{5, 5, opNotEqual, false},
		{5, 4, opNotEqualSynonym, true},
		{5, 5, opNotEqualSynonym, false},
		{5, 4, opNotIn, true},
		{5, 5, opNotIn, false},
		{5, 3, opLarger, true},
		{5, 5, opLarger, false},
		{5, 5, opLargerEqual, true},
		{5, 4, opLess, false},
		{5, 6, opLess, true},
		{5, 5, opLessEqual, true},
		{5, 4, opLessEqual, false},
		{5, 5, "unknown", false},
	}
	ass := assert.New(t)
	for _, tc := range testData {
		ass.Equal(tc.out, cmpInt(tc.actual, tc.expect, tc.op), "actual=%d||expect=%d||op=%s", tc.actual, tc.expect, tc.op)
	}
}

func TestCmpFloat_AllOps(t *testing.T) {
	var testData = []struct {
		actual, expect float64
		op             string
		out            bool
	}{
		{5, 5, opEqual, true},
		{5, 4, opNotEqual, true},
		{5, 4, opNotEqualSynonym, true},
		{5, 3, opLarger, true},
		{5, 5, opLargerEqual, true},
		{5, 4, opLess, false},
		{5, 6, opLess, true},
		{5, 5, opLessEqual, true},
		{5, 4, opLessEqual, false},
		{5, 5, "unknown", false},
	}
	ass := assert.New(t)
	for _, tc := range testData {
		ass.Equal(tc.out, cmpFloat(tc.actual, tc.expect, tc.op), "actual=%v||expect=%v||op=%s", tc.actual, tc.expect, tc.op)
	}
}

func TestCmpStr_AllOps(t *testing.T) {
	var testData = []struct {
		actual, expect string
		op             string
		out            bool
	}{
		{"b", "b", opEqual, true},
		{"b", "b", opIn, true},
		{"b", "a", opNotEqual, true},
		{"b", "b", opNotEqual, false},
		{"b", "a", opNotEqualSynonym, true},
		{"b", "b", opNotEqualSynonym, false},
		{"b", "a", opNotIn, true},
		{"b", "b", opNotIn, false},
		{"b", "a", opLarger, true},
		{"b", "c", opLarger, false},
		{"b", "b", opLargerEqual, true},
		{"b", "c", opLess, true},
		{"b", "b", opLess, false},
		{"b", "b", opLessEqual, true},
		{"b", "a", opLessEqual, false},
		{"b", "b", "unknown", false},
	}
	ass := assert.New(t)
	for _, tc := range testData {
		ass.Equal(tc.out, cmpStr(tc.actual, tc.expect, tc.op), "actual=%s||expect=%s||op=%s", tc.actual, tc.expect, tc.op)
	}
}

func TestShouldCompareSet(t *testing.T) {
	ass := assert.New(t)
	for _, op := range []string{opEqual, opNotEqual, opNotEqualSynonym, opInter, opNotInter, opIn, opNotIn} {
		ass.True(shouldCompareSet(op), "op=%s", op)
	}
	ass.False(shouldCompareSet("unknown"))
}

func TestCompareSet_InvalidInput(t *testing.T) {
	ass := assert.New(t)
	ass.False(compareSet(1, []string{"1"}, "unknown"))
	ass.False(compareSet(struct{}{}, []string{"1"}, opEqual))
}

func TestIntArr2i64Arr(t *testing.T) {
	ass := assert.New(t)
	ass.Nil(intArr2i64Arr(nil))
	ass.Equal([]int64{1, 2}, intArr2i64Arr([]int{1, 2}))
}

func TestCmpIntSet_Errors(t *testing.T) {
	ass := assert.New(t)
	ass.False(cmpIntSet([]int64{1}, []string{"1.5"}, opIn))
	ass.False(cmpIntSet([]int64{1}, []string{"1"}, opEqual))
}

func TestCmpFloatSet_Errors(t *testing.T) {
	ass := assert.New(t)
	ass.False(cmpFloatSet([]float64{1}, []string{"x"}, opIn))
	ass.False(cmpFloatSet([]float64{1}, []string{"1"}, opEqual))
}

func TestCmpStringSet_InvalidOp(t *testing.T) {
	assert.False(t, cmpStringSet([]string{"a"}, []string{"a"}, opEqual))
}

func TestSetComparisons_EmptyInputs(t *testing.T) {
	ass := assert.New(t)
	ass.False(floatSetInter(nil, nil))
	ass.False(floatSetInter([]float64{1}, nil))
	ass.False(strSetInter(nil, nil))
	ass.False(strSetInter([]string{"a"}, nil))
	ass.True(strSetNotInter([]string{"a"}, []string{"b"}))
	ass.False(strSetNotInter([]string{"a"}, []string{"a"}))
}
