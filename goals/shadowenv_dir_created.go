package goals

import (
	"os"
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
