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

func resourceOvirtTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceOvirtTagCreate,
		Read:   resourceOvirtTagRead,
		Update: resourceOvirtTagUpdate,
		Delete: resourceOvirtTagDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "",
			},
			"vm_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"host_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceOvirtTagCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)
	systemService := conn.SystemService()

	tagBuilder := ovirtsdk4.NewTagBuilder().
		Name(d.Get("name").(string)).
		Description(d.Get("description").(string))

	if pid, ok := d.GetOk("parent_id"); ok {
		tagBuilder.Parent(
			ovirtsdk4.NewTagBuilder().
				Id(pid.(string)).
				MustBuild())
	}

	addResp, err := systemService.TagsService().
		Add().
		Tag(tagBuilder.MustBuild()).
		Send()
	if err != nil {
		log.Printf("[DEBUG] Error adding Tag (%s): %s", d.Get("name").(string), err)
		return err
	}

	d.SetId(addResp.MustTag().MustId())

	// Attach it to vms
	if ids, ok := d.GetOk("vm_ids"); ok {
		attachTagToVMs(systemService.VmsService(), d.Id(), ids.(*schema.Set).List())
	}

	// Attach it to hosts
	if ids, ok := d.GetOk("host_ids"); ok {
		attachTagToHosts(systemService.HostsService(), d.Id(), ids.(*schema.Set).List())
	}

	return resourceOvirtTagRead(d, meta)
}

func resourceOvirtTagUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)
	systemService := conn.SystemService()

	d.Partial(true)
	attributeUpdate := false
	paramBuilder := ovirtsdk4.NewTagBuilder()
	if d.HasChange("name") {
		paramBuilder.Name(d.Get("name").(string))
		attributeUpdate = true
	}
	if d.HasChange("description") {
		paramBuilder.Description(d.Get("description").(string))
		attributeUpdate = true
	}
	if d.HasChange("parent_id") {
		paramBuilder.Parent(ovirtsdk4.NewTagBuilder().
			Id(d.Get("parent_id").(string)).
			MustBuild())
		attributeUpdate = true
	}
	if attributeUpdate {
		_, err := systemService.TagsService().
			TagService(d.Id()).
			Update().
			Tag(paramBuilder.MustBuild()).
			Send()
		if err != nil {
			return err
		}
		d.SetPartial("name")
		d.SetPartial("description")
		d.SetPartial("parent_id")
	}

	if d.HasChange("vm_ids") {
		vmsService := systemService.VmsService()
		o, n := d.GetChange("vm_ids")
		os, ns := o.(*schema.Set), n.(*schema.Set)
		removed, added := os.Difference(ns).List(), ns.Difference(os).List()
		if len(removed) > 0 {
			if err := detachTagFromVMs(vmsService, d.Id(), removed); err != nil {
				return err
			}
		}
		if len(added) > 0 {
			if err := attachTagToVMs(vmsService, d.Id(), added); err != nil {
				return err
			}
		}
		d.SetPartial("vm_ids")
	}

	if d.HasChange("host_ids") {
		hostsService := systemService.HostsService()
		o, n := d.GetChange("host_ids")
		os, ns := o.(*schema.Set), n.(*schema.Set)
		removed, added := os.Difference(ns).List(), ns.Difference(os).List()
		if len(removed) > 0 {
			if err := detachTagFromHosts(hostsService, d.Id(), removed); err != nil {
				return err
			}
		}
		if len(added) > 0 {
			if err := attachTagToHosts(hostsService, d.Id(), added); err != nil {
				return err
			}
		}
		d.SetPartial("host_ids")
	}

	d.Partial(false)
	return resourceOvirtTagRead(d, meta)
}

func resourceOvirtTagRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)
	systemService := conn.SystemService()
	getTagResp, err := systemService.TagsService().TagService(d.Id()).Get().Send()
	if err != nil {
		if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
			d.SetId("")
			return nil
		}
	}
	tag := getTagResp.MustTag()

	d.Set("name", tag.MustName())
	if parent, ok := tag.Parent(); ok {
		d.Set("parent_id", parent.MustId())
	}
	if desc, ok := tag.Description(); ok {
		d.Set("description", desc)
	}

	vmIDs, err := searchVmsByTag(systemService.VmsService(), tag.MustName())
	if err != nil {
		log.Printf("[DEBUG] Error searching VMs via Tag (%s): %s", tag.MustName(), err)
		return err
	}
	if len(vmIDs) > 0 {
		d.Set("vm_ids", vmIDs)
	}

	hostIDs, err := searchHostsByTag(systemService.HostsService(), tag.MustName())
	if err != nil {
		log.Printf("[DEBUG] Error searching Hosts via Tag (%s): %s", tag.MustName(), err)
	}
	if len(hostIDs) > 0 {
		d.Set("host_ids", hostIDs)
	}

	return nil
}

func resourceOvirtTagDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)
	_, err := conn.SystemService().TagsService().
		TagService(d.Id()).
		Remove().
		Send()
	if err != nil {
		if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
			d.SetId("")
			return nil
		}
		return err
	}
	return nil

}

func attachTagToVMs(vmsServices *ovirtsdk4.VmsService, tagID string, vmIDs []interface{}) error {
	for _, v := range vmIDs {
		_, err := vmsServices.
			VmService(v.(string)).
			TagsService().
			Add().
			Tag(ovirtsdk4.NewTagBuilder().
				Id(tagID).
				MustBuild()).
			Send()
		if err != nil {
			log.Printf("[DEBUG] Error attaching Tag (%s) to VM (%s): %s",
				tagID, v.(string), err)
			return err
		}
	}
	return nil
}

func detachTagFromVMs(vmsServices *ovirtsdk4.VmsService, tagID string, vmIDs []interface{}) error {
	for _, v := range vmIDs {
		_, err := vmsServices.VmService(v.(string)).
			TagsService().
			TagService(tagID).
			Remove().
			Send()
		if err != nil {
			log.Printf("[DEBUG] Error detaching Tag (%s) from VM (%s): %s", tagID, v.(string), err)
			return err
		}
	}
	return nil
}

func attachTagToHosts(hostsService *ovirtsdk4.HostsService, tagID string, hostIDs []interface{}) error {
	for _, v := range hostIDs {
		_, err := hostsService.
			HostService(v.(string)).
			TagsService().
			Add().
			Tag(ovirtsdk4.NewTagBuilder().
				Id(tagID).
				MustBuild()).
			Send()
		if err != nil {
			log.Printf("[DEBUG] Error attaching Tag (%s) to Host (%s): %s",
				tagID, v.(string), err)
			return err
		}
	}
	return nil
}

func detachTagFromHosts(hostsService *ovirtsdk4.HostsService, tagID string, hostIDs []interface{}) error {
	for _, v := range hostIDs {
		_, err := hostsService.
			HostService(v.(string)).
			TagsService().
			TagService(tagID).
			Remove().
			Send()
		if err != nil {
			log.Printf("[DEBUG] Error dettaching Tag (%s) from Host (%s): %s",
				tagID, v.(string), err)
			return err
		}
	}
	return nil
}

func searchVmsByTag(service *ovirtsdk4.VmsService, tagName string) ([]string, error) {
	var vmIDs []string
	resp, err := service.List().Search(fmt.Sprintf("tag=%s", tagName)).Send()
	if err != nil {
		return nil, err
	}
	for _, v := range resp.MustVms().Slice() {
		vmIDs = append(vmIDs, v.MustId())
	}
	return vmIDs, nil
}

func searchHostsByTag(service *ovirtsdk4.HostsService, tagName string) ([]string, error) {
	var hostIDs []string
	resp, err := service.List().Search(fmt.Sprintf("tag=%s", tagName)).Send()
	if err != nil {
		return nil, err
	}
	for _, v := range resp.MustHosts().Slice() {
		hostIDs = append(hostIDs, v.MustId())
	}
	return hostIDs, nil
}
