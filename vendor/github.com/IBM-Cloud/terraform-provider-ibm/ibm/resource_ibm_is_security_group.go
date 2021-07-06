// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/IBM/vpc-go-sdk/vpcclassicv1"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isSecurityGroupName          = "name"
	isSecurityGroupVPC           = "vpc"
	isSecurityGroupRules         = "rules"
	isSecurityGroupResourceGroup = "resource_group"
	isSecurityGroupTags          = "tags"
	isSecurityGroupCRN           = "crn"
)

func resourceIBMISSecurityGroup() *schema.Resource {

	return &schema.Resource{
		Create:   resourceIBMISSecurityGroupCreate,
		Read:     resourceIBMISSecurityGroupRead,
		Update:   resourceIBMISSecurityGroupUpdate,
		Delete:   resourceIBMISSecurityGroupDelete,
		Exists:   resourceIBMISSecurityGroupExists,
		Importer: &schema.ResourceImporter{},

		CustomizeDiff: customdiff.Sequence(
			func(diff *schema.ResourceDiff, v interface{}) error {
				return resourceTagsCustomizeDiff(diff)
			},
		),

		Schema: map[string]*schema.Schema{

			isSecurityGroupName: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Security group name",
				ValidateFunc: InvokeValidator("ibm_is_security_group", isSecurityGroupName),
			},
			isSecurityGroupVPC: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Security group's resource group id",
				ForceNew:    true,
			},

			isSecurityGroupTags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: InvokeValidator("ibm_is_security_group", "tag")},
				Set:         resourceIBMVPCHash,
				Description: "List of tags",
			},

			isSecurityGroupCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			isSecurityGroupRules: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Security Rules",
				Elem: &schema.Resource{
					Schema: makeIBMISSecurityRuleSchema(),
				},
			},

			isSecurityGroupResourceGroup: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Resource Group ID",
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

			ResourceCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},
		},
	}
}

func resourceIBMISSecurityGroupValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isSecurityGroupName,
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "tag",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128})

	ibmISSecurityGroupResourceValidator := ResourceValidator{ResourceName: "ibm_is_security_group", Schema: validateSchema}
	return &ibmISSecurityGroupResourceValidator
}

func resourceIBMISSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	vpc := d.Get(isSecurityGroupVPC).(string)
	if userDetails.generation == 1 {
		err := classicSgCreate(d, meta, vpc)
		if err != nil {
			return err
		}
	} else {
		err := sgCreate(d, meta, vpc)
		if err != nil {
			return err
		}
	}
	return resourceIBMISSecurityGroupRead(d, meta)
}

func classicSgCreate(d *schema.ResourceData, meta interface{}, vpc string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	createSecurityGroupOptions := &vpcclassicv1.CreateSecurityGroupOptions{
		VPC: &vpcclassicv1.VPCIdentity{
			ID: &vpc,
		},
	}
	var rg, name string
	if grp, ok := d.GetOk(isSecurityGroupResourceGroup); ok {
		rg = grp.(string)
		createSecurityGroupOptions.ResourceGroup = &vpcclassicv1.ResourceGroupIdentity{
			ID: &rg,
		}
	}
	if nm, ok := d.GetOk(isSecurityGroupName); ok {
		name = nm.(string)
		createSecurityGroupOptions.Name = &name
	}

	sg, response, err := sess.CreateSecurityGroup(createSecurityGroupOptions)
	if err != nil {
		return fmt.Errorf("Error while creating Security Group %s\n%s", err, response)
	}
	d.SetId(*sg.ID)
	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk(isSecurityGroupTags); ok || v != "" {
		oldList, newList := d.GetChange(isSecurityGroupTags)
		err = UpdateTagsUsingCRN(oldList, newList, meta, *sg.CRN)
		if err != nil {
			log.Printf("Error while creating Security Group tags %s\n%s", *sg.ID, err)
		}
	}
	return nil
}

