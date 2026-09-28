package usecase

import (
	"errors"
	"testing"

	"todo-api/domain"

	"github.com/oklog/ulid/v2"
)

type fakeUserRepository struct {
	user        *domain.User
	findByID    ulid.ULID
	updatedUser *domain.User
	findErr     error
	updateErr   error
}

func (f *fakeUserRepository) FindById(id ulid.ULID) (*domain.User, error) {
	f.findByID = id
	if f.findErr != nil {
		return nil, f.findErr
	}
	return f.user, nil
}

func (f *fakeUserRepository) Create(*domain.User) error {
	return nil
}

func (f *fakeUserRepository) Update(user *domain.User) error {
	f.updatedUser = user
	return f.updateErr
}

func (f *fakeUserRepository) Delete(ulid.ULID) error {
	return nil
}

func TestUpdateUserUsecase_Execute(t *testing.T) {
	userID := ulid.Make()
	user, err := domain.Reconstructor(userID, "old name", "old@example.com", "old-password")
	if err != nil {
		t.Fatalf("could not construct user: %v", err)
	}
	repository := &fakeUserRepository{user: user}
	usecase := NewUpdateUserUsecase(repository)

	name := "new name"
	email := "new@example.com"
	password := "new-password"
	err = usecase.Execute(UpdateUserInput{
		ID:       userID.String(),
		Name:     name,
		Email:    email,
		Password: password,
	})
	if err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}

	if repository.findByID != userID {
		t.Fatalf("FindById() received %v, want %v", repository.findByID, userID)
	}
	if got := repository.updatedUser.Name(); got != name {
		t.Errorf("updated name = %q, want %q", got, name)
	}
	if got := repository.updatedUser.Email(); got != email {
		t.Errorf("updated email = %q, want %q", got, email)
	}
	if got := repository.updatedUser.Password(); got != password {
		t.Errorf("updated password = %q, want %q", got, password)
	}
	if got := user.Name(); got != "old name" {
		t.Errorf("original name = %q, want %q", got, "old name")
	}
}

func TestUpdateUserUsecase_Execute_updatesAllFields(t *testing.T) {
	userID := ulid.Make()
	user, err := domain.Reconstructor(userID, "old name", "old@example.com", "old-password")
	if err != nil {
		t.Fatalf("could not construct user: %v", err)
	}
	repository := &fakeUserRepository{user: user}
	usecase := NewUpdateUserUsecase(repository)

	name := "new name"
	email := "new@example.com"
	password := "new-password"
	err = usecase.Execute(UpdateUserInput{
		ID:       userID.String(),
		Name:     name,
		Email:    email,
		Password: password,
	})
	if err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}

	if got := repository.updatedUser.Name(); got != name {
		t.Errorf("updated name = %q, want %q", got, name)
	}
	if got := repository.updatedUser.Email(); got != email {
		t.Errorf("updated email = %q, want %q", got, email)
	}
	if got := repository.updatedUser.Password(); got != password {
		t.Errorf("updated password = %q, want %q", got, password)
	}
}

func TestUpdateUserUsecase_Execute_invalidID(t *testing.T) {
	repository := &fakeUserRepository{}
	usecase := NewUpdateUserUsecase(repository)

	if err := usecase.Execute(UpdateUserInput{ID: "invalid"}); err == nil {
		t.Fatal("Execute() succeeded unexpectedly")
	}
	if repository.findByID != (ulid.ULID{}) {
		t.Fatal("FindById() was called for an invalid ID")
	}
}

func TestUpdateUserUsecase_Execute_propagatesFindError(t *testing.T) {
	findErr := errors.New("find failed")
	repository := &fakeUserRepository{findErr: findErr}
	usecase := NewUpdateUserUsecase(repository)

	err := usecase.Execute(UpdateUserInput{ID: ulid.Make().String()})
	if !errors.Is(err, findErr) {
		t.Fatalf("Execute() error = %v, want %v", err, findErr)
	}
}
