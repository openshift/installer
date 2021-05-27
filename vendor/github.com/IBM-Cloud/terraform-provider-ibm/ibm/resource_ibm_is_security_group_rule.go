// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/IBM/vpc-go-sdk/vpcclassicv1"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isSecurityGroupRuleCode             = "code"
	isSecurityGroupRuleDirection        = "direction"
	isSecurityGroupRuleIPVersion        = "ip_version"
	isSecurityGroupRuleIPVersionDefault = "ipv4"
	isSecurityGroupRulePortMax          = "port_max"
	isSecurityGroupRulePortMin          = "port_min"
	isSecurityGroupRuleProtocolICMP     = "icmp"
	isSecurityGroupRuleProtocolTCP      = "tcp"
	isSecurityGroupRuleProtocolUDP      = "udp"
	isSecurityGroupRuleProtocol         = "protocol"
	isSecurityGroupRuleRemote           = "remote"
	isSecurityGroupRuleType             = "type"
	isSecurityGroupID                   = "group"
	isSecurityGroupRuleID               = "rule_id"
)

func resourceIBMISSecurityGroupRule() *schema.Resource {

	return &schema.Resource{
		Create:   resourceIBMISSecurityGroupRuleCreate,
		Read:     resourceIBMISSecurityGroupRuleRead,
		Update:   resourceIBMISSecurityGroupRuleUpdate,
		Delete:   resourceIBMISSecurityGroupRuleDelete,
		Exists:   resourceIBMISSecurityGroupRuleExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{

			isSecurityGroupID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Security group id",
				ForceNew:    true,
			},

			isSecurityGroupRuleID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule id",
			},

			isSecurityGroupRuleDirection: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Direction of traffic to enforce, either inbound or outbound",
				ValidateFunc: InvokeValidator("ibm_is_security_group_rule", isSecurityGroupRuleDirection),
			},

			isSecurityGroupRuleIPVersion: {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "IP version: ipv4 or ipv6",
				Default:      isSecurityGroupRuleIPVersionDefault,
				ValidateFunc: InvokeValidator("ibm_is_security_group_rule", isSecurityGroupRuleIPVersion),
			},

			isSecurityGroupRuleRemote: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Security group id: an IP address, a CIDR block, or a single security group identifier",
			},

			isSecurityGroupRuleProtocolICMP: {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ForceNew:      true,
				MinItems:      1,
				ConflictsWith: []string{isSecurityGroupRuleProtocolTCP, isSecurityGroupRuleProtocolUDP},
				Description:   "protocol=icmp",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isSecurityGroupRuleType: {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     false,
							ValidateFunc: InvokeValidator("ibm_is_security_group_rule", isSecurityGroupRuleType),
						},
						isSecurityGroupRuleCode: {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     false,
							ValidateFunc: InvokeValidator("ibm_is_security_group_rule", isSecurityGroupRuleCode),
						},
					},
				},
			},

			isSecurityGroupRuleProtocolTCP: {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				MinItems:      1,
				ForceNew:      true,
				Description:   "protocol=tcp",
				ConflictsWith: []string{isSecurityGroupRuleProtocolUDP, isSecurityGroupRuleProtocolICMP},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isSecurityGroupRulePortMin: {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     false,
							Default:      1,
							ValidateFunc: InvokeValidator("ibm_is_security_group_rule", isSecurityGroupRulePortMin),
						},
						isSecurityGroupRulePortMax: {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     false,
							Default:      65535,
							ValidateFunc: InvokeValidator("ibm_is_security_group_rule", isSecurityGroupRulePortMax),
						},
					},
				},
			},

			isSecurityGroupRuleProtocolUDP: {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ForceNew:      true,
				MinItems:      1,
				Description:   "protocol=udp",
				ConflictsWith: []string{isSecurityGroupRuleProtocolTCP, isSecurityGroupRuleProtocolICMP},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isSecurityGroupRulePortMin: {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     false,
							Default:      1,
							ValidateFunc: InvokeValidator("ibm_is_security_group_rule", isSecurityGroupRulePortMin),
						},
						isSecurityGroupRulePortMax: {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     false,
							Default:      65535,
							ValidateFunc: InvokeValidator("ibm_is_security_group_rule", isSecurityGroupRulePortMax),
						},
					},
				},
			},

			RelatedCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the Security Group",
			},
			isSecurityGroupRuleProtocol: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Security Group Rule Protocol",
			},
		},
	}
}

func resourceIBMISSecurityGroupRuleValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 1)
	direction := "inbound, outbound"
	ip_version := "ipv4, ipv6"

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isSecurityGroupRuleDirection,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              direction})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isSecurityGroupRuleIPVersion,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              ip_version})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isSecurityGroupRuleType,
			ValidateFunctionIdentifier: IntBetween,
			Type:                       TypeInt,
			MinValue:                   "0",
			MaxValue:                   "254"})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isSecurityGroupRuleCode,
			ValidateFunctionIdentifier: IntBetween,
			Type:                       TypeInt,
			MinValue:                   "0",
			MaxValue:                   "255"})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isSecurityGroupRulePortMin,
			ValidateFunctionIdentifier: IntBetween,
			Type:                       TypeInt,
			MinValue:                   "1",
			MaxValue:                   "65535"})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isSecurityGroupRulePortMax,
			ValidateFunctionIdentifier: IntBetween,
			Type:                       TypeInt,
			MinValue:                   "1",
			MaxValue:                   "65535"})

	ibmISSecurityGroupRuleResourceValidator := ResourceValidator{ResourceName: "ibm_is_security_group_rule", Schema: validateSchema}
	return &ibmISSecurityGroupRuleResourceValidator
}

