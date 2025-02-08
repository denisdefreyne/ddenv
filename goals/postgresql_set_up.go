package goals

import (
	"fmt"

	"denisdefreyne.com/x/ddenv/core"
)

func init() {
	core.RegisterGoal("postgresql", func(value interface{}) (core.Goal, error) {
		detailsMap, ok := value.(map[interface{}]interface{})
		if !ok {
			return nil, fmt.Errorf("expected details map")
		}

		// Get version
		version, ok := detailsMap["version"]
		if !ok {
			return nil, fmt.Errorf("expected version")
		}
		intVersion, ok := version.(int)
		if !ok {
			return nil, fmt.Errorf("expected integer version")
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

		g := PostgresqlSetUp{Version: intVersion, Env: env}

		return g, nil
	})
}

type PostgresqlSetUp struct {
	Version int
	Env     map[string]string
}

func (g PostgresqlSetUp) Description() string {
	return fmt.Sprintf("Setting up PostgreSQL %v", g.Version)
}

func (g PostgresqlSetUp) HashIdentity() string {
	return fmt.Sprintf("PostgresqlSetUp %v", g)
}

func (g PostgresqlSetUp) PreGoals() []core.Goal {
	return []core.Goal{
		PostgresqlStarted{Version: g.Version, Env: g.Env},
		PostgresqlShadowenvCreated{Version: g.Version, Env: g.Env},
	}
}

func (g PostgresqlSetUp) homebrewPackageName() string {
	return fmt.Sprintf("postgresql@%v", g.Version)
}
