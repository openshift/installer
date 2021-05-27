// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"

	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	cisGLBID                 = "glb_id"
	cisGLBName               = "name"
	cisGLBFallbackPoolID     = "fallback_pool_id"
	cisGLBDefaultPoolIDs     = "default_pool_ids"
	cisGLBDesc               = "description"
	cisGLBProxied            = "proxied"
	cisGLBTTL                = "ttl"
	cisGLBSessionAffinity    = "session_affinity"
	cisGLBEnabled            = "enabled"
	cisGLBPopPools           = "pop_pools"
	cisGLBPopPoolsPop        = "pop"
	cisGLBPopPoolsPoolIDs    = "pool_ids"
	cisGLBRegionPools        = "region_pools"
	cisGLBRegionPoolsRegion  = "region"
	cisGLBRegionPoolsPoolIDs = "pool_ids"
	cisGLBCreatedOn          = "created_on"
	cisGLBModifiedOn         = "modified_on"
)

func resourceIBMCISGlb() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Description:      "Associated CIS domain",
				Required:         true,
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisGLBID: {
				Type:        schema.TypeString,
				Description: "global load balancer id",
				Computed:    true,
			},
			cisGLBName: {
				Type:        schema.TypeString,
				Description: "name",
				Required:    true,
			},
			cisGLBFallbackPoolID: {
				Type:        schema.TypeString,
				Description: "fallback pool ID",
				Required:    true,
			},
			cisGLBDefaultPoolIDs: {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of default Pool IDs",
			},
			cisGLBDesc: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description for the load balancer instance",
			},
			cisGLBTTL: {
				Type:          schema.TypeInt,
				Optional:      true,
				Default:       60,
				ConflictsWith: []string{"proxied"},
				Description:   "TTL value", // this is set to zero regardless of config when proxied=true

			},
			cisGLBProxied: {
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       false,
				ConflictsWith: []string{cisGLBTTL},
				Description:   "set to true if proxy needs to be enabled",
			},
			cisGLBSessionAffinity: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "none",
				// Set to cookie when proxy=true
				ValidateFunc: validateAllowedStringValue([]string{"none", "cookie"}),
				Description:  "Session affinity info",
			},
			cisGLBEnabled: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "set to true of LB needs to enabled",
			},
			cisGLBPopPools: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						cisGLBPopPoolsPop: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "pop pools region",
						},

						cisGLBPopPoolsPoolIDs: {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			cisGLBRegionPools: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						cisGLBRegionPoolsRegion: {
							Type:     schema.TypeString,
							Required: true,
						},

						cisGLBRegionPoolsPoolIDs: {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			cisGLBCreatedOn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Load balancer creation date",
			},
			cisGLBModifiedOn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Load balancer modified date",
			},
		},

		Create:   resourceCISGlbCreate,
		Read:     resourceCISGlbRead,
		Update:   resourceCISGlbUpdate,
		Delete:   resourceCISGlbDelete,
		Importer: &schema.ResourceImporter{},
	}
}

