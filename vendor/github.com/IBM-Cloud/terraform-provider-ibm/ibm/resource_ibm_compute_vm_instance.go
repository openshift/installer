// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/internal/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/helpers/product"
	"github.com/softlayer/softlayer-go/helpers/virtual"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
)

type storageIds []int

func (s storageIds) Storages(meta interface{}) ([]datatypes.Network_Storage, error) {
	storageService := services.GetNetworkStorageService(meta.(ClientSession).SoftLayerSession())
	storages := make([]datatypes.Network_Storage, len(s))

	for i, id := range s {
		var err error
		storages[i], err = storageService.Id(id).GetObject()
		if err != nil {
			return nil, err
		}
	}
	return storages, nil
}

const (
	staticIPRouted = "STATIC_IP_ROUTED"

	upgradeTransaction = "UPGRADE"
	pendingUpgrade     = "pending_upgrade"
	inProgressUpgrade  = "upgrade_started"

	activeTransaction = "active"
	idleTransaction   = "idle"

	virtualGuestAvailable    = "available"
	virtualGuestProvisioning = "provisioning"

	networkStorageMassAccessControlModificationException = "SoftLayer_Exception_Network_Storage_Group_MassAccessControlModification"
	retryDelayForModifyingStorageAccess                  = 10 * time.Second
)

func resourceIBMComputeVmInstance() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMComputeVmInstanceCreate,
		Read:     resourceIBMComputeVmInstanceRead,
		Update:   resourceIBMComputeVmInstanceUpdate,
		Delete:   resourceIBMComputeVmInstanceDelete,
		Exists:   resourceIBMComputeVmInstanceExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"bulk_vms"},
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
			},

			"domain": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"bulk_vms"},
			},

			"bulk_vms": {
				Type:          schema.TypeSet,
				Optional:      true,
				ForceNew:      true,
				MinItems:      2,
				ConflictsWith: []string{"hostname", "domain"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostname": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
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
						},

						"domain": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
				Set: resourceIBMBulkVMHostHash,
			},

			"os_reference_code": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, o, n string, d *schema.ResourceData) bool {
					if strings.HasSuffix(n, "_LATEST") {
						t := strings.Trim(n, "_LATEST")
						if strings.Contains(o, t) {
							return true
						}
					}
					return o == n
				},
				ConflictsWith: []string{"image_id"},
			},

			"hourly_billing": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},

			"private_network_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"datacenter": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"datacenter_choice"},
			},

			"datacenter_choice": {
				Type:          schema.TypeList,
				Description:   "The user provided datacenter options",
				Optional:      true,
				ConflictsWith: []string{"datacenter", "public_vlan_id", "private_vlan_id", "placement_group_name", "placement_group_id"},
				Elem:          &schema.Schema{Type: schema.TypeMap},
			},

			"placement_group_name": {
				Type:          schema.TypeString,
				Description:   "The placement group name",
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"datacenter_choice", "dedicated_acct_host_only", "dedicated_host_name", "dedicated_host_id", "placement_group_id"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					_, ok := d.GetOk("placement_group_id")
					return new == "" && ok
				},
			},

			"placement_group_id": {
				Type:          schema.TypeInt,
				Description:   "The placement group id",
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"datacenter_choice", "dedicated_acct_host_only", "dedicated_host_name", "dedicated_host_id", "placement_group_name"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					_, ok := d.GetOk("placement_group_name")
					return new == "0" && ok
				},
			},

			"flavor_key_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Description:   "Flavor key name used to provision vm.",
				ConflictsWith: []string{"cores", "memory"},
			},

			"cores": {

				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"flavor_key_name"},
			},

			"memory": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					memoryInMB := float64(v.(int))

					// Validate memory to match gigs format
					remaining := math.Mod(memoryInMB, 1024)
					if remaining > 0 {
						suggested := math.Ceil(memoryInMB/1024) * 1024
						errors = append(errors, fmt.Errorf(
							"Invalid 'memory' value %d megabytes, must be a multiple of 1024 (e.g. use %d)", int(memoryInMB), int(suggested)))
					}

					return
				},
				ConflictsWith: []string{"flavor_key_name"},
			},

			"dedicated_acct_host_only": {
				Type:          schema.TypeBool,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"dedicated_host_name", "dedicated_host_id", "placement_group_id", "placement_group_name"},
			},

			"dedicated_host_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"dedicated_acct_host_only", "dedicated_host_id", "placement_group_name", "placement_group_id"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					_, ok := d.GetOk("dedicated_host_id")
					return new == "" && ok
				},
			},

			"dedicated_host_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"dedicated_acct_host_only", "dedicated_host_name", "placement_group_name", "placement_group_id"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					_, ok := d.GetOk("dedicated_host_name")
					return new == "0" && ok
				},
			},

			"transient": {
				Type:          schema.TypeBool,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"dedicated_acct_host_only", "dedicated_host_name", "dedicated_host_id", "cores", "memory", "public_bandwidth_limited", "public_bandwidth_unlimited"},
			},

			"public_vlan_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"datacenter_choice"},
			},
			"public_interface_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"public_subnet": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"public_subnet_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"public_security_group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Set: func(v interface{}) int {
					return v.(int)
				},
				ForceNew: true,
				MaxItems: 5,
			},

			"private_vlan_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"datacenter_choice"},
			},
			"private_interface_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"private_subnet": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"private_subnet_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"private_security_group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Set: func(v interface{}) int {
					return v.(int)
				},
				ForceNew: true,
				MaxItems: 5,
			},

			"disks": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},

			"network_speed": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  100,
			},

			"ipv4_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ipv4_address_private": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ip_address_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"ip_address_id_private": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"ipv6_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"ipv6_static_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"ipv6_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ipv6_address_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			// SoftLayer does not support public_ipv6_subnet configuration in vm creation. So, public_ipv6_subnet
			// is defined as a computed parameter.
			"public_ipv6_subnet": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"public_ipv6_subnet_id": {
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
			},

			"secondary_ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"ssh_key_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Set: func(v interface{}) int {
					return v.(int)
				},
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
			"user_metadata": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"notes": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNotes,
			},

			"local_disk": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},

			"post_install_script_uri": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
				ForceNew: true,
			},

			"image_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"os_reference_code"},
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"wait_time_minutes": {
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "This field is deprecated. Use timeouts block instead",
				Default:    90,
			},
			// Monthly only
			// Limited BandWidth
			"public_bandwidth_limited": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				DiffSuppressFunc: applyOnce,
				ConflictsWith:    []string{"private_network_only", "public_bandwidth_unlimited"},
				ValidateFunc:     validatePublicBandwidth,
			},

			// Monthly only
			// Unlimited BandWidth
			"public_bandwidth_unlimited": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          false,
				ForceNew:         true,
				DiffSuppressFunc: applyOnce,
				ConflictsWith:    []string{"private_network_only", "public_bandwidth_limited"},
			},

			"evault": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: applyOnce,
			},

			// Quote based provisioning only
			"quote_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Quote ID for Quote based provisioning",
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
			ResourceStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the resource",
			},
		},
	}
}

type vmMember map[string]interface{}

func getSubnetID(subnet string, meta interface{}) (int, error) {
	service := services.GetAccountService(meta.(ClientSession).SoftLayerSession())

	subnetInfo := strings.Split(subnet, "/")
	if len(subnetInfo) != 2 {
		return 0, fmt.Errorf(
			"Unable to parse the provided subnet: %s", subnet)
	}

	networkIdentifier := subnetInfo[0]
	cidr := subnetInfo[1]

	subnets, err := service.
		Mask("id").
		Filter(
			filter.Build(
				filter.Path("subnets.cidr").Eq(cidr),
				filter.Path("subnets.networkIdentifier").Eq(networkIdentifier),
			),
		).
		GetSubnets()

	if err != nil {
		return 0, fmt.Errorf("Error looking up Subnet: %s", err)
	}

	if len(subnets) < 1 {
		return 0, fmt.Errorf(
			"Unable to locate a subnet matching the provided subnet: %s", subnet)
	}

	return *subnets[0].Id, nil
}

func resourceIBMBulkVMHostHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-",
		m["hostname"].(string)))

	return hashcode.String(buf.String())
}

func getNameForBlockDevice(i int) string {
	// skip 1, which is reserved for the swap disk.
	// so we get 0, 2, 3, 4, 5 ...
	if i == 0 {
		return "0"
	}

	return strconv.Itoa(i + 1)
}

func getNameForBlockDeviceWithFlavor(i int) string {
	// skip 0, which is taken from flavor.
	// skip 1, which is reserved for the swap disk.
	// so we get  2, 3, 4, 5 ...

	return strconv.Itoa(i + 2)
}

func getBlockDevices(d *schema.ResourceData) []datatypes.Virtual_Guest_Block_Device {
	numBlocks := d.Get("disks.#").(int)
	if numBlocks == 0 {
		return nil
	}
	blocks := make([]datatypes.Virtual_Guest_Block_Device, 0, numBlocks)
	for i := 0; i < numBlocks; i++ {
		var name string
		blockRef := fmt.Sprintf("disks.%d", i)
		if _, ok := d.GetOk("flavor_key_name"); ok {
			name = getNameForBlockDeviceWithFlavor(i)
		} else {
			name = getNameForBlockDevice(i)
		}
		capacity := d.Get(blockRef).(int)
		block := datatypes.Virtual_Guest_Block_Device{
			Device: &name,
			DiskImage: &datatypes.Virtual_Disk_Image{
				Capacity: &capacity,
			},
		}
		blocks = append(blocks, block)
	}

	return blocks
}

func expandSecurityGroupBindings(securityGroupsList []interface{}) ([]datatypes.Virtual_Network_SecurityGroup_NetworkComponentBinding, error) {
	if len(securityGroupsList) == 0 {
		return nil, nil
	}
	sgBindings := make([]datatypes.Virtual_Network_SecurityGroup_NetworkComponentBinding,
		len(securityGroupsList))
	for i, v := range securityGroupsList {
		sgid := v.(int)
		sgBindings[i] = datatypes.Virtual_Network_SecurityGroup_NetworkComponentBinding{
			SecurityGroupId: sl.Int(sgid),
		}
	}
	return sgBindings, nil
}

func getVirtualGuestTemplateFromResourceData(d *schema.ResourceData, meta interface{}, datacenter string, publicVlanID, privateVlanID, quote_id int) ([]datatypes.Virtual_Guest, error) {

	dc := datatypes.Location{
		Name: sl.String(datacenter),
	}
	// FIXME: Work around bug in terraform (?)
	// For properties that have a default value set and a diff suppress function,
	// it is not using the default value.
	networkSpeed := d.Get("network_speed").(int)
	if networkSpeed == 0 {
		networkSpeed = resourceIBMComputeVmInstance().Schema["network_speed"].Default.(int)
	}

	networkComponent := datatypes.Virtual_Guest_Network_Component{
		MaxSpeed: &networkSpeed,
	}
	members := []vmMember{}
	bulkVMs := d.Get("bulk_vms").(*schema.Set).List()
	if len(bulkVMs) > 0 {
		for _, vm := range bulkVMs {
			members = append(members, vm.(map[string]interface{}))
		}
	} else {
		member := vmMember{
			"hostname": d.Get("hostname").(string),
			"domain":   d.Get("domain").(string),
		}
		members = append(members, member)
	}

	vms := make([]datatypes.Virtual_Guest, 0)
	for _, member := range members {
		opts := datatypes.Virtual_Guest{
			Hostname:               sl.String(member["hostname"].(string)),
			Domain:                 sl.String(member["domain"].(string)),
			HourlyBillingFlag:      sl.Bool(d.Get("hourly_billing").(bool)),
			PrivateNetworkOnlyFlag: sl.Bool(d.Get("private_network_only").(bool)),
			Datacenter:             &dc,
			NetworkComponents:      []datatypes.Virtual_Guest_Network_Component{networkComponent},
			BlockDevices:           getBlockDevices(d),
			LocalDiskFlag:          sl.Bool(d.Get("local_disk").(bool)),
			PostInstallScriptUri:   sl.String(d.Get("post_install_script_uri").(string)),
		}

		if placementGroupID, ok := d.GetOk("placement_group_id"); ok {
			grpID := placementGroupID.(int)
			service := services.GetVirtualPlacementGroupService(meta.(ClientSession).SoftLayerSession())
			grp, err := service.Id(grpID).Mask("id,name,backendRouter[datacenter[name]]").GetObject()
			if err != nil {
				return vms, fmt.Errorf("Error looking up placement group: %s", err)
			}

			opts.PlacementGroupId = sl.Int(*grp.Id)

		} else if placementGroupName, ok := d.GetOk("placement_group_name"); ok {
			grpName := placementGroupName.(string)
			service := services.GetAccountService(meta.(ClientSession).SoftLayerSession())
			groups, err := service.
				Mask("id,name,backendRouter[hostname,datacenter[name]]").
				Filter(filter.Path("placementGroup.name").Eq(grpName).Build()).
				GetPlacementGroups()

			if err != nil {
				return vms, fmt.Errorf("Error looking up placement group '%s': %s", grpName, err)
			}
			grps := []datatypes.Virtual_PlacementGroup{}
			for _, g := range groups {
				if grpName == *g.Name {
					grps = append(grps, g)

				}
			}
			if len(grps) == 0 {
				return vms, fmt.Errorf("Error looking up placement group '%s'", grpName)
			}
			grp := grps[0]

			opts.PlacementGroupId = sl.Int(*grp.Id)
		}

		if startCPUs, ok := d.GetOk("cores"); ok {
			opts.StartCpus = sl.Int(startCPUs.(int))
		}
		if maxMemory, ok := d.GetOk("memory"); ok {
			opts.MaxMemory = sl.Int(maxMemory.(int))
		}

		if flavor, ok := d.GetOk("flavor_key_name"); ok {
			flavorComponenet := datatypes.Virtual_Guest_SupplementalCreateObjectOptions{
				FlavorKeyName: sl.String(flavor.(string)),
			}
			opts.SupplementalCreateObjectOptions = &flavorComponenet
		}

		if dedicatedAcctHostOnly, ok := d.GetOk("dedicated_acct_host_only"); ok {
			opts.DedicatedAccountHostOnlyFlag = sl.Bool(dedicatedAcctHostOnly.(bool))
		} else if dedicatedHostID, ok := d.GetOk("dedicated_host_id"); ok {
			opts.DedicatedHost = &datatypes.Virtual_DedicatedHost{
				Id: sl.Int(dedicatedHostID.(int)),
			}
		} else if dedicatedHostName, ok := d.GetOk("dedicated_host_name"); ok {
			hostName := dedicatedHostName.(string)
			service := services.GetAccountService(meta.(ClientSession).SoftLayerSession())
			hosts, err := service.
				Mask("id").
				Filter(filter.Path("dedicatedHosts.name").Eq(hostName).Build()).
				GetDedicatedHosts()

			if err != nil {
				return vms, fmt.Errorf("Error looking up dedicated host '%s': %s", hostName, err)
			} else if len(hosts) == 0 {
				return vms, fmt.Errorf("Error looking up dedicated host '%s'", hostName)
			}

			opts.DedicatedHost = &hosts[0]
		}

		if transientFlag, ok := d.GetOk("transient"); ok {
			if !*opts.HourlyBillingFlag || *opts.LocalDiskFlag {
				return vms, fmt.Errorf("Unable to provision a transient instance with a hourly_billing false or local_disk true")
			}
			opts.TransientGuestFlag = sl.Bool(transientFlag.(bool))
		}

		if quote_id == 0 {

			if imgID, ok := d.GetOk("image_id"); ok {
				imageID := imgID.(int)
				service := services.
					GetVirtualGuestBlockDeviceTemplateGroupService(meta.(ClientSession).SoftLayerSession())

				image, err := service.
					Mask("id,globalIdentifier").Id(imageID).
					GetObject()
				if err != nil {
					return vms, fmt.Errorf("Error looking up image %d: %s", imageID, err)
				} else if image.GlobalIdentifier == nil {
					return vms, fmt.Errorf(
						"Image template %d does not have a global identifier", imageID)
				}

				opts.BlockDeviceTemplateGroup = &datatypes.Virtual_Guest_Block_Device_Template_Group{
					GlobalIdentifier: image.GlobalIdentifier,
				}
			}

		}

		if operatingSystemReferenceCode, ok := d.GetOk("os_reference_code"); ok {
			opts.OperatingSystemReferenceCode = sl.String(operatingSystemReferenceCode.(string))
		}

		publicSubnet := d.Get("public_subnet").(string)
		privateSubnet := d.Get("private_subnet").(string)

		primaryNetworkComponent := datatypes.Virtual_Guest_Network_Component{
			NetworkVlan: &datatypes.Network_Vlan{},
		}

		usePrimaryNetworkComponent := false

		if publicVlanID > 0 {
			primaryNetworkComponent.NetworkVlan.Id = &publicVlanID
			usePrimaryNetworkComponent = true
		}

		// Apply public subnet if provided
		if publicSubnet != "" {
			primarySubnetID, err := getSubnetID(publicSubnet, meta)
			if err != nil {
				return vms, fmt.Errorf("Error creating virtual guest: %s", err)
			}
			primaryNetworkComponent.NetworkVlan.PrimarySubnetId = &primarySubnetID
			usePrimaryNetworkComponent = true
		}

		// Apply security groups if provided
		publicSecurityGroupIDList := d.Get("public_security_group_ids").(*schema.Set).List()
		sgb, err := expandSecurityGroupBindings(publicSecurityGroupIDList)
		if err != nil {
			return vms, err
		}
		if sgb != nil {
			primaryNetworkComponent.SecurityGroupBindings = sgb
			usePrimaryNetworkComponent = true
		}

		if usePrimaryNetworkComponent {
			opts.PrimaryNetworkComponent = &primaryNetworkComponent
		}

		primaryBackendNetworkComponent := datatypes.Virtual_Guest_Network_Component{
			NetworkVlan: &datatypes.Network_Vlan{},
		}

		usePrimaryBackendNetworkComponent := false

		if privateVlanID > 0 {
			primaryBackendNetworkComponent.NetworkVlan.Id = &privateVlanID
			usePrimaryBackendNetworkComponent = true
		}

		// Apply private subnet if provided
		if privateSubnet != "" {
			primarySubnetID, err := getSubnetID(privateSubnet, meta)
			if err != nil {
				return vms, fmt.Errorf("Error creating virtual guest: %s", err)
			}
			primaryBackendNetworkComponent.NetworkVlan.PrimarySubnetId = &primarySubnetID
			usePrimaryBackendNetworkComponent = true
		}

		// Apply security groups if provided
		privateSecurityGroupIDList := d.Get("private_security_group_ids").(*schema.Set).List()
		sgb, err = expandSecurityGroupBindings(privateSecurityGroupIDList)
		if err != nil {
			return vms, err
		}
		if sgb != nil {
			primaryBackendNetworkComponent.SecurityGroupBindings = sgb
			usePrimaryBackendNetworkComponent = true
		}

		if usePrimaryBackendNetworkComponent {
			opts.PrimaryBackendNetworkComponent = &primaryBackendNetworkComponent
		}

		if userData, ok := d.GetOk("user_metadata"); ok {
			opts.UserData = []datatypes.Virtual_Guest_Attribute{
				{
					Value: sl.String(userData.(string)),
				},
			}
		}

		if quote_id == 0 {

			// Get configured ssh_keys
			sshKeySet := d.Get("ssh_key_ids").(*schema.Set)
			sshKeys := sshKeySet.List()
			sshKeyLen := len(sshKeys)
			if sshKeyLen > 0 {
				opts.SshKeys = make([]datatypes.Security_Ssh_Key, 0, sshKeyLen)
				for _, sshKey := range sshKeys {
					opts.SshKeys = append(opts.SshKeys, datatypes.Security_Ssh_Key{
						Id: sl.Int(sshKey.(int)),
					})
				}
			}
		}

		vms = append(vms, opts)
	}

	return vms, nil
}

func resourceIBMComputeVmInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetVirtualGuestService(sess)

	var id int
	var receipt datatypes.Container_Product_Order_Receipt

	var err1 error
	var err error

	var dcName string
	var retryOptions []interface{}
	if dc, ok := d.GetOk("datacenter"); ok {
		dcName = dc.(string)
	}

	quote_id := d.Get("quote_id").(int)

	if options, ok := d.GetOk("datacenter_choice"); ok {
		retryOptions = options.([]interface{})
	}

	if dcName == "" && len(retryOptions) == 0 {
		return fmt.Errorf("Provide either `datacenter` or `datacenter_choice`")
	}

	if (d.Get("hostname").(string) == "" || d.Get("domain").(string) == "") && len(d.Get("bulk_vms").(*schema.Set).List()) == 0 {
		return fmt.Errorf("Provide either `hostname` and `domain` or `bulk_vms`")
	}

	if dcName != "" {
		publicVlan := 0
		privateVlan := 0
		if v, ok := d.GetOk("public_vlan_id"); ok {
			publicVlan = v.(int)

		}
		if v, ok := d.GetOk("private_vlan_id"); ok {
			privateVlan = v.(int)

		}

		receipt, err1 = placeOrder(d, meta, dcName, publicVlan, privateVlan, quote_id)
	} else if len(retryOptions) > 0 {

		err := validateDatacenterOption(retryOptions, []string{"datacenter", "public_vlan_id", "private_vlan_id"})
		if err != nil {
			return err
		}
		for _, option := range retryOptions {
			if option == nil {
				return fmt.Errorf("Provide a valid `datacenter_choice`")
			}
			center := option.(map[string]interface{})
			var publicVlan, privateVlan int
			var name string

			if v, ok := center["datacenter"]; ok {
				name = v.(string)
			} else {
				return fmt.Errorf("Missing datacenter in `datacenter_choice`")
			}

			if v, ok := center["public_vlan_id"]; ok {
				publicVlan, _ = strconv.Atoi(v.(string))
			}
			if v, ok := center["private_vlan_id"]; ok {
				privateVlan, _ = strconv.Atoi(v.(string))
			}

			receipt, err1 = placeOrder(d, meta, name, publicVlan, privateVlan, quote_id)
			if err1 == nil {
				break

			}
		}
	}

	if err1 != nil {
		return fmt.Errorf("Error ordering virtual guest: %s", err1)
	}

	var idStrings []string
	if quote_id > 0 {
		vmId := fmt.Sprintf("%d", *receipt.OrderDetails.VirtualGuests[0].Id)
		idStrings = append(idStrings, vmId)
		d.SetId(vmId)

	} else if len(receipt.OrderDetails.OrderContainers) > 1 {
		for i := 0; i < len(receipt.OrderDetails.OrderContainers); i++ {
			idStrings = append(idStrings, fmt.Sprintf("%d", *receipt.OrderDetails.OrderContainers[i].VirtualGuests[0].Id))
		}
		d.SetId(strings.Join(idStrings, "/"))
	} else {
		vmId := fmt.Sprintf("%d", *receipt.OrderDetails.OrderContainers[0].VirtualGuests[0].Id)
		idStrings = append(idStrings, vmId)
		d.SetId(vmId)
	}
	log.Printf("[INFO] Virtual Machine ID: %s", d.Id())
	for _, str := range idStrings {
		id, err = strconv.Atoi(str)
		if err != nil {
			return err
		}
		// Set tags
		tags := getTags(d)
		if tags != "" {
			//Try setting only when it is non empty as we are creating virtual guest
			err = setGuestTags(id, tags, meta)
			if err != nil {
				return err
			}
		}

		var storageIds []int
		if fileStorageSet := d.Get("file_storage_ids").(*schema.Set); len(fileStorageSet.List()) > 0 {
			storageIds = expandIntList(fileStorageSet.List())

		}
		if blockStorageSet := d.Get("block_storage_ids").(*schema.Set); len(blockStorageSet.List()) > 0 {
			storageIds = append(storageIds, expandIntList(blockStorageSet.List())...)
		}
		if len(storageIds) > 0 {
			err := addAccessToStorageList(service.Id(id), id, storageIds, meta)
			if err != nil {
				return err
			}
		}

		// Set notes
		err = setNotes(id, d, meta)
		if err != nil {
			return err
		}

		// wait for machine availability

		_, err = WaitForVirtualGuestAvailable(id, d, meta)

		if err != nil {
			return fmt.Errorf(
				"Error waiting for virtual machine (%s) to become ready: %s", d.Id(), err)
		}
	}

	return resourceIBMComputeVmInstanceRead(d, meta)
}

