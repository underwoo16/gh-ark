package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/underwoo16/gh-diffstack/gh"
	"github.com/underwoo16/gh-diffstack/git"
	"github.com/underwoo16/gh-diffstack/utils"
)

var prCmd = &cobra.Command{
	Use:   "npr",
	Short: "Create PR from latest commit",
	Long:  `Creates a pull request on GitHub which contains the latest commit and targets origin/master`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPrCmd()
	},
	Example: `gh-diffstack npr
gh-diffstack npr -l`,
}

// var branch string
var newPrList bool

func init() {
	// prCmd.Flags().StringVarP(&branch, "branch", "b", "main", "Branch to target PR")
	prCmd.Flags().BoolVarP(&newPrList, "list", "l", false, "Select commit from list")
}

func runPrCmd() error {
	gitService := git.NewGitService()
	ghService := gh.NewGitHubService()

	if newPrList {
		return createPullRequestList(gitService, ghService)
	}

	return createPullRequestLatest(gitService, ghService)
}

func createPullRequestList(gitService git.GitService, ghService gh.GitHubService) error {
	commits, _ := gitService.LogFromMainOrMaster()
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
	fmt.Printf("Creating PR from commit %s...\n", utils.Yellow(commit))

	trunk := gitService.CurrentBranch()

	branchName := gitService.BuildBranchNameFromCommit(commit)

	if !gitService.LocalBranchExists(branchName) {
		err := gitService.CreateBranch(branchName)
		if err != nil {
			log.Fatal(err)
		}
	}

	// TODO: stash changes or use worktree?
	err := gitService.Switch(branchName)
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

	err = gitService.PushNewBranch()
	if err != nil {
		fmt.Println("Push failed...")
		log.Fatal(err)
	}

	// TODO: handle when user cancels - error code 2?
	err = ghService.CreatePullRequest()
	if err != nil {
		fmt.Println("PR creation failed...")
		log.Fatal(err)
	}

	err = gitService.Switch(trunk)
	if err != nil {
		fmt.Printf("Switching back to %s failed...\n", utils.Green(trunk))
		log.Fatal(err)
	}

	return nil
}
