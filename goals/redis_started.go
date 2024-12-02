package goals

import (
	"fmt"
	"os/exec"

	"denisdefreyne.com/x/ddenv/core"
	"denisdefreyne.com/x/ddenv/homebrew"
)

func init() {
	core.RegisterGoal("redis", func(value interface{}) (core.Goal, error) {
		return RedisStarted{}, nil
	})
}

type RedisStarted struct {
}

func (g RedisStarted) Description() string {
	return fmt.Sprintf("Starting Redis")
}

func (g RedisStarted) HashIdentity() string {
	return fmt.Sprintf("RedisStarted %v", g)
}

func (g RedisStarted) IsAchieved() bool {
	brewServicesListEntries, err := homebrew.ServiceInfoFor("redis")
	if err != nil {
		return false
	}

	// Check
	if len(brewServicesListEntries) > 0 {
		return brewServicesListEntries[0].Status == "started"
	}

	return false
}

func (g RedisStarted) Achieve() error {
	cmd := exec.Command("brew", "services", "start", "redis")

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v\n\n%v", err, string(stdoutStderr))
	}

	return nil
}

func (g RedisStarted) PreGoals() []core.Goal {
	return []core.Goal{
		HomebrewPackageInstalled{PackageName: "redis"},
	}
}

func (g RedisStarted) PostGoals() []core.Goal {
	return []core.Goal{
		RedisShadowenvCreated{},
	}
}
