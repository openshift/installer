// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func ResourceIBMIsInstanceClusterNetworkAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsInstanceClusterNetworkAttachmentCreate,
		ReadContext:   resourceIBMIsInstanceClusterNetworkAttachmentRead,
		UpdateContext: resourceIBMIsInstanceClusterNetworkAttachmentUpdate,
		DeleteContext: resourceIBMIsInstanceClusterNetworkAttachmentDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"instance_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_instance_cluster_network_attachment", "instance_id"),
				Description:  "The virtual server instance identifier.",
			},
			"before": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "The instance cluster network attachment that is immediately before. If absent, this is thelast instance cluster network attachment.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this instance cluster network attachment.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The unique identifier for this instance cluster network attachment.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this instance cluster network attachment. The name is unique across all network attachments for the instance.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"cluster_network_interface": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "The cluster network interface for this instance cluster network attachment.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
							Description: "The URL for this cluster network interface.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this cluster network interface.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name for this cluster network interface. The name is unique across all interfaces in the cluster network.",
						},
						"primary_ip": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "The primary IP for this cluster network interface.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The IP address.If the address is pending allocation, the value will be `0.0.0.0`.This property may [expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) to support IPv6 addresses in the future.",
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
										Required:    true,
										Description: "The URL for this cluster network subnet reserved IP.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The unique identifier for this cluster network subnet reserved IP.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The name for this cluster network subnet reserved IP. The name is unique across all reserved IPs in a cluster network subnet.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						"subnet": &schema.Schema{
							Type:     schema.TypeList,
							MinItems: 1,
							MaxItems: 1,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
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
										Description: "The URL for this cluster network subnet.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The unique identifier for this cluster network subnet.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this cluster network subnet. The name is unique across all cluster network subnets in the cluster network.",
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
				},
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_instance_cluster_network_attachment", "name"),
				Description:  "The name for this instance cluster network attachment. The name is unique across all network attachments for the instance.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this instance cluster network attachment.",
			},
			"lifecycle_reasons": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reasons for the current `lifecycle_state` (if any).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A reason code for this lifecycle state:- `internal_error`: internal error (contact IBM support)- `resource_suspended_by_provider`: The resource has been suspended (contact IBM  support)The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
						},
						"message": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the reason for this lifecycle state.",
						},
						"more_info": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Link to documentation about the reason for this lifecycle state.",
						},
					},
				},
			},
			"lifecycle_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the instance cluster network attachment.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"instance_cluster_network_attachment_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this instance cluster network attachment.",
			},
		},
	}
}

func ResourceIBMIsInstanceClusterNetworkAttachmentValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "instance_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[-0-9a-z_]+$`,
			MinValueLength:             1,
			MaxValueLength:             64,
		},
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_instance_cluster_network_attachment", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMIsInstanceClusterNetworkAttachmentCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_cluster_network_attachment", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createClusterNetworkAttachmentOptions := &vpcv1.CreateClusterNetworkAttachmentOptions{}

	createClusterNetworkAttachmentOptions.SetInstanceID(d.Get("instance_id").(string))
	clusterNetworkInterfaceModel, err := ResourceIBMIsInstanceClusterNetworkAttachmentMapToInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterface(d.Get("cluster_network_interface.0").(map[string]interface{}))
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_cluster_network_attachment", "create", "parse-cluster_network_interface").GetDiag()
	}
	createClusterNetworkAttachmentOptions.SetClusterNetworkInterface(clusterNetworkInterfaceModel)
	if _, ok := d.GetOk("before"); ok {
		beforeModel, err := ResourceIBMIsInstanceClusterNetworkAttachmentMapToInstanceClusterNetworkAttachmentBeforePrototype(d.Get("before.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_cluster_network_attachment", "create", "parse-before").GetDiag()
		}
		createClusterNetworkAttachmentOptions.SetBefore(beforeModel)
	}
	if _, ok := d.GetOk("name"); ok {
		createClusterNetworkAttachmentOptions.SetName(d.Get("name").(string))
	}

	instanceClusterNetworkAttachment, _, err := vpcClient.CreateClusterNetworkAttachmentWithContext(context, createClusterNetworkAttachmentOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateClusterNetworkAttachmentWithContext failed: %s", err.Error()), "ibm_is_instance_cluster_network_attachment", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *createClusterNetworkAttachmentOptions.InstanceID, *instanceClusterNetworkAttachment.ID))

	return resourceIBMIsInstanceClusterNetworkAttachmentRead(context, d, meta)
}

