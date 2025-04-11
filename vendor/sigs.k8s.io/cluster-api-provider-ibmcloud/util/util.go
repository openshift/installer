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

package util

import (
	"fmt"

	regionUtil "github.com/ppc64le-cloud/powervs-utils"

	"k8s.io/utils/ptr"

	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/endpoints"
)

// GetTransitGatewayLocationAndRouting returns appropriate location and routing suitable for transit gateway.
// routing indicates whether to enable global routing on transit gateway or not.
// returns true when PowerVS and VPC region are not same otherwise false.
func GetTransitGatewayLocationAndRouting(powerVSZone *string, vpcRegion *string) (*string, *bool, error) {
	if powerVSZone == nil {
		return nil, nil, fmt.Errorf("powervs zone is not set")
	}
	powerVSRegion := endpoints.ConstructRegionFromZone(*powerVSZone)

	if vpcRegion != nil {
		routing := regionUtil.IsGlobalRoutingRequiredForTG(powerVSRegion, *vpcRegion)
		return vpcRegion, &routing, nil
	}
	location, err := regionUtil.VPCRegionForPowerVSRegion(powerVSRegion)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch vpc region associated with powervs region '%s': %w", powerVSRegion, err)
	}

	// since VPC region is not set and used PowerVS region to calculate the transit gateway location, hence returning local routing as default.
	return &location, ptr.To(false), nil
}
