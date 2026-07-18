package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository struct {
	db *pgxpool.Pool
}

func NewSessionRepository(db *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{}
}

func (r *SessionRepository) Create(ctx context.Context, s Session) error {
	_, err := r.db.Exec(ctx, `INSERT INTO sessions (id, user_id, expires_at) VALUES ($1, $2, $3)`,
		s.ID, s.UserID, s.ExpiresAt)
	return err
}

func (r *SessionRepository) Get(ctx context.Context, id string) (Session, error) {
	var s Session
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, expires_at FROM sessions WHERE id = $1`, id,
	).Scan(&s.ID, &s.UserID, &s.ExpiresAt)
	return s, err
}

func (r *SessionRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, "DELETE FROM sessions WHERE id = $1", id)
	return err
}
