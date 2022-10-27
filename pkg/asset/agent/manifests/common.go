package manifests

import (
	"github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/version"

	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
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

func getPullSecretName(ic *agent.OptionalInstallConfig) string {
	return ic.ClusterName() + "-pull-secret"
}

func getProxy(ic *agent.OptionalInstallConfig) *aiv1beta1.Proxy {
	return &aiv1beta1.Proxy{
		HTTPProxy:  ic.Config.Proxy.HTTPProxy,
		HTTPSProxy: ic.Config.Proxy.HTTPSProxy,
		NoProxy:    ic.Config.Proxy.NoProxy,
	}
}

func getObjectMetaNamespace(ic *agent.OptionalInstallConfig) string {
	if ic.Config != nil {
		return ic.Config.Namespace
	}
	return ""
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
