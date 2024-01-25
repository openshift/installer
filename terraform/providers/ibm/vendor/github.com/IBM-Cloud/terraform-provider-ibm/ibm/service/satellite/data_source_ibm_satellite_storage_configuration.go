// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package satellite

import (
	"fmt"

	"github.com/IBM-Cloud/container-services-go-sdk/kubernetesserviceapiv1"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMSatelliteStorageConfiguration() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMSatelliteStorageConfigurationRead,

		Schema: map[string]*schema.Schema{
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Location Name.",
			},
			"config_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Storage Configuration.",
			},
			"config_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version of the Storage Configuration.",
			},
			"storage_template_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Storage Template Name.",
			},
			"storage_template_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Storage Template Version.",
			},
			"user_config_parameters": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The storage configuration parameters depending on the storage template.",
			},
			"storage_class_parameters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Description: "The List of Storage Class Parameters as a list of a  Map of string key-value.",
				},
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Universally Unique IDentifier (UUID) of the Storage Configuration.",
			},
		},
	}
}

func dataSourceIBMSatelliteStorageConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return err
	}

	satLocation := d.Get("location").(string)
	d.Set("location", satLocation)
	getSatelliteLocation := &kubernetesserviceapiv1.GetSatelliteLocationOptions{
		Controller: &satLocation,
	}
	_, _, err = satClient.GetSatelliteLocation(getSatelliteLocation)
	if err != nil {
		return fmt.Errorf("[ERROR] Location not found! - %v", err)
	}

	storageConfigName := d.Get("config_name").(string)
	getStorageConfigurationOptions := &kubernetesserviceapiv1.GetStorageConfigurationOptions{
		Name: &storageConfigName,
	}
	result, _, err := satClient.GetStorageConfiguration(getStorageConfigurationOptions)
	if err != nil {
		return err
	}

	d.Set("config_name", *result.ConfigName)
	d.Set("config_version", *result.ConfigVersion)
	d.Set("storage_template_name", *result.StorageTemplateName)
	d.Set("storage_template_version", *result.StorageTemplateVersion)
	d.Set("user_config_parameters", result.UserConfigParameters)
	d.Set("storage_class_parameters", result.StorageClassParameters)
	d.Set("uuid", *result.UUID)
	d.SetId(*result.UUID + "/" + satLocation)

	return nil
}
