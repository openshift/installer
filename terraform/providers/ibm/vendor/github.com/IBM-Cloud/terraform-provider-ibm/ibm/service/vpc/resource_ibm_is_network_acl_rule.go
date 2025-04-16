// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isNwACLID         = "network_acl"
	isNwACLRuleId     = "rule_id"
	isNwACLRuleBefore = "before"
)

func ResourceIBMISNetworkACLRule() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISNetworkACLRuleCreate,
		Read:     resourceIBMISNetworkACLRuleRead,
		Update:   resourceIBMISNetworkACLRuleUpdate,
		Delete:   resourceIBMISNetworkACLRuleDelete,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: customdiff.Sequence(
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return flex.ResourceTagsCustomizeDiff(diff)
			},
		),

		Schema: map[string]*schema.Schema{
			isNwACLID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Network ACL id",
			},
			isNwACLRuleId: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The network acl rule id.",
			},
			isNwACLRuleBefore: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The rule that this rule is immediately before. If absent, this is the last rule.",
			},
			isNetworkACLRuleProtocol: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The protocol of the rule.",
			},
			isNetworkACLRuleHref: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The url of the rule.",
			},
			isNetworkACLRuleName: {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				Description:  "The user-defined name for this rule. Names must be unique within the network ACL the rule resides in. If unspecified, the name will be a hyphenated list of randomly-selected words.",
				ValidateFunc: validate.InvokeValidator("ibm_is_network_acl_rule", isNetworkACLRuleName),
			},
			isNetworkACLRuleAction: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				Description:  "Whether to allow or deny matching traffic",
				ValidateFunc: validate.InvokeValidator("ibm_is_network_acl_rule", isNetworkACLRuleAction),
			},
			isNetworkACLRuleIPVersion: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP version for this rule.",
			},
			isNetworkACLRuleSource: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				Description:  "The source CIDR block. The CIDR block 0.0.0.0/0 applies to all addresses.",
				ValidateFunc: validate.InvokeValidator("ibm_is_network_acl_rule", isNetworkACLRuleSource),
			},
			isNetworkACLRuleDestination: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validate.InvokeValidator("ibm_is_network_acl_rule", isNetworkACLRuleDestination),
				Description:  "The destination CIDR block. The CIDR block 0.0.0.0/0 applies to all addresses.",
			},
			isNetworkACLRuleDirection: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				Description:  "Direction of traffic to enforce, either inbound or outbound",
				ValidateFunc: validate.InvokeValidator("ibm_is_network_acl_rule", isNetworkACLRuleDirection),
			},
			isNetworkACLRuleICMP: {
				Type:          schema.TypeList,
				MinItems:      0,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: []string{isNetworkACLRuleTCP, isNetworkACLRuleUDP},
				ForceNew:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isNetworkACLRuleICMPCode: {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validate.InvokeValidator("ibm_is_network_acl_rule", isNetworkACLRuleICMPCode),
							Description:  "The ICMP traffic code to allow. Valid values from 0 to 255.",
						},
						isNetworkACLRuleICMPType: {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validate.InvokeValidator("ibm_is_network_acl_rule", isNetworkACLRuleICMPType),
							Description:  "The ICMP traffic type to allow. Valid values from 0 to 254.",
						},
					},
				},
			},

			isNetworkACLRuleTCP: {
				Type:          schema.TypeList,
				MinItems:      0,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: []string{isNetworkACLRuleICMP, isNetworkACLRuleUDP},
				ForceNew:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isNetworkACLRulePortMax: {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      65535,
							ValidateFunc: validate.InvokeValidator("ibm_is_network_acl_rule", isNetworkACLRulePortMax),
							Description:  "The highest port in the range of ports to be matched",
						},
						isNetworkACLRulePortMin: {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validate.InvokeValidator("ibm_is_network_acl_rule", isNetworkACLRulePortMin),
							Description:  "The lowest port in the range of ports to be matched",
						},
						isNetworkACLRuleSourcePortMax: {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      65535,
							ValidateFunc: validate.InvokeValidator("ibm_is_network_acl_rule", isNetworkACLRuleSourcePortMax),
							Description:  "The highest port in the range of ports to be matched",
						},
						isNetworkACLRuleSourcePortMin: {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validate.InvokeValidator("ibm_is_network_acl_rule", isNetworkACLRuleSourcePortMin),
							Description:  "The lowest port in the range of ports to be matched",
						},
					},
				},
			},

			isNetworkACLRuleUDP: {
				Type:          schema.TypeList,
				MinItems:      0,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: []string{isNetworkACLRuleICMP, isNetworkACLRuleTCP},
				ForceNew:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isNetworkACLRulePortMax: {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      65535,
							ValidateFunc: validate.InvokeValidator("ibm_is_network_acl_rule", isNetworkACLRulePortMax),
							Description:  "The highest port in the range of ports to be matched",
						},
						isNetworkACLRulePortMin: {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validate.InvokeValidator("ibm_is_network_acl_rule", isNetworkACLRulePortMin),
							Description:  "The lowest port in the range of ports to be matched",
						},
						isNetworkACLRuleSourcePortMax: {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      65535,
							ValidateFunc: validate.InvokeValidator("ibm_is_network_acl_rule", isNetworkACLRuleSourcePortMax),
							Description:  "The highest port in the range of ports to be matched",
						},
						isNetworkACLRuleSourcePortMin: {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validate.InvokeValidator("ibm_is_network_acl_rule", isNetworkACLRuleSourcePortMin),
							Description:  "The lowest port in the range of ports to be matched",
						},
					},
				},
			},
		},
	}
}

func ResourceIBMISNetworkACLRuleValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)
	direction := "inbound, outbound"
	action := "allow, deny"

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isNetworkACLRuleAction,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              action})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isNetworkACLRuleDirection,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              direction})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isNwACLID,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isNetworkACLRuleName,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isNetworkACLRuleDestination,
			ValidateFunctionIdentifier: validate.ValidateIPorCIDR,
			Type:                       validate.TypeString,
			Required:                   true})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isNetworkACLRuleSource,
			ValidateFunctionIdentifier: validate.ValidateIPorCIDR,
			Type:                       validate.TypeString,
			Required:                   true})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isNetworkACLRuleICMPType,
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			MinValue:                   "0",
			MaxValue:                   "254"})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isNetworkACLRuleICMPCode,
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			MinValue:                   "0",
			MaxValue:                   "255"})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isNetworkACLRulePortMin,
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			MinValue:                   "1",
			MaxValue:                   "65535"})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isNetworkACLRulePortMax,
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			MinValue:                   "1",
			MaxValue:                   "65535"})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isNetworkACLRuleSourcePortMin,
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			MinValue:                   "1",
			MaxValue:                   "65535"})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isNetworkACLRuleSourcePortMax,
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			MinValue:                   "1",
			MaxValue:                   "65535"})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "tags",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128})

	ibmISNetworkACLRuleResourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_network_acl_rule", Schema: validateSchema}
	return &ibmISNetworkACLRuleResourceValidator
}

func resourceIBMISNetworkACLRuleCreate(d *schema.ResourceData, meta interface{}) error {
	nwACLID := d.Get(isNwACLID).(string)

	err := nwaclRuleCreate(d, meta, nwACLID)
	if err != nil {
		return err
	}
	return resourceIBMISNetworkACLRuleRead(d, meta)

}

