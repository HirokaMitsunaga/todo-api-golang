package domain

type inProgressStatus struct{}

func (inProgressStatus) Name() string {
	return "IN_PROGRESS"
}

func (inProgressStatus) transitionTo(next TodoStatus) (TodoStatus, error) {
	if next != nil && next.Name() == "COMPLETED" {
		return next, nil
	}
	return nil, invalidTodoStatusTransitionError(inProgressStatus{}, next)
}
