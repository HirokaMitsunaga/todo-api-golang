package repository

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/oklog/ulid/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"todo-api/domain"
	"todo-api/domain/repository"
)

func newMockUserRepository(t *testing.T) (*userRepository, sqlmock.Sqlmock) {
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
			TranslateError:         true,
		},
	)
	if err != nil {
		t.Fatalf("failed to open gorm DB: %v", err)
	}

	return &userRepository{db: db}, mock
}

func newTestUser(t *testing.T) *domain.User {
	t.Helper()

	user, err := domain.Reconstructor(
		ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV"),
		"Taro",
		"taro@example.com",
		"hashed-password",
	)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	return user
}

func userSelectQuery() string {
	return `SELECT \* FROM "users" WHERE id\s*=\s*\$1 ORDER BY "users"\."id" LIMIT \$2`
}

func TestUserRepository_FindById(t *testing.T) {
	repo, mock := newMockUserRepository(t)
	user := newTestUser(t)

	mock.ExpectQuery(userSelectQuery()).
		WithArgs(user.ID().String(), 1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"name",
			"email",
			"password",
		}).AddRow(
			user.ID().String(),
			user.Name(),
			user.Email(),
			user.Password(),
		))

	got, err := repo.FindById(user.ID())
	if err != nil {
		t.Fatalf("FindById() error = %v", err)
	}

	if got.ID() != user.ID() {
		t.Fatalf("FindById() ID = %s, want %s", got.ID(), user.ID())
	}
	if got.Name() != user.Name() {
		t.Fatalf("FindById() Name = %s, want %s", got.Name(), user.Name())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}

func TestUserRepository_FindById_NotFound(t *testing.T) {
	repo, mock := newMockUserRepository(t)
	user := newTestUser(t)

	mock.ExpectQuery(userSelectQuery()).
		WithArgs(user.ID().String(), 1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"name",
			"email",
			"password",
		}))

	_, err := repo.FindById(user.ID())
	if !errors.Is(err, repository.ErrUserNotFoundRepository) {
		t.Fatalf("FindById() error = %v, want ErrUserNotFoundRepository", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}

func TestUserRepository_Create(t *testing.T) {
	repo, mock := newMockUserRepository(t)
	user := newTestUser(t)

	mock.ExpectExec(`INSERT INTO "users"`).
		WithArgs(
			user.ID().String(),
			user.Name(),
			user.Email(),
			user.Password(),
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.Create(user); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}

func TestUserRepository_Create_Duplicated(t *testing.T) {
	repo, mock := newMockUserRepository(t)
	user := newTestUser(t)

	mock.ExpectExec(`INSERT INTO "users"`).
		WithArgs(
			user.ID().String(),
			user.Name(),
			user.Email(),
			user.Password(),
		).
		WillReturnError(gorm.ErrDuplicatedKey)

	err := repo.Create(user)
	if !errors.Is(err, repository.ErrUserDuplicatedRepository) {
		t.Fatalf("Create() error = %v, want ErrUserDuplicatedRepository", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}

func TestUserRepository_Update(t *testing.T) {
	repo, mock := newMockUserRepository(t)
	user := newTestUser(t)

	mock.ExpectExec(`UPDATE "users" SET`).
		WithArgs(
			user.Name(),
			user.Email(),
			user.Password(),
			user.ID().String(),
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.Update(user); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}

func TestUserRepository_Update_Duplicated(t *testing.T) {
	repo, mock := newMockUserRepository(t)
	user := newTestUser(t)

	mock.ExpectExec(`UPDATE "users" SET`).
		WithArgs(
			user.Name(),
			user.Email(),
			user.Password(),
			user.ID().String(),
		).
		WillReturnError(gorm.ErrDuplicatedKey)

	err := repo.Update(user)
	if !errors.Is(err, repository.ErrUserDuplicatedRepository) {
		t.Fatalf("Update() error = %v, want ErrUserDuplicatedRepository", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}

func TestUserRepository_Update_NotFound(t *testing.T) {
	repo, mock := newMockUserRepository(t)
	user := newTestUser(t)

	mock.ExpectExec(`UPDATE "users" SET`).
		WithArgs(
			user.Name(),
			user.Email(),
			user.Password(),
			user.ID().String(),
		).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.Update(user)
	if !errors.Is(err, repository.ErrUserNotFoundRepository) {
		t.Fatalf("Update() error = %v, want ErrUserNotFoundRepository", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}

func TestUserRepository_Delete(t *testing.T) {
	repo, mock := newMockUserRepository(t)
	user := newTestUser(t)

	mock.ExpectExec(`DELETE FROM "users" WHERE id = \$1`).
		WithArgs(user.ID().String()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.Delete(user.ID()); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}

func TestUserRepository_Delete_NotFound(t *testing.T) {
	repo, mock := newMockUserRepository(t)
	user := newTestUser(t)

	mock.ExpectExec(`DELETE FROM "users" WHERE id = \$1`).
		WithArgs(user.ID().String()).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.Delete(user.ID())
	if !errors.Is(err, repository.ErrUserNotFoundRepository) {
		t.Fatalf("Delete() error = %v, want ErrUserNotFoundRepository", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}
