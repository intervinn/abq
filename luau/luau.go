package luau

// Base node
type Node interface {
	Render() string
}

// Statement
type Stmt interface {
	Node
}

// Expression
type Expr interface {
	Node
}
