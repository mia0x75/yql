%{
package grammar

type fieldExpr struct {
	name  string
	funcs []string
}
%}

%union {
	expr  Expr
	field *fieldExpr
	str   string
	strs  []string
}

%token <str> FIELDNAME STRING INT FLOAT TRUE FALSE FUNC
%token <str> CMPOP // = != <> > >= < <=
%token <str> SETOP // in !in intersect !intersect
%token AND OR LPAREN RPAREN COMMA DOT

%type <expr> query expr boolExpr
%type <field> leftExpr
%type <str> value
%type <strs> valueList

%left OR
%left AND

%%

query:
	expr
	{
		yylex.(*lexer).query = &Query{Expr: $1}
	}
;

expr:
	boolExpr
	{
		$$ = $1
	}
	| expr AND expr
	{
		$$ = &BinaryExpr{Op: "AND", Left: $1, Right: $3}
	}
	| expr OR expr
	{
		$$ = &BinaryExpr{Op: "OR", Left: $1, Right: $3}
	}
	| LPAREN expr RPAREN
	{
		$$ = $2
	}
;

boolExpr:
	leftExpr CMPOP value
	{
		$$ = &CompareExpr{
			Field:  $1.name,
			Funcs:  $1.funcs,
			Op:     $2,
			Values: []string{$3},
		}
	}
	| leftExpr SETOP LPAREN valueList RPAREN
	{
		$$ = &CompareExpr{
			Field:  $1.name,
			Funcs:  $1.funcs,
			Op:     $2,
			Values: $4,
		}
	}
;

leftExpr:
	FIELDNAME
	{
		$$ = &fieldExpr{name: $1}
	}
	| leftExpr DOT FUNC LPAREN RPAREN
	{
		$1.funcs = append($1.funcs, $3)
		$$ = $1
	}
;

value:
	STRING
	{
		$$ = $1
	}
	| INT
	{
		$$ = $1
	}
	| FLOAT
	{
		$$ = $1
	}
	| TRUE
	{
		$$ = $1
	}
	| FALSE
	{
		$$ = $1
	}
;

valueList:
	value
	{
		$$ = []string{$1}
	}
	| valueList COMMA value
	{
		$$ = append($1, $3)
	}
;

%%
