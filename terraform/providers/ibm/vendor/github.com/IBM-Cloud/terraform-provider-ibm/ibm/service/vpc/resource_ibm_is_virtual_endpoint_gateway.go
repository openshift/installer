// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
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
	isVirtualEndpointGatewayName                      = "name"
	isVirtualEndpointGatewayResourceType              = "resource_type"
	isVirtualEndpointGatewayCRN                       = "crn"
	isVirtualEndpointGatewayResourceGroupID           = "resource_group"
	isVirtualEndpointGatewayCreatedAt                 = "created_at"
	isVirtualEndpointGatewayIPs                       = "ips"
	isVirtualEndpointGatewayIPsID                     = "id"
	isVirtualEndpointGatewayIPsAddress                = "address"
	isVirtualEndpointGatewayIPsName                   = "name"
	isVirtualEndpointGatewayIPsSubnet                 = "subnet"
	isVirtualEndpointGatewayIPsResourceType           = "resource_type"
	isVirtualEndpointGatewayHealthState               = "health_state"
	isVirtualEndpointGatewayLifecycleState            = "lifecycle_state"
	isVirtualEndpointGatewayLifecycleReasons          = "lifecycle_reasons"
	isVirtualEndpointGatewayTarget                    = "target"
	isVirtualEndpointGatewayTargetName                = "name"
	isVirtualEndpointGatewayTargetCRN                 = "crn"
	isVirtualEndpointGatewayTargetResourceType        = "resource_type"
	isVirtualEndpointGatewayVpcID                     = "vpc"
	isVirtualEndpointGatewayTags                      = "tags"
	isVirtualEndpointGatewaySecurityGroups            = "security_groups"
	isVirtualEndpointGatewayServiceEndpoints          = "service_endpoints"
	isVirtualEndpointGatewayAccessTags                = "access_tags"
	isVirtualEndpointGatewayAllowDnsResolutionBinding = "allow_dns_resolution_binding"
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

		CustomizeDiff: customdiff.All(
			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.ResourceTagsCustomizeDiff(diff)
				},
			),
			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.ResourceValidateAccessTags(diff, v)
				}),
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
			isVirtualEndpointGatewayServiceEndpoints: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The fully qualified domain names for the target service. A fully qualified domain name for the target service",
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
			isVirtualEndpointGatewayLifecycleReasons: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reasons for the current lifecycle_state (if any).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A snake case string succinctly identifying the reason for this lifecycle state.",
						},

						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the reason for this lifecycle state.",
						},

						"more_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Link to documentation about the reason for this lifecycle state.",
						},
					},
				},
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
							Computed: true,
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
			isVirtualEndpointGatewayAllowDnsResolutionBinding: {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether to allow this endpoint gateway to participate in DNS resolution bindings with a VPC that has dns.enable_hub set to true.",
			},
			isVirtualEndpointGatewayTags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_virtual_endpoint_gateway", "tags")},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of tags for VPE",
			},
			isVirtualEndpointGatewayAccessTags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_virtual_endpoint_gateway", "accesstag")},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of access management tags",
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
			Identifier:                 "tags",
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
			AllowedValues:              "provider_cloud_service, provider_infrastructure_service, private_path_service_gateway"})

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "accesstag",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([A-Za-z0-9_.-]|[A-Za-z0-9_.-][A-Za-z0-9_ .-]*[A-Za-z0-9_.-]):([A-Za-z0-9_.-]|[A-Za-z0-9_.-][A-Za-z0-9_ .-]*[A-Za-z0-9_.-])$`,
			MinValueLength:             1,
			MaxValueLength:             128})

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
	// dns resolution binding change
	if allowDnsResolutionBindingOk, ok := d.GetOkExists(isVirtualEndpointGatewayAllowDnsResolutionBinding); ok {
		allowDnsResolutionBinding := allowDnsResolutionBindingOk.(bool)
		opt.AllowDnsResolutionBinding = &allowDnsResolutionBinding
	}
	endpointGateway, response, err := sess.CreateEndpointGateway(opt)
	if err != nil {
		log.Printf("Create Endpoint Gateway failed: %v", response)
		return fmt.Errorf("[ERROR] Create Endpoint Gateway failed %s\n%s", err, response)
	}

	d.SetId(*endpointGateway.ID)

	_, err = isWaitForVirtualEndpointGatewayAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		if d.Get(targetResourceTypeFmt).(string) == "private_path_service_gateway" {
			isAccessPending := false
			if strings.Contains(err.Error(), "timeout while waiting for state to become") {
				opt := sess.NewGetEndpointGatewayOptions(d.Id())
				endpointGateway, response, err := sess.GetEndpointGateway(opt)
				if err != nil {
					log.Printf("Get Endpoint Gateway failed: %v", response)
					return fmt.Errorf("[ERROR] Get Endpoint Gateway failed %s\n%s", err, response)
				}
				if len(endpointGateway.LifecycleReasons) > 0 {
					if endpointGateway.LifecycleReasons[0].Code != nil && strings.Compare(*endpointGateway.LifecycleReasons[0].Code, "access_pending") == 0 {
						isAccessPending = true
					}
				}
			}
			if !isAccessPending {
				return err
			}
		} else {
			return err
		}
	}

	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk(isVirtualEndpointGatewayTags); ok || v != "" {
		oldList, newList := d.GetChange(isVirtualEndpointGatewayTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *endpointGateway.CRN, "", isUserTagType)
		if err != nil {
			log.Printf(
				"Error on create of VPE (%s) tags: %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk(isVirtualEndpointGatewayAccessTags); ok {
		oldList, newList := d.GetChange(isVirtualEndpointGatewayAccessTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *endpointGateway.CRN, "", isAccessTagType)
		if err != nil {
			log.Printf(
				"Error on create of VPE (%s) access tags: %s", d.Id(), err)
		}
	}

	return resourceIBMisVirtualEndpointGatewayRead(d, meta)
}

func resourceIBMisVirtualEndpointGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	// create option
	endpointGatewayPatchModel := new(vpcv1.EndpointGatewayPatch)
	if d.HasChange(isVirtualEndpointGatewayName) {
		name := d.Get(isVirtualEndpointGatewayName).(string)
		endpointGatewayPatchModel.Name = core.StringPtr(name)
	}
	if d.HasChange(isVirtualEndpointGatewayAllowDnsResolutionBinding) {
		allowDnsResolutionBinding := d.Get(isVirtualEndpointGatewayAllowDnsResolutionBinding).(bool)
		endpointGatewayPatchModel.AllowDnsResolutionBinding = &allowDnsResolutionBinding
	}
	endpointGatewayPatchModelAsPatch, _ := endpointGatewayPatchModel.AsPatch()
	opt := sess.NewUpdateEndpointGatewayOptions(d.Id(), endpointGatewayPatchModelAsPatch)
	_, response, err := sess.UpdateEndpointGateway(opt)
	if err != nil {
		log.Printf("Update Endpoint Gateway failed: %v", response)
		return fmt.Errorf("Update Endpoint Gateway failed : %s\n%s", err, response)
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
	if d.HasChange(isVirtualEndpointGatewayTags) || d.HasChange(isVirtualEndpointGatewayAccessTags) {
		opt := sess.NewGetEndpointGatewayOptions(d.Id())
		endpointGateway, response, err := sess.GetEndpointGateway(opt)
		if err != nil {
			return fmt.Errorf("[ERROR] Error getting VPE: %s\n%s", err, response)
		}
		if d.HasChange(isVirtualEndpointGatewayTags) {
			oldList, newList := d.GetChange(isVirtualEndpointGatewayTags)
			err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *endpointGateway.CRN, "", isUserTagType)
			if err != nil {
				log.Printf(
					"Error on update of VPE (%s) tags: %s", d.Id(), err)
			}
		}

		if d.HasChange(isVirtualEndpointGatewayAccessTags) {
			oldList, newList := d.GetChange(isVirtualEndpointGatewayAccessTags)
			err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *endpointGateway.CRN, "", isAccessTagType)
			if err != nil {
				log.Printf(
					"Error on update of VPE (%s) access tags: %s", d.Id(), err)
			}
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
	endpointGateway, response, err := sess.GetEndpointGateway(opt)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("Get Endpoint Gateway failed: %v", response)
		return fmt.Errorf("[ERROR] Get Endpoint Gateway failed %s\n%s", err, response)
	}
	d.Set(isVirtualEndpointGatewayName, endpointGateway.Name)
	d.Set(isVirtualEndpointGatewayHealthState, endpointGateway.HealthState)
	d.Set(isVirtualEndpointGatewayCreatedAt, endpointGateway.CreatedAt.String())
	d.Set(isVirtualEndpointGatewayLifecycleState, endpointGateway.LifecycleState)
	if err := d.Set(isVirtualEndpointGatewayLifecycleReasons, resourceEGWFlattenLifecycleReasons(endpointGateway.LifecycleReasons)); err != nil {
		return fmt.Errorf("[ERROR] Error setting lifecycle_reasons: %s", err)
	}
	d.Set(isVirtualEndpointGatewayAllowDnsResolutionBinding, endpointGateway.AllowDnsResolutionBinding)
	d.Set(isVirtualEndpointGatewayResourceType, endpointGateway.ResourceType)
	d.Set(isVirtualEndpointGatewayCRN, endpointGateway.CRN)
	d.Set(isVirtualEndpointGatewayIPs, flattenIPs(endpointGateway.Ips))
	d.Set(isVirtualEndpointGatewayResourceGroupID, endpointGateway.ResourceGroup.ID)
	d.Set(isVirtualEndpointGatewayTarget,
		flattenEndpointGatewayTarget(endpointGateway.Target.(*vpcv1.EndpointGatewayTarget)))
	serviceEndpoints := []string{}
	if len(endpointGateway.ServiceEndpoints) > 0 {
		serviceEndpoints = endpointGateway.ServiceEndpoints
	}
	d.Set(isVirtualEndpointGatewayServiceEndpoints, serviceEndpoints)
	d.Set(isVirtualEndpointGatewayVpcID, endpointGateway.VPC.ID)
	if endpointGateway.SecurityGroups != nil {
		d.Set(isVirtualEndpointGatewaySecurityGroups, flattenDataSourceSecurityGroups(endpointGateway.SecurityGroups))
	}
	tags, err := flex.GetGlobalTagsUsingCRN(meta, *endpointGateway.CRN, "", isUserTagType)
	if err != nil {
		log.Printf(
			"Error on get of VPE (%s) tags: %s", d.Id(), err)
	}
	d.Set(isVirtualEndpointGatewayTags, tags)

	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *endpointGateway.CRN, "", isAccessTagType)
	if err != nil {
		log.Printf(
			"Error on get of VPE (%s) access tags: %s", d.Id(), err)
	}
	d.Set(isVirtualEndpointGatewayAccessTags, accesstags)

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

func isWaitForVirtualEndpointGatewayForPPSGAvailable(sess *vpcv1.VpcV1, endPointGatewayId string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for virtual endpoint gateway (%s) to be available.", endPointGatewayId)
	// When the target is PPSG, pending is a valid state when the endpoint gateway binding is not permitted within the terraform configuration.
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"waiting", "updating"},
		Target:                    []string{"stable", "failed", "pending", ""},
		Refresh:                   isVirtualEndpointGatewayRefreshFunc(sess, endPointGatewayId),
		Timeout:                   timeout,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 6,
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
	_, err = isWaitForEGWDelete(sess, d, d.Id())
	if err != nil {
		return err
	}
	return nil
}

func isWaitForEGWDelete(vpcClient *vpcv1.VpcV1, d *schema.ResourceData, id string) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending: []string{"deleting", "stable"},
		Target:  []string{"done", ""},
		Refresh: func() (interface{}, string, error) {
			getegwoptions := &vpcv1.GetEndpointGatewayOptions{
				ID: &id,
			}
			egw, response, err := vpcClient.GetEndpointGateway(getegwoptions)
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					return egw, "done", nil
				}
				return nil, "", fmt.Errorf("[ERROR] Error Getting EGW: %s\n%s", err, response)
			}
			if *egw.LifecycleState == "failed" {
				return egw, *egw.LifecycleState, fmt.Errorf("[ERROR] The egw %s failed to delete: %v", d.Id(), err)
			}
			return egw, "deleting", nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
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
			Subnet: ipsSubnetOpt,
		}
		if ipsID != "" {
			ipsOpt.ID = &ipsID
		}
		if ipsName != "" {
			ipsOpt.Name = &ipsName
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

func resourceEGWFlattenLifecycleReasons(lifecycleReasons []vpcv1.EndpointGatewayLifecycleReason) (lifecycleReasonsList []map[string]interface{}) {
	lifecycleReasonsList = make([]map[string]interface{}, 0)
	for _, lr := range lifecycleReasons {
		currentLR := map[string]interface{}{}
		if lr.Code != nil && lr.Message != nil {
			currentLR[isInstanceLifecycleReasonsCode] = *lr.Code
			currentLR[isInstanceLifecycleReasonsMessage] = *lr.Message
			if lr.MoreInfo != nil {
				currentLR[isInstanceLifecycleReasonsMoreInfo] = *lr.MoreInfo
			}
			lifecycleReasonsList = append(lifecycleReasonsList, currentLR)
		}
	}
	return lifecycleReasonsList
}
