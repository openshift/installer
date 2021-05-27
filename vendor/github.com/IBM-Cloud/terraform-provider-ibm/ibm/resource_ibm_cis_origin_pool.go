// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"log"

	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/IBM/networking-go-sdk/globalloadbalancerpoolsv0"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	cisGLBPoolID                   = "pool_id"
	cisGLBPoolName                 = "name"
	cisGLBPoolRegions              = "check_regions"
	cisGLBPoolDesc                 = "description"
	cisGLBPoolEnabled              = "enabled"
	cisGLBPoolMinimumOrigins       = "minimum_origins"
	cisGLBPoolMonitor              = "monitor"
	cisGLBPoolNotificationEMail    = "notification_email"
	cisGLBPoolOrigins              = "origins"
	cisGLBPoolHealth               = "health"
	cisGLBPoolHealthy              = "healthy"
	cisGLBPoolCreatedOn            = "created_on"
	cisGLBPoolModifiedOn           = "modified_on"
	cisGLBPoolOriginsName          = "name"
	cisGLBPoolOriginsAddress       = "address"
	cisGLBPoolOriginsEnabled       = "enabled"
	cisGLBPoolOriginsHealthy       = "healthy"
	cisGLBPoolOriginsWeight        = "weight"
	cisGLBPoolOriginsDisabledAt    = "disabled_at"
	cisGLBPoolOriginsFailureReason = "failure_reason"
)

func resourceIBMCISPool() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
			},
			cisGLBPoolID: {
				Type:     schema.TypeString,
				Computed: true,
			},
			cisGLBPoolName: {
				Type:        schema.TypeString,
				Description: "name",
				Required:    true,
			},
			cisGLBPoolRegions: {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of regions",
			},
			cisGLBPoolDesc: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the CIS Origin Pool",
			},
			cisGLBPoolEnabled: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Boolean value set to true if cis origin pool needs to be enabled",
			},
			cisGLBPoolMinimumOrigins: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Minimum number of Origins",
			},
			cisGLBPoolMonitor: {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Monitor value",
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisGLBPoolNotificationEMail: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Email address configured to recieve the notifications",
			},
			cisGLBPoolOrigins: {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Origins info",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						cisGLBPoolOriginsName: {
							Type:     schema.TypeString,
							Required: true,
						},
						cisGLBPoolOriginsAddress: {
							Type:     schema.TypeString,
							Required: true,
						},
						cisGLBPoolOriginsEnabled: {
							Type:     schema.TypeBool,
							Required: true,
						},
						cisGLBPoolOriginsWeight: {
							Type:     schema.TypeFloat,
							Default:  1,
							Optional: true,
						},
						cisGLBPoolOriginsHealthy: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						cisGLBPoolOriginsDisabledAt: {
							Type:     schema.TypeString,
							Computed: true,
						},
						cisGLBPoolOriginsFailureReason: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			cisGLBPoolHealth: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Health info",
			},
			cisGLBPoolHealthy: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Health status",
			},
			cisGLBPoolCreatedOn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation date info",
			},
			cisGLBPoolModifiedOn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Modified date info",
			},
		},

		Create:   resourceCISPoolCreate,
		Read:     resourceCISPoolRead,
		Update:   resourceCISPoolUpdate,
		Delete:   resourceCISPoolDelete,
		Exists:   resourceCISPoolExists,
		Importer: &schema.ResourceImporter{},
	}
}

