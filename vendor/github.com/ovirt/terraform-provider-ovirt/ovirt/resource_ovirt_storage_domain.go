// Copyright (C) 2018 Joey Ma <majunjiev@gmail.com>
// All rights reserved.
//
// This software may be modified and distributed under the terms
// of the BSD-2 license.  See the LICENSE file for details.

package ovirt

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

func resourceOvirtStorageDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceOvirtStorageDomainCreate,
		Read:   resourceOvirtStorageDomainRead,
		Delete: resourceOvirtStorageDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"host_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Host ID to be used to mount storage, this is not maintained by terraform",
			},
			"datacenter_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Datacenter ID where storage domain should be attached",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(ovirtsdk4.STORAGEDOMAINTYPE_DATA),
					string(ovirtsdk4.STORAGEDOMAINTYPE_ISO),
					string(ovirtsdk4.STORAGEDOMAINTYPE_EXPORT),
				}, false),
				Default:     string(ovirtsdk4.STORAGEDOMAINTYPE_DATA),
				Description: "The function of the storage domain",
			},
			// "localfs": {
			// 	Type:     schema.TypeList,
			// 	MinItems: 1,
			// 	MaxItems: 1,
			// 	Optional: true,
			// 	ForceNew: true,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"path": {
			// 				Type:     schema.TypeString,
			// 				Required: true,
			// 			},
			// 		},
			// 	},
			// 	ConflictsWith: []string{"nfs"},
			// 	Description:   "The attributes of localfs storage type",
			// },
			"nfs": {
				Type:     schema.TypeList,
				MinItems: 1,
				MaxItems: 1,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:     schema.TypeString,
							Required: true,
						},
						"path": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				// ConflictsWith: []string{"localfs"},
				Description: "The attributes of nfs storage type",
			},
			"wipe_after_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Boolean flag which indicates whether the storage domain should wipe the data after delete",
			},
		},
	}
}

func resourceOvirtStorageDomainCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	// Currently only nfs and localfs type of storage are supported
	storageType := ""
	for _, v := range []string{"nfs", "localfs"} {
		if _, ok := d.GetOk(v); ok {
			storageType = v
			break
		}
	}
	if storageType == "" {
		return fmt.Errorf("One and only one type of storage must be assigned")
	}

	// Build the storage domain entity
	sdBuilder := ovirtsdk4.NewStorageDomainBuilder().
		Name(d.Get("name").(string)).
		Type(ovirtsdk4.StorageDomainType(d.Get("type").(string))).
		DataCenter(
			ovirtsdk4.NewDataCenterBuilder().
				Id(d.Get("datacenter_id").(string)).
				MustBuild()).
		Host(
			ovirtsdk4.NewHostBuilder().
				Id(d.Get("host_id").(string)).
				MustBuild()).
		WipeAfterDelete(d.Get("wipe_after_delete").(bool))

	if v, ok := d.GetOk("description"); ok {
		sdBuilder.Description(v.(string))
	}

	// expand host storage
	s, err := expandOvirtHostStorage(storageType, d.Get(storageType).([]interface{}))
	if err != nil {
		return err
	}
	sdBuilder.Storage(s)

	sd, err := sdBuilder.Build()
	if err != nil {
		log.Printf("[DEBUG] Error building storage domain instance: %s", err)
		return err
	}

	resp, err := conn.SystemService().StorageDomainsService().Add().StorageDomain(sd).Send()
	if err != nil {
		log.Printf("[DEBUG] Error requesting for adding new storage domain: %s", err)
		return err
	}

	newSd, ok := resp.StorageDomain()
	if !ok {
		d.SetId("")
		return nil
	}
	d.SetId(newSd.MustId())

	// Attach the new StorageDomain
	log.Printf("[DEBUG] Check the Datacenter which the StorgeDomain will attached to")
	dcService := conn.SystemService().DataCentersService().
		DataCenterService(d.Get("datacenter_id").(string))
	dcResp, err := dcService.Get().Send()
	if err != nil {
		return fmt.Errorf("Error getting the Datacenter (%s): %s", d.Get("datacenter_id").(string), err)
	}
	if dcResp.MustDataCenter().MustStatus() != ovirtsdk4.DATACENTERSTATUS_UP {
		return fmt.Errorf("Error attaching the StorageDomain for the Datacenter (%s) status is not up", d.Get("datacenter_id").(string))
	}
	log.Printf("[DEBUG] Attach StorageDomain (%s) to Datacenter (%s)", newSd.MustId(), d.Get("datacenter_id").(string))
	_, err = dcService.StorageDomainsService().Add().
		StorageDomain(
			ovirtsdk4.NewStorageDomainBuilder().
				Id(newSd.MustId()).
				MustBuild()).
		Send()
	if err != nil {
		return err
	}

	// The storage domain is attached to the data center and is automatically activated.
	log.Printf("[DEBUG] Wait for StorageDomain (%s) status to become active", d.Id())
	activeStateConf := &resource.StateChangeConf{
		Target:     []string{string(ovirtsdk4.STORAGEDOMAINSTATUS_ACTIVE)},
		Refresh:    StorageDomainStateRefreshFunc(conn, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = activeStateConf.WaitForState()
	if err != nil {
		log.Printf("[DEBUG] Error waiting for StorageDomain (%s) to become active: %s", d.Id(), err)
		return err
	}
	log.Printf("[DEBUG] StorageDomain (%s) status has became to active", d.Id())

	return resourceOvirtStorageDomainRead(d, meta)
}

func resourceOvirtStorageDomainRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	resp, err := conn.SystemService().StorageDomainsService().
		StorageDomainService(d.Id()).
		Get().
		Send()
	if err != nil {
		if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] Error getting storage domain (%s): %s", d.Id(), err)
		return err
	}
	sd, ok := resp.StorageDomain()
	if !ok {
		d.SetId("")
		return nil
	}

	d.Set("name", sd.MustName())
	// d.Set("host_id", sd.MustHost().MustId())
	// TODO: use MustDataCenters to get a datacenter slice attached to
	d.Set("datacenter_id", sd.MustDataCenters().Slice()[0].MustId())
	d.Set("wipe_after_delete", sd.MustWipeAfterDelete())
	d.Set("type", sd.MustType())
	if v, ok := sd.Description(); ok {
		d.Set("description", v)
	}
	if err := d.Set(string(sd.MustStorage().MustType()), flattenOvirtHostStorage(sd.MustStorage())); err != nil {
		return fmt.Errorf("Error setting host storage: %s", err)
	}

	return nil
}

func resourceOvirtStorageDomainDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	dcID, err := findAttachedDatacenterByStorageDomain(conn, d.Id())
	if err != nil {
		if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
			d.SetId("")
			return nil
		}
		return err
	}

	// Only 'unattached' storage domain could be removed
	if dcID != "" {
		log.Printf("[DEBUG] Do maintenance for StorageDomain (%s) in Datacenter (%s)", d.Id(), dcID)
		// Maintain it
		err := maintenanceOvirtStorageDomain(conn, dcID, d.Id())
		if err != nil {
			if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
				d.SetId("")
				return nil
			}
			return err
		}
		log.Printf("[DEBUG] Dettach StorageDomain (%s) from Datacenter (%s)", d.Id(), dcID)
		// Unattach it
		err = unAttachedOvirtStorageDomain(conn, dcID, d.Id())
		if err != nil {
			return err
		}
	}

	log.Printf("[DEBUG] Determine the host to remove StorageDomain (%s)", d.Id())
	var hostID string
	if v, ok := d.GetOk("host_id"); ok {
		hostID = v.(string)
		log.Printf("[DEBUG] Host (%s by parameter) is to remove StorgeDomain (%s)", hostID, d.Id())
	} else {
		hostsResp, err := conn.SystemService().HostsService().List().Search("status=up").Send()
		if err != nil {
			log.Printf("[DEBUG] Error finding hosts with status=up: %s", err)
			return err
		}
		if hostSlice, ok := hostsResp.Hosts(); ok && len(hostSlice.Slice()) > 0 {
			hostID = hostSlice.Slice()[0].MustId()
			log.Printf("[DEBUG] Host (%s by searching) is to remove StorgeDomain (%s)", hostID, d.Id())
		} else {
			log.Printf("[DEBUG] No hosts found with status up")
			return fmt.Errorf("Not possible to remove StorgeDomain (%s) because no host found with status up", d.Id())
		}
	}

	log.Printf("[DEBUG] Remove StorageDomain (%s) with option destory=false", d.Id())
	_, err = conn.SystemService().StorageDomainsService().
		StorageDomainService(d.Id()).
		Remove().
		Destroy(false).
		Host(hostID).
		Send()
	if err != nil {
		log.Printf("[DEBUG] Error removing StorageDomain (%s): %s", d.Id(), err)
		return err
	}

	return nil
}

// StorageDomainStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an oVirt Storage Domain.
func StorageDomainStateRefreshFunc(conn *ovirtsdk4.Connection, sdID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		// First, try to get status from systemService.sdsService.sdService
		sdResp, err := conn.SystemService().StorageDomainsService().
			StorageDomainService(sdID).Get().Send()

		if err != nil {
			return nil, "", err
		}

		sd := sdResp.MustStorageDomain()
		if status, ok := sd.Status(); ok {
			return sd, string(status), nil
		}

		// Lack of status field indicates that the StorageDomain has been attached to
		// a Datacenter
		dcID, err := findAttachedDatacenterByStorageDomain(conn, sdID)
		if err != nil {
			return nil, "", err
		}

		attachedSdResp, err := conn.SystemService().DataCentersService().
			DataCenterService(dcID).
			StorageDomainsService().
			StorageDomainService(sdID).
			Get().
			Send()
		if err != nil {
			return nil, "", err
		}

		attachedSd := attachedSdResp.MustStorageDomain()
		return attachedSd, string(attachedSd.MustStatus()), nil
	}
}

// Finds the ID of the Datacenter that a given StorageDomain is attached to.
// The values returned:
// 	1. When the given StorageDomain does not exists:
//			<"", ovirtsdk4.NotFoundError>
//  2. When the given StorageDomain is not attached to any Datacenter:
//			<"", nil>
//  3. When everything is ok:
//			<DC-ID, nil>
func findAttachedDatacenterByStorageDomain(conn *ovirtsdk4.Connection, sdID string) (string, error) {
	r, err := conn.SystemService().
		StorageDomainsService().
		StorageDomainService(sdID).
		Get().
		Send()
	if err != nil {
		log.Printf("[DEBUG] Error getting StorageDomain (%s): %s", sdID, err)
		return "", err
	}
	sd := r.MustStorageDomain()
	dcs, ok := sd.DataCenters()
	if !ok || len(dcs.Slice()) == 0 {
		fmt.Printf("[DEBUG] StorageDomain (%s) has not been attached to any Datacenter", sdID)
		return "", nil
	}
	return dcs.Slice()[0].MustId(), nil
}

// This function fits all the kinds of storage types
func expandOvirtHostStorage(t string, l []interface{}) (*ovirtsdk4.HostStorage, error) {
	if len(l) == 0 {
		return nil, nil
	}
	s := l[0].(map[string]interface{})
	builder := ovirtsdk4.NewHostStorageBuilder()
	builder.Type(ovirtsdk4.StorageType(t))
	if v, ok := s["address"]; ok {
		builder.Address(v.(string))
	}
	if v, ok := s["path"]; ok {
		builder.Path(v.(string))
	}
	return builder.Build()
}

