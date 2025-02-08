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
			return RubyInstalledAndAvailable{Version: rubyVersionString}, nil
		}
	})
}

type RubyInstalledAndAvailable struct {
	Version string
}

func (g RubyInstalledAndAvailable) Description() string {
	return fmt.Sprintf("Setting up Ruby %v", g.Version)
}

func (g RubyInstalledAndAvailable) HashIdentity() string {
	return fmt.Sprintf("RubyInstalledAndAvailable %v", g)
}

func (g RubyInstalledAndAvailable) IsAchieved() bool {
	// TODO: weird non-goal
	return false
}

func (g RubyInstalledAndAvailable) Achieve() error {
	// TODO: weird non-goal
	return nil
}

func (g RubyInstalledAndAvailable) PreGoals() []core.Goal {
	rubyInstalledGoal := RubyInstalled{Version: g.Version}

	return []core.Goal{
		rubyInstalledGoal,
		RubyShadowenvCreated{Version: g.Version, Path: rubyInstalledGoal.path()},
	}
}
