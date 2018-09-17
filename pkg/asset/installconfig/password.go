package installconfig

import (
	"bufio"
	"fmt"
	"os"

	"github.com/openshift/installer/pkg/asset"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	passwordPrompt = fmt.Sprintf("Enter password:")
	reEnterPrompt  = fmt.Sprintf("Re-enter password:")
)

// password is an asset that queries the user for the password to the cluster
//
// Contents[0] is the actual string form of the password
type password struct {
	InputReader *bufio.Reader
}

var _ asset.Asset = (*password)(nil)

// Dependencies returns no dependencies.
func (a *password) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate queries for input from the user.
func (a *password) Generate(map[asset.Asset]*asset.State) (*asset.State, error) {
	password, err := a.queryUserForPassword()
	return assetStateForStringContents(password), err
}

// queryUserForPassword for prompts Stdin for password without echoing anything on screen
// TODO: mask the password characters
func (a *password) queryUserForPassword() (string, error) {
	for {
		fmt.Printf(passwordPrompt)
		input, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Printf("failed to read password %v", err)
			return "", err
		}
		fmt.Println()
		fmt.Printf(reEnterPrompt)
		confirmInput, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Printf("failed to read password %v", err)
			return "", err
		}
		fmt.Println()
		if string(input) == string(confirmInput) {
			return string(input), nil
		}
		fmt.Println("Password did not match.")
	}
}
