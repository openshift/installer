// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/scc-go-sdk/v3/posturemanagementv2"
)

func DataSourceIBMSccPostureCredentials() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMSccPostureListCredentialsRead,

		Schema: map[string]*schema.Schema{
			"first": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The URL of a page.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL of a page.",
						},
					},
				},
			},
			"last": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The URL of a page.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL of a page.",
						},
					},
				},
			},
			"previous": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The URL of a page.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL of a page.",
						},
					},
				},
			},
			"credentials": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The details of a credentials.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Credentials status enabled/disbaled.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Credentials ID.",
						},
						"type": &schema.Schema{
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
										Sensitive:   true,
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
										Sensitive:   true,
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
										Sensitive:   true,
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
										Sensitive:   true,
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
										Sensitive:   true,
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
				},
			},
		},
	}
}

func dataSourceIBMSccPostureListCredentialsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	postureManagementClient, err := meta.(conns.ClientSession).PostureManagementV2()
	if err != nil {
		return diag.FromErr(err)
	}

	listCredentialsOptions := &posturemanagementv2.ListCredentialsOptions{}
	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting userDetails %s", err))
	}

	accountID := userDetails.UserAccount
	listCredentialsOptions.SetAccountID(accountID)

	var credentialList *posturemanagementv2.CredentialList

	listCredentialsOptions.Limit = core.Int64Ptr(int64(2))
	result, response, err := postureManagementClient.ListCredentialsWithContext(context, listCredentialsOptions)
	credentialList = result
	if err != nil {
		log.Printf("[DEBUG] ListCredentialsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListCredentialsWithContext failed %s\n%s", err, response))
	}

	d.SetId(dataSourceIBMSccPostureListCredentialsID(d))

	if credentialList.First != nil {
		err = d.Set("first", dataSourceCredentialListFlattenFirst(*credentialList.First))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting first %s", err))
		}
	}

	if credentialList.Last != nil {
		err = d.Set("last", dataSourceCredentialListFlattenLast(*credentialList.Last))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting last %s", err))
		}
	}

	if credentialList.Previous != nil {
		err = d.Set("previous", dataSourceCredentialListFlattenPrevious(*credentialList.Previous))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting previous %s", err))
		}
	}

	if credentialList.Credentials != nil {
		err = d.Set("credentials", dataSourceCredentialListFlattenCredentials(credentialList.Credentials))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting credentials %s", err))
		}
	}

	return nil
}

// dataSourceIBMSccPostureListCredentialsID returns a reasonable ID for the list.
func dataSourceIBMSccPostureListCredentialsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceCredentialListFlattenFirst(result posturemanagementv2.PageLink) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceCredentialListFirstToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceCredentialListFirstToMap(firstItem posturemanagementv2.PageLink) (firstMap map[string]interface{}) {
	firstMap = map[string]interface{}{}

	if firstItem.Href != nil {
		firstMap["href"] = firstItem.Href
	}

	return firstMap
}

func dataSourceCredentialListFlattenLast(result posturemanagementv2.PageLink) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceCredentialListLastToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceCredentialListLastToMap(lastItem posturemanagementv2.PageLink) (lastMap map[string]interface{}) {
	lastMap = map[string]interface{}{}

	if lastItem.Href != nil {
		lastMap["href"] = lastItem.Href
	}

	return lastMap
}

func dataSourceCredentialListFlattenPrevious(result posturemanagementv2.PageLink) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceCredentialListPreviousToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceCredentialListPreviousToMap(previousItem posturemanagementv2.PageLink) (previousMap map[string]interface{}) {
	previousMap = map[string]interface{}{}

	if previousItem.Href != nil {
		previousMap["href"] = previousItem.Href
	}

	return previousMap
}

func dataSourceCredentialListFlattenCredentials(result []posturemanagementv2.Credential) (credentials []map[string]interface{}) {
	for _, credentialsItem := range result {
		credentials = append(credentials, dataSourceCredentialListCredentialsToMap(credentialsItem))
	}

	return credentials
}

func dataSourceCredentialListCredentialsToMap(credentialsItem posturemanagementv2.Credential) (credentialsMap map[string]interface{}) {
	credentialsMap = map[string]interface{}{}

	if credentialsItem.Enabled != nil {
		credentialsMap["enabled"] = credentialsItem.Enabled
	}
	if credentialsItem.ID != nil {
		credentialsMap["id"] = credentialsItem.ID
	}
	if credentialsItem.Type != nil {
		credentialsMap["type"] = credentialsItem.Type
	}
	if credentialsItem.Name != nil {
		credentialsMap["name"] = credentialsItem.Name
	}
	if credentialsItem.Description != nil {
		credentialsMap["description"] = credentialsItem.Description
	}
	if credentialsItem.DisplayFields != nil {
		displayFieldsList := []map[string]interface{}{}
		displayFieldsMap := dataSourceCredentialListCredentialsDisplayFieldsToMap(*credentialsItem.DisplayFields)
		displayFieldsList = append(displayFieldsList, displayFieldsMap)
		credentialsMap["display_fields"] = displayFieldsList
	}
	if credentialsItem.CreatedBy != nil {
		credentialsMap["created_by"] = credentialsItem.CreatedBy
	}
	if credentialsItem.CreatedAt != nil {
		credentialsMap["created_at"] = credentialsItem.CreatedAt.String()
	}
	if credentialsItem.UpdatedAt != nil {
		credentialsMap["updated_at"] = credentialsItem.UpdatedAt.String()
	}
	if credentialsItem.UpdatedBy != nil {
		credentialsMap["updated_by"] = credentialsItem.UpdatedBy
	}
	if credentialsItem.Group != nil {
		groupList := []map[string]interface{}{}
		groupMap := dataSourceCredentialListCredentialsGroupToMap(*credentialsItem.Group)
		groupList = append(groupList, groupMap)
		credentialsMap["group"] = groupList
	}
	if credentialsItem.Purpose != nil {
		credentialsMap["purpose"] = credentialsItem.Purpose
	}

	return credentialsMap
}

func dataSourceCredentialListCredentialsDisplayFieldsToMap(displayFieldsItem posturemanagementv2.CredentialDisplayFields) (displayFieldsMap map[string]interface{}) {
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

func dataSourceCredentialListCredentialsGroupToMap(groupItem posturemanagementv2.CredentialGroup) (groupMap map[string]interface{}) {
	groupMap = map[string]interface{}{}

	if groupItem.ID != nil {
		groupMap["id"] = groupItem.ID
	}
	if groupItem.Passphrase != nil {
		groupMap["passphrase"] = groupItem.Passphrase
	}

	return groupMap
}

func dataSourceCredentialListGetNext(next interface{}) int64 {
	if reflect.ValueOf(next).IsNil() {
		return 0
	}

	u, err := url.Parse(reflect.ValueOf(next).Elem().FieldByName("Href").Elem().String())
	if err != nil {
		return 0
	}

	q := u.Query()
	var page string

	if q.Get("start") != "" {
		page = q.Get("start")
	} else if q.Get("offset") != "" {
		page = q.Get("offset")
	}

	convertedVal, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		return 0
	}
	return convertedVal
}
