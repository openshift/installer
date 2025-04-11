// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func ResourceIBMIsBareMetalServerNetworkAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsBareMetalServerNetworkAttachmentCreate,
		ReadContext:   resourceIBMIsBareMetalServerNetworkAttachmentRead,
		UpdateContext: resourceIBMIsBareMetalServerNetworkAttachmentUpdate,
		DeleteContext: resourceIBMIsBareMetalServerNetworkAttachmentDelete,
		Importer:      &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bare_metal_server": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server_network_attachment", "bare_metal_server"),
				Description:  "The bare metal server identifier.",
			},
			"floating_bare_metal_server": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The bare metal server identifier of the server where the attachment is floated to(only applicated for allow_to_float true).",
			},

			"network_attachment": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The network attachment's id.",
			},
			"virtual_network_interface_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The virtual_network_interface's id.",
			},
			"interface_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server_network_attachment", "interface_type"),
				Description:  "The network attachment's interface type:- `pci`: a physical PCI device which can only be created or deleted when the bare metal  server is stopped  - Has an `allowed_vlans` property which controls the VLANs that will be permitted    to use the PCI attachment  - Cannot directly use an IEEE 802.1q VLAN tag.- `vlan`: a virtual device, used through a `pci` device that has the `vlan` in its  array of `allowed_vlans`.  - Must use an IEEE 802.1q tag.The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server_network_attachment", "name"),
				Description:  "The name for this bare metal server network attachment. The name is unique across all network attachments for the bare metal server.",
			},
			"virtual_network_interface": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Computed:    true,
				Description: "A virtual network interface for the bare metal server network attachment. This can be specified using an existing virtual network interface, or a prototype object for a new virtual network interface.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "The virtual network interface id for this bare metal server network attachment.",
						},
						"allow_ip_spoofing": &schema.Schema{
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"virtual_network_interface.0.id"},
							Computed:      true,
							Description:   "Indicates whether source IP spoofing is allowed on this interface. If `false`, source IP spoofing is prevented on this interface. If `true`, source IP spoofing is allowed on this interface.",
						},
						"auto_delete": &schema.Schema{
							Type:          schema.TypeBool,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"virtual_network_interface.0.id"},
							Description:   "Indicates whether this virtual network interface will be automatically deleted when`target` is deleted.",
						},
						"enable_infrastructure_nat": &schema.Schema{
							Type:          schema.TypeBool,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"virtual_network_interface.0.id"},
							Description:   "If `true`:- The VPC infrastructure performs any needed NAT operations.- `floating_ips` must not have more than one floating IP.If `false`:- Packets are passed unchanged to/from the network interface,  allowing the workload to perform any needed NAT operations.- `allow_ip_spoofing` must be `false`.- If the virtual network interface is attached:  - The target `resource_type` must be `bare_metal_server_network_attachment`.  - The target `interface_type` must not be `hipersocket`.",
						},
						"ips": &schema.Schema{
							Type:          schema.TypeSet,
							Optional:      true,
							Computed:      true,
							Set:           hashIpsList,
							ConflictsWith: []string{"virtual_network_interface.0.id"},
							Description:   "The reserved IPs bound to this virtual network interface.May be empty when `lifecycle_state` is `pending`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": &schema.Schema{
										Type: schema.TypeString,
										// Optional:    true,
										Computed:    true,
										Description: "The IP address.If the address has not yet been selected, the value will be `0.0.0.0`.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.",
									},
									"deleted": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									"auto_delete": &schema.Schema{
										Type: schema.TypeBool,
										// Optional:    true,
										Computed:    true,
										Description: "Indicates whether this reserved IP member will be automatically deleted when either target is deleted, or the reserved IP is unbound.",
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this reserved IP.",
									},
									"reserved_ip": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The unique identifier for this reserved IP.",
									},
									"name": &schema.Schema{
										Type: schema.TypeString,
										// Optional:    true,
										Computed:    true,
										Description: "The name for this reserved IP. The name is unique across all reserved IPs in a subnet.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"name": &schema.Schema{
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"virtual_network_interface.0.id"},
							ValidateFunc:  validate.InvokeValidator("ibm_is_virtual_network_interface", "vni_name"),
							Description:   "The name for this virtual network interface. The name is unique across all virtual network interfaces in the VPC.",
						},
						"protocol_state_filtering_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.InvokeValidator("ibm_is_virtual_network_interface", "protocol_state_filtering_mode"),
							Description:  "The protocol state filtering mode used for this virtual network interface.",
						},
						"primary_ip": &schema.Schema{
							Type:          schema.TypeList,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"virtual_network_interface.0.id"},
							Description:   "The primary IP address of the virtual network interface for the bare metal server networkattachment.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The IP address.If the address has not yet been selected, the value will be `0.0.0.0`.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.",
									},
									"deleted": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this reserved IP.",
									},
									"reserved_ip": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The unique identifier for this reserved IP.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The name for this reserved IP. The name is unique across all reserved IPs in a subnet.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"resource_group": &schema.Schema{
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"virtual_network_interface.0.id"},
							Computed:      true,
							Description:   "The resource group id for this virtual network interface.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						"security_groups": {
							Type:          schema.TypeSet,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"virtual_network_interface.0.id"},
							Elem:          &schema.Schema{Type: schema.TypeString},
							Set:           schema.HashString,
							Description:   "The security groups for this virtual network interface.",
						},
						"subnet": &schema.Schema{
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ForceNew:      true,
							ConflictsWith: []string{"virtual_network_interface.0.id"},
							Description:   "The associated subnet id.",
						},
					},
				},
			},
			"allowed_vlans": &schema.Schema{
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"vlan"},
				Set:           schema.HashInt,
				Description:   "Indicates what VLAN IDs (for VLAN type only) can use this physical (PCI type) attachment.",
				Elem:          &schema.Schema{Type: schema.TypeInt},
			},
			"allow_to_float": &schema.Schema{
				Type:          schema.TypeBool,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"allowed_vlans"},
				Description:   "Indicates if the bare metal server network attachment can automatically float to any other server within the same `resource_group`. The bare metal server network attachment will float automatically if the network detects a GARP or RARP on another bare metal server in the resource group. Applies only to bare metal server network attachments with `vlan` interface type.",
			},
			"vlan": &schema.Schema{
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"allowed_vlans"},
				Computed:      true,
				ValidateFunc:  validate.InvokeValidator("ibm_is_bare_metal_server_network_attachment", "vlan"),
				Description:   "Indicates the 802.1Q VLAN ID tag that must be used for all traffic on this attachment.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the bare metal server network attachment was created.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this bare metal server network attachment.",
			},
			"lifecycle_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the bare metal server network attachment.",
			},
			"port_speed": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The port speed for this bare metal server network attachment in Mbps.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The bare metal server network attachment type.",
			},
			isBareMetalServerHardStop: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Only used for PCI network attachments, whether to hard/immediately stop server",
			},
		},
	}
}

