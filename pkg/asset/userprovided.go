package asset

import (
	"bufio"
	"fmt"
	"io"
)

// UserProvided generates an asset that is supplied by a user.
type userProvided struct {
	inputReader *bufio.Reader
	prompt      string
	validation  func(string) (string, bool)
}

var _ Asset = (*userProvided)(nil)

// Dependencies returns no dependencies.
func (a *userProvided) Dependencies() []Asset {
	return []Asset{}
}

// Generate queries for input from the user.
func (a *userProvided) Generate(map[Asset]*State) (*State, error) {
	input := a.queryUser()
	return &State{
		Contents: []Content{
			{Data: []byte(input)},
		},
	}, nil
}

func (a *userProvided) queryUser() string {
	for {
		fmt.Println(a.prompt)
		input, err := a.inputReader.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Println("Could not understand response. Please retry.")
			continue
		}
		if input != "" && input[len(input)-1] == '\n' {
			input = input[:len(input)-1]
		}
		if a.validation != nil {
			validatedInput, ok := a.validation(input)
			if !ok {
				continue
			}
			input = validatedInput
		}
		return input
	}
}
