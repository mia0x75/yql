grammar Yql;

query:
	  expr
	;

expr:
	  booleanExpr       #boolExpr
	| expr AND expr     #andExpr
	| expr OR expr      #orExpr
	| '(' expr ')'      #embbedExpr
	;

booleanExpr:
	  leftExpr op=('='|'!='|'<>'|'>'|'<'|'>='|'<=') value
	| leftExpr op=('in'|'!in'|'∩'|'!∩') '(' value (',' value)* ')'
	;

leftExpr:
	  FIELDNAME (('.' FUNC '()')+)?
	;

value:
	  STRING
	| INT
	| FLOAT
	| TRUE
	| FALSE
	;

TRUE:
	  [Tt][Rr][Uu][Ee]
	;
FALSE:
	  [Ff][Aa][Ll][Ss][Ee]
	;
AND:
	  [Aa][Nn][Dd]
	;

OR:
	  [Oo][Rr]
	;

FUNC: 
	  COUNT
	| SUM
	| AVG
	| MAX
	| MIN
	;

COUNT:
	  [Cc][Oo][Uu][Nn][Tt]
	;

SUM:
	  [Ss][Uu][Mm]
	;

AVG:
	  [Aa][Vv][Gg]
	;

MAX:
	  [Mm][Aa][Xx]
	;

MIN:
	  [Mm][Ii][Nn]
	;

FIELDNAME:
	  ([a-zA-Z]|'_')+
	;

STRING:
	  '\'' .*? '\''
	| '"' .*? '"'
	;

fragment DIGIT:
	  [0-9]
	;

INT:
	  ('+'|'-')? DIGIT+
	;

FLOAT:
	  ('+'|'-')? DIGIT+ '.' DIGIT*
	;

WS:
	  [ \t\r\n]+ -> skip
	;