func resourceIBMComputeVmInstanceRead(d *schema.ResourceData, meta interface{}) error {
	service := services.GetVirtualGuestService(meta.(ClientSession).SoftLayerSession())
	parts, err := vmIdParts(d.Id())
	if err != nil {
		return err
	}
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	result, err := service.Id(id).Mask(
		"hostname,domain,blockDevices[diskImage],startCpus,maxMemory,dedicatedAccountHostOnlyFlag,operatingSystemReferenceCode,blockDeviceTemplateGroup[id],transientGuestFlag," +
			"billingItem[orderItem[preset[keyName]]]," +
			"primaryIpAddress,primaryBackendIpAddress,privateNetworkOnlyFlag," +
			"hourlyBillingFlag,localDiskFlag," +
			"allowedNetworkStorage[id,nasType]," +
			"notes,userData[value],tagReferences[id,tag[name]]," +
			"datacenter[id,name,longName]," +
			"sshKeys,status[keyName,name]," +
			"primaryNetworkComponent[networkVlan[id],subnets," +
			"primaryVersion6IpAddressRecord[subnet,guestNetworkComponentBinding[ipAddressId]]," +
			"primaryIpAddressRecord[subnet,guestNetworkComponentBinding[ipAddressId]]," +
			"securityGroupBindings[securityGroup]]," +
			"primaryBackendNetworkComponent[networkVlan[id]," +
			"primaryIpAddressRecord[subnet,guestNetworkComponentBinding[ipAddressId]]," +
			"securityGroupBindings[securityGroup]],evaultNetworkStorage[capacityGb]",
	).GetObject()

	if err != nil {
		return fmt.Errorf("Error retrieving virtual guest: %s", err)
	}

	if len(parts) == 1 {
		d.Set("hostname", *result.Hostname)
		d.Set("domain", *result.Domain)
	} else {
		members := make([]vmMember, 0)
		for _, part := range parts {
			vmId, err := strconv.Atoi(part)
			if err != nil {
				return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
			}
			vmResult, err := service.Id(vmId).Mask(
				"hostname,domain",
			).GetObject()
			member := vmMember{
				"hostname": *vmResult.Hostname,
				"domain":   *vmResult.Domain,
			}
			members = append(members, member)
		}
	}

	keyName, ok := sl.GrabOk(result, "BillingItem.OrderItem.Preset.KeyName")
	if ok {
		d.Set("flavor_key_name", keyName)
	}

	if result.BlockDeviceTemplateGroup != nil {
		d.Set("image_id", result.BlockDeviceTemplateGroup.Id)
	} else {
		if result.OperatingSystemReferenceCode != nil {
			d.Set("os_reference_code", result.OperatingSystemReferenceCode)
		}
	}

	if result.Datacenter != nil {
		d.Set("datacenter", *result.Datacenter.Name)
	}

	if result.DedicatedHost != nil {
		d.Set("dedicated_host_id", *result.DedicatedHost.Id)
		d.Set("dedicated_host_name", *result.DedicatedHost.Name)
	}

	if result.PlacementGroup != nil {
		d.Set("placement_group_id", *result.PlacementGroup.Id)
		d.Set("placement_group_name", *result.PlacementGroup.Name)
	}

	d.Set(
		"network_speed",
		sl.Grab(
			result,
			"PrimaryBackendNetworkComponent.MaxSpeed",
			d.Get("network_speed").(int),
		),
	)
	if result.OperatingSystemReferenceCode != nil && strings.HasPrefix(*result.OperatingSystemReferenceCode, "WIN") {
		d.Set("disks", flattenDisksForWindows(result))
	} else {
		d.Set("disks", flattenDisks(result))
	}
	d.Set("cores", *result.StartCpus)
	d.Set("memory", *result.MaxMemory)
	d.Set("dedicated_acct_host_only", *result.DedicatedAccountHostOnlyFlag)
	d.Set("transient", *result.TransientGuestFlag)
	d.Set("ipv4_address", result.PrimaryIpAddress)
	d.Set("ipv4_address_private", result.PrimaryBackendIpAddress)
	if result.PrimaryNetworkComponent != nil && result.PrimaryNetworkComponent.PrimaryIpAddressRecord != nil {
		d.Set("ip_address_id", *result.PrimaryNetworkComponent.PrimaryIpAddressRecord.GuestNetworkComponentBinding.IpAddressId)
	}
	if result.PrimaryNetworkComponent != nil {
		d.Set("public_interface_id", result.PrimaryNetworkComponent.Id)
	}
	if result.PrimaryBackendNetworkComponent != nil && result.PrimaryBackendNetworkComponent.PrimaryIpAddressRecord != nil {
		d.Set("ip_address_id_private",
			*result.PrimaryBackendNetworkComponent.PrimaryIpAddressRecord.GuestNetworkComponentBinding.IpAddressId)
	}
	if result.PrimaryBackendNetworkComponent != nil {
		d.Set("private_interface_id", result.PrimaryBackendNetworkComponent.Id)
	}
	d.Set("private_network_only", *result.PrivateNetworkOnlyFlag)
	d.Set("hourly_billing", *result.HourlyBillingFlag)
	d.Set("local_disk", *result.LocalDiskFlag)

	if result.PrimaryNetworkComponent != nil && result.PrimaryNetworkComponent.NetworkVlan != nil {
		d.Set("public_vlan_id", *result.PrimaryNetworkComponent.NetworkVlan.Id)
	}

	if result.PrimaryBackendNetworkComponent != nil && result.PrimaryBackendNetworkComponent.NetworkVlan != nil {
		d.Set("private_vlan_id", *result.PrimaryBackendNetworkComponent.NetworkVlan.Id)
	}

	if result.PrimaryNetworkComponent != nil && result.PrimaryNetworkComponent.PrimaryIpAddressRecord != nil {
		publicSubnet := result.PrimaryNetworkComponent.PrimaryIpAddressRecord.Subnet
		d.Set(
			"public_subnet",
			fmt.Sprintf("%s/%d", *publicSubnet.NetworkIdentifier, *publicSubnet.Cidr),
		)
		d.Set("public_subnet_id", result.PrimaryNetworkComponent.PrimaryIpAddressRecord.SubnetId)
	}

	if result.PrimaryNetworkComponent != nil && result.PrimaryNetworkComponent.SecurityGroupBindings != nil {
		var sgs []int
		for _, sg := range result.PrimaryNetworkComponent.SecurityGroupBindings {
			sgs = append(sgs, *sg.SecurityGroup.Id)
		}
		d.Set("public_security_group_ids", sgs)
	}

	if result.PrimaryBackendNetworkComponent != nil && result.PrimaryBackendNetworkComponent.PrimaryIpAddressRecord != nil {
		privateSubnet := result.PrimaryBackendNetworkComponent.PrimaryIpAddressRecord.Subnet
		d.Set(
			"private_subnet",
			fmt.Sprintf("%s/%d", *privateSubnet.NetworkIdentifier, *privateSubnet.Cidr),
		)
		d.Set("private_subnet_id", result.PrimaryBackendNetworkComponent.PrimaryIpAddressRecord.SubnetId)

	}

	if result.PrimaryBackendNetworkComponent != nil && result.PrimaryBackendNetworkComponent.SecurityGroupBindings != nil {
		var sgs []int
		for _, sg := range result.PrimaryBackendNetworkComponent.SecurityGroupBindings {
			sgs = append(sgs, *sg.SecurityGroup.Id)
		}
		d.Set("private_security_group_ids", sgs)
	}

	d.Set("ipv6_enabled", false)
	d.Set("ipv6_static_enabled", false)
	if result.PrimaryNetworkComponent != nil && result.PrimaryNetworkComponent.PrimaryVersion6IpAddressRecord != nil {
		d.Set("ipv6_enabled", true)
		d.Set("ipv6_address", *result.PrimaryNetworkComponent.PrimaryVersion6IpAddressRecord.IpAddress)
		d.Set("ipv6_address_id", *result.PrimaryNetworkComponent.PrimaryVersion6IpAddressRecord.GuestNetworkComponentBinding.IpAddressId)
		publicSubnet := result.PrimaryNetworkComponent.PrimaryVersion6IpAddressRecord.Subnet
		log.Println("DUDE", *publicSubnet, result.PrimaryNetworkComponent.PrimaryVersion6IpAddressRecord.SubnetId)
		d.Set(
			"public_ipv6_subnet",
			fmt.Sprintf("%s/%d", *publicSubnet.NetworkIdentifier, *publicSubnet.Cidr),
		)
		d.Set("public_ipv6_subnet_id", result.PrimaryNetworkComponent.PrimaryVersion6IpAddressRecord.SubnetId)
	}
	if result.PrimaryNetworkComponent != nil {
		for _, subnet := range result.PrimaryNetworkComponent.Subnets {
			if *subnet.SubnetType == "STATIC_IP_ROUTED_6" {
				d.Set("ipv6_static_enabled", true)
			}
		}
	}

	userData := result.UserData

	if userData != nil && len(userData) > 0 {
		d.Set("user_metadata", userData[0].Value)
	}

	d.Set("notes", sl.Get(result.Notes, nil))

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
	d.Set("block_storage_ids", flattenBlockStorageID(storages))
	d.Set("file_storage_ids", flattenFileStorageID(storages))

	sshKeys := result.SshKeys
	if len(sshKeys) > 0 {
		d.Set("ssh_key_ids", flattenSSHKeyIDs(sshKeys))
	}

	// Set connection info
	connInfo := map[string]string{"type": "ssh"}
	if !*result.PrivateNetworkOnlyFlag && result.PrimaryIpAddress != nil {
		connInfo["host"] = *result.PrimaryIpAddress
	} else if result.PrimaryBackendIpAddress != nil {
		connInfo["host"] = *result.PrimaryBackendIpAddress
	}
	d.SetConnInfo(connInfo)

	if _, ok := sl.GrabOk(result, "EvaultNetworkStorage"); ok {
		if len(result.EvaultNetworkStorage) > 0 {
			d.Set("evault", result.EvaultNetworkStorage[0].CapacityGb)
		}

	}
	d.Set(ResourceControllerURL, fmt.Sprintf("https://cloud.ibm.com/gen1/infrastructure/virtual-server/%s/details#main", d.Id()))
	d.Set(ResourceName, *result.Hostname)
	d.Set(ResourceStatus, *result.Status.Name)
	err = readSecondaryIPAddresses(d, meta, result.PrimaryIpAddress)
	return err
}

