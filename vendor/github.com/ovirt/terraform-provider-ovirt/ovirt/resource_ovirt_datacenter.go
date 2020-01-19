// Copyright (C) 2018 Joey Ma <majunjiev@gmail.com>
// All rights reserved.
//
// This software may be modified and distributed under the terms
// of the BSD-2 license.  See the LICENSE file for details.

package ovirt

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

func resourceOvirtDataCenter() *schema.Resource {
	return &schema.Resource{
		Create: resourceOvirtDataCenterCreate,
		Read:   resourceOvirtDataCenterRead,
		Update: resourceOvirtDataCenterUpdate,
		Delete: resourceOvirtDataCenterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			// This field identifies whether the datacenter uses local storage
			"local": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: false,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOvirtDataCenterCreate(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*ovirtsdk4.Connection)
	name := d.Get("name").(string)
	local := d.Get("local").(bool)

	//Name and Local are required when create a datacenter
	datacenterbuilder := ovirtsdk4.NewDataCenterBuilder().Name(name).Local(local)

	// Check if has description
	if description, ok := d.GetOk("description"); ok {
		datacenterbuilder = datacenterbuilder.Description(description.(string))
	}

	datacenter, err := datacenterbuilder.Build()
	if err != nil {
		return err
	}

	addResp, err := conn.SystemService().DataCentersService().Add().DataCenter(datacenter).Send()
	if err != nil {
		return err
	}

	d.SetId(addResp.MustDataCenter().MustId())
	return resourceOvirtDataCenterRead(d, meta)

}

func resourceOvirtDataCenterUpdate(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*ovirtsdk4.Connection)
	datacenterService := conn.SystemService().DataCentersService().DataCenterService(d.Id())
	datacenterBuilder := ovirtsdk4.NewDataCenterBuilder()

	if d.HasChange("name") {
		datacenterBuilder.Name(d.Get("name").(string))
	}

	if description, ok := d.GetOk("description"); ok && d.HasChange("description") {
		datacenterBuilder.Description(description.(string))
	}

	if d.HasChange("local") {
		datacenterBuilder.Local(d.Get("local").(bool))
	}

	datacenter, err := datacenterBuilder.Build()
	if err != nil {
		return err
	}

	_, err = datacenterService.Update().DataCenter(datacenter).Send()

	if err != nil {
		return err
	}

	return resourceOvirtDataCenterRead(d, meta)
}

func resourceOvirtDataCenterRead(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*ovirtsdk4.Connection)
	getDataCenterResp, err := conn.SystemService().DataCentersService().
		DataCenterService(d.Id()).Get().Send()
	if err != nil {
		if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
			d.SetId("")
			return nil
		}
		return err
	}

	datacenter, ok := getDataCenterResp.DataCenter()
	if !ok {
		d.SetId("")
		return nil
	}

	d.Set("name", datacenter.MustName())
	d.Set("local", datacenter.MustLocal())

	if description, ok := datacenter.Description(); ok {
		d.Set("description", description)
	}
	d.Set("status", string(datacenter.MustStatus()))

	return nil
}

func resourceOvirtDataCenterDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	_, err := conn.SystemService().DataCentersService().
		DataCenterService(d.Id()).Remove().Send()
	if err != nil {
		if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
			return nil
		}
		return err
	}
	return nil
}
