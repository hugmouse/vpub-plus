package validator

import (
	"errors"
	"vpub/model"
	"vpub/storage"
)

func ValidateBoardCreation(store *storage.Storage, request model.BoardRequest) error {
	if checkStringIsEmpty(request.Name) {
		return errors.New("board name can't be empty")
	}

	boardExists, err := store.BoardNameExists(request.Name)
	if err != nil {
		return err
	}
	if boardExists {
		return errors.New("board name already exists")
	}

	return nil
}

func ValidateBoardModification(store *storage.Storage, boardId int64, request model.BoardRequest) error {
	if checkStringIsEmpty(request.Name) {
		return errors.New("board name can't be empty")
	}

	anotherBoardExists, err := store.AnotherBoardExists(boardId, request.Name)
	if err != nil {
		return err
	}
	if anotherBoardExists {
		return errors.New("board name already exists")
	}

	return nil
}
