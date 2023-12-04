// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package satellite

import (
	"fmt"

	"github.com/IBM-Cloud/container-services-go-sdk/kubernetesserviceapiv1"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMSatelliteStorageAssignment() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMSatelliteStorageAssignmentRead,

		Schema: map[string]*schema.Schema{
			"assignment_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the Assignment.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Universally Unique IDentifier (UUID) of the Assignment.",
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Owner of the Assignment.",
			},
			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "One or more cluster groups on which you want to apply the configuration. Note that at least one cluster group is required. ",
			},
			"cluster": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the Satellite cluster or Service Cluster that you want to apply the configuration to.",
			},
			"svc_cluster": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the Service Cluster that you applied the configuration to.",
			},
			"sat_cluster": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the Satellite cluster that you applied the configuration to.",
			},
			"config": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Storage Configuration Name or ID.",
			},
			"config_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Universally Unique IDentifier (UUID) of the Storage Configuration.",
			},
			"config_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Storage Configuration Version.",
			},
			"config_version_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Universally Unique IDentifier (UUID) of the Storage Configuration Version.",
			},
			"assignment_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Type of Assignment.",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Time of Creation of the Assignment.",
			},
			"rollout_success_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The Rollout Success Count of the Assignment.",
			},
			"rollout_error_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The Rollout Error Count of the Assignment.",
			},
			"is_assignment_upgrade_available": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether an Upgrade is Available for the Assignment.",
			},
		},
	}
}

func dataSourceIBMSatelliteStorageAssignmentRead(d *schema.ResourceData, meta interface{}) error {

	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return err
	}

	uuid := d.Get("uuid").(string)
	getAssignmentOptions := &kubernetesserviceapiv1.GetAssignmentOptions{
		UUID: &uuid,
	}

	result, _, err := satClient.GetAssignment(getAssignmentOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error getting Assignment of UUID %s - %v", uuid, err)
	}

	d.SetId(uuid + "/" + *result.Name)
	d.Set("assignment_name", *result.Name)
	d.Set("uuid", *result.UUID)
	d.Set("owner", *result.Owner.Name)
	if result.Groups != nil {
		d.Set("groups", result.Groups)
	}
	if result.Cluster != nil {
		d.Set("cluster", *result.Cluster)
	}
	if result.SatSvcClusterID != nil {
		d.Set("svc_cluster", *result.SatSvcClusterID)
	}
	if result.Satcluster != nil {
		d.Set("sat_cluster", *result.Satcluster)
	}
	if result.ChannelName != nil {
		d.Set("config", *result.ChannelName)
	}
	if result.ChannelUUID != nil {
		d.Set("config_uuid", *result.ChannelUUID)
	}
	if result.Version != nil {
		d.Set("config_version", *result.Version)
	}
	if result.VersionUUID != nil {
		d.Set("config_version_uuid", *result.VersionUUID)
	}
	if result.SubscriptionType != nil {
		d.Set("assignment_type", *result.SubscriptionType)
	}
	if result.Created != nil {
		d.Set("created", *result.Created)
	}
	if result.IsAssignmentUpgradeAvailable != nil {
		d.Set("is_assignment_upgrade_available", *result.IsAssignmentUpgradeAvailable)
	}
	if result.RolloutStatus != nil {
		d.Set("rollout_success_count", *result.RolloutStatus.SuccessCount)
		d.Set("rollout_error_count", *result.RolloutStatus.ErrorCount)
	}
	return nil
}
