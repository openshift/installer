// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/helpers/location"
	"github.com/softlayer/softlayer-go/helpers/product"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
)

func resourceIBMComputeBareMetal() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMComputeBareMetalCreate,
		Read:     resourceIBMComputeBareMetalRead,
		Update:   resourceIBMComputeBareMetalUpdate,
		Delete:   resourceIBMComputeBareMetalDelete,
		Exists:   resourceIBMComputeBareMetalExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{

			"hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				DefaultFunc: genID,
				DiffSuppressFunc: func(k, o, n string, d *schema.ResourceData) bool {
					// FIXME: Work around another bug in terraform.
					// When a default function is used with an optional property,
					// terraform will always execute it on apply, even when the property
					// already has a value in the state for it. This causes a false diff.
					// Making the property Computed:true does not make a difference.
					if strings.HasPrefix(o, "terraformed-") && strings.HasPrefix(n, "terraformed-") {
						return true
					}

					return o == n
				},
				Description: "Host name",
			},

			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Domain name",
			},

			"ssh_key_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				ForceNew:    true,
				Description: "SSH KEY IDS list",
			},

			"user_metadata": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "User metadata info",
			},

			"notes": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional notes info",
			},

			"file_storage_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Set: func(v interface{}) int {
					return v.(int)
				},
			},

			"block_storage_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Set: func(v interface{}) int {
					return v.(int)
				},
			},

			"post_install_script_uri": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          nil,
				ForceNew:         true,
				DiffSuppressFunc: applyOnce,
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			// Hourly only
			"fixed_config_preset": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: applyOnce,
				Description:      "Fixed config preset value",
			},

			// Hourly only
			"os_reference_code": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ConflictsWith:    []string{"image_template_id"},
				DiffSuppressFunc: applyOnce,
				Description:      "OS refernece code value",
			},

			"image_template_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"os_reference_code"},
				Description:   "OS image template ID",
			},

			"datacenter": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"network_speed": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
				ForceNew:    true,
				Description: "Network speed in MBPS",
			},

			"hourly_billing": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				ForceNew:    true,
				Description: "Enables hourly billing",
			},

			"private_network_only": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "only private network configured if is true",
			},

			"tcp_monitoring": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          false,
				ForceNew:         true,
				DiffSuppressFunc: applyOnce,
				Description:      "TCP monitoring enabled if set as true",
			},

			"redundant_power_supply": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"software_guard_extensions": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          false,
				ForceNew:         true,
				DiffSuppressFunc: applyOnce,
			},

			// Monthly only
			"package_key_name": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: applyOnce,
			},

			// Monthly only
			"process_key_name": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: applyOnce,
			},

			// Monthly only
			"os_key_name": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: applyOnce,
			},

			// Monthly only
			"gpu_key_name": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: applyOnce,
			},

			// Monthly only
			"gpu_secondary_key_name": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: applyOnce,
			},

			// Monthly only
			"disk_key_names": {
				Type:             schema.TypeList,
				Optional:         true,
				ForceNew:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: applyOnce,
			},

			// Monthly/Hourly  only
			"redundant_network": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          false,
				ForceNew:         true,
				DiffSuppressFunc: applyOnce,
			},

			// Monthly/Hourly only
			"unbonded_network": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          false,
				ForceNew:         true,
				DiffSuppressFunc: applyOnce,
			},

			// Monthly only. For controlling datacenter restricted port speed
			"restricted_network": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"extended_hardware_testing": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			// Monthly only
			"public_bandwidth": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				DiffSuppressFunc: applyOnce,
			},

			// Monthly only
			"memory": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				//Sometime memory returns back as different. Since this resource is immutable at this point
				//and memory can't be really updated , suppress the change until we figure out how to handle it
				DiffSuppressFunc: applyOnce,
				Computed:         true,
			},
			// Monthly only
			"storage_groups": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"array_type_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Array type ID",
						},
						"hard_drives": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Required:    true,
							Description: "Hard-drives List",
						},
						"array_size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Array size of harddrives list",
						},
						"partition_template_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Partition template ID",
						},
					},
				},
				DiffSuppressFunc: applyOnce,
			},

			// Quote based provisioning only
			"quote_id": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: applyOnce,
				Description:      "Quote ID for Quote based provisioning",
			},

			// Quote based provisioning, Monthly
			"public_vlan_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			// Quote based provisioning, Monthly
			"public_subnet": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			// Quote based provisioning, Monthly
			"private_vlan_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			// Quote based provisioning, Monthly
			"private_subnet": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"public_ipv4_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"public_ipv4_address_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"private_ipv4_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"private_ipv4_address_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"secondary_ip_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateSecondaryIPCount,
				DiffSuppressFunc: func(k, o, n string, d *schema.ResourceData) bool {
					// secondary_ip_count is only used when a virtual_guest resource is created.
					if d.State() == nil {
						return false
					}
					return true
				},
				Description: "Secondary IP addresses count",
			},
			"secondary_ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ipv6_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Boolean value true if IPV6 ia enabled or false",
			},
			"ipv6_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_address_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ipv6_static_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "boolean value true if ipv6 static is enabled else false",
			},

			"global_identifier": &schema.Schema{
				Description: "The unique global identifier of the bare metal server",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func getBareMetalOrderFromResourceData(d *schema.ResourceData, meta interface{}) (datatypes.Hardware, error) {
	dc := datatypes.Location{
		Name: sl.String(d.Get("datacenter").(string)),
	}

	networkComponent := datatypes.Network_Component{
		MaxSpeed: sl.Int(d.Get("network_speed").(int)),
	}

	hardware := datatypes.Hardware{
		Hostname:               sl.String(d.Get("hostname").(string)),
		Domain:                 sl.String(d.Get("domain").(string)),
		HourlyBillingFlag:      sl.Bool(d.Get("hourly_billing").(bool)),
		PrivateNetworkOnlyFlag: sl.Bool(d.Get("private_network_only").(bool)),
		Datacenter:             &dc,
		NetworkComponents:      []datatypes.Network_Component{networkComponent},
		PostInstallScriptUri:   sl.String(d.Get("post_install_script_uri").(string)),
		BareMetalInstanceFlag:  sl.Int(1),
		FixedConfigurationPreset: &datatypes.Product_Package_Preset{
			KeyName: sl.String(d.Get("fixed_config_preset").(string)),
		},
	}

	if operatingSystemReferenceCode, ok := d.GetOk("os_reference_code"); ok {
		hardware.OperatingSystemReferenceCode = sl.String(operatingSystemReferenceCode.(string))
	}

	public_vlan_id := d.Get("public_vlan_id").(int)
	if public_vlan_id > 0 {
		hardware.PrimaryNetworkComponent = &datatypes.Network_Component{
			NetworkVlan: &datatypes.Network_Vlan{Id: sl.Int(public_vlan_id)},
		}
	}

	private_vlan_id := d.Get("private_vlan_id").(int)
	if private_vlan_id > 0 {
		hardware.PrimaryBackendNetworkComponent = &datatypes.Network_Component{
			NetworkVlan: &datatypes.Network_Vlan{Id: sl.Int(private_vlan_id)},
		}
	}

	if public_subnet, ok := d.GetOk("public_subnet"); ok {
		subnet := public_subnet.(string)
		subnetID, err := getSubnetID(subnet, meta)
		if err != nil {
			return hardware, fmt.Errorf("Error determining id for subnet %s: %s", subnet, err)
		}

		hardware.PrimaryNetworkComponent.NetworkVlan.PrimarySubnetId = sl.Int(subnetID)
	}

	if private_subnet, ok := d.GetOk("private_subnet"); ok {
		subnet := private_subnet.(string)
		subnetID, err := getSubnetID(subnet, meta)
		if err != nil {
			return hardware, fmt.Errorf("Error determining id for subnet %s: %s", subnet, err)
		}

		hardware.PrimaryBackendNetworkComponent.NetworkVlan.PrimarySubnetId = sl.Int(subnetID)
	}

	if userMetadata, ok := d.GetOk("user_metadata"); ok {
		hardware.UserData = []datatypes.Hardware_Attribute{
			{Value: sl.String(userMetadata.(string))},
		}
	}

	// Get configured ssh_keys
	ssh_key_ids := d.Get("ssh_key_ids").([]interface{})
	if len(ssh_key_ids) > 0 {
		hardware.SshKeys = make([]datatypes.Security_Ssh_Key, 0, len(ssh_key_ids))
		for _, ssh_key_id := range ssh_key_ids {
			hardware.SshKeys = append(hardware.SshKeys, datatypes.Security_Ssh_Key{
				Id: sl.Int(ssh_key_id.(int)),
			})
		}
	}

	return hardware, nil
}

func resourceIBMComputeBareMetalCreate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	hwService := services.GetHardwareService(sess)
	var order datatypes.Container_Product_Order
	var err error
	quote_id := d.Get("quote_id").(int)
	hardware := datatypes.Hardware{
		Hostname: sl.String(d.Get("hostname").(string)),
		Domain:   sl.String(d.Get("domain").(string)),
	}

	if quote_id > 0 {
		// Build a bare metal template from the quote.
		order, err = services.GetBillingOrderQuoteService(sess).
			Id(quote_id).GetRecalculatedOrderContainer(nil, sl.Bool(false))
		if err != nil {
			return fmt.Errorf(
				"Encountered problem trying to get the bare metal order template from quote: %s", err)
		}
		order.Quantity = sl.Int(1)
		order.Hardware = make([]datatypes.Hardware, 0, 1)
		order.Hardware = append(
			order.Hardware,
			hardware,
		)
	} else if _, ok := d.GetOk("fixed_config_preset"); ok {
		// Build an hourly bare metal server template using fixed_config_preset.
		hardware, err = getBareMetalOrderFromResourceData(d, meta)
		if err != nil {
			return err
		}
		order, err = services.GetHardwareService(sess).GenerateOrderTemplate(&hardware)
		if err != nil {
			return fmt.Errorf(
				"Encountered problem trying to get the bare metal order template: %s", err)
		}
		items, err := product.GetPackageProducts(sess, *order.PackageId, productItemMaskWithPriceLocationGroupID)
		if err != nil {
			return err
		}
		redundantNetwork := d.Get("redundant_network").(bool)
		unbondedNetwork := d.Get("unbonded_network").(bool)

		if redundantNetwork || unbondedNetwork {
			// Remove network price
			prices := make([]datatypes.Product_Item_Price, len(order.Prices))
			i := 0
			for _, p := range order.Prices {
				if !strings.Contains(*p.Item.Description, "Network Uplink") {
					prices[i] = p
					i++
				}
			}
			portSpeed, err := findNetworkItemPriceId(items, d)
			if err != nil {
				return err
			}
			prices[i] = portSpeed
			order.Prices = prices
		}
		err = setMonthlyHourlyCommonOrder(d, items, &order)
		if err != nil {
			return err
		}

	} else {
		// Build a monthly bare metal server template
		order, err = getMonthlyBareMetalOrder(d, meta)
		if err != nil {
			return fmt.Errorf(
				"Encountered problem trying to get the custom bare metal order template: %s", err)
		}
	}

	order, err = setCommonBareMetalOrderOptions(d, meta, order)
	if err != nil {
		return fmt.Errorf(
			"Encountered problem trying to configure bare metal server options: %s", err)
	}

	log.Println("[INFO] Ordering bare metal server")
	orderReceipt, err := services.GetProductOrderService(sess.SetRetries(0)).PlaceOrder(&order, sl.Bool(false))
	if err != nil {
		return fmt.Errorf("Error ordering bare metal server: %s\n%+v\n", err, order)
	}

	gID := *orderReceipt.OrderDetails.Hardware[0].GlobalIdentifier

	log.Printf("[INFO] Bare Metal Server ID: %s", d.Id())
	log.Printf("[INFO] Bare Metal Server global ID: %s", gID)

	// wait for machine availability
	bm, err := waitForBareMetalProvision(&hardware, d, meta, gID)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for bare metal server (%s) to become ready: %s", d.Id(), err)
	}

	id := *bm.(datatypes.Hardware).Id
	d.SetId(fmt.Sprintf("%d", id))

	// Set tags
	if _, ok := d.GetOk("tags"); ok {
		err = setHardwareTags(id, d, meta)
		if err != nil {
			return err
		}
	}

	var storageIds []int
	if storageIdsSet := d.Get("file_storage_ids").(*schema.Set); len(storageIdsSet.List()) > 0 {
		storageIds = expandIntList(storageIdsSet.List())

	}
	if storageIdsSet := d.Get("block_storage_ids").(*schema.Set); len(storageIdsSet.List()) > 0 {
		storageIds = append(storageIds, expandIntList(storageIdsSet.List())...)
	}
	if len(storageIds) > 0 {
		err := addAccessToStorageList(hwService.Id(id), id, storageIds, meta)
		if err != nil {
			return err
		}
	}

	// Set notes
	if d.Get("notes").(string) != "" {
		err = setHardwareNotes(id, d, meta)
		if err != nil {
			return err
		}
	}

	return resourceIBMComputeBareMetalRead(d, meta)
}

