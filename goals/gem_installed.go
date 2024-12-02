package goals

import (
	"fmt"
	"os/exec"
	"strings"
)

type GemInstalled struct {
	Name string
}

func (g GemInstalled) Description() string {
	return fmt.Sprintf("Installing Ruby gem %v", g.Name)
}

func (g GemInstalled) HashIdentity() string {
	return fmt.Sprintf("GemInstalled %v", g)
}

func (g GemInstalled) IsAchieved() bool {
	// Get raw output
	cmd := exec.Command("shadowenv", "exec", "--", "gem", "list", "-i", fmt.Sprintf("^%v$", g.Name))

	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.TrimSpace(string(out)) == "true"
}

func (g GemInstalled) Achieve() error {
	cmd := exec.Command("shadowenv", "exec", "--", "gem", "install", g.Name)

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v\n\n%v", err, string(stdoutStderr))
	}

	return nil
}
