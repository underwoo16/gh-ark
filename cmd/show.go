package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/spf13/cobra"
	"github.com/underwoo16/gh-diffstack/gh"
	"github.com/underwoo16/gh-diffstack/git"
	"github.com/underwoo16/gh-diffstack/utils"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the current stack of diffs",
	Long: `Show the current stack of diffs.
Includes the following:

- The branch name associated with the diff
- The commit sha
- The commit message
- The pull request url (if it exists)`,
	Run: runShowCmd,
}

// TODO: cleanup
func runShowCmd(cmd *cobra.Command, args []string) {
	gitService := git.NewGitService()

	logs, err := gitService.LogFromMainOrMaster()
	if err != nil {
		log.Fatal(err)
	}

	stacks := []*diffStack{}
	for _, log := range logs {
		parts := strings.Fields(log)
		sha := strings.TrimSpace(parts[0])

		commitMessage := strings.Join(parts[1:], " ")

		diffStack := diffStack{
			sha:           sha,
			commitMessage: commitMessage,
			branchName:    gitService.BuildBranchNameFromCommit(sha),
		}
		stacks = append(stacks, &diffStack)
	}

	ghService := gh.NewGitHubService()

	fmt.Println("Fetching pull requests")
	io := iostreams.System()
	io.StartProgressIndicator()

	// TODO: cache pull requests when created
	// only call api if no cache found
	pullRequests := ghService.GetPullRequests()

	for _, stack := range stacks {
		for _, pr := range pullRequests {
			if pr.HeadRefName == stack.branchName {
				stack.prUrl = pr.Url
			}
		}
	}

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s %s\n%s\n", head, "head *", vertical))
	for _, stack := range stacks {
		if stack.prUrl != "" {
			sb.WriteString(dot + " ")
		} else {
			sb.WriteString(circle + " ")
		}

		bn := utils.Green(stack.branchName)
		s := utils.Yellow(fmt.Sprintf(" (%s)", stack.sha))
		sb.WriteString(bn + s + "\n" + vertical)

		if stack.prUrl != "" {
			pr := utils.Blue(stack.prUrl)
			sb.WriteString(fmt.Sprintf("- %s\n%s", pr, vertical))
		}

		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("%s %s\n", trunk, "trunk"))

	io.StopProgressIndicator()
	fmt.Print(sb.String())
}

var vertical = "│"
var trunk = "┴"
var head = "┬"
var circle = "◯"
var dot = "●"

type diffStack struct {
	sha           string
	branchName    string
	commitMessage string
	prUrl         string
}
