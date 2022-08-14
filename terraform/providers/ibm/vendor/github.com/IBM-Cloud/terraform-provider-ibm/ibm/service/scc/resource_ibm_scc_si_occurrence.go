// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/scc-go-sdk/findingsv1"
)

func ResourceIBMSccSiOccurrence() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMSccSiOccurrenceCreate,
		ReadContext:   resourceIBMSccSiOccurrenceRead,
		UpdateContext: resourceIBMSccSiOccurrenceUpdate,
		DeleteContext: resourceIBMSccSiOccurrenceDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"provider_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Part of the parent. This field contains the provider ID. For example: providers/{provider_id}.",
			},
			"note_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "An analysis note associated with this image, in the form \"{account_id}/providers/{provider_id}/notes/{note_id}\" This field can be used as a filter in list requests.",
			},
			"kind": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_scc_si_occurrence", "kind"),
				Description:  "The type of note. Use this field to filter notes and occurences by kind. - FINDING&#58; The note and occurrence represent a finding. - KPI&#58; The note and occurrence represent a KPI value. - CARD&#58; The note represents a card showing findings and related metric values. - CARD_CONFIGURED&#58; The note represents a card configured for a user account. - SECTION&#58; The note represents a section in a dashboard.",
			},
			"occurrence_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the occurrence.",
			},
			"resource_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique URL of the resource, image or the container, for which the `Occurrence` applies. For example, https://gcr.io/provider/image@sha256:foo. This field can be used as a filter in list requests.",
			},
			"remediation": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A description of actions that can be taken to remedy the `Note`.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The time this `Occurrence` was created.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The time this `Occurrence` was last updated.",
			},
			"context": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The IBM Cloud region.",
						},
						"resource_crn": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The resource CRN (e.g. certificate CRN, image CRN).",
						},
						"resource_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The resource ID, in case the CRN is not available.",
						},
						"resource_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The user-friendly resource name.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The resource type name (e.g. Pod, Cluster, Certificate, Image).",
						},
						"service_crn": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The service CRN (e.g. CertMgr Instance CRN).",
						},
						"service_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The service name (e.g. CertMgr).",
						},
						"environment_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the environment the occurrence applies to.",
						},
						"component_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the component the occurrence applies to.",
						},
						"toolchain_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The id of the toolchain the occurrence applies to.",
						},
					},
				},
			},
			"finding": {
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				Description:  "Finding provides details about a finding occurrence.",
				ExactlyOneOf: []string{"finding", "kpi"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"severity": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Note provider-assigned severity/impact ranking- LOW&#58; Low Impact- MEDIUM&#58; Medium Impact- HIGH&#58; High Impact- CRITICAL&#58; Critical Impact.",
						},
						"certainty": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Note provider-assigned confidence on the validity of an occurrence- LOW&#58; Low Certainty- MEDIUM&#58; Medium Certainty- HIGH&#58; High Certainty.",
						},
						"next_steps": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Remediation steps for the issues reported in this finding. They override the note's next steps.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"title": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Title of this next step.",
									},
									"url": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The URL associated to this next steps.",
									},
								},
							},
						},
						"network_connection": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "It provides details about a network connection.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"direction": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The direction of this network connection.",
									},
									"protocol": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The protocol of this network connection.",
									},
									"client": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "It provides details about a socket address.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The IP address of this socket address.",
												},
												"port": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "The port number of this socket address.",
												},
											},
										},
									},
									"server": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "It provides details about a socket address.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The IP address of this socket address.",
												},
												"port": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "The port number of this socket address.",
												},
											},
										},
									},
								},
							},
						},
						"data_transferred": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "It provides details about data transferred between clients and servers.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"client_bytes": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The number of client bytes transferred.",
									},
									"server_bytes": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The number of server bytes transferred.",
									},
									"client_packets": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The number of client packets transferred.",
									},
									"server_packets": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The number of server packets transferred.",
									},
								},
							},
						},
					},
				},
			},
			"kpi": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Kpi provides details about a KPI occurrence.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:        schema.TypeFloat,
							Required:    true,
							Description: "The value of this KPI.",
						},
						"total": {
							Type:        schema.TypeFloat,
							Optional:    true,
							Description: "The total value of this KPI.",
						},
					},
				},
			},
			"replace_if_exists": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When set to true, an existing occurrence is replaced rather than duplicated.",
			},
		},
	}
}

func ResourceIBMSccSiOccurrenceValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "kind",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "FINDING, KPI",
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_scc_si_occurrence", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMSccSiOccurrenceCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	findingsClient, err := meta.(conns.ClientSession).FindingsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(err)
	}

	accountID := d.Get("account_id").(string)
	log.Println(fmt.Sprintf("[DEBUG] using specified AccountID %s", accountID))
	if accountID == "" {
		accountID = userDetails.UserAccount
		log.Println(fmt.Sprintf("[DEBUG] AccountID not spedified, using %s", accountID))
	}
	findingsClient.AccountID = &accountID

	createOccurrenceOptions := &findingsv1.CreateOccurrenceOptions{}

	createOccurrenceOptions.SetProviderID(d.Get("provider_id").(string))
	createOccurrenceOptions.SetNoteName(d.Get("note_name").(string))
	createOccurrenceOptions.SetKind(d.Get("kind").(string))
	createOccurrenceOptions.SetID(d.Get("occurrence_id").(string))
	if _, ok := d.GetOk("resource_url"); ok {
		createOccurrenceOptions.SetResourceURL(d.Get("resource_url").(string))
	}
	if _, ok := d.GetOk("remediation"); ok {
		createOccurrenceOptions.SetRemediation(d.Get("remediation").(string))
	}
	if _, ok := d.GetOk("context"); ok {
		context := resourceIBMSccSiOccurrenceMapToContext(d.Get("context.0").(map[string]interface{}))
		createOccurrenceOptions.SetContext(&context)
	}
	if _, ok := d.GetOk("finding"); ok {
		finding := resourceIBMSccSiOccurrenceMapToFinding(d.Get("finding.0").(map[string]interface{}))
		createOccurrenceOptions.SetFinding(&finding)
	}
	if _, ok := d.GetOk("kpi"); ok {
		kpi := resourceIBMSccSiOccurrenceMapToKpi(d.Get("kpi.0").(map[string]interface{}))
		createOccurrenceOptions.SetKpi(&kpi)
	}
	if _, ok := d.GetOk("replace_if_exists"); ok {
		createOccurrenceOptions.SetReplaceIfExists(d.Get("replace_if_exists").(bool))
	}

	apiOccurrence, response, err := findingsClient.CreateOccurrenceWithContext(context, createOccurrenceOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateOccurrenceWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateOccurrenceWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", *findingsClient.AccountID, *createOccurrenceOptions.ProviderID, *apiOccurrence.ID))

	return resourceIBMSccSiOccurrenceRead(context, d, meta)
}

func resourceIBMSccSiOccurrenceMapToContext(contextMap map[string]interface{}) findingsv1.Context {
	context := findingsv1.Context{}

	if contextMap["region"] != nil {
		context.Region = core.StringPtr(contextMap["region"].(string))
	}
	if contextMap["resource_crn"] != nil {
		context.ResourceCRN = core.StringPtr(contextMap["resource_crn"].(string))
	}
	if contextMap["resource_id"] != nil {
		context.ResourceID = core.StringPtr(contextMap["resource_id"].(string))
	}
	if contextMap["resource_name"] != nil {
		context.ResourceName = core.StringPtr(contextMap["resource_name"].(string))
	}
	if contextMap["resource_type"] != nil {
		context.ResourceType = core.StringPtr(contextMap["resource_type"].(string))
	}
	if contextMap["service_crn"] != nil {
		context.ServiceCRN = core.StringPtr(contextMap["service_crn"].(string))
	}
	if contextMap["service_name"] != nil {
		context.ServiceName = core.StringPtr(contextMap["service_name"].(string))
	}
	if contextMap["environment_name"] != nil {
		context.EnvironmentName = core.StringPtr(contextMap["environment_name"].(string))
	}
	if contextMap["component_name"] != nil {
		context.ComponentName = core.StringPtr(contextMap["component_name"].(string))
	}
	if contextMap["toolchain_id"] != nil {
		context.ToolchainID = core.StringPtr(contextMap["toolchain_id"].(string))
	}

	return context
}