func resourceIBMComputeBareMetalRead(d *schema.ResourceData, meta interface{}) error {
	service := services.GetHardwareService(meta.(ClientSession).SoftLayerSession())

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	result, err := service.Id(id).Mask(
		"hostname,domain,globalIdentifier," +
			"primaryIpAddress,primaryBackendIpAddress,privateNetworkOnlyFlag," +
			"notes,userData[value],tagReferences[id,tag[name]]," +
			"allowedNetworkStorage[id,nasType]," +
			"hourlyBillingFlag," +
			"datacenter[id,name,longName]," +
			"primaryNetworkComponent[primarySubnet[networkVlan[id,primaryRouter,vlanNumber],id]," +
			"primaryIpAddressRecord[id]," +
			"primaryVersion6IpAddressRecord[subnet,id]]," +
			"primaryBackendNetworkComponent[primarySubnet[networkVlan[id,primaryRouter,vlanNumber],id]," +
			"primaryIpAddressRecord[id]," +
			"maxSpeed,redundancyEnabledFlag]," +
			"memoryCapacity,powerSupplyCount," +
			"operatingSystem[softwareLicense[softwareDescription[referenceCode]]]",
	).GetObject()

	if err != nil {
		return fmt.Errorf("Error retrieving bare metal server: %s", err)
	}

	d.Set("hostname", *result.Hostname)
	d.Set("domain", *result.Domain)
	d.Set("global_identifier", result.GlobalIdentifier)

	if result.Datacenter != nil {
		d.Set("datacenter", *result.Datacenter.Name)
	}

	d.Set("network_speed", *result.PrimaryNetworkComponent.MaxSpeed)
	if result.PrimaryIpAddress != nil {
		d.Set("public_ipv4_address", *result.PrimaryIpAddress)
	}
	if result.PrimaryNetworkComponent.PrimaryIpAddressRecord != nil {
		d.Set("public_ipv4_address_id", *result.PrimaryNetworkComponent.PrimaryIpAddressRecord.Id)
	}
	d.Set("private_ipv4_address", *result.PrimaryBackendIpAddress)
	d.Set("private_ipv4_address_id",
		*result.PrimaryBackendNetworkComponent.PrimaryIpAddressRecord.Id)

	d.Set("private_network_only", *result.PrivateNetworkOnlyFlag)
	d.Set("hourly_billing", *result.HourlyBillingFlag)

	if result.PrimaryNetworkComponent.PrimarySubnet != nil {
		d.Set("public_vlan_id", *result.PrimaryNetworkComponent.PrimarySubnet.NetworkVlan.Id)
		d.Set("public_subnet", *result.PrimaryNetworkComponent.PrimarySubnet.Id)
	}

	if result.PrimaryBackendNetworkComponent.PrimarySubnet != nil {
		d.Set("private_vlan_id", *result.PrimaryBackendNetworkComponent.PrimarySubnet.NetworkVlan.Id)
		d.Set("private_subnet", *result.PrimaryBackendNetworkComponent.PrimarySubnet.Id)
	}

	userData := result.UserData
	if len(userData) > 0 && userData[0].Value != nil {
		d.Set("user_metadata", *userData[0].Value)
	}

	d.Set("notes", sl.Get(result.Notes, nil))
	d.Set("memory", *result.MemoryCapacity)

	d.Set("redundant_power_supply", false)

	if *result.PowerSupplyCount == 2 {
		d.Set("redundant_power_supply", true)
	}

	d.Set("redundant_network", false)
	d.Set("unbonded_network", false)

	backendNetworkComponent, err := service.Filter(
		filter.Build(
			filter.Path("backendNetworkComponents.status").Eq("ACTIVE"),
		),
	).Id(id).GetBackendNetworkComponents()

	if err != nil {
		return fmt.Errorf("Error retrieving bare metal server network: %s", err)
	}

	if len(backendNetworkComponent) > 2 && result.PrimaryBackendNetworkComponent != nil {
		if *result.PrimaryBackendNetworkComponent.RedundancyEnabledFlag {
			d.Set("redundant_network", true)
		} else {
			d.Set("unbonded_network", true)
		}
	}

	if result.OperatingSystem != nil &&
		result.OperatingSystem.SoftwareLicense != nil &&
		result.OperatingSystem.SoftwareLicense.SoftwareDescription != nil &&
		result.OperatingSystem.SoftwareLicense.SoftwareDescription.ReferenceCode != nil {
		d.Set("os_reference_code", *result.OperatingSystem.SoftwareLicense.SoftwareDescription.ReferenceCode)
	}

	tagReferences := result.TagReferences
	tagReferencesLen := len(tagReferences)
	if tagReferencesLen > 0 {
		tags := make([]string, 0, tagReferencesLen)
		for _, tagRef := range tagReferences {
			tags = append(tags, *tagRef.Tag.Name)
		}
		d.Set("tags", tags)
	}

	storages := result.AllowedNetworkStorage
	if len(storages) > 0 {
		d.Set("block_storage_ids", flattenBlockStorageID(storages))
		d.Set("file_storage_ids", flattenFileStorageID(storages))
	}

	connInfo := map[string]string{"type": "ssh"}
	if !*result.PrivateNetworkOnlyFlag && result.PrimaryIpAddress != nil {
		connInfo["host"] = *result.PrimaryIpAddress
	} else {
		connInfo["host"] = *result.PrimaryBackendIpAddress
	}
	d.SetConnInfo(connInfo)

	d.Set("ipv6_enabled", false)
	if result.PrimaryNetworkComponent.PrimaryVersion6IpAddressRecord != nil {
		d.Set("ipv6_enabled", true)
		d.Set("ipv6_address", *result.PrimaryNetworkComponent.PrimaryVersion6IpAddressRecord.IpAddress)
		d.Set("ipv6_address_id", *result.PrimaryNetworkComponent.PrimaryVersion6IpAddressRecord.Id)
	}
	err = readSecondaryIPAddresses(d, meta, result.PrimaryIpAddress)
	return err

}

