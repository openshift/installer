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
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"

	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_p_vm_instances"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"
	"k8s.io/utils/pointer"

	"sigs.k8s.io/controller-runtime/pkg/client"

	capiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
	capierrors "sigs.k8s.io/cluster-api/errors"
	"sigs.k8s.io/cluster-api/util/patch"

	infrav1beta2 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/powervs"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/resourcecontroller"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/endpoints"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/options"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/record"
)

// PowerVSMachineScopeParams defines the input parameters used to create a new PowerVSMachineScope.
type PowerVSMachineScopeParams struct {
	Logger            logr.Logger
	Client            client.Client
	Cluster           *capiv1beta1.Cluster
	Machine           *capiv1beta1.Machine
	IBMPowerVSCluster *infrav1beta2.IBMPowerVSCluster
	IBMPowerVSMachine *infrav1beta2.IBMPowerVSMachine
	IBMPowerVSImage   *infrav1beta2.IBMPowerVSImage
	ServiceEndpoint   []endpoints.ServiceEndpoint
	DHCPIPCacheStore  cache.Store
}

// PowerVSMachineScope defines a scope defined around a Power VS Machine.
type PowerVSMachineScope struct {
	logr.Logger
	Client      client.Client
	patchHelper *patch.Helper

	IBMPowerVSClient  powervs.PowerVS
	Cluster           *capiv1beta1.Cluster
	Machine           *capiv1beta1.Machine
	IBMPowerVSCluster *infrav1beta2.IBMPowerVSCluster
	IBMPowerVSMachine *infrav1beta2.IBMPowerVSMachine
	IBMPowerVSImage   *infrav1beta2.IBMPowerVSImage
	ServiceEndpoint   []endpoints.ServiceEndpoint
	DHCPIPCacheStore  cache.Store
}

// NewPowerVSMachineScope creates a new PowerVSMachineScope from the supplied parameters.
func NewPowerVSMachineScope(params PowerVSMachineScopeParams) (scope *PowerVSMachineScope, err error) {
	scope = &PowerVSMachineScope{}

	if params.Client == nil {
		err = errors.New("client is required when creating a MachineScope")
		return nil, err
	}
	scope.Client = params.Client

	if params.Machine == nil {
		err = errors.New("machine is required when creating a MachineScope")
		return nil, err
	}
	scope.Machine = params.Machine

	if params.Cluster == nil {
		err = errors.New("cluster is required when creating a MachineScope")
		return nil, err
	}
	scope.Cluster = params.Cluster

	if params.IBMPowerVSMachine == nil {
		err = errors.New("PowerVS machine is required when creating a MachineScope")
		return nil, err
	}
	scope.IBMPowerVSMachine = params.IBMPowerVSMachine
	scope.IBMPowerVSCluster = params.IBMPowerVSCluster
	scope.IBMPowerVSImage = params.IBMPowerVSImage

	if params.Logger == (logr.Logger{}) {
		params.Logger = klogr.New()
	}
	scope.Logger = params.Logger

	helper, err := patch.NewHelper(params.IBMPowerVSMachine, params.Client)
	if err != nil {
		err = errors.Wrap(err, "failed to init patch helper")
		return nil, err
	}
	scope.patchHelper = helper

	m := params.IBMPowerVSMachine

	rc, err := resourcecontroller.NewService(resourcecontroller.ServiceOptions{})
	if err != nil {
		return nil, err
	}

	// Fetch the resource controller endpoint.
	if rcEndpoint := endpoints.FetchRCEndpoint(params.ServiceEndpoint); rcEndpoint != "" {
		if err := rc.SetServiceURL(rcEndpoint); err != nil {
			return nil, errors.Wrap(err, "failed to set resource controller endpoint")
		}
		scope.Logger.V(3).Info("Overriding the default resource controller endpoint")
	}

	res, _, err := rc.GetResourceInstance(
		&resourcecontrollerv2.GetResourceInstanceOptions{
			ID: core.StringPtr(m.Spec.ServiceInstanceID),
		})
	if err != nil {
		err = errors.Wrap(err, "failed to get resource instance")
		return nil, err
	}

	region := endpoints.CostructRegionFromZone(*res.RegionID)
	scope.SetRegion(region)
	scope.SetZone(*res.RegionID)

	serviceOptions := powervs.ServiceOptions{
		IBMPIOptions: &ibmpisession.IBMPIOptions{
			Debug: params.Logger.V(DEBUGLEVEL).Enabled(),
			Zone:  *res.RegionID,
		},
		CloudInstanceID: m.Spec.ServiceInstanceID,
	}

	// Fetch the service endpoint.
	if svcEndpoint := endpoints.FetchPVSEndpoint(region, params.ServiceEndpoint); svcEndpoint != "" {
		serviceOptions.IBMPIOptions.URL = svcEndpoint
		scope.Logger.V(3).Info("Overriding the default powervs service endpoint")
	}

	c, err := powervs.NewService(serviceOptions)
	if err != nil {
		err = fmt.Errorf("failed to create PowerVS service")
		return nil, err
	}
	scope.IBMPowerVSClient = c
	scope.DHCPIPCacheStore = params.DHCPIPCacheStore
	return scope, nil
}

