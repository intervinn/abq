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

func TestTables(t *testing.T) {
	text := `
	package main

	func main() {
		x := map[string]string{
			"key": "value",
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

func TestStructs(t *testing.T) {
	text := `
	package main

	type Entity struct {
		Name string
	}

	func (e *Entity) Foo() {
	}

	func NewEntity() *Entity {
		return &Entity{
			Name: "bob",
		}
	}

	func main() {
		e := NewEntity()
		fmt.Println(e.Name)
		e.Foo()
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

func TestMod(t *testing.T) {
	text := `
	package main

	func main() {
		transform.Mod("print('aaa')")
		x := 10
		transform.Mod("assert(x == 10, 'what')")
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
