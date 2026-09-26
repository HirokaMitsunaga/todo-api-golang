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

// 名前を変更した新しいUserを返すようにしてイミュータブルになるようにしている
func (u *User) ChangeName(name string) (*User, error) {
	// User{} のようにゼロ値で生成された User は不変条件を満たさないため、
	// 利用時にエラーとして扱う。
	if u == nil || u.id == (ulid.ULID{}) {
		return nil, errors.New("invalid user: create with NewUser")
	}

	updated := *u
	updated.name = name
	return &updated, nil
}
