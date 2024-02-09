package main

import (
	"fmt"
	"path/filepath"
)

func main() {
	a := "1234"
	fmt.Printf("a before: %s\n", a)
	a = filepath.Join(a, "asdf")
	fmt.Printf("a after: %s\n", a)
}
