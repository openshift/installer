package manifests

import (
	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/version"
)

func getAgentClusterInstallName(ic *agent.OptionalInstallConfig) string {
	return ic.ClusterName()
}

func getClusterDeploymentName(ic *agent.OptionalInstallConfig) string {
	return ic.ClusterName()
}

func getInfraEnvName(ic *agent.OptionalInstallConfig) string {
	return ic.ClusterName()
}

func getPullSecretName(clusterName string) string {
	return clusterName + "-pull-secret"
}

func getProxy(ic *agent.OptionalInstallConfig) *aiv1beta1.Proxy {
	return &aiv1beta1.Proxy{
		HTTPProxy:  ic.Config.Proxy.HTTPProxy,
		HTTPSProxy: ic.Config.Proxy.HTTPSProxy,
		NoProxy:    ic.Config.Proxy.NoProxy,
	}
}

func getNMStateConfigName(ic *agent.OptionalInstallConfig) string {
	return ic.ClusterName()
}

func getNMStateConfigLabels(ic *agent.OptionalInstallConfig) map[string]string {
	return map[string]string{
		"infraenvs.agent-install.openshift.io": getInfraEnvName(ic),
	}
}

func getClusterImageSetReferenceName() string {
	versionString, _ := version.Version()
	return "openshift-" + versionString
}
