package azure

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/dns"
)

// editIgnition attempts to edit the contents of the bootstrap ignition when the user has selected
// a custom DNS configuration. Find the public and private load balancer addresses and fill in the
// infrastructure file within the ignition struct.
func editIgnition(ctx context.Context, in clusterapi.IgnitionInput, publicIP string) (*clusterapi.IgnitionOutput, error) {
	// ARO wants the ability to enable custom-dns on day-2. In that case, we might have to
	// add LB IPs to Infra CR and within bootstrap Ignition even when `UserProvisionedDNS` is
	// not enabled in install-config.
	if in.InstallConfig.Config.Azure.UserProvisionedDNS != dns.UserProvisionedDNSEnabled {
		return &clusterapi.IgnitionOutput{
			UpdatedBootstrapIgn: in.BootstrapIgnData,
			UpdatedMasterIgn:    in.MasterIgnData,
			UpdatedWorkerIgn:    in.WorkerIgnData}, nil
	}
	logrus.Debugf("Azure: Editing Ignition files to start in-cluster DNS when UserProvisionedDNS is enabled")
	azureCluster := &capz.AzureCluster{}
	key := client.ObjectKey{
		Name:      in.InfraID,
		Namespace: capiutils.Namespace,
	}
	if err := in.Client.Get(ctx, key, azureCluster); err != nil {
		return nil, fmt.Errorf("failed to get Azure cluster: %w", err)
	}
	if apiLB := azureCluster.Spec.NetworkSpec.APIServerLB; apiLB == nil || len(apiLB.FrontendIPs) == 0 {
		return nil, fmt.Errorf("failed to get Azure cluster LB frontend IPs")
	}

	apiIntLBIP := azureCluster.Spec.NetworkSpec.APIServerLB.FrontendIPs[0].PrivateIPAddress
	if apiIntLBIP == "" {
		return nil, fmt.Errorf("failed to get Azure cluster API Server Internal LB IP")
	}
	apiLBIP := apiIntLBIP
	// Update API LB IP for public clusters
	if in.InstallConfig.Config.PublicAPI() && publicIP != "" {
		apiLBIP = publicIP
	}
	logrus.Debugf("Azure: Editing Ignition files with API LB IP: %s and API Int LB IP: %s", apiLBIP, apiIntLBIP)
	return clusterapi.EditIgnition(in, azure.Name, []string{apiLBIP}, []string{apiIntLBIP})
}
