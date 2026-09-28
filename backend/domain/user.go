package domain

import (
	"errors"

	"github.com/oklog/ulid/v2"
)

// フィールドは非公開で、domainパッケージ外から直接読み書きできない(小文字はじまりのため)
type User struct {
	id       ulid.ULID
	name     string
	email    string
	password string
}

func NewUser(name string, email string, password string) (*User, error) {
	return &User{
		id:       ulid.Make(),
		name:     name,
		email:    email,
		password: password,
	}, nil
}

func Reconstructor(id ulid.ULID, name string, email string, password string) (*User, error) {
	return &User{
		id:       id,
		name:     name,
		email:    email,
		password: password,
	}, nil
}

// WithName は名前を変更した新しい User を返す。
// 元の User は変更せず、イミュータブルに扱えるようにする。
func (u *User) WithName(name string) (*User, error) {
	// User{} のようにゼロ値で生成された User は不変条件を満たさないため、
	// 利用時にエラーとして扱う。
	if u == nil || u.id == (ulid.ULID{}) {
		return nil, errors.New("invalid user: create with NewUser")
	}

	updated := *u
	updated.name = name
	return &updated, nil
}

// WithEmail はメールアドレスを変更した新しい User を返す。
func (u *User) WithEmail(email string) (*User, error) {
	if u == nil || u.id == (ulid.ULID{}) {
		return nil, errors.New("invalid user: create with NewUser")
	}

	updated := *u
	updated.email = email
	return &updated, nil
}

// WithPassword はパスワードを変更した新しい User を返す。
func (u *User) WithPassword(password string) (*User, error) {
	if u == nil || u.id == (ulid.ULID{}) {
		return nil, errors.New("invalid user: create with NewUser")
	}

	updated := *u
	updated.password = password
	return &updated, nil
}

func (u *User) ID() ulid.ULID {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Password() string {
	return u.password
}
