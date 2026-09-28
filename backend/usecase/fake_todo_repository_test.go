package usecase

import (
	"todo-api/domain"

	"github.com/oklog/ulid/v2"
)

type fakeTodoRepository struct {
	createdTodo *domain.Todo
	createErr   error
}

func (f *fakeTodoRepository) FindById(ulid.ULID) (*domain.Todo, error) {
	return nil, nil
}

func (f *fakeTodoRepository) Create(todo *domain.Todo) error {
	f.createdTodo = todo
	return f.createErr
}

func (f *fakeTodoRepository) Update(*domain.Todo) error {
	return nil
}

func (f *fakeTodoRepository) Delete(ulid.ULID) error {
	return nil
}
