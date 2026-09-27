package repository

import (
	"errors"
	"todo-api/domain"
	"todo-api/domain/repository"
	"todo-api/infra/gorm/model"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) repository.ITodoRepository {
	return &todoRepository{db}
}

func (td *todoRepository) FindById(id ulid.ULID) (*domain.Todo, error) {
	var todo model.Todo
	if err := td.db.Where("id=?", id.String()).First(&todo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrTodoNotFound
		}
		return nil, err
	}
	return domain.ReconstructTodo(
		id,
		*domain.ReconstructTitle(todo.Title),
		domain.ReconstructTodoStatus(todo.Status),
		ulid.MustParse(todo.UserID),
		*domain.ReconstructPriority(todo.Priority),
	), nil
}

func (td *todoRepository) Create(todo *domain.Todo) error {

	if err := td.db.Create(&model.Todo{
		ID:       todo.ID().String(),
		Title:    todo.Title().Value(),
		Status:   todo.Status().Name(),
		UserID:   todo.UserID().String(),
		Priority: todo.Priority().Value(),
	}).Error; err != nil {
		return err
	}
	return nil
}

func (td *todoRepository) Update(todo *domain.Todo) error {

	result := td.db.Updates(&model.Todo{
		ID:       todo.ID().String(),
		Title:    todo.Title().Value(),
		Status:   todo.Status().Name(),
		UserID:   todo.UserID().String(),
		Priority: todo.Priority().Value(),
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return repository.ErrTodoNotFound
	}

	return nil
}

func (td *todoRepository) Delete(id ulid.ULID) error {
	result := td.db.Where("id=?", id.String()).Delete(&model.Todo{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return repository.ErrTodoNotFound
	}
	return nil
}
