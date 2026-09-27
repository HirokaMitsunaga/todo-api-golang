package repository

import (
	"todo-api/domain"

	"github.com/oklog/ulid/v2"
)

type IUserRepository interface {
	findById(id ulid.ULID) (*domain.User, error)
	save(todo *domain.Todo) error
	delete(id ulid.ULID) error
}
