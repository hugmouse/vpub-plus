package model

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"os"
	"path"
	"regexp"
	"strings"
)

const UserDir = "users"

type User struct {
	Name     string
	Password string
	Hash     []byte
	About    string
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

func (u User) Folder() string {
	return path.Join(UserDir, strings.ToLower(u.Name))
}

func (u User) CreateFolder() error {
	if _, err := os.Stat(u.Folder()); !os.IsNotExist(err) {
		return err
	}
	return os.Mkdir(strings.ToLower(u.Folder()), 0700)
}
