/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ssh

import (
	"crypto/rand"
	"crypto/rsa"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

// GenerateSSHKey generates a private and public ssh key.
func GenerateSSHKey() (*rsa.PrivateKey, ssh.PublicKey, error) {
	privateKey, perr := rsa.GenerateKey(rand.Reader, 2048)
	if perr != nil {
		return nil, nil, errors.Wrap(perr, "Failed to generate private key")
	}

	publicRsaKey, perr := ssh.NewPublicKey(&privateKey.PublicKey)
	if perr != nil {
		return nil, nil, errors.Wrap(perr, "Failed to generate public key")
	}

	return privateKey, publicRsaKey, nil
}
