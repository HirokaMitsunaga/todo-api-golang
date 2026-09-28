package repository

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/oklog/ulid/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"todo-api/domain"
	domainRepository "todo-api/domain/repository"
)

func newMockTodoRepository(t *testing.T) (*todoRepository, sqlmock.Sqlmock) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	db, err := gorm.Open(
		postgres.New(postgres.Config{
			Conn:                 sqlDB,
			PreferSimpleProtocol: true,
		}),
		&gorm.Config{
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		t.Fatalf("failed to open gorm DB: %v", err)
	}

	return &todoRepository{db: db}, mock
}

func newTestTodo(t *testing.T) *domain.Todo {
	t.Helper()

	title := domain.ReconstructTitle("Buy milk")
	priority := domain.ReconstructPriority(5)

	return domain.ReconstructTodo(
		ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV"),
		title,
		domain.ReconstructTodoStatus("PENDING"),
		ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAW"),
		priority,
	)
}

func todoSelectQuery() string {
	return `SELECT \* FROM "todos" WHERE id\s*=\s*\$1 ORDER BY "todos"\."id" LIMIT \$2`
}

func TestTodoRepository_FindById(t *testing.T) {
	repo, mock := newMockTodoRepository(t)
	todo := newTestTodo(t)

	mock.ExpectQuery(todoSelectQuery()).
		WithArgs(todo.ID().String(), 1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"title",
			"status",
			"user_id",
			"priority",
		}).AddRow(
			todo.ID().String(),
			todo.Title().Value(),
			todo.Status().Name(),
			todo.UserID().String(),
			todo.Priority().Value(),
		))

	got, err := repo.FindById(todo.ID())
	if err != nil {
		t.Fatalf("FindById() error = %v", err)
	}

	if got.ID() != todo.ID() {
		t.Fatalf("FindById() ID = %s, want %s", got.ID(), todo.ID())
	}
	if got.Title().Value() != todo.Title().Value() {
		t.Fatalf("FindById() Title = %s, want %s", got.Title().Value(), todo.Title().Value())
	}
	if got.Status().Name() != todo.Status().Name() {
		t.Fatalf("FindById() Status = %s, want %s", got.Status().Name(), todo.Status().Name())
	}
	if got.UserID() != todo.UserID() {
		t.Fatalf("FindById() UserID = %s, want %s", got.UserID(), todo.UserID())
	}
	if got.Priority().Value() != todo.Priority().Value() {
		t.Fatalf("FindById() Priority = %d, want %d", got.Priority().Value(), todo.Priority().Value())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}

func TestTodoRepository_FindById_NotFound(t *testing.T) {
	repo, mock := newMockTodoRepository(t)
	todo := newTestTodo(t)

	mock.ExpectQuery(todoSelectQuery()).
		WithArgs(todo.ID().String(), 1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"title",
			"status",
			"user_id",
			"priority",
		}))

	_, err := repo.FindById(todo.ID())
	if !errors.Is(err, domainRepository.ErrTodoNotFoundRepository) {
		t.Fatalf("FindById() error = %v, want ErrTodoNotFoundRepository", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}

func TestTodoRepository_Create(t *testing.T) {
	repo, mock := newMockTodoRepository(t)
	todo := newTestTodo(t)

	mock.ExpectExec(`INSERT INTO "todos"`).
		WithArgs(
			todo.ID().String(),
			todo.Title().Value(),
			todo.Status().Name(),
			todo.UserID().String(),
			todo.Priority().Value(),
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.Create(todo); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}

func TestTodoRepository_Update(t *testing.T) {
	repo, mock := newMockTodoRepository(t)
	todo := newTestTodo(t)

	mock.ExpectExec(`UPDATE "todos" SET`).
		WithArgs(
			todo.Title().Value(),
			todo.Status().Name(),
			todo.UserID().String(),
			todo.Priority().Value(),
			todo.ID().String(),
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.Update(todo); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}

func TestTodoRepository_Update_NotFound(t *testing.T) {
	repo, mock := newMockTodoRepository(t)
	todo := newTestTodo(t)

	mock.ExpectExec(`UPDATE "todos" SET`).
		WithArgs(
			todo.Title().Value(),
			todo.Status().Name(),
			todo.UserID().String(),
			todo.Priority().Value(),
			todo.ID().String(),
		).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.Update(todo)
	if !errors.Is(err, domainRepository.ErrTodoNotFoundRepository) {
		t.Fatalf("Update() error = %v, want ErrTodoNotFoundRepository", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}

func TestTodoRepository_Delete(t *testing.T) {
	repo, mock := newMockTodoRepository(t)
	todo := newTestTodo(t)

	mock.ExpectExec(`DELETE FROM "todos" WHERE id\s*=\s*\$1`).
		WithArgs(todo.ID().String()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.Delete(todo.ID()); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}

func TestTodoRepository_Delete_NotFound(t *testing.T) {
	repo, mock := newMockTodoRepository(t)
	todo := newTestTodo(t)

	mock.ExpectExec(`DELETE FROM "todos" WHERE id\s*=\s*\$1`).
		WithArgs(todo.ID().String()).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.Delete(todo.ID())
	if !errors.Is(err, domainRepository.ErrTodoNotFoundRepository) {
		t.Fatalf("Delete() error = %v, want ErrTodoNotFoundRepository", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}
