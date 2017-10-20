// Package engine defines interface of an Engine. The engine would know how to truncate a table
package engine

// Engine is an interface for db interaction layer
type Engine interface {
	// Truncate a table
	Truncate(table string) error

	// Close db connection
	Close() error
}
