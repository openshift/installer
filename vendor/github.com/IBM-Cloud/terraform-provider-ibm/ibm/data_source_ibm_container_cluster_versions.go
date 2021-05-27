// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMContainerClusterVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMContainerClusterVersionsRead,

		Schema: map[string]*schema.Schema{
			"org_guid": {
				Description: "The bluemix organization guid this cluster belongs to",
				Type:        schema.TypeString,
				Optional:    true,
				Deprecated:  "This field is deprecated",
			},
			"space_guid": {
				Description: "The bluemix space guid this cluster belongs to",
				Type:        schema.TypeString,
				Optional:    true,
				Deprecated:  "This field is deprecated",
			},
			"account_guid": {
				Description: "The bluemix account guid this cluster belongs to",
				Type:        schema.TypeString,
				Optional:    true,
				Deprecated:  "This field is deprecated",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The cluster region",
				Deprecated:  "This field is deprecated",
			},
			"resource_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the resource group.",
			},
			"valid_kube_versions": {
				Description: "List supported kube-versions",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"valid_openshift_versions": {
				Description: "List of supported openshift-versions",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceIBMContainerClusterVersionsRead(d *schema.ResourceData, meta interface{}) error {
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return err
	}
	verAPI := csClient.KubeVersions()
	targetEnv, err := getClusterTargetHeader(d, meta)
	if err != nil {
		return err
	}

	availableVersions, _ := verAPI.ListV1(targetEnv)
	versions := make([]string, len(availableVersions["kubernetes"]))
	for i, version := range availableVersions["kubernetes"] {
		versions[i] = fmt.Sprintf("%d%s%d%s%d", version.Major, ".", version.Minor, ".", version.Patch)
	}

	openshiftVersions := make([]string, len(availableVersions["openshift"]))
	for i, version := range availableVersions["openshift"] {
		openshiftVersions[i] = fmt.Sprintf("%d%s%d%s%d", version.Major, ".", version.Minor, ".", version.Patch)
	}
	d.SetId(time.Now().UTC().String())
	d.Set("valid_kube_versions", versions)
	d.Set("valid_openshift_versions", openshiftVersions)
	return nil
}
