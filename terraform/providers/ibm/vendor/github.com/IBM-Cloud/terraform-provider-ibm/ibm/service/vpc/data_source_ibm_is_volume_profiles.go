// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isVolumeProfiles = "profiles"
)

func DataSourceIBMISVolumeProfiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISVolumeProfilesRead,

		Schema: map[string]*schema.Schema{

			isVolumeProfiles: {
				Type:        schema.TypeList,
				Description: "List of Volume profile maps",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"boot_capacity": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type for this profile field.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The value for this profile field.",
									},
									"default": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The default value for this profile field.",
									},
									"max": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The maximum value for this profile field.",
									},
									"min": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The minimum value for this profile field.",
									},
									"step": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The increment step value for this profile field.",
									},
									"values": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The permitted values for this profile field.",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
								},
							},
						},
						"capacity": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type for this profile field.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The value for this profile field.",
									},
									"default": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The default value for this profile field.",
									},
									"max": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The maximum value for this profile field.",
									},
									"min": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The minimum value for this profile field.",
									},
									"step": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The increment step value for this profile field.",
									},
									"values": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The permitted values for this profile field.",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
								},
							},
						},
						"family": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The product family this volume profile belongs to.The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this volume profile.",
						},
						"iops": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type for this profile field.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The value for this profile field.",
									},
									"default": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The default value for this profile field.",
									},
									"max": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The maximum value for this profile field.",
									},
									"min": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The minimum value for this profile field.",
									},
									"step": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The increment step value for this profile field.",
									},
									"values": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The permitted values for this profile field.",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
								},
							},
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique name for this volume profile.",
						},
						// defined_performance changes
						"adjustable_capacity_states": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type for this profile field.",
									},
									"values": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The attachment states that support adjustable capacity for a volume with this profile.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"adjustable_iops_states": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type for this profile field.",
									},
									"values": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The attachment states that support adjustable IOPS for a volume with this profile.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
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

func dataSourceIBMISVolumeProfilesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	err := volumeProfilesList(d, meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_volume_profiles", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	return nil
}

func volumeProfilesList(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	start := ""
	allrecs := []vpcv1.VolumeProfile{}
	for {
		listVolumeProfilesOptions := &vpcv1.ListVolumeProfilesOptions{}
		if start != "" {
			listVolumeProfilesOptions.Start = &start
		}
		availableProfiles, response, err := sess.ListVolumeProfiles(listVolumeProfilesOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Fetching Volume Profiles %s\n%s", err, response)
		}
		start = flex.GetNext(availableProfiles.Next)
		allrecs = append(allrecs, availableProfiles.Profiles...)
		if start == "" {
			break
		}
	}

	// listVolumeProfilesOptions := &vpcv1.ListVolumeProfilesOptions{}
	// availableProfiles, response, err := sess.ListVolumeProfiles(listVolumeProfilesOptions)
	// if err != nil {
	// 	return fmt.Errorf("[ERROR] Error Fetching Volume Profiles %s\n%s", err, response)
	// }
	profilesInfo := make([]map[string]interface{}, 0)
	for _, profile := range allrecs {
		modelMap, err := DataSourceIBMIsVolumeProfilesVolumeProfileToMap(&profile)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_volume_profiles", "read")
			return tfErr
		}
		profilesInfo = append(profilesInfo, modelMap)
	}
	d.SetId(dataSourceIBMISVolumeProfilesID(d))
	d.Set(isVolumeProfiles, profilesInfo)
	return nil
}

