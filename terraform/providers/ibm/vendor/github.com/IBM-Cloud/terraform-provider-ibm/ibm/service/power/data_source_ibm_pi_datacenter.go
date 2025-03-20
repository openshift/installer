// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPIDatacenter() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIDatacenterRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description: "The GUID of the service instance associated with an account.",
				Optional:    true,
				Type:        schema.TypeString,
			},
			Arg_DatacenterZone: {
				Description:  "Datacenter zone you want to retrieve.",
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_CapabilityDetails: {
				Computed:    true,
				Description: "Additional Datacenter Capability Details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_DisasterRecovery: {
							Computed:    true,
							Description: "Disaster Recovery Information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									Attr_AsynchronousReplication: {
										Computed:    true,
										Description: "Asynchronous Replication Target Information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												Attr_Enabled: {
													Computed:    true,
													Description: "Service Enabled.",
													Type:        schema.TypeBool,
												},
												Attr_TargetLocations: {
													Computed:    true,
													Description: "List of all replication targets.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															Attr_Region: {
																Computed:    true,
																Description: "regionZone of replication site.",
																Type:        schema.TypeString,
															},
															Attr_Status: {
																Computed:    true,
																Description: "the replication site is active / down.",
																Type:        schema.TypeString,
															},
														},
													},
													Type: schema.TypeList,
												},
											},
										},
										Type: schema.TypeList,
									},
									Attr_SynchronousReplication: {
										Computed:    true,
										Description: "Synchronous Replication Target Information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												Attr_Enabled: {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Service Enabled.",
												},
												Attr_TargetLocations: {
													Computed:    true,
													Description: "List of all replication targets.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															Attr_Region: {
																Computed:    true,
																Description: "regionZone of replication site.",
																Type:        schema.TypeString,
															},
															Attr_Status: {
																Computed:    true,
																Description: "the replication site is active / down.",
																Type:        schema.TypeString,
															},
														},
													},
													Type: schema.TypeList,
												},
											},
										},
										Type: schema.TypeList,
									},
								},
							},
							Type: schema.TypeList,
						},
						Attr_SupportedSystems: {
							Computed:    true,
							Description: "Datacenter System Types Information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									Attr_Dedicated: {
										Computed:    true,
										Description: "List of all available dedicated host types.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Type: schema.TypeList,
									},
									Attr_General: {
										Computed:    true,
										Description: "List of all available host types.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Type: schema.TypeList,
									},
								},
							},
							Type: schema.TypeList,
						},
					},
				},
				Type: schema.TypeList,
			},
			Attr_DatacenterCapabilities: {
				Computed:    true,
				Description: "Datacenter Capabilities.",
				Elem: &schema.Schema{
					Type: schema.TypeBool,
				},
				Type: schema.TypeMap,
			},
			Attr_DatacenterHref: {
				Computed:    true,
				Description: "Datacenter href.",
				Type:        schema.TypeString,
			},
			Attr_DatacenterLocation: {
				Computed:    true,
				Description: "Datacenter location.",
				Type:        schema.TypeMap,
			},
			Attr_DatacenterStatus: {
				Computed:    true,
				Description: "Datacenter status, active,maintenance or down.",
				Type:        schema.TypeString,
			},
			Attr_DatacenterType: {
				Computed:    true,
				Description: "Datacenter type, off-premises or on-premises.",
				Type:        schema.TypeString,
			},
		},
	}
}

func dataSourceIBMPIDatacenterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	datacenterZone := sess.Options.Zone
	if region, ok := d.GetOk(Arg_DatacenterZone); ok {
		datacenterZone = region.(string)
	}

	cloudInstanceID := ""
	if cloudInstance, ok := d.GetOk(Arg_CloudInstanceID); ok {
		cloudInstanceID = cloudInstance.(string)
	}
	client := instance.NewIBMPIDatacenterClient(ctx, sess, cloudInstanceID)

	dcData, err := client.Get(datacenterZone)
	if err != nil {
		return diag.FromErr(err)
	}
	var genID, _ = uuid.GenerateUUID()
	d.SetId(genID)
	d.Set(Attr_DatacenterCapabilities, dcData.Capabilities)
	dclocation := map[string]interface{}{
		Attr_Region: *dcData.Location.Region,
		Attr_Type:   *dcData.Location.Type,
		Attr_URL:    *dcData.Location.URL,
	}
	d.Set(Attr_DatacenterHref, dcData.Href)
	d.Set(Attr_DatacenterLocation, flex.Flatten(dclocation))
	d.Set(Attr_DatacenterStatus, dcData.Status)
	d.Set(Attr_DatacenterType, dcData.Type)
	capabilityDetails := make([]map[string]interface{}, 0, 10)
	if dcData.CapabilitiesDetails != nil {
		capabilityDetailsMap, err := capabilityDetailsToMap(dcData.CapabilitiesDetails)
		if err != nil {
			return diag.FromErr(err)
		}
		capabilityDetails = append(capabilityDetails, capabilityDetailsMap)
	}
	d.Set(Attr_CapabilityDetails, capabilityDetails)

	return nil
}

func capabilityDetailsToMap(cd *models.CapabilitiesDetails) (map[string]interface{}, error) {
	capabilityDetailsMap := make(map[string]interface{})
	disasterRecoveryMap := disasterRecoveryToMap(cd.DisasterRecovery)
	capabilityDetailsMap[Attr_DisasterRecovery] = []map[string]interface{}{disasterRecoveryMap}

	supportedSystemsMap := make(map[string]interface{})
	supportedSystemsMap[Attr_Dedicated] = cd.SupportedSystems.Dedicated
	supportedSystemsMap[Attr_General] = cd.SupportedSystems.General
	capabilityDetailsMap[Attr_SupportedSystems] = []map[string]interface{}{supportedSystemsMap}
	return capabilityDetailsMap, nil
}

func disasterRecoveryToMap(dr *models.DisasterRecovery) map[string]interface{} {
	disasterRecoveryMap := make(map[string]interface{})
	asynchronousReplicationMap := replicationServiceToMap(dr.AsynchronousReplication)
	disasterRecoveryMap[Attr_AsynchronousReplication] = []map[string]interface{}{asynchronousReplicationMap}
	if dr.SynchronousReplication != nil {
		synchronousReplicationMap := replicationServiceToMap(dr.SynchronousReplication)
		disasterRecoveryMap[Attr_SynchronousReplication] = []map[string]interface{}{synchronousReplicationMap}
	}

	return disasterRecoveryMap
}

func replicationServiceToMap(rs *models.ReplicationService) map[string]interface{} {
	replicationMap := make(map[string]interface{})
	replicationMap[Attr_Enabled] = rs.Enabled
	targetLocations := []map[string]interface{}{}
	for _, targetLocationsItem := range rs.TargetLocations {

		targetLocationsItemMap := make(map[string]interface{})
		if targetLocationsItem.Region != "" {
			targetLocationsItemMap[Attr_Region] = targetLocationsItem.Region
		}
		if targetLocationsItem.Status != "" {
			targetLocationsItemMap[Attr_Status] = targetLocationsItem.Status
		}
		targetLocations = append(targetLocations, targetLocationsItemMap)
	}
	replicationMap[Attr_TargetLocations] = targetLocations
	return replicationMap
}
