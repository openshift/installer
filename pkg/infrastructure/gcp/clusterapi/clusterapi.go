package clusterapi

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	capg "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openshift/installer/pkg/asset/cluster/tfvars"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap/gcp"
	icgcp "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	"github.com/openshift/installer/pkg/types"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

// Provider implements gcp infrastructure in conjunction with the
// GCP CAPI provider.
type Provider struct {
}

var _ clusterapi.PreProvider = (*Provider)(nil)
var _ clusterapi.IgnitionProvider = (*Provider)(nil)
var _ clusterapi.InfraReadyProvider = (*Provider)(nil)
var _ clusterapi.PostProvider = (*Provider)(nil)
var _ clusterapi.BootstrapDestroyer = (*Provider)(nil)

// Name returns the name for the platform.
func (p Provider) Name() string {
	return gcptypes.Name
}

// BootstrapHasPublicIP indicates that machine ready checks
// should wait for an ExternalIP in the status.
func (Provider) BootstrapHasPublicIP() bool { return true }

// PreProvision is called before provisioning using CAPI controllers has initiated.
// GCP resources that are not created by CAPG (and are required for other stages of the install) are
// created here using the gcp sdk.
func (p Provider) PreProvision(ctx context.Context, in clusterapi.PreProvisionInput) error {
	// Create ServiceAccounts which will be used for machines
	projectID := in.InstallConfig.Config.Platform.GCP.ProjectID

	// ServiceAccount for masters
	// Only create ServiceAccount for masters if a shared VPC install is not being done
	if len(in.InstallConfig.Config.Platform.GCP.NetworkProjectID) == 0 ||
		(in.InstallConfig.Config.ControlPlane != nil &&
			in.InstallConfig.Config.ControlPlane.Platform.GCP != nil &&
			in.InstallConfig.Config.ControlPlane.Platform.GCP.ServiceAccount == "") {
		masterSA, err := CreateServiceAccount(ctx, in.InfraID, projectID, "master")
		if err != nil {
			return fmt.Errorf("failed to create master serviceAccount: %w", err)
		}
		if err = AddServiceAccountRoles(ctx, projectID, masterSA, GetMasterRoles()); err != nil {
			return fmt.Errorf("failed to add master roles: %w", err)
		}
	}

	// ServiceAccount for workers
	workerSA, err := CreateServiceAccount(ctx, in.InfraID, projectID, "worker")
	if err != nil {
		return fmt.Errorf("failed to create worker serviceAccount: %w", err)
	}
	if err = AddServiceAccountRoles(ctx, projectID, workerSA, GetWorkerRoles()); err != nil {
		return fmt.Errorf("failed to add worker roles: %w", err)
	}

	return nil
}

// Ignition provisions the GCP bucket and url that points to the bucket. Bootstrap ignition data cannot
// populate the metadata field of the bootstrap instance as the data can be too large. Instead, the data is
// added to a bucket. A signed url is generated to point to the bucket and the ignition data will be
// updated to point to the url. This is also allows for bootstrap data to be edited after its initial creation.
func (p Provider) Ignition(ctx context.Context, in clusterapi.IgnitionInput) ([]byte, error) {
	// Create the bucket and presigned url. The url is generated using a known/expected name so that the
	// url can be retrieved from the api by this name.
	ctx, cancel := context.WithTimeout(ctx, time.Minute*2)
	defer cancel()

	bucketName := gcp.GetBootstrapStorageName(in.InfraID)
	bucketHandle, err := gcp.CreateBucketHandle(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to create bucket handle %s: %w", bucketName, err)
	}

	if err := gcp.CreateStorage(ctx, in.InstallConfig, bucketHandle, in.InfraID); err != nil {
		return nil, fmt.Errorf("failed to create bucket %s: %w", bucketName, err)
	}

	editedIgnitionBytes, err := EditIgnition(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed to edit bootstrap ignition: %w", err)
	}

	ignitionBytes := in.BootstrapIgnData
	if editedIgnitionBytes != nil {
		ignitionBytes = editedIgnitionBytes
	}

	if err := gcp.FillBucket(ctx, bucketHandle, string(ignitionBytes)); err != nil {
		return nil, fmt.Errorf("ignition failed to fill bucket: %w", err)
	}

	for _, file := range in.TFVarsAsset.Files() {
		if file.Filename == tfvars.TfPlatformVarsFileName {
			var found bool
			tfvarsData := make(map[string]interface{})
			err = json.Unmarshal(file.Data, &tfvarsData)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal %s to json: %w", tfvars.TfPlatformVarsFileName, err)
			}

			ignShim, found := tfvarsData["gcp_ignition_shim"].(string)
			if !found {
				return nil, fmt.Errorf("failed to find ignition shim")
			}

			return []byte(ignShim), nil
		}
	}

	return nil, fmt.Errorf("failed to complete ignition process")
}

