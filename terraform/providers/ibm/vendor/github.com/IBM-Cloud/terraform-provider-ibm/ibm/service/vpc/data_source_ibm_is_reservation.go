// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsReservation() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsReservationRead,

		Schema: map[string]*schema.Schema{

			"identifier": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"name", "identifier"},
				Description:  "The reservation identifier.",
			},

			"name": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ExactlyOneOf: []string{"name", "identifier"},
				Description:  "The unique user-defined name for this reservation.",
			},

			"affinity_policy": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The affinity policy to use for this reservation.",
			},
			"capacity": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: "The capacity configuration for this reservation. If absent, this reservation has no assigned capacity.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allocated": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The amount allocated to this capacity reservation.",
						},
						"available": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The amount of this capacity reservation available for new attachments.",
						},
						"status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the capacity reservation.",
						},
						"total": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total amount of this capacity reservation.",
						},
						"used": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The amount of this capacity reservation used by existing attachments.",
						},
					},
				},
			},
			"committed_use": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: "The committed use configuration for this reservation. If absent, this reservation has no commitment for use.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expiration_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The expiration date and time for this committed use reservation.",
						},
						"expiration_policy": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The policy to apply when the committed use term expires.",
						},
						"term": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The term for this committed use reservation.",
						},
					},
				},
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the reservation was created.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this reservation.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this reservation.",
			},
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this reservation.",
			},
			"lifecycle_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the reservation.",
			},
			"profile": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The profile for this reservation.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this virtual server instance profile.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique name for this virtual server instance profile.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"resource_group": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The resource group for this reservation.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this resource group.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this resource group.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this resource group.",
						},
					},
				},
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of resource referenced.",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the reservation.",
			},
			"status_reasons": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reasons for the current status (if any).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A snake case string succinctly identifying the status reason.",
						},
						"message": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the status reason.",
						},
						"more_info": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "Link to documentation about this status reason.",
						},
					},
				},
			},
			"zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The zone for this reservation.",
			},
		},
	}
}

func dataSourceIBMIsReservationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	var reservation *vpcv1.Reservation

	if v, ok := d.GetOk("identifier"); ok {
		id := v.(string)
		getReservationOptions := &vpcv1.GetReservationOptions{}
		getReservationOptions.SetID(id)
		reservationInfo, response, err := sess.GetReservationWithContext(context, getReservationOptions)
		if err != nil {
			log.Printf("[DEBUG] GetReservationWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("[ERROR] GetReservationWithContext failed %s\n%s", err, response))
		}
		reservation = reservationInfo

	}
	if v, ok := d.GetOk("name"); ok {

		name := v.(string)
		start := ""
		allrecs := []vpcv1.Reservation{}
		for {
			listReservationsOptions := &vpcv1.ListReservationsOptions{}
			if start != "" {
				listReservationsOptions.Start = &start
			}
			reservationCollection, response, err := sess.ListReservationsWithContext(context, listReservationsOptions)
			if err != nil {
				log.Printf("[DEBUG] ListReservationsWithContext failed %s\n%s", err, response)
				return diag.FromErr(fmt.Errorf("[ERROR] ListReservationsWithContext failed %s\n%s", err, response))
			}
			if reservationCollection != nil && *reservationCollection.TotalCount == int64(0) {
				break
			}
			start = flex.GetNext(reservationCollection.Next)
			allrecs = append(allrecs, reservationCollection.Reservations...)
			if start == "" {
				break
			}
		}
		for _, reservationInfo := range allrecs {
			if *reservationInfo.Name == name {
				reservation = &reservationInfo
				break
			}
		}
		if reservation == nil {
			return diag.FromErr(fmt.Errorf("[ERROR] No reservation found with name (%s)", name))
		}
	}

	d.SetId(*reservation.ID)
	if err = d.Set("affinity_policy", reservation.AffinityPolicy); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting affinity_policy: %s", err))
	}
	if reservation.Capacity != nil {
		capacityList := []map[string]interface{}{}
		capacityMap := dataSourceReservationCapacityToMap(*reservation.Capacity)
		capacityList = append(capacityList, capacityMap)
		if err = d.Set("capacity", capacityList); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting capacity: %s", err))
		}
	}
	if reservation.CommittedUse != nil {
		committedUseList := []map[string]interface{}{}
		committedUseMap := dataSourceReservationCommittedUseToMap(*reservation.CommittedUse)
		committedUseList = append(committedUseList, committedUseMap)
		if err = d.Set("committed_use", committedUseList); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting committed_use: %s", err))
		}
	}
	if err = d.Set("created_at", reservation.CreatedAt.String()); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
	}
	if err = d.Set("crn", reservation.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting crn: %s", err))
	}
	if err = d.Set("href", reservation.Href); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting href: %s", err))
	}
	if err = d.Set("lifecycle_state", reservation.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting lifecycle_state: %s", err))
	}
	if err = d.Set("name", reservation.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}
	if reservation.Profile != nil {
		profileList := []map[string]interface{}{}
		profile := reservation.Profile
		profileMap := dataSourceReservationProfileToMap(profile)
		profileList = append(profileList, profileMap)
		if err = d.Set("profile", profileList); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting profile: %s", err))
		}
	}
	if reservation.ResourceGroup != nil {
		resourceGroupList := []map[string]interface{}{}
		resourceGroupMap := dataSourceReservationResourceGroupToMap(*reservation.ResourceGroup)
		resourceGroupList = append(resourceGroupList, resourceGroupMap)
		if err = d.Set("resource_group", resourceGroupList); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_group: %s", err))
		}
	}
	if err = d.Set("resource_type", reservation.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_type: %s", err))
	}
	if err = d.Set("status", reservation.Status); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting status: %s", err))
	}
	if reservation.StatusReasons != nil {
		statusReasonsList := []map[string]interface{}{}
		for _, statusReasonsItem := range reservation.StatusReasons {
			statusReasonsList = append(statusReasonsList, dataSourceReservationStatusReasonsToMap(statusReasonsItem))
		}
		if err = d.Set("status_reasons", statusReasonsList); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting status_reasons: %s", err))
		}
	}
	zone := ""
	if reservation.Zone != nil && reservation.Zone.Name != nil {
		zone = *reservation.Zone.Name
	}
	if err = d.Set("zone", zone); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting zone: %s", err))
	}
	return nil
}

