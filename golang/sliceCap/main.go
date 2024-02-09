package main

import (
	"fmt"
)

// Demonstrates that a slice created from an array
// has its capacity set set to the number of elements
// from its start to the end of the array.
func main() {
	var myArr [4]string
	mySlice := myArr[1:1]
	fmt.Printf("myArr: %#v\n", myArr)
	fmt.Printf("mySlice: %#v\n", mySlice)
	fmt.Printf("mySlice len: %d\n", len(mySlice))
	fmt.Printf("mySlice cap: %d\n", cap(mySlice))
}
