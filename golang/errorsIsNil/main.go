package main

import (
	"errors"
	"fmt"
)

var testErr = errors.New("this is a test error")

func main() {
	fmt.Println(errors.Is(nil, testErr))
}
