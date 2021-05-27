package core

// (C) Copyright IBM Corp. 2021.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"encoding/json"
	"fmt"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

// coreJWTClaims are the fields within a JWT's "claims" segment that we're interested in.
type coreJWTClaims struct {
	ExpiresAt int64 `json:"exp,omitempty"`
	IssuedAt  int64 `json:"iat,omitempty"`
}

// parseJWT parses the specified JWT token string and returns an instance of the coreJWTClaims struct.
func parseJWT(tokenString string) (claims *coreJWTClaims, err error) {
	// A JWT consists of three .-separated segments
	segments := strings.Split(tokenString, ".")
	if len(segments) != 3 {
		err = fmt.Errorf("token contains an invalid number of segments")
		return
	}

	// Parse Claims segment.
	var claimBytes []byte
	claimBytes, err = jwt.DecodeSegment(segments[1])
	if err != nil {
		err = fmt.Errorf("error decoding claims segment: %s", err.Error())
		return
	}

	// Now deserialize the claims segment into our coreClaims struct.
	claims = &coreJWTClaims{}
	err = json.Unmarshal(claimBytes, claims)
	if err != nil {
		err = fmt.Errorf("error unmarshalling token: %s", err.Error())
		return
	}

	return
}
