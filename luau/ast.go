package luau

// Block is a basic container of nodes
type Block struct {
	List []Node
}

func (c *Block) Render(w Writer) {
	w.IncIndent()
	for _, n := range c.List {
		n.Render(w)
	}
	w.DecIndent()
}

// Identifier
type Ident struct {
	Name string
}

func (i *Ident) Render(w Writer) {
	w.Write(i.Name)
}

// Do-end block
// ex: do print("hi") end
type DoStmt struct {
	Block *Block
}

func (d *DoStmt) Render(w Writer) {
	w.End("do")
	d.Block.Render(w)
	w.End("end")
}

// While block
// ex: while true do print("hi") end
type WhileStmt struct {
	Exp   Node
	Block *Block
}

func (wh *WhileStmt) Render(w Writer) {
	w.Write("while ")
	wh.Exp.Render(w)
	w.End("do")

	wh.Block.Render(w)
	w.End("end")
}

// For block
// ex: for i = 1,10,1 do end
// ex: for i,v in pairs({1,2,3}) do end
type ForStmt struct {
	Block *Block
	Exp   Node
}

func (f *ForStmt) Render(w Writer) {

}

// Var declaration
// ex: local foo,bar = 4,2
type VarDecl struct {
	Scope  Scope
	Names  []*Ident
	Values []Node
}

func (v *VarDecl) Render(w Writer) {
	if v.Scope == LOCAL {
		w.Pre("local ")
	}
	for i, n := range v.Names {
		n.Render(w)
		if i != len(v.Names)-1 {
			w.Write(",")
		}
	}
	w.Write(" = ")
	for i, n := range v.Values {
		n.Render(w)
		if i != len(v.Values)-1 {
			w.Write(",")
		}
	}
	w.Write("\n")
}

// Function declaraion
// ex: function foo() end
type FuncDecl struct {
	Name   *Ident
	Params []*Ident
	Block  *Block
	Scope  Scope
}

func (f *FuncDecl) Render(w Writer) {
	if f.Scope == LOCAL {
		w.Pre("local ")
	}

	w.Write("function ")

	f.Name.Render(w)

	w.Write("(")
	for i, p := range f.Params {
		p.Render(w)
		if i != len(f.Params)-1 {
			w.Write(",")
		}
	}
	w.End(")")

	f.Block.Render(w)
	w.End("end")
}

// Function literal
// ex: local x = function() end
type FuncLit struct {
	Name   *Ident
	Params []*Ident
	Block  *Block
}

func (f *FuncLit) Render(w Writer) {
	w.Pre("function ")
	w.Write("(")

	for i, p := range f.Params {
		p.Render(w)
		if i != len(f.Params) {
			w.Write(",")
		}
	}
	w.End(")")

	f.Block.Render(w)
	w.End("end")
}

// Numeric literal
// ex: 152.123
type NumericLit struct {
	Value string
}

func (n *NumericLit) Render(w Writer) {
	w.Write(n.Value)
}

// String literal
// ex: "aa"
type StringLit struct {
	Value string
}

func (s *StringLit) Render(w Writer) {
	w.Write("\"" + s.Value + "\"")
}

// Call expression
// ex: foo()
type CallExpr struct {
	Args []Node
	Fun  Ident
}

func (c *CallExpr) Render(w Writer) {
	c.Fun.Render(w)
	for i, a := range c.Args {
		a.Render(w)
		if i != len(c.Args)-1 {
			w.Write(",")
		}
	}
}

// Index expression
// ex: table["index"]
type IndexExpr struct {
	Sub   Node
	Index Node
}

func (i *IndexExpr) Render(w Writer) {
	i.Sub.Render(w)
	w.Write("[")
	i.Index.Render(w)
	w.Write("]\n")
}

// Selector expression
// ex: table.property
type SelectorExpr struct {
	Sub Node
	Sel *Ident
}

func (s *SelectorExpr) Render(w Writer) {
	s.Sub.Render(w)
	w.Write(".")
	s.Sel.Render(w)
}

// Binary expression
// ex: 2 + 2
type BinaryExpr struct {
	Left  Node
	Right Node
	Op    Operator
}

func (b *BinaryExpr) Render(w Writer) {
	b.Left.Render(w)
	w.Write(FormatOperator(b.Op))
	b.Right.Render(w)
}

// Parenthesized expression
// ex: (2 + 2) * 2
type ParenExpr struct {
	Sub Node
}

func (p *ParenExpr) Render(w Writer) {
	w.Write("(")
	p.Sub.Render(w)
	w.Write(")")
}
