// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func ResourceIbmIsDedicatedHostGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmIsDedicatedHostGroupCreate,
		ReadContext:   resourceIbmIsDedicatedHostGroupRead,
		UpdateContext: resourceIbmIsDedicatedHostGroupUpdate,
		DeleteContext: resourceIbmIsDedicatedHostGroupDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"class": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The dedicated host profile class for hosts in this group.",
			},
			"family": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_dedicated_host_group", "family"),
				Description:  "The dedicated host profile family for hosts in this group.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_dedicated_host_group", "name"),
				Description:  "The unique user-defined name for this dedicated host group. If unspecified, the name will be a hyphenated list of randomly-selected words.",
			},
			"resource_group": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The unique identifier of the resource group to use. If unspecified, the account's [default resourcegroup](https://cloud.ibm.com/apidocs/resource-manager#introduction) is used.",
			},
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The globally unique name of the zone this dedicated host group will reside in.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the dedicated host group was created.",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this dedicated host group.",
			},
			"dedicated_hosts": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The dedicated hosts that are in this dedicated host group.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this dedicated host.",
						},
						"deleted": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this dedicated host.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this dedicated host.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique user-defined name for this dedicated host. If unspecified, the name will be a hyphenated list of randomly-selected words.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of resource referenced.",
						},
					},
				},
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this dedicated host group.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of resource referenced.",
			},
			"supported_instance_profiles": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Array of instance profiles that can be used by instances placed on this dedicated host group.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this virtual server instance profile.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique name for this virtual server instance profile.",
						},
					},
				},
			},
		},
	}
}

func ResourceIbmIsDedicatedHostGroupValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "family",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "balanced, compute, memory",
		})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		})

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_dedicated_host_group", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmIsDedicatedHostGroupCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	createDedicatedHostGroupOptions := &vpcv1.CreateDedicatedHostGroupOptions{}

	if class, ok := d.GetOk("class"); ok {
		createDedicatedHostGroupOptions.SetClass(class.(string))
	}
	if family, ok := d.GetOk("family"); ok {
		createDedicatedHostGroupOptions.SetFamily(family.(string))
	}
	if name, ok := d.GetOk("name"); ok {
		createDedicatedHostGroupOptions.SetName(name.(string))
	}
	if resgroup, ok := d.GetOk("resource_group"); ok {
		resgroupstr := resgroup.(string)
		resourceGroup := vpcv1.ResourceGroupIdentity{
			ID: &resgroupstr,
		}
		createDedicatedHostGroupOptions.SetResourceGroup(&resourceGroup)
	}
	if zone, ok := d.GetOk("zone"); ok {
		zonestr := zone.(string)
		zoneidentity := vpcv1.ZoneIdentity{
			Name: &zonestr,
		}
		createDedicatedHostGroupOptions.SetZone(&zoneidentity)
	}

	dedicatedHostGroup, response, err := vpcClient.CreateDedicatedHostGroupWithContext(context, createDedicatedHostGroupOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateDedicatedHostGroupWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	d.SetId(*dedicatedHostGroup.ID)

	return resourceIbmIsDedicatedHostGroupRead(context, d, meta)
}

func resourceIbmIsDedicatedHostGroupMapToResourceGroupIdentity(resourceGroupIdentityMap map[string]interface{}) vpcv1.ResourceGroupIdentity {
	resourceGroupIdentity := vpcv1.ResourceGroupIdentity{}

	if resourceGroupIdentityMap["id"] != nil {
		resourceGroupIdentity.ID = core.StringPtr(resourceGroupIdentityMap["id"].(string))
	}

	return resourceGroupIdentity
}

func resourceIbmIsDedicatedHostGroupMapToResourceGroupIdentityByID(resourceGroupIdentityByIDMap map[string]interface{}) vpcv1.ResourceGroupIdentityByID {
	resourceGroupIdentityByID := vpcv1.ResourceGroupIdentityByID{}

	resourceGroupIdentityByID.ID = core.StringPtr(resourceGroupIdentityByIDMap["id"].(string))

	return resourceGroupIdentityByID
}

func resourceIbmIsDedicatedHostGroupRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getDedicatedHostGroupOptions := &vpcv1.GetDedicatedHostGroupOptions{}

	getDedicatedHostGroupOptions.SetID(d.Id())

	dedicatedHostGroup, response, err := vpcClient.GetDedicatedHostGroupWithContext(context, getDedicatedHostGroupOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetDedicatedHostGroupWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	if err = d.Set("class", dedicatedHostGroup.Class); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting class: %s", err))
	}
	if err = d.Set("family", dedicatedHostGroup.Family); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting family: %s", err))
	}
	if err = d.Set("name", dedicatedHostGroup.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}
	if dedicatedHostGroup.ResourceGroup != nil {
		resourceGroupID := *dedicatedHostGroup.ResourceGroup.ID
		if err = d.Set("resource_group", resourceGroupID); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_group: %s", err))
		}
	}
	if dedicatedHostGroup.Zone != nil {
		zoneName := *dedicatedHostGroup.Zone.Name
		if err = d.Set("zone", zoneName); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting zone: %s", err))
		}
	}
	if err = d.Set("created_at", dedicatedHostGroup.CreatedAt.String()); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
	}
	if err = d.Set("crn", dedicatedHostGroup.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting crn: %s", err))
	}
	dedicatedHosts := []map[string]interface{}{}
	for _, dedicatedHostsItem := range dedicatedHostGroup.DedicatedHosts {
		dedicatedHostsItemMap := resourceIbmIsDedicatedHostGroupDedicatedHostReferenceToMap(dedicatedHostsItem)
		dedicatedHosts = append(dedicatedHosts, dedicatedHostsItemMap)
	}
	if err = d.Set("dedicated_hosts", dedicatedHosts); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting dedicated_hosts: %s", err))
	}
	if err = d.Set("href", dedicatedHostGroup.Href); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting href: %s", err))
	}
	if err = d.Set("resource_type", dedicatedHostGroup.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_type: %s", err))
	}
	supportedInstanceProfiles := []map[string]interface{}{}
	for _, supportedInstanceProfilesItem := range dedicatedHostGroup.SupportedInstanceProfiles {
		supportedInstanceProfilesItemMap := resourceIbmIsDedicatedHostGroupInstanceProfileReferenceToMap(supportedInstanceProfilesItem)
		supportedInstanceProfiles = append(supportedInstanceProfiles, supportedInstanceProfilesItemMap)
	}
	if err = d.Set("supported_instance_profiles", supportedInstanceProfiles); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting supported_instance_profiles: %s", err))
	}

	return nil
}

func resourceIbmIsDedicatedHostGroupResourceGroupIdentityToMap(resourceGroupIdentity vpcv1.ResourceGroupIdentity) map[string]interface{} {
	resourceGroupIdentityMap := map[string]interface{}{}

	resourceGroupIdentityMap["id"] = resourceGroupIdentity.ID

	return resourceGroupIdentityMap
}

func resourceIbmIsDedicatedHostGroupResourceGroupIdentityByIDToMap(resourceGroupIdentityByID vpcv1.ResourceGroupIdentityByID) map[string]interface{} {
	resourceGroupIdentityByIDMap := map[string]interface{}{}

	resourceGroupIdentityByIDMap["id"] = resourceGroupIdentityByID.ID

	return resourceGroupIdentityByIDMap
}

func resourceIbmIsDedicatedHostGroupZoneIdentityToMap(zoneIdentity vpcv1.ZoneIdentity) map[string]interface{} {
	zoneIdentityMap := map[string]interface{}{}

	zoneIdentityMap["name"] = zoneIdentity.Name
	zoneIdentityMap["href"] = zoneIdentity.Href

	return zoneIdentityMap
}

func resourceIbmIsDedicatedHostGroupZoneIdentityByNameToMap(zoneIdentityByName vpcv1.ZoneIdentityByName) map[string]interface{} {
	zoneIdentityByNameMap := map[string]interface{}{}

	zoneIdentityByNameMap["name"] = zoneIdentityByName.Name

	return zoneIdentityByNameMap
}

