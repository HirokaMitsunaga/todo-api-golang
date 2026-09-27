package repository

import "todo-api/domain"

type ITodoRepository interface {
	FindById(id string) error
	Save(todo *domain.Todo) error
	Delete(id string) error
}