func resourceIBMComputeBareMetalUpdate(d *schema.ResourceData, meta interface{}) error {
	id, _ := strconv.Atoi(d.Id())
	service := services.GetHardwareService(meta.(ClientSession).SoftLayerSession())

	if d.HasChange("tags") {
		err := setHardwareTags(id, d, meta)
		if err != nil {
			return err
		}
	}

	if d.HasChange("notes") {
		err := setHardwareNotes(id, d, meta)
		if err != nil {
			return err
		}
	}
	err := modifyStorageAccess(service.Id(id), id, meta, d)
	if err != nil {
		return err
	}

	return nil
}

func resourceIBMComputeBareMetalDelete(d *schema.ResourceData, meta interface{}) error {
	return deleteHardware(d, meta)
}

func deleteHardware(d dataRetriever, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetHardwareService(sess)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	_, err = waitForNoBareMetalActiveTransactions(id, meta)
	if err != nil {
		return fmt.Errorf("Error deleting bare metal server while waiting for zero active transactions: %s", err)
	}

	billingItem, err := service.Id(id).GetBillingItem()
	if err != nil {
		return fmt.Errorf("Error getting billing item for bare metal server: %s", err)
	}

	// Monthly bare metal servers only support an anniversary date cancellation option.
	billingItemService := services.GetBillingItemService(sess)
	_, err = billingItemService.Id(*billingItem.Id).CancelItem(
		sl.Bool(d.Get("hourly_billing").(bool)), sl.Bool(true), sl.String("No longer required"), sl.String("Please cancel this server"),
	)
	if err != nil {
		return fmt.Errorf("Error canceling the bare metal server (%d): %s", id, err)
	}

	return nil
}

func resourceIBMComputeBareMetalExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	service := services.GetHardwareService(meta.(ClientSession).SoftLayerSession())

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	result, err := service.Id(id).GetObject()
	if err != nil {
		if apiErr, ok := err.(sl.Error); !ok || apiErr.StatusCode != 404 {
			return false, fmt.Errorf("Error trying to retrieve the Bare Metal server: %s", err)
		}
	}

	return result.Id != nil && *result.Id == id, nil
}

// Bare metal creation does not return a bare metal object with an Id.
// Have to wait on provision date to become available on server that matches
// hostname and domain.
// http://sldn.softlayer.com/blog/bpotter/ordering-bare-metal-servers-using-softlayer-api
func waitForBareMetalProvision(hw *datatypes.Hardware, d *schema.ResourceData, meta interface{}, globalIdentifier string) (interface{}, error) {
	hostname := *hw.Hostname
	domain := *hw.Domain
	log.Printf("Waiting for server (%s.%s) to have to be provisioned", hostname, domain)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"retry", "pending"},
		Target:  []string{"provisioned"},
		Refresh: func() (interface{}, string, error) {
			sess := meta.(ClientSession).SoftLayerSession()
			service := services.GetAccountService(sess)
			bms, err := service.Filter(
				filter.Build(
					filter.Path("hardware.globalIdentifier").Eq(globalIdentifier))).Mask("id,provisionDate").GetHardware()
			if err != nil {
				return false, "retry", nil
			}

			if len(bms) == 0 || bms[0].ProvisionDate == nil {
				return datatypes.Hardware{}, "pending", nil
			}
			// Check Secondary IP address availability.
			if d.Get("secondary_ip_count").(int) > 0 {
				log.Println("Refreshing secondary IPs state.")
				secondarySubnetResult, err := services.GetAccountService(sess).
					Mask("ipAddresses[id,ipAddress]").
					Filter(filter.Build(filter.Path("publicSubnets.endPointIpAddress.hardware.id").Eq(bms[0].Id))).
					GetPublicSubnets()
				if err != nil {
					return nil, "", fmt.Errorf("Error retrieving secondary ip address: %s", err)
				}
				if len(secondarySubnetResult) == 0 {
					return datatypes.Hardware{}, "pending", nil
				}
			}

			return bms[0], "provisioned", nil

		},
		Timeout:        24 * time.Hour,
		Delay:          10 * time.Second,
		MinTimeout:     1 * time.Minute,
		NotFoundChecks: 24 * 60,
	}

	return stateConf.WaitForState()
}

func waitForNoBareMetalActiveTransactions(id int, meta interface{}) (interface{}, error) {
	log.Printf("Waiting for server (%d) to have zero active transactions", id)
	service := services.GetHardwareServerService(meta.(ClientSession).SoftLayerSession())

	stateConf := &resource.StateChangeConf{
		Pending: []string{"retry", "active"},
		Target:  []string{"idle"},
		Refresh: func() (interface{}, string, error) {
			bm, err := service.Id(id).Mask("id,activeTransactionCount").GetObject()
			if err != nil {
				return false, "retry", nil
			}

			if bm.ActiveTransactionCount != nil && *bm.ActiveTransactionCount == 0 {
				return bm, "idle", nil
			}
			return bm, "active", nil

		},
		Timeout:        24 * time.Hour,
		Delay:          10 * time.Second,
		MinTimeout:     1 * time.Minute,
		NotFoundChecks: 24 * 60,
	}

	return stateConf.WaitForState()
}

