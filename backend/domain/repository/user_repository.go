package repository

import (
	"todo-api/domain"

	"github.com/oklog/ulid/v2"
)

type IUserRepository interface {
	FindById(id ulid.ULID) (*domain.User, error)
	Save(todo *domain.Todo) error
	Delete(id ulid.ULID) error
}
