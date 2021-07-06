// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"

	"github.com/IBM/go-sdk-core/v4/core"
	cislockdownv1 "github.com/IBM/networking-go-sdk/zonelockdownv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	ibmCISFirewall                                 = "ibm_cis_firewall"
	cisFirewallType                                = "firewall_type"
	cisFirewallTypeLockdowns                       = "lockdowns"
	cisFirewallTypeAccessRules                     = "access_rules"
	cisFirewallTypeUARules                         = "ua_rules"
	cisFirewallLockdown                            = "lockdown"
	cisFirewallLockdownID                          = "lockdown_id"
	cisFirewallLockdownName                        = "name"
	cisFirewallLockdownPaused                      = "paused"
	cisFirewallLockdownDesc                        = "description"
	cisFirewallLockdownPriority                    = "priority"
	cisFirewallLockdownURLs                        = "urls"
	cisFirewallLockdownConfigurations              = "configurations"
	cisFirewallLockdownConfigurationsTarget        = "target"
	cisFirewallLockdownConfigurationsTargetIP      = "ip"
	cisFirewallLockdownConfigurationsTargetIPRange = "ip_range"
	cisFirewallLockdownConfigurationsValue         = "value"
	cisFirewallAccessRule                          = "access_rule"
	cisFirewallAccessRuleID                        = "access_rule_id"
	cisFirewallAccessRuleMode                      = "mode"
	cisFirewallAccessRuleModeBlock                 = "block"
	cisFirewallAccessRuleModeChallenge             = "challenge"
	cisFirewallAccessRuleModeWhitelist             = "whitelist"
	cisFirewallAccessRuleModeJSChallenge           = "js_challenge"
	cisFirewallAccessRuleNotes                     = "notes"
	cisFirewallAccessRuleConfiguration             = "configuration"
	cisFirewallAccessRuleConfigurationTarget       = "target"
	cisFirewallAccessRuleConfigurationValue        = "value"
	cisFirewallUARule                              = "ua_rule"
	cisFirewallUARuleID                            = "ua_rule_id"
	cisFirewallUARulePaused                        = "paused"
	cisFirewallUARuleDesc                          = "description"
	cisFirewallUARuleMode                          = "mode"
	cisFirewallUARuleModeBlock                     = "block"
	cisFirewallUARuleModeChallenge                 = "challenge"
	cisFirewallUARuleModeJSChallenge               = "js_challenge"
	cisFirewallUARuleConfiguration                 = "configuration"
	cisFirewallUARuleConfigurationTarget           = "target"
	cisFirewallUARuleConfigurationTargetUA         = "ua"
	cisFirewallUARuleConfigurationValue            = "value"
)

