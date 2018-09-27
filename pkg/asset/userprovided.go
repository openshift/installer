package asset

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// GenerateUserProvidedAsset queries for input from the user.
func GenerateUserProvidedAsset(inputName string, question *survey.Question, envVarName string) (string, error) {
	return generateUserProvidedAsset(inputName, question, envVarName, "")
}

// GenerateUserProvidedAssetForPath queries for input from the user. The input can
// be read from a file specified in an environment variable.
func GenerateUserProvidedAssetForPath(inputName string, question *survey.Question, envVarName, pathEnvVarName string) (string, error) {
	return generateUserProvidedAsset(inputName, question, envVarName, pathEnvVarName)
}

func generateUserProvidedAsset(inputName string, question *survey.Question, envVarName, pathEnvVarName string) (response string, err error) {
	defer func() {
		if err != nil {
			err = errors.Wrapf(err, "failed to acquire user-provided input %s", inputName)
		}
	}()

	if value, ok := os.LookupEnv(envVarName); ok {
		response = value
	} else if path, ok := os.LookupEnv(pathEnvVarName); ok {
		value, err := ioutil.ReadFile(path)
		if err != nil {
			return "", errors.Wrapf(err, "failed to read file from %s", pathEnvVarName)
		}
		response = string(value)
	}

	if response == "" {
		if err := survey.Ask([]*survey.Question{question}, &response); err != nil {
			return "", errors.Wrap(err, "failed to Ask")
		}
	} else if question.Validate != nil {
		if err := question.Validate(response); err != nil {
			return "", errors.Wrap(err, "validation failed")
		}
	}

	return response, nil
}
