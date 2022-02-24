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

func dataSourceIBMSccSiNote() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMSccSiNoteRead,

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
			"note_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Second part of note `name`: providers/{provider_id}/notes/{note_id}.",
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
	}
}

func dataSourceIBMSccSiNoteRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	getNoteOptions := &findingsv1.GetNoteOptions{}

	getNoteOptions.SetProviderID(d.Get("provider_id").(string))
	getNoteOptions.SetNoteID(d.Get("note_id").(string))

	apiNote, response, err := findingsClient.GetNoteWithContext(context, getNoteOptions)
	if err != nil {
		log.Printf("[DEBUG] GetNoteWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetNoteWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", *findingsClient.AccountID, *getNoteOptions.ProviderID, *getNoteOptions.NoteID))
	if err = d.Set("short_description", apiNote.ShortDescription); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting short_description: %s", err))
	}
	if err = d.Set("long_description", apiNote.LongDescription); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting long_description: %s", err))
	}
	if err = d.Set("kind", apiNote.Kind); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting kind: %s", err))
	}

	if apiNote.RelatedURL != nil {
		err = d.Set("related_url", dataSourceAPINoteFlattenRelatedURL(apiNote.RelatedURL))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting related_url %s", err))
		}
	}
	if err = d.Set("create_time", dateTimeToString(apiNote.CreateTime)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting create_time: %s", err))
	}
	if err = d.Set("update_time", dateTimeToString(apiNote.UpdateTime)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting update_time: %s", err))
	}
	if err = d.Set("shared", apiNote.Shared); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting shared: %s", err))
	}

	if apiNote.ReportedBy != nil {
		err = d.Set("reported_by", dataSourceAPINoteFlattenReportedBy(*apiNote.ReportedBy))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting reported_by %s", err))
		}
	}

	if apiNote.Finding != nil {
		err = d.Set("finding", dataSourceAPINoteFlattenFinding(*apiNote.Finding))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting finding %s", err))
		}
	}

	if apiNote.Kpi != nil {
		err = d.Set("kpi", dataSourceAPINoteFlattenKpi(*apiNote.Kpi))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting kpi %s", err))
		}
	}

	if apiNote.Card != nil {
		err = d.Set("card", dataSourceAPINoteFlattenCard(*apiNote.Card))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting card %s", err))
		}
	}

	if apiNote.Section != nil {
		err = d.Set("section", dataSourceAPINoteFlattenSection(*apiNote.Section))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting section %s", err))
		}
	}

	return nil
}

func dataSourceAPINoteFlattenRelatedURL(result []findingsv1.APINoteRelatedURL) (relatedURL []map[string]interface{}) {
	for _, relatedURLItem := range result {
		relatedURL = append(relatedURL, dataSourceAPINoteRelatedURLToMap(relatedURLItem))
	}

	return relatedURL
}

func dataSourceAPINoteRelatedURLToMap(relatedURLItem findingsv1.APINoteRelatedURL) (relatedURLMap map[string]interface{}) {
	relatedURLMap = map[string]interface{}{}

	if relatedURLItem.Label != nil {
		relatedURLMap["label"] = relatedURLItem.Label
	}
	if relatedURLItem.URL != nil {
		relatedURLMap["url"] = relatedURLItem.URL
	}

	return relatedURLMap
}

