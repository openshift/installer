/*
Copyright (c) 2018 Red Hat, Inc.

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

// This file contains the implementations of the methods of the connection that handle OpenID
// authentication tokens.

package sdk

import (
	"context"

	"time"
)

// Tokens returns the access and refresh tokens that are currently in use by the connection. If it
// is necessary to request new tokens because they weren't requested yet, or because they are
// expired, this method will do it and will return an error if it fails.
//
// If new tokens are needed the request will be retried with an exponential backoff.
//
// This operation is potentially lengthy, as it may require network communication. Consider using a
// context and the TokensContext method.
// The returned access and refresh tokens are empty strings if the
// connection does not use authentication. In that case no error is
// returned either
func (c *Connection) Tokens(expiresIn ...time.Duration) (access, refresh string, err error) {
	if len(expiresIn) == 1 {
		access, refresh, err = c.TokensContext(context.Background(), expiresIn[0])
	} else {
		access, refresh, err = c.TokensContext(context.Background())
	}
	return

}

// TokensContext returns the access and refresh tokens that are currently in use by the connection.
// If it is necessary to request new tokens because they weren't requested yet, or because they are
// expired, this method will do it and will return an error if it fails.
//
// If new tokens are needed the request will be retried with an exponential backoff.
// The returned access and refresh tokens are empty strings if the
// connection does not use authentication. In that case no error is
// returned either
func (c *Connection) TokensContext(ctx context.Context, expiresIn ...time.Duration) (access,
	refresh string, err error) {
	if c.authnWrapper != nil {
		access, refresh, err = c.authnWrapper.Tokens(ctx, expiresIn...)
	}
	return
}
