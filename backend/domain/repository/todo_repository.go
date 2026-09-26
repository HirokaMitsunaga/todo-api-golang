package repository

import "todo-api/domain"

type ITodoRepository interface {
	save(todo *domain.Todo) error
	delete(id string) error
}
