// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/scc-go-sdk/v3/posturemanagementv2"
)

func DataSourceIBMSccPostureScopes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMSccPostureListScopesRead,

		Schema: map[string]*schema.Schema{
			"offset": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The offset of the page.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of scopes displayed per page.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of scopes. This value is 0 if no scopes are available and below fields will not be available in that case.",
			},
			"first": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The URL of a page.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL of a page.",
						},
					},
				},
			},
			"last": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The URL of a page.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL of a page.",
						},
					},
				},
			},
			"previous": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The URL of a page.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL of a page.",
						},
					},
				},
			},
			"next": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The URL of a page.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL of a page.",
						},
					},
				},
			},
			"scopes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Scopes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A detailed description of the scope.",
						},
						"created_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user who created the scope.",
						},
						"modified_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user who most recently modified the scope.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An auto-generated unique identifier for the scope.",
						},
						"uuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stores the value of scope_uuid .",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A unique name for your scope.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether scope is enabled/disabled.",
						},
						"credential_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The environment that the scope is targeted to.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time that the scope was created in UTC.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time that the scope was last modified in UTC.",
						},
						"collectors": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Stores the value of collectors .Will be displayed only when value exists.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The id of the collector.",
									},
									"display_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-friendly name of the collector.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the collector.",
									},
									"public_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The public key of the collector.Will be used for ssl communciation between collector and orchestrator .This will be populated when collector is installed.",
									},
									"last_heartbeat": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Stores the heartbeat time of a controller . This value exists when collector is installed and running.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The status of collector.",
									},
									"collector_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The collector version. This field is populated when collector is installed.",
									},
									"image_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The image version of the collector. This field is populated when collector is installed. \".",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The description of the collector.",
									},
									"created_by": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The id of the user that created the collector.",
									},
									"created_at": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ISO Date/Time the collector was created.",
									},
									"updated_by": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The id of the user that modified the collector.",
									},
									"updated_at": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ISO Date/Time the collector was modified.",
									},
									"enabled": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Identifies whether the collector is enabled or not(deleted).",
									},
									"registration_code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The registration code of the collector.This is will be used for initial authentication during installation of collector.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of the collector.",
									},
									"credential_public_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The credential public key.",
									},
									"failure_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of times the collector has failed.",
									},
									"approved_local_gateway_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The approved local gateway ip of the collector. This field will be populated only when collector is installed.",
									},
									"approved_internet_gateway_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The approved internet gateway ip of the collector. This field will be populated only when collector is installed.",
									},
									"last_failed_local_gateway_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The failed local gateway ip. This field will be populated only when collector is installed.",
									},
									"reset_reason": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The reason for the collector reset .User resets the collector with a reason for reset. The reason entered by the user is saved in this field .",
									},
									"hostname": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The collector host name. This field will be populated when collector is installed.This will have fully qualified domain name.",
									},
									"install_path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The installation path of the collector. This field will be populated when collector is installed.The value will be folder path.",
									},
									"use_private_endpoint": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether the collector should use a public or private endpoint. This value is generated based on is_public field value during collector creation. If is_public is set to true, this value will be false.",
									},
									"managed_by": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The entity that manages the collector.",
									},
									"trial_expiry": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The trial expiry. This holds the expiry date of registration_code. This field will be populated when collector is installed.",
									},
									"last_failed_internet_gateway_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The failed internet gateway ip of the collector.",
									},
									"status_description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The collector status.",
									},
									"reset_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ISO Date/Time of the collector reset. This value will be populated when a collector is reset. The data-time when the reset event is occured is captured in this field.",
									},
									"is_public": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Determines whether the collector endpoint is accessible on a public network.If set to `true`, the collector connects to resources in your account over a public network. If set to `false`, the collector connects to resources by using a private IP that is accessible only through the IBM Cloud private network.",
									},
									"is_ubi_image": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Determines whether the collector has a Ubi image.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMSccPostureListScopesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	postureManagementClient, err := meta.(conns.ClientSession).PostureManagementV2()
	if err != nil {
		return diag.FromErr(err)
	}

	listScopesOptions := &posturemanagementv2.ListScopesOptions{}
	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting userDetails %s", err))
	}

	accountID := userDetails.UserAccount
	listScopesOptions.SetAccountID(accountID)

	finalList, response, err := postureManagementClient.ListScopesWithContext(context, listScopesOptions)
	if err != nil {
		log.Printf("[DEBUG] ListScopesWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListScopesWithContext failed %s\n%s", err, response))
	}

	scopeList := finalList

	d.SetId(dataSourceIBMSccPostureListScopesID(d))
	if err = d.Set("offset", flex.IntValue(scopeList.Offset)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting offset: %s", err))
	}
	if err = d.Set("limit", flex.IntValue(scopeList.Limit)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting limit: %s", err))
	}
	if err = d.Set("total_count", flex.IntValue(scopeList.TotalCount)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting total_count: %s", err))
	}

	if scopeList.First != nil {
		err = d.Set("first", dataSourceScopeListFlattenFirst(*scopeList.First))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting first %s", err))
		}
	}

	if scopeList.Last != nil {
		err = d.Set("last", dataSourceScopeListFlattenLast(*scopeList.Last))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting last %s", err))
		}
	}

	if scopeList.Previous != nil {
		err = d.Set("previous", dataSourceScopeListFlattenPrevious(*scopeList.Previous))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting previous %s", err))
		}
	}

	if scopeList.Next != nil {
		err = d.Set("next", dataSourceScopeListFlattenNext(*scopeList.Next))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting next %s", err))
		}
	}

	if scopeList.Scopes != nil {
		err = d.Set("scopes", dataSourceScopeListFlattenScopes(scopeList.Scopes))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting scopes %s", err))
		}
	}

	return nil
}

