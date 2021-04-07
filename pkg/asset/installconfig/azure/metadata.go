package azure

import (
	"sync"

	"github.com/pkg/errors"

	typesazure "github.com/openshift/installer/pkg/types/azure"
)

// Metadata holds additional metadata for InstallConfig resources that
// does not need to be user-supplied (e.g. because it can be retrieved
// from external APIs).
type Metadata struct {
	session *Session
	client  *Client
	dnsCfg  *DNSConfig

	// CloudName indicates the Azure cloud environment (e.g. public, gov't).
	CloudName typesazure.CloudEnvironment `json:"cloudName,omitempty"`

	// ARMEndpoint indicates the resource management API endpoint used by AzureStack.
	ARMEndpoint string `json:"armEndpoint,omitempty"`

	mutex sync.Mutex
}

// NewMetadata initializes a new Metadata object.
func NewMetadata(cloudName typesazure.CloudEnvironment, armEndpoint string) *Metadata {
	return &Metadata{CloudName: cloudName, ARMEndpoint: armEndpoint}
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
		m.session, err = GetSession(m.CloudName, m.ARMEndpoint)
		if err != nil {
			return nil, errors.Wrap(err, "creating Azure session")
		}
	}

	return m.session, nil
}

// Client holds an Azure Client that implements calls to the Azure API.
func (m *Metadata) Client() (*Client, error) {
	if m.client == nil {
		ssn, err := m.Session()
		if err != nil {
			return nil, err
		}
		m.client = NewClient(ssn)
	}
	return m.client, nil
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
