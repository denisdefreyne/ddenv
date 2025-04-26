package goals

import (
	"fmt"
	"os/exec"

	"denisdefreyne.com/x/ddenv/core"
)

func init() {
	core.RegisterGoal("yarn", func(value any) (core.Goal, error) {
		return YarnPackagesInstalled{}, nil
	})
}

type YarnPackagesInstalled struct {
}

func (g YarnPackagesInstalled) Description() string {
	return "Installing yarn packages"
}

func (g YarnPackagesInstalled) HashIdentity() string {
	return fmt.Sprintf("YarnPackagesInstalled %v", g)
}

func (g YarnPackagesInstalled) IsAchieved() bool {
	// Get raw output
	cmd := exec.Command("shadowenv", "exec", "--", "npx", "check-dependencies")

	err := cmd.Run()
	return err == nil
}

func (g YarnPackagesInstalled) Achieve() error {
	cmd := exec.Command("shadowenv", "exec", "--", "yarn", "install")

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v\n\n%v", err, string(stdoutStderr))
	}

	return nil
}
