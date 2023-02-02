package ovirtclient

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"testing"
	"time"

	ovirtclientlog "github.com/ovirt/go-ovirt-client-log/v3"
)

// TestHelper is a helper to run tests against an oVirt engine. When created it scans the oVirt Engine and tries to find
// working resources (hosts, clusters, etc) for running tests against. Tests should clean up after themselves.
type TestHelper interface {
	// GetClient returns the goVirt client.
	GetClient() Client

	// GetClusterID returns the ID for the cluster.
	GetClusterID() ClusterID

	// GetBlankTemplateID returns the ID of the blank template that can be used for creating dummy VMs.
	GetBlankTemplateID() TemplateID

	// GetStorageDomainID returns the ID of the storage domain to create the images on.
	GetStorageDomainID() StorageDomainID

	// GetSecondaryStorageDomainID returns a second ID of a storage domain to create images on. If no secondary
	// storage domain is available, the test will be skipped.
	GetSecondaryStorageDomainID(t *testing.T) StorageDomainID

	// GenerateRandomID generates a random ID for testing.
	GenerateRandomID(length uint) string

	// GetVNICProfileID returns a VNIC profile ID for testing.
	GetVNICProfileID() VNICProfileID

	// GetTLS returns the TLS provider used for this test helper.
	GetTLS() TLSProvider

	// GetUsername returns the oVirt username.
	GetUsername() string

	// GetPassword returns the oVirt password.
	GetPassword() string

	// GenerateTestResourceName generates a test resource name from the test name.
	GenerateTestResourceName(t *testing.T) string
}

// MustNewTestHelper is identical to NewTestHelper, but panics instead of returning an error.
func MustNewTestHelper(
	username string,
	password string,
	url string,
	tlsProvider TLSProvider,
	mock bool,
	logger ovirtclientlog.Logger,
	params TestHelperParameters,
) TestHelper {
	helper, err := NewTestHelper(
		url,
		username,
		password,
		params,
		tlsProvider,
		mock,
		logger,
	)
	if err != nil {
		panic(err)
	}
	return helper
}

// NewTestHelper creates a helper for executing tests. Depending on the mock parameter, this either sets up a mock oVirt
// client with test fixtures (host, cluster, etc), or connects a real oVirt cluster and finds usable infrastructure
// to test on. The returned object provides helper functions to access these parameters.
//
// The ID parameters are optional and trigger automatic detection if left empty.
func NewTestHelper(
	url string,
	username string,
	password string,
	params TestHelperParameters,
	tlsProvider TLSProvider,
	mock bool,
	logger ovirtclientlog.Logger,
) (TestHelper, error) {
	client, err := createTestClient(url, username, password, tlsProvider, mock, logger)
	if err != nil {
		return nil, err
	}

	if params == nil {
		params = &testHelperParameters{}
	}

	clusterID, err := setupTestClusterID(params.ClusterID(), client)
	if err != nil {
		return nil, err
	}

	storageDomainID, err := setupTestStorageDomainID(params.StorageDomainID(), client)
	if err != nil {
		return nil, err
	}

	secondaryStorageDomainID, err := setupSecondaryStorageDomainID(params.SecondaryStorageDomainID(), storageDomainID, client)
	if err != nil {
		return nil, err
	}

	blankTemplateID, err := setupBlankTemplateID(params.BlankTemplateID(), client)
	if err != nil {
		return nil, err
	}

	vnicProfileID, err := setupVNICProfileID(params.VNICProfileID(), clusterID, client)
	if err != nil {
		return nil, err
	}

	return &testHelper{
		username:                 username,
		password:                 password,
		client:                   client,
		tls:                      tlsProvider,
		clusterID:                clusterID,
		storageDomainID:          storageDomainID,
		secondaryStorageDomainID: secondaryStorageDomainID,
		blankTemplateID:          blankTemplateID,
		vnicProfileID:            vnicProfileID,
		// We are suppressing gosec linting here since rand is not used in a security-relevant context,
		// only to generate random ID's for testing.
		rand: rand.New(rand.NewSource(time.Now().UnixNano())), //nolint:gosec
	}, nil
}