func resourceIBMCISFirewallRecord() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMCISFirewallRecordCreate,
		Read:     resourceIBMCISFirewallRecordRead,
		Update:   resourceIBMCISFirewallRecordUpdate,
		Delete:   resourceIBMCISFirewallRecordDelete,
		Exists:   resourceIBMCISFirewallRecordExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS object id",
				Required:    true,
				ForceNew:    true,
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Description:      "Associated CIS domain",
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisFirewallType: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Type of firewall.Allowable values are access-rules,ua-rules,lockdowns",
				ValidateFunc: InvokeValidator(ibmCISFirewall, cisFirewallType),
			},

			cisFirewallLockdown: {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				ExactlyOneOf: []string{
					cisFirewallLockdown,
					cisFirewallAccessRule,
					cisFirewallUARule},
				Description: "Lockdown Data",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						cisFirewallLockdownID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "firewall identifier",
						},
						cisFirewallLockdownPaused: {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Firewall rule paused or enabled",
						},
						cisFirewallLockdownDesc: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "description",
						},
						cisFirewallLockdownPriority: {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Firewall priority",
						},
						cisFirewallLockdownURLs: {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "URL in which firewall rule is applied",
						},
						cisFirewallLockdownConfigurations: {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									cisFirewallLockdownConfigurationsTarget: {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Target type",
										ValidateFunc: InvokeValidator(
											ibmCISFirewall,
											cisFirewallLockdownConfigurationsTarget),
									},
									cisFirewallLockdownConfigurationsValue: {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Target value",
									},
								},
							},
						},
					},
				},
			},
			cisFirewallAccessRule: {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				ExactlyOneOf: []string{
					cisFirewallLockdown,
					cisFirewallAccessRule,
					cisFirewallUARule},
				Description: "Access Rule Data",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						cisFirewallAccessRuleID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "access rule firewall identifier",
						},
						cisFirewallAccessRuleNotes: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "description",
						},
						cisFirewallAccessRuleMode: {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Access rule mode",
							ValidateFunc: InvokeValidator(ibmCISFirewall, cisFirewallAccessRuleMode),
						},
						cisFirewallAccessRuleConfiguration: {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									cisFirewallUARuleConfigurationTarget: {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "Target type",
										ValidateFunc: InvokeValidator(ibmCISFirewall,
											cisFirewallAccessRuleConfigurationTarget),
									},
									cisFirewallUARuleConfigurationValue: {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "Target value",
									},
								},
							},
						},
					},
				},
			},
			cisFirewallUARule: {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				ExactlyOneOf: []string{
					cisFirewallLockdown,
					cisFirewallAccessRule,
					cisFirewallUARule},
				Description: "User Agent Rule Data",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						cisFirewallUARuleID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User Agent firewall identifier",
						},
						cisFirewallUARulePaused: {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Rule whether paused or not",
						},
						cisFirewallUARuleDesc: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "description",
						},
						cisFirewallUARuleMode: {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "user agent rule mode",
							ValidateFunc: InvokeValidator(ibmCISFirewall, cisFirewallUARuleMode),
						},
						cisFirewallUARuleConfiguration: {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									cisFirewallUARuleConfigurationTarget: {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Target type",
										ValidateFunc: InvokeValidator(ibmCISFirewall,
											cisFirewallUARuleConfigurationTarget),
									},
									cisFirewallUARuleConfigurationValue: {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Target value",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceIBMCISFirewallValidator() *ResourceValidator {
	firewallTypes := "access_rules, ua_rules, lockdowns"
	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisFirewallType,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              firewallTypes})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisFirewallLockdownConfigurationsTarget,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              "ip, ip_range"})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisFirewallAccessRuleConfigurationTarget,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              "ip, ip_range, asn, country"})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisFirewallUARuleConfigurationTarget,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              "ua"})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisFirewallAccessRuleMode,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              "block, challenge, whitelist, js_challenge"})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisFirewallUARuleMode,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              "block, challenge, js_challenge"})
	cisFirewallValidator := ResourceValidator{ResourceName: ibmCISHealthCheck, Schema: validateSchema}
	return &cisFirewallValidator
}