func resourceIBMISSecurityGroupRuleCreate(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	if userDetails.generation == 1 {
		err := classicSgRuleCreate(d, meta)
		if err != nil {
			return err
		}
	} else {
		err := sgRuleCreate(d, meta)
		if err != nil {
			return err
		}
	}
	return err
}

func classicSgRuleCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	parsed, sgTemplate, _, err := parseIBMISClassicSecurityGroupRuleDictionary(d, "create", sess)
	if err != nil {
		return err
	}
	isSecurityGroupRuleKey := "security_group_rule_key_" + parsed.secgrpID
	ibmMutexKV.Lock(isSecurityGroupRuleKey)
	defer ibmMutexKV.Unlock(isSecurityGroupRuleKey)

	options := &vpcclassicv1.CreateSecurityGroupRuleOptions{
		SecurityGroupID:            &parsed.secgrpID,
		SecurityGroupRulePrototype: sgTemplate,
	}

	rule, response, err := sess.CreateSecurityGroupRule(options)
	if err != nil {
		return fmt.Errorf("Error while creating Security Group Rule %s\n%s", err, response)
	}
	switch reflect.TypeOf(rule).String() {
	case "*vpcclassicv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp":
		{
			sgrule := rule.(*vpcclassicv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp)
			d.Set(isSecurityGroupRuleID, *sgrule.ID)
			tfID := makeTerraformRuleID(parsed.secgrpID, *sgrule.ID)
			d.SetId(tfID)
		}
	case "*vpcclassicv1.SecurityGroupRuleSecurityGroupRuleProtocolAll":
		{
			sgrule := rule.(*vpcclassicv1.SecurityGroupRuleSecurityGroupRuleProtocolAll)
			d.Set(isSecurityGroupRuleID, *sgrule.ID)
			tfID := makeTerraformRuleID(parsed.secgrpID, *sgrule.ID)
			d.SetId(tfID)
		}
	case "*vpcclassicv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp":
		{
			sgrule := rule.(*vpcclassicv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp)
			d.Set(isSecurityGroupRuleID, *sgrule.ID)
			tfID := makeTerraformRuleID(parsed.secgrpID, *sgrule.ID)
			d.SetId(tfID)
		}
	}
	return nil
}

func sgRuleCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	parsed, sgTemplate, _, err := parseIBMISSecurityGroupRuleDictionary(d, "create", sess)
	if err != nil {
		return err
	}
	isSecurityGroupRuleKey := "security_group_rule_key_" + parsed.secgrpID
	ibmMutexKV.Lock(isSecurityGroupRuleKey)
	defer ibmMutexKV.Unlock(isSecurityGroupRuleKey)

	options := &vpcv1.CreateSecurityGroupRuleOptions{
		SecurityGroupID:            &parsed.secgrpID,
		SecurityGroupRulePrototype: sgTemplate,
	}

	rule, response, err := sess.CreateSecurityGroupRule(options)
	if err != nil {
		return fmt.Errorf("Error while creating Security Group Rule %s\n%s", err, response)
	}
	switch reflect.TypeOf(rule).String() {
	case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp":
		{
			sgrule := rule.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp)
			d.Set(isSecurityGroupRuleID, *sgrule.ID)
			tfID := makeTerraformRuleID(parsed.secgrpID, *sgrule.ID)
			d.SetId(tfID)
		}
	case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll":
		{
			sgrule := rule.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll)
			d.Set(isSecurityGroupRuleID, *sgrule.ID)
			tfID := makeTerraformRuleID(parsed.secgrpID, *sgrule.ID)
			d.SetId(tfID)
		}
	case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp":
		{
			sgrule := rule.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp)
			d.Set(isSecurityGroupRuleID, *sgrule.ID)
			tfID := makeTerraformRuleID(parsed.secgrpID, *sgrule.ID)
			d.SetId(tfID)
		}
	}
	return nil
}

func resourceIBMISSecurityGroupRuleRead(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	secgrpID, ruleID, err := parseISTerraformID(d.Id())
	if err != nil {
		return err
	}
	if userDetails.generation == 1 {
		err := classicSgRuleGet(d, meta, secgrpID, ruleID)
		if err != nil {
			return err
		}
	} else {
		err := sgRuleGet(d, meta, secgrpID, ruleID)
		if err != nil {
			return err
		}
	}
	return nil
}

