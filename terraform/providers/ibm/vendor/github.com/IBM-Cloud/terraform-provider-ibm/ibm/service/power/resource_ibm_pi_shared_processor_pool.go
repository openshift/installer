// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
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
		CustomizeDiff: customdiff.Sequence(
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return flex.ResourcePowerUserTagsCustomizeDiff(diff)
			},
		),

		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description: "The GUID of the service instance associated with an account.",
				ForceNew:    true,
				Required:    true,
				Type:        schema.TypeString,
			},
			Arg_HostID: {
				Description: "The host id of a host in a host group (only available for dedicated hosts).",
				Optional:    true,
				Type:        schema.TypeString,
			},
			Arg_SharedProcessorPoolHostGroup: {
				Description: "Host group of the shared processor pool. Valid values are 's922', 'e980' and 's1022'.",
				ForceNew:    true,
				Required:    true,
				Type:        schema.TypeString,
			},
			Arg_SharedProcessorPoolName: {
				Description: "The name of the shared processor pool.",
				Required:    true,
				Type:        schema.TypeString,
			},
			Arg_SharedProcessorPoolPlacementGroupID: {
				ConflictsWith: []string{Arg_SharedProcessorPoolPlacementGroups},
				Deprecated:    "This field is deprecated, use pi_shared_processor_pool_placement_groups instead",
				Description:   "The ID of the placement group the shared processor pool is created in.",
				Optional:      true,
				Type:          schema.TypeString,
			},
			Arg_SharedProcessorPoolPlacementGroups: {
				ConflictsWith: []string{Arg_SharedProcessorPoolPlacementGroupID, Attr_SharedProcessorPoolPlacementGroups},
				Description:   "The list of shared processor pool placement groups that the shared processor pool is in.",
				Elem:          &schema.Schema{Type: schema.TypeString},
				Optional:      true,
				Type:          schema.TypeList,
			},
			Arg_SharedProcessorPoolReservedCores: {
				Description: "The amount of reserved cores for the shared processor pool.",
				Required:    true,
				Type:        schema.TypeInt,
			},
			Arg_UserTags: {
				Computed:    true,
				Description: "The user tags attached to this resource.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Set:         schema.HashString,
				Type:        schema.TypeSet,
			},

			// Attributes
			Attr_AllocatedCores: {
				Computed:    true,
				Description: "The allocated cores in the shared processor pool.",
				Type:        schema.TypeFloat,
			},
			Attr_AvailableCores: {
				Computed:    true,
				Description: "The available cores in the shared processor pool.",
				Type:        schema.TypeInt,
			},
			Attr_CRN: {
				Computed:    true,
				Description: "The CRN of this resource.",
				Type:        schema.TypeString,
			},
			Attr_DedicatedHostID: {
				Computed:    true,
				Description: "The dedicated host ID where the shared processor pool resides.",
				Type:        schema.TypeString,
			},
			Attr_HostID: {
				Computed:    true,
				Description: "The host ID where the shared processor pool resides.",
				Type:        schema.TypeInt,
			},
			Attr_Instances: {
				Computed:    true,
				Description: "The list of server instances that are deployed in the shared processor pool.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_AvailabilityZone: {
							Computed:    true,
							Description: "Availability zone for the server instances.",
							Type:        schema.TypeString,
						},
						Attr_CPUs: {
							Computed:    true,
							Description: "The amount of cpus for the server instance.",
							Type:        schema.TypeInt,
						},
						Attr_ID: {
							Computed:    true,
							Description: "The server instance ID.",
							Type:        schema.TypeString,
						},
						Attr_Memory: {
							Computed:    true,
							Description: "The amount of memory for the server instance.",
							Type:        schema.TypeInt,
						},
						Attr_Name: {
							Computed:    true,
							Description: "The server instance name.",
							Type:        schema.TypeString,
						},
						Attr_Status: {
							Computed:    true,
							Description: "Status of the instance.",
							Type:        schema.TypeString,
						},
						Attr_Uncapped: {
							Computed:    true,
							Description: "Identifies if uncapped or not.",
							Type:        schema.TypeBool,
						},
						Attr_VCPUs: {
							Computed:    true,
							Description: "The amout of vcpus for the server instance.",
							Type:        schema.TypeFloat,
						},
					},
				},
				Type: schema.TypeList,
			},
			Attr_SharedProcessorPoolID: {
				Computed:    true,
				Description: "The shared processor pool's unique ID.",
				Type:        schema.TypeString,
			},
			Attr_SharedProcessorPoolPlacementGroups: {
				ConflictsWith: []string{Arg_SharedProcessorPoolPlacementGroups},
				Deprecated:    "This field is deprecated, use pi_shared_processor_pool_placement_groups instead",
				Description:   "The list of shared processor pool placement groups that the shared processor pool is in.",
				Elem:          &schema.Schema{Type: schema.TypeString},
				Optional:      true,
				Type:          schema.TypeList,
			},
			Attr_Status: {
				Computed:    true,
				Description: "The status of the shared processor pool.",
				Type:        schema.TypeString,
			},
			Attr_StatusDetail: {
				Computed:    true,
				Description: "The status details of the shared processor pool.",
				Type:        schema.TypeString,
			},
		},
	}
}

func resourceIBMPISharedProcessorPoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	hostGroup := d.Get(Arg_SharedProcessorPoolHostGroup).(string)
	hostID := d.Get(Arg_HostID).(string)
	name := d.Get(Arg_SharedProcessorPoolName).(string)
	reservedCores := d.Get(Arg_SharedProcessorPoolReservedCores).(int)
	cores := int64(reservedCores)
	client := instance.NewIBMPISharedProcessorPoolClient(ctx, sess, cloudInstanceID)
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

	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, *spp.ID))
	_, err = isWaitForPISharedProcessorPoolAvailable(ctx, d, client, *spp.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	diagErr := detectSPPPlacementGroupChange(ctx, sess, cloudInstanceID, d, *spp.ID)
	if diagErr != nil {
		return diagErr
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

func isWaitForPISharedProcessorPoolAvailable(ctx context.Context, d *schema.ResourceData, client *instance.IBMPISharedProcessorPoolClient, id string) (interface{}, error) {
	log.Printf("Waiting for PISharedProcessorPool (%s) to be active ", id)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Configuring},
		Target:     []string{State_Active, State_Failed, ""},
		Refresh:    isPISharedProcessorPoolRefreshFunc(client, id),
		Delay:      20 * time.Second,
		MinTimeout: Timeout_Active,
		Timeout:    d.Timeout(schema.TimeoutCreate),
	}

	return stateConf.WaitForStateContext(ctx)
}

func isPISharedProcessorPoolRefreshFunc(client *instance.IBMPISharedProcessorPoolClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		pool, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}
		if pool.SharedProcessorPool.Status == State_Active {
			return pool, State_Active, nil
		}
		if pool.SharedProcessorPool.Status == State_Failed {
			err = fmt.Errorf("failed to create the shared processor pool")
			return pool, pool.SharedProcessorPool.Status, err
		}
		return pool, State_Configuring, nil
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
	client := instance.NewIBMPISharedProcessorPoolClient(ctx, sess, cloudInstanceID)

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
	if response.SharedProcessorPool.AllocatedCores != nil {
		d.Set(Attr_AllocatedCores, response.SharedProcessorPool.AllocatedCores)
	}
	if response.SharedProcessorPool.AvailableCores != nil {
		d.Set(Attr_AvailableCores, response.SharedProcessorPool.AvailableCores)
	}
	if response.SharedProcessorPool.ReservedCores != nil {
		d.Set(Arg_SharedProcessorPoolReservedCores, response.SharedProcessorPool.ReservedCores)
	}
	if response.SharedProcessorPool.SharedProcessorPoolPlacementGroups != nil {
		pgIDs := make([]string, len(response.SharedProcessorPool.SharedProcessorPoolPlacementGroups))
		for i, pg := range response.SharedProcessorPool.SharedProcessorPoolPlacementGroups {
			pgIDs[i] = *pg.ID
		}
		if _, ok := d.GetOk(Attr_SharedProcessorPoolPlacementGroups); ok {
			d.Set(Attr_SharedProcessorPoolPlacementGroups, pgIDs)
		} else {
			d.Set(Arg_SharedProcessorPoolPlacementGroups, pgIDs)
		}
	}
	d.Set(Attr_DedicatedHostID, response.SharedProcessorPool.DedicatedHostID)
	d.Set(Attr_HostID, response.SharedProcessorPool.HostID)
	d.Set(Attr_Status, response.SharedProcessorPool.Status)
	d.Set(Attr_StatusDetail, response.SharedProcessorPool.StatusDetail)

	serversMap := []map[string]interface{}{}
	if response.Servers != nil {
		for _, s := range response.Servers {
			if s != nil {
				v := map[string]interface{}{
					Attr_AvailabilityZone: s.AvailabilityZone,
					Attr_CPUs:             s.Cpus,
					Attr_ID:               s.ID,
					Attr_Memory:           s.Memory,
					Attr_Name:             s.Name,
					Attr_Status:           s.Status,
					Attr_Uncapped:         s.Uncapped,
					Attr_VCPUs:            s.Vcpus,
				}
				serversMap = append(serversMap, v)
			}
		}
	}
	d.Set(Attr_Instances, serversMap)

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

	client := instance.NewIBMPISharedProcessorPoolClient(ctx, sess, cloudInstanceID)
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

	diagErr := detectSPPPlacementGroupChange(ctx, sess, cloudInstanceID, d, sppID)
	if diagErr != nil {
		return diagErr
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

func detectSPPPlacementGroupChange(ctx context.Context, sess *ibmpisession.IBMPISession, cloudInstanceID string, d *schema.ResourceData, sppID string) diag.Diagnostics {
	if d.HasChanges(Arg_SharedProcessorPoolPlacementGroups, Attr_SharedProcessorPoolPlacementGroups) {

		pgClient := instance.NewIBMPISPPPlacementGroupClient(ctx, sess, cloudInstanceID)

		var oldRaw, newRaw interface{}
		if d.HasChange(Arg_SharedProcessorPoolPlacementGroups) {
			oldRaw, newRaw = d.GetChange(Arg_SharedProcessorPoolPlacementGroups)
		} else {
			oldRaw, newRaw = d.GetChange(Attr_SharedProcessorPoolPlacementGroups)
		}
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
	return nil
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
	client := instance.NewIBMPISharedProcessorPoolClient(ctx, sess, cloudInstanceID)
	err = client.Delete(parts[1])

	if err != nil {
		return diag.Errorf("error deleting the shared processor pool: %v", err)
	}
	d.SetId("")
	return nil
}
