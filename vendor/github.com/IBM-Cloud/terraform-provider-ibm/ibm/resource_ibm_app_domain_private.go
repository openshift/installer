// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	v2 "github.com/IBM-Cloud/bluemix-go/api/mccp/mccpv2"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIBMAppDomainPrivate() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMAppDomainPrivateCreate,
		Read:     resourceIBMAppDomainPrivateRead,
		Update:   resourceIBMAppDomainPrivateUpdate,
		Delete:   resourceIBMAppDomainPrivateDelete,
		Exists:   resourceIBMAppDomainPrivateExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The name of the domain",
				ValidateFunc: validateDomainName,
			},

			"org_guid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The organization that owns the domain.",
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

func resourceIBMAppDomainPrivateCreate(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(ClientSession).MccpAPI()
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	orgGUID := d.Get("org_guid").(string)

	params := v2.PrivateDomainRequest{
		Name:    name,
		OrgGUID: orgGUID,
	}

	prdomain, err := cfClient.PrivateDomains().Create(params)
	if err != nil {
		return fmt.Errorf("Error creating private domain: %s", err)
	}

	d.SetId(prdomain.Metadata.GUID)

	return resourceIBMAppDomainPrivateRead(d, meta)
}

func resourceIBMAppDomainPrivateUpdate(d *schema.ResourceData, meta interface{}) error {
	//Only tags are updated and that too locally hence nothing to validate and update in terms of real API at this point
	return nil
}

func resourceIBMAppDomainPrivateRead(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(ClientSession).MccpAPI()
	if err != nil {
		return err
	}
	prdomainGUID := d.Id()

	prdomain, err := cfClient.PrivateDomains().Get(prdomainGUID)
	if err != nil {
		return fmt.Errorf("Error retrieving private domain: %s", err)
	}
	d.Set("name", prdomain.Entity.Name)
	d.Set("org_guid", prdomain.Entity.OwningOrganizationGUID)

	return nil
}

func resourceIBMAppDomainPrivateDelete(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(ClientSession).MccpAPI()
	if err != nil {
		return err
	}

	prdomainGUID := d.Id()

	err = cfClient.PrivateDomains().Delete(prdomainGUID, false)
	if err != nil {
		return fmt.Errorf("Error deleting private domain: %s", err)
	}

	d.SetId("")

	return nil
}

func resourceIBMAppDomainPrivateExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	cfClient, err := meta.(ClientSession).MccpAPI()
	if err != nil {
		return false, err
	}
	prdomainGUID := d.Id()

	prdomain, err := cfClient.PrivateDomains().Get(prdomainGUID)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}

	return prdomain.Metadata.GUID == prdomainGUID, nil
}
