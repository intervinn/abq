package luau

type ImportDecl struct {
	Path string
	As   string
}

// Raw renderable string
type Raw struct {
	Content string
}

func (r *Raw) Render(w Writer) {
	w.Write(r.Content)
	w.Write("\n")
}

// Chunk is a basic container of nodes with indentation considered
type Chunk struct {
	List []Node
}

func (b *Chunk) Render(w Writer) {
	w.IncIndent()
	for _, n := range b.List {
		n.Render(w)
	}
	w.DecIndent()
}

// Block is meant for internal use, whenever one node should be transformed into multiple
type Block struct {
	List []Node
}

func (c *Block) Render(w Writer) {
	for _, n := range c.List {
		n.Render(w)
	}
}

// Identifier
type Ident struct {
	Name string
}

func (i *Ident) Render(w Writer) {
	w.Write(i.Name)
}

// If statement
// ex: if true then end
type IfStmt struct {
	Cond Node
	Body *Chunk
	Else Node
}

func (i *IfStmt) Render(w Writer) {
	w.Pre("if ")
	i.Cond.Render(w)
	w.Write(" then\n")

	i.Body.Render(w)

	if elseif, ok := i.Else.(*IfStmt); ok {
		w.Pre("elseif ")
		elseif.Cond.Render(w)
		w.Write(" then\n")
		elseif.Body.Render(w)
		w.Pre("else")
		elseif.Else.Render(w)
	} else {
		w.Pre("else\n")
		i.Else.Render(w)
	}

	w.Pre("end\n")
}

// Do-end Chunk
// ex: do print("hi") end
type DoStmt struct {
	Chunk *Chunk
}

func (d *DoStmt) Render(w Writer) {
	w.Pre("do\n")
	d.Chunk.Render(w)
	w.Pre("end\n")
}

// While Chunk
// ex: while true do print("hi") end
type WhileStmt struct {
	Exp   Node
	Chunk *Chunk
}

func (wh *WhileStmt) Render(w Writer) {
	w.Pre("while ")
	wh.Exp.Render(w)
	w.Write("do\n")

	wh.Chunk.Render(w)
	w.Pre("end\n")
}

// Numeric For Statement
// ex: for i = 1,10,1 do end
type NumericForStmt struct {
	Chunk *Chunk
	Init  Node
	Cond  Node
	End   Node
}

func (n *NumericForStmt) Render(w Writer) {
	w.Pre("for ")
	n.Init.Render(w)
	w.Write(",")
	n.Cond.Render(w)
	w.Write(",")
	n.End.Render(w)
	w.Write(" do\n")

	n.Chunk.Render(w)
	w.Pre("end\n")
}

// Generic for statement
// ex: for i,v in pairs({1,2,3}) do end
type GenericForStmt struct {
	Chunk  *Chunk
	Idents []*Ident
	Iter   Node // iterator, ex: pairs({1,2,3})
}

func (g *GenericForStmt) Render(w Writer) {
	w.Pre("for ")
	for i, v := range g.Idents {
		v.Render(w)
		if i != len(g.Idents)-1 {
			w.Write(",")
		}
	}

	w.Write(" in ")
	g.Iter.Render(w)
	w.Write(" do\n")
	g.Chunk.Render(w)
	w.Pre("end")
}

// Var declaration
// ex: local foo,bar = 4,2
type DeclStmt struct {
	Scope  Scope
	Names  []Node
	Values []Node
}

