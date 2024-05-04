package goals

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

const ShadowenvInitialized_Contents = "shadowenv init fish | source"

type ShadowenvInitialized struct {
	Version string
	Path    string
}

func (g ShadowenvInitialized) Description() string {
	return "Adding Shadowenv to shell"
}

func (g ShadowenvInitialized) IsAchieved() bool {
	fishShellConfigPath := g.fishShellConfigPath()

	actualContents, err := os.ReadFile(fishShellConfigPath)
	if err != nil {
		// Kinda unrecoverable I guess...
		return false
	}

	actualLines := strings.Split(string(actualContents), "\n")
	for _, line := range actualLines {
		if line == ShadowenvInitialized_Contents {
			return true
		}
	}

	return false
}

func (g ShadowenvInitialized) Achieve() error {
	fishShellConfigPath := g.fishShellConfigPath()

	f, err := os.OpenFile(fishShellConfigPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	toAppend := fmt.Sprintf("\n# ddenv integration\n%v\n", ShadowenvInitialized_Contents)

	if _, err := f.Write([]byte(toAppend)); err != nil {
		f.Close()
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

func (g ShadowenvInitialized) fishShellConfigPath() string {
	usr, _ := user.Current()
	homeDir := usr.HomeDir
	return filepath.Join(homeDir, ".config", "fish", "config.fish")
}
