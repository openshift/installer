package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/openshift/installer/pkg/terraform/exec/plugins/ibms3presign"
)

func init() {
	ibmS3PresignProvider := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: ibms3presign.Provider,
		})
	}
	KnownPlugins["terraform-provider-ibms3presign"] = ibmS3PresignProvider
}
