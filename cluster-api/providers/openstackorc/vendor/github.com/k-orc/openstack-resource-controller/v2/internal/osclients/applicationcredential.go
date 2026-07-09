/*
Copyright The ORC Authors.

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

package osclients

import (
	"context"
	"errors"
	"fmt"
	"iter"

	tokens3 "github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/applicationcredentials"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/users"
	"github.com/gophercloud/utils/v2/openstack/clientconfig"
)

type ApplicationCredentialClient interface {
	ListApplicationCredentials(ctx context.Context, userID string, listOpts applicationcredentials.ListOptsBuilder) iter.Seq2[*applicationcredentials.ApplicationCredential, error]
	CreateApplicationCredential(ctx context.Context, userID string, opts applicationcredentials.CreateOptsBuilder) (*applicationcredentials.ApplicationCredential, error)
	DeleteApplicationCredential(ctx context.Context, userID string, resourceID string) error
	GetApplicationCredential(ctx context.Context, resourceID string) (*applicationcredentials.ApplicationCredential, error)
}

type applicationcredentialClient struct{ client *gophercloud.ServiceClient }

// NewApplicationCredentialClient returns a new OpenStack client.
func NewApplicationCredentialClient(providerClient *gophercloud.ProviderClient, providerClientOpts *clientconfig.ClientOpts) (ApplicationCredentialClient, error) {
	client, err := openstack.NewIdentityV3(providerClient, gophercloud.EndpointOpts{
		Region:       providerClientOpts.RegionName,
		Availability: clientconfig.GetEndpointType(providerClientOpts.EndpointType),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create applicationcredential service client: %v", err)
	}

	return &applicationcredentialClient{client}, nil
}

func (c applicationcredentialClient) ListApplicationCredentials(ctx context.Context, userID string, listOpts applicationcredentials.ListOptsBuilder) iter.Seq2[*applicationcredentials.ApplicationCredential, error] {
	pager := applicationcredentials.List(c.client, userID, listOpts)
	return func(yield func(*applicationcredentials.ApplicationCredential, error) bool) {
		_ = pager.EachPage(ctx, yieldPage(applicationcredentials.ExtractApplicationCredentials, yield))
	}
}

func (c applicationcredentialClient) CreateApplicationCredential(ctx context.Context, userID string, opts applicationcredentials.CreateOptsBuilder) (*applicationcredentials.ApplicationCredential, error) {
	return applicationcredentials.Create(ctx, c.client, userID, opts).Extract()
}

func (c applicationcredentialClient) DeleteApplicationCredential(ctx context.Context, userID string, resourceID string) error {
	return applicationcredentials.Delete(ctx, c.client, userID, resourceID).ExtractErr()
}

func (c applicationcredentialClient) GetApplicationCredential(ctx context.Context, resourceID string) (*applicationcredentials.ApplicationCredential, error) {
	// The unique ID of an application credential is not enough to query it from OpenStack
	// OpenStack actually also requires a unique user ID.
	// We can not provide the user ID here, as the function signatures of ORC interfaces
	// expect us to return an OpenStack resource based on a single string.

	// To work around this, we first query ApplicationCredentials for the currently
	// authenticated user which ORC is connected as. If that fails, we iterate over
	// all users we have access to and query their ApplicationCredentials.

	// Currently authenticated user
	userID, err := GetAuthenticatedUserID(c.client.ProviderClient)
	if err == nil {
		appCred, appCredErr := applicationcredentials.Get(ctx, c.client, userID, resourceID).Extract()

		if appCred != nil {
			return appCred, appCredErr
		}
	}

	// If not found in currently authenticated user, try iterating over all users
	userPager := users.List(c.client, nil)
	userIterator := func(yield func(*users.User, error) bool) {
		_ = userPager.EachPage(ctx, yieldPage(users.ExtractUsers, yield))
	}

	for user, userErr := range userIterator {
		if userErr != nil {
			continue
		}

		appCred, appCredErr := applicationcredentials.Get(ctx, c.client, user.ID, resourceID).Extract()

		if appCred != nil {
			return appCred, appCredErr
		}
	}

	return nil, gophercloud.ErrResourceNotFound{
		Name:         resourceID,
		ResourceType: "ApplicationCredential",
	}
}

func GetAuthenticatedUserID(providerClient *gophercloud.ProviderClient) (string, error) {
	r := providerClient.GetAuthResult()
	if r == nil {
		return "", errors.New("no AuthResult available")
	}
	switch r := r.(type) {
	case tokens3.CreateResult:
		u, err := r.ExtractUser()
		if err != nil {
			return "", err
		}
		return u.ID, nil
	default:
		return "", errors.New("wrong AuthResult version")
	}
}

type applicationcredentialErrorClient struct{ error }

// NewApplicationCredentialErrorClient returns a ApplicationCredentialClient in which every method returns the given error.
func NewApplicationCredentialErrorClient(e error) ApplicationCredentialClient {
	return applicationcredentialErrorClient{e}
}

func (e applicationcredentialErrorClient) ListApplicationCredentials(_ context.Context, _ string, _ applicationcredentials.ListOptsBuilder) iter.Seq2[*applicationcredentials.ApplicationCredential, error] {
	return func(yield func(*applicationcredentials.ApplicationCredential, error) bool) {
		yield(nil, e.error)
	}
}

func (e applicationcredentialErrorClient) CreateApplicationCredential(_ context.Context, _ string, _ applicationcredentials.CreateOptsBuilder) (*applicationcredentials.ApplicationCredential, error) {
	return nil, e.error
}

func (e applicationcredentialErrorClient) DeleteApplicationCredential(_ context.Context, _ string, _ string) error {
	return e.error
}

func (e applicationcredentialErrorClient) GetApplicationCredential(_ context.Context, _ string) (*applicationcredentials.ApplicationCredential, error) {
	return nil, e.error
}
