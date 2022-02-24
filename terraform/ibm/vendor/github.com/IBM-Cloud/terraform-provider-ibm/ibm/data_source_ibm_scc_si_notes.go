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

func ValidatePageSize(val interface{}, key string) (warns []string, errs []error) {
	v := int64(val.(int))
	if v < 2 {
		errs = append(errs, fmt.Errorf("%q must be atleast 2, got: %d", key, v))
	}
	return
}

func dataSourceIBMSccSiNotes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMSccSiNotesRead,

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
			"notes": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The notes requested.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"note_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the note.",
						},
						"short_description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A one sentence description of your note.",
						},
						"long_description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A more detailed description of your note.",
						},
						"kind": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of note. Use this field to filter notes and occurences by kind. - FINDING&#58; The note and occurrence represent a finding. - KPI&#58; The note and occurrence represent a KPI value. - CARD&#58; The note represents a card showing findings and related metric values. - CARD_CONFIGURED&#58; The note represents a card configured for a user account. - SECTION&#58; The note represents a section in a dashboard.",
						},
						"related_url": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"label": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Label to describe usage of the URL.",
									},
									"url": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL that you want to associate with the note.",
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
						"shared": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "True if this note can be shared by multiple accounts.",
						},
						"reported_by": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The entity reporting a note.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The id of this reporter.",
									},
									"title": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The title of this reporter.",
									},
									"url": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The url of this reporter.",
									},
								},
							},
						},
						"finding": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "FindingType provides details about a finding note.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"severity": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Note provider-assigned severity/impact ranking- LOW&#58; Low Impact- MEDIUM&#58; Medium Impact- HIGH&#58; High Impact- CRITICAL&#58; Critical Impact.",
									},
									"next_steps": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Common remediation steps for the finding of this type.",
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
								},
							},
						},
						"kpi": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "KpiType provides details about a KPI note.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"aggregation_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The aggregation type of the KPI values. - SUM&#58; A single-value metrics aggregation type that sums up numeric values  that are extracted from KPI occurrences.",
									},
								},
							},
						},
						"card": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Card provides details about a card kind of note.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"section": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The section this card belongs to.",
									},
									"title": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The title of this card.",
									},
									"subtitle": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The subtitle of this card.",
									},
									"order": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The order of the card in which it will appear on SA dashboard in the mentioned section.",
									},
									"finding_note_names": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The finding note names associated to this card.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"requires_configuration": &schema.Schema{
										Type:     schema.TypeBool,
										Computed: true,
									},
									"badge_text": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The text associated to the card's badge.",
									},
									"badge_image": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The base64 content of the image associated to the card's badge.",
									},
									"elements": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The elements of this card.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"text": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The text of this card element.",
												},
												"kind": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Kind of element- NUMERIC&#58; Single numeric value- BREAKDOWN&#58; Breakdown of numeric values- TIME_SERIES&#58; Time-series of numeric values.",
												},
												"default_time_range": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The default time range of this card element.",
												},
												"value_type": &schema.Schema{
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"kind": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Kind of element- KPI&#58; Kind of value derived from a KPI occurrence.",
															},
															"kpi_note_name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The name of the kpi note associated to the occurrence with the value for this card element value type.",
															},
															"text": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The text of this element type.",
															},
															"finding_note_names": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "the names of the finding note associated that act as filters for counting the occurrences.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
														},
													},
												},
												"value_types": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "the value types associated to this card element.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"kind": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Kind of element- KPI&#58; Kind of value derived from a KPI occurrence.",
															},
															"kpi_note_name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The name of the kpi note associated to the occurrence with the value for this card element value type.",
															},
															"text": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The text of this element type.",
															},
															"finding_note_names": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "the names of the finding note associated that act as filters for counting the occurrences.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
														},
													},
												},
												"default_interval": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
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
							Computed:    true,
							Description: "Card provides details about a card kind of note.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"title": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The title of this section.",
									},
									"image": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The image of this section.",
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

func dataSourceIBMSccSiNotesValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 0)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "page_size",
			ValidateFunctionIdentifier: IntBetween,
			Required:                   false,
			MinValue:                   "2"})

	ibmSccSiNotesDataSourceValidator := ResourceValidator{ResourceName: "ibm_scc_si_notes", Schema: validateSchema}
	return &ibmSccSiNotesDataSourceValidator
}

func dataSourceIBMSccSiNotesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	listNoteOptions := &findingsv1.ListNotesOptions{}

	if pageSize, ok := d.GetOk("page_size"); ok {
		listNoteOptions.SetPageSize(int64(pageSize.(int)))
	}
	if pageToken, ok := d.GetOk("page_token"); ok {
		listNoteOptions.SetPageToken(pageToken.(string))
	}
	listNoteOptions.SetProviderID(d.Get("provider_id").(string))

	apiNotes := []findingsv1.APINote{}

	if listNoteOptions.PageToken != nil {
		apiNotes, err = collectSpecificNotes(findingsClient, context, listNoteOptions)
		if err != nil {
			log.Printf("[DEBUG] GetNoteWithContext failed %s", err)
			return diag.FromErr(fmt.Errorf("GetNoteWithContext failed %s", err))
		}
	} else {
		apiNotes, err = collectAllNotes(findingsClient, context, listNoteOptions)
		if err != nil {
			log.Printf("[DEBUG] GetNoteWithContext failed %s", err)
			return diag.FromErr(fmt.Errorf("GetNoteWithContext failed %s", err))
		}
	}

	d.SetId(dataSourceIBMSccSiNotesID(d))

	if apiNotes != nil {
		err = d.Set("notes", dataSourceAPIListNotesResponseFlattenProviders(apiNotes))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting notes %s", err))
		}
	}

	return nil
}

// dataSourceIBMSccSiNotesID returns a reasonable ID for the list.
func dataSourceIBMSccSiNotesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func collectSpecificNotes(findingsClient *findingsv1.FindingsV1, ctx context.Context, options *findingsv1.ListNotesOptions) ([]findingsv1.APINote, error) {
	apiListNotesResponse, response, err := findingsClient.ListNotesWithContext(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("%s\n%s", err, response)
	}

	return apiListNotesResponse.Notes, nil
}

func collectAllNotes(findingsClient *findingsv1.FindingsV1, ctx context.Context, options *findingsv1.ListNotesOptions) ([]findingsv1.APINote, error) {
	finalList := []findingsv1.APINote{}

	for {
		apiListNotesResponse, response, err := findingsClient.ListNotesWithContext(ctx, options)
		if err != nil {
			return nil, fmt.Errorf("%s\n%s", err, response)
		}

		finalList = append(finalList, apiListNotesResponse.Notes...)

		if options.PageSize != nil && int64(len(finalList)) == *options.PageSize {
			break
		}

		options.PageToken = apiListNotesResponse.NextPageToken

		if *apiListNotesResponse.NextPageToken == "" {
			break
		}
	}

	return finalList, nil
}

func dataSourceAPIListNotesResponseFlattenProviders(result []findingsv1.APINote) (notes []map[string]interface{}) {
	for _, notesItem := range result {
		notes = append(notes, dataSourceAPIListNotesResponseProvidersToMap(notesItem))
	}

	return notes
}

func dataSourceAPIListNotesResponseProvidersToMap(notesItem findingsv1.APINote) (notesMap map[string]interface{}) {
	notesMap = map[string]interface{}{}

	if notesItem.ID != nil {
		notesMap["note_id"] = notesItem.ID
	}
	if notesItem.ShortDescription != nil {
		notesMap["short_description"] = notesItem.ShortDescription
	}
	if notesItem.LongDescription != nil {
		notesMap["long_description"] = notesItem.LongDescription
	}
	if notesItem.Kind != nil {
		notesMap["kind"] = notesItem.Kind
	}

	if notesItem.RelatedURL != nil {
		notesMap["related_url"] = dataSourceAPINotesFlattenRelatedURL(notesItem.RelatedURL)
	}

	if notesItem.CreateTime != nil {
		notesMap["create_time"] = dateTimeToString(notesItem.CreateTime)
	}
	if notesItem.UpdateTime != nil {
		notesMap["update_time"] = dateTimeToString(notesItem.UpdateTime)
	}
	if notesItem.Shared != nil {
		notesMap["shared"] = notesItem.Shared
	}

	if notesItem.ReportedBy != nil {
		notesMap["reported_by"] = dataSourceAPINotesFlattenReportedBy(*notesItem.ReportedBy)
	}

	if notesItem.Finding != nil {
		notesMap["finding"] = dataSourceAPINotesFlattenFinding(*notesItem.Finding)
	}

	if notesItem.Kpi != nil {
		notesMap["kpi"] = dataSourceAPINotesFlattenKpi(*notesItem.Kpi)
	}

	if notesItem.Card != nil {
		notesMap["card"] = dataSourceAPINotesFlattenCard(*notesItem.Card)
	}

	if notesItem.Section != nil {
		notesMap["section"] = dataSourceAPINotesFlattenSection(*notesItem.Section)
	}

	return notesMap
}

