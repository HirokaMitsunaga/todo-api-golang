package domain

import (
	"reflect"
	"testing"
)

func TestUser_chnageName(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		cid       string
		cname     string
		cemail    string
		cpassword string
		// Named input parameters for target function.
		newName string
		want    *User
		wantErr bool
	}{
		{
			name:      "changes the user name",
			cid:       "user-1",
			cname:     "old name",
			cemail:    "user@example.com",
			cpassword: "password",
			newName:   "new name",
			want: &User{
				id:       "user-1",
				name:     "new name",
				email:    "user@example.com",
				password: "password",
			},
		},
		{
			name:      "changes the user name to an empty string",
			cid:       "user-2",
			cname:     "old name",
			cemail:    "another@example.com",
			cpassword: "another-password",
			newName:   "",
			want: &User{
				id:       "user-2",
				name:     "",
				email:    "another@example.com",
				password: "another-password",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := NewUser(tt.cid, tt.cname, tt.cemail, tt.cpassword)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got, gotErr := u.ChangeName(tt.newName)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("chnageName() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("chnageName() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("chnageName() = %v, want %v", got, tt.want)
			}

			wantOriginal := &User{
				id:       tt.cid,
				name:     tt.cname,
				email:    tt.cemail,
				password: tt.cpassword,
			}
			if !reflect.DeepEqual(u, wantOriginal) {
				t.Errorf("ChangeName() mutated the receiver: got %v, want %v", u, wantOriginal)
			}
		})
	}
}

func TestUser_chnageName_zeroValue(t *testing.T) {
	var u User

	got, err := u.ChangeName("new name")
	if err == nil {
		t.Fatal("ChangeName() succeeded unexpectedly for a zero-value User")
	}
	if got != nil {
		t.Errorf("ChangeName() = %v, want nil for a zero-value User", got)
	}
}
