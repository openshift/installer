// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/scc-go-sdk/findingsv1"
)

func dataSourceIBMSccSiOccurrence() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMSccSiOccurrenceRead,

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
			"occurrence_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Second part of occurrence `name`: providers/{provider_id}/occurrences/{occurrence_id}.",
			},
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
	}
}

func dataSourceIBMSccSiOccurrenceRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	getOccurrenceOptions := &findingsv1.GetOccurrenceOptions{}

	getOccurrenceOptions.SetProviderID(d.Get("provider_id").(string))
	getOccurrenceOptions.SetOccurrenceID(d.Get("occurrence_id").(string))

	apiOccurrence, response, err := findingsClient.GetOccurrenceWithContext(context, getOccurrenceOptions)
	if err != nil {
		log.Printf("[DEBUG] GetOccurrenceWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetOccurrenceWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", *findingsClient.AccountID, *getOccurrenceOptions.ProviderID, *getOccurrenceOptions.OccurrenceID))
	if err = d.Set("resource_url", apiOccurrence.ResourceURL); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_url: %s", err))
	}
	if err = d.Set("note_name", apiOccurrence.NoteName); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting note_name: %s", err))
	}
	if err = d.Set("kind", apiOccurrence.Kind); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting kind: %s", err))
	}
	if err = d.Set("remediation", apiOccurrence.Remediation); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting remediation: %s", err))
	}
	if err = d.Set("create_time", dateTimeToString(apiOccurrence.CreateTime)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting create_time: %s", err))
	}
	if err = d.Set("update_time", dateTimeToString(apiOccurrence.UpdateTime)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting update_time: %s", err))
	}

	if apiOccurrence.Context != nil {
		err = d.Set("context", dataSourceAPIOccurrenceFlattenContext(*apiOccurrence.Context))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting context %s", err))
		}
	}

	if apiOccurrence.Finding != nil {
		err = d.Set("finding", dataSourceAPIOccurrenceFlattenFinding(*apiOccurrence.Finding))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting finding %s", err))
		}
	}

	if apiOccurrence.Kpi != nil {
		err = d.Set("kpi", dataSourceAPIOccurrenceFlattenKpi(*apiOccurrence.Kpi))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting kpi %s", err))
		}
	}

	return nil
}

func dataSourceAPIOccurrenceFlattenContext(result findingsv1.Context) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceAPIOccurrenceContextToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceAPIOccurrenceContextToMap(contextItem findingsv1.Context) (contextMap map[string]interface{}) {
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

func dataSourceAPIOccurrenceFlattenFinding(result findingsv1.Finding) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceAPIOccurrenceFindingToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceAPIOccurrenceFindingToMap(findingItem findingsv1.Finding) (findingMap map[string]interface{}) {
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
			nextStepsList = append(nextStepsList, dataSourceAPIOccurrenceFindingNextStepsToMap(nextStepsItem))
		}
		findingMap["next_steps"] = nextStepsList
	}
	if findingItem.NetworkConnection != nil {
		networkConnectionList := []map[string]interface{}{}
		networkConnectionMap := dataSourceAPIOccurrenceFindingNetworkConnectionToMap(*findingItem.NetworkConnection)
		networkConnectionList = append(networkConnectionList, networkConnectionMap)
		findingMap["network_connection"] = networkConnectionList
	}
	if findingItem.DataTransferred != nil {
		dataTransferredList := []map[string]interface{}{}
		dataTransferredMap := dataSourceAPIOccurrenceFindingDataTransferredToMap(*findingItem.DataTransferred)
		dataTransferredList = append(dataTransferredList, dataTransferredMap)
		findingMap["data_transferred"] = dataTransferredList
	}

	return findingMap
}

func dataSourceAPIOccurrenceFindingNextStepsToMap(nextStepsItem findingsv1.RemediationStep) (nextStepsMap map[string]interface{}) {
	nextStepsMap = map[string]interface{}{}

	if nextStepsItem.Title != nil {
		nextStepsMap["title"] = nextStepsItem.Title
	}
	if nextStepsItem.URL != nil {
		nextStepsMap["url"] = nextStepsItem.URL
	}

	return nextStepsMap
}

func dataSourceAPIOccurrenceFindingNetworkConnectionToMap(networkConnectionItem findingsv1.NetworkConnection) (networkConnectionMap map[string]interface{}) {
	networkConnectionMap = map[string]interface{}{}

	if networkConnectionItem.Direction != nil {
		networkConnectionMap["direction"] = networkConnectionItem.Direction
	}
	if networkConnectionItem.Protocol != nil {
		networkConnectionMap["protocol"] = networkConnectionItem.Protocol
	}
	if networkConnectionItem.Client != nil {
		clientList := []map[string]interface{}{}
		clientMap := dataSourceAPIOccurrenceNetworkConnectionClientToMap(*networkConnectionItem.Client)
		clientList = append(clientList, clientMap)
		networkConnectionMap["client"] = clientList
	}
	if networkConnectionItem.Server != nil {
		serverList := []map[string]interface{}{}
		serverMap := dataSourceAPIOccurrenceNetworkConnectionServerToMap(*networkConnectionItem.Server)
		serverList = append(serverList, serverMap)
		networkConnectionMap["server"] = serverList
	}

	return networkConnectionMap
}

func dataSourceAPIOccurrenceNetworkConnectionClientToMap(clientItem findingsv1.SocketAddress) (clientMap map[string]interface{}) {
	clientMap = map[string]interface{}{}

	if clientItem.Address != nil {
		clientMap["address"] = clientItem.Address
	}
	if clientItem.Port != nil {
		clientMap["port"] = clientItem.Port
	}

	return clientMap
}

func dataSourceAPIOccurrenceNetworkConnectionServerToMap(serverItem findingsv1.SocketAddress) (serverMap map[string]interface{}) {
	serverMap = map[string]interface{}{}

	if serverItem.Address != nil {
		serverMap["address"] = serverItem.Address
	}
	if serverItem.Port != nil {
		serverMap["port"] = serverItem.Port
	}

	return serverMap
}

func dataSourceAPIOccurrenceFindingDataTransferredToMap(dataTransferredItem findingsv1.DataTransferred) (dataTransferredMap map[string]interface{}) {
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

func dataSourceAPIOccurrenceFlattenKpi(result findingsv1.Kpi) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceAPIOccurrenceKpiToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceAPIOccurrenceKpiToMap(kpiItem findingsv1.Kpi) (kpiMap map[string]interface{}) {
	kpiMap = map[string]interface{}{}

	if kpiItem.Value != nil {
		kpiMap["value"] = kpiItem.Value
	}
	if kpiItem.Total != nil {
		kpiMap["total"] = kpiItem.Total
	}

	return kpiMap
}