func resourceIbmIsDedicatedHostGroupZoneIdentityByHrefToMap(zoneIdentityByHref vpcv1.ZoneIdentityByHref) map[string]interface{} {
	zoneIdentityByHrefMap := map[string]interface{}{}

	zoneIdentityByHrefMap["href"] = zoneIdentityByHref.Href

	return zoneIdentityByHrefMap
}

func resourceIbmIsDedicatedHostGroupDedicatedHostReferenceToMap(dedicatedHostReference vpcv1.DedicatedHostReference) map[string]interface{} {
	dedicatedHostReferenceMap := map[string]interface{}{}

	dedicatedHostReferenceMap["crn"] = dedicatedHostReference.CRN
	if dedicatedHostReference.Deleted != nil {
		DeletedMap := resourceIbmIsDedicatedHostGroupDedicatedHostReferenceDeletedToMap(*dedicatedHostReference.Deleted)
		dedicatedHostReferenceMap["deleted"] = []map[string]interface{}{DeletedMap}
	}
	dedicatedHostReferenceMap["href"] = dedicatedHostReference.Href
	dedicatedHostReferenceMap["id"] = dedicatedHostReference.ID
	dedicatedHostReferenceMap["name"] = dedicatedHostReference.Name
	dedicatedHostReferenceMap["resource_type"] = dedicatedHostReference.ResourceType

	return dedicatedHostReferenceMap
}

func resourceIbmIsDedicatedHostGroupDedicatedHostReferenceDeletedToMap(dedicatedHostReferenceDeleted vpcv1.Deleted) map[string]interface{} {
	dedicatedHostReferenceDeletedMap := map[string]interface{}{}

	dedicatedHostReferenceDeletedMap["more_info"] = dedicatedHostReferenceDeleted.MoreInfo

	return dedicatedHostReferenceDeletedMap
}

func resourceIbmIsDedicatedHostGroupInstanceProfileReferenceToMap(instanceProfileReference vpcv1.InstanceProfileReference) map[string]interface{} {
	instanceProfileReferenceMap := map[string]interface{}{}

	instanceProfileReferenceMap["href"] = instanceProfileReference.Href
	instanceProfileReferenceMap["name"] = instanceProfileReference.Name

	return instanceProfileReferenceMap
}

func resourceIbmIsDedicatedHostGroupUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	updateDedicatedHostGroupOptions := &vpcv1.UpdateDedicatedHostGroupOptions{}

	updateDedicatedHostGroupOptions.SetID(d.Id())

	hasChange := false

	if d.HasChange("name") {
		groupnamestr := d.Get("name").(string)
		dedicatedHostGroupPatchModel := vpcv1.DedicatedHostGroupPatch{
			Name: &groupnamestr,
		}
		dedicatedHostGroupPatch, err := dedicatedHostGroupPatchModel.AsPatch()
		if err != nil {
			log.Printf("[DEBUG] Error calling asPatch for DedicatedHostGroupPatch: %s", err)
			return diag.FromErr(err)
		}
		updateDedicatedHostGroupOptions.DedicatedHostGroupPatch = dedicatedHostGroupPatch
		hasChange = true
	}

	if hasChange {
		_, response, err := vpcClient.UpdateDedicatedHostGroupWithContext(context, updateDedicatedHostGroupOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateDedicatedHostGroupWithContext failed %s\n%s", err, response)
			return diag.FromErr(err)
		}
	}

	return resourceIbmIsDedicatedHostGroupRead(context, d, meta)
}

func resourceIbmIsDedicatedHostGroupDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getDedicatedHostGroupOptions := &vpcv1.GetDedicatedHostGroupOptions{}

	getDedicatedHostGroupOptions.SetID(d.Id())

	_, response, err := vpcClient.GetDedicatedHostGroupWithContext(context, getDedicatedHostGroupOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetDedicatedHostGroupWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	deleteDedicatedHostGroupOptions := &vpcv1.DeleteDedicatedHostGroupOptions{}

	deleteDedicatedHostGroupOptions.SetID(d.Id())

	response, err = vpcClient.DeleteDedicatedHostGroupWithContext(context, deleteDedicatedHostGroupOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteDedicatedHostGroupWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
