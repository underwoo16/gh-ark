package main

import (
	"fmt"
	"log"

	"github.com/underwoo16/gh-diffstack/gh"
	"github.com/underwoo16/gh-diffstack/git"
)

func main() {
	gitService := git.NewGitService()
	latestCommit := gitService.LatestCommit()
	fmt.Println(latestCommit)

	branchName := gitService.BuildBranchNameFromCommit(latestCommit)
	fmt.Println(branchName)

	err := gitService.CreateBranch(branchName)
	if err != nil {
		log.Fatal(err)
	}

	err = gitService.Switch(branchName)
	if err != nil {
		log.Fatal(err)
	}

	err = gitService.CherryPick(latestCommit)
	if err != nil {
		// TODO: abort chery-pick and switch back to main
		log.Fatal(err)
	}

	err = gitService.Push()
	if err != nil {
		log.Fatal(err)
	}

	// create pr
	ghService := gh.NewGitHubService()
	err = ghService.CreatePullRequest(branchName, "main")
	if err != nil {
		log.Fatal(err)
	}

	//switch back to main
	err = gitService.Switch("main")
	if err != nil {
		log.Fatal(err)
	}
}

// For more examples of using go-gh, see:
// https://github.com/cli/go-gh/blob/trunk/example_gh_test.go
