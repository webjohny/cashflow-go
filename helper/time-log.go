package helper

import (
	"fmt"
	"time"
)

func TimeLog(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}