// dataSourceIBMListScopesID returns a reasonable ID for the list.
func dataSourceIBMSccPostureListScopesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceScopeListFlattenFirst(result posturemanagementv2.PageLink) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceScopeListFirstToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceScopeListFirstToMap(firstItem posturemanagementv2.PageLink) (firstMap map[string]interface{}) {
	firstMap = map[string]interface{}{}

	if firstItem.Href != nil {
		firstMap["href"] = firstItem.Href
	}

	return firstMap
}

func dataSourceScopeListFlattenLast(result posturemanagementv2.PageLink) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceScopeListLastToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceScopeListLastToMap(lastItem posturemanagementv2.PageLink) (lastMap map[string]interface{}) {
	lastMap = map[string]interface{}{}

	if lastItem.Href != nil {
		lastMap["href"] = lastItem.Href
	}

	return lastMap
}

func dataSourceScopeListFlattenPrevious(result posturemanagementv2.PageLink) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceScopeListPreviousToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceScopeListPreviousToMap(previousItem posturemanagementv2.PageLink) (previousMap map[string]interface{}) {
	previousMap = map[string]interface{}{}

	if previousItem.Href != nil {
		previousMap["href"] = previousItem.Href
	}

	return previousMap
}

func dataSourceScopeListFlattenNext(result posturemanagementv2.PageLink) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceScopeListNextToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceScopeListNextToMap(nextItem posturemanagementv2.PageLink) (nextMap map[string]interface{}) {
	nextMap = map[string]interface{}{}

	if nextItem.Href != nil {
		nextMap["href"] = nextItem.Href
	}

	return nextMap
}

func dataSourceScopeListFlattenScopes(result []posturemanagementv2.ScopeItem) (scopes []map[string]interface{}) {
	for _, scopesItem := range result {
		scopes = append(scopes, dataSourceScopeListScopesToMap(scopesItem))
	}

	return scopes
}

func dataSourceScopeListScopesToMap(scopesItem posturemanagementv2.ScopeItem) (scopesMap map[string]interface{}) {
	scopesMap = map[string]interface{}{}

	if scopesItem.Description != nil {
		scopesMap["description"] = scopesItem.Description
	}
	if scopesItem.CreatedBy != nil {
		scopesMap["created_by"] = scopesItem.CreatedBy
	}
	if scopesItem.ModifiedBy != nil {
		scopesMap["modified_by"] = scopesItem.ModifiedBy
	}
	if scopesItem.ID != nil {
		scopesMap["id"] = scopesItem.ID
	}
	if scopesItem.UUID != nil {
		scopesMap["uuid"] = scopesItem.UUID
	}
	if scopesItem.Name != nil {
		scopesMap["name"] = scopesItem.Name
	}
	if scopesItem.Enabled != nil {
		scopesMap["enabled"] = scopesItem.Enabled
	}
	if scopesItem.CredentialType != nil {
		scopesMap["credential_type"] = scopesItem.CredentialType
	}
	if scopesItem.CreatedAt != nil {
		scopesMap["created_at"] = scopesItem.CreatedAt.String()
	}
	if scopesItem.UpdatedAt != nil {
		scopesMap["updated_at"] = scopesItem.UpdatedAt.String()
	}
	if scopesItem.Collectors != nil {
		collectorsList := []map[string]interface{}{}
		for _, collectorsItem := range scopesItem.Collectors {
			collectorsList = append(collectorsList, dataSourceScopeListScopesCollectorsToMap(collectorsItem))
		}
		scopesMap["collectors"] = collectorsList
	}

	return scopesMap
}

