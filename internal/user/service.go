package user

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetProfile(ctx context.Context, id int64) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) UpdateProfile(ctx context.Context, u *User) error {
	return s.repo.Update(ctx, u)
}

func (s *Service) CreateUser(ctx context.Context, u *User) (int64, error) {
	return s.repo.Create(ctx, u)
}
