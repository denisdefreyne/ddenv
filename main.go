package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

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

type Goal interface {
	Description() string
	IsAchieved() bool
	Achieve() error
}

func ReadGoals() ([]Goal, error) {
	// Get config
	config, err := ReadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	var gs []Goal

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
				gs = append(gs, goals.HomebrewPackageInstalled{PackageName: packageName})
			} else {
				return nil, fmt.Errorf("homebrew goal: expected string package name")
			}

		case "ruby":
			if rubyVersionBytes, err := os.ReadFile(".ruby-version"); err != nil {
				return nil, fmt.Errorf("ruby goal: expected .ruby-version to exist")
			} else {
				rubyVersionString := strings.TrimSpace(string(rubyVersionBytes))
				gs = append(gs, goals.RubyInstalled{Version: rubyVersionString})
			}
		}
	}

	return gs, nil
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
