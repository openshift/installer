/*
Copyright 2022 The Kubernetes Authors.

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

package privatedns

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/privatedns/armprivatedns"
	"github.com/pkg/errors"

	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

// azureRecordsClient contains the Azure go-sdk Client for record sets.
type azureRecordsClient struct {
	recordsets *armprivatedns.RecordSetsClient
}

// newRecordSetsClient creates a record sets client from an authorizer.
func newRecordSetsClient(auth azure.Authorizer) (*azureRecordsClient, error) {
	opts, err := azure.ARMClientOptions(auth.CloudEnvironment(), auth.BaseURI())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create recordsets client options")
	}
	factory, err := armprivatedns.NewClientFactory(auth.SubscriptionID(), auth.Token(), opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create armprivatedns client factory")
	}
	return &azureRecordsClient{factory.NewRecordSetsClient()}, nil
}

// Get gets the specified record set. Noop for records.
func (arc *azureRecordsClient) Get(_ context.Context, _ azure.ResourceSpecGetter) (result interface{}, err error) {
	return nil, nil
}

// CreateOrUpdateAsync creates or updates a record asynchronously.
// Creating a record set is not a long-running operation, so we don't ever return a future.
func (arc *azureRecordsClient) CreateOrUpdateAsync(ctx context.Context, spec azure.ResourceSpecGetter, _ string, parameters interface{}) (result interface{}, poller *runtime.Poller[armprivatedns.RecordSetsClientCreateOrUpdateResponse], err error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "privatedns.azureRecordsClient.CreateOrUpdateAsync")
	defer done()

	set, ok := parameters.(armprivatedns.RecordSet)
	if !ok && parameters != nil {
		return nil, nil, errors.Errorf("%T is not an armprivatedns.RecordSet", parameters)
	}

	// Determine record type.
	var (
		recordType armprivatedns.RecordType
		aRecords   = set.Properties.ARecords
		aaaRecords = set.Properties.AaaaRecords
	)
	if len(aRecords) > 0 && (aRecords)[0].IPv4Address != nil {
		recordType = armprivatedns.RecordTypeA
	} else if len(aaaRecords) > 0 && (aaaRecords)[0].IPv6Address != nil {
		recordType = armprivatedns.RecordTypeAAAA
	}

	recordSet, err := arc.recordsets.CreateOrUpdate(ctx, spec.ResourceGroupName(), spec.OwnerResourceName(), recordType, spec.ResourceName(), set, nil)
	if err != nil {
		return nil, nil, err
	}
	return recordSet, nil, err
}

// DeleteAsync deletes a record asynchronously. Noop for records.
func (arc *azureRecordsClient) DeleteAsync(_ context.Context, _ azure.ResourceSpecGetter, _ string) (poller *runtime.Poller[armprivatedns.RecordSetsClientDeleteResponse], err error) {
	return nil, nil
}
