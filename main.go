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

func main() {
	gs, err := ReadGoals()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	for _, goal := range gs {
		fmt.Printf("--- %v\n", goal.Description())

		fmt.Printf("      checking...\n")
		isAchieved := goal.IsAchieved()

		if isAchieved {
			fmt.Printf("      previous achieved\n")
		} else {
			fmt.Printf("      achieving...\n")
			err = goal.Achieve()
			if err != nil {
				fmt.Printf("      failed\n")
				fmt.Printf("        %v\n", err)
			} else {
				fmt.Printf("      newly achieved\n")
			}
		}
	}
}
