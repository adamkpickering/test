package main

import (
	"fmt"
	"os"
)

func main() {
	var err error
	if _, err := os.Stat("nonexistentpath1"); err != nil {
		err = fmt.Errorf("failed to print nonexistentpath1: %w", err)
	}
	fmt.Printf("err after nonexistentpath1: %q\n", err) // prints nil err
	if _, err = os.Stat("nonexistentpath2"); err != nil {
		err = fmt.Errorf("failed to print nonexistentpath2: %w", err)
	}
	fmt.Printf("err after nonexistentpath2: %q\n", err) // prints wrapped err from os.Stat call
}
