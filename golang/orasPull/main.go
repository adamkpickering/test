package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	oras "oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/oci"
	"oras.land/oras-go/v2/registry/remote"
)

type Image struct {
	SourceImage string
	Tags        []string
}

func main() {
	errs := make([]error, 0)
	run(&errs)
	if len(errs) != 0 {
		fmt.Println(errors.Join(errs...))
		os.Exit(1)
	}
}

func run(errs *[]error) {
	imagesWithNewTags := []Image{
		{
			SourceImage: "library/alpineeeeee",
			Tags:        []string{"latest"},
		},
		{
			SourceImage: "quay.io/calico/apiserver",
			Tags:        []string{"v3.20.1", "v3.20.2", "v3.21.2"},
		},
	}

	dirPath, err := os.MkdirTemp("", "image-mirror-validation-*")
	if err != nil {
		*errs = append(*errs, fmt.Errorf("failed to create temp dir: %w", err))
		return
	}
	defer os.RemoveAll(dirPath)
	store, err := oci.New(dirPath)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("failed to instantiate oras store: %w", err))
		return
	}

	for _, newTagImage := range imagesWithNewTags {
		repo, err := parseRepository(newTagImage.SourceImage)
		if err != nil {
			wrappedErr := fmt.Errorf("failed to parse %s as repository: %w", newTagImage.SourceImage, err)
			*errs = append(*errs, wrappedErr)
			continue
		}
		for _, newTag := range newTagImage.Tags {
			descriptor, err := oras.Copy(context.Background(), repo, newTag, store, newTag, oras.DefaultCopyOptions)
			if err != nil {
				*errs = append(*errs, fmt.Errorf("failed to pull %s:%s: %w", newTagImage.SourceImage, newTag, err))
				continue
			}
			fmt.Printf("%+v\n", descriptor)
		}
	}
}

func parseRepository(repository string) (*remote.Repository, error) {
	preparedRepository := repository
	parts := strings.SplitN(repository, "/", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid format")
	}
	if !strings.Contains(parts[0], ".") {
		preparedRepository = "docker.io/" + preparedRepository
	}
	repo, err := remote.NewRepository(preparedRepository)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate repository: %w", err)
	}
	return repo, nil
}
