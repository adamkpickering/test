package main

import (
	"fmt"
	"os"
	"context"
	"errors"

	"github.com/google/go-github/v71/github"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	owner := "rancher"
	repo := "image-mirror"
	prNumber := 1016

	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		return errors.New("must define GITHUB_TOKEN")
	}
	ghClient := github.NewClient(nil).WithAuthToken(githubToken)

	reviewersRequest := github.ReviewersRequest{
		Reviewers: []string{"diogoasouza"},
		TeamReviewers: []string{"observation-backup"},
	}
	_, _, err := ghClient.PullRequests.RequestReviewers(context.Background(), owner, repo, prNumber, reviewersRequest)
	if err != nil {
		fmt.Printf("warning: failed to request reviewers for PR %d: %v\n", prNumber, err)
	}

	return nil
}
