package azure

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/openshift/installer/pkg/types"
	typesazure "github.com/openshift/installer/pkg/types/azure"
	azuredefaults "github.com/openshift/installer/pkg/types/azure/defaults"
)

// Metadata holds additional metadata for InstallConfig resources that
// does not need to be user-supplied (e.g. because it can be retrieved
// from external APIs).
type Metadata struct {
	session           *Session
	client            API
	dnsCfg            *DNSConfig
	availabilityZones []string
	vmZones           []string
	region            string
	ZonesSubnetMap    map[string][]string

	controlPlane    *types.MachinePool
	compute         *types.MachinePool
	defaultPlatform *typesazure.MachinePool

	controlPlaneCapabilities map[string]string
	computeCapabilities      map[string]string

	// CloudName indicates the Azure cloud environment (e.g. public, gov't).
	CloudName typesazure.CloudEnvironment `json:"cloudName,omitempty"`

	// ARMEndpoint indicates the resource management API endpoint used by AzureStack.
	ARMEndpoint string `json:"armEndpoint,omitempty"`

	// Credentials hold prepopulated Azure credentials.
	// At the moment the installer doesn't use it and reads credentials
	// from the file system, but external consumers of the package can
	// provide credentials. This is useful when we run the installer
	// as a service (Azure Red Hat OpenShift, for example): in this case
	// we do not want to rely on the filesystem or user input as we
	// serve multiple users with different credentials via a web server.
	Credentials *Credentials `json:"credentials,omitempty"`

	mutex sync.Mutex
}

// NewMetadata initializes a new Metadata object.
func NewMetadata(az *typesazure.Platform, controlPlane, compute *types.MachinePool) *Metadata {
	return NewMetadataWithCredentials(az, controlPlane, compute, nil)
}

// NewMetadataWithCredentials initializes a new Metadata object
// with prepopulated Azure credentials.
func NewMetadataWithCredentials(az *typesazure.Platform, controlPlane, compute *types.MachinePool, credentials *Credentials) *Metadata {
	return &Metadata{
		CloudName:       az.CloudName,
		ARMEndpoint:     az.ARMEndpoint,
		Credentials:     credentials,
		controlPlane:    controlPlane,
		compute:         compute,
		region:          az.Region,
		defaultPlatform: az.DefaultMachinePlatform,
	}
}

// Session holds an Azure session which can be used for Azure API calls
// during asset generation.
func (m *Metadata) Session() (*Session, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.unlockedSession()
}

func (m *Metadata) unlockedSession() (*Session, error) {
	if m.session == nil {
		var err error
		m.session, err = GetSessionWithCredentials(m.CloudName, m.ARMEndpoint, m.Credentials)
		if err != nil {
			return nil, fmt.Errorf("creating Azure session: %w", err)
		}
	}

	return m.session, nil
}

// Client holds an Azure Client that implements calls to the Azure API.
func (m *Metadata) Client() (API, error) {
	if m.client == nil {
		ssn, err := m.Session()
		if err != nil {
			return nil, err
		}
		m.client = NewClient(ssn)
	}
	return m.client, nil
}

// UseMockClient returns the provided client from Client() instead of creating
// a new one.
func (m *Metadata) UseMockClient(client API) {
	m.client = client
}

// DNSConfig holds an Azure DNSConfig Client that implements calls to the Azure API.
func (m *Metadata) DNSConfig() (*DNSConfig, error) {
	if m.dnsCfg == nil {
		ssn, err := m.Session()
		if err != nil {
			return nil, err
		}
		m.dnsCfg = NewDNSConfig(ssn)
	}
	return m.dnsCfg, nil
}

// AvailabilityZones retrieves a list of availability zones for the configured region.
func (m *Metadata) AvailabilityZones(ctx context.Context) ([]string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.availabilityZones) == 0 {
		zones, err := m.client.GetRegionAvailabilityZones(ctx, m.region)
		if err != nil {
			return nil, fmt.Errorf("error retrieving Availability Zones: %w", err)
		}
		if zones != nil {
			sort.Strings(zones)
			m.availabilityZones = zones
		}
	}

	return m.availabilityZones, nil
}

// VMAvailabilityZones retrieves a list of availability zones for the configured region and instance type.
func (m *Metadata) VMAvailabilityZones(ctx context.Context, instanceType string) ([]string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.vmZones) == 0 {
		zones, err := m.client.GetAvailabilityZones(ctx, m.region, instanceType)
		if err != nil {
			return nil, fmt.Errorf("error retrieving Availability Zones: %w", err)
		}
		if zones != nil {
			sort.Strings(zones)
			m.vmZones = zones
		}
	}

	return m.vmZones, nil
}

