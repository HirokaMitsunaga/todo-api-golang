package domain

type pendingStatus struct{}

func (pendingStatus) Name() string {
	return "PENDING"
}

func (pendingStatus) transitionTo(next TodoStatus) (TodoStatus, error) {
	if next != nil && next.Name() == "IN_PROGRESS" {
		return next, nil
	}
	return nil, invalidTodoStatusTransitionError(pendingStatus{}, next)
}