func setHardwareTags(id int, d dataRetriever, meta interface{}) error {
	service := services.GetHardwareService(meta.(ClientSession).SoftLayerSession())

	tags := getTags(d)
	_, err := service.Id(id).SetTags(sl.String(tags))
	if err != nil {
		return fmt.Errorf("Could not set tags on bare metal server %d", id)
	}

	return nil
}

func setHardwareNotes(id int, d dataRetriever, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetHardwareServerService(sess)

	result, err := service.Id(id).GetObject()
	if err != nil {
		return err
	}

	result.Notes = sl.String(d.Get("notes").(string))

	_, err = service.Id(id).EditObject(&result)
	if err != nil {
		return err
	}

	return nil
}

// Returns a price from an item list.
// Example usage : getItemPriceId(items, 'server', 'INTEL_XEON_2690_2_60')
func getItemPriceId(items []datatypes.Product_Item, categoryCode string, keyName string, capacity ...int) (datatypes.Product_Item_Price, error) {
	availableItems := ""
	for _, item := range items {
		for _, itemCategory := range item.Categories {
			if *itemCategory.CategoryCode == categoryCode {
				availableItems = availableItems + *item.KeyName + " ( " + *item.Description + " ) , "
				if *item.KeyName == keyName {
					for _, price := range item.Prices {
						capacityMin := -1
						capacityMax := -1
						var err error

						if price.CapacityRestrictionMinimum != nil && price.CapacityRestrictionMaximum != nil && (*price.CapacityRestrictionType == "CORE" || *price.CapacityRestrictionType == "PROCESSOR") {
							capacityMin, err = strconv.Atoi(*price.CapacityRestrictionMinimum)
							if err != nil {
								return datatypes.Product_Item_Price{}, err
							}
							capacityMax, err = strconv.Atoi(*price.CapacityRestrictionMaximum)
							if err != nil {
								return datatypes.Product_Item_Price{}, err
							}

						}
						for _, category := range price.Categories {
							if *category.CategoryCode == categoryCode && price.LocationGroupId == nil && !*price.BareMetalReservedCapacityFlag {

								if len(capacity) > 0 && capacityMin != -1 && capacityMax != -1 {

									if capacity[0] >= capacityMin && capacity[0] <= capacityMax {
										return datatypes.Product_Item_Price{Id: price.Id}, nil
									}

								} else {
									return datatypes.Product_Item_Price{Id: price.Id}, nil
								}

							}
						}
					}
				}
			}
		}
	}
	return datatypes.Product_Item_Price{},
		fmt.Errorf("Could not find the matching item with categorycode %s and keyName %s. Available item(s) is(are) %s", categoryCode, keyName, availableItems)
}

func getMonthlyBareMetalOrder(d *schema.ResourceData, meta interface{}) (datatypes.Container_Product_Order, error) {
	sess := meta.(ClientSession).SoftLayerSession()
	// Validate attributes for monthly bare metal server ordering.
	if d.Get("hourly_billing").(bool) {
		return datatypes.Container_Product_Order{}, fmt.Errorf("Monthly bare metal server only supports monthly billing.")
	}

	model, ok := d.GetOk("package_key_name")
	if !ok {
		return datatypes.Container_Product_Order{}, fmt.Errorf("The attribute 'package_key_name' is not defined.")
	}

	datacenter, ok := d.GetOk("datacenter")
	if !ok {
		return datatypes.Container_Product_Order{}, fmt.Errorf("The attribute 'datacenter' is not defined.")
	}

	osKeyName, ok := d.GetOk("os_key_name")
	if !ok {
		return datatypes.Container_Product_Order{}, fmt.Errorf("The attribute 'os_key_name' is not defined.")
	}

	dc, err := location.GetDatacenterByName(sess, datacenter.(string), "id")
	if err != nil {
		return datatypes.Container_Product_Order{}, err
	}

	// 1. Find a package id using monthly bare metal package key name.
	pkg, err := getPackageByModel(sess, model.(string))
	if err != nil {
		return datatypes.Container_Product_Order{}, err
	}

	if pkg.Id == nil {
		return datatypes.Container_Product_Order{}, err
	}

	// 2. Get all prices for the package
	items, err := product.GetPackageProducts(sess, *pkg.Id, productItemMaskWithPriceLocationGroupID)
	if err != nil {
		return datatypes.Container_Product_Order{}, err
	}
	// 3. Build price items
	server, err := getItemPriceId(items, "server", d.Get("process_key_name").(string))
	if err != nil {
		return datatypes.Container_Product_Order{}, err
	}

	coreCapacity := -1

	restrictionType := getCapacityRestrictionType(items, "os", osKeyName.(string))

	coreCapacity = getCoreCapacity(items, "server", d.Get("process_key_name").(string), restrictionType)
	if err != nil {
		return datatypes.Container_Product_Order{}, err
	}
	var os datatypes.Product_Item_Price

	if coreCapacity == -1 {
		os, err = getItemPriceId(items, "os", osKeyName.(string))
		if err != nil {
			return datatypes.Container_Product_Order{}, err
		}

	} else {
		os, err = getItemPriceId(items, "os", osKeyName.(string), coreCapacity)
		if err != nil {
			return datatypes.Container_Product_Order{}, err
		}

	}

	ram, err := findMemoryItemPriceId(items, d)
	if err != nil {
		return datatypes.Container_Product_Order{}, err
	}

	portSpeed, err := findNetworkItemPriceId(items, d)
	if err != nil {
		return datatypes.Container_Product_Order{}, err
	}
	var order datatypes.Container_Product_Order

	order = datatypes.Container_Product_Order{
		Quantity: sl.Int(1),
		Hardware: []datatypes.Hardware{{
			Hostname: sl.String(d.Get("hostname").(string)),
			Domain:   sl.String(d.Get("domain").(string)),
		},
		},
		Location:  sl.String(strconv.Itoa(*dc.Id)),
		PackageId: pkg.Id,
		Prices: []datatypes.Product_Item_Price{
			server,
			os,
			ram,
			portSpeed,
		},
	}

	if d.Get("tcp_monitoring").(bool) {
		monitoring, err := getItemPriceId(items, "monitoring", "MONITORING_HOST_PING_AND_TCP_SERVICE")
		if err != nil {
			return datatypes.Container_Product_Order{}, err
		}
		order.Prices = append(order.Prices, monitoring)

	}

	order = addCommomDefaultPrices(d, meta, order, items)

	// Add optional price ids.
	// Add public bandwidth
	privateNetworkOnly := d.Get("private_network_only").(bool)
	if publicBandwidth, ok := d.GetOk("public_bandwidth"); ok || privateNetworkOnly {
		publicBandwidthStr := "BANDWIDTH_" + strconv.Itoa(publicBandwidth.(int)) + "_GB"
		bandwidth, err := getItemPriceId(items, "bandwidth", publicBandwidthStr)
		if err != nil {
			return datatypes.Container_Product_Order{}, err
		}
		order.Prices = append(order.Prices, bandwidth)
	}

	// Add prices of disks.
	disks := d.Get("disk_key_names").([]interface{})
	diskLen := len(disks)
	if diskLen > 0 {
		for i, disk := range disks {
			diskPrice, err := getItemPriceId(items, "disk"+strconv.Itoa(i), disk.(string))
			if err != nil {
				return datatypes.Container_Product_Order{}, err
			}
			order.Prices = append(order.Prices, diskPrice)
		}
	}

	if _, ok := d.GetOk("storage_groups"); ok {
		order.StorageGroups = getStorageGroupsFromResourceData(d)
		diskController, err := getItemPriceId(items, "disk_controller", "DISK_CONTROLLER_RAID")
		if err != nil {
			return datatypes.Container_Product_Order{}, err
		}
		order.Prices = append(order.Prices, diskController)
	}

	err = setMonthlyHourlyCommonOrder(d, items, &order)
	if err != nil {
		return order, err
	}
	return order, nil
}

