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

	"github.com/go-logr/logr"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"

	"sigs.k8s.io/controller-runtime/pkg/client"

	capiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/patch"

	infrav1beta2 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
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

	IBMVPCClient    vpc.Vpc
	Cluster         *capiv1beta1.Cluster
	Machine         *capiv1beta1.Machine
	IBMVPCCluster   *infrav1beta2.IBMVPCCluster
	IBMVPCMachine   *infrav1beta2.IBMVPCMachine
	ServiceEndpoint []endpoints.ServiceEndpoint
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

	return &MachineScope{
		Logger:        params.Logger,
		Client:        params.Client,
		IBMVPCClient:  vpcClient,
		Cluster:       params.Cluster,
		IBMVPCCluster: params.IBMVPCCluster,
		patchHelper:   helper,
		Machine:       params.Machine,
		IBMVPCMachine: params.IBMVPCMachine,
	}, nil
}

// CreateMachine creates a vpc machine.
func (m *MachineScope) CreateMachine() (*vpcv1.Instance, error) {
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

	imageID, err := fetchImageID(m.IBMVPCMachine.Spec.Image, m)
	if err != nil {
		record.Warnf(m.IBMVPCMachine, "FailedRetriveImage", "Failed image retrival - %v", err)
		return nil, fmt.Errorf("error while fetching image ID: %v", err)
	}

	options := &vpcv1.CreateInstanceOptions{}
	instancePrototype := &vpcv1.InstancePrototype{
		Name: &m.IBMVPCMachine.Name,
		Image: &vpcv1.ImageIdentity{
			ID: imageID,
		},
		Profile: &vpcv1.InstanceProfileIdentity{
			Name: &m.IBMVPCMachine.Spec.Profile,
		},
		Zone: &vpcv1.ZoneIdentity{
			Name: &m.IBMVPCMachine.Spec.Zone,
		},
		PrimaryNetworkInterface: &vpcv1.NetworkInterfacePrototype{
			Subnet: &vpcv1.SubnetIdentity{
				ID: &m.IBMVPCMachine.Spec.PrimaryNetworkInterface.Subnet,
			},
		},
		ResourceGroup: &vpcv1.ResourceGroupIdentity{
			ID: &m.IBMVPCCluster.Spec.ResourceGroup,
		},
		UserData: &cloudInitData,
	}

	if m.IBMVPCMachine.Spec.SSHKeys != nil {
		instancePrototype.Keys = []vpcv1.KeyIdentityIntf{}
		for _, sshKey := range m.IBMVPCMachine.Spec.SSHKeys {
			keyID, err := fetchKeyID(sshKey, m)
			if err != nil {
				return nil, fmt.Errorf("error while fetching SSHKey: %v error: %v", sshKey, err)
			}
			key := &vpcv1.KeyIdentity{
				ID: keyID,
			}
			instancePrototype.Keys = append(instancePrototype.Keys, key)
		}
	}

	if m.IBMVPCMachine.Spec.BootVolume != nil {
		instancePrototype.BootVolumeAttachment = volumeToVPCVolumeAttachment(m.IBMVPCMachine.Spec.BootVolume)
	}

	options.SetInstancePrototype(instancePrototype)
	instance, _, err := m.IBMVPCClient.CreateInstance(options)
	if err != nil {
		record.Warnf(m.IBMVPCMachine, "FailedCreateInstance", "Failed instance creation - %v", err)
	} else {
		record.Eventf(m.IBMVPCMachine, "SuccessfulCreateInstance", "Created Instance %q", *instance.Name)
	}
	return instance, err
}

func volumeToVPCVolumeAttachment(volume *infrav1beta2.VPCVolume) *vpcv1.VolumeAttachmentPrototypeInstanceByImageContext {
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

// CreateVPCLoadBalancerPoolMember creates a new pool member and adds it to the load balancer pool.
func (m *MachineScope) CreateVPCLoadBalancerPoolMember(internalIP *string, targetPort int64) (*vpcv1.LoadBalancerPoolMember, error) {
	loadBalancer, _, err := m.IBMVPCClient.GetLoadBalancer(&vpcv1.GetLoadBalancerOptions{
		ID: m.IBMVPCCluster.Status.VPCEndpoint.LBID,
	})
	if err != nil {
		return nil, err
	}

	if *loadBalancer.ProvisioningStatus != string(infrav1beta2.VPCLoadBalancerStateActive) {
		return nil, fmt.Errorf("load balancer is not in active state")
	}

	if len(loadBalancer.Pools) == 0 {
		return nil, fmt.Errorf("no pools exist for the load balancer")
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
		return nil, fmt.Errorf("failed to bind ListLoadBalancerPoolMembers to control plane %s/%s: %w", m.IBMVPCMachine.Namespace, m.IBMVPCMachine.Name, err)
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
		listImagesOptions := &vpcv1.ListImagesOptions{
			ResourceGroupID: &m.IBMVPCCluster.Spec.ResourceGroup,
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

// SetProviderID will set the provider id for the machine.
func (m *MachineScope) SetProviderID(id *string) error {
	// Based on the ProviderIDFormat version the providerID format will be decided.
	if options.ProviderIDFormatType(options.ProviderIDFormat) == options.ProviderIDFormatV2 {
		accountID, err := utils.GetAccountID()
		if err != nil {
			m.Logger.Error(err, "failed to get cloud account id", err.Error())
			return err
		}
		m.IBMVPCMachine.Spec.ProviderID = ptr.To(fmt.Sprintf("ibm://%s///%s/%s", accountID, m.Machine.Spec.ClusterName, *id))
	} else {
		m.IBMVPCMachine.Spec.ProviderID = ptr.To(fmt.Sprintf("ibmvpc://%s/%s", m.Machine.Spec.ClusterName, m.IBMVPCMachine.Name))
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
