// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/helpers/network"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/sl"
)

func resourceIBMStorageBlock() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMStorageBlockCreate,
		Read:     resourceIBMStorageBlockRead,
		Update:   resourceIBMStorageBlockUpdate,
		Delete:   resourceIBMStorageBlockDelete,
		Exists:   resourceIBMStorageBlockExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
			Update: schema.DefaultTimeout(45 * time.Minute),
			Delete: schema.DefaultTimeout(45 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Storage block type",
			},

			"datacenter": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Datacenter name",
			},

			"capacity": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Storage block size",
			},

			"iops": {
				Type:        schema.TypeFloat,
				Required:    true,
				Description: "IOPS value required",
			},

			"volumename": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Volume name",
			},

			"hostname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Hostname",
			},

			"lunid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "LUN Id",
			},

			"snapshot_capacity": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Snapshot capacity in GB",
			},

			"os_format_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "OS formatr type",
			},

			"allowed_virtual_guest_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Set: func(v interface{}) int {
					return v.(int)
				},
				Description: "List of allowed virtual guest IDs",
			},

			"notes": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Additional note info",
			},
			//TODO in v0.9.0
			"allowed_virtual_guest_info": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"password": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_iqn": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Set: func(v interface{}) int {
					virtualGuest := v.(map[string]interface{})
					return virtualGuest["id"].(int)
				},
				Deprecated: "Please use 'allowed_host_info' instead",
			},

			"allowed_hardware_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Set: func(v interface{}) int {
					return v.(int)
				},
				Description: "List of allowe hardware IDs",
			},

			//TODO in v0.9.0
			"allowed_hardware_info": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"password": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_iqn": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Set: func(v interface{}) int {
					baremetal := v.(map[string]interface{})
					return baremetal["id"].(int)
				},
				Deprecated: "Please use 'allowed_host_info' instead",
			},

			"allowed_ip_addresses": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Allowed IP addresses",
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "List of tags associated with the resource",
			},
			"hourly_billing": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Billing done hourly, if set to true",
			},
			"allowed_host_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"password": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_iqn": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"target_address": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of target Addresses",
			},
			ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about this instance",
			},
			ResourceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource",
			},
		},
	}
}

func resourceIBMStorageBlockCreate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()

	storageType := d.Get("type").(string)
	iops := d.Get("iops").(float64)
	datacenter := d.Get("datacenter").(string)
	capacity := d.Get("capacity").(int)
	snapshotCapacity := d.Get("snapshot_capacity").(int)
	osFormatType := d.Get("os_format_type").(string)
	osType, err := network.GetOsTypeByName(sess, osFormatType)
	hourlyBilling := d.Get("hourly_billing").(bool)

	if err != nil {
		return err
	}

	storageOrderContainer, err := buildStorageProductOrderContainer(sess, storageType, iops, capacity, snapshotCapacity, blockStorage, datacenter, hourlyBilling)
	if err != nil {
		return fmt.Errorf("Error while creating storage:%s", err)
	}

	log.Println("[INFO] Creating storage")

	var receipt datatypes.Container_Product_Order_Receipt

	switch storageType {
	case enduranceType:
		receipt, err = services.GetProductOrderService(sess.SetRetries(0)).PlaceOrder(
			&datatypes.Container_Product_Order_Network_Storage_AsAService{
				Container_Product_Order: storageOrderContainer,
				OsFormatType: &datatypes.Network_Storage_Iscsi_OS_Type{
					Id:      osType.Id,
					KeyName: osType.KeyName,
				},
				VolumeSize: &capacity,
			}, sl.Bool(false))
	case performanceType:
		receipt, err = services.GetProductOrderService(sess.SetRetries(0)).PlaceOrder(
			&datatypes.Container_Product_Order_Network_Storage_AsAService{
				Container_Product_Order: storageOrderContainer,
				OsFormatType: &datatypes.Network_Storage_Iscsi_OS_Type{
					Id:      osType.Id,
					KeyName: osType.KeyName,
				},
				Iops:       sl.Int(int(iops)),
				VolumeSize: &capacity,
			}, sl.Bool(false))
	default:
		return fmt.Errorf("Error during creation of storage: Invalid storageType %s", storageType)
	}

	if err != nil {
		return fmt.Errorf("Error during creation of storage: %s", err)
	}

	// Find the storage device
	blockStorage, err := findStorageByOrderId(sess, *receipt.OrderId, d.Timeout(schema.TimeoutCreate))

	if err != nil {
		return fmt.Errorf("Error during creation of storage: %s", err)
	}
	d.SetId(fmt.Sprintf("%d", *blockStorage.Id))

	// Wait for storage availability
	_, err = WaitForStorageAvailable(d, meta)

	if err != nil {
		return fmt.Errorf(
			"Error waiting for storage (%s) to become ready: %s", d.Id(), err)
	}

	// SoftLayer changes the device ID after completion of provisioning. It is necessary to refresh device ID.
	blockStorage, err = findStorageByOrderId(sess, *receipt.OrderId, d.Timeout(schema.TimeoutCreate))

	if err != nil {
		return fmt.Errorf("Error during creation of storage: %s", err)
	}
	d.SetId(fmt.Sprintf("%d", *blockStorage.Id))

	log.Printf("[INFO] Storage ID: %s", d.Id())

	return resourceIBMStorageBlockUpdate(d, meta)
}

func resourceIBMStorageBlockRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	storageId, _ := strconv.Atoi(d.Id())

	storage, err := services.GetNetworkStorageService(sess).
		Id(storageId).
		Mask(storageDetailMask).
		GetObject()

	if err != nil {
		return fmt.Errorf("Error retrieving storage information: %s", err)
	}

	storageType := strings.Fields(*storage.StorageType.Description)[0]

	// Calculate IOPS
	iops, err := getIops(storage, storageType)
	if err != nil {
		return fmt.Errorf("Error retrieving storage information: %s", err)
	}

	d.Set("type", storageType)
	d.Set("capacity", *storage.CapacityGb)
	d.Set("volumename", *storage.Username)
	d.Set("hostname", *storage.ServiceResourceBackendIpAddress)
	d.Set("lunid", *storage.LunId)
	d.Set("iops", iops)
	if storage.SnapshotCapacityGb != nil {
		snapshotCapacity, _ := strconv.Atoi(*storage.SnapshotCapacityGb)
		d.Set("snapshot_capacity", snapshotCapacity)
	}

	// Parse data center short name from ServiceResourceName. For example,
	// if SoftLayer API returns "'serviceResourceName': 'PerfStor Aggr aggr_staasdal0601_p01'",
	// the data center short name is "dal06".
	r, _ := regexp.Compile("[a-zA-Z]{3}[0-9]{2}")
	d.Set("datacenter", r.FindString(*storage.ServiceResourceName))

	allowedHostInfoList := make([]map[string]interface{}, 0)

	// Read allowed_ip_addresses
	allowedIpaddressesList := make([]string, 0, len(storage.AllowedIpAddresses))
	for _, allowedIpaddress := range storage.AllowedIpAddresses {
		singleHost := make(map[string]interface{})
		singleHost["id"] = *allowedIpaddress.SubnetId
		singleHost["username"] = *allowedIpaddress.AllowedHost.Credential.Username
		singleHost["password"] = *allowedIpaddress.AllowedHost.Credential.Password
		singleHost["host_iqn"] = *allowedIpaddress.AllowedHost.Name
		allowedHostInfoList = append(allowedHostInfoList, singleHost)
		allowedIpaddressesList = append(allowedIpaddressesList, *allowedIpaddress.IpAddress)
	}
	d.Set("allowed_ip_addresses", allowedIpaddressesList)

	// Read allowed_virtual_guest_ids and allowed_host_info
	allowedVirtualGuestInfoList := make([]map[string]interface{}, 0)
	allowedVirtualGuestIdsList := make([]int, 0, len(storage.AllowedVirtualGuests))

	for _, allowedVirtualGuest := range storage.AllowedVirtualGuests {
		singleVirtualGuest := make(map[string]interface{})
		singleVirtualGuest["id"] = *allowedVirtualGuest.Id
		singleVirtualGuest["username"] = *allowedVirtualGuest.AllowedHost.Credential.Username
		singleVirtualGuest["password"] = *allowedVirtualGuest.AllowedHost.Credential.Password
		singleVirtualGuest["host_iqn"] = *allowedVirtualGuest.AllowedHost.Name
		allowedHostInfoList = append(allowedHostInfoList, singleVirtualGuest)
		allowedVirtualGuestInfoList = append(allowedVirtualGuestInfoList, singleVirtualGuest)
		allowedVirtualGuestIdsList = append(allowedVirtualGuestIdsList, *allowedVirtualGuest.Id)
	}
	d.Set("allowed_virtual_guest_ids", allowedVirtualGuestIdsList)
	d.Set("allowed_virtual_guest_info", allowedVirtualGuestInfoList)

	// Read allowed_hardware_ids and allowed_host_info
	allowedHardwareInfoList := make([]map[string]interface{}, 0)
	allowedHardwareIdsList := make([]int, 0, len(storage.AllowedHardware))
	for _, allowedHW := range storage.AllowedHardware {
		singleHardware := make(map[string]interface{})
		singleHardware["id"] = *allowedHW.Id
		singleHardware["username"] = *allowedHW.AllowedHost.Credential.Username
		singleHardware["password"] = *allowedHW.AllowedHost.Credential.Password
		singleHardware["host_iqn"] = *allowedHW.AllowedHost.Name
		allowedHostInfoList = append(allowedHostInfoList, singleHardware)
		allowedHardwareInfoList = append(allowedHardwareInfoList, singleHardware)
		allowedHardwareIdsList = append(allowedHardwareIdsList, *allowedHW.Id)
	}
	d.Set("allowed_hardware_ids", allowedHardwareIdsList)
	d.Set("allowed_hardware_info", allowedHardwareInfoList)
	d.Set("allowed_host_info", allowedHostInfoList)

	if storage.OsType != nil {
		d.Set("os_format_type", *storage.OsType.Name)
	}

	if storage.Notes != nil {
		d.Set("notes", *storage.Notes)
	}

	if storage.BillingItem != nil {
		d.Set("hourly_billing", storage.BillingItem.HourlyFlag)
	}

	d.Set("target_address", storage.IscsiTargetIpAddresses)
	d.Set(ResourceControllerURL, fmt.Sprintf("https://cloud.ibm.com/classic/storage/block/%s", d.Id()))
	d.Set(ResourceName, *storage.ServiceResourceName)

	return nil
}

