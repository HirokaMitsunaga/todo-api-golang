package domain

import (
	"errors"
	"testing"

	"github.com/oklog/ulid/v2"
)

func TestTodo_Update(t *testing.T) {
	id := ulid.Make()
	userID := ulid.Make()
	todo := ReconstructTodo(
		id,
		ReconstructTitle("Buy milk"),
		pendingStatus{},
		userID,
		ReconstructPriority(1),
	)

	updated, err := todo.Update(
		ReconstructTitle("Buy bread"),
		inProgressStatus{},
		userID,
		ReconstructPriority(5),
	)
	if err != nil {
		t.Fatalf("Update() failed: %v", err)
	}

	if updated == todo {
		t.Fatal("Update() returned the original Todo")
	}
	if updated.ID() != id {
		t.Errorf("updated ID = %v, want %v", updated.ID(), id)
	}
	if updated.Title().Value() != "Buy bread" {
		t.Errorf("updated title = %q, want %q", updated.Title().Value(), "Buy bread")
	}
	if updated.Status().Name() != "IN_PROGRESS" {
		t.Errorf("updated status = %q, want %q", updated.Status().Name(), "IN_PROGRESS")
	}
	if updated.UserID() != userID {
		t.Errorf("updated user ID = %v, want %v", updated.UserID(), userID)
	}
	if updated.Priority().Value() != 5 {
		t.Errorf("updated priority = %d, want %d", updated.Priority().Value(), 5)
	}

	if todo.Title().Value() != "Buy milk" {
		t.Errorf("original title = %q, want %q", todo.Title().Value(), "Buy milk")
	}
	if todo.Status().Name() != "PENDING" {
		t.Errorf("original status = %q, want %q", todo.Status().Name(), "PENDING")
	}
	if todo.Priority().Value() != 1 {
		t.Errorf("original priority = %d, want %d", todo.Priority().Value(), 1)
	}
}

func TestTodo_Update_invalidStatusTransition(t *testing.T) {
	userId := ulid.Make()
	todo := ReconstructTodo(
		ulid.Make(),
		ReconstructTitle("Buy milk"),
		pendingStatus{},
		userId,
		ReconstructPriority(1),
	)

	updated, err := todo.Update(
		ReconstructTitle("Buy bread"),
		completedStatus{},
		userId,
		ReconstructPriority(5),
	)
	if err == nil {
		t.Fatal("Update() succeeded unexpectedly")
	}
	if updated != nil {
		t.Errorf("Update() = %v, want nil on invalid transition", updated)
	}
	if !errors.Is(err, ErrInvalidTodoStatusTransition) {
		t.Errorf("Update() error = %v, want ErrInvalidTodoStatusTransition", err)
	}
	if todo.Status().Name() != "PENDING" {
		t.Errorf("original status = %q, want %q", todo.Status().Name(), "PENDING")
	}
}
