package goals

import (
	"os/exec"
)

type NpmPackagesInstalled struct {
}

func (g NpmPackagesInstalled) Description() string {
	return "Installing npm packages"
}

func (g NpmPackagesInstalled) IsAchieved() bool {
	// Get raw output
	cmd := exec.Command("shadowenv", "exec", "--", "npx", "check-dependencies")

	err := cmd.Run()
	return err == nil
}

func (g NpmPackagesInstalled) Achieve() error {
	cmd := exec.Command("shadowenv", "exec", "--", "npm", "install")
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
