package validator

import (
	"errors"
	"vpub/model"
)

func ValidateTopicRequest(request model.TopicRequest) error {
	if checkStringIsEmpty(request.Subject) {
		return errors.New("topic subject can't be empty")
	}

	if checkStringIsEmpty(request.Content) {
		return errors.New("topic content can't be empty")
	}

	return nil
}
