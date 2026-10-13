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

// This file contains helper functions used in several places in the package.

package authentication

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// tokenRemaining determines if the given token will eventually expire (offile access tokens and
// opaque tokens, for example, never expire) and the time till it expires. That time will be
// positive if the token isn't expired, and negative if the token has already expired.
//
// For tokens that don't have the `exp` claim, or that have it with value zero (typical for offline
// access tokens) the result will always be `false` and zero.
func tokenRemaining(token *tokenInfo, now time.Time) (expires bool, duration time.Duration,
	err error) {
	// For opaque tokens we can't use the claims to determine when they expire, so we will
	// assume that they never expire.
	if token == nil || token.object == nil {
		return
	}

	// For JSON web tokens we use tthe `exp` claim to determine when they expire.
	claims, ok := token.object.Claims.(jwt.MapClaims)
	if !ok {
		err = fmt.Errorf("expected map claims but got %T", claims)
		return
	}
	var exp float64
	claim, ok := claims["exp"]
	if !ok {
		return
	}
	exp, ok = claim.(float64)
	if !ok {
		err = fmt.Errorf("expected floating point 'exp' but got %T", claim)
		return
	}
	if exp == 0 {
		return
	}
	duration = time.Unix(int64(exp), 0).Sub(now)
	expires = true
	return
}