func resourceIBMCISFirewallRecordCreate(d *schema.ResourceData, meta interface{}) error {
	crn := d.Get(cisID).(string)
	zoneID, _, _ := convertTftoCisTwoVar(d.Get(cisDomainID).(string))
	firewallType := d.Get(cisFirewallType).(string)

	if firewallType == cisFirewallTypeLockdowns {
		// Firewall Type : Lockdowns

		cisClient, err := meta.(ClientSession).CisLockdownClientSession()
		if err != nil {
			return err
		}
		lockdown := d.Get(cisFirewallLockdown).([]interface{})[0].(map[string]interface{})

		opt := cisClient.NewCreateZoneLockdownRuleOptions()
		// not able to check bool variable availability
		v, _ := lockdown[cisFirewallLockdownPaused]
		opt.SetPaused(v.(bool))
		if v, ok := lockdown[cisFirewallLockdownDesc]; ok && v.(string) != "" {
			opt.SetDescription(v.(string))
		}
		if v, ok := lockdown[cisFirewallLockdownPriority]; ok && v.(int) > 0 {
			opt.SetPriority(int64(v.(int)))
		}
		urls := expandStringList(lockdown[cisFirewallLockdownURLs].([]interface{}))
		configurations, err := expandLockdownsTypeConfiguration(
			lockdown[cisFirewallLockdownConfigurations].([]interface{}))
		if err != nil {
			return err
		}
		opt.SetUrls(urls)
		opt.SetConfigurations(configurations)

		cisClient.Crn = core.StringPtr(crn)
		cisClient.ZoneIdentifier = core.StringPtr(zoneID)
		result, response, err := cisClient.CreateZoneLockdownRule(opt)
		if err != nil {
			log.Printf("Create zone firewall lockdown failed: %v", response)
			return err
		}
		d.SetId(convertCisToTfFourVar(firewallType, *result.Result.ID, zoneID, crn))

	} else if firewallType == cisFirewallTypeAccessRules {

		cisClient, err := meta.(ClientSession).CisAccessRuleClientSession()
		if err != nil {
			return err
		}
		cisClient.Crn = core.StringPtr(crn)
		cisClient.ZoneIdentifier = core.StringPtr(zoneID)
		accessRule := d.Get(cisFirewallAccessRule).([]interface{})[0].(map[string]interface{})

		mode := accessRule[cisFirewallAccessRuleMode].(string)

		configList := accessRule[cisFirewallAccessRuleConfiguration].([]interface{})

		config := configList[0].(map[string]interface{})
		target := config[cisFirewallAccessRuleConfigurationTarget].(string)
		value := config[cisFirewallAccessRuleConfigurationValue].(string)

		configOpt, err := cisClient.NewZoneAccessRuleInputConfiguration(target, value)
		if err != nil {
			log.Printf("Error in firewall type %s input: %s", firewallType, err)
			return err
		}

		opt := cisClient.NewCreateZoneAccessRuleOptions()
		opt.SetMode(mode)
		opt.SetConfiguration(configOpt)
		if v, ok := accessRule[cisFirewallAccessRuleNotes]; ok && v.(string) != "" {
			opt.SetNotes(v.(string))
		}

		result, response, err := cisClient.CreateZoneAccessRule(opt)
		if err != nil {
			log.Printf("Create zone firewall access rule failed: %v", response)
			return err
		}
		d.SetId(convertCisToTfFourVar(firewallType, *result.Result.ID, zoneID, crn))

	} else if firewallType == cisFirewallTypeUARules {
		cisClient, err := meta.(ClientSession).CisUARuleClientSession()
		if err != nil {
			return err
		}
		cisClient.Crn = core.StringPtr(crn)
		cisClient.ZoneIdentifier = core.StringPtr(zoneID)
		uaRule := d.Get(cisFirewallUARule).([]interface{})[0].(map[string]interface{})

		mode := uaRule[cisFirewallUARuleMode].(string)
		configList := uaRule[cisFirewallUARuleConfiguration].([]interface{})
		if len(configList) > 1 {
			return fmt.Errorf("Only one configuration is allowed for %s type", firewallType)
		}
		config := configList[0].(map[string]interface{})
		target := config[cisFirewallLockdownConfigurationsTarget].(string)
		value := config[cisFirewallLockdownConfigurationsValue].(string)

		configOpt, err := cisClient.NewUseragentRuleInputConfiguration(target, value)
		if err != nil {
			log.Printf("Error in firewall type %s input: %s", firewallType, err)
			return err
		}

		opt := cisClient.NewCreateZoneUserAgentRuleOptions()
		opt.SetMode(mode)
		opt.SetConfiguration(configOpt)

		if v, ok := uaRule[cisFirewallUARuleDesc]; ok && v.(string) != "" {
			opt.SetDescription(v.(string))
		}
		// not able to check bool attribute availablity
		v, _ := uaRule[cisFirewallUARulePaused]
		opt.SetPaused(v.(bool))

		result, response, err := cisClient.CreateZoneUserAgentRule(opt)
		if err != nil {
			log.Printf("Create zone user agent rule failed: %v", response)
			return err
		}
		d.SetId(convertCisToTfFourVar(firewallType, *result.Result.ID, zoneID, crn))
	}

	return resourceIBMCISFirewallRecordRead(d, meta)
}

