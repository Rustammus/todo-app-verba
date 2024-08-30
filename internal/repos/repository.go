package repos

import (
	"ToDoVerba/internal/crud"
	"ToDoVerba/pkg/logging"
)

type Repositories struct {
	Task TaskRepository
}

func NewRepositories(pool crud.Client, logger logging.Logger) Repositories {
	return Repositories{
		Task: crud.NewTaskCRUD(pool, logger),
	}
}
