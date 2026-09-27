package repository

import (
	"errors"
	"todo-api/domain"

	"github.com/oklog/ulid/v2"
)

type ITodoRepository interface {
	FindById(id ulid.ULID) (*domain.Todo, error)
	Create(todo *domain.Todo) error
	Update(todo *domain.Todo) error
	Delete(id ulid.ULID) error
}

var ErrTodoNotFoundRepository = errors.New("todo not found")