func resourceIBMCISFirewallRecordRead(d *schema.ResourceData, meta interface{}) error {
	firewallType, lockdownID, zoneID, crn, _ := convertTfToCisFourVar(d.Id())

	if firewallType == cisFirewallTypeLockdowns {
		// Firewall Type : Lockdowns
		cisClient, err := meta.(ClientSession).CisLockdownClientSession()
		if err != nil {
			return err
		}

		cisClient.Crn = core.StringPtr(crn)
		cisClient.ZoneIdentifier = core.StringPtr(zoneID)

		opt := cisClient.NewGetLockdownOptions(lockdownID)

		result, response, err := cisClient.GetLockdown(opt)
		if err != nil {
			log.Printf("Get zone firewall lockdown failed: %v", response)
			return err
		}
		lockdownList := []interface{}{}
		lockdown := map[string]interface{}{}
		lockdown[cisFirewallLockdownID] = *result.Result.ID
		lockdown[cisFirewallLockdownPaused] = *result.Result.Paused
		lockdown[cisFirewallLockdownURLs] = flattenStringList(result.Result.Urls)
		lockdown[cisFirewallLockdownConfigurations] =
			flattenLockdownsTypeConfiguration(result.Result.Configurations)
		if result.Result.Description != nil {
			lockdown[cisFirewallLockdownDesc] = *result.Result.Description
		}
		if result.Result.Priority != nil {
			lockdown[cisFirewallLockdownPriority] = *result.Result.Priority
		}
		lockdownList = append(lockdownList, lockdown)
		d.Set(cisFirewallLockdown, lockdownList)

	} else if firewallType == cisFirewallTypeAccessRules {

		// Firewall Type : Zone Access firewall rules
		cisClient, err := meta.(ClientSession).CisAccessRuleClientSession()
		if err != nil {
			return err
		}
		cisClient.Crn = core.StringPtr(crn)
		cisClient.ZoneIdentifier = core.StringPtr(zoneID)

		opt := cisClient.NewGetZoneAccessRuleOptions(lockdownID)

		result, response, err := cisClient.GetZoneAccessRule(opt)
		if err != nil {
			log.Printf("Get zone firewall lockdown failed: %v", response)
			return err
		}

		config := map[string]interface{}{}
		configList := []interface{}{}
		config[cisFirewallUARuleConfigurationTarget] = *result.Result.Configuration.Target
		config[cisFirewallUARuleConfigurationValue] = *result.Result.Configuration.Value
		configList = append(configList, config)

		accessRuleList := []interface{}{}
		accessRule := map[string]interface{}{}
		accessRule[cisFirewallAccessRuleID] = *result.Result.ID
		accessRule[cisFirewallAccessRuleNotes] = *result.Result.Notes
		accessRule[cisFirewallAccessRuleMode] = *result.Result.Mode
		accessRule[cisFirewallAccessRuleConfiguration] = configList
		accessRuleList = append(accessRuleList, accessRule)
		d.Set(cisFirewallAccessRule, accessRuleList)

	} else if firewallType == cisFirewallTypeUARules {
		// Firewall Type: User Agent access rules
		cisClient, err := meta.(ClientSession).CisUARuleClientSession()
		if err != nil {
			return err
		}
		cisClient.Crn = core.StringPtr(crn)
		cisClient.ZoneIdentifier = core.StringPtr(zoneID)

		opt := cisClient.NewGetUserAgentRuleOptions(lockdownID)
		result, response, err := cisClient.GetUserAgentRule(opt)
		if err != nil {
			log.Printf("Get zone user agent rule failed: %v", response)
			return err
		}

		config := map[string]interface{}{}
		configList := []interface{}{}
		config[cisFirewallUARuleConfigurationTarget] = *result.Result.Configuration.Target
		config[cisFirewallUARuleConfigurationValue] = *result.Result.Configuration.Value
		configList = append(configList, config)

		uaRuleList := []interface{}{}
		uaRule := map[string]interface{}{}
		uaRule[cisFirewallUARuleID] = *result.Result.ID
		uaRule[cisFirewallUARulePaused] = *result.Result.Paused
		uaRule[cisFirewallUARuleMode] = *result.Result.Mode
		uaRule[cisFirewallUARuleConfiguration] = configList
		if result.Result.Description != nil {
			uaRule[cisFirewallUARuleDesc] = *result.Result.Description
		}
		uaRuleList = append(uaRuleList, uaRule)
		d.Set(cisFirewallUARule, uaRuleList)
	}
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisFirewallType, firewallType)
	return nil
}