func (v *DeclStmt) Render(w Writer) {
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

// Assignment statement
// ex: a = 5
type AssignStmt struct {
	Left  []Node
	Right []Node
}

func (a *AssignStmt) Render(w Writer) {
	for i, p := range a.Left {
		p.Render(w)
		if i != len(a.Left)-1 {
			w.Write(",")
		}
	}
	w.Write(" = ")
	for i, p := range a.Right {
		p.Render(w)
		if i != len(a.Right)-1 {
			w.Write(",")
		}
	}
	w.Write("\n")
}

// Expression statement
// ex: print("hello")
type ExprStmt struct {
	X Node
}

func (e *ExprStmt) Render(w Writer) {
	w.Pre("")
	e.X.Render(w)
	w.Write("\n")
}

// Return statement
// ex: return 4,2
type ReturnStmt struct {
	Results []Node
}

func (r *ReturnStmt) Render(w Writer) {
	w.Pre("return ")
	for i, rs := range r.Results {
		rs.Render(w)
		if i != len(r.Results)-1 {
			w.Write(",")
		}
	}
	w.Write("\n")
}

// Function declaraion
// ex: function foo() end
type FuncStmt struct {
	Name   *Ident
	Params []*Ident
	Chunk  *Chunk
	Scope  Scope
}

func (f *FuncStmt) Render(w Writer) {
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
	w.Write(")\n")

	f.Chunk.Render(w)
	w.Write("end\n")
}

// Function literal
// ex: local x = function() end
type FuncLit struct {
	Name   *Ident
	Params []*Ident
	Chunk  *Chunk
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
	w.Write(")\n")

	f.Chunk.Render(w)
	w.Write("end\n")
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

// Table literal
// ex: {a = 5}
type TableLit struct {
	Elts []Node
}

// check if theres a key-value expr
func haskv(t *TableLit) bool {
	for _, v := range t.Elts {
		if _, ok := v.(*KeyValueExpr); ok {
			return true
		}
	}
	return false
}

func (t *TableLit) Render(w Writer) {
	w.Write("{")
	if haskv(t) {
		w.Write("\n")
		w.IncIndent()
		for i, e := range t.Elts {
			w.Pre("")
			e.Render(w)
			if i != len(t.Elts)-1 {
				w.Write(",")
			}
			w.Write("\n")
		}
		w.DecIndent()
		w.Pre("}")
	} else {
		for i, e := range t.Elts {
			e.Render(w)
			if i != len(t.Elts)-1 {
				w.Write(", ")
			}
		}
		w.Write("}")
	}
}

// Call expression
// ex: foo()
type CallExpr struct {
	Args []Node
	Fun  Node
}

func (c *CallExpr) Render(w Writer) {
	c.Fun.Render(w)
	w.Write("(")
	for i, a := range c.Args {
		a.Render(w)
		if i != len(c.Args)-1 {
			w.Write(",")
		}
	}
	w.Write(")")
}

// Index expression
// ex: table["index"]
type IndexExpr struct {
	X     Node
	Index Node
}

func (i *IndexExpr) Render(w Writer) {
	i.X.Render(w)
	w.Write("[")
	i.Index.Render(w)
	w.Write("]")
}

// Selector expression
// ex: table.property
type SelectorExpr struct {
	X   Node
	Sel *Ident
}

func (s *SelectorExpr) Render(w Writer) {
	s.X.Render(w)
	w.Write(".")
	s.Sel.Render(w)
}

// Binary expression
// ex: 2 + 2
type BinaryExpr struct {
	Left  Node
	Right Node
	Op    Token
}

func (b *BinaryExpr) Render(w Writer) {
	b.Left.Render(w)
	w.Write(FormatToken(b.Op))
	b.Right.Render(w)
}

// Parenthesized expression
// ex: (2 + 2) * 2
type ParenExpr struct {
	X Node
}

func (p *ParenExpr) Render(w Writer) {
	w.Write("(")
	p.X.Render(w)
	w.Write(")")
}

// Key value expression
// ex: ["foo"] = 5
type KeyValueExpr struct {
	Key   Node
	Value Node
}

func (k *KeyValueExpr) Render(w Writer) {
	if ident, ok := k.Key.(*Ident); ok {
		ident.Render(w)
	} else {
		w.Write("[")
		k.Key.Render(w)
		w.Write("]")
	}

	w.Write(" = ")
	k.Value.Render(w)
}
