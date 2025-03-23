-- Use the beanbag database
\c beanbag

-- USERS TABLE
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME
);

-- QUIZZES TABLE
CREATE TABLE quizzes (
    quiz_id SERIAL PRIMARY KEY,
    creator_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
    quiz_title TEXT NOT NULL,
    description TEXT,
    is_priv BOOLEAN NOT NULL DEFAULT true,
    timer INTEGER NOT NULL DEFAULT 30,
    created_at DATETIME,
    updated_at DATETIME
);

-- QUESTIONS TABLE (Depends on QUIZZES)
CREATE TABLE questions (
    ques_id SERIAL PRIMARY KEY,
    quiz_id INTEGER REFERENCES quizzes(quiz_id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    timer_option BOOLEAN NOT NULL DEFAULT false,
    timer INTEGER NOT NULL DEFAULT 30,
    created_at DATETIME,
    updated_at DATETIME
);

-- ANWSERS TABLE (Depends on QUESTIONS)
CREATE TABLE answers (
    ans_id SERIAL PRIMARY KEY,
    ques_id INTEGER REFERENCES questions(ques_id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    is_correct BOOLEAN NOT NULL,
    created_at DATETIME,
    update_at DATETIME,
);

ALTER TABLE quizzes OWNER TO postgres;
ALTER TABLE questions OWNER TO postgres;
ALTER TABLE answers OWNER TO postgres;

-- Ensure privileges for future tables
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON TABLES TO postgres;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON SEQUENCES TO postgres;

-- Explicitly grant privileges to the user
grant all privileges on all tables in schema public to postgres;
grant all privileges on all sequences in schema public to postgres;
