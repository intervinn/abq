package main

import (
	"os"
	"path"

	"github.com/intervinn/abq/transform/pack"
	"github.com/spf13/cobra"
)

var Build = &cobra.Command{
	Use:  "build",
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		arg := args[0]
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		root := path.Join(cwd, arg)
		out := path.Join(cwd, "out")

		os.Mkdir(out, 0700)

		p := pack.NewPack(out)
		err = p.Dir(root)
		if err != nil {
			return err
		}

		err = p.Render()
		if err != nil {
			return err
		}

		return nil
	},
}
