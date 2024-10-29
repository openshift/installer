package manifests

import (
	"fmt"

	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/installer/pkg/asset/agent"
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

func getProxy(proxy *types.Proxy, rendezvousIP string) *aiv1beta1.Proxy {
	// if proxy set add the rendezvousIP to noproxy
	noProxy := proxy.NoProxy
	if (proxy.HTTPProxy != "" || proxy.HTTPSProxy != "") && rendezvousIP != "" {
		if noProxy == "" {
			noProxy = rendezvousIP
		} else {
			noProxy = fmt.Sprintf("%s,%s", noProxy, rendezvousIP)
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
