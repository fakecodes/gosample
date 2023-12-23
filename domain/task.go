package domain

import (
	"context"
	"time"
)

// Task is representing the role data struct
type Task struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	DueDate     time.Time `json:"due_date"`
	Priority    string    `json:"priority"`
	Completed   bool      `json:"completed"`
}

// TaskUsecase represent the role's usecases
type TaskUsecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]Task, string, error)
	GetByID(ctx context.Context, id int64) (Task, error)
	Create(context.Context, *Task) error
}

// TaskRepository represent the task's repository contract
type TaskRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []Task, nextCursor string, err error)
	GetByID(ctx context.Context, id int64) (Task, error)
	Create(ctx context.Context, a *Task) error
}
