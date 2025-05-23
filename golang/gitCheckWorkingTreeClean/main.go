package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	workingTreeClean, err := IsWorkingTreeClean()
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Println("working tree clean:", workingTreeClean)
}

func IsWorkingTreeClean() (bool, error) {
	cmd := exec.Command("git", "diff", "--quiet")
	err := cmd.Run()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return false, nil
		} else {
			return false, fmt.Errorf("failed to run git diff: %w", err)
		}
	}
	return true, nil
}
