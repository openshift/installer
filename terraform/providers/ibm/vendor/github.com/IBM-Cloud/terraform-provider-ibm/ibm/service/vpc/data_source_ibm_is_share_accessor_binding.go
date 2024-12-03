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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsShareAccessorBinding() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsShareAccessorBindingRead,

		Schema: map[string]*schema.Schema{
			"share": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The file share identifier.",
			},
			"accessor_binding": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The file share accessor binding identifier.",
			},
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
	}
}

func dataSourceIBMIsShareAccessorBindingRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		// Error is coming from SDK client, so it doesn't need to be discriminated.
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_share_accessor_binding", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getShareAccessorBindingOptions := &vpcv1.GetShareAccessorBindingOptions{}

	getShareAccessorBindingOptions.SetShareID(d.Get("share").(string))
	getShareAccessorBindingOptions.SetID(d.Get("accessor_binding").(string))

	shareAccessorBinding, _, err := vpcClient.GetShareAccessorBindingWithContext(context, getShareAccessorBindingOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetShareAccessorBindingWithContext failed: %s", err.Error()), "(Data) ibm_is_share_accessor_binding", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*shareAccessorBinding.ID)

	accessor := []map[string]interface{}{}
	if shareAccessorBinding.Accessor != nil {
		modelMap, err := DataSourceIBMIsShareAccessorBindingShareAccessorBindingAccessorToMap(shareAccessorBinding.Accessor)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_share_accessor_binding", "read", "accessor-to-map").GetDiag()
		}
		accessor = append(accessor, modelMap)
	}
	if err = d.Set("accessor", accessor); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting accessor: %s", err), "(Data) ibm_is_share_accessor_binding", "read", "set-accessor").GetDiag()
	}

	if err = d.Set("created_at", flex.DateTimeToString(shareAccessorBinding.CreatedAt)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting created_at: %s", err), "(Data) ibm_is_share_accessor_binding", "read", "set-created_at").GetDiag()
	}

	if err = d.Set("href", shareAccessorBinding.Href); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_is_share_accessor_binding", "read", "set-href").GetDiag()
	}

	if err = d.Set("lifecycle_state", shareAccessorBinding.LifecycleState); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting lifecycle_state: %s", err), "(Data) ibm_is_share_accessor_binding", "read", "set-lifecycle_state").GetDiag()
	}

	if err = d.Set("resource_type", shareAccessorBinding.ResourceType); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting resource_type: %s", err), "(Data) ibm_is_share_accessor_binding", "read", "set-resource_type").GetDiag()
	}

	return nil
}

func DataSourceIBMIsShareAccessorBindingShareAccessorBindingAccessorToMap(model vpcv1.ShareAccessorBindingAccessorIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.ShareAccessorBindingAccessorShareReference); ok {
		return DataSourceIBMIsShareAccessorBindingShareAccessorBindingAccessorShareReferenceToMap(model.(*vpcv1.ShareAccessorBindingAccessorShareReference))
	} else if _, ok := model.(*vpcv1.ShareAccessorBindingAccessorWatsonxMachineLearningReference); ok {
		return DataSourceIBMIsShareAccessorBindingShareAccessorBindingAccessorWatsonxMachineLearningReferenceToMap(model.(*vpcv1.ShareAccessorBindingAccessorWatsonxMachineLearningReference))
	} else if _, ok := model.(*vpcv1.ShareAccessorBindingAccessor); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.ShareAccessorBindingAccessor)
		if model.CRN != nil {
			modelMap["crn"] = *model.CRN
		}
		if model.Deleted != nil {
			deletedMap, err := DataSourceIBMIsShareAccessorBindingShareReferenceDeletedToMap(model.Deleted)
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
			remoteMap, err := DataSourceIBMIsShareAccessorBindingShareRemoteToMap(model.Remote)
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

func DataSourceIBMIsShareAccessorBindingShareReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func DataSourceIBMIsShareAccessorBindingShareRemoteToMap(model *vpcv1.ShareRemote) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Account != nil {
		accountMap, err := DataSourceIBMIsShareAccessorBindingAccountReferenceToMap(model.Account)
		if err != nil {
			return modelMap, err
		}
		modelMap["account"] = []map[string]interface{}{accountMap}
	}
	if model.Region != nil {
		regionMap, err := DataSourceIBMIsShareAccessorBindingRegionReferenceToMap(model.Region)
		if err != nil {
			return modelMap, err
		}
		modelMap["region"] = []map[string]interface{}{regionMap}
	}
	return modelMap, nil
}

func DataSourceIBMIsShareAccessorBindingAccountReferenceToMap(model *vpcv1.AccountReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMIsShareAccessorBindingRegionReferenceToMap(model *vpcv1.RegionReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func DataSourceIBMIsShareAccessorBindingShareAccessorBindingAccessorShareReferenceToMap(model *vpcv1.ShareAccessorBindingAccessorShareReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsShareAccessorBindingShareReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	if model.Remote != nil {
		remoteMap, err := DataSourceIBMIsShareAccessorBindingShareRemoteToMap(model.Remote)
		if err != nil {
			return modelMap, err
		}
		modelMap["remote"] = []map[string]interface{}{remoteMap}
	}
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMIsShareAccessorBindingShareAccessorBindingAccessorWatsonxMachineLearningReferenceToMap(model *vpcv1.ShareAccessorBindingAccessorWatsonxMachineLearningReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}
