package goals

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

const ShadowenvInitialized_Line_Bash = "eval \"$(shadowenv init bash)\""
const ShadowenvInitialized_Line_Zsh = "eval \"$(shadowenv init zsh)\""
const ShadowenvInitialized_Line_Fish = "shadowenv init fish | source"

type ShadowenvInitialized struct{}

func (g ShadowenvInitialized) Description() string {
	return "Adding Shadowenv to shell"
}

func (g ShadowenvInitialized) HashIdentity() string {
	return fmt.Sprintf("ShadowenvInitialized %v", g)
}

func (g ShadowenvInitialized) IsAchieved() bool {
	return g.isAchievedForFishShell() && g.isAchievedForBashShell() && g.isAchievedForZshShell()
}

func (g ShadowenvInitialized) Achieve() error {
	if err := g.achieveForFishShell(); err != nil {
		return err
	}

	if err := g.achieveForBashShell(); err != nil {
		return err
	}

	if err := g.achieveForZshShell(); err != nil {
		return err
	}

	return nil
}

// Utils

func (g ShadowenvInitialized) isAchievedFor(configPath string, expectedLine string) bool {
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		// If file does not exist, then it is not relevant, and thus achieved.
		return true
	}

	actualContents, err := os.ReadFile(configPath)
	if err != nil {
		// Kinda unrecoverable I guess...
		return false
	}

	actualLines := strings.Split(string(actualContents), "\n")
	for _, actualLine := range actualLines {
		if actualLine == expectedLine {
			return true
		}
	}

	return false
}

func (g ShadowenvInitialized) achieveFor(configPath string, expectedLine string) error {
	if g.isAchievedFor(configPath, expectedLine) {
		return nil
	}

	f, err := os.OpenFile(configPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	toAppend := fmt.Sprintf(
		"\n# ddenv integration\n%v\n",
		expectedLine,
	)

	if _, err := f.Write([]byte(toAppend)); err != nil {
		f.Close()
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

func (g ShadowenvInitialized) isAchievedForFishShell() bool {
	return g.isAchievedFor(g.fishShellConfigPath(), ShadowenvInitialized_Line_Fish)
}

func (g ShadowenvInitialized) isAchievedForBashShell() bool {
	return g.isAchievedFor(g.bashShellConfigPath(), ShadowenvInitialized_Line_Bash)
}

func (g ShadowenvInitialized) isAchievedForZshShell() bool {
	return g.isAchievedFor(g.zshShellConfigPath(), ShadowenvInitialized_Line_Zsh)
}

func (g ShadowenvInitialized) achieveForFishShell() error {
	return g.achieveFor(g.fishShellConfigPath(), ShadowenvInitialized_Line_Fish)
}

func (g ShadowenvInitialized) achieveForBashShell() error {
	return g.achieveFor(g.bashShellConfigPath(), ShadowenvInitialized_Line_Bash)
}

func (g ShadowenvInitialized) achieveForZshShell() error {
	return g.achieveFor(g.zshShellConfigPath(), ShadowenvInitialized_Line_Zsh)
}

func (g ShadowenvInitialized) fishShellConfigPath() string {
	usr, _ := user.Current()
	homeDir := usr.HomeDir
	return filepath.Join(homeDir, ".config", "fish", "config.fish")
}

func (g ShadowenvInitialized) bashShellConfigPath() string {
	usr, _ := user.Current()
	homeDir := usr.HomeDir
	return filepath.Join(homeDir, ".bash_profile")
}

func (g ShadowenvInitialized) zshShellConfigPath() string {
	usr, _ := user.Current()
	homeDir := usr.HomeDir
	return filepath.Join(homeDir, ".zshrc")
}
