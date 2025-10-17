/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package scope

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-logr/logr"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/globaltaggingv1"
	"github.com/IBM/vpc-go-sdk/vpcv1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"

	"sigs.k8s.io/controller-runtime/pkg/client"

	capiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/patch"

	infrav1beta2 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/authenticator"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/globaltagging"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/utils"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/vpc"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/endpoints"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/options"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/record"
)

// MachineScopeParams defines the input parameters used to create a new MachineScope.
type MachineScopeParams struct {
	IBMVPCClient    vpc.Vpc
	Client          client.Client
	Logger          logr.Logger
	Cluster         *capiv1beta1.Cluster
	Machine         *capiv1beta1.Machine
	IBMVPCCluster   *infrav1beta2.IBMVPCCluster
	IBMVPCMachine   *infrav1beta2.IBMVPCMachine
	ServiceEndpoint []endpoints.ServiceEndpoint
}

// MachineScope defines a scope defined around a machine and its cluster.
type MachineScope struct {
	logr.Logger
	Client      client.Client
	patchHelper *patch.Helper

	IBMVPCClient        vpc.Vpc
	GlobalTaggingClient globaltagging.GlobalTagging
	Cluster             *capiv1beta1.Cluster
	Machine             *capiv1beta1.Machine
	IBMVPCCluster       *infrav1beta2.IBMVPCCluster
	IBMVPCMachine       *infrav1beta2.IBMVPCMachine
	ServiceEndpoint     []endpoints.ServiceEndpoint
}

// NewMachineScope creates a new MachineScope from the supplied parameters.
func NewMachineScope(params MachineScopeParams) (*MachineScope, error) {
	if params.Machine == nil {
		return nil, errors.New("failed to generate new scope from nil Machine")
	}
	if params.IBMVPCMachine == nil {
		return nil, errors.New("failed to generate new scope from nil IBMVPCMachine")
	}

	if params.Logger == (logr.Logger{}) {
		params.Logger = klog.Background()
	}

	helper, err := patch.NewHelper(params.IBMVPCMachine, params.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to init patch helper: %w", err)
	}

	// Fetch the service endpoint.
	svcEndpoint := endpoints.FetchVPCEndpoint(params.IBMVPCCluster.Spec.Region, params.ServiceEndpoint)

	vpcClient, err := vpc.NewService(svcEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to create IBM VPC session: %w", err)
	}

	if params.Logger.V(DEBUGLEVEL).Enabled() {
		core.SetLoggingLevel(core.LevelDebug)
	}

	auth, err := authenticator.GetAuthenticator()
	if err != nil {
		return nil, fmt.Errorf("error failed to create authenticator: %w", err)
	}

	// Create Global Tagging client.
	gtOptions := globaltagging.ServiceOptions{
		GlobalTaggingV1Options: &globaltaggingv1.GlobalTaggingV1Options{
			Authenticator: auth,
		},
	}
	// Override the Global Tagging endpoint if provided.
	if gtEndpoint := endpoints.FetchEndpoints(string(endpoints.GlobalTagging), params.ServiceEndpoint); gtEndpoint != "" {
		gtOptions.URL = gtEndpoint
		params.Logger.Info("Overriding the default global tagging endpoint", "GlobalTaggingEndpoint", gtEndpoint)
	}
	globalTaggingClient, err := globaltagging.NewService(gtOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create global tagging client: %w", err)
	}

	return &MachineScope{
		Logger:              params.Logger,
		Client:              params.Client,
		IBMVPCClient:        vpcClient,
		GlobalTaggingClient: globalTaggingClient,
		Cluster:             params.Cluster,
		IBMVPCCluster:       params.IBMVPCCluster,
		patchHelper:         helper,
		Machine:             params.Machine,
		IBMVPCMachine:       params.IBMVPCMachine,
	}, nil
}

