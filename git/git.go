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
	PushNewBranch() error
	BuildBranchNameFromCommit(commitSha string) string
	AbortCherryPick() error
	AmendCommitWithFixup(commitSha string) error
	RebaseInteractiveAutosquash(commitSha string) error
	CurrentBranch() string
	LogFromMainFormatted() ([]string, error)
}

type gitService struct{}

func NewGitService() GitService {
	return &gitService{}
}

func (g *gitService) CurrentBranch() string {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		log.Fatal(err)
	}
	result := strings.TrimSpace(string(out))
	return result
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
	return exec.Command("git", "branch", "--no-track", branchName, "origin/master").Run()
}

func (g *gitService) Switch(branch string) error {
	return exec.Command("git", "switch", branch).Run()
}

func (g *gitService) CherryPick(commit string) error {
	return exec.Command("git", "cherry-pick", commit).Run()
}

func (g *gitService) AbortCherryPick() error {
	return exec.Command("git", "cherry-pick", "--abort").Run()
}

func (g *gitService) Push() error {
	return exec.Command("git", "push").Run()
}

func (g *gitService) PushNewBranch() error {
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

func (g *gitService) AmendCommitWithFixup(commitSha string) error {
	fixupArg := "--fixup=" + commitSha
	return exec.Command("git", "commit", "--amend", fixupArg).Run()
}

func (g *gitService) RebaseInteractiveAutosquash(commitSha string) error {
	commitArg := commitSha + "^"
	cmd := exec.Command("git", "rebase", "--interactive", "--autosquash", commitArg)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GIT_SEQUENCE_EDITOR=true")
	return cmd.Run()
}

func (g *gitService) LogFromMainFormatted() ([]string, error) {
	out, err := exec.Command("git", "log", "--oneline", "--no-decorate", "origin/master..HEAD").Output()
	if err != nil {
		return nil, err
	}

	outputString := string(out)
	logs := strings.FieldsFunc(outputString, func(r rune) bool {
		return r == '\n'
	})

	for i, log := range logs {
		parts := strings.Fields(log)
		commitMessage := strings.Join(parts[1:], " ")
		log = parts[0] + " - " + commitMessage
		logs[i] = strings.TrimSpace(log)
	}

	return logs, nil
}
