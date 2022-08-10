// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/scc-go-sdk/v3/posturemanagementv2"
)

func DataSourceIBMSccPostureCredential() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMSccPostureCredentialRead,

		Schema: map[string]*schema.Schema{
			"credential_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id for the given API.",
			},
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Credentials status enabled/disbaled.",
			},
			"credential_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Credentials type.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Credentials name.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Credentials description.",
			},
			"display_fields": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Details the fields on the credential. This will change as per credential type selected.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ibm_api_key": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IBM Cloud API Key. This is mandatory for IBM Credential Type.",
						},
						"aws_client_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AWS client Id.This is mandatory for AWS Cloud.",
						},
						"aws_client_secret": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AWS client secret.This is mandatory for AWS Cloud.",
						},
						"aws_region": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AWS region.",
						},
						"aws_arn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AWS arn value.",
						},
						"username": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "username of the user.This is mandatory for DataBase, Kerbros,OpenStack Credentials.",
						},
						"password": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "password of the user.This is mandatory for DataBase, Kerbros,OpenStack Credentials.",
						},
						"azure_client_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Azure client Id. This is mandatory for Azure Credential type.",
						},
						"azure_client_secret": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Azure client secret.This is mandatory for Azure Credential type.",
						},
						"azure_subscription_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Azure subscription Id.This is mandatory for Azure Credential type.",
						},
						"azure_resource_group": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Azure resource group.",
						},
						"database_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database name.This is mandatory for Database Credential type.",
						},
						"winrm_authtype": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Kerberos windows auth type.This is mandatory for Windows Kerberos Credential type.",
						},
						"winrm_usessl": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Kerberos windows ssl.This is mandatory for Windows Kerberos Credential type.",
						},
						"winrm_port": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Kerberos windows port.This is mandatory for Windows Kerberos Credential type.",
						},
						"ms_365_client_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The MS365 client Id.This is mandatory for Windows MS365 Credential type.",
						},
						"ms_365_client_secret": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The MS365 client secret.This is mandatory for Windows MS365 Credential type.",
						},
						"ms_365_tenant_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The MS365 tenantId.This is mandatory for Windows MS365 Credential type.",
						},
						"auth_url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "auth url of the Open Stack cloud.This is mandatory for Open Stack Credential type.",
						},
						"project_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project name of the Open Stack cloud.This is mandatory for Open Stack Credential type.",
						},
						"user_domain_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "user domain name of the Open Stack cloud.This is mandatory for Open Stack Credential type.",
						},
						"project_domain_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "project domain name of the Open Stack cloud.This is mandatory for Open Stack Credential type.",
						},
						"pem_file_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the PEM file.",
						},
						"pem_data": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The base64 encoded data to associate with the PEM file.",
						},
					},
				},
			},
			"created_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the user who created the credentials.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time that the credentials was created in UTC.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The modified time that the credentials was modified in UTC.",
			},
			"updated_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the user who modified the credentials.",
			},
			"group": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Credential group details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "credential group id.",
						},
						"passphrase": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "passphase of the credential.",
						},
					},
				},
			},
			"purpose": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Purpose for which the credential is created.",
			},
		},
	}
}

func dataSourceIBMSccPostureCredentialRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	postureManagementClient, err := meta.(conns.ClientSession).PostureManagementV2()
	if err != nil {
		return diag.FromErr(err)
	}

	CredentialOptions := &posturemanagementv2.GetCredentialOptions{}
	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error getting userDetails %s", err))
	}

	accountID := userDetails.UserAccount
	CredentialOptions.SetAccountID(accountID)
	CredentialOptions.SetID(d.Get("credential_id").(string))

	credential, response, err := postureManagementClient.GetCredentialWithContext(context, CredentialOptions)
	if err != nil {
		log.Printf("[DEBUG] GetCredentialWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetCredentialWithContext failed %s\n%s", err, response))
	}

	d.SetId(*credential.ID)
	if err = d.Set("enabled", credential.Enabled); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting enabled: %s", err))
	}
	if err = d.Set("credential_type", credential.Type); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting type: %s", err))
	}
	if err = d.Set("name", credential.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("description", credential.Description); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
	}

	if credential.DisplayFields != nil {
		err = d.Set("display_fields", dataSourceCredentialFlattenDisplayFields(*credential.DisplayFields))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting display_fields %s", err))
		}
	}
	if err = d.Set("created_by", credential.CreatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_by: %s", err))
	}
	if err = d.Set("created_at", flex.DateTimeToString(credential.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("updated_at", flex.DateTimeToString(credential.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}
	if err = d.Set("updated_by", credential.UpdatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_by: %s", err))
	}

	if credential.Group != nil {
		err = d.Set("group", dataSourceCredentialFlattenGroup(*credential.Group))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting group %s", err))
		}
	}
	if err = d.Set("purpose", credential.Purpose); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting purpose: %s", err))
	}

	return nil
}

