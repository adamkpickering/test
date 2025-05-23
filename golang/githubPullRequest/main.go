package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v71/github"
	"os"
)

func main() {
	ghClient := github.NewClient(nil)
	githubOwner := "adamkpickering"
	githubRepo := "image-mirror"
	branchName := "auto-update/flannel/v0.26.7"
	pullRequests, _, err := ghClient.PullRequests.List(context.Background(), githubOwner, githubRepo, &github.PullRequestListOptions{
		Head:  githubOwner + ":" + branchName,
		State: "all",
	})
	handleError(err)
	for _, pullRequest := range pullRequests {
		pretty := map[string]string{
			"ID":    fmt.Sprintf("%d", *pullRequest.ID),
			"Title": *pullRequest.Title,
			"Head":  *pullRequest.Head.Label,
		}
		fmt.Println(pretty)
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}
