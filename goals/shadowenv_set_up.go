package goals

import (
	"fmt"

	"denisdefreyne.com/x/ddenv/core"
)

type ShadowenvSetUp struct{}

func (g ShadowenvSetUp) Description() string {
	return "Setting up Shadowenv"
}

func (g ShadowenvSetUp) HashIdentity() string {
	return fmt.Sprintf("ShadowenvSetUp %v", g)
}

func (g ShadowenvSetUp) SubGoals() []core.Goal {
	return []core.Goal{
		ShadowenvInitialized{},
		ShadowenvDirCreated{},
		ShadowenvDirGitIgnored{},
		ShadowenvTrusted{},
	}
}
