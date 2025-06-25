package main

import (
	"fmt"
)

func main() {
	slice1 := []string{"asdf", "qwer", "zxcv"}
	slice2 := make([]string, len(slice1))
	num := copy(slice2, slice1)

	fmt.Printf("elements copied: %d\n", num)
	fmt.Printf("slice1: %+v\n", slice1)
	fmt.Printf("slice2: %+v\n", slice2)
}
