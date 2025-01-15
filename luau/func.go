package luau

import (
	"strings"
)

type FuncType struct {
	Params  []Field
	Results []Field
}

func (f *FuncType) Render() string {
	params := make([]string, len(f.Params))
	for i, p := range f.Params {
		params[i] = p.Render()
	}

	results := make([]string, len(f.Params))
	for i, r := range f.Results {
		params[i] = r.Render()
	}

	return "(" + strings.Join(params, ",") + ")" + " -> (" + strings.Join(results, ",") + ")"
}
