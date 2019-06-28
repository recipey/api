package userregistration

import (
	"github.com/recipey/api"
)

// Service exposes API for this bounded context
type Service interface {
	RegisterUser(user *api.User) error
}

type service struct {
	users api.UserRepository
}

// RegisterUser registers user in storage
func (s service) RegisterUser(user *api.User) error {
	return s.users.Store(user)
}

// NewUserRegistrationService returns a new user registration service
func NewUserRegistrationService(ur api.UserRepository) Service {
	return &service{
		users: ur,
	}
}
