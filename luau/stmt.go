package luau

import (
	"strings"
)

type AssignStmt struct {
	Tok Token
	Rhs []Node
	Lhs []Node
}

func (a *AssignStmt) Render() string {
	rhs := make([]string, len(a.Rhs))
	for i, h := range a.Rhs {
		rhs[i] = h.Render()
	}

	lhs := make([]string, len(a.Lhs))
	for i, h := range a.Lhs {
		lhs[i] = h.Render()
	}

	return strings.Join(lhs, ",") + " = " + strings.Join(rhs, ",")
}
