package goals

import (
	"bytes"
	"fmt"
	"os"
)

const ShadowenvDirGitIgnored_Path = ".shadowenv.d/.gitignore"
const ShadowenvDirGitIgnored_Contents = "*"

type ShadowenvDirGitIgnored struct {
	Version string
	Path    string
}

func (g ShadowenvDirGitIgnored) Description() string {
	return "Adding Shadowenv dir to .gitignore"
}

func (g ShadowenvDirGitIgnored) HashIdentity() string {
	return fmt.Sprintf("ShadowenvDirGitIgnored %v", g)
}

func (g ShadowenvDirGitIgnored) IsAchieved() bool {
	_, err := os.Lstat(ShadowenvDirGitIgnored_Path)
	if err != nil {
		return false
	}

	oldContents, err := os.ReadFile(ShadowenvDirGitIgnored_Path)
	if err != nil {
		return false
	}

	if !bytes.Equal(oldContents, []byte(ShadowenvDirGitIgnored_Contents)) {
		return false
	}

	return true
}

func (g ShadowenvDirGitIgnored) Achieve() error {
	err := os.WriteFile(ShadowenvDirGitIgnored_Path, []byte(ShadowenvDirGitIgnored_Contents), 0755)
	if err != nil {
		return err
	}

	return nil
}
