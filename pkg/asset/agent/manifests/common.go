package manifests

import (
	"github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/asset/agent/agentconfig"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/version"
)

func getAgentClusterInstallName(ic *agent.OptionalInstallConfig) string {
	return ic.Config.ObjectMeta.Name
}

func getClusterDeploymentName(ic *agent.OptionalInstallConfig) string {
	return ic.Config.ObjectMeta.Name
}

func getInfraEnvName(ic *agent.OptionalInstallConfig) string {
	return ic.Config.ObjectMeta.Name
}

func getPullSecretName(ic *agent.OptionalInstallConfig) string {
	return ic.Config.ObjectMeta.Name + "-pull-secret"
}

func getObjectMetaNamespace(ic *agent.OptionalInstallConfig) string {
	return ic.Config.Namespace
}

func getNMStateConfigName(a *agentconfig.AgentConfig) string {
	return a.Config.ObjectMeta.Name
}

func getNMStateConfigNamespace(a *agentconfig.AgentConfig) string {
	return a.Config.Namespace
}

func getNMStateConfigLabelsFromOptionalInstallConfig(ic *agent.OptionalInstallConfig) map[string]string {
	return map[string]string{
		"infraenvs.agent-install.openshift.io": getInfraEnvName(ic),
	}
}

func getNMStateConfigLabelsFromAgentConfig(a *agentconfig.AgentConfig) map[string]string {
	return map[string]string{
		"infraenvs.agent-install.openshift.io": getNMStateConfigName(a),
	}
}

func getClusterImageSetReferenceName() string {
	versionString, _ := version.Version()
	return "openshift-" + versionString
}

// getVIPs returns a string representation of the platform's API VIP and Ingress VIP.
// It returns an empty string if the platform does not configure a VIP
func getVIPs(p *types.Platform) (string, string) {
	switch {
	case p == nil:
		return "", ""
	case p.BareMetal != nil:
		return p.BareMetal.APIVIPs[0], p.BareMetal.IngressVIPs[0]
	case p.VSphere != nil:
		return p.VSphere.APIVIPs[0], p.VSphere.IngressVIPs[0]
	default:
		return "", ""
	}
}
