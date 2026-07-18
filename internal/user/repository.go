package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("user not found")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, u *User) (int64, error) {
	var id int64
	err := r.db.QueryRow(ctx, `INSERT INTO users (username, email, github_id, bio)
         VALUES ($1, $2, $3, $4) RETURNING id`, u.Username, u.Email, u.GitHubID, u.Bio).Scan(&id)
	return id, err
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*User, error) {
	var u User
	err := r.db.QueryRow(ctx,
		`SELECT id, username, email, github_id, bio, created_at
         FROM users WHERE id=$1`, id,
	).Scan(&u.ID, &u.Username, &u.Email, &u.GitHubID, &u.Bio, &u.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) Update(ctx context.Context, u *User, id int64) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE users SET username=$1, bio=$2 WHERE id=$3`,
		u.Username, u.Bio, id,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	err := r.db.QueryRow(ctx, `DELETE FROM users WHERE id=$1`, id)

	if err != nil {
		return pgx.ErrNoRows
	}
	return nil
}
