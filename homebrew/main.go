package homebrew

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type ServicesListEntry struct {
	Name   string
	Status string
}

func ServiceList() ([]ServicesListEntry, error) {
	// Get raw output
	brewServicesListCmd := exec.Command("brew", "services", "list", "--json")
	brewServicesListData, err := brewServicesListCmd.Output()
	if err != nil {
		return nil, err
	}

	// Parse JSON
	var brewServicesListEntries []ServicesListEntry
	if err := json.Unmarshal(brewServicesListData, &brewServicesListEntries); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return brewServicesListEntries, nil
}

func ServiceInfoFor(packageName string) ([]ServicesListEntry, error) {
	// Get raw output
	brewServicesListCmd := exec.Command("brew", "services", "info", "--json", packageName)
	brewServicesListData, err := brewServicesListCmd.Output()
	if err != nil {
		return nil, err
	}

	// Parse JSON
	var brewServicesListEntries []ServicesListEntry
	if err := json.Unmarshal(brewServicesListData, &brewServicesListEntries); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return brewServicesListEntries, nil
}
