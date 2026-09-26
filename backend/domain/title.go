package domain

import (
	"fmt"
)

const (
	minTitle = 3
	maxTitle = 100
)

type Title struct {
	title string
}

func NewTitle(value string) (*Title, error) {
	if len(value) < minTitle || len(value) > maxTitle {
		return &Title{}, fmt.Errorf("３文字以内、１００文字以下で入力してください")
	}
	return &Title{title: value}, nil
}
