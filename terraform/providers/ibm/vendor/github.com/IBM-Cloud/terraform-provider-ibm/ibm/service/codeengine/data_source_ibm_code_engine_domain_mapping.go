// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package codeengine

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/code-engine-go-sdk/codeenginev2"
)

func DataSourceIbmCodeEngineDomainMapping() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmCodeEngineDomainMappingRead,

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the project.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of your domain mapping.",
			},
			"cname_target": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Exposes the value of the CNAME record that needs to be configured in the DNS settings of the domain, to route traffic properly to the target Code Engine region.",
			},
			"component": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A reference to another component.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the referenced component.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the referenced resource.",
						},
					},
				},
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the resource was created.",
			},
			"domain_mapping_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The identifier of the resource.",
			},
			"entity_tag": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the domain mapping instance, which is used to achieve optimistic locking.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When you provision a new domain mapping, a URL is created identifying the location of the instance.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the CE Resource.",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the domain mapping.",
			},
			"status_details": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The detailed status of the domain mapping.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reason": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Optional information to provide more context in case of a 'failed' or 'warning' status.",
						},
					},
				},
			},
			"tls_secret": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the TLS secret that holds the certificate and private key of this domain mapping.",
			},
			"user_managed": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Exposes whether the domain mapping is managed by the user or by Code Engine.",
			},
			"visibility": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Exposes whether the domain mapping is reachable through the public internet, or private IBM network, or only through other components within the same Code Engine project.",
			},
		},
	}
}

func dataSourceIbmCodeEngineDomainMappingRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getDomainMappingOptions := &codeenginev2.GetDomainMappingOptions{}

	getDomainMappingOptions.SetProjectID(d.Get("project_id").(string))
	getDomainMappingOptions.SetName(d.Get("name").(string))

	domainMapping, response, err := codeEngineClient.GetDomainMappingWithContext(context, getDomainMappingOptions)
	if err != nil {
		log.Printf("[DEBUG] GetDomainMappingWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetDomainMappingWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *getDomainMappingOptions.ProjectID, *getDomainMappingOptions.Name))

	if err = d.Set("cname_target", domainMapping.CnameTarget); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting cname_target: %s", err))
	}

	component := []map[string]interface{}{}
	if domainMapping.Component != nil {
		modelMap, err := dataSourceIbmCodeEngineDomainMappingComponentRefToMap(domainMapping.Component)
		if err != nil {
			return diag.FromErr(err)
		}
		component = append(component, modelMap)
	}
	if err = d.Set("component", component); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting component %s", err))
	}

	if err = d.Set("created_at", domainMapping.CreatedAt); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}

	if err = d.Set("entity_tag", domainMapping.EntityTag); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting entity_tag: %s", err))
	}

	if err = d.Set("href", domainMapping.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}

	if err = d.Set("domain_mapping_id", domainMapping.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting domain_mapping_id: %s", err))
	}

	if err = d.Set("resource_type", domainMapping.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}

	if err = d.Set("status", domainMapping.Status); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting status: %s", err))
	}

	statusDetails := []map[string]interface{}{}
	if domainMapping.StatusDetails != nil {
		modelMap, err := dataSourceIbmCodeEngineDomainMappingDomainMappingStatusToMap(domainMapping.StatusDetails)
		if err != nil {
			return diag.FromErr(err)
		}
		statusDetails = append(statusDetails, modelMap)
	}
	if err = d.Set("status_details", statusDetails); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting status_details %s", err))
	}

	if err = d.Set("tls_secret", domainMapping.TlsSecret); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting tls_secret: %s", err))
	}

	if err = d.Set("user_managed", domainMapping.UserManaged); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting user_managed: %s", err))
	}

	if err = d.Set("visibility", domainMapping.Visibility); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting visibility: %s", err))
	}

	return nil
}

func dataSourceIbmCodeEngineDomainMappingComponentRefToMap(model *codeenginev2.ComponentRef) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func dataSourceIbmCodeEngineDomainMappingDomainMappingStatusToMap(model *codeenginev2.DomainMappingStatus) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Reason != nil {
		modelMap["reason"] = model.Reason
	}
	return modelMap, nil
}
