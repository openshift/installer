// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
)

const (
	enablePrivateSECmdAction = "enablePrivateServiceEndpoint"
	enablePublicSECmdAction  = "enablePublicServiceEndpoint"
	disablePublicSECmdAction = "disablePublicServiceEndpoint"
	reloadAction             = "reload"
)

func resourceIBMContainerClusterFeature() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMContainerClusterFeatureCreate,
		Read:     resourceIBMContainerClusterFeatureRead,
		Update:   resourceIBMContainerClusterFeatureUpdate,
		Delete:   resourceIBMContainerClusterFeatureDelete,
		Importer: &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			"cluster": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster name of ID",
			},

			"public_service_endpoint": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"private_service_endpoint": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"public_service_endpoint_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"refresh_api_servers": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Boolean value true of API server to be refreshed in K8S cluster",
			},

			"reload_workers": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Boolean value set true if worker nodes to be reloaded",
			},

			"private_service_endpoint_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"resource_group_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "ID of the resource group.",
				Computed:         true,
				DiffSuppressFunc: applyOnce,
			},
		},
	}
}

func resourceIBMContainerClusterFeatureCreate(d *schema.ResourceData, meta interface{}) error {

	cluster := d.Get("cluster").(string)
	var isOptionSet bool

	if v, ok := d.GetOkExists("private_service_endpoint"); ok {
		if v.(bool) {
			err := updateCluster(cluster, enablePrivateSECmdAction, d.Timeout(schema.TimeoutCreate), d, meta)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("The `private_service_endpoint` can not be disabled")
		}
		d.SetId(cluster)
		err := reloadCluster(cluster, d.Timeout(schema.TimeoutCreate), d, meta)
		if err != nil {
			return err
		}
		isOptionSet = true
	}

	if v, ok := d.GetOkExists("public_service_endpoint"); ok {
		var cmd string
		if v.(bool) {
			cmd = enablePublicSECmdAction
		} else {
			cmd = disablePublicSECmdAction
		}
		log.Printf("Started enabling the public ep %s", cmd)
		err := updateCluster(cluster, cmd, d.Timeout(schema.TimeoutCreate), d, meta)
		if err != nil {
			return err
		}
		d.SetId(cluster)
		err = reloadCluster(cluster, d.Timeout(schema.TimeoutCreate), d, meta)
		if err != nil {
			return err
		}
		isOptionSet = true
	}

	if !isOptionSet {
		return fmt.Errorf("Provide either `public_service_endpoint` or  `private_service_endpoint` or both.")
	}
	return resourceIBMContainerClusterFeatureRead(d, meta)
}

func reloadCluster(cluster string, timeout time.Duration, d *schema.ResourceData, meta interface{}) error {
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return err
	}
	targetEnv, err := getWorkerPoolTargetHeader(d, meta)
	if err != nil {
		return err
	}
	if v, ok := d.GetOkExists("refresh_api_servers"); ok {
		if v.(bool) {
			err = csClient.Clusters().RefreshAPIServers(cluster, targetEnv)
			if err != nil {
				return err
			}
		}
	}
	if v, ok := d.GetOkExists("reload_workers"); ok {
		if v.(bool) {
			log.Printf("Waiting for cluster (%s) to be available.", cluster)
			_, err = WaitForClusterAvailableForFeatureUpdate(cluster, timeout, meta, targetEnv)
			if err != nil {
				return fmt.Errorf(
					"Error waiting for cluster (%s) to become ready: %s", cluster, err)
			}
			log.Printf("Waiting for workers (%s) to be available.", cluster)
			_, err = WaitForWorkerAvailableForFeatureUpdate(cluster, timeout, meta, targetEnv)
			if err != nil {
				return fmt.Errorf(
					"Error waiting for workers of cluster (%s) to become ready: %s", cluster, err)
			}
			params := v1.UpdateWorkerCommand{
				Action: reloadAction,
			}
			workerFields, err := csClient.Workers().List(cluster, targetEnv)
			if err != nil {
				return fmt.Errorf("Error retrieving workers for cluster: %s", err)
			}
			workers := make([]string, len(workerFields))
			for i, worker := range workerFields {
				workers[i] = worker.ID
			}
			err = csClient.Clusters().UpdateClusterWorkers(cluster, workers, params, targetEnv)
			if err != nil {
				return err
			}
			_, err = WaitForClusterAvailableForFeatureUpdate(cluster, timeout, meta, targetEnv)
			if err != nil {
				return fmt.Errorf(
					"Error waiting for cluster (%s) to become ready: %s", d.Id(), err)
			}
			_, err = WaitForWorkerAvailableForFeatureUpdate(cluster, timeout, meta, targetEnv)
			if err != nil {
				return fmt.Errorf(
					"Error waiting for workers of cluster (%s) to become ready: %s", d.Id(), err)
			}
		}
	}

	return nil
}

