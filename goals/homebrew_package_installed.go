package goals

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type HomebrewPackageInstalled struct {
	PackageName string
}

type brewInfoFormula struct {
	Installed []interface{}
}

type brewInfoCask struct {
	Installed string
}

type brewInfo struct {
	Formulae []brewInfoFormula
	Casks []brewInfoCask
}

type brewInfoEntry struct {
	Installed []interface{}
}

func (g HomebrewPackageInstalled) Description() string {
	return fmt.Sprintf("Installing Homebrew package ‘%v’", g.PackageName)
}

func (g HomebrewPackageInstalled) IsAchieved() bool {
	// Get raw output
	brewInfoCmd := exec.Command("brew", "info", "--json=v2", g.PackageName)
	brewInfoOut, err := brewInfoCmd.Output()
	if err != nil {
		return false
	}

	// Parse JSON
	var info brewInfo
	if err := json.Unmarshal(brewInfoOut, &info); err != nil {
		fmt.Println(err.Error())
		return false
	}

	// Check
	if len(info.Formulae) > 0 {
		if len(info.Formulae[0].Installed) > 0 {
			return true
		} else {
			return false
		}
	} else if len(info.Casks) > 0 {
		if info.Casks[0].Installed != "" {
			return true
		} else {
			return false
		}
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
