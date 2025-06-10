package main

import (
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"io"
	"strings"
)

func main() {
	hasher := sha256.New()
	io.WriteString(hasher, "1234")
	output := hasher.Sum(nil)
	fmt.Println(output)

	strOutput := base32.StdEncoding.EncodeToString(output)
	lowerStrOutput := strings.ToLower(strOutput)
	fmt.Println(lowerStrOutput[:8])
}
