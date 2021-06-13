package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/kubevirt/terraform-provider-kubevirt/kubevirt"

	kubevirttypes "github.com/openshift/installer/pkg/types/kubevirt"
)

func init() {
	kubevirtProvider := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: kubevirt.Provider})
	}
	KnownPlugins["terraform-provider-kubevirt"] = &TFPlugin{
		Name: kubevirttypes.Name,
		Exec: kubevirtProvider,
	}
}
