package types

import (
	"fmt"
	"net/url"
)

// Pre-defined environment keys. Environment settings can be added w/o being defined
// in this section.
const (
	CategoriesKey = "categories"
)

// ErrNotFound is returned by Get() for missing keys
var ErrNotFound = fmt.Errorf("environment key not found")

// ApiCredentials is set of identifiers and secrets used to authenticate with
// the underlying infrastructure.
type ApiCredentials struct {
	// Username for basic authentication
	Username string `json:"username,omitempty"`
	// Password for basic authentication
	Password string `json:"password,omitempty"`
	// KeyPair is JSON-encoded key pair for TLS client authentication
	KeyPair string `json:"keyPair,omitempty"`
}

// ManagementEndpoint specifies API endpoint used for interacting with underlying
// infrastructure.
type ManagementEndpoint struct {
	// ApiCredentials embedded into endpoint
	ApiCredentials
	// Address is URL of management endpoint
	Address *url.URL `json:"address,omitempty"`
	// Whether to authenticate TLS endpoint in case of HTTPS as transport.
	// HTTPS is used for encryption independent of this setting. An
	// unauthenticated TLS endpoint is prone to man-in-the-middle attacks.
	Insecure bool `json:"insecure,omitempty"`
	// AdditionalTrustBundle is a PEM-encoded certificate bundle to be used
	// in addition to system trust store
	AdditionalTrustBundle string `json:"additionalTrustBundle,omitempty"`
}

// Topology is a map of topological domains to topological segments.
// A topological domain is a sub-division of a cluster, like "region",
// "zone", "rack", etc
type Topology map[string]string

// Environment in which a Kubernetes has been deployed.
type Environment interface {
	// GetManagementEndpoint retrieves management endpoint
	GetManagementEndpoint(topology Topology) (*ManagementEndpoint, error)
	// Get retrieves settings applicable to the environment like project
	// or category to which resources created on behalf of Kubernetes cluster are
	// assigned to.
	// These settings might have to be explicitly propagated to resources.
	// Return whether lookup was successful to distinguish from nil as value.
	Get(topology Topology, key string) (interface{}, error)
}

// Provider of an environment
type Provider interface {
	Environment
}
