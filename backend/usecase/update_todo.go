package usecase

import (
	"errors"
	"todo-api/domain"
	"todo-api/domain/repository"

	"github.com/oklog/ulid/v2"
)

type InputUpdateTodoUsecase struct {
	Id       string
	Title    string
	Status   string
	UserId   string
	priority uint16
}

type updateTodoUsecase struct {
	tr repository.ITodoRepository
}

func NewUpdateTodoUsecase(tr repository.ITodoRepository) *updateTodoUsecase {
	return &updateTodoUsecase{tr}
}

func (ut *updateTodoUsecase) Execute(input InputUpdateTodoUsecase) error {
	id, err := ulid.Parse(input.Id)
	if err != nil {
		return err
	}
	todo, err := ut.tr.FindById(id)
	if err != nil {
		if errors.Is(err, repository.ErrTodoNotFoundRepository) {
			return ErrTodoNotFoundUsecase
		}
		return err
	}

	title, err := domain.NewTitle(input.Title)
	if err != nil {
		return err
	}
	status, err := domain.NewTodoStatus(input.Status)
	if err != nil {
		return err
	}
	userId, err := ulid.Parse(input.UserId)
	if err != nil {
		return err
	}
	priority, err := domain.NewPriority(input.priority)
	if err != nil {
		return err
	}

	updateTodo, err := todo.Update(title, status, userId, priority)
	if err != nil {
		return err
	}

	if err = ut.tr.Update(updateTodo); err != nil {
		if errors.Is(err, repository.ErrTodoNotFoundRepository) {
			return ErrTodoNotFoundUsecase
		}
		return err
	}

	return nil
}

var ErrTodoNotFoundUsecase = errors.New("todo is not Found")
