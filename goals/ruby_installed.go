package goals

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"denisdefreyne.com/x/ddenv/core"
)

type RubyInstalled struct {
	Version string
}

func (g RubyInstalled) Description() string {
	return fmt.Sprintf("Installing Ruby %v", g.Version)
}

func (g RubyInstalled) IsAchieved() bool {
	// Get path to Ruby installation directory
	usr, _ := user.Current()
	homeDir := usr.HomeDir
	path := filepath.Join(homeDir, ".rubies", fmt.Sprintf("ruby-%v", g.Version))

	// Check
	_, err := os.Lstat(path)
	return err == nil
}

func (g RubyInstalled) Achieve() error {
	rubyInstallCmd := exec.Command("ruby-install", "--cleanup", g.Version)
	if err := rubyInstallCmd.Run(); err != nil {
		return err
	}

	return nil
}

func (g RubyInstalled) PreGoals() []core.Goal {
	return []core.Goal{
		HomebrewPackageInstalled{PackageName: "ruby-install"},
	}
}
