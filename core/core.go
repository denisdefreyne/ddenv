package core

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
