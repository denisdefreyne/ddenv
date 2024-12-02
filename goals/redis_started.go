package goals

import (
	"fmt"
	"os/exec"

	"denisdefreyne.com/x/ddenv/core"
	"denisdefreyne.com/x/ddenv/homebrew"
)

func init() {
	core.RegisterGoal("redis", func(value interface{}) (core.Goal, error) {
		detailsMap, ok := value.(map[interface{}]interface{})
		if !ok {
			return nil, fmt.Errorf("expected details map")
		}

		// Get env
		rawEnv, ok := detailsMap["env"]
		env := make(map[string]string)
		if ok {
			if typedEnv, ok := rawEnv.(map[interface{}]interface{}); !ok {
				return nil, fmt.Errorf("expected env to be a map")
			} else {
				for rawKey, rawValue := range typedEnv {
					if key, ok := rawKey.(string); ok {
						if value, ok := rawValue.(string); ok {
							env[key] = value
						} else {
							return nil, fmt.Errorf("expected env values to be strings")
						}
					} else {
						return nil, fmt.Errorf("expected env keys to be strings")
					}
				}
			}
		}

		g := RedisStarted{Env: env}

		return g, nil
	})
}

type RedisStarted struct {
	Env map[string]string
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
		RedisShadowenvCreated{Env: g.Env},
	}
}
