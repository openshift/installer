// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IBM/platform-services-go-sdk/catalogmanagementv1"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

func dataSourceIBMISEndpointGatewayTargets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISEndpointGatewayTargetsRead,

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

func dataSourceIBMISEndpointGatewayTargetsRead(d *schema.ResourceData, meta interface{}) error {
	bmxSess, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	region := bmxSess.Config.Region
	catalogManagementClient, err := meta.(ClientSession).CatalogManagementV1()
	if err != nil {
		return err
	}

	getCatalogOptions := &catalogmanagementv1.SearchObjectsOptions{}
	// query := "kind%3Avpe+AND+svc+AND+parent_id%3Aus-south"
	query := fmt.Sprintf("kind:vpe AND svc AND parent_id:%s", region)
	getCatalogOptions.Query = &query
	digest := false
	getCatalogOptions.Digest = &digest
	catalog, response, err := catalogManagementClient.SearchObjectsWithContext(context.TODO(), getCatalogOptions)
	if err != nil {
		log.Printf("[DEBUG] GetCatalogWithContext failed %s\n%s", err, response)
		return err
	}
	if catalog != nil && *catalog.ResourceCount > 0 && catalog.Resources != nil {
		resourceInfo := make([]map[string]interface{}, 0)
		for _, res := range catalog.Resources {
			l := map[string]interface{}{}
			if res.ParentID != nil {
				l[isVPEResourceParent] = *res.ParentID
			}
			l[isVPEResourceName] = "provider_cloud_service"
			if res.Label != nil {
				l[isVPEResourceType] = *res.Label
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
		d.Set(isVPEResources, resourceInfo)
		d.SetId(dataSourceIBMISEndpointGatewayTargetsId(d))
	}
	return nil
}
func dataSourceIBMISEndpointGatewayTargetsId(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
