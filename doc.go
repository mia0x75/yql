/*
Package yql implements a small, SQL-like boolean expression language for rule
engines.

# Syntax

Expressions combine comparisons with and, or, and parentheses. The and
operator binds more tightly than or.

Comparisons support:

	=  !=  <>  >  >=  <  <=
	in  !in
	intersect  !intersect

Examples:

	name='deen' and age>=23 and (hobby in ('soccer', 'swim') or score>90)
	score intersect (7, 1, 9, 5, 3)
	score.sum() > 10

String values must be quoted, as in name='deen'. Numbers and booleans are
written directly, as in age>=23 and enabled=true.

# Data types

Values in the data map must be one of int, int64, float64, string, bool, or
json.Number. Set operators also accept slices of int, int64, float64, or string.
Aggregate helpers accept numeric slices; count also accepts any slice.

Comparison is type-aware: a value on the right of a comparison is parsed
according to the dynamic type of the field on the left.

# Set operators

For slice values, in reports whether every element of the field is present in
the expected set, while intersect reports whether at least one element is
present. !in and !intersect negate those results.

# Helpers

Helpers use field.func() syntax and can be chained. The supported aggregates are
sum, avg, max, min, and count. sum, avg, max, and min operate on numeric slices;
count accepts any slice and returns 1 for scalar values.

See the package examples and tests for more details.
*/
package yql
