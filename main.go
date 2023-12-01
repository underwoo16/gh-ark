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
		fmt.Println("Cherry-pick failed, need to --quit...")
		log.Fatal(err)
	}

	fmt.Println("Cherry-pick succeeded, pushing to remote...")
	err = gitService.Push()
	if err != nil {
		fmt.Println("Push failed...")
		log.Fatal(err)
	}

	fmt.Println("Push succeeded, creating PR...")
	// create pr
	ghService := gh.NewGitHubService()
	err = ghService.CreatePullRequest(branchName, "main")
	if err != nil {
		fmt.Println("PR creation failed...")
		log.Fatal(err)
	}

	fmt.Println("PR created, switching back to main...")
	//switch back to main
	err = gitService.Switch("main")
	if err != nil {
		fmt.Println("Switching back to main failed...")
		log.Fatal(err)
	}
	fmt.Println("DONE!")
}

// For more examples of using go-gh, see:
// https://github.com/cli/go-gh/blob/trunk/example_gh_test.go