// TestHelperParams creates a new copy of BuildableTestHelperParameters usable for building test parameters.
func TestHelperParams() BuildableTestHelperParameters {
	return &testHelperParameters{}
}

// TestHelperParameters contains the optional parameters for the test helper.
type TestHelperParameters interface {
	// ClusterID returns the cluster ID used for testing. It can return an empty string if no
	// test cluster is designated, in which case a cluster is selected.
	ClusterID() ClusterID

	// StorageDomainID returns the storage domain ID usable for testing. It can return an empty
	// string if no test storage domain is designated for testing, in which case a working
	// storage domain is selected.
	StorageDomainID() StorageDomainID

	// SecondaryStorageDomainID returns a second ID for a storage domain to test features that require
	// two storage domains. This may be empty, in which case such tests should be skipped.
	SecondaryStorageDomainID() StorageDomainID

	// BlankTemplateID returns an ID to a template that is blank and can be used as a basis
	// for testing. It may return an empty string if no template is provided.
	BlankTemplateID() TemplateID

	// VNICProfileID returns an ID to a VNIC profile designated for testing. It may return
	// an empty string, in which case an arbitrary VNIC profile is selected.
	VNICProfileID() VNICProfileID
}

// BuildableTestHelperParameters is a buildable version of the TestHelperParameters.
type BuildableTestHelperParameters interface {
	TestHelperParameters

	// WithClusterID sets the cluster ID usable for testing.
	WithClusterID(ClusterID) BuildableTestHelperParameters
	// WithStorageDomainID sets the storage domain that can be used for testing.
	WithStorageDomainID(StorageDomainID) BuildableTestHelperParameters
	// WithSecondaryStorageDomainID sets the storage domain that can be used for testing, which is not identical to
	// the primary storage domain ID.
	WithSecondaryStorageDomainID(StorageDomainID) BuildableTestHelperParameters
	// WithBlankTemplateID sets the blank template that can be used for testing.
	WithBlankTemplateID(TemplateID) BuildableTestHelperParameters
	// WithVNICProfileID sets the ID of the VNIC profile that can be used for testing.
	WithVNICProfileID(VNICProfileID) BuildableTestHelperParameters
}

type testHelperParameters struct {
	clusterID                ClusterID
	storageDomainID          StorageDomainID
	secondaryStorageDomainID StorageDomainID
	blankTemplateID          TemplateID
	vnicProfileID            VNICProfileID
}

func (t *testHelperParameters) WithSecondaryStorageDomainID(s StorageDomainID) BuildableTestHelperParameters {
	t.secondaryStorageDomainID = s
	return t
}

func (t *testHelperParameters) ClusterID() ClusterID {
	return t.clusterID
}

func (t *testHelperParameters) StorageDomainID() StorageDomainID {
	return t.storageDomainID
}

func (t *testHelperParameters) SecondaryStorageDomainID() StorageDomainID {
	return t.secondaryStorageDomainID
}

func (t *testHelperParameters) BlankTemplateID() TemplateID {
	return t.blankTemplateID
}

func (t *testHelperParameters) VNICProfileID() VNICProfileID {
	return t.vnicProfileID
}

func (t *testHelperParameters) WithClusterID(s ClusterID) BuildableTestHelperParameters {
	t.clusterID = s
	return t
}

func (t *testHelperParameters) WithStorageDomainID(s StorageDomainID) BuildableTestHelperParameters {
	t.storageDomainID = s
	return t
}

func (t *testHelperParameters) WithBlankTemplateID(s TemplateID) BuildableTestHelperParameters {
	t.blankTemplateID = s
	return t
}

func (t *testHelperParameters) WithVNICProfileID(s VNICProfileID) BuildableTestHelperParameters {
	t.vnicProfileID = s
	return t
}

