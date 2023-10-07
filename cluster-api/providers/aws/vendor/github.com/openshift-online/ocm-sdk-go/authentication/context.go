/*
Copyright (c) 2019 Red Hat, Inc.

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

// This file contains functions that extract information from the context.

package authentication

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

// ContextWithToken creates a new context containing the given token.
func ContextWithToken(parent context.Context, token *jwt.Token) context.Context {
	return context.WithValue(parent, tokenKeyValue, token)
}

// TokenFromContext extracts the JSON web token of the user from the context. If no token is found
// in the context then the result will be nil.
func TokenFromContext(ctx context.Context) (result *jwt.Token, err error) {
	switch token := ctx.Value(tokenKeyValue).(type) {
	case nil:
	case *jwt.Token:
		result = token
	default:
		err = fmt.Errorf(
			"expected a token in the '%s' context value, but got '%T'",
			tokenKeyValue, token,
		)
	}
	return
}

// BearerFromContext extracts the bearer token of the user from the context. If no user is found in
// the context then the result will be the empty string.
func BearerFromContext(ctx context.Context) (result string, err error) {
	token, err := TokenFromContext(ctx)
	if err != nil {
		return
	}
	if token == nil {
		return
	}
	result = token.Raw
	return
}

// tokenKeyType is the type of the key used to store the token in the context.
type tokenKeyType string

// tokenKeyValue is the key used to store the token in the context:
const tokenKeyValue tokenKeyType = "token"
