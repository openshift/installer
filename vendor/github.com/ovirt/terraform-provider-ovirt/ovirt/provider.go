// Copyright (C) 2017 Battelle Memorial Institute
// Copyright (C) 2018 Joey Ma <majunjiev@gmail.com>
// All rights reserved.
//
// This software may be modified and distributed under the terms
// of the BSD-2 license.  See the LICENSE file for details.

package ovirt

import (
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

// Provider returns oVirt provider configuration
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVIRT_USERNAME", os.Getenv("OVIRT_USERNAME")),
				Description: "Login username",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVIRT_PASSWORD", os.Getenv("OVIRT_PASSWORD")),
				Description: "Login password",
				Sensitive:   true,
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVIRT_URL", os.Getenv("OVIRT_URL")),
				Description: "Ovirt server url",
			},
			"headers": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Additional headers to be added to each API call",
			},
		},
		ConfigureFunc: ConfigureProvider,
		ResourcesMap: map[string]*schema.Resource{
			"ovirt_vm":              resourceOvirtVM(),
			"ovirt_template":        resourceOvirtTemplate(),
			"ovirt_disk":            resourceOvirtDisk(),
			"ovirt_disk_attachment": resourceOvirtDiskAttachment(),
			"ovirt_datacenter":      resourceOvirtDataCenter(),
			"ovirt_network":         resourceOvirtNetwork(),
			"ovirt_vnic":            resourceOvirtVnic(),
			"ovirt_vnic_profile":    resourceOvirtVnicProfile(),
			"ovirt_snapshot":        resourceOvirtSnapshot(),
			"ovirt_storage_domain":  resourceOvirtStorageDomain(),
			"ovirt_tag":             resourceOvirtTag(),
			"ovirt_user":            resourceOvirtUser(),
			"ovirt_cluster":         resourceOvirtCluster(),
			"ovirt_mac_pool":        resourceOvirtMacPool(),
			"ovirt_host":            resourceOvirtHost(),
			"ovirt_image_transfer":  resourceOvirtImageTransfer(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ovirt_disks":          dataSourceOvirtDisks(),
			"ovirt_datacenters":    dataSourceOvirtDataCenters(),
			"ovirt_networks":       dataSourceOvirtNetworks(),
			"ovirt_clusters":       dataSourceOvirtClusters(),
			"ovirt_storagedomains": dataSourceOvirtStorageDomains(),
			"ovirt_vnic_profiles":  dataSourceOvirtVNicProfiles(),
			"ovirt_authzs":         dataSourceOvirtAuthzs(),
			"ovirt_users":          dataSourceOvirtUsers(),
			"ovirt_mac_pools":      dataSourceOvirtMacPools(),
			"ovirt_vms":            dataSourceOvirtVMs(),
			"ovirt_hosts":          dataSourceOvirtHosts(),
			"ovirt_nics":           dataSourceOvirtNics(),
			"ovirt_templates":      dataSourceOvirtTemplates(),
		},
	}
}

// ConfigureProvider initializes the API connection object by config items
func ConfigureProvider(d *schema.ResourceData) (interface{}, error) {
	connBuilder := ovirtsdk4.NewConnectionBuilder().
		URL(d.Get("url").(string)).
		Username(d.Get("username").(string)).
		Password(d.Get("password").(string)).
		Insecure(true)

	// Set headers if needed
	if v, ok := d.GetOk("headers"); ok {
		headers := map[string]string{}
		for k, v := range v.(map[string]interface{}) {
			headers[k] = v.(string)
		}
		connBuilder.Headers(headers)
	}

	return connBuilder.Build()
}
