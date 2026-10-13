package kubeconfig

import (
	"context"
	"fmt"
	"strings"

	"github.com/openshift/installer/pkg/asset"
	agentmanifests "github.com/openshift/installer/pkg/asset/agent/manifests"
	"github.com/openshift/installer/pkg/asset/tls"
)

// AgentAdminClient is the asset for the agent admin kubeconfig.
type AgentAdminClient struct {
	AdminClient
}

// Dependencies returns the dependency of the kubeconfig.
func (k *AgentAdminClient) Dependencies() []asset.Asset {
	return []asset.Asset{
		&tls.AdminKubeConfigClientCertKey{},
		&tls.KubeAPIServerCompleteCABundle{},
		&agentmanifests.ClusterDeployment{},
	}
}

// Generate generates the kubeconfig.
func (k *AgentAdminClient) Generate(_ context.Context, parents asset.Parents) error {
	ca := &tls.KubeAPIServerCompleteCABundle{}
	clientCertKey := &tls.AdminKubeConfigClientCertKey{}
	parents.Get(ca, clientCertKey)

	clusterDeployment := &agentmanifests.ClusterDeployment{}
	parents.Get(clusterDeployment)

	clusterName := clusterDeployment.Config.Spec.ClusterName
	extAPIServerURL := fmt.Sprintf("https://api.%s.%s:6443", clusterName, strings.TrimSuffix(clusterDeployment.Config.Spec.BaseDomain, "."))

	return k.kubeconfig.generate(
		ca,
		clientCertKey,
		extAPIServerURL,
		clusterName,
		"admin",
		kubeconfigAdminPath,
	)
}
