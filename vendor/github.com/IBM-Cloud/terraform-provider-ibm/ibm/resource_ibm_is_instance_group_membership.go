// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isInstanceGroupMembership                                = "instance_group_membership"
	isInstanceGroup                                          = "instance_group"
	isInstanceGroupMembershipName                            = "name"
	isInstanceGroupMemershipActionDelete                     = "action_delete"
	isInstanceGroupMemershipDeleteInstanceOnMembershipDelete = "delete_instance_on_membership_delete"
	isInstanceGroupMemershipInstance                         = "instance"
	isInstanceGroupMemershipInstanceName                     = "name"
	isInstanceGroupMemershipInstanceTemplate                 = "instance_template"
	isInstanceGroupMemershipInstanceTemplateName             = "name"
	isInstanceGroupMembershipCrn                             = "crn"
	isInstanceGroupMembershipVirtualServerInstance           = "virtual_server_instance"
	isInstanceGroupMembershipLoadBalancerPoolMember          = "load_balancer_pool_member"
	isInstanceGroupMembershipStatus                          = "status"
)

func resourceIBMISInstanceGroupMembership() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISInstanceGroupMembershipUpdate,
		Read:     resourceIBMISInstanceGroupMembershipRead,
		Update:   resourceIBMISInstanceGroupMembershipUpdate,
		Delete:   resourceIBMISInstanceGroupMembershipDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{

			isInstanceGroup: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: InvokeValidator("ibm_is_instance_group_membership", isInstanceGroup),
				Description:  "The instance group identifier.",
			},
			isInstanceGroupMembership: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: InvokeValidator("ibm_is_instance_group_membership", isInstanceGroupMembership),
				Description:  "The unique identifier for this instance group membership.",
			},
			isInstanceGroupMembershipName: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: InvokeValidator("ibm_is_instance_group_membership", isInstanceGroupMembershipName),
				Description:  "The user-defined name for this instance group membership. Names must be unique within the instance group.",
			},
			isInstanceGroupMemershipActionDelete: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "The delete flag for this instance group membership. Must be set to true to delete instance group membership.",
			},
			isInstanceGroupMemershipDeleteInstanceOnMembershipDelete: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set to true, when deleting the membership the instance will also be deleted.",
			},
			isInstanceGroupMemershipInstance: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceGroupMembershipCrn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this virtual server instance.",
						},
						isInstanceGroupMembershipVirtualServerInstance: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this virtual server instance.",
						},
						isInstanceGroupMemershipInstanceName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this virtual server instance (and default system hostname).",
						},
					},
				},
			},
			isInstanceGroupMemershipInstanceTemplate: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceGroupMembershipCrn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this instance template.",
						},
						isInstanceGroupMemershipInstanceTemplate: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this instance template.",
						},
						isInstanceGroupMemershipInstanceTemplateName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique user-defined name for this instance template.",
						},
					},
				},
			},
			isInstanceGroupMembershipLoadBalancerPoolMember: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this load balancer pool member.",
			},
			isInstanceGroupMembershipStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the instance group membership- `deleting`: Membership is deleting dependent resources- `failed`: Membership was unable to maintain dependent resources- `healthy`: Membership is active and serving in the group- `pending`: Membership is waiting for dependent resources- `unhealthy`: Membership has unhealthy dependent resources.",
			},
		},
	}
}

func resourceIBMISInstanceGroupMembershipValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isInstanceGroupMembershipName,
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isInstanceGroup,
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Required:                   true,
			Regexp:                     `^[-0-9a-z_]+$`,
			MinValueLength:             1,
			MaxValueLength:             64})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isInstanceGroupMembership,
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Required:                   true,
			Regexp:                     `^[-0-9a-z_]+$`,
			MinValueLength:             1,
			MaxValueLength:             64})
	ibmISInstanceGroupMembershipResourceValidator := ResourceValidator{ResourceName: "ibm_is_instance_group_membership", Schema: validateSchema}
	return &ibmISInstanceGroupMembershipResourceValidator
}