func (m *PowerVSMachineScope) ensureInstanceUnique(instanceName string) (*models.PVMInstanceReference, error) {
	instances, err := m.IBMPowerVSClient.GetAllInstance()
	if err != nil {
		return nil, err
	}
	for _, ins := range instances.PvmInstances {
		if *ins.ServerName == instanceName {
			return ins, nil
		}
	}
	return nil, nil
}

// CreateMachine creates a powervs machine.
func (m *PowerVSMachineScope) CreateMachine() (*models.PVMInstanceReference, error) {
	s := m.IBMPowerVSMachine.Spec

	instanceReply, err := m.ensureInstanceUnique(m.IBMPowerVSMachine.Name)
	if err != nil {
		return nil, err
	} else if instanceReply != nil {
		// TODO need a reasonable wrapped error.
		return instanceReply, nil
	}

	// Check if create request has been already triggered.
	// If InstanceReadyCondition is Unknown then return and wait for it to get updated.
	for _, con := range m.IBMPowerVSMachine.Status.Conditions {
		if con.Type == infrav1beta2.InstanceReadyCondition && con.Status == corev1.ConditionUnknown {
			return nil, nil
		}
	}

	cloudInitData, err := m.GetBootstrapData()
	if err != nil {
		return nil, err
	}

	memory := float64(s.MemoryGiB)

	var processors float64
	switch s.Processors.Type {
	case intstr.Int:
		processors = float64(s.Processors.IntVal)
	case intstr.String:
		processors, err = strconv.ParseFloat(s.Processors.StrVal, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to convert Processors(%s) to float64", s.Processors.StrVal)
		}
	}

	var imageID *string
	if m.IBMPowerVSImage != nil {
		imageID = &m.IBMPowerVSImage.Status.ImageID
	} else {
		imageID, err = getImageID(s.Image, m)
		if err != nil {
			record.Warnf(m.IBMPowerVSMachine, "FailedRetriveImage", "Failed image retrival - %v", err)
			return nil, fmt.Errorf("error getting image ID: %v", err)
		}
	}

	networkID, err := getNetworkID(s.Network, m)
	if err != nil {
		record.Warnf(m.IBMPowerVSMachine, "FailedRetrieveNetwork", "Failed network retrieval - %v", err)
		return nil, fmt.Errorf("error getting network ID: %v", err)
	}

	procType := strings.ToLower(string(s.ProcessorType))

	params := &p_cloud_p_vm_instances.PcloudPvminstancesPostParams{
		Body: &models.PVMInstanceCreate{
			ImageID:     imageID,
			KeyPairName: s.SSHKey,
			Networks: []*models.PVMInstanceAddNetwork{
				{
					NetworkID: networkID,
					//IPAddress: address,
				},
			},
			ServerName: &m.IBMPowerVSMachine.Name,
			Memory:     &memory,
			Processors: &processors,
			ProcType:   &procType,
			SysType:    s.SystemType,
			UserData:   cloudInitData,
		},
	}
	_, err = m.IBMPowerVSClient.CreateInstance(params.Body)
	if err != nil {
		record.Warnf(m.IBMPowerVSMachine, "FailedCreateInstance", "Failed instance creation - %v", err)
		return nil, err
	}
	record.Eventf(m.IBMPowerVSMachine, "SuccessfulCreateInstance", "Created Instance %q", m.IBMPowerVSMachine.Name)
	return nil, nil
}

// Close closes the current scope persisting the cluster configuration and status.
func (m *PowerVSMachineScope) Close() error {
	return m.PatchObject()
}

// PatchObject persists the cluster configuration and status.
func (m *PowerVSMachineScope) PatchObject() error {
	return m.patchHelper.Patch(context.TODO(), m.IBMPowerVSMachine)
}

