package transform

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func parse(name string, src string) (*ast.File, error) {
	return parser.ParseFile(token.NewFileSet(), name, src, parser.AllErrors)
}

/*
func Source(name string, src string) {
	f, err := parse(name, src)
	if err != nil {
		panic(err)
	}

}

func decls(d []ast.Decl) {

}

func fdecl(f *ast.FuncDecl) *luau.FuncDecl {

}

func block(b *ast.BlockStmt) *luau.Block {

}

func stmt(s *ast.Stmt) {

}
*/
