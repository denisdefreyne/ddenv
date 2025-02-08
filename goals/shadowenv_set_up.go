package goals

import (
	"fmt"
	"os"

	"denisdefreyne.com/x/ddenv/core"
)

const ShadowenvSetUp_Path = ".shadowenv.d"

type ShadowenvSetUp struct{}

func (g ShadowenvSetUp) Description() string {
	return "Setting up Shadowenv"
}

func (g ShadowenvSetUp) HashIdentity() string {
	return fmt.Sprintf("ShadowenvSetUp %v", g)
}

func (g ShadowenvSetUp) IsAchieved() bool {
	_, err := os.Lstat(ShadowenvSetUp_Path)
	return err == nil
}

func (g ShadowenvSetUp) Achieve() error {
	err := os.Mkdir(ShadowenvSetUp_Path, 0755)
	return err
}

func (g ShadowenvSetUp) PreGoals() []core.Goal {
	return []core.Goal{
		ShadowenvInitialized{},
		ShadowenvDirCreated{},
		ShadowenvDirGitIgnored{},
		ShadowenvTrusted{},
	}
}
