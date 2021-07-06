// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"strings"
	"time"

	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	obMonitoringCluster         = "cluster"
	obMonitoringInstanceID      = "instance_id"
	obMonitoringInstanceName    = "instance_name"
	obMonitoringIngestionkey    = "sysdig_access_key"
	obMonitoringPrivateEndpoint = "private_endpoint"
	obMonitoringDaemonSetName   = "daemonset_name"
	obMonitoringAgentKey        = "agent_key"
	obMonitoringAgentNamespace  = "agent_namespace"
	obMonitoringCrn             = "crn"
	obMonitoringDiscoveredAgent = "discovered_agent"
	obMonitoringNamespace       = "namespace"
)

func resourceIBMObMonitoring() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMMonitoringCreate,
		Read:     resourceIBMMonitoringRead,
		Update:   resourceIBMMonitoringUpdate,
		Delete:   resourceIBMMonitoringDelete,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			obMonitoringCluster: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name or ID of the cluster to be used.",
			},

			obMonitoringInstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the Sysdig service instance to latch",
			},

			obMonitoringIngestionkey: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Sysdig ingestion key",
			},

			obMonitoringPrivateEndpoint: {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Add this option to connect to your Sysdig service instance through the private service endpoint",
			},

			obMonitoringDaemonSetName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Daemon Set Name",
			},

			obMonitoringInstanceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Sysdig instance Name",
			},

			obMonitoringAgentKey: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Agent key name",
			},

			obMonitoringAgentNamespace: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Agent Namespace",
			},

			obMonitoringCrn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRN",
			},

			obMonitoringDiscoveredAgent: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Discovered agent",
			},

			obMonitoringNamespace: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Namespace",
			},
		},
	}
}

func resourceIBMMonitoringCreate(d *schema.ResourceData, meta interface{}) error {

	client, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}

	var ingestionkey string
	var privateEndpoint bool

	//Read cluster ID and sysdif instanceID
	clusterName := d.Get(obMonitoringCluster).(string)
	sysdigInstanceID := d.Get(obMonitoringInstanceID).(string)

	//Read Ingestionkey
	if iKey, ok := d.GetOk(obMonitoringIngestionkey); ok {
		ingestionkey = iKey.(string)
	}

	//Read private enpoint
	if endPoint, ok := d.GetOk(obMonitoringPrivateEndpoint); ok {
		privateEndpoint = endPoint.(bool)
	}

	//populate sysdig configure create request
	params := v2.MonitoringCreateRequest{
		Cluster:         clusterName,
		IngestionKey:    ingestionkey,
		SysidigInstance: sysdigInstanceID,
		PrivateEndpoint: privateEndpoint,
	}

	targetEnv, err := getMonitoringTargetHeader(d, meta)
	if err != nil {
		return err
	}

	var monitoring v2.MonitoringCreateResponse
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		var err error
		monitoring, err = client.Monitoring().CreateMonitoringConfig(params, targetEnv)
		if err != nil {
			log.Printf("[DEBUG] monitoring Instance err %s", err)
			if strings.Contains(err.Error(), "The user doesn't have enough privileges to perform this action") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if isResourceTimeoutError(err) {
		monitoring, err = client.Monitoring().CreateMonitoringConfig(params, targetEnv)
	}
	if err != nil {
		return fmt.Errorf("error latching monitoring instance to cluster: %w", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", clusterName, monitoring.InstanceID))

	return resourceIBMMonitoringRead(d, meta)
}

func getMonitoringTargetHeader(d *schema.ResourceData, meta interface{}) (v2.MonitoringTargetHeader, error) {
	_, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return v2.MonitoringTargetHeader{}, err
	}

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return v2.MonitoringTargetHeader{}, err
	}
	accountID := userDetails.userAccount

	targetEnv := v2.MonitoringTargetHeader{
		AccountID: accountID,
	}

	return targetEnv, nil
}

func resourceIBMMonitoringRead(d *schema.ResourceData, meta interface{}) error {

	client, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	clusterName := parts[0]
	monitoringID := parts[1]

	targetEnv, err := getMonitoringTargetHeader(d, meta)
	if err != nil {
		return err
	}

	config, err := client.Monitoring().GetMonitoringConfig(clusterName, monitoringID, targetEnv)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("Error in GetMonitoringConfig: %s", err)
	}

	d.Set(obMonitoringPrivateEndpoint, config.PrivateEndpoint)
	d.Set(obMonitoringDaemonSetName, config.DaemonsetName)
	d.Set(obMonitoringInstanceName, config.InstanceName)
	d.Set(obMonitoringAgentKey, config.AgentKey)
	d.Set(obMonitoringAgentNamespace, config.AgentNamespace)
	d.Set(obMonitoringDiscoveredAgent, config.DiscoveredAgent)
	d.Set(obMonitoringCrn, config.CRN)
	d.Set(obMonitoringNamespace, config.Namespace)

	return nil

}

func resourceIBMMonitoringUpdate(d *schema.ResourceData, meta interface{}) error {

	hasChanged := false
	idChanged := false

	client, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}

	targetEnv, err := getMonitoringTargetHeader(d, meta)
	if err != nil {
		return err
	}

	monitoringUpdateModel := v2.MonitoringUpdateRequest{}

	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	cluster := parts[0]
	monitoringID := parts[1]

	clusterName := cluster
	monitoringUpdateModel.Cluster = clusterName

	//if d.HasChange(obMonitoringInstanceID) && !d.IsNewResource() {
	if d.HasChange(obMonitoringInstanceID) {
		hasChanged = true
		idChanged = true
		old, new := d.GetChange(obMonitoringInstanceID)
		monitoringUpdateModel.Instance = old.(string)
		monitoringUpdateModel.NewInstance = new.(string)
	} else {
		monitoringUpdateModel.Instance = monitoringID
	}

	if d.HasChange(obMonitoringIngestionkey) {
		key := d.Get(obMonitoringIngestionkey).(string)
		monitoringUpdateModel.IngestionKey = key
		hasChanged = true
	}

	if d.HasChange(obMonitoringPrivateEndpoint) {
		endpoint := d.Get(obMonitoringPrivateEndpoint).(bool)
		monitoringUpdateModel.PrivateEndpoint = endpoint
		hasChanged = true
	}

	if hasChanged {
		_, err := client.Monitoring().UpdateMonitoringConfig(monitoringUpdateModel, targetEnv)
		if err != nil {
			return err
		} else if idChanged {
			d.SetId(fmt.Sprintf("%s/%s", clusterName, monitoringUpdateModel.NewInstance))
		}
	}

	return resourceIBMMonitoringRead(d, meta)
}

func resourceIBMMonitoringDelete(d *schema.ResourceData, meta interface{}) error {

	client, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}

	targetEnv, err := getMonitoringTargetHeader(d, meta)
	if err != nil {
		return err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	clusterName := parts[0]
	monitoringID := parts[1]

	//populate sysdig configure create request
	params := v2.MonitoringDeleteRequest{
		Cluster:  clusterName,
		Instance: monitoringID,
	}

	_, err = client.Monitoring().DeleteMonitoringConfig(params, targetEnv)

	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("Error in DeleteMonitoringConfig: %s", err)
	}
	d.SetId("")
	return nil

}
