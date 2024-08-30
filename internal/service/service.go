package service

import (
	"ToDoVerba/internal/dto"
	"ToDoVerba/internal/repos"
	"ToDoVerba/internal/service/taskService"
	"ToDoVerba/pkg/logging"
)

type Deps struct {
	Repos  repos.Repositories
	Logger logging.Logger
}

type Services struct {
	Task ITaskService
}

func NewServices(d Deps) Services {
	return Services{
		Task: taskService.NewTaskService(taskService.Deps{
			Repo:   d.Repos.Task,
			Logger: d.Logger,
		}),
	}
}

//go:generate mockgen -source=service.go -destination=mocks\mock.go

type ITaskService interface {
	Create(cTask *dto.TaskCreate) (*dto.TaskRead, error)
	FindByID(id int) (*dto.TaskRead, error)
	List() ([]dto.TaskRead, error)
	UpdateById(id int, update *dto.TaskUpdate) (*dto.TaskRead, error)
	DeleteById(id int) error
}
