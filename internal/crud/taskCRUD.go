package crud

import (
	"ToDoVerba/internal/dto"
	"ToDoVerba/pkg/logging"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type TaskCRUD struct {
	client Client
	logger logging.Logger
}

func (c *TaskCRUD) Create(ctx context.Context, cTask *dto.TaskCreate) (*dto.TaskRead, error) {
	q := `INSERT INTO public.tasks (title, description, due_date, created_at, updated_at) 
		  VALUES ($1, $2, $3, $4, $5) 
          RETURNING id, title, description, due_date, created_at, updated_at`

	curTime := pgtype.Timestamptz{
		Time:             time.Now().UTC(),
		InfinityModifier: 0,
		Valid:            true,
	}
	rTask := &dto.TaskRead{}

	err := c.client.QueryRow(ctx, q, cTask.Title, cTask.Description, cTask.DueDate, curTime, curTime).
		Scan(&rTask.Id, &rTask.Title, &rTask.Description, &rTask.DueDate, &rTask.CreatedAt, &rTask.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return rTask, nil
}

func (c *TaskCRUD) FindById(ctx context.Context, id int) (*dto.TaskRead, error) {
	q := `SELECT id, title, description, due_date, created_at, updated_at 
		  FROM public.tasks 
		  WHERE id = $1`

	rTask := &dto.TaskRead{}

	err := c.client.QueryRow(ctx, q, id).
		Scan(&rTask.Id, &rTask.Title, &rTask.Description, &rTask.DueDate, &rTask.CreatedAt, &rTask.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return rTask, nil
}

func (c *TaskCRUD) List(ctx context.Context) ([]dto.TaskRead, error) {
	q := `SELECT id, title, description, due_date, created_at, updated_at 
		  FROM public.tasks`

	rows, err := c.client.Query(ctx, q)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var tasks []dto.TaskRead

	for rows.Next() {
		rTask := dto.TaskRead{}
		err := rows.Scan(&rTask.Id, &rTask.Title, &rTask.Description, &rTask.DueDate, &rTask.CreatedAt, &rTask.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, rTask)
	}
	if len(tasks) == 0 {
		return nil, pgx.ErrNoRows
	}

	return tasks, err
}

func (c *TaskCRUD) UpdateByID(ctx context.Context, id int, update *dto.TaskUpdate) (*dto.TaskRead, error) {
	q := `UPDATE public.tasks 
		  SET (title, description, due_date, updated_at) = ($2, $3, $4, $5) 
		  WHERE id = $1 
		  RETURNING id, title, description, due_date, created_at, updated_at`

	curTime := pgtype.Timestamptz{
		Time:             time.Now().UTC(),
		InfinityModifier: 0,
		Valid:            true,
	}
	rTask := &dto.TaskRead{}

	err := c.client.QueryRow(ctx, q, id, update.Title, update.Description, update.DueDate, curTime).
		Scan(&rTask.Id, &rTask.Title, &rTask.Description, &rTask.DueDate, &rTask.CreatedAt, &rTask.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return rTask, nil
}

func (c *TaskCRUD) DeleteByID(ctx context.Context, id int) (int, error) {
	q := `DELETE FROM public.tasks WHERE id = $1 RETURNING id`

	err := c.client.QueryRow(ctx, q, id).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func NewTaskCRUD(client Client, logger logging.Logger) *TaskCRUD {
	return &TaskCRUD{
		client: client,
		logger: logger,
	}
}
