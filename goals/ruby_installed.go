package goals

import (
	"bytes"
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
	_, err := os.Lstat(g.path())
	return err == nil
}

func (g RubyInstalled) Achieve() error {
	rubyInstallCmd := exec.Command("ruby-install", "--cleanup", g.Version)

	var stdoutBuf, stderrBuf bytes.Buffer
	rubyInstallCmd.Stdout = &stdoutBuf
	rubyInstallCmd.Stderr = &stderrBuf

	if err := rubyInstallCmd.Run(); err != nil {
		return fmt.Errorf(
			"%v:\n===[ stderr ]=======\n%v\n\n===[ stdout ]=======\n%v",
			err,
			stderrBuf.String(),
			stdoutBuf.String(),
		)
	}

	return nil
}

func (g RubyInstalled) PreGoals() []core.Goal {
	return []core.Goal{
		HomebrewPackageInstalled{PackageName: "ruby-install"},
	}
}

func (g RubyInstalled) PostGoals() []core.Goal {
	return []core.Goal{
		RubyShadowenvCreated{Version: g.Version, Path: g.path()},
	}
}

func (g RubyInstalled) path() string {
	usr, _ := user.Current()
	homeDir := usr.HomeDir
	return filepath.Join(homeDir, ".rubies", fmt.Sprintf("ruby-%v", g.Version))
}
