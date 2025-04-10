/*
Copyright 2020 The Kubernetes Authors.

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

package roleassignments

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v2"
	"github.com/pkg/errors"

	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

// azureClient contains the Azure go-sdk Client.
type azureClient struct {
	roleassignments armauthorization.RoleAssignmentsClient
}

// newClient creates a new role assignments client from an authorizer.
func newClient(auth azure.Authorizer) (*azureClient, error) {
	opts, err := azure.ARMClientOptions(auth.CloudEnvironment(), auth.BaseURI())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create roleassignments client options")
	}
	factory, err := armauthorization.NewClientFactory(auth.SubscriptionID(), auth.Token(), opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create armauthorization client factory")
	}
	return &azureClient{*factory.NewRoleAssignmentsClient()}, nil
}

// Get gets the specified role assignment.
func (ac *azureClient) Get(ctx context.Context, spec azure.ResourceSpecGetter) (interface{}, error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "roleassignments.azureClient.Get")
	defer done()

	resp, err := ac.roleassignments.Get(ctx, spec.ResourceGroupName(), spec.ResourceName(), nil)
	if err != nil {
		return nil, err
	}
	return resp.RoleAssignment, nil
}

// CreateOrUpdateAsync creates a roleassignment.
// Creating a roleassignment is not a long running operation, so we don't ever return a poller.
func (ac *azureClient) CreateOrUpdateAsync(ctx context.Context, spec azure.ResourceSpecGetter, resumeToken string, parameters interface{}) (result interface{}, poller *runtime.Poller[armauthorization.RoleAssignmentsClientCreateResponse], err error) { //nolint:revive // keeping resumeToken for readability
	ctx, _, done := tele.StartSpanWithLogger(ctx, "roleassignments.azureClient.CreateOrUpdateAsync")
	defer done()

	createParams, ok := parameters.(armauthorization.RoleAssignmentCreateParameters)
	if !ok {
		return nil, nil, errors.Errorf("%T is not an armauthorization.RoleAssignmentCreateParameters", parameters)
	}
	resp, err := ac.roleassignments.Create(ctx, spec.OwnerResourceName(), spec.ResourceName(), createParams, nil)
	return resp.RoleAssignment, nil, err
}
