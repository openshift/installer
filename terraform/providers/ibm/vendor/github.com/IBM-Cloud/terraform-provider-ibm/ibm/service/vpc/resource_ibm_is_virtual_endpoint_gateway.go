// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isVirtualEndpointGatewayName               = "name"
	isVirtualEndpointGatewayResourceType       = "resource_type"
	isVirtualEndpointGatewayCRN                = "crn"
	isVirtualEndpointGatewayResourceGroupID    = "resource_group"
	isVirtualEndpointGatewayCreatedAt          = "created_at"
	isVirtualEndpointGatewayIPs                = "ips"
	isVirtualEndpointGatewayIPsID              = "id"
	isVirtualEndpointGatewayIPsAddress         = "address"
	isVirtualEndpointGatewayIPsName            = "name"
	isVirtualEndpointGatewayIPsSubnet          = "subnet"
	isVirtualEndpointGatewayIPsResourceType    = "resource_type"
	isVirtualEndpointGatewayHealthState        = "health_state"
	isVirtualEndpointGatewayLifecycleState     = "lifecycle_state"
	isVirtualEndpointGatewayTarget             = "target"
	isVirtualEndpointGatewayTargetName         = "name"
	isVirtualEndpointGatewayTargetCRN          = "crn"
	isVirtualEndpointGatewayTargetResourceType = "resource_type"
	isVirtualEndpointGatewayVpcID              = "vpc"
	isVirtualEndpointGatewayTags               = "tags"
	isVirtualEndpointGatewaySecurityGroups     = "security_groups"
)

func ResourceIBMISEndpointGateway() *schema.Resource {
	targetNameFmt := fmt.Sprintf("%s.0.%s", isVirtualEndpointGatewayTarget, isVirtualEndpointGatewayTargetName)
	targetCRNFmt := fmt.Sprintf("%s.0.%s", isVirtualEndpointGatewayTarget, isVirtualEndpointGatewayTargetCRN)
	return &schema.Resource{
		Create:   resourceIBMisVirtualEndpointGatewayCreate,
		Read:     resourceIBMisVirtualEndpointGatewayRead,
		Update:   resourceIBMisVirtualEndpointGatewayUpdate,
		Delete:   resourceIBMisVirtualEndpointGatewayDelete,
		Exists:   resourceIBMisVirtualEndpointGatewayExists,
		Importer: &schema.ResourceImporter{},

		CustomizeDiff: customdiff.Sequence(
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return flex.ResourceTagsCustomizeDiff(diff)
			},
		),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			isVirtualEndpointGatewayName: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_virtual_endpoint_gateway", isVirtualEndpointGatewayName),
				Description:  "Endpoint gateway name",
			},
			isVirtualEndpointGatewayResourceType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Endpoint gateway resource type",
			},
			isVirtualEndpointGatewayCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this Endpoint gateway",
			},
			isVirtualEndpointGatewayResourceGroupID: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The resource group id",
			},
			isVirtualEndpointGatewayCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Endpoint gateway created date and time",
			},
			isVirtualEndpointGatewayHealthState: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Endpoint gateway health state",
			},
			isVirtualEndpointGatewayLifecycleState: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Endpoint gateway lifecycle state",
			},
			isVirtualEndpointGatewaySecurityGroups: {
				Type:        schema.TypeSet,
				Computed:    true,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "Endpoint gateway securitygroups list",
			},
			isVirtualEndpointGatewayIPs: {
				Type:             schema.TypeList,
				Optional:         true,
				Computed:         true,
				Description:      "Endpoint gateway IPs",
				DiffSuppressFunc: flex.ApplyOnce,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isVirtualEndpointGatewayIPsID: {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The IPs id",
						},
						isVirtualEndpointGatewayIPsName: {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The IPs name",
						},
						isVirtualEndpointGatewayIPsSubnet: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The Subnet id",
						},
						isVirtualEndpointGatewayIPsResourceType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPE Resource Type",
						},
						isVirtualEndpointGatewayIPsAddress: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP Address",
						},
					},
				},
			},
			isVirtualEndpointGatewayTarget: {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				MaxItems:    1,
				Description: "Endpoint gateway target",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isVirtualEndpointGatewayTargetName: {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							AtLeastOneOf: []string{
								targetNameFmt,
								targetCRNFmt,
							},
							Description: "The target name",
						},
						isVirtualEndpointGatewayTargetCRN: {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							AtLeastOneOf: []string{
								targetNameFmt,
								targetCRNFmt,
							},
							Description: "The target crn",
						},
						isVirtualEndpointGatewayTargetResourceType: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.InvokeValidator("ibm_is_virtual_endpoint_gateway", isVirtualEndpointGatewayTargetResourceType),
							Description:  "The target resource type",
						},
					},
				},
			},
			isVirtualEndpointGatewayVpcID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The VPC id",
			},
			isVirtualEndpointGatewayTags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_virtual_endpoint_gateway", "tag")},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of tags for VPE",
			},
		},
	}
}

func ResourceIBMISEndpointGatewayValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isVirtualEndpointGatewayName,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "tag",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isVirtualEndpointGatewayTargetResourceType,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "provider_cloud_service, provider_infrastructure_service"})

	ibmEndpointGatewayResourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_virtual_endpoint_gateway", Schema: validateSchema}
	return &ibmEndpointGatewayResourceValidator
}

func resourceIBMisVirtualEndpointGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	name := d.Get(isVirtualEndpointGatewayName).(string)

	// target opiton
	targetOpt := &vpcv1.EndpointGatewayTargetPrototype{}
	targetNameFmt := fmt.Sprintf("%s.0.%s", isVirtualEndpointGatewayTarget, isVirtualEndpointGatewayTargetName)
	targetCRNFmt := fmt.Sprintf("%s.0.%s", isVirtualEndpointGatewayTarget, isVirtualEndpointGatewayTargetCRN)
	targetResourceTypeFmt := fmt.Sprintf("%s.0.%s", isVirtualEndpointGatewayTarget, isVirtualEndpointGatewayTargetResourceType)
	targetOpt.ResourceType = core.StringPtr(d.Get(targetResourceTypeFmt).(string))
	if v, ok := d.GetOk(targetNameFmt); ok {
		targetOpt.Name = core.StringPtr(v.(string))
	}
	if v, ok := d.GetOk(targetCRNFmt); ok {
		targetOpt.CRN = core.StringPtr(v.(string))
	}

	// vpc option
	vpcID := d.Get(isVirtualEndpointGatewayVpcID).(string)
	vpcOpt := &vpcv1.VPCIdentity{
		ID: core.StringPtr(vpcID),
	}

	// update option
	opt := sess.NewCreateEndpointGatewayOptions(targetOpt, vpcOpt)
	opt.SetName(name)
	opt.SetTarget(targetOpt)
	opt.SetVPC(vpcOpt)

	// IPs option
	if ips, ok := d.GetOk(isVirtualEndpointGatewayIPs); ok {
		opt.SetIps(expandIPs(ips.([]interface{})))
	}
	// Security group option
	var securityGroups *schema.Set
	if sg, ok := d.GetOk(isVirtualEndpointGatewaySecurityGroups); ok {
		securityGroups = sg.(*schema.Set)
		if securityGroups != nil && securityGroups.Len() != 0 {
			securityGroupobjs := make([]vpcv1.SecurityGroupIdentityIntf, securityGroups.Len())
			for i, securityGroup := range securityGroups.List() {
				securityGroupstr := securityGroup.(string)
				securityGroupobjs[i] = &vpcv1.SecurityGroupIdentity{
					ID: &securityGroupstr,
				}
			}
			opt.SecurityGroups = securityGroupobjs
		}
	}
	// Resource group option
	if resourceGroup, ok := d.GetOk(isVirtualEndpointGatewayResourceGroupID); ok {
		resourceGroupID := resourceGroup.(string)

		resourceGroupOpt := &vpcv1.ResourceGroupIdentity{
			ID: core.StringPtr(resourceGroupID),
		}
		opt.SetResourceGroup(resourceGroupOpt)

	}
	result, response, err := sess.CreateEndpointGateway(opt)
	if err != nil {
		log.Printf("Create Endpoint Gateway failed: %v", response)
		return fmt.Errorf("Create Endpoint Gateway failed %s\n%s", err, response)
	}

	d.SetId(*result.ID)
	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk(isVirtualEndpointGatewayTags); ok || v != "" {
		oldList, newList := d.GetChange(isVirtualEndpointGatewayTags)
		err = flex.UpdateTagsUsingCRN(oldList, newList, meta, *result.CRN)
		if err != nil {
			log.Printf(
				"Error on create of VPE (%s) tags: %s", d.Id(), err)
		}
	}
	return resourceIBMisVirtualEndpointGatewayRead(d, meta)
}

func resourceIBMisVirtualEndpointGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	if d.HasChange(isVirtualEndpointGatewayName) {
		name := d.Get(isVirtualEndpointGatewayName).(string)

		// create option
		endpointGatewayPatchModel := new(vpcv1.EndpointGatewayPatch)
		endpointGatewayPatchModel.Name = core.StringPtr(name)
		endpointGatewayPatchModelAsPatch, _ := endpointGatewayPatchModel.AsPatch()
		opt := sess.NewUpdateEndpointGatewayOptions(d.Id(), endpointGatewayPatchModelAsPatch)
		_, response, err := sess.UpdateEndpointGateway(opt)
		if err != nil {
			log.Printf("Update Endpoint Gateway failed: %v", response)
			return fmt.Errorf("Update Endpoint Gateway failed : %s\n%s", err, response)
		}

	}
	id := d.Id()
	var remove, add []string
	if d.HasChange(isVirtualEndpointGatewaySecurityGroups) {
		o, n := d.GetChange(isVirtualEndpointGatewaySecurityGroups)
		oSecurityGroups := o.(*schema.Set)
		nSecurityGroups := n.(*schema.Set)
		remove = flex.ExpandStringList(oSecurityGroups.Difference(nSecurityGroups).List())
		add = flex.ExpandStringList(nSecurityGroups.Difference(oSecurityGroups).List())
		if len(add) > 0 {
			for _, sgId := range add {
				createSecurityGroupTargetBindingOptions := &vpcv1.CreateSecurityGroupTargetBindingOptions{}
				createSecurityGroupTargetBindingOptions.SecurityGroupID = &sgId
				createSecurityGroupTargetBindingOptions.ID = &id
				_, response, err := sess.CreateSecurityGroupTargetBinding(createSecurityGroupTargetBindingOptions)
				if err != nil {
					return fmt.Errorf("Error while creating Security Group Target Binding %s\n%s", err, response)
				}
				_, err = isWaitForVirtualEndpointGatewayAvailable(sess, d.Id(), d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return err
				}
			}
		}
		if len(remove) > 0 {
			for _, sgId := range remove {
				getSecurityGroupTargetOptions := &vpcv1.GetSecurityGroupTargetOptions{
					SecurityGroupID: &sgId,
					ID:              &id,
				}
				_, response, err := sess.GetSecurityGroupTarget(getSecurityGroupTargetOptions)
				if err != nil {
					if response != nil && response.StatusCode == 404 {
						continue
					}
					return fmt.Errorf("Error Getting Security Group Target for this endpoint gateway (%s): %s\n%s", sgId, err, response)
				}
				deleteSecurityGroupTargetBindingOptions := sess.NewDeleteSecurityGroupTargetBindingOptions(sgId, id)
				response, err = sess.DeleteSecurityGroupTargetBinding(deleteSecurityGroupTargetBindingOptions)
				if err != nil {
					return fmt.Errorf("Error Deleting Security Group Target for this endpoint gateway : %s\n%s", err, response)
				}
				_, err = isWaitForVirtualEndpointGatewayAvailable(sess, d.Id(), d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return err
				}
			}
		}

	}
	if d.HasChange(isVirtualEndpointGatewayTags) {
		opt := sess.NewGetEndpointGatewayOptions(d.Id())
		result, response, err := sess.GetEndpointGateway(opt)
		if err != nil {
			return fmt.Errorf("[ERROR] Error getting VPE: %s\n%s", err, response)
		}
		oldList, newList := d.GetChange(isVirtualEndpointGatewayTags)
		err = flex.UpdateTagsUsingCRN(oldList, newList, meta, *result.CRN)
		if err != nil {
			log.Printf(
				"Error on update of VPE (%s) tags: %s", d.Id(), err)
		}
	}
	return resourceIBMisVirtualEndpointGatewayRead(d, meta)
}

func resourceIBMisVirtualEndpointGatewayRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	// read option
	opt := sess.NewGetEndpointGatewayOptions(d.Id())
	result, response, err := sess.GetEndpointGateway(opt)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("Get Endpoint Gateway failed: %v", response)
		return fmt.Errorf("Get Endpoint Gateway failed %s\n%s", err, response)
	}
	d.Set(isVirtualEndpointGatewayName, result.Name)
	d.Set(isVirtualEndpointGatewayHealthState, result.HealthState)
	d.Set(isVirtualEndpointGatewayCreatedAt, result.CreatedAt.String())
	d.Set(isVirtualEndpointGatewayLifecycleState, result.LifecycleState)
	d.Set(isVirtualEndpointGatewayResourceType, result.ResourceType)
	d.Set(isVirtualEndpointGatewayCRN, result.CRN)
	d.Set(isVirtualEndpointGatewayIPs, flattenIPs(result.Ips))
	d.Set(isVirtualEndpointGatewayResourceGroupID, result.ResourceGroup.ID)
	d.Set(isVirtualEndpointGatewayTarget,
		flattenEndpointGatewayTarget(result.Target.(*vpcv1.EndpointGatewayTarget)))
	d.Set(isVirtualEndpointGatewayVpcID, result.VPC.ID)
	if result.SecurityGroups != nil {
		d.Set(isVirtualEndpointGatewaySecurityGroups, flattenDataSourceSecurityGroups(result.SecurityGroups))
	}
	tags, err := flex.GetTagsUsingCRN(meta, *result.CRN)
	if err != nil {
		log.Printf(
			"Error on get of VPE (%s) tags: %s", d.Id(), err)
	}
	d.Set(isVirtualEndpointGatewayTags, tags)
	return nil
}

