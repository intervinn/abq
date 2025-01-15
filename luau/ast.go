package luau

import "strings"

// A basic literal for strings and numberss
type BasicLit struct {
	Kind  Token
	Value string
}

func (b *BasicLit) Render() string {
	return b.Value
}

// Anonymous function literal
type FuncLit struct {
	Type *FuncType
	Body *BlockStmt
}

// Assign statement, do not confuse with declaration statement
// Ex: foo = 5
type AssignStmt struct {
	Tok Token
	Rhs []Node   // lit
	Lhs []*Ident // ident
}

func (a *AssignStmt) Render() string {
	rhs := make([]string, len(a.Rhs))
	for i, h := range a.Rhs {
		rhs[i] = h.Render()
	}

	lhs := make([]string, len(a.Lhs))
	for i, h := range a.Lhs {
		lhs[i] = h.Render()
	}

	return strings.Join(lhs, ",") + " = " + strings.Join(rhs, ",")
}

// Variable declaration statement
// Ex: local foo = 5
type DeclStmt struct {
	Idt  *Ident // ident
	Val  Node
	Type Node
}

func (d *DeclStmt) Render() string {
	return "local" + d.Idt.Render() + " = " + d.Val.Render()
}

// An expression statement
type ExprStmt struct {
	Expr Expr
}

func (e *ExprStmt) Render() string {
	return e.Expr.Render()
}

// Block statement holds a list of statements within
// and is usable as a function body, loop body, if statement body, and a do-end block
type BlockStmt struct {
	List []Stmt
}

func (b *BlockStmt) Render() {

}

// An identifier that holds the name of variables, fields, and whatever
type Ident struct {
	Name string
}

func (i *Ident) Render() string {
	return i.Name
}

// Do not confuse with FuncLit, implements rendering of a function type definition
type FuncType struct {
	Params  []Field
	Results []Field
}

func (f *FuncType) Render() string {
	params := make([]string, len(f.Params))
	for i, p := range f.Params {
		params[i] = p.Render()
	}

	results := make([]string, len(f.Params))
	for i, r := range f.Results {
		params[i] = r.Render()
	}

	return "(" + strings.Join(params, ",") + ")" + " -> (" + strings.Join(results, ",") + ")"
}

// Field is a common syntax used when declaring function parameters
// and table type keys
type Field struct {
	Type Node
	Name *Ident
}

func (f *Field) Render() string {
	return f.Name.Render() + ": " + f.Type.Render()
}

// Basic Types
type (
	StringType   struct{}
	NumberType   struct{}
	BoolType     struct{}
	NilType      struct{}
	TableType    struct{}
	ThreadType   struct{}
	UserDataType struct{}
	VectorType   struct{}
	BufferType   struct{}
)

func (s *StringType) Render() string {
	return "string"
}
func (n *NumberType) Render() string {
	return "number"
}
func (b *BoolType) Render() string {
	return "boolean"
}
func (u *UserDataType) Render() string {
	return "userdata"
}
func (v *VectorType) Render() string {
	return "vector"
}
func (b *BufferType) Render() string {
	return "buffer"
}
func (t *ThreadType) Render() string {
	return "thread"
}
func (t *TableType) Render() string {
	return "table"
}
func (n *NilType) Render() string {
	return "nil"
}
