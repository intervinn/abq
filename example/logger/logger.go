package logger

import "github.com/intervinn/abq/transform"

var _ = transform.Mod[any]("local logger = require(\"logger\")")

type Logger struct{}

func NewLogger() *Logger {
	return transform.Mod[*Logger]("logger.new()")
}

func (l *Logger) Log(msg string) {
	_ = transform.Mod[any]("logger.msg(l, msg)")
}
