package main

import "time"

type Task struct {
	ID          int
	Description string
	CreatedAt   time.Time
	IsComplete  bool
}

func NewTask(id int, description string) Task {
	return Task{
		ID:          id,
		Description: description,
		CreatedAt:   time.Now(),
		IsComplete:  false,
	}
}
