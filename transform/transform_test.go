package transform

import (
	"fmt"
	"testing"

	"github.com/intervinn/abq/luau"
)

func TestBasic(t *testing.T) {
	text := `
	package main

	func main() {
		fmt.Println("hello, world")
	}
	`

	src, err := Source("main.go", text)
	if err != nil {
		t.Error(err)
	}

	for _, s := range src {
		w := luau.NewStringWriter()
		s.Render(w)
		fmt.Println(w.Content)
	}
}

func TestIf(t *testing.T) {
	text := `
	package main

	func main() {
		if true {
			fmt.Println("true")
		} else if false {
			fmt.Println("false")
		} else {
			fmt.Println("else")
		}
	}
	`

	src, err := Source("main.go", text)
	if err != nil {
		t.Error(err)
	}

	for _, s := range src {
		w := luau.NewStringWriter()
		s.Render(w)
		fmt.Println(w.Content)
	}
}
