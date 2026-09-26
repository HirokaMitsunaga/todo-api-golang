package domain

import "fmt"

const (
	minPriority = 1
	maxPriority = 10
)

type Priority struct {
	priority uint16
}

func NewPriority(value uint16) (*Priority, error) {
	if value < minPriority || value > maxPriority {
		return &Priority{}, fmt.Errorf("優先度は1以上10以下で入力してください")
	}
	return &Priority{priority: value}, nil
}
