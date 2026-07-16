package postdestroy

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	"github.com/sirupsen/logrus"
)

// deletePorts deletes all ports that reference the given security group.
// This cleans up orphaned bootstrap ports that the ORC controller could
// not delete because envtest was torn down first.
func deletePorts(ctx context.Context, networkClient *gophercloud.ServiceClient, secGroupID string) error {
	allPages, err := ports.List(networkClient, ports.ListOpts{
		SecurityGroups: []string{secGroupID},
	}).AllPages(ctx)
	if err != nil {
		return fmt.Errorf("failed to list ports for security group %s: %w", secGroupID, err)
	}

	portList, err := ports.ExtractPorts(allPages)
	if err != nil {
		return fmt.Errorf("failed to extract ports: %w", err)
	}

	if len(portList) == 0 {
		return nil
	}

	var errs []error
	for _, port := range portList {
		logrus.Infof("Deleting orphaned bootstrap port %s (ID: %s)", port.Name, port.ID)

		err = ports.Delete(ctx, networkClient, port.ID).ExtractErr()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				logrus.Debugf("Bootstrap port %s already deleted", port.ID)
				continue
			}
			logrus.Errorf("Failed to delete bootstrap port %s: %v", port.Name, err)
			errs = append(errs, fmt.Errorf("failed to delete bootstrap port %s: %w", port.Name, err))
			continue
		}

		logrus.Infof("Successfully deleted bootstrap port %s", port.Name)
	}

	return errors.Join(errs...)
}
