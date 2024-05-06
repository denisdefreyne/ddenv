package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"denisdefreyne.com/x/ddenv/core"
	"denisdefreyne.com/x/ddenv/goals"
)

type Config struct {
	Up []interface{}
}

func ReadConfig() (Config, error) {
	filename, _ := filepath.Abs("ddenv.yaml")

	// Read file
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}

	// Parse file
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func ReadGoals() ([]core.Goal, error) {
	// Get config
	config, err := ReadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	var simpleGoals []core.Goal

	for _, entry := range config.Up {
		var key string
		var value interface{}

		if s, ok := entry.(string); ok {
			key = s
		} else if h, ok := entry.(map[interface{}]interface{}); ok {
			if len(h) != 1 {
				return nil, fmt.Errorf("too many keys (expected only 1)")
			}

			// Get key
			for k, v := range h {
				if sk := k.(string); ok {
					key = sk
					value = v
					break
				} else {
					return nil, fmt.Errorf("expected a string key")
				}
			}
		} else {
			return nil, fmt.Errorf("unknown type (expected string or map)")
		}

		switch key {
		case "homebrew":
			if packageName, ok := value.(string); ok {
				simpleGoals = append(simpleGoals, goals.HomebrewPackageInstalled{PackageName: packageName})
			} else {
				return nil, fmt.Errorf("homebrew goal: expected string package name")
			}

		case "ruby":
			if rubyVersionBytes, err := os.ReadFile(".ruby-version"); err != nil {
				return nil, fmt.Errorf("ruby goal: expected .ruby-version to exist")
			} else {
				rubyVersionString := strings.TrimSpace(string(rubyVersionBytes))
				simpleGoals = append(simpleGoals, goals.RubyInstalled{Version: rubyVersionString})
			}

		case "bundle":
			simpleGoals = append(simpleGoals, goals.BundleInstalled{})

		case "node":
			if nodeVersion, ok := value.(string); ok {
				simpleGoals = append(simpleGoals, goals.NodeInstalled{Version: nodeVersion})
			} else {
				return nil, fmt.Errorf("node goal: expected string version")
			}

		case "npm":
			simpleGoals = append(simpleGoals, goals.NpmPackagesInstalled{})

			// TODO: handle unknown key
		}
	}

	// Flatten goals
	var flattenedGoals []core.Goal
	for _, goal := range simpleGoals {
		flattenedGoals = flattenGoals(flattenedGoals, goal)
	}

	// TODO: Remove duplicate goals

	return flattenedGoals, nil
}

func flattenGoals(out []core.Goal, inGoal core.Goal) []core.Goal {
	if withPreGoals, ok := inGoal.(core.WithPreGoals); ok {
		for _, g := range withPreGoals.PreGoals() {
			out = flattenGoals(out, g)
		}
	}

	out = append(out, inGoal)

	if withPostGoals, ok := inGoal.(core.WithPostGoals); ok {
		for _, g := range withPostGoals.PostGoals() {
			out = flattenGoals(out, g)
		}
	}

	return out
}

func updateStatus(rowDelta int, colDelta int, status string) {
	// Save cursor position
	fmt.Printf("\033[s")

	// Move to line and column
	fmt.Printf("\033[%dA", rowDelta)
	fmt.Printf("\033[%dC", colDelta)

	// Update status
	fmt.Printf("\033[K")
	fmt.Print(status)

	// Restore cursor position
	fmt.Printf("\033[u")
}

func main() {
	// Get goals
	gs, err := ReadGoals()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	// Get goal message size
	var maxDescriptionLen int
	for _, goal := range gs {
		thisLen := len(goal.Description())
		if thisLen > maxDescriptionLen {
			maxDescriptionLen = thisLen
		}
	}

	// Print pending status
	for _, goal := range gs {
		fmt.Printf("%-[2]*[1]s  pending\n", goal.Description(), maxDescriptionLen)
	}

	for idx, goal := range gs {
		rowDelta := len(gs) - idx
		colDelta := maxDescriptionLen + 2

		updateStatus(rowDelta, colDelta, "checking...")
		isAchieved := goal.IsAchieved()

		if isAchieved {
			updateStatus(rowDelta, colDelta, "previously achieved")
		} else {
			updateStatus(rowDelta, colDelta, "achieving...")
			err = goal.Achieve()
			if err != nil {
				updateStatus(rowDelta, colDelta, "failed")
				fmt.Printf("        %v\n", err)
			} else {
				updateStatus(rowDelta, colDelta, "newly achieved")
			}
		}
	}
}
