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

// Package bytes provides utilities for working with byte arrays.
package bytes

import (
	"bytes"
	"encoding/base64"
)

// Split takes a byte array, optionally encodes it in base64. Should it be encoded,
// the result of encoding is split. Otherwise, the array is simply split as is.
func Split(data []byte, encode bool, maxSize int, iterFunc func([]byte)) {
	if encode {
		encoded := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
		base64.StdEncoding.Encode(encoded, data)
		data = encoded
	}
	buff := bytes.NewBuffer(data)
	for {
		chunk := buff.Next(maxSize)
		if len(chunk) == 0 {
			return
		}
		iterFunc(chunk)
	}
}
