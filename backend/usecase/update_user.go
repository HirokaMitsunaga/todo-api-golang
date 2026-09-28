package usecase

import (
	"todo-api/domain/repository"

	"github.com/oklog/ulid/v2"
)

// UpdateUserInput は User の更新内容を表す。
// 名前・メールアドレス・パスワードはすべて必須とする。
type UpdateUserInput struct {
	ID       string
	Name     string
	Email    string
	Password string
}

type updateUserUsecase struct {
	ur repository.IUserRepository
}

func NewUpdateUserUsecase(ur repository.IUserRepository) *updateUserUsecase {
	return &updateUserUsecase{ur: ur}
}

func (uu *updateUserUsecase) Execute(input UpdateUserInput) error {
	id, err := ulid.Parse(input.ID)
	if err != nil {
		return err
	}

	user, err := uu.ur.FindById(id)
	if err != nil {
		return err
	}

	user, err = user.WithName(input.Name)
	if err != nil {
		return err
	}

	user, err = user.WithEmail(input.Email)
	if err != nil {
		return err
	}

	user, err = user.WithPassword(input.Password)
	if err != nil {
		return err
	}

	return uu.ur.Update(user)
}
