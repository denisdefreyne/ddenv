package goals

import (
	"os/exec"

	"denisdefreyne.com/x/ddenv/core"
)

func init() {
	core.RegisterGoal("bundle", func(value interface{}) (core.Goal, error) {
		return BundleInstalled{}, nil
	})
}

type BundleInstalled struct {
}

func (g BundleInstalled) Description() string {
	return "Installing bundle"
}

func (g BundleInstalled) IsAchieved() bool {
	// Get raw output
	cmd := exec.Command("shadowenv", "exec", "--", "bundle", "check")

	err := cmd.Run()
	return err == nil
}

func (g BundleInstalled) Achieve() error {
	cmd := exec.Command("shadowenv", "exec", "--", "bundle", "install")
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (g BundleInstalled) PreGoals() []core.Goal {
	return []core.Goal{
		GemInstalled{Name: "bundler"},
	}
}