func classicSgRuleGet(d *schema.ResourceData, meta interface{}, secgrpID, ruleID string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	getSecurityGroupRuleOptions := &vpcclassicv1.GetSecurityGroupRuleOptions{
		SecurityGroupID: &secgrpID,
		ID:              &ruleID,
	}
	sgrule, response, err := sess.GetSecurityGroupRule(getSecurityGroupRuleOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting Security Group Rule (%s): %s\n%s", ruleID, err, response)
	}

	d.Set(isSecurityGroupID, secgrpID)
	getSecurityGroupOptions := &vpcclassicv1.GetSecurityGroupOptions{
		ID: &secgrpID,
	}
	sg, response, err := sess.GetSecurityGroup(getSecurityGroupOptions)
	if err != nil {
		return fmt.Errorf("Error Getting Security Group : %s\n%s", err, response)
	}
	d.Set(RelatedCRN, *sg.CRN)
	switch reflect.TypeOf(sgrule).String() {
	case "*vpcclassicv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp":
		{
			rule := sgrule.(*vpcclassicv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp)
			d.Set(isSecurityGroupRuleID, *rule.ID)
			tfID := makeTerraformRuleID(secgrpID, *rule.ID)
			d.SetId(tfID)
			d.Set(isSecurityGroupRuleIPVersion, *rule.IPVersion)
			d.Set(isSecurityGroupRuleProtocol, *rule.Protocol)
			icmpProtocol := map[string]interface{}{}

			if rule.Type != nil {
				icmpProtocol["type"] = *rule.Type
			}
			if rule.Code != nil {
				icmpProtocol["code"] = *rule.Code
			}
			protocolList := make([]map[string]interface{}, 0)
			protocolList = append(protocolList, icmpProtocol)
			d.Set(isSecurityGroupRuleProtocolICMP, protocolList)
			remote, ok := rule.Remote.(*vpcclassicv1.SecurityGroupRuleRemote)
			if ok {
				if remote != nil && reflect.ValueOf(remote).IsNil() == false {
					if remote.ID != nil {
						d.Set(isSecurityGroupRuleRemote, remote.ID)
					} else if remote.Address != nil {
						d.Set(isSecurityGroupRuleRemote, remote.Address)
					} else if remote.CIDRBlock != nil {
						d.Set(isSecurityGroupRuleRemote, remote.CIDRBlock)
					}
				}
			}
		}
	case "*vpcclassicv1.SecurityGroupRuleSecurityGroupRuleProtocolAll":
		{
			rule := sgrule.(*vpcclassicv1.SecurityGroupRuleSecurityGroupRuleProtocolAll)
			d.Set(isSecurityGroupRuleID, *rule.ID)
			tfID := makeTerraformRuleID(secgrpID, *rule.ID)
			d.SetId(tfID)
			d.Set(isSecurityGroupRuleIPVersion, *rule.IPVersion)
			d.Set(isSecurityGroupRuleProtocol, *rule.Protocol)
			remote, ok := rule.Remote.(*vpcclassicv1.SecurityGroupRuleRemote)
			if ok {
				if remote != nil && reflect.ValueOf(remote).IsNil() == false {
					if remote.ID != nil {
						d.Set(isSecurityGroupRuleRemote, remote.ID)
					} else if remote.Address != nil {
						d.Set(isSecurityGroupRuleRemote, remote.Address)
					} else if remote.CIDRBlock != nil {
						d.Set(isSecurityGroupRuleRemote, remote.CIDRBlock)
					}
				}
			}
		}
	case "*vpcclassicv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp":
		{
			rule := sgrule.(*vpcclassicv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp)
			d.Set(isSecurityGroupRuleID, *rule.ID)
			tfID := makeTerraformRuleID(secgrpID, *rule.ID)
			d.SetId(tfID)
			d.Set(isSecurityGroupRuleIPVersion, *rule.IPVersion)
			d.Set(isSecurityGroupRuleProtocol, *rule.Protocol)
			tcpProtocol := map[string]interface{}{}

			if rule.PortMin != nil {
				tcpProtocol["port_min"] = *rule.PortMin
			}
			if rule.PortMax != nil {
				tcpProtocol["port_max"] = *rule.PortMax
			}
			protocolList := make([]map[string]interface{}, 0)
			protocolList = append(protocolList, tcpProtocol)
			if *rule.Protocol == isSecurityGroupRuleProtocolTCP {
				d.Set(isSecurityGroupRuleProtocolTCP, protocolList)
			} else {
				d.Set(isSecurityGroupRuleProtocolUDP, protocolList)
			}
			remote, ok := rule.Remote.(*vpcclassicv1.SecurityGroupRuleRemote)
			if ok {
				if remote != nil && reflect.ValueOf(remote).IsNil() == false {
					if remote.ID != nil {
						d.Set(isSecurityGroupRuleRemote, remote.ID)
					} else if remote.Address != nil {
						d.Set(isSecurityGroupRuleRemote, remote.Address)
					} else if remote.CIDRBlock != nil {
						d.Set(isSecurityGroupRuleRemote, remote.CIDRBlock)
					}
				}
			}
		}
	}
	return nil
}

