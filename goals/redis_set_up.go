package goals

import (
	"fmt"

	"denisdefreyne.com/x/ddenv/core"
)

func init() {
	core.RegisterGoal("redis", func(value any) (core.Goal, error) {
		detailsMap, ok := value.(map[any]any)
		if !ok {
			return nil, fmt.Errorf("expected details map")
		}

		// Get env
		rawEnv, ok := detailsMap["env"]
		env := make(map[string]string)
		if ok {
			if typedEnv, ok := rawEnv.(map[any]any); !ok {
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

		g := RedisSetUp{Env: env}

		return g, nil
	})
}

type RedisSetUp struct {
	Env map[string]string
}

func (g RedisSetUp) Description() string {
	return fmt.Sprintf("Setting up Redis")
}

func (g RedisSetUp) HashIdentity() string {
	return fmt.Sprintf("RedisSetUp %v", g)
}

func (g RedisSetUp) SubGoals() []core.Goal {
	return []core.Goal{
		RedisStarted{Env: g.Env},
		RedisShadowenvCreated{Env: g.Env},
	}
}
