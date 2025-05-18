package multiple

import "fmt"

type Logger struct{}

func (l *Logger) Log(msg string) {
	fmt.Println(msg)
}
