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

func ResourceIBMIsClusterNetworkSubnet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsClusterNetworkSubnetCreate,
		ReadContext:   resourceIBMIsClusterNetworkSubnetRead,
		UpdateContext: resourceIBMIsClusterNetworkSubnetUpdate,
		DeleteContext: resourceIBMIsClusterNetworkSubnetDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"cluster_network_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_cluster_network_subnet", "cluster_network_id"),
				Description:  "The cluster network identifier.",
			},
			"ip_version": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_cluster_network_subnet", "ip_version"),
				Description:  "The IP version for this cluster network subnet.The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
			},
			"ipv4_cidr_block": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"ipv4_cidr_block", "total_ipv4_address_count"},
				ValidateFunc: validate.InvokeValidator("ibm_is_cluster_network_subnet", "ipv4_cidr_block"),
				Description:  "The IPv4 range of this cluster network subnet, expressed in CIDR format.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_cluster_network_subnet", "name"),
				Description:  "The name for this cluster network subnet. The name is unique across all cluster network subnets in the cluster network.",
			},
			"total_ipv4_address_count": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"ipv4_cidr_block", "total_ipv4_address_count"},
				Description:  "The total number of IPv4 addresses in this cluster network subnet.Note: This is calculated as 2<sup>(32 - prefix length)</sup>. For example, the prefix length `/24` gives:<br> 2<sup>(32 - 24)</sup> = 2<sup>8</sup> = 256 addresses.",
			},
			"available_ipv4_address_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of IPv4 addresses in this cluster network subnet that are not in use, and have not been reserved by the user or the provider.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the cluster network subnet was created.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this cluster network subnet.",
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
				Description: "The lifecycle state of the cluster network subnet.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"cluster_network_subnet_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this cluster network subnet.",
			},
			"etag": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceIBMIsClusterNetworkSubnetValidator() *validate.ResourceValidator {
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
			Identifier:                 "ip_version",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "ipv4",
			Regexp:                     `^[a-z][a-z0-9]*(_[a-z0-9]+)*$`,
			MinValueLength:             1,
			MaxValueLength:             128,
		},
		validate.ValidateSchema{
			Identifier:                 "ipv4_cidr_block",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])(\/(3[0-2]|[1-2][0-9]|[0-9]))$`,
			MinValueLength:             9,
			MaxValueLength:             18,
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

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_cluster_network_subnet", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMIsClusterNetworkSubnetCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	bodyModelMap := map[string]interface{}{}
	createClusterNetworkSubnetOptions := &vpcv1.CreateClusterNetworkSubnetOptions{}

	if _, ok := d.GetOk("ip_version"); ok {
		bodyModelMap["ip_version"] = d.Get("ip_version")
	}
	if _, ok := d.GetOk("name"); ok {
		bodyModelMap["name"] = d.Get("name")
	}
	if _, ok := d.GetOk("total_ipv4_address_count"); ok {
		bodyModelMap["total_ipv4_address_count"] = d.Get("total_ipv4_address_count")
	}
	if _, ok := d.GetOk("ipv4_cidr_block"); ok {
		bodyModelMap["ipv4_cidr_block"] = d.Get("ipv4_cidr_block")
	}
	createClusterNetworkSubnetOptions.SetClusterNetworkID(d.Get("cluster_network_id").(string))
	convertedModel, err := ResourceIBMIsClusterNetworkSubnetMapToClusterNetworkSubnetPrototype(bodyModelMap)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet", "create", "parse-request-body").GetDiag()
	}
	createClusterNetworkSubnetOptions.ClusterNetworkSubnetPrototype = convertedModel

	clusterNetworkSubnet, _, err := vpcClient.CreateClusterNetworkSubnetWithContext(context, createClusterNetworkSubnetOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateClusterNetworkSubnetWithContext failed: %s", err.Error()), "ibm_is_cluster_network_subnet", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *createClusterNetworkSubnetOptions.ClusterNetworkID, *clusterNetworkSubnet.ID))

	return resourceIBMIsClusterNetworkSubnetRead(context, d, meta)
}

func resourceIBMIsClusterNetworkSubnetRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getClusterNetworkSubnetOptions := &vpcv1.GetClusterNetworkSubnetOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet", "read", "sep-id-parts").GetDiag()
	}

	getClusterNetworkSubnetOptions.SetClusterNetworkID(parts[0])
	getClusterNetworkSubnetOptions.SetID(parts[1])

	clusterNetworkSubnet, response, err := vpcClient.GetClusterNetworkSubnetWithContext(context, getClusterNetworkSubnetOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetClusterNetworkSubnetWithContext failed: %s", err.Error()), "ibm_is_cluster_network_subnet", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if !core.IsNil(clusterNetworkSubnet.IPVersion) {
		if err = d.Set("ip_version", clusterNetworkSubnet.IPVersion); err != nil {
			err = fmt.Errorf("Error setting ip_version: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet", "read", "set-ip_version").GetDiag()
		}
	}
	if !core.IsNil(clusterNetworkSubnet.Ipv4CIDRBlock) {
		if err = d.Set("ipv4_cidr_block", clusterNetworkSubnet.Ipv4CIDRBlock); err != nil {
			err = fmt.Errorf("Error setting ipv4_cidr_block: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet", "read", "set-ipv4_cidr_block").GetDiag()
		}
	}
	if !core.IsNil(clusterNetworkSubnet.Name) {
		if err = d.Set("name", clusterNetworkSubnet.Name); err != nil {
			err = fmt.Errorf("Error setting name: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet", "read", "set-name").GetDiag()
		}
	}
	if !core.IsNil(clusterNetworkSubnet.TotalIpv4AddressCount) {
		if err = d.Set("total_ipv4_address_count", flex.IntValue(clusterNetworkSubnet.TotalIpv4AddressCount)); err != nil {
			err = fmt.Errorf("Error setting total_ipv4_address_count: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet", "read", "set-total_ipv4_address_count").GetDiag()
		}
	}
	if err = d.Set("available_ipv4_address_count", flex.IntValue(clusterNetworkSubnet.AvailableIpv4AddressCount)); err != nil {
		err = fmt.Errorf("Error setting available_ipv4_address_count: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet", "read", "set-available_ipv4_address_count").GetDiag()
	}
	if err = d.Set("created_at", flex.DateTimeToString(clusterNetworkSubnet.CreatedAt)); err != nil {
		err = fmt.Errorf("Error setting created_at: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet", "read", "set-created_at").GetDiag()
	}
	if err = d.Set("href", clusterNetworkSubnet.Href); err != nil {
		err = fmt.Errorf("Error setting href: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet", "read", "set-href").GetDiag()
	}
	lifecycleReasons := []map[string]interface{}{}
	for _, lifecycleReasonsItem := range clusterNetworkSubnet.LifecycleReasons {
		lifecycleReasonsItemMap, err := ResourceIBMIsClusterNetworkSubnetClusterNetworkSubnetLifecycleReasonToMap(&lifecycleReasonsItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet", "read", "lifecycle_reasons-to-map").GetDiag()
		}
		lifecycleReasons = append(lifecycleReasons, lifecycleReasonsItemMap)
	}
	if err = d.Set("lifecycle_reasons", lifecycleReasons); err != nil {
		err = fmt.Errorf("Error setting lifecycle_reasons: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet", "read", "set-lifecycle_reasons").GetDiag()
	}
	if err = d.Set("lifecycle_state", clusterNetworkSubnet.LifecycleState); err != nil {
		err = fmt.Errorf("Error setting lifecycle_state: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet", "read", "set-lifecycle_state").GetDiag()
	}
	if err = d.Set("resource_type", clusterNetworkSubnet.ResourceType); err != nil {
		err = fmt.Errorf("Error setting resource_type: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet", "read", "set-resource_type").GetDiag()
	}
	if err = d.Set("cluster_network_subnet_id", clusterNetworkSubnet.ID); err != nil {
		err = fmt.Errorf("Error setting cluster_network_subnet_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet", "read", "set-cluster_network_subnet_id").GetDiag()
	}
	if err = d.Set("etag", response.Headers.Get("Etag")); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting etag: %s", err), "ibm_is_cluster_network_subnet", "read", "set-etag").GetDiag()
	}

	return nil
}

func resourceIBMIsClusterNetworkSubnetUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateClusterNetworkSubnetOptions := &vpcv1.UpdateClusterNetworkSubnetOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet", "update", "sep-id-parts").GetDiag()
	}

	updateClusterNetworkSubnetOptions.SetClusterNetworkID(parts[0])
	updateClusterNetworkSubnetOptions.SetID(parts[1])

	hasChange := false

	patchVals := &vpcv1.ClusterNetworkSubnetPatch{}
	if d.HasChange("cluster_network_id") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "cluster_network_id")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_is_cluster_network_subnet", "update", "cluster_network_id-forces-new").GetDiag()
	}
	if d.HasChange("name") {
		newName := d.Get("name").(string)
		patchVals.Name = &newName
		hasChange = true
	}
	// updateClusterNetworkSubnetOptions.SetIfMatch(d.Get("etag").(string))

	if hasChange {
		// Fields with `nil` values are omitted from the generic map,
		// so we need to re-add them to support removing arguments
		// in merge-patch operations sent to the service.
		updateClusterNetworkSubnetOptions.ClusterNetworkSubnetPatch = ResourceIBMIsClusterNetworkSubnetClusterNetworkSubnetPatchAsPatch(patchVals, d)

		_, _, err = vpcClient.UpdateClusterNetworkSubnetWithContext(context, updateClusterNetworkSubnetOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateClusterNetworkSubnetWithContext failed: %s", err.Error()), "ibm_is_cluster_network_subnet", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIBMIsClusterNetworkSubnetRead(context, d, meta)
}

func resourceIBMIsClusterNetworkSubnetDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteClusterNetworkSubnetOptions := &vpcv1.DeleteClusterNetworkSubnetOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_subnet", "delete", "sep-id-parts").GetDiag()
	}

	deleteClusterNetworkSubnetOptions.SetClusterNetworkID(parts[0])
	deleteClusterNetworkSubnetOptions.SetID(parts[1])

	_, _, err = vpcClient.DeleteClusterNetworkSubnetWithContext(context, deleteClusterNetworkSubnetOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteClusterNetworkSubnetWithContext failed: %s", err.Error()), "ibm_is_cluster_network_subnet", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIBMIsClusterNetworkSubnetMapToClusterNetworkSubnetPrototype(modelMap map[string]interface{}) (vpcv1.ClusterNetworkSubnetPrototypeIntf, error) {
	model := &vpcv1.ClusterNetworkSubnetPrototype{}
	if modelMap["ip_version"] != nil && modelMap["ip_version"].(string) != "" {
		model.IPVersion = core.StringPtr(modelMap["ip_version"].(string))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["total_ipv4_address_count"] != nil {
		model.TotalIpv4AddressCount = core.Int64Ptr(int64(modelMap["total_ipv4_address_count"].(int)))
	}
	if modelMap["ipv4_cidr_block"] != nil && modelMap["ipv4_cidr_block"].(string) != "" {
		model.Ipv4CIDRBlock = core.StringPtr(modelMap["ipv4_cidr_block"].(string))
	}
	return model, nil
}

func ResourceIBMIsClusterNetworkSubnetMapToClusterNetworkSubnetPrototypeClusterNetworkSubnetByTotalCountPrototype(modelMap map[string]interface{}) (*vpcv1.ClusterNetworkSubnetPrototypeClusterNetworkSubnetByTotalCountPrototype, error) {
	model := &vpcv1.ClusterNetworkSubnetPrototypeClusterNetworkSubnetByTotalCountPrototype{}
	if modelMap["ip_version"] != nil && modelMap["ip_version"].(string) != "" {
		model.IPVersion = core.StringPtr(modelMap["ip_version"].(string))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	model.TotalIpv4AddressCount = core.Int64Ptr(int64(modelMap["total_ipv4_address_count"].(int)))
	return model, nil
}

func ResourceIBMIsClusterNetworkSubnetMapToClusterNetworkSubnetPrototypeClusterNetworkSubnetByIPv4CIDRBlockPrototype(modelMap map[string]interface{}) (*vpcv1.ClusterNetworkSubnetPrototypeClusterNetworkSubnetByIPv4CIDRBlockPrototype, error) {
	model := &vpcv1.ClusterNetworkSubnetPrototypeClusterNetworkSubnetByIPv4CIDRBlockPrototype{}
	if modelMap["ip_version"] != nil && modelMap["ip_version"].(string) != "" {
		model.IPVersion = core.StringPtr(modelMap["ip_version"].(string))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	model.Ipv4CIDRBlock = core.StringPtr(modelMap["ipv4_cidr_block"].(string))
	return model, nil
}

func ResourceIBMIsClusterNetworkSubnetClusterNetworkSubnetLifecycleReasonToMap(model *vpcv1.ClusterNetworkSubnetLifecycleReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = *model.Code
	modelMap["message"] = *model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}

func ResourceIBMIsClusterNetworkSubnetClusterNetworkSubnetPatchAsPatch(patchVals *vpcv1.ClusterNetworkSubnetPatch, d *schema.ResourceData) map[string]interface{} {
	patch, _ := patchVals.AsPatch()
	var path string

	path = "name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["name"] = nil
	}

	return patch
}
