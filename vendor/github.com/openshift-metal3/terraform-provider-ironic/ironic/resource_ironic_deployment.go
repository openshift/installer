package ironic

import (
	"fmt"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/nodes"
	utils "github.com/gophercloud/utils/openstack/baremetal/v1/nodes"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Schema resource definition for an Ironic deployment.
func resourceDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceDeploymentCreate,
		Read:   resourceDeploymentRead,
		Delete: resourceDeploymentDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"node_uuid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_info": {
				Type:     schema.TypeMap,
				Required: true,
				ForceNew: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"network_data": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"provision_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_error": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// Create an deployment, including driving Ironic's state machine
func resourceDeploymentCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Clients).GetIronicClient()
	if err != nil {
		return err
	}

	// Reload the resource before returning
	defer resourceDeploymentRead(d, meta)

	// Set instance info
	instanceInfo := d.Get("instance_info").(map[string]interface{})
	if instanceInfo != nil {
		_, err := UpdateNode(client, d.Get("node_uuid").(string), nodes.UpdateOpts{
			nodes.UpdateOperation{
				Op:    nodes.AddOp,
				Path:  "/instance_info",
				Value: instanceInfo,
			},
		})
		if err != nil {
			return fmt.Errorf("could not update instance info: %s", err)
		}
	}

	d.SetId(d.Get("node_uuid").(string))

	configDrive, err := buildConfigDrive(client.Microversion,
		d.Get("user_data").(string),
		d.Get("network_data").(map[string]interface{}),
		d.Get("metadata").(map[string]interface{}))
	if err != nil {
		return err
	}

	// Deploy the node - drive Ironic state machine until node is 'active'
	return ChangeProvisionStateToTarget(client, d.Id(), "active", &configDrive)
}

// buildConfigDrive handles building a config drive appropriate for the Ironic version we are using.  Newer versions
// support sending the user data directly, otherwise we need to build an ISO image
func buildConfigDrive(apiVersion, userData string, networkData, metaData map[string]interface{}) (configDrive interface{}, err error) {
	actual, err := version.NewVersion(apiVersion)
	minimum, err := version.NewVersion("1.56")

	if minimum.GreaterThan(actual) {
		// Create config drive ISO directly with gophercloud/utils
		configDriveData := utils.ConfigDrive{
			UserData:    utils.UserDataString(userData),
			NetworkData: networkData,
			MetaData:    metaData,
		}
		configDriveISO, err := configDriveData.ToConfigDrive()
		if err != nil {
			return nil, err
		}
		configDrive = &configDriveISO
	} else {
		// Let Ironic handle creating the config drive
		configDrive = &nodes.ConfigDrive{
			UserData:    userData,
			NetworkData: networkData,
			MetaData:    metaData,
		}
	}

	return
}

// Read the deployment's data from Ironic
func resourceDeploymentRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Clients).GetIronicClient()
	if err != nil {
		return err
	}

	// Ensure node exists first
	id := d.Get("node_uuid").(string)
	result, err := nodes.Get(client, id).Extract()
	if err != nil {
		return fmt.Errorf("could not find node %s: %s", id, err)
	}

	d.Set("provision_state", result.ProvisionState)
	d.Set("last_error", result.LastError)

	return nil
}

// Delete an deployment from Ironic - this cleans the node and returns it's state to 'available'
func resourceDeploymentDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Clients).GetIronicClient()
	if err != nil {
		return err
	}

	return ChangeProvisionStateToTarget(client, d.Id(), "deleted", nil)
}
