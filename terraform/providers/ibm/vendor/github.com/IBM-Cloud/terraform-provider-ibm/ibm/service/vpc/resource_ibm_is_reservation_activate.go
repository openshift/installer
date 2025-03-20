// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"log"
	"runtime/debug"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMISReservationActivate() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISReservationActivateCreate,
		Read:     resourceIBMISReservationActivateRead,
		Delete:   resourceIBMISReservationActivateDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			isReservation: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The unique identifier for this reservation.",
			},
			isReservationAffinityPolicy: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The affinity policy to use for this reservation",
			},
			isReservationCapacity: &schema.Schema{
				Type:        schema.TypeList,
				ForceNew:    true,
				Computed:    true,
				Description: "The capacity reservation configuration to use",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isReservationCapacityTotal: &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total amount to use for this capacity reservation.",
						},
						isReservationCapacityAllocated: &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The amount allocated to this capacity reservation.",
						},
						isReservationCapacityAvailable: &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The amount of this capacity reservation available for new attachments.",
						},
						isReservationCapacityUsed: &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The amount of this capacity reservation used by existing attachments.",
						},
						isReservationCapacityStatus: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the capacity reservation.",
						},
					},
				},
			},
			isReservationCommittedUse: &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The committed use configuration to use for this reservation",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isReservationComittedUseExpirationPolicy: &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The maximum number of days to keep each backup after creation.",
						},
						isReservationComittedUseTerm: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The maximum number of recent backups to keep. If unspecified, there will be no maximum.",
						},
						isReservationCommittedUseExpirationAt: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The expiration date and time for this committed use reservation.",
						},
					},
				},
			},
			isReservationCreatedAt: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the reservation was created.",
			},
			isReservationCrn: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this reservation.",
			},
			isReservationHref: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this reservation.",
			},
			isReservationId: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this reservation.",
			},
			isReservationLifecycleState: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of this reservation.",
			},
			isReservationName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reservation name",
			},
			isReservationProfile: &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The profile used for this reservation.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isReservationProfileName: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique name for this virtual server instance profile.",
						},
						isReservationProfileResourceType: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type of the profile.",
						},
						isReservationProfileHref: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this virtual server instance profile.",
						},
					},
				},
			},
			isReservationResourceGroup: &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "The committed use configuration to use for this reservation",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isReservationResourceGroupHref: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this resource group.",
						},
						isReservationResourceGroupId: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this resource group",
						},
						isReservationResourceGroupName: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this resource group.",
						},
					},
				},
			},
			isReservationResourceType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			isReservationStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the reservation.",
			},
			isReservationStatusReasons: &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The committed use configuration to use for this reservation",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isReservationStatusReasonCode: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: " snake case string succinctly identifying the status reason.",
						},
						isReservationStatusReasonMessage: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the status reason",
						},
						isReservationStatusReasonMoreInfo: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "Link to documentation about this status reason.",
						},
					},
				},
			},
			isReservationZone: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The globally unique name for this zone.",
			},
		},
	}
}
func resourceIBMISReservationActivateCreate(d *schema.ResourceData, meta interface{}) error {

	id := d.Get(isReservation).(string)
	activateReservationOptions := &vpcv1.ActivateReservationOptions{
		ID: core.StringPtr(id),
	}

	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	response, err := sess.ActivateReservation(activateReservationOptions)
	if err != nil {
		log.Printf("[DEBUG] Reservation activation err %s\n%s", err, response)
		return fmt.Errorf("[ERROR] Error while activating Reservation %s\n%v", err, response)
	}
	log.Printf("[INFO] Reservation activated: %s", id)
	d.SetId(id)

	return resourceIBMISReservationActivateRead(d, meta)
}