func resourceIBMIsInstanceClusterNetworkAttachmentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_cluster_network_attachment", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getInstanceClusterNetworkAttachmentOptions := &vpcv1.GetInstanceClusterNetworkAttachmentOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_cluster_network_attachment", "read", "sep-id-parts").GetDiag()
	}

	getInstanceClusterNetworkAttachmentOptions.SetInstanceID(parts[0])
	getInstanceClusterNetworkAttachmentOptions.SetID(parts[1])

	instanceClusterNetworkAttachment, response, err := vpcClient.GetInstanceClusterNetworkAttachmentWithContext(context, getInstanceClusterNetworkAttachmentOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetInstanceClusterNetworkAttachmentWithContext failed: %s", err.Error()), "ibm_is_instance_cluster_network_attachment", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if !core.IsNil(instanceClusterNetworkAttachment.Before) {
		beforeMap, err := ResourceIBMIsInstanceClusterNetworkAttachmentInstanceClusterNetworkAttachmentBeforeToMap(instanceClusterNetworkAttachment.Before)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_cluster_network_attachment", "read", "before-to-map").GetDiag()
		}
		if err = d.Set("before", []map[string]interface{}{beforeMap}); err != nil {
			err = fmt.Errorf("Error setting before: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_cluster_network_attachment", "read", "set-before").GetDiag()
		}
	}
	clusterNetworkInterfaceMap, err := ResourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkInterfaceReferenceToMap(instanceClusterNetworkAttachment.ClusterNetworkInterface)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_cluster_network_attachment", "read", "cluster_network_interface-to-map").GetDiag()
	}
	if err = d.Set("cluster_network_interface", []map[string]interface{}{clusterNetworkInterfaceMap}); err != nil {
		err = fmt.Errorf("Error setting cluster_network_interface: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_cluster_network_attachment", "read", "set-cluster_network_interface").GetDiag()
	}
	if !core.IsNil(instanceClusterNetworkAttachment.Name) {
		if err = d.Set("name", instanceClusterNetworkAttachment.Name); err != nil {
			err = fmt.Errorf("Error setting name: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_cluster_network_attachment", "read", "set-name").GetDiag()
		}
	}
	if err = d.Set("href", instanceClusterNetworkAttachment.Href); err != nil {
		err = fmt.Errorf("Error setting href: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_cluster_network_attachment", "read", "set-href").GetDiag()
	}
	lifecycleReasons := []map[string]interface{}{}
	for _, lifecycleReasonsItem := range instanceClusterNetworkAttachment.LifecycleReasons {
		lifecycleReasonsItemMap, err := ResourceIBMIsInstanceClusterNetworkAttachmentInstanceClusterNetworkAttachmentLifecycleReasonToMap(&lifecycleReasonsItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_cluster_network_attachment", "read", "lifecycle_reasons-to-map").GetDiag()
		}
		lifecycleReasons = append(lifecycleReasons, lifecycleReasonsItemMap)
	}
	if err = d.Set("lifecycle_reasons", lifecycleReasons); err != nil {
		err = fmt.Errorf("Error setting lifecycle_reasons: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_cluster_network_attachment", "read", "set-lifecycle_reasons").GetDiag()
	}
	if err = d.Set("lifecycle_state", instanceClusterNetworkAttachment.LifecycleState); err != nil {
		err = fmt.Errorf("Error setting lifecycle_state: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_cluster_network_attachment", "read", "set-lifecycle_state").GetDiag()
	}
	if err = d.Set("resource_type", instanceClusterNetworkAttachment.ResourceType); err != nil {
		err = fmt.Errorf("Error setting resource_type: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_cluster_network_attachment", "read", "set-resource_type").GetDiag()
	}
	if err = d.Set("instance_cluster_network_attachment_id", instanceClusterNetworkAttachment.ID); err != nil {
		err = fmt.Errorf("Error setting instance_cluster_network_attachment_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_cluster_network_attachment", "read", "set-instance_cluster_network_attachment_id").GetDiag()
	}

	return nil
}

