package domain

import (
	"errors"
	"testing"
)

func TestNewTodoStatus(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "creates PENDING",
			input: "PENDING",
			want:  "PENDING",
		},
		{
			name:  "creates IN_PROGRESS",
			input: "IN_PROGRESS",
			want:  "IN_PROGRESS",
		},
		{
			name:  "creates COMPLETED",
			input: "COMPLETED",
			want:  "COMPLETED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTodoStatus(tt.input)
			if err != nil {
				t.Fatalf("NewTodoStatus() failed: %v", err)
			}
			if got.Name() != tt.want {
				t.Errorf("NewTodoStatus().Name() = %q, want %q", got.Name(), tt.want)
			}
		})
	}
}

func TestNewTodoStatus_invalid(t *testing.T) {
	got, err := NewTodoStatus("UNKNOWN")
	if err == nil {
		t.Fatal("NewTodoStatus() succeeded unexpectedly")
	}
	if !errors.Is(err, ErrInvalidTodoStatus) {
		t.Errorf("NewTodoStatus() error = %v, want ErrInvalidTodoStatus", err)
	}
	if got != nil {
		t.Errorf("NewTodoStatus() = %v, want nil", got)
	}
}

func TestTodoStatus_transitionTo(t *testing.T) {
	tests := []struct {
		name        string
		current     TodoStatus
		next        TodoStatus
		want        TodoStatus
		wantErr     bool
		wantMessage string
	}{
		{
			name:    "PENDING can transition to IN_PROGRESS",
			current: pendingStatus{},
			next:    inProgressStatus{},
			want:    inProgressStatus{},
		},
		{
			name:        "PENDING cannot transition to COMPLETED",
			current:     pendingStatus{},
			next:        completedStatus{},
			wantErr:     true,
			wantMessage: "invalid todo status transition: PENDING -> COMPLETED",
		},
		{
			name:    "IN_PROGRESS can transition to COMPLETED",
			current: inProgressStatus{},
			next:    completedStatus{},
			want:    completedStatus{},
		},
		{
			name:        "IN_PROGRESS cannot transition to PENDING",
			current:     inProgressStatus{},
			next:        pendingStatus{},
			wantErr:     true,
			wantMessage: "invalid todo status transition: IN_PROGRESS -> PENDING",
		},
		{
			name:        "COMPLETED cannot transition to PENDING",
			current:     completedStatus{},
			next:        pendingStatus{},
			wantErr:     true,
			wantMessage: "invalid todo status transition: COMPLETED -> PENDING",
		},
		{
			name:        "nil is not a valid next status",
			current:     pendingStatus{},
			wantErr:     true,
			wantMessage: "invalid todo status transition: PENDING -> <nil>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.current.transitionTo(tt.next)
			if tt.wantErr {
				if err == nil {
					t.Fatal("transitionTo() succeeded unexpectedly")
				}
				if !errors.Is(err, ErrInvalidTodoStatusTransition) {
					t.Errorf("transitionTo() error = %v, want ErrInvalidTodoStatusTransition", err)
				}
				if err.Error() != tt.wantMessage {
					t.Errorf("transitionTo() error = %q, want %q", err.Error(), tt.wantMessage)
				}
				return
			}

			if err != nil {
				t.Fatalf("transitionTo() failed: %v", err)
			}
			if got != tt.want {
				t.Errorf("transitionTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTodo_updateStatus(t *testing.T) {
	todo := &Todo{status: pendingStatus{}}

	updated, err := todo.updateStatus(inProgressStatus{})
	if err != nil {
		t.Fatalf("updateStatus() failed: %v", err)
	}
	if updated.status != (inProgressStatus{}) {
		t.Errorf("updateStatus() status = %v, want IN_PROGRESS", updated.status.Name())
	}
	if todo.status != (pendingStatus{}) {
		t.Errorf("updateStatus() mutated the receiver: got %v, want PENDING", todo.status.Name())
	}

	updated, err = todo.updateStatus(completedStatus{})
	if err == nil {
		t.Fatal("updateStatus() succeeded unexpectedly")
	}
	if updated != nil {
		t.Errorf("updateStatus() = %v, want nil on invalid transition", updated)
	}
	if !errors.Is(err, ErrInvalidTodoStatusTransition) {
		t.Errorf("updateStatus() error = %v, want ErrInvalidTodoStatusTransition", err)
	}
}
