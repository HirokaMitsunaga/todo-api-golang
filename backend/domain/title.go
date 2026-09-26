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

func NewTitle(value string) (*Title, error) {
	if len(value) < min || len(value) > max {
		return &Title{}, fmt.Errorf("３文字以内、１００文字以下で入力してください")
	}
	return &Title{title: value}, nil
}