func resourceCISGlbCreate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisGLBClientSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	zoneID, _, err := convertTftoCisTwoVar(d.Get(cisDomainID).(string))
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	tfDefaultPoolIds := expandStringList(d.Get(cisGLBDefaultPoolIDs).(*schema.Set).List())
	defaultPoolIds, _, err := convertTfToCisTwoVarSlice(tfDefaultPoolIds)
	fbPoolID := d.Get(cisGLBFallbackPoolID).(string)
	fallbackPool, _, err := convertTftoCisTwoVar(fbPoolID)

	opt := cisClient.NewCreateLoadBalancerOptions()
	opt.SetName(d.Get(cisGLBName).(string))
	opt.SetDefaultPools(defaultPoolIds)
	opt.SetFallbackPool(fallbackPool)
	opt.SetProxied(d.Get(cisGLBProxied).(bool))
	opt.SetSessionAffinity(d.Get(cisGLBSessionAffinity).(string))

	if description, ok := d.GetOk(cisGLBDesc); ok {
		opt.SetDescription(description.(string))
	}
	if ttl, ok := d.GetOk(cisGLBTTL); ok {
		opt.SetTTL(int64(ttl.(int)))
	}
	if regionPools, ok := d.GetOk(cisGLBRegionPools); ok {
		expandedRegionPools, err := expandGeoPools(regionPools, cisGLBRegionPoolsRegion)
		if err != nil {
			return err
		}
		opt.SetRegionPools(expandedRegionPools)
	}
	if popPools, ok := d.GetOk(cisGLBPopPools); ok {
		expandedPopPools, err := expandGeoPools(popPools, cisGLBPopPoolsPop)
		if err != nil {
			return err
		}
		opt.SetPopPools(expandedPopPools)
	}

	result, resp, err := cisClient.CreateLoadBalancer(opt)
	if err != nil {
		log.Printf("Create GLB failed %s\n", resp)
		return err
	}
	d.SetId(convertCisToTfThreeVar(*result.Result.ID, zoneID, crn))
	return resourceCISGlbUpdate(d, meta)
}

func resourceCISGlbRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisGLBClientSession()
	if err != nil {
		return err
	}

	// Extract CIS Ids from TF Id
	glbID, zoneID, crn, err := convertTfToCisThreeVar(d.Id())
	if err != nil {
		return err
	}

	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	opt := cisClient.NewGetLoadBalancerSettingsOptions(glbID)

	result, resp, err := cisClient.GetLoadBalancerSettings(opt)
	if err != nil {
		log.Printf("[WARN] GLB Read failed: %v\n", resp)
		return err
	}
	glbObj := result.Result
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisGLBID, glbObj.ID)
	d.Set(cisGLBName, glbObj.Name)
	d.Set(cisGLBDefaultPoolIDs, convertCisToTfTwoVarSlice(glbObj.DefaultPools, crn))
	d.Set(cisGLBDesc, glbObj.Description)
	d.Set(cisGLBFallbackPoolID, convertCisToTfTwoVar(*glbObj.FallbackPool, crn))
	d.Set(cisGLBTTL, glbObj.TTL)
	d.Set(cisGLBProxied, glbObj.Proxied)
	d.Set(cisGLBEnabled, glbObj.Enabled)
	d.Set(cisGLBSessionAffinity, glbObj.SessionAffinity)
	flattenPopPools := flattenPools(
		glbObj.PopPools, cisGLBPopPoolsPop, crn)
	d.Set(cisGLBPopPools, flattenPopPools)
	flattenRegionPools := flattenPools(
		glbObj.RegionPools, cisGLBRegionPoolsRegion, crn)
	d.Set(cisGLBRegionPools, flattenRegionPools)

	return nil
}