// InfraReady is called once cluster.Status.InfrastructureReady
// is true, typically after load balancers have been provisioned. It can be used
// to create DNS records.
func (p Provider) InfraReady(ctx context.Context, in clusterapi.InfraReadyInput) error {
	gcpCluster := &capg.GCPCluster{}
	key := client.ObjectKey{
		Name:      in.InfraID,
		Namespace: capiutils.Namespace,
	}
	if err := in.Client.Get(ctx, key, gcpCluster); err != nil {
		return fmt.Errorf("failed to get GCP cluster: %w", err)
	}

	// public load balancer is created by CAPG. The health check for this load balancer is also created by
	// the CAPG.
	apiIPAddress := gcpCluster.Spec.ControlPlaneEndpoint.Host
	if apiIPAddress == "" && in.InstallConfig.Config.Publish == types.ExternalPublishingStrategy {
		logrus.Debugf("publish strategy is set to external but api address is empty")
	}

	if err := createBootstrapFirewallRules(ctx, in, *gcpCluster.Status.Network.SelfLink); err != nil {
		return fmt.Errorf("failed to add bootstrap firewall rule: %w", err)
	}

	client, err := icgcp.NewClient(context.TODO())
	if err != nil {
		return err
	}

	networkProjectID := in.InstallConfig.Config.GCP.NetworkProjectID
	if networkProjectID == "" {
		networkProjectID = in.InstallConfig.Config.GCP.ProjectID
	}

	networkSelfLink := *gcpCluster.Status.Network.SelfLink
	networkName := path.Base(networkSelfLink)
	masterSubnetName := gcptypes.DefaultSubnetName(in.InfraID, "master")
	if in.InstallConfig.Config.GCP.ControlPlaneSubnet != "" {
		masterSubnetName = in.InstallConfig.Config.GCP.ControlPlaneSubnet
	}

	subnets, err := client.GetSubnetworks(context.TODO(), networkName, networkProjectID, in.InstallConfig.Config.GCP.Region)
	if err != nil {
		return fmt.Errorf("failed to retrieve subnets: %w", err)
	}

	var masterSubnetSelflink string
	for _, s := range subnets {
		if strings.EqualFold(s.Name, masterSubnetName) {
			masterSubnetSelflink = s.SelfLink
			break
		}
	}

	if masterSubnetSelflink == "" {
		return fmt.Errorf("could not find master subnet %s in subnets %v", masterSubnetName, subnets)
	}

	// The firewall for masters, aka control-plane, is created by CAPG
	// Create the ones needed for worker to master communication
	if err = createFirewallRules(ctx, in, *gcpCluster.Status.Network.SelfLink); err != nil {
		return fmt.Errorf("failed to add firewall rules: %w", err)
	}

	if in.InstallConfig.Config.GCP.UserProvisionedDNS != gcptypes.UserProvisionedDNSEnabled {
		// Get the network from the GCP Cluster. The network is used to create the private managed zone.
		if gcpCluster.Status.Network.SelfLink == nil {
			return fmt.Errorf("failed to get GCP network: %w", err)
		}

		// Create the private zone if one does not exist
		if err := createPrivateManagedZone(ctx, in.InstallConfig, in.InfraID, *gcpCluster.Status.Network.SelfLink); err != nil {
			return fmt.Errorf("failed to create the private managed zone: %w", err)
		}

		apiIntIPAddress, err := getInternalLBAddress(ctx, in.InstallConfig.Config.GCP.ProjectID, in.InstallConfig.Config.GCP.Region, getAPIAddressName(in.InfraID))
		if err != nil {
			return fmt.Errorf("failed to get the internal load balancer address: %w", err)
		}

		// Create the public (optional) and private dns records
		if err := createDNSRecords(ctx, in.InstallConfig, in.InfraID, apiIPAddress, apiIntIPAddress); err != nil {
			return fmt.Errorf("failed to create DNS records: %w", err)
		}
	}

	return nil
}

// DestroyBootstrap destroys the temporary bootstrap resources.
func (p Provider) DestroyBootstrap(ctx context.Context, in clusterapi.BootstrapDestroyInput) error {
	logrus.Warnf("Destroying GCP Bootstrap Resources")
	if err := gcp.DestroyStorage(ctx, in.Metadata.InfraID); err != nil {
		return fmt.Errorf("failed to destroy storage: %w", err)
	}

	projectID := in.Metadata.GCP.ProjectID
	if in.Metadata.GCP.NetworkProjectID != "" {
		projectID = in.Metadata.GCP.NetworkProjectID
	}
	if err := removeBootstrapFirewallRules(ctx, in.Metadata.InfraID, projectID); err != nil {
		return fmt.Errorf("failed to remove bootstrap firewall rules: %w", err)
	}

	return nil
}

// PostProvision should be called to add or update and GCP resources after provisioning has completed.
func (p Provider) PostProvision(ctx context.Context, in clusterapi.PostProvisionInput) error {
	return nil
}
