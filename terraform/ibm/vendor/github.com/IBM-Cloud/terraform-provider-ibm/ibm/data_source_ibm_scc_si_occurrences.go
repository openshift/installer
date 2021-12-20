// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/scc-go-sdk/findingsv1"
)

func dataSourceIBMSccSiOccurrences() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMSccSiOccurrencesRead,

		Schema: map[string]*schema.Schema{
			"account_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"provider_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Part of the parent. This field contains the provider ID. For example: providers/{provider_id}.",
			},
			"page_size": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: ValidatePageSize,
				Description:  "Number of notes to return in the list.",
			},
			"page_token": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Token to provide to skip to a particular spot in the list.",
			},
			"occurrences": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The occurrences requested.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique URL of the resource, image or the container, for which the `Occurrence` applies. For example, https://gcr.io/provider/image@sha256:foo. This field can be used as a filter in list requests.",
						},
						"note_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An analysis note associated with this image, in the form \"{account_id}/providers/{provider_id}/notes/{note_id}\" This field can be used as a filter in list requests.",
						},
						"kind": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of note. Use this field to filter notes and occurences by kind. - FINDING&#58; The note and occurrence represent a finding. - KPI&#58; The note and occurrence represent a KPI value. - CARD&#58; The note represents a card showing findings and related metric values. - CARD_CONFIGURED&#58; The note represents a card configured for a user account. - SECTION&#58; The note represents a section in a dashboard.",
						},
						"remediation": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A description of actions that can be taken to remedy the `Note`.",
						},
						"create_time": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Output only. The time this `Occurrence` was created.",
						},
						"update_time": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Output only. The time this `Occurrence` was last updated.",
						},
						"occurrence_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the occurrence.",
						},
						"context": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IBM Cloud region.",
									},
									"resource_crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource CRN (e.g. certificate CRN, image CRN).",
									},
									"resource_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource ID, in case the CRN is not available.",
									},
									"resource_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-friendly resource name.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type name (e.g. Pod, Cluster, Certificate, Image).",
									},
									"service_crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The service CRN (e.g. CertMgr Instance CRN).",
									},
									"service_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The service name (e.g. CertMgr).",
									},
									"environment_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the environment the occurrence applies to.",
									},
									"component_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the component the occurrence applies to.",
									},
									"toolchain_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The id of the toolchain the occurrence applies to.",
									},
								},
							},
						},
						"finding": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Finding provides details about a finding occurrence.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"severity": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Note provider-assigned severity/impact ranking- LOW&#58; Low Impact- MEDIUM&#58; Medium Impact- HIGH&#58; High Impact- CRITICAL&#58; Critical Impact.",
									},
									"certainty": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Note provider-assigned confidence on the validity of an occurrence- LOW&#58; Low Certainty- MEDIUM&#58; Medium Certainty- HIGH&#58; High Certainty.",
									},
									"next_steps": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Remediation steps for the issues reported in this finding. They override the note's next steps.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"title": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Title of this next step.",
												},
												"url": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The URL associated to this next steps.",
												},
											},
										},
									},
									"network_connection": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "It provides details about a network connection.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"direction": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The direction of this network connection.",
												},
												"protocol": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The protocol of this network connection.",
												},
												"client": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "It provides details about a socket address.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"address": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The IP address of this socket address.",
															},
															"port": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "The port number of this socket address.",
															},
														},
													},
												},
												"server": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "It provides details about a socket address.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"address": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The IP address of this socket address.",
															},
															"port": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "The port number of this socket address.",
															},
														},
													},
												},
											},
										},
									},
									"data_transferred": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "It provides details about data transferred between clients and servers.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"client_bytes": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The number of client bytes transferred.",
												},
												"server_bytes": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The number of server bytes transferred.",
												},
												"client_packets": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The number of client packets transferred.",
												},
												"server_packets": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The number of server packets transferred.",
												},
											},
										},
									},
								},
							},
						},
						"kpi": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Kpi provides details about a KPI occurrence.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": &schema.Schema{
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "The value of this KPI.",
									},
									"total": &schema.Schema{
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "The total value of this KPI.",
									},
								},
							},
						},
					},
				},
			},
			"next_page_token": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The next pagination token in the list response. It should be used as`page_token` for the following request. An empty value means no more results.",
			},
		},
	}
}

func dataSourceIBMSccSiOccurrencesValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 0)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "page_size",
			ValidateFunctionIdentifier: IntBetween,
			Required:                   false,
			MinValue:                   "2"})

	ibmSccSiOccurrencesDataSourceValidator := ResourceValidator{ResourceName: "ibm_scc_si_occurrences", Schema: validateSchema}
	return &ibmSccSiOccurrencesDataSourceValidator
}

func dataSourceIBMSccSiOccurrencesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	findingsClient, err := meta.(ClientSession).FindingsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(err)
	}

	accountID := d.Get("account_id").(string)
	log.Println(fmt.Sprintf("[DEBUG] using specified AccountID %s", accountID))
	if accountID == "" {
		accountID = userDetails.userAccount
		log.Println(fmt.Sprintf("[DEBUG] AccountID not spedified, using %s", accountID))
	}
	findingsClient.AccountID = &accountID

	listOccurrencesOptions := &findingsv1.ListOccurrencesOptions{}

	if pageSize, ok := d.GetOk("page_size"); ok {
		listOccurrencesOptions.SetPageSize(int64(pageSize.(int)))
	}
	if pageToken, ok := d.GetOk("page_token"); ok {
		listOccurrencesOptions.SetPageToken(pageToken.(string))
	}

	listOccurrencesOptions.SetProviderID(d.Get("provider_id").(string))

	apiListOccurrencesResponse, response, err := findingsClient.ListOccurrencesWithContext(context, listOccurrencesOptions)
	if err != nil {
		log.Printf("[DEBUG] ListOccurrencesWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListOccurrencesWithContext failed %s\n%s", err, response))
	}

	d.SetId(dataSourceIBMSccSiOccurrencesID(d))

	if apiListOccurrencesResponse.Occurrences != nil {
		err = d.Set("occurrences", dataSourceAPIListOccurrencesResponseFlattenOccurrences(apiListOccurrencesResponse.Occurrences))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting occurrences %s", err))
		}
	}
	if err = d.Set("next_page_token", apiListOccurrencesResponse.NextPageToken); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting next_page_token: %s", err))
	}

	return nil
}

// dataSourceIBMSccSiOccurrencesID returns a reasonable ID for the list.
func dataSourceIBMSccSiOccurrencesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceAPIListOccurrencesResponseFlattenOccurrences(result []findingsv1.APIOccurrence) (occurrences []map[string]interface{}) {
	for _, occurrencesItem := range result {
		occurrences = append(occurrences, dataSourceAPIListOccurrencesResponseOccurrencesToMap(occurrencesItem))
	}

	return occurrences
}