func resourceIBMStorageBlockUpdate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	storage, err := services.GetNetworkStorageService(sess).
		Id(id).
		Mask(storageDetailMask).
		GetObject()

	if err != nil {
		return fmt.Errorf("Error updating storage information: %s", err)
	}

	// Update allowed_ip_addresses
	if d.HasChange("allowed_ip_addresses") {
		err := updateAllowedIpAddresses(d, sess, storage)
		if err != nil {
			return fmt.Errorf("Error updating storage information: %s", err)
		}
	}

	// Update allowed_subnets
	if d.HasChange("allowed_subnets") {
		err := updateAllowedSubnets(d, sess, storage)
		if err != nil {
			return fmt.Errorf("Error updating storage information: %s", err)
		}
	}

	// Update allowed_virtual_guest_ids
	if d.HasChange("allowed_virtual_guest_ids") {
		err := updateAllowedVirtualGuestIds(d, sess, storage)
		if err != nil {
			return fmt.Errorf("Error updating storage information: %s", err)
		}
	}

	// Update allowed_hardware_ids
	if d.HasChange("allowed_hardware_ids") {
		err := updateAllowedHardwareIds(d, sess, storage)
		if err != nil {
			return fmt.Errorf("Error updating storage information: %s", err)
		}
	}

	// Update notes
	if d.HasChange("notes") {
		err := updateNotes(d, sess, storage)
		if err != nil {
			return fmt.Errorf("Error updating storage information: %s", err)
		}
	}

	if (d.HasChange("capacity") || d.HasChange("iops")) && !d.IsNewResource() {
		size := d.Get("capacity").(int)
		iops := d.Get("iops").(float64)

		modifyOrder, err := prepareModifyOrder(sess, storage, iops, size)
		if err != nil {
			return fmt.Errorf("Error updating storage: %s", err)
		}

		_, err = services.GetProductOrderService(sess.SetRetries(0)).PlaceOrder(
			&datatypes.Container_Product_Order_Network_Storage_AsAService_Upgrade{
				Container_Product_Order_Network_Storage_AsAService: modifyOrder,
				Volume: &datatypes.Network_Storage{
					Id: sl.Int(id),
				},
			}, sl.Bool(false))
		// Wait for storage availability
		_, err = WaitForStorageUpdate(d, meta)

		if err != nil {
			return fmt.Errorf(
				"Error waiting for storage (%s) to update: %s", d.Id(), err)
		}
	}

	return resourceIBMStorageBlockRead(d, meta)
}

func resourceIBMStorageBlockDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceIBMStorageFileDelete(d, meta)
}

func resourceIBMStorageBlockExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	return resourceIBMStorageFileExists(d, meta)
}
