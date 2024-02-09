package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var jsonContents string = `{
	"Port": 1234,
	"User": "adam",
	"Password": "adamspassword"
}`

type ConnectionInfo struct {
	Host     string
	Port     int
	User     string
	Password string
}

func main() {
	connectionInfo := &ConnectionInfo{
		Host: "127.0.0.1",
	}

	fmt.Printf("connectionInfo before: %+v\n", *connectionInfo)

	if err := json.Unmarshal([]byte(jsonContents), connectionInfo); err != nil {
		fmt.Printf("failed to unmarshal contents: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("connectionInfo after: %+v\n", *connectionInfo)
}
