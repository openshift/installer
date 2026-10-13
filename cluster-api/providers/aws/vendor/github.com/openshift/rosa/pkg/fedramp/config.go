/*
Copyright (c) 2022 Red Hat, Inc.

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

// This file contains the types and functions used to manage the configuration of the command line
// client when running in FedRAMP mode.

package fedramp

import "fmt"

var regions = []string{"us-gov-west-1", "us-gov-east-1"}

func IsGovRegion(region string) bool {
	for _, r := range regions {
		if r == region {
			return true
		}
	}
	return false
}

func IsValidEnv(env string) bool {
	for urlAlias := range URLAliases {
		if urlAlias == env {
			return true
		}
	}
	return false
}

// JumpAccounts are the various of AWS accounts used for the installer jump role in the various OCM environments
var JumpAccounts = map[string]string{
	"production":  "448648337690",
	"staging":     "448870092490",
	"staging01":   "448870092490",
	"integration": "449053620653",
}

// LoginURLs allows the value of the `--env` option to map to the various login URLs.
var LoginURLs = map[string]string{
	"production":  "https://api.openshiftusgov.com/auth",
	"staging":     "https://api.stage.openshiftusgov.com/auth",
	"staging01":   "https://api01.stage.openshiftusgov.com/auth",
	"integration": "https://api.int.openshiftusgov.com/auth",
}

// AdminLoginURLs allows the value of the `--env` option to map to the various Admin login URLs.
var AdminLoginURLs = map[string]string{
	"production":  "https://api-admin.openshiftusgov.com/auth",
	"staging":     "https://api-admin.stage.openshiftusgov.com/auth",
	"staging01":   "https://api.stage.openshiftusgov.com/auth",
	"integration": "https://api-admin.int.openshiftusgov.com/auth",
}

// URLAliases allows the value of the `--env` option to map to the various API URLs.
var URLAliases = map[string]string{
	"production":  "https://api.openshiftusgov.com",
	"staging":     "https://api.stage.openshiftusgov.com",
	"staging01":   "https://api01.stage.openshiftusgov.com",
	"integration": "https://api.int.openshiftusgov.com",
}

// AdminURLAliases allows the value of the `--env` option to map to the various Admin API URLs.
var AdminURLAliases = map[string]string{
	"production":  "https://api-admin.openshiftusgov.com",
	"staging":     "https://api-admin.stage.openshiftusgov.com",
	"staging01":   "https://api01.stage.openshiftusgov.com",
	"integration": "https://api-admin.int.openshiftusgov.com",
}

const cognitoURL = "auth-fips.us-gov-west-1.amazoncognito.com/oauth2/token"
const keycloakURL = "realms/redhat-external/protocol/openid-connect/token"

// TokenURLs allows the value of the `--env` option to map to the various AWS Cognito token URLs.
var TokenURLs = map[string]string{
	"production":  fmt.Sprintf("https://sso.openshiftusgov.com/%s", keycloakURL),
	"staging":     fmt.Sprintf("https://sso.stage.openshiftusgov.com/%s", keycloakURL),
	"staging01":   fmt.Sprintf("https://sso01.stage.openshiftusgov.com/%s", keycloakURL),
	"integration": fmt.Sprintf("https://sso.int.openshiftusgov.com/%s", keycloakURL),
}

// AdminTokenURLs allows the value of the `--env` option to map to the various Admin AWS Cognito token URLs.
var AdminTokenURLs = map[string]string{
	"production":  fmt.Sprintf("https://ocm-ra-production-domain.%s", cognitoURL),
	"staging":     fmt.Sprintf("https://ocm-ra-stage-domain.%s", cognitoURL),
	"staging01":   fmt.Sprintf("https://ocm-ra-stage-domain.%s", cognitoURL),
	"integration": fmt.Sprintf("https://rh-ocm-appsre-integration.%s", cognitoURL),
}

// ClientID stores the client id for use with all `--env` options for Keycloak authentication flow.
// Value is the same for all env's
var ClientID = "console-dot"

// AdminClientIDs allows the value of the `--env` option to map to the various Admin AWS Cognito user pool clients.
var AdminClientIDs = map[string]string{
	"production":  "72ekjh5laouap6qcfis521jlgi",
	"staging":     "1lb687dlpsmsfuj53r3je06vpp",
	"staging01":   "1lb687dlpsmsfuj53r3je06vpp",
	"integration": "20fbrpgl28f8oehp6709mk3nnr",
}
