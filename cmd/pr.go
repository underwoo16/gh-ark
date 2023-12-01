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
	Example: `  gh-slack read <slack-permalink>
  gh-slack read -i <issue-url> <slack-permalink>`,
}

func runPrCmd() error {
	fmt.Println("Creating PR from current commit...")

	gitService := git.NewGitService()
	ghService := gh.NewGitHubService()
	return createPullRequest(gitService, ghService)
}

func createPullRequest(gitService git.GitService, ghService gh.GitHubService) error {
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
	err = ghService.CreatePullRequest()
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

	return nil
}
