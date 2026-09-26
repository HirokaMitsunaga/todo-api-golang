package domain

import (
	"reflect"
	"testing"

	"github.com/oklog/ulid/v2"
)

func TestUser_ChangeName(t *testing.T) {
	tests := []struct {
		name     string
		cname    string
		cemail   string
		password string
		newName  string
	}{
		{
			name:     "changes the user name",
			cname:    "old name",
			cemail:   "user@example.com",
			password: "password",
			newName:  "new name",
		},
		{
			name:     "changes the user name to an empty string",
			cname:    "old name",
			cemail:   "another@example.com",
			password: "another-password",
			newName:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := NewUser(tt.cname, tt.cemail, tt.password)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			if u.id == (ulid.ULID{}) {
				t.Fatal("NewUser() generated a zero ULID")
			}

			got, err := u.ChangeName(tt.newName)
			if err != nil {
				t.Fatalf("ChangeName() failed: %v", err)
			}

			want := &User{
				id:       u.id,
				name:     tt.newName,
				email:    tt.cemail,
				password: tt.password,
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("ChangeName() = %v, want %v", got, want)
			}

			wantOriginal := &User{
				id:       u.id,
				name:     tt.cname,
				email:    tt.cemail,
				password: tt.password,
			}
			if !reflect.DeepEqual(u, wantOriginal) {
				t.Errorf("ChangeName() mutated the receiver: got %v, want %v", u, wantOriginal)
			}
		})
	}
}

func TestUser_ChangeName_zeroValue(t *testing.T) {
	var u User

	got, err := u.ChangeName("new name")
	if err == nil {
		t.Fatal("ChangeName() succeeded unexpectedly for a zero-value User")
	}
	if got != nil {
		t.Errorf("ChangeName() = %v, want nil for a zero-value User", got)
	}
}