// DeleteMachine deletes the power vs machine associated with machine instance id and service instance id.
func (m *PowerVSMachineScope) DeleteMachine() error {
	if err := m.IBMPowerVSClient.DeleteInstance(m.IBMPowerVSMachine.Status.InstanceID); err != nil {
		record.Warnf(m.IBMPowerVSMachine, "FailedDeleteInstance", "Failed instance deletion - %v", err)
		return err
	}
	record.Eventf(m.IBMPowerVSMachine, "SuccessfulDeleteInstance", "Deleted Instance %q", m.IBMPowerVSMachine.Name)
	return nil
}

// GetBootstrapData returns the base64 encoded bootstrap data from the secret in the Machine's bootstrap.dataSecretName.
func (m *PowerVSMachineScope) GetBootstrapData() (string, error) {
	if m.Machine.Spec.Bootstrap.DataSecretName == nil {
		return "", errors.New("error retrieving bootstrap data: linked Machine's bootstrap.dataSecretName is nil")
	}

	secret := &corev1.Secret{}
	key := types.NamespacedName{Namespace: m.Machine.Namespace, Name: *m.Machine.Spec.Bootstrap.DataSecretName}
	if err := m.Client.Get(context.TODO(), key, secret); err != nil {
		return "", errors.Wrapf(err, "failed to retrieve bootstrap data secret for IBMPowerVSMachine %v", klog.KObj(m.Machine))
	}

	value, ok := secret.Data["value"]
	if !ok {
		return "", errors.New("error retrieving bootstrap data: secret value key is missing")
	}

	return base64.StdEncoding.EncodeToString(value), nil
}

func getImageID(image *infrav1beta2.IBMPowerVSResourceReference, m *PowerVSMachineScope) (*string, error) {
	if image.ID != nil {
		return image.ID, nil
	} else if image.Name != nil {
		images, err := m.GetImages()
		if err != nil {
			m.Logger.Error(err, "Failed to get images")
			return nil, err
		}
		for _, img := range images.Images {
			if *image.Name == *img.Name {
				m.Logger.Info("Image found with ID", "Image", *image.Name, "ID", *img.ImageID)
				return img.ImageID, nil
			}
		}
	} else {
		return nil, fmt.Errorf("both ID and Name can't be nil")
	}
	return nil, fmt.Errorf("failed to find an image ID")
}

// GetImages will get list of images for the powervs service instance.
func (m *PowerVSMachineScope) GetImages() (*models.Images, error) {
	return m.IBMPowerVSClient.GetAllImage()
}

func getNetworkID(network infrav1beta2.IBMPowerVSResourceReference, m *PowerVSMachineScope) (*string, error) {
	if network.ID != nil {
		return network.ID, nil
	} else if network.Name != nil {
		networks, err := m.GetNetworks()
		if err != nil {
			m.Logger.Error(err, "Failed to get networks")
			return nil, err
		}
		for _, nw := range networks.Networks {
			if *network.Name == *nw.Name {
				m.Logger.Info("Network found with ID", "Network", *network.Name, "ID", *nw.NetworkID)
				return nw.NetworkID, nil
			}
		}
		return nil, fmt.Errorf("failed to find a network ID with name %s", *network.Name)
	} else if network.RegEx != nil {
		networks, err := m.GetNetworks()
		if err != nil {
			m.Logger.Error(err, "Failed to get networks")
			return nil, err
		}
		re, err := regexp.Compile(*network.RegEx)
		if err != nil {
			m.Logger.Error(err, "Failed to compile regular expression", "regex", *network.RegEx)
			return nil, err
		}
		// In case of multiple network names matches the provided regular expression the first matched network will be selected.
		for _, nw := range networks.Networks {
			if match := re.Match([]byte(*nw.Name)); match {
				m.Logger.Info("Network found with ID", "Network", *nw.Name, "ID", *nw.NetworkID)
				return nw.NetworkID, nil
			}
		}
		return nil, fmt.Errorf("failed to find a network ID with RegEx %s", *network.RegEx)
	} else {
		return nil, fmt.Errorf("ID, Name and RegEx can't be nil")
	}
}

// GetNetworks will get list of networks for the powervs service instance.
func (m *PowerVSMachineScope) GetNetworks() (*models.Networks, error) {
	return m.IBMPowerVSClient.GetAllNetwork()
}

// SetReady will set the status as ready for the machine.
func (m *PowerVSMachineScope) SetReady() {
	m.IBMPowerVSMachine.Status.Ready = true
}

