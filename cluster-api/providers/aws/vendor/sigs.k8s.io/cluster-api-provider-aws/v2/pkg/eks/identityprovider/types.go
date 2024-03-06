/*
Copyright 2021 The Kubernetes Authors.

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

package identityprovider

import (
	"reflect"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
)

// OidcIdentityProviderConfig represents a normalized version of the configuration for an OpenID Connect (OIDC)
// identity provider configuration. To reconcile the config we are going to get the version from EKS and
// AWSManagedControlPlane and will need to have one consistent version of string values from each API.
type OidcIdentityProviderConfig struct {
	ClientID                   string
	GroupsClaim                string
	GroupsPrefix               string
	IdentityProviderConfigArn  string
	IdentityProviderConfigName string
	IssuerURL                  string
	RequiredClaims             map[string]string
	Status                     string
	Tags                       infrav1.Tags
	UsernameClaim              string
	UsernamePrefix             string
}

// IsEqual returns true if the OidcIdentityProviderConfig is equal to the supplied one.
func (o *OidcIdentityProviderConfig) IsEqual(other *OidcIdentityProviderConfig) bool {
	if o == other {
		return true
	}

	if o.ClientID != other.ClientID {
		return false
	}

	if o.GroupsClaim != other.GroupsClaim {
		return false
	}

	if o.GroupsPrefix != other.GroupsPrefix {
		return false
	}

	if o.IdentityProviderConfigName != other.IdentityProviderConfigName {
		return false
	}

	if o.IssuerURL != other.IssuerURL {
		return false
	}

	if !reflect.DeepEqual(o.RequiredClaims, other.RequiredClaims) {
		return false
	}

	if o.UsernameClaim != other.UsernameClaim {
		return false
	}

	if o.UsernamePrefix != other.UsernamePrefix {
		return false
	}

	return true
}
