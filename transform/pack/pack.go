package pack

import (
	"fmt"
	"io/fs"
	"os"
	"slices"
	"strings"
)

var Except = []string{
	".git",
	"out",
}

func Dir(p string, dir fs.DirEntry) error {
	if slices.Contains(Except, dir.Name()) {
		return nil
	}

	entries, err := os.ReadDir(p)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.IsDir() {
			fmt.Println(e.Name())
		}

		if !strings.HasPrefix(e.Name(), ".go") {
			continue
		}
	}

	return nil
}