// CreateMachine creates a vpc machine.
func (m *MachineScope) CreateMachine() (*vpcv1.Instance, error) { //nolint: gocyclo
	instanceReply, err := m.ensureInstanceUnique(m.IBMVPCMachine.Name)
	if err != nil {
		return nil, err
	} else if instanceReply != nil {
		// TODO need a reasonable wrapped error.
		return instanceReply, nil
	}

	cloudInitData, err := m.GetBootstrapData()
	if err != nil {
		return nil, err
	}

	options := &vpcv1.CreateInstanceOptions{}
	// Build common field resources, as unique InstancePrototype's are defined based on machine source.
	// TODO(cjschaef): Replace with webhook validation
	if m.IBMVPCMachine.Spec.Profile == "" {
		return nil, fmt.Errorf("error profile is empty for machine %s", m.IBMVPCMachine.Name)
	}
	profile := &vpcv1.InstanceProfileIdentity{
		Name: &m.IBMVPCMachine.Spec.Profile,
	}

	subnetIdentity := &vpcv1.SubnetIdentity{}
	// If Network Status is available, attempt to retrieve subnet ID from there.
	if m.IBMVPCCluster.Status.Network != nil {
		if m.IBMVPCCluster.Status.Network.ControlPlaneSubnets != nil {
			if subnet, ok := m.IBMVPCCluster.Status.Network.ControlPlaneSubnets[m.IBMVPCMachine.Spec.PrimaryNetworkInterface.Subnet]; ok {
				subnetIdentity.ID = ptr.To(subnet.ID)
			}
		}
		if m.IBMVPCCluster.Status.Network.WorkerSubnets != nil {
			if subnet, ok := m.IBMVPCCluster.Status.Network.WorkerSubnets[m.IBMVPCMachine.Spec.PrimaryNetworkInterface.Subnet]; ok {
				subnetIdentity.ID = ptr.To(subnet.ID)
			}
		}
	}
	// If the ID hasn't been set yet, rely on Machine Spec for lookup, and finally falling back to previous logic of using the subnet value directly as an ID.
	if subnetIdentity.ID == nil {
		// For Machines not reliant directly on Cluster managed subnets, lookup subnet ID by name.
		subnetDetails, err := m.IBMVPCClient.GetVPCSubnetByName(m.IBMVPCMachine.Spec.PrimaryNetworkInterface.Subnet)
		if err != nil {
			return nil, fmt.Errorf("error retrieving subnet ID for machine %s: %w", m.IBMVPCMachine.Name, err)
		} else if subnetDetails != nil {
			subnetIdentity.ID = subnetDetails.ID
		} else {
			subnetIdentity.ID = &m.IBMVPCMachine.Spec.PrimaryNetworkInterface.Subnet
		}
	}
	primaryNetworkInterface := &vpcv1.NetworkInterfacePrototype{
		Subnet: subnetIdentity,
	}

	// Populate the PrimaryNetworkInterface's SecurityGroups, if provided.
	if len(m.IBMVPCMachine.Spec.PrimaryNetworkInterface.SecurityGroups) > 0 {
		securityGroups := make([]vpcv1.SecurityGroupIdentityIntf, 0, len(m.IBMVPCMachine.Spec.PrimaryNetworkInterface.SecurityGroups))
		for _, sg := range m.IBMVPCMachine.Spec.PrimaryNetworkInterface.SecurityGroups {
			// Try using Security Group name if provided.
			if sg.Name != nil {
				// If Network Status is available, attempt to retrieve Security Group ID from there.
				if m.IBMVPCCluster.Status.Network != nil {
					if sgStatus, ok := m.IBMVPCCluster.Status.Network.SecurityGroups[*sg.Name]; ok {
						securityGroups = append(securityGroups, &vpcv1.SecurityGroupIdentityByID{
							ID: ptr.To(sgStatus.ID),
						})
						continue
					}
				}
				// If not found in Network Status, try looking up the Security Group via API.
				sgDetails, err := m.IBMVPCClient.GetSecurityGroupByName(*sg.Name)
				if err != nil {
					return nil, fmt.Errorf("error retrieving security group id with name %s for machine %s: %w", *sg.Name, m.IBMVPCMachine.Name, err)
				} else if sgDetails != nil {
					securityGroups = append(securityGroups, &vpcv1.SecurityGroupIdentityByID{
						ID: sgDetails.ID,
					})
					continue
				}
				// If Name was provided but it cannot be found in Network Status or via API, return an error.
				return nil, fmt.Errorf("error cannot find security group %s for machine %s", *sg.Name, m.IBMVPCMachine.Name)
			}
			// If ID is provided for Security Group, attempt lookup to confirm it exists.
			if sg.ID != nil {
				sgOptions := &vpcv1.GetSecurityGroupOptions{
					ID: sg.ID,
				}
				sgDetails, _, err := m.IBMVPCClient.GetSecurityGroup(sgOptions)
				if err != nil {
					return nil, fmt.Errorf("error retrieving security by id %s for machine %s: %w", *sg.ID, m.IBMVPCMachine.Name, err)
				} else if sgDetails == nil {
					return nil, fmt.Errorf("error security group not found with id %s for machine %s", *sg.ID, m.IBMVPCMachine.Name)
				}
				securityGroups = append(securityGroups, &vpcv1.SecurityGroupIdentityByID{
					ID: sg.ID,
				})
				continue
			}
			// TODO(cjschaef): Replace with webhook validation check.
			return nil, fmt.Errorf("error no name or id provided for security group for machine %s", m.IBMVPCMachine.Name)
		}
		// After processing all Security Groups, add them to the PrimaryNetworkInterface.
		primaryNetworkInterface.SecurityGroups = securityGroups
	}

	var resourceGroupIdentity *vpcv1.ResourceGroupIdentity
	if m.IBMVPCCluster.Status.ResourceGroup != nil {
		resourceGroupIdentity = &vpcv1.ResourceGroupIdentity{
			ID: &m.IBMVPCCluster.Status.ResourceGroup.ID,
		}
	} else {
		resourceGroupIdentity = &vpcv1.ResourceGroupIdentity{
			ID: &m.IBMVPCCluster.Spec.ResourceGroup,
		}
	}

	var vpcIdentity *vpcv1.VPCIdentityByID
	if m.IBMVPCCluster.Status.Network != nil && m.IBMVPCCluster.Status.Network.VPC != nil {
		vpcIdentity = &vpcv1.VPCIdentityByID{
			ID: ptr.To(m.IBMVPCCluster.Status.Network.VPC.ID),
		}
	}

	zone := &vpcv1.ZoneIdentity{
		Name: &m.IBMVPCMachine.Spec.Zone,
	}

	// Populate Placement target details, if provided.
	var placementTarget vpcv1.InstancePlacementTargetPrototypeIntf
	if m.IBMVPCMachine.Spec.PlacementTarget != nil {
		placementTarget, err = m.configurePlacementTarget()
		if err != nil {
			return nil, fmt.Errorf("error configuration machine placement target: %w", err)
		}
	}

	// Populate any SSH Keys, if provided.
	sshKeys := make([]vpcv1.KeyIdentityIntf, 0)
	if m.IBMVPCMachine.Spec.SSHKeys != nil {
		for _, sshKey := range m.IBMVPCMachine.Spec.SSHKeys {
			keyID, err := fetchKeyID(sshKey, m)
			if err != nil {
				return nil, fmt.Errorf("error while fetching SSHKey: %v error: %v", sshKey, err)
			}
			key := &vpcv1.KeyIdentity{
				ID: keyID,
			}
			sshKeys = append(sshKeys, key)
		}
	}

	// Populate boot volume attachment, if provided.
	var bootVolumeAttachment *vpcv1.VolumeAttachmentPrototypeInstanceByImageContext
	if m.IBMVPCMachine.Spec.BootVolume != nil {
		bootVolumeAttachment = m.volumeToVPCVolumeAttachment(m.IBMVPCMachine.Spec.BootVolume)
	}

	// Configure the Machine's Image or CatalogOffering based on provided fields.
	// If an Image was provided, use that, if a Catalog Offering was provided use that (based on details provided), otherwise return an error.
	if m.IBMVPCMachine.Spec.Image != nil {
		imageInstancePrototype := &vpcv1.InstancePrototype{
			Name:                    ptr.To(m.IBMVPCMachine.Name),
			Profile:                 profile,
			PrimaryNetworkInterface: primaryNetworkInterface,
			ResourceGroup:           resourceGroupIdentity,
			UserData:                ptr.To(cloudInitData),
			VPC:                     vpcIdentity,
			Zone:                    zone,
		}
		imageID, err := fetchImageID(m.IBMVPCMachine.Spec.Image, m)
		if err != nil {
			record.Warnf(m.IBMVPCMachine, "FailedRetrieveImage", "Failed image retrieval - %w", err)
			return nil, fmt.Errorf("error while fetching image ID: %w", err)
		}
		imageInstancePrototype.Image = &vpcv1.ImageIdentity{
			ID: imageID,
		}

		// Configure additional fields if they were populated.
		if placementTarget != nil {
			imageInstancePrototype.PlacementTarget = placementTarget
		}
		if len(sshKeys) > 0 {
			imageInstancePrototype.Keys = sshKeys
		}
		if bootVolumeAttachment != nil {
			imageInstancePrototype.BootVolumeAttachment = bootVolumeAttachment
		}

		m.Logger.Info("machine creation configured with existing image", "machineName", m.IBMVPCMachine.Name, "imageID", *imageID)
		options.SetInstancePrototype(imageInstancePrototype)
	} else if m.IBMVPCMachine.Spec.CatalogOffering != nil {
		catalogInstancePrototype := &vpcv1.InstancePrototypeInstanceByCatalogOffering{
			Name:                    ptr.To(m.IBMVPCMachine.Name),
			Profile:                 profile,
			PrimaryNetworkInterface: primaryNetworkInterface,
			ResourceGroup:           resourceGroupIdentity,
			UserData:                ptr.To(cloudInitData),
			VPC:                     vpcIdentity,
			Zone:                    zone,
		}
		catalogOfferingPrototype := &vpcv1.InstanceCatalogOfferingPrototype{}
		if m.IBMVPCMachine.Spec.CatalogOffering.OfferingCRN != nil {
			// TODO(cjschaef): Perform lookup or use webhook validation to confirm Catalog Offering CRN.
			catalogOfferingPrototype.Offering = &vpcv1.CatalogOfferingIdentityCatalogOfferingByCRN{
				CRN: m.IBMVPCMachine.Spec.CatalogOffering.OfferingCRN,
			}
			m.Logger.Info("machine creation configured with catalog offering", "machineName", m.IBMVPCMachine.Name, "offeringCRN", *m.IBMVPCMachine.Spec.CatalogOffering.OfferingCRN)
		} else if m.IBMVPCMachine.Spec.CatalogOffering.VersionCRN != nil {
			// TODO(cjschaef): Perform lookup or use webhook validation to confirm Catalog Offering Version CRN.
			catalogOfferingPrototype.Version = &vpcv1.CatalogOfferingVersionIdentityCatalogOfferingVersionByCRN{
				CRN: m.IBMVPCMachine.Spec.CatalogOffering.VersionCRN,
			}
			m.Logger.Info("machine creation configured with catalog version", "machineName", m.IBMVPCMachine.Name, "versionCRN", *m.IBMVPCMachine.Spec.CatalogOffering.VersionCRN)
		} else {
			// TODO(cjschaef): Look to add webhook validation to ensure one is provided.
			return nil, fmt.Errorf("error catalog offering missing offering crn and version crn, one must be provided")
		}
		if m.IBMVPCMachine.Spec.CatalogOffering.PlanCRN != nil {
			// TODO(cjschaef): Perform lookup or use webhook validation to confirm Catalog Offering Plan CRN.
			catalogOfferingPrototype.Plan = &vpcv1.CatalogOfferingVersionPlanIdentityCatalogOfferingVersionPlanByCRN{
				CRN: m.IBMVPCMachine.Spec.CatalogOffering.PlanCRN,
			}
			m.Logger.Info("machine creation configured with catalog plan", "machineName", m.IBMVPCMachine.Name, "planCRN", *m.IBMVPCMachine.Spec.CatalogOffering.PlanCRN)
		}

		// Configure additional fields if they were populated.
		if placementTarget != nil {
			catalogInstancePrototype.PlacementTarget = placementTarget
		}
		if len(sshKeys) > 0 {
			catalogInstancePrototype.Keys = sshKeys
		}
		if bootVolumeAttachment != nil {
			catalogInstancePrototype.BootVolumeAttachment = bootVolumeAttachment
		}

		catalogInstancePrototype.CatalogOffering = catalogOfferingPrototype
		options.SetInstancePrototype(catalogInstancePrototype)
	} else {
		// TODO(cjschaef): Move this to webhook validation.
		return nil, fmt.Errorf("error no machine image or catalog offering provided to build: %s", m.IBMVPCMachine.Spec.Name)
	}

	m.Logger.Info("creating instance", "createOptions", options, "name", m.IBMVPCMachine.Name, "profile", *profile.Name, "resourceGroup", resourceGroupIdentity, "vpc", vpcIdentity, "zone", zone)
	instance, _, err := m.IBMVPCClient.CreateInstance(options)
	if err != nil {
		record.Warnf(m.IBMVPCMachine, "FailedCreateInstance", "Failed instance creation - %s, %v", options, err)
	} else {
		record.Eventf(m.IBMVPCMachine, "SuccessfulCreateInstance", "Created Instance %q", *instance.Name)
	}
	return instance, err
}

