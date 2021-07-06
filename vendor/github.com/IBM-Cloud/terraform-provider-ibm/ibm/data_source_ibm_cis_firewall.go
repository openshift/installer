// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"log"
	"time"

	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataIBMCISFirewallsRecord() *schema.Resource {
	return &schema.Resource{
		Read: dataIBMCISFirewallRecordRead,
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS object id",
				Required:    true,
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Description:      "Associated CIS domain",
				Required:         true,
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisFirewallType: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of firewall.Allowable values are access-rules,ua-rules,lockdowns",
			},
			cisFirewallLockdown: {
				Type:        schema.TypeList,
				Computed:    true,
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
							Computed:    true,
							Description: "Firewall rule paused or enabled",
						},
						cisFirewallLockdownDesc: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "description",
						},
						cisFirewallLockdownPriority: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Firewall priority",
						},
						cisFirewallLockdownURLs: {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "URL in which firewall rule is applied",
						},
						cisFirewallLockdownConfigurations: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									cisFirewallLockdownConfigurationsTarget: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Target type",
									},
									cisFirewallLockdownConfigurationsValue: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Target value",
									},
								},
							},
						},
					},
				},
			},
			cisFirewallAccessRule: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Access Rule Data",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						cisFirewallAccessRuleID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "firewall identifier",
						},
						cisFirewallAccessRuleNotes: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "description",
						},
						cisFirewallAccessRuleMode: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access rule mode",
						},
						cisFirewallAccessRuleConfiguration: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									cisFirewallUARuleConfigurationTarget: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Target type",
									},
									cisFirewallUARuleConfigurationValue: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Target value",
									},
								},
							},
						},
					},
				},
			},
			cisFirewallUARule: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "User Agent Rule Data",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						cisFirewallUARuleID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "firewall identifier",
						},
						cisFirewallUARulePaused: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Rule whether paused or not",
						},
						cisFirewallUARuleDesc: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "description",
						},
						cisFirewallUARuleMode: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "user agent rule mode",
						},
						cisFirewallUARuleConfiguration: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									cisFirewallUARuleConfigurationTarget: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Target type",
									},
									cisFirewallUARuleConfigurationValue: {
										Type:        schema.TypeString,
										Computed:    true,
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
func dataIBMCISFirewallRecordRead(d *schema.ResourceData, meta interface{}) error {
	crn := d.Get(cisID).(string)
	zoneID, _, _ := convertTftoCisTwoVar(d.Get(cisDomainID).(string))
	firewallType := d.Get(cisFirewallType).(string)

	if firewallType == cisFirewallTypeLockdowns {
		cisClient, err := meta.(ClientSession).CisLockdownClientSession()
		if err != nil {
			return err
		}
		cisClient.Crn = core.StringPtr(crn)
		cisClient.ZoneIdentifier = core.StringPtr(zoneID)
		opt := cisClient.NewListAllZoneLockownRulesOptions()
		result, response, err := cisClient.ListAllZoneLockownRules(opt)
		if err != nil {
			log.Printf("List all zone lockdown rules failed: %v", response)
			return err
		}
		lockdownList := make([]map[string]interface{}, 0)
		for _, instance := range result.Result {
			configurationList := []interface{}{}
			for _, c := range instance.Configurations {
				configuration := make(map[string]interface{}, 0)
				configuration[cisFirewallLockdownConfigurationsTarget] = c.Target
				configuration[cisFirewallLockdownConfigurationsValue] = c.Value
				configurationList = append(configurationList, configuration)
			}
			lockdown := make(map[string]interface{})
			lockdown[cisFirewallLockdownID] = *instance.ID
			lockdown[cisFirewallLockdownPaused] = *instance.Paused
			if instance.Priority != nil {
				lockdown[cisFirewallLockdownPriority] = *instance.Priority
			}
			lockdown[cisFirewallLockdownURLs] = instance.Urls
			lockdown[cisFirewallLockdownConfigurations] = configurationList
			if instance.Description != nil {
				lockdown[cisFirewallLockdownDesc] = *instance.Description
			}
			lockdownList = append(lockdownList, lockdown)
		}
		d.Set(cisFirewallLockdown, lockdownList)
	} else if firewallType == cisFirewallTypeAccessRules {
		cisClient, err := meta.(ClientSession).CisAccessRuleClientSession()
		if err != nil {
			return err
		}
		cisClient.Crn = core.StringPtr(crn)
		cisClient.ZoneIdentifier = core.StringPtr(zoneID)
		opt := cisClient.NewListAllZoneAccessRulesOptions()
		result, response, err := cisClient.ListAllZoneAccessRules(opt)
		if err != nil {
			log.Printf("List all zone access rules failed: %v", response)
			return err
		}
		accessRuleList := make([]interface{}, 0)
		for _, instance := range result.Result {
			configurations := []interface{}{}
			configuration := map[string]interface{}{}
			configuration[cisFirewallAccessRuleConfigurationTarget] = *instance.Configuration.Target
			configuration[cisFirewallAccessRuleConfigurationValue] = *instance.Configuration.Value
			configurations = append(configurations, configuration)
			accessRule := make(map[string]interface{}, 0)
			accessRule[cisFirewallAccessRuleID] = *instance.ID
			accessRule[cisFirewallAccessRuleMode] = *instance.Mode
			accessRule[cisFirewallAccessRuleNotes] = *instance.Notes
			accessRule[cisFirewallAccessRuleConfiguration] = configurations
			accessRuleList = append(accessRuleList, accessRule)
		}
		d.Set(cisFirewallAccessRule, accessRuleList)
	} else if firewallType == cisFirewallTypeUARules {
		cisClient, err := meta.(ClientSession).CisUARuleClientSession()
		if err != nil {
			return err
		}
		cisClient.Crn = core.StringPtr(crn)
		cisClient.ZoneIdentifier = core.StringPtr(zoneID)
		opt := cisClient.NewListAllZoneUserAgentRulesOptions()
		result, response, err := cisClient.ListAllZoneUserAgentRules(opt)
		if err != nil {
			log.Printf("List all zone ua rules failed: %v", response)
			return err
		}
		uaRuleList := make([]interface{}, 0)
		for _, instance := range result.Result {
			configurations := []interface{}{}
			configuration := map[string]interface{}{}
			configuration[cisFirewallUARuleConfigurationTarget] = *instance.Configuration.Target
			configuration[cisFirewallUARuleConfigurationValue] = *instance.Configuration.Value
			configurations = append(configurations, configuration)
			uaRule := make(map[string]interface{}, 0)
			uaRule[cisFirewallUARuleID] = *instance.ID
			uaRule[cisFirewallUARuleMode] = *instance.Mode
			uaRule[cisFirewallUARulePaused] = *instance.Paused
			if instance.Description != nil {
				uaRule[cisFirewallUARuleDesc] = *instance.Description
			}
			uaRule[cisFirewallUARuleConfiguration] = configurations
			uaRuleList = append(uaRuleList, uaRule)
		}
		d.Set(cisFirewallUARule, uaRuleList)
	}

	d.SetId(dataIBMCISFirewallRecordsID(d))
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisFirewallType, firewallType)

	return nil
}

func dataIBMCISFirewallRecordsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