func resourceIBMSccSiOccurrenceMapToFinding(findingMap map[string]interface{}) findingsv1.Finding {
	finding := findingsv1.Finding{}

	if findingMap["severity"] != nil {
		finding.Severity = core.StringPtr(findingMap["severity"].(string))
	}
	if findingMap["certainty"] != nil {
		finding.Certainty = core.StringPtr(findingMap["certainty"].(string))
	}
	if findingMap["next_steps"] != nil {
		nextSteps := []findingsv1.RemediationStep{}
		for _, nextStepsItem := range findingMap["next_steps"].([]interface{}) {
			nextStepsItemModel := resourceIBMSccSiOccurrenceMapToRemediationStep(nextStepsItem.(map[string]interface{}))
			nextSteps = append(nextSteps, nextStepsItemModel)
		}
		finding.NextSteps = nextSteps
	}
	if findingMap["network_connection"] != nil {
		if len(findingMap["network_connection"].([]interface{})) > 0 {
			networkConnection := resourceIBMSccSiOccurrenceMapToNetworkConnection(findingMap["network_connection"].([]interface{})[0].(map[string]interface{}))
			finding.NetworkConnection = &networkConnection
		}
	}
	if findingMap["data_transferred"] != nil {
		if len(findingMap["data_transferred"].([]interface{})) > 0 {
			dataTransferred := resourceIBMSccSiOccurrenceMapToDataTransferred(findingMap["data_transferred"].([]interface{})[0].(map[string]interface{}))
			finding.DataTransferred = &dataTransferred
		}
	}

	return finding
}

func resourceIBMSccSiOccurrenceMapToRemediationStep(remediationStepMap map[string]interface{}) findingsv1.RemediationStep {
	remediationStep := findingsv1.RemediationStep{}

	if remediationStepMap["title"] != nil {
		remediationStep.Title = core.StringPtr(remediationStepMap["title"].(string))
	}
	if remediationStepMap["url"] != nil {
		remediationStep.URL = core.StringPtr(remediationStepMap["url"].(string))
	}

	return remediationStep
}

func resourceIBMSccSiOccurrenceMapToNetworkConnection(networkConnectionMap map[string]interface{}) findingsv1.NetworkConnection {
	networkConnection := findingsv1.NetworkConnection{}

	if networkConnectionMap["direction"] != nil {
		networkConnection.Direction = core.StringPtr(networkConnectionMap["direction"].(string))
	}
	if networkConnectionMap["protocol"] != nil {
		networkConnection.Protocol = core.StringPtr(networkConnectionMap["protocol"].(string))
	}
	if networkConnectionMap["client"] != nil {
		client := resourceIBMSccSiOccurrenceMapToSocketAddress(networkConnectionMap["client"].([]interface{})[0].(map[string]interface{}))
		networkConnection.Client = &client
	}
	if networkConnectionMap["server"] != nil {
		server := resourceIBMSccSiOccurrenceMapToSocketAddress(networkConnectionMap["server"].([]interface{})[0].(map[string]interface{}))
		networkConnection.Server = &server
	}

	return networkConnection
}

func resourceIBMSccSiOccurrenceMapToSocketAddress(socketAddressMap map[string]interface{}) findingsv1.SocketAddress {
	socketAddress := findingsv1.SocketAddress{}

	socketAddress.Address = core.StringPtr(socketAddressMap["address"].(string))
	if socketAddressMap["port"] != nil {
		socketAddress.Port = core.Int64Ptr(int64(socketAddressMap["port"].(int)))
	}

	return socketAddress
}

