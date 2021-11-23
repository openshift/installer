package plugins

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-provider-openstack/terraform-provider-openstack/openstack"
)

func init() {
	// Fix for Bug 2015837: terraform-provider-openstack obeys the
	// environment variable OS_CLOUD with precedence.
	if err := os.Unsetenv("OS_CLOUD"); err != nil {
		panic(fmt.Errorf("failed to unset OS_CLOUD: %w", err))
	}
	exec := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: openstack.Provider,
		})
	}
	KnownPlugins["terraform-provider-openstack"] = exec
}
