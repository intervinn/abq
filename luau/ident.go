package luau

type Ident struct {
	Name string
}

func (i *Ident) Render() string {
	return i.Name
}
