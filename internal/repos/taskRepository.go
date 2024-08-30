package repos

import (
	"ToDoVerba/internal/dto"
	"context"
)

type TaskRepository interface {
	Create(ctx context.Context, cTask *dto.TaskCreate) (*dto.TaskRead, error)
	FindById(ctx context.Context, id int) (*dto.TaskRead, error)
	List(ctx context.Context) ([]dto.TaskRead, error)
	UpdateByID(ctx context.Context, id int, update *dto.TaskUpdate) (*dto.TaskRead, error)
	DeleteByID(ctx context.Context, id int) (int, error)
}