func setMonthlyHourlyCommonOrder(d *schema.ResourceData, items []datatypes.Product_Item, order *datatypes.Container_Product_Order) error {
	if d.Get("redundant_power_supply").(bool) {
		powerSupply, err := getItemPriceId(items, "power_supply", "REDUNDANT_POWER_SUPPLY")
		if err != nil {
			return err
		}
		order.Prices = append(order.Prices, powerSupply)
	}
	if d.Get("software_guard_extensions").(bool) {
		sgx, err := getItemPriceId(items, "software_guard_extensions", "SOFTWARE_GUARD_EXTENSIONS")
		if err != nil {
			return err
		}
		order.Prices = append(order.Prices, sgx)
	}
	if gpu0, ok := d.GetOk("gpu_key_name"); ok {
		gpu0Price, err := getItemPriceId(items, "gpu0", gpu0.(string))
		if err != nil {
			return err
		}
		order.Prices = append(order.Prices, gpu0Price)
	}

	if gpu1, ok := d.GetOk("gpu_secondary_key_name"); ok {
		gpu1Price, err := getItemPriceId(items, "gpu1", gpu1.(string))
		if err != nil {
			return err
		}
		order.Prices = append(order.Prices, gpu1Price)
	}

	secondaryIPCount := d.Get("secondary_ip_count").(int)
	privateNetworkOnly := d.Get("private_network_only").(bool)
	if secondaryIPCount > 0 {
		if privateNetworkOnly {
			return fmt.Errorf("Unable to configure public secondary addresses with a private_network_only option")
		}
		keyName := strconv.Itoa(secondaryIPCount) + "_PUBLIC_IP_ADDRESSES"

		price, err := getItemPriceId(items, "sec_ip_addresses", keyName)
		if err != nil {
			return err
		}
		order.Prices = append(order.Prices, price)
	}

	if d.Get("ipv6_enabled").(bool) {
		if privateNetworkOnly {
			return fmt.Errorf("Unable to configure a public IPv6 address with a private_network_only option")
		}
		keyName := "1_IPV6_ADDRESS"

		price, err := getItemPriceId(items, "pri_ipv6_addresses", keyName)
		if err != nil {
			return err
		}
		order.Prices = append(order.Prices, price)
	}

	if d.Get("ipv6_static_enabled").(bool) {
		if privateNetworkOnly {
			return fmt.Errorf("Unable to configure a public static IPv6 address with a private_network_only option")
		}
		keyName := "64_BLOCK_STATIC_PUBLIC_IPV6_ADDRESSES"

		price, err := getItemPriceId(items, "static_ipv6_addresses", keyName)
		if err != nil {
			return err
		}
		order.Prices = append(order.Prices, price)
	}
	return nil
}

