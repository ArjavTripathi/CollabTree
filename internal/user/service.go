package user

import "context"

type globalRepository interface {
	GetByID(ctx context.Context, id int64) (*User, error)
	Create(ctx context.Context, u *User) (int64, error)
	Update(ctx context.Context, u *User, id int64) error
	Delete(ctx context.Context, id int64) error
	FindOrCreateGithubID(ctx context.Context, id string, login string, email string) (int64, error)
}

type Service struct {
	repo globalRepository
}

func (s *Service) FindOrCreateByGitHubID(ctx context.Context, githubID string, username string, email string) (int64, error) {
	return s.repo.FindOrCreateGithubID(ctx, githubID, username, email)
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetProfile(ctx context.Context, id int64) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) UpdateProfile(ctx context.Context, u *User, id int64) error {
	return s.repo.Update(ctx, u, id)
}

func (s *Service) CreateUser(ctx context.Context, u *User) (int64, error) {
	return s.repo.Create(ctx, u)

}

func (s *Service) DeleteUser(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
