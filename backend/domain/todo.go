package domain

import "github.com/oklog/ulid/v2"

type Todo struct {
	id       ulid.ULID
	title    Title
	status   TodoStatus
	userId   ulid.ULID
	priority Priority
}

func NewTodo(title Title, status TodoStatus, userId ulid.ULID, priority Priority) (*Todo, error) {
	return &Todo{
		id:       ulid.Make(),
		title:    title,
		status:   status,
		userId:   userId,
		priority: priority,
	}, nil
}

func (t *Todo) updateStatus(status TodoStatus) (*Todo, error) {
	nextStatus, err := t.status.transitionTo(status)
	if err != nil {
		return nil, err
	}

	updated := *t
	updated.status = nextStatus
	return &updated, nil
}
