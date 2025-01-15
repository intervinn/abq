package luau

import "testing"

func TestAssignStmt(t *testing.T) {
	s := &AssignStmt{
		Tok: ASSIGN,
		Lhs: []*Ident{
			&Ident{
				Name: "foo",
			},
		},
		Rhs: []Node{
			&BasicLit{
				Kind:  NUMBER,
				Value: "123",
			},
		},
	}

	r := s.Render()
	t.Log(r)
	if r != "foo = 123" {
		t.Fail()
	}
}
