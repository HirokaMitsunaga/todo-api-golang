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

var ErrInvalidTodoStatus = errors.New(
	"invalid todo status",
)

func NewTodoStatus(name string) (TodoStatus, error) {
	switch name {
	case "PENDING":
		return pendingStatus{}, nil
	case "IN_PROGRESS":
		return inProgressStatus{}, nil
	case "COMPLETED":
		return completedStatus{}, nil
	default:
		return nil, ErrInvalidTodoStatus
	}
}

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