func resourceIBMCISFirewallRecordUpdate(d *schema.ResourceData, meta interface{}) error {

	firewallType, lockdownID, zoneID, crn, _ := convertTfToCisFourVar(d.Id())

	if d.HasChange(cisFirewallLockdown) ||
		d.HasChange(cisFirewallAccessRule) ||
		d.HasChange(cisFirewallUARule) {

		if firewallType == cisFirewallTypeLockdowns {
			// Firewall Type : Lockdowns
			lockdown := d.Get(cisFirewallLockdown).([]interface{})[0].(map[string]interface{})

			cisClient, err := meta.(ClientSession).CisLockdownClientSession()
			if err != nil {
				return err
			}

			opt := cisClient.NewUpdateLockdownRuleOptions(lockdownID)
			// not able to check bool variable availability
			v, _ := lockdown[cisFirewallLockdownPaused]
			opt.SetPaused(v.(bool))
			if v, ok := lockdown[cisFirewallLockdownDesc]; ok && v.(string) != "" {
				opt.SetDescription(v.(string))
			}
			if v, ok := lockdown[cisFirewallLockdownPriority]; ok && v.(int) > 0 {
				opt.SetPriority(int64(v.(int)))
			}
			urls := expandStringList(lockdown[cisFirewallLockdownURLs].([]interface{}))
			configurations, err := expandLockdownsTypeConfiguration(lockdown[cisFirewallLockdownConfigurations].([]interface{}))
			if err != nil {
				return err
			}
			opt.SetUrls(urls)
			opt.SetConfigurations(configurations)

			cisClient.Crn = core.StringPtr(crn)
			cisClient.ZoneIdentifier = core.StringPtr(zoneID)
			_, response, err := cisClient.UpdateLockdownRule(opt)
			if err != nil {
				log.Printf("Update zone firewall lockdown failed: %v", response)
				return err
			}

		} else if firewallType == cisFirewallTypeAccessRules {

			accessRule := d.Get(cisFirewallAccessRule).([]interface{})[0].(map[string]interface{})

			// Firewall Type : Zone Access firewall rules
			cisClient, err := meta.(ClientSession).CisAccessRuleClientSession()
			if err != nil {
				return err
			}
			cisClient.Crn = core.StringPtr(crn)
			cisClient.ZoneIdentifier = core.StringPtr(zoneID)

			mode := accessRule[cisFirewallAccessRuleMode].(string)
			opt := cisClient.NewUpdateZoneAccessRuleOptions(lockdownID)
			if v, ok := accessRule[cisFirewallAccessRuleNotes]; ok && v.(string) != "" {
				opt.SetNotes(v.(string))
			}
			opt.SetMode(mode)

			_, response, err := cisClient.UpdateZoneAccessRule(opt)
			if err != nil {
				log.Printf("Update zone firewall access rule failed: %v", response)
				return err
			}

		} else if firewallType == cisFirewallTypeUARules {
			// Firewall Type: User Agent access rules
			uaRule := d.Get(cisFirewallUARule).([]interface{})[0].(map[string]interface{})
			cisClient, err := meta.(ClientSession).CisUARuleClientSession()
			if err != nil {
				return err
			}
			cisClient.Crn = core.StringPtr(crn)
			cisClient.ZoneIdentifier = core.StringPtr(zoneID)

			mode := uaRule[cisFirewallUARuleMode].(string)
			config := uaRule[cisFirewallUARuleConfiguration].([]interface{})[0].(map[string]interface{})
			target := config[cisFirewallUARuleConfigurationTarget].(string)
			value := config[cisFirewallUARuleConfigurationValue].(string)

			configOpt, err := cisClient.NewUseragentRuleInputConfiguration(target, value)
			if err != nil {
				log.Printf("Error in firewall type %s input: %s", firewallType, err)
				return err
			}

			opt := cisClient.NewUpdateUserAgentRuleOptions(lockdownID)
			opt.SetMode(mode)
			opt.SetConfiguration(configOpt)

			if v, ok := uaRule[cisFirewallUARuleDesc]; ok && v.(string) != "" {
				opt.SetDescription(v.(string))
			}
			// not able to check bool attribute availablity
			v, _ := uaRule[cisFirewallUARulePaused]
			opt.SetPaused(v.(bool))

			_, response, err := cisClient.UpdateUserAgentRule(opt)
			if err != nil {
				log.Printf("Update zone user agent rule failed: %v", response)
				return err
			}
		}

	}
	return resourceIBMCISFirewallRecordRead(d, meta)
}

