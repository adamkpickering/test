package main

import (
	"encoding/base64"
	"fmt"
)

func main() {
	encoded := "dGVzdCBERVIgZGF0YQ=="
	fmt.Printf("encoded: %s\n", encoded)
	// unencoded := "test DER data"

	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("decoded: %s\n", decoded)
}