func sgRuleGet(d *schema.ResourceData, meta interface{}, secgrpID, ruleID string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	getSecurityGroupRuleOptions := &vpcv1.GetSecurityGroupRuleOptions{
		SecurityGroupID: &secgrpID,
		ID:              &ruleID,
	}
	sgrule, response, err := sess.GetSecurityGroupRule(getSecurityGroupRuleOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting Security Group Rule (%s): %s\n%s", ruleID, err, response)
	}
	d.Set(isSecurityGroupID, secgrpID)
	getSecurityGroupOptions := &vpcv1.GetSecurityGroupOptions{
		ID: &secgrpID,
	}
	sg, response, err := sess.GetSecurityGroup(getSecurityGroupOptions)
	if err != nil {
		return fmt.Errorf("Error Getting Security Group : %s\n%s", err, response)
	}
	d.Set(RelatedCRN, *sg.CRN)
	switch reflect.TypeOf(sgrule).String() {
	case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp":
		{
			rule := sgrule.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp)
			d.Set(isSecurityGroupRuleID, *rule.ID)
			tfID := makeTerraformRuleID(secgrpID, *rule.ID)
			d.SetId(tfID)
			d.Set(isSecurityGroupRuleIPVersion, *rule.IPVersion)
			d.Set(isSecurityGroupRuleProtocol, *rule.Protocol)
			icmpProtocol := map[string]interface{}{}

			if rule.Type != nil {
				icmpProtocol["type"] = *rule.Type
			}
			if rule.Code != nil {
				icmpProtocol["code"] = *rule.Code
			}
			protocolList := make([]map[string]interface{}, 0)
			protocolList = append(protocolList, icmpProtocol)
			d.Set(isSecurityGroupRuleProtocolICMP, protocolList)
			remote, ok := rule.Remote.(*vpcv1.SecurityGroupRuleRemote)
			if ok {
				if remote != nil && reflect.ValueOf(remote).IsNil() == false {
					if remote.ID != nil {
						d.Set(isSecurityGroupRuleRemote, remote.ID)
					} else if remote.Address != nil {
						d.Set(isSecurityGroupRuleRemote, remote.Address)
					} else if remote.CIDRBlock != nil {
						d.Set(isSecurityGroupRuleRemote, remote.CIDRBlock)
					}
				}
			}
		}
	case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll":
		{
			rule := sgrule.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll)
			d.Set(isSecurityGroupRuleID, *rule.ID)
			tfID := makeTerraformRuleID(secgrpID, *rule.ID)
			d.SetId(tfID)
			d.Set(isSecurityGroupRuleIPVersion, *rule.IPVersion)
			d.Set(isSecurityGroupRuleProtocol, *rule.Protocol)
			remote, ok := rule.Remote.(*vpcv1.SecurityGroupRuleRemote)
			if ok {
				if remote != nil && reflect.ValueOf(remote).IsNil() == false {
					if remote.ID != nil {
						d.Set(isSecurityGroupRuleRemote, remote.ID)
					} else if remote.Address != nil {
						d.Set(isSecurityGroupRuleRemote, remote.Address)
					} else if remote.CIDRBlock != nil {
						d.Set(isSecurityGroupRuleRemote, remote.CIDRBlock)
					}
				}
			}
		}
	case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp":
		{
			rule := sgrule.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp)
			d.Set(isSecurityGroupRuleID, *rule.ID)
			tfID := makeTerraformRuleID(secgrpID, *rule.ID)
			d.SetId(tfID)
			d.Set(isSecurityGroupRuleIPVersion, *rule.IPVersion)
			d.Set(isSecurityGroupRuleProtocol, *rule.Protocol)
			tcpProtocol := map[string]interface{}{}

			if rule.PortMin != nil {
				tcpProtocol["port_min"] = *rule.PortMin
			}
			if rule.PortMax != nil {
				tcpProtocol["port_max"] = *rule.PortMax
			}
			protocolList := make([]map[string]interface{}, 0)
			protocolList = append(protocolList, tcpProtocol)
			if *rule.Protocol == isSecurityGroupRuleProtocolTCP {
				d.Set(isSecurityGroupRuleProtocolTCP, protocolList)
			} else {
				d.Set(isSecurityGroupRuleProtocolUDP, protocolList)
			}
			remote, ok := rule.Remote.(*vpcv1.SecurityGroupRuleRemote)
			if ok {
				if remote != nil && reflect.ValueOf(remote).IsNil() == false {
					if remote.ID != nil {
						d.Set(isSecurityGroupRuleRemote, remote.ID)
					} else if remote.Address != nil {
						d.Set(isSecurityGroupRuleRemote, remote.Address)
					} else if remote.CIDRBlock != nil {
						d.Set(isSecurityGroupRuleRemote, remote.CIDRBlock)
					}
				}
			}
		}
	}
	return nil
}

func resourceIBMISSecurityGroupRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	if userDetails.generation == 1 {
		err := classicSgRuleUpdate(d, meta)
		if err != nil {
			return err
		}
	} else {
		err := sgRuleUpdate(d, meta)
		if err != nil {
			return err
		}
	}
	return resourceIBMISSecurityGroupRuleRead(d, meta)
}

func classicSgRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	parsed, _, sgTemplate, err := parseIBMISClassicSecurityGroupRuleDictionary(d, "update", sess)
	if err != nil {
		return err
	}
	isSecurityGroupRuleKey := "security_group_rule_key_" + parsed.secgrpID
	ibmMutexKV.Lock(isSecurityGroupRuleKey)
	defer ibmMutexKV.Unlock(isSecurityGroupRuleKey)
	securityGroupRulePatchBody, _ := sgTemplate.AsPatch()
	updateSecurityGroupRuleOptions := &vpcclassicv1.UpdateSecurityGroupRuleOptions{
		SecurityGroupID:        &parsed.secgrpID,
		ID:                     &parsed.ruleID,
		SecurityGroupRulePatch: securityGroupRulePatchBody,
	}
	_, response, err := sess.UpdateSecurityGroupRule(updateSecurityGroupRuleOptions)
	if err != nil {
		return fmt.Errorf("Error Updating Security Group Rule : %s\n%s", err, response)
	}
	return nil
}

func sgRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	parsed, _, sgTemplate, err := parseIBMISSecurityGroupRuleDictionary(d, "update", sess)
	if err != nil {
		return err
	}
	isSecurityGroupRuleKey := "security_group_rule_key_" + parsed.secgrpID
	ibmMutexKV.Lock(isSecurityGroupRuleKey)
	defer ibmMutexKV.Unlock(isSecurityGroupRuleKey)

	updateSecurityGroupRuleOptions := sgTemplate
	_, response, err := sess.UpdateSecurityGroupRule(updateSecurityGroupRuleOptions)
	if err != nil {
		return fmt.Errorf("Error Updating Security Group Rule : %s\n%s", err, response)
	}
	return nil
}

func resourceIBMISSecurityGroupRuleDelete(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	secgrpID, ruleID, err := parseISTerraformID(d.Id())
	if err != nil {
		return err
	}

	isSecurityGroupRuleKey := "security_group_rule_key_" + secgrpID
	ibmMutexKV.Lock(isSecurityGroupRuleKey)
	defer ibmMutexKV.Unlock(isSecurityGroupRuleKey)

	if userDetails.generation == 1 {
		err := classicSgRuleDelete(d, meta, secgrpID, ruleID)
		if err != nil {
			return err
		}
	} else {
		err := sgRuleDelete(d, meta, secgrpID, ruleID)
		if err != nil {
			return err
		}
	}
	d.SetId("")
	return nil
}

func classicSgRuleDelete(d *schema.ResourceData, meta interface{}, secgrpID, ruleID string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}

	getSecurityGroupRuleOptions := &vpcclassicv1.GetSecurityGroupRuleOptions{
		SecurityGroupID: &secgrpID,
		ID:              &ruleID,
	}
	_, response, err := sess.GetSecurityGroupRule(getSecurityGroupRuleOptions)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting Security Group Rule (%s): %s\n%s", ruleID, err, response)
	}

	deleteSecurityGroupRuleOptions := &vpcclassicv1.DeleteSecurityGroupRuleOptions{
		SecurityGroupID: &secgrpID,
		ID:              &ruleID,
	}
	response, err = sess.DeleteSecurityGroupRule(deleteSecurityGroupRuleOptions)
	if err != nil {
		return fmt.Errorf("Error Deleting Security Group Rule : %s\n%s", err, response)
	}
	d.SetId("")
	return nil
}

func sgRuleDelete(d *schema.ResourceData, meta interface{}, secgrpID, ruleID string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	getSecurityGroupRuleOptions := &vpcv1.GetSecurityGroupRuleOptions{
		SecurityGroupID: &secgrpID,
		ID:              &ruleID,
	}
	_, response, err := sess.GetSecurityGroupRule(getSecurityGroupRuleOptions)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting Security Group Rule (%s): %s\n%s", ruleID, err, response)
	}

	deleteSecurityGroupRuleOptions := &vpcv1.DeleteSecurityGroupRuleOptions{
		SecurityGroupID: &secgrpID,
		ID:              &ruleID,
	}
	response, err = sess.DeleteSecurityGroupRule(deleteSecurityGroupRuleOptions)
	if err != nil {
		return fmt.Errorf("Error Deleting Security Group Rule : %s\n%s", err, response)
	}
	d.SetId("")
	return nil
}

func resourceIBMISSecurityGroupRuleExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return false, err
	}
	secgrpID, ruleID, err := parseISTerraformID(d.Id())
	if err != nil {
		return false, err
	}
	if userDetails.generation == 1 {
		exists, err := classicSgRuleExists(d, meta, secgrpID, ruleID)
		return exists, err
	} else {
		exists, err := sgRuleExists(d, meta, secgrpID, ruleID)
		return exists, err
	}
}

