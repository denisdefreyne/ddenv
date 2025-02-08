package goals

import (
	"fmt"
	"os"
)

const ShadowenvDirCreated_Path = ".shadowenv.d"

type ShadowenvDirCreated struct{}

func (g ShadowenvDirCreated) Description() string {
	return "Creating Shadowenv dir"
}

func (g ShadowenvDirCreated) HashIdentity() string {
	return fmt.Sprintf("ShadowenvDirCreated %v", g)
}

func (g ShadowenvDirCreated) IsAchieved() bool {
	_, err := os.Lstat(ShadowenvDirCreated_Path)
	return err == nil
}

func (g ShadowenvDirCreated) Achieve() error {
	err := os.Mkdir(ShadowenvDirCreated_Path, 0755)
	return err
}