func setupVNICProfileID(vnicProfileID VNICProfileID, clusterID ClusterID, client Client) (VNICProfileID, error) {
	if vnicProfileID != "" {
		_, err := client.GetVNICProfile(vnicProfileID)
		if err != nil {
			return "", fmt.Errorf("failed to verify VNIC profile ID %s", vnicProfileID)
		}
		return vnicProfileID, nil
	}
	vnicProfiles, err := client.ListVNICProfiles()
	if err != nil {
		return "", fmt.Errorf("failed to list VNIC profiles (%w)", err)
	}
	for _, vnicProfile := range vnicProfiles {
		network, err := vnicProfile.Network()
		if err != nil {
			return "", fmt.Errorf("failed to fetch network %s (%w)", vnicProfile.NetworkID(), err)
		}
		dc, err := network.Datacenter()
		if err != nil {
			return "", fmt.Errorf("failed to fetch datacenter from network %s (%w)", network.ID(), err)
		}
		hasCluster, err := dc.HasCluster(clusterID)
		if err != nil {
			return "", fmt.Errorf("failed to get datacenter clusters for %s", dc.ID())
		}
		if hasCluster {
			return vnicProfile.ID(), nil
		}
	}
	return "", fmt.Errorf("failed to find a valid VNIC profile ID for testing")
}

func setupBlankTemplateID(blankTemplateID TemplateID, client Client) (id TemplateID, err error) {
	if blankTemplateID == "" {
		blankTemplateID, err = findBlankTemplateID(client)
		if err != nil {
			return "", fmt.Errorf("failed to find blank template (%w)", err)
		}
	} else if err := verifyBlankTemplateID(client, blankTemplateID); err != nil {
		return "", fmt.Errorf("failed to verify blank template ID %s (%w)", blankTemplateID, err)
	}
	return blankTemplateID, nil
}

func setupTestStorageDomainID(storageDomainID StorageDomainID, client Client) (id StorageDomainID, err error) {
	if storageDomainID == "" {
		storageDomainID, err = findTestStorageDomainID("", client)
		if err != nil && !errors.Is(err, errNoTestStorageDomainFound) {
			return "", fmt.Errorf("failed to find storage domain to test on (%w)", err)
		}
	} else if err := verifyTestStorageDomainID(client, storageDomainID); err != nil {
		return "", fmt.Errorf("failed to verify storage domain ID %s (%w)", storageDomainID, err)
	}
	return storageDomainID, nil
}

func setupSecondaryStorageDomainID(
	storageDomainID StorageDomainID,
	skipStorageDomainID StorageDomainID,
	client Client,
) (id StorageDomainID, err error) {
	if storageDomainID == "" {
		storageDomainID, err = findTestStorageDomainID(skipStorageDomainID, client)
		if err != nil && !errors.Is(err, errNoTestStorageDomainFound) {
			return "", fmt.Errorf("failed to find storage domain to test on (%w)", err)
		}
	} else if err := verifyTestStorageDomainID(client, storageDomainID); err != nil {
		return "", fmt.Errorf("failed to verify storage domain ID %s (%w)", storageDomainID, err)
	}
	return storageDomainID, nil
}

func setupTestClusterID(clusterID ClusterID, client Client) (id ClusterID, err error) {
	if clusterID == "" {
		clusterID, err = findTestClusterID(client)
		if err != nil {
			return "", fmt.Errorf("failed to find a cluster to test on (%w)", err)
		}
	} else if err := verifyTestClusterID(client, clusterID); err != nil {
		return "", fmt.Errorf("failed to verify cluster ID %s (%w)", clusterID, err)
	}
	return clusterID, nil
}

func createTestClient(
	url string,
	username string,
	password string,
	tlsProvider TLSProvider,
	mock bool,
	logger Logger,
) (Client, error) {
	var client Client
	var err error
	if mock {
		client = NewMockWithLogger(logger)
	} else {
		client, err = New(
			url,
			username,
			password,
			tlsProvider,
			logger,
			nil,
		)
		if err != nil {
			return nil, err
		}
	}
	return client, err
}

func findBlankTemplateID(client Client) (TemplateID, error) {
	template, err := client.GetBlankTemplate()
	if err != nil {
		return "", fmt.Errorf("failed to find blank template for testing (%w)", err)
	}
	return template.ID(), nil
}

func verifyBlankTemplateID(client Client, templateID TemplateID) error {
	_, err := client.GetTemplate(templateID)
	return err
}

