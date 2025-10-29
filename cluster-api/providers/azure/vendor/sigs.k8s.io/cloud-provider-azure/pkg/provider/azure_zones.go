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

package provider

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog/v2"

	azcache "sigs.k8s.io/cloud-provider-azure/pkg/cache"
	"sigs.k8s.io/cloud-provider-azure/pkg/consts"
)

var _ cloudprovider.Zones = (*Cloud)(nil)

func (az *Cloud) refreshZones(ctx context.Context, refreshFunc func(ctx context.Context) error) {
	klog.V(2).Info("refreshZones: refreshing zones every 30 minutes.")
	err := wait.PollUntilContextCancel(ctx, consts.ZoneFetchingInterval, false, func(ctx context.Context) (bool, error) {
		_ = refreshFunc(ctx)
		return false, nil
	})
	klog.V(2).Infof("refreshZones: refresh zones finished with error: %s", err.Error())
}

func (az *Cloud) syncRegionZonesMap(ctx context.Context) error {
	klog.V(2).Infof("syncRegionZonesMap: starting to fetch all available zones for the subscription %s", az.SubscriptionID)
	zones, err := az.zoneRepo.ListZones(ctx)
	if err != nil {
		return fmt.Errorf("list zones: %w", err)
	}

	if len(zones) == 0 {
		klog.Warning("syncRegionZonesMap: empty zone list")
	}

	az.updateRegionZonesMap(zones)

	return nil
}

func (az *Cloud) updateRegionZonesMap(zones map[string][]string) {
	az.refreshZonesLock.Lock()
	defer az.refreshZonesLock.Unlock()

	if az.regionZonesMap == nil {
		az.regionZonesMap = make(map[string][]string)
	}

	for region, z := range zones {
		az.regionZonesMap[region] = z
	}
}

func (az *Cloud) getRegionZonesBackoff(ctx context.Context, region string) ([]*string, error) {
	if az.IsStackCloud() {
		// Azure Stack does not support zone at the moment
		// https://docs.microsoft.com/en-us/azure-stack/user/azure-stack-network-differences?view=azs-2102
		klog.V(3).Infof("getRegionZonesMapWrapper: Azure Stack does not support Zones at the moment, skipping")
		return to.SliceOfPtrs(az.regionZonesMap[region]...), nil
	}

	if len(az.regionZonesMap) != 0 {
		az.refreshZonesLock.RLock()
		defer az.refreshZonesLock.RUnlock()

		return to.SliceOfPtrs(az.regionZonesMap[region]...), nil
	}

	klog.V(2).Infof("getRegionZonesMapWrapper: the region-zones map is not initialized successfully, retrying immediately")

	var (
		zones    map[string][]string
		innerErr error
	)
	err := wait.ExponentialBackoffWithContext(ctx, az.RequestBackoff(), func(ctx context.Context) (done bool, err error) {
		zones, innerErr = az.zoneRepo.ListZones(ctx)
		if innerErr != nil {
			klog.ErrorS(err, "Failed to list zones")
			return false, nil
		}

		return true, nil
	})

	if wait.Interrupted(err) {
		return []*string{}, innerErr
	}

	az.updateRegionZonesMap(zones)

	if len(az.regionZonesMap) != 0 {
		az.refreshZonesLock.RLock()
		defer az.refreshZonesLock.RUnlock()

		return to.SliceOfPtrs(az.regionZonesMap[region]...), nil
	}

	return []*string{}, nil
}

// makeZone returns the zone value in format of <region>-<zone-id>.
func (az *Cloud) makeZone(location string, zoneID int) string {
	return fmt.Sprintf("%s-%d", strings.ToLower(location), zoneID)
}

// isAvailabilityZone returns true if the zone is in format of <region>-<zone-id>.
func (az *Cloud) isAvailabilityZone(zone string) bool {
	return strings.HasPrefix(zone, fmt.Sprintf("%s-", az.Location))
}

// GetZoneID returns the ID of zone from node's zone label.
func (az *Cloud) GetZoneID(zoneLabel string) string {
	if !az.isAvailabilityZone(zoneLabel) {
		return ""
	}

	return strings.TrimPrefix(zoneLabel, fmt.Sprintf("%s-", az.Location))
}

