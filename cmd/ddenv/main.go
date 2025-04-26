package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"

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
	Up []any
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
func removeDuplicateGoals(input []core.Goal) []core.Goal {
	res := make([]core.Goal, 0, len(input))
	seen := make(map[string]bool, len(input))

	for _, goal := range input {
		goalHashIdentity := goal.HashIdentity()

		if !seen[goalHashIdentity] {
			res = append(res, goal)
			seen[goalHashIdentity] = true
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
		var value any

		if s, ok := entry.(string); ok {
			key = s
		} else if h, ok := entry.(map[any]any); ok {
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
				return nil, fmt.Errorf("%v goal: %v", key, err)
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
	uniqueGoals := removeDuplicateGoals(flattenedGoals)

	return uniqueGoals, nil
}

func flattenGoals(out []core.Goal, inGoal core.Goal) []core.Goal {
	if withSubGoals, ok := inGoal.(core.WithSubGoals); ok {
		for _, g := range withSubGoals.SubGoals() {
			out = flattenGoals(out, g)
		}
	}

	out = append(out, inGoal)

	return out
}

func updateStatus(rowDelta int, colDelta int, status string) {
	// Save cursor position
	fmt.Printf("\033[s")

	// Move to line and column
	fmt.Printf("\r") // beginning of line
	fmt.Printf("\033[%dA", rowDelta)
	fmt.Printf("\033[%dC", colDelta)

	// Update status
	fmt.Printf("\033[K")
	fmt.Print(status)

	// Restore cursor position
	fmt.Printf("\033[u")
}

func calcUseColor() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	if os.Getenv("TERM") == "dumb" {
		return false
	}

	return true
}

var useColor = calcUseColor()

func gcManagedShadowenvFiles(goals []core.Goal) {
	// Get Shadowenv dir contents.
	f, err := os.Open(".shadowenv.d")
	if err != nil {
		return
	}
	files, err := f.Readdir(0)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Collect all file paths for files that exist. This array will be modified
	// later, so that it only contains the file paths that need to be removed.
	filePaths := make(map[string]struct{})
	for _, v := range files {
		name := ".shadowenv.d/" + v.Name()
		if strings.HasSuffix(name, ".lisp") {
			filePaths[name] = struct{}{}
		}
	}

	// Remove correctly managed file paths from `filePaths`.
	for _, goal := range goals {
		if g, ok := goal.(core.WithManagedShadowenvFilePaths); ok {
			for _, filePath := range g.ManagedShadowenvFilePaths() {
				_, found := filePaths[filePath]
				if found {
					delete(filePaths, filePath)
				}
			}
		}
	}

	// Delete obsolete files.
	if len(filePaths) > 0 {
		for filePath := range filePaths {
			os.Remove(filePath)
		}
	}
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

	red := ""
	reset := ""
	if useColor {
		red = "\u001B[31m"
		reset = "\u001B[0m"
	}

	// Execute
	for idx, goal := range gs {
		rowDelta := len(gs) - idx
		colDelta := maxDescriptionLen + 2

		updateStatus(rowDelta, colDelta, "checking...")
		var isAchieved = true
		if withAchieve, ok := goal.(core.WithAchieve); ok {
			isAchieved = withAchieve.IsAchieved()
		}

		if isAchieved {
			updateStatus(rowDelta, colDelta, "skipped")
		} else {
			updateStatus(rowDelta, colDelta, "working...")
			if withAchieve, ok := goal.(core.WithAchieve); ok {
				err = withAchieve.Achieve()
			}
			if err != nil {
				updateStatus(rowDelta, colDelta, fmt.Sprintf("%v%v%v", red, "failed", reset))
				fmt.Printf("\n%v%v%v", red, err, reset)
				return
			} else {
				updateStatus(rowDelta, colDelta, "done")
			}
		}
	}

	gcManagedShadowenvFiles(gs)
}
