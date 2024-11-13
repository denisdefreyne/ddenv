package main

import (
	flag "github.com/spf13/pflag"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"denisdefreyne.com/x/ddenv/core"
	_ "denisdefreyne.com/x/ddenv/goals"
)

var (
	buildVersion = "dev"
	buildCommit  = "none"
	buildDate    = "unknown"
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

// Returns a slice with all duplicates removed. Only the first unique element is
// retained.
func removeDuplicates[T comparable](input []T) []T {
	res := make([]T, 0, len(input))
	seen := make(map[T]bool, len(input))

	for _, element := range input {
		if !seen[element] {
			res = append(res, element)
			seen[element] = true
		}
	}

	return res
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

		goalFn := core.FindGoal(key)
		if goalFn != nil {
			if goal, err := goalFn(value); err != nil {
				return nil, fmt.Errorf("%v goal: %v",key, err)
			} else {
				simpleGoals = append(simpleGoals, goal)
			}
		} else {
			return nil, fmt.Errorf("unknown goal: %s", key)
		}
	}

	// Flatten goals
	var flattenedGoals []core.Goal
	for _, goal := range simpleGoals {
		flattenedGoals = flattenGoals(flattenedGoals, goal)
	}

	// Remove duplicate goals
	uniqueGoals := removeDuplicates(flattenedGoals)

	return uniqueGoals, nil
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
	// Parse args
	var showVersion = flag.BoolP("version", "v", false, "print version")
	var showHelp = flag.BoolP("help", "h", false, "print help")
	flag.Parse()

	// Handle args
	if *showVersion {
		fmt.Printf(
			"ddenv %s\ncommit %s\nbuilt at %s\n",
			buildVersion, buildCommit, buildDate)
		os.Exit(0)
	}
	if *showHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}

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
			updateStatus(rowDelta, colDelta, "skipped")
		} else {
			updateStatus(rowDelta, colDelta, "working...")
			err = goal.Achieve()
			if err != nil {
				updateStatus(rowDelta, colDelta, "failed")
				fmt.Printf("        %v\n", err)
				return
			} else {
				updateStatus(rowDelta, colDelta, "done")
			}
		}
	}
}
