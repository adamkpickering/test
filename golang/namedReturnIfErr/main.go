package main

import (
	"fmt"
	"os"
)

func testFunc() (err error) {
	_, err = os.Getwd()
	if _, err = os.ReadFile("doesnotexist"); err != nil {
		return
	}
	return
}

func main() {
	fmt.Println(testFunc())
}
