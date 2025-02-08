package core

import "fmt"

// The super-type for all goals.
//
// An instance of `Goal` also needs to implement either `WithAchieve` or
// `WithSubGoals`.
type Goal interface {
	// A one-line description of this goal, in imperative form, e.g. `Setting up
	// Shadowenv`.
	Description() string

	// A string that uniquely identifies this goal, both by its type and its
	// arguments. This is used for deduplication.
	HashIdentity() string
}

// An extension to `Goal` for goals that can explicitly be achieved.
type WithAchieve interface {
	// Returns true if the goal is already achieved, false otherwise.
	IsAchieved() bool

	// Attempt to achieve the goal.
	Achieve() error
}

// An extension to `Goal` for goals that have sub-goals. Sub-goals will be
// attempted to be achieved before the main goal itself.
type WithSubGoals interface {
	SubGoals() []Goal
}

var goalFnsByName map[string]func(value interface{}) (Goal, error)

func init() {
	goalFnsByName = make(map[string]func(value interface{}) (Goal, error))
}

// Register a goal with the given name. The `name` parameter is the identifier
// used in the `ddenv.yaml` file. The second argument is a function that will be
// called to instantiate the goal with the given argument value.
//
// For example:
//
//	core.RegisterGoal("bundle", func(value interface{}) (core.Goal, error) {
//	  return BundleInstalled{}, nil
//	})
func RegisterGoal(name string, fn func(value interface{}) (Goal, error)) {
	if _, ok := goalFnsByName[name]; ok {
		panic(fmt.Sprintf("a goal named %q is already registered", name))
	}

	goalFnsByName[name] = fn
}

func FindGoal(name string) func(value interface{}) (Goal, error) {
	return goalFnsByName[name]
}
