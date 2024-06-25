package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	testUrl := "file://test.txt"
	response, err := http.Get(testUrl)
	if err != nil {
		fmt.Printf("got error: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(response.Body)
}