func resourceCISPoolCreate(d *schema.ResourceData, meta interface{}) error {
	var regions []string
	cisClient, err := meta.(ClientSession).CisGLBPoolClientSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	cisClient.Crn = core.StringPtr(crn)
	name := d.Get(cisGLBPoolName).(string)
	origins := d.Get(cisGLBPoolOrigins).(*schema.Set).List()
	checkRegions := d.Get(cisGLBPoolRegions).(*schema.Set).List()

	for _, region := range checkRegions {
		regions = append(regions, region.(string))
	}

	glbOrigins := []globalloadbalancerpoolsv0.LoadBalancerPoolReqOriginsItem{}

	for _, origin := range origins {
		orig := origin.(map[string]interface{})
		glbOrigin := globalloadbalancerpoolsv0.LoadBalancerPoolReqOriginsItem{
			Name:    core.StringPtr(orig[cisGLBPoolOriginsName].(string)),
			Address: core.StringPtr(orig[cisGLBPoolOriginsAddress].(string)),
			Enabled: core.BoolPtr(orig[cisGLBPoolOriginsEnabled].(bool)),
			Weight:  core.Float64Ptr(orig[cisGLBPoolOriginsWeight].(float64)),
		}
		glbOrigins = append(glbOrigins, glbOrigin)
	}

	opt := cisClient.NewCreateLoadBalancerPoolOptions()
	opt.SetName(name)
	opt.SetCheckRegions(regions)
	opt.SetOrigins(glbOrigins)
	opt.SetEnabled(d.Get(cisGLBPoolEnabled).(bool))

	if notifEmail, ok := d.GetOk(cisGLBPoolNotificationEMail); ok {
		opt.SetNotificationEmail(notifEmail.(string))
	}
	if monitor, ok := d.GetOk(cisGLBPoolMonitor); ok {
		monitorID, _, _ := convertTftoCisTwoVar(monitor.(string))
		opt.SetMonitor(monitorID)
	}
	if minOrigins, ok := d.GetOk(cisGLBPoolMinimumOrigins); ok {
		opt.SetMinimumOrigins(int64(minOrigins.(int)))
	}
	if description, ok := d.GetOk(cisGLBPoolDesc); ok {
		opt.SetDescription(description.(string))
	}

	result, resp, err := cisClient.CreateLoadBalancerPool(opt)
	if err != nil {
		log.Printf("[WARN] Create GLB Pools failed %s\n", resp)
		return err
	}
	//Set unique TF Id from concatenated CIS Ids
	d.SetId(convertCisToTfTwoVar(*result.Result.ID, crn))
	return resourceCISPoolRead(d, meta)
}

func resourceCISPoolRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisGLBPoolClientSession()
	if err != nil {
		return err
	}
	poolID, crn, err := convertTftoCisTwoVar(d.Id())
	if err != nil {
		return err
	}
	cisClient.Crn = core.StringPtr(crn)
	opt := cisClient.NewGetLoadBalancerPoolOptions(poolID)
	result, resp, err := cisClient.GetLoadBalancerPool(opt)
	if err != nil {
		log.Printf("[WARN] Create GLB Pools failed %s\n", resp)
		return err
	}

	poolObj := *result.Result
	d.Set(cisID, crn)
	d.Set(cisGLBPoolID, poolObj.ID)
	d.Set(cisGLBPoolName, poolObj.Name)
	d.Set(cisGLBPoolOrigins, flattenOrigins(poolObj.Origins))
	d.Set(cisGLBPoolRegions, poolObj.CheckRegions)
	d.Set(cisGLBPoolDesc, poolObj.Description)
	d.Set(cisGLBPoolEnabled, poolObj.Enabled)
	d.Set(cisGLBPoolNotificationEMail, poolObj.NotificationEmail)
	d.Set(cisGLBPoolHealthy, poolObj.Healthy)
	d.Set(cisGLBPoolMinimumOrigins, poolObj.MinimumOrigins)
	d.Set(cisGLBPoolCreatedOn, poolObj.CreatedOn)
	d.Set(cisGLBPoolModifiedOn, poolObj.ModifiedOn)
	if poolObj.Monitor != nil {
		d.Set(cisGLBPoolMonitor, *poolObj.Monitor)
	}
	return nil
}