func dataSourceAPIListOccurrencesResponseOccurrencesToMap(occurrencesItem findingsv1.APIOccurrence) (occurrencesMap map[string]interface{}) {
	occurrencesMap = map[string]interface{}{}

	if occurrencesItem.ID != nil {
		occurrencesMap["occurrence_id"] = occurrencesItem.ID
	}
	if occurrencesItem.ResourceURL != nil {
		occurrencesMap["resource_url"] = occurrencesItem.ResourceURL
	}
	if occurrencesItem.NoteName != nil {
		occurrencesMap["note_name"] = occurrencesItem.NoteName
	}
	if occurrencesItem.Kind != nil {
		occurrencesMap["kind"] = occurrencesItem.Kind
	}
	if occurrencesItem.Remediation != nil {
		occurrencesMap["remediation"] = occurrencesItem.Remediation
	}
	if occurrencesItem.CreateTime != nil {
		occurrencesMap["create_time"] = occurrencesItem.CreateTime.String()
	}
	if occurrencesItem.UpdateTime != nil {
		occurrencesMap["update_time"] = occurrencesItem.UpdateTime.String()
	}
	if occurrencesItem.ID != nil {
		occurrencesMap["occurrence_id"] = occurrencesItem.ID
	}
	if occurrencesItem.Context != nil {
		contextList := []map[string]interface{}{}
		contextMap := dataSourceAPIListOccurrencesResponseOccurrencesContextToMap(*occurrencesItem.Context)
		contextList = append(contextList, contextMap)
		occurrencesMap["context"] = contextList
	}
	if occurrencesItem.Finding != nil {
		findingList := []map[string]interface{}{}
		findingMap := dataSourceAPIListOccurrencesResponseOccurrencesFindingToMap(*occurrencesItem.Finding)
		findingList = append(findingList, findingMap)
		occurrencesMap["finding"] = findingList
	}
	if occurrencesItem.Kpi != nil {
		kpiList := []map[string]interface{}{}
		kpiMap := dataSourceAPIListOccurrencesResponseOccurrencesKpiToMap(*occurrencesItem.Kpi)
		kpiList = append(kpiList, kpiMap)
		occurrencesMap["kpi"] = kpiList
	}

	return occurrencesMap
}

func dataSourceAPIListOccurrencesResponseOccurrencesContextToMap(contextItem findingsv1.Context) (contextMap map[string]interface{}) {
	contextMap = map[string]interface{}{}

	if contextItem.Region != nil {
		contextMap["region"] = contextItem.Region
	}
	if contextItem.ResourceCRN != nil {
		contextMap["resource_crn"] = contextItem.ResourceCRN
	}
	if contextItem.ResourceID != nil {
		contextMap["resource_id"] = contextItem.ResourceID
	}
	if contextItem.ResourceName != nil {
		contextMap["resource_name"] = contextItem.ResourceName
	}
	if contextItem.ResourceType != nil {
		contextMap["resource_type"] = contextItem.ResourceType
	}
	if contextItem.ServiceCRN != nil {
		contextMap["service_crn"] = contextItem.ServiceCRN
	}
	if contextItem.ServiceName != nil {
		contextMap["service_name"] = contextItem.ServiceName
	}
	if contextItem.EnvironmentName != nil {
		contextMap["environment_name"] = contextItem.EnvironmentName
	}
	if contextItem.ComponentName != nil {
		contextMap["component_name"] = contextItem.ComponentName
	}
	if contextItem.ToolchainID != nil {
		contextMap["toolchain_id"] = contextItem.ToolchainID
	}

	return contextMap
}

func dataSourceAPIListOccurrencesResponseOccurrencesFindingToMap(findingItem findingsv1.Finding) (findingMap map[string]interface{}) {
	findingMap = map[string]interface{}{}

	if findingItem.Severity != nil {
		findingMap["severity"] = findingItem.Severity
	}
	if findingItem.Certainty != nil {
		findingMap["certainty"] = findingItem.Certainty
	}
	if findingItem.NextSteps != nil {
		nextStepsList := []map[string]interface{}{}
		for _, nextStepsItem := range findingItem.NextSteps {
			nextStepsList = append(nextStepsList, dataSourceAPIListOccurrencesResponseFindingNextStepsToMap(nextStepsItem))
		}
		findingMap["next_steps"] = nextStepsList
	}
	if findingItem.NetworkConnection != nil {
		networkConnectionList := []map[string]interface{}{}
		networkConnectionMap := dataSourceAPIListOccurrencesResponseFindingNetworkConnectionToMap(*findingItem.NetworkConnection)
		networkConnectionList = append(networkConnectionList, networkConnectionMap)
		findingMap["network_connection"] = networkConnectionList
	}
	if findingItem.DataTransferred != nil {
		dataTransferredList := []map[string]interface{}{}
		dataTransferredMap := dataSourceAPIListOccurrencesResponseFindingDataTransferredToMap(*findingItem.DataTransferred)
		dataTransferredList = append(dataTransferredList, dataTransferredMap)
		findingMap["data_transferred"] = dataTransferredList
	}

	return findingMap
}

