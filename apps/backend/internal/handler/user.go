package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/reche13/habitum/internal/errs"
	"github.com/reche13/habitum/internal/middleware"
	"github.com/reche13/habitum/internal/model"
	"github.com/reche13/habitum/internal/model/user"
	"github.com/reche13/habitum/internal/service"
	"github.com/rs/zerolog"
)

type UserHandler struct {
	logger      zerolog.Logger
	userService *service.UserService
}

func NewUserHandler(
	userService *service.UserService,
) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var payload user.CreateUserPayload

	if err := c.Bind(&payload); err != nil {
		return errs.NewBadRequestError("Invalid request payload")
	}

	if fieldErrors := middleware.ValidateStruct(&payload); fieldErrors != nil {
		return errs.NewValidationError(fieldErrors)
	}

	createdUser, err := h.userService.CreateUser(c.Request().Context(), &payload)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, model.SuccessResponse(createdUser))
}

func (h *UserHandler) GetUser(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return errs.NewBadRequestError("Invalid user ID format")
	}

	u, err := h.userService.GetUser(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, model.SuccessResponse(u))
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	users, err := h.userService.GetUsers(c.Request().Context())
	if err != nil {
		return err
	}

	meta := &model.Meta{
		RequestID: middleware.GetRequestID(c),
		Total:     len(users),
	}

	return c.JSON(http.StatusOK, model.SuccessResponseWithMeta(users, meta))
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return errs.NewBadRequestError("Invalid user ID format")
	}

	var payload user.UpdateUserPayload
	if err := c.Bind(&payload); err != nil {
		return errs.NewBadRequestError("Invalid request payload")
	}

	if fieldErrors := middleware.ValidateStruct(&payload); fieldErrors != nil {
		return errs.NewValidationError(fieldErrors)
	}

	updatedUser, err := h.userService.UpdateUser(c.Request().Context(), id, &payload)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, model.SuccessResponse(updatedUser))
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return errs.NewBadRequestError("Invalid user ID format")
	}

	if err := h.userService.DeleteUser(c.Request().Context(), id); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
