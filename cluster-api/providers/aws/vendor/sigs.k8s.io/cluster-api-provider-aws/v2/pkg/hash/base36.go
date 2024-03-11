/*
Copyright 2019 The Kubernetes Authors.

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

package hash

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2b"
	_ "k8s.io/apimachinery/pkg/util/intstr" // keep the blank import to include intstr.
)

const base36set = "0123456789abcdefghijklmnopqrstuvwxyz"

// Base36TruncatedHash returns a consistent hash using blake2b
// and truncating the byte values to alphanumeric only
// of a fixed length specified by the consumer.
func Base36TruncatedHash(str string, length int) (string, error) {
	hasher, err := blake2b.New(length, nil)
	if err != nil {
		return "", errors.Wrap(err, "unable to create hash function")
	}

	if _, err := hasher.Write([]byte(str)); err != nil {
		return "", errors.Wrap(err, "unable to write hash")
	}
	return base36Truncate(hasher.Sum(nil)), nil
}

// base36Truncate returns a string that is base36 compliant
// It is not an encoding since it returns a same-length string
// for any byte value.
func base36Truncate(bytes []byte) string {
	var chars string
	for _, bite := range bytes {
		idx := int(bite) % 36
		chars += string(base36set[idx])
	}

	return chars
}
