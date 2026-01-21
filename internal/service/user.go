package service

import (
	"Lolopenza/CRUD-F/internal/models"
	"Lolopenza/CRUD-F/internal/repository"
	"context"
	"fmt"
)

type UserService struct {
	Repo *repository.Repo
}

func NewUserService(repo *repository.Repo) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, email, name, surname string) (int, error) {

	if email == "" {
		return 0, fmt.Errorf("email is required")
	}

	id, err := s.Repo.CreateUser(ctx, email, name, surname)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	return id, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.Repo.GetAllUsers(ctx)
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (models.User, error) {
	if id < 1 {
		return models.User{}, fmt.Errorf("invalid id")
	}
	return s.Repo.GetUserByID(ctx, id)
}

func (s *UserService) UpdateUser(ctx context.Context, id int, email, name, surname string) (models.User, error) {
	if id < 1 {
		return models.User{}, fmt.Errorf("invalid id")
	}

	if email == "" {
		return models.User{}, fmt.Errorf("email is required")
	}

	user, err := s.Repo.UpdateUser(ctx, id, email, name, surname)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to update user: %w", err)
	}
	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	if id < 1 {
		return fmt.Errorf("invalid id")
	}
	if err := s.Repo.DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
