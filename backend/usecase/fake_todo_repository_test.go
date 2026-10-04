package usecase

import (
	"todo-api/domain"

	"github.com/oklog/ulid/v2"
)

type fakeTodoRepository struct {
	todo        *domain.Todo
	findByID    ulid.ULID
	createdTodo *domain.Todo
	updatedTodo *domain.Todo
	findErr     error
	createErr   error
	updateErr   error
}

func (f *fakeTodoRepository) FindById(id ulid.ULID) (*domain.Todo, error) {
	f.findByID = id
	if f.findErr != nil {
		return nil, f.findErr
	}
	return f.todo, nil
}

func (f *fakeTodoRepository) Create(todo *domain.Todo) error {
	f.createdTodo = todo
	return f.createErr
}

func (f *fakeTodoRepository) Update(todo *domain.Todo) error {
	f.updatedTodo = todo
	return f.updateErr
}

func (f *fakeTodoRepository) Delete(ulid.ULID) error {
	return nil
}
