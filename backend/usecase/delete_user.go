package usecase

import (
	"errors"
	"todo-api/domain/repository"

	"github.com/oklog/ulid/v2"
)

type deleteUserUsecase struct {
	ur repository.IUserRepository
}

func NewDeleteUserUseCase(ur repository.IUserRepository) *deleteUserUsecase {
	return &deleteUserUsecase{ur}
}

func (du *deleteUserUsecase) Execute(id string) error {
	perseId, err := ulid.Parse(id)
	if err != nil {
		return err
	}

	if err := du.ur.Delete(perseId); err != nil {
		if errors.Is(err, repository.ErrUserNotFoundRepository) {
			return ErrUserNotFoundUsecase
		}
		return err
	}

	return nil
}

var ErrUserNotFoundUsecase = errors.New(
	"user is not Found",
)
