package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type testType struct {
	MyProperty string
}

func main() {
	if err := doThing(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func doThing() error {
	data, err := os.ReadFile("cert.der.base64")
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	a := testType{
		MyProperty: string(data),
	}
	fmt.Printf("a as go struct:\n%s\n\n", a)
	marshaledData, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf("failed to marshal a: %w", err)
	}
	fmt.Printf("a as JSON:\n%s\n\n", marshaledData)
	b := testType{}
	if err := json.Unmarshal(marshaledData, &b); err != nil {
		return fmt.Errorf("failed to unmarshal marshaledData: %w", err)
	}
	fmt.Printf("b after unmarshaling from marshaled a:\n%s\n\n", b)
	return nil
}
