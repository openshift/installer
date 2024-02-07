package clusterapi

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/compute/v1"
	capg "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
)

func getApiInternalResourceName(infraID string) string {
	return fmt.Sprintf("%s-api-internal", infraID)
}

func getApiAddressName(infraID string) string {
	return fmt.Sprintf("%s-cluster-ip", infraID)
}

func NewService() (*compute.Service, error) {
	ctx := context.Background()

	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to created compute service: %w", err)
	}

	return service, nil
}

func createInternalLBBackendService(ctx context.Context, in clusterapi.PostProvisionInput) error {
	service, err := NewService()
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	name := getApiInternalResourceName(in.InfraID)

	ctx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	healthCheck, err := service.HealthChecks.Get(in.InstallConfig.Config.GCP.ProjectID, name).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to find health check: %w", err)
	}

	address, err := service.Addresses.Get(in.InstallConfig.Config.GCP.ProjectID, in.InstallConfig.Config.GCP.Region, getApiAddressName(in.InfraID)).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to find gcp address: %w", err)
	}

	// TODO: Back ends can/should be retrieved from the provider
	backendService := &compute.BackendService{
		Name:                name,
		Description:         "Created By OpenShift Installer",
		LoadBalancingScheme: "INTERNAL",
		Protocol:            "TCP",
		TimeoutSec:          120,
		HealthChecks:        []string{healthCheck.SelfLink},
		// TODO: create the backends
		Backends: []*compute.Backend{},
	}

	if _, err := service.BackendServices.Insert(in.InstallConfig.Config.GCP.ProjectID, backendService).Context(ctx).Do(); err != nil {
		return fmt.Errorf("failed to create backend service for internal load balancer: %w", err)
	}

	gcpCluster := &capg.GCPCluster{}
	key := client.ObjectKey{
		Name:      in.InfraID,
		Namespace: capiutils.Namespace,
	}
	if err := in.Client.Get(ctx, key, gcpCluster); err != nil {
		return fmt.Errorf("failed to get GCP cluster: %w", err)
	}

	forwardingRule := &compute.ForwardingRule{
		Name:                name,
		Description:         "Created By OpenShift Installer",
		IPAddress:           address.Address,
		BackendService:      "", // TODO: Backend service self link
		Ports:               []string{"6443", "22623"},
		Subnetwork:          "", // TODO: find from the client: master subnet
		Network:             *gcpCluster.Status.Network.SelfLink,
		Labels:              mergeLabels(in.InstallConfig, in.InfraID),
		LoadBalancingScheme: "INTERNAL",
	}

	_, err = service.ForwardingRules.Insert(in.InstallConfig.Config.GCP.ProjectID, in.InstallConfig.Config.GCP.Region, forwardingRule).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to create a forwarding rule for the internal load balancer: %w", err)
	}

	return nil
}

func getInternalLBAddress(ctx context.Context, project, region, name string) (string, error) {
	service, err := NewService()
	if err != nil {
		return "", fmt.Errorf("failed to create compute service: %w", err)
	}

	addrOutput, err := service.Addresses.Get(project, region, name).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to get compute address %s: %w", name, err)
	}
	return addrOutput.Address, nil
}

// createInternalLBAddress creates a static ip address for the internal load balancer.
func createInternalLBAddress(ctx context.Context, in clusterapi.IgnitionInput) (string, error) {
	name := getApiAddressName(in.InfraID)

	// TODO: Find the self link for a subnet from the in.Client
	// TODO: can we pick one returned from the service.Get() ?
	subnetSelfLink := ""

	// TODO: the subnet is only relevant for internal load balancer ???
	addr := &compute.Address{
		Name:        name,
		AddressType: "INTERNAL",
		Subnetwork:  subnetSelfLink,
		Description: "Created By OpenShift Installer",
		Labels:      mergeLabels(in.InstallConfig, in.InfraID),
		Region:      in.InstallConfig.Config.GCP.Region,
	}

	service, err := NewService()
	if err != nil {
		return "", fmt.Errorf("failed to create service: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	// TODO: Do we need a child context ?
	_, err = service.Addresses.Insert(in.InstallConfig.Config.GCP.ProjectID, in.InstallConfig.Config.GCP.Region, addr).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to create internal compute address: %w", err)
	}

	ipAddress, err := getInternalLBAddress(ctx, in.InstallConfig.Config.GCP.ProjectID, in.InstallConfig.Config.GCP.Region, name)
	if err != nil {
		return "", fmt.Errorf("failed to get internal load balancer IP address: %w", err)
	}

	healthCheck := &compute.HealthCheck{
		Name:               getApiInternalResourceName(in.InfraID),
		Description:        "Created By OpenShift Installer",
		HealthyThreshold:   3,
		UnhealthyThreshold: 3,
		CheckIntervalSec:   2,
		TimeoutSec:         2,
		HttpHealthCheck: &compute.HTTPHealthCheck{
			Port:        6443,
			RequestPath: "/readyz",
		},
	}

	if _, err := service.HealthChecks.Insert(in.InstallConfig.Config.GCP.ProjectID, healthCheck).Context(ctx).Do(); err != nil {
		return "", fmt.Errorf("failed to create api-internal health check: %w", err)
	}

	return ipAddress, nil
}
