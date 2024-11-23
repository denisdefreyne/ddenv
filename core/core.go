package core

import "fmt"

type Goal interface {
	Description() string
	IsAchieved() bool
	Achieve() error
}

type WithPreGoals interface {
	PreGoals() []Goal
}

type WithPostGoals interface {
	PostGoals() []Goal
}

var goalFnsByName map[string]func(value interface{}) (Goal, error)

func init() {
	goalFnsByName = make(map[string]func(value interface{}) (Goal, error))
}

func RegisterGoal(name string, fn func(value interface{}) (Goal, error)) {
	if _, ok := goalFnsByName[name]; ok {
		panic(fmt.Sprintf("a goal named %q is already registered", name))
	}

	goalFnsByName[name] = fn
}

func FindGoal(name string) func(value interface{}) (Goal, error) {
	return goalFnsByName[name]
}
