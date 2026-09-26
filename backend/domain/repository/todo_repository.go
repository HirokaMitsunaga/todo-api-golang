package repository

import "todo-api/domain"

type ITodoRepository interface {
	findById(id string) error
	save(todo *domain.Todo) error
	delete(id string) error
}
