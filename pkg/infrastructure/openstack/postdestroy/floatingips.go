package postdestroy

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/sirupsen/logrus"

	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
)

// FloatingIPs deletes the bootstrap floating IP that was previously disassociated
// (but not deleted) when the bootstrap machine was destroyed by CAPO.
func FloatingIPs(ctx context.Context, cloud string, infraID string) error {
	clusterTag := fmt.Sprintf("openshiftClusterID=%s", infraID)
	logrus.Debugf("Searching for bootstrap floating IP with tags: %s, %s", clusterTag, roleTag)

	networkClient, err := openstackdefaults.NewServiceClient(ctx, "network", openstackdefaults.DefaultClientOpts(cloud))
	if err != nil {
		return fmt.Errorf("failed to create network client: %w", err)
	}

	// Search for floating IPs with both cluster ID and role tags
	// OpenStack tags filter requires ALL specified tags to match (AND logic)
	allPages, err := floatingips.List(networkClient, floatingips.ListOpts{
		Tags: fmt.Sprintf("%s,%s", clusterTag, roleTag),
	}).AllPages(ctx)
	if err != nil {
		return fmt.Errorf("failed to list floating IPs: %w", err)
	}

	fips, err := floatingips.ExtractFloatingIPs(allPages)
	if err != nil {
		return fmt.Errorf("failed to extract floating IPs: %w", err)
	}

	if len(fips) == 0 {
		logrus.Debug("No bootstrap floating IP found (may have already been deleted)")
		return nil
	}

	// Should only find one FIP with both tags, but delete all if multiple exist
	var errs []error
	for _, fip := range fips {
		logrus.Infof("Deleting bootstrap floating IP %s (ID: %s)", fip.FloatingIP, fip.ID)

		err = floatingips.Delete(ctx, networkClient, fip.ID).ExtractErr()
		if err != nil {
			// Check if it's a "not found" error, which is acceptable
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				logrus.Debugf("Bootstrap floating IP %s already deleted", fip.ID)
				continue
			}
			logrus.Errorf("Failed to delete bootstrap floating IP %s: %v", fip.FloatingIP, err)
			errs = append(errs, fmt.Errorf("failed to delete bootstrap floating IP %s: %w", fip.FloatingIP, err))
			continue
		}

		logrus.Infof("Successfully deleted bootstrap floating IP %s", fip.FloatingIP)
	}

	return errors.Join(errs...)
}
