package vsphere

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/virtualdevice"
	"github.com/vmware/govmomi/pbm"
	types2 "github.com/vmware/govmomi/pbm/types"
	"github.com/vmware/govmomi/vapi/tags"
)

const TagNamespace = "http://www.vmware.com/storage/tag"
const TagPlacement = "Tag based placement"

func resourceVMStoragePolicy() *schema.Resource {
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
		Create: resourceVMStoragePolicyCreate,
		Read:   resourceVMStoragePolicyRead,
		Update: resourceVMStoragePolicyUpdate,
		Delete: resourceVMStoragePolicyDelete,
		Schema: sch,
	}
}

func resourceVMStoragePolicyCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Beginning create storage policy profile %s", d.Get("name").(string))
	client := meta.(*Client).vimClient
	rc := meta.(*Client).restClient

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
			Namespace:    TagNamespace,
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
		SubProfileName: TagPlacement,
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

	return resourceVMStoragePolicyRead(d, meta)
}

func resourceVMStoragePolicyRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Reading vm storage policy profile", resourceVSphereVirtualMachineIDString(d))
	client := meta.(*Client).vimClient
	pbmClient, err := pbm.NewClient(context.Background(), client.Client)
	if err != nil {
		return fmt.Errorf("error while creating pbm client %s", err)
	}
	profileID := types2.PbmProfileId{
		UniqueId: d.Id(),
	}
	pbmProfileIds := []types2.PbmProfileId{profileID}

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
	_ = d.Set("name", pbmCapabilityProfile.Name)
	_ = d.Set("description", pbmCapabilityProfile.Description)

	var tagRules []map[string]interface{}
	if pbmCapabilityProfile.Constraints != nil {
		pbmSubProfileConstraints := pbmCapabilityProfile.Constraints.(*types2.PbmCapabilitySubProfileConstraints)
		if pbmSubProfileConstraints == nil {
			return nil
		}
		pbmSubProfiles := pbmSubProfileConstraints.SubProfiles
		for _, subProfile := range pbmSubProfiles {
			if subProfile.Name != TagPlacement {
				continue
			}
			capabilities := subProfile.Capability
			for _, capability := range capabilities {
				if capability.Id.Namespace != TagNamespace {
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
	_ = d.Set("tag_rules", tagRules)
	return nil
}

func resourceVMStoragePolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Print("[DEBUG] :  Performing update")
	client := meta.(*Client).vimClient
	rc := meta.(*Client).restClient
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
				Namespace:    TagNamespace,
				PropertyList: properties,
			}
			capabilities = append(capabilities, capability)
		}

		capabilityProfileCreateSpec := pbm.CapabilityProfileCreateSpec{
			CapabilityList: capabilities,
			SubProfileName: TagPlacement,
		}
		pbmCapabilityProfileCreateSpec, err := pbm.CreateCapabilityProfileSpec(capabilityProfileCreateSpec)
		if err != nil {
			return fmt.Errorf("error while creating profile updatespec %s", err)
		}
		updateSpec.Constraints = pbmCapabilityProfileCreateSpec.Constraints
	}

	policyIDToUpdate := types2.PbmProfileId{
		UniqueId: d.Id(),
	}
	err = pbmClient.UpdateProfile(context.Background(), policyIDToUpdate, updateSpec)
	if err != nil {
		return fmt.Errorf("error while performing profile update for Id %s %s", d.Id(), err)
	}
	log.Print("[DEBUG] : update complete")
	return resourceVMStoragePolicyRead(d, meta)
}

func resourceVMStoragePolicyDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Performing create of VM storage policy with ID %s", d.Id())
	client := meta.(*Client).vimClient
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
