package engine

// Engine is an interface for db interaction layer
type Engine interface {
	// Truncate a table
	Truncate(table string) error

	// Close db connection
	Close() error
}
