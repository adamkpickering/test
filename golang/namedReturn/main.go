package main

import (
	"fmt"
)

func testFunc() (n int) {
	n = 4
	return
}

func main() {
	fmt.Println(testFunc())
}
