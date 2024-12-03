// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func ResourceIBMIsInstanceNetworkAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsInstanceNetworkAttachmentCreate,
		ReadContext:   resourceIBMIsInstanceNetworkAttachmentRead,
		UpdateContext: resourceIBMIsInstanceNetworkAttachmentUpdate,
		DeleteContext: resourceIBMIsInstanceNetworkAttachmentDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"instance": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_instance_network_attachment", "instance"),
				Description:  "The virtual server instance identifier.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_instance_network_attachment", "name"),
				Description:  "The name for this instance network attachment. The name is unique across all network attachments for the instance.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the instance network attachment was created.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this instance network attachment.",
			},
			"lifecycle_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the instance network attachment.",
			},
			"port_speed": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The port speed for this instance network attachment in Mbps.",
			},

			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The instance network attachment type.",
			},
			"network_attachment": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this instance network attachment.",
			},
			// vni properties
			"virtual_network_interface": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "A virtual network interface for the instance network attachment. This can be specified using an existing virtual network interface, or a prototype object for a new virtual network interface.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"virtual_network_interface.0.allow_ip_spoofing", "virtual_network_interface.0.auto_delete", "virtual_network_interface.0.enable_infrastructure_nat", "virtual_network_interface.0.ips", "virtual_network_interface.0.name", "virtual_network_interface.0.primary_ip", "virtual_network_interface.0.resource_group", "virtual_network_interface.0.security_groups", "virtual_network_interface.0.security_groups"},
							Description:   "The virtual network interface id for this instance network attachment.",
						},
						"allow_ip_spoofing": &schema.Schema{
							Type:          schema.TypeBool,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"virtual_network_interface.0.id"},
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
							ValidateFunc:  validate.InvokeValidator("ibm_is_virtual_network_interface", "name"),
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
							ConflictsWith: []string{"virtual_network_interface.0.id"},
							Computed:      true,
							Description:   "The primary IP address of the virtual network interface for the instance networkattachment.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The IP address.If the address has not yet been selected, the value will be `0.0.0.0`.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.",
									},
									"auto_delete": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     true,
										Description: "Indicates whether this primary_ip will be automatically deleted when `vni` is deleted.",
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
							ForceNew:      true,
							ConflictsWith: []string{"virtual_network_interface.0.id"},
							Elem:          &schema.Schema{Type: schema.TypeString},
							Set:           schema.HashString,
							Description:   "The security groups for this virtual network interface.",
						},
						"subnet": &schema.Schema{
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"virtual_network_interface.0.id"},
							ForceNew:      true,
							Description:   "The associated subnet id.",
						},
					},
				},
			},
		},
	}
}

