package repository

import (
	"todo-api/domain"

	"github.com/oklog/ulid/v2"
)

type IUserRepository interface {
	FindById(id ulid.ULID) (*domain.User, error)
	Create(user *domain.User) error
	Update(user *domain.User) error
	Delete(id ulid.ULID) error
}
