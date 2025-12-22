package service

import (
	"github.com/reche13/habitum/internal/repository"
)


type Services struct{
	User *UserService
}

func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		User: NewUserService(repos.User),
	}
}