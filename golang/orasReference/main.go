package main

import (
	"fmt"
	"os"

	"oras.land/oras-go/v2/registry"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	rawRefs := []string{
		"flannelcni/flannel",
		"docker.io/flannelcni/flannel",
		"docker.io/rancher/mirrored-flannelcni-flannel:v0.16.0",
		"registry.suse.com/rancher/mirrored-flannelcni-flannel:v0.16.0",
	}
	for _, rawRef := range rawRefs {
		ref, err := registry.ParseReference(rawRef)
		if err != nil {
			return fmt.Errorf("failed to parse %q: %w", rawRef, err)
		}
		fmt.Printf("%s -> {Host: %s, Registry: %s, Repository: %s, Reference: %s}\n", rawRef, ref.Host(), ref.Registry, ref.Repository, ref.Reference)
	}

	return nil
}
