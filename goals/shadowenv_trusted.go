package goals

import (
	"fmt"
	"os/exec"
)

type ShadowenvTrusted struct {
}

func (g ShadowenvTrusted) Description() string {
	return "Trusting Shadowenv"
}

func (g ShadowenvTrusted) HashIdentity() string {
	return fmt.Sprintf("ShadowenvTrusted %v", g)
}

func (g ShadowenvTrusted) IsAchieved() bool {
	cmd := exec.Command("shadowenv", "exec", "ls")

	err := cmd.Run()
	return err == nil
}

func (g ShadowenvTrusted) Achieve() error {
	cmd := exec.Command("shadowenv", "trust")

	err := cmd.Run()
	return err
}
