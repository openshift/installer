// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/sl"
)

func resourceIBMComputeProvisioningHook() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMComputeProvisioningHookCreate,
		Read:     resourceIBMComputeProvisioningHookRead,
		Update:   resourceIBMComputeProvisioningHookUpdate,
		Delete:   resourceIBMComputeProvisioningHookDelete,
		Exists:   resourceIBMComputeProvisioningHookExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Provision hook name",
			},

			"uri": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "URI of the hook",
			},

			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "Tags associated with resource",
			},
		},
	}
}

func resourceIBMComputeProvisioningHookCreate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetProvisioningHookService(sess.SetRetries(0))

	opts := datatypes.Provisioning_Hook{
		Name: sl.String(d.Get("name").(string)),
		Uri:  sl.String(d.Get("uri").(string)),
	}

	hook, err := service.CreateObject(&opts)
	if err != nil {
		return fmt.Errorf("Error creating Provisioning Hook: %s", err)
	}

	d.SetId(strconv.Itoa(*hook.Id))
	log.Printf("[INFO] Provisioning Hook ID: %d", *hook.Id)

	return resourceIBMComputeProvisioningHookRead(d, meta)
}

func resourceIBMComputeProvisioningHookRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetProvisioningHookService(sess)

	hookId, _ := strconv.Atoi(d.Id())

	hook, err := service.Id(hookId).GetObject()
	if err != nil {
		if err, ok := err.(sl.Error); ok {
			if err.StatusCode == http.StatusNotFound {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("Error retrieving Provisioning Hook: %s", err)
	}

	d.Set("name", hook.Name)
	d.Set("uri", hook.Uri)

	return nil
}

func resourceIBMComputeProvisioningHookUpdate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetProvisioningHookService(sess.SetRetries(0))

	hookId, _ := strconv.Atoi(d.Id())

	opts := datatypes.Provisioning_Hook{}

	if d.HasChange("name") {
		opts.Name = sl.String(d.Get("name").(string))
	}

	if d.HasChange("uri") {
		opts.Uri = sl.String(d.Get("uri").(string))
	}

	opts.TypeId = sl.Int(1)
	_, err := service.Id(hookId).EditObject(&opts)

	if err != nil {
		return fmt.Errorf("Error editing Provisioning Hook: %s", err)
	}
	return nil
}

func resourceIBMComputeProvisioningHookDelete(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetProvisioningHookService(sess)

	hookId, err := strconv.Atoi(d.Id())
	log.Printf("[INFO] Deleting Provisioning Hook: %d", hookId)
	_, err = service.Id(hookId).DeleteObject()
	if err != nil {
		return fmt.Errorf("Error deleting Provisioning Hook: %s", err)
	}

	return nil
}

func resourceIBMComputeProvisioningHookExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetProvisioningHookService(sess)

	hookId, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	result, err := service.Id(hookId).GetObject()
	if err != nil {
		if apiErr, ok := err.(sl.Error); ok {
			if apiErr.StatusCode == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}
	return result.Id != nil && *result.Id == hookId, nil
}
