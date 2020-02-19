// Copyright (C) 2018 Joey Ma <majunjiev@gmail.com>
// All rights reserved.
//
// This software may be modified and distributed under the terms
// of the BSD-2 license.  See the LICENSE file for details.

package ovirt

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

func resourceOvirtVnicProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceOvirtVnicProfileCreate,
		Read:   resourceOvirtVnicProfileRead,
		Update: resourceOvirtVnicProfileUpdate,
		Delete: resourceOvirtVnicProfileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"migratable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"port_mirroring": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceOvirtVnicProfileCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	network := ovirtsdk4.NewNetworkBuilder().Id(d.Get("network_id").(string)).MustBuild()
	builder := ovirtsdk4.NewVnicProfileBuilder()

	builder.Network(network)
	builder.Name(d.Get("name").(string))
	if v, ok := d.GetOk("migratable"); ok {
		builder.Migratable(v.(bool))
	}
	if v, ok := d.GetOk("port_mirroring"); ok {
		builder.PortMirroring(v.(bool))
	}

	addResp, err := conn.SystemService().VnicProfilesService().
		Add().
		Profile(builder.MustBuild()).
		Send()
	if err != nil {
		return err
	}

	d.SetId(addResp.MustProfile().MustId())

	return resourceOvirtVnicProfileRead(d, meta)
}

func resourceOvirtVnicProfileRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)
	getResp, err := conn.SystemService().VnicProfilesService().
		ProfileService(d.Id()).
		Get().
		Send()
	if err != nil {
		if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
			d.SetId("")
			return nil
		}
		return err
	}
	profile, ok := getResp.Profile()
	if !ok {
		d.SetId("")
		return nil
	}
	d.Set("name", profile.MustName())
	d.Set("network_id", profile.MustNetwork().MustId())
	if v, ok := profile.Migratable(); ok {
		d.Set("migratable", v)
	}
	if v, ok := profile.PortMirroring(); ok {
		d.Set("port_mirroring", v)
	}

	return nil
}

func resourceOvirtVnicProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)
	profileService := conn.SystemService().VnicProfilesService().ProfileService(d.Id())

	newBuilder := ovirtsdk4.NewVnicProfileBuilder()
	if d.HasChange("name") {
		newBuilder.Name(d.Get("name").(string))
	}
	if d.HasChange("migratable") {
		newBuilder.Migratable(d.Get("migratable").(bool))
	}
	if d.HasChange("port_mirroring") {
		newBuilder.PortMirroring(d.Get("port_mirroring").(bool))
	}
	_, err := profileService.Update().
		Profile(newBuilder.MustBuild()).
		Send()
	if err != nil {
		return err
	}

	return resourceOvirtVnicProfileRead(d, meta)
}

func resourceOvirtVnicProfileDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	err := resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, e := conn.SystemService().VnicProfilesService().
			ProfileService(d.Id()).
			Remove().
			Send()
		if e != nil {
			if _, ok := e.(*ovirtsdk4.NotFoundError); ok {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("failed to delete vnicpfole: %s, wait for next run", e))
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