func ResourceIBMIsInstanceNetworkAttachmentValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "instance",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[-0-9a-z_]+$`,
			MinValueLength:             1,
			MaxValueLength:             64,
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
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_instance_network_attachment", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMIsInstanceNetworkAttachmentCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	createInstanceNetworkAttachmentOptions := &vpcv1.CreateInstanceNetworkAttachmentOptions{}

	createInstanceNetworkAttachmentOptions.SetInstanceID(d.Get("instance").(string))
	virtualNetworkInterfaceModel, err := resourceIBMIsInstanceNetworkAttachmentMapToInstanceNetworkAttachmentPrototypeVirtualNetworkInterface(d.Get("virtual_network_interface.0").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	createInstanceNetworkAttachmentOptions.SetVirtualNetworkInterface(virtualNetworkInterfaceModel)
	if _, ok := d.GetOk("name"); ok {
		createInstanceNetworkAttachmentOptions.SetName(d.Get("name").(string))
	}

	instanceNetworkAttachment, response, err := vpcClient.CreateInstanceNetworkAttachmentWithContext(context, createInstanceNetworkAttachmentOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateInstanceNetworkAttachmentWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateInstanceNetworkAttachmentWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *createInstanceNetworkAttachmentOptions.InstanceID, *instanceNetworkAttachment.ID))
	_, err = isWaitForInstanceNetworkAttachmentStable(vpcClient, *createInstanceNetworkAttachmentOptions.InstanceID, *instanceNetworkAttachment.ID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(fmt.Errorf("isWaitForInstanceNetworkAttachmentStable failed %s", err))
	}
	return resourceIBMIsInstanceNetworkAttachmentRead(context, d, meta)
}

func resourceIBMIsInstanceNetworkAttachmentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	getInstanceNetworkAttachmentOptions := &vpcv1.GetInstanceNetworkAttachmentOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getInstanceNetworkAttachmentOptions.SetInstanceID(parts[0])
	getInstanceNetworkAttachmentOptions.SetID(parts[1])

	instanceNetworkAttachment, response, err := vpcClient.GetInstanceNetworkAttachmentWithContext(context, getInstanceNetworkAttachmentOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetInstanceNetworkAttachmentWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetInstanceNetworkAttachmentWithContext failed %s\n%s", err, response))
	}
	// attachment details
	if !core.IsNil(instanceNetworkAttachment.Name) {
		if err = d.Set("name", instanceNetworkAttachment.Name); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
		}
	}
	if err = d.Set("created_at", flex.DateTimeToString(instanceNetworkAttachment.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
	}

	if err = d.Set("href", instanceNetworkAttachment.Href); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting href: %s", err))
	}
	if err = d.Set("lifecycle_state", instanceNetworkAttachment.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting lifecycle_state: %s", err))
	}
	if err = d.Set("port_speed", flex.IntValue(instanceNetworkAttachment.PortSpeed)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting port_speed: %s", err))
	}
	if err = d.Set("resource_type", instanceNetworkAttachment.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_type: %s", err))
	}
	if err = d.Set("type", instanceNetworkAttachment.Type); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting type: %s", err))
	}
	if err = d.Set("network_attachment", instanceNetworkAttachment.ID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting network_attachment: %s", err))
	}
	// vni details
	vniId := *instanceNetworkAttachment.VirtualNetworkInterface.ID
	vniMap := make(map[string]interface{})
	vniMap["id"] = vniId
	getVniOptions := &vpcv1.GetVirtualNetworkInterfaceOptions{
		ID: &vniId,
	}
	vniDetails, response, err := vpcClient.GetVirtualNetworkInterface(getVniOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetVirtualNetworkInterface failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetVirtualNetworkInterface failed %s\n%s", err, response))
	}
	vniMap["allow_ip_spoofing"] = vniDetails.AllowIPSpoofing
	vniMap["auto_delete"] = vniDetails.AutoDelete
	vniMap["enable_infrastructure_nat"] = vniDetails.EnableInfrastructureNat
	vniMap["name"] = vniDetails.Name
	vniMap["resource_group"] = vniDetails.ResourceGroup.ID
	vniMap["resource_type"] = vniDetails.ResourceType
	vniMap["protocol_state_filtering_mode"] = vniDetails.ProtocolStateFilteringMode
	primaryipId := *instanceNetworkAttachment.PrimaryIP.ID
	if !core.IsNil(vniDetails.Ips) {
		ips := []map[string]interface{}{}
		for _, ipsItem := range vniDetails.Ips {
			if *ipsItem.ID != primaryipId {
				ipsItemMap, err := resourceIBMIsVirtualNetworkInterfaceReservedIPReferenceToMap(&ipsItem, true)
				if err != nil {
					return diag.FromErr(err)
				}
				ips = append(ips, ipsItemMap)
			}
		}
		vniMap["ips"] = ips
	}

	if !core.IsNil(vniDetails.SecurityGroups) {
		securityGroups := make([]string, 0)
		for _, securityGroupsItem := range vniDetails.SecurityGroups {
			if securityGroupsItem.ID != nil {
				securityGroups = append(securityGroups, *securityGroupsItem.ID)
			}
		}
		vniMap["security_groups"] = securityGroups
	}
	autoDelete := true
	if autoDeleteOk, ok := d.GetOkExists("virtual_network_interface.0.primary_ip.0.auto_delete"); ok {
		autoDelete = autoDeleteOk.(bool)
	}
	primaryIPMap, err := resourceIBMIsInstanceNetworkAttachmentReservedIPReferenceToMap(instanceNetworkAttachment.PrimaryIP, autoDelete)
	if err != nil {
		return diag.FromErr(err)
	}
	vniMap["primary_ip"] = []map[string]interface{}{primaryIPMap}

	vniMap["subnet"] = *instanceNetworkAttachment.Subnet.ID
	if err = d.Set("virtual_network_interface", []map[string]interface{}{vniMap}); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting virtual_network_interface: %s", err))
	}

	return nil
}

func resourceIBMIsInstanceNetworkAttachmentUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	updateInstanceNetworkAttachmentOptions := &vpcv1.UpdateInstanceNetworkAttachmentOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	updateInstanceNetworkAttachmentOptions.SetInstanceID(parts[0])
	updateInstanceNetworkAttachmentOptions.SetID(parts[1])

	hasChange := false

	patchVals := &vpcv1.InstanceNetworkAttachmentPatch{}
	if d.HasChange("instance") {
		return diag.FromErr(fmt.Errorf("[ERROR] Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "instance"))
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
		if d.HasChange("virtual_network_interface.0.protocol_state_filtering_mode") {
			pStateFilteringMode := d.Get("virtual_network_interface.0.protocol_state_filtering_mode").(string)
			virtualNetworkInterfacePatch.ProtocolStateFilteringMode = &pStateFilteringMode
		}
		virtualNetworkInterfacePatchAsPatch, err := virtualNetworkInterfacePatch.AsPatch()
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error encountered while apply as patch for virtualNetworkInterfacePatch of instance(%s) vni (%s) %s", d.Id(), vniId, err))
		}
		updateVirtualNetworkInterfaceOptions.VirtualNetworkInterfacePatch = virtualNetworkInterfacePatchAsPatch
		_, response, err := vpcClient.UpdateVirtualNetworkInterfaceWithContext(context, updateVirtualNetworkInterfaceOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateVirtualNetworkInterfaceWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateVirtualNetworkInterfaceWithContext failed during instance(%s) network attachment patch %s\n%s", d.Id(), err, response))
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
							log.Printf("[DEBUG] AddVirtualNetworkInterfaceIPWithContext failed in VirtualNetworkInterface patch during instance nac patch %s\n%s", err, response)
							return diag.FromErr(fmt.Errorf("AddVirtualNetworkInterfaceIPWithContext failed in VirtualNetworkInterface patch during instance nac patch %s\n%s", err, response))
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
							log.Printf("[DEBUG] RemoveVirtualNetworkInterfaceIPWithContext failed in VirtualNetworkInterface patch during instance nac patch %s\n%s", err, response)
							return diag.FromErr(fmt.Errorf("RemoveVirtualNetworkInterfaceIPWithContext failed in VirtualNetworkInterface patch during instance nac patch %s\n%s", err, response))
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

	if d.HasChange("name") {
		newName := d.Get("name").(string)
		patchVals.Name = &newName
		hasChange = true
	}

	if hasChange {
		updateInstanceNetworkAttachmentOptions.InstanceNetworkAttachmentPatch, _ = patchVals.AsPatch()
		_, response, err := vpcClient.UpdateInstanceNetworkAttachmentWithContext(context, updateInstanceNetworkAttachmentOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateInstanceNetworkAttachmentWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateInstanceNetworkAttachmentWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIBMIsInstanceNetworkAttachmentRead(context, d, meta)
}

func resourceIBMIsInstanceNetworkAttachmentDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}
	getInstanceNetworkAttachmentOptions := &vpcv1.GetInstanceNetworkAttachmentOptions{}
	getInstanceNetworkAttachmentOptions.SetInstanceID(parts[0])
	getInstanceNetworkAttachmentOptions.SetID(parts[1])

	ina, response, err := vpcClient.GetInstanceNetworkAttachmentWithContext(context, getInstanceNetworkAttachmentOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetInstanceNetworkAttachmentWithContext failed while deleting %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetInstanceNetworkAttachmentWithContext failed %s\n%s", err, response))
	}

	deleteInstanceNetworkAttachmentOptions := &vpcv1.DeleteInstanceNetworkAttachmentOptions{}
	deleteInstanceNetworkAttachmentOptions.SetInstanceID(parts[0])
	deleteInstanceNetworkAttachmentOptions.SetID(parts[1])

	response, err = vpcClient.DeleteInstanceNetworkAttachmentWithContext(context, deleteInstanceNetworkAttachmentOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteInstanceNetworkAttachmentWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteInstanceNetworkAttachmentWithContext failed %s\n%s", err, response))
	}
	_, err = isWaitForInstanceNetworkAttachmentDeleted(vpcClient, parts[0], parts[1], ina, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(fmt.Errorf("isWaitForInstanceNetworkAttachmentDeleted failed %s", err))
	}

	d.SetId("")

	return nil
}

func resourceIBMIsInstanceNetworkAttachmentMapToInstanceNetworkAttachmentPrototypeVirtualNetworkInterface(modelMap map[string]interface{}) (vpcv1.InstanceNetworkAttachmentPrototypeVirtualNetworkInterfaceIntf, error) {
	model := &vpcv1.InstanceNetworkAttachmentPrototypeVirtualNetworkInterface{}
	if modelMap["allow_ip_spoofing"] != nil {
		model.AllowIPSpoofing = core.BoolPtr(modelMap["allow_ip_spoofing"].(bool))
	}
	if modelMap["auto_delete"] != nil {
		model.AutoDelete = core.BoolPtr(modelMap["auto_delete"].(bool))
	}
	if modelMap["enable_infrastructure_nat"] != nil {
		model.EnableInfrastructureNat = core.BoolPtr(modelMap["enable_infrastructure_nat"].(bool))
	}
	if modelMap["ips"] != nil {
		ips := []vpcv1.VirtualNetworkInterfaceIPPrototypeIntf{}
		for _, ipsItem := range modelMap["ips"].(*schema.Set).List() {
			ipsItemModel, err := resourceIBMIsInstanceNetworkAttachmentMapToVirtualNetworkInterfaceIPPrototype(ipsItem.(map[string]interface{}))
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
		PrimaryIPModel, err := resourceIBMIsInstanceNetworkAttachmentMapToVirtualNetworkInterfacePrimaryIPPrototype(modelMap["primary_ip"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.PrimaryIP = PrimaryIPModel
	}
	if modelMap["resource_group"] != nil && modelMap["resource_group"].(string) != "" {
		rgId := modelMap["resource_group"].(string)
		model.ResourceGroup = &vpcv1.ResourceGroupIdentity{
			ID: &rgId,
		}
	}
	if modelMap["security_groups"] != nil && modelMap["security_groups"].(*schema.Set).Len() > 0 {
		securityGroups := []vpcv1.SecurityGroupIdentityIntf{}
		sg := modelMap["security_groups"].(*schema.Set)
		for _, v := range sg.List() {
			value := v.(string)
			securityGroupsItem := &vpcv1.SecurityGroupIdentity{
				ID: &value,
			}
			securityGroups = append(securityGroups, securityGroupsItem)
		}
		model.SecurityGroups = securityGroups
	}
	if modelMap["subnet"] != nil && modelMap["subnet"].(string) != "" {
		subnetid := modelMap["subnet"].(string)
		model.Subnet = &vpcv1.SubnetIdentity{
			ID: &subnetid,
		}
	}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	return model, nil
}

func resourceIBMIsInstanceNetworkAttachmentMapToVirtualNetworkInterfaceIPPrototype(modelMap map[string]interface{}) (vpcv1.VirtualNetworkInterfaceIPPrototypeIntf, error) {
	model := &vpcv1.VirtualNetworkInterfaceIPPrototype{}
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

func resourceIBMIsInstanceNetworkAttachmentMapToVirtualNetworkInterfaceIPPrototypeReservedIPIdentityVirtualNetworkInterfaceIPsContext(modelMap map[string]interface{}) (vpcv1.VirtualNetworkInterfaceIPPrototypeReservedIPIdentityVirtualNetworkInterfaceIPsContextIntf, error) {
	model := &vpcv1.VirtualNetworkInterfaceIPPrototypeReservedIPIdentityVirtualNetworkInterfaceIPsContext{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}

func resourceIBMIsInstanceNetworkAttachmentMapToVirtualNetworkInterfaceIPPrototypeReservedIPIdentityVirtualNetworkInterfaceIPsContextByID(modelMap map[string]interface{}) (*vpcv1.VirtualNetworkInterfaceIPPrototypeReservedIPIdentityVirtualNetworkInterfaceIPsContextByID, error) {
	model := &vpcv1.VirtualNetworkInterfaceIPPrototypeReservedIPIdentityVirtualNetworkInterfaceIPsContextByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func resourceIBMIsInstanceNetworkAttachmentMapToVirtualNetworkInterfaceIPPrototypeReservedIPIdentityVirtualNetworkInterfaceIPsContextByHref(modelMap map[string]interface{}) (*vpcv1.VirtualNetworkInterfaceIPPrototypeReservedIPIdentityVirtualNetworkInterfaceIPsContextByHref, error) {
	model := &vpcv1.VirtualNetworkInterfaceIPPrototypeReservedIPIdentityVirtualNetworkInterfaceIPsContextByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func resourceIBMIsInstanceNetworkAttachmentMapToVirtualNetworkInterfaceIPPrototypeReservedIPPrototypeVirtualNetworkInterfaceIPsContext(modelMap map[string]interface{}) (*vpcv1.VirtualNetworkInterfaceIPPrototypeReservedIPPrototypeVirtualNetworkInterfaceIPsContext, error) {
	model := &vpcv1.VirtualNetworkInterfaceIPPrototypeReservedIPPrototypeVirtualNetworkInterfaceIPsContext{}
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

func resourceIBMIsInstanceNetworkAttachmentMapToVirtualNetworkInterfacePrimaryIPPrototype(modelMap map[string]interface{}) (vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeIntf, error) {
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

func resourceIBMIsInstanceNetworkAttachmentMapToVirtualNetworkInterfacePrimaryIPPrototypeReservedIPIdentityVirtualNetworkInterfacePrimaryIPContext(modelMap map[string]interface{}) (vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeReservedIPIdentityVirtualNetworkInterfacePrimaryIPContextIntf, error) {
	model := &vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeReservedIPIdentityVirtualNetworkInterfacePrimaryIPContext{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}

func resourceIBMIsInstanceNetworkAttachmentMapToVirtualNetworkInterfacePrimaryIPPrototypeReservedIPIdentityVirtualNetworkInterfacePrimaryIPContextByID(modelMap map[string]interface{}) (*vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeReservedIPIdentityVirtualNetworkInterfacePrimaryIPContextByID, error) {
	model := &vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeReservedIPIdentityVirtualNetworkInterfacePrimaryIPContextByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func resourceIBMIsInstanceNetworkAttachmentMapToVirtualNetworkInterfacePrimaryIPPrototypeReservedIPIdentityVirtualNetworkInterfacePrimaryIPContextByHref(modelMap map[string]interface{}) (*vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeReservedIPIdentityVirtualNetworkInterfacePrimaryIPContextByHref, error) {
	model := &vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeReservedIPIdentityVirtualNetworkInterfacePrimaryIPContextByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func resourceIBMIsInstanceNetworkAttachmentMapToVirtualNetworkInterfacePrimaryIPPrototypeReservedIPPrototypeVirtualNetworkInterfacePrimaryIPContext(modelMap map[string]interface{}) (*vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeReservedIPPrototypeVirtualNetworkInterfacePrimaryIPContext, error) {
	model := &vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeReservedIPPrototypeVirtualNetworkInterfacePrimaryIPContext{}
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
func resourceIBMIsInstanceNetworkAttachmentMapToSubnetIdentityByID(modelMap map[string]interface{}) (*vpcv1.SubnetIdentityByID, error) {
	model := &vpcv1.SubnetIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func resourceIBMIsInstanceNetworkAttachmentMapToSubnetIdentityByCRN(modelMap map[string]interface{}) (*vpcv1.SubnetIdentityByCRN, error) {
	model := &vpcv1.SubnetIdentityByCRN{}
	model.CRN = core.StringPtr(modelMap["crn"].(string))
	return model, nil
}

func resourceIBMIsInstanceNetworkAttachmentMapToSubnetIdentityByHref(modelMap map[string]interface{}) (*vpcv1.SubnetIdentityByHref, error) {
	model := &vpcv1.SubnetIdentityByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func resourceIBMIsInstanceNetworkAttachmentMapToInstanceNetworkAttachmentPrototypeVirtualNetworkInterfaceVirtualNetworkInterfacePrototypeInstanceNetworkAttachmentContext(modelMap map[string]interface{}) (*vpcv1.InstanceNetworkAttachmentPrototypeVirtualNetworkInterfaceVirtualNetworkInterfacePrototypeInstanceNetworkAttachmentContext, error) {
	model := &vpcv1.InstanceNetworkAttachmentPrototypeVirtualNetworkInterfaceVirtualNetworkInterfacePrototypeInstanceNetworkAttachmentContext{}
	if modelMap["allow_ip_spoofing"] != nil {
		model.AllowIPSpoofing = core.BoolPtr(modelMap["allow_ip_spoofing"].(bool))
	}
	if modelMap["auto_delete"] != nil {
		model.AutoDelete = core.BoolPtr(modelMap["auto_delete"].(bool))
	}
	if modelMap["enable_infrastructure_nat"] != nil {
		model.EnableInfrastructureNat = core.BoolPtr(modelMap["enable_infrastructure_nat"].(bool))
	}
	if modelMap["ips"] != nil {
		ips := []vpcv1.VirtualNetworkInterfaceIPPrototypeIntf{}
		for _, ipsItem := range modelMap["ips"].(*schema.Set).List() {
			ipsItemModel, err := resourceIBMIsInstanceNetworkAttachmentMapToVirtualNetworkInterfaceIPPrototype(ipsItem.(map[string]interface{}))
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
	if modelMap["primary_ip"] != nil && len(modelMap["primary_ip"].([]interface{})) > 0 {
		PrimaryIPModel, err := resourceIBMIsInstanceNetworkAttachmentMapToVirtualNetworkInterfacePrimaryIPPrototype(modelMap["primary_ip"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.PrimaryIP = PrimaryIPModel
	}
	if modelMap["resource_group"] != nil && modelMap["resource_group"].(string) != "" {
		rgId := modelMap["resource_group"].(string)
		model.ResourceGroup = &vpcv1.ResourceGroupIdentity{
			ID: &rgId,
		}
	}
	if modelMap["security_groups"] != nil && modelMap["security_groups"].(*schema.Set).Len() > 0 {
		securityGroups := []vpcv1.SecurityGroupIdentityIntf{}
		sg := modelMap["security_groups"].(*schema.Set)
		for _, v := range sg.List() {
			value := v.(string)
			securityGroupsItem := &vpcv1.SecurityGroupIdentity{
				ID: &value,
			}
			securityGroups = append(securityGroups, securityGroupsItem)
		}
		model.SecurityGroups = securityGroups
	}
	if modelMap["subnet"] != nil && modelMap["subnet"].(string) != "" {
		subnetid := modelMap["subnet"].(string)
		model.Subnet = &vpcv1.SubnetIdentity{
			ID: &subnetid,
		}
	}
	return model, nil
}

func resourceIBMIsInstanceNetworkAttachmentReservedIPReferenceToMap(model *vpcv1.ReservedIPReference, autodelete bool) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = model.Address
	if model.Deleted != nil {
		deletedMap, err := resourceIBMIsInstanceNetworkAttachmentReservedIPReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = model.Href
	modelMap["auto_delete"] = autodelete
	modelMap["reserved_ip"] = model.ID
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func resourceIBMIsInstanceNetworkAttachmentReservedIPReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func resourceIBMIsInstanceNetworkAttachmentSubnetReferenceToMap(model *vpcv1.SubnetReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.CRN
	if model.Deleted != nil {
		deletedMap, err := resourceIBMIsInstanceNetworkAttachmentSubnetReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func resourceIBMIsInstanceNetworkAttachmentSubnetReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func isWaitForInstanceNetworkAttachmentStable(instanceC *vpcv1.VpcV1, instanceId, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for instance network attachment (%s) to be stable.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"deleting", "waiting", "updating", "pending"},
		Target:     []string{"stable", "failed", "suspended", ""},
		Refresh:    isInstanceNetworkAttachmentRefreshFunc(instanceC, instanceId, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}
func isInstanceNetworkAttachmentRefreshFunc(instanceC *vpcv1.VpcV1, instanceId, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getInstanceNetworkAttachmentOptions := &vpcv1.GetInstanceNetworkAttachmentOptions{
			InstanceID: &instanceId,
			ID:         &id,
		}
		networkAttachment, response, err := instanceC.GetInstanceNetworkAttachment(getInstanceNetworkAttachmentOptions)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error getting network attachment: %s\n%s", err, response)
		}

		if *networkAttachment.LifecycleState == "failed" || *networkAttachment.LifecycleState == "suspended" {
			return networkAttachment, *networkAttachment.LifecycleState, fmt.Errorf("[ERROR] Error network attachment(%s) in (%s) state", id, *networkAttachment.LifecycleState)
		}

		return networkAttachment, *networkAttachment.LifecycleState, nil
	}
}
func isWaitForInstanceNetworkAttachmentDeleted(instanceC *vpcv1.VpcV1, instanceId, id string, ina *vpcv1.InstanceNetworkAttachment, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for instance network attachment (%s) to be deleted.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"deleting", "waiting", "updating", "pending"},
		Target:     []string{"deleted", "failed", "suspended", ""},
		Refresh:    isInstanceNetworkAttachmentDeleteRefreshFunc(instanceC, instanceId, id, ina),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}
func isInstanceNetworkAttachmentDeleteRefreshFunc(instanceC *vpcv1.VpcV1, instanceId, id string, ina *vpcv1.InstanceNetworkAttachment) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getInstanceNetworkAttachmentOptions := &vpcv1.GetInstanceNetworkAttachmentOptions{
			InstanceID: &instanceId,
			ID:         &id,
		}
		networkAttachment, response, err := instanceC.GetInstanceNetworkAttachment(getInstanceNetworkAttachmentOptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return ina, "deleted", nil
			}
			return ina, "", fmt.Errorf("[ERROR] Error deleting network attachment: %s\n%s", err, response)
		}

		if *networkAttachment.LifecycleState == "failed" || *networkAttachment.LifecycleState == "suspended" {
			return networkAttachment, *networkAttachment.LifecycleState, fmt.Errorf("[ERROR] Error network attachment(%s) in (%s) state", id, *networkAttachment.LifecycleState)
		}

		return networkAttachment, *networkAttachment.LifecycleState, nil
	}
}
