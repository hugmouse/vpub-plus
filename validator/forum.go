package validator

import (
	"errors"
	"strings"
	"vpub/model"
	"vpub/storage"
)

func ValidateForumCreation(store *storage.Storage, forum model.Forum) error {
	if len(strings.TrimSpace(forum.Name)) == 0 {
		return errors.New("Forum name can't be empty")
	}

	if store.ForumNameExists(forum.Name) {
		return errors.New("A forum with that name already exists")
	}

	return nil
}

func ValidateForumModification(store *storage.Storage, forum model.Forum) error {
	if len(strings.TrimSpace(forum.Name)) == 0 {
		return errors.New("Forum name can't be empty")
	}

	if store.AnotherForumExists(forum.Id, forum.Name) {
		return errors.New("A forum with that name already exists")
	}

	return nil
}
