package taskService

import (
	"ToDoVerba/internal/dto"
	"ToDoVerba/internal/repos"
	"ToDoVerba/pkg/logging"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"time"
)

type Deps struct {
	Repo   repos.TaskRepository
	Logger logging.Logger
}

type TaskService struct {
	repo   repos.TaskRepository
	logger logging.Logger
}

func (s *TaskService) Create(cTask *dto.TaskCreate) (*dto.TaskRead, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rTask, err := s.repo.Create(ctx, cTask)
	if err != nil {
		s.logger.Errorf("service error on create task: %s", err)
		return nil, err
	}
	s.logger.Debugf("service task created: %+v", rTask)
	return rTask, nil
}

func (s *TaskService) FindByID(id int) (*dto.TaskRead, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rTask, err := s.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.Debugf("No rows found with task id %d", id)
		} else {
			s.logger.Errorf("service error on find task with id %d : %s", id, err)
		}
		return nil, err
	}

	s.logger.Debugf("service task found: %+v", rTask)
	return rTask, nil
}

func (s *TaskService) List() ([]dto.TaskRead, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rTasks, err := s.repo.List(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.Debug("No rows found on list task")
		} else {
			s.logger.Errorf("service error on list task: %s", err)
		}
		return nil, err
	}

	s.logger.Debugf("service found %d tasks", len(rTasks))
	return rTasks, nil
}

func (s *TaskService) UpdateById(id int, update *dto.TaskUpdate) (*dto.TaskRead, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rTask, err := s.repo.UpdateByID(ctx, id, update)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.Debugf("no rows found with task id %d", id)
		} else {
			s.logger.Errorf("service error on list task: %s", err)
		}
		return nil, err
	}

	s.logger.Debugf("service task updated: %+v", rTask)
	return rTask, nil
}

func (s *TaskService) DeleteById(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.repo.DeleteByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.Debugf("no rows found with task id %d", id)
		} else {
			s.logger.Errorf("service error on delete task: %s", err)
		}
	}

	s.logger.Debugf("service task deleted: %d", id)
	return err
}

func NewTaskService(d Deps) *TaskService {
	return &TaskService{
		repo:   d.Repo,
		logger: d.Logger,
	}
}
