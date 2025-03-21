package pack

import (
	"io"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/intervinn/abq/luau"
	"github.com/intervinn/abq/transform"
)

var Except = []string{
	".git",
	"out",
}

type Pack struct {
	Assembly []*luau.File
	Out      string // outdir root
}

func (pc *Pack) Assembled(name string) *luau.File {
	for _, v := range pc.Assembly {
		if v.Name == name {
			return v
		}
	}
	return nil
}

func NewPack(out string) *Pack {
	return &Pack{
		Out: out,
	}
}

func (pc *Pack) File(p string) (string, error) {
	src, err := os.ReadFile(p)
	if err != nil {
		return "", err
	}
	return string(src), nil
}

func (pc *Pack) Dir(p string) error {
	dir := path.Base(p)
	if slices.Contains(Except, dir) {
		return nil
	}

	entries, err := os.ReadDir(p)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.IsDir() {
			if slices.Contains(Except, e.Name()) {
				continue
			}
			err = pc.Dir(path.Join(p, e.Name()))
			if err != nil {
				return err
			}
		}

		if strings.HasSuffix(e.Name(), ".go") {
			asm := pc.Assembled(dir)
			if asm == nil {
				asm = luau.NewFile(dir, path.Join(pc.Out, dir))
				pc.Assembly = append(pc.Assembly, asm)
			}

			str, err := pc.File(path.Join(p, e.Name()))
			if err != nil {
				return err
			}

			src, err := transform.Source(e.Name(), str)
			if err != nil {
				return err
			}
			asm.Decls = append(asm.Decls, src...)
		}

		if strings.HasSuffix(e.Name(), ".luau") {
			asm := pc.Assembled(dir)
			if asm == nil {
				asm = luau.NewFile(dir, path.Join(pc.Out, dir))
				pc.Assembly = append(pc.Assembly, asm)
			}

			asm.Include = append(asm.Include, path.Join(p, e.Name()))
		}
	}

	return nil
}

func (p *Pack) Render() error {
	for _, a := range p.Assembly {
		root := a.Out
		err := os.MkdirAll(root, 0700)
		if err != nil {
			return err
		}

		// init.luau

		// final transformations
		// TODO: move in a more suitable place
		a.Decls = append(a.Decls, transform.Exports(a)) // export table

		// create file
		init, err := os.Create(path.Join(root, "init.luau"))
		if err != nil {
			return err
		}
		defer init.Close()

		w := luau.NewStringWriter()
		a.Render(w)

		_, err = init.WriteString(w.Content)
		if err != nil {
			return err
		}

		// included files
		for _, v := range a.Include {
			src, err := os.Open(v)
			if err != nil {
				return err
			}

			defer src.Close()
			dst, err := os.Create(path.Join(root, path.Base(v)))
			if err != nil {
				return err
			}
			defer dst.Close()

			_, err = io.Copy(dst, src)
			if err != nil {
				return err
			}

			return dst.Sync()
		}
	}
	return nil
}
