package repository

import "todo-api/domain"

type IUserRepository interface {
	findById(id string) error
	save(todo *domain.Todo) error
	delete(id string) error
}
