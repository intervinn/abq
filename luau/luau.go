package luau

// Base node
type Node interface {
	Render(w Writer) string
}

type Writer interface {
	Write(s string) error
	Line(s string) error
	Indent() int
}

type Scope int

const (
	GLOBAL Scope = iota
	LOCAL
	NONE
)
