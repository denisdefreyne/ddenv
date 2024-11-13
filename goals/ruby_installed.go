package goals

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"denisdefreyne.com/x/ddenv/core"
)

func init() {
	core.RegisterGoal("ruby", func(value interface{}) (core.Goal, error) {
		if rubyVersionBytes, err := os.ReadFile(".ruby-version"); err != nil {
			return nil, fmt.Errorf("expected .ruby-version to exist")
		} else {
			rubyVersionString := strings.TrimSpace(string(rubyVersionBytes))
			return RubyInstalled{Version: rubyVersionString}, nil
		}
	})
}

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

	versionStartsWithDigit, err := regexp.Match(`^[0-9]`, []byte(g.Version))
	if err != nil {
		panic(err)
	}

	prefix := ""
	if versionStartsWithDigit {
		prefix = "ruby-"
	}

	result := filepath.Join(homeDir, ".rubies", fmt.Sprintf("%v%v", prefix, g.Version))

	return result
}
