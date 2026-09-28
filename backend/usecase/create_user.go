package usecase

import (
	"errors"
	"todo-api/domain"
	"todo-api/domain/repository"
)

type createUserUsecase struct {
	ur repository.IUserRepository
}

func NewCreateUserUsecase(ur repository.IUserRepository) *createUserUsecase {
	return &createUserUsecase{ur}
}

func (cu *createUserUsecase) Execute(name string, email string, password string) error {
	user, err := domain.NewUser(
		name,
		email,
		password,
	)
	if err != nil {
		return err
	}
	if err := cu.ur.Create(user); err != nil {
		if errors.Is(err, repository.ErrUserDuplicatedRepository) {
			return ErrUserDuplicatedUsecae
		}
		return err
	}
	return nil

}

var ErrUserDuplicatedUsecae = errors.New(
	"user is Duplicated",
)