func resourceCISGlbUpdate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisGLBClientSession()
	if err != nil {
		return err
	}
	// Extract CIS Ids from TF Id
	glbID, zoneID, crn, err := convertTfToCisThreeVar(d.Id())
	if err != nil {
		return err
	}
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	if d.HasChange(cisGLBName) || d.HasChange(cisGLBDefaultPoolIDs) ||
		d.HasChange(cisGLBFallbackPoolID) || d.HasChange(cisGLBProxied) ||
		d.HasChange(cisGLBSessionAffinity) || d.HasChange(cisGLBDesc) ||
		d.HasChange(cisGLBTTL) || d.HasChange(cisGLBEnabled) ||
		d.HasChange(cisGLBPopPools) || d.HasChange(cisGLBRegionPools) {

		tfDefaultPools := expandStringList(d.Get(cisGLBDefaultPoolIDs).(*schema.Set).List())
		defaultPoolIds, _, err := convertTfToCisTwoVarSlice(tfDefaultPools)
		fbPoolID := d.Get(cisGLBFallbackPoolID).(string)
		fallbackPool, _, _ := convertTftoCisTwoVar(fbPoolID)

		opt := cisClient.NewEditLoadBalancerOptions(glbID)
		opt.SetName(d.Get(cisGLBName).(string))
		opt.SetProxied(d.Get(cisGLBProxied).(bool))
		opt.SetSessionAffinity(d.Get(cisGLBSessionAffinity).(string))
		opt.SetDefaultPools(defaultPoolIds)
		opt.SetFallbackPool(fallbackPool)
		if description, ok := d.GetOk(cisGLBDesc); ok {
			opt.SetDescription(description.(string))
		}
		if ttl, ok := d.GetOk(cisGLBTTL); ok {
			opt.SetTTL(int64(ttl.(int)))
		}
		if enabled, ok := d.GetOk(cisGLBEnabled); ok {
			opt.SetEnabled(enabled.(bool))
		}
		if regionPools, ok := d.GetOk(cisGLBRegionPools); ok {
			expandedRegionPools, err := expandGeoPools(regionPools, cisGLBRegionPoolsRegion)
			if err != nil {
				return err
			}
			opt.SetRegionPools(expandedRegionPools)
		}
		if popPools, ok := d.GetOk(cisGLBPopPools); ok {
			expandedPopPools, err := expandGeoPools(popPools, cisGLBPopPoolsPop)
			if err != nil {
				return err
			}
			opt.SetPopPools(expandedPopPools)
		}

		_, resp, err := cisClient.EditLoadBalancer(opt)
		if err != nil {
			log.Printf("[WARN] Error updating GLB %v\n", resp)
			return err
		}
	}

	return resourceCISGlbRead(d, meta)
}

func resourceCISGlbDelete(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisGLBClientSession()
	if err != nil {
		return err
	}
	// Extract CIS Ids from TF Id
	glbID, zoneID, crn, err := convertTfToCisThreeVar(d.Id())
	if err != nil {
		return err
	}
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)
	opt := cisClient.NewDeleteLoadBalancerOptions(glbID)

	result, resp, err := cisClient.DeleteLoadBalancer(opt)
	if err != nil {
		log.Printf("[WARN] Error deleting GLB %v\n", resp)
		return err
	}
	log.Printf("Deletion successful : %s", *result.Result.ID)
	return nil
}

func resourceCISGlbExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	cisClient, err := meta.(ClientSession).CisGLBClientSession()
	if err != nil {
		return false, err
	}
	// Extract CIS Ids from TF Id
	glbID, zoneID, crn, err := convertTfToCisThreeVar(d.Id())
	if err != nil {
		return false, err
	}
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)
	opt := cisClient.NewGetLoadBalancerSettingsOptions(glbID)

	_, response, err := cisClient.GetLoadBalancerSettings(opt)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			log.Printf("global load balancer does not exist.")
			return false, nil
		}
		log.Printf("[WARN] Error getting GLB %v\n", response)
		return false, err
	}
	return true, nil
}

func expandGeoPools(pool interface{}, geoType string) (map[string][]string, error) {
	pools := pool.(*schema.Set).List()
	expandPool := make(map[string][]string)
	for _, v := range pools {
		locationConfig := v.(map[string]interface{})
		location := locationConfig[geoType].(string)
		if _, p := expandPool[location]; !p {
			geoPools := expandStringList(locationConfig[cisGLBRegionPoolsPoolIDs].([]interface{}))
			expandPool[location], _, _ = convertTfToCisTwoVarSlice(geoPools)
		} else {
			return nil, fmt.Errorf("duplicate entry specified for %s pool in location %q. "+
				"each location must only be specified once", geoType, location)
		}
	}
	return expandPool, nil
}

func flattenPools(pools interface{}, geoType string, cisID string) []interface{} {
	result := make([]interface{}, 0)
	for k, v := range pools.(map[string]interface{}) {
		poolIds := convertCisToTfTwoVarSlice(expandStringList(v.([]interface{})), cisID)
		pool := map[string]interface{}{
			geoType:               k,
			cisGLBPopPoolsPoolIDs: poolIds,
		}
		result = append(result, pool)
	}
	return result
}
