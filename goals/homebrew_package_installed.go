package goals

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"denisdefreyne.com/x/ddenv/core"
)

func init() {
	core.RegisterGoal("homebrew", func(value any) (core.Goal, error) {
		if packageName, ok := value.(string); ok {
			return HomebrewPackageInstalled{PackageName: packageName, IsCask: false}, nil
		} else {
			return nil, fmt.Errorf("expected string formula name")
		}
	})
}

func init() {
	core.RegisterGoal("homebrew_cask", func(value any) (core.Goal, error) {
		if packageName, ok := value.(string); ok {
			return HomebrewPackageInstalled{PackageName: packageName, IsCask: true}, nil
		} else {
			return nil, fmt.Errorf("expected string cask name")
		}
	})
}

type HomebrewPackageInstalled struct {
	PackageName string
	IsCask      bool
}

type brewInfoFormula struct {
	Installed []any
}

type brewInfoCask struct {
	Installed string
}

type brewInfo struct {
	Formulae []brewInfoFormula
	Casks    []brewInfoCask
}

func (g HomebrewPackageInstalled) Description() string {
	if g.IsCask {
		return fmt.Sprintf("Installing Homebrew cask ‘%v’", g.PackageName)
	} else {
		return fmt.Sprintf("Installing Homebrew formula ‘%v’", g.PackageName)
	}
}

func (g HomebrewPackageInstalled) HashIdentity() string {
	return fmt.Sprintf("HomebrewPackageInstalled %v", g)
}

func (g HomebrewPackageInstalled) IsAchieved() bool {
	// Get raw output
	var brewInfoCmd *exec.Cmd
	if g.IsCask {
		brewInfoCmd = exec.Command("brew", "info", "--json=v2", "--cask", g.PackageName)
	} else {
		brewInfoCmd = exec.Command("brew", "info", "--json=v2", g.PackageName)
	}
	brewInfoOut, err := brewInfoCmd.Output()
	if err != nil {
		return false
	}

	// Parse JSON
	var info brewInfo
	if err := json.Unmarshal(brewInfoOut, &info); err != nil {
		fmt.Println(err.Error())
		return false
	}

	// Check
	if len(info.Formulae) > 0 {
		return len(info.Formulae[0].Installed) > 0
	} else if len(info.Casks) > 0 {
		return info.Casks[0].Installed != ""
	}

	// This can’t really happen: either there are formulae or casks, but not both
	// can be missing. If an unknown name is given, the `brew info` command will
	// fail. Still, return `false` is a safe fallback.
	return false
}

func (g HomebrewPackageInstalled) Achieve() error {
	var cmd *exec.Cmd
	if g.IsCask {
		cmd = exec.Command("brew", "install", "--cask", g.PackageName)
	} else {
		cmd = exec.Command("brew", "install", g.PackageName)
	}

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v\n\n%v", err, string(stdoutStderr))
	}

	return nil
}
