package types

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	configv1 "github.com/openshift/api/config/v1"
	features "github.com/openshift/api/features"
)

// StringsToIPs is used to convert list of strings to list of IP addresses.
func StringsToIPs(ips []string) []configv1.IP {
	res := []configv1.IP{}

	if ips == nil {
		return res
	}

	for _, ip := range ips {
		res = append(res, configv1.IP(ip))
	}

	return res
}

// MachineNetworksToCIDRs is used to convert list of Machine Network Entries to
// list of CIDRs.
func MachineNetworksToCIDRs(nets []MachineNetworkEntry) []configv1.CIDR {
	res := []configv1.CIDR{}

	if nets == nil {
		return res
	}

	for _, net := range nets {
		res = append(res, configv1.CIDR(net.CIDR.String()))
	}

	return res
}

// GetClusterProfileName utility method to retrieve the cluster profile setting.  This is used
// when dealing with openshift api to get FeatureSets.
func GetClusterProfileName() features.ClusterProfileName {
	// Get cluster profile for new FeatureGate access.  Blank is no longer an option, so default to
	// SelfManaged.
	clusterProfile := features.SelfManaged
	if cp := os.Getenv("OPENSHIFT_INSTALL_EXPERIMENTAL_CLUSTER_PROFILE"); cp != "" {
		logrus.Warnf("Found override for Cluster Profile: %q", cp)
		// All profiles when getting FeatureSets need to have "include.release.openshift.io/" at the beginning.
		// See vendor/openshift/api/config/v1/feature_gates.go for more info.
		clusterProfile = features.ClusterProfileName(fmt.Sprintf("%s%s", "include.release.openshift.io/", cp))
	}
	return clusterProfile
}