func readSecondaryIPAddresses(d *schema.ResourceData, meta interface{}, primaryIPAddress *string) error {
	d.Set("secondary_ip_addresses", nil)
	if primaryIPAddress != nil {
		secondarySubnetResult, err := services.GetAccountService(meta.(ClientSession).SoftLayerSession()).
			Mask("ipAddresses[id,ipAddress],subnetType").
			Filter(filter.Build(filter.Path("publicSubnets.endPointIpAddress.ipAddress").Eq(*primaryIPAddress))).
			GetPublicSubnets()
		if err != nil {
			log.Printf("Error getting secondary Ip addresses: %s", err)
		}

		secondaryIps := make([]string, 0)
		for _, subnet := range secondarySubnetResult {
			// Count static secondary ip addresses.
			if *subnet.SubnetType == staticIPRouted {
				for _, ipAddressObj := range subnet.IpAddresses {
					secondaryIps = append(secondaryIps, *ipAddressObj.IpAddress)
				}
			}
		}
		if len(secondaryIps) > 0 {
			d.Set("secondary_ip_addresses", secondaryIps)
			d.Set("secondary_ip_count", len(secondaryIps))
		}
	}
	return nil
}
func resourceIBMComputeVmInstanceUpdate(d *schema.ResourceData, meta interface{}) error {

	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetVirtualGuestService(sess)

	parts, err := vmIdParts(d.Id())
	if err != nil {
		return err
	}
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	result, err := service.Id(id).GetObject()
	if err != nil {
		return fmt.Errorf("Error retrieving virtual guest: %s", err)
	}

	isChanged := false

	// Update "hostname" and "domain" fields if present and changed
	// Those are the only fields, which could be updated
	if d.HasChange("hostname") || d.HasChange("domain") {
		result.Hostname = sl.String(d.Get("hostname").(string))
		result.Domain = sl.String(d.Get("domain").(string))
		isChanged = true
	}

	if d.HasChange("notes") {
		result.Notes = sl.String(d.Get("notes").(string))
		isChanged = true
	}

	if isChanged {
		_, err = service.Id(id).EditObject(&result)
		if err != nil {
			return fmt.Errorf("Couldn't update virtual guest: %s", err)
		}
	}

	// Update tags
	if d.HasChange("tags") {
		tags := getTags(d)
		err := setGuestTags(id, tags, meta)
		if err != nil {
			return err
		}
	}

	err = modifyStorageAccess(service.Id(id), id, meta, d)
	if err != nil {
		return err
	}

	// Upgrade "cores", "memory" and "network_speed" if provided and changed
	upgradeOptions := map[string]float64{}
	if d.HasChange("cores") {
		upgradeOptions[product.CPUCategoryCode] = float64(d.Get("cores").(int))
	}

	if d.HasChange("memory") {
		memoryInMB := float64(d.Get("memory").(int))

		// Convert memory to GB, as softlayer only allows to upgrade RAM in Gigs
		// Must be already validated at this step
		upgradeOptions[product.MemoryCategoryCode] = float64(int(memoryInMB / 1024))
	}

	if d.HasChange("network_speed") {
		upgradeOptions[product.NICSpeedCategoryCode] = float64(d.Get("network_speed").(int))
	}

	if d.HasChange("disks") {
		oldDisks, newDisks := d.GetChange("disks")
		oldDisk := oldDisks.([]interface{})
		newDisk := newDisks.([]interface{})

		//Remove is not supported for now.
		if len(oldDisk) > len(newDisk) {
			return fmt.Errorf("Removing drives is not supported.")
		}

		var diskName string
		//Update the disks if any change
		for i := 0; i < len(oldDisk); i++ {
			if newDisk[i].(int) != oldDisk[i].(int) {

				if _, ok := d.GetOk("flavor_key_name"); ok {
					diskName = fmt.Sprintf("guest_disk%d", i+1)
				} else {
					diskName = fmt.Sprintf("guest_disk%d", i)
				}
				capacity := newDisk[i].(int)
				upgradeOptions[diskName] = float64(capacity)
			}
		}
		//Add new disks
		for i := len(oldDisk); i < len(newDisk); i++ {
			if _, ok := d.GetOk("flavor_key_name"); ok {
				diskName = fmt.Sprintf("guest_disk%d", i+1)
			} else {
				diskName = fmt.Sprintf("guest_disk%d", i)
			}
			capacity := newDisk[i].(int)
			upgradeOptions[diskName] = float64(capacity)
		}

	}

	if len(upgradeOptions) > 0 || d.HasChange("flavor_key_name") {

		if _, ok := d.GetOk("flavor_key_name"); ok {
			presetKeyName := d.Get("flavor_key_name").(string)
			_, err = virtual.UpgradeVirtualGuestWithPreset(sess.SetRetries(0), &result, presetKeyName, upgradeOptions)
			if err != nil {
				return fmt.Errorf("Couldn't upgrade virtual guest: %s", err)
			}

		} else {
			_, err = virtual.UpgradeVirtualGuest(sess.SetRetries(0), &result, upgradeOptions)
			if err != nil {
				return fmt.Errorf("Couldn't upgrade virtual guest: %s", err)
			}
		}

		// Wait for softlayer to start upgrading...
		_, err = WaitForUpgradeTransactionsToAppear(d, meta)
		if err != nil {
			return err
		}
		// Wait for upgrade transactions to finish
		_, err = WaitForNoActiveTransactions(id, d, d.Timeout(schema.TimeoutUpdate), meta)
		if err != nil {
			return err
		}

	}

	return resourceIBMComputeVmInstanceRead(d, meta)
}

