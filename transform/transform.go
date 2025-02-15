package transform

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/intervinn/abq/luau"
)

func parse(name string, src string) (*ast.File, error) {
	return parser.ParseFile(token.NewFileSet(), name, src, parser.AllErrors)
}

func Source(name string, src string) {
	f, err := parse(name, src)
	if err != nil {
		panic(err)
	}

	for _, d := range f.Decls {
		decl(d)
	}
}

func decl(decl ast.Decl) (luau.Node, error) {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		return fdecl(d), nil
	default:
		return nil, errors.New("unknown declaration")
	}
}

func fdecl(f *ast.FuncDecl) *luau.FuncStmt {
	plist := f.Type.Params.List
	params := []*luau.Ident{}
	for _, ls := range plist {
		for _, p := range ls.Names {
			params = append(params, &luau.Ident{Name: p.Name})
		}
	}

	return &luau.FuncStmt{
		Name: &luau.Ident{
			Name: f.Name.Name,
		},
		Params: params,
		Block:  block(f.Body),
		Scope:  luau.GLOBAL,
	}
}

func block(b *ast.BlockStmt) *luau.Block {
	return nil
}

/*
func stmt(s *ast.Stmt) {

}
*/
