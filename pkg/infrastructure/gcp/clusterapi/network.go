package clusterapi

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/compute/v1"

	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
)

func getAPIInternalResourceName(infraID string) string {
	return fmt.Sprintf("%s-api-internal", infraID)
}

func getAPIAddressName(infraID string) string {
	return fmt.Sprintf("%s-cluster-ip", infraID)
}

func getInternalLBAddress(ctx context.Context, project, region, name string) (string, error) {
	service, err := NewComputeService()
	if err != nil {
		return "", err
	}

	addrOutput, err := service.Addresses.Get(project, region, name).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to get compute address %s: %w", name, err)
	}
	return addrOutput.Address, nil
}

// createInternalLBAddress creates a static ip address for the internal load balancer.
func createInternalLBAddress(ctx context.Context, in clusterapi.InfraReadyInput) (string, error) {
	name := getAPIAddressName(in.InfraID)

	// TODO: Find the self link for a subnet from the in.Client
	// TODO: can we pick one returned from the service.Get() ?
	subnetSelfLink := ""

	// TODO: the subnet is only relevant for internal load balancer
	addr := &compute.Address{
		Name:        name,
		AddressType: "INTERNAL",
		Subnetwork:  subnetSelfLink,
		Description: resourceDescription,
		Labels:      mergeLabels(in.InstallConfig, in.InfraID),
		Region:      in.InstallConfig.Config.GCP.Region,
	}

	service, err := NewComputeService()
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	op, err := service.Addresses.Insert(in.InstallConfig.Config.GCP.ProjectID, in.InstallConfig.Config.GCP.Region, addr).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to create internal compute address: %w", err)
	}

	if err := WaitForOperationRegional(ctx, in.InstallConfig.Config.GCP.ProjectID, in.InstallConfig.Config.GCP.Region, op); err != nil {
		return "", fmt.Errorf("failed to wait for compute address creation: %w", err)
	}

	ipAddress, err := getInternalLBAddress(ctx, in.InstallConfig.Config.GCP.ProjectID, in.InstallConfig.Config.GCP.Region, name)
	if err != nil {
		return "", fmt.Errorf("failed to get internal load balancer IP address: %w", err)
	}

	healthCheck := &compute.HealthCheck{
		Name:               getAPIInternalResourceName(in.InfraID),
		Description:        resourceDescription,
		HealthyThreshold:   3,
		UnhealthyThreshold: 3,
		CheckIntervalSec:   2,
		TimeoutSec:         2,
		Type:               "HTTPS",
		HttpsHealthCheck: &compute.HTTPSHealthCheck{
			Port:        6443,
			RequestPath: "/readyz",
		},
	}

	if _, err := service.HealthChecks.Insert(in.InstallConfig.Config.GCP.ProjectID, healthCheck).Context(ctx).Do(); err != nil {
		return "", fmt.Errorf("failed to create api-internal health check: %w", err)
	}

	return ipAddress, nil
}
