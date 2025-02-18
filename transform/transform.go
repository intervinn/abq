package transform

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/intervinn/abq/luau"
)

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
		decl, err := Decl(d)
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

func Decl(d ast.Decl) (luau.Node, error) {
	switch decl := d.(type) {
	case *ast.FuncDecl:
		return FuncDecl(decl)
	case *ast.GenDecl:
		return GenDecl(decl)
	}
	return nil, fmt.Errorf("unknown declaration: %#v", d)
}

func GenDecl(g *ast.GenDecl) (luau.Node, error) {
	block := &luau.Block{}

	for _, s := range g.Specs {
		spec, err := Spec(s)
		if err != nil {
			return nil, err
		}
		block.List = append(block.List, spec)
	}

	return block, nil
}

func Spec(s ast.Spec) (luau.Node, error) {
	switch spec := s.(type) {
	case *ast.ValueSpec:
		return ValueSpec(spec)
	case *ast.TypeSpec:
		return TypeSpec(spec)
	}
	return nil, fmt.Errorf("unknown spec: %#v", s)
}

func ValueSpec(v *ast.ValueSpec) (*luau.DeclStmt, error) {
	names := make([]luau.Node, len(v.Names))
	for i, v := range v.Names {
		names[i] = Ident(v)
	}

	values := make([]luau.Node, len(v.Values))
	for i, v := range v.Values {
		e, err := Expr(v)
		if err != nil {
			return nil, err
		}

		values[i] = e
	}

	return &luau.DeclStmt{
		Scope:  luau.GLOBAL,
		Names:  names,
		Values: values,
	}, nil
}

func TypeSpec(t *ast.TypeSpec) (*luau.DeclStmt, error) {
	i := Ident(t.Name)
	return &luau.DeclStmt{
		Scope:  luau.LOCAL,
		Names:  []luau.Node{i},
		Values: []luau.Node{&luau.TableLit{Elts: []luau.Node{}}},
	}, nil
}

