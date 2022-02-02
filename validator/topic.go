package validator

import (
	"errors"
	"vpub/model"
)

func ValidateTopicRequest(request model.TopicRequest) error {
	if checkStringIsEmpty(request.Subject) {
		return errors.New("Topic subject can't be empty")
	}

	if checkStringIsEmpty(request.Content) {
		return errors.New("Topic content can't be empty")
	}

	return nil
}
