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

func ResourceIBMIsClusterNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsClusterNetworkCreate,
		ReadContext:   resourceIBMIsClusterNetworkRead,
		UpdateContext: resourceIBMIsClusterNetworkUpdate,
		DeleteContext: resourceIBMIsClusterNetworkDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_cluster_network", "name"),
				Description:  "The name for this cluster network. The name must not be used by another cluster network in the region.",
			},
			"profile": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The globally unique name for this cluster network profile.",
			},
			"resource_group": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The unique identifier for this resource group for this cluster network.",
			},
			"subnet_prefixes": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				// Default:     [{"cidr":"10.0.0.0/9"}],
				Description: "The IP address ranges available for subnets for this cluster network.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allocation_policy": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The allocation policy for this subnet prefix:- `auto`: Subnets created by total count in this cluster network can use this prefix.",
						},
						"cidr": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The CIDR block for this prefix.",
						},
					},
				},
			},
			"vpc": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				ForceNew:    true,
				Description: "The VPC this cluster network resides in.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this VPC.",
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
							Description: "The URL for this VPC.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The unique identifier for this VPC.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this VPC. The name is unique across all VPCs in the region.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"zone": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The globally unique name for the zone this cluster network resides in.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the cluster network was created.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this cluster network.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this cluster network.",
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
				Description: "The lifecycle state of the cluster network.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"etag": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceIBMIsClusterNetworkValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
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

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_cluster_network", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMIsClusterNetworkCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createClusterNetworkOptions := &vpcv1.CreateClusterNetworkOptions{}

	createClusterNetworkOptions.Profile = &vpcv1.ClusterNetworkProfileIdentity{
		Name: core.StringPtr(d.Get("profile").(string)),
	}
	vpcModel, err := ResourceIBMIsClusterNetworkMapToVPCIdentity(d.Get("vpc.0").(map[string]interface{}))
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network", "create", "parse-vpc").GetDiag()
	}
	createClusterNetworkOptions.SetVPC(vpcModel)
	createClusterNetworkOptions.Zone = &vpcv1.ZoneIdentity{
		Name: core.StringPtr(d.Get("zone").(string)),
	}
	if _, ok := d.GetOk("name"); ok {
		createClusterNetworkOptions.SetName(d.Get("name").(string))
	}
	if rgOk, ok := d.GetOk("resource_group"); ok {

		createClusterNetworkOptions.ResourceGroup = &vpcv1.ResourceGroupIdentity{
			ID: core.StringPtr(rgOk.(string)),
		}
	}
	if _, ok := d.GetOk("subnet_prefixes"); ok {
		var subnetPrefixes []vpcv1.ClusterNetworkSubnetPrefixPrototype
		for _, v := range d.Get("subnet_prefixes").([]interface{}) {
			value := v.(map[string]interface{})
			subnetPrefixesItem, err := ResourceIBMIsClusterNetworkMapToClusterNetworkSubnetPrefixPrototype(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network", "create", "parse-subnet_prefixes").GetDiag()
			}
			subnetPrefixes = append(subnetPrefixes, *subnetPrefixesItem)
		}
		createClusterNetworkOptions.SetSubnetPrefixes(subnetPrefixes)
	}

	clusterNetwork, _, err := vpcClient.CreateClusterNetworkWithContext(context, createClusterNetworkOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateClusterNetworkWithContext failed: %s", err.Error()), "ibm_is_cluster_network", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*clusterNetwork.ID)

	return resourceIBMIsClusterNetworkRead(context, d, meta)
}

func resourceIBMIsClusterNetworkRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getClusterNetworkOptions := &vpcv1.GetClusterNetworkOptions{}

	getClusterNetworkOptions.SetID(d.Id())

	clusterNetwork, response, err := vpcClient.GetClusterNetworkWithContext(context, getClusterNetworkOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetClusterNetworkWithContext failed: %s", err.Error()), "ibm_is_cluster_network", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if !core.IsNil(clusterNetwork.Name) {
		if err = d.Set("name", clusterNetwork.Name); err != nil {
			err = fmt.Errorf("Error setting name: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network", "read", "set-name").GetDiag()
		}
	}
	if err = d.Set("profile", clusterNetwork.Profile.Name); err != nil {
		err = fmt.Errorf("Error setting profile: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network", "read", "set-profile").GetDiag()
	}
	if !core.IsNil(clusterNetwork.ResourceGroup) {
		if err = d.Set("resource_group", clusterNetwork.ResourceGroup.ID); err != nil {
			err = fmt.Errorf("Error setting resource_group: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network", "read", "set-resource_group").GetDiag()
		}
	}
	if !core.IsNil(clusterNetwork.SubnetPrefixes) {
		subnetPrefixes := []map[string]interface{}{}
		for _, subnetPrefixesItem := range clusterNetwork.SubnetPrefixes {
			subnetPrefixesItemMap, err := ResourceIBMIsClusterNetworkClusterNetworkSubnetPrefixToMap(&subnetPrefixesItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network", "read", "subnet_prefixes-to-map").GetDiag()
			}
			subnetPrefixes = append(subnetPrefixes, subnetPrefixesItemMap)
		}
		if err = d.Set("subnet_prefixes", subnetPrefixes); err != nil {
			err = fmt.Errorf("Error setting subnet_prefixes: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network", "read", "set-subnet_prefixes").GetDiag()
		}
	}
	vpcMap, err := ResourceIBMIsClusterNetworkVPCReferenceToMap(clusterNetwork.VPC)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network", "read", "vpc-to-map").GetDiag()
	}
	if err = d.Set("vpc", []map[string]interface{}{vpcMap}); err != nil {
		err = fmt.Errorf("Error setting vpc: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network", "read", "set-vpc").GetDiag()
	}
	if err = d.Set("zone", clusterNetwork.Zone.Name); err != nil {
		err = fmt.Errorf("Error setting zone: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network", "read", "set-zone").GetDiag()
	}
	if err = d.Set("created_at", flex.DateTimeToString(clusterNetwork.CreatedAt)); err != nil {
		err = fmt.Errorf("Error setting created_at: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network", "read", "set-created_at").GetDiag()
	}
	if err = d.Set("crn", clusterNetwork.CRN); err != nil {
		err = fmt.Errorf("Error setting crn: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network", "read", "set-crn").GetDiag()
	}
	if err = d.Set("href", clusterNetwork.Href); err != nil {
		err = fmt.Errorf("Error setting href: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network", "read", "set-href").GetDiag()
	}
	lifecycleReasons := []map[string]interface{}{}
	for _, lifecycleReasonsItem := range clusterNetwork.LifecycleReasons {
		lifecycleReasonsItemMap, err := ResourceIBMIsClusterNetworkClusterNetworkLifecycleReasonToMap(&lifecycleReasonsItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network", "read", "lifecycle_reasons-to-map").GetDiag()
		}
		lifecycleReasons = append(lifecycleReasons, lifecycleReasonsItemMap)
	}
	if err = d.Set("lifecycle_reasons", lifecycleReasons); err != nil {
		err = fmt.Errorf("Error setting lifecycle_reasons: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network", "read", "set-lifecycle_reasons").GetDiag()
	}
	if err = d.Set("lifecycle_state", clusterNetwork.LifecycleState); err != nil {
		err = fmt.Errorf("Error setting lifecycle_state: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network", "read", "set-lifecycle_state").GetDiag()
	}
	if err = d.Set("resource_type", clusterNetwork.ResourceType); err != nil {
		err = fmt.Errorf("Error setting resource_type: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network", "read", "set-resource_type").GetDiag()
	}
	if err = d.Set("etag", response.Headers.Get("Etag")); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting etag: %s", err), "ibm_is_cluster_network", "read", "set-etag").GetDiag()
	}

	return nil
}

func resourceIBMIsClusterNetworkUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateClusterNetworkOptions := &vpcv1.UpdateClusterNetworkOptions{}

	updateClusterNetworkOptions.SetID(d.Id())

	hasChange := false

	patchVals := &vpcv1.ClusterNetworkPatch{}
	if d.HasChange("name") {
		newName := d.Get("name").(string)
		patchVals.Name = &newName
		hasChange = true
	}
	// updateClusterNetworkOptions.SetIfMatch(d.Get("etag").(string))

	if hasChange {
		updateClusterNetworkOptions.ClusterNetworkPatch = ResourceIBMIsClusterNetworkClusterNetworkPatchAsPatch(patchVals, d)

		_, response, err := vpcClient.UpdateClusterNetworkWithContext(context, updateClusterNetworkOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateClusterNetworkWithContext failed: %s", err.Error()), "ibm_is_cluster_network", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("etag", response.Headers.Get("Etag")); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting etag: %s", err), "ibm_is_cluster_network", "update", "set-etag").GetDiag()
		}
	}

	return resourceIBMIsClusterNetworkRead(context, d, meta)
}

func resourceIBMIsClusterNetworkDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteClusterNetworkOptions := &vpcv1.DeleteClusterNetworkOptions{}

	deleteClusterNetworkOptions.SetID(d.Id())

	_, _, err = vpcClient.DeleteClusterNetworkWithContext(context, deleteClusterNetworkOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteClusterNetworkWithContext failed: %s", err.Error()), "ibm_is_cluster_network", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIBMIsClusterNetworkMapToClusterNetworkProfileIdentityByName(modelMap map[string]interface{}) (*vpcv1.ClusterNetworkProfileIdentityByName, error) {
	model := &vpcv1.ClusterNetworkProfileIdentityByName{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	return model, nil
}

func ResourceIBMIsClusterNetworkMapToClusterNetworkProfileIdentityByHref(modelMap map[string]interface{}) (*vpcv1.ClusterNetworkProfileIdentityByHref, error) {
	model := &vpcv1.ClusterNetworkProfileIdentityByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func ResourceIBMIsClusterNetworkMapToVPCIdentity(modelMap map[string]interface{}) (vpcv1.VPCIdentityIntf, error) {
	model := &vpcv1.VPCIdentity{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	return model, nil
}

func ResourceIBMIsClusterNetworkMapToVPCIdentityByID(modelMap map[string]interface{}) (*vpcv1.VPCIdentityByID, error) {
	model := &vpcv1.VPCIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func ResourceIBMIsClusterNetworkMapToVPCIdentityByCRN(modelMap map[string]interface{}) (*vpcv1.VPCIdentityByCRN, error) {
	model := &vpcv1.VPCIdentityByCRN{}
	model.CRN = core.StringPtr(modelMap["crn"].(string))
	return model, nil
}

func ResourceIBMIsClusterNetworkMapToVPCIdentityByHref(modelMap map[string]interface{}) (*vpcv1.VPCIdentityByHref, error) {
	model := &vpcv1.VPCIdentityByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func ResourceIBMIsClusterNetworkMapToZoneIdentityByName(modelMap map[string]interface{}) (*vpcv1.ZoneIdentityByName, error) {
	model := &vpcv1.ZoneIdentityByName{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	return model, nil
}

func ResourceIBMIsClusterNetworkMapToZoneIdentityByHref(modelMap map[string]interface{}) (*vpcv1.ZoneIdentityByHref, error) {
	model := &vpcv1.ZoneIdentityByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func ResourceIBMIsClusterNetworkMapToResourceGroupIdentityByID(modelMap map[string]interface{}) (*vpcv1.ResourceGroupIdentityByID, error) {
	model := &vpcv1.ResourceGroupIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func ResourceIBMIsClusterNetworkMapToClusterNetworkSubnetPrefixPrototype(modelMap map[string]interface{}) (*vpcv1.ClusterNetworkSubnetPrefixPrototype, error) {
	model := &vpcv1.ClusterNetworkSubnetPrefixPrototype{}
	if modelMap["cidr"] != nil && modelMap["cidr"].(string) != "" {
		model.CIDR = core.StringPtr(modelMap["cidr"].(string))
	}
	return model, nil
}

func ResourceIBMIsClusterNetworkClusterNetworkSubnetPrefixToMap(model *vpcv1.ClusterNetworkSubnetPrefix) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["allocation_policy"] = *model.AllocationPolicy
	modelMap["cidr"] = *model.CIDR
	return modelMap, nil
}

func ResourceIBMIsClusterNetworkVPCReferenceToMap(model *vpcv1.VPCReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := ResourceIBMIsClusterNetworkDeletedToMap(model.Deleted)
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

func ResourceIBMIsClusterNetworkDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func ResourceIBMIsClusterNetworkZoneReferenceToMap(model *vpcv1.ZoneReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func ResourceIBMIsClusterNetworkClusterNetworkLifecycleReasonToMap(model *vpcv1.ClusterNetworkLifecycleReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = *model.Code
	modelMap["message"] = *model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}

func ResourceIBMIsClusterNetworkClusterNetworkPatchAsPatch(patchVals *vpcv1.ClusterNetworkPatch, d *schema.ResourceData) map[string]interface{} {
	patch, _ := patchVals.AsPatch()
	var path string

	path = "name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["name"] = nil
	}

	return patch
}
