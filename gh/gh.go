package gh

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/cli/go-gh/v2"
	"github.com/cli/go-gh/v2/pkg/prompter"
)

type GitHubService interface {
	GetPullRequests() []PullRequest
	CreatePullRequest() error
	Prompt(promptText string, defaultVal string, values []string) (int, error)
}

type PullRequest struct {
	BaseRefName string `json:"baseRefName"`
	HeadRefName string `json:"headRefName"`
	Id          string `json:"id"`
	Number      int    `json:"number"`
	Url         string `json:"url"`
}

type gitHubService struct{}

// TODO: Add func to test if gh is installed

func NewGitHubService() *gitHubService {
	return &gitHubService{}
}

func (g *gitHubService) Prompt(promptText string, defaultVal string, values []string) (int, error) {
	return prompter.New(os.Stdin, os.Stdout, os.Stderr).Select(promptText, defaultVal, values)
}

func (g *gitHubService) GetPullRequests() []PullRequest {
	out, _, err := gh.Exec("pr", "list", "--author", "@me", "--json", "number,baseRefName,headRefName,url")
	if err != nil {
		log.Fatal(err)
	}

	pullRequests := []PullRequest{}
	err = json.Unmarshal(out.Bytes(), &pullRequests)
	if err != nil {
		log.Fatal(err)
	}

	return pullRequests
}

func (g *gitHubService) CreatePullRequest() error {
	return gh.ExecInteractive(context.Background(), "pr", "create")
}
