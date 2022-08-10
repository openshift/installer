// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/scc-go-sdk/v3/posturemanagementv2"
)

func DataSourceIBMSccPostureCollector() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMSccPostureCollectorRead,

		Schema: map[string]*schema.Schema{
			"collector_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id for the given API.",
			},
			"display_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user-friendly name of the collector.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the collector.",
			},
			"public_key": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The public key of the collector.Will be used for ssl communciation between collector and orchestrator .This will be populated when collector is installed.",
			},
			"last_heartbeat": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Stores the heartbeat time of a controller . This value exists when collector is installed and running.",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of collector.",
			},
			"collector_version": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The collector version. This field is populated when collector is installed.",
			},
			"image_version": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The image version of the collector. This field is populated when collector is installed. \".",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the collector.",
			},
			"created_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the user that created the collector.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ISO Date/Time the collector was created.",
			},
			"updated_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the user that modified the collector.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ISO Date/Time the collector was modified.",
			},
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Identifies whether the collector is enabled or not(deleted).",
			},
			"registration_code": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The registration code of the collector.This is will be used for initial authentication during installation of collector.",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the collector.",
			},
			"credential_public_key": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The credential public key.",
			},
			"failure_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of times the collector has failed.",
			},
			"approved_local_gateway_ip": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The approved local gateway ip of the collector. This field will be populated only when collector is installed.",
			},
			"approved_internet_gateway_ip": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The approved internet gateway ip of the collector. This field will be populated only when collector is installed.",
			},
			"last_failed_local_gateway_ip": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The failed local gateway ip. This field will be populated only when collector is installed.",
			},
			"reset_reason": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The reason for the collector reset .User resets the collector with a reason for reset. The reason entered by the user is saved in this field .",
			},
			"hostname": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The collector host name. This field will be populated when collector is installed.This will have fully qualified domain name.",
			},
			"install_path": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The installation path of the collector. This field will be populated when collector is installed.The value will be folder path.",
			},
			"use_private_endpoint": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the collector should use a public or private endpoint. This value is generated based on is_public field value during collector creation. If is_public is set to true, this value will be false.",
			},
			"managed_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The entity that manages the collector.",
			},
			"trial_expiry": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The trial expiry. This holds the expiry date of registration_code. This field will be populated when collector is installed.",
			},
			"last_failed_internet_gateway_ip": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The failed internet gateway ip of the collector.",
			},
			"status_description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The collector status.",
			},
			"reset_time": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ISO Date/Time of the collector reset. This value will be populated when a collector is reset. The data-time when the reset event is occured is captured in this field.",
			},
			"is_public": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Determines whether the collector endpoint is accessible on a public network.If set to `true`, the collector connects to resources in your account over a public network. If set to `false`, the collector connects to resources by using a private IP that is accessible only through the IBM Cloud private network.",
			},
			"is_ubi_image": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Determines whether the collector has a Ubi image.",
			},
		},
	}
}

func dataSourceIBMSccPostureCollectorRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	postureManagementClient, err := meta.(conns.ClientSession).PostureManagementV2()
	if err != nil {
		return diag.FromErr(err)
	}

	CollectorOptions := &posturemanagementv2.GetCollectorOptions{}
	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error getting userDetails %s", err))
	}

	accountID := userDetails.UserAccount
	CollectorOptions.SetAccountID(accountID)
	CollectorOptions.SetID(d.Get("collector_id").(string))

	collector, response, err := postureManagementClient.GetCollectorWithContext(context, CollectorOptions)
	if err != nil {
		log.Printf("[DEBUG] GetCollectorWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetCollectorWithContext failed %s\n%s", err, response))
	}

	d.SetId(*collector.ID)
	if err = d.Set("display_name", collector.DisplayName); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting display_name: %s", err))
	}
	if err = d.Set("name", collector.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("public_key", collector.PublicKey); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting public_key: %s", err))
	}
	if err = d.Set("last_heartbeat", flex.DateTimeToString(collector.LastHeartbeat)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting last_heartbeat: %s", err))
	}
	if err = d.Set("status", collector.Status); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting status: %s", err))
	}
	if err = d.Set("collector_version", collector.CollectorVersion); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting collector_version: %s", err))
	}
	if err = d.Set("image_version", collector.ImageVersion); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting image_version: %s", err))
	}
	if err = d.Set("description", collector.Description); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
	}
	if err = d.Set("created_by", collector.CreatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_by: %s", err))
	}
	if err = d.Set("created_at", flex.DateTimeToString(collector.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("updated_by", collector.UpdatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_by: %s", err))
	}
	if err = d.Set("updated_at", flex.DateTimeToString(collector.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}
	if err = d.Set("enabled", collector.Enabled); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting enabled: %s", err))
	}
	if err = d.Set("registration_code", collector.RegistrationCode); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting registration_code: %s", err))
	}
	if err = d.Set("type", collector.Type); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting type: %s", err))
	}
	if err = d.Set("credential_public_key", collector.CredentialPublicKey); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting credential_public_key: %s", err))
	}
	if err = d.Set("failure_count", flex.IntValue(collector.FailureCount)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting failure_count: %s", err))
	}
	if err = d.Set("approved_local_gateway_ip", collector.ApprovedLocalGatewayIP); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting approved_local_gateway_ip: %s", err))
	}
	if err = d.Set("approved_internet_gateway_ip", collector.ApprovedInternetGatewayIP); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting approved_internet_gateway_ip: %s", err))
	}
	if err = d.Set("last_failed_local_gateway_ip", collector.LastFailedLocalGatewayIP); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting last_failed_local_gateway_ip: %s", err))
	}
	if err = d.Set("reset_reason", collector.ResetReason); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting reset_reason: %s", err))
	}
	if err = d.Set("hostname", collector.Hostname); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting hostname: %s", err))
	}
	if err = d.Set("install_path", collector.InstallPath); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting install_path: %s", err))
	}
	if err = d.Set("use_private_endpoint", collector.UsePrivateEndpoint); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting use_private_endpoint: %s", err))
	}
	if err = d.Set("managed_by", collector.ManagedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting managed_by: %s", err))
	}
	if err = d.Set("trial_expiry", flex.DateTimeToString(collector.TrialExpiry)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting trial_expiry: %s", err))
	}
	if err = d.Set("last_failed_internet_gateway_ip", collector.LastFailedInternetGatewayIP); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting last_failed_internet_gateway_ip: %s", err))
	}
	if err = d.Set("status_description", collector.StatusDescription); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting status_description: %s", err))
	}
	if err = d.Set("reset_time", flex.DateTimeToString(collector.ResetTime)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting reset_time: %s", err))
	}
	if err = d.Set("is_public", collector.IsPublic); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting is_public: %s", err))
	}
	if err = d.Set("is_ubi_image", collector.IsUbiImage); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting is_ubi_image: %s", err))
	}

	return nil
}
