package dto

import "github.com/jackc/pgx/v5/pgtype"

type TaskCreate struct {
	Title       string
	Description string
	DueDate     pgtype.Timestamptz
}

type TaskRead struct {
	Id          int
	Title       string
	Description string
	DueDate     pgtype.Timestamptz
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}

type TaskUpdate struct {
	Title       string
	Description string
	DueDate     pgtype.Timestamptz
}