func resourceIBMSccSiOccurrenceMapToDataTransferred(dataTransferredMap map[string]interface{}) findingsv1.DataTransferred {
	dataTransferred := findingsv1.DataTransferred{}

	if dataTransferredMap["client_bytes"] != nil {
		dataTransferred.ClientBytes = core.Int64Ptr(int64(dataTransferredMap["client_bytes"].(int)))
	}
	if dataTransferredMap["server_bytes"] != nil {
		dataTransferred.ServerBytes = core.Int64Ptr(int64(dataTransferredMap["server_bytes"].(int)))
	}
	if dataTransferredMap["client_packets"] != nil {
		dataTransferred.ClientPackets = core.Int64Ptr(int64(dataTransferredMap["client_packets"].(int)))
	}
	if dataTransferredMap["server_packets"] != nil {
		dataTransferred.ServerPackets = core.Int64Ptr(int64(dataTransferredMap["server_packets"].(int)))
	}

	return dataTransferred
}

func resourceIBMSccSiOccurrenceMapToKpi(kpiMap map[string]interface{}) findingsv1.Kpi {
	kpi := findingsv1.Kpi{}

	kpi.Value = core.Float64Ptr(kpiMap["value"].(float64))
	if kpiMap["total"] != nil {
		kpi.Total = core.Float64Ptr(kpiMap["total"].(float64))
	}

	return kpi
}

func resourceIBMSccSiOccurrenceRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	findingsClient, err := meta.(conns.ClientSession).FindingsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getOccurrenceOptions := &findingsv1.GetOccurrenceOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	findingsClient.AccountID = &parts[0]

	d.Set("account_id", &parts[0])

	getOccurrenceOptions.SetProviderID(parts[1])
	getOccurrenceOptions.SetOccurrenceID(parts[2])

	apiOccurrence, response, err := findingsClient.GetOccurrenceWithContext(context, getOccurrenceOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetOccurrenceWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetOccurrenceWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("provider_id", getOccurrenceOptions.ProviderID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting provider_id: %s", err))
	}
	// TODO: handle argument of type bool
	if err = d.Set("note_name", apiOccurrence.NoteName); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting note_name: %s", err))
	}
	if err = d.Set("kind", apiOccurrence.Kind); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting kind: %s", err))
	}
	if err = d.Set("occurrence_id", apiOccurrence.ID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting occurrence_id: %s", err))
	}
	if err = d.Set("resource_url", apiOccurrence.ResourceURL); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_url: %s", err))
	}
	if err = d.Set("remediation", apiOccurrence.Remediation); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting remediation: %s", err))
	}
	if err = d.Set("create_time", flex.DateTimeToString(apiOccurrence.CreateTime)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting create_time: %s", err))
	}
	if err = d.Set("update_time", flex.DateTimeToString(apiOccurrence.UpdateTime)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting update_time: %s", err))
	}
	if apiOccurrence.Context != nil {
		contextMap := resourceIBMSccSiOccurrenceContextToMap(*apiOccurrence.Context)
		if len(contextMap) > 0 {
			if err = d.Set("context", []map[string]interface{}{contextMap}); err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error setting context: %s", err))
			}
		}
	}
	if apiOccurrence.Finding != nil {
		findingMap := resourceIBMSccSiOccurrenceFindingToMap(*apiOccurrence.Finding)
		if err = d.Set("finding", []map[string]interface{}{findingMap}); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting finding: %s", err))
		}
	}
	if apiOccurrence.Kpi != nil {
		kpiMap := resourceIBMSccSiOccurrenceKpiToMap(*apiOccurrence.Kpi)
		if err = d.Set("kpi", []map[string]interface{}{kpiMap}); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting kpi: %s", err))
		}
	}
	if _, ok := d.GetOk("replace_if_exists"); ok {
		replaceIfExists := d.Get("replace_if_exists").(bool)
		d.Set("replace_if_exists", replaceIfExists)
	}

	return nil
}

