package validator

import (
	"errors"
	"vpub/model"
)

func ValidatePostRequest(request model.PostRequest) error {
	if checkStringIsEmpty(request.Subject) {
		return errors.New("Post subject can't be empty")
	}

	if checkStringIsEmpty(request.Content) {
		return errors.New("Post content can't be empty")
	}

	return nil
}
