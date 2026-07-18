CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username TEXT NOT NULL UNIQUE,
                       email TEXT NOT NULL UNIQUE,
                       github_id TEXT UNIQUE,
                       bio TEXT DEFAULT '',
                       created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);