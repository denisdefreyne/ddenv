package goals

import (
	"fmt"
	"os/exec"

	"denisdefreyne.com/x/ddenv/core"
)

func init() {
	core.RegisterGoal("pnpm", func(value any) (core.Goal, error) {
		return PnpmPackagesInstalled{}, nil
	})
}

type PnpmPackagesInstalled struct {
}

func (g PnpmPackagesInstalled) Description() string {
	return "Installing pnpm packages"
}

func (g PnpmPackagesInstalled) HashIdentity() string {
	return fmt.Sprintf("PnpmPackagesInstalled %v", g)
}

func (g PnpmPackagesInstalled) IsAchieved() bool {
	// Get raw output
	cmd := exec.Command("shadowenv", "exec", "--", "npx", "check-dependencies")

	err := cmd.Run()
	return err == nil
}

func (g PnpmPackagesInstalled) Achieve() error {
	cmd := exec.Command("shadowenv", "exec", "--", "pnpm", "install")

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v\n\n%v", err, string(stdoutStderr))
	}

	return nil
}
