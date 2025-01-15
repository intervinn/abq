package luau

type Field struct {
	Type Node
	Name *Ident
}

func (f *Field) Render() string {
	return f.Name.Render() + ": " + f.Type.Render()
}
