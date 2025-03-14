package luau

import "slices"

// ==================================
// Base node
type Node interface {
	Render(w Writer)
}

// ==================================
type Writer interface {
	// Pre will add the indent but won't add the newline
	Pre(s string) error
	// Write will just add the string
	Write(s string) error

	Indent() int
	IncIndent()
	DecIndent()
}

// ==================================
type File struct {
	Name  string
	Decls []Node
}

func (f *File) Render(w Writer) {
	// decl order: raw -> decl -> func -> [the rest]
	rendered := []int{}

	for i, v := range f.Decls {
		if _, ok := v.(*Raw); ok {
			v.Render(w)
			rendered = append(rendered, i)
		}
	}

	for i, v := range f.Decls {
		if _, ok := v.(*DeclStmt); ok {
			v.Render(w)
			rendered = append(rendered, i)
		}
	}

	for i, v := range f.Decls {
		if _, ok := v.(*FuncStmt); ok {
			v.Render(w)
			rendered = append(rendered, i)
		}
	}

	for i, v := range f.Decls {
		if !slices.Contains(rendered, i) {
			v.Render(w)
		}
	}

}

// ==================================
type Scope int

const (
	GLOBAL Scope = iota
	LOCAL
	NONE
)

// ==================================
type Token int

const (
	ILLEGAL Token = iota

	ADD  // +
	SUB  // -
	MUL  // *
	DIV  // /
	FDIV // //
	REM  // %
	POW  // ^
	CCT  // ..

	ADD_ASSIGN  // +=
	SUB_ASSIGN  // -=
	MUL_ASSIGN  // *=
	DIV_ASSIGN  // /=
	FDIV_ASSIGN // //=
	REM_ASSIGN  // %=
	POW_ASSIGN  // ^=
	CCT_ASSIGN  // ..=

	AND // and
	OR  // or
	NOT // not

	NEQ // ~=
	LEQ // <=
	GEQ // >=

	EQL // ==
	LSS // <
	GTR // >

	LEN // #

	BREAK    // break
	CONTINUE // continue
)

func FormatToken(o Token) string {
	switch o {
	case ADD:
		return "+"
	case SUB:
		return "-"
	case MUL:
		return "*"
	case DIV:
		return "/"
	case FDIV:
		return "//"
	case REM:
		return "%"
	case POW:
		return "^"
	case CCT:
		return ".."
	case ADD_ASSIGN:
		return "+="
	case SUB_ASSIGN:
		return "-="
	case MUL_ASSIGN:
		return "*="
	case DIV_ASSIGN:
		return "/="
	case FDIV_ASSIGN:
		return "//="
	case REM_ASSIGN:
		return "%="
	case POW_ASSIGN:
		return "^="
	case CCT_ASSIGN:
		return "..="
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
	case LEN:
		return "#"
	default:
		return ""
	}
}
