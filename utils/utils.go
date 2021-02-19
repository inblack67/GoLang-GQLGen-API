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

// func main() {
// 	defer Elapsed("page")()
// 	time.Sleep(time.Second * 2)
// }