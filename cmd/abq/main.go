package main

import (
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "abq",
	Short: "Transpile Go to Luau",
	Long:  "",

	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	root.AddCommand(Compile)
}

func main() {
	root.Execute()
}
