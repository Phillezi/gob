package gob

import (
	"fmt"
	"os/exec"
	"strings"
)

func GitVersion() string {
	tag, _ := runGit("describe", "--tags", "--abbrev=0")
	commit, _ := runGit("rev-parse", "HEAD")
	tagCommit, _ := runGit("rev-list", "-n", "1", tag)

	dirty := detectDirty()

	short := commit
	if len(commit) >= 7 {
		short = commit[:7]
	}

	if !dirty {
		if tagCommit == commit {
			return tag
		}
		return fmt.Sprintf("%s-%s", tag, short)
	}

	return fmt.Sprintf("%s-dirty-%s", tag, short)
}

func runGit(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}

func detectDirty() bool {
	cmd := exec.Command("git", "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) != ""
}
