package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/underwoo16/gh-diffstack/gh"
	"github.com/underwoo16/gh-diffstack/git"
)

var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Create PR from latest commit",
	Long:  `Creates a pull request on GitHub which contains the latest commit and targets origin/main`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPrCmd()
	},
	Example: `gh-diffstack pr`,
}

func init() {
	// TODO: add flag to specify branch to target PR instead of defaulting to main
}

func runPrCmd() error {
	fmt.Println("Creating PR from current commit...")

	gitService := git.NewGitService()
	ghService := gh.NewGitHubService()
	return createPullRequest(gitService, ghService)
}

// TODO: check for existing PR before creating new one
func createPullRequest(gitService git.GitService, ghService gh.GitHubService) error {
	trunk := gitService.CurrentBranch()
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
