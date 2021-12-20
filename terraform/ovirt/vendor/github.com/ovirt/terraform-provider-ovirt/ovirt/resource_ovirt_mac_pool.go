// Copyright (C) 2018 Joey Ma <majunjiev@gmail.com>
// All rights reserved.
//
// This software may be modified and distributed under the terms
// of the BSD-2 license.  See the LICENSE file for details.

package ovirt

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

func resourceOvirtMacPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceOvirtMacPoolCreate,
		Read:   resourceOvirtMacPoolRead,
		Update: resourceOvirtMacPoolUpdate,
		Delete: resourceOvirtMacPoolDelete,
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
				Default:  "",
			},
			"allow_duplicates": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  false,
			},
			"ranges": {
				Type:     schema.TypeSet,
				MinItems: 1,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: macRange(),
				},
				ForceNew:    false,
				Description: "MAC ranges that the from and to should be split by comma",
			},
		},
	}
}

func resourceOvirtMacPoolCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)
	mpBuilder := ovirtsdk4.NewMacPoolBuilder().
		Name(d.Get("name").(string)).
		Description(d.Get("description").(string)).
		AllowDuplicates(d.Get("allow_duplicates").(bool)).
		RangesOfAny(expandMacPoolRanges(d.Get("ranges").(*schema.Set))...)

	resp, err := conn.SystemService().
		MacPoolsService().
		Add().
		Pool(mpBuilder.MustBuild()).
		Send()
	if err != nil {
		log.Printf("[DEBUG] Error adding new MacPool (%s): %s", d.Get("name").(string), err)
		return err
	}
	d.SetId(resp.MustPool().MustId())

	return resourceOvirtMacPoolRead(d, meta)
}

func resourceOvirtMacPoolRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)
	resp, err := conn.SystemService().
		MacPoolsService().
		MacPoolService(d.Id()).
		Get().
		Send()
	if err != nil {
		if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
			d.SetId("")
			return nil
		}
		return err
	}

	mp := resp.MustPool()

	d.Set("name", mp.MustName())
	d.Set("allow_duplicates", mp.MustAllowDuplicates())
	d.Set("ranges", flattenMacPoolRanges(mp.MustRanges()))
	if desc, ok := mp.Description(); ok {
		d.Set("description", desc)
	}

	return nil
}

func resourceOvirtMacPoolUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	d.Partial(true)
	paramBuilder := ovirtsdk4.NewMacPoolBuilder()
	attributeUpdate := false

	if d.HasChange("name") {
		paramBuilder.Name(d.Get("name").(string))
		attributeUpdate = true
	}
	if d.HasChange("description") {
		paramBuilder.Description(d.Get("description").(string))
		attributeUpdate = true
	}
	if d.HasChange("allow_duplicates") {
		paramBuilder.AllowDuplicates(d.Get("allow_duplicates").(bool))
		attributeUpdate = true
	}
	if d.HasChange("ranges") {
		paramBuilder.RangesOfAny(expandMacPoolRanges(d.Get("ranges").(*schema.Set))...)
		attributeUpdate = true
	}

	if attributeUpdate {
		_, err := conn.SystemService().
			MacPoolsService().
			MacPoolService(d.Id()).
			Update().
			Pool(paramBuilder.MustBuild()).
			Send()
		if err != nil {
			log.Printf("[DEBUG] Error updating MacPool (%s): %s", d.Id(), err)
			return err
		}
	}

	d.Partial(false)
	return resourceOvirtMacPoolRead(d, meta)
}

func resourceOvirtMacPoolDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)
	_, err := conn.SystemService().
		MacPoolsService().
		MacPoolService(d.Id()).
		Remove().
		Send()
	if err != nil {
		if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
			return nil
		}
		return err
	}

	return nil
}

func expandMacPoolRanges(s *schema.Set) []*ovirtsdk4.Range {
	var ranges []*ovirtsdk4.Range
	for _, v := range s.List() {
		r := ovirtsdk4.NewRangeBuilder().
			From(strings.Split(v.(string), ",")[0]).
			To(strings.Split(v.(string), ",")[1]).
			MustBuild()
		ranges = append(ranges, r)
	}
	return ranges
}

func flattenMacPoolRanges(rs *ovirtsdk4.RangeSlice) []string {
	var ranges []string
	for _, r := range rs.Slice() {
		ranges = append(ranges, fmt.Sprintf("%s,%s", r.MustFrom(), r.MustTo()))
	}
	return ranges
}
