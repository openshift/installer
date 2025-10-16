package gcp

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

const (
	// CreateGCPFirewallPermission is the role/permission to create or skip the creation of
	// firewall rules for GCP during a xpn installation.
	CreateGCPFirewallPermission = "compute.firewalls.create"

	// DeleteGCPFirewallPermission is the role/permission to delete or skip the delete of
	// firewall rules for GCP during a bootstrap firewall rule deletion during a xpn installation.
	DeleteGCPFirewallPermission = "compute.firewalls.delete"
)

// HasPermissions will check if the user has the permissions in the project.
func HasPermissions(ctx context.Context, projectID string, permissions []string, endpoints *gcptypes.PSCEndpoint) (bool, error) {
	client, err := NewClient(ctx, endpoints)
	if err != nil {
		return false, fmt.Errorf("failed to create client during permission check: %w", err)
	}

	foundPermissions, err := client.GetProjectPermissions(ctx, projectID, permissions)
	if err != nil {
		return false, fmt.Errorf("failed to find project permissions during permission check: %w", err)
	}

	permissionsValid := true
	for _, permission := range permissions {
		if hasPermission := foundPermissions.Has(permission); !hasPermission {
			logrus.Warnf("failed to find permission %s", permission)
			permissionsValid = false
		}
	}

	return permissionsValid, nil
}