func classicSgRuleExists(d *schema.ResourceData, meta interface{}, secgrpID, ruleID string) (bool, error) {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return false, err
	}
	getSecurityGroupRuleOptions := &vpcclassicv1.GetSecurityGroupRuleOptions{
		SecurityGroupID: &secgrpID,
		ID:              &ruleID,
	}
	_, response, err := sess.GetSecurityGroupRule(getSecurityGroupRuleOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error getting Security Group Rule (%s): %s\n%s", ruleID, err, response)
	}
	return true, nil
}

func sgRuleExists(d *schema.ResourceData, meta interface{}, secgrpID, ruleID string) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}
	getSecurityGroupRuleOptions := &vpcv1.GetSecurityGroupRuleOptions{
		SecurityGroupID: &secgrpID,
		ID:              &ruleID,
	}
	_, response, err := sess.GetSecurityGroupRule(getSecurityGroupRuleOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error getting Security Group Rule (%s): %s\n%s", ruleID, err, response)
	}
	return true, nil
}

func parseISTerraformID(s string) (string, string, error) {
	segments := strings.Split(s, ".")
	if len(segments) != 2 {
		return "", "", fmt.Errorf("invalid terraform Id %s (incorrect number of segments)", s)
	}
	if segments[0] == "" || segments[1] == "" {
		return "", "", fmt.Errorf("invalid terraform Id %s (one or more empty segments)", s)
	}
	return segments[0], segments[1], nil
}

type parsedIBMISSecurityGroupRuleDictionary struct {
	// After parsing, unused string fields are set to
	// "" and unused int64 fields will be set to -1.
	// This ("" for unused strings and -1 for unused int64s)
	// is expected by our riaas API client.
	secgrpID       string
	ruleID         string
	direction      string
	ipversion      string
	remote         string
	remoteAddress  string
	remoteCIDR     string
	remoteSecGrpID string
	protocol       string
	icmpType       int64
	icmpCode       int64
	portMin        int64
	portMax        int64
}

func inferRemoteSecurityGroup(s string) (address, cidr, id string, err error) {
	if isSecurityGroupAddress(s) {
		address = s
		return
	} else if isSecurityGroupCIDR(s) {
		cidr = s
		return
	} else {
		id = s
		return
	}
}

