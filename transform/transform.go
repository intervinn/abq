package transform

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/intervinn/abq/luau"
)

// transform.Mod is a reserved call expression
// for rendering raw luau strings in where its written
func Mod[T any](value string) T {
	var a T
	return a
}

// transform.Require is reserved
// for making sure a certain compiled package is imported.
// It should be used as a top-level declaration
func Require(value string) any {
	return nil
}

func Token(t token.Token) luau.Token {
	switch t {
	case token.ADD:
		return luau.ADD
	case token.ADD_ASSIGN:
		return luau.ADD_ASSIGN
	case token.SUB:
		return luau.SUB
	case token.SUB_ASSIGN:
		return luau.SUB_ASSIGN
	case token.MUL:
		return luau.MUL
	case token.MUL_ASSIGN:
		return luau.MUL_ASSIGN
	case token.QUO:
		return luau.DIV
	case token.QUO_ASSIGN:
		return luau.DIV_ASSIGN
	case token.REM:
		return luau.REM
	case token.REM_ASSIGN:
		return luau.REM_ASSIGN
	}
	return luau.ILLEGAL
}

func Source(name string, src string) ([]luau.Node, error) {
	f, err := Parse(name, src)
	if err != nil {
		panic(err)
	}

	res := []luau.Node{}
	for _, d := range f.Decls {
		decl, err := Decl(d, f)
		if err != nil {
			return nil, err
		}

		res = append(res, decl)
	}
	return res, nil
}

func Parse(name string, src string) (*ast.File, error) {
	return parser.ParseFile(token.NewFileSet(), name, src, parser.AllErrors)
}

func Decl(d ast.Decl, f *ast.File) (luau.Node, error) {
	switch decl := d.(type) {
	case *ast.FuncDecl:
		return FuncDecl(decl, f)
	case *ast.GenDecl:
		return GenDecl(decl, f)
	}
	return nil, fmt.Errorf("unknown declaration: %#v", d)
}

func GenDecl(g *ast.GenDecl, f *ast.File) (luau.Node, error) {
	block := &luau.Block{}

	for _, s := range g.Specs {
		spec, err := Spec(s, f)
		if err != nil {
			return nil, err
		}
		block.List = append(block.List, spec)
	}

	return block, nil
}

func Spec(s ast.Spec, f *ast.File) (luau.Node, error) {
	switch spec := s.(type) {
	case *ast.ValueSpec:
		return ValueSpec(spec, f)
	case *ast.TypeSpec:
		return TypeSpec(spec, f)
	case *ast.ImportSpec:
		return &luau.Raw{Content: ""}, nil
	}
	return nil, fmt.Errorf("unknown spec: %#v", s)
}

func ValueSpec(v *ast.ValueSpec, f *ast.File) (luau.Node, error) {
	names := make([]luau.Node, len(v.Names))
	for i, v := range v.Names {
		names[i] = Ident(v, f)
	}

	values := make([]luau.Node, len(v.Values))
	for i, v := range v.Values {
		e, err := Expr(v, f)
		if err != nil {
			return nil, err
		}

		values[i] = e
	}

	// check if its a transform.Mod
	if len(v.Names) == 1 && len(v.Values) == 1 {
		value := v.Values[0]

		if c, ok := value.(*ast.CallExpr); ok {
			expr, err := CallExpr(c, f)

			if raw, ok := expr.(*luau.Raw); ok && err == nil {
				return raw, nil
			}
		}
	}

	return &luau.DeclStmt{
		Scope:  luau.GLOBAL,
		Names:  names,
		Values: values,
	}, nil
}

func TypeSpec(t *ast.TypeSpec, f *ast.File) (*luau.DeclStmt, error) {
	i := Ident(t.Name, f)
	return &luau.DeclStmt{
		Scope:  luau.LOCAL,
		Names:  []luau.Node{i},
		Values: []luau.Node{&luau.TableLit{Elts: []luau.Node{}}},
	}, nil
}

func FuncDecl(f *ast.FuncDecl, file *ast.File) (*luau.FuncStmt, error) {
	plist := f.Type.Params.List
	params := []*luau.Ident{}
	for _, ls := range plist {
		for _, p := range ls.Names {
			params = append(params, &luau.Ident{Name: p.Name})
		}
	}

	c, err := Chunk(f.Body, file)
	if err != nil {
		return nil, err
	}

	name := &luau.Ident{
		Name: f.Name.Name,
	}

	if f.Recv != nil {
		r := f.Recv.List[0]
		i := Ident(r.Names[0], file)

		params = append([]*luau.Ident{i}, params...)
		rtype, ok := r.Type.(*ast.StarExpr)
		if !ok {
			return nil, errors.New("function receivers must be pointers")
		}
		recver, ok := rtype.X.(*ast.Ident)
		if !ok {
			return nil, fmt.Errorf("receiver type must be identifier, got %#v", recver)
		}

		name.Name = recver.Name + "." + name.Name
	}

	return &luau.FuncStmt{
		Name:   name,
		Params: params,
		Chunk:  c,
		Scope:  luau.GLOBAL,
	}, nil
}

