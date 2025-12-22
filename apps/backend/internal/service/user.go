package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/reche13/habitum/internal/model/user"
	"github.com/reche13/habitum/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(
	userRepo *repository.UserRepository,
) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}


func (s *UserService) CreateUser(
	ctx context.Context,
	payload *user.CreateUserPayload,
) (*user.User, error) {
	createdUser, err := s.userRepo.Create(ctx, payload)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*user.User, error) {
	u, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *UserService) GetUsers(ctx context.Context) ([]user.User, error) {
	users, err := s.userRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) UpdateUser(
	ctx context.Context,
	id uuid.UUID,
	payload *user.UpdateUserPayload,
) (*user.User, error) {
	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	updatedUser, err := s.userRepo.Update(ctx, id, payload)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
