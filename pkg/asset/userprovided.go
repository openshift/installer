package asset

import (
	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// GenerateUserProvidedAsset queries for input from the user.
func GenerateUserProvidedAsset(inputName string, question *survey.Question) (string, error) {
	var response string
	err := survey.Ask([]*survey.Question{question}, &response)
	return response, errors.Wrapf(err, "failed to acquire user-provided input %s", inputName)
}