// Set common parameters for server ordering.
func setCommonBareMetalOrderOptions(d *schema.ResourceData, meta interface{}, order datatypes.Container_Product_Order) (datatypes.Container_Product_Order, error) {

	extendedHardwareTesting := d.Get("extended_hardware_testing").(bool)
	order.ExtendedHardwareTesting = sl.Bool(extendedHardwareTesting)

	public_vlan_id := d.Get("public_vlan_id").(int)
	if public_vlan_id > 0 {
		order.Hardware[0].PrimaryNetworkComponent = &datatypes.Network_Component{
			NetworkVlan: &datatypes.Network_Vlan{Id: sl.Int(public_vlan_id)},
		}
	}

	private_vlan_id := d.Get("private_vlan_id").(int)
	if private_vlan_id > 0 {
		order.Hardware[0].PrimaryBackendNetworkComponent = &datatypes.Network_Component{
			NetworkVlan: &datatypes.Network_Vlan{Id: sl.Int(private_vlan_id)},
		}
	}

	if public_subnet, ok := d.GetOk("public_subnet"); ok {
		subnet := public_subnet.(string)
		subnetId, err := getSubnetId(subnet, meta)
		if err != nil {
			return datatypes.Container_Product_Order{}, fmt.Errorf("Error determining id for subnet %s: %s", subnet, err)
		}

		order.Hardware[0].PrimaryNetworkComponent.NetworkVlan.PrimarySubnetId = sl.Int(subnetId)
	}

	if private_subnet, ok := d.GetOk("private_subnet"); ok {
		subnet := private_subnet.(string)
		subnetId, err := getSubnetId(subnet, meta)
		if err != nil {
			return datatypes.Container_Product_Order{}, fmt.Errorf("Error determining id for subnet %s: %s", subnet, err)
		}

		order.Hardware[0].PrimaryBackendNetworkComponent.NetworkVlan.PrimarySubnetId = sl.Int(subnetId)
	}

	if userMetadata, ok := d.GetOk("user_metadata"); ok {
		order.Hardware[0].UserData = []datatypes.Hardware_Attribute{
			{Value: sl.String(userMetadata.(string))},
		}
	}

	// Get configured ssh_keys
	ssh_key_ids := d.Get("ssh_key_ids").([]interface{})
	if len(ssh_key_ids) > 0 {
		order.SshKeys = make([]datatypes.Container_Product_Order_SshKeys, 0, len(ssh_key_ids))
		for _, ssh_key_id := range ssh_key_ids {
			sshKeyA := make([]int, 1)
			sshKeyA[0] = ssh_key_id.(int)
			order.SshKeys = append(order.SshKeys, datatypes.Container_Product_Order_SshKeys{
				SshKeyIds: sshKeyA,
			})
		}
	}

	// Set image template id if it exists
	if rawImageTemplateId, ok := d.GetOk("image_template_id"); ok {
		imageTemplateId := rawImageTemplateId.(int)
		order.ImageTemplateId = sl.Int(imageTemplateId)
	}

	if postInstallURI, ok := d.GetOk("post_install_script_uri"); ok {
		postInstallURIA := make([]string, 1)
		postInstallURIA[0] = postInstallURI.(string)
		order.ProvisionScripts = postInstallURIA
	}

	return order, nil
}

// Find price item using network options
func findNetworkItemPriceId(items []datatypes.Product_Item, d dataRetriever) (datatypes.Product_Item_Price, error) {
	networkSpeed := d.Get("network_speed").(int)
	redundantNetwork := d.Get("redundant_network").(bool)
	unbondedNetwork := d.Get("unbonded_network").(bool)
	restrictedNetwork := d.Get("restricted_network").(bool)
	privateNetworkOnly := d.Get("private_network_only").(bool)

	networkSpeedStr := "_MBPS_"
	redundantNetworkStr := ""
	unbondedNetworkStr := ""
	restrictedNetworkStr := ""

	if networkSpeed < 1000 {
		networkSpeedStr = strconv.Itoa(networkSpeed) + networkSpeedStr
	} else {
		networkSpeedStr = strconv.Itoa(networkSpeed/1000) + "_GBPS"
	}
	if redundantNetwork {
		redundantNetworkStr = "_REDUNDANT"
	}

	if unbondedNetwork {
		unbondedNetworkStr = "_UNBONDED"
	}

	if restrictedNetwork {
		restrictedNetworkStr = "_NON_DATACENTER_RESTRICTED"
	}

	for _, item := range items {
		for _, itemCategory := range item.Categories {
			if *itemCategory.CategoryCode == "port_speed" &&
				strings.HasPrefix(*item.KeyName, networkSpeedStr) &&
				strings.Contains(*item.KeyName, redundantNetworkStr) &&
				strings.Contains(*item.KeyName, restrictedNetworkStr) &&
				strings.Contains(*item.KeyName, unbondedNetworkStr) {
				if (privateNetworkOnly && strings.Contains(*item.KeyName, "_PUBLIC_PRIVATE")) ||
					(!privateNetworkOnly && !strings.Contains(*item.KeyName, "_PUBLIC_PRIVATE")) ||
					(!unbondedNetwork && strings.Contains(*item.KeyName, "_UNBONDED")) ||
					(!redundantNetwork && strings.Contains(*item.KeyName, "_REDUNDANT")) ||
					(!restrictedNetwork && strings.Contains(*item.KeyName, "_DATACENTER_RESTRICTED")) {
					break
				}
				for _, price := range item.Prices {
					if price.LocationGroupId == nil && !*price.BareMetalReservedCapacityFlag {
						return datatypes.Product_Item_Price{Id: price.Id}, nil
					}
				}
			}
		}
	}
	return datatypes.Product_Item_Price{},
		fmt.Errorf("Could not find the network with %s, %s, %s, and private_network_only = %t",
			networkSpeedStr, redundantNetworkStr, unbondedNetworkStr, privateNetworkOnly)
}

// Find memory price item using memory size.
func findMemoryItemPriceId(items []datatypes.Product_Item, d dataRetriever) (datatypes.Product_Item_Price, error) {
	memory := d.Get("memory").(int)
	availableMemories := ""

	for _, item := range items {
		for _, itemCategory := range item.Categories {
			if *itemCategory.CategoryCode == "ram" {
				availableMemories = availableMemories + *item.KeyName + "(" + *item.Description + ")" + ", "
				if int(*item.Capacity) == memory {
					for _, price := range item.Prices {
						if price.LocationGroupId == nil && !*price.BareMetalReservedCapacityFlag {
							return datatypes.Product_Item_Price{Id: price.Id}, nil
						}
					}
				}
			}
		}
	}

	return datatypes.Product_Item_Price{},
		fmt.Errorf("Could not find the price item for %d GB memory. Available items are %s", memory, availableMemories)
}

// Find a bare metal package object using a package key name
func getPackageByModel(sess *session.Session, model string) (datatypes.Product_Package, error) {
	objectMask := "id,keyName,name,description,isActive,type[keyName]"
	service := services.GetProductPackageService(sess)
	availableModels := ""

	// Get package id
	packages, err := service.Mask(objectMask).
		Filter(
			filter.Build(
				filter.Path("type.keyName").Eq("BARE_METAL_CPU"),
			),
		).GetAllObjects()
	if err != nil {
		return datatypes.Product_Package{}, err
	}

	for _, pkg := range packages {
		availableModels = availableModels + *pkg.KeyName
		if pkg.Description != nil {
			availableModels = availableModels + " ( " + *pkg.Description + " ), "
		} else {
			availableModels = availableModels + ", "
		}
		if *pkg.KeyName == model {
			return pkg, nil
		}
	}

	return datatypes.Product_Package{}, fmt.Errorf("No custom bare metal package key name for %s. Available package key name(s) is(are) %s", model, availableModels)
}