func Ident(i *ast.Ident, f *ast.File) *luau.Ident {
	return &luau.Ident{Name: i.Name}
}

func BasicLit(l *ast.BasicLit) (luau.Node, error) {
	switch l.Kind {
	case token.INT, token.FLOAT, token.IMAG:
		return &luau.NumericLit{Value: l.Value}, nil
	case token.CHAR, token.STRING:
		trim := l.Value[1 : len(l.Value)-1]
		return &luau.StringLit{Value: trim}, nil
	}
	return nil, fmt.Errorf("unknown literal: %#v", l)
}

func CompositeLit(l *ast.CompositeLit, f *ast.File) (luau.Node, error) {
	switch l.Type.(type) {
	case *ast.ArrayType, *ast.MapType, *ast.Ident:
		elts := make([]luau.Node, len(l.Elts))
		for i, v := range l.Elts {
			e, err := Expr(v, f)
			if err != nil {
				return nil, err
			}
			elts[i] = e
		}

		return &luau.TableLit{
			Elts: elts,
		}, nil
	case *ast.ChanType:
		return nil, errors.New("channels not supported")
	}
	return nil, fmt.Errorf("unknown composite literal: %#v", l)
}

func Chunk(b *ast.BlockStmt, f *ast.File) (*luau.Chunk, error) {
	ls := b.List
	result := make([]luau.Node, len(ls))
	for i, v := range ls {
		s, err := Stmt(v, f)
		if err != nil {
			return nil, err
		}
		result[i] = s
	}

	return &luau.Chunk{
		List: result,
	}, nil
}

func Expr(e ast.Expr, f *ast.File) (luau.Node, error) {
	switch expr := e.(type) {
	case *ast.BadExpr:
		return nil, errors.New("bad expression")
	case *ast.BinaryExpr:
		return BinaryExpr(expr, f)
	case *ast.CallExpr:
		return CallExpr(expr, f)
	case *ast.IndexExpr:
		return IndexExpr(expr, f)
	case *ast.ParenExpr:
		return ParenExpr(expr, f)
	case *ast.SelectorExpr:
		return SelectorExpr(expr, f)
	case *ast.BasicLit:
		return BasicLit(expr)
	case *ast.Ident:
		return Ident(expr, f), nil
	case *ast.KeyValueExpr:
		return KeyValueExpr(expr, f)
	case *ast.CompositeLit:
		return CompositeLit(expr, f)
	case *ast.UnaryExpr:
		return UnaryExpr(expr, f)
	case *ast.StarExpr:
		return nil, nil
	}
	return nil, fmt.Errorf("unknown expression: %#v", e)
}

func UnaryExpr(u *ast.UnaryExpr, f *ast.File) (luau.Node, error) {
	x, err := Expr(u.X, f)
	if err != nil {
		return nil, err
	}

	// if constructing a struct
	if _, ok := u.X.(*ast.CompositeLit); ok && u.Op == token.AND {
		return &luau.CallExpr{
			Fun: &luau.Ident{
				Name: "setmetatable",
			},
			Args: []luau.Node{
				&luau.TableLit{
					Elts: []luau.Node{},
				},
				x,
			},
		}, nil
	}

	return x, nil
}

func BinaryExpr(e *ast.BinaryExpr, f *ast.File) (*luau.BinaryExpr, error) {
	op := Token(e.Op)
	left, err := Expr(e.X, f)
	if err != nil {
		return nil, err
	}

	right, err := Expr(e.Y, f)
	if err != nil {
		return nil, err
	}

	{
		if _, srok := right.(*luau.StringLit); srok {
			if op == luau.ADD_ASSIGN {
				op = luau.CCT_ASSIGN
			}

			if op == luau.ADD {
				op = luau.CCT
			}
		}
	}

	return &luau.BinaryExpr{
		Left:  left,
		Right: right,
		Op:    op,
	}, nil
}

func KeyValueExpr(k *ast.KeyValueExpr, f *ast.File) (*luau.KeyValueExpr, error) {
	key, err := Expr(k.Key, f)
	if err != nil {
		return nil, err
	}
	value, err := Expr(k.Value, f)
	if err != nil {
		return nil, err
	}

	return &luau.KeyValueExpr{
		Key:   key,
		Value: value,
	}, nil
}

