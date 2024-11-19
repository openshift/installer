// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	models "github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMPISharedProcessorPool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPISharedProcessorPoolCreate,
		ReadContext:   resourceIBMPISharedProcessorPoolRead,
		UpdateContext: resourceIBMPISharedProcessorPoolUpdate,
		DeleteContext: resourceIBMPISharedProcessorPoolDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			// Required Arguments
			Arg_SharedProcessorPoolName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the shared processor pool",
			},

			Arg_SharedProcessorPoolHostGroup: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Host group of the shared processor pool",
			},

			Arg_SharedProcessorPoolReservedCores: {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The amount of reserved cores for the shared processor pool",
			},

			Arg_CloudInstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "PI cloud instance ID",
			},
			Arg_HostID: {
				Description: "The host id of a host in a host group (only available for dedicated hosts)",
				Optional:    true,
				Type:        schema.TypeString,
			},

			// Optional Arguments
			Arg_SharedProcessorPoolPlacementGroupID: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Placement group the shared processor pool is created in",
			},
			Arg_UserTags: {
				Description: "The user tags attached to this resource.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Set:         schema.HashString,
				Type:        schema.TypeSet,
			},

			// Attributes
			Attr_CRN: {
				Computed:    true,
				Description: "The CRN of this resource.",
				Type:        schema.TypeString,
			},
			Attr_SharedProcessorPoolID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Shared processor pool ID",
			},

			Attr_SharedProcessorPoolAvailableCores: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Shared processor pool available cores",
			},

			Attr_SharedProcessorPoolAllocatedCores: {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Shared processor pool allocated cores",
			},

			Attr_SharedProcessorPoolHostID: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The host ID where the shared processor pool resides",
			},

			Attr_Status: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the shared processor pool",
			},

			Attr_StatusDetail: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status details of the shared processor pool",
			},

			Attr_SharedProcessorPoolPlacementGroups: {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "SPP placement groups the shared processor pool are in",
			},

			Attr_SharedProcessorPoolInstances: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of server instances deployed in the shared processor pool",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_SharedProcessorPoolInstanceCpus: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The amount of cpus for the server instance",
						},
						Attr_SharedProcessorPoolInstanceUncapped: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Identifies if uncapped or not",
						},
						Attr_SharedProcessorPoolInstanceAvailabilityZone: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability zone for the server instances",
						},
						Attr_SharedProcessorPoolInstanceId: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The server instance ID",
						},
						Attr_SharedProcessorPoolInstanceMemory: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The amount of memory for the server instance",
						},
						Attr_SharedProcessorPoolInstanceName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The server instance name",
						},
						Attr_Status: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the server",
						},
						Attr_SharedProcessorPoolInstanceVcpus: {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The amout of vcpus for the server instance",
						},
					},
				},
			},
		},
	}
}

func resourceIBMPISharedProcessorPoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	name := d.Get(Arg_SharedProcessorPoolName).(string)
	hostGroup := d.Get(Arg_SharedProcessorPoolHostGroup).(string)
	hostID := d.Get(Arg_HostID).(string)
	reservedCores := d.Get(Arg_SharedProcessorPoolReservedCores).(int)
	cores := int64(reservedCores)
	client := st.NewIBMPISharedProcessorPoolClient(ctx, sess, cloudInstanceID)
	body := &models.SharedProcessorPoolCreate{
		HostGroup:     &hostGroup,
		HostID:        hostID,
		Name:          &name,
		ReservedCores: &cores,
	}

	if pg, ok := d.GetOk(Arg_SharedProcessorPoolPlacementGroupID); ok {
		body.PlacementGroupID = pg.(string)
	}
	if tags, ok := d.GetOk(Arg_UserTags); ok {
		body.UserTags = flex.FlattenSet(tags.(*schema.Set))
	}

	spp, err := client.Create(body)
	if err != nil || spp == nil {
		return diag.Errorf("error creating the shared processor pool: %v", err)
	}

	var sharedProcessorPoolReadyStatus string
	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, *spp.ID))
	_, err = isWaitForPISharedProcessorPoolAvailable(ctx, d, client, *spp.ID, sharedProcessorPoolReadyStatus)
	if err != nil {
		return diag.FromErr(err)
	}

	if _, ok := d.GetOk(Arg_UserTags); ok {
		if spp.Crn != "" {
			oldList, newList := d.GetChange(Arg_UserTags)
			err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, string(spp.Crn), "", UserTagType)
			if err != nil {
				log.Printf("Error on update of pi shared processor pool (%s) pi_user_tags during creation: %s", *spp.ID, err)
			}
		}
	}

	return resourceIBMPISharedProcessorPoolRead(ctx, d, meta)

}

