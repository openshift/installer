package vsphere

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/virtualdevice"
	"github.com/vmware/govmomi/pbm"
	types2 "github.com/vmware/govmomi/pbm/types"
	"github.com/vmware/govmomi/vapi/tags"
	"log"
	"strings"
)

const TAG_NAMESPACE = "http://www.vmware.com/storage/tag"
const TAG_PLACEMENT = "Tag based placement"

func resourceVmStoragePolicy() *schema.Resource {
	sch := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of the storage policy.",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Description of the storage policy.",
		},
		"tag_rules": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Tag rules to filter datastores to be used for placement of VMs.",
			Elem:        &schema.Resource{Schema: virtualdevice.VirtualMachineTagRulesSchema()},
		},
	}

	return &schema.Resource{
		Create: resourceVmStoragePolicyCreate,
		Read:   resourceVmStoragePolicyRead,
		Update: resourceVmStoragePolicyUpdate,
		Delete: resourceVmStoragePolicyDelete,
		Schema: sch,
	}
}

func resourceVmStoragePolicyCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Beginning create storage policy profile %s", d.Get("name").(string))
	client := meta.(*VSphereClient).vimClient
	rc := meta.(*VSphereClient).restClient

	pbmClient, err := pbm.NewClient(context.Background(), client.Client)
	if err != nil {
		return fmt.Errorf("error while creating pbm client %s", err)
	}

	tagsManager := tags.NewManager(rc)
	tagRules := d.Get("tag_rules").([]interface{})

	var capabilities []pbm.Capability

	for _, tagRule := range tagRules {

		tagCategory := tagRule.(map[string]interface{})["tag_category"].(string)
		_, err := tagCategoryByName(tagsManager, tagCategory)
		if err != nil {
			return fmt.Errorf("error while getting the tag %s %s", tagCategory, err)
		}

		tagValues := tagRule.(map[string]interface{})["tags"].([]interface{})
		tagValuesArr := make([]string, len(tagValues))
		for i, v := range tagValues {
			tagValuesArr[i] = fmt.Sprint(v)
		}
		tagsStr := strings.Join(tagValuesArr, ",")

		var includeTags string
		if !tagRule.(map[string]interface{})["include_datastores_with_tags"].(bool) {
			includeTags = "NOT"
		}

		var properties []pbm.Property
		properties = append(properties, pbm.Property{
			ID:       "com.vmware.storage.tag." + tagCategory + ".property",
			DataType: "Set",
			Value:    tagsStr,
			Operator: includeTags,
		})
		capability := pbm.Capability{
			ID:           tagCategory,
			Namespace:    TAG_NAMESPACE,
			PropertyList: properties,
		}
		capabilities = append(capabilities, capability)
	}

	name := d.Get("name").(string)
	description := d.Get("description").(string)

	capabilityProfileCreateSpec := pbm.CapabilityProfileCreateSpec{
		Name:           name,
		Description:    description,
		CapabilityList: capabilities,
		SubProfileName: TAG_PLACEMENT,
	}
	pbmCapabilityProfileCreateSpec, err := pbm.CreateCapabilityProfileSpec(capabilityProfileCreateSpec)
	if err != nil {
		return fmt.Errorf("error while creating storage policy spec %s", err)
	}

	profileID, err := pbmClient.CreateProfile(context.Background(), *pbmCapabilityProfileCreateSpec)
	if err != nil {
		return fmt.Errorf("error while creating storage policy %s", err)
	}

	d.SetId(profileID.UniqueId)
	log.Printf("[DEBUG] Storage policy create complete with Id %s", profileID.UniqueId)

	return resourceVmStoragePolicyRead(d, meta)
}

func resourceVmStoragePolicyRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Reading vm storage policy profile", resourceVSphereVirtualMachineIDString(d))
	client := meta.(*VSphereClient).vimClient
	pbmClient, err := pbm.NewClient(context.Background(), client.Client)
	if err != nil {
		return fmt.Errorf("error while creating pbm client %s", err)
	}
	profileId := types2.PbmProfileId{
		UniqueId: d.Id(),
	}
	pbmProfileIds := []types2.PbmProfileId{profileId}

	vmStoragePolicies, err := pbmClient.RetrieveContent(context.Background(), pbmProfileIds)
	if err != nil {
		if strings.Contains(err.Error(), "Profile not found") {
			log.Printf("[DEBUG] storage policy profile %s: Resource has been deleted", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error while reading vm storage policy profile with Id %s %s", d.Id(), err.Error())
	}
	if len(vmStoragePolicies) == 0 {
		d.SetId("")
		return nil
	}

	pbmCapabilityProfile := vmStoragePolicies[0].(*types2.PbmCapabilityProfile)
	d.SetId(pbmCapabilityProfile.ProfileId.UniqueId)
	d.Set("name", pbmCapabilityProfile.Name)
	d.Set("description", pbmCapabilityProfile.Description)

	var tagRules []map[string]interface{}
	if pbmCapabilityProfile.Constraints != nil {
		pbmSubProfileConstraints := pbmCapabilityProfile.Constraints.(*types2.PbmCapabilitySubProfileConstraints)
		if pbmSubProfileConstraints == nil {
			return nil
		}
		pbmSubProfiles := pbmSubProfileConstraints.SubProfiles
		for _, subProfile := range pbmSubProfiles {
			if subProfile.Name != TAG_PLACEMENT {
				continue
			}
			capabilities := subProfile.Capability
			for _, capability := range capabilities {

				if capability.Id.Namespace != TAG_NAMESPACE {
					continue
				}
				tagCategory := capability.Id.Id

				constraints := capability.Constraint
				if len(constraints) == 0 {
					continue
				}
				propertyInstances := constraints[0].PropertyInstance
				if len(propertyInstances) == 0 {
					continue
				}
				tagsSet := propertyInstances[0].Value.(types2.PbmCapabilityDiscreteSet)
				tagRule := make(map[string]interface{})
				tagRule["tag_category"] = tagCategory
				tagRule["tags"] = tagsSet.Values

				includeTags := true
				if propertyInstances[0].Operator == "NOT" {
					includeTags = false
				}
				tagRule["include_datastores_with_tags"] = includeTags

				tagRules = append(tagRules, tagRule)
			}
		}
	}
	d.Set("tag_rules", tagRules)
	return nil
}

func resourceVmStoragePolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Print("[DEBUG] :  Performing update")
	client := meta.(*VSphereClient).vimClient
	rc := meta.(*VSphereClient).restClient
	pbmClient, err := pbm.NewClient(context.Background(), client.Client)
	if err != nil {
		return fmt.Errorf("error while creating pbm client %s", err)
	}

	updateSpec := types2.PbmCapabilityProfileUpdateSpec{}

	updateSpec.Name = d.Get("name").(string)
	updateSpec.Description = d.Get("description").(string)
	log.Print("update spec ", updateSpec)

	if d.HasChange("tag_rules") {

		tagsManager := tags.NewManager(rc)
		tagRules := d.Get("tag_rules").([]interface{})
		var capabilities []pbm.Capability

		for _, tagRule := range tagRules {

			tagCategory := tagRule.(map[string]interface{})["tag_category"].(string)
			_, err := tagCategoryByName(tagsManager, tagCategory)
			if err != nil {
				return fmt.Errorf("error while getting the tag %s %s", tagCategory, err)
			}

			tagValues := tagRule.(map[string]interface{})["tags"].([]interface{})
			tagValuesArr := make([]string, len(tagValues))
			for i, v := range tagValues {
				tagValuesArr[i] = fmt.Sprint(v)
			}
			tagsStr := strings.Join(tagValuesArr, ",")

			var includeTags string
			if !tagRule.(map[string]interface{})["include_datastores_with_tags"].(bool) {
				includeTags = "NOT"
			}

			var properties []pbm.Property
			properties = append(properties, pbm.Property{
				ID:       "com.vmware.storage.tag." + tagCategory + ".property",
				DataType: "Set",
				Value:    tagsStr,
				Operator: includeTags,
			})
			capability := pbm.Capability{
				ID:           tagCategory,
				Namespace:    TAG_NAMESPACE,
				PropertyList: properties,
			}
			capabilities = append(capabilities, capability)
		}

		capabilityProfileCreateSpec := pbm.CapabilityProfileCreateSpec{
			CapabilityList: capabilities,
			SubProfileName: TAG_PLACEMENT,
		}
		pbmCapabilityProfileCreateSpec, err := pbm.CreateCapabilityProfileSpec(capabilityProfileCreateSpec)
		if err != nil {
			return fmt.Errorf("error while creating profile updatespec %s", err)
		}
		updateSpec.Constraints = pbmCapabilityProfileCreateSpec.Constraints
	}

	policyIdToUpdate := types2.PbmProfileId{
		UniqueId: d.Id(),
	}
	err = pbmClient.UpdateProfile(context.Background(), policyIdToUpdate, updateSpec)
	if err != nil {
		return fmt.Errorf("error while performing profile update for Id %s %s", d.Id(), err)

	}
	log.Print("[DEBUG] : update complete")
	return resourceVmStoragePolicyRead(d, meta)
}

func resourceVmStoragePolicyDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Performing create of VM storage policy with ID %s", d.Id())
	client := meta.(*VSphereClient).vimClient
	pbmClient, err := pbm.NewClient(context.Background(), client.Client)
	if err != nil {
		return fmt.Errorf("error while creating pbm client %s", err)
	}
	var policyIdsToDelete []types2.PbmProfileId
	policyIdsToDelete = append(policyIdsToDelete, types2.PbmProfileId{
		UniqueId: d.Id(),
	})

	_, err = pbmClient.DeleteProfile(context.Background(), policyIdsToDelete)
	if err != nil {
		return fmt.Errorf("error while deleting policy with ID %s %s", d.Id(), err)
	}
	d.SetId("")
	log.Printf("[DEBUG] %s: Delete complete", d.Id())
	return nil
}
