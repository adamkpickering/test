package main

import (
	"context"
	"fmt"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	fmt.Printf("before cancel: %t\n", contextIsDone(ctx))
	cancel()
	fmt.Printf("after cancel: %t\n", contextIsDone(ctx))
}

func contextIsDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