func dataSourceAPINoteFlattenReportedBy(result findingsv1.Reporter) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceAPINoteReportedByToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceAPINoteReportedByToMap(reportedByItem findingsv1.Reporter) (reportedByMap map[string]interface{}) {
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

func dataSourceAPINoteFlattenFinding(result findingsv1.FindingType) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceAPINoteFindingToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceAPINoteFindingToMap(findingItem findingsv1.FindingType) (findingMap map[string]interface{}) {
	findingMap = map[string]interface{}{}

	if findingItem.Severity != nil {
		findingMap["severity"] = findingItem.Severity
	}
	if findingItem.NextSteps != nil {
		nextStepsList := []map[string]interface{}{}
		for _, nextStepsItem := range findingItem.NextSteps {
			nextStepsList = append(nextStepsList, dataSourceAPINoteFindingNextStepsToMap(nextStepsItem))
		}
		findingMap["next_steps"] = nextStepsList
	}

	return findingMap
}

func dataSourceAPINoteFindingNextStepsToMap(nextStepsItem findingsv1.RemediationStep) (nextStepsMap map[string]interface{}) {
	nextStepsMap = map[string]interface{}{}

	if nextStepsItem.Title != nil {
		nextStepsMap["title"] = nextStepsItem.Title
	}
	if nextStepsItem.URL != nil {
		nextStepsMap["url"] = nextStepsItem.URL
	}

	return nextStepsMap
}

func dataSourceAPINoteFlattenKpi(result findingsv1.KpiType) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceAPINoteKpiToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceAPINoteKpiToMap(kpiItem findingsv1.KpiType) (kpiMap map[string]interface{}) {
	kpiMap = map[string]interface{}{}

	if kpiItem.AggregationType != nil {
		kpiMap["aggregation_type"] = kpiItem.AggregationType
	}

	return kpiMap
}

func dataSourceAPINoteFlattenCard(result findingsv1.Card) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceAPINoteCardToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceAPINoteCardToMap(cardItem findingsv1.Card) (cardMap map[string]interface{}) {
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
			elementsList = append(elementsList, dataSourceAPINoteCardElementsToMap(elementsItem))
		}
		cardMap["elements"] = elementsList
	}

	return cardMap
}

func dataSourceAPINoteCardElementsToMap(elementsItem findingsv1.CardElementIntf) (elementsMap map[string]interface{}) {
	cardElementMap := map[string]interface{}{}

	switch v := elementsItem.(type) {
	case *findingsv1.CardElementNumericCardElement:
		cardElementMap["value_type"] = dataSourceAPINoteElementsValueTypeToMap(*v.ValueType)
	case *findingsv1.CardElementBreakdownCardElement:
		cardElementMap["value_types"] = dataSourceAPINoteElementsValueTypesToMap(v.ValueTypes)
	case *findingsv1.CardElementTimeSeriesCardElement:
		cardElementMap["value_types"] = dataSourceAPINoteElementsValueTypesToMap(v.ValueTypes)
	}

	return cardElementMap
}

func dataSourceAPINoteElementsValueTypeToMap(valueTypeItem findingsv1.NumericCardElementValueType) (valueTypeMap map[string]interface{}) {
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

func dataSourceAPINoteElementsFindingCountValueTypeToMap(valueTypeItem findingsv1.ValueTypeFindingCountValueType) (valueTypeMap map[string]interface{}) {
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

func dataSourceAPINoteElementsKpiValueTypeToMap(valueTypeItem findingsv1.ValueTypeKpiValueType) (valueTypeMap map[string]interface{}) {
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

func dataSourceAPINoteElementsValueTypesToMap(valueTypesItem []findingsv1.ValueTypeIntf) (valueTypesMap []map[string]interface{}) {
	valueTypesMap = []map[string]interface{}{}

	valueTypeMap := map[string]interface{}{}

	for _, valueType := range valueTypesItem {

		switch v := valueType.(type) {
		case *findingsv1.NumericCardElementValueType:
			valueTypeMap = dataSourceAPINoteElementsValueTypeToMap(*v)
		case *findingsv1.ValueTypeFindingCountValueType:
			valueTypeMap = dataSourceAPINoteElementsFindingCountValueTypeToMap(*v)
		case *findingsv1.ValueTypeKpiValueType:
			valueTypeMap = dataSourceAPINoteElementsKpiValueTypeToMap(*v)
		}

		valueTypesMap = append(valueTypesMap, valueTypeMap)
	}

	return valueTypesMap
}

func dataSourceAPINoteFlattenSection(result findingsv1.Section) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceAPINoteSectionToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceAPINoteSectionToMap(sectionItem findingsv1.Section) (sectionMap map[string]interface{}) {
	sectionMap = map[string]interface{}{}

	if sectionItem.Title != nil {
		sectionMap["title"] = sectionItem.Title
	}
	if sectionItem.Image != nil {
		sectionMap["image"] = sectionItem.Image
	}

	return sectionMap
}