func dataSourceCredentialFlattenDisplayFields(result posturemanagementv2.CredentialDisplayFields) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceCredentialDisplayFieldsToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceCredentialDisplayFieldsToMap(displayFieldsItem posturemanagementv2.CredentialDisplayFields) (displayFieldsMap map[string]interface{}) {
	displayFieldsMap = map[string]interface{}{}

	if displayFieldsItem.IBMAPIKey != nil {
		displayFieldsMap["ibm_api_key"] = displayFieldsItem.IBMAPIKey
	}
	if displayFieldsItem.AwsClientID != nil {
		displayFieldsMap["aws_client_id"] = displayFieldsItem.AwsClientID
	}
	if displayFieldsItem.AwsClientSecret != nil {
		displayFieldsMap["aws_client_secret"] = displayFieldsItem.AwsClientSecret
	}
	if displayFieldsItem.AwsRegion != nil {
		displayFieldsMap["aws_region"] = displayFieldsItem.AwsRegion
	}
	if displayFieldsItem.AwsArn != nil {
		displayFieldsMap["aws_arn"] = displayFieldsItem.AwsArn
	}
	if displayFieldsItem.Username != nil {
		displayFieldsMap["username"] = displayFieldsItem.Username
	}
	if displayFieldsItem.Password != nil {
		displayFieldsMap["password"] = displayFieldsItem.Password
	}
	if displayFieldsItem.AzureClientID != nil {
		displayFieldsMap["azure_client_id"] = displayFieldsItem.AzureClientID
	}
	if displayFieldsItem.AzureClientSecret != nil {
		displayFieldsMap["azure_client_secret"] = displayFieldsItem.AzureClientSecret
	}
	if displayFieldsItem.AzureSubscriptionID != nil {
		displayFieldsMap["azure_subscription_id"] = displayFieldsItem.AzureSubscriptionID
	}
	if displayFieldsItem.AzureResourceGroup != nil {
		displayFieldsMap["azure_resource_group"] = displayFieldsItem.AzureResourceGroup
	}
	if displayFieldsItem.DatabaseName != nil {
		displayFieldsMap["database_name"] = displayFieldsItem.DatabaseName
	}
	if displayFieldsItem.WinrmAuthtype != nil {
		displayFieldsMap["winrm_authtype"] = displayFieldsItem.WinrmAuthtype
	}
	if displayFieldsItem.WinrmUsessl != nil {
		displayFieldsMap["winrm_usessl"] = displayFieldsItem.WinrmUsessl
	}
	if displayFieldsItem.WinrmPort != nil {
		displayFieldsMap["winrm_port"] = displayFieldsItem.WinrmPort
	}
	if displayFieldsItem.Ms365ClientID != nil {
		displayFieldsMap["ms_365_client_id"] = displayFieldsItem.Ms365ClientID
	}
	if displayFieldsItem.Ms365ClientSecret != nil {
		displayFieldsMap["ms_365_client_secret"] = displayFieldsItem.Ms365ClientSecret
	}
	if displayFieldsItem.Ms365TenantID != nil {
		displayFieldsMap["ms_365_tenant_id"] = displayFieldsItem.Ms365TenantID
	}
	if displayFieldsItem.AuthURL != nil {
		displayFieldsMap["auth_url"] = displayFieldsItem.AuthURL
	}
	if displayFieldsItem.ProjectName != nil {
		displayFieldsMap["project_name"] = displayFieldsItem.ProjectName
	}
	if displayFieldsItem.UserDomainName != nil {
		displayFieldsMap["user_domain_name"] = displayFieldsItem.UserDomainName
	}
	if displayFieldsItem.ProjectDomainName != nil {
		displayFieldsMap["project_domain_name"] = displayFieldsItem.ProjectDomainName
	}
	if displayFieldsItem.PemFileName != nil {
		displayFieldsMap["pem_file_name"] = displayFieldsItem.PemFileName
	}
	if displayFieldsItem.PemData != nil {
		displayFieldsMap["pem_data"] = displayFieldsItem.PemData
	}

	return displayFieldsMap
}

func dataSourceCredentialFlattenGroup(result posturemanagementv2.CredentialGroup) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceCredentialGroupToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceCredentialGroupToMap(groupItem posturemanagementv2.CredentialGroup) (groupMap map[string]interface{}) {
	groupMap = map[string]interface{}{}

	if groupItem.ID != nil {
		groupMap["id"] = groupItem.ID
	}
	if groupItem.Passphrase != nil {
		groupMap["passphrase"] = groupItem.Passphrase
	}

	return groupMap
}
