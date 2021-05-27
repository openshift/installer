// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"strconv"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIBMServiceKey() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMServiceKeyCreate,
		Read:     resourceIBMServiceKeyRead,
		Update:   resourceIBMServiceKeyUpdate,
		Delete:   resourceIBMServiceKeyDelete,
		Exists:   resourceIBMServiceKeyExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the service key",
			},

			"service_instance_guid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The guid of the service instance for which to create service key",
			},
			"parameters": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Arbitrary parameters to pass along to the service broker. Must be a JSON object",
			},
			"credentials": {
				Description: "Credentials asociated with the key",
				Type:        schema.TypeMap,
				Sensitive:   true,
				Computed:    true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceIBMServiceKeyCreate(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(ClientSession).MccpAPI()
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	serviceInstanceGUID := d.Get("service_instance_guid").(string)

	var keyParams map[string]interface{}

	if parameters, ok := d.GetOk("parameters"); ok {
		temp := parameters.(map[string]interface{})
		for k, v := range temp {
			if v == "true" || v == "false" {
				b, _ := strconv.ParseBool(v.(string))
				keyParams[k] = b

			} else {
				keyParams[k] = v
			}
		}
	}

	serviceKey, err := cfClient.ServiceKeys().Create(serviceInstanceGUID, name, keyParams)
	if err != nil {
		return fmt.Errorf("Error creating service key: %s", err)
	}

	d.SetId(serviceKey.Metadata.GUID)

	return resourceIBMServiceKeyRead(d, meta)
}

func resourceIBMServiceKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	//Only tags are updated and that too locally hence nothing to validate and update in terms of real API at this point
	return nil
}

func resourceIBMServiceKeyRead(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(ClientSession).MccpAPI()
	if err != nil {
		return err
	}
	serviceKeyGUID := d.Id()

	serviceKey, err := cfClient.ServiceKeys().Get(serviceKeyGUID)
	if err != nil {
		return fmt.Errorf("Error retrieving service key: %s", err)
	}
	d.Set("credentials", Flatten(serviceKey.Entity.Credentials))
	d.Set("service_instance_guid", serviceKey.Entity.ServiceInstanceGUID)
	d.Set("name", serviceKey.Entity.Name)

	return nil
}

func resourceIBMServiceKeyDelete(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(ClientSession).MccpAPI()
	if err != nil {
		return err
	}

	serviceKeyGUID := d.Id()

	err = cfClient.ServiceKeys().Delete(serviceKeyGUID)
	if err != nil {
		return fmt.Errorf("Error deleting service key: %s", err)
	}

	d.SetId("")

	return nil
}

func resourceIBMServiceKeyExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	cfClient, err := meta.(ClientSession).MccpAPI()
	if err != nil {
		return false, err
	}
	serviceKeyGUID := d.Id()

	serviceKey, err := cfClient.ServiceKeys().Get(serviceKeyGUID)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}

	return serviceKey.Metadata.GUID == serviceKeyGUID, nil
}
