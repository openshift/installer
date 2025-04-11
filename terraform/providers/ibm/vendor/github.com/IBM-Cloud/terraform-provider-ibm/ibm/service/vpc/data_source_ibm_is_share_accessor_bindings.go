// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.90.0-5aad763d-20240506-203857
 */

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsShareAccessorBindings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsShareAccessorBindingsRead,

		Schema: map[string]*schema.Schema{
			"share": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The file share identifier.",
			},
			"accessor_bindings": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of share accessor bindings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accessor": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The accessor for this share accessor binding.The resources supported by this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this file share.",
									},
									"deleted": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this file share.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this file share.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this share. The name is unique across all shares in the region.",
									},
									"remote": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates that the resource associated with this referenceis remote and therefore may not be directly retrievable.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"account": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "If present, this property indicates that the referenced resource is remote to thisaccount, and identifies the owning account.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The unique identifier for this account.",
															},
															"resource_type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The resource type.",
															},
														},
													},
												},
												"region": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "If present, this property indicates that the referenced resource is remote to thisregion, and identifies the native region.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"href": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The URL for this region.",
															},
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The globally unique name for this region.",
															},
														},
													},
												},
											},
										},
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the share accessor binding was created.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this share accessor binding.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this share accessor binding.",
						},
						"lifecycle_state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of the file share accessor binding.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsShareAccessorBindingsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_share_accessor_bindings", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	listShareAccessorBindingsOptions := &vpcv1.ListShareAccessorBindingsOptions{}

	listShareAccessorBindingsOptions.SetID(d.Get("share").(string))

	var pager *vpcv1.ShareAccessorBindingsPager
	pager, err = vpcClient.NewShareAccessorBindingsPager(listShareAccessorBindingsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_share_accessor_bindings", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	allItems, err := pager.GetAll()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ShareAccessorBindingsPager.GetAll() failed %s", err), "(Data) ibm_is_share_accessor_bindings", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIBMIsShareAccessorBindingsID(d))

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelMap, err := DataSourceIBMIsShareAccessorBindingsShareAccessorBindingToMap(&modelItem)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_share_accessor_bindings", "read")
			return tfErr.GetDiag()
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("accessor_bindings", mapSlice); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting accessor_bindings %s", err), "(Data) ibm_is_share_accessor_bindings", "read")
		return tfErr.GetDiag()
	}

	return nil
}

// dataSourceIBMIsShareAccessorBindingsID returns a reasonable ID for the list.
func dataSourceIBMIsShareAccessorBindingsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIBMIsShareAccessorBindingsShareAccessorBindingToMap(model *vpcv1.ShareAccessorBinding) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	accessorMap, err := DataSourceIBMIsShareAccessorBindingsShareAccessorBindingAccessorToMap(model.Accessor)
	if err != nil {
		return modelMap, err
	}
	modelMap["accessor"] = []map[string]interface{}{accessorMap}
	modelMap["created_at"] = model.CreatedAt.String()
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["lifecycle_state"] = *model.LifecycleState
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMIsShareAccessorBindingsShareAccessorBindingAccessorToMap(model vpcv1.ShareAccessorBindingAccessorIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.ShareAccessorBindingAccessorShareReference); ok {
		return DataSourceIBMIsShareAccessorBindingsShareAccessorBindingAccessorShareReferenceToMap(model.(*vpcv1.ShareAccessorBindingAccessorShareReference))
	} else if _, ok := model.(*vpcv1.ShareAccessorBindingAccessorWatsonxMachineLearningReference); ok {
		return DataSourceIBMIsShareAccessorBindingsShareAccessorBindingAccessorWatsonxMachineLearningReferenceToMap(model.(*vpcv1.ShareAccessorBindingAccessorWatsonxMachineLearningReference))
	} else if _, ok := model.(*vpcv1.ShareAccessorBindingAccessor); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.ShareAccessorBindingAccessor)
		if model.CRN != nil {
			modelMap["crn"] = *model.CRN
		}
		if model.Deleted != nil {
			deletedMap, err := DataSourceIBMIsShareAccessorBindingsShareReferenceDeletedToMap(model.Deleted)
			if err != nil {
				return modelMap, err
			}
			modelMap["deleted"] = []map[string]interface{}{deletedMap}
		}
		if model.Href != nil {
			modelMap["href"] = *model.Href
		}
		if model.ID != nil {
			modelMap["id"] = *model.ID
		}
		if model.Name != nil {
			modelMap["name"] = *model.Name
		}
		if model.Remote != nil {
			remoteMap, err := DataSourceIBMIsShareAccessorBindingsShareRemoteToMap(model.Remote)
			if err != nil {
				return modelMap, err
			}
			modelMap["remote"] = []map[string]interface{}{remoteMap}
		}
		if model.ResourceType != nil {
			modelMap["resource_type"] = *model.ResourceType
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized vpcv1.ShareAccessorBindingAccessorIntf subtype encountered")
	}
}

func DataSourceIBMIsShareAccessorBindingsShareReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func DataSourceIBMIsShareAccessorBindingsShareRemoteToMap(model *vpcv1.ShareRemote) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Account != nil {
		accountMap, err := DataSourceIBMIsShareAccessorBindingsAccountReferenceToMap(model.Account)
		if err != nil {
			return modelMap, err
		}
		modelMap["account"] = []map[string]interface{}{accountMap}
	}
	if model.Region != nil {
		regionMap, err := DataSourceIBMIsShareAccessorBindingsRegionReferenceToMap(model.Region)
		if err != nil {
			return modelMap, err
		}
		modelMap["region"] = []map[string]interface{}{regionMap}
	}
	return modelMap, nil
}

func DataSourceIBMIsShareAccessorBindingsAccountReferenceToMap(model *vpcv1.AccountReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMIsShareAccessorBindingsRegionReferenceToMap(model *vpcv1.RegionReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func DataSourceIBMIsShareAccessorBindingsShareAccessorBindingAccessorShareReferenceToMap(model *vpcv1.ShareAccessorBindingAccessorShareReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsShareAccessorBindingsShareReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	if model.Remote != nil {
		remoteMap, err := DataSourceIBMIsShareAccessorBindingsShareRemoteToMap(model.Remote)
		if err != nil {
			return modelMap, err
		}
		modelMap["remote"] = []map[string]interface{}{remoteMap}
	}
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMIsShareAccessorBindingsShareAccessorBindingAccessorWatsonxMachineLearningReferenceToMap(model *vpcv1.ShareAccessorBindingAccessorWatsonxMachineLearningReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}
