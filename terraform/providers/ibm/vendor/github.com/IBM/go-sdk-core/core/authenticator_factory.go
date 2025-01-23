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
	"strings"
)

// GetAuthenticatorFromEnvironment instantiates an Authenticator using service properties
// retrieved from external config sources.
func GetAuthenticatorFromEnvironment(credentialKey string) (authenticator Authenticator, err error) {
	properties, err := getServiceProperties(credentialKey)
	if len(properties) == 0 {
		return
	}

	// Default the authentication type to IAM if not specified.
	authType := properties[PROPNAME_AUTH_TYPE]
	if authType == "" {
		authType = AUTHTYPE_IAM
	}

	// Create the authenticator appropriate for the auth type.
	if strings.EqualFold(authType, AUTHTYPE_BASIC) {
		authenticator, err = newBasicAuthenticatorFromMap(properties)
	} else if strings.EqualFold(authType, AUTHTYPE_BEARER_TOKEN) {
		authenticator, err = newBearerTokenAuthenticatorFromMap(properties)
	} else if strings.EqualFold(authType, AUTHTYPE_IAM) {
		authenticator, err = newIamAuthenticatorFromMap(properties)
	} else if strings.EqualFold(authType, AUTHTYPE_CP4D) {
		authenticator, err = newCloudPakForDataAuthenticatorFromMap(properties)
	} else if strings.EqualFold(authType, AUTHTYPE_NOAUTH) {
		authenticator, err = NewNoAuthAuthenticator()
	} else {
		err = fmt.Errorf(ERRORMSG_AUTHTYPE_UNKNOWN, authType)
	}

	return
}
