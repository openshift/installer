package clusterapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	capg "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	"github.com/openshift/installer/pkg/types/dns"
	"github.com/openshift/installer/pkg/types/gcp"
)

// editIgnition attempts to edit the contents of the bootstrap ignition when the user has selected
// a custom DNS configuration. Find the public and private load balancer addresses and fill in the
// infrastructure file within the ignition struct.
func editIgnition(ctx context.Context, in clusterapi.IgnitionInput) (*clusterapi.IgnitionOutput, error) {
	if in.InstallConfig.Config.GCP.UserProvisionedDNS != dns.UserProvisionedDNSEnabled {
		return &clusterapi.IgnitionOutput{
			UpdatedBootstrapIgn: in.BootstrapIgnData,
			UpdatedMasterIgn:    in.MasterIgnData,
			UpdatedWorkerIgn:    in.WorkerIgnData}, nil
	}

	gcpCluster := &capg.GCPCluster{}
	key := client.ObjectKey{
		Name:      in.InfraID,
		Namespace: capiutils.Namespace,
	}
	if err := in.Client.Get(ctx, key, gcpCluster); err != nil {
		return nil, fmt.Errorf("failed to get GCP cluster: %w", err)
	}

	opts := []option.ClientOption{}
	pscEndpoint := in.InstallConfig.Config.GCP.Endpoint
	if gcp.ShouldUseEndpointForInstaller(pscEndpoint) {
		opts = append(opts, gcpconfig.CreateEndpointOption(pscEndpoint.Name, gcpconfig.ServiceNameGCPCompute))
	}
	svc, err := gcpconfig.GetComputeService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	project := in.InstallConfig.Config.GCP.ProjectID

	apiIntIPAddress := *gcpCluster.Status.Network.APIInternalAddress
	addressIntCut := apiIntIPAddress[strings.LastIndex(apiIntIPAddress, "/")+1:]
	// The LoadBalancer's IP address which is part of the frontend configuration, would be located in
	// the service project even in the case of Shared VPC (XPN) installs.
	computeIntAddressObj, err := svc.Addresses.Get(project, in.InstallConfig.Config.GCP.Region, addressIntCut).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get compute address: %w", err)
	}
	computeIntAddress := computeIntAddressObj.Address

	computeAddress := ""
	if in.InstallConfig.Config.PublicAPI() {
		apiIPAddress := *gcpCluster.Status.Network.APIServerAddress
		addressCut := apiIPAddress[strings.LastIndex(apiIPAddress, "/")+1:]
		computeAddressObj, err := svc.GlobalAddresses.Get(project, addressCut).Context(ctx).Do()
		if err != nil {
			return nil, fmt.Errorf("failed to get global compute address: %w", err)
		}

		computeAddress = computeAddressObj.Address
	} else {
		// In private clusters, the API and API-Int servers both point to the same internal load balancer
		computeAddress = computeIntAddressObj.Address
	}

	logrus.Debugf("GCP: Editing Ignition files to start in-cluster DNS when UserProvisionedDNS is enabled")
	return clusterapi.EditIgnition(in, gcp.Name, []string{computeAddress}, []string{computeIntAddress})
}
