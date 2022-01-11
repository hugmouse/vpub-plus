package model

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

type User struct {
	Id       int64
	Name     string
	Password string
	Hash     []byte
	About    string
	IsAdmin  bool
}

func (u User) Validate() error {
	if u.Name == "" {
		return errors.New("username is mandatory")
	}
	if u.Password == "" {
		return errors.New("password is mandatory")
	}
	match, _ := regexp.MatchString("^[a-z0-9-_]+$", u.Name)
	if !match {
		return errors.New("username should match [a-z0-9-_]")
	}
	return nil
}

func (u User) HashPassword() ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
}

func (u User) CompareHashToPassword(hash []byte) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(u.Password))
}
