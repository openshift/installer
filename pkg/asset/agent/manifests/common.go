package manifests

import (
	"fmt"
	"net"
	"strings"

	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/version"
)

func getAgentClusterInstallName(ic *agent.OptionalInstallConfig) string {
	return ic.ClusterName()
}

func getClusterDeploymentName(ic *agent.OptionalInstallConfig) string {
	return ic.ClusterName()
}

func getPullSecretName(clusterName string) string {
	return clusterName + "-pull-secret"
}

func getProxy(proxy *types.Proxy, machineNetwork *[]types.MachineNetworkEntry, rendezvousIP string) *aiv1beta1.Proxy {
	noProxy := proxy.NoProxy
	if (proxy.HTTPProxy != "" || proxy.HTTPSProxy != "") && rendezvousIP != "" {
		// if proxy set, add the machineNetwork corresponding to rendezvousIP to noproxy
		cidr := ""
		if machineNetwork != nil {
			for _, mn := range *machineNetwork {
				ipNet, err := ipnet.ParseCIDR(mn.CIDR.String())
				if err != nil {
					continue
				}
				ip := net.ParseIP(rendezvousIP)
				if ipNet.Contains(ip) {
					cidr = mn.CIDR.String()
					break
				}
			}
		}

		if cidr != "" {
			if noProxy == "" {
				noProxy = cidr
			} else if !strings.Contains(noProxy, cidr) {
				noProxy = fmt.Sprintf("%s,%s", noProxy, cidr)
			}
		}
	}
	return &aiv1beta1.Proxy{
		HTTPProxy:  proxy.HTTPProxy,
		HTTPSProxy: proxy.HTTPSProxy,
		NoProxy:    noProxy,
	}
}

func getNMStateConfigLabels(clusterName string) map[string]string {
	return map[string]string{
		"infraenvs.agent-install.openshift.io": clusterName,
	}
}

func getClusterImageSetReferenceName() string {
	versionString, _ := version.Version()
	return "openshift-" + versionString
}
