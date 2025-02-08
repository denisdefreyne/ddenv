package goals

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"denisdefreyne.com/x/ddenv/core"
)

type NodeInstalled struct {
	Version string
}

func (g NodeInstalled) Description() string {
	return fmt.Sprintf("Installing Node %v", g.Version)
}

func (g NodeInstalled) HashIdentity() string {
	return fmt.Sprintf("NodeInstalled %v", g)
}

func (g NodeInstalled) IsAchieved() bool {
	_, err := os.Lstat(g.path())
	return err == nil
}

func (g NodeInstalled) Achieve() error {
	cmd := exec.Command("node-build", g.Version, g.path())

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v\n\n%v", err, string(stdoutStderr))
	}

	return nil
}

func (g NodeInstalled) SubGoals() []core.Goal {
	return []core.Goal{
		HomebrewPackageInstalled{PackageName: "node-build"},
	}
}

func (g NodeInstalled) path() string {
	usr, _ := user.Current()
	homeDir := usr.HomeDir
	return filepath.Join(homeDir, ".node-versions", g.Version)
}
