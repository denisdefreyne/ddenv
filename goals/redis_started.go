package goals

import (
	"fmt"

	"denisdefreyne.com/x/ddenv/core"
)

type RedisStarted struct {
	Env map[string]string
}

func (g RedisStarted) Description() string {
	return fmt.Sprintf("Starting Redis")
}

func (g RedisStarted) HashIdentity() string {
	return fmt.Sprintf("RedisStarted %v", g)
}

func (g RedisStarted) SubGoals() []core.Goal {
	return []core.Goal{
		HomebrewPackageInstalled{PackageName: "redis"},
	}
}
