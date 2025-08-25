package goals

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"denisdefreyne.com/x/ddenv/core"
)

func init() {
	core.RegisterGoal("homebrew_tap", func(value any) (core.Goal, error) {
		// Try with a name string instead of a JSON objct first
		if simpleName, ok := value.(string); ok {
			return HomebrewRepositoryTapped{Url: "", Name: simpleName}, nil
		}

		detailsMap, ok := value.(map[any]any)
		if !ok {
			return nil, fmt.Errorf("expected details map")
		}

		// Extract URL
		rawUrl, ok := detailsMap["url"]
		if !ok {
			return nil, fmt.Errorf("expected url")
		}
		url, ok := rawUrl.(string)
		if !ok {
			return nil, fmt.Errorf("expected string url")
		}

		// Extract user and repository
		rawName, ok := detailsMap["name"]
		if !ok {
			return nil, fmt.Errorf("expected name")
		}
		name, ok := rawName.(string)
		if !ok {
			return nil, fmt.Errorf("expected string name")
		}

		return HomebrewRepositoryTapped{Url: url, Name: name}, nil
	})
}

type HomebrewRepositoryTapped struct {
	Url  string
	Name string
}

type brewTapInfo struct {
	Installed bool
}

func (g HomebrewRepositoryTapped) Description() string {
	return fmt.Sprintf("Tapping Homebrew repository ‘%v’", g.Name)
}

func (g HomebrewRepositoryTapped) HashIdentity() string {
	return fmt.Sprintf("HomebrewRepositoryTapped %v", g)
}

func (g HomebrewRepositoryTapped) IsAchieved() bool {
	// Get raw output
	var brewTapInfoCmd *exec.Cmd
	brewTapInfoCmd = exec.Command("brew", "tap-info", "--json", g.Name)
	brewTapInfoOut, err := brewTapInfoCmd.Output()
	if err != nil {
		return false
	}

	// Parse JSON
	var tapInfo []brewTapInfo
	if err := json.Unmarshal(brewTapInfoOut, &tapInfo); err != nil {
		fmt.Println(err.Error())
		return false
	}

	// Check
	if len(tapInfo) < 1 {
		return false
	}
	return tapInfo[0].Installed
}

func (g HomebrewRepositoryTapped) Achieve() error {
	var cmd *exec.Cmd
	if len(g.Url) > 0 {
		cmd = exec.Command("brew", "tap", g.Name, g.Url)
	} else {
		cmd = exec.Command("brew", "tap", g.Name)
	}

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v\n\n%v", err, string(stdoutStderr))
	}

	return nil
}
