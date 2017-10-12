package engine

import "fmt"

type NoOp struct{}

func (n *NoOp) Truncate(table string) error {
	fmt.Printf("Truncate table %s with noop engine\n", table)

	return nil
}
