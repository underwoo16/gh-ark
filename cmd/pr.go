package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/underwoo16/gh-diffstack/gh"
	"github.com/underwoo16/gh-diffstack/git"
)

var prCmd = &cobra.Command{
	Use:   "npr",
	Short: "Create PR from latest commit",
	Long:  `Creates a pull request on GitHub which contains the latest commit and targets origin/main`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPrCmd()
	},
	Example: `gh-diffstack npr
gh-diffstack npr -l`,
}

// var branch string
var list bool

func init() {
	// prCmd.Flags().StringVarP(&branch, "branch", "b", "main", "Branch to target PR")
	prCmd.Flags().BoolVarP(&list, "list", "l", false, "Select commit from list")
}

func runPrCmd() error {
	fmt.Println("Creating PR from current commit...")

	gitService := git.NewGitService()
	ghService := gh.NewGitHubService()

	if list {
		return createPullRequestList(gitService, ghService)
	}

	return createPullRequestLatest(gitService, ghService)
}

// TODO: error handling
func createPullRequestList(gitService git.GitService, ghService gh.GitHubService) error {
	commits, _ := gitService.LogFromMainFormatted()
	choice, _ := ghService.Prompt("Select commit to create PR from", commits[0], commits)

	commit := strings.Fields(commits[choice])[0]
	return createPullRequest(commit, gitService, ghService)
}

func createPullRequestLatest(gitService git.GitService, ghService gh.GitHubService) error {
	latestCommit := gitService.LatestCommit()
	return createPullRequest(latestCommit, gitService, ghService)
}

// TODO: check for existing PR before creating new one
func createPullRequest(commit string, gitService git.GitService, ghService gh.GitHubService) error {
	trunk := gitService.CurrentBranch()

	branchName := gitService.BuildBranchNameFromCommit(commit)
	fmt.Println(branchName)

	err := gitService.CreateBranch(branchName)
	if err != nil {
		log.Fatal(err)
	}

	err = gitService.Switch(branchName)
	if err != nil {
		log.Fatal(err)
	}

	err = gitService.CherryPick(commit)

	if err != nil {
		fmt.Println("Cherry-pick failed, aborting...")
		gitService.AbortCherryPick()
		gitService.Switch(trunk)
		log.Fatal(err)
	}

	fmt.Println("Cherry-pick succeeded, pushing to remote...")

	err = gitService.PushNewBranch()
	if err != nil {
		fmt.Println("Push failed...")
		log.Fatal(err)
	}

	fmt.Println("Push succeeded, creating PR...")

	err = ghService.CreatePullRequest()
	if err != nil {
		fmt.Println("PR creation failed...")
		log.Fatal(err)
	}

	fmt.Printf("PR created, switching back to %s...\n", trunk)

	err = gitService.Switch(trunk)
	if err != nil {
		fmt.Printf("Switching back to %s failed...\n", trunk)
		log.Fatal(err)
	}

	return nil
}
