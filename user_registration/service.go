package userregistration

import (
	"github.com/recipey/api"
)

// Service exposes API for this bounded context
type Service interface {
	RegisterUser(user api.User) error
}

type service struct {
	users api.UserRepository
}

// RegisterUser registers user in storage
func (s service) RegisterUser(username, email string) error {
	user := api.NewUser(username, email)
	return s.users.Store(user)
}
