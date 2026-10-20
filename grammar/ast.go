package grammar

// Expr is a parsed YQL expression.
type Expr interface {
	isExpr()
}

// Query is the root node of a parsed YQL expression.
type Query struct {
	Expr Expr
}

// BinaryExpr combines two expressions with AND or OR.
type BinaryExpr struct {
	Op    string
	Left  Expr
	Right Expr
}

func (*BinaryExpr) isExpr() {}

// CompareExpr compares a field against one or more values.
type CompareExpr struct {
	Field  string
	Funcs  []string
	Op     string
	Values []string
}

func (*CompareExpr) isExpr() {}