func sgCreate(d *schema.ResourceData, meta interface{}, vpc string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	createSecurityGroupOptions := &vpcv1.CreateSecurityGroupOptions{
		VPC: &vpcv1.VPCIdentity{
			ID: &vpc,
		},
	}
	var rg, name string
	if grp, ok := d.GetOk(isSecurityGroupResourceGroup); ok {
		rg = grp.(string)
		createSecurityGroupOptions.ResourceGroup = &vpcv1.ResourceGroupIdentity{
			ID: &rg,
		}
	}
	if nm, ok := d.GetOk(isSecurityGroupName); ok {
		name = nm.(string)
		createSecurityGroupOptions.Name = &name
	}
	sg, response, err := sess.CreateSecurityGroup(createSecurityGroupOptions)
	if err != nil {
		return fmt.Errorf("Error while creating Security Group %s\n%s", err, response)
	}
	d.SetId(*sg.ID)
	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk(isSecurityGroupTags); ok || v != "" {
		oldList, newList := d.GetChange(isSecurityGroupTags)
		err = UpdateTagsUsingCRN(oldList, newList, meta, *sg.CRN)
		if err != nil {
			log.Printf(
				"Error while creating Security Group tags : %s\n%s", *sg.ID, err)
		}
	}
	return nil
}

func resourceIBMISSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	id := d.Id()
	if userDetails.generation == 1 {
		err := classicSgGet(d, meta, id)
		if err != nil {
			return err
		}
	} else {
		err := sgGet(d, meta, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func classicSgGet(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	getSecurityGroupOptions := &vpcclassicv1.GetSecurityGroupOptions{
		ID: &id,
	}
	group, response, err := sess.GetSecurityGroup(getSecurityGroupOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting Security Group : %s\n%s", err, response)
	}
	tags, err := GetTagsUsingCRN(meta, *group.CRN)
	if err != nil {
		log.Printf(
			"Error getting Security Group tags : %s\n%s", d.Id(), err)
	}
	d.Set(isSecurityGroupTags, tags)
	d.Set(isSecurityGroupCRN, *group.CRN)
	d.Set(isSecurityGroupName, *group.Name)
	d.Set(isSecurityGroupVPC, *group.VPC.ID)
	rules := make([]map[string]interface{}, 0)
	if len(group.Rules) > 0 {
		for _, rule := range group.Rules {
			switch reflect.TypeOf(rule).String() {
			case "*vpcclassicv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp":
				{
					rule := rule.(*vpcclassicv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp)
					r := make(map[string]interface{})
					if rule.Code != nil {
						r[isSecurityGroupRuleCode] = int(*rule.Code)
					}
					if rule.Type != nil {
						r[isSecurityGroupRuleType] = int(*rule.Type)
					}
					r[isSecurityGroupRuleDirection] = *rule.Direction
					r[isSecurityGroupRuleIPVersion] = *rule.IPVersion
					if rule.Protocol != nil {
						r[isSecurityGroupRuleProtocol] = *rule.Protocol
					}
					remote, ok := rule.Remote.(*vpcclassicv1.SecurityGroupRuleRemote)
					if ok {
						if remote != nil && reflect.ValueOf(remote).IsNil() == false {
							if remote.ID != nil {
								r[isSecurityGroupRuleRemote] = remote.ID
							} else if remote.Address != nil {
								r[isSecurityGroupRuleRemote] = remote.Address
							} else if remote.CIDRBlock != nil {
								r[isSecurityGroupRuleRemote] = remote.CIDRBlock
							}
						}
					}
					rules = append(rules, r)
				}
			case "*vpcclassicv1.SecurityGroupRuleSecurityGroupRuleProtocolAll":
				{
					rule := rule.(*vpcclassicv1.SecurityGroupRuleSecurityGroupRuleProtocolAll)
					r := make(map[string]interface{})
					r[isSecurityGroupRuleDirection] = *rule.Direction
					r[isSecurityGroupRuleIPVersion] = *rule.IPVersion
					if rule.Protocol != nil {
						r[isSecurityGroupRuleProtocol] = *rule.Protocol
					}
					remote, ok := rule.Remote.(*vpcclassicv1.SecurityGroupRuleRemote)
					if ok {
						if remote != nil && reflect.ValueOf(remote).IsNil() == false {
							if remote.ID != nil {
								r[isSecurityGroupRuleRemote] = remote.ID
							} else if remote.Address != nil {
								r[isSecurityGroupRuleRemote] = remote.Address
							} else if remote.CIDRBlock != nil {
								r[isSecurityGroupRuleRemote] = remote.CIDRBlock
							}
						}
					}
					rules = append(rules, r)
				}
			case "*vpcclassicv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp":
				{
					rule := rule.(*vpcclassicv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp)
					r := make(map[string]interface{})
					if rule.PortMin != nil {
						r[isSecurityGroupRulePortMin] = int(*rule.PortMin)
					}
					if rule.PortMax != nil {
						r[isSecurityGroupRulePortMax] = int(*rule.PortMax)
					}
					r[isSecurityGroupRuleDirection] = *rule.Direction
					r[isSecurityGroupRuleIPVersion] = *rule.IPVersion
					if rule.Protocol != nil {
						r[isSecurityGroupRuleProtocol] = *rule.Protocol
					}
					remote, ok := rule.Remote.(*vpcclassicv1.SecurityGroupRuleRemote)
					if ok {
						if remote != nil && reflect.ValueOf(remote).IsNil() == false {
							if remote.ID != nil {
								r[isSecurityGroupRuleRemote] = remote.ID
							} else if remote.Address != nil {
								r[isSecurityGroupRuleRemote] = remote.Address
							} else if remote.CIDRBlock != nil {
								r[isSecurityGroupRuleRemote] = remote.CIDRBlock
							}
						}
					}
					rules = append(rules, r)
				}
			}
		}
	}
	d.Set(isSecurityGroupRules, rules)
	d.SetId(*group.ID)
	if group.ResourceGroup != nil {
		d.Set(isSecurityGroupResourceGroup, group.ResourceGroup.ID)
		rsMangClient, err := meta.(ClientSession).ResourceManagementAPIv2()
		if err != nil {
			return err
		}
		grp, err := rsMangClient.ResourceGroup().Get(*group.ResourceGroup.ID)
		if err != nil {
			return err
		}
		d.Set(ResourceGroupName, grp.Name)
	}
	controller, err := getBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(ResourceControllerURL, controller+"/vpc/network/securityGroups")
	d.Set(ResourceName, *group.Name)
	d.Set(ResourceCRN, *group.CRN)
	return nil
}

func sgGet(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	getSecurityGroupOptions := &vpcv1.GetSecurityGroupOptions{
		ID: &id,
	}
	group, response, err := sess.GetSecurityGroup(getSecurityGroupOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting Security Group : %s\n%s", err, response)
	}
	tags, err := GetTagsUsingCRN(meta, *group.CRN)
	if err != nil {
		log.Printf(
			"Error getting Security Group tags : %s\n%s", d.Id(), err)
	}
	d.Set(isSecurityGroupTags, tags)
	d.Set(isSecurityGroupCRN, *group.CRN)
	d.Set(isSecurityGroupName, *group.Name)
	d.Set(isSecurityGroupVPC, *group.VPC.ID)
	rules := make([]map[string]interface{}, 0)
	if len(group.Rules) > 0 {
		for _, rule := range group.Rules {
			switch reflect.TypeOf(rule).String() {
			case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp":
				{
					rule := rule.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp)
					r := make(map[string]interface{})
					if rule.Code != nil {
						r[isSecurityGroupRuleCode] = int(*rule.Code)
					}
					if rule.Type != nil {
						r[isSecurityGroupRuleType] = int(*rule.Type)
					}
					r[isSecurityGroupRuleDirection] = *rule.Direction
					r[isSecurityGroupRuleIPVersion] = *rule.IPVersion
					if rule.Protocol != nil {
						r[isSecurityGroupRuleProtocol] = *rule.Protocol
					}
					remote, ok := rule.Remote.(*vpcv1.SecurityGroupRuleRemote)
					if ok {
						if remote != nil && reflect.ValueOf(remote).IsNil() == false {
							if remote.ID != nil {
								r[isSecurityGroupRuleRemote] = remote.ID
							} else if remote.Address != nil {
								r[isSecurityGroupRuleRemote] = remote.Address
							} else if remote.CIDRBlock != nil {
								r[isSecurityGroupRuleRemote] = remote.CIDRBlock
							}
						}
					}
					rules = append(rules, r)
				}
			case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll":
				{
					rule := rule.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll)
					r := make(map[string]interface{})
					r[isSecurityGroupRuleDirection] = *rule.Direction
					r[isSecurityGroupRuleIPVersion] = *rule.IPVersion
					if rule.Protocol != nil {
						r[isSecurityGroupRuleProtocol] = *rule.Protocol
					}
					remote, ok := rule.Remote.(*vpcv1.SecurityGroupRuleRemote)
					if ok {
						if remote != nil && reflect.ValueOf(remote).IsNil() == false {
							if remote.ID != nil {
								r[isSecurityGroupRuleRemote] = remote.ID
							} else if remote.Address != nil {
								r[isSecurityGroupRuleRemote] = remote.Address
							} else if remote.CIDRBlock != nil {
								r[isSecurityGroupRuleRemote] = remote.CIDRBlock
							}
						}
					}
					rules = append(rules, r)
				}
			case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp":
				{
					rule := rule.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp)
					r := make(map[string]interface{})
					if rule.PortMin != nil {
						r[isSecurityGroupRulePortMin] = int(*rule.PortMin)
					}
					if rule.PortMax != nil {
						r[isSecurityGroupRulePortMax] = int(*rule.PortMax)
					}
					r[isSecurityGroupRuleDirection] = *rule.Direction
					r[isSecurityGroupRuleIPVersion] = *rule.IPVersion
					if rule.Protocol != nil {
						r[isSecurityGroupRuleProtocol] = *rule.Protocol
					}
					remote, ok := rule.Remote.(*vpcv1.SecurityGroupRuleRemote)
					if ok {
						if remote != nil && reflect.ValueOf(remote).IsNil() == false {
							if remote.ID != nil {
								r[isSecurityGroupRuleRemote] = remote.ID
							} else if remote.Address != nil {
								r[isSecurityGroupRuleRemote] = remote.Address
							} else if remote.CIDRBlock != nil {
								r[isSecurityGroupRuleRemote] = remote.CIDRBlock
							}
						}
					}
					rules = append(rules, r)
				}
			}
		}
	}
	d.Set(isSecurityGroupRules, rules)
	d.SetId(*group.ID)
	if group.ResourceGroup != nil {
		d.Set(isSecurityGroupResourceGroup, group.ResourceGroup.ID)
		d.Set(ResourceGroupName, group.ResourceGroup.Name)
	}
	controller, err := getBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(ResourceControllerURL, controller+"/vpc-ext/network/securityGroups")
	d.Set(ResourceName, *group.Name)
	d.Set(ResourceCRN, *group.CRN)
	return nil
}

func resourceIBMISSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	id := d.Id()
	name := ""
	hasChanged := false

	if d.HasChange(isSecurityGroupTags) {
		oldList, newList := d.GetChange(isSecurityGroupTags)
		err = UpdateTagsUsingCRN(oldList, newList, meta, d.Get(isSecurityGroupCRN).(string))
		if err != nil {
			log.Printf(
				"Error Updating Security Group tags: %s\n%s", d.Id(), err)
		}
	}

	if d.HasChange(isSecurityGroupName) {
		name = d.Get(isSecurityGroupName).(string)
		hasChanged = true
	} else {
		return resourceIBMISSecurityGroupRead(d, meta)
	}
	if userDetails.generation == 1 {
		err := classicSgUpdate(d, meta, id, name, hasChanged)
		if err != nil {
			return err
		}
	} else {
		err := sgUpdate(d, meta, id, name, hasChanged)
		if err != nil {
			return err
		}
	}
	return resourceIBMISSecurityGroupRead(d, meta)
}

func classicSgUpdate(d *schema.ResourceData, meta interface{}, id, name string, hasChanged bool) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	if hasChanged {
		updateSecurityGroupOptions := &vpcclassicv1.UpdateSecurityGroupOptions{
			ID: &id,
		}
		securityGroupPatchModel := &vpcclassicv1.SecurityGroupPatch{
			Name: &name,
		}
		securityGroupPatch, err := securityGroupPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("Error calling asPatch for SecurityGroupPatch: %s", err)
		}
		updateSecurityGroupOptions.SecurityGroupPatch = securityGroupPatch
		_, response, err := sess.UpdateSecurityGroup(updateSecurityGroupOptions)
		if err != nil {
			return fmt.Errorf("Error Updating Security Group : %s\n%s", err, response)
		}
	}
	return nil
}

