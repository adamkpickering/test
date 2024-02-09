package main

import (
	"fmt"
)

func main() {
	toTruncate := "this string 🔥✨🎉 is 40 bytes\nlong"
	fmt.Printf("non-truncated: %q\n", toTruncate)
	fmt.Printf("[]byte length: %d\n", len(toTruncate))
	fmt.Printf("[]rune length: %d\n", len([]rune(toTruncate)))
	truncatedAsString := toTruncate[0:14]
	fmt.Printf("truncatedAstString: %s\n", truncatedAsString)
	truncated := string([]rune(toTruncate)[0:14])
	fmt.Printf("truncated: %q\n", truncated)
	pastLength := toTruncate[0:50]
	fmt.Printf("pastLength: %q\n", pastLength)
}
