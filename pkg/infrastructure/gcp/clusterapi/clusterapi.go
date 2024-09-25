package clusterapi

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
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

// PublicGatherEndpoint indicates that machine ready checks should wait for an ExternalIP
// in the status and use that when gathering bootstrap log bundles.
func (Provider) PublicGatherEndpoint() clusterapi.GatherEndpoint { return clusterapi.ExternalIP }

// PreProvision is called before provisioning using CAPI controllers has initiated.
// GCP resources that are not created by CAPG (and are required for other stages of the install) are
// created here using the gcp sdk.
func (p Provider) PreProvision(ctx context.Context, in clusterapi.PreProvisionInput) error {
	// Create ServiceAccounts which will be used for machines
	platform := in.InstallConfig.Config.Platform.GCP
	projectID := platform.ProjectID

	// Only create ServiceAccounts for machines if a pre-created Service Account is not defined
	controlPlaneMpool := &gcptypes.MachinePool{}
	controlPlaneMpool.Set(in.InstallConfig.Config.GCP.DefaultMachinePlatform)
	if in.InstallConfig.Config.ControlPlane != nil {
		controlPlaneMpool.Set(in.InstallConfig.Config.ControlPlane.Platform.GCP)
	}

	if controlPlaneMpool.ServiceAccount != "" {
		logrus.Debugf("Using pre-created ServiceAccount for control plane nodes")
	} else {
		// Create ServiceAccount for control plane nodes
		logrus.Debugf("Creating ServiceAccount for control plane nodes")
		masterSA, err := CreateServiceAccount(ctx, in.InfraID, projectID, "master")
		if err != nil {
			return fmt.Errorf("failed to create master serviceAccount: %w", err)
		}
		if err = AddServiceAccountRoles(ctx, projectID, masterSA, GetMasterRoles()); err != nil {
			return fmt.Errorf("failed to add master roles: %w", err)
		}

		// Add additional roles for shared VPC
		if len(in.InstallConfig.Config.Platform.GCP.NetworkProjectID) > 0 {
			projID := in.InstallConfig.Config.Platform.GCP.NetworkProjectID
			// Add roles needed for creating firewalls
			roles := GetSharedVPCRoles()
			if err = AddServiceAccountRoles(ctx, projID, masterSA, roles); err != nil {
				return fmt.Errorf("failed to add roles for shared VPC: %w", err)
			}
		}
	}

	createSA := false
	for _, compute := range in.InstallConfig.Config.Compute {
		computeMpool := compute.Platform.GCP
		if gcptypes.GetConfiguredServiceAccount(platform, computeMpool) == "" {
			// If any compute nodes aren't using defined service account then create the service account
			createSA = true
		}
	}
	if createSA {
		// Create ServiceAccount for workers
		logrus.Debugf("Creating ServiceAccount for compute nodes")
		workerSA, err := CreateServiceAccount(ctx, in.InfraID, projectID, "worker")
		if err != nil {
			return fmt.Errorf("failed to create worker serviceAccount: %w", err)
		}
		if err = AddServiceAccountRoles(ctx, projectID, workerSA, GetWorkerRoles()); err != nil {
			return fmt.Errorf("failed to add worker roles: %w", err)
		}
	}

	return nil
}

// Ignition provisions the GCP bucket and url that points to the bucket. Bootstrap ignition data cannot
// populate the metadata field of the bootstrap instance as the data can be too large. Instead, the data is
// added to a bucket. A signed url is generated to point to the bucket and the ignition data will be
// updated to point to the url. This is also allows for bootstrap data to be edited after its initial creation.
func (p Provider) Ignition(ctx context.Context, in clusterapi.IgnitionInput) ([]*corev1.Secret, error) {
	// Create the bucket and presigned url. The url is generated using a known/expected name so that the
	// url can be retrieved from the api by this name.
	bucketName := gcp.GetBootstrapStorageName(in.InfraID)
	storageClient, err := gcp.NewStorageClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage client: %w", err)
	}

	bucketHandle := storageClient.Bucket(bucketName)
	if err := gcp.CreateStorage(ctx, in.InstallConfig, bucketHandle, in.InfraID); err != nil {
		return nil, fmt.Errorf("failed to create bucket %s: %w", bucketName, err)
	}

	editedIgnitionBytes, err := editIgnition(ctx, in)
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

	var ignShim string
	for _, file := range in.TFVarsAsset.Files() {
		if file.Filename == tfvars.TfPlatformVarsFileName {
			var found bool
			tfvarsData := make(map[string]interface{})
			err = json.Unmarshal(file.Data, &tfvarsData)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal %s to json: %w", tfvars.TfPlatformVarsFileName, err)
			}

			ignShim, found = tfvarsData["gcp_ignition_shim"].(string)
			if !found {
				return nil, fmt.Errorf("failed to find ignition shim")
			}
		}
	}

	ignSecrets := []*corev1.Secret{
		clusterapi.IgnitionSecret([]byte(ignShim), in.InfraID, "bootstrap"),
		clusterapi.IgnitionSecret(in.MasterIgnData, in.InfraID, "master"),
	}

	return ignSecrets, nil
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

		createFwRules, err := hasFirewallPermission(ctx, projectID)
		if err != nil {
			return fmt.Errorf("failed to remove bootstrap firewall rules: %w", err)
		}
		if !createFwRules {
			return nil
		}
	}
	if err := removeBootstrapFirewallRules(ctx, in.Metadata.InfraID, projectID); err != nil {
		return fmt.Errorf("failed to remove bootstrap firewall rules: %w", err)
	}

	if in.Metadata.GCP.NetworkProjectID == "" {
		// Remove the overly permissive firewall rules created by CAPG that are redundant with those created by installer
		// These are not created in a shared VPC installation
		logrus.Infof("Removing firewall rules created by cluster-api-provider-gcp")
		if err := removeCAPGFirewallRules(ctx, in.Metadata.InfraID, in.Metadata.GCP.ProjectID); err != nil {
			return fmt.Errorf("failed to remove firewall rules created by cluster-api-provider-gcp: %w", err)
		}
	}

	return nil
}

// PostProvision should be called to add or update and GCP resources after provisioning has completed.
func (p Provider) PostProvision(ctx context.Context, in clusterapi.PostProvisionInput) error {
	return nil
}
