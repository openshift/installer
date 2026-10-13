package gcp

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

const (
	// CreateFirewallPermission is the permission to create firewall rules in the google cloud provider.
	CreateFirewallPermission = "compute.firewalls.create" // required to create GCP firewall rules

	// DeleteFirewallPermission is the permission to delete firewall rules in the google cloud provider.
	DeleteFirewallPermission = "compute.firewalls.delete"

	// UpdateNetworksPermission is the permission to update networks and the network resources in the google cloud provider.
	UpdateNetworksPermission = "compute.networks.updatePolicy"
)

// HasPermission determines if the permission exists for the service account in the project.
func HasPermission(ctx context.Context, projectID string, permissions []string, endpoint *gcptypes.PSCEndpoint) (bool, error) {
	client, err := NewClient(ctx, endpoint)
	if err != nil {
		return false, fmt.Errorf("failed to create client permission check: %w", err)
	}

	foundPermissions, err := client.GetProjectPermissions(ctx, projectID, permissions)
	if err != nil {
		return false, fmt.Errorf("failed to find project permissions: %w", err)
	}

	permissionsValid := true
	for _, permission := range permissions {
		if hasPermission := foundPermissions.Has(permission); !hasPermission {
			logrus.Debugf("permission %s not found", permission)
			permissionsValid = false
		}
	}

	return permissionsValid, nil
}
