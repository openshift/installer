/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package genericarmclient

import (
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"

	"github.com/Azure/azure-service-operator/v2/internal/version"
)

type userAgentPolicy struct {
	userAgent string
}

var userAgent = "aso-controller/" + version.BuildVersion

// NewUserAgentPolicy creates a new policy.Policy appending the specified user agent to each request
func NewUserAgentPolicy(userAgent string) policy.Policy {
	return userAgentPolicy{
		userAgent: userAgent,
	}
}

func (p userAgentPolicy) Do(req *policy.Request) (*http.Response, error) {
	if p.userAgent == "" {
		return req.Next()
	}

	newUserAgent := p.userAgent
	// preserve the existing User-Agent string
	if ua := req.Raw().Header.Get("User-Agent"); ua != "" {
		newUserAgent = fmt.Sprintf("%s %s", ua, newUserAgent)
	}
	req.Raw().Header.Set("user-Agent", newUserAgent)
	return req.Next()
}