func resourceIBMISInstanceGroupMembershipUpdate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	instanceGroupID := d.Get(isInstanceGroup).(string)
	instanceGroupMembershipID := d.Get(isInstanceGroupMembership).(string)

	getInstanceGroupMembershipOptions := vpcv1.GetInstanceGroupMembershipOptions{
		ID:              &instanceGroupMembershipID,
		InstanceGroupID: &instanceGroupID,
	}

	instanceGroupMembership, response, err := sess.GetInstanceGroupMembership(&getInstanceGroupMembershipOptions)
	if err != nil || instanceGroupMembership == nil {
		return fmt.Errorf("Error Getting InstanceGroup Membership: %s\n%s", err, response)
	}
	d.SetId(fmt.Sprintf("%s/%s", instanceGroupID, instanceGroupMembershipID))

	if v, ok := d.GetOk(isInstanceGroupMemershipActionDelete); ok {
		actionDelete := v.(bool)
		if actionDelete {
			return resourceIBMISInstanceGroupMembershipDelete(d, meta)
		}
	}

	if v, ok := d.GetOk(isInstanceGroupMembershipName); ok {
		name := v.(string)
		if name != *instanceGroupMembership.Name {

			updateInstanceGroupMembershipOptions := vpcv1.UpdateInstanceGroupMembershipOptions{}
			instanceGroupMembershipPatchModel := &vpcv1.InstanceGroupMembershipPatch{}
			instanceGroupMembershipPatchModel.Name = &name

			updateInstanceGroupMembershipOptions.ID = &instanceGroupMembershipID
			updateInstanceGroupMembershipOptions.InstanceGroupID = &instanceGroupID
			instanceGroupMembershipPatch, err := instanceGroupMembershipPatchModel.AsPatch()
			if err != nil {
				return fmt.Errorf("Error calling asPatch for InstanceGroupMembershipPatch: %s", err)
			}
			updateInstanceGroupMembershipOptions.InstanceGroupMembershipPatch = instanceGroupMembershipPatch
			_, response, err := sess.UpdateInstanceGroupMembership(&updateInstanceGroupMembershipOptions)
			if err != nil {
				return fmt.Errorf("Error updating InstanceGroup Membership: %s\n%s", err, response)
			}
		}
	}
	return resourceIBMISInstanceGroupMembershipRead(d, meta)
}

func resourceIBMISInstanceGroupMembershipRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	instanceGroupID := parts[0]
	instanceGroupMembershipID := parts[1]

	getInstanceGroupMembershipOptions := vpcv1.GetInstanceGroupMembershipOptions{
		ID:              &instanceGroupMembershipID,
		InstanceGroupID: &instanceGroupID,
	}
	instanceGroupMembership, response, err := sess.GetInstanceGroupMembership(&getInstanceGroupMembershipOptions)
	if err != nil || instanceGroupMembership == nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting InstanceGroup Membership: %s\n%s", err, response)
	}
	d.Set(isInstanceGroupMemershipDeleteInstanceOnMembershipDelete, *instanceGroupMembership.DeleteInstanceOnMembershipDelete)
	d.Set(isInstanceGroupMembership, *instanceGroupMembership.ID)
	d.Set(isInstanceGroupMembershipStatus, *instanceGroupMembership.Status)

	instances := make([]map[string]interface{}, 0)
	if instanceGroupMembership.Instance != nil {
		instance := map[string]interface{}{
			isInstanceGroupMembershipCrn:                   *instanceGroupMembership.Instance.CRN,
			isInstanceGroupMembershipVirtualServerInstance: *instanceGroupMembership.Instance.ID,
			isInstanceGroupMemershipInstanceName:           *instanceGroupMembership.Instance.Name,
		}
		instances = append(instances, instance)
	}
	d.Set(isInstanceGroupMemershipInstance, instances)

	instance_templates := make([]map[string]interface{}, 0)
	if instanceGroupMembership.InstanceTemplate != nil {
		instance_template := map[string]interface{}{
			isInstanceGroupMembershipCrn:                 *instanceGroupMembership.InstanceTemplate.CRN,
			isInstanceGroupMemershipInstanceTemplate:     *instanceGroupMembership.InstanceTemplate.ID,
			isInstanceGroupMemershipInstanceTemplateName: *instanceGroupMembership.InstanceTemplate.Name,
		}
		instance_templates = append(instance_templates, instance_template)
	}
	d.Set(isInstanceGroupMemershipInstanceTemplate, instance_templates)

	if instanceGroupMembership.PoolMember != nil && instanceGroupMembership.PoolMember.ID != nil {
		d.Set(isInstanceGroupMembershipLoadBalancerPoolMember, *instanceGroupMembership.PoolMember.ID)
	}
	return nil
}

func resourceIBMISInstanceGroupMembershipDelete(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	instanceGroupID := parts[0]
	instanceGroupMembershipID := parts[1]

	deleteInstanceGroupMembershipOptions := vpcv1.DeleteInstanceGroupMembershipOptions{
		ID:              &instanceGroupMembershipID,
		InstanceGroupID: &instanceGroupID,
	}
	response, err := sess.DeleteInstanceGroupMembership(&deleteInstanceGroupMembershipOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Deleting the InstanceGroup Membership: %s\n%s", err, response)
	}
	return nil
}
