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

	NEQ // ~=
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

// Usable for only certain token types
func FormatToken(t Token) string {
	switch t {
	case ADD:
		return "+"
	case SUB:
		return "-"
	case MUL:
		return "*"
	case REM:
		return "%"
	case ADD_ASSIGN:
		return "+="
	case SUB_ASSIGN:
		return "-="
	case MUL_ASSIGN:
		return "*="
	case DIV_ASSIGN:
		return "/="
	case REM_ASSIGN:
		return "%="
	case AND:
		return "and"
	case OR:
		return "or"
	case NOT:
		return "not"
	case NEQ:
		return "~="
	case LEQ:
		return "<="
	case GEQ:
		return ">="
	case EQL:
		return "=="
	case LSS:
		return "<"
	case GTR:
		return ">"
	case ASSIGN:
		return "="
	case RPAREN:
		return ")"
	case LPAREN:
		return "("
	case RBRACE:
		return "}"
	case LBRACE:
		return "{"
	}
	return ""
}
