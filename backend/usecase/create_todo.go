package usecase

import (
	"todo-api/domain"
	"todo-api/domain/repository"

	"github.com/oklog/ulid/v2"
)

type CreateTodoInput struct {
	Title    string
	Status   string
	UserID   string
	Priority uint16
}

type createTodoUsecase struct {
	tr repository.ITodoRepository
}

func NewCreateTodoUsecase(tr repository.ITodoRepository) *createTodoUsecase {
	return &createTodoUsecase{tr}
}

func (ct *createTodoUsecase) Execute(input CreateTodoInput) error {
	title, err := domain.NewTitle(input.Title)
	if err != nil {
		return err
	}
	status, err := domain.NewTodoStatus(input.Status)
	if err != nil {
		return err
	}
	userID, err := ulid.Parse(input.UserID)
	if err != nil {
		return err
	}
	priority, err := domain.NewPriority(input.Priority)
	if err != nil {
		return err
	}

	todo, err := domain.NewTodo(title, status, userID, priority)
	if err != nil {
		return err
	}

	return ct.tr.Create(todo)
}
