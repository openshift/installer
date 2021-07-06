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
	obLoggingCluster         = "cluster"
	obLoggingInstanceID      = "instance_id"
	obLoggingInstanceName    = "instance_name"
	obLoggingIngestionkey    = "logdna_ingestion_key"
	obLoggingPrivateEndpoint = "private_endpoint"
	obLoggingDaemonSetName   = "daemonset_name"
	obLoggingAgentKey        = "agent_key"
	obLoggingAgentNamespace  = "agent_namespace"
	obLoggingCrn             = "crn"
	obLoggingDiscoveredAgent = "discovered_agent"
	obLoggingNamespace       = "namespace"
)

func resourceIBMObLogging() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMLoggingCreate,
		Read:     resourceIBMLoggingRead,
		Update:   resourceIBMLoggingUpdate,
		Delete:   resourceIBMLoggingDelete,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			obLoggingCluster: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name or ID of the cluster to be used.",
			},

			obLoggingInstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the LogDNA service instance to latch",
			},

			obLoggingIngestionkey: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "LogDNA ingestion key",
			},

			obLoggingPrivateEndpoint: {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Add this option to connect to your LogDNA service instance through the private service endpoint",
			},

			obLoggingDaemonSetName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Daemon Set Name",
			},

			obLoggingInstanceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "LogDNA instance Name",
			},

			obLoggingAgentKey: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Agent key name",
			},

			obLoggingAgentNamespace: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Agent Namespace",
			},

			obLoggingCrn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRN",
			},

			obLoggingDiscoveredAgent: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Discovered agent",
			},

			obLoggingNamespace: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Namespace",
			},
		},
	}
}

func resourceIBMLoggingCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}

	var ingestionkey string
	var privateEndpoint bool

	//Read cluster ID and logging instanceID
	clusterName := d.Get(obLoggingCluster).(string)
	loggingInstanceID := d.Get(obLoggingInstanceID).(string)

	//Read Ingestionkey
	if iKey, ok := d.GetOk(obLoggingIngestionkey); ok {
		ingestionkey = iKey.(string)
	}

	//Read private enpoint
	if endPoint, ok := d.GetOk(obLoggingPrivateEndpoint); ok {
		privateEndpoint = endPoint.(bool)
	}

	//populate sysdig configure create request
	params := v2.LoggingCreateRequest{
		Cluster:         clusterName,
		IngestionKey:    ingestionkey,
		LoggingInstance: loggingInstanceID,
		PrivateEndpoint: privateEndpoint,
	}

	targetEnv, err := getLoggingTargetHeader(d, meta)
	if err != nil {
		return err
	}

	var logging v2.LoggingCreateResponse
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		var err error
		logging, err = client.Logging().CreateLoggingConfig(params, targetEnv)
		if err != nil {
			log.Printf("[DEBUG] logging Instance err %s", err)
			if strings.Contains(err.Error(), "The user doesn't have enough privileges to perform this action") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if isResourceTimeoutError(err) {
		logging, err = client.Logging().CreateLoggingConfig(params, targetEnv)
	}
	if err != nil {
		return fmt.Errorf("error latching logging instance to cluster: %w", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", clusterName, logging.InstanceID))

	return resourceIBMLoggingRead(d, meta)
}

func getLoggingTargetHeader(d *schema.ResourceData, meta interface{}) (v2.LoggingTargetHeader, error) {
	_, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return v2.LoggingTargetHeader{}, err
	}

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return v2.LoggingTargetHeader{}, err
	}
	accountID := userDetails.userAccount

	targetEnv := v2.LoggingTargetHeader{
		AccountID: accountID,
	}

	return targetEnv, nil
}

func resourceIBMLoggingRead(d *schema.ResourceData, meta interface{}) error {

	client, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	clusterName := parts[0]
	loggingID := parts[1]

	targetEnv, err := getLoggingTargetHeader(d, meta)
	if err != nil {
		return err
	}

	config, err := client.Logging().GetLoggingConfig(clusterName, loggingID, targetEnv)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("Error in GetLoggingConfig: %s", err)
	}

	d.Set(obLoggingPrivateEndpoint, config.PrivateEndpoint)
	d.Set(obLoggingDaemonSetName, config.DaemonsetName)
	d.Set(obLoggingInstanceName, config.InstanceName)
	d.Set(obLoggingAgentKey, config.AgentKey)
	d.Set(obLoggingAgentNamespace, config.AgentNamespace)
	d.Set(obLoggingDiscoveredAgent, config.DiscoveredAgent)
	d.Set(obLoggingCrn, config.CRN)
	d.Set(obLoggingNamespace, config.Namespace)

	return nil

}

func resourceIBMLoggingUpdate(d *schema.ResourceData, meta interface{}) error {

	hasChanged := false
	idChanged := false

	client, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}

	targetEnv, err := getLoggingTargetHeader(d, meta)
	if err != nil {
		return err
	}

	loggingUpdateModel := v2.LoggingUpdateRequest{}

	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	cluster := parts[0]
	loggingID := parts[1]

	clusterName := cluster
	loggingUpdateModel.Cluster = clusterName

	if d.HasChange(obLoggingInstanceID) {
		hasChanged = true
		idChanged = true
		old, new := d.GetChange(obLoggingInstanceID)
		loggingUpdateModel.Instance = old.(string)
		loggingUpdateModel.NewInstance = new.(string)
	} else {
		loggingUpdateModel.Instance = loggingID
	}

	if d.HasChange(obLoggingIngestionkey) {
		key := d.Get(obLoggingIngestionkey).(string)
		loggingUpdateModel.IngestionKey = key
		hasChanged = true
	}

	if d.HasChange(obLoggingPrivateEndpoint) {
		endpoint := d.Get(obLoggingPrivateEndpoint).(bool)
		loggingUpdateModel.PrivateEndpoint = endpoint
		hasChanged = true
	}

	if hasChanged {

		_, err := client.Logging().UpdateLoggingConfig(loggingUpdateModel, targetEnv)
		if err != nil {
			return err
		} else if idChanged {
			d.SetId(fmt.Sprintf("%s/%s", clusterName, loggingUpdateModel.NewInstance))
		}
	}

	return resourceIBMLoggingRead(d, meta)
}

func resourceIBMLoggingDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}

	targetEnv, err := getLoggingTargetHeader(d, meta)
	if err != nil {
		return err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	clusterName := parts[0]
	loggingID := parts[1]

	//populate logging logDNA configure create request
	params := v2.LoggingDeleteRequest{
		Cluster:  clusterName,
		Instance: loggingID,
	}

	_, err = client.Logging().DeleteLoggingConfig(params, targetEnv)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("Error in DeleteLoggingConfig: %s", err)
	}
	d.SetId("")
	return nil

}