func resourceIBMIsInstanceClusterNetworkAttachmentUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_cluster_network_attachment", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateInstanceClusterNetworkAttachmentOptions := &vpcv1.UpdateInstanceClusterNetworkAttachmentOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_cluster_network_attachment", "update", "sep-id-parts").GetDiag()
	}

	updateInstanceClusterNetworkAttachmentOptions.SetInstanceID(parts[0])
	updateInstanceClusterNetworkAttachmentOptions.SetID(parts[1])

	hasChange := false

	patchVals := &vpcv1.InstanceClusterNetworkAttachmentPatch{}
	if d.HasChange("instance_id") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "instance_id")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_is_instance_cluster_network_attachment", "update", "instance_id-forces-new").GetDiag()
	}
	if d.HasChange("name") {
		newName := d.Get("name").(string)
		patchVals.Name = &newName
		hasChange = true
	}

	if hasChange {
		// Fields with `nil` values are omitted from the generic map,
		// so we need to re-add them to support removing arguments
		// in merge-patch operations sent to the service.
		updateInstanceClusterNetworkAttachmentOptions.InstanceClusterNetworkAttachmentPatch = ResourceIBMIsInstanceClusterNetworkAttachmentInstanceClusterNetworkAttachmentPatchAsPatch(patchVals, d)

		_, _, err = vpcClient.UpdateInstanceClusterNetworkAttachmentWithContext(context, updateInstanceClusterNetworkAttachmentOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateInstanceClusterNetworkAttachmentWithContext failed: %s", err.Error()), "ibm_is_instance_cluster_network_attachment", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIBMIsInstanceClusterNetworkAttachmentRead(context, d, meta)
}

func resourceIBMIsInstanceClusterNetworkAttachmentDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_cluster_network_attachment", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteInstanceClusterNetworkAttachmentOptions := &vpcv1.DeleteInstanceClusterNetworkAttachmentOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_cluster_network_attachment", "delete", "sep-id-parts").GetDiag()
	}

	deleteInstanceClusterNetworkAttachmentOptions.SetInstanceID(parts[0])
	deleteInstanceClusterNetworkAttachmentOptions.SetID(parts[1])

	_, _, err = vpcClient.DeleteInstanceClusterNetworkAttachmentWithContext(context, deleteInstanceClusterNetworkAttachmentOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteInstanceClusterNetworkAttachmentWithContext failed: %s", err.Error()), "ibm_is_instance_cluster_network_attachment", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIBMIsInstanceClusterNetworkAttachmentMapToInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterface(modelMap map[string]interface{}) (vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceIntf, error) {
	model := &vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterface{}
	if modelMap["auto_delete"] != nil {
		model.AutoDelete = core.BoolPtr(modelMap["auto_delete"].(bool))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["primary_ip"] != nil && len(modelMap["primary_ip"].([]interface{})) > 0 {
		PrimaryIPModel, err := ResourceIBMIsInstanceClusterNetworkAttachmentMapToClusterNetworkInterfacePrimaryIPPrototype(modelMap["primary_ip"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.PrimaryIP = PrimaryIPModel
	}
	if modelMap["subnet"] != nil && len(modelMap["subnet"].([]interface{})) > 0 {
		SubnetModel, err := ResourceIBMIsInstanceClusterNetworkAttachmentMapToClusterNetworkSubnetIdentity(modelMap["subnet"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Subnet = SubnetModel
	}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}

func ResourceIBMIsInstanceClusterNetworkAttachmentMapToClusterNetworkInterfacePrimaryIPPrototype(modelMap map[string]interface{}) (vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeIntf, error) {
	model := &vpcv1.ClusterNetworkInterfacePrimaryIPPrototype{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	if modelMap["address"] != nil && modelMap["address"].(string) != "" {
		model.Address = core.StringPtr(modelMap["address"].(string))
	}
	if modelMap["auto_delete"] != nil {
		model.AutoDelete = core.BoolPtr(modelMap["auto_delete"].(bool))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	return model, nil
}

func ResourceIBMIsInstanceClusterNetworkAttachmentMapToClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContext(modelMap map[string]interface{}) (vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextIntf, error) {
	model := &vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContext{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}

func ResourceIBMIsInstanceClusterNetworkAttachmentMapToClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByID(modelMap map[string]interface{}) (*vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByID, error) {
	model := &vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func ResourceIBMIsInstanceClusterNetworkAttachmentMapToClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByHref(modelMap map[string]interface{}) (*vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByHref, error) {
	model := &vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func ResourceIBMIsInstanceClusterNetworkAttachmentMapToClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPPrototypeClusterNetworkInterfacePrimaryIPContext(modelMap map[string]interface{}) (*vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPPrototypeClusterNetworkInterfacePrimaryIPContext, error) {
	model := &vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPPrototypeClusterNetworkInterfacePrimaryIPContext{}
	if modelMap["address"] != nil && modelMap["address"].(string) != "" {
		model.Address = core.StringPtr(modelMap["address"].(string))
	}
	if modelMap["auto_delete"] != nil {
		model.AutoDelete = core.BoolPtr(modelMap["auto_delete"].(bool))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	return model, nil
}

func ResourceIBMIsInstanceClusterNetworkAttachmentMapToClusterNetworkSubnetIdentity(modelMap map[string]interface{}) (vpcv1.ClusterNetworkSubnetIdentityIntf, error) {
	model := &vpcv1.ClusterNetworkSubnetIdentity{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}

func ResourceIBMIsInstanceClusterNetworkAttachmentMapToClusterNetworkSubnetIdentityByID(modelMap map[string]interface{}) (*vpcv1.ClusterNetworkSubnetIdentityByID, error) {
	model := &vpcv1.ClusterNetworkSubnetIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func ResourceIBMIsInstanceClusterNetworkAttachmentMapToClusterNetworkSubnetIdentityByHref(modelMap map[string]interface{}) (*vpcv1.ClusterNetworkSubnetIdentityByHref, error) {
	model := &vpcv1.ClusterNetworkSubnetIdentityByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func ResourceIBMIsInstanceClusterNetworkAttachmentMapToInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceInstanceClusterNetworkInterfacePrototypeInstanceClusterNetworkAttachment(modelMap map[string]interface{}) (*vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceInstanceClusterNetworkInterfacePrototypeInstanceClusterNetworkAttachment, error) {
	model := &vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceInstanceClusterNetworkInterfacePrototypeInstanceClusterNetworkAttachment{}
	if modelMap["auto_delete"] != nil {
		model.AutoDelete = core.BoolPtr(modelMap["auto_delete"].(bool))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["primary_ip"] != nil && len(modelMap["primary_ip"].([]interface{})) > 0 {
		PrimaryIPModel, err := ResourceIBMIsInstanceClusterNetworkAttachmentMapToClusterNetworkInterfacePrimaryIPPrototype(modelMap["primary_ip"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.PrimaryIP = PrimaryIPModel
	}
	if modelMap["subnet"] != nil && len(modelMap["subnet"].([]interface{})) > 0 {
		SubnetModel, err := ResourceIBMIsInstanceClusterNetworkAttachmentMapToClusterNetworkSubnetIdentity(modelMap["subnet"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Subnet = SubnetModel
	}
	return model, nil
}

func ResourceIBMIsInstanceClusterNetworkAttachmentMapToInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentity(modelMap map[string]interface{}) (vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityIntf, error) {
	model := &vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentity{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}

func ResourceIBMIsInstanceClusterNetworkAttachmentMapToInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByID(modelMap map[string]interface{}) (*vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByID, error) {
	model := &vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func ResourceIBMIsInstanceClusterNetworkAttachmentMapToInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByHref(modelMap map[string]interface{}) (*vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByHref, error) {
	model := &vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func ResourceIBMIsInstanceClusterNetworkAttachmentMapToInstanceClusterNetworkAttachmentBeforePrototype(modelMap map[string]interface{}) (vpcv1.InstanceClusterNetworkAttachmentBeforePrototypeIntf, error) {
	model := &vpcv1.InstanceClusterNetworkAttachmentBeforePrototype{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}

func ResourceIBMIsInstanceClusterNetworkAttachmentMapToInstanceClusterNetworkAttachmentBeforePrototypeInstanceClusterNetworkAttachmentIdentityByID(modelMap map[string]interface{}) (*vpcv1.InstanceClusterNetworkAttachmentBeforePrototypeInstanceClusterNetworkAttachmentIdentityByID, error) {
	model := &vpcv1.InstanceClusterNetworkAttachmentBeforePrototypeInstanceClusterNetworkAttachmentIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func ResourceIBMIsInstanceClusterNetworkAttachmentMapToInstanceClusterNetworkAttachmentBeforePrototypeInstanceClusterNetworkAttachmentIdentityByHref(modelMap map[string]interface{}) (*vpcv1.InstanceClusterNetworkAttachmentBeforePrototypeInstanceClusterNetworkAttachmentIdentityByHref, error) {
	model := &vpcv1.InstanceClusterNetworkAttachmentBeforePrototypeInstanceClusterNetworkAttachmentIdentityByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func ResourceIBMIsInstanceClusterNetworkAttachmentInstanceClusterNetworkAttachmentBeforeToMap(model *vpcv1.InstanceClusterNetworkAttachmentBefore) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func ResourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkInterfaceReferenceToMap(model *vpcv1.ClusterNetworkInterfaceReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := ResourceIBMIsInstanceClusterNetworkAttachmentDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	primaryIPMap, err := ResourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkSubnetReservedIPReferenceToMap(model.PrimaryIP)
	if err != nil {
		return modelMap, err
	}
	modelMap["primary_ip"] = []map[string]interface{}{primaryIPMap}
	modelMap["resource_type"] = *model.ResourceType
	subnetMap, err := ResourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkSubnetReferenceToMap(model.Subnet)
	if err != nil {
		return modelMap, err
	}
	modelMap["subnet"] = []map[string]interface{}{subnetMap}
	return modelMap, nil
}

func ResourceIBMIsInstanceClusterNetworkAttachmentDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func ResourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkSubnetReservedIPReferenceToMap(model *vpcv1.ClusterNetworkSubnetReservedIPReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = *model.Address
	if model.Deleted != nil {
		deletedMap, err := ResourceIBMIsInstanceClusterNetworkAttachmentDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func ResourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkSubnetReferenceToMap(model *vpcv1.ClusterNetworkSubnetReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := ResourceIBMIsInstanceClusterNetworkAttachmentDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func ResourceIBMIsInstanceClusterNetworkAttachmentInstanceClusterNetworkAttachmentLifecycleReasonToMap(model *vpcv1.InstanceClusterNetworkAttachmentLifecycleReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = *model.Code
	modelMap["message"] = *model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}

func ResourceIBMIsInstanceClusterNetworkAttachmentInstanceClusterNetworkAttachmentPatchAsPatch(patchVals *vpcv1.InstanceClusterNetworkAttachmentPatch, d *schema.ResourceData) map[string]interface{} {
	patch, _ := patchVals.AsPatch()
	var path string

	path = "name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["name"] = nil
	}

	return patch
}
