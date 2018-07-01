package logging

import (
	"fmt"
	"os"
)

// Logger ...
type Logger interface {
	Println(msg string, args ...interface{})
}

// Noop logger does nothing
type Noop struct{}

func (n *Noop) Println(msg string, args ...interface{}) {}

// Stdout logger prints log to stdout
type Stdout struct{}

func (s *Stdout) Println(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, msg+"\n", args...)
}
