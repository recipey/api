package api

import (
	"errors"
)

// User is the user domain model.
type User struct {
	ID       string
	Username string
	Email    string
}

// UserRepository provides access to a user store.
type UserRepository interface {
	Store(user *User) error
	Find(id string) (*User, error)
	FindAll() ([]*User, error)
}

// ErrUnknownUser is used when a user could not be found.
var ErrUnknownUser = errors.New("unknown user")

// NewUser creates a new user.
func NewUser(username, email string) *User {
	return &User{
		Username: username,
		Email:    email,
	}
}
