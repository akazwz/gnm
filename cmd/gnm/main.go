package main

import (
	"fmt"
	"os"

	"gnm/internal/manager"
)

func main() {
	if err := manager.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
