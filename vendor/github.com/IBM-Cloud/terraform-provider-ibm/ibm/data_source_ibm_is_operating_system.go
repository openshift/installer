// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isOperatingSystemName         = "name"
	isOperatingSystemArchitecture = "architecture"
	isOperatingSystemDHOnly       = "dedicated_host_only"
	isOperatingSystemDisplayName  = "display_name"
	isOperatingSystemFamily       = "family"
	isOperatingSystemHref         = "href"
	isOperatingSystemVendor       = "vendor"
	isOperatingSystemVersion      = "version"
)

func dataSourceIBMISOperatingSystem() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISOperatingSystemRead,

		Schema: map[string]*schema.Schema{

			isOperatingSystemName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The globally unique name for this operating system",
			},

			isOperatingSystemArchitecture: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The operating system architecture",
			},

			isOperatingSystemVersion: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The major release version of this operating system",
			},
			isOperatingSystemDHOnly: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag which shows images with this operating system can only be used on dedicated hosts or dedicated host groups",
			},
			isOperatingSystemDisplayName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A unique, display-friendly name for the operating system",
			},
			isOperatingSystemFamily: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the software family this operating system belongs to",
			},
			isOperatingSystemHref: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this operating system",
			},

			isOperatingSystemVendor: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The vendor of the operating system",
			},
		},
	}
}

func dataSourceIBMISOperatingSystemRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get(isOperatingSystemName).(string)
	err := osGet(d, meta, name)
	if err != nil {
		return err
	}
	return nil
}

func osGet(d *schema.ResourceData, meta interface{}, name string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	getOperatingSystemOptions := &vpcv1.GetOperatingSystemOptions{
		Name: &name,
	}
	os, response, err := sess.GetOperatingSystem(getOperatingSystemOptions)
	if err != nil || os == nil {
		return fmt.Errorf("Error Getting Operating System Details %s , %s", err, response)
	}
	d.Set(isOperatingSystemName, *os.Name)
	d.SetId(*os.Name)
	d.Set(isOperatingSystemDHOnly, *os.DedicatedHostOnly)
	d.Set(isOperatingSystemArchitecture, *os.Architecture)
	d.Set(isOperatingSystemDisplayName, *os.DisplayName)
	d.Set(isOperatingSystemFamily, *os.Family)
	d.Set(isOperatingSystemHref, *os.Href)
	d.Set(isOperatingSystemVendor, *os.Vendor)
	d.Set(isOperatingSystemVersion, *os.Version)
	return nil
}
