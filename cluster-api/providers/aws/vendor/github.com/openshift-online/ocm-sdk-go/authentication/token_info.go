/*
Copyright (c) 2021 Red Hat, Inc.

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

package authentication

import (
	"github.com/golang-jwt/jwt/v4"
)

// tokenInfo stores information about a token. We need to store both the original text and the
// parsed objects because some tokens (refresh tokens in particular) may be opaque strings instead
// of JSON web tokens.
type tokenInfo struct {
	text   string
	object *jwt.Token
}
