package core

// (C) Copyright IBM Corp. 2019.
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
	"fmt"
	"net/http"
)

// The BearerTokenAuthenticator will authenticate requests by adding
// a "Bearer"-type Authorization header that contains the configured bearer token value.
type BearerTokenAuthenticator struct {

	// [Required] - the bearer token value to be used to authenticate request.
	BearerToken string
}

// NewBearerTokenAuthenticator: Constructs a new BearerTokenAuthenticator instance.
func NewBearerTokenAuthenticator(bearerToken string) (*BearerTokenAuthenticator, error) {
	obj := &BearerTokenAuthenticator{
		BearerToken: bearerToken,
	}
	if err := obj.Validate(); err != nil {
		return nil, err
	}
	return obj, nil
}

// newBearerTokenAuthenticator : Constructs a new BearerTokenAuthenticator instance from a map.
func newBearerTokenAuthenticatorFromMap(properties map[string]string) (*BearerTokenAuthenticator, error) {
	if properties == nil {
		return nil, fmt.Errorf(ERRORMSG_PROPS_MAP_NIL)
	}

	return NewBearerTokenAuthenticator(properties[PROPNAME_BEARER_TOKEN])
}

func (BearerTokenAuthenticator) AuthenticationType() string {
	return AUTHTYPE_BEARER_TOKEN
}

// Authenticate: authenticates the specified request by adding an Authorization header
// that contains the bearer token value.
func (this BearerTokenAuthenticator) Authenticate(request *http.Request) error {
	request.Header.Set("Authorization", fmt.Sprintf(`Bearer %s`, this.BearerToken))
	return nil
}

// Validate: validates the configuration
func (this BearerTokenAuthenticator) Validate() error {
	if this.BearerToken == "" {
		return fmt.Errorf(ERRORMSG_PROP_MISSING, "BearerToken")
	}
	return nil
}
