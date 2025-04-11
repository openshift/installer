package satellite

import (
	"fmt"
	"time"

	"github.com/IBM-Cloud/container-services-go-sdk/kubernetesserviceapiv1"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMSatelliteStorageAssignment() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMContainerStorageAssignmentCreate,
		Read:     resourceIBMContainerStorageAssignmentRead,
		Update:   resourceIBMContainerStorageAssignmentUpdate,
		Delete:   resourceIBMContainerStorageAssignmentDelete,
		Importer: &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"assignment_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Assignment.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: "The Universally Unique IDentifier (UUID) of the Assignment.",
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Owner of the Assignment.",
			},
			"groups": {
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"cluster"},
				RequiredWith:  []string{"controller"},
				Description:   "One or more cluster groups on which you want to apply the configuration. Note that at least one cluster group is required. ",
			},
			"cluster": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "ID of the Satellite cluster or Service Cluster that you want to apply the configuration to.",
				DiffSuppressFunc: flex.ApplyOnce,
				RequiredWith:     []string{"controller"},
			},
			"svc_cluster": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the Service Cluster that you applied the configuration to.",
			},
			"sat_cluster": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the Satellite cluster that you applied the configuration to.",
			},
			"config": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "Storage Configuration Name or ID.",
			},
			"config_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Universally Unique IDentifier (UUID) of the Storage Configuration.",
			},
			"config_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Storage Configuration Version.",
			},
			"config_version_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Universally Unique IDentifier (UUID) of the Storage Configuration Version.",
			},
			"assignment_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Type of Assignment.",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Time of Creation of the Assignment.",
			},
			"rollout_success_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The Rollout Success Count of the Assignment.",
			},
			"rollout_error_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The Rollout Error Count of the Assignment.",
			},
			"is_assignment_upgrade_available": {
				Type:        schema.TypeBool,
				Computed:    true,
				ForceNew:    false,
				Description: "Whether an Upgrade is Available for the Assignment.",
			},
			"update_config_revision": {
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				Description: "Updating an assignment to the latest available storage configuration version.",
			},
			"controller": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Name or ID of the Satellite Location.",
			},
		},
	}
}

func resourceIBMContainerStorageAssignmentCreate(d *schema.ResourceData, meta interface{}) error {

	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return err
	}

	var result *kubernetesserviceapiv1.CreateSubscriptionData

	assignmentOptions := &kubernetesserviceapiv1.CreateAssignmentOptions{}

	if v, ok := d.GetOk("assignment_name"); ok {
		name := v.(string)
		assignmentOptions.Name = &name
	}

	if v, ok := d.GetOk("config"); ok {
		config := v.(string)
		assignmentOptions.Config = &config
	}

	if v, ok := d.GetOk("cluster"); ok {
		cluster := v.(string)
		assignmentOptions.Cluster = &cluster
	}

	if v, ok := d.GetOk("controller"); ok {
		controller := v.(string)
		assignmentOptions.Controller = &controller
	}
	// If Groups are defined, create assignment to group function is called or else an assignment is made to a cluster
	if v, groupsOk := d.GetOk("groups"); groupsOk {
		groups := v.([]interface{})
		assignmentOptions.Groups = flex.ExpandStringList(groups)
		result, _, err = satClient.CreateAssignment(assignmentOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Creating Assignment - %v", err)
		}
	} else {
		result, _, err = satClient.CreateAssignmentByCluster(assignmentOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Creating Assignment by Cluster - %v", err)
		}
	}

	d.Set("uuid", *result.AddSubscription.UUID)

	getAssignmentOptions := &kubernetesserviceapiv1.GetAssignmentOptions{
		UUID: result.AddSubscription.UUID,
	}
	_, err = waitForAssignmentCreationStatus(getAssignmentOptions, meta, d)
	if err != nil {
		return err
	}
	d.SetId(*result.AddSubscription.UUID)

	return resourceIBMContainerStorageAssignmentRead(d, meta)
}

func resourceIBMContainerStorageAssignmentRead(d *schema.ResourceData, meta interface{}) error {

	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return err
	}

	uuid := d.Get("uuid").(string)
	controller := d.Get("controller").(string)
	d.Set("controller", controller)

	getAssignmentOptions := &kubernetesserviceapiv1.GetAssignmentOptions{
		UUID: &uuid,
	}

	result, _, err := satClient.GetAssignment(getAssignmentOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error getting Assignment of UUID %s - %v", uuid, err)
	}
	d.Set("assignment_name", *result.Name)
	d.Set("uuid", *result.UUID)
	d.Set("owner", *result.Owner.Name)
	if result.Groups != nil {
		d.Set("groups", result.Groups)
	}
	if result.Cluster != nil {
		d.Set("cluster", *result.Cluster)
	}
	if result.SatSvcClusterID != nil {
		d.Set("svc_cluster", *result.SatSvcClusterID)
	}
	if result.Satcluster != nil {
		d.Set("sat_cluster", *result.Satcluster)
	}
	if result.ChannelName != nil {
		d.Set("config", *result.ChannelName)
	}
	if result.ChannelUUID != nil {
		d.Set("config_uuid", *result.ChannelUUID)
	}
	if result.Version != nil {
		d.Set("config_version", *result.Version)
	}
	if result.VersionUUID != nil {
		d.Set("config_version_uuid", *result.VersionUUID)
	}
	if result.SubscriptionType != nil {
		d.Set("assignment_type", *result.SubscriptionType)
	}
	if result.Created != nil {
		d.Set("created", *result.Created)
	}
	if result.IsAssignmentUpgradeAvailable != nil {
		d.Set("is_assignment_upgrade_available", *result.IsAssignmentUpgradeAvailable)
	}
	if result.RolloutStatus != nil {
		d.Set("rollout_success_count", *result.RolloutStatus.SuccessCount)
		d.Set("rollout_error_count", *result.RolloutStatus.ErrorCount)
	}
	d.Set("update_config_revision", false)
	return nil
}

