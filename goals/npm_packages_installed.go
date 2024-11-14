package goals

import (
	"fmt"
	"os/exec"

	"denisdefreyne.com/x/ddenv/core"
)

func init() {
	core.RegisterGoal("npm", func(value interface{}) (core.Goal, error) {
		return NpmPackagesInstalled{}, nil
	})
}

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

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v\n\n%v", err, string(stdoutStderr))
	}

	return nil
}