func dataSourceScopeListScopesCollectorsToMap(collectorsItem posturemanagementv2.Collector) (collectorsMap map[string]interface{}) {
	collectorsMap = map[string]interface{}{}

	if collectorsItem.ID != nil {
		collectorsMap["id"] = collectorsItem.ID
	}
	if collectorsItem.DisplayName != nil {
		collectorsMap["display_name"] = collectorsItem.DisplayName
	}
	if collectorsItem.Name != nil {
		collectorsMap["name"] = collectorsItem.Name
	}
	if collectorsItem.PublicKey != nil {
		collectorsMap["public_key"] = collectorsItem.PublicKey
	}
	if collectorsItem.LastHeartbeat != nil {
		collectorsMap["last_heartbeat"] = collectorsItem.LastHeartbeat.String()
	}
	if collectorsItem.Status != nil {
		collectorsMap["status"] = collectorsItem.Status
	}
	if collectorsItem.CollectorVersion != nil {
		collectorsMap["collector_version"] = collectorsItem.CollectorVersion
	}
	if collectorsItem.ImageVersion != nil {
		collectorsMap["image_version"] = collectorsItem.ImageVersion
	}
	if collectorsItem.Description != nil {
		collectorsMap["description"] = collectorsItem.Description
	}
	if collectorsItem.CreatedBy != nil {
		collectorsMap["created_by"] = collectorsItem.CreatedBy
	}
	if collectorsItem.CreatedAt != nil {
		collectorsMap["created_at"] = collectorsItem.CreatedAt.String()
	}
	if collectorsItem.UpdatedBy != nil {
		collectorsMap["updated_by"] = collectorsItem.UpdatedBy
	}
	if collectorsItem.UpdatedAt != nil {
		collectorsMap["updated_at"] = collectorsItem.UpdatedAt.String()
	}
	if collectorsItem.Enabled != nil {
		collectorsMap["enabled"] = collectorsItem.Enabled
	}
	if collectorsItem.RegistrationCode != nil {
		collectorsMap["registration_code"] = collectorsItem.RegistrationCode
	}
	if collectorsItem.Type != nil {
		collectorsMap["type"] = collectorsItem.Type
	}
	if collectorsItem.CredentialPublicKey != nil {
		collectorsMap["credential_public_key"] = collectorsItem.CredentialPublicKey
	}
	if collectorsItem.FailureCount != nil {
		collectorsMap["failure_count"] = collectorsItem.FailureCount
	}
	if collectorsItem.ApprovedLocalGatewayIP != nil {
		collectorsMap["approved_local_gateway_ip"] = collectorsItem.ApprovedLocalGatewayIP
	}
	if collectorsItem.ApprovedInternetGatewayIP != nil {
		collectorsMap["approved_internet_gateway_ip"] = collectorsItem.ApprovedInternetGatewayIP
	}
	if collectorsItem.LastFailedLocalGatewayIP != nil {
		collectorsMap["last_failed_local_gateway_ip"] = collectorsItem.LastFailedLocalGatewayIP
	}
	if collectorsItem.ResetReason != nil {
		collectorsMap["reset_reason"] = collectorsItem.ResetReason
	}
	if collectorsItem.Hostname != nil {
		collectorsMap["hostname"] = collectorsItem.Hostname
	}
	if collectorsItem.InstallPath != nil {
		collectorsMap["install_path"] = collectorsItem.InstallPath
	}
	if collectorsItem.UsePrivateEndpoint != nil {
		collectorsMap["use_private_endpoint"] = collectorsItem.UsePrivateEndpoint
	}
	if collectorsItem.ManagedBy != nil {
		collectorsMap["managed_by"] = collectorsItem.ManagedBy
	}
	if collectorsItem.TrialExpiry != nil {
		collectorsMap["trial_expiry"] = collectorsItem.TrialExpiry.String()
	}
	if collectorsItem.LastFailedInternetGatewayIP != nil {
		collectorsMap["last_failed_internet_gateway_ip"] = collectorsItem.LastFailedInternetGatewayIP
	}
	if collectorsItem.StatusDescription != nil {
		collectorsMap["status_description"] = collectorsItem.StatusDescription
	}
	if collectorsItem.ResetTime != nil {
		collectorsMap["reset_time"] = collectorsItem.ResetTime.String()
	}
	if collectorsItem.IsPublic != nil {
		collectorsMap["is_public"] = collectorsItem.IsPublic
	}
	if collectorsItem.IsUbiImage != nil {
		collectorsMap["is_ubi_image"] = collectorsItem.IsUbiImage
	}

	return collectorsMap
}
