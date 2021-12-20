// Copyright (C) 2018 Joey Ma <majunjiev@gmail.com>
// All rights reserved.
//
// This software may be modified and distributed under the terms
// of the BSD-2 license.  See the LICENSE file for details.

package ovirt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

func resourceOvirtUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceOvirtUserCreate,
		Read:   resourceOvirtUserRead,
		Delete: resourceOvirtUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"authz_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceOvirtUserCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	builder := ovirtsdk4.NewUserBuilder()

	builder.
		Principal(d.Get("name").(string)).
		Namespace(d.Get("namespace").(string)).
		UserName(
			fmt.Sprintf("%s@%s",
				d.Get("name").(string),
				d.Get("authz_name").(string))).
		Domain(
			ovirtsdk4.NewDomainBuilder().
				Name(d.Get("authz_name").(string)).
				MustBuild())

	resp, err := conn.SystemService().UsersService().
		Add().
		User(builder.MustBuild()).
		Send()
	if err != nil {
		log.Printf("[DEBUG] Error adding user (%s): %s", d.Get("name").(string), err)
		return err
	}

	d.SetId(resp.MustUser().MustId())
	return resourceOvirtUserRead(d, meta)
}

func resourceOvirtUserRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	resp, err := conn.SystemService().
		UsersService().
		UserService(d.Id()).
		Get().
		Send()
	if err != nil {
		if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", resp.MustUser().MustPrincipal())
	d.Set("namespace", resp.MustUser().MustNamespace())
	d.Set("authz_name", resp.MustUser().MustDomain().MustName())

	return nil
}

func resourceOvirtUserDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	_, err := conn.SystemService().
		UsersService().
		UserService(d.Id()).
		Remove().
		Send()
	if err != nil {
		if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
			return nil
		}
		log.Printf("[DEBUG] Error deleting User (%s): %s", d.Id(), err)
		return err
	}

	return nil
}
