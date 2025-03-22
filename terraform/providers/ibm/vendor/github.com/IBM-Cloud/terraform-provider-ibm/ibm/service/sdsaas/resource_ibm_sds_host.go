// Copyright IBM Corp. 2025 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.96.0-d6dec9d7-20241008-212902
 */

package sdsaas

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/sds-go-sdk/sdsaasv1"
)

func ResourceIBMSdsHost() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMSdsHostCreate,
		ReadContext:   resourceIBMSdsHostRead,
		UpdateContext: resourceIBMSdsHostUpdate,
		DeleteContext: resourceIBMSdsHostDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"sds_endpoint": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The endpoint to use for operations",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_sds_host", "name"),
				Description:  "The name for this host. The name must not be used by another host.  If unspecified, the name will be a hyphenated list of randomly-selected words.",
			},
			"nqn": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_sds_host", "nqn"),
				Description:  "The NQN of the host configured in customer's environment.",
			},
			"volumes": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The host-to-volume map.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The current status of a volume/host mapping attempt.",
						},
						"volume_id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The volume ID that needs to be mapped with a host.",
						},
						"volume_name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The volume name.",
						},
						"storage_identifiers": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Storage network and ID information associated with a volume/host mapping.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The storage ID associated with a volume/host mapping.",
									},
									"namespace_id": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The namespace ID associated with a volume/host mapping.",
									},
									"namespace_uuid": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The namespace UUID associated with a volume/host mapping.",
									},
									"network_info": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The IP and port for volume/host mappings.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"gateway_ip": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Network information for volume/host mappings.",
												},
												"port": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Network information for volume/host mappings.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the host was created.",
			},
		},
	}
}

func ResourceIBMSdsHostValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$`,
			MinValueLength:             1,
			MaxValueLength:             64,
		},
		validate.ValidateSchema{
			Identifier:                 "nqn",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^nqn\.\d{4}-\d{2}\.[a-z0-9-]+(?:\.[a-z0-9-]+)*:[a-zA-Z0-9.\-:]+$`,
			MinValueLength:             1,
			MaxValueLength:             64,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_sds_host", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMSdsHostCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	endpoint := d.Get("sds_endpoint").(string)
	sdsaasClient, err := getSDSConfigClient(meta, endpoint)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_sds_host", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	hostCreateOptions := &sdsaasv1.HostCreateOptions{}

	hostCreateOptions.SetNqn(d.Get("nqn").(string))
	if _, ok := d.GetOk("name"); ok {
		hostCreateOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("volumes"); ok {
		var volumes []sdsaasv1.VolumeMappingIdentity
		for _, v := range d.Get("volumes").([]interface{}) {
			value := v.(map[string]interface{})
			volumesItem, err := ResourceIBMSdsHostMapToVolumeMappingIdentity(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_sds_host", "create", "parse-volumes").GetDiag()
			}
			volumes = append(volumes, *volumesItem)
		}
		hostCreateOptions.SetVolumes(volumes)
	}

	host, _, err := sdsaasClient.HostCreateWithContext(context, hostCreateOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("HostCreateWithContext failed: %s", err.Error()), "ibm_sds_host", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*host.ID)

	return resourceIBMSdsHostRead(context, d, meta)
}

func resourceIBMSdsHostRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	endpoint := d.Get("sds_endpoint").(string)
	sdsaasClient, err := getSDSConfigClient(meta, endpoint)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_sds_host", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	hostOptions := &sdsaasv1.HostOptions{}

	hostOptions.SetHostID(d.Id())

	host, response, err := sdsaasClient.HostWithContext(context, hostOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("HostWithContext failed: %s", err.Error()), "ibm_sds_host", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if !core.IsNil(host.Name) {
		if err = d.Set("name", host.Name); err != nil {
			err = fmt.Errorf("Error setting name: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_sds_host", "read", "set-name").GetDiag()
		}
	}
	if err = d.Set("nqn", host.Nqn); err != nil {
		err = fmt.Errorf("Error setting nqn: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_sds_host", "read", "set-nqn").GetDiag()
	}
	if !core.IsNil(host.Volumes) {
		volumes := []map[string]interface{}{}
		for _, volumesItem := range host.Volumes {
			volumesItemMap, err := ResourceIBMSdsHostVolumeMappingReferenceToMap(&volumesItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_sds_host", "read", "volumes-to-map").GetDiag()
			}
			volumes = append(volumes, volumesItemMap)
		}
		if err = d.Set("volumes", volumes); err != nil {
			err = fmt.Errorf("Error setting volumes: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_sds_host", "read", "set-volumes").GetDiag()
		}
	}
	if !core.IsNil(host.CreatedAt) {
		if err = d.Set("created_at", host.CreatedAt); err != nil {
			err = fmt.Errorf("Error setting created_at: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_sds_host", "read", "set-created_at").GetDiag()
		}
	}

	return nil
}

func resourceIBMSdsHostUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	endpoint := d.Get("sds_endpoint").(string)
	sdsaasClient, err := getSDSConfigClient(meta, endpoint)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_sds_host", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	hostUpdateOptions := &sdsaasv1.HostUpdateOptions{}

	hostUpdateOptions.SetHostID(d.Id())

	hasChange := false

	patchVals := &sdsaasv1.HostPatch{}
	if d.HasChange("name") {
		newName := d.Get("name").(string)
		patchVals.Name = &newName
		hasChange = true
	}

	if hasChange {
		// Fields with `nil` values are omitted from the generic map,
		// so we need to re-add them to support removing arguments
		// in merge-patch operations sent to the service.
		hostUpdateOptions.HostPatch = ResourceIBMSdsHostHostPatchAsPatch(patchVals, d)

		_, _, err = sdsaasClient.HostUpdateWithContext(context, hostUpdateOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("HostUpdateWithContext failed: %s", err.Error()), "ibm_sds_host", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIBMSdsHostRead(context, d, meta)
}

func resourceIBMSdsHostDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	endpoint := d.Get("sds_endpoint").(string)
	sdsaasClient, err := getSDSConfigClient(meta, endpoint)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_sds_host", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	hostDeleteOptions := &sdsaasv1.HostDeleteOptions{}

	hostDeleteOptions.SetHostID(d.Id())

	_, err = sdsaasClient.HostDeleteWithContext(context, hostDeleteOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("HostDeleteWithContext failed: %s", err.Error()), "ibm_sds_host", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIBMSdsHostMapToVolumeMappingIdentity(modelMap map[string]interface{}) (*sdsaasv1.VolumeMappingIdentity, error) {
	model := &sdsaasv1.VolumeMappingIdentity{}
	model.VolumeID = core.StringPtr(modelMap["volume_id"].(string))
	return model, nil
}

func ResourceIBMSdsHostVolumeMappingReferenceToMap(model *sdsaasv1.VolumeMappingReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Status != nil {
		modelMap["status"] = *model.Status
	}
	modelMap["volume_id"] = *model.VolumeID
	modelMap["volume_name"] = *model.VolumeName
	if model.StorageIdentifiers != nil {
		storageIdentifiersMap, err := ResourceIBMSdsHostStorageIdentifiersReferenceToMap(model.StorageIdentifiers)
		if err != nil {
			return modelMap, err
		}
		modelMap["storage_identifiers"] = []map[string]interface{}{storageIdentifiersMap}
	}
	return modelMap, nil
}

func ResourceIBMSdsHostStorageIdentifiersReferenceToMap(model *sdsaasv1.StorageIdentifiersReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.NamespaceID != nil {
		modelMap["namespace_id"] = flex.IntValue(model.NamespaceID)
	}
	if model.NamespaceUUID != nil {
		modelMap["namespace_uuid"] = *model.NamespaceUUID
	}
	if model.NetworkInfo != nil {
		networkInfo := []map[string]interface{}{}
		for _, networkInfoItem := range model.NetworkInfo {
			networkInfoItemMap, err := ResourceIBMSdsHostNetworkInfoReferenceToMap(&networkInfoItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			networkInfo = append(networkInfo, networkInfoItemMap)
		}
		modelMap["network_info"] = networkInfo
	}
	return modelMap, nil
}

func ResourceIBMSdsHostNetworkInfoReferenceToMap(model *sdsaasv1.NetworkInfoReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.GatewayIP != nil {
		modelMap["gateway_ip"] = *model.GatewayIP
	}
	if model.Port != nil {
		modelMap["port"] = flex.IntValue(model.Port)
	}
	return modelMap, nil
}

func ResourceIBMSdsHostHostPatchAsPatch(patchVals *sdsaasv1.HostPatch, d *schema.ResourceData) map[string]interface{} {
	patch, _ := patchVals.AsPatch()
	var path string

	path = "name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["name"] = nil
	}

	return patch
}