// SetNotReady will set status as not ready for the machine.
func (m *PowerVSMachineScope) SetNotReady() {
	m.IBMPowerVSMachine.Status.Ready = false
}

// SetFailureReason will set status FailureReason for the machine.
func (m *PowerVSMachineScope) SetFailureReason(reason capierrors.MachineStatusError) {
	m.IBMPowerVSMachine.Status.FailureReason = &reason
}

// SetFailureMessage will set status FailureMessage for the machine.
func (m *PowerVSMachineScope) SetFailureMessage(message string) {
	m.IBMPowerVSMachine.Status.FailureMessage = &message
}

// IsReady will return the status for the machine.
func (m *PowerVSMachineScope) IsReady() bool {
	return m.IBMPowerVSMachine.Status.Ready
}

// SetInstanceID will set the instance id for the machine.
func (m *PowerVSMachineScope) SetInstanceID(id *string) {
	if id != nil {
		m.IBMPowerVSMachine.Status.InstanceID = *id
	}
}

// GetInstanceID will get the instance id for the machine.
func (m *PowerVSMachineScope) GetInstanceID() string {
	return m.IBMPowerVSMachine.Status.InstanceID
}

// SetHealth will set the health status for the machine.
func (m *PowerVSMachineScope) SetHealth(health *models.PVMInstanceHealth) {
	if health != nil {
		m.IBMPowerVSMachine.Status.Health = health.Status
	}
}

// SetAddresses will set the addresses for the machine.
func (m *PowerVSMachineScope) SetAddresses(instance *models.PVMInstance) {
	var addresses []corev1.NodeAddress
	// Setting the name of the vm to the InternalDNS and Hostname as the vm uses that as hostname.
	addresses = append(addresses, corev1.NodeAddress{
		Type:    corev1.NodeInternalDNS,
		Address: *instance.ServerName,
	})
	addresses = append(addresses, corev1.NodeAddress{
		Type:    corev1.NodeHostName,
		Address: *instance.ServerName,
	})
	for _, network := range instance.Networks {
		if strings.TrimSpace(network.IPAddress) != "" {
			addresses = append(addresses, corev1.NodeAddress{
				Type:    corev1.NodeInternalIP,
				Address: strings.TrimSpace(network.IPAddress),
			})
		}
		if strings.TrimSpace(network.ExternalIP) != "" {
			addresses = append(addresses, corev1.NodeAddress{
				Type:    corev1.NodeExternalIP,
				Address: strings.TrimSpace(network.ExternalIP),
			})
		}
	}
	m.IBMPowerVSMachine.Status.Addresses = addresses
	if len(addresses) > 2 {
		// If the address length is more than 2 means either NodeInternalIP or NodeExternalIP is updated so return
		return
	}
	// In this case there is no IP found under instance.Networks, So try to fetch the IP from cache or DHCP server
	// Look for DHCP IP from the cache
	obj, exists, err := m.DHCPIPCacheStore.GetByKey(*instance.ServerName)
	if err != nil {
		m.Error(err, "Failed to fetch the DHCP IP address from cache store", "VM", *instance.ServerName)
	}
	if exists {
		m.Info("Found IP for VM from DHCP cache", "IP", obj.(powervs.VMip).IP, "VM", *instance.ServerName)
		addresses = append(addresses, corev1.NodeAddress{
			Type:    corev1.NodeInternalIP,
			Address: obj.(powervs.VMip).IP,
		})
		m.IBMPowerVSMachine.Status.Addresses = addresses
		return
	}
	// Fetch the VM network ID
	networkID, err := getNetworkID(m.IBMPowerVSMachine.Spec.Network, m)
	if err != nil {
		m.Error(err, "Failed to fetch network id from network resource", "VM", *instance.ServerName)
		return
	}
	// Fetch the details of the network attached to the VM
	var pvmNetwork *models.PVMInstanceNetwork
	for _, network := range instance.Networks {
		if network.NetworkID == *networkID {
			pvmNetwork = network
			m.Info("Found network attached to VM", "Network ID", network.NetworkID, "VM", *instance.ServerName)
		}
	}
	if pvmNetwork == nil {
		m.Info("Failed to get network attached to VM", "VM", *instance.ServerName, "Network ID", *networkID)
		return
	}
	// Get all the DHCP servers
	dhcpServer, err := m.IBMPowerVSClient.GetAllDHCPServers()
	if err != nil {
		m.Error(err, "Failed to get DHCP server")
		return
	}
	// Get the Details of DHCP server associated with the network
	var dhcpServerDetails *models.DHCPServerDetail
	for _, server := range dhcpServer {
		if *server.Network.ID == *networkID {
			m.Info("found DHCP server for network", "DHCP server ID", *server.ID, "network ID", *networkID)
			dhcpServerDetails, err = m.IBMPowerVSClient.GetDHCPServer(*server.ID)
			if err != nil {
				m.Error(err, "Failed to get DHCP server details", "DHCP server ID", *server.ID)
				return
			}
			break
		}
	}
	if dhcpServerDetails == nil {
		errStr := fmt.Errorf("DHCP server details is nil")
		m.Error(errStr, "DHCP server associated with network is nil", "Network ID", *networkID)
		return
	}

	// Fetch the VM IP using VM's mac from DHCP server lease
	var internalIP *string
	for _, lease := range dhcpServerDetails.Leases {
		if *lease.InstanceMacAddress == pvmNetwork.MacAddress {
			m.Info("Found internal ip for VM from DHCP lease", "IP", *lease.InstanceIP, "VM", *instance.ServerName)
			internalIP = lease.InstanceIP
			break
		}
	}
	if internalIP == nil {
		errStr := fmt.Errorf("internal IP is nil")
		m.Error(errStr, "Failed to get internal IP, DHCP lease not found for VM with MAC in DHCP network", "VM", *instance.ServerName,
			"MAC", pvmNetwork.MacAddress, "DHCP server ID", *dhcpServerDetails.ID)
		return
	}
	m.Info("found internal IP for VM from DHCP lease", "IP", *internalIP, "VM", *instance.ServerName)
	addresses = append(addresses, corev1.NodeAddress{
		Type:    corev1.NodeInternalIP,
		Address: *internalIP,
	})
	// Update the cache with the ip and VM name
	err = m.DHCPIPCacheStore.Add(powervs.VMip{
		Name: *instance.ServerName,
		IP:   *internalIP,
	})
	if err != nil {
		m.Error(err, "Failed to update the DHCP cache store with the IP", "VM", *instance.ServerName, "IP", *internalIP)
	}
	m.IBMPowerVSMachine.Status.Addresses = addresses
}

