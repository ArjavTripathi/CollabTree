CREATE TABLE sessions (
        id text PRIMARY KEY,           -- random token, not sequential
        user_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        expires_at timestamptz NOT NULL
);