package gh

import (
	"encoding/json"
	"log"
	"os"

	"github.com/cli/go-gh/v2"
)

type GitHubService interface {
	GetPullRequests() []PullRequest
	CreatePullRequest(baseBranch string, headBranch string)
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

// TODO: fill pr title and body automatically
func (g *gitHubService) CreatePullRequest(baseBranch string, headBranch string) error {
	out, _, err := gh.Exec("pr", "create", "-B", baseBranch, "-H", headBranch)

	os.Stdout.Write(out.Bytes())
	return err
}