// dataSourceIBMISVolumeProfilesID returns a reasonable ID for a Volume Profile list.
func dataSourceIBMISVolumeProfilesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIBMIsVolumeProfilesVolumeProfileToMap(model *vpcv1.VolumeProfile) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.BootCapacity != nil {
		bootCapacityMap, err := DataSourceIBMIsVolumeProfilesVolumeProfileBootCapacityToMap(model.BootCapacity)
		if err != nil {
			return modelMap, err
		}
		modelMap["boot_capacity"] = []map[string]interface{}{bootCapacityMap}
	}
	if model.Capacity != nil {
		capacityMap, err := DataSourceIBMIsVolumeProfilesVolumeProfileCapacityToMap(model.Capacity)
		if err != nil {
			return modelMap, err
		}
		modelMap["capacity"] = []map[string]interface{}{capacityMap}
	}
	modelMap["family"] = *model.Family
	modelMap["href"] = *model.Href
	if model.Iops != nil {
		iopsMap, err := DataSourceIBMIsVolumeProfilesVolumeProfileIopsToMap(model.Iops)
		if err != nil {
			return modelMap, err
		}
		modelMap["iops"] = []map[string]interface{}{iopsMap}
	}
	modelMap["name"] = *model.Name
	if model.AdjustableCapacityStates != nil {
		adjustableCapacityStates, err := DataSourceIBMIsVolumeProfileVolumeProfileAdjustableCapacityStatesToMap(model.AdjustableCapacityStates)
		if err != nil {
			return modelMap, err
		}
		modelMap["adjustable_capacity_states"] = []map[string]interface{}{adjustableCapacityStates}
	}
	if model.AdjustableIopsStates != nil {
		adjustableIopsStates, err := DataSourceIBMIsVolumeProfileVolumeProfileAdjustableIopsStatesToMap(model.AdjustableIopsStates)
		if err != nil {
			return modelMap, err
		}
		modelMap["adjustable_iops_states"] = []map[string]interface{}{adjustableIopsStates}
	}
	return modelMap, nil
}

func DataSourceIBMIsVolumeProfilesVolumeProfileBootCapacityToMap(model vpcv1.VolumeProfileBootCapacityIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.VolumeProfileBootCapacityFixed); ok {
		return DataSourceIBMIsVolumeProfilesVolumeProfileBootCapacityFixedToMap(model.(*vpcv1.VolumeProfileBootCapacityFixed))
	} else if _, ok := model.(*vpcv1.VolumeProfileBootCapacityRange); ok {
		return DataSourceIBMIsVolumeProfilesVolumeProfileBootCapacityRangeToMap(model.(*vpcv1.VolumeProfileBootCapacityRange))
	} else if _, ok := model.(*vpcv1.VolumeProfileBootCapacityEnum); ok {
		return DataSourceIBMIsVolumeProfilesVolumeProfileBootCapacityEnumToMap(model.(*vpcv1.VolumeProfileBootCapacityEnum))
	} else if _, ok := model.(*vpcv1.VolumeProfileBootCapacityDependentRange); ok {
		return DataSourceIBMIsVolumeProfilesVolumeProfileBootCapacityDependentRangeToMap(model.(*vpcv1.VolumeProfileBootCapacityDependentRange))
	} else if _, ok := model.(*vpcv1.VolumeProfileBootCapacity); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.VolumeProfileBootCapacity)
		if model.Type != nil {
			modelMap["type"] = *model.Type
		}
		if model.Value != nil {
			modelMap["value"] = flex.IntValue(model.Value)
		}
		if model.Default != nil {
			modelMap["default"] = flex.IntValue(model.Default)
		}
		if model.Max != nil {
			modelMap["max"] = flex.IntValue(model.Max)
		}
		if model.Min != nil {
			modelMap["min"] = flex.IntValue(model.Min)
		}
		if model.Step != nil {
			modelMap["step"] = flex.IntValue(model.Step)
		}
		if model.Values != nil {
			modelMap["values"] = model.Values
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized vpcv1.VolumeProfileBootCapacityIntf subtype encountered")
	}
}

func DataSourceIBMIsVolumeProfilesVolumeProfileBootCapacityFixedToMap(model *vpcv1.VolumeProfileBootCapacityFixed) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = *model.Type
	modelMap["value"] = flex.IntValue(model.Value)
	return modelMap, nil
}

