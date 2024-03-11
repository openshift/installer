package clusterapi

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/compute/v1"

	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types/gcp"
)

// createInternalLoadBalancer creates the backend service and forwarding rule for the internal load balancer.
func createInternalLoadBalancer(ctx context.Context, installConfig *installconfig.InstallConfig, clusterID string) error {
	// Adding time, because there are several operations here.
	ctx, cancel := context.WithTimeout(ctx, time.Minute*3)
	defer cancel()

	service, err := NewComputeService()
	if err != nil {
		return fmt.Errorf("create internal load balancer failed to create compute service: %w", err)
	}

	apiAddressName := getAPIAddressName(clusterID)
	healthCheck, err := service.HealthChecks.Get(installConfig.Config.GCP.ProjectID, apiAddressName).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to find health check %s: %w", apiAddressName, err)
	}

	backendService := &compute.BackendService{
		Name:                fmt.Sprintf("%s-api-internal", clusterID),
		Description:         resourceDescription,
		LoadBalancingScheme: "INTERNAL", // Internal pass through network load balancer
		Protocol:            "TCP",
		TimeoutSec:          120,
		Region:              installConfig.Config.GCP.Region,
		Backends:            []*compute.Backend{}, // TODO: backends needs to be figured out
		HealthChecks:        []string{healthCheck.SelfLink},
	}

	backendServiceOp, err := service.BackendServices.Insert(installConfig.Config.GCP.ProjectID, backendService).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to create internal backend service: %w", err)
	}
	if err := WaitForOperationRegional(ctx, installConfig.Config.GCP.ProjectID, installConfig.Config.GCP.Region, backendServiceOp); err != nil {
		return fmt.Errorf("failed to wait for internal backend service creation: %w", err)
	}

	returnedBackendService, err := service.RegionBackendServices.Get(installConfig.Config.GCP.ProjectID, installConfig.Config.GCP.Region, backendService.Name).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to get the regional backend service %s: %w", backendService.Name, err)
	}

	ipAddress, err := getInternalLBAddress(ctx, installConfig.Config.GCP.ProjectID, installConfig.Config.GCP.Region, apiAddressName)
	if err != nil {
		return fmt.Errorf("failed to get internal load balancer IP address: %w", err)
	}

	networkName := fmt.Sprintf("%s-network", clusterID)
	if installConfig.Config.GCP.Network != "" {
		networkName = installConfig.Config.GCP.Network
	}

	subnetName := gcp.DefaultSubnetName(clusterID, "master")
	if installConfig.Config.GCP.ControlPlaneSubnet != "" {
		subnetName = installConfig.Config.GCP.ControlPlaneSubnet
	}

	// Create the forwarding rule for the Backend Service
	forwardingRule := &compute.ForwardingRule{
		Name:                fmt.Sprintf("%s-api-internal", clusterID),
		Description:         resourceDescription,
		IPAddress:           ipAddress,
		BackendService:      returnedBackendService.SelfLink,
		Region:              installConfig.Config.GCP.Region,
		Ports:               []string{"6443", "22623"},
		Subnetwork:          subnetName,
		Network:             networkName,
		Labels:              mergeLabels(installConfig, clusterID),
		LoadBalancingScheme: "INTERNAL",
	}

	forwardingRuleOp, err := service.ForwardingRules.Insert(installConfig.Config.GCP.ProjectID, installConfig.Config.GCP.Region, forwardingRule).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to create the forwarding rule for the backend service %s: %w", backendService.Name, err)
	}
	if err := WaitForOperationRegional(ctx, installConfig.Config.GCP.ProjectID, installConfig.Config.GCP.Region, forwardingRuleOp); err != nil {
		return fmt.Errorf("failed to wait for fowarding rule creation: %w", err)
	}

	return nil
}
