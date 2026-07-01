package types

import (
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"
)

// BuildNoProxySet creates a set of NoProxy entries from install configuration.
// It includes localhost entries (.svc, .cluster.local, 127.0.0.1, localhost),
// all network CIDRs (cluster, service, machine), the internal API server hostname,
// and user-provided NoProxy values.
// The caller can add platform-specific entries as needed.
func BuildNoProxySet(config *InstallConfig) sets.Set[string] {
	set := sets.New[string](
		"127.0.0.1",
		"localhost",
		".svc",
		".cluster.local",
	)

	for _, network := range config.Networking.ServiceNetwork {
		set.Insert(network.String())
	}

	for _, network := range config.Networking.MachineNetwork {
		set.Insert(network.CIDR.String())
	}

	for _, network := range config.Networking.ClusterNetwork {
		set.Insert(network.CIDR.String())
	}

	// Add internal API server hostname
	set.Insert("api-int." + config.ClusterDomain())

	if config.Proxy != nil {
		for _, userValue := range strings.Split(config.Proxy.NoProxy, ",") {
			trimmed := strings.TrimSpace(userValue)
			if trimmed != "" {
				set.Insert(trimmed)
			}
		}
	}

	return set
}
