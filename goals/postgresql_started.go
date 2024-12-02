package goals

import (
	"fmt"
	"os/exec"
	"strings"

	"denisdefreyne.com/x/ddenv/core"
	"denisdefreyne.com/x/ddenv/homebrew"
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

		g := PostgresqlStarted{Version: intVersion, Env: env}

		return g, nil
	})
}

type PostgresqlStarted struct {
	Version int
	Env     map[string]string
}

func (g PostgresqlStarted) Description() string {
	return fmt.Sprintf("Starting PostgreSQL %v", g.Version)
}

func (g PostgresqlStarted) HashIdentity() string {
	return fmt.Sprintf("PostgresqlStarted %v", g)
}

func (g PostgresqlStarted) IsAchieved() bool {
	brewServicesListEntries, err := homebrew.ServiceInfoFor(g.homebrewPackageName())
	if err != nil {
		return false
	}

	// Check
	if len(brewServicesListEntries) > 0 {
		return brewServicesListEntries[0].Status == "started"
	}

	return false
}

func (g PostgresqlStarted) Achieve() error {
	// Find existing PostgreSQL servers of other versions
	brewServicesListEntries, err := homebrew.ServiceList()
	if err == nil {
		for _, entry := range brewServicesListEntries {
			if strings.HasPrefix(entry.Name, "postgresql@") && entry.Name != g.homebrewPackageName() {
				return fmt.Errorf("A PostgreSQL server with a different version (%v) is already running, and so the requested PostgreSQL server version (%v) cannot be started. ddenv cannot safely resolve this problem.\n", entry.Name, g.homebrewPackageName())
			}
		}
	}

	cmd := exec.Command("brew", "services", "start", g.homebrewPackageName())

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v\n\n%v", err, string(stdoutStderr))
	}

	return nil
}

func (g PostgresqlStarted) PreGoals() []core.Goal {
	return []core.Goal{
		HomebrewPackageInstalled{PackageName: g.homebrewPackageName()},
	}
}

func (g PostgresqlStarted) PostGoals() []core.Goal {
	return []core.Goal{
		PostgresqlShadowenvCreated{Version: g.Version, Env: g.Env},
	}
}

func (g PostgresqlStarted) homebrewPackageName() string {
	return fmt.Sprintf("postgresql@%v", g.Version)
}
