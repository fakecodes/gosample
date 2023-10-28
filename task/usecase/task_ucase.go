package usecase

import (
	"context"
	"time"

	"github.com/fakecodes/gosample/domain"
)

type taskUsecase struct {
	taskRepo       domain.TaskRepository
	contextTimeout time.Duration
}

// NewTaskUsecase will create new an taskUsecase object representation of domain.TaskUsecase interface
func NewTaskUsecase(a domain.TaskRepository, timeout time.Duration) domain.TaskUsecase {
	return &taskUsecase{
		taskRepo:       a,
		contextTimeout: timeout,
	}
}

func (a *taskUsecase) Create(c context.Context, m *domain.Task) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	err = a.taskRepo.Create(ctx, m)
	return
}
