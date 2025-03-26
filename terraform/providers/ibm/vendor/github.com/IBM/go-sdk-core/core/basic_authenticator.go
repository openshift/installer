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

// The BasicAuthenticator will perform authentication on outbound requests by adding
// a "Basic" type Authorization header that contains the base64-encoded username and password.
type BasicAuthenticator struct {
	// [Required] - the basic auth username and password.
	Username string
	Password string
}

// NewBasicAuthenticator: Constructs a new BasicAuthenticator instance.
func NewBasicAuthenticator(username string, password string) (*BasicAuthenticator, error) {
	obj := &BasicAuthenticator{
		Username: username,
		Password: password,
	}
	if err := obj.Validate(); err != nil {
		return nil, err
	}
	return obj, nil
}

// newBasicAuthenticatorFromMap: Constructs a new BasicAuthenticator instance from a map.
func newBasicAuthenticatorFromMap(properties map[string]string) (*BasicAuthenticator, error) {
	if properties == nil {
		return nil, fmt.Errorf(ERRORMSG_PROPS_MAP_NIL)
	}

	return NewBasicAuthenticator(properties[PROPNAME_USERNAME], properties[PROPNAME_PASSWORD])
}

func (BasicAuthenticator) AuthenticationType() string {
	return AUTHTYPE_BASIC
}

// Authenticate: authenticates the specified request by adding an Authorizatin header.
func (this BasicAuthenticator) Authenticate(request *http.Request) error {
	request.SetBasicAuth(this.Username, this.Password)
	return nil
}

// Validate: validates the configuration
func (this BasicAuthenticator) Validate() error {
	if this.Username == "" {
		return fmt.Errorf(ERRORMSG_PROP_MISSING, "Username")
	}

	if this.Password == "" {
		return fmt.Errorf(ERRORMSG_PROP_MISSING, "Password")
	}

	if HasBadFirstOrLastChar(this.Username) {
		return fmt.Errorf(ERRORMSG_PROP_INVALID, "Username")
	}

	if HasBadFirstOrLastChar(this.Password) {
		return fmt.Errorf(ERRORMSG_PROP_INVALID, "Password")
	}

	return nil
}
