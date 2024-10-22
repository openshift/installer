package clusterapi

import (
	"context"
	"fmt"
	"strings"

	capg "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/dns"
	"github.com/openshift/installer/pkg/types/gcp"
)

// editIgnition attempts to edit the contents of the bootstrap ignition when the user has selected
// a custom DNS configuration. Find the public and private load balancer addresses and fill in the
// infrastructure file within the ignition struct.
func editIgnition(ctx context.Context, in clusterapi.IgnitionInput) ([]byte, []byte, error) {
	if in.InstallConfig.Config.GCP.UserProvisionedDNS != dns.UserProvisionedDNSEnabled {
		return in.BootstrapIgnData, in.MasterIgnData, nil
	}

	gcpCluster := &capg.GCPCluster{}
	key := client.ObjectKey{
		Name:      in.InfraID,
		Namespace: capiutils.Namespace,
	}
	if err := in.Client.Get(ctx, key, gcpCluster); err != nil {
		return nil, nil, fmt.Errorf("failed to get GCP cluster: %w", err)
	}

	svc, err := NewComputeService()
	if err != nil {
		return nil, nil, err
	}

	project := in.InstallConfig.Config.GCP.ProjectID
	if in.InstallConfig.Config.GCP.NetworkProjectID != "" {
		project = in.InstallConfig.Config.GCP.NetworkProjectID
	}

	computeAddress := ""
	if in.InstallConfig.Config.Publish == types.ExternalPublishingStrategy {
		apiIPAddress := *gcpCluster.Status.Network.APIServerAddress
		addressCut := apiIPAddress[strings.LastIndex(apiIPAddress, "/")+1:]
		computeAddressObj, err := svc.GlobalAddresses.Get(project, addressCut).Context(ctx).Do()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get global compute address: %w", err)
		}

		computeAddress = computeAddressObj.Address
	}

	apiIntIPAddress := *gcpCluster.Status.Network.APIInternalAddress
	addressIntCut := apiIntIPAddress[strings.LastIndex(apiIntIPAddress, "/")+1:]
	computeIntAddress, err := svc.Addresses.Get(project, in.InstallConfig.Config.GCP.Region, addressIntCut).Context(ctx).Do()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get compute address: %w", err)
	}

	return clusterapi.EditIgnition(in, gcp.Name, []string{computeAddress}, []string{computeIntAddress.Address})
}
