package domain

type completedStatus struct{}

func (completedStatus) Name() string {
	return "COMPLETED"
}

func (completedStatus) transitionTo(next TodoStatus) (TodoStatus, error) {
	return nil, invalidTodoStatusTransitionError(completedStatus{}, next)
}
