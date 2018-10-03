package asset

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// UserProvided generates an asset that is supplied by a user.
type UserProvided struct {
	AssetName      string
	Question       *survey.Question
	EnvVarName     string
	PathEnvVarName string
}

var _ Asset = (*UserProvided)(nil)

// Dependencies returns no dependencies.
func (a *UserProvided) Dependencies() []Asset {
	return []Asset{}
}

// Generate queries for input from the user.
func (a *UserProvided) Generate(map[Asset]*State) (state *State, err error) {
	defer func() {
		if err != nil {
			err = errors.Wrapf(err, "failed to acquire user-provided input %s", a.AssetName)
		}
	}()

	var response string

	if value, ok := os.LookupEnv(a.EnvVarName); ok {
		response = value
	} else if path, ok := os.LookupEnv(a.PathEnvVarName); ok {
		value, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read file from %s", a.PathEnvVarName)
		}
		response = string(value)
	}

	if response == "" {
		if err := survey.Ask([]*survey.Question{a.Question}, &response); err != nil {
			return nil, errors.Wrap(err, "failed to Ask")
		}
	} else if a.Question.Validate != nil {
		if err := a.Question.Validate(response); err != nil {
			return nil, errors.Wrap(err, "validation failed")
		}
	}

	return &State{
		Contents: []Content{{
			Data: []byte(response),
		}},
	}, nil
}

// Name returns the human-friendly name of the asset.
func (a *UserProvided) Name() string {
	return a.AssetName
}