func nwaclRuleCreate(d *schema.ResourceData, meta interface{}, nwACLID string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	action := d.Get(isNetworkACLRuleAction).(string)

	direction := d.Get(isNetworkACLRuleDirection).(string)
	// creating rule
	name := d.Get(isNetworkACLRuleName).(string)
	source := d.Get(isNetworkACLRuleSource).(string)
	destination := d.Get(isNetworkACLRuleDestination).(string)
	icmp := d.Get(isNetworkACLRuleICMP).([]interface{})
	tcp := d.Get(isNetworkACLRuleTCP).([]interface{})
	udp := d.Get(isNetworkACLRuleUDP).([]interface{})
	icmptype := int64(-1)
	icmpcode := int64(-1)
	minport := int64(-1)
	maxport := int64(-1)
	sourceminport := int64(-1)
	sourcemaxport := int64(-1)
	protocol := "all"

	ruleTemplate := &vpcv1.NetworkACLRulePrototype{
		Action:      &action,
		Destination: &destination,
		Direction:   &direction,
		Source:      &source,
		Name:        &name,
	}

	if before, ok := d.GetOk(isNwACLRuleBefore); ok {
		beforeStr := before.(string)
		ruleTemplate.Before = &vpcv1.NetworkACLRuleBeforePrototype{
			ID: &beforeStr,
		}
	}

	if len(icmp) > 0 {
		protocol = "icmp"
		ruleTemplate.Protocol = &protocol
		if !isNil(icmp[0]) {
			icmpval := icmp[0].(map[string]interface{})
			if val, ok := icmpval[isNetworkACLRuleICMPType]; ok {
				icmptype = int64(val.(int))
				ruleTemplate.Type = &icmptype
			}
			if val, ok := icmpval[isNetworkACLRuleICMPCode]; ok {
				icmpcode = int64(val.(int))
				ruleTemplate.Code = &icmpcode
			}
		}
	} else if len(tcp) > 0 {
		protocol = "tcp"
		ruleTemplate.Protocol = &protocol
		tcpval := tcp[0].(map[string]interface{})
		if val, ok := tcpval[isNetworkACLRulePortMin]; ok {
			minport = int64(val.(int))
			ruleTemplate.DestinationPortMin = &minport
		}
		if val, ok := tcpval[isNetworkACLRulePortMax]; ok {
			maxport = int64(val.(int))
			ruleTemplate.DestinationPortMax = &maxport
		}
		if val, ok := tcpval[isNetworkACLRuleSourcePortMin]; ok {
			sourceminport = int64(val.(int))
			ruleTemplate.SourcePortMin = &sourceminport
		}
		if val, ok := tcpval[isNetworkACLRuleSourcePortMax]; ok {
			sourcemaxport = int64(val.(int))
			ruleTemplate.SourcePortMax = &sourcemaxport
		}
	} else if len(udp) > 0 {
		protocol = "udp"
		ruleTemplate.Protocol = &protocol
		udpval := udp[0].(map[string]interface{})
		if val, ok := udpval[isNetworkACLRulePortMin]; ok {
			minport = int64(val.(int))
			ruleTemplate.DestinationPortMin = &minport
		}
		if val, ok := udpval[isNetworkACLRulePortMax]; ok {
			maxport = int64(val.(int))
			ruleTemplate.DestinationPortMax = &maxport
		}
		if val, ok := udpval[isNetworkACLRuleSourcePortMin]; ok {
			sourceminport = int64(val.(int))
			ruleTemplate.SourcePortMin = &sourceminport
		}
		if val, ok := udpval[isNetworkACLRuleSourcePortMax]; ok {
			sourcemaxport = int64(val.(int))
			ruleTemplate.SourcePortMax = &sourcemaxport
		}
	}
	if protocol == "all" {
		ruleTemplate.Protocol = &protocol
	}

	createNetworkAclRuleOptions := &vpcv1.CreateNetworkACLRuleOptions{
		NetworkACLID:            &nwACLID,
		NetworkACLRulePrototype: ruleTemplate,
	}
	nwaclRule, response, err := sess.CreateNetworkACLRule(createNetworkAclRuleOptions)
	if err != nil || nwaclRule == nil {
		return fmt.Errorf("[ERROR] Error Creating network ACL rule : %s\n%s", err, response)
	}
	err = nwaclRuleGet(d, meta, nwACLID, nwaclRule)
	if err != nil {
		return err
	}
	return nil
}

func resourceIBMISNetworkACLRuleRead(d *schema.ResourceData, meta interface{}) error {
	nwACLID, ruleId, err := parseNwACLTerraformID(d.Id())
	if err != nil {
		return err
	}
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	getNetworkAclRuleOptions := &vpcv1.GetNetworkACLRuleOptions{
		NetworkACLID: &nwACLID,
		ID:           &ruleId,
	}
	nwaclRule, response, err := sess.GetNetworkACLRule(getNetworkAclRuleOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error getting Network ACL Rule (%s) : %s\n%s", ruleId, err, response)
	}
	err = nwaclRuleGet(d, meta, nwACLID, nwaclRule)
	if err != nil {
		return err
	}
	return nil
}