func isWaitForPISharedProcessorPoolAvailable(ctx context.Context, d *schema.ResourceData, client *st.IBMPISharedProcessorPoolClient, id string, sharedProcessorPoolReadyStatus string) (interface{}, error) {
	log.Printf("Waiting for PISharedProcessorPool (%s) to be active ", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"configuring"},
		Target:     []string{"active", "failed", ""},
		Refresh:    isPISharedProcessorPoolRefreshFunc(client, id, sharedProcessorPoolReadyStatus),
		Delay:      20 * time.Second,
		MinTimeout: Timeout_Active,
		Timeout:    d.Timeout(schema.TimeoutCreate),
	}

	return stateConf.WaitForStateContext(ctx)
}

func isPISharedProcessorPoolRefreshFunc(client *st.IBMPISharedProcessorPoolClient, id, sharedProcessorPoolReadyStatus string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		pool, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}
		// Check for `sharedProcessorPoolReadyStatus` status
		if pool.SharedProcessorPool.Status == "active" {
			return pool, "active", nil
		}
		if pool.SharedProcessorPool.Status == "failed" {
			err = fmt.Errorf("failed to create the shared processor pool")
			return pool, pool.SharedProcessorPool.Status, err
		}

		return pool, "configuring", nil
	}
}

func resourceIBMPISharedProcessorPoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := parts[0]
	client := st.NewIBMPISharedProcessorPoolClient(ctx, sess, cloudInstanceID)

	response, err := client.Get(parts[1])
	if err != nil || response == nil {
		return diag.Errorf("error reading the shared processor pool: %v", err)
	}

	d.Set(Arg_CloudInstanceID, cloudInstanceID)
	if response.SharedProcessorPool.Crn != "" {
		d.Set(Attr_CRN, response.SharedProcessorPool.Crn)
		tags, err := flex.GetGlobalTagsUsingCRN(meta, string(response.SharedProcessorPool.Crn), "", UserTagType)
		if err != nil {
			log.Printf("Error on get of pi shared processor pool (%s) pi_user_tags: %s", *response.SharedProcessorPool.ID, err)
		}
		d.Set(Arg_UserTags, tags)
	}
	d.Set(Arg_SharedProcessorPoolHostGroup, response.SharedProcessorPool.HostGroup)

	if response.SharedProcessorPool.Name != nil {
		d.Set(Arg_SharedProcessorPoolName, response.SharedProcessorPool.Name)
	}
	if response.SharedProcessorPool.ID != nil {
		d.Set(Attr_SharedProcessorPoolID, response.SharedProcessorPool.ID)
	}
	if response.SharedProcessorPool.ReservedCores != nil {
		d.Set(Arg_SharedProcessorPoolReservedCores, response.SharedProcessorPool.ReservedCores)
	}
	if response.SharedProcessorPool.AllocatedCores != nil {
		d.Set(Attr_SharedProcessorPoolAllocatedCores, response.SharedProcessorPool.AllocatedCores)
	}
	if response.SharedProcessorPool.AvailableCores != nil {
		d.Set(Attr_SharedProcessorPoolAvailableCores, response.SharedProcessorPool.AvailableCores)
	}
	if response.SharedProcessorPool.AvailableCores != nil {
		d.Set(Attr_SharedProcessorPoolAvailableCores, response.SharedProcessorPool.AvailableCores)
	}
	if response.SharedProcessorPool.SharedProcessorPoolPlacementGroups != nil {
		pgIDs := make([]string, len(response.SharedProcessorPool.SharedProcessorPoolPlacementGroups))
		for i, pg := range response.SharedProcessorPool.SharedProcessorPoolPlacementGroups {
			pgIDs[i] = *pg.ID
		}
		d.Set(Attr_SharedProcessorPoolPlacementGroups, pgIDs)
	}
	d.Set(Attr_SharedProcessorPoolHostID, response.SharedProcessorPool.HostID)
	d.Set(Attr_Status, response.SharedProcessorPool.Status)
	d.Set(Attr_SharedProcessorPoolStatusDetail, response.SharedProcessorPool.StatusDetail)

	serversMap := []map[string]interface{}{}
	if response.Servers != nil {
		for _, s := range response.Servers {
			if s != nil {
				v := map[string]interface{}{
					Attr_SharedProcessorPoolInstanceCpus:             s.Cpus,
					Attr_SharedProcessorPoolInstanceUncapped:         s.Uncapped,
					Attr_SharedProcessorPoolInstanceAvailabilityZone: s.AvailabilityZone,
					Attr_SharedProcessorPoolInstanceId:               s.ID,
					Attr_SharedProcessorPoolInstanceMemory:           s.Memory,
					Attr_SharedProcessorPoolInstanceName:             s.Name,
					Attr_Status:                                      s.Status,
					Attr_SharedProcessorPoolInstanceVcpus:            s.Vcpus,
				}
				serversMap = append(serversMap, v)
			}
		}
	}
	d.Set(Attr_SharedProcessorPoolInstances, serversMap)

	return nil
}

func resourceIBMPISharedProcessorPoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, sppID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := st.NewIBMPISharedProcessorPoolClient(ctx, sess, cloudInstanceID)
	body := &models.SharedProcessorPoolUpdate{}

	if d.HasChange(Arg_SharedProcessorPoolName) {
		name := d.Get(Arg_SharedProcessorPoolName).(string)
		body.Name = name
	}
	if d.HasChange(Arg_SharedProcessorPoolReservedCores) {
		reservedCores := int64(d.Get(Arg_SharedProcessorPoolReservedCores).(int))
		body.ReservedCores = &reservedCores
	}

	_, err = client.Update(sppID, body)
	if err != nil {
		return diag.Errorf("error updating the shared processor pool: %v", err)
	}

	if d.HasChange(Attr_SharedProcessorPoolPlacementGroups) {

		pgClient := st.NewIBMPISPPPlacementGroupClient(ctx, sess, cloudInstanceID)

		oldRaw, newRaw := d.GetChange(Attr_SharedProcessorPoolPlacementGroups)
		old := oldRaw.([]interface{})
		new := newRaw.([]interface{})

		var oldPGs []string
		for _, o := range old {
			oldPGs = append(oldPGs, o.(string))
		}
		var newPGs []string
		for _, n := range new {
			newPGs = append(newPGs, n.(string))
		}
		// find removed pgs and remove them
		pgsToRemove := getDifferences(oldPGs, newPGs)

		for _, pgToRemove := range pgsToRemove {
			if len(strings.TrimSpace(pgToRemove)) > 0 {
				placementGroupID := pgToRemove
				//remove spp from old placement group
				_, err := pgClient.DeleteMember(placementGroupID, sppID)
				if err != nil {
					// ignore delete member error where the spp is already not in the PG
					if !strings.Contains(err.Error(), "is not part of spp placement group") {
						return diag.FromErr(err)
					}
				}
			}
		}

		// find added pgs and then add them
		pgsToAdd := getDifferences(newPGs, oldPGs)

		for _, pgToAdd := range pgsToAdd {
			if len(strings.TrimSpace(pgToAdd)) > 0 {
				placementGroupID := pgToAdd
				// add spp to a new placement group
				_, err := pgClient.AddMember(placementGroupID, sppID)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	if d.HasChange(Arg_UserTags) {
		if crn, ok := d.GetOk(Attr_CRN); ok {
			oldList, newList := d.GetChange(Arg_UserTags)
			err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, crn.(string), "", UserTagType)
			if err != nil {
				log.Printf("Error on update of pi shared processor pool (%s) pi_user_tags: %s", sppID, err)
			}
		}
	}

	return resourceIBMPISharedProcessorPoolRead(ctx, d, meta)
}

// returns the elements in string array a that are not in array z
func getDifferences(a, z []string) []string {
	mb := make(map[string]struct{}, len(z))
	for _, x := range z {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func resourceIBMPISharedProcessorPoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID := parts[0]
	client := st.NewIBMPISharedProcessorPoolClient(ctx, sess, cloudInstanceID)
	err = client.Delete(parts[1])

	if err != nil {
		return diag.Errorf("error deleting the shared processor pool: %v", err)
	}
	d.SetId("")
	return nil
}
