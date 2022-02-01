package validator

import (
	"errors"
	"vpub/model"
	"vpub/storage"
)

func ValidateForumCreation(store *storage.Storage, request model.ForumRequest) error {
	if checkStringIsEmpty(request.Name) {
		return errors.New("Forum name can't be empty")
	}

	if store.ForumNameExists(request.Name) {
		return errors.New("Forum name already exists")
	}

	return nil
}

func ValidateForumModification(store *storage.Storage, forumId int64, request model.ForumRequest) error {
	if checkStringIsEmpty(request.Name) {
		return errors.New("Forum name can't be empty")
	}

	if store.AnotherForumExists(forumId, request.Name) {
		return errors.New("A forum with that name already exists")
	}

	return nil
}
