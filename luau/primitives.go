package luau

type (
	StringType struct{}

	NumberType struct{}

	BoolType struct{}

	NilType struct{}

	TableType struct{}

	ThreadType struct{}

	UserDataType struct{}

	VectorType struct{}

	BufferType struct{}
)

func (s *StringType) Render() string {
	return "string"
}

func (n *NumberType) Render() string {
	return "number"
}

func (b *BoolType) Render() string {
	return "boolean"
}

func (u *UserDataType) Render() string {
	return "userdata"
}

func (v *VectorType) Render() string {
	return "vector"
}

func (b *BufferType) Render() string {
	return "buffer"
}

func (t *ThreadType) Render() string {
	return "thread"
}

func (t *TableType) Render() string {
	return "table"
}

func (n *NilType) Render() string {
	return "nil"
}
