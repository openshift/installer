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
}

var _ Asset = (*UserProvided)(nil)

// Dependencies returns no dependencies.
func (a *UserProvided) Dependencies() []Asset {
	return []Asset{}
}

// Generate queries for input from the user.
func (a *UserProvided) Generate(map[Asset]*State) (*State, error) {
	input := QueryUser(a.InputReader, a.Prompt)
	return &State{
		Contents: []Content{
			{Data: []byte(input)},
		},
	}, nil
}

// QueryUser queries the user for input.
func QueryUser(inputReader *bufio.Reader, prompt string) string {
	for {
		fmt.Printf("%s: ", prompt)
		input, err := inputReader.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Println("Could not understand response. Please retry.")
			continue
		}
		if input != "" && input[len(input)-1] == '\n' {
			input = input[:len(input)-1]
		}
		return input
	}
}