func nwaclRuleGet(d *schema.ResourceData, meta interface{}, nwACLID string, nwaclRule interface{}) error {

	log.Println("[DEBUG] Type of the Rule", reflect.TypeOf(nwaclRule))
	d.Set(isNwACLID, nwACLID)
	switch reflect.TypeOf(nwaclRule).String() {
	case "*vpcv1.NetworkACLRuleNetworkACLRuleProtocolIcmp":
		{
			rulex := nwaclRule.(*vpcv1.NetworkACLRuleNetworkACLRuleProtocolIcmp)
			d.SetId(makeTerraformACLRuleID(nwACLID, *rulex.ID))
			d.Set(isNwACLRuleId, *rulex.ID)
			if rulex.Before != nil {
				d.Set(isNwACLRuleBefore, *rulex.Before.ID)
			}
			d.Set(isNetworkACLRuleName, *rulex.Name)
			d.Set(isNetworkACLRuleHref, *rulex.Href)
			d.Set(isNetworkACLRuleProtocol, *rulex.Protocol)
			d.Set(isNetworkACLRuleAction, *rulex.Action)
			d.Set(isNetworkACLRuleIPVersion, *rulex.IPVersion)
			d.Set(isNetworkACLRuleSource, *rulex.Source)
			d.Set(isNetworkACLRuleDestination, *rulex.Destination)
			d.Set(isNetworkACLRuleDirection, *rulex.Direction)
			d.Set(isNetworkACLRuleTCP, make([]map[string]int, 0, 0))
			d.Set(isNetworkACLRuleUDP, make([]map[string]int, 0, 0))
			icmp := make([]map[string]int, 1, 1)
			if rulex.Code != nil && rulex.Type != nil {
				icmp[0] = map[string]int{
					isNetworkACLRuleICMPCode: int(*rulex.Code),
					isNetworkACLRuleICMPType: int(*rulex.Type),
				}
			}
			d.Set(isNetworkACLRuleICMP, icmp)
		}
	case "*vpcv1.NetworkACLRuleNetworkACLRuleProtocolTcpudp":
		{
			rulex := nwaclRule.(*vpcv1.NetworkACLRuleNetworkACLRuleProtocolTcpudp)
			d.SetId(makeTerraformACLRuleID(nwACLID, *rulex.ID))
			d.Set(isNwACLRuleId, *rulex.ID)
			if rulex.Before != nil {
				d.Set(isNwACLRuleBefore, *rulex.Before.ID)
			}
			d.Set(isNetworkACLRuleHref, *rulex.Href)
			d.Set(isNetworkACLRuleProtocol, *rulex.Protocol)
			d.Set(isNetworkACLRuleName, *rulex.Name)
			d.Set(isNetworkACLRuleAction, *rulex.Action)
			d.Set(isNetworkACLRuleIPVersion, *rulex.IPVersion)
			d.Set(isNetworkACLRuleSource, *rulex.Source)
			d.Set(isNetworkACLRuleDestination, *rulex.Destination)
			d.Set(isNetworkACLRuleDirection, *rulex.Direction)
			if *rulex.Protocol == "tcp" {
				d.Set(isNetworkACLRuleICMP, make([]map[string]int, 0, 0))
				d.Set(isNetworkACLRuleUDP, make([]map[string]int, 0, 0))
				tcp := make([]map[string]int, 1, 1)
				tcp[0] = map[string]int{
					isNetworkACLRuleSourcePortMax: checkNetworkACLNil(rulex.SourcePortMax),
					isNetworkACLRuleSourcePortMin: checkNetworkACLNil(rulex.SourcePortMin),
				}
				tcp[0][isNetworkACLRulePortMax] = checkNetworkACLNil(rulex.DestinationPortMax)
				tcp[0][isNetworkACLRulePortMin] = checkNetworkACLNil(rulex.DestinationPortMin)
				d.Set(isNetworkACLRuleTCP, tcp)
			} else if *rulex.Protocol == "udp" {
				d.Set(isNetworkACLRuleICMP, make([]map[string]int, 0, 0))
				d.Set(isNetworkACLRuleTCP, make([]map[string]int, 0, 0))
				udp := make([]map[string]int, 1, 1)
				udp[0] = map[string]int{
					isNetworkACLRuleSourcePortMax: checkNetworkACLNil(rulex.SourcePortMax),
					isNetworkACLRuleSourcePortMin: checkNetworkACLNil(rulex.SourcePortMin),
				}
				udp[0][isNetworkACLRulePortMax] = checkNetworkACLNil(rulex.DestinationPortMax)
				udp[0][isNetworkACLRulePortMin] = checkNetworkACLNil(rulex.DestinationPortMin)
				d.Set(isNetworkACLRuleUDP, udp)
			}
		}
	case "*vpcv1.NetworkACLRuleNetworkACLRuleProtocolAll":
		{
			rulex := nwaclRule.(*vpcv1.NetworkACLRuleNetworkACLRuleProtocolAll)
			d.SetId(makeTerraformACLRuleID(nwACLID, *rulex.ID))
			d.Set(isNwACLRuleId, *rulex.ID)
			if rulex.Before != nil {
				d.Set(isNwACLRuleBefore, *rulex.Before.ID)
			}
			d.Set(isNetworkACLRuleHref, *rulex.Href)
			d.Set(isNetworkACLRuleProtocol, *rulex.Protocol)
			d.Set(isNetworkACLRuleName, *rulex.Name)
			d.Set(isNetworkACLRuleAction, *rulex.Action)
			d.Set(isNetworkACLRuleIPVersion, *rulex.IPVersion)
			d.Set(isNetworkACLRuleSource, *rulex.Source)
			d.Set(isNetworkACLRuleDestination, *rulex.Destination)
			d.Set(isNetworkACLRuleDirection, *rulex.Direction)
			d.Set(isNetworkACLRuleICMP, make([]map[string]int, 0, 0))
			d.Set(isNetworkACLRuleTCP, make([]map[string]int, 0, 0))
			d.Set(isNetworkACLRuleUDP, make([]map[string]int, 0, 0))
		}
	}
	return nil
}

func resourceIBMISNetworkACLRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	nwACLId, ruleId, err := parseNwACLTerraformID(id)

	err = nwaclRuleUpdate(d, meta, ruleId, nwACLId)
	if err != nil {
		return err
	}
	return resourceIBMISNetworkACLRuleRead(d, meta)
}

func nwaclRuleUpdate(d *schema.ResourceData, meta interface{}, id, nwACLId string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	updateNetworkACLRuleOptions := &vpcv1.UpdateNetworkACLRuleOptions{
		NetworkACLID: &nwACLId,
		ID:           &id,
	}
	updateNetworkACLOptionsPatchModel := &vpcv1.NetworkACLRulePatch{}

	hasChanged := false

	if d.HasChange(isNetworkACLRuleAction) {
		hasChanged = true
		if actionVar, ok := d.GetOk(isNetworkACLRuleAction); ok {
			action := actionVar.(string)
			updateNetworkACLOptionsPatchModel.Action = &action
		}
	}
	aclRuleBeforeNull := false
	if d.HasChange(isNwACLRuleBefore) {
		hasChanged = true
		beforeVar := d.Get(isNwACLRuleBefore).(string)
		if beforeVar == "null" {
			aclRuleBeforeNull = true
		} else if beforeVar != "" {
			updateNetworkACLOptionsPatchModel.Before = &vpcv1.NetworkACLRuleBeforePatchNetworkACLRuleIdentityByID{
				ID: &beforeVar,
			}
		}
	}

	if d.HasChange(isNetworkACLRuleName) {
		hasChanged = true
		if nameVar, ok := d.GetOk(isNetworkACLRuleName); ok {
			nameStr := nameVar.(string)
			updateNetworkACLOptionsPatchModel.Name = &nameStr
		}
	}
	if d.HasChange(isNetworkACLRuleDirection) {
		hasChanged = true
		if directionVar, ok := d.GetOk(isNetworkACLRuleDirection); ok {
			directionStr := directionVar.(string)
			updateNetworkACLOptionsPatchModel.Direction = &directionStr
		}
	}
	if d.HasChange(isNetworkACLRuleDestination) {
		hasChanged = true
		if destinationVar, ok := d.GetOk(isNetworkACLRuleDestination); ok {
			destination := destinationVar.(string)
			updateNetworkACLOptionsPatchModel.Destination = &destination
		}
	}
	if d.HasChange(isNetworkACLRuleICMP) {
		icmpCode := fmt.Sprint(isNetworkACLRuleICMP, ".0.", isNetworkACLRuleICMPCode)
		icmpType := fmt.Sprint(isNetworkACLRuleICMP, ".0.", isNetworkACLRuleICMPType)
		if d.HasChange(icmpCode) {
			hasChanged = true
			if codeVar, ok := d.GetOk(icmpCode); ok {
				code := int64(codeVar.(int))
				updateNetworkACLOptionsPatchModel.Code = &code
			}
		}
		if d.HasChange(icmpType) {
			hasChanged = true
			if typeVar, ok := d.GetOk(icmpType); ok {
				typeInt := int64(typeVar.(int))
				updateNetworkACLOptionsPatchModel.Type = &typeInt
			}
		}
	}
	if d.HasChange(isNetworkACLRuleTCP) {
		tcp := d.Get(isNetworkACLRuleTCP).([]interface{})
		tcpval := tcp[0].(map[string]interface{})
		max := fmt.Sprint(isNetworkACLRuleTCP, ".0.", isNetworkACLRulePortMax)
		min := fmt.Sprint(isNetworkACLRuleTCP, ".0.", isNetworkACLRulePortMin)
		maxSource := fmt.Sprint(isNetworkACLRuleTCP, ".0.", isNetworkACLRuleSourcePortMax)
		minSource := fmt.Sprint(isNetworkACLRuleTCP, ".0.", isNetworkACLRuleSourcePortMin)
		if d.HasChange(max) {
			hasChanged = true
			if destinationVar, ok := tcpval[isNetworkACLRulePortMax]; ok {
				destination := int64(destinationVar.(int))
				updateNetworkACLOptionsPatchModel.DestinationPortMax = &destination
			}
		}
		if d.HasChange(min) {
			hasChanged = true
			if destinationVar, ok := tcpval[isNetworkACLRulePortMin]; ok {
				destination := int64(destinationVar.(int))
				updateNetworkACLOptionsPatchModel.DestinationPortMin = &destination
			}
		}
		if d.HasChange(maxSource) {
			hasChanged = true
			if sourceVar, ok := tcpval[isNetworkACLRuleSourcePortMax]; ok {
				source := int64(sourceVar.(int))
				updateNetworkACLOptionsPatchModel.SourcePortMax = &source
			}
		}
		if d.HasChange(minSource) {
			hasChanged = true
			if sourceVar, ok := tcpval[isNetworkACLRuleSourcePortMin]; ok {
				source := int64(sourceVar.(int))
				updateNetworkACLOptionsPatchModel.SourcePortMin = &source
			}
		}
	}
	if d.HasChange(isNetworkACLRuleUDP) {
		udp := d.Get(isNetworkACLRuleUDP).([]interface{})
		udpval := udp[0].(map[string]interface{})
		max := fmt.Sprint(isNetworkACLRuleUDP, ".0.", isNetworkACLRulePortMax)
		min := fmt.Sprint(isNetworkACLRuleUDP, ".0.", isNetworkACLRulePortMin)
		maxSource := fmt.Sprint(isNetworkACLRuleUDP, ".0.", isNetworkACLRuleSourcePortMax)
		minSource := fmt.Sprint(isNetworkACLRuleUDP, ".0.", isNetworkACLRuleSourcePortMin)

		if d.HasChange(max) {
			hasChanged = true
			if destinationVar, ok := udpval[isNetworkACLRulePortMax]; ok {
				destination := int64(destinationVar.(int))
				updateNetworkACLOptionsPatchModel.DestinationPortMax = &destination
			}
		}
		if d.HasChange(min) {
			hasChanged = true
			if destinationVar, ok := udpval[isNetworkACLRulePortMin]; ok {
				destination := int64(destinationVar.(int))
				updateNetworkACLOptionsPatchModel.DestinationPortMin = &destination
			}
		}
		if d.HasChange(maxSource) {
			hasChanged = true
			if sourceVar, ok := udpval[isNetworkACLRuleSourcePortMax]; ok {
				source := int64(sourceVar.(int))
				updateNetworkACLOptionsPatchModel.SourcePortMax = &source
			}
		}
		if d.HasChange(minSource) {
			hasChanged = true
			if sourceVar, ok := udpval[isNetworkACLRuleSourcePortMin]; ok {
				source := int64(sourceVar.(int))
				updateNetworkACLOptionsPatchModel.SourcePortMin = &source
			}
		}
	}

	if d.HasChange(isNetworkACLRuleSource) {
		hasChanged = true
		if sourceVar, ok := d.GetOk(isNetworkACLRuleSource); ok {
			source := sourceVar.(string)
			updateNetworkACLOptionsPatchModel.Source = &source
		}
	}

	if hasChanged {
		updateNetworkACLOptionsPatch, err := updateNetworkACLOptionsPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("[ERROR] Error calling asPatch for NetworkACLOptionsPatch : %s", err)
		}
		if aclRuleBeforeNull {
			updateNetworkACLOptionsPatch["before"] = nil
		}
		updateNetworkACLRuleOptions.NetworkACLRulePatch = updateNetworkACLOptionsPatch
		_, response, err := sess.UpdateNetworkACLRule(updateNetworkACLRuleOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Updating Network ACL Rule : %s\n%s", err, response)
		}
	}
	return nil
}

