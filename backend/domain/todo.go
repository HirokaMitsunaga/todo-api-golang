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

func ReconstructTodo(id ulid.ULID, title Title, status TodoStatus, userId ulid.ULID, priority Priority) *Todo {
	return &Todo{
		id:       id,
		title:    title,
		status:   status,
		userId:   userId,
		priority: priority,
	}
}

func (t *Todo) updateStatus(status TodoStatus) (*Todo, error) {
	nextStatus, err := t.status.transitionTo(status)
	if err != nil {
		return nil, err
	}
	return &Todo{
		id:       t.id,
		title:    t.title,
		status:   nextStatus,
		userId:   t.userId,
		priority: t.priority,
	}, nil

}

func (t *Todo) Update(title Title, status TodoStatus, priority Priority) (*Todo, error) {
	updated, err := t.updateStatus(status)
	if err != nil {
		return nil, err
	}
	return &Todo{
		id:       t.id,
		title:    title,
		status:   updated.status,
		userId:   t.userId,
		priority: priority,
	}, nil

}

func (t *Todo) ID() ulid.ULID {
	return t.id
}

func (t *Todo) Title() Title {
	return t.title
}

func (t *Todo) Status() TodoStatus {
	return t.status
}

func (t *Todo) UserID() ulid.ULID {
	return t.userId
}

func (t *Todo) Priority() Priority {
	return t.priority
}
