// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/IBM-Cloud/bluemix-go/api/globalsearch/globalsearchv2"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/platform-services-go-sdk/catalogmanagementv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isVPEResources                        = "resources"
	isVPEResourceCRN                      = "crn"
	isVPEResourceParent                   = "parent"
	isVPEResourceName                     = "name"
	isVPEResourceEndpointType             = "endpoint_type"
	isVPEResourceType                     = "resource_type"
	isVPEResourceFullQualifiedDomainNames = "full_qualified_domain_names"
	isVPEResourceServiceLocation          = "location"
)

func DataSourceIBMISEndpointGatewayTargets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISEndpointGatewayTargetsRead,

		Schema: map[string]*schema.Schema{
			isVPEResources: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of resources",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isVPEResourceCRN: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CRN for this specific object",
						},
						isVPEResourceParent: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parent for this specific object",
						},
						isVPEResourceName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Display name in the requested language",
						},
						isVPEResourceType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource type of this offering.",
						},
						isVPEResourceEndpointType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data endpoint type of this offering",
						},
						isVPEResourceFullQualifiedDomainNames: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Fully qualified domain names",
						},
						isVPEResourceServiceLocation: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service location of this offering",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMISEndpointGatewayTargetsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	bmxSess, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return diag.FromErr(err)
	}
	region := bmxSess.Config.Region
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		return diag.FromErr(err)
	}
	resourceInfo := make([]map[string]interface{}, 0)
	getCatalogOptions := &catalogmanagementv1.SearchObjectsOptions{}
	// query := "kind%3Avpe+AND+svc+AND+parent_id%3Aus-south"
	query := fmt.Sprintf("kind:vpe AND svc AND parent_id:%s", region)
	getCatalogOptions.Query = &query
	digest := false
	getCatalogOptions.Digest = &digest

	start := int64(0)
	catalog := []catalogmanagementv1.CatalogObject{}
	for {
		if start != int64(0) {
			getCatalogOptions.Offset = &start
		}
		search, response, err := catalogManagementClient.SearchObjectsWithContext(context, getCatalogOptions)
		if err != nil {
			log.Printf("[DEBUG] GetCatalogWithContext failed %s\n%s", err, response)
			return diag.FromErr(err)
		}
		next := search.Next
		if next == nil {
			start = int64(0)
		} else {
			u, _ := url.Parse(fmt.Sprintf("%s", *next))
			q := u.Query()
			start, _ = strconv.ParseInt(q.Get("offset"), 10, 64)
		}
		catalog = append(catalog, search.Resources...)
		if start == int64(0) {
			break
		}
	}
	if catalog != nil {
		for _, res := range catalog {
			l := map[string]interface{}{}
			if res.ParentID != nil {
				l[isVPEResourceParent] = *res.ParentID
			}
			l[isVPEResourceType] = "provider_cloud_service"
			if res.Label != nil {
				l[isVPEResourceName] = *res.Label
			}
			sl := ""
			data := res.Data
			if data != nil {
				if serviceCrn, ok := data["service_crn"].(string); ok {
					if serviceCrn != "" {
						l[isVPEResourceCRN] = serviceCrn
						crnFs := strings.Split(serviceCrn, ":")
						if len(crnFs) > 5 {
							sl = crnFs[5]
						}
						l[isVPEResourceServiceLocation] = sl
					}
				}
				if data["endpoint_type"] != nil {
					l[isVPEResourceEndpointType] = data["endpoint_type"]
				}
				if data["fully_qualified_domain_names"] != nil {
					l[isVPEResourceFullQualifiedDomainNames] = data["fully_qualified_domain_names"]
				}
			}
			resourceInfo = append(resourceInfo, l)
		}
		staticService := map[string]interface{}{}
		staticService[isVPEResourceName] = "ibm-ntp-server"
		staticService[isVPEResourceType] = "provider_infrastructure_service"
		resourceInfo = append(resourceInfo, staticService)
	}
	queryString := "doc.extensions.virtual_private_endpoints.endpoints.ip_address:*"
	fields := []string{"name", "region", "family", "type", "crn", "tags", "organization_guid", "doc.extensions", "doc.resource_group_id", "doc.space_guid", "resource_id"}
	getSearchOptions := globalsearchv2.SearchBody{
		Query:  queryString,
		Fields: fields,
	}
	globalSearchClient, err := meta.(conns.ClientSession).GlobalSearchAPI()
	if err != nil {
		return diag.FromErr(err)
	}
	searchResult, err := globalSearchClient.Searches().PostQuery(getSearchOptions)
	if err != nil {
		log.Printf("[DEBUG] PostQuery on globalSearchApi for query string %s failed %s", queryString, err)
		return diag.FromErr(err)
	}
	searchItems := searchResult.Items
	for _, item := range searchItems {
		info := map[string]interface{}{}
		info[isVPEResourceName] = item.Name
		info[isVPEResourceType] = item.Type
		info[isVPEResourceCRN] = item.CRN
		resourceInfo = append(resourceInfo, info)

	}

	d.Set(isVPEResources, resourceInfo)
	d.SetId(dataSourceIBMISEndpointGatewayTargetsId(d))
	return nil
}
func dataSourceIBMISEndpointGatewayTargetsId(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
