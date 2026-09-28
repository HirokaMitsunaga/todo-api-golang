package usecase

import (
	"errors"
	"testing"

	"github.com/oklog/ulid/v2"
)

func TestCreateTodoUsecase_Execute(t *testing.T) {
	userID := ulid.Make()
	repositoryMock := &fakeTodoRepository{}
	usecase := NewCreateTodoUsecase(repositoryMock)

	err := usecase.Execute(CreateTodoInput{
		Title:    "Buy milk",
		Status:   "PENDING",
		UserID:   userID.String(),
		Priority: 5,
	})
	if err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}

	if repositoryMock.createdTodo == nil {
		t.Fatal("Create() was not called")
	}
	if got := repositoryMock.createdTodo.Title().Value(); got != "Buy milk" {
		t.Errorf("created title = %q, want %q", got, "Buy milk")
	}
	if got := repositoryMock.createdTodo.Status().Name(); got != "PENDING" {
		t.Errorf("created status = %q, want %q", got, "PENDING")
	}
	if got := repositoryMock.createdTodo.UserID(); got != userID {
		t.Errorf("created user ID = %v, want %v", got, userID)
	}
	if got := repositoryMock.createdTodo.Priority().Value(); got != 5 {
		t.Errorf("created priority = %d, want %d", got, 5)
	}
	if repositoryMock.createdTodo.ID() == (ulid.ULID{}) {
		t.Error("created todo ID is zero")
	}
}

func TestCreateTodoUsecase_Execute_invalidInput(t *testing.T) {
	validUserID := ulid.Make().String()
	tests := []struct {
		name  string
		input CreateTodoInput
	}{
		{
			name: "invalid title",
			input: CreateTodoInput{
				Title:    "ab",
				Status:   "PENDING",
				UserID:   validUserID,
				Priority: 5,
			},
		},
		{
			name: "invalid status",
			input: CreateTodoInput{
				Title:    "Buy milk",
				Status:   "UNKNOWN",
				UserID:   validUserID,
				Priority: 5,
			},
		},
		{
			name: "invalid user ID",
			input: CreateTodoInput{
				Title:    "Buy milk",
				Status:   "PENDING",
				UserID:   "invalid",
				Priority: 5,
			},
		},
		{
			name: "invalid priority",
			input: CreateTodoInput{
				Title:    "Buy milk",
				Status:   "PENDING",
				UserID:   validUserID,
				Priority: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repositoryMock := &fakeTodoRepository{}
			usecase := NewCreateTodoUsecase(repositoryMock)

			if err := usecase.Execute(tt.input); err == nil {
				t.Fatal("Execute() succeeded unexpectedly")
			}
			if repositoryMock.createdTodo != nil {
				t.Fatal("Create() was called for invalid input")
			}
		})
	}
}

func TestCreateTodoUsecase_Execute_propagatesRepositoryError(t *testing.T) {
	createErr := errors.New("create failed")
	repositoryMock := &fakeTodoRepository{createErr: createErr}
	usecase := NewCreateTodoUsecase(repositoryMock)

	err := usecase.Execute(CreateTodoInput{
		Title:    "Buy milk",
		Status:   "PENDING",
		UserID:   ulid.Make().String(),
		Priority: 5,
	})
	if !errors.Is(err, createErr) {
		t.Fatalf("Execute() error = %v, want %v", err, createErr)
	}
}
