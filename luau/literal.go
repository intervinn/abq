package luau

type BasicLit struct {
	Kind  Token
	Value string
}

func (b *BasicLit) Render() string {
	return b.Value
}

type FuncLit struct {
}
