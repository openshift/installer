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

func ResourceIBMIsClusterNetworkSubnetReservedIP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsClusterNetworkSubnetReservedIPCreate,
		ReadContext:   resourceIBMIsClusterNetworkSubnetReservedIPRead,
		UpdateContext: resourceIBMIsClusterNetworkSubnetReservedIPUpdate,
		DeleteContext: resourceIBMIsClusterNetworkSubnetReservedIPDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"cluster_network_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_cluster_network_subnet_reserved_ip", "cluster_network_id"),
				Description:  "The cluster network identifier.",
			},
			"cluster_network_subnet_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_cluster_network_subnet_reserved_ip", "cluster_network_subnet_id"),
				Description:  "The cluster network subnet identifier.",
			},
			"address": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_cluster_network_subnet_reserved_ip", "address"),
				Description:  "The IP address.If the address is pending allocation, the value will be `0.0.0.0`.This property may [expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) to support IPv6 addresses in the future.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_cluster_network_subnet_reserved_ip", "name"),
				Description:  "The name for this cluster network subnet reserved IP. The name is unique across all reserved IPs in a cluster network subnet.",
			},
			"auto_delete": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether this cluster network subnet reserved IP member will be automatically deleted when either `target` is deleted, or the cluster network subnet reserved IP is unbound.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the cluster network subnet reserved IP was created.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this cluster network subnet reserved IP.",
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
				Description: "The lifecycle state of the cluster network subnet reserved IP.",
			},
			"owner": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The owner of the cluster network subnet reserved IPThe enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"target": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The target this cluster network subnet reserved IP is bound to.If absent, this cluster network subnet reserved IP is provider-owned or unbound.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
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
							Optional:    true,
							Computed:    true,
							Description: "The URL for this cluster network interface.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The unique identifier for this cluster network interface.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The name for this cluster network interface. The name is unique across all interfaces in the cluster network.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"cluster_network_subnet_reserved_ip_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this cluster network subnet reserved IP.",
			},
			"etag": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceIBMIsClusterNetworkSubnetReservedIPValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cluster_network_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[-0-9a-z_]+$`,
			MinValueLength:             1,
			MaxValueLength:             64,
		},
		validate.ValidateSchema{
			Identifier:                 "cluster_network_subnet_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[-0-9a-z_]+$`,
			MinValueLength:             1,
			MaxValueLength:             64,
		},
		validate.ValidateSchema{
			Identifier:                 "address",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`,
			MinValueLength:             7,
			MaxValueLength:             15,
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

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_cluster_network_subnet_reserved_ip", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMIsClusterNetworkSubnetReservedIPCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet_reserved_ip", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createClusterNetworkSubnetReservedIPOptions := &vpcv1.CreateClusterNetworkSubnetReservedIPOptions{}

	createClusterNetworkSubnetReservedIPOptions.SetClusterNetworkID(d.Get("cluster_network_id").(string))
	createClusterNetworkSubnetReservedIPOptions.SetClusterNetworkSubnetID(d.Get("cluster_network_subnet_id").(string))
	if _, ok := d.GetOk("address"); ok {
		createClusterNetworkSubnetReservedIPOptions.SetAddress(d.Get("address").(string))
	}
	if _, ok := d.GetOk("name"); ok {
		createClusterNetworkSubnetReservedIPOptions.SetName(d.Get("name").(string))
	}

	clusterNetworkSubnetReservedIP, _, err := vpcClient.CreateClusterNetworkSubnetReservedIPWithContext(context, createClusterNetworkSubnetReservedIPOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateClusterNetworkSubnetReservedIPWithContext failed: %s", err.Error()), "ibm_is_cluster_network_subnet_reserved_ip", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", *createClusterNetworkSubnetReservedIPOptions.ClusterNetworkID, *createClusterNetworkSubnetReservedIPOptions.ClusterNetworkSubnetID, *clusterNetworkSubnetReservedIP.ID))

	return resourceIBMIsClusterNetworkSubnetReservedIPRead(context, d, meta)
}

func resourceIBMIsClusterNetworkSubnetReservedIPRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet_reserved_ip", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getClusterNetworkSubnetReservedIPOptions := &vpcv1.GetClusterNetworkSubnetReservedIPOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet_reserved_ip", "read", "sep-id-parts").GetDiag()
	}

	getClusterNetworkSubnetReservedIPOptions.SetClusterNetworkID(parts[0])
	getClusterNetworkSubnetReservedIPOptions.SetClusterNetworkSubnetID(parts[1])
	getClusterNetworkSubnetReservedIPOptions.SetID(parts[2])

	clusterNetworkSubnetReservedIP, response, err := vpcClient.GetClusterNetworkSubnetReservedIPWithContext(context, getClusterNetworkSubnetReservedIPOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetClusterNetworkSubnetReservedIPWithContext failed: %s", err.Error()), "ibm_is_cluster_network_subnet_reserved_ip", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if !core.IsNil(clusterNetworkSubnetReservedIP.Address) {
		if err = d.Set("address", clusterNetworkSubnetReservedIP.Address); err != nil {
			err = fmt.Errorf("Error setting address: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet_reserved_ip", "read", "set-address").GetDiag()
		}
	}
	if !core.IsNil(clusterNetworkSubnetReservedIP.Name) {
		if err = d.Set("name", clusterNetworkSubnetReservedIP.Name); err != nil {
			err = fmt.Errorf("Error setting name: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet_reserved_ip", "read", "set-name").GetDiag()
		}
	}
	if err = d.Set("auto_delete", clusterNetworkSubnetReservedIP.AutoDelete); err != nil {
		err = fmt.Errorf("Error setting auto_delete: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet_reserved_ip", "read", "set-auto_delete").GetDiag()
	}
	if err = d.Set("created_at", flex.DateTimeToString(clusterNetworkSubnetReservedIP.CreatedAt)); err != nil {
		err = fmt.Errorf("Error setting created_at: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet_reserved_ip", "read", "set-created_at").GetDiag()
	}
	if err = d.Set("href", clusterNetworkSubnetReservedIP.Href); err != nil {
		err = fmt.Errorf("Error setting href: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet_reserved_ip", "read", "set-href").GetDiag()
	}
	lifecycleReasons := []map[string]interface{}{}
	for _, lifecycleReasonsItem := range clusterNetworkSubnetReservedIP.LifecycleReasons {
		lifecycleReasonsItemMap, err := ResourceIBMIsClusterNetworkSubnetReservedIPClusterNetworkSubnetReservedIPLifecycleReasonToMap(&lifecycleReasonsItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet_reserved_ip", "read", "lifecycle_reasons-to-map").GetDiag()
		}
		lifecycleReasons = append(lifecycleReasons, lifecycleReasonsItemMap)
	}
	if err = d.Set("lifecycle_reasons", lifecycleReasons); err != nil {
		err = fmt.Errorf("Error setting lifecycle_reasons: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet_reserved_ip", "read", "set-lifecycle_reasons").GetDiag()
	}
	if err = d.Set("lifecycle_state", clusterNetworkSubnetReservedIP.LifecycleState); err != nil {
		err = fmt.Errorf("Error setting lifecycle_state: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet_reserved_ip", "read", "set-lifecycle_state").GetDiag()
	}
	if err = d.Set("owner", clusterNetworkSubnetReservedIP.Owner); err != nil {
		err = fmt.Errorf("Error setting owner: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet_reserved_ip", "read", "set-owner").GetDiag()
	}
	if err = d.Set("resource_type", clusterNetworkSubnetReservedIP.ResourceType); err != nil {
		err = fmt.Errorf("Error setting resource_type: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet_reserved_ip", "read", "set-resource_type").GetDiag()
	}
	targetMap := make(map[string]interface{})
	if !core.IsNil(clusterNetworkSubnetReservedIP.Target) {
		targetMap, err = ResourceIBMIsClusterNetworkSubnetReservedIPClusterNetworkSubnetReservedIPTargetToMap(clusterNetworkSubnetReservedIP.Target)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet_reserved_ip", "read", "target-to-map").GetDiag()
		}
	}
	if err = d.Set("target", []map[string]interface{}{targetMap}); err != nil {
		err = fmt.Errorf("Error setting target: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet_reserved_ip", "read", "set-target").GetDiag()
	}
	if err = d.Set("cluster_network_subnet_reserved_ip_id", clusterNetworkSubnetReservedIP.ID); err != nil {
		err = fmt.Errorf("Error setting cluster_network_subnet_reserved_ip_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet_reserved_ip", "read", "set-cluster_network_subnet_reserved_ip_id").GetDiag()
	}
	if err = d.Set("etag", response.Headers.Get("Etag")); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting etag: %s", err), "ibm_is_cluster_network_subnet_reserved_ip", "read", "set-etag").GetDiag()
	}

	return nil
}

func resourceIBMIsClusterNetworkSubnetReservedIPUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet_reserved_ip", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateClusterNetworkSubnetReservedIPOptions := &vpcv1.UpdateClusterNetworkSubnetReservedIPOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet_reserved_ip", "update", "sep-id-parts").GetDiag()
	}

	updateClusterNetworkSubnetReservedIPOptions.SetClusterNetworkID(parts[0])
	updateClusterNetworkSubnetReservedIPOptions.SetClusterNetworkSubnetID(parts[1])
	updateClusterNetworkSubnetReservedIPOptions.SetID(parts[2])

	hasChange := false

	patchVals := &vpcv1.ClusterNetworkSubnetReservedIPPatch{}
	if d.HasChange("cluster_network_id") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "cluster_network_id")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_is_cluster_network_subnet_reserved_ip", "update", "cluster_network_id-forces-new").GetDiag()
	}
	if d.HasChange("cluster_network_subnet_id") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "cluster_network_subnet_id")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_is_cluster_network_subnet_reserved_ip", "update", "cluster_network_subnet_id-forces-new").GetDiag()
	}
	if d.HasChange("auto_delete") {
		newAutoDelete := d.Get("auto_delete").(bool)
		patchVals.AutoDelete = &newAutoDelete
		hasChange = true
	}
	if d.HasChange("name") {
		newName := d.Get("name").(string)
		patchVals.Name = &newName
		hasChange = true
	}
	// updateClusterNetworkSubnetReservedIPOptions.SetIfMatch(d.Get("etag").(string))

	if hasChange {
		updateClusterNetworkSubnetReservedIPOptions.ClusterNetworkSubnetReservedIPPatch = ResourceIBMIsClusterNetworkSubnetReservedIPClusterNetworkSubnetReservedIPPatchAsPatch(patchVals, d)

		_, _, err = vpcClient.UpdateClusterNetworkSubnetReservedIPWithContext(context, updateClusterNetworkSubnetReservedIPOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateClusterNetworkSubnetReservedIPWithContext failed: %s", err.Error()), "ibm_is_cluster_network_subnet_reserved_ip", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIBMIsClusterNetworkSubnetReservedIPRead(context, d, meta)
}

func resourceIBMIsClusterNetworkSubnetReservedIPDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet_reserved_ip", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteClusterNetworkSubnetReservedIPOptions := &vpcv1.DeleteClusterNetworkSubnetReservedIPOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet_reserved_ip", "delete", "sep-id-parts").GetDiag()
	}

	deleteClusterNetworkSubnetReservedIPOptions.SetClusterNetworkID(parts[0])
	deleteClusterNetworkSubnetReservedIPOptions.SetClusterNetworkSubnetID(parts[1])
	deleteClusterNetworkSubnetReservedIPOptions.SetID(parts[2])

	_, _, err = vpcClient.DeleteClusterNetworkSubnetReservedIPWithContext(context, deleteClusterNetworkSubnetReservedIPOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteClusterNetworkSubnetReservedIPWithContext failed: %s", err.Error()), "ibm_is_cluster_network_subnet_reserved_ip", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIBMIsClusterNetworkSubnetReservedIPClusterNetworkSubnetReservedIPLifecycleReasonToMap(model *vpcv1.ClusterNetworkSubnetReservedIPLifecycleReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = *model.Code
	modelMap["message"] = *model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}

func ResourceIBMIsClusterNetworkSubnetReservedIPClusterNetworkSubnetReservedIPTargetToMap(model vpcv1.ClusterNetworkSubnetReservedIPTargetIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.ClusterNetworkSubnetReservedIPTargetClusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContext); ok {
		return ResourceIBMIsClusterNetworkSubnetReservedIPClusterNetworkSubnetReservedIPTargetClusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContextToMap(model.(*vpcv1.ClusterNetworkSubnetReservedIPTargetClusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContext))
	} else if _, ok := model.(*vpcv1.ClusterNetworkSubnetReservedIPTarget); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.ClusterNetworkSubnetReservedIPTarget)
		if model.Deleted != nil {
			deletedMap, err := ResourceIBMIsClusterNetworkSubnetReservedIPDeletedToMap(model.Deleted)
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
		if model.ResourceType != nil {
			modelMap["resource_type"] = *model.ResourceType
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized vpcv1.ClusterNetworkSubnetReservedIPTargetIntf subtype encountered")
	}
}

func ResourceIBMIsClusterNetworkSubnetReservedIPDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func ResourceIBMIsClusterNetworkSubnetReservedIPClusterNetworkSubnetReservedIPTargetClusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContextToMap(model *vpcv1.ClusterNetworkSubnetReservedIPTargetClusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := ResourceIBMIsClusterNetworkSubnetReservedIPDeletedToMap(model.Deleted)
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

func ResourceIBMIsClusterNetworkSubnetReservedIPClusterNetworkSubnetReservedIPPatchAsPatch(patchVals *vpcv1.ClusterNetworkSubnetReservedIPPatch, d *schema.ResourceData) map[string]interface{} {
	patch, _ := patchVals.AsPatch()
	var path string

	path = "auto_delete"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["auto_delete"] = nil
	}
	path = "name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["name"] = nil
	}

	return patch
}
