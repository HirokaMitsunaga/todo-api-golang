package repository

import "todo-api/domain"

type ITodoRepository interface {
	FindById(id string) error
	Create(todo *domain.Todo) error
	Update(todo *domain.Todo) error
	Delete(id string) error
}
