// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-beta-go-sdk/vpcbetav1"
)

func ResourceIBMIsShareMountTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsShareMountTargetCreate,
		ReadContext:   resourceIBMIsShareMountTargetRead,
		UpdateContext: resourceIBMIsShareMountTargetUpdate,
		DeleteContext: resourceIBMIsShareMountTargetDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"share": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The file share identifier.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_share_mount_target", "name"),
				Description:  "The user-defined name for this share target. Names must be unique within the share the share target resides in. If unspecified, the name will be a hyphenated list of randomly-selected words.",
			},
			"transit_encryption": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The transit encryption mode.",
			},
			"access_control_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The access control mode for the share",
			},
			"virtual_network_interface": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				MinItems:      1,
				MaxItems:      1,
				ConflictsWith: []string{"vpc"},
				ExactlyOneOf:  []string{"virtual_network_interface", "vpc"},
				Description:   "VNI for mount target.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "href of virtual network interface",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of this VNI",
						},
						"crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CRN of this VNI",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of this VNI",
						},
						"primary_ip": {
							Type:        schema.TypeList,
							MinItems:    0,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "VNI for mount target.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"reserved_ip": {
										Type:          schema.TypeString,
										Optional:      true,
										ForceNew:      true,
										Computed:      true,
										ConflictsWith: []string{"virtual_network_interface.0.primary_ip.0.name", "virtual_network_interface.0.primary_ip.0.address", "virtual_network_interface.0.primary_ip.0.auto_delete"},
										AtLeastOneOf:  []string{"virtual_network_interface.0.primary_ip.0.reserved_ip", "virtual_network_interface.0.primary_ip.0.name", "virtual_network_interface.0.primary_ip.0.address", "virtual_network_interface.0.primary_ip.0.auto_delete"},
										Description:   "ID of reserved IP",
									},
									"address": {
										Type:          schema.TypeString,
										Optional:      true,
										Computed:      true,
										ForceNew:      true,
										ConflictsWith: []string{"virtual_network_interface.0.primary_ip.0.reserved_ip"},
										AtLeastOneOf:  []string{"virtual_network_interface.0.primary_ip.0.reserved_ip", "virtual_network_interface.0.primary_ip.0.name", "virtual_network_interface.0.primary_ip.0.address", "virtual_network_interface.0.primary_ip.0.auto_delete"},
										Description:   "The IP address to reserve, which must not already be reserved on the subnet.",
									},
									"auto_delete": {
										Type:          schema.TypeBool,
										Optional:      true,
										Computed:      true,
										AtLeastOneOf:  []string{"virtual_network_interface.0.primary_ip.0.reserved_ip", "virtual_network_interface.0.primary_ip.0.name", "virtual_network_interface.0.primary_ip.0.address", "virtual_network_interface.0.primary_ip.0.auto_delete"},
										ConflictsWith: []string{"virtual_network_interface.0.primary_ip.0.reserved_ip"},
										Description:   "Indicates whether this reserved IP member will be automatically deleted when either target is deleted, or the reserved IP is unbound.",
									},
									"name": {
										Type:          schema.TypeString,
										Optional:      true,
										Computed:      true,
										ConflictsWith: []string{"virtual_network_interface.0.primary_ip.0.reserved_ip"},
										AtLeastOneOf:  []string{"virtual_network_interface.0.primary_ip.0.reserved_ip", "virtual_network_interface.0.primary_ip.0.name", "virtual_network_interface.0.primary_ip.0.address", "virtual_network_interface.0.primary_ip.0.auto_delete"},
										Description:   "Name for reserved IP",
									},
									"resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Resource type of primary ip",
									},
									"href": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "href of primary ip",
									},
								},
							},
						},
						"resource_group": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Resource group id",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource type of VNI",
						},
						"security_groups": {
							Type:        schema.TypeSet,
							Computed:    true,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         schema.HashString,
							Description: "The security groups to use for this virtual network interface.",
						},
						"subnet": {
							Type:     schema.TypeString,
							Optional: true,
							//ConflictsWith: []string{"virtual_network_interface.0.primary_ip"},
							Description: "The associated subnet. Required if primary_ip is not specified.",
						},
					},
				},
			},
			"vpc": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"virtual_network_interface"},
				ExactlyOneOf:  []string{"virtual_network_interface", "vpc"},
				Description:   "The unique identifier of the VPC in which instances can mount the file share using this share target.This property will be removed in a future release.The `subnet` property should be used instead.",
			},
			"mount_target": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of this target",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the share target was created.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this share target.",
			},
			"lifecycle_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the mount target.",
			},
			"mount_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The mount path for the share.The IP addresses used in the mount path are currently within the IBM services IP range, but are expected to change to be within one of the VPC's subnets in the future.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of resource referenced.",
			},
		},
	}
}

func ResourceIBMIsShareMountTargetValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 1)
	validateSchema = append(validateSchema,
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

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_share_mount_target", Schema: validateSchema}
	return &resourceValidator
}

func isWaitForOldTargetDelete(context context.Context, vpcClient *vpcbetav1.VpcbetaV1, d *schema.ResourceData, shareid, targetid string) {

	shareTargetOptions := &vpcbetav1.GetShareMountTargetOptions{}

	shareTargetOptions.SetShareID(shareid)
	shareTargetOptions.SetID(targetid)
	for i := 0; i < 6; i++ {
		target, _, err := vpcClient.GetShareMountTargetWithContext(context, shareTargetOptions)
		if err != nil {
			return
		}
		if target != nil && *target.LifecycleState != "deleting" {
			time.Sleep(10 * time.Second)
		} else {
			break
		}
	}

	return
}

func resourceIBMIsShareMountTargetCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1BetaAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	//Temporary code to fix concurrent mount target issues.
	listShareMountTargetOptions := &vpcbetav1.ListShareMountTargetsOptions{}
	shareId := d.Get("share").(string)
	vpcId := d.Get("vpc").(string)
	listShareMountTargetOptions.SetShareID(shareId)

	shareTargets, response, err := vpcClient.ListShareMountTargetsWithContext(context, listShareMountTargetOptions)
	if err != nil || shareTargets == nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] ListShareMountTargetsWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}
	for _, mountTargets := range shareTargets.MountTargets {
		if mountTargets.VPC != nil && *mountTargets.VPC.ID == vpcId {
			isWaitForOldTargetDelete(context, vpcClient, d, shareId, *mountTargets.ID)
			if *mountTargets.LifecycleState == "deleting" {
				_, err = isWaitForTargetDelete(context, vpcClient, d, shareId, *mountTargets.ID)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	createShareMountTargetOptions := &vpcbetav1.CreateShareMountTargetOptions{}

	createShareMountTargetOptions.SetShareID(d.Get("share").(string))
	shareMountTargetPrototype := &vpcbetav1.ShareMountTargetPrototype{}
	if vpcIdIntf, ok := d.GetOk("vpc"); ok {
		vpcId := vpcIdIntf.(string)
		vpc := &vpcbetav1.VPCIdentity{
			ID: &vpcId,
		}
		shareMountTargetPrototype.VPC = vpc
	} else if vniIntf, ok := d.GetOk("virtual_network_interface"); ok {
		vniPrototype := vpcbetav1.ShareMountTargetVirtualNetworkInterfacePrototype{}
		vniMap := vniIntf.([]interface{})[0].(map[string]interface{})
		vniPrototype, err = ShareMountTargetMapToShareMountTargetPrototype(d, vniMap)
		if err != nil {
			return diag.FromErr(err)
		}
		shareMountTargetPrototype.VirtualNetworkInterface = &vniPrototype
	}
	if nameIntf, ok := d.GetOk("name"); ok {
		name := nameIntf.(string)
		shareMountTargetPrototype.Name = &name
	}
	if transitEncryptionIntf, ok := d.GetOk("transit_encryption"); ok {
		transitEncryption := transitEncryptionIntf.(string)
		shareMountTargetPrototype.TransitEncryption = &transitEncryption
	}
	createShareMountTargetOptions.ShareMountTargetPrototype = shareMountTargetPrototype
	shareTarget, response, err := vpcClient.CreateShareMountTargetWithContext(context, createShareMountTargetOptions)
	if err != nil || shareTarget == nil {
		log.Printf("[DEBUG] CreateShareMountTargetWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s/%s", *createShareMountTargetOptions.ShareID, *shareTarget.ID))
	if shareTarget.VirtualNetworkInterface != nil {
		_, err = WaitForVNIAvailable(vpcClient, *shareTarget.VirtualNetworkInterface.ID, d, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	_, err = WaitForMountTargetAvailable(context, vpcClient, *createShareMountTargetOptions.ShareID, *shareTarget.ID, d, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("mount_target", *shareTarget.ID)
	return resourceIBMIsShareMountTargetRead(context, d, meta)
}

func resourceIBMIsShareMountTargetRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1BetaAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	getShareMountTargetOptions := &vpcbetav1.GetShareMountTargetOptions{}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	getShareMountTargetOptions.SetShareID(parts[0])
	getShareMountTargetOptions.SetID(parts[1])

	shareTarget, response, err := vpcClient.GetShareMountTargetWithContext(context, getShareMountTargetOptions)
	if err != nil || shareTarget == nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetShareMountTargetWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}
	if shareTarget.AccessControlMode != nil {
		d.Set("access_control_mode", *shareTarget.AccessControlMode)
	}
	d.Set("mount_target", *shareTarget.ID)
	if shareTarget.VPC != nil && shareTarget.VPC.ID != nil {
		if err = d.Set("vpc", *shareTarget.VPC.ID); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
		}
	}
	if shareTarget.VirtualNetworkInterface != nil {
		vniList, err := ShareMountTargetVirtualNetworkInterfaceToMap(context, vpcClient, d, *shareTarget.VirtualNetworkInterface.ID)
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("virtual_network_interface", vniList)
	}
	if err = d.Set("name", *shareTarget.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	if shareTarget.TransitEncryption != nil {
		if err = d.Set("transit_encryption", *shareTarget.TransitEncryption); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting transit_encryption: %s", err))
		}
	}

	if err = d.Set("created_at", shareTarget.CreatedAt.String()); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("href", shareTarget.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("lifecycle_state", shareTarget.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lifecycle_state: %s", err))
	}
	if err = d.Set("mount_path", shareTarget.MountPath); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting mount_path: %s", err))
	}
	if err = d.Set("resource_type", shareTarget.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}

	return nil
}

func resourceIBMIsShareMountTargetUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1BetaAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	updateShareMountTargetOptions := &vpcbetav1.UpdateShareMountTargetOptions{}

	parts, err := flex.IdParts(d.Id())
	shareId := parts[0]
	mountTargetId := parts[1]
	if err != nil {
		return diag.FromErr(err)
	}

	updateShareMountTargetOptions.SetShareID(shareId)
	updateShareMountTargetOptions.SetID(mountTargetId)

	hasChange := false

	shareTargetPatchModel := &vpcbetav1.ShareMountTargetPatch{}

	if d.HasChange("name") {
		name := d.Get("name").(string)
		shareTargetPatchModel.Name = &name
		hasChange = true
	}

	if d.HasChange("virtual_network_interface.0.name") {
		vniName := d.Get("virtual_network_interface.0.name").(string)
		vniPatchModel := &vpcbetav1.VirtualNetworkInterfacePatch{
			Name: &vniName,
		}
		vniPatch, err := vniPatchModel.AsPatch()
		if err != nil {
			log.Printf("[DEBUG] Virtual network interface AsPatch failed %s", err)
			return diag.FromErr(err)
		}
		shareTargetOptions := &vpcbetav1.GetShareMountTargetOptions{}

		shareTargetOptions.SetShareID(shareId)
		shareTargetOptions.SetID(mountTargetId)
		shareTarget, _, err := vpcClient.GetShareMountTargetWithContext(context, shareTargetOptions)
		if err != nil {
			diag.FromErr(err)
		}
		vniId := *shareTarget.VirtualNetworkInterface.ID
		updateVNIOptions := &vpcbetav1.UpdateVirtualNetworkInterfaceOptions{
			ID:                           &vniId,
			VirtualNetworkInterfacePatch: vniPatch,
		}
		_, response, err := vpcClient.UpdateVirtualNetworkInterfaceWithContext(context, updateVNIOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateShareTargetWithContext failed %s\n%s", err, response)
			return diag.FromErr(err)
		}
		_, err = WaitForVNIAvailable(vpcClient, *shareTarget.VirtualNetworkInterface.ID, d, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("virtual_network_interface.0.security_groups") && !d.IsNewResource() {
		ovs, nvs := d.GetChange("virtual_network_interface.0.security_groups")
		ov := ovs.(*schema.Set)
		nv := nvs.(*schema.Set)
		remove := flex.ExpandStringList(ov.Difference(nv).List())
		add := flex.ExpandStringList(nv.Difference(ov).List())
		networkID := d.Get("virtual_network_interface.0.id").(string)
		if len(add) > 0 {

			for i := range add {
				createsgnicoptions := &vpcbetav1.CreateSecurityGroupTargetBindingOptions{
					SecurityGroupID: &add[i],
					ID:              &networkID,
				}
				_, response, err := vpcClient.CreateSecurityGroupTargetBinding(createsgnicoptions)
				if err != nil {
					return diag.FromErr(fmt.Errorf("[ERROR] Error while creating security group %q for virtual network interface of share mount target %s\n%s: %q", add[i], d.Id(), err, response))
				}
				_, err = WaitForVNIAvailable(vpcClient, networkID, d, d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return diag.FromErr(err)
				}

				_, err = WaitForTargetAvailable(context, vpcClient, shareId, mountTargetId, d, d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return diag.FromErr(err)
				}
			}

		}
		if len(remove) > 0 {
			for i := range remove {
				deletesgnicoptions := &vpcbetav1.DeleteSecurityGroupTargetBindingOptions{
					SecurityGroupID: &remove[i],
					ID:              &networkID,
				}
				response, err := vpcClient.DeleteSecurityGroupTargetBinding(deletesgnicoptions)
				if err != nil {
					return diag.FromErr(fmt.Errorf("[ERROR] Error while removing security group %q for virtual network interface of share mount target %s\n%s: %q", remove[i], d.Id(), err, response))
				}
				_, err = WaitForVNIAvailable(vpcClient, networkID, d, d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return diag.FromErr(err)
				}

				_, err = WaitForTargetAvailable(context, vpcClient, shareId, mountTargetId, d, d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	if !d.IsNewResource() && (d.HasChange("virtual_network_interface.0.primary_ip.0.name") || d.HasChange("virtual_network_interface.0.primary_ip.0.auto_delete")) {
		sess, err := meta.(conns.ClientSession).VpcV1API()
		if err != nil {
			return diag.FromErr(err)
		}
		subnetId := d.Get("virtual_network_interface.0.subnet").(string)
		ripId := d.Get("virtual_network_interface.0.primary_ip.0.reserved_ip").(string)
		updateripoptions := &vpcbetav1.UpdateSubnetReservedIPOptions{
			SubnetID: &subnetId,
			ID:       &ripId,
		}
		reservedIpPath := &vpcbetav1.ReservedIPPatch{}
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
			return diag.FromErr(fmt.Errorf("[ERROR] Error calling reserved ip as patch \n%s", err))
		}
		updateripoptions.ReservedIPPatch = reservedIpPathAsPatch
		_, response, err := vpcClient.UpdateSubnetReservedIP(updateripoptions)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error updating instance network interface reserved ip(%s): %s\n%s", ripId, err, response))
		}
		_, err = isWaitForReservedIpAvailable(sess, subnetId, ripId, d.Timeout(schema.TimeoutUpdate), d)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error waiting for the reserved IP to be available: %s", err))
		}
	}

	if hasChange {
		shareTargetPatch, err := shareTargetPatchModel.AsPatch()
		if err != nil {
			log.Printf("[DEBUG] ShareMountTargetPatch AsPatch failed %s", err)
			return diag.FromErr(err)
		}
		updateShareMountTargetOptions.SetShareMountTargetPatch(shareTargetPatch)
		_, response, err := vpcClient.UpdateShareMountTargetWithContext(context, updateShareMountTargetOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateShareMountTargetWithContext failed %s\n%s", err, response)
			return diag.FromErr(err)
		}
		_, err = WaitForMountTargetAvailable(context, vpcClient, shareId, mountTargetId, d, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceIBMIsShareMountTargetRead(context, d, meta)
}

func resourceIBMIsShareMountTargetDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1BetaAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteShareMountTargetOptions := &vpcbetav1.DeleteShareMountTargetOptions{}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	deleteShareMountTargetOptions.SetShareID(parts[0])
	deleteShareMountTargetOptions.SetID(parts[1])

	_, response, err := vpcClient.DeleteShareMountTargetWithContext(context, deleteShareMountTargetOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteShareMountTargetWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}
	_, err = isWaitForMountTargetDelete(context, vpcClient, d, parts[0], parts[1])
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func WaitForMountTargetAvailable(context context.Context, vpcClient *vpcbetav1.VpcbetaV1, shareid, targetid string, d *schema.ResourceData, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for target (%s) to be available.", targetid)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"updating", "pending", "waiting"},
		Target:     []string{"stable", "failed"},
		Refresh:    mountTargetRefresh(context, vpcClient, shareid, targetid, d),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func mountTargetRefresh(context context.Context, vpcClient *vpcbetav1.VpcbetaV1, shareid, targetid string, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		shareTargetOptions := &vpcbetav1.GetShareMountTargetOptions{}

		shareTargetOptions.SetShareID(shareid)
		shareTargetOptions.SetID(targetid)

		target, response, err := vpcClient.GetShareMountTargetWithContext(context, shareTargetOptions)
		if err != nil {
			return nil, "", fmt.Errorf("Error Getting target: %s\n%s", err, response)
		}
		d.Set("lifecycle_state", *target.LifecycleState)
		if *target.LifecycleState == "stable" || *target.LifecycleState == "failed" {

			return target, *target.LifecycleState, nil

		}
		return target, "pending", nil
	}
}

func isWaitForMountTargetDelete(context context.Context, vpcClient *vpcbetav1.VpcbetaV1, d *schema.ResourceData, shareid, targetid string) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending: []string{"deleting", "stable"},
		Target:  []string{"done"},
		Refresh: func() (interface{}, string, error) {
			shareTargetOptions := &vpcbetav1.GetShareMountTargetOptions{}

			shareTargetOptions.SetShareID(shareid)
			shareTargetOptions.SetID(targetid)

			target, response, err := vpcClient.GetShareMountTargetWithContext(context, shareTargetOptions)
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					return target, "done", nil
				}
				return nil, "", fmt.Errorf("Error Getting Target: %s\n%s", err, response)
			}
			if *target.LifecycleState == isInstanceFailed {
				return target, *target.LifecycleState, fmt.Errorf("The  target %s failed to delete: %v", targetid, err)
			}
			return target, "deleting", nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}
