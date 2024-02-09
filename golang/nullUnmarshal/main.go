package main

import (
	"encoding/json"
	"fmt"
)

var inputJson string = "null"

type TestStruct struct {
	X int
	Y int
}

func main() {
	t1 := TestStruct{}
	err := json.Unmarshal([]byte(inputJson), &t1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("t1: %#v\n", t1)
}