func dataSourceReservationProfileToMap(profileItem vpcv1.ReservationProfileIntf) (profileMap map[string]interface{}) {
	var profile *vpcv1.ReservationProfile
	if profileItem == nil {
		return
	}
	if profile = profileItem.(*vpcv1.ReservationProfile); profile == nil {
		return
	}

	profileMap = make(map[string]interface{})

	if profile.Href != nil {
		profileMap["href"] = profile.Href
	}
	if profile.Name != nil {
		profileMap["name"] = profile.Name
	}
	if profile.ResourceType != nil {
		profileMap["resource_type"] = profile.ResourceType
	}
	return profileMap
}

func dataSourceReservationResourceGroupToMap(resourceGroupItem vpcv1.ResourceGroupReference) (resourceGroupMap map[string]interface{}) {
	resourceGroupMap = map[string]interface{}{}

	if resourceGroupItem.Href != nil {
		resourceGroupMap["href"] = resourceGroupItem.Href
	}
	if resourceGroupItem.ID != nil {
		resourceGroupMap["id"] = resourceGroupItem.ID
	}
	if resourceGroupItem.Name != nil {
		resourceGroupMap["name"] = resourceGroupItem.Name
	}
	return resourceGroupMap
}

func dataSourceReservationCapacityToMap(capacityItem vpcv1.ReservationCapacity) (capacityMap map[string]interface{}) {
	capacityMap = map[string]interface{}{}

	if capacityItem.Allocated != nil {
		capacityMap["allocated"] = capacityItem.Allocated
	}
	if capacityItem.Available != nil {
		capacityMap["available"] = capacityItem.Available
	}
	if capacityItem.Total != nil {
		capacityMap["total"] = capacityItem.Total
	}
	if capacityItem.Used != nil {
		capacityMap["used"] = capacityItem.Used
	}
	if capacityItem.Status != nil {
		capacityMap["status"] = capacityItem.Status
	}
	return capacityMap
}

func dataSourceReservationCommittedUseToMap(committedUseItem vpcv1.ReservationCommittedUse) (committedUseMap map[string]interface{}) {
	committedUseMap = map[string]interface{}{}

	if committedUseItem.ExpirationAt != nil {
		committedUseMap["expiration_at"] = committedUseItem.ExpirationAt.String()
	}
	if committedUseItem.ExpirationPolicy != nil {
		committedUseMap["expiration_policy"] = committedUseItem.ExpirationPolicy
	}
	if committedUseItem.Term != nil {
		committedUseMap["term"] = committedUseItem.Term
	}
	return committedUseMap
}

func dataSourceReservationStatusReasonsToMap(statusReasonsItem vpcv1.ReservationStatusReason) (statusReasonsMap map[string]interface{}) {
	statusReasonsMap = map[string]interface{}{}

	if statusReasonsItem.Code != nil {
		statusReasonsMap["code"] = statusReasonsItem.Code
	}
	if statusReasonsItem.Message != nil {
		statusReasonsMap["message"] = statusReasonsItem.Message
	}
	if statusReasonsItem.MoreInfo != nil {
		statusReasonsMap["more_info"] = statusReasonsItem.MoreInfo
	}
	return statusReasonsMap
}
