package goals

import (
	"fmt"
	"os"
	"strings"

	"denisdefreyne.com/x/ddenv/core"
)

func init() {
	core.RegisterGoal("ruby", func(value interface{}) (core.Goal, error) {
		if rubyVersionBytes, err := os.ReadFile(".ruby-version"); err != nil {
			return nil, fmt.Errorf("expected .ruby-version to exist")
		} else {
			rubyVersionString := strings.TrimSpace(string(rubyVersionBytes))
			return RubySetUp{Version: rubyVersionString}, nil
		}
	})
}

type RubySetUp struct {
	Version string
}

func (g RubySetUp) Description() string {
	return fmt.Sprintf("Setting up Ruby %v", g.Version)
}

func (g RubySetUp) HashIdentity() string {
	return fmt.Sprintf("RubySetUp %v", g)
}

func (g RubySetUp) PreGoals() []core.Goal {
	rubyInstalledGoal := RubyInstalled{Version: g.Version}

	return []core.Goal{
		rubyInstalledGoal,
		RubyShadowenvCreated{Version: g.Version, Path: rubyInstalledGoal.path()},
	}
}
