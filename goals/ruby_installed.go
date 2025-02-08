package goals

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"

	"denisdefreyne.com/x/ddenv/core"
)

type RubyInstalled struct {
	Version string
}

func (g RubyInstalled) Description() string {
	return fmt.Sprintf("Installing Ruby %v", g.Version)
}

func (g RubyInstalled) HashIdentity() string {
	return fmt.Sprintf("RubyInstalled %v", g)
}

func (g RubyInstalled) IsAchieved() bool {
	_, err := os.Lstat(g.path())
	return err == nil
}

func (g RubyInstalled) Achieve() error {
	cmd := exec.Command("ruby-install", "--cleanup", g.Version)

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v\n\n%v", err, string(stdoutStderr))
	}

	return nil
}

func (g RubyInstalled) SubGoals() []core.Goal {
	return []core.Goal{
		HomebrewPackageInstalled{PackageName: "ruby-install"},
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
