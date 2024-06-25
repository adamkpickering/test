package main

import (
	"fmt"
)

type MyStruct struct {
	Annotations map[string]string
}

func main() {
	a := MyStruct{}
	fmt.Printf("a.Annotations: %#v\n", a.Annotations)

	b := MyStruct{}
	b.Annotations = make(map[string]string)
	fmt.Printf("b.Annotations: %#v\n", b.Annotations)
}