// SetInstanceState will set the state for the machine.
func (m *PowerVSMachineScope) SetInstanceState(status *string) {
	m.IBMPowerVSMachine.Status.InstanceState = infrav1beta2.PowerVSInstanceState(*status)
}

// GetInstanceState will get the state for the machine.
func (m *PowerVSMachineScope) GetInstanceState() infrav1beta2.PowerVSInstanceState {
	return m.IBMPowerVSMachine.Status.InstanceState
}

// SetRegion will set the region for the machine.
func (m *PowerVSMachineScope) SetRegion(region string) {
	m.IBMPowerVSMachine.Status.Region = &region
}

// GetRegion will get the region for the machine.
func (m *PowerVSMachineScope) GetRegion() string {
	if m.IBMPowerVSMachine.Status.Region == nil {
		return ""
	}
	return *m.IBMPowerVSMachine.Status.Region
}

// SetZone will set the zone for the machine.
func (m *PowerVSMachineScope) SetZone(zone string) {
	m.IBMPowerVSMachine.Status.Zone = &zone
}

// GetZone will get the zone for the machine.
func (m *PowerVSMachineScope) GetZone() string {
	if m.IBMPowerVSMachine.Status.Zone == nil {
		return ""
	}
	return *m.IBMPowerVSMachine.Status.Zone
}

// SetProviderID will set the provider id for the machine.
func (m *PowerVSMachineScope) SetProviderID(id *string) {
	// Based on the ProviderIDFormat version the providerID format will be decided.
	if options.ProviderIDFormatType(options.PowerVSProviderIDFormat) == options.PowerVSProviderIDFormatV2 ||
		options.ProviderIDFormatType(options.ProviderIDFormat) == options.ProviderIDFormatV2 {
		if id != nil {
			m.IBMPowerVSMachine.Spec.ProviderID = pointer.String(fmt.Sprintf("ibmpowervs://%s/%s/%s/%s", m.GetRegion(), m.GetZone(), m.IBMPowerVSMachine.Spec.ServiceInstanceID, *id))
		}
	} else {
		m.IBMPowerVSMachine.Spec.ProviderID = pointer.String(fmt.Sprintf("ibmpowervs://%s/%s", m.Machine.Spec.ClusterName, m.IBMPowerVSMachine.Name))
	}
}
