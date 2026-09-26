package domain

import (
	"fmt"
)

const (
	min = 3
	max = 100
)

type Title struct {
	title string
}

func NewTitle(name string) (*Title, error) {
	if len(name) < min || len(name) > max {
		return &Title{}, fmt.Errorf("３文字以内、１００文字以下で入力してください")
	}
	return &Title{title: name}, nil
}
