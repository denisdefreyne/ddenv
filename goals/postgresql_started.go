package goals

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"denisdefreyne.com/x/ddenv/core"
)

func init() {
	core.RegisterGoal("postgresql", func(value interface{}) (core.Goal, error) {
		if postgresqlVersion, ok := value.(int); ok {
			return PostgresqlStarted{Version: postgresqlVersion}, nil
		} else {
			return nil, fmt.Errorf("expected string version")
		}
	})
}

type brewServicesListEntry struct {
	Name   string
	Status string
}

type PostgresqlStarted struct {
	Version int
}

func (g PostgresqlStarted) Description() string {
	return fmt.Sprintf("Starting PostgreSQL %v", g.Version)
}

func (g PostgresqlStarted) IsAchieved() bool {
	brewServicesListEntries, err := g.homebrewServiceInfoForThisPostgresqlVersion()
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
	brewServicesListEntries, err := g.homebrewServiceList()
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
		PostgresqlShadowenvCreated{Version: g.Version},
	}
}

func (g PostgresqlStarted) homebrewPackageName() string {
	return fmt.Sprintf("postgresql@%v", g.Version)
}

func (g PostgresqlStarted) homebrewServiceInfoForThisPostgresqlVersion() ([]brewServicesListEntry, error) {
	// Get raw output
	brewServicesListCmd := exec.Command("brew", "services", "info", "--json", g.homebrewPackageName())
	brewServicesListData, err := brewServicesListCmd.Output()
	if err != nil {
		return nil, err
	}

	// Parse JSON
	var brewServicesListEntries []brewServicesListEntry
	if err := json.Unmarshal(brewServicesListData, &brewServicesListEntries); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return brewServicesListEntries, nil
}

func (g PostgresqlStarted) homebrewServiceList() ([]brewServicesListEntry, error) {
	// Get raw output
	brewServicesListCmd := exec.Command("brew", "services", "list", "--json")
	brewServicesListData, err := brewServicesListCmd.Output()
	if err != nil {
		return nil, err
	}

	// Parse JSON
	var brewServicesListEntries []brewServicesListEntry
	if err := json.Unmarshal(brewServicesListData, &brewServicesListEntries); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return brewServicesListEntries, nil
}