func getStorageGroupsFromResourceData(d dataRetriever) []datatypes.Container_Product_Order_Storage_Group {
	storageGroupLists := d.Get("storage_groups").([]interface{})
	storageGroups := make([]datatypes.Container_Product_Order_Storage_Group, 0)

	for _, storageGroupList := range storageGroupLists {
		storageGroup := storageGroupList.(map[string]interface{})
		var storageGroupObj datatypes.Container_Product_Order_Storage_Group
		storageGroupObj.ArrayTypeId = sl.Int(storageGroup["array_type_id"].(int))
		hardDrives := storageGroup["hard_drives"].([]interface{})
		storageGroupObj.HardDrives = make([]int, 0, len(hardDrives))
		for _, hardDrive := range hardDrives {
			storageGroupObj.HardDrives = append(storageGroupObj.HardDrives, hardDrive.(int))
		}
		arraySize := storageGroup["array_size"].(int)
		if arraySize > 0 {
			storageGroupObj.ArraySize = sl.Float(float64(arraySize))
		}
		partitionTemplateId := storageGroup["partition_template_id"].(int)
		if partitionTemplateId > 0 {
			storageGroupObj.PartitionTemplateId = sl.Int(partitionTemplateId)
		}
		storageGroups = append(storageGroups, storageGroupObj)
	}
	return storageGroups
}

// Use this function for attributes which only should be applied in resource creation time.
func applyOnce(k, o, n string, d *schema.ResourceData) bool {
	if len(d.Id()) == 0 {
		return false
	}
	return true
}

func addCommomDefaultPrices(d *schema.ResourceData, meta interface{}, order datatypes.Container_Product_Order, items []datatypes.Product_Item) datatypes.Container_Product_Order {

	if !d.Get("tcp_monitoring").(bool) {
		monExists, moniotring := getCommonItemPriceID(items, "monitoring", "MONITORING_HOST_PING")

		if monExists {
			order.Prices = append(order.Prices, moniotring)
		}

	}

	priExists, priIPAddress := getCommonItemPriceID(items, "pri_ip_addresses", "1_IP_ADDRESS")
	if priExists {
		order.Prices = append(order.Prices, priIPAddress)
	}

	remotExists, remoteManagement := getCommonItemPriceID(items, "remote_management", "REBOOT_KVM_OVER_IP")
	if remotExists {
		order.Prices = append(order.Prices, remoteManagement)
	}

	vpnExists, vpnManagement := getCommonItemPriceID(items, "vpn_management", "UNLIMITED_SSL_VPN_USERS_1_PPTP_VPN_USER_PER_ACCOUNT")
	if vpnExists {
		order.Prices = append(order.Prices, vpnManagement)
	}

	notificationExists, notification := getCommonItemPriceID(items, "notification", "NOTIFICATION_EMAIL_AND_TICKET")
	if notificationExists {
		order.Prices = append(order.Prices, notification)
	}

	resExists, response := getCommonItemPriceID(items, "response", "AUTOMATED_NOTIFICATION")
	if resExists {
		order.Prices = append(order.Prices, response)
	}

	vulExists, vulnerabilityScanner := getCommonItemPriceID(items, "vulnerability_scanner", "NESSUS_VULNERABILITY_ASSESSMENT_REPORTING")
	if vulExists {
		order.Prices = append(order.Prices, vulnerabilityScanner)
	}

	sapExists, sapCertified := getCommonItemPriceID(items, "sap_certified_server", "SAP_HANA_NETWEAVER_CERTIFIED_SERVERS")
	if sapExists {
		order.Prices = append(order.Prices, sapCertified)
	}

	if _, ok := d.GetOk("storage_groups"); !ok {
		diskExists, diskController := getCommonItemPriceID(items, "disk_controller", "DISK_CONTROLLER_NONRAID")
		if diskExists {
			order.Prices = append(order.Prices, diskController)
		}

	}

	return order
}

// Returns a common default sprice from an item list.
// Example usage : getItemPriceId(items, 'server', 'INTEL_XEON_2690_2_60')
func getCommonItemPriceID(items []datatypes.Product_Item, categoryCode string, keyName string) (bool, datatypes.Product_Item_Price) {
	availableItems := ""
	for _, item := range items {
		for _, itemCategory := range item.Categories {
			if *itemCategory.CategoryCode == categoryCode {
				availableItems = availableItems + *item.KeyName + " ( " + *item.Description + " ) , "
				if *item.KeyName == keyName {
					for _, price := range item.Prices {
						for _, category := range price.Categories {
							if *category.CategoryCode == categoryCode && price.LocationGroupId == nil {
								return true, datatypes.Product_Item_Price{Id: price.Id}
							}
						}
					}
				}
			}
		}
	}
	return false, datatypes.Product_Item_Price{}
}

func getCoreCapacity(items []datatypes.Product_Item, categoryCode string, keyName string, restrictionType string) int {
	availableItems := ""
	for _, item := range items {
		for _, itemCategory := range item.Categories {
			if *itemCategory.CategoryCode == categoryCode {
				availableItems = availableItems + *item.KeyName + " ( " + *item.Description + " ) , "
				if *item.KeyName == keyName {
					if restrictionType == "PROCESSOR" && item.TotalProcessorCapacity != nil {
						return *item.TotalProcessorCapacity
					}
					if restrictionType == "CORE" && item.TotalPhysicalCoreCapacity != nil {
						return *item.TotalPhysicalCoreCapacity
					}

				}
			}
		}
	}
	return -1
}

func getCapacityRestrictionType(items []datatypes.Product_Item, categoryCode string, keyName string) string {
	for _, item := range items {
		for _, itemCategory := range item.Categories {
			if *itemCategory.CategoryCode == categoryCode {
				if *item.KeyName == keyName {
					for _, price := range item.Prices {
						if price.CapacityRestrictionType != nil {
							return *price.CapacityRestrictionType
						}

					}
				}
			}
		}
	}
	return ""
}
