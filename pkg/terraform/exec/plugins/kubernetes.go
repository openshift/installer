package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-provider-kubernetes/kubernetes"
)

func init() {
	exec := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: kubernetes.Provider})
	}
	KnownPlugins["terraform-provider-kubernetes"] = &TFPlugin{
		Name:    "kubernetes",
		Exec:    exec,
		Version: GetKubernetesVersion(),
	}
}