// GetZone returns the Zone containing the current availability zone and locality region that the program is running in.
// DEPRECATED: Zones is deprecated in favor of retrieving zone/region information from InstancesV2.
// This interface will not be called if InstancesV2 is enabled.
// If the node is not running with availability zones, then it will fall back to fault domain.
func (az *Cloud) GetZone(ctx context.Context) (cloudprovider.Zone, error) {
	if az.UseInstanceMetadata {
		metadata, err := az.Metadata.GetMetadata(ctx, azcache.CacheReadTypeUnsafe)
		if err != nil {
			return cloudprovider.Zone{}, err
		}

		if metadata.Compute == nil {
			_ = az.Metadata.imsCache.Delete(consts.MetadataCacheKey)
			return cloudprovider.Zone{}, fmt.Errorf("failure of getting compute information from instance metadata")
		}

		zone := ""
		location := metadata.Compute.Location
		if metadata.Compute.Zone != "" {
			zoneID, err := strconv.Atoi(metadata.Compute.Zone)
			if err != nil {
				return cloudprovider.Zone{}, fmt.Errorf("failed to parse zone ID %q: %w", metadata.Compute.Zone, err)
			}
			zone = az.makeZone(location, zoneID)
		} else {
			klog.V(3).Infof("Availability zone is not enabled for the node, falling back to fault domain")
			zone = metadata.Compute.FaultDomain
		}

		return cloudprovider.Zone{
			FailureDomain: strings.ToLower(zone),
			Region:        strings.ToLower(location),
		}, nil
	}
	// if UseInstanceMetadata is false, get Zone name by calling ARM
	hostname, err := os.Hostname()
	if err != nil {
		return cloudprovider.Zone{}, fmt.Errorf("failure getting hostname from kernel")
	}
	if az.VMSet == nil {
		return cloudprovider.Zone{}, fmt.Errorf("VMSet is not initialized")
	}
	return az.VMSet.GetZoneByNodeName(ctx, strings.ToLower(hostname))
}

// GetZoneByProviderID implements Zones.GetZoneByProviderID
// This is particularly useful in external cloud providers where the kubelet
// does not initialize node data.
// DEPRECATED: Zones is deprecated in favor of retrieving zone/region information from InstancesV2.
// This interface will not be called if InstancesV2 is enabled.
func (az *Cloud) GetZoneByProviderID(ctx context.Context, providerID string) (cloudprovider.Zone, error) {
	if providerID == "" {
		return cloudprovider.Zone{}, errNodeNotInitialized
	}

	// Returns nil for unmanaged nodes because azure cloud provider couldn't fetch information for them.
	if az.IsNodeUnmanagedByProviderID(providerID) {
		klog.V(2).Infof("GetZoneByProviderID: omitting unmanaged node %q", providerID)
		return cloudprovider.Zone{}, nil
	}

	nodeName, err := az.VMSet.GetNodeNameByProviderID(ctx, providerID)
	if err != nil {
		return cloudprovider.Zone{}, err
	}

	return az.GetZoneByNodeName(ctx, nodeName)
}

// GetZoneByNodeName implements Zones.GetZoneByNodeName
// This is particularly useful in external cloud providers where the kubelet
// does not initialize node data.
// DEPRECATED: Zones is deprecated in favor of retrieving zone/region information from InstancesV2.
// This interface will not be called if InstancesV2 is enabled.
func (az *Cloud) GetZoneByNodeName(ctx context.Context, nodeName types.NodeName) (cloudprovider.Zone, error) {
	// Returns "" for unmanaged nodes because azure cloud provider couldn't fetch information for them.
	unmanaged, err := az.IsNodeUnmanaged(string(nodeName))
	if err != nil {
		return cloudprovider.Zone{}, err
	}
	if unmanaged {
		klog.V(2).Infof("GetZoneByNodeName: omitting unmanaged node %q", nodeName)
		return cloudprovider.Zone{}, nil
	}

	return az.VMSet.GetZoneByNodeName(ctx, string(nodeName))
}
