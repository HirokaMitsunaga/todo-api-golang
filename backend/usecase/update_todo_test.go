package usecase

import (
	"errors"
	"testing"

	"todo-api/domain"
	"todo-api/domain/repository"

	"github.com/oklog/ulid/v2"
)

func TestUpdateTodoUsecase_Execute(t *testing.T) {
	todoID := ulid.Make()
	userID := ulid.Make()
	todo := domain.ReconstructTodo(
		todoID,
		domain.ReconstructTitle("Buy milk"),
		domain.ReconstructTodoStatus("PENDING"),
		userID,
		domain.ReconstructPriority(1),
	)
	repositoryMock := &fakeTodoRepository{todo: todo}
	usecase := NewUpdateTodoUsecase(repositoryMock)

	err := usecase.Execute(InputUpdateTodoUsecase{
		Id:       todoID.String(),
		Title:    "Buy bread",
		Status:   "IN_PROGRESS",
		UserId:   userID.String(),
		priority: 5,
	})
	if err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}

	if repositoryMock.findByID != todoID {
		t.Fatalf("FindById() received %v, want %v", repositoryMock.findByID, todoID)
	}
	if repositoryMock.updatedTodo == nil {
		t.Fatal("Update() was not called")
	}
	if got := repositoryMock.updatedTodo.ID(); got != todoID {
		t.Errorf("updated ID = %v, want %v", got, todoID)
	}
	if got := repositoryMock.updatedTodo.Title().Value(); got != "Buy bread" {
		t.Errorf("updated title = %q, want %q", got, "Buy bread")
	}
	if got := repositoryMock.updatedTodo.Status().Name(); got != "IN_PROGRESS" {
		t.Errorf("updated status = %q, want %q", got, "IN_PROGRESS")
	}
	if got := repositoryMock.updatedTodo.UserID(); got != userID {
		t.Errorf("updated user ID = %v, want %v", got, userID)
	}
	if got := repositoryMock.updatedTodo.Priority().Value(); got != 5 {
		t.Errorf("updated priority = %d, want %d", got, 5)
	}
}

func TestUpdateTodoUsecase_Execute_invalidID(t *testing.T) {
	repositoryMock := &fakeTodoRepository{}
	usecase := NewUpdateTodoUsecase(repositoryMock)

	if err := usecase.Execute(InputUpdateTodoUsecase{Id: "invalid"}); err == nil {
		t.Fatal("Execute() succeeded unexpectedly")
	}
	if repositoryMock.findByID != (ulid.ULID{}) {
		t.Fatal("FindById() was called for an invalid ID")
	}
}

func TestUpdateTodoUsecase_Execute_invalidInput(t *testing.T) {
	todoID := ulid.Make()
	userID := ulid.Make()
	validInput := InputUpdateTodoUsecase{
		Id:       todoID.String(),
		Title:    "Buy bread",
		Status:   "IN_PROGRESS",
		UserId:   userID.String(),
		priority: 5,
	}

	tests := []struct {
		name  string
		input InputUpdateTodoUsecase
	}{
		{
			name: "invalid title",
			input: func() InputUpdateTodoUsecase {
				input := validInput
				input.Title = "ab"
				return input
			}(),
		},
		{
			name: "invalid status",
			input: func() InputUpdateTodoUsecase {
				input := validInput
				input.Status = "UNKNOWN"
				return input
			}(),
		},
		{
			name: "invalid user ID",
			input: func() InputUpdateTodoUsecase {
				input := validInput
				input.UserId = "invalid"
				return input
			}(),
		},
		{
			name: "invalid priority",
			input: func() InputUpdateTodoUsecase {
				input := validInput
				input.priority = 0
				return input
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repositoryMock := &fakeTodoRepository{
				todo: domain.ReconstructTodo(
					todoID,
					domain.ReconstructTitle("Buy milk"),
					domain.ReconstructTodoStatus("PENDING"),
					userID,
					domain.ReconstructPriority(1),
				),
			}
			usecase := NewUpdateTodoUsecase(repositoryMock)

			if err := usecase.Execute(tt.input); err == nil {
				t.Fatal("Execute() succeeded unexpectedly")
			}
			if repositoryMock.updatedTodo != nil {
				t.Fatal("Update() was called for invalid input")
			}
		})
	}
}

