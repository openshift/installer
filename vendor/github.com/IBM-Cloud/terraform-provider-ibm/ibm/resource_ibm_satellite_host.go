// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.ibm.com/ibmcloud/kubernetesservice-go-sdk/kubernetesserviceapiv1"
)

const (
	hostCluster    = "cluster"
	hostLocation   = "location"
	hostID         = "host_id"
	hostState      = "host_state"
	hostLabels     = "labels"
	hostZone       = "zone"
	hostWorkerPool = "worker_pool"
	hostProvider   = "host_provider"

	rsHostNormalStatus       = "normal"
	rsHostProvisioningStatus = "provisioning"
	rsHostReadyStatus        = "ready"
	rsHostUnknownStatus      = "unknown"
)

func resourceIBMSatelliteHost() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMSatelliteHostCreate,
		Read:     resourceIBMSatelliteHostRead,
		Update:   resourceIBMSatelliteHostUpdate,
		Delete:   resourceIBMSatelliteHostDelete,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(75 * time.Minute),
			Read:   schema.DefaultTimeout(75 * time.Minute),
			Delete: schema.DefaultTimeout(45 * time.Minute),
			Update: schema.DefaultTimeout(45 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			hostLocation: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name or ID of the Satellite location",
			},
			hostCluster: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name or ID of a Satellite location or cluster to assign the host to",
			},
			hostID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The specific host ID to assign to a Satellite location or cluster",
			},
			hostLabels: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "List of labels for the host",
			},
			hostZone: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The zone within the cluster to assign the host to",
			},
			hostWorkerPool: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name or ID of the worker pool within the cluster to assign the host to",
			},
			hostProvider: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Host Provider",
			},
			hostState: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Health status of the host",
			},
		},
	}
}

func resourceIBMSatelliteHostCreate(d *schema.ResourceData, meta interface{}) error {
	hostName := d.Get(hostID).(string)
	location := d.Get(hostLocation).(string)

	satClient, err := meta.(ClientSession).SatelliteClientSession()
	if err != nil {
		return err
	}

	hostAssignOptions := &kubernetesserviceapiv1.CreateSatelliteAssignmentOptions{}
	hostAssignOptions.Controller = ptrToString(location)

	if _, ok := d.GetOk(hostCluster); ok {
		hostAssignOptions.Cluster = ptrToString(d.Get(hostCluster).(string))
	} else {
		hostAssignOptions.Cluster = ptrToString(location)
	}
	hostAssignOptions.HostID = ptrToString(hostName)

	//Check host attached to location
	hostStatus, err := waitForHostAttachment(hostName, location, d, meta)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for attaching host (%s) to be succeeded: %s", hostName, err)
	}

	labels := make(map[string]string)
	if _, ok := d.GetOk(hostLabels); ok {
		l := d.Get(hostLabels).(*schema.Set)
		labels = flattenHostLabels(l.List())
		hostAssignOptions.Labels = labels
	} else {
		hostAssignOptions.Labels = labels
	}

	if _, ok := d.GetOk(hostWorkerPool); ok {
		hostAssignOptions.Workerpool = ptrToString(d.Get(hostWorkerPool).(string))
	}

	if _, ok := d.GetOk(hostZone); ok {
		hostAssignOptions.Zone = ptrToString(d.Get(hostZone).(string))
	}

	if hostStatus == rsHostReadyStatus {
		_, response, err := satClient.CreateSatelliteAssignment(hostAssignOptions)
		if err != nil {
			return fmt.Errorf("Error Assigning Satellite Host: %s\n%s", err, response)
		}
	}

	d.SetId(fmt.Sprintf("%s/%s", location, hostName))

	//Wait for host to reach normal state
	_, err = waitForHostAttachment(hostName, location, d, meta)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for host (%s) to get normal state: %s", hostName, err)
	}

	return resourceIBMSatelliteHostRead(d, meta)
}