func FuncDecl(f *ast.FuncDecl) (*luau.FuncStmt, error) {
	plist := f.Type.Params.List
	params := []*luau.Ident{}
	for _, ls := range plist {
		for _, p := range ls.Names {
			params = append(params, &luau.Ident{Name: p.Name})
		}
	}

	c, err := Chunk(f.Body)
	if err != nil {
		return nil, err
	}

	name := &luau.Ident{
		Name: f.Name.Name,
	}

	if f.Recv != nil {
		r := f.Recv.List[0]
		i := Ident(r.Names[0])

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

func Ident(i *ast.Ident) *luau.Ident {
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

func CompositeLit(l *ast.CompositeLit) (luau.Node, error) {
	switch l.Type.(type) {
	case *ast.ArrayType, *ast.MapType, *ast.Ident:
		elts := make([]luau.Node, len(l.Elts))
		for i, v := range l.Elts {
			e, err := Expr(v)
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

func Chunk(b *ast.BlockStmt) (*luau.Chunk, error) {
	ls := b.List
	result := make([]luau.Node, len(ls))
	for i, v := range ls {
		s, err := Stmt(v)
		if err != nil {
			return nil, err
		}
		result[i] = s
	}

	return &luau.Chunk{
		List: result,
	}, nil
}

func Expr(e ast.Expr) (luau.Node, error) {
	switch expr := e.(type) {
	case *ast.BadExpr:
		return nil, errors.New("bad expression")
	case *ast.BinaryExpr:
		return BinaryExpr(expr)
	case *ast.CallExpr:
		return CallExpr(expr)
	case *ast.IndexExpr:
		return IndexExpr(expr)
	case *ast.ParenExpr:
		return ParenExpr(expr)
	case *ast.SelectorExpr:
		return SelectorExpr(expr)
	case *ast.BasicLit:
		return BasicLit(expr)
	case *ast.Ident:
		return Ident(expr), nil
	case *ast.KeyValueExpr:
		return KeyValueExpr(expr)
	case *ast.CompositeLit:
		return CompositeLit(expr)
	case *ast.UnaryExpr:
		return UnaryExpr(expr)
	}
	return nil, fmt.Errorf("unknown expression: %#v", e)
}

func UnaryExpr(u *ast.UnaryExpr) (luau.Node, error) {
	x, err := Expr(u.X)
	if err != nil {
		return nil, err
	}

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

func BinaryExpr(e *ast.BinaryExpr) (*luau.BinaryExpr, error) {
	op := Token(e.Op)
	left, err := Expr(e.X)
	if err != nil {
		return nil, err
	}

	right, err := Expr(e.Y)
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

func KeyValueExpr(k *ast.KeyValueExpr) (*luau.KeyValueExpr, error) {
	key, err := Expr(k.Key)
	if err != nil {
		return nil, err
	}
	value, err := Expr(k.Value)
	if err != nil {
		return nil, err
	}

	return &luau.KeyValueExpr{
		Key:   key,
		Value: value,
	}, nil
}

func CallExpr(c *ast.CallExpr) (*luau.CallExpr, error) {
	f, err := Expr(c.Fun)
	if err != nil {
		return nil, err
	}

	args := make([]luau.Node, len(c.Args))
	for i, v := range c.Args {
		e, err := Expr(v)
		if err != nil {
			return nil, err
		}
		args[i] = e
	}

	return &luau.CallExpr{
		Fun:  f,
		Args: args,
	}, nil
}
func IndexExpr(i *ast.IndexExpr) (*luau.IndexExpr, error) {
	x, err := Expr(i.X)
	if err != nil {
		return nil, err
	}

	index, err := Expr(i.Index)
	if err != nil {
		return nil, err
	}

	return &luau.IndexExpr{
		Index: index,
		X:     x,
	}, nil
}
func ParenExpr(p *ast.ParenExpr) (*luau.ParenExpr, error) {
	x, err := Expr(p.X)
	if err != nil {
		return nil, err
	}
	return &luau.ParenExpr{
		X: x,
	}, nil
}

func SelectorExpr(s *ast.SelectorExpr) (*luau.SelectorExpr, error) {
	sel := Ident(s.Sel)
	x, err := Expr(s.X)
	if err != nil {
		return nil, err
	}

	return &luau.SelectorExpr{
		Sel: sel,
		X:   x,
	}, nil
}

func Stmt(s ast.Stmt) (luau.Node, error) {
	switch stmt := s.(type) {
	case *ast.AssignStmt:
		return AssignStmt(stmt)
	case *ast.BlockStmt:
		return BlockStmt(stmt)
	case *ast.ExprStmt:
		return ExprStmt(stmt)
	case *ast.IfStmt:
		return IfStmt(stmt)
	case *ast.ReturnStmt:
		return ReturnStmt(stmt)
	}
	return nil, fmt.Errorf("unknown statement: %#v", s)
}

func IfStmt(i *ast.IfStmt) (*luau.IfStmt, error) {
	cond, err := Expr(i.Cond)
	if err != nil {
		return nil, err
	}

	chunk, err := Chunk(i.Body)
	if err != nil {
		return nil, err
	}

	els, err := Stmt(i.Else)
	if err != nil {
		return nil, err
	}

	return &luau.IfStmt{
		Cond: cond,
		Body: chunk,
		Else: els,
	}, nil
}

func AssignStmt(a *ast.AssignStmt) (luau.Node, error) {
	left := make([]luau.Node, len(a.Lhs))
	for i, v := range a.Lhs {
		e, err := Expr(v)
		if err != nil {
			return nil, err
		}

		left[i] = e
	}

	right := make([]luau.Node, len(a.Rhs))
	for i, v := range a.Rhs {
		e, err := Expr(v)
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

func BlockStmt(b *ast.BlockStmt) (*luau.DoStmt, error) {
	c, err := Chunk(b)
	if err != nil {
		return nil, err
	}
	return &luau.DoStmt{
		Chunk: c,
	}, nil
}

func ExprStmt(e *ast.ExprStmt) (*luau.ExprStmt, error) {
	expr, err := Expr(e.X)
	if err != nil {
		return nil, err
	}
	return &luau.ExprStmt{
		X: expr,
	}, nil
}

func ReturnStmt(r *ast.ReturnStmt) (*luau.ReturnStmt, error) {
	res := make([]luau.Node, len(r.Results))
	for i, v := range r.Results {
		e, err := Expr(v)
		if err != nil {
			return nil, err
		}

		res[i] = e
	}

	return &luau.ReturnStmt{
		Results: res,
	}, nil
}