func TestUpdateTodoUsecase_Execute_todoNotFound(t *testing.T) {
	repositoryMock := &fakeTodoRepository{findErr: repository.ErrTodoNotFoundRepository}
	usecase := NewUpdateTodoUsecase(repositoryMock)

	err := usecase.Execute(InputUpdateTodoUsecase{Id: ulid.Make().String()})
	if !errors.Is(err, ErrTodoNotFoundUsecase) {
		t.Fatalf("Execute() error = %v, want %v", err, ErrTodoNotFoundUsecase)
	}
}

func TestUpdateTodoUsecase_Execute_propagatesFindError(t *testing.T) {
	findErr := errors.New("find failed")
	repositoryMock := &fakeTodoRepository{findErr: findErr}
	usecase := NewUpdateTodoUsecase(repositoryMock)

	err := usecase.Execute(InputUpdateTodoUsecase{Id: ulid.Make().String()})
	if !errors.Is(err, findErr) {
		t.Fatalf("Execute() error = %v, want %v", err, findErr)
	}
}

func TestUpdateTodoUsecase_Execute_invalidStatusTransition(t *testing.T) {
	todoID := ulid.Make()
	userID := ulid.Make()
	repositoryMock := &fakeTodoRepository{
		todo: domain.ReconstructTodo(
			todoID,
			domain.ReconstructTitle("Buy milk"),
			domain.ReconstructTodoStatus("PENDING"),
			userID,
			domain.ReconstructPriority(1),
		),
	}
	usecase := NewUpdateTodoUsecase(repositoryMock)

	err := usecase.Execute(InputUpdateTodoUsecase{
		Id:       todoID.String(),
		Title:    "Buy bread",
		Status:   "COMPLETED",
		UserId:   userID.String(),
		priority: 5,
	})
	if !errors.Is(err, domain.ErrInvalidTodoStatusTransition) {
		t.Fatalf("Execute() error = %v, want %v", err, domain.ErrInvalidTodoStatusTransition)
	}
	if repositoryMock.updatedTodo != nil {
		t.Fatal("Update() was called after an invalid status transition")
	}
}

func TestUpdateTodoUsecase_Execute_updateTodoNotFound(t *testing.T) {
	todoID := ulid.Make()
	userID := ulid.Make()
	repositoryMock := &fakeTodoRepository{
		todo: domain.ReconstructTodo(
			todoID,
			domain.ReconstructTitle("Buy milk"),
			domain.ReconstructTodoStatus("PENDING"),
			userID,
			domain.ReconstructPriority(1),
		),
		updateErr: repository.ErrTodoNotFoundRepository,
	}
	usecase := NewUpdateTodoUsecase(repositoryMock)

	err := usecase.Execute(InputUpdateTodoUsecase{
		Id:       todoID.String(),
		Title:    "Buy bread",
		Status:   "IN_PROGRESS",
		UserId:   userID.String(),
		priority: 5,
	})
	if !errors.Is(err, ErrTodoNotFoundUsecase) {
		t.Fatalf("Execute() error = %v, want %v", err, ErrTodoNotFoundUsecase)
	}
}

func TestUpdateTodoUsecase_Execute_propagatesUpdateError(t *testing.T) {
	todoID := ulid.Make()
	userID := ulid.Make()
	updateErr := errors.New("update failed")
	repositoryMock := &fakeTodoRepository{
		todo: domain.ReconstructTodo(
			todoID,
			domain.ReconstructTitle("Buy milk"),
			domain.ReconstructTodoStatus("PENDING"),
			userID,
			domain.ReconstructPriority(1),
		),
		updateErr: updateErr,
	}
	usecase := NewUpdateTodoUsecase(repositoryMock)

	err := usecase.Execute(InputUpdateTodoUsecase{
		Id:       todoID.String(),
		Title:    "Buy bread",
		Status:   "IN_PROGRESS",
		UserId:   userID.String(),
		priority: 5,
	})
	if !errors.Is(err, updateErr) {
		t.Fatalf("Execute() error = %v, want %v", err, updateErr)
	}
}
