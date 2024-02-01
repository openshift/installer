package clusterapi

import (
	"context"
	"fmt"
	"github.com/openshift/installer/pkg/types"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"
	"time"

	"github.com/openshift/installer/pkg/asset/installconfig"
)

type loadBalancerType string

const (
	externalLoadBalancer loadBalancerType = "EXTERNAL"
	internalLoadBalancer loadBalancerType = "INTERNAL"

	internalLoadBalancerPort = 6443
	externalLoadBalancerPort = 6080
)

var (
	apiStrLookup = map[loadBalancerType]string{
		externalLoadBalancer: "api",
		internalLoadBalancer: "api-internal",
	}
)

func createHealthCheck(ctx context.Context, ic *installconfig.InstallConfig, clusterID string, lbType loadBalancerType, port int64) error {
	healthCheck := &compute.HealthCheck{
		Name:               fmt.Sprintf("%s-%s", clusterID, apiStrLookup[lbType]),
		Description:        "Created By OpenShift Installer",
		HealthyThreshold:   3,
		UnhealthyThreshold: 3,
		CheckIntervalSec:   2,
		TimeoutSec:         2,
		HttpHealthCheck: &compute.HTTPHealthCheck{
			Port:        port,
			RequestPath: "/readyz",
		},
	}

	ctx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	// TODO: Look into scopes by using options or opts below
	service, err := compute.NewService(ctx)
	if err != nil {
		return fmt.Errorf("failed to created compute service: %w", err)
	}

	// TODO: Do we need a child context ?
	if _, err := service.HealthChecks.Insert(ic.Config.GCP.ProjectID, healthCheck).Context(ctx).Do(); err != nil {
		return fmt.Errorf("failed to create %s health check: %w", apiStrLookup[lbType], err)
	}

	return nil
}

// createLoadBalancerAddress creates a static ip address for the load balancer.
func createLoadBalancerAddress(ctx context.Context, ic *installconfig.InstallConfig, clusterID string, subnetSelfLink string, lbType loadBalancerType) (string, error) {
	if lbType == externalLoadBalancer && ic.Config.Publish != types.ExternalPublishingStrategy {
		return "", nil
	}

	name := fmt.Sprintf("%s-cluster-ip", clusterID)
	if lbType == internalLoadBalancer {
		name = fmt.Sprintf("%s-cluster-public-ip", clusterID)
	}

	labels := mergeLabels(ic, clusterID)

	// TODO: the subnet is only relevant for internal load balancer ???
	addr := &compute.Address{
		Name:        name,
		AddressType: string(lbType),
		Subnetwork:  subnetSelfLink,
		Description: "Created By OpenShift Installer",
		Labels:      labels,
		Region:      ic.Config.GCP.Region,
	}

	ctx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	// TODO: Look into scopes by using options or opts below
	service, err := compute.NewService(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to created compute service: %w", err)
	}

	// TODO: Do we need a child context ?
	opt, err := service.Addresses.Insert(ic.Config.GCP.ProjectID, ic.Config.GCP.Region, addr).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to create %s compute address: %w", lbType, err)
	}

	logrus.Warnf("output for creating the gcp static address: %s", opt.Status)

	ipAddress := ""
	if opt.HTTPStatusCode == 200 && opt.Status == "DONE" {
		// TODO: figure out if we need to wait here ...
		addrOutput, err := service.Addresses.Get(ic.Config.GCP.ProjectID, ic.Config.GCP.Region, name).Context(ctx).Do()
		if err != nil {
			return "", fmt.Errorf("failed to get compute address %s: %w", name, err)
		}
		ipAddress = addrOutput.Address
	}

	return ipAddress, nil
}
