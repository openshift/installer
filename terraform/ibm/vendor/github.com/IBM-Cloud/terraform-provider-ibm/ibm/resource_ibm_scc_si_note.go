// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/scc-go-sdk/findingsv1"
)

func resourceIBMSccSiNote() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMSccSiNoteCreate,
		ReadContext:   resourceIBMSccSiNoteRead,
		UpdateContext: resourceIBMSccSiNoteUpdate,
		DeleteContext: resourceIBMSccSiNoteDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"account_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"provider_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Part of the parent. This field contains the provider ID. For example: providers/{provider_id}.",
			},
			"short_description": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "A one sentence description of your note.",
			},
			"long_description": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "A more detailed description of your note.",
			},
			"kind": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: InvokeValidator("ibm_scc_si_note", "kind"),
				Description:  "The type of note. Use this field to filter notes and occurences by kind. - FINDING&#58; The note and occurrence represent a finding. - KPI&#58; The note and occurrence represent a KPI value. - CARD&#58; The note represents a card showing findings and related metric values. - CARD_CONFIGURED&#58; The note represents a card configured for a user account. - SECTION&#58; The note represents a section in a dashboard.",
			},
			"note_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the note.",
			},
			"reported_by": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "The entity reporting a note.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The id of this reporter.",
						},
						"title": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The title of this reporter.",
						},
						"url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The url of this reporter.",
						},
					},
				},
			},
			"related_url": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Label to describe usage of the URL.",
						},
						"url": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The URL that you want to associate with the note.",
						},
					},
				},
			},
			"shared": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "True if this note can be shared by multiple accounts.",
			},
			"finding": &schema.Schema{
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				Description:  "FindingType provides details about a finding note.",
				ExactlyOneOf: []string{"finding", "kpi", "card", "section"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"severity": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Note provider-assigned severity/impact ranking- LOW&#58; Low Impact- MEDIUM&#58; Medium Impact- HIGH&#58; High Impact- CRITICAL&#58; Critical Impact.",
						},
						"next_steps": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Common remediation steps for the finding of this type.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"title": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Title of this next step.",
									},
									"url": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The URL associated to this next steps.",
									},
								},
							},
						},
					},
				},
			},
			"kpi": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "KpiType provides details about a KPI note.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"aggregation_type": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The aggregation type of the KPI values. - SUM&#58; A single-value metrics aggregation type that sums up numeric values  that are extracted from KPI occurrences.",
						},
					},
				},
			},
			"card": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Card provides details about a card kind of note.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"section": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The section this card belongs to.",
						},
						"title": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The title of this card.",
						},
						"subtitle": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The subtitle of this card.",
						},
						"order": &schema.Schema{
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: InvokeValidator("ibm_scc_si_note", "order"),
							Description:  "The order of the card in which it will appear on SA dashboard in the mentioned section.",
						},
						"finding_note_names": &schema.Schema{
							Type:        schema.TypeList,
							Required:    true,
							Description: "The finding note names associated to this card.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"requires_configuration": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"badge_text": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The text associated to the card's badge.",
						},
						"badge_image": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The base64 content of the image associated to the card's badge.",
						},
						"elements": &schema.Schema{
							Type:        schema.TypeList,
							Required:    true,
							Description: "The elements of this card.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"text": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The text of this card element.",
									},
									"kind": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "NUMERIC",
										Description: "Kind of element- NUMERIC&#58; Single numeric value- BREAKDOWN&#58; Breakdown of numeric values- TIME_SERIES&#58; Time-series of numeric values.",
									},
									"default_time_range": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "4d",
										Description: "The default time range of this card element.",
									},
									"value_type": &schema.Schema{
										Type:     schema.TypeList,
										MaxItems: 1,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"kind": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Kind of element- KPI&#58; Kind of value derived from a KPI occurrence.",
												},
												"kpi_note_name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The name of the kpi note associated to the occurrence with the value for this card element value type.",
												},
												"text": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Default:     "label",
													Description: "The text of this element type.",
												},
												"finding_note_names": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "the names of the finding note associated that act as filters for counting the occurrences.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"value_types": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "the value types associated to this card element.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"kind": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Kind of element- KPI&#58; Kind of value derived from a KPI occurrence.",
												},
												"kpi_note_name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The name of the kpi note associated to the occurrence with the value for this card element value type.",
												},
												"text": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Default:     "label",
													Description: "The text of this element type.",
												},
												"finding_note_names": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "the names of the finding note associated that act as filters for counting the occurrences.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"default_interval": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "d",
										Description: "The default interval of the time series.",
									},
								},
							},
						},
					},
				},
			},
			"section": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Card provides details about a card kind of note.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"title": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The title of this section.",
						},
						"image": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The image of this section.",
						},
					},
				},
			},
			"create_time": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The time this note was created. This field can be used as a filter in list requests.",
			},
			"update_time": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The time this note was last updated. This field can be used as a filter in list requests.",
			},
		},
	}
}

func resourceIBMSccSiNoteValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 2)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "kind",
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              "CARD, CARD_CONFIGURED, FINDING, KPI, SECTION",
		},
		ValidateSchema{
			Identifier:                 "order",
			ValidateFunctionIdentifier: IntBetween,
			Type:                       TypeInt,
			Required:                   false,
			MinValue:                   "1",
			MaxValue:                   "6"},
	)

	resourceValidator := ResourceValidator{ResourceName: "ibm_scc_si_note", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMSccSiNoteCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	createNoteOptions := &findingsv1.CreateNoteOptions{}

	createNoteOptions.SetProviderID(d.Get("provider_id").(string))
	createNoteOptions.SetShortDescription(d.Get("short_description").(string))
	createNoteOptions.SetLongDescription(d.Get("long_description").(string))
	createNoteOptions.SetKind(d.Get("kind").(string))
	createNoteOptions.SetID(d.Get("note_id").(string))
	reportedBy := resourceIBMSccSiNoteMapToReporter(d.Get("reported_by.0").(map[string]interface{}))
	createNoteOptions.SetReportedBy(&reportedBy)
	if _, ok := d.GetOk("related_url"); ok {
		var relatedURL []findingsv1.APINoteRelatedURL
		for _, e := range d.Get("related_url").([]interface{}) {
			value := e.(map[string]interface{})
			relatedURLItem := resourceIBMSccSiNoteMapToAPINoteRelatedURL(value)
			relatedURL = append(relatedURL, relatedURLItem)
		}
		createNoteOptions.SetRelatedURL(relatedURL)
	}
	if _, ok := d.GetOk("shared"); ok {
		createNoteOptions.SetShared(d.Get("shared").(bool))
	}
	if _, ok := d.GetOk("finding"); ok {
		finding := resourceIBMSccSiNoteMapToFindingType(d.Get("finding.0").(map[string]interface{}))
		createNoteOptions.SetFinding(&finding)
	}
	if _, ok := d.GetOk("kpi"); ok {
		kpi := resourceIBMSccSiNoteMapToKpiType(d.Get("kpi.0").(map[string]interface{}))
		createNoteOptions.SetKpi(&kpi)
	}
	if _, ok := d.GetOk("card"); ok {
		card := resourceIBMSccSiNoteMapToCard(d.Get("card.0").(map[string]interface{}))
		createNoteOptions.SetCard(&card)
	}
	if _, ok := d.GetOk("section"); ok {
		section := resourceIBMSccSiNoteMapToSection(d.Get("section.0").(map[string]interface{}))
		createNoteOptions.SetSection(&section)
	}

	apiNote, response, err := findingsClient.CreateNoteWithContext(context, createNoteOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateNoteWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateNoteWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", *findingsClient.AccountID, *createNoteOptions.ProviderID, *apiNote.ID))

	return resourceIBMSccSiNoteRead(context, d, meta)
}

func resourceIBMSccSiNoteMapToReporter(reporterMap map[string]interface{}) findingsv1.Reporter {
	reporter := findingsv1.Reporter{}

	reporter.ID = core.StringPtr(reporterMap["id"].(string))
	reporter.Title = core.StringPtr(reporterMap["title"].(string))
	if reporterMap["url"] != nil {
		reporter.URL = core.StringPtr(reporterMap["url"].(string))
	}

	return reporter
}

func resourceIBMSccSiNoteMapToAPINoteRelatedURL(apiNoteRelatedURLMap map[string]interface{}) findingsv1.APINoteRelatedURL {
	apiNoteRelatedURL := findingsv1.APINoteRelatedURL{}

	apiNoteRelatedURL.Label = core.StringPtr(apiNoteRelatedURLMap["label"].(string))
	apiNoteRelatedURL.URL = core.StringPtr(apiNoteRelatedURLMap["url"].(string))

	return apiNoteRelatedURL
}

func resourceIBMSccSiNoteMapToFindingType(findingTypeMap map[string]interface{}) findingsv1.FindingType {
	findingType := findingsv1.FindingType{}

	findingType.Severity = core.StringPtr(findingTypeMap["severity"].(string))
	if findingTypeMap["next_steps"] != nil {
		nextSteps := []findingsv1.RemediationStep{}
		for _, nextStepsItem := range findingTypeMap["next_steps"].([]interface{}) {
			nextStepsItemModel := resourceIBMSccSiNoteMapToRemediationStep(nextStepsItem.(map[string]interface{}))
			nextSteps = append(nextSteps, nextStepsItemModel)
		}
		findingType.NextSteps = nextSteps
	}

	return findingType
}

func resourceIBMSccSiNoteMapToRemediationStep(remediationStepMap map[string]interface{}) findingsv1.RemediationStep {
	remediationStep := findingsv1.RemediationStep{}

	if remediationStepMap["title"] != nil {
		remediationStep.Title = core.StringPtr(remediationStepMap["title"].(string))
	}
	if remediationStepMap["url"] != nil {
		remediationStep.URL = core.StringPtr(remediationStepMap["url"].(string))
	}

	return remediationStep
}

func resourceIBMSccSiNoteMapToKpiType(kpiTypeMap map[string]interface{}) findingsv1.KpiType {
	kpiType := findingsv1.KpiType{}

	kpiType.AggregationType = core.StringPtr(kpiTypeMap["aggregation_type"].(string))

	return kpiType
}

func resourceIBMSccSiNoteMapToCard(cardMap map[string]interface{}) findingsv1.Card {
	card := findingsv1.Card{}

	card.Section = core.StringPtr(cardMap["section"].(string))
	card.Title = core.StringPtr(cardMap["title"].(string))
	card.Subtitle = core.StringPtr(cardMap["subtitle"].(string))
	if cardMap["order"] != nil && cardMap["order"].(int) > 0 {
		card.Order = core.Int64Ptr(int64(cardMap["order"].(int)))

	}
	findingNoteNames := []string{}
	for _, findingNoteNamesItem := range cardMap["finding_note_names"].([]interface{}) {
		findingNoteNames = append(findingNoteNames, findingNoteNamesItem.(string))
	}
	card.FindingNoteNames = findingNoteNames
	if cardMap["requires_configuration"] != nil {
		card.RequiresConfiguration = core.BoolPtr(cardMap["requires_configuration"].(bool))
	}
	if cardMap["badge_text"] != nil {
		card.BadgeText = core.StringPtr(cardMap["badge_text"].(string))
	}
	if cardMap["badge_image"] != nil {
		card.BadgeImage = core.StringPtr(cardMap["badge_image"].(string))
	}
	elements := []findingsv1.CardElementIntf{}
	for _, elementsItem := range cardMap["elements"].([]interface{}) {
		elementsItemModel := resourceIBMSccSiNoteMapToCardElement(elementsItem.(map[string]interface{}))
		elements = append(elements, elementsItemModel)
	}
	card.Elements = elements

	return card
}

func resourceIBMSccSiNoteMapToCardElement(cardElementMap map[string]interface{}) findingsv1.CardElementIntf {
	cardElement := findingsv1.CardElement{}

	if cardElementMap["text"] != nil {
		cardElement.Text = core.StringPtr(cardElementMap["text"].(string))
	}
	if cardElementMap["kind"] != nil {
		cardElement.Kind = core.StringPtr(cardElementMap["kind"].(string))
	}
	if cardElementMap["default_time_range"] != nil {
		cardElement.DefaultTimeRange = core.StringPtr(cardElementMap["default_time_range"].(string))
	}

	if cardElementMap["value_type"] != nil && len(cardElementMap["value_type"].([]interface{})) > 0 {
		cardElementValueType := findingsv1.NumericCardElementValueType{}

		if cardElementMap["value_type"].([]interface{})[0].(map[string]interface{})["kind"] != nil {
			cardElementValueType.Kind = core.StringPtr(cardElementMap["value_type"].([]interface{})[0].(map[string]interface{})["kind"].(string))
		}
		if cardElementMap["value_type"].([]interface{})[0].(map[string]interface{})["text"] != nil {
			cardElementValueType.Text = core.StringPtr(cardElementMap["value_type"].([]interface{})[0].(map[string]interface{})["text"].(string))
		}
		if cardElementMap["value_type"].([]interface{})[0].(map[string]interface{})["kpi_note_name"] != nil && cardElementMap["value_type"].([]interface{})[0].(map[string]interface{})["kpi_note_name"] != "" {
			cardElementValueType.KpiNoteName = core.StringPtr(cardElementMap["value_type"].([]interface{})[0].(map[string]interface{})["kpi_note_name"].(string))
		}
		if cardElementMap["value_type"].([]interface{})[0].(map[string]interface{})["finding_note_names"] != nil && len(cardElementMap["value_type"].([]interface{})[0].(map[string]interface{})["finding_note_names"].([]interface{})) > 0 {
			findingNoteNames := []string{}
			for _, findingNoteNamesItem := range cardElementMap["value_type"].([]interface{})[0].(map[string]interface{})["finding_note_names"].([]interface{}) {
				findingNoteNames = append(findingNoteNames, findingNoteNamesItem.(string))
			}
			cardElementValueType.FindingNoteNames = findingNoteNames
		}
		cardElement.ValueType = &cardElementValueType
	}

	if cardElementMap["value_types"] != nil {
		valueTypes := []findingsv1.ValueTypeIntf{}
		for _, valueTypesItem := range cardElementMap["value_types"].([]interface{}) {
			valueTypesItemModel := resourceIBMSccSiNoteMapToValueType(valueTypesItem.(map[string]interface{}))
			valueTypes = append(valueTypes, valueTypesItemModel)
		}
		cardElement.ValueTypes = valueTypes
	}
	if cardElementMap["default_interval"] != nil {
		cardElement.DefaultInterval = core.StringPtr(cardElementMap["default_interval"].(string))
	}

	return &cardElement
}

func resourceIBMSccSiNoteMapToNumericCardElementValueType(numericCardElementValueTypeMap map[string]interface{}) findingsv1.NumericCardElementValueType {
	numericCardElementValueType := findingsv1.NumericCardElementValueType{}

	if numericCardElementValueTypeMap["kind"] != nil {
		numericCardElementValueType.Kind = core.StringPtr(numericCardElementValueTypeMap["kind"].(string))
	}
	if numericCardElementValueTypeMap["kpi_note_name"] != nil {
		numericCardElementValueType.KpiNoteName = core.StringPtr(numericCardElementValueTypeMap["kpi_note_name"].(string))
	}
	if numericCardElementValueTypeMap["text"] != nil {
		numericCardElementValueType.Text = core.StringPtr(numericCardElementValueTypeMap["text"].(string))
	}
	if numericCardElementValueTypeMap["finding_note_names"] != nil {
		findingNoteNames := []string{}
		for _, findingNoteNamesItem := range numericCardElementValueTypeMap["finding_note_names"].([]interface{}) {
			findingNoteNames = append(findingNoteNames, findingNoteNamesItem.(string))
		}
		numericCardElementValueType.FindingNoteNames = findingNoteNames
	}

	return numericCardElementValueType
}

func resourceIBMSccSiNoteMapToValueType(valueTypeMap map[string]interface{}) findingsv1.ValueTypeIntf {
	valueType := findingsv1.ValueType{}

	if valueTypeMap["kind"] != nil {
		valueType.Kind = core.StringPtr(valueTypeMap["kind"].(string))
	}
	if valueTypeMap["kpi_note_name"] != nil && len(valueTypeMap["kpi_note_name"].(string)) > 0 {
		valueType.KpiNoteName = core.StringPtr(valueTypeMap["kpi_note_name"].(string))
	}
	if valueTypeMap["text"] != nil {
		valueType.Text = core.StringPtr(valueTypeMap["text"].(string))
	}
	if valueTypeMap["finding_note_names"] != nil {
		findingNoteNames := []string{}
		for _, findingNoteNamesItem := range valueTypeMap["finding_note_names"].([]interface{}) {
			findingNoteNames = append(findingNoteNames, findingNoteNamesItem.(string))
		}
		valueType.FindingNoteNames = findingNoteNames
	}

	return &valueType
}

func resourceIBMSccSiNoteMapToValueTypeFindingCountValueType(valueTypeFindingCountValueTypeMap map[string]interface{}) findingsv1.ValueTypeFindingCountValueType {
	valueTypeFindingCountValueType := findingsv1.ValueTypeFindingCountValueType{}

	valueTypeFindingCountValueType.Kind = core.StringPtr(valueTypeFindingCountValueTypeMap["kind"].(string))
	findingNoteNames := []string{}
	for _, findingNoteNamesItem := range valueTypeFindingCountValueTypeMap["finding_note_names"].([]interface{}) {
		findingNoteNames = append(findingNoteNames, findingNoteNamesItem.(string))
	}
	valueTypeFindingCountValueType.FindingNoteNames = findingNoteNames
	valueTypeFindingCountValueType.Text = core.StringPtr(valueTypeFindingCountValueTypeMap["text"].(string))

	return valueTypeFindingCountValueType
}

func resourceIBMSccSiNoteMapToValueTypeKpiValueType(valueTypeKpiValueTypeMap map[string]interface{}) findingsv1.ValueTypeKpiValueType {
	valueTypeKpiValueType := findingsv1.ValueTypeKpiValueType{}

	valueTypeKpiValueType.Kind = core.StringPtr(valueTypeKpiValueTypeMap["kind"].(string))
	valueTypeKpiValueType.KpiNoteName = core.StringPtr(valueTypeKpiValueTypeMap["kpi_note_name"].(string))
	valueTypeKpiValueType.Text = core.StringPtr(valueTypeKpiValueTypeMap["text"].(string))

	return valueTypeKpiValueType
}

func resourceIBMSccSiNoteMapToCardElementTimeSeriesCardElement(cardElementTimeSeriesCardElementMap map[string]interface{}) findingsv1.CardElementTimeSeriesCardElement {
	cardElementTimeSeriesCardElement := findingsv1.CardElementTimeSeriesCardElement{}

	cardElementTimeSeriesCardElement.Text = core.StringPtr(cardElementTimeSeriesCardElementMap["text"].(string))
	if cardElementTimeSeriesCardElementMap["default_interval"] != nil {
		cardElementTimeSeriesCardElement.DefaultInterval = core.StringPtr(cardElementTimeSeriesCardElementMap["default_interval"].(string))
	}
	cardElementTimeSeriesCardElement.Kind = core.StringPtr(cardElementTimeSeriesCardElementMap["kind"].(string))
	if cardElementTimeSeriesCardElementMap["default_time_range"] != nil {
		cardElementTimeSeriesCardElement.DefaultTimeRange = core.StringPtr(cardElementTimeSeriesCardElementMap["default_time_range"].(string))
	}
	valueTypes := []findingsv1.ValueTypeIntf{}
	for _, valueTypesItem := range cardElementTimeSeriesCardElementMap["value_types"].([]interface{}) {
		valueTypesItemModel := resourceIBMSccSiNoteMapToValueType(valueTypesItem.(map[string]interface{}))
		valueTypes = append(valueTypes, valueTypesItemModel)
	}
	cardElementTimeSeriesCardElement.ValueTypes = valueTypes

	return cardElementTimeSeriesCardElement
}

func resourceIBMSccSiNoteMapToCardElementBreakdownCardElement(cardElementBreakdownCardElementMap map[string]interface{}) findingsv1.CardElementBreakdownCardElement {
	cardElementBreakdownCardElement := findingsv1.CardElementBreakdownCardElement{}

	cardElementBreakdownCardElement.Text = core.StringPtr(cardElementBreakdownCardElementMap["text"].(string))
	cardElementBreakdownCardElement.Kind = core.StringPtr(cardElementBreakdownCardElementMap["kind"].(string))
	if cardElementBreakdownCardElementMap["default_time_range"] != nil {
		cardElementBreakdownCardElement.DefaultTimeRange = core.StringPtr(cardElementBreakdownCardElementMap["default_time_range"].(string))
	}
	valueTypes := []findingsv1.ValueTypeIntf{}
	for _, valueTypesItem := range cardElementBreakdownCardElementMap["value_types"].([]interface{}) {
		valueTypesItemModel := resourceIBMSccSiNoteMapToValueType(valueTypesItem.(map[string]interface{}))
		valueTypes = append(valueTypes, valueTypesItemModel)
	}
	cardElementBreakdownCardElement.ValueTypes = valueTypes

	return cardElementBreakdownCardElement
}

func resourceIBMSccSiNoteMapToCardElementNumericCardElement(cardElementNumericCardElementMap map[string]interface{}) findingsv1.CardElementNumericCardElement {
	cardElementNumericCardElement := findingsv1.CardElementNumericCardElement{}

	cardElementNumericCardElement.Text = core.StringPtr(cardElementNumericCardElementMap["text"].(string))
	cardElementNumericCardElement.Kind = core.StringPtr(cardElementNumericCardElementMap["kind"].(string))
	if cardElementNumericCardElementMap["default_time_range"] != nil {
		cardElementNumericCardElement.DefaultTimeRange = core.StringPtr(cardElementNumericCardElementMap["default_time_range"].(string))
	}
	// TODO: handle ValueType of type NumericCardElementValueType -- not primitive type, not list

	return cardElementNumericCardElement
}

func resourceIBMSccSiNoteMapToSection(sectionMap map[string]interface{}) findingsv1.Section {
	section := findingsv1.Section{}

	section.Title = core.StringPtr(sectionMap["title"].(string))
	section.Image = core.StringPtr(sectionMap["image"].(string))

	return section
}

func resourceIBMSccSiNoteRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	findingsClient, err := meta.(ClientSession).FindingsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getNoteOptions := &findingsv1.GetNoteOptions{}

	parts, err := sepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	findingsClient.AccountID = &parts[0]

	d.Set("account_id", &parts[0])

	getNoteOptions.SetProviderID(parts[1])
	getNoteOptions.SetNoteID(parts[2])

	apiNote, response, err := findingsClient.GetNoteWithContext(context, getNoteOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetNoteWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetNoteWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("provider_id", getNoteOptions.ProviderID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting provider_id: %s", err))
	}
	if err = d.Set("short_description", apiNote.ShortDescription); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting short_description: %s", err))
	}
	if err = d.Set("long_description", apiNote.LongDescription); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting long_description: %s", err))
	}
	if err = d.Set("kind", apiNote.Kind); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting kind: %s", err))
	}
	if err = d.Set("note_id", apiNote.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting note_id: %s", err))
	}
	reportedByMap := resourceIBMSccSiNoteReporterToMap(*apiNote.ReportedBy)
	if err = d.Set("reported_by", []map[string]interface{}{reportedByMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting reported_by: %s", err))
	}
	if apiNote.RelatedURL != nil {
		relatedURL := []map[string]interface{}{}
		for _, relatedURLItem := range apiNote.RelatedURL {
			relatedURLItemMap := resourceIBMSccSiNoteAPINoteRelatedURLToMap(relatedURLItem)
			relatedURL = append(relatedURL, relatedURLItemMap)
		}
		if err = d.Set("related_url", relatedURL); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting related_url: %s", err))
		}
	}
	if err = d.Set("shared", apiNote.Shared); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting shared: %s", err))
	}
	if apiNote.Finding != nil {
		findingMap := resourceIBMSccSiNoteFindingTypeToMap(*apiNote.Finding)
		if err = d.Set("finding", []map[string]interface{}{findingMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting finding: %s", err))
		}
	}
	if apiNote.Kpi != nil {
		kpiMap := resourceIBMSccSiNoteKpiTypeToMap(*apiNote.Kpi)
		if err = d.Set("kpi", []map[string]interface{}{kpiMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting kpi: %s", err))
		}
	}
	if apiNote.Card != nil {
		cardIntf := d.Get("card")
		cardMap := resourceIBMSccSiNoteCardToMap(*apiNote.Card, cardIntf)
		if err = d.Set("card", []map[string]interface{}{cardMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting card: %s", err))
		}
	}
	if apiNote.Section != nil {
		sectionMap := resourceIBMSccSiNoteSectionToMap(*apiNote.Section)
		if err = d.Set("section", []map[string]interface{}{sectionMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting section: %s", err))
		}
	}
	if err = d.Set("create_time", dateTimeToString(apiNote.CreateTime)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting create_time: %s", err))
	}
	if err = d.Set("update_time", dateTimeToString(apiNote.UpdateTime)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting update_time: %s", err))
	}

	return nil
}

func resourceIBMSccSiNoteReporterToMap(reporter findingsv1.Reporter) map[string]interface{} {
	reporterMap := map[string]interface{}{}

	reporterMap["id"] = reporter.ID
	reporterMap["title"] = reporter.Title
	if reporter.URL != nil {
		reporterMap["url"] = reporter.URL
	}

	return reporterMap
}

func resourceIBMSccSiNoteAPINoteRelatedURLToMap(apiNoteRelatedURL findingsv1.APINoteRelatedURL) map[string]interface{} {
	apiNoteRelatedURLMap := map[string]interface{}{}

	apiNoteRelatedURLMap["label"] = apiNoteRelatedURL.Label
	apiNoteRelatedURLMap["url"] = apiNoteRelatedURL.URL

	return apiNoteRelatedURLMap
}

func resourceIBMSccSiNoteFindingTypeToMap(findingType findingsv1.FindingType) map[string]interface{} {
	findingTypeMap := map[string]interface{}{}

	findingTypeMap["severity"] = findingType.Severity
	if findingType.NextSteps != nil {
		nextSteps := []map[string]interface{}{}
		for _, nextStepsItem := range findingType.NextSteps {
			nextStepsItemMap := resourceIBMSccSiNoteRemediationStepToMap(nextStepsItem)
			nextSteps = append(nextSteps, nextStepsItemMap)
			// TODO: handle NextSteps of type TypeList -- list of non-primitive, not model items
		}
		findingTypeMap["next_steps"] = nextSteps
	}

	return findingTypeMap
}

func resourceIBMSccSiNoteRemediationStepToMap(remediationStep findingsv1.RemediationStep) map[string]interface{} {
	remediationStepMap := map[string]interface{}{}

	if remediationStep.Title != nil {
		remediationStepMap["title"] = remediationStep.Title
	}
	if remediationStep.URL != nil {
		remediationStepMap["url"] = remediationStep.URL
	}

	return remediationStepMap
}

func resourceIBMSccSiNoteKpiTypeToMap(kpiType findingsv1.KpiType) map[string]interface{} {
	kpiTypeMap := map[string]interface{}{}

	kpiTypeMap["aggregation_type"] = kpiType.AggregationType

	return kpiTypeMap
}

func resourceIBMSccSiNoteCardToMap(card findingsv1.Card, cardIntf interface{}) map[string]interface{} {
	cardMap := map[string]interface{}{}

	cardMap["section"] = card.Section
	cardMap["title"] = card.Title
	cardMap["subtitle"] = card.Subtitle
	if card.Order != nil {
		order := intValue(card.Order)
		if order != 0 {
			cardMap["order"] = intValue(card.Order)
		}
	}
	cardMap["finding_note_names"] = card.FindingNoteNames
	if card.RequiresConfiguration != nil {
		cardMap["requires_configuration"] = card.RequiresConfiguration
	}
	if card.BadgeText != nil {
		cardMap["badge_text"] = card.BadgeText
	}
	if card.BadgeImage != nil {
		cardMap["badge_image"] = card.BadgeImage
	}
	elements := []map[string]interface{}{}
	for i, elementsItem := range card.Elements {
		var elemResource interface{}
		if cardIntf != nil && len(cardIntf.([]interface{})) > 0 {
			elemResource = cardIntf.([]interface{})[0].(map[string]interface{})["elements"].([]interface{})[i]
		}
		elementsItemMap := resourceIBMSccSiNoteCardElementToMap(elementsItem, elemResource)
		elements = append(elements, elementsItemMap)
		// TODO: handle Elements of type TypeList -- list of non-primitive, not model items
	}
	cardMap["elements"] = elements

	return cardMap
}

func resourceIBMSccSiNoteCardElementToMap(cardElement findingsv1.CardElementIntf, elemResource interface{}) map[string]interface{} {
	cardElementMap := map[string]interface{}{}

	switch v := cardElement.(type) {
	case *findingsv1.CardElementNumericCardElement:
		cardElementMap = resourceIBMSccSiNoteCardElementNumericCardElementToMap(*v, elemResource)
	case *findingsv1.CardElementBreakdownCardElement:
		cardElementMap = resourceIBMSccSiNoteCardElementBreakdownCardElementToMap(*v, elemResource)
	case *findingsv1.CardElementTimeSeriesCardElement:
		cardElementMap = resourceIBMSccSiNoteCardElementTimeSeriesCardElementToMap(*v, elemResource)
	default:
		log.Printf("[DEBUG] Unknown card element type")
	}

	return cardElementMap
}

func resourceIBMSccSiNoteNumericCardElementValueTypeToMap(numericCardElementValueType findingsv1.NumericCardElementValueType, elemResource interface{}) map[string]interface{} {
	numericCardElementValueTypeMap := map[string]interface{}{}

	if numericCardElementValueType.Kind != nil {
		numericCardElementValueTypeMap["kind"] = numericCardElementValueType.Kind
	}

	if numericCardElementValueType.KpiNoteName != nil {
		numericCardElementValueTypeMap["kpi_note_name"] = numericCardElementValueType.KpiNoteName
		if elemResource != nil {
			findingNoteNamesMap := make([]string, 0)
			findingNoteNames := elemResource.(map[string]interface{})["finding_note_names"]
			if findingNoteNames != nil {
				for _, findingNoteName := range findingNoteNames.([]interface{}) {
					findingNoteNamesMap = append(findingNoteNamesMap, findingNoteName.(string))
				}
				numericCardElementValueTypeMap["finding_note_names"] = findingNoteNamesMap
			}
		}
	}

	if numericCardElementValueType.FindingNoteNames != nil {
		numericCardElementValueTypeMap["finding_note_names"] = numericCardElementValueType.FindingNoteNames
		if elemResource != nil {
			kpiNoteName := elemResource.(map[string]interface{})["kpi_note_name"]
			if kpiNoteName != nil {
				numericCardElementValueTypeMap["kpi_note_name"] = kpiNoteName
			}
		}
	}

	if numericCardElementValueType.Text != nil {
		numericCardElementValueTypeMap["text"] = numericCardElementValueType.Text
	}

	return numericCardElementValueTypeMap
}

func resourceIBMSccSiNoteValueTypeToMap(valueType findingsv1.ValueTypeIntf, elemResource interface{}) map[string]interface{} {
	valueTypeMap := map[string]interface{}{}

	switch v := valueType.(type) {

	case *findingsv1.ValueTypeFindingCountValueType:
		valueTypeMap["kind"] = v.Kind
		valueTypeMap["finding_note_names"] = v.FindingNoteNames
		valueTypeMap["text"] = v.Text
		kpiNoteName := elemResource.(map[string]interface{})["kpi_note_name"]
		if kpiNoteName == nil {
			valueTypeMap["kpi_note_name"] = ""
		} else {
			valueTypeMap["kpi_note_name"] = kpiNoteName
		}
	case *findingsv1.ValueTypeKpiValueType:
		valueTypeMap["kind"] = v.Kind
		valueTypeMap["kpi_note_name"] = v.KpiNoteName
		valueTypeMap["text"] = v.Text
		findingNoteNames := elemResource.(map[string]interface{})["finding_note_names"]
		if findingNoteNames == nil {
			valueTypeMap["finding_note_names"] = []string{}
		} else {
			findingNoteNamesMap := make([]string, 0)
			for _, findingNoteName := range findingNoteNames.([]interface{}) {
				findingNoteNamesMap = append(findingNoteNamesMap, findingNoteName.(string))
			}
			valueTypeMap["finding_note_names"] = findingNoteNamesMap
		}
	default:
		log.Printf("[DEBUG] Unknown card element value_type type")
	}

	return valueTypeMap
}

func resourceIBMSccSiNoteValueTypeFindingCountValueTypeToMap(valueTypeFindingCountValueType findingsv1.ValueTypeFindingCountValueType) map[string]interface{} {
	valueTypeFindingCountValueTypeMap := map[string]interface{}{}

	valueTypeFindingCountValueTypeMap["kind"] = valueTypeFindingCountValueType.Kind
	valueTypeFindingCountValueTypeMap["finding_note_names"] = valueTypeFindingCountValueType.FindingNoteNames
	valueTypeFindingCountValueTypeMap["text"] = valueTypeFindingCountValueType.Text

	return valueTypeFindingCountValueTypeMap
}

func resourceIBMSccSiNoteValueTypeKpiValueTypeToMap(valueTypeKpiValueType findingsv1.ValueTypeKpiValueType) map[string]interface{} {
	valueTypeKpiValueTypeMap := map[string]interface{}{}

	valueTypeKpiValueTypeMap["kind"] = valueTypeKpiValueType.Kind
	valueTypeKpiValueTypeMap["kpi_note_name"] = valueTypeKpiValueType.KpiNoteName
	valueTypeKpiValueTypeMap["text"] = valueTypeKpiValueType.Text

	return valueTypeKpiValueTypeMap
}

func resourceIBMSccSiNoteCardElementTimeSeriesCardElementToMap(cardElementTimeSeriesCardElement findingsv1.CardElementTimeSeriesCardElement, elemResource interface{}) map[string]interface{} {
	cardElementTimeSeriesCardElementMap := map[string]interface{}{}

	cardElementTimeSeriesCardElementMap["text"] = cardElementTimeSeriesCardElement.Text
	if cardElementTimeSeriesCardElement.DefaultInterval != nil {
		cardElementTimeSeriesCardElementMap["default_interval"] = cardElementTimeSeriesCardElement.DefaultInterval
	}
	cardElementTimeSeriesCardElementMap["kind"] = cardElementTimeSeriesCardElement.Kind
	if cardElementTimeSeriesCardElement.DefaultTimeRange != nil {
		cardElementTimeSeriesCardElementMap["default_time_range"] = cardElementTimeSeriesCardElement.DefaultTimeRange
	}
	valueTypes := []map[string]interface{}{}
	for _, valueTypesItem := range cardElementTimeSeriesCardElement.ValueTypes {
		valueTypesItemMap := resourceIBMSccSiNoteValueTypeToMap(valueTypesItem, elemResource)
		valueTypes = append(valueTypes, valueTypesItemMap)
		// TODO: handle ValueTypes of type TypeList -- list of non-primitive, not model items
	}
	cardElementTimeSeriesCardElementMap["value_types"] = valueTypes

	return cardElementTimeSeriesCardElementMap
}

func resourceIBMSccSiNoteCardElementBreakdownCardElementToMap(cardElementBreakdownCardElement findingsv1.CardElementBreakdownCardElement, elemResource interface{}) map[string]interface{} {
	cardElementBreakdownCardElementMap := map[string]interface{}{}

	cardElementBreakdownCardElementMap["text"] = cardElementBreakdownCardElement.Text
	cardElementBreakdownCardElementMap["kind"] = cardElementBreakdownCardElement.Kind
	if cardElementBreakdownCardElement.DefaultTimeRange != nil {
		cardElementBreakdownCardElementMap["default_time_range"] = cardElementBreakdownCardElement.DefaultTimeRange
	}
	valueTypes := []map[string]interface{}{}
	for i, valueTypesItem := range cardElementBreakdownCardElement.ValueTypes {
		valueType := elemResource.(map[string]interface{})["value_types"].([]interface{})[i]
		valueTypesItemMap := resourceIBMSccSiNoteValueTypeToMap(valueTypesItem, valueType)
		valueTypes = append(valueTypes, valueTypesItemMap)
		// TODO: handle ValueTypes of type TypeList -- list of non-primitive, not model items
	}
	cardElementBreakdownCardElementMap["value_types"] = valueTypes

	if elemResource != nil {
		if elemResource.(map[string]interface{})["default_interval"] != nil {
			cardElementBreakdownCardElementMap["default_interval"] = elemResource.(map[string]interface{})["default_interval"].(string)
		}
	}

	return cardElementBreakdownCardElementMap
}

func resourceIBMSccSiNoteCardElementNumericCardElementToMap(cardElementNumericCardElement findingsv1.CardElementNumericCardElement, elemResource interface{}) map[string]interface{} {
	cardElementNumericCardElementMap := map[string]interface{}{}

	cardElementNumericCardElementMap["text"] = cardElementNumericCardElement.Text
	cardElementNumericCardElementMap["kind"] = cardElementNumericCardElement.Kind
	if cardElementNumericCardElement.DefaultTimeRange != nil {
		cardElementNumericCardElementMap["default_time_range"] = cardElementNumericCardElement.DefaultTimeRange
	}
	ValueTypeMap := resourceIBMSccSiNoteNumericCardElementValueTypeToMap(*cardElementNumericCardElement.ValueType, elemResource)
	cardElementNumericCardElementMap["value_type"] = []map[string]interface{}{ValueTypeMap}

	if elemResource != nil {
		if elemResource.(map[string]interface{})["default_interval"] != nil {
			cardElementNumericCardElementMap["default_interval"] = elemResource.(map[string]interface{})["default_interval"].(string)
		}
	}

	return cardElementNumericCardElementMap
}

func resourceIBMSccSiNoteSectionToMap(section findingsv1.Section) map[string]interface{} {
	sectionMap := map[string]interface{}{}

	sectionMap["title"] = section.Title
	sectionMap["image"] = section.Image

	return sectionMap
}

func resourceIBMSccSiNoteUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	findingsClient, err := meta.(ClientSession).FindingsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	updateNoteOptions := &findingsv1.UpdateNoteOptions{}

	parts, err := sepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	findingsClient.AccountID = &parts[0]

	updateNoteOptions.SetProviderID(parts[1])
	updateNoteOptions.SetNoteID(parts[2])
	updateNoteOptions.SetProviderID(d.Get("provider_id").(string))
	updateNoteOptions.SetShortDescription(d.Get("short_description").(string))
	updateNoteOptions.SetLongDescription(d.Get("long_description").(string))
	updateNoteOptions.SetKind(d.Get("kind").(string))
	updateNoteOptions.SetID(d.Get("note_id").(string))
	reportedBy := resourceIBMSccSiNoteMapToReporter(d.Get("reported_by.0").(map[string]interface{}))
	updateNoteOptions.SetReportedBy(&reportedBy)
	if _, ok := d.GetOk("related_url"); ok {
		var relatedURL []findingsv1.APINoteRelatedURL
		for _, e := range d.Get("related_url").([]interface{}) {
			value := e.(map[string]interface{})
			relatedURLItem := resourceIBMSccSiNoteMapToAPINoteRelatedURL(value)
			relatedURL = append(relatedURL, relatedURLItem)
		}
		updateNoteOptions.SetRelatedURL(relatedURL)
	}
	if _, ok := d.GetOk("shared"); ok {
		updateNoteOptions.SetShared(d.Get("shared").(bool))
	}
	if _, ok := d.GetOk("finding"); ok {
		finding := resourceIBMSccSiNoteMapToFindingType(d.Get("finding.0").(map[string]interface{}))
		updateNoteOptions.SetFinding(&finding)
	}
	if _, ok := d.GetOk("kpi"); ok {
		kpi := resourceIBMSccSiNoteMapToKpiType(d.Get("kpi.0").(map[string]interface{}))
		updateNoteOptions.SetKpi(&kpi)
	}
	if _, ok := d.GetOk("card"); ok {
		card := resourceIBMSccSiNoteMapToCard(d.Get("card.0").(map[string]interface{}))
		updateNoteOptions.SetCard(&card)
	}
	if _, ok := d.GetOk("section"); ok {
		section := resourceIBMSccSiNoteMapToSection(d.Get("section.0").(map[string]interface{}))
		updateNoteOptions.SetSection(&section)
	}

	_, response, err := findingsClient.UpdateNoteWithContext(context, updateNoteOptions)
	if err != nil {
		log.Printf("[DEBUG] UpdateNoteWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("UpdateNoteWithContext failed %s\n%s", err, response))
	}

	return resourceIBMSccSiNoteRead(context, d, meta)
}

func resourceIBMSccSiNoteDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	findingsClient, err := meta.(ClientSession).FindingsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteNoteOptions := &findingsv1.DeleteNoteOptions{}

	parts, err := sepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	findingsClient.AccountID = &parts[0]

	deleteNoteOptions.SetProviderID(parts[1])
	deleteNoteOptions.SetNoteID(parts[2])

	response, err := findingsClient.DeleteNoteWithContext(context, deleteNoteOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteNoteWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteNoteWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
