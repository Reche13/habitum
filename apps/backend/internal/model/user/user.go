package user

import "github.com/reche13/habitum/internal/model"

type User struct {
	model.Base

	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
}