func modifyStorageAccess(sam storageAccessModifier, deviceID int, meta interface{}, d *schema.ResourceData) error {
	var remove, add []int
	if d.HasChange("file_storage_ids") {
		o, n := d.GetChange("file_storage_ids")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		remove = expandIntList(os.Difference(ns).List())
		add = expandIntList(ns.Difference(os).List())
	}
	if d.HasChange("block_storage_ids") {
		o, n := d.GetChange("block_storage_ids")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		remove = append(remove, expandIntList(os.Difference(ns).List())...)
		add = append(add, expandIntList(ns.Difference(os).List())...)
	}

	if len(add) > 0 {
		err := addAccessToStorageList(sam, deviceID, add, meta)
		if err != nil {
			return err
		}
	}
	if len(remove) > 0 {
		err := removeAccessToStorageList(sam, deviceID, remove, meta)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceIBMComputeVmInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetVirtualGuestService(sess)
	parts, err := vmIdParts(d.Id())
	if err != nil {
		return err
	}
	for _, part := range parts {
		id, err := strconv.Atoi(part)
		if err != nil {
			return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
		}

		_, err = WaitForNoActiveTransactions(id, d, d.Timeout(schema.TimeoutDelete), meta)

		if err != nil {
			return fmt.Errorf("Error deleting virtual guest, couldn't wait for zero active transactions: %s", err)
		}
		err = detachSecurityGroupNetworkComponentBindings(d, meta, id)
		if err != nil {
			return err
		}
		ok, err := service.Id(id).DeleteObject()
		if err != nil {
			return fmt.Errorf("Error deleting virtual guest: %s", err)
		}

		if !ok {
			return fmt.Errorf(
				"API reported it was unsuccessful in removing the virtual guest '%d'", id)
		}
	}

	return nil
}

func detachSecurityGroupNetworkComponentBindings(d *schema.ResourceData, meta interface{}, id int) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetVirtualGuestService(sess)
	publicSgIDs := d.Get("public_security_group_ids").(*schema.Set).List()
	privateSgIDS := d.Get("private_security_group_ids").(*schema.Set).List()
	if len(publicSgIDs) == 0 && len(privateSgIDS) == 0 {
		log.Println("No security groups specified, hence no detachment required before delete operation")
		return nil
	}
	vsi, err := service.Id(id).Mask(
		"primaryNetworkComponent[id,securityGroupBindings[securityGroupId,networkComponentId]]," +
			"primaryBackendNetworkComponent[id,securityGroupBindings[securityGroupId,networkComponentId]]",
	).GetObject()

	if err != nil {
		return err
	}
	sgService := services.GetNetworkSecurityGroupService(sess)
	//Detach security group as destroy might fail if the security group is attempted
	//to be destroyed in the same terraform configuration file. VSI destroy takes
	//some time andif during the same time security group which was referred in the VSI
	//is attempted to be destroyed it will fail.
	for _, v := range publicSgIDs {
		sgID := v.(int)
		for _, v := range vsi.PrimaryNetworkComponent.SecurityGroupBindings {
			if sgID == *v.SecurityGroupId {
				_, err := sgService.Id(sgID).DetachNetworkComponents([]int{*v.NetworkComponentId})
				if err != nil {
					return err
				}
			}
		}
	}
	for _, v := range privateSgIDS {
		sgID := v.(int)
		for _, v := range vsi.PrimaryBackendNetworkComponent.SecurityGroupBindings {
			if sgID == *v.SecurityGroupId {
				_, err := sgService.Id(sgID).DetachNetworkComponents([]int{*v.NetworkComponentId})
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

//genID generates a random string to be used for the optional
//hostname
func genID() (interface{}, error) {
	numBytes := 8
	bytes := make([]byte, numBytes)
	n, err := rand.Reader.Read(bytes)
	if err != nil {
		return nil, err
	}

	if n != numBytes {
		return nil, errors.New("generated insufficient random bytes")
	}

	hexStr := hex.EncodeToString(bytes)
	return fmt.Sprintf("terraformed-%s", hexStr), nil
}

// WaitForUpgradeTransactionsToAppear Wait for upgrade transactions
func WaitForUpgradeTransactionsToAppear(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	log.Printf("Waiting for server (%s) to have upgrade transactions", d.Id())

	parts, err := vmIdParts(d.Id())
	if err != nil {
		return nil, err
	}
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("The instance ID %s must be numeric", d.Id())
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"retry", pendingUpgrade},
		Target:  []string{inProgressUpgrade},
		Refresh: func() (interface{}, string, error) {
			service := services.GetVirtualGuestService(meta.(ClientSession).SoftLayerSession())
			transactions, err := service.Id(id).GetActiveTransactions()
			if err != nil {
				if apiErr, ok := err.(sl.Error); ok && apiErr.StatusCode == 404 {
					return nil, "", fmt.Errorf("Couldn't fetch active transactions: %s", err)
				}
				return false, "retry", nil
			}
			for _, transaction := range transactions {
				if strings.Contains(*transaction.TransactionStatus.Name, upgradeTransaction) {
					return transactions, inProgressUpgrade, nil
				}
			}
			return transactions, pendingUpgrade, nil
		},
		Timeout:    10 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	return stateConf.WaitForState()
}

// WaitForNoActiveTransactions Wait for no active transactions
func WaitForNoActiveTransactions(id int, d *schema.ResourceData, timeout time.Duration, meta interface{}) (interface{}, error) {
	log.Printf("Waiting for server (%s) to have zero active transactions", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending: []string{"retry", activeTransaction},
		Target:  []string{idleTransaction},
		Refresh: func() (interface{}, string, error) {
			service := services.GetVirtualGuestService(meta.(ClientSession).SoftLayerSession())
			transactions, err := service.Id(id).GetActiveTransactions()
			if err != nil {
				if apiErr, ok := err.(sl.Error); ok && apiErr.StatusCode == 404 {
					return nil, "", nil
				}
				return false, "retry", fmt.Errorf("Couldn't get active transactions: %s", err)
			}
			if len(transactions) == 0 {
				return transactions, idleTransaction, nil
			}
			return transactions, activeTransaction, nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

// WaitForVirtualGuestAvailable Waits for virtual guest creation
func WaitForVirtualGuestAvailable(id int, d *schema.ResourceData, meta interface{}) (interface{}, error) {
	log.Printf("Waiting for server (%s) to be available.", d.Id())
	sess := meta.(ClientSession).SoftLayerSession()
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", virtualGuestProvisioning},
		Target:     []string{virtualGuestAvailable},
		Refresh:    virtualGuestStateRefreshFunc(sess, id, d),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func virtualGuestStateRefreshFunc(sess *session.Session, instanceID int, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		// Check active transactions
		publicNetwork := !d.Get("private_network_only").(bool)
		service := services.GetVirtualGuestService(sess)
		result, err := service.Id(instanceID).Mask("activeTransaction,primaryBackendIpAddress,primaryIpAddress").GetObject()
		if err != nil {
			if apiErr, ok := err.(sl.Error); ok && apiErr.StatusCode == 404 {
				return nil, "", fmt.Errorf("Error retrieving virtual guest: %s", err)
			}
			return false, "retry", nil
		}
		// Check active transactions
		log.Println("Checking active transactions.")
		if result.ActiveTransaction != nil {
			return result, virtualGuestProvisioning, nil
		}

		// Check Primary IP address availability.
		log.Println("Checking primary backend IP address.")
		if result.PrimaryBackendIpAddress == nil {
			return result, virtualGuestProvisioning, nil
		}

		log.Println("Checking primary IP address.")
		if publicNetwork && result.PrimaryIpAddress == nil {
			return result, virtualGuestProvisioning, nil
		}

		// Check Secondary IP address availability.
		if d.Get("secondary_ip_count").(int) > 0 {
			log.Println("Refreshing secondary IPs state.")
			secondarySubnetResult, err := services.GetAccountService(sess).
				Mask("ipAddresses[id,ipAddress]").
				Filter(filter.Build(filter.Path("publicSubnets.endPointIpAddress.virtualGuest.id").Eq(fmt.Sprintf("%d", instanceID)))).
				GetPublicSubnets()
			if err != nil {
				return nil, "", fmt.Errorf("Error retrieving secondary ip address: %s", err)
			}
			if len(secondarySubnetResult) == 0 {
				return result, virtualGuestProvisioning, nil
			}
		}

		return result, virtualGuestAvailable, nil
	}
}

func resourceIBMComputeVmInstanceExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	service := services.GetVirtualGuestService(meta.(ClientSession).SoftLayerSession())
	parts, err := vmIdParts(d.Id())
	if err != nil {
		return false, err
	}
	guestID, err := strconv.Atoi(parts[0])
	if err != nil {
		return false, fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	result, err := service.Id(guestID).GetObject()
	if err != nil {
		if apiErr, ok := err.(sl.Error); ok {
			if apiErr.StatusCode == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}

	return result.Id != nil && *result.Id == guestID, nil
}

func getTags(d dataRetriever) string {
	tagSet := d.Get("tags").(*schema.Set)

	if tagSet.Len() == 0 {
		return ""
	}

	tags := make([]string, 0, tagSet.Len())
	for _, elem := range tagSet.List() {
		tag := elem.(string)
		tags = append(tags, tag)
	}
	return strings.Join(tags, ",")
}

func setGuestTags(id int, tags string, meta interface{}) error {
	service := services.GetVirtualGuestService(meta.(ClientSession).SoftLayerSession())
	_, err := service.Id(id).SetTags(sl.String(tags))
	if err != nil {
		return fmt.Errorf("Could not set tags on virtual guest %d", id)
	}
	return nil
}

type storageAccessModifier interface {
	AllowAccessToNetworkStorageList([]datatypes.Network_Storage) (resp bool, err error)
	RemoveAccessToNetworkStorageList([]datatypes.Network_Storage) (resp bool, err error)
}

func addAccessToStorageList(sam storageAccessModifier, deviceID int, ids storageIds, meta interface{}) error {
	s, err := ids.Storages(meta)
	if err != nil {
		return err
	}
	for {
		_, err := sam.AllowAccessToNetworkStorageList(s)
		if err != nil {
			if apiErr, ok := err.(sl.Error); ok && apiErr.Exception == networkStorageMassAccessControlModificationException {
				log.Printf("[DEBUG]  Allow access to storage failed with error %q. Will retry again after %q", err, retryDelayForModifyingStorageAccess)
				time.Sleep(retryDelayForModifyingStorageAccess)
				continue
			}
			return fmt.Errorf("Could not authorize Device %d, access to the following storages %q, %q", deviceID, ids, err)
		}
		log.Printf("[INFO] Device authorized to access %q", ids)
		break
	}
	return nil
}

func removeAccessToStorageList(sam storageAccessModifier, deviceID int, ids storageIds, meta interface{}) error {
	s, err := ids.Storages(meta)
	if err != nil {
		return err
	}
	for {
		_, err := sam.RemoveAccessToNetworkStorageList(s)
		if err != nil {
			if apiErr, ok := err.(sl.Error); ok && apiErr.Exception == networkStorageMassAccessControlModificationException {
				log.Printf("[DEBUG]  Remove access to storage failed with error %q. Will retry again after %q", err, retryDelayForModifyingStorageAccess)
				time.Sleep(retryDelayForModifyingStorageAccess)
				continue
			}
			return fmt.Errorf("Could not remove Device %d, access to the following storages %q, %q", deviceID, ids, err)
		}
		log.Printf("[INFO] Devices's access to %q have been removed", ids)
		break
	}
	return nil
}

func setNotes(id int, d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetVirtualGuestService(sess)

	if notes := d.Get("notes").(string); notes != "" {
		result, err := service.Id(id).GetObject()
		if err != nil {
			return fmt.Errorf("Error retrieving virtual guest: %s", err)
		}

		result.Notes = sl.String(notes)

		_, err = service.Id(id).EditObject(&result)
		if err != nil {
			return fmt.Errorf("Could not set note on virtual guest %d", id)
		}
	}

	return nil
}

func placeOrder(d *schema.ResourceData, meta interface{}, name string, publicVlanID, privateVlanID, quote_id int) (datatypes.Container_Product_Order_Receipt, error) {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetVirtualGuestService(sess)

	options, err := getVirtualGuestTemplateFromResourceData(d, meta, name, publicVlanID, privateVlanID, quote_id)
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}
	guestOrders := make([]datatypes.Container_Product_Order, 0)
	var template datatypes.Container_Product_Order
	if quote_id > 0 {
		// Build a virtual instance template from the quote.
		template, err = services.GetBillingOrderQuoteService(sess).
			Id(quote_id).GetRecalculatedOrderContainer(nil, sl.Bool(false))
		if err != nil {
			return datatypes.Container_Product_Order_Receipt{}, fmt.Errorf(
				"Encountered problem trying to get the virtual machine order template from quote: %s", err)
		}
		template.Quantity = sl.Int(1)
		template.ComplexType = sl.String("SoftLayer_Container_Product_Order_Virtual_Guest")
		template.VirtualGuests = make([]datatypes.Virtual_Guest, 0, 1)
		template.VirtualGuests = append(
			template.VirtualGuests,
			options[0],
		)
		// Get configured ssh_keys
		sshKeySet := d.Get("ssh_key_ids").(*schema.Set)
		sshKeys := sshKeySet.List()
		sshKeyLen := len(sshKeys)
		if sshKeyLen > 0 {
			sshKeyA := make([]int, sshKeyLen)
			template.SshKeys = make([]datatypes.Container_Product_Order_SshKeys, 0, sshKeyLen)
			for i, sshKey := range sshKeys {
				sshKeyA[i] = sshKey.(int)

			}
			template.SshKeys = append(template.SshKeys, datatypes.Container_Product_Order_SshKeys{
				SshKeyIds: sshKeyA,
			})
		}
		if rawImageTemplateId, ok := d.GetOk("image_id"); ok {
			imageTemplateId := rawImageTemplateId.(int)
			template.ImageTemplateId = sl.Int(imageTemplateId)
		}

		if postInstallURI, ok := d.GetOk("post_install_script_uri"); ok {
			postInstallURIA := make([]string, 1)
			postInstallURIA[0] = postInstallURI.(string)
			template.ProvisionScripts = postInstallURIA
		}

		guestOrders = append(guestOrders, template)
		order := &datatypes.Container_Product_Order{
			OrderContainers: guestOrders,
		}
		receipt, err1 := services.GetBillingOrderQuoteService(sess).
			Id(quote_id).PlaceOrder(order)
		return receipt, err1
	}
	for i := 0; i < len(options); i++ {
		opts := options[i]

		log.Println("[INFO] Creating virtual machine")

		// Build an order template with a custom image.
		if opts.BlockDevices != nil && opts.BlockDeviceTemplateGroup != nil {
			bd := *opts.BlockDeviceTemplateGroup
			opts.BlockDeviceTemplateGroup = nil
			opts.OperatingSystemReferenceCode = sl.String("UBUNTU_LATEST")
			template, err = service.GenerateOrderTemplate(&opts)
			if err != nil {
				return datatypes.Container_Product_Order_Receipt{}, fmt.Errorf("Error generating order template: %s", err)
			}

			// Remove temporary OS from actual order
			prices := make([]datatypes.Product_Item_Price, len(template.Prices))
			i := 0
			for _, p := range template.Prices {
				if !strings.Contains(*p.Item.Description, "Ubuntu") {
					prices[i] = p
					i++
				}
			}
			template.Prices = prices[:i]

			template.ImageTemplateId = sl.Int(d.Get("image_id").(int))
			template.VirtualGuests[0].BlockDeviceTemplateGroup = &bd
			template.VirtualGuests[0].OperatingSystemReferenceCode = nil
		} else {
			// Build an order template with os_reference_code
			template, err = service.GenerateOrderTemplate(&opts)
			if err != nil {
				return datatypes.Container_Product_Order_Receipt{}, fmt.Errorf("Error generating order template: %s", err)
			}
		}

		items, err := product.GetPackageProducts(sess, *template.PackageId, productItemMaskWithPriceLocationGroupID)
		if err != nil {
			return datatypes.Container_Product_Order_Receipt{}, fmt.Errorf("Error generating order template: %s", err)
		}

		privateNetworkOnly := d.Get("private_network_only").(bool)

		secondaryIPCount := d.Get("secondary_ip_count").(int)
		if secondaryIPCount > 0 {
			if privateNetworkOnly {
				return datatypes.Container_Product_Order_Receipt{}, fmt.Errorf("Unable to configure public secondary addresses with a private_network_only option")
			}
			keyName := strconv.Itoa(secondaryIPCount) + "_PUBLIC_IP_ADDRESSES"
			price, err := getItemPriceId(items, "sec_ip_addresses", keyName)
			if err != nil {
				return datatypes.Container_Product_Order_Receipt{}, err
			}
			template.Prices = append(template.Prices, price)
		}

		if d.Get("ipv6_enabled").(bool) {
			if privateNetworkOnly {
				return datatypes.Container_Product_Order_Receipt{}, fmt.Errorf("Unable to configure a public IPv6 address with a private_network_only option")
			}
			price, err := getItemPriceId(items, "pri_ipv6_addresses", "1_IPV6_ADDRESS")
			if err != nil {
				return datatypes.Container_Product_Order_Receipt{}, fmt.Errorf("Error generating order template: %s", err)
			}
			template.Prices = append(template.Prices, price)
		}

		if d.Get("ipv6_static_enabled").(bool) {
			if privateNetworkOnly {
				return datatypes.Container_Product_Order_Receipt{}, fmt.Errorf("Unable to configure a public static IPv6 address with a private_network_only option")
			}
			price, err := getItemPriceId(items, "static_ipv6_addresses", "64_BLOCK_STATIC_PUBLIC_IPV6_ADDRESSES")
			if err != nil {
				return datatypes.Container_Product_Order_Receipt{}, fmt.Errorf("Error generating order template: %s", err)
			}
			template.Prices = append(template.Prices, price)
		}

		// Add optional price ids.
		// Add public bandwidth limited
		if publicBandwidth, ok := d.GetOk("public_bandwidth_limited"); ok {
			if *opts.HourlyBillingFlag {
				return datatypes.Container_Product_Order_Receipt{}, fmt.Errorf("Unable to configure a public bandwidth with a hourly_billing true")
			}
			// Remove Default bandwidth price
			prices := make([]datatypes.Product_Item_Price, len(template.Prices))
			i := 0
			for _, p := range template.Prices {
				item := p.Item
				if item != nil {
					if strings.Contains(*item.Description, "Bandwidth") {
						continue
					}
				}
				prices[i] = p
				i++
			}
			template.Prices = prices[:i]
			keyName := "BANDWIDTH_" + strconv.Itoa(publicBandwidth.(int)) + "_GB"
			price, err := getItemPriceId(items, "bandwidth", keyName)
			if err != nil {
				return datatypes.Container_Product_Order_Receipt{}, fmt.Errorf("Error generating order template: %s", err)
			}
			template.Prices = append(template.Prices, price)
		}

		// Add public bandwidth unlimited
		publicUnlimitedBandwidth := d.Get("public_bandwidth_unlimited").(bool)
		if publicUnlimitedBandwidth {
			if *opts.HourlyBillingFlag {
				return datatypes.Container_Product_Order_Receipt{}, fmt.Errorf("Unable to configure a public bandwidth with a hourly_billing true")
			}
			networkSpeed := d.Get("network_speed").(int)
			if networkSpeed != 100 {
				return datatypes.Container_Product_Order_Receipt{}, fmt.Errorf("Network speed must be 100 Mbps to configure public bandwidth unlimited")
			}
			// Remove Default bandwidth price
			prices := make([]datatypes.Product_Item_Price, len(template.Prices))
			i := 0
			for _, p := range template.Prices {
				item := p.Item
				if item != nil {
					if strings.Contains(*item.Description, "Bandwidth") {
						continue
					}
				}
				prices[i] = p
				i++
			}
			template.Prices = prices[:i]
			price, err := getItemPriceId(items, "bandwidth", "BANDWIDTH_UNLIMITED_100_MBPS_UPLINK")
			if err != nil {
				return datatypes.Container_Product_Order_Receipt{}, fmt.Errorf("Error generating order template: %s", err)
			}
			template.Prices = append(template.Prices, price)
		}

		if evault, ok := d.GetOk("evault"); ok {
			if *opts.HourlyBillingFlag {
				return datatypes.Container_Product_Order_Receipt{}, fmt.Errorf("Unable to configure a evault with hourly_billing true")
			}

			keyName := "EVAULT_" + strconv.Itoa(evault.(int)) + "_GB"
			price, err := getItemPriceId(items, "evault", keyName)
			if err != nil {
				return datatypes.Container_Product_Order_Receipt{}, fmt.Errorf("Error generating order template: %s", err)
			}
			template.Prices = append(template.Prices, price)
		}
		// GenerateOrderTemplate omits UserData, subnet, and maxSpeed, so configure virtual_guest.
		template.VirtualGuests[0] = opts
		guestOrders = append(guestOrders, template)

	}
	order := &datatypes.Container_Product_Order{
		OrderContainers: guestOrders,
	}

	orderService := services.GetProductOrderService(sess.SetRetries(0))
	receipt, err1 := orderService.PlaceOrder(order, sl.Bool(false))
	return receipt, err1

}
