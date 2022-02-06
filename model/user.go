package model

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int64
	Name     string
	Password string
	Hash     string
	IsAdmin  bool
	About    string
	Picture  string
}

// UserCreationRequest represents the request to create a user.
type UserCreationRequest struct {
	Name     string
	Password string
	IsAdmin  bool
}

func (u User) CompareHashToPassword(hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(u.Password))
}