func flattenDataSourceSecurityGroups(securityGroupList []vpcv1.SecurityGroupReference) interface{} {
	securitygroupList := make([]string, 0)
	for _, securityGroup := range securityGroupList {
		if securityGroup.ID != nil {
			securityGroupID := *securityGroup.ID
			securitygroupList = append(securitygroupList, securityGroupID)
		}
	}
	return securitygroupList
}

func isWaitForVirtualEndpointGatewayAvailable(sess *vpcv1.VpcV1, endPointGatewayId string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for virtual endpoint gateway (%s) to be available.", endPointGatewayId)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"waiting", "pending", "updating"},
		Target:     []string{"stable", "failed", ""},
		Refresh:    isVirtualEndpointGatewayRefreshFunc(sess, endPointGatewayId),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isVirtualEndpointGatewayRefreshFunc(sess *vpcv1.VpcV1, endPointGatewayId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		opt := sess.NewGetEndpointGatewayOptions(endPointGatewayId)
		result, response, err := sess.GetEndpointGateway(opt)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return nil, "", fmt.Errorf("Error Getting Virtual Endpoint Gateway : %s\n%s", err, response)
			}
		}
		if *result.LifecycleState == "stable" || *result.LifecycleState == "failed" {
			return result, *result.LifecycleState, nil
		}
		return result, *result.LifecycleState, nil
	}
}

func resourceIBMisVirtualEndpointGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	opt := sess.NewDeleteEndpointGatewayOptions(d.Id())
	response, err := sess.DeleteEndpointGateway(opt)
	if err != nil {
		log.Printf("Delete Endpoint Gateway failed: %v", response)
		return fmt.Errorf("Delete Endpoint Gateway failed : %s\n%s", err, response)
	}
	return nil
}

func resourceIBMisVirtualEndpointGatewayExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}

	opt := sess.NewGetEndpointGatewayOptions(d.Id())
	_, response, err := sess.GetEndpointGateway(opt)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			log.Printf("Endpoint Gateway does not exist.")
			return false, nil
		}
		log.Printf("Error : %s", response)
		return false, err
	}
	return true, nil
}

func expandIPs(ipsSet []interface{}) (ipsOptions []vpcv1.EndpointGatewayReservedIPIntf) {
	ipsList := ipsSet
	for _, item := range ipsList {
		ips := item.(map[string]interface{})
		// IPs option
		ipsID := ips[isVirtualEndpointGatewayIPsID].(string)
		ipsName := ips[isVirtualEndpointGatewayIPsName].(string)

		// IPs subnet option
		ipsSubnetID := ips[isVirtualEndpointGatewayIPsSubnet].(string)

		ipsSubnetOpt := &vpcv1.SubnetIdentity{
			ID: &ipsSubnetID,
		}

		ipsOpt := &vpcv1.EndpointGatewayReservedIP{
			ID:     core.StringPtr(ipsID),
			Name:   core.StringPtr(ipsName),
			Subnet: ipsSubnetOpt,
		}
		ipsOptions = append(ipsOptions, ipsOpt)
	}
	return ipsOptions
}

func flattenIPs(ipsList []vpcv1.ReservedIPReference) interface{} {
	ipsListOutput := make([]interface{}, 0)
	for _, item := range ipsList {
		ips := make(map[string]interface{}, 0)
		ips[isVirtualEndpointGatewayIPsID] = *item.ID
		ips[isVirtualEndpointGatewayIPsName] = *item.Name
		ips[isVirtualEndpointGatewayIPsResourceType] = *item.ResourceType
		ips[isVirtualEndpointGatewayIPsAddress] = *item.Address

		ipsListOutput = append(ipsListOutput, ips)
	}
	return ipsListOutput
}

func flattenEndpointGatewayTarget(target *vpcv1.EndpointGatewayTarget) interface{} {
	targetSlice := []interface{}{}
	targetOutput := map[string]string{}
	if target == nil {
		return targetOutput
	}
	if target.Name != nil {
		targetOutput[isVirtualEndpointGatewayTargetName] = *target.Name
	}
	if target.CRN != nil {
		targetOutput[isVirtualEndpointGatewayTargetCRN] = *target.CRN
	}
	if target.ResourceType != nil {
		targetOutput[isVirtualEndpointGatewayTargetResourceType] = *target.ResourceType
	}
	targetSlice = append(targetSlice, targetOutput)
	return targetSlice
}
