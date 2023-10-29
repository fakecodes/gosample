package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Task is representing the role data struct
type Task struct {
	ID          primitive.ObjectID `bson:"id" json:"id"`
	Name        string             `bson:"name" json:"name" validate:"required"`
	Description string             `bson:"description" json:"description"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	DueDate     time.Time          `bson:"due_date" json:"due_date"`
	Priority    string             `bson:"priority" json:"priority"`
	Completed   bool               `bson:"completed" json:"completed"`
}

// TaskUsecase represent the role's usecases
type TaskUsecase interface {
	Create(context.Context, *Task) error
}

// ArticleRepository represent the article's repository contract
type TaskRepository interface {
	Create(ctx context.Context, a *Task) error
}
