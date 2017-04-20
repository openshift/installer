package server

import (
	"errors"
	"net"
	"strings"

	"golang.org/x/net/context"

	"github.com/coreos/tectonic-installer/installer/server/asset"
)

// Cluster errors
var (
	errInvalidClusterType      = errors.New("installer: Invalid cluster kind")
	errInvalidExternalETCD     = errors.New("installer: External etcd must be a valid host:port combination")
	errMissingControllerDomain = errors.New("installer: Missing master domain name")
	errTooFewControllers       = errors.New("installer: At least one master is required")
	errClusterTooSmall         = errors.New("installer: Cluster must have at least one master and worker")
	errMissingMACAddress       = errors.New("installer: Missing MAC address")
	errMissingNodeName         = errors.New("installer: Missing node name")
	errMissingMatchboxEndpoint = errors.New("installer: Missing matchbox endpoint")
	errMissingChannel          = errors.New("installer: Missing CoreOS channel")
	errMissingVersion          = errors.New("installer: Missing CoreOS version")
)

// A Cluster defines cluster setup operations and steps.
type Cluster interface {
	// Initialize validates cluster fields and sets any defaults.
	Initialize() error
	// GenerateAssets generates cluster provisioning assets.
	GenerateAssets() ([]asset.Asset, error)
	// StatusChecker returns a checker for the status of cluster components.
	StatusChecker() (StatusChecker, error)
	// Kind returns the kind name of a cluster.
	Kind() string
	// Publish writes configs to a provisioning service.
	Publish(context.Context) error
}

// The Node type can simplify generation of cluster manifests.
type Node struct {
	// FQDN
	Name string `json:"name"`
	// MAC Address
	MAC *macAddr `json:"mac"`
}

// macAddr is a net.HardwareAddr which can be JSON marshalled/unmarshalled.
type macAddr net.HardwareAddr

// ParseMACAddr parses s into a macAddr by calling through to net.ParseMAC.
func parseMACAddr(s string) (macAddr, error) {
	addr, err := net.ParseMAC(s)
	if err != nil {
		return nil, err
	}
	return macAddr(addr), nil
}

// String returns the ':' separated MAC address.
func (m *macAddr) String() string {
	return net.HardwareAddr(*m).String()
}

// DashString returns the '-' separated MAC address.
func (m *macAddr) DashString() string {
	return strings.Replace(m.String(), ":", "-", -1)
}

func (m *macAddr) MarshalJSON() ([]byte, error) {
	hwAddr := net.HardwareAddr(*m)
	return []byte(`"` + hwAddr.String() + `"`), nil
}

func (m *macAddr) UnmarshalJSON(data []byte) error {
	raw := strings.Replace(string(data), "\"", "", -1)
	addr, err := net.ParseMAC(raw)
	if err != nil {
		return err
	}
	*m = macAddr(addr)
	return nil
}

// Bool is a helper that allocates a new bool, stores v, and returns the
// pointer.
func Bool(v bool) *bool { return &v }

// Int is a helper that allocates a new int, stores v, and returns the
// pointer.
func Int(v int) *int { return &v }