func DataSourceIBMIsVolumeProfilesVolumeProfileBootCapacityRangeToMap(model *vpcv1.VolumeProfileBootCapacityRange) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["default"] = flex.IntValue(model.Default)
	modelMap["max"] = flex.IntValue(model.Max)
	modelMap["min"] = flex.IntValue(model.Min)
	modelMap["step"] = flex.IntValue(model.Step)
	modelMap["type"] = *model.Type
	return modelMap, nil
}

func DataSourceIBMIsVolumeProfilesVolumeProfileBootCapacityEnumToMap(model *vpcv1.VolumeProfileBootCapacityEnum) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["default"] = flex.IntValue(model.Default)
	modelMap["type"] = *model.Type
	modelMap["values"] = model.Values
	return modelMap, nil
}

func DataSourceIBMIsVolumeProfilesVolumeProfileBootCapacityDependentRangeToMap(model *vpcv1.VolumeProfileBootCapacityDependentRange) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["max"] = flex.IntValue(model.Max)
	modelMap["min"] = flex.IntValue(model.Min)
	modelMap["step"] = flex.IntValue(model.Step)
	modelMap["type"] = *model.Type
	return modelMap, nil
}

func DataSourceIBMIsVolumeProfilesVolumeProfileCapacityToMap(model vpcv1.VolumeProfileCapacityIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.VolumeProfileCapacityFixed); ok {
		return DataSourceIBMIsVolumeProfilesVolumeProfileCapacityFixedToMap(model.(*vpcv1.VolumeProfileCapacityFixed))
	} else if _, ok := model.(*vpcv1.VolumeProfileCapacityRange); ok {
		return DataSourceIBMIsVolumeProfilesVolumeProfileCapacityRangeToMap(model.(*vpcv1.VolumeProfileCapacityRange))
	} else if _, ok := model.(*vpcv1.VolumeProfileCapacityEnum); ok {
		return DataSourceIBMIsVolumeProfilesVolumeProfileCapacityEnumToMap(model.(*vpcv1.VolumeProfileCapacityEnum))
	} else if _, ok := model.(*vpcv1.VolumeProfileCapacityDependentRange); ok {
		return DataSourceIBMIsVolumeProfilesVolumeProfileCapacityDependentRangeToMap(model.(*vpcv1.VolumeProfileCapacityDependentRange))
	} else if _, ok := model.(*vpcv1.VolumeProfileCapacity); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.VolumeProfileCapacity)
		if model.Type != nil {
			modelMap["type"] = *model.Type
		}
		if model.Value != nil {
			modelMap["value"] = flex.IntValue(model.Value)
		}
		if model.Default != nil {
			modelMap["default"] = flex.IntValue(model.Default)
		}
		if model.Max != nil {
			modelMap["max"] = flex.IntValue(model.Max)
		}
		if model.Min != nil {
			modelMap["min"] = flex.IntValue(model.Min)
		}
		if model.Step != nil {
			modelMap["step"] = flex.IntValue(model.Step)
		}
		if model.Values != nil {
			modelMap["values"] = model.Values
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized vpcv1.VolumeProfileCapacityIntf subtype encountered")
	}
}

func DataSourceIBMIsVolumeProfilesVolumeProfileCapacityFixedToMap(model *vpcv1.VolumeProfileCapacityFixed) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = *model.Type
	modelMap["value"] = flex.IntValue(model.Value)
	return modelMap, nil
}

func DataSourceIBMIsVolumeProfilesVolumeProfileCapacityRangeToMap(model *vpcv1.VolumeProfileCapacityRange) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["default"] = flex.IntValue(model.Default)
	modelMap["max"] = flex.IntValue(model.Max)
	modelMap["min"] = flex.IntValue(model.Min)
	modelMap["step"] = flex.IntValue(model.Step)
	modelMap["type"] = *model.Type
	return modelMap, nil
}

func DataSourceIBMIsVolumeProfilesVolumeProfileCapacityEnumToMap(model *vpcv1.VolumeProfileCapacityEnum) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["default"] = flex.IntValue(model.Default)
	modelMap["type"] = *model.Type
	modelMap["values"] = model.Values
	return modelMap, nil
}

