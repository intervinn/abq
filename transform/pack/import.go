package pack

import (
	"log"
	"os"
	"path"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

/*
	/shared/things.go -> /shared/go_include/github.com/.../shared/things.go
	github.com/gofiber/fiber -> /shared/go_include/github.com/fiber
*/

func ReadModfile(src string) (*modfile.File, error) {
	modsrc, err := os.ReadFile(src)
	if err != nil {
		return nil, err
	}

	mod, err := modfile.Parse("go.mod", modsrc, nil)
	return mod, err
}

func ModPath(name, version string) (string, error) {
	cache, ok := os.LookupEnv("GOMODCACHE")
	if !ok {
		cache = path.Join(os.Getenv("GOPATH"), "pkg", "mod")
	}

	escPath, err := module.EscapePath(name)
	if err != nil {
		return "", err
	}

	escVer, err := module.EscapeVersion(version)
	if err != nil {
		return "", err
	}

	return path.Join(cache, escPath+"@"+escVer), nil
}

func ResolveImports(mod *modfile.File, out string) ([]*Pack, error) {
	res := []*Pack{}
	for _, r := range mod.Require {
		modPath, err := ModPath(r.Mod.Path, r.Mod.Version)
		if err != nil {
			return nil, err
		}

		log.Printf("building module %s...\n", r.Mod.String())
		p := NewPack(path.Join(out, r.Mod.Path))

		err = p.Dir(modPath)
		if err != nil {
			return nil, err
		}

		res = append(res, p)
	}

	return res, nil
}