func findTestClusterID(client Client) (ClusterID, error) {
	clusters, err := client.ListClusters()
	if err != nil {
		return "", err
	}
	hosts, err := client.ListHosts()
	if err != nil {
		return "", err
	}
	for _, cluster := range clusters {
		for _, host := range hosts {
			if host.Status() == HostStatusUp && host.ClusterID() == cluster.ID() {
				return cluster.ID(), nil
			}
		}
	}
	return "", fmt.Errorf("failed to find cluster suitable for testing")
}

func verifyTestClusterID(client Client, clusterID ClusterID) error {
	_, err := client.GetCluster(clusterID)
	return err
}

var errNoTestStorageDomainFound = fmt.Errorf("failed to find a working storage domain for testing")

func findTestStorageDomainID(skipID StorageDomainID, client Client) (StorageDomainID, error) {
	storageDomains, err := client.ListStorageDomains()
	if err != nil {
		return "", err
	}
	storageDomains = storageDomains.Filter(func(sd StorageDomain) bool {
		return sd.ID() != skipID
	}).Filter(func(sd StorageDomain) bool {
		// Assume 2GB will be enough for testing
		return sd.Available() >= 2*1024*1024*1024
	}).Filter(func(sd StorageDomain) bool {
		return (sd.Status() == StorageDomainStatusActive) || (sd.Status() == StorageDomainStatusNA && sd.ExternalStatus() == StorageDomainExternalStatusOk)
	}).Filter(func(sd StorageDomain) bool {
		return sd.StorageType() == StorageDomainTypeNFS
	})

	if len(storageDomains) > 0 {
		return storageDomains[0].ID(), nil
	}
	return "", errNoTestStorageDomainFound
}

func verifyTestStorageDomainID(client Client, storageDomainID StorageDomainID) error {
	storageDomain, err := client.GetStorageDomain(storageDomainID)
	if err != nil {
		return err
	}
	if storageDomain.Status() == StorageDomainStatusActive ||
		(storageDomain.Status() == StorageDomainStatusNA && storageDomain.ExternalStatus() == StorageDomainExternalStatusOk) {
		return nil
	}
	return fmt.Errorf("storage domain %s is %s, not active, external status is %s, not ok", storageDomain.ID(), storageDomain.Status(), storageDomain.ExternalStatus())
}

type testHelper struct {
	client                   Client
	tls                      TLSProvider
	rand                     *rand.Rand
	clusterID                ClusterID
	storageDomainID          StorageDomainID
	blankTemplateID          TemplateID
	vnicProfileID            VNICProfileID
	secondaryStorageDomainID StorageDomainID
	password                 string
	username                 string
}

var resourceNameRe = regexp.MustCompile(`[^a-zA-Z\d_.-]`)

func (t *testHelper) GenerateTestResourceName(test *testing.T) string {
	return resourceNameRe.ReplaceAllString(fmt.Sprintf("%s-%s", test.Name(), t.GenerateRandomID(5)), "_")
}

func (t *testHelper) GetUsername() string {
	return t.username
}

func (t *testHelper) GetPassword() string {
	return t.password
}

func (t *testHelper) GetSecondaryStorageDomainID(te *testing.T) StorageDomainID {
	if t.secondaryStorageDomainID == "" {
		te.Skipf("No secondary storage domain available, skipping test.")
	}
	return t.secondaryStorageDomainID
}

func (t *testHelper) GetTLS() TLSProvider {
	return t.tls
}

func (t *testHelper) GetVNICProfileID() VNICProfileID {
	return t.vnicProfileID
}

func (t *testHelper) GetClient() Client {
	return t.client
}

func (t *testHelper) GetClusterID() ClusterID {
	return t.clusterID
}

func (t *testHelper) GetBlankTemplateID() TemplateID {
	return t.blankTemplateID
}

func (t *testHelper) GetStorageDomainID() StorageDomainID {
	return t.storageDomainID
}

func (t *testHelper) GenerateRandomID(length uint) string {
	return generateRandomID(length, t.rand)
}

// NewMockTestHelper creates a test helper with a mocked oVirt engine in the backend.
func NewMockTestHelper(logger ovirtclientlog.Logger) (TestHelper, error) {
	helper, err := NewTestHelper(
		"https://localhost/ovirt-engine/api",
		"admin@internal",
		"",
		nil,
		TLS().Insecure(),
		true,
		logger,
	)
	if err != nil {
		return nil, err
	}
	return helper, nil
}

