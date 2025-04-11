// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.94.1-71478489-20240820-161623
 */

package codeengine

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/code-engine-go-sdk/codeenginev2"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIbmCodeEngineDomainMapping() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmCodeEngineDomainMappingRead,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the project.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of your domain mapping.",
			},
			"cname_target": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The value of the CNAME record that must be configured in the DNS settings of the domain, to route traffic properly to the target Code Engine region.",
			},
			"component": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A reference to another component.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the referenced component.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the referenced resource.",
						},
					},
				},
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the resource was created.",
			},
			"domain_mapping_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The identifier of the resource.",
			},
			"entity_tag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the domain mapping instance, which is used to achieve optimistic locking.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When you provision a new domain mapping, a URL is created identifying the location of the instance.",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The region of the project the resource is located in. Possible values: 'au-syd', 'br-sao', 'ca-tor', 'eu-de', 'eu-gb', 'jp-osa', 'jp-tok', 'us-east', 'us-south'.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the Code Engine resource.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the domain mapping.",
			},
			"status_details": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The detailed status of the domain mapping.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Optional information to provide more context in case of a 'failed' or 'warning' status.",
						},
					},
				},
			},
			"tls_secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the TLS secret that includes the certificate and private key of this domain mapping.",
			},
			"user_managed": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Specifies whether the domain mapping is managed by the user or by Code Engine.",
			},
			"visibility": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies whether the domain mapping is reachable through the public internet, or private IBM network, or only through other components within the same Code Engine project.",
			},
		},
	}
}

func dataSourceIbmCodeEngineDomainMappingRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_code_engine_domain_mapping", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getDomainMappingOptions := &codeenginev2.GetDomainMappingOptions{}

	getDomainMappingOptions.SetProjectID(d.Get("project_id").(string))
	getDomainMappingOptions.SetName(d.Get("name").(string))

	domainMapping, _, err := codeEngineClient.GetDomainMappingWithContext(context, getDomainMappingOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetDomainMappingWithContext failed: %s", err.Error()), "(Data) ibm_code_engine_domain_mapping", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *getDomainMappingOptions.ProjectID, *getDomainMappingOptions.Name))

	if !core.IsNil(domainMapping.CnameTarget) {
		if err = d.Set("cname_target", domainMapping.CnameTarget); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting cname_target: %s", err), "(Data) ibm_code_engine_domain_mapping", "read", "set-cname_target").GetDiag()
		}
	}

	component := []map[string]interface{}{}
	componentMap, err := DataSourceIbmCodeEngineDomainMappingComponentRefToMap(domainMapping.Component)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_code_engine_domain_mapping", "read", "component-to-map").GetDiag()
	}
	component = append(component, componentMap)
	if err = d.Set("component", component); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting component: %s", err), "(Data) ibm_code_engine_domain_mapping", "read", "set-component").GetDiag()
	}

	if !core.IsNil(domainMapping.CreatedAt) {
		if err = d.Set("created_at", domainMapping.CreatedAt); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting created_at: %s", err), "(Data) ibm_code_engine_domain_mapping", "read", "set-created_at").GetDiag()
		}
	}

	if err = d.Set("entity_tag", domainMapping.EntityTag); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting entity_tag: %s", err), "(Data) ibm_code_engine_domain_mapping", "read", "set-entity_tag").GetDiag()
	}

	if !core.IsNil(domainMapping.Href) {
		if err = d.Set("href", domainMapping.Href); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_code_engine_domain_mapping", "read", "set-href").GetDiag()
		}
	}

	if !core.IsNil(domainMapping.ID) {
		if err = d.Set("domain_mapping_id", domainMapping.ID); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting domain_mapping_id: %s", err), "(Data) ibm_code_engine_domain_mapping", "read", "set-domain_mapping_id").GetDiag()
		}
	}

	if !core.IsNil(domainMapping.Region) {
		if err = d.Set("region", domainMapping.Region); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting region: %s", err), "(Data) ibm_code_engine_domain_mapping", "read", "set-region").GetDiag()
		}
	}

	if !core.IsNil(domainMapping.ResourceType) {
		if err = d.Set("resource_type", domainMapping.ResourceType); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting resource_type: %s", err), "(Data) ibm_code_engine_domain_mapping", "read", "set-resource_type").GetDiag()
		}
	}

	if !core.IsNil(domainMapping.Status) {
		if err = d.Set("status", domainMapping.Status); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting status: %s", err), "(Data) ibm_code_engine_domain_mapping", "read", "set-status").GetDiag()
		}
	}

	if !core.IsNil(domainMapping.StatusDetails) {
		statusDetails := []map[string]interface{}{}
		statusDetailsMap, err := DataSourceIbmCodeEngineDomainMappingDomainMappingStatusToMap(domainMapping.StatusDetails)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_code_engine_domain_mapping", "read", "status_details-to-map").GetDiag()
		}
		statusDetails = append(statusDetails, statusDetailsMap)
		if err = d.Set("status_details", statusDetails); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting status_details: %s", err), "(Data) ibm_code_engine_domain_mapping", "read", "set-status_details").GetDiag()
		}
	}

	if err = d.Set("tls_secret", domainMapping.TlsSecret); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting tls_secret: %s", err), "(Data) ibm_code_engine_domain_mapping", "read", "set-tls_secret").GetDiag()
	}

	if !core.IsNil(domainMapping.UserManaged) {
		if err = d.Set("user_managed", domainMapping.UserManaged); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting user_managed: %s", err), "(Data) ibm_code_engine_domain_mapping", "read", "set-user_managed").GetDiag()
		}
	}

	if !core.IsNil(domainMapping.Visibility) {
		if err = d.Set("visibility", domainMapping.Visibility); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting visibility: %s", err), "(Data) ibm_code_engine_domain_mapping", "read", "set-visibility").GetDiag()
		}
	}

	return nil
}

func DataSourceIbmCodeEngineDomainMappingComponentRefToMap(model *codeenginev2.ComponentRef) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIbmCodeEngineDomainMappingDomainMappingStatusToMap(model *codeenginev2.DomainMappingStatus) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Reason != nil {
		modelMap["reason"] = *model.Reason
	}
	return modelMap, nil
}
