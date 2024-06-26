package goals

import (
	"os"

	"denisdefreyne.com/x/ddenv/core"
)

const ShadowenvDirCreated_Path = ".shadowenv.d"

type ShadowenvDirCreated struct {
	Version string
	Path    string
}

func (g ShadowenvDirCreated) Description() string {
	return "Creating Shadowenv dir"
}

func (g ShadowenvDirCreated) IsAchieved() bool {
	_, err := os.Lstat(ShadowenvDirCreated_Path)
	return err == nil
}

func (g ShadowenvDirCreated) Achieve() error {
	err := os.Mkdir(ShadowenvDirCreated_Path, 0755)
	return err
}

func (g ShadowenvDirCreated) PreGoals() []core.Goal {
	return []core.Goal{
		ShadowenvInitialized{},
	}
}

func (g ShadowenvDirCreated) PostGoals() []core.Goal {
	return []core.Goal{
		ShadowenvDirGitIgnored{},
	}
}
