package clusterapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"

	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

// deleteFirewallRule deletes the firewall rule identified by name.
func deleteFirewallRule(ctx context.Context, svc *compute.Service, name, projectID string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*3)
	defer cancel()

	op, err := svc.Firewalls.Delete(projectID, name).Context(ctx).Do()
	if err != nil {
		var gErr *googleapi.Error
		if errors.As(err, &gErr) && gErr.Code == http.StatusNotFound {
			logrus.Debugf("Firewall rule %s already deleted, skipping", name)
			return nil
		}
		return fmt.Errorf("failed to delete %s firewall rule: %w", name, err)
	}

	if err := WaitForOperationGlobal(ctx, svc, projectID, op); err != nil {
		return fmt.Errorf("failed to wait for delete %s firewall rule: %w", name, err)
	}

	return nil
}

// removeBootstrapFirewallRules removes the rules created for the bootstrap node.
func removeBootstrapFirewallRules(ctx context.Context, infraID, projectID string, endpoint *gcptypes.PSCEndpoint) error {
	opts := []option.ClientOption{option.WithScopes(compute.CloudPlatformScope)}
	if gcptypes.ShouldUseEndpointForInstaller(endpoint) {
		opts = append(opts, gcpconfig.CreateEndpointOption(endpoint.Name, gcpconfig.ServiceNameGCPCompute))
	}
	svc, err := gcpconfig.GetComputeService(ctx, opts...)
	if err != nil {
		return err
	}

	firewallName := fmt.Sprintf("%s-bootstrap-in-ssh", infraID)
	if err := deleteFirewallRule(ctx, svc, firewallName, projectID); err != nil {
		return err
	}

	firewallName = fmt.Sprintf("%s-bootstrap-in-konnectivity", infraID)
	return deleteFirewallRule(ctx, svc, firewallName, projectID)
}

// removeCAPGFirewallRules removes the overly permissive firewall rules created by cluster-api-provider-gcp.
func removeCAPGFirewallRules(ctx context.Context, infraID, projectID string, endpoint *gcptypes.PSCEndpoint) error {
	opts := []option.ClientOption{option.WithScopes(compute.CloudPlatformScope)}
	if gcptypes.ShouldUseEndpointForInstaller(endpoint) {
		opts = append(opts, gcpconfig.CreateEndpointOption(endpoint.Name, gcpconfig.ServiceNameGCPCompute))
	}
	svc, err := gcpconfig.GetComputeService(ctx, opts...)
	if err != nil {
		return err
	}

	firewallName := fmt.Sprintf("allow-%s-cluster", infraID)
	if err := deleteFirewallRule(ctx, svc, firewallName, projectID); err != nil {
		return err
	}

	firewallName = fmt.Sprintf("allow-%s-healthchecks", infraID)
	return deleteFirewallRule(ctx, svc, firewallName, projectID)
}