func resourceIBMISNetworkACLRuleDelete(d *schema.ResourceData, meta interface{}) error {
	nwACLID, ruleId, err := parseNwACLTerraformID(d.Id())
	if err != nil {
		return err
	}

	err = nwaclRuleDelete(d, meta, ruleId, nwACLID)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func nwaclRuleDelete(d *schema.ResourceData, meta interface{}, id, nwACLId string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	getNetworkAclRuleOptions := &vpcv1.GetNetworkACLRuleOptions{
		NetworkACLID: &nwACLId,
		ID:           &id,
	}
	_, response, err := sess.GetNetworkACLRule(getNetworkAclRuleOptions)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error Getting Network ACL Rule  (%s): %s\n%s", id, err, response)
	}

	deleteNetworkAclRuleOptions := &vpcv1.DeleteNetworkACLRuleOptions{
		NetworkACLID: &nwACLId,
		ID:           &id,
	}
	response, err = sess.DeleteNetworkACLRule(deleteNetworkAclRuleOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Deleting Network ACL Rule : %s\n%s", err, response)
	}
	d.SetId("")
	return nil
}

func nwaclRuleExists(d *schema.ResourceData, meta interface{}, id, nwACLId string) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}
	getNetworkAclRuleOptions := &vpcv1.GetNetworkACLRuleOptions{
		NetworkACLID: &nwACLId,
		ID:           &id,
	}
	_, response, err := sess.GetNetworkACLRule(getNetworkAclRuleOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error getting Network ACL Rule: %s\n%s", err, response)
	}
	return true, nil
}

func makeTerraformACLRuleID(id1, id2 string) string {
	// Include both network acl id and rule id to create a unique Terraform id.  As a bonus,
	// we can extract the network acl id as needed for API calls such as READ.
	return fmt.Sprintf("%s/%s", id1, id2)
}

func parseNwACLTerraformID(s string) (string, string, error) {
	segments := strings.Split(s, "/")
	if len(segments) != 2 {
		return "", "", fmt.Errorf("invalid terraform Id %s (incorrect number of segments)", s)
	}
	if segments[0] == "" || segments[1] == "" {
		return "", "", fmt.Errorf("invalid terraform Id %s (one or more empty segments)", s)
	}
	return segments[0], segments[1], nil
}