func resourceCISPoolUpdate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisGLBPoolClientSession()
	if err != nil {
		return err
	}
	poolID, crn, err := convertTftoCisTwoVar(d.Id())
	if err != nil {
		return err
	}
	cisClient.Crn = core.StringPtr(crn)
	if d.HasChange(cisGLBPoolName) ||
		d.HasChange(cisGLBPoolOrigins) ||
		d.HasChange(cisGLBPoolRegions) ||
		d.HasChange(cisGLBPoolNotificationEMail) ||
		d.HasChange(cisGLBPoolMonitor) ||
		d.HasChange(cisGLBPoolEnabled) ||
		d.HasChange(cisGLBPoolMinimumOrigins) ||
		d.HasChange(cisGLBPoolDesc) {

		opt := cisClient.NewEditLoadBalancerPoolOptions(poolID)
		if monitor, ok := d.GetOk(cisGLBPoolMonitor); ok {
			monitorID, _, _ := convertTftoCisTwoVar(monitor.(string))
			opt.SetMonitor(monitorID)
		}

		if name, ok := d.GetOk(cisGLBPoolName); ok {
			opt.SetName(name.(string))
		}
		if origins, ok := d.GetOk(cisGLBPoolOrigins); ok {
			glbOrigins := []globalloadbalancerpoolsv0.LoadBalancerPoolReqOriginsItem{}

			for _, origin := range origins.(*schema.Set).List() {
				orig := origin.(map[string]interface{})
				glbOrigin := globalloadbalancerpoolsv0.LoadBalancerPoolReqOriginsItem{
					Name:    core.StringPtr(orig[cisGLBPoolOriginsName].(string)),
					Address: core.StringPtr(orig[cisGLBPoolOriginsAddress].(string)),
					Enabled: core.BoolPtr(orig[cisGLBPoolOriginsEnabled].(bool)),
					Weight:  core.Float64Ptr(orig[cisGLBPoolOriginsWeight].(float64)),
				}
				glbOrigins = append(glbOrigins, glbOrigin)
			}
			opt.SetOrigins(glbOrigins)
		}
		if checkregions, ok := d.GetOk(cisGLBPoolRegions); ok {
			checkRegions := checkregions.(*schema.Set).List()
			var regions []string
			for _, region := range checkRegions {
				regions = append(regions, region.(string))
			}
			opt.SetCheckRegions(regions)
		}
		if notEmail, ok := d.GetOk(cisGLBPoolNotificationEMail); ok {
			opt.SetNotificationEmail(notEmail.(string))
		}

		if enabled, ok := d.GetOk(cisGLBPoolEnabled); ok {
			opt.SetEnabled(enabled.(bool))
		}
		if minOrigins, ok := d.GetOk(cisGLBPoolMinimumOrigins); ok {
			opt.SetMinimumOrigins(int64(minOrigins.(int)))
		}
		if description, ok := d.GetOk(cisGLBPoolDesc); ok {
			opt.SetDescription(description.(string))
		}
		_, resp, err := cisClient.EditLoadBalancerPool(opt)
		if err != nil {
			log.Printf("[WARN] Error getting zone during PoolUpdate %v\n", resp)
			return err
		}
	}
	return resourceCISPoolRead(d, meta)
}

func resourceCISPoolDelete(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisGLBPoolClientSession()
	if err != nil {
		return err
	}
	poolID, crn, err := convertTftoCisTwoVar(d.Id())
	if err != nil {
		return err
	}
	cisClient.Crn = core.StringPtr(crn)
	opt := cisClient.NewDeleteLoadBalancerPoolOptions(poolID)
	result, resp, err := cisClient.DeleteLoadBalancerPool(opt)
	if err != nil {
		log.Printf("[WARN] Delete GLB Pools failed %s\n", resp)
		return err
	}
	log.Printf("Pool %s deleted successfully.", *result.Result.ID)
	return nil
}

func resourceCISPoolExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	cisClient, err := meta.(ClientSession).CisGLBPoolClientSession()
	if err != nil {
		return false, err
	}
	poolID, cisID, err := convertTftoCisTwoVar(d.Id())
	if err != nil {
		return false, err
	}
	cisClient.Crn = core.StringPtr(cisID)
	opt := cisClient.NewGetLoadBalancerPoolOptions(poolID)
	result, response, err := cisClient.GetLoadBalancerPool(opt)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			log.Printf("global load balancer pool does not exist.")
			return false, nil
		}
		log.Printf("Error : %s", response)
		return false, err
	}
	log.Printf("global load balancer pool exist : %s", *result.Result.ID)
	return true, nil
}

// Cloud Internet Services
func flattenOrigins(list []globalloadbalancerpoolsv0.LoadBalancerPoolPackOriginsItem) []map[string]interface{} {
	origins := []map[string]interface{}{}
	for _, origin := range list {
		l := map[string]interface{}{
			cisGLBPoolOriginsName:    origin.Name,
			cisGLBPoolOriginsAddress: origin.Address,
			cisGLBPoolOriginsEnabled: origin.Enabled,
			cisGLBPoolOriginsHealthy: origin.Healthy,
			cisGLBPoolOriginsWeight:  origin.Weight,
		}
		if origin.DisabledAt != nil {
			l[cisGLBPoolOriginsDisabledAt] = *origin.DisabledAt
		}
		if origin.FailureReason != nil {
			l[cisGLBPoolOriginsFailureReason] = *origin.FailureReason
		}
		origins = append(origins, l)
	}
	return origins
}