// configurePlacementTarget will configure a Machine's Placement Target based on the Machine's provided configuration, if supplied.
func (m *MachineScope) configurePlacementTarget() (vpcv1.InstancePlacementTargetPrototypeIntf, error) {
	// TODO(cjschaef): We currently don't support the other placement target options (Dedicated Host Group, Placement Group), they need to be added.
	if m.IBMVPCMachine.Spec.PlacementTarget.DedicatedHost != nil {
		// Lookup Dedicated Host ID by Name if it was provided.
		var dedicatedHostID *string
		if m.IBMVPCMachine.Spec.PlacementTarget.DedicatedHost.ID != nil {
			dedicatedHostID = m.IBMVPCMachine.Spec.PlacementTarget.DedicatedHost.ID
		} else if m.IBMVPCMachine.Spec.PlacementTarget.DedicatedHost.Name != nil {
			dHost, err := m.IBMVPCClient.GetDedicatedHostByName(*m.IBMVPCMachine.Spec.PlacementTarget.DedicatedHost.Name)
			if err != nil {
				return nil, fmt.Errorf("error failed lookup of dedicated host by name %s: %w", *m.IBMVPCMachine.Spec.PlacementTarget.DedicatedHost.Name, err)
			} else if dHost == nil {
				return nil, fmt.Errorf("error no dedicated host found with name %s", *m.IBMVPCMachine.Spec.PlacementTarget.DedicatedHost.Name)
			}
			dedicatedHostID = dHost.ID
		}

		m.Logger.Info("machine creation configured with dedicated host placement", "machineName", m.IBMVPCMachine.Name, "dedicatedHostID", *dedicatedHostID)
		return &vpcv1.InstancePlacementTargetPrototypeDedicatedHostGroupIdentityDedicatedHostGroupIdentityByID{
			ID: dedicatedHostID,
		}, nil
	}
	return nil, nil
}

