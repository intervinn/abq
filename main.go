package main

import (
	"fmt"

	"github.com/intervinn/abq/luau"
)

func main() {
	w := luau.NewStringWriter()
	t := &luau.FuncStmt{
		Name: &luau.Ident{Name: "foo"},
		Params: []*luau.Ident{
			{Name: "x"},
		},
		Chunk: &luau.Chunk{
			List: []luau.Node{
				&luau.DeclStmt{
					Scope: luau.LOCAL,
					Values: []luau.Node{
						&luau.NumericLit{Value: "123.512"},
						&luau.StringLit{Value: "booo"},
					},
				},
			},
		},
		Scope: luau.LOCAL,
	}

	t.Render(w)
	fmt.Print(w.Content)
}
