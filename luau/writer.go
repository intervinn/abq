package luau

import "strings"

type StringWriter struct {
	Content string
	indent  int
}

func NewStringWriter() *StringWriter {
	return &StringWriter{
		Content: "",
		indent:  0,
	}
}

func (sw *StringWriter) Indent() int {
	return sw.indent
}

func (sw *StringWriter) Pre(s string) error {
	return sw.Write(strings.Repeat("\t", sw.indent) + s)
}

func (sw *StringWriter) Write(s string) error {
	sw.Content += s
	return nil
}

func (sw *StringWriter) End(s string) error {
	return sw.Pre(s + "\n")
}

func (sw *StringWriter) IncIndent() {
	sw.indent++
}

func (sw *StringWriter) DecIndent() {
	sw.indent--
}
