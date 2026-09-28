package usecase

import (
	"errors"
	"testing"

	"todo-api/domain/repository"

	"github.com/oklog/ulid/v2"
)

func TestDeleteUserUsecase_Execute(t *testing.T) {
	userID := ulid.Make()
	repositoryMock := &fakeUserRepository{}
	usecase := NewDeleteUserUseCase(repositoryMock)

	if err := usecase.Execute(userID.String()); err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}

	if repositoryMock.deletedID != userID {
		t.Fatalf("Delete() received %v, want %v", repositoryMock.deletedID, userID)
	}
}

func TestDeleteUserUsecase_Execute_invalidID(t *testing.T) {
	repositoryMock := &fakeUserRepository{}
	usecase := NewDeleteUserUseCase(repositoryMock)

	if err := usecase.Execute("invalid"); err == nil {
		t.Fatal("Execute() succeeded unexpectedly")
	}

	if repositoryMock.deletedID != (ulid.ULID{}) {
		t.Fatal("Delete() was called for an invalid ID")
	}
}

func TestDeleteUserUsecase_Execute_userNotFound(t *testing.T) {
	repositoryMock := &fakeUserRepository{
		deleteErr: repository.ErrUserNotFoundRepository,
	}
	usecase := NewDeleteUserUseCase(repositoryMock)

	err := usecase.Execute(ulid.Make().String())
	if !errors.Is(err, ErrUserNotFoundUsecase) {
		t.Fatalf("Execute() error = %v, want %v", err, ErrUserNotFoundUsecase)
	}
}

func TestDeleteUserUsecase_Execute_propagatesRepositoryError(t *testing.T) {
	deleteErr := errors.New("delete failed")
	repositoryMock := &fakeUserRepository{deleteErr: deleteErr}
	usecase := NewDeleteUserUseCase(repositoryMock)

	err := usecase.Execute(ulid.Make().String())
	if !errors.Is(err, deleteErr) {
		t.Fatalf("Execute() error = %v, want %v", err, deleteErr)
	}
}
