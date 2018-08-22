package asset

import (
	"bufio"
	"fmt"
	"io"
)

// UserProvided generates an asset that is supplied by a user.
type UserProvided struct {
	InputReader *bufio.Reader
	Prompt      string
	Validation  func(string) (string, bool)
}

var _ Asset = (*UserProvided)(nil)

// Dependencies returns no dependencies.
func (a *UserProvided) Dependencies() []Asset {
	return []Asset{}
}

// Generate queries for input from the user.
func (a *UserProvided) Generate(map[Asset]*State) (*State, error) {
	input := a.queryUser()
	return &State{
		Contents: []Content{
			{Data: []byte(input)},
		},
	}, nil
}

func (a *UserProvided) queryUser() string {
	for {
		fmt.Println(a.Prompt)
		input, err := a.InputReader.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Println("Could not understand response. Please retry.")
			continue
		}
		if input != "" && input[len(input)-1] == '\n' {
			input = input[:len(input)-1]
		}
		if a.Validation != nil {
			validatedInput, ok := a.Validation(input)
			if !ok {
				continue
			}
			input = validatedInput
		}
		return input
	}
}