func (m *MachineScope) volumeToVPCVolumeAttachment(volume *infrav1beta2.VPCVolume) *vpcv1.VolumeAttachmentPrototypeInstanceByImageContext {
	bootVolume := &vpcv1.VolumeAttachmentPrototypeInstanceByImageContext{
		DeleteVolumeOnInstanceDelete: core.BoolPtr(volume.DeleteVolumeOnInstanceDelete),
		Volume:                       &vpcv1.VolumePrototypeInstanceByImageContext{},
	}

	if volume.Name != "" {
		bootVolume.Volume.Name = core.StringPtr(volume.Name)
	}

	if volume.Profile != "" {
		bootVolume.Volume.Profile = &vpcv1.VolumeProfileIdentity{
			Name: core.StringPtr(volume.Profile),
		}
	}

	if volume.SizeGiB != 0 {
		bootVolume.Volume.Capacity = core.Int64Ptr(volume.SizeGiB)
	}

	if volume.Iops != 0 {
		bootVolume.Volume.Iops = core.Int64Ptr(volume.Iops)
	}

	if volume.EncryptionKeyCRN != "" {
		bootVolume.Volume.EncryptionKey = &vpcv1.EncryptionKeyIdentity{
			CRN: core.StringPtr(volume.EncryptionKeyCRN),
		}
		m.Logger.Info("machine creation configured with volumn encryption key", "machineName", m.IBMVPCMachine.Name, "encryptionKeyCRN", volume.EncryptionKeyCRN)
	}

	return bootVolume
}

// DeleteMachine deletes the vpc machine associated with machine instance id.
func (m *MachineScope) DeleteMachine() error {
	if m.IBMVPCMachine.Status.InstanceID == "" {
		return nil
	}
	options := &vpcv1.DeleteInstanceOptions{}
	options.SetID(m.IBMVPCMachine.Status.InstanceID)
	_, err := m.IBMVPCClient.DeleteInstance(options)
	if err != nil {
		record.Warnf(m.IBMVPCMachine, "FailedDeleteInstance", "Failed instance deletion - %v", err)
	} else {
		record.Eventf(m.IBMVPCMachine, "SuccessfulDeleteInstance", "Deleted Instance %q", m.IBMVPCMachine.Name)
	}
	return err
}

func (m *MachineScope) ensureInstanceUnique(instanceName string) (*vpcv1.Instance, error) {
	var instance *vpcv1.Instance
	f := func(start string) (bool, string, error) {
		// check for existing instances
		listInstancesOptions := &vpcv1.ListInstancesOptions{}
		if start != "" {
			listInstancesOptions.Start = &start
		}

		instancesList, _, err := m.IBMVPCClient.ListInstances(listInstancesOptions)
		if err != nil {
			return false, "", err
		}

		if instancesList == nil {
			return false, "", fmt.Errorf("instance list returned is nil")
		}

		for i, ins := range instancesList.Instances {
			if (*ins.Name) == instanceName {
				instance = &instancesList.Instances[i]
				return true, "", nil
			}
		}

		if instancesList.Next != nil && *instancesList.Next.Href != "" {
			return false, *instancesList.Next.Href, nil
		}
		return true, "", nil
	}

	if err := utils.PagingHelper(f); err != nil {
		return nil, err
	}

	return instance, nil
}

// getLoadBalancerID will return the ID of a Load Balancer.
func (m *MachineScope) getLoadBalancerID(loadBalancer *infrav1beta2.VPCResource) (*string, error) {
	// Lookup Load Balancer ID by Name if necessary
	if loadBalancer.ID != nil {
		return loadBalancer.ID, nil
	} else if loadBalancer.Name != nil {
		loadBalancerDetails, err := m.IBMVPCClient.GetLoadBalancerByName(*loadBalancer.Name)
		if err != nil {
			return nil, fmt.Errorf("error failed to lookup load balancer id by name %s: %w", *loadBalancer.Name, err)
		} else if loadBalancerDetails == nil || loadBalancerDetails.ID == nil {
			return nil, fmt.Errorf("error unable to find load balancer id with name: %s", *loadBalancer.Name)
		}
		return loadBalancerDetails.ID, nil
	}

	return nil, fmt.Errorf("error no load balancer id or name provided")
}

// getLoadBalancerPoolID will return the ID of a Load Balancer Pool.
func (m *MachineScope) getLoadBalancerPoolID(pool *infrav1beta2.VPCResource, loadBalancerID string) (*string, error) {
	// Lookup Load Balancer Pool ID by Name if necessary
	if pool.ID != nil {
		return pool.ID, nil
	} else if pool.Name != nil {
		loadBalancerPoolDetails, err := m.IBMVPCClient.GetLoadBalancerPoolByName(loadBalancerID, *pool.Name)
		if err != nil {
			return nil, fmt.Errorf("error failed to lookup load balancer pool id by name %s: %w", *pool.Name, err)
		} else if loadBalancerPoolDetails == nil || loadBalancerPoolDetails.ID == nil {
			return nil, fmt.Errorf("error unable to find load balancer pool id with name: %s", *pool.Name)
		}
		return loadBalancerPoolDetails.ID, nil
	}

	return nil, fmt.Errorf("error no load balancer pool id or name provided")
}

// ReconcileVPCLoadBalancerPoolMember reconciles a Machine's Load Balancer Pool membership.
func (m *MachineScope) ReconcileVPCLoadBalancerPoolMember(poolMember infrav1beta2.VPCLoadBalancerBackendPoolMember) (bool, error) {
	// Collect the Machine's internal IP.
	internalIP := m.GetMachineInternalIP()
	if internalIP == nil {
		// TODO(cjschaef): Allow options for adding Machines to the pool without an internal IP
		return false, fmt.Errorf("error unable to find machine's internal ip to use for load balancer pool")
	}

	// Check if Instance is already a member of Load Balancer Backend Pool.
	existingMember, err := m.checkVPCLoadBalancerPoolMemberExists(poolMember, internalIP)
	if err != nil {
		return false, fmt.Errorf("error failed to check if member exists in pool")
	} else if existingMember != nil {
		// If the member already exists in the pool, check whether it is ready (active).
		if *existingMember.ProvisioningStatus == vpcv1.LoadBalancerPoolMemberProvisioningStatusActiveConst {
			return false, nil
		}
		// If not ready, trigger requeue.
		return true, nil
	}

	// Otherwise, create VPC Load Balancer Backend Pool Member
	return m.createVPCLoadBalancerPoolMember(poolMember, internalIP)
}

