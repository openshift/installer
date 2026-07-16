package postdestroy

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/groups"
	"github.com/sirupsen/logrus"

	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
)

// SecurityGroups deletes the bootstrap security group which is no longer
// required after bootstrapping is complete. It first cleans up any
// orphaned bootstrap ports that still reference the security group,
// since the ORC port controller may have been terminated before it
// could delete them.
func SecurityGroups(ctx context.Context, cloud string, infraID string) error {
	clusterTag := fmt.Sprintf("openshiftClusterID=%s", infraID)
	logrus.Debugf("Searching for bootstrap security group with tags: %s, %s", clusterTag, roleTag)

	networkClient, err := openstackdefaults.NewServiceClient(ctx, "network", openstackdefaults.DefaultClientOpts(cloud))
	if err != nil {
		return fmt.Errorf("failed to create network client: %w", err)
	}

	// Search for security groups with both cluster ID and role tags
	// OpenStack tags filter requires ALL specified tags to match (AND logic)
	allPages, err := groups.List(networkClient, groups.ListOpts{
		Tags: fmt.Sprintf("%s,%s", clusterTag, roleTag),
	}).AllPages(ctx)
	if err != nil {
		return fmt.Errorf("failed to list security groups: %w", err)
	}

	secGroups, err := groups.ExtractGroups(allPages)
	if err != nil {
		return fmt.Errorf("failed to extract security groups: %w", err)
	}

	if len(secGroups) == 0 {
		logrus.Debug("No bootstrap security group found (may have already been deleted)")
		return nil
	}

	// Should only find one security group with both tags, but handle multiple.
	var errs []error
	for _, secGroup := range secGroups {
		// Delete any orphaned ports still referencing this security
		// group. The ORC port controller may have been killed when
		// envtest was torn down, leaving the Neutron port behind.
		if err := deletePorts(ctx, networkClient, secGroup.ID); err != nil {
			errs = append(errs, err)
			continue
		}

		logrus.Infof("Deleting bootstrap security group %s (ID: %s)", secGroup.Name, secGroup.ID)

		err = groups.Delete(ctx, networkClient, secGroup.ID).ExtractErr()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				logrus.Debugf("Bootstrap security group %s already deleted", secGroup.ID)
				continue
			}
			logrus.Errorf("Failed to delete bootstrap security group %s: %v", secGroup.Name, err)
			errs = append(errs, fmt.Errorf("failed to delete bootstrap security group %s: %w", secGroup.Name, err))
			continue
		}

		logrus.Infof("Successfully deleted bootstrap security group %s", secGroup.Name)
	}

	return errors.Join(errs...)
}
