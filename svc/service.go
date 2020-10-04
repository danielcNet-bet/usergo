package user

import (
	"context"

	om "github.com/thelotter-enterprise/usergo/shared"
)

// Service API
type Service interface {
	GetUserByID(ctx context.Context, userID int) (om.User, error)
}

type service struct {
	repo Repository
}

// NewService creates a new instance of service
// service is where we define all the business logic.
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// GetUserByID will execute business logic for getting user information by id
func (s *service) GetUserByID(ctx context.Context, userID int) (om.User, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	return user, err
}