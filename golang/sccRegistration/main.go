package main

import (
	"fmt"

	"github.com/SUSE/connect-ng/pkg/registration"
)

func main() {
	metadata := registration.Metadata
	fmt.Println(metadata)
}
