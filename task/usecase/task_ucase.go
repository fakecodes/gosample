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

func (a *taskUsecase) Fetch(c context.Context, cursor string, num int64) (res []domain.Task, nextCursor string, err error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, nextCursor, err = a.taskRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}
	return
}

func (a *taskUsecase) GetByID(c context.Context, id int64) (res domain.Task, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err = a.taskRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	return
}

func (a *taskUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedTask, err := a.taskRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	if existedTask == (domain.Task{}) {
		return domain.ErrNotFound
	}
	return a.taskRepo.Delete(ctx, id)
}