func resourceIBMSatelliteHostRead(d *schema.ResourceData, meta interface{}) error {
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	location := parts[0]
	hostName := parts[1]

	satClient, err := meta.(ClientSession).SatelliteClientSession()
	if err != nil {
		return err
	}

	hostOptions := &kubernetesserviceapiv1.GetSatelliteHostsOptions{
		Controller: &location,
	}
	hostList, resp, err := satClient.GetSatelliteHosts(hostOptions)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Println("resourceIBMSatelliteHostRead : error in getting hostlist :", err, resp)
		return err
	}

	for _, h := range hostList {
		if hostName == *h.Name || hostName == *h.ID {
			d.Set(hostLocation, location)
			d.Set("host_id", hostName)

			if _, ok := d.GetOk(hostLabels); ok {
				l := d.Get(hostLabels).(*schema.Set)
				d.Set(hostLabels, l)
			}

			if h.Health != nil {
				d.Set(hostState, *h.Health.Status)
			}

			if _, ok := d.GetOk(hostCluster); ok {
				d.Set(hostCluster, d.Get(hostCluster).(string))
			} else {
				d.Set(hostCluster, location)
			}

			if h.Assignment != nil {
				d.Set(hostWorkerPool, *h.Assignment.WorkerPoolName)
				d.Set(hostZone, *h.Assignment.Zone)
			}
		}
	}

	return nil
}

func resourceIBMSatelliteHostUpdate(d *schema.ResourceData, meta interface{}) error {
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}

	locationName := parts[0]
	hostID := parts[1]
	satClient, err := meta.(ClientSession).SatelliteClientSession()
	if err != nil {
		return err
	}

	updateHostOptions := &kubernetesserviceapiv1.UpdateSatelliteHostOptions{}
	updateHostOptions.Controller = &locationName
	updateHostOptions.HostID = &hostID

	if v, ok := d.GetOk(hostState); ok && v != nil && v.(string) == rsHostReadyStatus {
		labels := make(map[string]string)
		if _, ok := d.GetOk(hostLabels); ok {
			l := d.Get(hostLabels).(*schema.Set)
			labels = flattenHostLabels(l.List())
			updateHostOptions.Labels = labels
		}
		response, err := satClient.UpdateSatelliteHost(updateHostOptions)
		if err != nil {
			return fmt.Errorf("Error Updating Satellite Host: %s\n%s", err, response)
		}
	}

	return resourceIBMSatelliteHostRead(d, meta)
}

func resourceIBMSatelliteHostDelete(d *schema.ResourceData, meta interface{}) error {
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}

	location := parts[0]
	hostID := parts[1]
	satClient, err := meta.(ClientSession).SatelliteClientSession()
	if err != nil {
		return err
	}

	removeSatHostOptions := &kubernetesserviceapiv1.RemoveSatelliteHostOptions{}
	removeSatHostOptions.Controller = &location
	removeSatHostOptions.HostID = &hostID

	response, err := satClient.RemoveSatelliteHost(removeSatHostOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("Error Deleting Satellite Host: %s\n%s", err, response)
	}

	d.SetId("")
	return nil
}

func waitForHostAttachment(hostName, location string, d *schema.ResourceData, meta interface{}) (interface{}, error) {
	satClient, err := meta.(ClientSession).SatelliteClientSession()
	if err != nil {
		return false, err
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{rsHostProvisioningStatus, rsHostUnknownStatus},
		Target:  []string{rsHostReadyStatus, rsHostNormalStatus},
		Refresh: func() (interface{}, string, error) {
			attachOptions := &kubernetesserviceapiv1.GetSatelliteHostsOptions{
				Controller: &location,
			}
			hostList, resp, err := satClient.GetSatelliteHosts(attachOptions)
			if err != nil {
				if apiErr, ok := err.(bmxerror.RequestFailure); ok && apiErr.StatusCode() != 404 {
					return nil, "", fmt.Errorf("The satellite host (%s) failed to attached: %v\n%s", hostName, err, resp)
				}
			}

			if hostList != nil {
				for _, h := range hostList {
					if h.Health != nil {
						if (hostName == *h.Name) && (*h.Health.Status == rsHostNormalStatus || *h.Health.Status == rsHostReadyStatus) {
							return *h.Health.Status, *h.Health.Status, err
						}
					}
				}
			}
			return hostName, rsHostProvisioningStatus, nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      60 * time.Second,
		MinTimeout: 60 * time.Second,
	}

	return stateConf.WaitForState()
}
