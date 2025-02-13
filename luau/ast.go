package luau

// Block is a basic container of nodes
type Block struct {
	List []Node
}

func (c *Block) Render(w Writer) {
	for _, n := range c.List {
		w.Line(n.Render(w))
	}
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
type DoBlock struct {
	Block *Block
}

func (d *DoBlock) Render(w Writer) {
	w.Line("do")
	d.Block.Render(w)
	w.Line("end")
}

// While block
// ex: while true do print("hi") end
type WhileBlock struct {
	Exp   Node
	Block *Block
}

func (wh *WhileBlock) Render(w Writer) {
	w.Write("while ")
	wh.Exp.Render(w)
	w.Line("do")

	wh.Block.Render(w)
	w.Line("end")
}

// For block
// ex: for i = 1,10,1 do end
// ex: for i,v in pairs({1,2,3}) do end
type ForBlock struct {
	Block *Block
	Exp   Node
}

func (f *ForBlock) Render(w Writer) {

}

// Var declaration
// ex: local foo,bar = 4,2
type VarDecl struct {
	Scope  Scope
	Names  []Ident
	Values []Node
}

func (v *VarDecl) Render(w Writer) {
	if v.Scope == LOCAL {
		w.Write("local ")
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
}

// Function
// Works both as a declaration and a literal
// ex: local foo = function()
// ex: function foo()
type Function struct {
	Name   Node
	Params []Ident
	Block  *Block
	Scope  Scope
}

func (f *Function) Render(w Writer) {
	if f.Scope == LOCAL {
		w.Write("local ")
	}

	w.Write("function")

	if f.Scope != NONE {
		w.Write(f.Name.Render(w))
	}

	w.Write("(")
	for i, p := range f.Params {
		p.Render(w)
		if i != len(f.Params)-1 {
			w.Write(",")
		}
	}
	w.Line(")")

	f.Block.Render(w)
	w.Line("end")
}
