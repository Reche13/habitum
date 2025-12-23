package service

import (
	"github.com/reche13/habitum/internal/repository"
	"github.com/reche13/habitum/internal/sqlerr"
)

type BaseService struct {
	resourceName string
}

func (b *BaseService) wrapError(err error) error {
	return sqlerr.WrapError(err, b.resourceName)
}

type Services struct{
	User *UserService
	Habit *HabitService
}

func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		User: NewUserService(repos.User),
		Habit: NewHabitService(repos.Habit),
	}
}