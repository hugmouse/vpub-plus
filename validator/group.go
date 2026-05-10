package validator

import (
	"errors"
	"vpub/model"
	"vpub/storage"
)

func ValidateGroupCreation(store *storage.Storage, request model.GroupRequest) error {
	if checkStringIsEmpty(request.Name) {
		return errors.New("group name can't be empty")
	}
	exists, err := store.GroupNameExists(request.Name)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("group name already exists")
	}
	return nil
}

func ValidateGroupModification(store *storage.Storage, groupID int64, request model.GroupRequest) error {
	if checkStringIsEmpty(request.Name) {
		return errors.New("group name can't be empty")
	}
	exists, err := store.AnotherGroupExists(groupID, request.Name)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("a group with that name already exists")
	}
	return nil
}
