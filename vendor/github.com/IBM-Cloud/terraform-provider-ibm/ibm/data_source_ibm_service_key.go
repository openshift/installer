// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMServiceKey() *schema.Resource {
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
	}
}

func dataSourceIBMServiceKeyRead(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(ClientSession).MccpAPI()
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
		return fmt.Errorf("Error retrieving service: %s", err)
	}
	serviceKey, err := skAPI.FindByName(serviceInstance.Metadata.GUID, name)
	if err != nil {
		return fmt.Errorf("Error retrieving service key: %s", err)
	}
	d.SetId(serviceKey.GUID)
	d.Set("credentials", Flatten(serviceKey.Credentials))
	return nil
}