func DataSourceIBMIsVolumeProfilesVolumeProfileCapacityDependentRangeToMap(model *vpcv1.VolumeProfileCapacityDependentRange) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["max"] = flex.IntValue(model.Max)
	modelMap["min"] = flex.IntValue(model.Min)
	modelMap["step"] = flex.IntValue(model.Step)
	modelMap["type"] = *model.Type
	return modelMap, nil
}

func DataSourceIBMIsVolumeProfilesVolumeProfileIopsToMap(model vpcv1.VolumeProfileIopsIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.VolumeProfileIopsFixed); ok {
		return DataSourceIBMIsVolumeProfilesVolumeProfileIopsFixedToMap(model.(*vpcv1.VolumeProfileIopsFixed))
	} else if _, ok := model.(*vpcv1.VolumeProfileIopsRange); ok {
		return DataSourceIBMIsVolumeProfilesVolumeProfileIopsRangeToMap(model.(*vpcv1.VolumeProfileIopsRange))
	} else if _, ok := model.(*vpcv1.VolumeProfileIopsEnum); ok {
		return DataSourceIBMIsVolumeProfilesVolumeProfileIopsEnumToMap(model.(*vpcv1.VolumeProfileIopsEnum))
	} else if _, ok := model.(*vpcv1.VolumeProfileIopsDependentRange); ok {
		return DataSourceIBMIsVolumeProfilesVolumeProfileIopsDependentRangeToMap(model.(*vpcv1.VolumeProfileIopsDependentRange))
	} else if _, ok := model.(*vpcv1.VolumeProfileIops); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.VolumeProfileIops)
		if model.Type != nil {
			modelMap["type"] = *model.Type
		}
		if model.Value != nil {
			modelMap["value"] = flex.IntValue(model.Value)
		}
		if model.Default != nil {
			modelMap["default"] = flex.IntValue(model.Default)
		}
		if model.Max != nil {
			modelMap["max"] = flex.IntValue(model.Max)
		}
		if model.Min != nil {
			modelMap["min"] = flex.IntValue(model.Min)
		}
		if model.Step != nil {
			modelMap["step"] = flex.IntValue(model.Step)
		}
		if model.Values != nil {
			modelMap["values"] = model.Values
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized vpcv1.VolumeProfileIopsIntf subtype encountered")
	}
}

func DataSourceIBMIsVolumeProfilesVolumeProfileIopsFixedToMap(model *vpcv1.VolumeProfileIopsFixed) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = *model.Type
	modelMap["value"] = flex.IntValue(model.Value)
	return modelMap, nil
}

func DataSourceIBMIsVolumeProfilesVolumeProfileIopsRangeToMap(model *vpcv1.VolumeProfileIopsRange) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["default"] = flex.IntValue(model.Default)
	modelMap["max"] = flex.IntValue(model.Max)
	modelMap["min"] = flex.IntValue(model.Min)
	modelMap["step"] = flex.IntValue(model.Step)
	modelMap["type"] = *model.Type
	return modelMap, nil
}

func DataSourceIBMIsVolumeProfilesVolumeProfileIopsEnumToMap(model *vpcv1.VolumeProfileIopsEnum) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["default"] = flex.IntValue(model.Default)
	modelMap["type"] = *model.Type
	modelMap["values"] = model.Values
	return modelMap, nil
}

func DataSourceIBMIsVolumeProfilesVolumeProfileIopsDependentRangeToMap(model *vpcv1.VolumeProfileIopsDependentRange) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["max"] = flex.IntValue(model.Max)
	modelMap["min"] = flex.IntValue(model.Min)
	modelMap["step"] = flex.IntValue(model.Step)
	modelMap["type"] = *model.Type
	return modelMap, nil
}

func printfull(response interface{}) string {
	output, err := json.MarshalIndent(response, "", "    ")
	if err == nil {
		return fmt.Sprintf("%+v\n", string(output))
	}
	return fmt.Sprintf("Error : %#v", response)
}