func resourceIBMSccSiOccurrenceContextToMap(context findingsv1.Context) map[string]interface{} {
	contextMap := map[string]interface{}{}

	if context.Region != nil {
		contextMap["region"] = context.Region
	}
	if context.ResourceCRN != nil {
		contextMap["resource_crn"] = context.ResourceCRN
	}
	if context.ResourceID != nil {
		contextMap["resource_id"] = context.ResourceID
	}
	if context.ResourceName != nil {
		contextMap["resource_name"] = context.ResourceName
	}
	if context.ResourceType != nil {
		contextMap["resource_type"] = context.ResourceType
	}
	if context.ServiceCRN != nil {
		contextMap["service_crn"] = context.ServiceCRN
	}
	if context.ServiceName != nil {
		contextMap["service_name"] = context.ServiceName
	}
	if context.EnvironmentName != nil {
		contextMap["environment_name"] = context.EnvironmentName
	}
	if context.ComponentName != nil {
		contextMap["component_name"] = context.ComponentName
	}
	if context.ToolchainID != nil {
		contextMap["toolchain_id"] = context.ToolchainID
	}

	return contextMap
}

func resourceIBMSccSiOccurrenceFindingToMap(finding findingsv1.Finding) map[string]interface{} {
	findingMap := map[string]interface{}{}

	if finding.Severity != nil {
		findingMap["severity"] = finding.Severity
	}
	if finding.Certainty != nil {
		findingMap["certainty"] = finding.Certainty
	}
	if finding.NextSteps != nil {
		nextSteps := []map[string]interface{}{}
		for _, nextStepsItem := range finding.NextSteps {
			nextStepsItemMap := resourceIBMSccSiOccurrenceRemediationStepToMap(nextStepsItem)
			nextSteps = append(nextSteps, nextStepsItemMap)
			// TODO: handle NextSteps of type TypeList -- list of non-primitive, not model items
		}
		findingMap["next_steps"] = nextSteps
	}
	if finding.NetworkConnection != nil {
		NetworkConnectionMap := resourceIBMSccSiOccurrenceNetworkConnectionToMap(*finding.NetworkConnection)
		findingMap["network_connection"] = []map[string]interface{}{NetworkConnectionMap}
	}
	if finding.DataTransferred != nil {
		DataTransferredMap := resourceIBMSccSiOccurrenceDataTransferredToMap(*finding.DataTransferred)
		findingMap["data_transferred"] = []map[string]interface{}{DataTransferredMap}
	}

	return findingMap
}

func resourceIBMSccSiOccurrenceRemediationStepToMap(remediationStep findingsv1.RemediationStep) map[string]interface{} {
	remediationStepMap := map[string]interface{}{}

	if remediationStep.Title != nil {
		remediationStepMap["title"] = remediationStep.Title
	}
	if remediationStep.URL != nil {
		remediationStepMap["url"] = remediationStep.URL
	}

	return remediationStepMap
}

func resourceIBMSccSiOccurrenceNetworkConnectionToMap(networkConnection findingsv1.NetworkConnection) map[string]interface{} {
	networkConnectionMap := map[string]interface{}{}

	if networkConnection.Direction != nil {
		networkConnectionMap["direction"] = networkConnection.Direction
	}
	if networkConnection.Protocol != nil {
		networkConnectionMap["protocol"] = networkConnection.Protocol
	}
	if networkConnection.Client != nil {
		ClientMap := resourceIBMSccSiOccurrenceSocketAddressToMap(*networkConnection.Client)
		networkConnectionMap["client"] = []map[string]interface{}{ClientMap}
	}
	if networkConnection.Server != nil {
		ServerMap := resourceIBMSccSiOccurrenceSocketAddressToMap(*networkConnection.Server)
		networkConnectionMap["server"] = []map[string]interface{}{ServerMap}
	}

	return networkConnectionMap
}

func resourceIBMSccSiOccurrenceSocketAddressToMap(socketAddress findingsv1.SocketAddress) map[string]interface{} {
	socketAddressMap := map[string]interface{}{}

	socketAddressMap["address"] = socketAddress.Address
	if socketAddress.Port != nil {
		socketAddressMap["port"] = flex.IntValue(socketAddress.Port)
	}

	return socketAddressMap
}

