// Copyright IBM Corp. 2022 All Rights Reserved.
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

func DataSourceIBMSccPostureCollectors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMSccPostureListCollectorsRead,

		Schema: map[string]*schema.Schema{
			"offset": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The offset from the start of the list (0-based).",
			},
			"limit": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of items to return.",
			},
			"total_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of items in the list. This will have value as 0 when no collectors are available and below values will not be populated in that case.",
			},
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
			"next": &schema.Schema{
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
			"collectors": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The array of items returned.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the collector.",
						},
						"display_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-friendly name of the collector.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the collector.",
						},
						"public_key": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The public key of the collector.Will be used for ssl communciation between collector and orchestrator .This will be populated when collector is installed.",
						},
						"last_heartbeat": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stores the heartbeat time of a controller . This value exists when collector is installed and running.",
						},
						"status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of collector.",
						},
						"collector_version": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The collector version. This field is populated when collector is installed.",
						},
						"image_version": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The image version of the collector. This field is populated when collector is installed. \".",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the collector.",
						},
						"created_by": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the user that created the collector.",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ISO Date/Time the collector was created.",
						},
						"updated_by": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the user that modified the collector.",
						},
						"updated_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ISO Date/Time the collector was modified.",
						},
						"enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Identifies whether the collector is enabled or not(deleted).",
						},
						"registration_code": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The registration code of the collector.This is will be used for initial authentication during installation of collector.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the collector.",
						},
						"credential_public_key": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The credential public key.",
						},
						"failure_count": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of times the collector has failed.",
						},
						"approved_local_gateway_ip": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The approved local gateway ip of the collector. This field will be populated only when collector is installed.",
						},
						"approved_internet_gateway_ip": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The approved internet gateway ip of the collector. This field will be populated only when collector is installed.",
						},
						"last_failed_local_gateway_ip": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The failed local gateway ip. This field will be populated only when collector is installed.",
						},
						"reset_reason": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The reason for the collector reset .User resets the collector with a reason for reset. The reason entered by the user is saved in this field .",
						},
						"hostname": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The collector host name. This field will be populated when collector is installed.This will have fully qualified domain name.",
						},
						"install_path": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The installation path of the collector. This field will be populated when collector is installed.The value will be folder path.",
						},
						"use_private_endpoint": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the collector should use a public or private endpoint. This value is generated based on is_public field value during collector creation. If is_public is set to true, this value will be false.",
						},
						"managed_by": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The entity that manages the collector.",
						},
						"trial_expiry": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The trial expiry. This holds the expiry date of registration_code. This field will be populated when collector is installed.",
						},
						"last_failed_internet_gateway_ip": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The failed internet gateway ip of the collector.",
						},
						"status_description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The collector status.",
						},
						"reset_time": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ISO Date/Time of the collector reset. This value will be populated when a collector is reset. The data-time when the reset event is occured is captured in this field.",
						},
						"is_public": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Determines whether the collector endpoint is accessible on a public network.If set to `true`, the collector connects to resources in your account over a public network. If set to `false`, the collector connects to resources by using a private IP that is accessible only through the IBM Cloud private network.",
						},
						"is_ubi_image": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Determines whether the collector has a Ubi image.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMSccPostureListCollectorsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	postureManagementClient, err := meta.(conns.ClientSession).PostureManagementV2()
	if err != nil {
		return diag.FromErr(err)
	}

	listCollectorsOptions := &posturemanagementv2.ListCollectorsOptions{}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting userDetails %s", err))
	}

	accountID := userDetails.UserAccount
	listCollectorsOptions.SetAccountID(accountID)

	collectorList, response, err := postureManagementClient.ListCollectorsWithContext(context, listCollectorsOptions)
	if err != nil {
		log.Printf("[DEBUG] ListCollectorsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListCollectorsWithContext failed %s\n%s", err, response))
	}

	d.SetId(dataSourceIBMSccPostureListCollectorsID(d))
	if err = d.Set("offset", flex.IntValue(collectorList.Offset)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting offset: %s", err))
	}
	if err = d.Set("limit", flex.IntValue(collectorList.Limit)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting limit: %s", err))
	}
	if err = d.Set("total_count", flex.IntValue(collectorList.TotalCount)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting total_count: %s", err))
	}

	if collectorList.First != nil {
		err = d.Set("first", dataSourceCollectorListFlattenFirst(*collectorList.First))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting first %s", err))
		}
	}

	if collectorList.Last != nil {
		err = d.Set("last", dataSourceCollectorListFlattenLast(*collectorList.Last))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting last %s", err))
		}
	}

	if collectorList.Next != nil {
		err = d.Set("next", dataSourceCollectorListFlattenNext(*collectorList.Next))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting next %s", err))
		}
	}

	if collectorList.Previous != nil {
		err = d.Set("previous", dataSourceCollectorListFlattenPrevious(*collectorList.Previous))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting previous %s", err))
		}
	}

	if collectorList.Collectors != nil {
		err = d.Set("collectors", dataSourceCollectorListFlattenCollectors(collectorList.Collectors))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting collectors %s", err))
		}
	}

	return nil
}

// dataSourceIBMSccPostureListCollectorsID returns a reasonable ID for the list.
func dataSourceIBMSccPostureListCollectorsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceCollectorListFlattenFirst(result posturemanagementv2.PageLink) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceCollectorListFirstToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceCollectorListFirstToMap(firstItem posturemanagementv2.PageLink) (firstMap map[string]interface{}) {
	firstMap = map[string]interface{}{}

	if firstItem.Href != nil {
		firstMap["href"] = firstItem.Href
	}

	return firstMap
}

func dataSourceCollectorListFlattenLast(result posturemanagementv2.PageLink) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceCollectorListLastToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceCollectorListLastToMap(lastItem posturemanagementv2.PageLink) (lastMap map[string]interface{}) {
	lastMap = map[string]interface{}{}

	if lastItem.Href != nil {
		lastMap["href"] = lastItem.Href
	}

	return lastMap
}

func dataSourceCollectorListFlattenNext(result posturemanagementv2.PageLink) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceCollectorListNextToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceCollectorListNextToMap(nextItem posturemanagementv2.PageLink) (nextMap map[string]interface{}) {
	nextMap = map[string]interface{}{}

	if nextItem.Href != nil {
		nextMap["href"] = nextItem.Href
	}

	return nextMap
}

func dataSourceCollectorListFlattenPrevious(result posturemanagementv2.PageLink) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceCollectorListPreviousToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceCollectorListPreviousToMap(previousItem posturemanagementv2.PageLink) (previousMap map[string]interface{}) {
	previousMap = map[string]interface{}{}

	if previousItem.Href != nil {
		previousMap["href"] = previousItem.Href
	}

	return previousMap
}

func dataSourceCollectorListFlattenCollectors(result []posturemanagementv2.Collector) (collectors []map[string]interface{}) {
	for _, collectorsItem := range result {
		collectors = append(collectors, dataSourceCollectorListCollectorsToMap(collectorsItem))
	}

	return collectors
}

func dataSourceCollectorListCollectorsToMap(collectorsItem posturemanagementv2.Collector) (collectorsMap map[string]interface{}) {
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