func resourceIBMISReservationActivateRead(d *schema.ResourceData, meta interface{}) error {

	id := d.Id()

	sess, err := vpcClient(meta)
	defer func() {

		log.Println("stacktrace from panic: \n", err, string(debug.Stack()))

	}()
	if err != nil {
		return err
	}
	getReservationOptions := &vpcv1.GetReservationOptions{
		ID: &id,
	}
	reservation, response, err := sess.GetReservation(getReservationOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error Getting Reservation (%s): %s\n%s", id, err, response)
	}

	if reservation.AffinityPolicy != nil {
		if err = d.Set(isReservationAffinityPolicy, reservation.AffinityPolicy); err != nil {
			log.Printf("[ERROR] Error setting %s: %s", isReservationAffinityPolicy, err)
			return fmt.Errorf("[ERROR] Error setting %s: %s", isReservationAffinityPolicy, err)
		}
	}

	if reservation.Capacity != nil {
		capacityMap := []map[string]interface{}{}
		finalList := map[string]interface{}{}

		if reservation.Capacity.Allocated != nil {
			finalList[isReservationCapacityAllocated] = flex.IntValue(reservation.Capacity.Allocated)
		}
		if reservation.Capacity.Available != nil {
			finalList[isReservationCapacityAvailable] = flex.IntValue(reservation.Capacity.Available)
		}
		if reservation.Capacity.Total != nil {
			finalList[isReservationCapacityTotal] = flex.IntValue(reservation.Capacity.Total)
		}
		if reservation.Capacity.Used != nil {
			finalList[isReservationCapacityUsed] = flex.IntValue(reservation.Capacity.Used)
		}
		if reservation.Capacity.Status != nil {
			finalList[isReservationCapacityStatus] = reservation.Capacity.Status
		}
		capacityMap = append(capacityMap, finalList)
		d.Set(isReservationCapacity, capacityMap)
	}

	if reservation.CommittedUse != nil {
		committedUseMap := []map[string]interface{}{}
		finalList := map[string]interface{}{}

		if reservation.CommittedUse.ExpirationAt != nil {
			finalList[isReservationCommittedUseExpirationAt] = flex.DateTimeToString(reservation.CommittedUse.ExpirationAt)
		}
		if reservation.CommittedUse.ExpirationPolicy != nil {
			finalList[isReservationComittedUseExpirationPolicy] = *reservation.CommittedUse.ExpirationPolicy
		}
		if reservation.CommittedUse.Term != nil {
			finalList[isReservationComittedUseTerm] = *reservation.CommittedUse.Term
		}
		committedUseMap = append(committedUseMap, finalList)
		d.Set(isReservationCommittedUse, committedUseMap)
	}

	if reservation.CreatedAt != nil {
		if err = d.Set(isReservationCreatedAt, flex.DateTimeToString(reservation.CreatedAt)); err != nil {
			log.Printf("[ERROR] Error setting %s: %s", isReservationCreatedAt, err)
			return fmt.Errorf("[ERROR] Error setting %s: %s", isReservationCreatedAt, err)
		}
	}

	if reservation.CRN != nil {
		if err = d.Set(isReservationCrn, reservation.CRN); err != nil {
			log.Printf("[ERROR] Error setting %s: %s", isReservationCrn, err)
			return fmt.Errorf("[ERROR] Error setting %s: %s", isReservationCrn, err)
		}
	}

	if reservation.Href != nil {
		if err = d.Set(isReservationHref, reservation.Href); err != nil {
			log.Printf("[ERROR] Error setting %s: %s", isReservationHref, err)
			return fmt.Errorf("[ERROR] Error setting %s: %s", isReservationHref, err)
		}
	}

	if reservation.LifecycleState != nil {
		if err = d.Set(isReservationLifecycleState, reservation.LifecycleState); err != nil {
			log.Printf("[ERROR] Error setting %s: %s", isReservationLifecycleState, err)
			return fmt.Errorf("[ERROR] Error setting %s: %s", isReservationLifecycleState, err)
		}
	}

	if reservation.Name != nil {
		if err = d.Set(isReservationName, reservation.Name); err != nil {
			log.Printf("[ERROR] Error setting %s: %s", isReservationName, err)
			return fmt.Errorf("[ERROR] Error setting %s: %s", isReservationName, err)
		}
	}

	if reservation.Profile != nil {
		profileMap := []map[string]interface{}{}
		finalList := map[string]interface{}{}

		profileItem := reservation.Profile.(*vpcv1.ReservationProfile)

		if profileItem.Href != nil {
			finalList[isReservationProfileHref] = profileItem.Href
		}
		if profileItem.Name != nil {
			finalList[isReservationProfileName] = profileItem.Name
		}
		if profileItem.ResourceType != nil {
			finalList[isReservationProfileResourceType] = profileItem.ResourceType
		}
		profileMap = append(profileMap, finalList)
		d.Set(isReservationProfile, profileMap)
	}

	if reservation.ResourceGroup != nil {
		rgMap := []map[string]interface{}{}
		finalList := map[string]interface{}{}

		if reservation.ResourceGroup.Href != nil {
			finalList[isReservationResourceGroupHref] = reservation.ResourceGroup.Href
		}
		if reservation.ResourceGroup.ID != nil {
			finalList[isReservationResourceGroupId] = reservation.ResourceGroup.ID
		}
		if reservation.ResourceGroup.Name != nil {
			finalList[isReservationResourceGroupName] = reservation.ResourceGroup.Name
		}
		rgMap = append(rgMap, finalList)
		d.Set(isReservationResourceGroup, rgMap)
	}

	if reservation.ResourceType != nil {
		if err = d.Set(isReservationResourceType, reservation.ResourceType); err != nil {
			log.Printf("[ERROR] Error setting %s: %s", isReservationResourceType, err)
			return fmt.Errorf("[ERROR] Error setting %s: %s", isReservationResourceType, err)
		}
	}

	if reservation.Status != nil {
		if err = d.Set(isReservationStatus, reservation.Status); err != nil {
			log.Printf("[ERROR] Error setting %s: %s", isReservationStatus, err)
			return fmt.Errorf("[ERROR] Error setting %s: %s", isReservationStatus, err)
		}
	}

	if reservation.StatusReasons != nil {
		srLen := len(reservation.StatusReasons)
		srList := []vpcv1.ReservationStatusReason{}

		for i := 0; i < srLen; i++ {
			srList = append(srList, reservation.StatusReasons[i])
		}
		d.Set(isReservationStatusReasons, srList)
	}

	if reservation.Zone != nil && reservation.Zone.Name != nil {
		if err = d.Set(isReservationZone, reservation.Zone.Name); err != nil {
			log.Printf("[ERROR] Error setting %s: %s", isReservationZone, err)
			return fmt.Errorf("[ERROR] Error setting %s: %s", isReservationZone, err)
		}
	}
	return nil
}

func resourceIBMISReservationActivateDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}