func flattenOvirtHostStorage(configured *ovirtsdk4.HostStorage) []map[string]interface{} {
	if configured == nil {
		attrs := make([]map[string]interface{}, 0)
		return attrs
	}
	attrs := make([]map[string]interface{}, 1)
	attr := make(map[string]interface{})

	if v, ok := configured.Address(); ok {
		attr["address"] = v
	}
	if v, ok := configured.Path(); ok {
		attr["path"] = v
	}

	attrs[0] = attr
	return attrs
}

func maintenanceOvirtStorageDomain(conn *ovirtsdk4.Connection, attachedDcID, sdID string) error {
	var dcID = attachedDcID
	if dcID == "" {
		var err error
		dcID, err = findAttachedDatacenterByStorageDomain(conn, sdID)
		if err != nil {
			return err
		}
	}
	if dcID == "" {
		// When the StorageDomain is not attached to any Datacenter, just return
		return nil
	}

	attachedSdService := conn.SystemService().DataCentersService().
		DataCenterService(dcID).
		StorageDomainsService().
		StorageDomainService(sdID)

	resp, err := attachedSdService.Get().Send()
	if err != nil {
		if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
			log.Printf("[DEBUG] StorageDomain (%s) has been unattached", sdID)
			return nil
		}
		return err
	}

	if resp.MustStorageDomain().MustStatus() == ovirtsdk4.STORAGEDOMAINSTATUS_MAINTENANCE {
		return nil
	}

	// Deactivate the StorageDomain
	log.Printf("[DEBUG] Deactivate (in maintenance) the StorageDomain (%s)", sdID)
	attachedSdService.Deactivate().Send()
	// The StorageDomain is attached to the Datacenter and is automatically activated.
	log.Printf("[DEBUG] Wait for StorageDomain (%s) status to become maintenance", sdID)
	maintenanceStateConf := &resource.StateChangeConf{
		Target:     []string{string(ovirtsdk4.STORAGEDOMAINSTATUS_MAINTENANCE)},
		Refresh:    StorageDomainStateRefreshFunc(conn, sdID),
		Timeout:    5 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = maintenanceStateConf.WaitForState()
	if err != nil {
		log.Printf("[DEBUG] Error waiting for StorageDomain (%s) to become maintenance: %s", sdID, err)
		return err
	}
	log.Printf("[DEBUG] StorageDomain (%s) status has became to maintenance", sdID)
	return nil
}

func unAttachedOvirtStorageDomain(conn *ovirtsdk4.Connection, attachedDcID, sdID string) error {
	var dcID = attachedDcID
	if dcID == "" {
		var err error
		dcID, err = findAttachedDatacenterByStorageDomain(conn, sdID)
		if err != nil {
			return err
		}
	}

	if dcID == "" {
		// When the StorageDomain is not attached to any Datacenter, just return
		return nil
	}

	attachedSdService := conn.SystemService().DataCentersService().
		DataCenterService(dcID).
		StorageDomainsService().
		StorageDomainService(sdID)

	resp, err := attachedSdService.Get().Send()
	if err != nil {
		if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
			log.Printf("[DEBUG] StorageDomain (%s) has been unattached", sdID)
			return nil
		}
		return err
	}

	if resp.MustStorageDomain().MustStatus() == ovirtsdk4.STORAGEDOMAINSTATUS_MAINTENANCE {
		attachedSdService.Remove().Send()
		log.Printf("[DEBUG] Wait for StorageDomain (%s) to be unattached", sdID)
		return resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err := attachedSdService.Get().Send()
			if err != nil {
				if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
					log.Printf("[DEBUG] StorageDomain (%s) has been unattached", sdID)
					return nil
				}
				return resource.RetryableError(fmt.Errorf("Error unattaching StorageDomain (%s): %s", sdID, err))
			}
			return resource.RetryableError(fmt.Errorf("StorageDomain (%s) is being unattached", sdID))
		})
	}
	return nil
}
