package usecase

import (
	"todo-api/domain"

	"github.com/oklog/ulid/v2"
)

type fakeUserRepository struct {
	user        *domain.User
	findByID    ulid.ULID
	updatedUser *domain.User
	deletedID   ulid.ULID
	findErr     error
	updateErr   error
	deleteErr   error
}

func (f *fakeUserRepository) FindById(id ulid.ULID) (*domain.User, error) {
	f.findByID = id
	if f.findErr != nil {
		return nil, f.findErr
	}
	return f.user, nil
}

func (f *fakeUserRepository) Create(*domain.User) error {
	return nil
}

func (f *fakeUserRepository) Update(user *domain.User) error {
	f.updatedUser = user
	return f.updateErr
}

func (f *fakeUserRepository) Delete(id ulid.ULID) error {
	f.deletedID = id
	return f.deleteErr
}
