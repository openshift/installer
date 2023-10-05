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
	"github.com/google/go-cmp/cmp"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
)

// OidcIdentityProviderConfig represents the configuration for an OpenID Connect (OIDC)
// identity provider.
type OidcIdentityProviderConfig struct {
	ClientID                   string
	GroupsClaim                *string
	GroupsPrefix               *string
	IdentityProviderConfigArn  *string
	IdentityProviderConfigName string
	IssuerURL                  string
	RequiredClaims             map[string]*string
	Status                     *string
	Tags                       infrav1.Tags
	UsernameClaim              *string
	UsernamePrefix             *string
}

func (o *OidcIdentityProviderConfig) IsEqual(other *OidcIdentityProviderConfig) bool {
	if o == other {
		return true
	}

	if !cmp.Equal(o.ClientID, other.ClientID) {
		return false
	}

	if !cmp.Equal(o.GroupsClaim, other.GroupsClaim) {
		return false
	}

	if !cmp.Equal(o.GroupsPrefix, other.GroupsPrefix) {
		return false
	}

	if !cmp.Equal(o.IdentityProviderConfigName, other.IdentityProviderConfigName) {
		return false
	}

	if !cmp.Equal(o.IssuerURL, other.IssuerURL) {
		return false
	}

	if !cmp.Equal(o.RequiredClaims, other.RequiredClaims) {
		return false
	}

	if !cmp.Equal(o.UsernameClaim, other.UsernameClaim) {
		return false
	}

	if !cmp.Equal(o.UsernamePrefix, other.UsernamePrefix) {
		return false
	}

	return true
}
