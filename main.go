package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

var (
	FILE = os.Args[1]
	node *ast.File
	err  error
	file *os.File
	//filestr string
	fset  *token.FileSet
	types = map[string]string{
		"int":    "number",
		"int32":  "number",
		"int64":  "number",
		"string": "string",
	}
	fmain *ast.FuncDecl = nil
)

func init() {
	file, err = os.Create("./out.luau")

	if err != nil {
		panic(err)
	}

	fset = token.NewFileSet()
	node, err = parser.ParseFile(fset, FILE, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	_, err := os.ReadFile(FILE)
	if err != nil {
		panic(err)
	}

	//filestr = string(dat)
	//lines = strings.Split(filestr, "\n")
}

func main() {
	defer file.Close()

	for _, v := range node.Decls {
		if funcdecl, ok := v.(*ast.FuncDecl); ok {

			if funcdecl.Name.Name == "main" {
				println("main found")
				fmain = funcdecl
				continue
			}

			file.WriteString(fmt.Sprintf("function %s(", funcdecl.Name.Name))
			pcount := 0 // param count
			ptype := "any"

			if len(funcdecl.Type.Params.List) == 0 {
				file.WriteString(")\n")
			} else {
				for _, t := range funcdecl.Type.Params.List {

					if typename, ok := t.Type.(*ast.Ident); ok {
						ptype = typename.Name
					}
					for _, p := range t.Names {
						pcount++

						if len(funcdecl.Type.Params.List) == pcount {
							file.WriteString(fmt.Sprintf("%s : %s) \n", p.Name, types[ptype]))
						} else {
							file.WriteString(fmt.Sprintf("%s : %s, ", p.Name, types[ptype]))
						}
					}
				}
			}

			writeStatements(funcdecl)

			file.WriteString("end \n")

		}
	}
	// main
	writeStatements(fmain)
}

func writeStatements(fundecl *ast.FuncDecl) {
	for _, stmt := range fundecl.Body.List {
		if expr, ok := stmt.(*ast.ExprStmt); ok {
			if callexpr, ok := expr.X.(*ast.CallExpr); ok {
				if fun, ok := callexpr.Fun.(*ast.Ident); ok {
					file.WriteString(fmt.Sprintf("%s(", fun.Name))
				}

				if len(callexpr.Args) == 0 {
					file.WriteString(")\n")
					break
				}

				argcount := 0

				for _, arg := range callexpr.Args {
					if val, ok := arg.(*ast.BasicLit); ok {
						argcount++
						if len(callexpr.Args) == argcount {
							file.WriteString(fmt.Sprintf("%s)\n", val.Value))
						} else {
							file.WriteString(fmt.Sprintf("%s, ", val.Value))
						}

					}
					if val, ok := arg.(*ast.Ident); ok {
						argcount++
						if len(callexpr.Args) == argcount {
							file.WriteString(fmt.Sprintf("%s)\n", val.Name))
						} else {
							file.WriteString(fmt.Sprintf("%s, ", val.Name))
						}
					}
				}
			}
		}

	}
}
