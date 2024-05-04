package main

import (
	"fmt"
	"os"
	"path/filepath"

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

func main() {
	// Get config
	config, err := ReadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	var gs []Goal

	for _, entry := range config.Up {
		fmt.Printf("%v, %t\n", entry, entry)

		var key string
		var value interface{}

		if s, ok := entry.(string); ok {
			key = s
		} else if h, ok := entry.(map[interface{}]interface{}); ok {
			if len(h) != 1 {
				// FIXME: error (too many keys)
			}

			// Get key
			for k, v := range h {
				if sk := k.(string); ok {
					key = sk
					value = v
					break
				} else {
					// FIXME: error (not a string key)
				}
			}
		} else {
			fmt.Printf("  error: %v\n", entry)
			// FIXME: error
		}

		fmt.Printf("  key: %v\n", key)
		fmt.Printf("  value: %v\n", value)

		switch key {
		case "homebrew":
			if packageName, ok := value.(string); ok {
				gs = append(gs, goals.HomebrewPackageInstalled{PackageName: packageName})
			} else {
				// FIXME: error
			}
		}
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