// checkVPCLoadBalancerPoolMemberExists determines whether a Machine's Load Balancer Pool membership already exists.
func (m *MachineScope) checkVPCLoadBalancerPoolMemberExists(poolMember infrav1beta2.VPCLoadBalancerBackendPoolMember, internalIP *string) (*vpcv1.LoadBalancerPoolMember, error) {
	loadBalancerID, err := m.getLoadBalancerID(&poolMember.LoadBalancer)
	if err != nil {
		return nil, fmt.Errorf("error checking if load balancer pool member exists: %w", err)
	}

	poolID, err := m.getLoadBalancerPoolID(&poolMember.Pool, *loadBalancerID)
	if err != nil {
		return nil, fmt.Errorf("error checking if load balancer pool member exists: %w", err)
	}

	// Check if the Pool has a member matching the Machine's internal IP.
	listLoadBalancerPoolMembersOptions := &vpcv1.ListLoadBalancerPoolMembersOptions{}
	listLoadBalancerPoolMembersOptions.SetLoadBalancerID(*loadBalancerID)
	listLoadBalancerPoolMembersOptions.SetPoolID(*poolID)

	poolMembers, detailedResponse, err := m.IBMVPCClient.ListLoadBalancerPoolMembers(listLoadBalancerPoolMembersOptions)
	if err != nil {
		return nil, fmt.Errorf("error listing members for load balancer pool: %w", err)
	} else if detailedResponse != nil && detailedResponse.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("error unable to find load balancer or pool members")
	} else if poolMembers == nil {
		return nil, fmt.Errorf("error no load balancer pool members returned")
	}

	for _, member := range poolMembers.Members {
		if target, ok := member.Target.(*vpcv1.LoadBalancerPoolMemberTarget); ok {
			// Verify the target address matches the Machine's internal IP.
			if *target.Address == *internalIP {
				m.Logger.Info("found existing load balancer pool member for machine", "machineName", m.IBMVPCMachine.Spec.Name, "internalIP", *internalIP, "poolID", *poolID, "loadBalancerID", *loadBalancerID)
				return ptr.To(member), nil
			}
		}
	}

	// If a match was not found at this point, expect that it doesn't exist.
	return nil, nil
}

// createVPCLoadBalancerPoolMember will create a new member within a Load Balancer Pool for the Machine's internal IP.
func (m *MachineScope) createVPCLoadBalancerPoolMember(poolMember infrav1beta2.VPCLoadBalancerBackendPoolMember, internalIP *string) (bool, error) {
	// Retrieve the Load Balancer ID.
	loadBalancerID, err := m.getLoadBalancerID(&poolMember.LoadBalancer)
	if err != nil {
		return false, fmt.Errorf("error creating load balancer pool member: %w", err)
	}

	loadBalancerBackendPoolID, err := m.getLoadBalancerPoolID(&poolMember.Pool, *loadBalancerID)
	if err != nil {
		return false, fmt.Errorf("error creating load balancer pool member: %w", err)
	}

	// Populate the LB Pool Member options.
	options := &vpcv1.CreateLoadBalancerPoolMemberOptions{
		LoadBalancerID: loadBalancerID,
		PoolID:         loadBalancerBackendPoolID,
		Port:           ptr.To(poolMember.Port),
		Target: &vpcv1.LoadBalancerPoolMemberTargetPrototypeIP{
			Address: internalIP,
		},
	}

	// Set the weight if it was provided.
	// TODO(cjschaef): Weight only affects weightroundrobin algorithm on a LB. We may wish to validate this via webhook, unless API ignores this field for other algorithms (and it doesn't matter if we provide it).
	if poolMember.Weight != nil {
		options.Weight = poolMember.Weight
	}

	// Create Machine Load Balancer Pool Member.
	loadBalancerPoolMember, _, err := m.IBMVPCClient.CreateLoadBalancerPoolMember(options)
	if err != nil {
		return false, fmt.Errorf("error failed creating load balancer backend pool member: %w", err)
	}
	m.Logger.Info("created load balancer backend pool member", "instanceID", m.IBMVPCMachine.Status.InstanceID, "loadBalancerID", loadBalancerID, "loadBalancerBackendPoolID", loadBalancerBackendPoolID, "port", poolMember.Port, "loadBalancerBackendPoolMemberID", loadBalancerPoolMember.ID)

	// Add the new pool member details to the Machine Status.
	// To prevent additional API calls, only use ID's and not Name's, as reconciliation does not rely on Name's for these resources in Status.
	newMember := infrav1beta2.VPCLoadBalancerBackendPoolMember{
		LoadBalancer: infrav1beta2.VPCResource{
			ID: loadBalancerID,
		},
		Pool: infrav1beta2.VPCResource{
			ID: loadBalancerBackendPoolID,
		},
		Port: poolMember.Port,
	}
	if poolMember.Weight != nil {
		newMember.Weight = poolMember.Weight
	}

	m.IBMVPCMachine.Status.LoadBalancerPoolMembers = append(m.IBMVPCMachine.Status.LoadBalancerPoolMembers, newMember)

	// TODO(cjschaef): Tagging does not appear valid for this resource, so we currently skip it.

	// Trigger a requeue after creating the pool member to confirm it exists on next round.
	return true, nil
}

// CreateVPCLoadBalancerPoolMember creates a new pool member and adds it to the load balancer pool.
func (m *MachineScope) CreateVPCLoadBalancerPoolMember(internalIP *string, targetPort int64) (*vpcv1.LoadBalancerPoolMember, error) {
	loadBalancer, _, err := m.IBMVPCClient.GetLoadBalancer(&vpcv1.GetLoadBalancerOptions{
		ID: m.IBMVPCCluster.Status.VPCEndpoint.LBID,
	})
	if err != nil {
		return nil, err
	}

	if *loadBalancer.ProvisioningStatus != string(infrav1beta2.VPCLoadBalancerStateActive) {
		return nil, fmt.Errorf("error load balancer is not in active state")
	}

	if len(loadBalancer.Pools) == 0 {
		return nil, fmt.Errorf("error no pools exist for the load balancer")
	}

	options := &vpcv1.CreateLoadBalancerPoolMemberOptions{}
	options.SetLoadBalancerID(*loadBalancer.ID)
	options.SetPoolID(*loadBalancer.Pools[0].ID)
	options.SetTarget(&vpcv1.LoadBalancerPoolMemberTargetPrototype{
		Address: internalIP,
	})
	options.SetPort(targetPort)

	listOptions := &vpcv1.ListLoadBalancerPoolMembersOptions{}
	listOptions.SetLoadBalancerID(*loadBalancer.ID)
	listOptions.SetPoolID(*loadBalancer.Pools[0].ID)
	listLoadBalancerPoolMembers, _, err := m.IBMVPCClient.ListLoadBalancerPoolMembers(listOptions)
	if err != nil {
		return nil, fmt.Errorf("error failed to bind ListLoadBalancerPoolMembers to control plane %s/%s: %w", m.IBMVPCMachine.Namespace, m.IBMVPCMachine.Name, err)
	}

	for _, member := range listLoadBalancerPoolMembers.Members {
		if _, ok := member.Target.(*vpcv1.LoadBalancerPoolMemberTarget); ok {
			mtarget := member.Target.(*vpcv1.LoadBalancerPoolMemberTarget)
			if *mtarget.Address == *internalIP && *member.Port == targetPort {
				m.Logger.V(3).Info("PoolMember already exist")
				return nil, nil
			}
		}
	}

	loadBalancerPoolMember, _, err := m.IBMVPCClient.CreateLoadBalancerPoolMember(options)
	if err != nil {
		return nil, err
	}
	return loadBalancerPoolMember, nil
}

