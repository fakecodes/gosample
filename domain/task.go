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
	Create(context.Context, *Task) error
}

// ArticleRepository represent the article's repository contract
type TaskRepository interface {
	Create(ctx context.Context, a *Task) error
}
