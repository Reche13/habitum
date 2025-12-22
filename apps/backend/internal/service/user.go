package service

import (
	"github.com/labstack/echo/v4"
	"github.com/reche13/habitum/internal/errs"
	"github.com/reche13/habitum/internal/model"
	"github.com/reche13/habitum/internal/model/user"
	"github.com/reche13/habitum/internal/repository"
	"github.com/reche13/habitum/internal/server"
)

type UserService struct {
	server  *server.Server
	userRepo *repository.UserRepository
}

func NewUserService(
	server *server.Server,
	userRepo *repository.UserRepository,
) *UserService {
	return &UserService{
		server:   server,
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(
	ctx echo.Context,
	payload *user.CreateUserPayload,
) (*user.User, error) {
	if payload.Name == "" {
		s.server.Logger.Warn().Msg("user creation failed: name is empty")
		return nil, errs.NewBadRequestError(
			"name is required",
			false, nil, nil, nil,
		)
	}

	if payload.Email == "" {
		s.server.Logger.Warn().Msg("user creation failed: email is empty")
		return nil, errs.NewBadRequestError(
			"email is required",
			false, nil, nil, nil,
		)
	}

	createdUser, err := s.userRepo.Create(
		ctx.Request().Context(),
		payload,
	)
	if err != nil {
		s.server.Logger.Error().Err(err).Msg("failed to create user")
		return nil, err
	}

	s.server.Logger.Info().
		Str("event", "user_created").
		Str("user_id", createdUser.ID.String()).
		Str("email", createdUser.Email).
		Msg("User created successfully")

	return createdUser, nil
}


func (s *UserService) GetUsers(
	ctx echo.Context,
) (*model.PaginatedResponse[user.User], error) {

	users, err := s.userRepo.List(ctx.Request().Context())
	if err != nil {
		s.server.Logger.Error().Err(err).Msg("failed to fetch users")
		return nil, err
	}

	return users, nil
}
