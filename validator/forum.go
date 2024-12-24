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

	forumNameExists, err := store.ForumNameExists(request.Name)
	if err != nil {
		return err
	}
	if forumNameExists {
		return errors.New("Forum name already exists")
	}

	return nil
}

func ValidateForumModification(store *storage.Storage, forumId int64, request model.ForumRequest) error {
	if checkStringIsEmpty(request.Name) {
		return errors.New("Forum name can't be empty")
	}

	anotherForumNameExists, err := store.AnotherForumExists(forumId, request.Name)
	if err != nil {
		return err
	}
	if anotherForumNameExists {
		return errors.New("A forum with that name already exists")
	}

	return nil
}