func dataSourceAPINotesFlattenRelatedURL(result []findingsv1.APINoteRelatedURL) (relatedURL []map[string]interface{}) {
	for _, relatedURLItem := range result {
		relatedURL = append(relatedURL, dataSourceAPINotesRelatedURLToMap(relatedURLItem))
	}

	return relatedURL
}

func dataSourceAPINotesRelatedURLToMap(relatedURLItem findingsv1.APINoteRelatedURL) (relatedURLMap map[string]interface{}) {
	relatedURLMap = map[string]interface{}{}

	if relatedURLItem.Label != nil {
		relatedURLMap["label"] = relatedURLItem.Label
	}
	if relatedURLItem.URL != nil {
		relatedURLMap["url"] = relatedURLItem.URL
	}

	return relatedURLMap
}

func dataSourceAPINotesFlattenReportedBy(result findingsv1.Reporter) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceAPINotesReportedByToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceAPINotesReportedByToMap(reportedByItem findingsv1.Reporter) (reportedByMap map[string]interface{}) {
	reportedByMap = map[string]interface{}{}

	if reportedByItem.ID != nil {
		reportedByMap["id"] = reportedByItem.ID
	}
	if reportedByItem.Title != nil {
		reportedByMap["title"] = reportedByItem.Title
	}
	if reportedByItem.URL != nil {
		reportedByMap["url"] = reportedByItem.URL
	}

	return reportedByMap
}

func dataSourceAPINotesFlattenFinding(result findingsv1.FindingType) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceAPINotesFindingToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceAPINotesFindingToMap(findingItem findingsv1.FindingType) (findingMap map[string]interface{}) {
	findingMap = map[string]interface{}{}

	if findingItem.Severity != nil {
		findingMap["severity"] = findingItem.Severity
	}
	if findingItem.NextSteps != nil {
		nextStepsList := []map[string]interface{}{}
		for _, nextStepsItem := range findingItem.NextSteps {
			nextStepsList = append(nextStepsList, dataSourceAPINotesFindingNextStepsToMap(nextStepsItem))
		}
		findingMap["next_steps"] = nextStepsList
	}

	return findingMap
}

func dataSourceAPINotesFindingNextStepsToMap(nextStepsItem findingsv1.RemediationStep) (nextStepsMap map[string]interface{}) {
	nextStepsMap = map[string]interface{}{}

	if nextStepsItem.Title != nil {
		nextStepsMap["title"] = nextStepsItem.Title
	}
	if nextStepsItem.URL != nil {
		nextStepsMap["url"] = nextStepsItem.URL
	}

	return nextStepsMap
}

func dataSourceAPINotesFlattenKpi(result findingsv1.KpiType) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceAPINotesKpiToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceAPINotesKpiToMap(kpiItem findingsv1.KpiType) (kpiMap map[string]interface{}) {
	kpiMap = map[string]interface{}{}

	if kpiItem.AggregationType != nil {
		kpiMap["aggregation_type"] = kpiItem.AggregationType
	}

	return kpiMap
}

func dataSourceAPINotesFlattenCard(result findingsv1.Card) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceAPINotesCardToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceAPINotesCardToMap(cardItem findingsv1.Card) (cardMap map[string]interface{}) {
	cardMap = map[string]interface{}{}

	if cardItem.Section != nil {
		cardMap["section"] = cardItem.Section
	}
	if cardItem.Title != nil {
		cardMap["title"] = cardItem.Title
	}
	if cardItem.Subtitle != nil {
		cardMap["subtitle"] = cardItem.Subtitle
	}
	if cardItem.Order != nil {
		cardMap["order"] = cardItem.Order
	}
	if cardItem.FindingNoteNames != nil {
		cardMap["finding_note_names"] = cardItem.FindingNoteNames
	}
	if cardItem.RequiresConfiguration != nil {
		cardMap["requires_configuration"] = cardItem.RequiresConfiguration
	}
	if cardItem.BadgeText != nil {
		cardMap["badge_text"] = cardItem.BadgeText
	}
	if cardItem.BadgeImage != nil {
		cardMap["badge_image"] = cardItem.BadgeImage
	}
	if cardItem.Elements != nil {
		elementsList := []map[string]interface{}{}
		for _, elementsItem := range cardItem.Elements {
			elementsList = append(elementsList, dataSourceAPINotesCardElementsToMap(elementsItem))
		}
		cardMap["elements"] = elementsList
	}

	return cardMap
}

