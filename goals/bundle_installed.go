package goals

import (
	"fmt"
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

func (g BundleInstalled) HashIdentity() string {
	return fmt.Sprintf("BundleInstalled %v", g)
}

func (g BundleInstalled) IsAchieved() bool {
	// Get raw output
	cmd := exec.Command("shadowenv", "exec", "--", "bundle", "check")

	err := cmd.Run()
	return err == nil
}

func (g BundleInstalled) Achieve() error {
	cmd := exec.Command("shadowenv", "exec", "--", "bundle", "install")

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v\n\n%v", err, string(stdoutStderr))
	}

	return nil
}

func (g BundleInstalled) SubGoals() []core.Goal {
	return []core.Goal{
		GemInstalled{Name: "bundler"},
	}
}
