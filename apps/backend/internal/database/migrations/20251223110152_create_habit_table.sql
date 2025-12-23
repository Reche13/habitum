-- +goose Up
-- +goose StatementBegin
CREATE TYPE habit_frequency AS ENUM ('daily', 'weekly');

CREATE TABLE habits (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    name TEXT NOT NULL,
    description TEXT,
    icon TEXT,
    color TEXT,
    category TEXT NOT NULL,

    frequency habit_frequency NOT NULL,
    times_per_week INT,

    archived_at TIMESTAMPTZ NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT habit_frequency_check CHECK (
        (frequency = 'daily' AND times_per_week IS NULL)
        OR
        (frequency = 'weekly' AND times_per_week BETWEEN 1 AND 7)
    )
);

CREATE INDEX idx_habits_user_id ON habits(user_id);

CREATE INDEX idx_habits_user_category ON habits(user_id, category);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS habits;
-- +goose StatementEnd
