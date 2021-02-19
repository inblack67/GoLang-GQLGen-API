package utils

import (
	"fmt"
	"time"
)

// Elapsed ...
func Elapsed(what string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", what, time.Since(start))
	}
}