func parseIBMISClassicSecurityGroupRuleDictionary(d *schema.ResourceData, tag string, sess *vpcclassicv1.VpcClassicV1) (*parsedIBMISSecurityGroupRuleDictionary, *vpcclassicv1.SecurityGroupRulePrototype, *vpcclassicv1.SecurityGroupRulePatch, error) {
	parsed := &parsedIBMISSecurityGroupRuleDictionary{}
	sgTemplate := &vpcclassicv1.SecurityGroupRulePrototype{}
	sgTemplateUpdate := &vpcclassicv1.SecurityGroupRulePatch{}
	var err error
	parsed.icmpType = -1
	parsed.icmpCode = -1
	parsed.portMin = -1
	parsed.portMax = -1

	parsed.secgrpID, parsed.ruleID, err = parseISTerraformID(d.Id())
	if err != nil {
		parsed.secgrpID = d.Get(isSecurityGroupID).(string)
	}

	parsed.direction = d.Get(isSecurityGroupRuleDirection).(string)
	sgTemplate.Direction = &parsed.direction
	sgTemplateUpdate.Direction = &parsed.direction

	if version, ok := d.GetOk(isSecurityGroupRuleIPVersion); ok {
		parsed.ipversion = version.(string)
		sgTemplate.IPVersion = &parsed.ipversion
		sgTemplateUpdate.IPVersion = &parsed.ipversion
	} else {
		parsed.ipversion = "IPv4"
		sgTemplate.IPVersion = &parsed.ipversion
		sgTemplateUpdate.IPVersion = &parsed.ipversion
	}

	parsed.remote = ""
	if pr, ok := d.GetOk(isSecurityGroupRuleRemote); ok {
		parsed.remote = pr.(string)
	}
	parsed.remoteAddress = ""
	parsed.remoteCIDR = ""
	parsed.remoteSecGrpID = ""
	err = nil
	if parsed.remote != "" {
		parsed.remoteAddress, parsed.remoteCIDR, parsed.remoteSecGrpID, err = inferRemoteSecurityGroup(parsed.remote)
		remoteTemplate := &vpcclassicv1.SecurityGroupRuleRemotePrototype{}
		remoteTemplateUpdate := &vpcclassicv1.SecurityGroupRuleRemotePatch{}
		if parsed.remoteAddress != "" {
			remoteTemplate.Address = &parsed.remoteAddress
			remoteTemplateUpdate.Address = &parsed.remoteAddress
		} else if parsed.remoteCIDR != "" {
			remoteTemplate.CIDRBlock = &parsed.remoteCIDR
			remoteTemplateUpdate.CIDRBlock = &parsed.remoteCIDR
		} else if parsed.remoteSecGrpID != "" {
			remoteTemplate.ID = &parsed.remoteSecGrpID
			remoteTemplateUpdate.ID = &parsed.remoteSecGrpID
		}
		sgTemplate.Remote = remoteTemplate
		sgTemplateUpdate.Remote = remoteTemplateUpdate
	}
	if err != nil {
		return nil, nil, nil, err
	}
	parsed.protocol = "all"

	if icmpInterface, ok := d.GetOk("icmp"); ok {
		if icmpInterface.([]interface{})[0] != nil {
			haveType := false
			icmp := icmpInterface.([]interface{})[0].(map[string]interface{})
			if value, ok := icmp["type"]; ok {
				parsed.icmpType = int64(value.(int))
				haveType = true
			}
			if value, ok := icmp["code"]; ok {
				if !haveType {
					return nil, nil, nil, fmt.Errorf("icmp code requires icmp type")
				}
				parsed.icmpCode = int64(value.(int))
			}
		}
		parsed.protocol = "icmp"
		if icmpInterface.([]interface{})[0] == nil {
			parsed.icmpType = 0
			parsed.icmpCode = 0
		} else {
			sgTemplate.Type = &parsed.icmpType
			sgTemplate.Code = &parsed.icmpCode
		}
		sgTemplate.Protocol = &parsed.protocol
		sgTemplateUpdate.Type = &parsed.icmpType
		sgTemplateUpdate.Code = &parsed.icmpCode
	}
	for _, prot := range []string{"tcp", "udp"} {
		if tcpInterface, ok := d.GetOk(prot); ok {
			if tcpInterface.([]interface{})[0] != nil {
				haveMin := false
				haveMax := false
				ports := tcpInterface.([]interface{})[0].(map[string]interface{})
				if value, ok := ports["port_min"]; ok {
					parsed.portMin = int64(value.(int))
					haveMin = true
				}
				if value, ok := ports["port_max"]; ok {
					parsed.portMax = int64(value.(int))
					haveMax = true
				}

				// If only min or max is set, ensure that both min and max are set to the same value
				if haveMin && !haveMax {
					parsed.portMax = parsed.portMin
				}
				if haveMax && !haveMin {
					parsed.portMin = parsed.portMax
				}
			}
			parsed.protocol = prot
			sgTemplate.Protocol = &parsed.protocol
			if tcpInterface.([]interface{})[0] == nil {
				parsed.portMax = 65535
				parsed.portMin = 1
			}
			sgTemplate.PortMax = &parsed.portMax
			sgTemplate.PortMin = &parsed.portMin
			sgTemplateUpdate.PortMax = &parsed.portMax
			sgTemplateUpdate.PortMin = &parsed.portMin
		}
	}

	if parsed.protocol == "all" {
		sgTemplate.Protocol = &parsed.protocol
		sgTemplateUpdate.Protocol = &parsed.protocol
	}
	//	log.Printf("[DEBUG] parse tag=%s\n\t%v  \n\t%v  \n\t%v  \n\t%v  \n\t%v \n\t%v \n\t%v \n\t%v  \n\t%v  \n\t%v  \n\t%v  \n\t%v ",
	//		tag, parsed.secgrpID, parsed.ruleID, parsed.direction, parsed.ipversion, parsed.protocol, parsed.remoteAddress,
	//		parsed.remoteCIDR, parsed.remoteSecGrpID, parsed.icmpType, parsed.icmpCode, parsed.portMin, parsed.portMax)
	return parsed, sgTemplate, sgTemplateUpdate, nil
}

