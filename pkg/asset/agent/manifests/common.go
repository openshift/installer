package manifests

import (
	"github.com/openshift/installer/pkg/asset/agent"
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

func getNMStateConfigLabels(ic *agent.OptionalInstallConfig) map[string]string {
	return map[string]string{
		"infraenvs.agent-install.openshift.io": getInfraEnvName(ic),
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
		return p.BareMetal.APIVIP, p.BareMetal.IngressVIP
	case p.VSphere != nil:
		return p.VSphere.APIVIP, p.VSphere.IngressVIP
	default:
		return "", ""
	}
}
