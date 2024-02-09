package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	fmt.Printf("result of time.Now(): %v\n", now)

	truncated := now.Truncate(24 * time.Hour)
	fmt.Printf("truncated: %v\n", truncated)
}