func resourceIBMSccSiOccurrenceDataTransferredToMap(dataTransferred findingsv1.DataTransferred) map[string]interface{} {
	dataTransferredMap := map[string]interface{}{}

	if dataTransferred.ClientBytes != nil {
		dataTransferredMap["client_bytes"] = flex.IntValue(dataTransferred.ClientBytes)
	}
	if dataTransferred.ServerBytes != nil {
		dataTransferredMap["server_bytes"] = flex.IntValue(dataTransferred.ServerBytes)
	}
	if dataTransferred.ClientPackets != nil {
		dataTransferredMap["client_packets"] = flex.IntValue(dataTransferred.ClientPackets)
	}
	if dataTransferred.ServerPackets != nil {
		dataTransferredMap["server_packets"] = flex.IntValue(dataTransferred.ServerPackets)
	}

	return dataTransferredMap
}

func resourceIBMSccSiOccurrenceKpiToMap(kpi findingsv1.Kpi) map[string]interface{} {
	kpiMap := map[string]interface{}{}

	kpiMap["value"] = kpi.Value
	if kpi.Total != nil {
		kpiMap["total"] = kpi.Total
	}

	return kpiMap
}

func resourceIBMSccSiOccurrenceUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	findingsClient, err := meta.(conns.ClientSession).FindingsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	updateOccurrenceOptions := &findingsv1.UpdateOccurrenceOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	findingsClient.AccountID = &parts[0]

	d.Set("account_id", parts[0])

	updateOccurrenceOptions.SetProviderID(parts[1])
	updateOccurrenceOptions.SetOccurrenceID(parts[2])
	updateOccurrenceOptions.SetProviderID(d.Get("provider_id").(string))
	updateOccurrenceOptions.SetNoteName(d.Get("note_name").(string))
	updateOccurrenceOptions.SetKind(d.Get("kind").(string))
	updateOccurrenceOptions.SetID(d.Get("occurrence_id").(string))
	if _, ok := d.GetOk("resource_url"); ok {
		updateOccurrenceOptions.SetResourceURL(d.Get("resource_url").(string))
	}
	if _, ok := d.GetOk("remediation"); ok {
		updateOccurrenceOptions.SetRemediation(d.Get("remediation").(string))
	}
	if _, ok := d.GetOk("context"); ok {
		context := resourceIBMSccSiOccurrenceMapToContext(d.Get("context.0").(map[string]interface{}))
		updateOccurrenceOptions.SetContext(&context)
	}
	if _, ok := d.GetOk("finding"); ok {
		finding := resourceIBMSccSiOccurrenceMapToFinding(d.Get("finding.0").(map[string]interface{}))
		updateOccurrenceOptions.SetFinding(&finding)
	}
	if _, ok := d.GetOk("kpi"); ok {
		kpi := resourceIBMSccSiOccurrenceMapToKpi(d.Get("kpi.0").(map[string]interface{}))
		updateOccurrenceOptions.SetKpi(&kpi)
	}
	_, response, err := findingsClient.UpdateOccurrenceWithContext(context, updateOccurrenceOptions)
	if err != nil {
		log.Printf("[DEBUG] UpdateOccurrenceWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("UpdateOccurrenceWithContext failed %s\n%s", err, response))
	}

	return resourceIBMSccSiOccurrenceRead(context, d, meta)
}

func resourceIBMSccSiOccurrenceDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	findingsClient, err := meta.(conns.ClientSession).FindingsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteOccurrenceOptions := &findingsv1.DeleteOccurrenceOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	findingsClient.AccountID = &parts[0]

	deleteOccurrenceOptions.SetProviderID(parts[1])
	deleteOccurrenceOptions.SetOccurrenceID(parts[2])

	response, err := findingsClient.DeleteOccurrenceWithContext(context, deleteOccurrenceOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteOccurrenceWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteOccurrenceWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
