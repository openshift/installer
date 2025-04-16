// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const cisGLB = "cis_glb"

func DataSourceIBMCISGlbs() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
				ValidateFunc: validate.InvokeDataSourceValidator(
					"ibm_cis_global_load_balancers",
					"cis_id"),
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
						cisGLBSteeringPolicy: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Steering policy info",
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
func DataSourceIBMCISGlbsValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})

	iBMCISGLBsValidator := validate.ResourceValidator{
		ResourceName: "ibm_cis_global_load_balancers",
		Schema:       validateSchema}
	return &iBMCISGLBsValidator
}

func dataSourceCISGlbsRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisGLBClientSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	zoneID, _, _ := flex.ConvertTftoCisTwoVar(d.Get(cisDomainID).(string))
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
		glbOutput["id"] = flex.ConvertCisToTfThreeVar(*glbObj.ID, zoneID, crn)
		glbOutput[cisGLBID] = *glbObj.ID
		glbOutput[cisGLBName] = *glbObj.Name
		glbOutput[cisGLBDefaultPoolIDs] = flex.ConvertCisToTfTwoVarSlice(glbObj.DefaultPools, crn)
		glbOutput[cisGLBDesc] = *glbObj.Description
		glbOutput[cisGLBFallbackPoolID] = flex.ConvertCisToTfTwoVar(*glbObj.FallbackPool, crn)
		glbOutput[cisGLBTTL] = *glbObj.TTL
		glbOutput[cisGLBSteeringPolicy] = *glbObj.SteeringPolicy
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
		poolIds := flex.ConvertCisToTfTwoVarSlice(flex.ExpandStringList(v.([]interface{})), cisID)
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
		poolIds := flex.ConvertCisToTfTwoVarSlice(flex.ExpandStringList(v.([]interface{})), cisID)
		pool := map[string]interface{}{
			cisGLBRegionPoolsRegion:  k,
			cisGLBRegionPoolsPoolIDs: poolIds,
		}
		result = append(result, pool)
	}
	return result
}
