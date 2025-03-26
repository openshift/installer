// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cloudfoundry

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMServiceKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMServiceKeyRead,

		Schema: map[string]*schema.Schema{
			"credentials": {
				Description: "Credentials asociated with the key",
				Sensitive:   true,
				Type:        schema.TypeMap,
				Computed:    true,
			},
			"name": {
				Description: "The name of the service key",
				Type:        schema.TypeString,
				Required:    true,
			},
			"service_instance_name": {
				Description: "Service instance name for example, speech_to_text",
				Type:        schema.TypeString,
				Required:    true,
			},
			"space_guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The guid of the space in which the service instance is present",
			},
		},
		DeprecationMessage: "This service is deprecated.",
	}
}

func dataSourceIBMServiceKeyRead(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(conns.ClientSession).MccpAPI()
	if err != nil {
		return err
	}
	siAPI := cfClient.ServiceInstances()
	skAPI := cfClient.ServiceKeys()
	serviceInstanceName := d.Get("service_instance_name").(string)
	spaceGUID := d.Get("space_guid").(string)
	name := d.Get("name").(string)
	inst, err := siAPI.FindByNameInSpace(spaceGUID, serviceInstanceName)
	if err != nil {
		return err
	}
	serviceInstance, err := siAPI.Get(inst.GUID)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving service: %s", err)
	}
	serviceKey, err := skAPI.FindByName(serviceInstance.Metadata.GUID, name)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving service key: %s", err)
	}
	d.SetId(serviceKey.GUID)
	d.Set("credentials", flex.Flatten(serviceKey.Credentials))
	return nil
}