// DeleteVPCLoadBalancerPoolMember deletes a pool member from the load balancer pool.
func (m *MachineScope) DeleteVPCLoadBalancerPoolMember() error {
	if m.IBMVPCMachine.Status.InstanceID == "" {
		m.Info("instance is not created, ignore deleting load balancer pool member")
		return nil
	}

	// If the Machine has Load Balancer Pool Members defined in its Status (part of extended VPC Machine support), process the removal of those members versus the legacy single LB design.
	if len(m.IBMVPCMachine.Status.LoadBalancerPoolMembers) > 0 {
		return m.deleteVPCLoadBalancerPoolMembers()
	}

	loadBalancer, _, err := m.IBMVPCClient.GetLoadBalancer(&vpcv1.GetLoadBalancerOptions{
		ID: m.IBMVPCCluster.Status.VPCEndpoint.LBID,
	})
	if err != nil {
		return err
	}

	if len(loadBalancer.Pools) == 0 {
		return nil
	}

	instance, _, err := m.IBMVPCClient.GetInstance(&vpcv1.GetInstanceOptions{
		ID: core.StringPtr(m.IBMVPCMachine.Status.InstanceID),
	})
	if err != nil {
		return err
	}

	listOptions := &vpcv1.ListLoadBalancerPoolMembersOptions{}
	listOptions.SetLoadBalancerID(*loadBalancer.ID)
	listOptions.SetPoolID(*loadBalancer.Pools[0].ID)
	listLoadBalancerPoolMembers, _, err := m.IBMVPCClient.ListLoadBalancerPoolMembers(listOptions)
	if err != nil {
		return err
	}

	for _, member := range listLoadBalancerPoolMembers.Members {
		if _, ok := member.Target.(*vpcv1.LoadBalancerPoolMemberTarget); ok {
			mtarget := member.Target.(*vpcv1.LoadBalancerPoolMemberTarget)
			if *mtarget.Address == *instance.PrimaryNetworkInterface.PrimaryIP.Address {
				if *loadBalancer.ProvisioningStatus != string(infrav1beta2.VPCLoadBalancerStateActive) {
					return fmt.Errorf("load balancer is not in active state")
				}

				deleteOptions := &vpcv1.DeleteLoadBalancerPoolMemberOptions{}
				deleteOptions.SetLoadBalancerID(*loadBalancer.ID)
				deleteOptions.SetPoolID(*loadBalancer.Pools[0].ID)
				deleteOptions.SetID(*member.ID)

				if _, err := m.IBMVPCClient.DeleteLoadBalancerPoolMember(deleteOptions); err != nil {
					return err
				}
				return nil
			}
		}
	}
	return nil
}

// deleteVPCLoadBalancerPoolMembers provides support to delete Load Balancer Pools Members for a Machine that are tracked in the Machine's Status, which is part of the extended VPC Machine support.
// This new support allows a Machine to have members in multiple Load Balancers, as defined by the Machine Spec, rather than defaulting (legacy) to the single Cluster Load Balancer.
func (m *MachineScope) deleteVPCLoadBalancerPoolMembers() error {
	// Retrieve the Instance details immediately (without them the member cannot be safely deleted).
	instanceOptions := &vpcv1.GetInstanceOptions{
		ID: ptr.To(m.IBMVPCMachine.Status.InstanceID),
	}
	instanceDetails, _, err := m.IBMVPCClient.GetInstance(instanceOptions)
	if err != nil {
		return fmt.Errorf("error retrieving instance for load balancer pool member deletion for machine %s: %w", m.IBMVPCMachine.Name, err)
	}
	// Verify the instance has a primary network interface IP address needed to delete the correct Load Balancer Pool Member.
	if instanceDetails.PrimaryNetworkInterface == nil || instanceDetails.PrimaryNetworkInterface.PrimaryIP == nil || instanceDetails.PrimaryNetworkInterface.PrimaryIP.Address == nil {
		return fmt.Errorf("error instance is missing the primary network interface IP address for load balancer pool member deletion for machine: %s", m.IBMVPCMachine.Name)
	}
	m.Logger.V(5).Info("collected instance details for load balancer pool member deletion", "machienName", m.IBMVPCMachine.Name, "instanceID", *instanceDetails.ID, "instanceIP", *instanceDetails.PrimaryNetworkInterface.PrimaryIP.Address)

	cleanupIncomplete := false
	for _, member := range m.IBMVPCMachine.Status.LoadBalancerPoolMembers {
		// Retrieve the Load Balancer.
		var loadBalancerDetails *vpcv1.LoadBalancer
		if member.LoadBalancer.ID != nil {
			loadBalancerDetails, _, err = m.IBMVPCClient.GetLoadBalancer(&vpcv1.GetLoadBalancerOptions{
				ID: member.LoadBalancer.ID,
			})
		} else if member.LoadBalancer.Name != nil {
			loadBalancerDetails, err = m.IBMVPCClient.GetLoadBalancerByName(*member.LoadBalancer.Name)
		} else {
			return fmt.Errorf("error load balancer has no id or name for load balancer pool member deleteion")
		}
		if err != nil {
			return fmt.Errorf("error retrieving load balancer for load balancer pool member deletion for machine %s: %w", m.IBMVPCMachine.Name, err)
		}
		m.Logger.V(5).Info("collected load balancer for load balancer pool member deletion", "machineName", m.IBMVPCMachine.Name, "loadBalancerID", *loadBalancerDetails.ID)

		// Lookup the Load Balancer Pool ID, if only name is available.
		loadBalancerPoolID, err := m.getLoadBalancerPoolID(ptr.To(member.Pool), *loadBalancerDetails.ID)
		if err != nil {
			return fmt.Errorf("error retrieving load balancer pool id for load balancer pool member deletion for machine %s: %w", m.IBMVPCMachine.Name, err)
		}
		m.Logger.V(5).Info("collected load balancer pool id for load balancer pool member deletion", "machineName", m.IBMVPCMachine.Name, "loadBalancerPoolID", *loadBalancerPoolID)

		listMembersOptions := &vpcv1.ListLoadBalancerPoolMembersOptions{
			LoadBalancerID: loadBalancerDetails.ID,
			PoolID:         loadBalancerPoolID,
		}

		m.Logger.V(5).Info("list load balancer pool members options", "machineName", m.IBMVPCMachine.Name, "options", *listMembersOptions)
		poolMembers, _, err := m.IBMVPCClient.ListLoadBalancerPoolMembers(listMembersOptions)
		if err != nil {
			return fmt.Errorf("error retrieving load balancer pool members for load balancer pool member deletion for machine %s: %w", m.IBMVPCMachine.Name, err)
		}

		for _, poolMember := range poolMembers.Members {
			poolMemberTarget, ok := poolMember.Target.(*vpcv1.LoadBalancerPoolMemberTarget)
			// If the member isn't a LoadBalancerPoolMemberTarget, has no Address, or the Address doesn't match the Machine's Primary IP Address, move to the next member.
			if !ok || poolMemberTarget.Address == nil || *poolMemberTarget.Address != *instanceDetails.PrimaryNetworkInterface.PrimaryIP.Address {
				continue
			}

			m.Logger.V(3).Info("found load balancer pool member to delete", "machineName", m.IBMVPCMachine.Name, "poolMemberID", *poolMember.ID)
			// Make LB status check now that it has been determined a change is required.
			if *loadBalancerDetails.ProvisioningStatus != string(infrav1beta2.VPCLoadBalancerStateActive) {
				m.Logger.V(5).Info("load balancer not in active status prior to load balancer pool member deletion", "machineName", m.IBMVPCMachine.Name, "loadBalancerID", *loadBalancerDetails.ID, "loadBalancerProvisioningStatus", *loadBalancerDetails.ProvisioningStatus)
				// Set flag that some cleanup was not completed, and break out of member target loop, to try next member from Machine Status.
				cleanupIncomplete = true
				break
			}

			deleteOptions := &vpcv1.DeleteLoadBalancerPoolMemberOptions{
				ID:             poolMember.ID,
				LoadBalancerID: loadBalancerDetails.ID,
				PoolID:         loadBalancerPoolID,
			}

			m.Logger.V(5).Info("delete load balancer pool member options", "machineName", m.IBMVPCMachine.Name, "options", *deleteOptions)
			// Delete the matching Load Balancer Pool Member.
			_, err := m.IBMVPCClient.DeleteLoadBalancerPoolMember(deleteOptions)
			if err != nil {
				return fmt.Errorf("error deleting load balancer pool member for machine: %s: %w", m.IBMVPCMachine.Name, err)
			}
			m.Logger.V(3).Info("deleted load balancer pool member", "machineName", m.IBMVPCMachine.Name, "loadBalancerID", *loadBalancerDetails.ID, "loadBalancerPoolID", *loadBalancerPoolID, "loadBalancerPoolMemberID", *poolMember.ID)
		}
	}

	if cleanupIncomplete {
		return fmt.Errorf("error load balancer pool member deletion could not complete as a load balancer was not active for machine: %s", m.IBMVPCMachine.Name)
	}

	// Assume if no errors were encountered, and if either all delete calls were successful or none were required, all Load Balancer Pool Members have been successfully deleted.
	return nil
}

