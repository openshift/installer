// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.96.0-d6dec9d7-20241008-212902
 */

package db2

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/cloud-db2-go-sdk/db2saasv1"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
)

func DataSourceIbmDb2ConnectionInfo() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmDb2ConnectionInfoRead,

		Schema: map[string]*schema.Schema{
			"deployment_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Encoded CRN deployment id.",
			},
			"x_deployment_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "CRN deployment id.",
			},
			"public": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostname": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"database_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssl_port": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssl": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"database_version": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"private": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostname": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"database_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssl_port": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssl": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"database_version": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_service_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloud_service_offering": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpe_service_crn": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_vpc_endpoint_service": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmDb2ConnectionInfoRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	db2saasClient, err := meta.(conns.ClientSession).Db2saasV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_db2_connection_info", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getDb2SaasConnectionInfoOptions := &db2saasv1.GetDb2SaasConnectionInfoOptions{}

	getDb2SaasConnectionInfoOptions.SetDeploymentID(d.Get("deployment_id").(string))
	getDb2SaasConnectionInfoOptions.SetXDeploymentID(d.Get("x_deployment_id").(string))

	successConnectionInfo, _, err := db2saasClient.GetDb2SaasConnectionInfoWithContext(context, getDb2SaasConnectionInfoOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetDb2SaasConnectionInfoWithContext failed: %s", err.Error()), "(Data) ibm_db2_connection_info", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmDb2ConnectionInfoID(d))

	if !core.IsNil(successConnectionInfo.Public) {
		public := []map[string]interface{}{}
		publicMap, err := DataSourceIbmDb2ConnectionInfoSuccessConnectionInfoPublicToMap(successConnectionInfo.Public)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_db2_connection_info", "read", "public-to-map").GetDiag()
		}
		public = append(public, publicMap)
		if err = d.Set("public", public); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting public: %s", err), "(Data) ibm_db2_connection_info", "read", "set-public").GetDiag()
		}
	}

	if !core.IsNil(successConnectionInfo.Private) {
		private := []map[string]interface{}{}
		privateMap, err := DataSourceIbmDb2ConnectionInfoSuccessConnectionInfoPrivateToMap(successConnectionInfo.Private)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_db2_connection_info", "read", "private-to-map").GetDiag()
		}
		private = append(private, privateMap)
		if err = d.Set("private", private); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting private: %s", err), "(Data) ibm_db2_connection_info", "read", "set-private").GetDiag()
		}
	}

	return nil
}

// dataSourceIbmDb2SaasConnectionInfoID returns a reasonable ID for the list.
func dataSourceIbmDb2ConnectionInfoID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIbmDb2ConnectionInfoSuccessConnectionInfoPublicToMap(model *db2saasv1.SuccessConnectionInfoPublic) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Hostname != nil {
		modelMap["hostname"] = *model.Hostname
	}
	if model.DatabaseName != nil {
		modelMap["database_name"] = *model.DatabaseName
	}
	if model.SslPort != nil {
		modelMap["ssl_port"] = *model.SslPort
	}
	if model.Ssl != nil {
		modelMap["ssl"] = *model.Ssl
	}
	if model.DatabaseVersion != nil {
		modelMap["database_version"] = *model.DatabaseVersion
	}
	return modelMap, nil
}

func DataSourceIbmDb2ConnectionInfoSuccessConnectionInfoPrivateToMap(model *db2saasv1.SuccessConnectionInfoPrivate) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Hostname != nil {
		modelMap["hostname"] = *model.Hostname
	}
	if model.DatabaseName != nil {
		modelMap["database_name"] = *model.DatabaseName
	}
	if model.SslPort != nil {
		modelMap["ssl_port"] = *model.SslPort
	}
	if model.Ssl != nil {
		modelMap["ssl"] = *model.Ssl
	}
	if model.DatabaseVersion != nil {
		modelMap["database_version"] = *model.DatabaseVersion
	}
	if model.PrivateServiceName != nil {
		modelMap["private_service_name"] = *model.PrivateServiceName
	}
	if model.CloudServiceOffering != nil {
		modelMap["cloud_service_offering"] = *model.CloudServiceOffering
	}
	if model.VpeServiceCrn != nil {
		modelMap["vpe_service_crn"] = *model.VpeServiceCrn
	}
	if model.DbVpcEndpointService != nil {
		modelMap["db_vpc_endpoint_service"] = *model.DbVpcEndpointService
	}
	return modelMap, nil
}
