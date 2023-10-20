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

package azure

import (
	"fmt"
	"net/url"
	"path"

	"github.com/pkg/errors"
)

// aliasAuth helps to embed the interface Authorize since the Authorizer interface also defines an Authorizer method and
// the compiler gets confused without a type alias (or renaming the method Authorizer).
type aliasAuth Authorizer

// baseURIAdapter wraps an azure.Authorizer and adds a region to the BaseURI. This is useful if you need to make direct
// calls to a specific Azure region. One possible case is to avoid replication delay when listing resources within a
// resource group. For example, listing the VMSSes within a resource group.
type baseURIAdapter struct {
	aliasAuth
	Region    string
	parsedURL *url.URL
}

// WithRegionalBaseURI returns an authorizer that has a regional base URI, like `https://{region}.management.azure.com`.
func WithRegionalBaseURI(authorizer Authorizer, region string) (Authorizer, error) {
	parsedURI, err := url.Parse(authorizer.BaseURI())
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse the base URI of client")
	}

	return &baseURIAdapter{
		aliasAuth: authorizer,
		Region:    region,
		parsedURL: parsedURI,
	}, nil
}

// BaseURI return a regional base URI, like `https://{region}.management.azure.com`.
func (a *baseURIAdapter) BaseURI() string {
	if a == nil || a.parsedURL == nil || a.Region == "" {
		return a.aliasAuth.BaseURI()
	}

	sansScheme := path.Join(fmt.Sprintf("%s.%s", a.Region, a.parsedURL.Host), a.parsedURL.Path)
	return fmt.Sprintf("%s://%s", a.parsedURL.Scheme, sansScheme)
}