// PatchObject persists the cluster configuration and status.
func (m *MachineScope) PatchObject() error {
	return m.patchHelper.Patch(context.TODO(), m.IBMVPCMachine)
}

// Close closes the current scope persisting the cluster configuration and status.
func (m *MachineScope) Close() error {
	return m.PatchObject()
}

// GetBootstrapData returns the bootstrap data from the secret in the Machine's bootstrap.dataSecretName.
func (m *MachineScope) GetBootstrapData() (string, error) {
	if m.Machine.Spec.Bootstrap.DataSecretName == nil {
		return "", errors.New("error retrieving bootstrap data: linked Machine's bootstrap.dataSecretName is nil")
	}

	secret := &corev1.Secret{}
	key := types.NamespacedName{Namespace: m.Machine.Namespace, Name: *m.Machine.Spec.Bootstrap.DataSecretName}
	if err := m.Client.Get(context.TODO(), key, secret); err != nil {
		return "", fmt.Errorf("failed to retrieve bootstrap data secret for IBMVPCMachine %s/%s: %w", m.Machine.Namespace, m.Machine.Name, err)
	}

	value, ok := secret.Data["value"]
	if !ok {
		return "", errors.New("error retrieving bootstrap data: secret value key is missing")
	}
	return string(value), nil
}

func fetchKeyID(key *infrav1beta2.IBMVPCResourceReference, m *MachineScope) (*string, error) {
	if key.ID == nil && key.Name == nil {
		return nil, fmt.Errorf("both ID and Name can't be nil")
	}

	if key.ID != nil {
		return key.ID, nil
	}

	var k *vpcv1.Key
	f := func(start string) (bool, string, error) {
		// check for existing keys
		listKeysOptions := &vpcv1.ListKeysOptions{}
		if start != "" {
			listKeysOptions.Start = &start
		}

		keysList, _, err := m.IBMVPCClient.ListKeys(listKeysOptions)
		if err != nil {
			m.Logger.Error(err, "Failed to get keys")
			return false, "", err
		}

		if keysList == nil {
			return false, "", fmt.Errorf("key list returned is nil")
		}

		for i, ks := range keysList.Keys {
			if *ks.Name == *key.Name {
				m.Logger.V(3).Info("Key found with ID", "Key", *ks.Name, "ID", *ks.ID)
				k = &keysList.Keys[i]
				return true, "", nil
			}
		}

		if keysList.Next != nil && *keysList.Next.Href != "" {
			return false, *keysList.Next.Href, nil
		}
		return true, "", nil
	}

	if err := utils.PagingHelper(f); err != nil {
		return nil, err
	}

	if k != nil {
		return k.ID, nil
	}

	return nil, fmt.Errorf("sshkey does not exist - failed to find Key ID")
}

