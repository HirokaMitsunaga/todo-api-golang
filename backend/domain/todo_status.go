package domain

import (
	"errors"
	"fmt"
)

type TodoStatus interface {
	Name() string
	transitionTo(next TodoStatus) (TodoStatus, error)
}

var ErrInvalidTodoStatusTransition = errors.New(
	"invalid todo status transition",
)

func invalidTodoStatusTransitionError(from TodoStatus, to TodoStatus) error {
	fromName := "<nil>"
	if from != nil {
		fromName = from.Name()
	}

	toName := "<nil>"
	if to != nil {
		toName = to.Name()
	}

	return fmt.Errorf(
		"%w: %s -> %s",
		ErrInvalidTodoStatusTransition,
		fromName,
		toName,
	)
}
