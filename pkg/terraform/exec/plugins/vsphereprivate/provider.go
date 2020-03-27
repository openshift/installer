package vsphereprivate

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere"
)

// defaultAPITimeout is a default timeout value that is passed to functions
// requiring contexts, and other various waiters.
const defaultAPITimeout = time.Minute * 5

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	vsphereProvider := vsphere.Provider()

	vsphereProvider.(*schema.Provider).DataSourcesMap = map[string]*schema.Resource{}

	vsphereProvider.(*schema.Provider).ResourcesMap = map[string]*schema.Resource{
		"vsphereprivate_import_ova": resourceVSpherePrivateImportOva(),
	}

	vsphereProvider.(*schema.Provider).ConfigureFunc = providerConfigure

	return vsphereProvider
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	c, err := NewConfig(d)
	if err != nil {
		return nil, err
	}
	return c.Client()
}
