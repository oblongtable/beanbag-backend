-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS answers (
    ans_id SERIAL PRIMARY KEY,
    ques_id INTEGER REFERENCES questions(ques_id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    is_correct BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS answers;
-- +goose StatementEnd
