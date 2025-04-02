// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsReservations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsReservationsRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Filters the collection to resources with the exact specified name",
				Optional:    true,
			},
			"zone_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to resources with a zone name property matching the exact specified name.",
			},
			"resource_group": {
				Type:        schema.TypeString,
				Description: "Filters the collection to resources in the resource group with the specified identifier",
				Optional:    true,
			},
			"reservations": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of reservations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique user-defined name for this reservation.",
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
				},
			},
		},
	}
}

func dataSourceIBMIsReservationsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	zoneName := d.Get("zone_name").(string)
	resourceGroupId := d.Get("resource_group").(string)

	start := ""
	reservations := []vpcv1.Reservation{}

	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	for {
		listReservationsOptions := &vpcv1.ListReservationsOptions{}
		if start != "" {
			listReservationsOptions.Start = &start
		}
		if name != "" {
			listReservationsOptions.SetName(name)
		}
		if zoneName != "" {
			listReservationsOptions.SetZoneName(zoneName)
		}
		if resourceGroupId != "" {
			listReservationsOptions.SetResourceGroupID(resourceGroupId)
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
		reservations = append(reservations, reservationCollection.Reservations...)
		if start == "" {
			break
		}
	}

	d.SetId(dataSourceIBMIsReservationsID(d))

	if reservations != nil {
		err = d.Set("reservations", dataSourceReservationCollectionFlattenReservations(reservations))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting reservations %s", err))
		}
	}

	return nil
}

// dataSourceIBMIsReservationsID returns a reasonable ID for the list.
func dataSourceIBMIsReservationsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceReservationCollectionFlattenReservations(result []vpcv1.Reservation) (reservations []map[string]interface{}) {
	for _, reservationsItem := range result {
		reservations = append(reservations, dataSourceReservationCollectionReservationsToMap(reservationsItem))
	}
	return reservations
}

func dataSourceReservationCollectionReservationsToMap(reservationsItem vpcv1.Reservation) (reservationsMap map[string]interface{}) {
	reservationsMap = map[string]interface{}{}

	if reservationsItem.AffinityPolicy != nil {
		reservationsMap["affinity_policy"] = reservationsItem.AffinityPolicy
	}
	if reservationsItem.Capacity != nil {
		capacityList := []map[string]interface{}{}
		capacityMap := dataSourceReservationCollectionReservationsCapacityToMap(*reservationsItem.Capacity)
		capacityList = append(capacityList, capacityMap)
		reservationsMap["capacity"] = capacityList
	}
	if reservationsItem.CommittedUse != nil {
		committedUseList := []map[string]interface{}{}
		committedUseMap := dataSourceReservationCollectionReservationsCommittedUseToMap(*reservationsItem.CommittedUse)
		committedUseList = append(committedUseList, committedUseMap)
		reservationsMap["committed_use"] = committedUseList
	}
	if reservationsItem.CreatedAt != nil {
		reservationsMap["created_at"] = reservationsItem.CreatedAt.String()
	}
	if reservationsItem.CRN != nil {
		reservationsMap["crn"] = reservationsItem.CRN
	}
	if reservationsItem.Href != nil {
		reservationsMap["href"] = reservationsItem.Href
	}
	if reservationsItem.ID != nil {
		reservationsMap["id"] = reservationsItem.ID
	}
	if reservationsItem.LifecycleState != nil {
		reservationsMap["lifecycle_state"] = reservationsItem.LifecycleState
	}
	if reservationsItem.Name != nil {
		reservationsMap["name"] = reservationsItem.Name
	}
	if reservationsItem.Profile != nil {
		profileList := []map[string]interface{}{}
		profile := reservationsItem.Profile
		profileMap := dataSourceReservationCollectionReservationsProfileToMap(profile)
		profileList = append(profileList, profileMap)
		reservationsMap["profile"] = profileList
	}
	if reservationsItem.ResourceGroup != nil {
		resourceGroupList := []map[string]interface{}{}
		resourceGroupMap := dataSourceReservationCollectionReservationsResourceGroupToMap(*reservationsItem.ResourceGroup)
		resourceGroupList = append(resourceGroupList, resourceGroupMap)
		reservationsMap["resource_group"] = resourceGroupList
	}
	if reservationsItem.ResourceType != nil {
		reservationsMap["resource_type"] = reservationsItem.ResourceType
	}
	if reservationsItem.Status != nil {
		reservationsMap["status"] = reservationsItem.Status
	}
	if reservationsItem.StatusReasons != nil {
		statusReasonsList := []map[string]interface{}{}
		for _, statusReasonsItem := range reservationsItem.StatusReasons {
			statusReasonsList = append(statusReasonsList, dataSourceReservationCollectionReservationsStatusReasonsToMap(statusReasonsItem))
		}
		reservationsMap["status_reasons"] = statusReasonsList
	}
	if reservationsItem.Zone != nil && reservationsItem.Zone.Name != nil {
		reservationsMap["zone"] = *reservationsItem.Zone.Name
	}
	return reservationsMap
}

func dataSourceReservationCollectionReservationsProfileToMap(profileItem vpcv1.ReservationProfileIntf) (profileMap map[string]interface{}) {
	if profileItem == nil {
		return
	}
	var profile *vpcv1.ReservationProfile
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

func dataSourceReservationCollectionReservationsResourceGroupToMap(resourceGroupItem vpcv1.ResourceGroupReference) (resourceGroupMap map[string]interface{}) {
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

func dataSourceReservationCollectionReservationsCapacityToMap(capacityItem vpcv1.ReservationCapacity) (capacityMap map[string]interface{}) {
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

func dataSourceReservationCollectionReservationsCommittedUseToMap(committedUseItem vpcv1.ReservationCommittedUse) (committedUseMap map[string]interface{}) {
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

func dataSourceReservationCollectionReservationsStatusReasonsToMap(statusReasonsItem vpcv1.ReservationStatusReason) (statusReasonsMap map[string]interface{}) {
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