func CallExpr(c *ast.CallExpr, f *ast.File) (luau.Node, error) {
	fn, err := Expr(c.Fun, f)
	if err != nil {
		return nil, err
	}

	args := []luau.Node{}
	method := func(sl *ast.SelectorExpr) (luau.Node, error) {
		// check if its transform.Mod
		if id, ok := sl.X.(*ast.Ident); ok && id.Name == "transform" && sl.Sel.Name == "Mod" {
			ermsg := errors.New("transform.Mod must have exactly one string argument")
			if len(c.Args) != 1 {
				return nil, ermsg
			}
			if blit, ok := c.Args[0].(*ast.BasicLit); ok {
				if blit.Kind != token.STRING {
					return nil, ermsg
				}
				return &luau.Raw{
					Content: blit.Value[1 : len(blit.Value)-1],
				}, nil
			}
		}

		// if method has an object, add self arg
		if id, ok := sl.X.(*ast.Ident); ok && id.Obj != nil {
			args = append(args, Ident(id, f))
		}
		return nil, nil
	}

	// check if its a struct method
	if sl, ok := c.Fun.(*ast.SelectorExpr); ok {
		if node, err := method(sl); node != nil {
			return node, err
		}
	}

	if il, ok := c.Fun.(*ast.IndexExpr); ok {
		if sl, ok := il.X.(*ast.SelectorExpr); ok {
			if node, err := method(sl); node != nil {
				return node, err
			}
		}
	}

	for _, v := range c.Args {
		e, err := Expr(v, f)
		if err != nil {
			return nil, err
		}
		args = append(args, e)
	}

	return &luau.CallExpr{
		Fun:  fn,
		Args: args,
	}, nil
}
func IndexExpr(i *ast.IndexExpr, f *ast.File) (*luau.IndexExpr, error) {
	x, err := Expr(i.X, f)
	if err != nil {
		return nil, err
	}

	index, err := Expr(i.Index, f)
	if err != nil {
		return nil, err
	}

	return &luau.IndexExpr{
		Index: index,
		X:     x,
	}, nil
}
func ParenExpr(p *ast.ParenExpr, f *ast.File) (*luau.ParenExpr, error) {
	x, err := Expr(p.X, f)
	if err != nil {
		return nil, err
	}
	return &luau.ParenExpr{
		X: x,
	}, nil
}

func SelectorExpr(s *ast.SelectorExpr, f *ast.File) (*luau.SelectorExpr, error) {
	sel := Ident(s.Sel, f)
	x, err := Expr(s.X, f)
	if err != nil {
		return nil, err
	}

	return &luau.SelectorExpr{
		Sel: sel,
		X:   x,
	}, nil
}

func Stmt(s ast.Stmt, f *ast.File) (luau.Node, error) {
	switch stmt := s.(type) {
	case *ast.AssignStmt:
		return AssignStmt(stmt, f)
	case *ast.BlockStmt:
		return BlockStmt(stmt, f)
	case *ast.ExprStmt:
		return ExprStmt(stmt, f)
	case *ast.IfStmt:
		return IfStmt(stmt, f)
	case *ast.ReturnStmt:
		return ReturnStmt(stmt, f)
	}
	return nil, fmt.Errorf("unknown statement: %#v", s)
}

func IfStmt(i *ast.IfStmt, f *ast.File) (*luau.IfStmt, error) {
	cond, err := Expr(i.Cond, f)
	if err != nil {
		return nil, err
	}

	chunk, err := Chunk(i.Body, f)
	if err != nil {
		return nil, err
	}

	els, err := Stmt(i.Else, f)
	if err != nil {
		return nil, err
	}

	return &luau.IfStmt{
		Cond: cond,
		Body: chunk,
		Else: els,
	}, nil
}

func AssignStmt(a *ast.AssignStmt, f *ast.File) (luau.Node, error) {
	left := make([]luau.Node, len(a.Lhs))
	for i, v := range a.Lhs {
		e, err := Expr(v, f)
		if err != nil {
			return nil, err
		}

		left[i] = e
	}

	right := make([]luau.Node, len(a.Rhs))
	for i, v := range a.Rhs {
		e, err := Expr(v, f)
		if err != nil {
			return nil, err
		}

		right[i] = e
	}

	if a.Tok == token.DEFINE {
		return &luau.DeclStmt{
			Scope:  luau.LOCAL,
			Names:  left,
			Values: right,
		}, nil
	}

	return &luau.AssignStmt{
		Left:  left,
		Right: right,
	}, nil
}

func BlockStmt(b *ast.BlockStmt, f *ast.File) (*luau.DoStmt, error) {
	c, err := Chunk(b, f)
	if err != nil {
		return nil, err
	}
	return &luau.DoStmt{
		Chunk: c,
	}, nil
}

func ExprStmt(e *ast.ExprStmt, f *ast.File) (*luau.ExprStmt, error) {
	expr, err := Expr(e.X, f)
	if err != nil {
		return nil, err
	}
	return &luau.ExprStmt{
		X: expr,
	}, nil
}

func ReturnStmt(r *ast.ReturnStmt, f *ast.File) (*luau.ReturnStmt, error) {
	res := make([]luau.Node, len(r.Results))
	for i, v := range r.Results {
		e, err := Expr(v, f)
		if err != nil {
			return nil, err
		}

		res[i] = e
	}

	return &luau.ReturnStmt{
		Results: res,
	}, nil
}
