package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("user not found")

type Repository struct {
	db *pgxpool.Pool
}

func (r *Repository) FindOrCreateGithubID(ctx context.Context, githubId string, login string, email string) (int64, error) {
	var id int64
	err := r.db.QueryRow(ctx, `
        INSERT INTO users (username, github_id, email)
        VALUES ($1, $2, $3)
        ON CONFLICT (github_id) DO UPDATE SET github_id = EXCLUDED.github_id
        RETURNING id`,
		login, githubId, email,
	).Scan(&id)

	if err != nil {
		return id, err
	}
	return id, nil
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
	var setClauses []string
	var args []interface{}
	argIdx := 1

	if u.Username != "" {
		setClauses = append(setClauses, fmt.Sprintf("username=$%d", argIdx))
		args = append(args, u.Username)
		argIdx++
	}

	if u.Bio != "" {
		setClauses = append(setClauses, fmt.Sprintf("bio=$%d", argIdx))
		args = append(args, u.Bio)
		argIdx++
	}

	if len(setClauses) == 0 {
		return nil
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE users SET %s WHERE id=$%d",
		strings.Join(setClauses, ", "), argIdx)

	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
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
