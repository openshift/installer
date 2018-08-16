package asset

import (
	"fmt"
)

// UserProvided generates an asset that is supplied by a user.
type UserProvided struct {
	Prompt string
}

var _ Asset = (*UserProvided)(nil)

// Dependencies returns no dependencies.
func (a *UserProvided) Dependencies() []Asset {
	return []Asset{}
}

// Generate queries for input from the user.
func (a *UserProvided) Generate(map[Asset]*State) (*State, error) {
	input := queryUser(a.Prompt)
	return &State{
		Contents: []Content{
			{Data: []byte(input)},
		},
	}, nil
}

func queryUser(prompt string) string {
	for {
		fmt.Print(prompt)
		var input string
		if _, err := fmt.Scanln(&input); err != nil {
			fmt.Println("Could not understand response. Please retry.")
			continue
		}
		return input
	}
}
