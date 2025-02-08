package goals

import (
	"fmt"
	"os"
	"strings"

	"denisdefreyne.com/x/ddenv/core"
)

func init() {
	core.RegisterGoal("node", func(value interface{}) (core.Goal, error) {
		if nodeVersionBytes, err := os.ReadFile(".node-version"); err != nil {
			return nil, fmt.Errorf("expected .node-version to exist")
		} else {
			nodeVersionString := strings.TrimSpace(string(nodeVersionBytes))
			return NodeSetUp{Version: nodeVersionString}, nil
		}
	})
}

type NodeSetUp struct {
	Version string
}

func (g NodeSetUp) Description() string {
	return fmt.Sprintf("Setting up Node %v", g.Version)
}

func (g NodeSetUp) HashIdentity() string {
	return fmt.Sprintf("NodeSetUp %v", g)
}

func (g NodeSetUp) PreGoals() []core.Goal {
	nodeInstalledGoal := NodeInstalled{Version: g.Version}

	return []core.Goal{
		nodeInstalledGoal,
		NodeShadowenvCreated{Version: g.Version, Path: nodeInstalledGoal.path()},
	}
}
