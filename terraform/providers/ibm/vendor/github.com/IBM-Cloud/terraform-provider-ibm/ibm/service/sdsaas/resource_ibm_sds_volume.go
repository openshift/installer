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

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/sds-go-sdk/sdsaasv1"
)

func ResourceIBMSdsVolume() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMSdsVolumeCreate,
		ReadContext:   resourceIBMSdsVolumeRead,
		UpdateContext: resourceIBMSdsVolumeUpdate,
		DeleteContext: resourceIBMSdsVolumeDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"sds_endpoint": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The endpoint to use for operations",
			},
			"hostnqnstring": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_sds_volume", "hostnqnstring"),
				Description:  "The host nqn.",
			},
			"capacity": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The capacity of the volume (in gigabytes).",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_sds_volume", "name"),
				Description:  "The name of the volume.",
			},
			"bandwidth": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum bandwidth (in megabits per second) for the volume.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the volume was created.",
			},
			"hosts": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of host details that volume is mapped to.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Unique identifer of the host.",
						},
						"host_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Unique name of the host.",
						},
						"host_nqn": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The NQN of the host configured in customer's environment.",
						},
					},
				},
			},
			"iops": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Iops The maximum I/O operations per second (IOPS) for this volume.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type of the volume.",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the volume.",
			},
			"status_reasons": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Reasons for the current status of the volume.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func getSDSConfigClient(meta interface{}, endpoint string) (*sdsaasv1.SdsaasV1, error) {
	sdsconfigClient, err := meta.(conns.ClientSession).SdsaasV1()
	if err != nil {
		return nil, err
	}
	url := conns.EnvFallBack([]string{"IBMCLOUD_SDS_ENDPOINT"}, endpoint)
	sdsconfigClient.Service.Options.URL = url
	return sdsconfigClient, nil
}

func ResourceIBMSdsVolumeValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "hostnqnstring",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^nqn\.\d{4}-\d{2}\.[a-z0-9-]+(?:\.[a-z0-9-]+)*:[a-zA-Z0-9.\-:]+$`,
			MinValueLength:             1,
			MaxValueLength:             200,
		},
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^.*$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_sds_volume", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMSdsVolumeCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	endpoint := d.Get("sds_endpoint").(string)
	sdsaasClient, err := getSDSConfigClient(meta, endpoint)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_sds_volume", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	volumeCreateOptions := &sdsaasv1.VolumeCreateOptions{}

	volumeCreateOptions.SetCapacity(int64(d.Get("capacity").(int)))
	if _, ok := d.GetOk("name"); ok {
		volumeCreateOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("hostnqnstring"); ok {
		volumeCreateOptions.SetHostnqnstring(d.Get("hostnqnstring").(string))
	}

	volume, _, err := sdsaasClient.VolumeCreateWithContext(context, volumeCreateOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("VolumeCreateWithContext failed: %s", err.Error()), "ibm_sds_volume", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*volume.ID)

	return resourceIBMSdsVolumeRead(context, d, meta)
}

func resourceIBMSdsVolumeRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	endpoint := d.Get("sds_endpoint").(string)
	sdsaasClient, err := getSDSConfigClient(meta, endpoint)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_sds_volume", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	volumeOptions := &sdsaasv1.VolumeOptions{}

	volumeOptions.SetVolumeID(d.Id())

	volume, response, err := sdsaasClient.VolumeWithContext(context, volumeOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("VolumeWithContext failed: %s", err.Error()), "ibm_sds_volume", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("capacity", flex.IntValue(volume.Capacity)); err != nil {
		err = fmt.Errorf("Error setting capacity: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_sds_volume", "read", "set-capacity").GetDiag()
	}
	if !core.IsNil(volume.Name) {
		if err = d.Set("name", volume.Name); err != nil {
			err = fmt.Errorf("Error setting name: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_sds_volume", "read", "set-name").GetDiag()
		}
	}
	if !core.IsNil(volume.Bandwidth) {
		if err = d.Set("bandwidth", flex.IntValue(volume.Bandwidth)); err != nil {
			err = fmt.Errorf("Error setting bandwidth: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_sds_volume", "read", "set-bandwidth").GetDiag()
		}
	}
	if !core.IsNil(volume.CreatedAt) {
		if err = d.Set("created_at", volume.CreatedAt); err != nil {
			err = fmt.Errorf("Error setting created_at: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_sds_volume", "read", "set-created_at").GetDiag()
		}
	}
	if !core.IsNil(volume.Hosts) {
		hosts := []map[string]interface{}{}
		for _, hostsItem := range volume.Hosts {
			hostsItemMap, err := ResourceIBMSdsVolumeHostMappingToMap(&hostsItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_sds_volume", "read", "hosts-to-map").GetDiag()
			}
			hosts = append(hosts, hostsItemMap)
		}
		if err = d.Set("hosts", hosts); err != nil {
			err = fmt.Errorf("Error setting hosts: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_sds_volume", "read", "set-hosts").GetDiag()
		}
	}
	if !core.IsNil(volume.Iops) {
		if err = d.Set("iops", flex.IntValue(volume.Iops)); err != nil {
			err = fmt.Errorf("Error setting iops: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_sds_volume", "read", "set-iops").GetDiag()
		}
	}
	if !core.IsNil(volume.ResourceType) {
		if err = d.Set("resource_type", volume.ResourceType); err != nil {
			err = fmt.Errorf("Error setting resource_type: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_sds_volume", "read", "set-resource_type").GetDiag()
		}
	}
	if !core.IsNil(volume.Status) {
		if err = d.Set("status", volume.Status); err != nil {
			err = fmt.Errorf("Error setting status: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_sds_volume", "read", "set-status").GetDiag()
		}
	}
	if !core.IsNil(volume.StatusReasons) {
		if err = d.Set("status_reasons", volume.StatusReasons); err != nil {
			err = fmt.Errorf("Error setting status_reasons: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_sds_volume", "read", "set-status_reasons").GetDiag()
		}
	}

	return nil
}

func resourceIBMSdsVolumeUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	endpoint := d.Get("sds_endpoint").(string)
	sdsaasClient, err := getSDSConfigClient(meta, endpoint)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_sds_volume", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	volumeUpdateOptions := &sdsaasv1.VolumeUpdateOptions{}

	volumeUpdateOptions.SetVolumeID(d.Id())

	hasChange := false

	patchVals := &sdsaasv1.VolumePatch{}
	if d.HasChange("capacity") {
		newCapacity := int64(d.Get("capacity").(int))
		patchVals.Capacity = &newCapacity
		hasChange = true
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
		volumeUpdateOptions.VolumePatch = ResourceIBMSdsVolumeVolumePatchAsPatch(patchVals, d)

		_, _, err = sdsaasClient.VolumeUpdateWithContext(context, volumeUpdateOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("VolumeUpdateWithContext failed: %s", err.Error()), "ibm_sds_volume", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIBMSdsVolumeRead(context, d, meta)
}

func resourceIBMSdsVolumeDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	endpoint := d.Get("sds_endpoint").(string)
	sdsaasClient, err := getSDSConfigClient(meta, endpoint)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_sds_volume", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	volumeDeleteOptions := &sdsaasv1.VolumeDeleteOptions{}

	volumeDeleteOptions.SetVolumeID(d.Id())

	_, err = sdsaasClient.VolumeDeleteWithContext(context, volumeDeleteOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("VolumeDeleteWithContext failed: %s", err.Error()), "ibm_sds_volume", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIBMSdsVolumeHostMappingToMap(model *sdsaasv1.HostMapping) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.HostID != nil {
		modelMap["host_id"] = *model.HostID
	}
	if model.HostName != nil {
		modelMap["host_name"] = *model.HostName
	}
	if model.HostNqn != nil {
		modelMap["host_nqn"] = *model.HostNqn
	}
	return modelMap, nil
}

func ResourceIBMSdsVolumeVolumePatchAsPatch(patchVals *sdsaasv1.VolumePatch, d *schema.ResourceData) map[string]interface{} {
	patch, _ := patchVals.AsPatch()
	var path string

	path = "capacity"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["capacity"] = nil
	}
	path = "name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["name"] = nil
	}

	return patch
}
