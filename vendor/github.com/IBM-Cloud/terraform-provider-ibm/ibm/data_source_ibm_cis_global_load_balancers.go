// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"log"
	"time"

	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const cisGLB = "cis_glb"

func dataSourceIBMCISGlbs() *schema.Resource {
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
			cisGLB: {
				Type:        schema.TypeList,
				Description: "Collection of GLB detail",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "identifier with zone id",
							Computed:    true,
						},
						cisGLBID: {
							Type:        schema.TypeString,
							Description: "global load balancer id",
							Computed:    true,
						},
						cisGLBName: {
							Type:        schema.TypeString,
							Description: "name",
							Computed:    true,
						},
						cisGLBFallbackPoolID: {
							Type:        schema.TypeString,
							Description: "fallback pool ID",
							Computed:    true,
						},
						cisGLBDefaultPoolIDs: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "List of default Pool IDs",
						},
						cisGLBDesc: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description for the load balancer instance",
						},
						cisGLBTTL: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "TTL value",
						},
						cisGLBProxied: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "set to true if proxy needs to be enabled",
						},
						cisGLBSessionAffinity: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Session affinity info",
						},
						cisGLBEnabled: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "set to true of LB needs to enabled",
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
						cisGLBPopPools: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									cisGLBPopPoolsPop: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "pop pools region",
									},

									cisGLBPopPoolsPoolIDs: {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						cisGLBRegionPools: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									cisGLBRegionPoolsRegion: {
										Type:     schema.TypeString,
										Computed: true,
									},
									cisGLBRegionPoolsPoolIDs: {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(10 * time.Minute),
		},
		Read:     dataSourceCISGlbsRead,
		Importer: &schema.ResourceImporter{},
	}
}

func dataSourceCISGlbsRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisGLBClientSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	zoneID, _, err := convertTftoCisTwoVar(d.Get(cisDomainID).(string))
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	opt := cisClient.NewListAllLoadBalancersOptions()

	result, resp, err := cisClient.ListAllLoadBalancers(opt)
	if err != nil {
		log.Printf("[WARN] List all GLB failed: %v\n", resp)
		return err
	}
	glbs := result.Result

	glbList := make([]map[string]interface{}, 0)
	for _, glbObj := range glbs {
		glbOutput := map[string]interface{}{}
		glbOutput["id"] = convertCisToTfThreeVar(*glbObj.ID, zoneID, crn)
		glbOutput[cisGLBID] = *glbObj.ID
		glbOutput[cisGLBName] = *glbObj.Name
		glbOutput[cisGLBDefaultPoolIDs] = convertCisToTfTwoVarSlice(glbObj.DefaultPools, crn)
		glbOutput[cisGLBDesc] = *glbObj.Description
		glbOutput[cisGLBFallbackPoolID] = convertCisToTfTwoVar(*glbObj.FallbackPool, crn)
		glbOutput[cisGLBTTL] = *glbObj.TTL
		glbOutput[cisGLBProxied] = *glbObj.Proxied
		glbOutput[cisGLBEnabled] = *glbObj.Enabled
		glbOutput[cisGLBSessionAffinity] = *glbObj.SessionAffinity
		glbOutput[cisGLBCreatedOn] = *glbObj.CreatedOn
		glbOutput[cisGLBModifiedOn] = *glbObj.ModifiedOn
		flattenPopPools := flattenDataSourcePopPools(
			glbObj.PopPools, cisGLBPopPoolsPop, crn)
		glbOutput[cisGLBPopPools] = flattenPopPools
		flattenRegionPools := flattenDataSourceRegionPools(
			glbObj.RegionPools, cisGLBRegionPoolsRegion, crn)
		glbOutput[cisGLBRegionPools] = flattenRegionPools
		glbList = append(glbList, glbOutput)
	}
	d.SetId(dataSourceCISGlbsCheckID(d))
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisGLB, glbList)
	return nil
}

// dataSourceCISGlbCheckID returns a reasonable ID glb list
func dataSourceCISGlbsCheckID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func flattenDataSourcePopPools(pools interface{}, geoType string, cisID string) []interface{} {
	result := make([]interface{}, 0)
	for k, v := range pools.(map[string]interface{}) {
		poolIds := convertCisToTfTwoVarSlice(expandStringList(v.([]interface{})), cisID)
		pool := map[string]interface{}{
			cisGLBPopPoolsPop:     k,
			cisGLBPopPoolsPoolIDs: poolIds,
		}
		result = append(result, pool)
	}
	return result
}

func flattenDataSourceRegionPools(pools interface{}, geoType string, cisID string) []interface{} {
	result := make([]interface{}, 0)
	for k, v := range pools.(map[string]interface{}) {
		poolIds := convertCisToTfTwoVarSlice(expandStringList(v.([]interface{})), cisID)
		pool := map[string]interface{}{
			cisGLBRegionPoolsRegion:  k,
			cisGLBRegionPoolsPoolIDs: poolIds,
		}
		result = append(result, pool)
	}
	return result
}