// NewLiveTestHelperFromEnv is a function that creates a test helper working against a live (not mock)
// oVirt engine using environment variables. The environment variables are as follows:
//
//	OVIRT_URL
//
// URL of the oVirt engine. Mandatory.
//
//	OVIRT_USERNAME
//
// The username for the oVirt engine. Mandatory.
//
//	OVIRT_PASSWORD
//
// The password for the oVirt engine. Mandatory.
//
//	OVIRT_CAFILE
//
// A file containing the CA certificate in PEM format.
//
//	OVIRT_CA_BUNDLE
//
// Provide the CA certificate in PEM format directly.
//
//	OVIRT_INSECURE
//
// Disable certificate verification if set. Not recommended.
//
//	OVIRT_CLUSTER_ID
//
// The cluster to use for testing. Will be automatically chosen if not provided.
//
//	OVIRT_BLANK_TEMPLATE_ID
//
// ID of the blank template. Will be automatically chosen if not provided.
//
//	OVIRT_STORAGE_DOMAIN_ID
//
// Storage domain to use for testing. Will be automatically chosen if not provided.
//
//	OVIRT_VNIC_PROFILE_ID
//
// VNIC profile to use for testing. Will be automatically chosen if not provided.
func NewLiveTestHelperFromEnv(logger ovirtclientlog.Logger) (TestHelper, error) {
	// Note: if this function changes please update the documentation above and also doc.go.
	url, tls, err := getConnectionParametersForLiveTesting()
	if err != nil {
		return nil, err
	}
	user := os.Getenv("OVIRT_USERNAME")
	if user == "" {
		return nil, fmt.Errorf("the OVIRT_USER environment variable must not be empty")
	}
	password := os.Getenv("OVIRT_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("the OVIRT_PASSWORD environment variable must not be empty")
	}

	params := TestHelperParams()
	params.WithClusterID(ClusterID(os.Getenv("OVIRT_CLUSTER_ID")))
	params.WithBlankTemplateID(TemplateID(os.Getenv("OVIRT_BLANK_TEMPLATE_ID")))
	params.WithStorageDomainID(StorageDomainID(os.Getenv("OVIRT_STORAGE_DOMAIN_ID")))
	params.WithSecondaryStorageDomainID(StorageDomainID(os.Getenv("OVIRT_SECONDARY_STORAGE_DOMAIN_ID")))
	params.WithVNICProfileID(VNICProfileID(os.Getenv("OVIRT_VNIC_PROFILE_ID")))

	helper, err := NewTestHelper(
		url,
		user,
		password,
		params,
		tls,
		false,
		logger,
	)
	if err != nil {
		return nil, err
	}
	return helper, nil
}

func getConnectionParametersForLiveTesting() (string, TLSProvider, error) {
	// Note: if this function changes please update the documentation above and also doc.go.
	url := os.Getenv("OVIRT_URL")
	if url == "" {
		return "", nil, fmt.Errorf("the OVIRT_URL environment variable must not be empty")
	}
	tls := TLS()
	configured := false
	if caFile := os.Getenv("OVIRT_CAFILE"); caFile != "" {
		configured = true
		tls.CACertsFromFile(caFile)
	}
	if caDir := os.Getenv("OVIRT_CA_DIR"); caDir != "" {
		configured = true
		tls.CACertsFromDir(caDir)
	}
	if caFile := os.Getenv("OVIRT_CA_FILE"); caFile != "" {
		configured = true
		tls.CACertsFromFile(caFile)
	}
	if caCert := os.Getenv("OVIRT_CA_CERT"); caCert != "" {
		configured = true
		tls.CACertsFromMemory([]byte(caCert))
	}
	if os.Getenv("OVIRT_INSECURE") != "" {
		configured = true
		tls.Insecure()
	}
	if system := os.Getenv("OVIRT_SYSTEM"); system != "" {
		configured = true
		tls.CACertsFromSystem()
	}
	if !configured {
		return "", nil, fmt.Errorf("one of OVIRT_CAFILE, OVIRT_CA_CERT, or OVIRT_INSECURE must be set")
	}
	return url, tls, nil
}
