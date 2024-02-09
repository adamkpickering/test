package main

import (
	"fmt"
)

func main() {
	myStr := "我\n"
	chineseChar := []rune(myStr)[0]
	newline := []rune(myStr)[1]
	fmt.Printf("%q\n", myStr)
	fmt.Printf("%q\n", chineseChar)
	fmt.Printf("%q\n", newline)
}
