package main

import (
	"io/fs"
	"os"
	"path"

	"github.com/intervinn/abq/transform/pack"
	"github.com/spf13/cobra"
)

var Compile = &cobra.Command{
	Use:  "compile",
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		arg := args[0]
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		root := path.Join(cwd, arg)
		f := os.DirFS(root)
		fs.WalkDir(f, ".", func(p string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				pack.Dir(p, d)
			}

			return nil
		})

		return nil
	},
}
