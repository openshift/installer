package azure

import (
	"context"
	"fmt"
	"sort"
	"sync"

	typesazure "github.com/openshift/installer/pkg/types/azure"
)

// Metadata holds additional metadata for InstallConfig resources that
// does not need to be user-supplied (e.g. because it can be retrieved
// from external APIs).
type Metadata struct {
	session           *Session
	client            API
	dnsCfg            *DNSConfig
	availabilityZones []string
	region            string

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
func NewMetadata(cloudName typesazure.CloudEnvironment, armEndpoint string, region string) *Metadata {
	return NewMetadataWithCredentials(cloudName, armEndpoint, nil, region)
}

// NewMetadataWithCredentials initializes a new Metadata object
// with prepopulated Azure credentials.
func NewMetadataWithCredentials(cloudName typesazure.CloudEnvironment, armEndpoint string, credentials *Credentials, region string) *Metadata {
	return &Metadata{
		CloudName:   cloudName,
		ARMEndpoint: armEndpoint,
		Credentials: credentials,
		region:      region,
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
