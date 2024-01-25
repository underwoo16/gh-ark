package git

import (
	"fmt"
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
	LogFrom(branch string) ([]string, error)
	LogFromMainOrMaster() ([]string, error)
	LocalBranchExists(branch string) bool
	RemoteBranchExists(branch string) bool
	Pull() error
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

func (g *gitService) LocalBranchExists(branch string) bool {
	_, err := exec.Command("git", "rev-parse", "--verify", branch).Output()
	return err == nil
}

func (g *gitService) RemoteBranchExists(branch string) bool {
	_, err := exec.Command("git", "ls-remote", "--exit-code", "origin", branch).Output()
	return err == nil
}

func (g *gitService) GetTrunk() string {
	trunk := "master"
	if !g.LocalBranchExists(trunk) {
		trunk = "main"
	}
	return trunk
}

func (g *gitService) CreateBranch(branchName string) error {
	trunkArg := fmt.Sprintf("origin/%s", g.GetTrunk())
	return exec.Command("git", "branch", "--no-track", branchName, trunkArg).Run()
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

func (g *gitService) LogFromMainOrMaster() ([]string, error) {
	return g.LogFrom(g.GetTrunk())
}

func (g *gitService) LogFrom(branch string) ([]string, error) {
	branchArg := fmt.Sprintf("origin/%s..HEAD", branch)
	out, err := exec.Command("git", "log", "--oneline", "--no-decorate", branchArg).Output()
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

func (g *gitService) Pull() error {
	return exec.Command("git", "pull").Run()
}
