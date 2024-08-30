package model

import "github.com/jackc/pgx/v5/pgtype"

type Task struct {
	Id          int
	Title       string
	Description string
	DueDate     pgtype.Timestamptz
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}
