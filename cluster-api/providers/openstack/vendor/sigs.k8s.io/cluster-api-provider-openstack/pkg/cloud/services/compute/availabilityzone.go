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

package compute

import (
	"fmt"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/availabilityzones"
)

func (s *Service) GetAvailabilityZones() ([]availabilityzones.AvailabilityZone, error) {
	availabilityZoneList, err := s.getComputeClient().ListAvailabilityZones()
	if err != nil {
		return nil, fmt.Errorf("error extracting availability zone list: %v", err)
	}

	available := make([]availabilityzones.AvailabilityZone, 0, len(availabilityZoneList))

	for _, az := range availabilityZoneList {
		if az.ZoneState.Available {
			available = append(available, az)
		}
	}

	return available, nil
}
