package logging

// Logger ...
type Logger interface {
	Println(msg string, args ...interface{})
}

// Noop logger does nothing
type Noop struct{}

func (n *Noop) Println(msg string, args ...interface{}) {}
