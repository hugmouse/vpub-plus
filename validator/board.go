package validator

import (
	"errors"
	"vpub/model"
	"vpub/storage"
)

func ValidateBoardCreation(store *storage.Storage, request model.BoardRequest) error {
	if checkStringIsEmpty(request.Name) {
		return errors.New("Board name can't be empty")
	}

	if store.BoardNameExists(request.Name) {
		return errors.New("Board name already exists")
	}

	return nil
}

func ValidateBoardModification(store *storage.Storage, boardId int64, request model.BoardRequest) error {
	if checkStringIsEmpty(request.Name) {
		return errors.New("Board name can't be empty")
	}

	if store.AnotherBoardExists(boardId, request.Name) {
		return errors.New("Board name already exists")
	}

	return nil
}