func updateCluster(cluster, actionCmd string, timeout time.Duration, d *schema.ResourceData, meta interface{}) error {
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return err
	}
	params := v1.ClusterUpdateParam{
		Action: actionCmd,
	}
	targetEnv, err := getWorkerPoolTargetHeader(d, meta)
	if err != nil {
		return err
	}
	log.Printf("Waiting for cluster (%s) to be available.", cluster)
	_, err = WaitForClusterAvailableForFeatureUpdate(cluster, timeout, meta, targetEnv)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for cluster (%s) to become ready: %s", d.Id(), err)
	}
	log.Printf("Waiting for workers (%s) to be available.", cluster)
	_, err = WaitForWorkerAvailableForFeatureUpdate(cluster, timeout, meta, targetEnv)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for workers of cluster (%s) to become ready: %s", d.Id(), err)
	}
	log.Printf("Calling update with action cmd %s", actionCmd)
	err = csClient.Clusters().Update(cluster, params, targetEnv)
	if err != nil {
		return err
	}
	log.Printf("success with action cmd %s", actionCmd)

	return nil
}

func resourceIBMContainerClusterFeatureRead(d *schema.ResourceData, meta interface{}) error {
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return err
	}

	clusterID := d.Id()
	targetEnv, err := getWorkerPoolTargetHeader(d, meta)
	if err != nil {
		return err
	}
	cls, err := csClient.Clusters().Find(clusterID, targetEnv)
	if err != nil {
		return fmt.Errorf("Error retrieving armada cluster: %s", err)
	}

	d.Set("cluster", clusterID)
	d.Set("public_service_endpoint", cls.PublicServiceEndpointEnabled)
	d.Set("private_service_endpoint_url", cls.PrivateServiceEndpointURL)
	d.Set("public_service_endpoint_url", cls.PublicServiceEndpointURL)
	d.Set("private_service_endpoint", cls.PrivateServiceEndpointEnabled)

	return nil
}

func resourceIBMContainerClusterFeatureDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

func resourceIBMContainerClusterFeatureUpdate(d *schema.ResourceData, meta interface{}) error {

	cluster := d.Get("cluster").(string)
	var isOptionSet bool
	if d.HasChange("private_service_endpoint") {
		if v, ok := d.GetOkExists("private_service_endpoint"); ok {
			if v.(bool) {
				err := updateCluster(cluster, enablePrivateSECmdAction, d.Timeout(schema.TimeoutUpdate), d, meta)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("The `private_service_endpoint` can not be disabled")
			}
			err := reloadCluster(cluster, d.Timeout(schema.TimeoutUpdate), d, meta)
			if err != nil {
				return err
			}
			isOptionSet = true
		}
	}
	if d.HasChange("public_service_endpoint") {
		if v, ok := d.GetOkExists("public_service_endpoint"); ok {
			var cmd string
			if v.(bool) {
				cmd = enablePublicSECmdAction
			} else {
				cmd = disablePublicSECmdAction
			}
			err := updateCluster(cluster, cmd, d.Timeout(schema.TimeoutUpdate), d, meta)
			if err != nil {
				return err
			}
			err = reloadCluster(cluster, d.Timeout(schema.TimeoutUpdate), d, meta)
			if err != nil {
				return err
			}
			isOptionSet = true
		}
	}

	if !isOptionSet {
		return fmt.Errorf("Provide either `public_service_endpoint` or  `private_service_endpoint` or both.")
	}

	return resourceIBMContainerClusterFeatureRead(d, meta)
}

// WaitForClusterAvailable Waits for cluster creation
func WaitForClusterAvailableForFeatureUpdate(cluster string, timeout time.Duration, meta interface{}, target v1.ClusterTargetHeader) (interface{}, error) {
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return nil, err
	}
	log.Printf("Waiting for cluster (%s) to be available.", cluster)
	id := cluster

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", clusterProvisioning},
		Target:     []string{clusterNormal},
		Refresh:    clusterStateRefreshFunc(csClient.Clusters(), id, target),
		Timeout:    timeout,
		Delay:      60 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func WaitForWorkerAvailableForFeatureUpdate(cluster string, timeout time.Duration, meta interface{}, target v1.ClusterTargetHeader) (interface{}, error) {
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return nil, err
	}
	log.Printf("Waiting for worker of the cluster (%s) to be available.", cluster)
	id := cluster

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", workerProvisioning},
		Target:     []string{workerNormal},
		Refresh:    workerStateRefreshFunc(csClient.Workers(), id, target),
		Timeout:    timeout,
		Delay:      60 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}
