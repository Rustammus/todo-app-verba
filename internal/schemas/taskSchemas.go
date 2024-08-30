package schemas

import (
	"ToDoVerba/internal/dto"
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type RequestTaskCreate struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
}

func (t *RequestTaskCreate) ToDTO() *dto.TaskCreate {
	parsedTime, _ := time.Parse(time.RFC3339, t.DueDate)
	dueDate := pgtype.Timestamptz{
		Time:             parsedTime,
		InfinityModifier: 0,
		Valid:            true,
	}

	return &dto.TaskCreate{
		Title:       t.Title,
		Description: t.Description,
		DueDate:     dueDate,
	}
}

func (t *RequestTaskCreate) Valid() error {
	errStr := ""
	if t.Title == "" {
		errStr += "Title is required;"
	}
	if t.Description == "" {
		errStr += "Description is required;"
	}
	if _, err := time.Parse(time.RFC3339, t.DueDate); err != nil {
		errStr += "DueDate is required and must be in RFC3339 format;"
	}
	if len(errStr) > 0 {
		return errors.New(errStr)
	}
	return nil
}

type RequestTaskUpdate struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
}

func (t *RequestTaskUpdate) ToDTO() *dto.TaskUpdate {
	parsedTime, _ := time.Parse(time.RFC3339, t.DueDate)
	dueDate := pgtype.Timestamptz{
		Time:             parsedTime,
		InfinityModifier: 0,
		Valid:            true,
	}

	return &dto.TaskUpdate{
		Title:       t.Title,
		Description: t.Description,
		DueDate:     dueDate,
	}
}

func (t *RequestTaskUpdate) Valid() error {
	errStr := ""
	if t.Title == "" {
		errStr += "Title is required;"
	}
	if t.Description == "" {
		errStr += "Description is required;"
	}
	if _, err := time.Parse(time.RFC3339, t.DueDate); err != nil {
		errStr += "DueDate is required and must be in RFC3339 format;"
	}
	if len(errStr) > 0 {
		return errors.New(errStr)
	}
	return nil
}

type ResponseTaskRead struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func (t *ResponseTaskRead) ScanDTO(task *dto.TaskRead) {
	//TODO out validation
	t.Id = task.Id
	t.Title = task.Title
	t.Description = task.Description
	t.DueDate = task.DueDate.Time.Format(time.RFC3339)
	t.CreatedAt = task.CreatedAt.Time.Format(time.RFC3339)
	t.UpdatedAt = task.UpdatedAt.Time.Format(time.RFC3339)
}
