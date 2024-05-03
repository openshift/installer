package clusterapi

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"

	"github.com/openshift/installer/pkg/asset/manifests/gcp"
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

// createInternalLB creates a static ip address for the internal load balancer.
// Returns the IP address of the created load balancer.
func createInternalLB(ctx context.Context, in clusterapi.InfraReadyInput, subnetSelfLink, networkSelfLink string, zones []*string) (string, error) {
	projectID := in.InstallConfig.Config.GCP.ProjectID
	region := in.InstallConfig.Config.GCP.Region
	name := getAPIAddressName(in.InfraID)
	labels := mergeLabels(in.InstallConfig, in.InfraID)

	service, err := NewComputeService()
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Minute*3)
	defer cancel()

	// Patch the balancing mode on CAPG proxy classic load balancer backends
	// to match the CONNECTION balancing mode used by installer-created
	// internal passthrough LB, because:
	// "all backend services that reference the instance group must use the same balancing mode"
	// cf: https://cloud.google.com/load-balancing/docs/backend-service
	logrus.Debug("Patching external load balancer")
	extBesvcName := fmt.Sprintf("%s-apiserver", in.InfraID)
	extBesvc, err := service.BackendServices.Get(projectID, extBesvcName).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to get backend service: %w", err)
	}

	for _, be := range extBesvc.Backends {
		be.BalancingMode = "CONNECTION"
		be.MaxConnections = int64(2 ^ 32)
	}

	op, err := service.BackendServices.Patch(projectID, extBesvcName, extBesvc).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to patch external load balancer: %w", err)
	}

	if err := WaitForOperationGlobal(ctx, projectID, op); err != nil {
		return "", fmt.Errorf("failed to wait for patching external load balancer: %w", err)
	}
	logrus.Debug("Successfully patched external load balancer")

	logrus.Debug("Creating internal load balancer")
	addr := &compute.Address{
		Name:        name,
		AddressType: "INTERNAL",
		Subnetwork:  subnetSelfLink,
		Description: resourceDescription,
		Labels:      labels,
		Region:      region,
	}

	op, err = service.Addresses.Insert(projectID, region, addr).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to create internal compute address: %w", err)
	}

	if err := WaitForOperationRegional(ctx, projectID, region, op); err != nil {
		return "", fmt.Errorf("failed to wait for compute address creation: %w", err)
	}

	ipAddress, err := getInternalLBAddress(ctx, projectID, region, name)
	if err != nil {
		return "", fmt.Errorf("failed to get internal load balancer IP address: %w", err)
	}

	hcName := getAPIInternalResourceName(in.InfraID)
	healthCheck := &compute.HealthCheck{
		Region:             region,
		Name:               hcName,
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

	_, err = service.RegionHealthChecks.Insert(projectID, region, healthCheck).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to create api-internal health check: %w", err)
	}

	if err := WaitForOperationRegional(ctx, projectID, region, op); err != nil {
		return "", fmt.Errorf("failed to wait for health check creation: %w", err)
	}

	hc, err := service.RegionHealthChecks.Get(projectID, region, hcName).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("error getting health check: %w", err)
	}
	backends := []*compute.Backend{}
	for _, zone := range zones {
		igName := fmt.Sprintf("%s-%s-%s", in.InfraID, gcp.InstanceGroupRoleTag, *zone)
		ig, err := service.InstanceGroups.Get(projectID, *zone, igName).Context(ctx).Do()
		if err != nil {
			return "", fmt.Errorf("error getting instance group %s in zone %s: %w", igName, *zone, err)
		}
		backends = append(backends, &compute.Backend{
			BalancingMode: "CONNECTION",
			Group:         ig.SelfLink,
		})
	}

	besvcName := fmt.Sprintf("%s-api-internal", in.InfraID)
	op, err = service.RegionBackendServices.Insert(projectID, region, &compute.BackendService{
		Backends:            backends,
		Name:                besvcName,
		LoadBalancingScheme: "INTERNAL",
		Protocol:            "TCP",
		TimeoutSec:          int64((10 * time.Minute).Seconds()),
		HealthChecks:        []string{hc.SelfLink},
		Region:              region,
		Network:             networkSelfLink,
	}).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to create internal backend service: %w", err)
	}

	if err := WaitForOperationRegional(ctx, projectID, region, op); err != nil {
		return "", fmt.Errorf("failed to wait for internal backend service creation: %w", err)
	}

	besvc, err := service.RegionBackendServices.Get(projectID, region, besvcName).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to get backend service: %w", err)
	}

	op, err = service.ForwardingRules.Insert(projectID, region, &compute.ForwardingRule{
		Name:                fmt.Sprintf("%s-api-internal", in.InfraID),
		IPProtocol:          "TCP",
		IPAddress:           ipAddress,
		LoadBalancingScheme: "INTERNAL",
		Ports:               []string{"6443", "22623"},
		BackendService:      besvc.SelfLink,
		Network:             networkSelfLink,
		Subnetwork:          subnetSelfLink,
		Region:              region,
		Labels:              labels,
	}).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to create forwarding rule: %w", err)
	}

	if err := WaitForOperationRegional(ctx, projectID, region, op); err != nil {
		return "", fmt.Errorf("failed to wait for forwarding rule creation: %w", err)
	}
	logrus.Debug("Successfully created internal load balancer")
	return ipAddress, nil
}