func fetchImageID(image *infrav1beta2.IBMVPCResourceReference, m *MachineScope) (*string, error) {
	if image.ID == nil && image.Name == nil {
		return nil, fmt.Errorf("both ID and Name can't be nil")
	}

	if image.ID != nil {
		return image.ID, nil
	}

	var img *vpcv1.Image
	f := func(start string) (bool, string, error) {
		// check for existing images
		resourceGroupID := ptr.To(m.IBMVPCCluster.Spec.ResourceGroup)
		if m.IBMVPCCluster.Status.ResourceGroup != nil {
			resourceGroupID = ptr.To(m.IBMVPCCluster.Status.ResourceGroup.ID)
		}
		listImagesOptions := &vpcv1.ListImagesOptions{
			ResourceGroupID: resourceGroupID,
		}
		if start != "" {
			listImagesOptions.Start = &start
		}

		imagesList, _, err := m.IBMVPCClient.ListImages(listImagesOptions)
		if err != nil {
			m.Logger.Error(err, "Failed to get images")
			return false, "", err
		}

		if imagesList == nil {
			return false, "", fmt.Errorf("image list returned is nil")
		}

		for j, i := range imagesList.Images {
			if *image.Name == *i.Name {
				m.Logger.Info("Image found with ID", "Image", *i.Name, "ID", *i.ID)
				img = &imagesList.Images[j]
				return true, "", nil
			}
		}

		if imagesList.Next != nil && *imagesList.Next.Href != "" {
			return false, *imagesList.Next.Href, nil
		}
		return true, "", nil
	}

	if err := utils.PagingHelper(f); err != nil {
		return nil, err
	}

	if img != nil {
		return img.ID, nil
	}

	return nil, fmt.Errorf("image does not exist - failed to find an image ID")
}

// GetInstanceID will return the Machine's Instance ID.
func (m *MachineScope) GetInstanceID() string {
	return m.IBMVPCMachine.Status.InstanceID
}

// GetInstanceStatus will return the Machine's Instance Status.
func (m *MachineScope) GetInstanceStatus() string {
	return m.IBMVPCMachine.Status.InstanceStatus
}

// GetMachineInternalIP returns the machine's internal IP.
func (m *MachineScope) GetMachineInternalIP() *string {
	for _, address := range m.IBMVPCMachine.Status.Addresses {
		if address.Type == corev1.NodeInternalIP {
			return ptr.To(address.Address)
		}
	}
	return nil
}

// IsReady returns whether the machine is ready.
func (m *MachineScope) IsReady() bool {
	return m.IBMVPCMachine.Status.Ready
}

// SetAddresses sets the Machine's addresses.
func (m *MachineScope) SetAddresses(instance *vpcv1.Instance) {
	addresses := make([]corev1.NodeAddress, 0)
	addresses = append(addresses, corev1.NodeAddress{
		Type:    corev1.NodeInternalDNS,
		Address: *instance.Name,
	})
	addresses = append(addresses, corev1.NodeAddress{
		Type:    corev1.NodeHostName,
		Address: *instance.Name,
	})

	// Currently, only the single network interface designated by a subnet, is expected for the Instance (as its primary/internal IP).
	addresses = append(addresses, corev1.NodeAddress{
		Type:    corev1.NodeInternalIP,
		Address: *instance.PrimaryNetworkInterface.PrimaryIP.Address,
	})

	m.IBMVPCMachine.Status.Addresses = addresses
}

// SetFailureMessage will set the Machine's Failure Message.
func (m *MachineScope) SetFailureMessage(message string) {
	m.IBMVPCMachine.Status.FailureMessage = ptr.To(message)
}

// SetFailureReason will set the Machine's Failure Reason.
func (m *MachineScope) SetFailureReason(reason string) {
	m.IBMVPCMachine.Status.FailureReason = ptr.To(reason)
}

// SetInstanceID sets the Machine's Instance ID.
func (m *MachineScope) SetInstanceID(id string) {
	m.IBMVPCMachine.Status.InstanceID = id
}

// SetInstanceStatus sets the Machine's Instance Status.
func (m *MachineScope) SetInstanceStatus(status string) {
	m.IBMVPCMachine.Status.InstanceStatus = status
}

// SetNotReady sets the Machine Status as not ready.
func (m *MachineScope) SetNotReady() {
	m.IBMVPCMachine.Status.Ready = false
}

// SetProviderID will set the provider id for the machine.
func (m *MachineScope) SetProviderID(id *string) error {
	// Based on the ProviderIDFormat version the providerID format will be decided.
	if options.ProviderIDFormatType(options.ProviderIDFormat) == options.ProviderIDFormatV2 {
		accountID, err := utils.GetAccountIDWrapper()
		if err != nil {
			m.Logger.Error(err, "failed to get cloud account id", err.Error())
			return err
		}
		m.IBMVPCMachine.Spec.ProviderID = ptr.To(fmt.Sprintf("ibm://%s///%s/%s", accountID, m.Machine.Spec.ClusterName, *id))
	} else {
		return fmt.Errorf("invalid value for ProviderIDFormat")
	}
	return nil
}

// SetReady sets the Machine Status as ready.
func (m *MachineScope) SetReady() {
	m.IBMVPCMachine.Status.Ready = true
}

// CheckTagExists checks whether a user tag already exists.
func (m *MachineScope) CheckTagExists(tagName string) (bool, error) {
	exists, err := m.GlobalTaggingClient.GetTagByName(tagName)
	if err != nil {
		return false, fmt.Errorf("failed checking for tag: %w", err)
	}

	return exists != nil, nil
}

// TagResource will attach a user Tag to a resource.
func (m *MachineScope) TagResource(tagName string, resourceCRN string) error {
	// Verify the Tag to use exists, otherwise create it.
	exists, err := m.CheckTagExists(tagName)
	if err != nil {
		return fmt.Errorf("failure checking if tag %s exists: %w", tagName, err)
	}

	// Create tag if it doesn't exist.
	if !exists {
		createOptions := &globaltaggingv1.CreateTagOptions{}
		createOptions.SetTagNames([]string{tagName})
		if _, _, err := m.GlobalTaggingClient.CreateTag(createOptions); err != nil {
			return fmt.Errorf("failure creating tag: %w", err)
		}
	}

	// Finally, tag the resource.
	tagOptions := &globaltaggingv1.AttachTagOptions{}
	tagOptions.SetResources([]globaltaggingv1.Resource{
		{
			ResourceID: ptr.To(resourceCRN),
		},
	})
	tagOptions.SetTagName(tagName)
	tagOptions.SetTagType(globaltaggingv1.AttachTagOptionsTagTypeUserConst)

	if _, _, err = m.GlobalTaggingClient.AttachTag(tagOptions); err != nil {
		return fmt.Errorf("failure tagging resource: %w", err)
	}

	return nil
}

// APIServerPort returns the APIServerPort.
func (m *MachineScope) APIServerPort() int32 {
	if m.Cluster.Spec.ClusterNetwork != nil && m.Cluster.Spec.ClusterNetwork.APIServerPort != nil {
		return *m.Cluster.Spec.ClusterNetwork.APIServerPort
	}
	return infrav1beta2.DefaultAPIServerPort
}
