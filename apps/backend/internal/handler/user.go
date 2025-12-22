package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/reche13/habitum/internal/errs"
	"github.com/reche13/habitum/internal/model/user"
	"github.com/reche13/habitum/internal/server"
	"github.com/reche13/habitum/internal/service"
)

type UserHandler struct {
	server     *server.Server
	userService *service.UserService
}

func NewUserHandler(
	s *server.Server,
	userService *service.UserService,
) *UserHandler {
	return &UserHandler{
		server:      s,
		userService: userService,
	}
}


func (h *UserHandler) CreateUser(c echo.Context) error {
	var payload user.CreateUserPayload

	if err := c.Bind(&payload); err != nil {
		return errs.NewBadRequestError(
			"invalid request payload",
			false, nil, nil, nil,
		)
	}

	createdUser, err := h.userService.CreateUser(c, &payload)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, createdUser)
}


func (h *UserHandler) GetUsers(c echo.Context) error {
	users, err := h.userService.GetUsers(c)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, users)
}
