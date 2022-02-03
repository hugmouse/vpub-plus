package validator

import (
	"errors"
	"regexp"
	"vpub/model"
	"vpub/storage"
)

func ValidateUserCreation(store *storage.Storage, key string, r model.UserCreationRequest) error {
	if len(r.Name) < 3 {
		return errors.New("Username needs to be at least 3 characters")
	}

	if len(r.Name) > 20 {
		return errors.New("Username should be 20 characters or less")
	}

	if match, _ := regexp.MatchString("^[a-z0-9-_]+$", r.Name); !match {
		return errors.New("Only lowercase letters and digits are accepted for username")
	}

	if store.UserExists(r.Name) {
		return errors.New("Username already exists")
	}

	if !store.KeyExists(key) {
		return errors.New("Key not found or already used")
	}

	return nil
}
