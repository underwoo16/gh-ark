package git

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

type GitService interface {
	LatestCommit() string
	RevParse(thing string) string
	CreateBranch(branchName string) error
	Switch(branch string) error
	CherryPick(commit string) error
	Push() error
	BuildBranchNameFromCommit(commitSha string) string
	AbortCherryPick() error
}

type gitService struct{}

func NewGitService() *gitService {
	return &gitService{}
}

func (g *gitService) LatestCommit() string {
	return g.RevParse("HEAD")
}

func (g *gitService) RevParse(thing string) string {
	out, err := exec.Command("git", "rev-parse", thing).Output()
	if err != nil {
		log.Fatal(err)
	}
	result := strings.TrimSpace(string(out))
	return result
}

func (g *gitService) CreateBranch(branchName string) error {
	return exec.Command("git", "branch", "--no-track", branchName, "origin/main").Run()
}

func (g *gitService) Switch(branch string) error {
	return exec.Command("git", "switch", branch).Run()
}

func (g *gitService) CherryPick(commit string) error {
	cmd := exec.Command("git", "cherry-pick", commit)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func (g *gitService) AbortCherryPick() error {
	return exec.Command("git", "cherry-pick", "--abort").Run()
}

func (g *gitService) Push() error {
	return exec.Command("git", "-c", "push.default=current", "push").Run()
}

func (g *gitService) BuildBranchNameFromCommit(commitSha string) string {
	out, err := exec.Command("git", "show", "--no-patch", "--format=%f", commitSha).Output()
	if err != nil {
		log.Fatal(err)
	}
	result := strings.TrimSpace(string(out))
	return result
}

func (g *gitService) BuildBranchNameFromLastCommit() string {
	return g.BuildBranchNameFromCommit("HEAD")
}