func dataSourceAPINotesCardElementsToMap(elementsItem findingsv1.CardElementIntf) (elementsMap map[string]interface{}) {
	cardElementMap := map[string]interface{}{}

	switch v := elementsItem.(type) {
	case *findingsv1.CardElementNumericCardElement:
		cardElementMap["value_type"] = []map[string]interface{}{dataSourceAPINotesElementsValueTypeToMap(*v.ValueType)}
	case *findingsv1.CardElementBreakdownCardElement:
		cardElementMap["value_types"] = dataSourceAPINotesElementsValueTypesToMap(v.ValueTypes)
	case *findingsv1.CardElementTimeSeriesCardElement:
		cardElementMap["value_types"] = dataSourceAPINotesElementsValueTypesToMap(v.ValueTypes)
	}

	return cardElementMap
}

func dataSourceAPINotesElementsValueTypeToMap(valueTypeItem findingsv1.NumericCardElementValueType) (valueTypeMap map[string]interface{}) {
	valueTypeMap = map[string]interface{}{}

	if valueTypeItem.Kind != nil {
		valueTypeMap["kind"] = valueTypeItem.Kind
	}
	if valueTypeItem.KpiNoteName != nil {
		valueTypeMap["kpi_note_name"] = valueTypeItem.KpiNoteName
	}
	if valueTypeItem.Text != nil {
		valueTypeMap["text"] = valueTypeItem.Text
	}
	if valueTypeItem.FindingNoteNames != nil {
		valueTypeMap["finding_note_names"] = valueTypeItem.FindingNoteNames
	}

	return valueTypeMap
}

func dataSourceAPINotesElementsFindingCountValueTypeToMap(valueTypeItem findingsv1.ValueTypeFindingCountValueType) (valueTypeMap map[string]interface{}) {
	valueTypeMap = map[string]interface{}{}

	if valueTypeItem.Kind != nil {
		valueTypeMap["kind"] = valueTypeItem.Kind
	}
	if valueTypeItem.Text != nil {
		valueTypeMap["text"] = valueTypeItem.Text
	}
	if valueTypeItem.FindingNoteNames != nil {
		valueTypeMap["finding_note_names"] = valueTypeItem.FindingNoteNames
	}

	return valueTypeMap
}

func dataSourceAPINotesElementsKpiValueTypeToMap(valueTypeItem findingsv1.ValueTypeKpiValueType) (valueTypeMap map[string]interface{}) {
	valueTypeMap = map[string]interface{}{}

	if valueTypeItem.Kind != nil {
		valueTypeMap["kind"] = valueTypeItem.Kind
	}
	if valueTypeItem.Text != nil {
		valueTypeMap["text"] = valueTypeItem.Text
	}
	if valueTypeItem.KpiNoteName != nil {
		valueTypeMap["kpi_note_name"] = valueTypeItem.KpiNoteName
	}

	return valueTypeMap
}

func dataSourceAPINotesElementsValueTypesToMap(valueTypesItem []findingsv1.ValueTypeIntf) (valueTypesMap []map[string]interface{}) {
	valueTypesMap = []map[string]interface{}{}

	valueTypeMap := map[string]interface{}{}

	for _, valueType := range valueTypesItem {

		switch v := valueType.(type) {
		case *findingsv1.NumericCardElementValueType:
			valueTypeMap = dataSourceAPINotesElementsValueTypeToMap(*v)
		case *findingsv1.ValueTypeFindingCountValueType:
			valueTypeMap = dataSourceAPINotesElementsFindingCountValueTypeToMap(*v)
		case *findingsv1.ValueTypeKpiValueType:
			valueTypeMap = dataSourceAPINotesElementsKpiValueTypeToMap(*v)
		}

		valueTypesMap = append(valueTypesMap, valueTypeMap)
	}

	return valueTypesMap
}

func dataSourceAPINotesFlattenSection(result findingsv1.Section) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceAPINotesSectionToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceAPINotesSectionToMap(sectionItem findingsv1.Section) (sectionMap map[string]interface{}) {
	sectionMap = map[string]interface{}{}

	if sectionItem.Title != nil {
		sectionMap["title"] = sectionItem.Title
	}
	if sectionItem.Image != nil {
		sectionMap["image"] = sectionItem.Image
	}

	return sectionMap
}
