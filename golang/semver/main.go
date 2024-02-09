package main

import (
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"os"
)

func main() {
	// demonstrate the difference between strict and non-strict semver parsing
	versions := []string{
		"v1.2.3",
		"1.2.3",
	}
	for _, version := range versions {
		parsedVersion, err := semver.NewVersion(version)
		if err != nil {
			fmt.Printf("failed to parse %q non-strict: %s\n", version, err)
		} else {
			fmt.Printf("parsedVersion: %v\n", parsedVersion)
		}

		parsedStrictVersion, err := semver.StrictNewVersion(version)
		if err != nil {
			fmt.Printf("failed to parse %q strict: %s\n", version, err)
		} else {
			fmt.Printf("parsedStrictVersion: %v\n", parsedStrictVersion)
		}
	}

	// demonstrate that constraints with v in them work the same as regular constraints
	constraints := []string{
		">v1.0.0",
		">1.0.0",
	}
	version, _ := semver.NewVersion("v1.2.3")
	for _, constraint := range constraints {
		parsedConstraint, err := semver.NewConstraint(constraint)
		if err != nil {
			fmt.Printf("failed to parse constraint %q non-strict: %s\n", constraint, err)
		}
		fmt.Printf("parsedConstraint: %v\n", parsedConstraint)

		satisfied := parsedConstraint.Check(version)
		fmt.Printf("constraint %q satisfied by %q: %t\n", parsedConstraint, version, satisfied)
	}

	// see what parsed semvers get serialized as
	versions = []string{
		"v1.2.3",
		"1.2.3",
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	for _, version := range versions {
		parsedVersion, err := semver.NewVersion(version)
		if err != nil {
			fmt.Printf("failed to parse version %q: %s\n", version, err)
		} else {
			fmt.Printf("parsedVersion as JSON: ")
			encoder.Encode(parsedVersion)
		}
	}
}