func sgUpdate(d *schema.ResourceData, meta interface{}, id, name string, hasChanged bool) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	if hasChanged {
		updateSecurityGroupOptions := &vpcv1.UpdateSecurityGroupOptions{
			ID: &id,
		}
		securityGroupPatchModel := &vpcv1.SecurityGroupPatch{
			Name: &name,
		}
		securityGroupPatch, err := securityGroupPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("Error calling asPatch for SecurityGroupPatch: %s", err)
		}
		updateSecurityGroupOptions.SecurityGroupPatch = securityGroupPatch
		_, response, err := sess.UpdateSecurityGroup(updateSecurityGroupOptions)
		if err != nil {
			return fmt.Errorf("Error Updating Security Group : %s\n%s", err, response)
		}
	}
	return nil
}

func resourceIBMISSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	id := d.Id()
	if userDetails.generation == 1 {
		err := classicSgDelete(d, meta, id)
		if err != nil {
			return err
		}
	} else {
		err := sgDelete(d, meta, id)
		if err != nil {
			return err
		}
	}
	d.SetId("")
	return nil
}

func classicSgDelete(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	getSecurityGroupOptions := &vpcclassicv1.GetSecurityGroupOptions{
		ID: &id,
	}
	_, response, err := sess.GetSecurityGroup(getSecurityGroupOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting Security Group (%s): %s\n%s", id, err, response)
	}

	deleteSecurityGroupOptions := &vpcclassicv1.DeleteSecurityGroupOptions{
		ID: &id,
	}
	response, err = sess.DeleteSecurityGroup(deleteSecurityGroupOptions)
	if err != nil {
		return fmt.Errorf("Error Deleting Security Group : %s\n%s", err, response)
	}
	d.SetId("")
	return nil
}

func sgDelete(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	getSecurityGroupOptions := &vpcv1.GetSecurityGroupOptions{
		ID: &id,
	}
	_, response, err := sess.GetSecurityGroup(getSecurityGroupOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting Security Group (%s): %s\n%s", id, err, response)
	}

	deleteSecurityGroupOptions := &vpcv1.DeleteSecurityGroupOptions{
		ID: &id,
	}
	response, err = sess.DeleteSecurityGroup(deleteSecurityGroupOptions)
	if err != nil {
		return fmt.Errorf("Error Deleting Security Group : %s\n%s", err, response)
	}
	d.SetId("")
	return nil
}

func resourceIBMISSecurityGroupExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return false, err
	}
	id := d.Id()
	if userDetails.generation == 1 {
		exists, err := classicSgExists(d, meta, id)
		return exists, err
	} else {
		exists, err := sgExists(d, meta, id)
		return exists, err
	}
}

func classicSgExists(d *schema.ResourceData, meta interface{}, id string) (bool, error) {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return false, err
	}
	getSecurityGroupOptions := &vpcclassicv1.GetSecurityGroupOptions{
		ID: &id,
	}
	_, response, err := sess.GetSecurityGroup(getSecurityGroupOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error getting Security Group: %s\n%s", err, response)
	}
	return true, nil
}

func sgExists(d *schema.ResourceData, meta interface{}, id string) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}
	getSecurityGroupOptions := &vpcv1.GetSecurityGroupOptions{
		ID: &id,
	}
	_, response, err := sess.GetSecurityGroup(getSecurityGroupOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error getting Security Group: %s\n%s", err, response)
	}
	return true, nil
}

func makeIBMISSecurityRuleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		isSecurityGroupRuleDirection: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Direction of traffic to enforce, either inbound or outbound",
		},

		isSecurityGroupRuleIPVersion: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "IP version: ipv4 or ipv6",
		},

		isSecurityGroupRuleRemote: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Security group id: an IP address, a CIDR block, or a single security group identifier",
		},

		isSecurityGroupRuleType: {
			Type:     schema.TypeInt,
			Computed: true,
		},

		isSecurityGroupRuleCode: {
			Type:     schema.TypeInt,
			Computed: true,
		},

		isSecurityGroupRulePortMin: {
			Type:     schema.TypeInt,
			Computed: true,
		},

		isSecurityGroupRulePortMax: {
			Type:     schema.TypeInt,
			Computed: true,
		},

		isSecurityGroupRuleProtocol: {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
