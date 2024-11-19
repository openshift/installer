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

const (
	// Supported authentication types.
	AUTHTYPE_BASIC        = "basic"
	AUTHTYPE_BEARER_TOKEN = "bearerToken"
	AUTHTYPE_NOAUTH       = "noAuth"
	AUTHTYPE_IAM          = "iam"
	AUTHTYPE_CP4D         = "cp4d"

	// Names of properties that can be defined as part of an external configuration (credential file, env vars, etc.).
	// Example:  export MYSERVICE_URL=https://myurl

	// Service client properties.
	PROPNAME_SVC_URL         = "URL"
	PROPNAME_SVC_DISABLE_SSL = "DISABLE_SSL"

	// Authenticator properties.
	PROPNAME_AUTH_TYPE        = "AUTH_TYPE"
	PROPNAME_USERNAME         = "USERNAME"
	PROPNAME_PASSWORD         = "PASSWORD"
	PROPNAME_BEARER_TOKEN     = "BEARER_TOKEN"
	PROPNAME_AUTH_URL         = "AUTH_URL"
	PROPNAME_AUTH_DISABLE_SSL = "AUTH_DISABLE_SSL"
	PROPNAME_APIKEY           = "APIKEY"
	PROPNAME_CLIENT_ID        = "CLIENT_ID"
	PROPNAME_CLIENT_SECRET    = "CLIENT_SECRET"

	// SSL error
	SSL_CERTIFICATION_ERROR = "x509: certificate"

	// Common error messages.
	ERRORMSG_PROP_MISSING            = "The %s property is required but was not specified."
	ERRORMSG_PROP_INVALID            = "The %s property is invalid. Please remove any surrounding {, }, or \" characters."
	ERRORMSG_NO_AUTHENTICATOR        = "Authentication information was not properly configured."
	ERRORMSG_AUTHTYPE_UNKNOWN        = "Unrecognized authentication type: %s"
	ERRORMSG_PROPS_MAP_NIL           = "The 'properties' map cannot be nil."
	ERRORMSG_SSL_VERIFICATION_FAILED = "The connection failed because the SSL certificate is not valid. To use a " +
		"self-signed certificate, disable verification of the server's SSL certificate " +
		"by invoking the DisableSSLVerification() function on your service instance " +
		"and/or use the DisableSSLVerification option of the authenticator."
)
