package goals

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type HomebrewPackageInstalled struct {
	PackageName string
}

type brewInfoEntry struct {
	Installed []interface{}
}

func (g HomebrewPackageInstalled) Description() string {
	return fmt.Sprintf("Installing Homebrew package ‘%v’", g.PackageName)
}

func (g HomebrewPackageInstalled) IsAchieved() bool {
	// Get raw output
	brewInfoCmd := exec.Command("brew", "info", "--json", g.PackageName)
	brewInfoOut, err := brewInfoCmd.Output()
	if err != nil {
		return false
	}

	// Parse JSON
	var brewInfoEntries []brewInfoEntry
	if err := json.Unmarshal(brewInfoOut, &brewInfoEntries); err != nil {
		fmt.Println(err.Error())
		return false
	}

	// Check
	if len(brewInfoEntries) < 1 {
		return false
	}
	if len(brewInfoEntries[0].Installed) < 1 {
		return false
	}

	return true
}

func (g HomebrewPackageInstalled) Achieve() error {
	brewInstallCmd := exec.Command("brew", "install", g.PackageName)
	if err := brewInstallCmd.Run(); err != nil {
		return err
	}

	return nil
}