func parseIBMISSecurityGroupRuleDictionary(d *schema.ResourceData, tag string, sess *vpcv1.VpcV1) (*parsedIBMISSecurityGroupRuleDictionary, *vpcv1.SecurityGroupRulePrototype, *vpcv1.UpdateSecurityGroupRuleOptions, error) {
	parsed := &parsedIBMISSecurityGroupRuleDictionary{}
	sgTemplate := &vpcv1.SecurityGroupRulePrototype{}
	sgTemplateUpdate := &vpcv1.UpdateSecurityGroupRuleOptions{}
	var err error
	parsed.icmpType = -1
	parsed.icmpCode = -1
	parsed.portMin = -1
	parsed.portMax = -1

	parsed.secgrpID, parsed.ruleID, err = parseISTerraformID(d.Id())
	if err != nil {
		parsed.secgrpID = d.Get(isSecurityGroupID).(string)
	} else {
		sgTemplateUpdate.SecurityGroupID = &parsed.secgrpID
		sgTemplateUpdate.ID = &parsed.ruleID
	}

	securityGroupRulePatchModel := &vpcv1.SecurityGroupRulePatch{}

	parsed.direction = d.Get(isSecurityGroupRuleDirection).(string)
	sgTemplate.Direction = &parsed.direction
	securityGroupRulePatchModel.Direction = &parsed.direction

	if version, ok := d.GetOk(isSecurityGroupRuleIPVersion); ok {
		parsed.ipversion = version.(string)
		sgTemplate.IPVersion = &parsed.ipversion
		securityGroupRulePatchModel.IPVersion = &parsed.ipversion
	} else {
		parsed.ipversion = "IPv4"
		sgTemplate.IPVersion = &parsed.ipversion
		securityGroupRulePatchModel.IPVersion = &parsed.ipversion
	}

	parsed.remote = ""
	if pr, ok := d.GetOk(isSecurityGroupRuleRemote); ok {
		parsed.remote = pr.(string)
	}
	parsed.remoteAddress = ""
	parsed.remoteCIDR = ""
	parsed.remoteSecGrpID = ""
	err = nil
	if parsed.remote != "" {
		parsed.remoteAddress, parsed.remoteCIDR, parsed.remoteSecGrpID, err = inferRemoteSecurityGroup(parsed.remote)
		remoteTemplate := &vpcv1.SecurityGroupRuleRemotePrototype{}
		remoteTemplateUpdate := &vpcv1.SecurityGroupRuleRemotePatch{}
		if parsed.remoteAddress != "" {
			remoteTemplate.Address = &parsed.remoteAddress
			remoteTemplateUpdate.Address = &parsed.remoteAddress
		} else if parsed.remoteCIDR != "" {
			remoteTemplate.CIDRBlock = &parsed.remoteCIDR
			remoteTemplateUpdate.CIDRBlock = &parsed.remoteCIDR
		} else if parsed.remoteSecGrpID != "" {
			remoteTemplate.ID = &parsed.remoteSecGrpID
			remoteTemplateUpdate.ID = &parsed.remoteSecGrpID
		}
		sgTemplate.Remote = remoteTemplate
		securityGroupRulePatchModel.Remote = remoteTemplateUpdate
	}
	if err != nil {
		return nil, nil, nil, err
	}
	parsed.protocol = "all"

	if icmpInterface, ok := d.GetOk("icmp"); ok {
		if icmpInterface.([]interface{})[0] != nil {
			haveType := false
			icmp := icmpInterface.([]interface{})[0].(map[string]interface{})
			if value, ok := icmp["type"]; ok {
				parsed.icmpType = int64(value.(int))
				haveType = true
			}
			if value, ok := icmp["code"]; ok {
				if !haveType {
					return nil, nil, nil, fmt.Errorf("icmp code requires icmp type")
				}
				parsed.icmpCode = int64(value.(int))
			}
		}
		parsed.protocol = "icmp"
		if icmpInterface.([]interface{})[0] == nil {
			parsed.icmpType = 0
			parsed.icmpCode = 0
		} else {
			sgTemplate.Type = &parsed.icmpType
			sgTemplate.Code = &parsed.icmpCode
		}
		sgTemplate.Protocol = &parsed.protocol
		securityGroupRulePatchModel.Type = &parsed.icmpType
		securityGroupRulePatchModel.Code = &parsed.icmpCode
	}
	for _, prot := range []string{"tcp", "udp"} {
		if tcpInterface, ok := d.GetOk(prot); ok {
			if tcpInterface.([]interface{})[0] != nil {
				haveMin := false
				haveMax := false
				ports := tcpInterface.([]interface{})[0].(map[string]interface{})
				if value, ok := ports["port_min"]; ok {
					parsed.portMin = int64(value.(int))
					haveMin = true
				}
				if value, ok := ports["port_max"]; ok {
					parsed.portMax = int64(value.(int))
					haveMax = true
				}

				// If only min or max is set, ensure that both min and max are set to the same value
				if haveMin && !haveMax {
					parsed.portMax = parsed.portMin
				}
				if haveMax && !haveMin {
					parsed.portMin = parsed.portMax
				}
			}
			parsed.protocol = prot
			sgTemplate.Protocol = &parsed.protocol
			if tcpInterface.([]interface{})[0] == nil {
				parsed.portMax = 65535
				parsed.portMin = 1
			}
			sgTemplate.PortMax = &parsed.portMax
			sgTemplate.PortMin = &parsed.portMin
			securityGroupRulePatchModel.PortMax = &parsed.portMax
			securityGroupRulePatchModel.PortMin = &parsed.portMin
		}
	}
	if parsed.protocol == "all" {
		sgTemplate.Protocol = &parsed.protocol
	}
	securityGroupRulePatch, err := securityGroupRulePatchModel.AsPatch()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Error calling asPatch for SecurityGroupRulePatch: %s", err)
	}
	sgTemplateUpdate.SecurityGroupRulePatch = securityGroupRulePatch
	//	log.Printf("[DEBUG] parse tag=%s\n\t%v  \n\t%v  \n\t%v  \n\t%v  \n\t%v \n\t%v \n\t%v \n\t%v  \n\t%v  \n\t%v  \n\t%v  \n\t%v ",
	//		tag, parsed.secgrpID, parsed.ruleID, parsed.direction, parsed.ipversion, parsed.protocol, parsed.remoteAddress,
	//		parsed.remoteCIDR, parsed.remoteSecGrpID, parsed.icmpType, parsed.icmpCode, parsed.portMin, parsed.portMax)
	return parsed, sgTemplate, sgTemplateUpdate, nil
}

func makeTerraformRuleID(id1, id2 string) string {
	// Include both group and rule id to create a unique Terraform id.  As a bonus,
	// we can extract the group id as needed for API calls such as READ.
	return id1 + "." + id2
}
