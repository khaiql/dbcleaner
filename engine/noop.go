package engine

import "fmt"

type NoOp struct{}

func (n *NoOp) Close() error {
	fmt.Println("Closed NoOp engine")
	return nil
}

func (n *NoOp) Truncate(table string) error {
	fmt.Printf("Truncate table %s with noop engine\n", table)

	return nil
}
