/*
Copyright 2018 The Kubernetes Authors.

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

package network

import (
	"context"
	"sort"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/filter"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
)

func (s *Service) getAvailableZones() ([]string, error) {
	out, err := s.EC2Client.DescribeAvailabilityZonesWithContext(context.TODO(), &ec2.DescribeAvailabilityZonesInput{
		Filters: []*ec2.Filter{
			filter.EC2.Available(),
			filter.EC2.IgnoreLocalZones(),
		},
	})
	if err != nil {
		record.Eventf(s.scope.InfraCluster(), "FailedDescribeAvailableZone", "Failed getting available zones: %v", err)
		return nil, errors.Wrap(err, "failed to describe availability zones")
	}

	zones := make([]string, 0, len(out.AvailabilityZones))
	for _, zone := range out.AvailabilityZones {
		zones = append(zones, *zone.ZoneName)
	}

	sort.Strings(zones)
	return zones, nil
}
