// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/softlayer/softlayer-go/filter"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/sl"
)

func resourceIBMComputePlacementGroup() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMComputePlacementGroupCreate,
		Read:     resourceIBMComputePlacementGroupRead,
		Update:   resourceIBMComputePlacementGroupUpdate,
		Delete:   resourceIBMComputePlacementGroupDelete,
		Exists:   resourceIBMComputePlacementGroupExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"datacenter": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Dataceneter name",
			},

			"pod": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.TrimSpace(old) == strings.TrimSpace(new)
				},
				Description: "Pod name",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name",
			},

			"rule": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "SPREAD",
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue([]string{"SPREAD"}),
				Description:  "Rule info",
			},

			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "List of tags",
			},
		},
	}
}

func resourceIBMComputePlacementGroupCreate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	name := d.Get("name").(string)
	datacenter := d.Get("datacenter").(string)
	pod := d.Get("pod").(string)
	podName := datacenter + "." + pod
	PodService := services.GetNetworkPodService(sess)
	podMask := `backendRouterId,name`
	rule := d.Get("rule").(string)

	// 1.Getting the router ID
	routerids, err := PodService.Filter(filter.Path("datacenterName").Eq(datacenter).Build()).Mask(podMask).GetAllObjects()
	if err != nil {
		return fmt.Errorf("Encountered problem trying to get the router ID: %s", err)
	}
	var routerid int
	for _, iterate := range routerids {
		if *iterate.Name == podName {
			routerid = *iterate.BackendRouterId
		}
	}

	ruleService := services.GetVirtualPlacementGroupRuleService(sess)
	ruleObject, err := ruleService.Id(1).
		Mask("id,name").
		Filter(filter.Path("name").Eq(rule).Build()).GetObject()
	if err != nil {
		return fmt.Errorf("Encountered problem trying to get the placement group rule ID: %s", err)
	}

	opts := datatypes.Virtual_PlacementGroup{
		Name:            sl.String(name),
		BackendRouterId: &routerid,
		RuleId:          ruleObject.Id,
	}

	service := services.GetVirtualPlacementGroupService(sess)

	pgrp, err := service.CreateObject(&opts)
	if err != nil {
		return fmt.Errorf("Error creating Placement Group: %s", err)
	}

	d.SetId(strconv.Itoa(*pgrp.Id))
	log.Printf("[INFO] Placement Group ID: %d", *pgrp.Id)

	return resourceIBMComputePlacementGroupRead(d, meta)
}

func resourceIBMComputePlacementGroupRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetVirtualPlacementGroupService(sess)

	pgrpID, _ := strconv.Atoi(d.Id())

	pgrp, err := service.Id(pgrpID).Mask("id,name,rule[name],backendRouter[hostname,datacenter[name]]").GetObject()
	if err != nil {
		if err, ok := err.(sl.Error); ok {
			if err.StatusCode == http.StatusNotFound {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("Error retrieving Placement Group: %s", err)
	}

	d.Set("name", pgrp.Name)
	d.Set("datacenter", pgrp.BackendRouter.Datacenter.Name)
	pod := strings.SplitAfter(*pgrp.BackendRouter.Hostname, ".")[0]
	r, _ := regexp.Compile("[0-9]{2}")
	pod = "pod" + r.FindString(pod)
	d.Set("pod", pod)
	d.Set("rule", pgrp.Rule.Name)

	return nil
}

func resourceIBMComputePlacementGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetVirtualPlacementGroupService(sess.SetRetries(0))

	pgrpID, _ := strconv.Atoi(d.Id())

	opts := datatypes.Virtual_PlacementGroup{}

	if d.HasChange("name") {
		opts.Name = sl.String(d.Get("name").(string))
		_, err := service.Id(pgrpID).EditObject(&opts)

		if err != nil {
			return fmt.Errorf("Error editing Placement Group: %s", err)
		}
	}

	return nil
}

func resourceIBMComputePlacementGroupExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetVirtualPlacementGroupService(sess)

	pgrpID, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	result, err := service.Id(pgrpID).GetObject()
	if err != nil {
		if apiErr, ok := err.(sl.Error); ok {
			if apiErr.StatusCode == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}
	return result.Id != nil && *result.Id == pgrpID, nil
}

func resourceIBMComputePlacementGroupDelete(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetVirtualPlacementGroupService(sess)

	pgrpID, err := strconv.Atoi(d.Id())
	log.Printf("[INFO] Deleting Placement Group: %d", pgrpID)

	const (
		noVms                    = "There are no vms on the Placement Group"
		vmsStillOnPlacementGroup = "VMs are still present on the Placement Group"
	)

	//Wait till all the VMs are disconnected before trying to delete
	stateConf := &resource.StateChangeConf{
		Target:     []string{noVms},
		Pending:    []string{vmsStillOnPlacementGroup},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second, // Wait 30 secs before starting
		Refresh: func() (interface{}, string, error) {
			vms, err := service.Id(pgrpID).GetGuests()
			if err != nil {
				log.Printf("[ERROR] Received error while fetching virtual guests on placement group to see if placement group can be cancelled now: %#v", err)
				return vms, "Error", err
			}
			if len(vms) != 0 {
				return vms, vmsStillOnPlacementGroup, nil
			}
			return vms, noVms, nil
		},
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return err
	}

	_, err = service.Id(pgrpID).DeleteObject()
	if err != nil {
		return fmt.Errorf("Error deleting Placement Group: %s", err)
	}

	return nil
}
