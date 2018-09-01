package installconfig

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/openshift/installer/installer/pkg/validate"
	"github.com/openshift/installer/pkg/asset"
)

type sshPublicKey struct {
	inputReader *bufio.Reader
}

// Dependencies returns no dependencies.
func (a *sshPublicKey) Dependencies() []asset.Asset {
	return nil
}

// Generate generates the SSH public key asset.
func (a *sshPublicKey) Generate(map[asset.Asset]*asset.State) (state *asset.State, err error) {
	var paths []string
	var pubKeys map[string]string
	home := os.Getenv("HOME")
	if home != "" {
		paths, err = filepath.Glob(filepath.Join(home, ".ssh", "*.pub"))
		if err != nil {
			return nil, err
		}

		if len(paths) > 0 {
			pubKeys = map[string]string{}
		}

		for _, path := range paths {
			pubKeyBytes, err := ioutil.ReadFile(path)
			if err != nil {
				continue
			}
			pubKey := string(pubKeyBytes)

			err = validate.OpenSSHPublicKey(pubKey)
			if err != nil {
				continue
			}

			pubKeys[path] = pubKey
		}

		paths = []string{}
		for path := range pubKeys {
			paths = append(paths, path)
		}
		sort.Strings(paths)
	}

	promptLines := []string{"SSH Public Key:"}
	if len(paths) == 0 {
		promptLines = append(
			promptLines,
			"Enter an empty string or your public key (e.g. 'ssh-rsa AAAA...')",
		)
	} else {
		promptLines = append(
			promptLines,
			"Enter an empty string, your public key (e.g. 'ssh-rsa AAAA...'), or one of the following numbers:",
		)
		for i, path := range paths {
			promptLines = append(
				promptLines,
				fmt.Sprintf("%d: %s", i+1, path),
			)
		}
	}
	prompt := strings.Join(promptLines, "\n")

	var input string
	for {
		fmt.Println(prompt)
		input, err = a.inputReader.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Println("Could not understand response. Please retry.")
			continue
		}
		if input != "" && input[len(input)-1] == '\n' {
			input = input[:len(input)-1]
		}
		var i int
		n, err := fmt.Sscanf(input, "%d", &i)
		if n == len(input) && err == nil && i > 0 && i <= len(paths) {
			path := paths[i-1]
			input = pubKeys[path]
		} else {
			err = validate.OpenSSHPublicKey(input)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
		break
	}

	return &asset.State{
		Contents: []asset.Content{
			{Data: []byte(input)},
		},
	}, nil
}
