package luau

type Token int

const (
	ILLEGAL Token = iota

	NUMBER
	STRING

	IDENT // foo

	ADD // +
	SUB // -
	MUL // *
	DIV // /
	REM // %

	ADD_ASSIGN // +=
	SUB_ASSIGN // -=
	MUL_ASSIGN // *=
	DIV_ASSIGN // /=
	REM_ASSIGN // %=

	AND // and
	OR  // or
	NOT // not

	NEQ // !=
	LEQ // <=
	GEQ // >=

	EQL    // ==
	LSS    // <
	GTR    // >
	ASSIGN // =

	RPAREN // )
	LPAREN // (
	RBRACE // }
	LBRACE // {
)