func resourceIBMCISFirewallRecordDelete(d *schema.ResourceData, meta interface{}) error {
	firewallType, lockdownID, zoneID, crn, _ := convertTfToCisFourVar(d.Id())

	if firewallType == cisFirewallTypeLockdowns {
		// Firewall Type : Lockdowns
		cisClient, err := meta.(ClientSession).CisLockdownClientSession()
		if err != nil {
			return err
		}

		cisClient.Crn = core.StringPtr(crn)
		cisClient.ZoneIdentifier = core.StringPtr(zoneID)

		opt := cisClient.NewDeleteZoneLockdownRuleOptions(lockdownID)

		_, response, err := cisClient.DeleteZoneLockdownRule(opt)
		if err != nil {
			log.Printf("Delete zone firewall lockdown failed: %v", response)
			return err
		}

	} else if firewallType == cisFirewallTypeAccessRules {

		// Firewall Type : Zone Access firewall rules
		cisClient, err := meta.(ClientSession).CisAccessRuleClientSession()
		if err != nil {
			return err
		}
		cisClient.Crn = core.StringPtr(crn)
		cisClient.ZoneIdentifier = core.StringPtr(zoneID)

		opt := cisClient.NewDeleteZoneAccessRuleOptions(lockdownID)

		_, response, err := cisClient.DeleteZoneAccessRule(opt)
		if err != nil {
			log.Printf("Delete zone firewall access rule failed: %v", response)
			return err
		}

	} else if firewallType == cisFirewallTypeUARules {
		// Firewall Type: User Agent access rules
		cisClient, err := meta.(ClientSession).CisUARuleClientSession()
		if err != nil {
			return err
		}
		cisClient.Crn = core.StringPtr(crn)
		cisClient.ZoneIdentifier = core.StringPtr(zoneID)

		opt := cisClient.NewDeleteZoneUserAgentRuleOptions(lockdownID)
		_, response, err := cisClient.DeleteZoneUserAgentRule(opt)
		if err != nil {
			log.Printf("Delete zone user agent rule failed: %v", response)
			return err
		}
	}

	return nil
}
func resourceIBMCISFirewallRecordExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	firewallType, lockdownID, zoneID, crn, _ := convertTfToCisFourVar(d.Id())

	if firewallType == cisFirewallTypeLockdowns {
		// Firewall Type : Lockdowns
		cisClient, err := meta.(ClientSession).CisLockdownClientSession()
		if err != nil {
			return false, err
		}

		cisClient.Crn = core.StringPtr(crn)
		cisClient.ZoneIdentifier = core.StringPtr(zoneID)

		opt := cisClient.NewGetLockdownOptions(lockdownID)

		_, response, err := cisClient.GetLockdown(opt)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				log.Printf("Zone Firewall Lockdown is not found")
				return false, nil
			}
			log.Printf("Get zone firewall lockdown failed: %v", response)
			return false, err
		}

	} else if firewallType == cisFirewallTypeAccessRules {

		// Firewall Type : Zone Access firewall rules
		cisClient, err := meta.(ClientSession).CisAccessRuleClientSession()
		if err != nil {
			return false, err
		}
		cisClient.Crn = core.StringPtr(crn)
		cisClient.ZoneIdentifier = core.StringPtr(zoneID)

		opt := cisClient.NewGetZoneAccessRuleOptions(lockdownID)

		_, response, err := cisClient.GetZoneAccessRule(opt)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				log.Printf("Zone Firewall Access Rule is not found")
				return false, nil
			}
			log.Printf("Get zone firewall lockdown failed: %v", response)
			return false, err
		}

	} else if firewallType == cisFirewallTypeUARules {
		// Firewall Type: User Agent access rules
		cisClient, err := meta.(ClientSession).CisUARuleClientSession()
		if err != nil {
			return false, err
		}
		cisClient.Crn = core.StringPtr(crn)
		cisClient.ZoneIdentifier = core.StringPtr(zoneID)

		opt := cisClient.NewGetUserAgentRuleOptions(lockdownID)
		_, response, err := cisClient.GetUserAgentRule(opt)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				log.Printf("Zone Firewall User Agent Rule does not found")
				return false, nil
			}
			log.Printf("Get zone user agent rule failed: %v", response)
			return false, err
		}

	}

	return true, nil
}

func expandLockdownsTypeConfiguration(lockdownConfigs []interface{}) ([]cislockdownv1.LockdownInputConfigurationsItem, error) {
	var configListOutput = make([]cislockdownv1.LockdownInputConfigurationsItem, 0)

	for _, lockdownConfig := range lockdownConfigs {
		configMap, _ := lockdownConfig.(map[string]interface{})
		target := configMap[cisFirewallLockdownConfigurationsTarget].(string)
		value := configMap[cisFirewallLockdownConfigurationsValue].(string)
		configOutput := cislockdownv1.LockdownInputConfigurationsItem{
			Target: core.StringPtr(target),
			Value:  core.StringPtr(value),
		}
		configListOutput = append(configListOutput, configOutput)
	}
	return configListOutput, nil
}

func flattenLockdownsTypeConfiguration(lockdownConfigs []cislockdownv1.LockdownObjectConfigurationsItem) interface{} {
	configListOutput := []interface{}{}

	for _, lockdownConfig := range lockdownConfigs {
		configOutput := map[string]string{}
		configOutput[cisFirewallLockdownConfigurationsTarget] = *lockdownConfig.Target
		configOutput[cisFirewallLockdownConfigurationsValue] = *lockdownConfig.Value
		configListOutput = append(configListOutput, configOutput)
	}
	return configListOutput
}
