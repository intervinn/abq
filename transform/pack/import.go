package pack

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"golang.org/x/mod/modfile"
)

/*
	/shared/things.go -> /shared/go_include/github.com/.../shared/things.go
	github.com/gofiber/fiber -> /shared/go_include/github.com/fiber
*/

func pkgDir() string {
	return "~/go/pkg/mod"
}

func dirExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	if errors.Is(err, fs.ErrNotExist) {
		return false
	}
	return false
}

/*
func isVersioned() {

}

func externalExists(path string, version string) bool {
	root := pkgDir()
	pieces := strings.Split(path, "/")

	for _, p := range pieces {

		if !dirExists(root) {
			return false
		}
		root += "/" + p
	}
	return false
}

func resolveExternal(i *luau.ImportDecl, field modfile.Require) *Pack {
	asm := []*luau.File{}

	return &Pack{
		Package:  "TODO: resolve package names from source",
		Assembly: asm,
		Out:      "TODO: add outdir/shared/go_include",
	}
}
*/

func ReadModfile(src string) (*modfile.File, error) {
	modsrc, err := os.ReadFile(src)
	if err != nil {
		return nil, err
	}

	mod, err := modfile.Parse("go.mod", modsrc, nil)
	return mod, err
}

func ResolveImports(mod *modfile.File) ([]*Pack, error) {
	res := []*Pack{}
	for _, r := range mod.Require {
		fmt.Println(r.Syntax.Comments.Suffix)
	}

	return res, nil
}
