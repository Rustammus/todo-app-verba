CREATE TABLE public.tasks
(
    id   SERIAL PRIMARY KEY ,
    title   TEXT NOT NULL ,
    description   TEXT NOT NULL ,
    due_date timestamptz NOT NULL,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP(0) NOT NULL,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP(0) NOT NULL
);