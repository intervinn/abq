package pack

import "github.com/intervinn/abq/luau"

func ResolveImport(i *luau.ImportDecl) luau.Node {
	return &luau.Raw{
		Content: "",
	}
}
