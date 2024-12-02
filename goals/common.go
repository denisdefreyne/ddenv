package goals

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// TODO: Move all this into a `homebrew` package or so

type brewServicesListEntry struct {
	Name   string
	Status string
}

func homebrewServiceList() ([]brewServicesListEntry, error) {
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

func homebrewServiceInfoFor(packageName string) ([]brewServicesListEntry, error) {
	// Get raw output
	brewServicesListCmd := exec.Command("brew", "services", "info", "--json", packageName)
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