// GenerateZonesSubnetMap creates a map of all the zones that are supported for nat gateways and vms and
// sets it to the subnets provided. If no subnets are provided, it creates subnets for multi zone
// functionality.
func (m *Metadata) GenerateZonesSubnetMap(subnetSpec []typesazure.SubnetSpec, defaultComputeSubnet string) (map[string][]string, error) {
	if m.ZonesSubnetMap == nil {
		// Get the availability zones.
		if m.availabilityZones == nil {
			_, err := m.AvailabilityZones(context.TODO())
			if err != nil {
				return nil, err
			}
		}
		subnetZones := m.availabilityZones
		computeSubnets := []string{}

		// Get all the byo subnets or generate subnet per az.
		if len(subnetSpec) != 0 {
			sort.Slice(subnetSpec, func(i, j int) bool {
				return subnetSpec[i].Name < subnetSpec[j].Name
			})
			for _, subnet := range subnetSpec {
				if subnet.Role == typesazure.SubnetNode {
					computeSubnets = append(computeSubnets, subnet.Name)
				}
			}
		} else {
			for idx := range subnetZones {
				computeName := fmt.Sprintf("%s-%d", defaultComputeSubnet, idx+1)
				if idx == 0 {
					computeName = defaultComputeSubnet
				}
				computeSubnets = append(computeSubnets, computeName)
			}
		}

		// Assign zone to subnets.
		subnetMap := map[string][]string{}
		zoneIndex := 0
		for _, subnet := range computeSubnets {
			if _, ok := subnetMap[subnetZones[zoneIndex]]; !ok {
				subnetMap[subnetZones[zoneIndex]] = []string{}
			}
			subnetMap[subnetZones[zoneIndex]] = append(subnetMap[subnetZones[zoneIndex]], subnet)
			zoneIndex++
			if zoneIndex >= len(subnetZones) {
				zoneIndex = 0
			}
		}
		m.ZonesSubnetMap = subnetMap
	}
	return m.ZonesSubnetMap, nil
}

// ControlPlaneCapabilities returns the capabilities for the instance type of control-plane
// nodes from the Azure API.
func (m *Metadata) ControlPlaneCapabilities() (map[string]string, error) {
	if m.controlPlaneCapabilities == nil {
		caps, err := m.getCapabilities(m.controlPlane, azuredefaults.ControlPlaneInstanceType)
		if err != nil {
			return nil, fmt.Errorf("failed to get control plane capabilities: %w", err)
		}
		m.controlPlaneCapabilities = caps
	}
	return m.controlPlaneCapabilities, nil
}

// ComputeCapabilities returns the capabilities for the instance type of compute nodes
// from the Azure API.
func (m *Metadata) ComputeCapabilities() (map[string]string, error) {
	if m.computeCapabilities == nil {
		caps, err := m.getCapabilities(m.compute, azuredefaults.ComputeInstanceType)
		if err != nil {
			return nil, fmt.Errorf("failed to get compute capabilities: %w", err)
		}
		m.computeCapabilities = caps
	}
	return m.computeCapabilities, nil
}

// ControlPlaneHyperVGeneration returns the HyperVGeneration for the control-plane
// instances. If the instance type supports both V1 & V2, V2 is returned.
func (m *Metadata) ControlPlaneHyperVGeneration() (string, error) {
	caps, err := m.ControlPlaneCapabilities()
	if err != nil {
		return "", fmt.Errorf("unable to get control plane capabilities: %w", err)
	}
	return GetHyperVGenerationVersion(caps, "")
}

// ComputeHyperVGeneration returns the HyperVGeneration for the compute instances.
// If the instance type supports both V1 & V2, V2 is returned.
func (m *Metadata) ComputeHyperVGeneration() (string, error) {
	caps, err := m.ComputeCapabilities()
	if err != nil {
		return "", fmt.Errorf("unable to get compute capabilities: %w", err)
	}
	return GetHyperVGenerationVersion(caps, "")
}

type machinePoolInstanceTypeFunc func(typesazure.CloudEnvironment, string, types.Architecture) string

func (m *Metadata) getCapabilities(mPool *types.MachinePool, defaultInstanceType machinePoolInstanceTypeFunc) (map[string]string, error) {
	if mPool == nil {
		return nil, fmt.Errorf("unable to get capabilities because machinepool is not populated in metadata")
	}
	instType := defaultInstanceType(m.CloudName, m.region, mPool.Architecture)
	if dmp := m.defaultPlatform; dmp != nil && dmp.InstanceType != "" {
		instType = dmp.InstanceType
	}
	if mPool.Platform.Azure != nil && mPool.Platform.Azure.InstanceType != "" {
		instType = mPool.Platform.Azure.InstanceType
	}
	client, err := m.Client()
	if err != nil {
		return nil, fmt.Errorf("failed to get azure client %w", err)
	}
	caps, err := client.GetVMCapabilities(context.TODO(), instType, m.region)
	if err != nil {
		return nil, err
	}
	return caps, nil
}
