package pack

import (
	"errors"
	"io/fs"
	"os"
	"strings"

	"github.com/intervinn/abq/luau"
	"golang.org/x/mod/modfile"
)

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

func ResolveImport(i *luau.ImportDecl) luau.Node {
	return &luau.Raw{}
}
