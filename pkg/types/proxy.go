package types

import (
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"
)

// BuildNoProxySet creates a set of NoProxy entries from networking configuration.
// It includes localhost entries (.svc, .cluster.local, 127.0.0.1, localhost),
// all network CIDRs (cluster, service, machine), and user-provided NoProxy values.
// The caller can add platform-specific and API-server entries as needed.
func BuildNoProxySet(networking *Networking, userNoProxy string) sets.String {
	set := sets.NewString(
		"127.0.0.1",
		"localhost",
		".svc",
		".cluster.local",
	)

	for _, network := range networking.ServiceNetwork {
		set.Insert(network.String())
	}

	for _, network := range networking.MachineNetwork {
		set.Insert(network.CIDR.String())
	}

	for _, network := range networking.ClusterNetwork {
		set.Insert(network.CIDR.String())
	}

	for _, userValue := range strings.Split(userNoProxy, ",") {
		trimmed := strings.TrimSpace(userValue)
		if trimmed != "" {
			set.Insert(trimmed)
		}
	}

	return set
}
