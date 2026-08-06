package clusterapi

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"

	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

const (
	// Health check names follow CAPG naming: {infraID}-{lbname}.
	apiServerHealthCheckSuffix   = "apiserver"
	apiInternalHealthCheckSuffix = "api-internal"

	// Target health check values matching legacy Terraform IPI.
	// Worst-case unhealthy detection: interval * threshold = 2s * 3 = 6s,
	// well within the 30s requirement from docs/dev/kube-apiserver-health-check.md.
	hcCheckIntervalSec   = int64(2)
	hcTimeoutSec         = int64(2)
	hcHealthyThreshold   = int64(3)
	hcUnhealthyThreshold = int64(3)
)

// updateHealthChecks patches the load balancer health checks created by CAPG
// to use faster timing values that match the legacy Terraform IPI configuration.
func updateHealthChecks(ctx context.Context, in clusterapi.InfraReadyInput) error {
	projectID := in.InstallConfig.Config.GCP.ProjectID
	region := in.InstallConfig.Config.GCP.Region

	opts := []option.ClientOption{option.WithScopes(compute.CloudPlatformScope)}
	pscEndpoint := in.InstallConfig.Config.GCP.Endpoint
	if gcptypes.ShouldUseEndpointForInstaller(pscEndpoint) {
		opts = append(opts, gcpconfig.CreateEndpointOption(pscEndpoint.Name, gcpconfig.ServiceNameGCPCompute))
	}
	svc, err := gcpconfig.GetComputeService(ctx, opts...)
	if err != nil {
		return fmt.Errorf("failed to create compute service for health check updates: %w", err)
	}

	hcPatch := &compute.HealthCheck{
		CheckIntervalSec:   hcCheckIntervalSec,
		TimeoutSec:         hcTimeoutSec,
		HealthyThreshold:   hcHealthyThreshold,
		UnhealthyThreshold: hcUnhealthyThreshold,
	}

	if in.InstallConfig.Config.PublicAPI() {
		hcName := fmt.Sprintf("%s-%s", in.InfraID, apiServerHealthCheckSuffix)
		logrus.Infof("Updating global health check %s", hcName)
		if err := patchGlobalHealthCheck(ctx, svc, projectID, hcName, hcPatch); err != nil {
			return fmt.Errorf("failed to update global health check %s: %w", hcName, err)
		}
	}

	hcName := fmt.Sprintf("%s-%s", in.InfraID, apiInternalHealthCheckSuffix)
	logrus.Infof("Updating regional health check %s", hcName)
	if err := patchRegionalHealthCheck(ctx, svc, projectID, region, hcName, hcPatch); err != nil {
		return fmt.Errorf("failed to update regional health check %s: %w", hcName, err)
	}

	return nil
}

func patchGlobalHealthCheck(ctx context.Context, svc *compute.Service, projectID, hcName string, hcPatch *compute.HealthCheck) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*3)
	defer cancel()

	op, err := svc.HealthChecks.Patch(projectID, hcName, hcPatch).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to patch global health check %s: %w", hcName, err)
	}

	if err := WaitForOperationGlobal(ctx, svc, projectID, op); err != nil {
		return fmt.Errorf("failed waiting for global health check patch %s: %w", hcName, err)
	}

	return nil
}

func patchRegionalHealthCheck(ctx context.Context, svc *compute.Service, projectID, region, hcName string, hcPatch *compute.HealthCheck) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*3)
	defer cancel()

	op, err := svc.RegionHealthChecks.Patch(projectID, region, hcName, hcPatch).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to patch regional health check %s: %w", hcName, err)
	}

	if err := WaitForOperationRegional(ctx, svc, projectID, region, op); err != nil {
		return fmt.Errorf("failed waiting for regional health check patch %s: %w", hcName, err)
	}

	return nil
}
