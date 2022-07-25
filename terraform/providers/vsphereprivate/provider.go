package main

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere"
)

// defaultAPITimeout is a default timeout value that is passed to functions
// requiring contexts, and other various waiters.
const defaultAPITimeout = time.Minute * 5

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	// reusing the terraform-provider-vsphere Provider() function
	vsphereProvider := vsphere.Provider()

	// Since we are not using Data, empty that map
	vsphereProvider.DataSourcesMap = map[string]*schema.Resource{}

	//Add our resources
	vsphereProvider.ResourcesMap = map[string]*schema.Resource{
		"vsphereprivate_import_ova": resourceVSpherePrivateImportOva(),
	}

	vsphereProvider.ConfigureFunc = providerConfigure

	return vsphereProvider
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	c, err := NewConfig(d)
	if err != nil {
		return nil, err
	}
	return c.Client()
}