func ResourceIBMIsBareMetalServerNetworkAttachmentValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "bare_metal_server",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[-0-9a-z_]+$`,
			MinValueLength:             1,
			MaxValueLength:             64,
		},
		validate.ValidateSchema{
			Identifier:                 "interface_type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "pci, vlan",
		},
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		},
		validate.ValidateSchema{
			Identifier:                 "vlan",
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			Optional:                   true,
			MinValue:                   "1",
			MaxValue:                   "4094",
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_bare_metal_server_network_attachment", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMIsBareMetalServerNetworkAttachmentCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	rebootNeeded := true
	rebooted := false
	bodyModelMap := map[string]interface{}{}
	createBareMetalServerNetworkAttachmentOptions := &vpcv1.CreateBareMetalServerNetworkAttachmentOptions{}

	bodyModelMap["interface_type"] = d.Get("interface_type")
	if _, ok := d.GetOk("name"); ok {
		bodyModelMap["name"] = d.Get("name")
	}
	bodyModelMap["virtual_network_interface"] = d.Get("virtual_network_interface")
	_, bodyModelMap["allow_ip_spoofing_exists"] = d.GetOkExists("virtual_network_interface.0.allow_ip_spoofing")
	_, bodyModelMap["auto_delete_exists"] = d.GetOkExists("virtual_network_interface.0.auto_delete")
	_, bodyModelMap["enable_infrastructure_nat_exists"] = d.GetOkExists("virtual_network_interface.0.enable_infrastructure_nat")
	if _, ok := d.GetOk("allow_to_float"); ok {
		bodyModelMap["allow_to_float"] = d.Get("allow_to_float")
	}
	if _, ok := d.GetOk("vlan"); ok {
		bodyModelMap["vlan"] = d.Get("vlan")
		if int64(d.Get("vlan").(int)) != 0 {
			rebootNeeded = false
		}
	}
	if _, ok := d.GetOk("allowed_vlans"); ok {
		bodyModelMap["allowed_vlans"] = d.Get("allowed_vlans")
	}
	bareMetalServerId := d.Get("bare_metal_server").(string)
	createBareMetalServerNetworkAttachmentOptions.SetBareMetalServerID(bareMetalServerId)
	convertedModel, err := resourceIBMIsBareMetalServerNetworkAttachmentMapToBareMetalServerNetworkAttachmentPrototype(bodyModelMap)
	if err != nil {
		return diag.FromErr(err)
	}
	if rebootNeeded {
		getbmsoptions := &vpcv1.GetBareMetalServerOptions{
			ID: &bareMetalServerId,
		}

		bms, response, err := vpcClient.GetBareMetalServerWithContext(context, getbmsoptions)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error fetching bare metal server (%s) during network attachment create err %s\n%s", bareMetalServerId, err, response))
		}
		// failed, pending, restarting, running, starting, stopped, stopping, maintenance
		if *bms.Status == "failed" {
			return diag.FromErr(fmt.Errorf("[ERROR] Error cannot attach network attachment to a failed bare metal server"))
		} else if *bms.Status == "running" && rebootNeeded {
			log.Printf("[DEBUG] Stopping bare metal server (%s) to create a PCI network attachment", bareMetalServerId)
			stopType := "hard"
			if _, ok := d.GetOk(isBareMetalServerHardStop); ok && !d.Get(isBareMetalServerHardStop).(bool) {
				stopType = "soft"
			}
			createstopaction := &vpcv1.StopBareMetalServerOptions{
				ID:   &bareMetalServerId,
				Type: &stopType,
			}
			res, err := vpcClient.StopBareMetalServerWithContext(context, createstopaction)
			if err != nil || res.StatusCode != 204 {
				return diag.FromErr(fmt.Errorf("[ERROR] Error stopping bare metal server (%s) during network attachment create err %s\n%s", bareMetalServerId, err, response))
			}
			_, err = isWaitForBareMetalServerStoppedForNIC(vpcClient, bareMetalServerId, d.Timeout(schema.TimeoutCreate), d)
			if err != nil {
				return diag.FromErr(err)
			}
			rebooted = true
		} else if *bms.Status != "stopped" {
			return diag.FromErr(fmt.Errorf("[ERROR] Error bare metal server in %s state, please try after some time", *bms.Status))
		}
	}
	createBareMetalServerNetworkAttachmentOptions.BareMetalServerNetworkAttachmentPrototype = convertedModel

	bareMetalServerNetworkAttachmentIntf, response, err := vpcClient.CreateBareMetalServerNetworkAttachmentWithContext(context, createBareMetalServerNetworkAttachmentOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateBareMetalServerNetworkAttachmentWithContext failed %s\n%s", err, response)
		if rebooted {
			createstartaction := &vpcv1.StartBareMetalServerOptions{
				ID: &bareMetalServerId,
			}
			res, err := vpcClient.StartBareMetalServerWithContext(context, createstartaction)
			if err != nil || res.StatusCode != 204 {
				return diag.FromErr(fmt.Errorf("[ERROR] Error starting bare metal server (%s) after attachment creation failed err %s\n%s", bareMetalServerId, err, response))
			}
			_, err = isWaitForBareMetalServerAvailableForNIC(vpcClient, bareMetalServerId, d.Timeout(schema.TimeoutCreate), d)
			if err != nil {
				return diag.FromErr(err)
			}
		}
		return diag.FromErr(fmt.Errorf("CreateBareMetalServerNetworkAttachmentWithContext failed %s\n%s", err, response))
	}

	if _, ok := bareMetalServerNetworkAttachmentIntf.(*vpcv1.BareMetalServerNetworkAttachmentByVlan); ok {
		bareMetalServerNetworkAttachment := bareMetalServerNetworkAttachmentIntf.(*vpcv1.BareMetalServerNetworkAttachmentByVlan)
		if bareMetalServerNetworkAttachment.VirtualNetworkInterface != nil && bareMetalServerNetworkAttachment.VirtualNetworkInterface.ID != nil {
			d.Set("virtual_network_interface_id", *bareMetalServerNetworkAttachment.VirtualNetworkInterface.ID)
		}
		d.SetId(fmt.Sprintf("%s/%s", *createBareMetalServerNetworkAttachmentOptions.BareMetalServerID, *bareMetalServerNetworkAttachment.ID))
		d.Set("floating_bare_metal_server", *createBareMetalServerNetworkAttachmentOptions.BareMetalServerID)
	} else if _, ok := bareMetalServerNetworkAttachmentIntf.(*vpcv1.BareMetalServerNetworkAttachmentByPci); ok {
		bareMetalServerNetworkAttachment := bareMetalServerNetworkAttachmentIntf.(*vpcv1.BareMetalServerNetworkAttachmentByPci)
		if bareMetalServerNetworkAttachment.VirtualNetworkInterface != nil && bareMetalServerNetworkAttachment.VirtualNetworkInterface.ID != nil {
			d.Set("virtual_network_interface_id", *bareMetalServerNetworkAttachment.VirtualNetworkInterface.ID)
		}
		d.SetId(fmt.Sprintf("%s/%s", *createBareMetalServerNetworkAttachmentOptions.BareMetalServerID, *bareMetalServerNetworkAttachment.ID))
		d.Set("floating_bare_metal_server", *createBareMetalServerNetworkAttachmentOptions.BareMetalServerID)
	} else if _, ok := bareMetalServerNetworkAttachmentIntf.(*vpcv1.BareMetalServerNetworkAttachment); ok {
		bareMetalServerNetworkAttachment := bareMetalServerNetworkAttachmentIntf.(*vpcv1.BareMetalServerNetworkAttachment)
		if bareMetalServerNetworkAttachment.VirtualNetworkInterface != nil && bareMetalServerNetworkAttachment.VirtualNetworkInterface.ID != nil {
			d.Set("virtual_network_interface_id", *bareMetalServerNetworkAttachment.VirtualNetworkInterface.ID)
		}
		d.SetId(fmt.Sprintf("%s/%s", *createBareMetalServerNetworkAttachmentOptions.BareMetalServerID, *bareMetalServerNetworkAttachment.ID))
		d.Set("floating_bare_metal_server", *createBareMetalServerNetworkAttachmentOptions.BareMetalServerID)
	} else {
		return diag.FromErr(fmt.Errorf("Unrecognized vpcv1.BareMetalServerNetworkAttachmentIntf subtype encountered"))
	}

	if rebooted {
		createstartaction := &vpcv1.StartBareMetalServerOptions{
			ID: &bareMetalServerId,
		}
		res, err := vpcClient.StartBareMetalServerWithContext(context, createstartaction)
		if err != nil || res.StatusCode != 204 {
			return diag.FromErr(fmt.Errorf("[ERROR] Error starting bare metal server (%s) after attachment creation err %s\n%s", bareMetalServerId, err, response))
		}
		_, err = isWaitForBareMetalServerAvailableForNIC(vpcClient, bareMetalServerId, d.Timeout(schema.TimeoutCreate), d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceIBMIsBareMetalServerNetworkAttachmentRead(context, d, meta)
}

func resourceIBMIsBareMetalServerNetworkAttachmentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	getBareMetalServerNetworkAttachmentOptions := &vpcv1.GetBareMetalServerNetworkAttachmentOptions{}
	bmId := ""
	nacId := ""

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}
	bmId = parts[0]
	nacId = parts[1]
	getBareMetalServerNetworkAttachmentOptions.SetBareMetalServerID(bmId)
	getBareMetalServerNetworkAttachmentOptions.SetID(nacId)
	bareMetalServerNetworkAttachmentIntf, response, err := vpcClient.GetBareMetalServerNetworkAttachmentWithContext(context, getBareMetalServerNetworkAttachmentOptions)
	if err != nil {

		if response != nil && response.StatusCode == 404 {

			allowToFloatIntf := d.Get("allow_to_float")
			if allowToFloatIntf != nil && allowToFloatIntf.(bool) {

				vniid := d.Get("virtual_network_interface.0.id").(string)
				getVirtualNetworkInterfaceOptions := &vpcv1.GetVirtualNetworkInterfaceOptions{
					ID: &vniid,
				}
				vniDetails, response, err := vpcClient.GetVirtualNetworkInterface(getVirtualNetworkInterfaceOptions)
				if err != nil {
					if response != nil && response.StatusCode == 404 {
						d.SetId("")
						return nil
					}
					return diag.FromErr(fmt.Errorf("[ERROR] Error on GetVirtualNetworkInterface in BareMetalServer : %s\n%s", err, response))
				}

				vniTargetIntf := vniDetails.Target

				if _, ok := vniTargetIntf.(*vpcv1.VirtualNetworkInterfaceTargetShareMountTargetReference); ok {
					vniTarget := vniTargetIntf.(*vpcv1.VirtualNetworkInterfaceTargetShareMountTargetReference)
					ree := regexp.MustCompile(`([^/]+)/network_attachments/([^/]+)`)

					// Find the matches
					matches := ree.FindStringSubmatch(*vniTarget.Href)
					if len(matches) >= 3 {
						bmId = (matches[1])
						nacId = matches[2]
					}

				} else if _, ok := vniTargetIntf.(*vpcv1.VirtualNetworkInterfaceTargetInstanceNetworkAttachmentReferenceVirtualNetworkInterfaceContext); ok {
					vniTarget := vniTargetIntf.(*vpcv1.VirtualNetworkInterfaceTargetInstanceNetworkAttachmentReferenceVirtualNetworkInterfaceContext)
					ree := regexp.MustCompile(`([^/]+)/network_attachments/([^/]+)`)

					// Find the matches
					matches := ree.FindStringSubmatch(*vniTarget.Href)
					if len(matches) >= 3 {
						bmId = (matches[1])
						nacId = matches[2]
					}

				} else if _, ok := vniTargetIntf.(*vpcv1.VirtualNetworkInterfaceTargetBareMetalServerNetworkAttachmentReferenceVirtualNetworkInterfaceContext); ok {
					vniTarget := vniTargetIntf.(*vpcv1.VirtualNetworkInterfaceTargetBareMetalServerNetworkAttachmentReferenceVirtualNetworkInterfaceContext)
					ree := regexp.MustCompile(`([^/]+)/network_attachments/([^/]+)`)

					// Find the matches
					matches := ree.FindStringSubmatch(*vniTarget.Href)
					if len(matches) >= 3 {
						bmId = (matches[1])
						nacId = matches[2]
					}

				} else if _, ok := vniTargetIntf.(*vpcv1.VirtualNetworkInterfaceTarget); ok {
					vniTarget := vniTargetIntf.(*vpcv1.VirtualNetworkInterfaceTarget)
					ree := regexp.MustCompile(`([^/]+)/network_attachments/([^/]+)`)

					// Find the matches
					matches := ree.FindStringSubmatch(*vniTarget.Href)
					if len(matches) >= 3 {
						bmId = (matches[1])
						nacId = matches[2]
					}

				}

				getBareMetalServerNetworkAttachmentOptions.SetBareMetalServerID(bmId)
				getBareMetalServerNetworkAttachmentOptions.SetID(nacId)

				bareMetalServerNetworkAttachmentIntf, response, err = vpcClient.GetBareMetalServerNetworkAttachmentWithContext(context, getBareMetalServerNetworkAttachmentOptions)
				if err != nil {

					if response != nil && response.StatusCode == 404 {
						d.SetId("")
						return nil
					}
				}
			}

			log.Printf("[DEBUG] GetBareMetalServerNetworkAttachmentWithContext failed %s\n%s", err, response)
			// return diag.FromErr(fmt.Errorf("GetBareMetalServerNetworkAttachmentWithContext failed %s\n%s", err, response))
		} else {
			log.Printf("[DEBUG] GetBareMetalServerNetworkAttachmentWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("GetBareMetalServerNetworkAttachmentWithContext failed %s\n%s", err, response))
		}
	}

	if _, ok := bareMetalServerNetworkAttachmentIntf.(*vpcv1.BareMetalServerNetworkAttachmentByVlan); ok {
		bareMetalServerNetworkAttachment := bareMetalServerNetworkAttachmentIntf.(*vpcv1.BareMetalServerNetworkAttachmentByVlan)
		d.SetId(fmt.Sprintf("%s/%s", bmId, *bareMetalServerNetworkAttachment.ID))
		if err = d.Set("interface_type", bareMetalServerNetworkAttachment.InterfaceType); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting interface_type: %s", err))
		}
		if !core.IsNil(bareMetalServerNetworkAttachment.Name) {
			if err = d.Set("name", bareMetalServerNetworkAttachment.Name); err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
			}
		}
		virtualNetworkInterfaceMap, err := resourceIBMIsBareMetalServerNetworkAttachmentVirtualNetworkInterfaceReferenceAttachmentContextToMap(bareMetalServerNetworkAttachment.VirtualNetworkInterface, vpcClient)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("virtual_network_interface", []map[string]interface{}{virtualNetworkInterfaceMap}); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting virtual_network_interface: %s", err))
		}
		if !core.IsNil(bareMetalServerNetworkAttachment.AllowToFloat) {
			if err = d.Set("allow_to_float", bareMetalServerNetworkAttachment.AllowToFloat); err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error setting allow_to_float: %s", err))
			}
		}
		if !core.IsNil(bareMetalServerNetworkAttachment.Vlan) {
			if err = d.Set("vlan", flex.IntValue(bareMetalServerNetworkAttachment.Vlan)); err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error setting vlan: %s", err))
			}
		}
		if err = d.Set("created_at", flex.DateTimeToString(bareMetalServerNetworkAttachment.CreatedAt)); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
		}
		if err = d.Set("href", bareMetalServerNetworkAttachment.Href); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting href: %s", err))
		}
		if err = d.Set("lifecycle_state", bareMetalServerNetworkAttachment.LifecycleState); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting lifecycle_state: %s", err))
		}
		if err = d.Set("port_speed", flex.IntValue(bareMetalServerNetworkAttachment.PortSpeed)); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting port_speed: %s", err))
		}
		if err = d.Set("resource_type", bareMetalServerNetworkAttachment.ResourceType); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_type: %s", err))
		}
		if err = d.Set("type", bareMetalServerNetworkAttachment.Type); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting type: %s", err))
		}
	} else if _, ok := bareMetalServerNetworkAttachmentIntf.(*vpcv1.BareMetalServerNetworkAttachmentByPci); ok {
		bareMetalServerNetworkAttachment := bareMetalServerNetworkAttachmentIntf.(*vpcv1.BareMetalServerNetworkAttachmentByPci)
		if err = d.Set("interface_type", bareMetalServerNetworkAttachment.InterfaceType); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting interface_type: %s", err))
		}
		if !core.IsNil(bareMetalServerNetworkAttachment.Name) {
			if err = d.Set("name", bareMetalServerNetworkAttachment.Name); err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
			}
		}
		virtualNetworkInterfaceMap, err := resourceIBMIsBareMetalServerNetworkAttachmentVirtualNetworkInterfaceReferenceAttachmentContextToMap(bareMetalServerNetworkAttachment.VirtualNetworkInterface, vpcClient)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("virtual_network_interface", []map[string]interface{}{virtualNetworkInterfaceMap}); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting virtual_network_interface: %s", err))
		}
		if !core.IsNil(bareMetalServerNetworkAttachment.AllowedVlans) {
			allowedVlans := []interface{}{}
			for _, allowedVlansItem := range bareMetalServerNetworkAttachment.AllowedVlans {
				allowedVlans = append(allowedVlans, int64(allowedVlansItem))
			}
			if err = d.Set("allowed_vlans", allowedVlans); err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error setting allowed_vlans: %s", err))
			}
		}
		if err = d.Set("created_at", flex.DateTimeToString(bareMetalServerNetworkAttachment.CreatedAt)); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
		}
		if err = d.Set("href", bareMetalServerNetworkAttachment.Href); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting href: %s", err))
		}
		if err = d.Set("lifecycle_state", bareMetalServerNetworkAttachment.LifecycleState); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting lifecycle_state: %s", err))
		}
		if err = d.Set("port_speed", flex.IntValue(bareMetalServerNetworkAttachment.PortSpeed)); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting port_speed: %s", err))
		}
		if err = d.Set("resource_type", bareMetalServerNetworkAttachment.ResourceType); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_type: %s", err))
		}

		if err = d.Set("type", bareMetalServerNetworkAttachment.Type); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting type: %s", err))
		}
		d.SetId(fmt.Sprintf("%s/%s", bmId, *bareMetalServerNetworkAttachment.ID))
	} else if _, ok := bareMetalServerNetworkAttachmentIntf.(*vpcv1.BareMetalServerNetworkAttachment); ok {
		bareMetalServerNetworkAttachment := bareMetalServerNetworkAttachmentIntf.(*vpcv1.BareMetalServerNetworkAttachment)
		d.SetId(fmt.Sprintf("%s/%s", bmId, *bareMetalServerNetworkAttachment.ID))

		// parent class argument: bare_metal_server string
		if err = d.Set("floating_bare_metal_server", getBareMetalServerNetworkAttachmentOptions.BareMetalServerID); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting bare_metal_server: %s", err))
		}
		// parent class argument: interface_type string
		// parent class argument: name string
		// parent class argument: virtual_network_interface VirtualNetworkInterfaceReferenceAttachmentContext
		// parent class argument: allowed_vlans []int64
		// parent class argument: allow_to_float bool
		// parent class argument: vlan int64
		if err = d.Set("interface_type", bareMetalServerNetworkAttachment.InterfaceType); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting interface_type: %s", err))
		}
		if !core.IsNil(bareMetalServerNetworkAttachment.Name) {
			if err = d.Set("name", bareMetalServerNetworkAttachment.Name); err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
			}
		}
		virtualNetworkInterfaceMap, err := resourceIBMIsBareMetalServerNetworkAttachmentVirtualNetworkInterfaceReferenceAttachmentContextToMap(bareMetalServerNetworkAttachment.VirtualNetworkInterface, vpcClient)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("virtual_network_interface", []map[string]interface{}{virtualNetworkInterfaceMap}); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting virtual_network_interface: %s", err))
		}
		if !core.IsNil(bareMetalServerNetworkAttachment.AllowedVlans) {
			allowedVlans := []interface{}{}
			for _, allowedVlansItem := range bareMetalServerNetworkAttachment.AllowedVlans {
				allowedVlans = append(allowedVlans, int64(allowedVlansItem))
			}
			if err = d.Set("allowed_vlans", allowedVlans); err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error setting allowed_vlans: %s", err))
			}
		}
		if !core.IsNil(bareMetalServerNetworkAttachment.AllowToFloat) {
			if err = d.Set("allow_to_float", bareMetalServerNetworkAttachment.AllowToFloat); err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error setting allow_to_float: %s", err))
			}
		}
		if !core.IsNil(bareMetalServerNetworkAttachment.Vlan) {
			if err = d.Set("vlan", flex.IntValue(bareMetalServerNetworkAttachment.Vlan)); err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error setting vlan: %s", err))
			}
		}
		if err = d.Set("created_at", flex.DateTimeToString(bareMetalServerNetworkAttachment.CreatedAt)); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
		}
		if err = d.Set("href", bareMetalServerNetworkAttachment.Href); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting href: %s", err))
		}
		if err = d.Set("lifecycle_state", bareMetalServerNetworkAttachment.LifecycleState); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting lifecycle_state: %s", err))
		}
		if err = d.Set("port_speed", flex.IntValue(bareMetalServerNetworkAttachment.PortSpeed)); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting port_speed: %s", err))
		}
		if err = d.Set("resource_type", bareMetalServerNetworkAttachment.ResourceType); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_type: %s", err))
		}
		if err = d.Set("type", bareMetalServerNetworkAttachment.Type); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting type: %s", err))
		}
		if err = d.Set("network_attachment", bareMetalServerNetworkAttachment.ID); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting network_attachment: %s", err))
		}
	} else {
		return diag.FromErr(fmt.Errorf("Unrecognized vpcv1.BareMetalServerNetworkAttachmentIntf subtype encountered"))
	}

	return nil
}

func resourceIBMIsBareMetalServerNetworkAttachmentUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	updateBareMetalServerNetworkAttachmentOptions := &vpcv1.UpdateBareMetalServerNetworkAttachmentOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	updateBareMetalServerNetworkAttachmentOptions.SetBareMetalServerID(parts[0])
	updateBareMetalServerNetworkAttachmentOptions.SetID(parts[1])

	hasChange := false

	patchVals := &vpcv1.BareMetalServerNetworkAttachmentPatch{}
	if d.HasChange("bare_metal_server") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "bare_metal_server"))
	}
	if d.HasChange("allowed_vlans") {
		var allowedVlans []int64
		for _, v := range d.Get("allowed_vlans").(*schema.Set).List() {
			allowedVlansItem := int64(v.(int))
			allowedVlans = append(allowedVlans, allowedVlansItem)
		}
		patchVals.AllowedVlans = allowedVlans
		hasChange = true
	}
	if d.HasChange("name") {
		newName := d.Get("name").(string)
		patchVals.Name = &newName
		hasChange = true
	}

	if hasChange {
		updateBareMetalServerNetworkAttachmentOptions.BareMetalServerNetworkAttachmentPatch, _ = patchVals.AsPatch()
		_, response, err := vpcClient.UpdateBareMetalServerNetworkAttachmentWithContext(context, updateBareMetalServerNetworkAttachmentOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateBareMetalServerNetworkAttachmentWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateBareMetalServerNetworkAttachmentWithContext failed %s\n%s", err, response))
		}
	}

	if d.HasChange("virtual_network_interface") && !d.IsNewResource() {
		vniId := d.Get("virtual_network_interface.0.id").(string)
		updateVirtualNetworkInterfaceOptions := &vpcv1.UpdateVirtualNetworkInterfaceOptions{
			ID: &vniId,
		}
		virtualNetworkInterfacePatch := &vpcv1.VirtualNetworkInterfacePatch{}
		if d.HasChange("virtual_network_interface.0.auto_delete") {
			autodelete := d.Get("virtual_network_interface.0.auto_delete").(bool)
			virtualNetworkInterfacePatch.AutoDelete = &autodelete
		}
		if d.HasChange("virtual_network_interface.0.name") {
			name := d.Get("virtual_network_interface.0.name").(string)
			virtualNetworkInterfacePatch.Name = &name
		}
		if d.HasChange("virtual_network_interface.0.enable_infrastructure_nat") {
			enableNat := d.Get("virtual_network_interface.0.enable_infrastructure_nat").(bool)
			virtualNetworkInterfacePatch.EnableInfrastructureNat = &enableNat
		}
		if d.HasChange("virtual_network_interface.0.allow_ip_spoofing") {
			allIpSpoofing := d.Get("virtual_network_interface.0.allow_ip_spoofing").(bool)
			virtualNetworkInterfacePatch.AllowIPSpoofing = &allIpSpoofing
		}
		virtualNetworkInterfacePatchAsPatch, err := virtualNetworkInterfacePatch.AsPatch()
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error encountered while apply as patch for virtualNetworkInterfacePatch of BareMetalServer(%s) vni (%s) %s", d.Id(), vniId, err))
		}
		updateVirtualNetworkInterfaceOptions.VirtualNetworkInterfacePatch = virtualNetworkInterfacePatchAsPatch
		_, response, err := vpcClient.UpdateVirtualNetworkInterfaceWithContext(context, updateVirtualNetworkInterfaceOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateVirtualNetworkInterfaceWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateVirtualNetworkInterfaceWithContext failed during BareMetalServer(%s) network attachment patch %s\n%s", d.Id(), err, response))
		}

		if d.HasChange("virtual_network_interface.0.ips") {
			oldips, newips := d.GetChange("virtual_network_interface.0.ips")
			os := oldips.(*schema.Set)
			ns := newips.(*schema.Set)
			var oldset, newset *schema.Set

			var out = make([]interface{}, ns.Len(), ns.Len())
			for i, nA := range ns.List() {
				newPack := nA.(map[string]interface{})
				out[i] = newPack["reserved_ip"].(string)
			}
			newset = schema.NewSet(schema.HashString, out)

			out = make([]interface{}, os.Len(), os.Len())
			for i, oA := range os.List() {
				oldPack := oA.(map[string]interface{})
				out[i] = oldPack["reserved_ip"].(string)
			}
			oldset = schema.NewSet(schema.HashString, out)

			remove := flex.ExpandStringList(oldset.Difference(newset).List())
			add := flex.ExpandStringList(newset.Difference(oldset).List())

			if add != nil && len(add) > 0 {
				for _, ipItem := range add {
					if ipItem != "" {

						addVirtualNetworkInterfaceIPOptions := &vpcv1.AddVirtualNetworkInterfaceIPOptions{}
						addVirtualNetworkInterfaceIPOptions.SetVirtualNetworkInterfaceID(vniId)
						addVirtualNetworkInterfaceIPOptions.SetID(ipItem)
						_, response, err := vpcClient.AddVirtualNetworkInterfaceIPWithContext(context, addVirtualNetworkInterfaceIPOptions)
						if err != nil {
							log.Printf("[DEBUG] AddVirtualNetworkInterfaceIPWithContext failed in VirtualNetworkInterface patch during BareMetalServer nac patch %s\n%s", err, response)
							return diag.FromErr(fmt.Errorf("AddVirtualNetworkInterfaceIPWithContext failed in VirtualNetworkInterface patch during BareMetalServer nac patch %s\n%s", err, response))
						}
					}
				}
			}
			if remove != nil && len(remove) > 0 {
				for _, ipItem := range remove {
					if ipItem != "" {

						removeVirtualNetworkInterfaceIPOptions := &vpcv1.RemoveVirtualNetworkInterfaceIPOptions{}
						removeVirtualNetworkInterfaceIPOptions.SetVirtualNetworkInterfaceID(vniId)
						removeVirtualNetworkInterfaceIPOptions.SetID(ipItem)
						response, err := vpcClient.RemoveVirtualNetworkInterfaceIPWithContext(context, removeVirtualNetworkInterfaceIPOptions)
						if err != nil {
							log.Printf("[DEBUG] RemoveVirtualNetworkInterfaceIPWithContext failed in VirtualNetworkInterface patch during BareMetalServer nac patch %s\n%s", err, response)
							return diag.FromErr(fmt.Errorf("RemoveVirtualNetworkInterfaceIPWithContext failed in VirtualNetworkInterface patch during BareMetalServer nac patch %s\n%s", err, response))
						}
					}
				}
			}

		}
		if d.HasChange("virtual_network_interface.0.primary_ip") {
			subnetId := d.Get("virtual_network_interface.0.subnet").(string)
			ripId := d.Get("virtual_network_interface.0.primary_ip.0.reserved_ip").(string)
			updateripoptions := &vpcv1.UpdateSubnetReservedIPOptions{
				SubnetID: &subnetId,
				ID:       &ripId,
			}
			reservedIpPath := &vpcv1.ReservedIPPatch{}
			if d.HasChange("virtual_network_interface.0.primary_ip.0.name") {
				name := d.Get("virtual_network_interface.0.primary_ip.0.name").(string)
				reservedIpPath.Name = &name
			}
			if d.HasChange("virtual_network_interface.0.primary_ip.0.auto_delete") {
				auto := d.Get("virtual_network_interface.0.primary_ip.0.auto_delete").(bool)
				reservedIpPath.AutoDelete = &auto
			}
			reservedIpPathAsPatch, err := reservedIpPath.AsPatch()
			if err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error calling reserved ip as patch on vni patch \n%s", err))
			}
			updateripoptions.ReservedIPPatch = reservedIpPathAsPatch
			_, response, err := vpcClient.UpdateSubnetReservedIP(updateripoptions)
			if err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error updating vni reserved ip(%s): %s\n%s", ripId, err, response))
			}
		}
		if d.HasChange("virtual_network_interface.0.protocol_state_filtering_mode") {
			psfMode := d.Get("virtual_network_interface.0.protocol_state_filtering_mode").(string)
			virtualNetworkInterfacePatch.ProtocolStateFilteringMode = &psfMode
		}
		if d.HasChange("virtual_network_interface.0.security_groups") {
			ovs, nvs := d.GetChange("virtual_network_interface.0.security_groups")
			ov := ovs.(*schema.Set)
			nv := nvs.(*schema.Set)
			remove := flex.ExpandStringList(ov.Difference(nv).List())
			add := flex.ExpandStringList(nv.Difference(ov).List())
			if len(add) > 0 {
				for i := range add {
					createsgnicoptions := &vpcv1.CreateSecurityGroupTargetBindingOptions{
						SecurityGroupID: &add[i],
						ID:              &vniId,
					}
					_, response, err := vpcClient.CreateSecurityGroupTargetBinding(createsgnicoptions)
					if err != nil {
						return diag.FromErr(fmt.Errorf("[ERROR] Error while creating security group %q for virtual network interface %s\n%s: %q", add[i], d.Id(), err, response))
					}
					_, err = isWaitForVirtualNetworkInterfaceAvailable(vpcClient, vniId, d.Timeout(schema.TimeoutUpdate))
					if err != nil {
						return diag.FromErr(err)
					}
				}

			}
			if len(remove) > 0 {
				for i := range remove {
					deletesgnicoptions := &vpcv1.DeleteSecurityGroupTargetBindingOptions{
						SecurityGroupID: &remove[i],
						ID:              &vniId,
					}
					response, err := vpcClient.DeleteSecurityGroupTargetBinding(deletesgnicoptions)
					if err != nil {
						return diag.FromErr(fmt.Errorf("[ERROR] Error while removing security group %q for virtual network interface %s\n%s: %q", remove[i], d.Id(), err, response))
					}
					_, err = isWaitForVirtualNetworkInterfaceAvailable(vpcClient, vniId, d.Timeout(schema.TimeoutUpdate))
					if err != nil {
						return diag.FromErr(err)
					}
				}
			}
		}

	}

	return resourceIBMIsBareMetalServerNetworkAttachmentRead(context, d, meta)
}

func resourceIBMIsBareMetalServerNetworkAttachmentDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	ifServerStopped := false
	interfaceType := d.Get("interface_type").(string)
	bareMetalServerId := d.Get("floating_bare_metal_server").(string)
	if interfaceType == "pci" {
		getbmsoptions := &vpcv1.GetBareMetalServerOptions{
			ID: &bareMetalServerId,
		}

		bms, response, err := vpcClient.GetBareMetalServerWithContext(context, getbmsoptions)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error fetching bare metal server (%s) during network attachment create err %s\n%s", bareMetalServerId, err, response))
		}
		// failed, pending, restarting, running, starting, stopped, stopping, maintenance
		if *bms.Status == "failed" {
			return diag.FromErr(fmt.Errorf("[ERROR] Error cannot attach network attachment to a failed bare metal server"))
		} else if *bms.Status == "running" {
			log.Printf("[DEBUG] Stopping bare metal server (%s) to create a PCI network attachment", bareMetalServerId)
			stopType := "hard"
			if _, ok := d.GetOk(isBareMetalServerHardStop); ok && !d.Get(isBareMetalServerHardStop).(bool) {
				stopType = "soft"
			}
			createstopaction := &vpcv1.StopBareMetalServerOptions{
				ID:   &bareMetalServerId,
				Type: &stopType,
			}
			res, err := vpcClient.StopBareMetalServerWithContext(context, createstopaction)
			ifServerStopped = true
			if err != nil || res.StatusCode != 204 {
				return diag.FromErr(fmt.Errorf("[ERROR] Error stopping bare metal server (%s) during network attachment create err %s\n%s", bareMetalServerId, err, response))
			}
			_, err = isWaitForBareMetalServerStoppedForNIC(vpcClient, bareMetalServerId, d.Timeout(schema.TimeoutCreate), d)
			if err != nil {
				return diag.FromErr(err)
			}
		} else if *bms.Status != "stopped" {
			return diag.FromErr(fmt.Errorf("[ERROR] Error bare metal server in %s state, please try after some time", *bms.Status))
		}
	}
	deleteBareMetalServerNetworkAttachmentOptions := &vpcv1.DeleteBareMetalServerNetworkAttachmentOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteBareMetalServerNetworkAttachmentOptions.SetBareMetalServerID(parts[0])
	deleteBareMetalServerNetworkAttachmentOptions.SetID(parts[1])

	response, err := vpcClient.DeleteBareMetalServerNetworkAttachmentWithContext(context, deleteBareMetalServerNetworkAttachmentOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteBareMetalServerNetworkAttachmentWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteBareMetalServerNetworkAttachmentWithContext failed %s\n%s", err, response))
	}
	if ifServerStopped {
		createstartaction := &vpcv1.StartBareMetalServerOptions{
			ID: &bareMetalServerId,
		}
		res, err := vpcClient.StartBareMetalServerWithContext(context, createstartaction)
		if err != nil || res.StatusCode != 204 {
			return diag.FromErr(fmt.Errorf("[ERROR] Error starting bare metal server (%s) after attachment creation err %s\n%s", bareMetalServerId, err, response))
		}
		_, err = isWaitForBareMetalServerAvailableForNIC(vpcClient, bareMetalServerId, d.Timeout(schema.TimeoutCreate), d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId("")

	return nil
}

func resourceIBMIsBareMetalServerNetworkAttachmentMapToBareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterface(modelMap map[string]interface{}, allow_ip_spoofing_exists, auto_delete_exists, enable_infrastructure_nat_exists interface{}) (vpcv1.BareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterfaceIntf, error) {
	model := &vpcv1.BareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterface{}
	if allow_ip_spoofing_exists.(bool) && modelMap["allow_ip_spoofing"] != nil {
		model.AllowIPSpoofing = core.BoolPtr(modelMap["allow_ip_spoofing"].(bool))
	}
	if auto_delete_exists.(bool) && modelMap["auto_delete"] != nil {
		model.AutoDelete = core.BoolPtr(modelMap["auto_delete"].(bool))
	}
	if enable_infrastructure_nat_exists.(bool) && modelMap["enable_infrastructure_nat"] != nil {
		model.EnableInfrastructureNat = core.BoolPtr(modelMap["enable_infrastructure_nat"].(bool))
	}
	if modelMap["ips"] != nil && modelMap["ips"].(*schema.Set).Len() > 0 {
		ips := []vpcv1.VirtualNetworkInterfaceIPPrototypeIntf{}
		for _, ipsItem := range modelMap["ips"].(*schema.Set).List() {
			ipsItemModel, err := resourceIBMIsBareMetalServerNetworkAttachmentMapToVirtualNetworkInterfaceIPPrototype(ipsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			ips = append(ips, ipsItemModel)
		}
		model.Ips = ips
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["protocol_state_filtering_mode"] != nil {
		if pStateFilteringInt, ok := modelMap["protocol_state_filtering_mode"]; ok && pStateFilteringInt.(string) != "" {
			model.ProtocolStateFilteringMode = core.StringPtr(pStateFilteringInt.(string))
		}
	}
	if modelMap["primary_ip"] != nil && len(modelMap["primary_ip"].([]interface{})) > 0 {
		PrimaryIPModel, err := resourceIBMIsBareMetalServerNetworkAttachmentMapToVirtualNetworkInterfacePrimaryIPPrototype(modelMap["primary_ip"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.PrimaryIP = PrimaryIPModel
	}
	if modelMap["resource_group"] != nil && modelMap["resource_group"].(string) != "" {

		model.ResourceGroup = &vpcv1.ResourceGroupIdentity{
			ID: core.StringPtr(modelMap["resource_group"].(string)),
		}
	}
	if modelMap["security_groups"] != nil && modelMap["security_groups"].(*schema.Set).Len() > 0 {
		securityGroups := []vpcv1.SecurityGroupIdentityIntf{}
		for _, securityGroupsItem := range modelMap["security_groups"].(*schema.Set).List() {
			securityGroupsItemModel := &vpcv1.SecurityGroupIdentity{
				ID: core.StringPtr(securityGroupsItem.(string)),
			}
			securityGroups = append(securityGroups, securityGroupsItemModel)
		}
		model.SecurityGroups = securityGroups
	}
	if modelMap["subnet"] != nil && modelMap["subnet"].(string) != "" {
		model.Subnet = &vpcv1.SubnetIdentity{
			ID: core.StringPtr(modelMap["subnet"].(string)),
		}
	}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	return model, nil
}

func resourceIBMIsBareMetalServerNetworkAttachmentMapToVirtualNetworkInterfaceIPPrototype(modelMap map[string]interface{}) (vpcv1.VirtualNetworkInterfaceIPPrototypeIntf, error) {
	model := &vpcv1.VirtualNetworkInterfaceIPPrototype{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	return model, nil
}

func resourceIBMIsBareMetalServerNetworkAttachmentMapToVirtualNetworkInterfacePrimaryIPPrototype(modelMap map[string]interface{}) (vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeIntf, error) {
	model := &vpcv1.VirtualNetworkInterfacePrimaryIPPrototype{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	if modelMap["address"] != nil && modelMap["address"].(string) != "" {
		model.Address = core.StringPtr(modelMap["address"].(string))
	}
	if modelMap["auto_delete"] != nil {
		model.AutoDelete = core.BoolPtr(modelMap["auto_delete"].(bool))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	return model, nil
}

func resourceIBMIsBareMetalServerNetworkAttachmentMapToBareMetalServerNetworkAttachmentPrototype(modelMap map[string]interface{}) (vpcv1.BareMetalServerNetworkAttachmentPrototypeIntf, error) {
	model := &vpcv1.BareMetalServerNetworkAttachmentPrototype{}
	if modelMap["vlan"] != nil && int64(modelMap["vlan"].(int)) != 0 {
		model.InterfaceType = core.StringPtr("vlan")
	} else {
		model.InterfaceType = core.StringPtr("pci")
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	VirtualNetworkInterfaceModel, err := resourceIBMIsBareMetalServerNetworkAttachmentMapToBareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterface(modelMap["virtual_network_interface"].([]interface{})[0].(map[string]interface{}), modelMap["allow_ip_spoofing_exists"], modelMap["auto_delete_exists"], modelMap["enable_infrastructure_nat_exists"])
	if err != nil {
		return model, err
	}
	model.VirtualNetworkInterface = VirtualNetworkInterfaceModel
	if modelMap["allowed_vlans"] != nil {
		allowedVlans := []int64{}
		for _, allowedVlansItem := range modelMap["allowed_vlans"].(*schema.Set).List() {
			allowedVlans = append(allowedVlans, int64(allowedVlansItem.(int)))
		}
		model.AllowedVlans = allowedVlans
	}
	if modelMap["allow_to_float"] != nil {
		model.AllowToFloat = core.BoolPtr(modelMap["allow_to_float"].(bool))
	}
	if modelMap["vlan"] != nil {
		model.Vlan = core.Int64Ptr(int64(modelMap["vlan"].(int)))
	}
	return model, nil
}

func resourceIBMIsBareMetalServerNetworkAttachmentVirtualNetworkInterfaceReferenceAttachmentContextToMap(model *vpcv1.VirtualNetworkInterfaceReferenceAttachmentContext, sess *vpcv1.VpcV1) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	vniid := *model.ID
	getVirtualNetworkInterfaceOptions := &vpcv1.GetVirtualNetworkInterfaceOptions{
		ID: &vniid,
	}
	vniDetails, response, err := sess.GetVirtualNetworkInterface(getVirtualNetworkInterfaceOptions)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error on GetBareMetalServerNetworkAttachment in BareMetalServer : %s\n%s", err, response)
	}
	modelMap["allow_ip_spoofing"] = vniDetails.AllowIPSpoofing
	modelMap["auto_delete"] = vniDetails.AutoDelete
	modelMap["enable_infrastructure_nat"] = vniDetails.EnableInfrastructureNat
	modelMap["resource_group"] = vniDetails.ResourceGroup.ID
	modelMap["protocol_state_filtering_mode"] = vniDetails.ProtocolStateFilteringMode
	primaryipId := *vniDetails.PrimaryIP.ID
	if !core.IsNil(vniDetails.Ips) {
		ips := []map[string]interface{}{}
		for _, ipsItem := range vniDetails.Ips {
			if *ipsItem.ID != primaryipId {
				ipsItemMap, err := resourceIBMIsVirtualNetworkInterfaceReservedIPReferenceToMap(&ipsItem, true)
				if err != nil {
					return nil, err
				}
				ips = append(ips, ipsItemMap)
			}
		}
		modelMap["ips"] = ips
	}
	primaryIPMap, err := resourceIBMIsBareMetalServerReservedIPReferenceToMap(vniDetails.PrimaryIP)
	if err != nil {
		return modelMap, err
	}
	modelMap["primary_ip"] = []map[string]interface{}{primaryIPMap}

	if !core.IsNil(vniDetails.SecurityGroups) {
		securityGroups := make([]string, 0)
		for _, securityGroupsItem := range vniDetails.SecurityGroups {
			if securityGroupsItem.ID != nil {
				securityGroups = append(securityGroups, *securityGroupsItem.ID)
			}
		}
		modelMap["security_groups"] = securityGroups
	}
	if model.ResourceType != nil {
		modelMap["resource_type"] = *model.ResourceType
	}
	if vniDetails.Subnet != nil {
		modelMap["subnet"] = *vniDetails.Subnet.ID
	}
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}
