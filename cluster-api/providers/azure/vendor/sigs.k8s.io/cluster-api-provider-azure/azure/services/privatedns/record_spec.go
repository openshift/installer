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

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/privatedns/armprivatedns"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure/converters"
)

// RecordSpec defines the specification for a record set.
type RecordSpec struct {
	Record        infrav1.AddressRecord
	ZoneName      string
	ResourceGroup string
}

// ResourceName returns the name of a record set.
func (s RecordSpec) ResourceName() string {
	return s.Record.Hostname
}

// OwnerResourceName returns the zone name of a record set.
func (s RecordSpec) OwnerResourceName() string {
	return s.ZoneName
}

// ResourceGroupName returns the name of the resource group of a record set.
func (s RecordSpec) ResourceGroupName() string {
	return s.ResourceGroup
}

// Parameters returns the parameters for a record set.
func (s RecordSpec) Parameters(ctx context.Context, existing interface{}) (params interface{}, err error) {
	if existing != nil {
		if _, ok := existing.(armprivatedns.RecordSet); !ok {
			return nil, errors.Errorf("%T is not an armprivatedns.RecordSet", existing)
		}
	}
	set := armprivatedns.RecordSet{
		Properties: &armprivatedns.RecordSetProperties{
			TTL: ptr.To[int64](300),
		},
	}
	recordType := converters.GetRecordType(s.Record.IP)
	switch recordType {
	case armprivatedns.RecordTypeA:
		set.Properties.ARecords = []*armprivatedns.ARecord{{
			IPv4Address: &s.Record.IP,
		}}
	case armprivatedns.RecordTypeAAAA:
		set.Properties.AaaaRecords = []*armprivatedns.AaaaRecord{{
			IPv6Address: &s.Record.IP,
		}}
	default:
		return nil, errors.Errorf("unknown record type %s", recordType)
	}

	return set, nil
}
