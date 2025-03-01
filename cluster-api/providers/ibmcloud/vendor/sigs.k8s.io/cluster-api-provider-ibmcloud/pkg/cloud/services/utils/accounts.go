/*
Copyright 2022 The Kubernetes Authors.

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

package utils

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"

	"github.com/IBM/go-sdk-core/v5/core"

	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/authenticator"
)

// GetAccount is function parses the account number from the token and returns it.
func GetAccount(auth core.Authenticator) (string, error) {
	// fake request to get a barer token from the request header
	ctx := context.TODO()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", http.NoBody)
	if err != nil {
		return "", err
	}
	err = auth.Authenticate(req)
	if err != nil {
		return "", err
	}
	bearerToken := req.Header.Get("Authorization")
	if strings.HasPrefix(bearerToken, "Bearer") {
		bearerToken = bearerToken[7:]
	}
	token, err := jwt.Parse(bearerToken, func(_ *jwt.Token) (interface{}, error) {
		return "", nil
	})
	if err != nil && !strings.Contains(err.Error(), "key is of invalid type") {
		return "", err
	}

	return token.Claims.(jwt.MapClaims)["account"].(map[string]interface{})["bss"].(string), nil
}

// GetAccountIDFunc is a variable that will hold the function reference.
var GetAccountIDFunc = GetAccountID // Default to the original function

// GetAccountIDWrapper is a function that calls GetAccountIDFunc.
func GetAccountIDWrapper() (string, error) {
	return GetAccountIDFunc() // Call the function that GetAccountIDFunc points to
}

// GetAccountID will parse and returns user cloud account ID.
func GetAccountID() (string, error) {
	auth, err := authenticator.GetAuthenticator()
	if err != nil {
		return "", err
	}
	return GetAccount(auth)
}
