-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS quizzes (
    quiz_id SERIAL PRIMARY KEY,
    creator_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
    quiz_title TEXT NOT NULL,
    description TEXT,
    is_priv BOOLEAN NOT NULL DEFAULT true,
    timer INTEGER NOT NULL DEFAULT 30,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS quizzes;
-- +goose StatementEnd