func resourceIBMContainerStorageAssignmentUpdate(d *schema.ResourceData, meta interface{}) error {
	uuid := d.Get("uuid").(string)
	updateAssignmentOptions := &kubernetesserviceapiv1.UpdateAssignmentOptions{}
	updateAssignmentOptions.UUID = &uuid

	if d.HasChange("assignment_name") || d.HasChange("groups") || d.HasChange("update_config_revision") && !d.IsNewResource() {
		assignmentName := d.Get("assignment_name").(string)
		updateAssignmentOptions.Name = &assignmentName

		groups := flex.ExpandStringList(d.Get("groups").([]interface{}))
		updateAssignmentOptions.Groups = groups

		updateConfigRevision := d.Get("update_config_revision").(bool)
		updateAssignmentOptions.UpdateConfigVersion = &updateConfigRevision

		_, err := waitForAssignmentUpdateStatus(updateAssignmentOptions, meta, d)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Updating Assignment with UUID %s - %v", uuid, err)
		}
	}
	return resourceIBMContainerStorageAssignmentRead(d, meta)
}

func resourceIBMContainerStorageAssignmentDelete(d *schema.ResourceData, meta interface{}) error {
	uuid := d.Get("uuid").(string)
	controller := d.Get("controller").(string)
	removeAssignmentOptions := &kubernetesserviceapiv1.RemoveAssignmentOptions{}
	removeAssignmentOptions.UUID = &uuid
	removeAssignmentOptions.Controller = &controller

	_, err := waitForAssignmentDeletionStatus(removeAssignmentOptions, meta, d)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Removing Assignment with UUID %s - %v", uuid, err)
	}

	d.SetId("")
	return nil
}

func waitForAssignmentCreationStatus(getAssignmentOptions *kubernetesserviceapiv1.GetAssignmentOptions, meta interface{}, d *schema.ResourceData) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:        []string{"NotReady"},
		Target:         []string{"Ready"},
		Refresh:        assignmentCreationStatusRefreshFunc(getAssignmentOptions, meta),
		Timeout:        d.Timeout(schema.TimeoutCreate),
		Delay:          10 * time.Second,
		MinTimeout:     10 * time.Second,
		NotFoundChecks: 100,
	}
	return stateConf.WaitForState()
}

func assignmentCreationStatusRefreshFunc(getAssignmentOptions *kubernetesserviceapiv1.GetAssignmentOptions, meta interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		satClient, err := meta.(conns.ClientSession).SatelliteClientSession()

		if err != nil {
			return nil, "NotReady", err
		}
		_, response, err := satClient.GetAssignment(getAssignmentOptions)

		if response.GetStatusCode() == 200 {
			return true, "Ready", nil
		}

		return nil, "NotReady", nil
	}
}

func waitForAssignmentUpdateStatus(updateAssignmentOptions *kubernetesserviceapiv1.UpdateAssignmentOptions, meta interface{}, d *schema.ResourceData) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:        []string{"NotReady"},
		Target:         []string{"Ready"},
		Refresh:        assignmentUpdateStatusRefreshFunc(updateAssignmentOptions, meta),
		Timeout:        d.Timeout(schema.TimeoutUpdate),
		Delay:          10 * time.Second,
		MinTimeout:     10 * time.Second,
		NotFoundChecks: 100,
	}
	return stateConf.WaitForState()
}

func assignmentUpdateStatusRefreshFunc(updateAssignmentOptions *kubernetesserviceapiv1.UpdateAssignmentOptions, meta interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		satClient, err := meta.(conns.ClientSession).SatelliteClientSession()

		if err != nil {
			return nil, "NotReady", err
		}

		result, _, err := satClient.UpdateAssignment(updateAssignmentOptions)

		if *result.EditSubscription.Success == true && err == nil {
			return true, "Ready", nil
		}

		return nil, "NotReady", nil
	}
}

func waitForAssignmentDeletionStatus(removeAssignmentOptions *kubernetesserviceapiv1.RemoveAssignmentOptions, meta interface{}, d *schema.ResourceData) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:        []string{"NotReady"},
		Target:         []string{"Ready"},
		Refresh:        assignmentDeletionStatusRefreshFunc(removeAssignmentOptions, meta),
		Timeout:        d.Timeout(schema.TimeoutDelete),
		Delay:          10 * time.Second,
		MinTimeout:     10 * time.Second,
		NotFoundChecks: 100,
	}
	return stateConf.WaitForState()
}

func assignmentDeletionStatusRefreshFunc(removeAssignmentOptions *kubernetesserviceapiv1.RemoveAssignmentOptions, meta interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		satClient, err := meta.(conns.ClientSession).SatelliteClientSession()

		if err != nil {
			return nil, "NotReady", err
		}

		response, _, err := satClient.RemoveAssignment(removeAssignmentOptions)
		if *response.RemoveSubscription.Success == true && err == nil {
			return true, "Ready", nil
		}

		return nil, "NotReady", nil
	}
}