func dataSourceAPIListOccurrencesResponseFindingNextStepsToMap(nextStepsItem findingsv1.RemediationStep) (nextStepsMap map[string]interface{}) {
	nextStepsMap = map[string]interface{}{}

	if nextStepsItem.Title != nil {
		nextStepsMap["title"] = nextStepsItem.Title
	}
	if nextStepsItem.URL != nil {
		nextStepsMap["url"] = nextStepsItem.URL
	}

	return nextStepsMap
}

func dataSourceAPIListOccurrencesResponseFindingNetworkConnectionToMap(networkConnectionItem findingsv1.NetworkConnection) (networkConnectionMap map[string]interface{}) {
	networkConnectionMap = map[string]interface{}{}

	if networkConnectionItem.Direction != nil {
		networkConnectionMap["direction"] = networkConnectionItem.Direction
	}
	if networkConnectionItem.Protocol != nil {
		networkConnectionMap["protocol"] = networkConnectionItem.Protocol
	}
	if networkConnectionItem.Client != nil {
		clientList := []map[string]interface{}{}
		clientMap := dataSourceAPIListOccurrencesResponseNetworkConnectionClientToMap(*networkConnectionItem.Client)
		clientList = append(clientList, clientMap)
		networkConnectionMap["client"] = clientList
	}
	if networkConnectionItem.Server != nil {
		serverList := []map[string]interface{}{}
		serverMap := dataSourceAPIListOccurrencesResponseNetworkConnectionServerToMap(*networkConnectionItem.Server)
		serverList = append(serverList, serverMap)
		networkConnectionMap["server"] = serverList
	}

	return networkConnectionMap
}

func dataSourceAPIListOccurrencesResponseNetworkConnectionClientToMap(clientItem findingsv1.SocketAddress) (clientMap map[string]interface{}) {
	clientMap = map[string]interface{}{}

	if clientItem.Address != nil {
		clientMap["address"] = clientItem.Address
	}
	if clientItem.Port != nil {
		clientMap["port"] = clientItem.Port
	}

	return clientMap
}

func dataSourceAPIListOccurrencesResponseNetworkConnectionServerToMap(serverItem findingsv1.SocketAddress) (serverMap map[string]interface{}) {
	serverMap = map[string]interface{}{}

	if serverItem.Address != nil {
		serverMap["address"] = serverItem.Address
	}
	if serverItem.Port != nil {
		serverMap["port"] = serverItem.Port
	}

	return serverMap
}

func dataSourceAPIListOccurrencesResponseFindingDataTransferredToMap(dataTransferredItem findingsv1.DataTransferred) (dataTransferredMap map[string]interface{}) {
	dataTransferredMap = map[string]interface{}{}

	if dataTransferredItem.ClientBytes != nil {
		dataTransferredMap["client_bytes"] = dataTransferredItem.ClientBytes
	}
	if dataTransferredItem.ServerBytes != nil {
		dataTransferredMap["server_bytes"] = dataTransferredItem.ServerBytes
	}
	if dataTransferredItem.ClientPackets != nil {
		dataTransferredMap["client_packets"] = dataTransferredItem.ClientPackets
	}
	if dataTransferredItem.ServerPackets != nil {
		dataTransferredMap["server_packets"] = dataTransferredItem.ServerPackets
	}

	return dataTransferredMap
}

func dataSourceAPIListOccurrencesResponseOccurrencesKpiToMap(kpiItem findingsv1.Kpi) (kpiMap map[string]interface{}) {
	kpiMap = map[string]interface{}{}

	if kpiItem.Value != nil {
		kpiMap["value"] = kpiItem.Value
	}
	if kpiItem.Total != nil {
		kpiMap["total"] = kpiItem.Total
	}

	return kpiMap
}
