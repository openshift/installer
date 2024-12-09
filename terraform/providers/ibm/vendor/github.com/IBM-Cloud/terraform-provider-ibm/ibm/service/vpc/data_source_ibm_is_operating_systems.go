// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isOperatingSystems = "operating_systems"
)

func DataSourceIBMISOperatingSystems() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISOperatingSystemsRead,

		Schema: map[string]*schema.Schema{
			isOperatingSystems: {
				Type:        schema.TypeList,
				Description: "List of operating systems",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isOperatingSystemAllowUserImageCreation: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Users may create new images with this operating system",
						},
						isOperatingSystemName: {
							Type:        schema.TypeString,
							Computed:    true,
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
						isOperatingSystemUserDataFormat: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user data format for this operating system",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMISOperatingSystemsRead(d *schema.ResourceData, meta interface{}) error {
	err := osList(d, meta)
	if err != nil {
		return err
	}
	return nil
}

func osList(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	start := ""
	allrecs := []vpcv1.OperatingSystem{}
	for {
		listOperatingSystemsOptions := &vpcv1.ListOperatingSystemsOptions{}
		if start != "" {
			listOperatingSystemsOptions.Start = &start
		}

		osList, response, err := sess.ListOperatingSystems(listOperatingSystemsOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Fetching operating systems %s\n%s", err, response)
		}
		start = flex.GetNext(osList.Next)
		allrecs = append(allrecs, osList.OperatingSystems...)
		if start == "" {
			break
		}
	}
	osInfo := make([]map[string]interface{}, 0)
	for _, os := range allrecs {
		l := map[string]interface{}{
			isOperatingSystemName:         *os.Name,
			isOperatingSystemArchitecture: *os.Architecture,
			isOperatingSystemDHOnly:       *os.DedicatedHostOnly,
			isOperatingSystemFamily:       *os.Family,
			isOperatingSystemHref:         *os.Href,
			isOperatingSystemDisplayName:  *os.DisplayName,
			isOperatingSystemVendor:       *os.Vendor,
			isOperatingSystemVersion:      *os.Version,
		}
		if os.AllowUserImageCreation != nil {
			l[isOperatingSystemAllowUserImageCreation] = *os.AllowUserImageCreation
		}
		if os.UserDataFormat != nil {
			l[isOperatingSystemUserDataFormat] = *os.UserDataFormat
		}
		osInfo = append(osInfo, l)
	}
	d.SetId(dataSourceIBMISOperatingSystemsId(d))
	d.Set(isOperatingSystems, osInfo)
	return nil
}

// dataSourceIBMISOperatingSystemsId returns a reasonable ID for a os list.
func dataSourceIBMISOperatingSystemsId(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
