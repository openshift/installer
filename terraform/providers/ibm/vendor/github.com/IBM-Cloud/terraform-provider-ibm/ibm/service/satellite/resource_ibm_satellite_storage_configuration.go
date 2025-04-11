package satellite

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/IBM-Cloud/container-services-go-sdk/kubernetesserviceapiv1"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"k8s.io/utils/strings/slices"
)

func ResourceIBMSatelliteStorageConfiguration() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMContainerStorageConfigurationCreate,
		Read:     resourceIBMContainerStorageConfigurationRead,
		Update:   resourceIBMContainerStorageConfigurationUpdate,
		Delete:   resourceIBMContainerStorageConfigurationDelete,
		Importer: &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Location ID.",
			},
			"config_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the Storage Configuration.",
			},
			"config_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version of the Storage Configuration.",
			},
			"storage_template_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Storage Template Name.",
			},
			"storage_template_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Storage Template Version.",
			},
			"user_config_parameters": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "User Config Parameters to pass as a Map of string key-value.",
			},
			"user_secret_parameters": {
				Type:      schema.TypeMap,
				Sensitive: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "User Secret Parameters to pass as a Map of string key-value.",
			},
			"storage_class_parameters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Description: "The List of Storage Class Parameters as a list of a  Map of string key-value.",
				},
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: "The Universally Unique IDentifier (UUID) of the Storage Configuration.",
			},
			"update_assignments": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Set to update all assignments during a configuration update.",
			},
			"delete_assignments": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Set to delete all assignments during a configuration destroy.",
			},
		},
	}
}

// Helper Function to convert map[string]interface to map[string]string
func convertToMapStringString(mapInterface map[string]interface{}) map[string]string {
	data := make(map[string]string)
	for k, v := range mapInterface {
		data[k] = v.(string)
	}
	return data
}

// Function to validate the keys of user_config_parameters and user_secrets_parameters
func validateStorageConfig(d *schema.ResourceData, meta interface{}) error {
	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return err
	}

	userconfigParams := convertToMapStringString(d.Get("user_config_parameters").(map[string]interface{}))
	usersecretParams := convertToMapStringString(d.Get("user_secret_parameters").(map[string]interface{}))
	storageTemplateName := d.Get("storage_template_name").(string)
	storageTemplateVersion := d.Get("storage_template_version").(string)
	storageresult := &kubernetesserviceapiv1.GetStorageTemplateOptions{
		Name:    &storageTemplateName,
		Version: &storageTemplateVersion,
	}
	// We get the details of the storage template i.e we get the parameter list for that specific template.
	result, _, err := satClient.GetStorageTemplate(storageresult)
	if err != nil {
		return err
	}

	var customparamList []string
	for _, v := range result.CustomParameters {
		var parameterOptions map[string]string
		inrec, _ := json.Marshal(v)
		json.Unmarshal(inrec, &parameterOptions)
		if parameterOptions["required"] == "true" {
			_, foundConfig := userconfigParams[parameterOptions["name"]]
			_, foundSecret := usersecretParams[parameterOptions["name"]]
			// if "required" key parameters are not present in the terraform schema
			if !(foundConfig || foundSecret) {
				// if they have default values, set it with those or else throw an error
				if len(parameterOptions["default"]) > 0 {
					userconfigParams[parameterOptions["name"]] = parameterOptions["default"]
				} else {
					return fmt.Errorf("%s Parameter missing - Required", parameterOptions["name"])
				}
			}
		}
		customparamList = append(customparamList, parameterOptions["name"])
	}
	// checks if the user has entered correct parameteric keys, if the key is not found an error is thrown
	for k, _ := range userconfigParams {
		if !slices.Contains(customparamList, k) {
			return fmt.Errorf("Config Parameter %s not found", k)
		}
	}

	// checks if the user has entered correct secret keys, if the key is not found an error is thrown
	for k, _ := range usersecretParams {
		if !slices.Contains(customparamList, k) {
			return fmt.Errorf("Secret Parameter %s not found", k)
		}
	}

	return nil
}

func resourceIBMContainerStorageConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return err
	}
	createStorageConfigurationOptions := &kubernetesserviceapiv1.CreateStorageConfigurationOptions{}
	satLocation := d.Get("location").(string)
	createStorageConfigurationOptions.Controller = &satLocation

	err = validateStorageConfig(d, meta)
	if err != nil {
		return err
	}

	var configName string
	if v, ok := d.GetOk("config_name"); ok {
		configName = v.(string)
		createStorageConfigurationOptions.SetConfigName(v.(string))
	}

	if v, ok := d.GetOk("storage_template_name"); ok {
		createStorageConfigurationOptions.SetStorageTemplateName(v.(string))
	}

	if v, ok := d.GetOk("storage_template_version"); ok {
		createStorageConfigurationOptions.SetStorageTemplateVersion(v.(string))
	}

	if v, ok := d.GetOk("user_config_parameters"); ok {
		userConfigParameters := convertToMapStringString(v.(map[string]interface{}))
		createStorageConfigurationOptions.SetUserConfigParameters(userConfigParameters)
	}

	if v, ok := d.GetOk("user_secret_parameters"); ok {
		userSecretParams := convertToMapStringString(v.(map[string]interface{}))
		createStorageConfigurationOptions.SetUserSecretParameters(userSecretParams)
	}
	// convert the storage class parameters into a list[map[string]string]
	if storageClassParamsList, ok := d.GetOk("storage_class_parameters"); ok {
		var scpList []map[string]string
		for _, value := range storageClassParamsList.([]interface{}) {
			storageclassParams := convertToMapStringString(value.(map[string]interface{}))
			scpList = append(scpList, storageclassParams)
		}
		createStorageConfigurationOptions.SetStorageClassParameters(scpList)
	}

	result, _, err := satClient.CreateStorageConfiguration(createStorageConfigurationOptions)
	if err != nil {
		return fmt.Errorf("Unable to Create Storage Configuration - %v", err)
	}
	getStorageConfigurationOptions := &kubernetesserviceapiv1.GetStorageConfigurationOptions{
		Name: createStorageConfigurationOptions.ConfigName,
	}
	// If we are able to successful get the configuration, then create is assumed to be a success
	_, err = waitForStorageConfigurationStatus(getStorageConfigurationOptions, meta, d)
	if err != nil {
		return err
	}

	d.SetId(*result.AddChannel.UUID + "/" + configName)
	return resourceIBMContainerStorageConfigurationRead(d, meta)
}

func resourceIBMContainerStorageConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return err
	}

	satLocation := d.Get("location").(string)
	d.Set("location", satLocation)
	// We will only set the user defined storage classes in the terraform state file
	var scDefinedList []map[string]string
	if storageClassParamsList, ok := d.GetOk("storage_class_parameters"); ok {
		for _, value := range storageClassParamsList.([]interface{}) {
			storageclassParams := convertToMapStringString(value.(map[string]interface{}))
			scDefinedList = append(scDefinedList, storageclassParams)
		}
	}

	storageConfigName := d.Get("config_name").(string)
	getStorageConfigurationOptions := &kubernetesserviceapiv1.GetStorageConfigurationOptions{
		Name: &storageConfigName,
	}

	result, _, err := satClient.GetStorageConfiguration(getStorageConfigurationOptions)
	if err != nil {
		return err
	}

	d.Set("config_name", *result.ConfigName)
	d.Set("config_version", *result.ConfigVersion)
	d.Set("storage_template_name", *result.StorageTemplateName)
	d.Set("storage_template_version", *result.StorageTemplateVersion)
	d.Set("user_config_parameters", result.UserConfigParameters)
	// The secret parameters from terraform are set directly in the state file, they cannot be retreived from the server. A local copy is kept for secret parameter refresh.
	userSecretParams := convertToMapStringString(d.Get("user_secret_parameters").(map[string]interface{}))
	for k, _ := range result.UserSecretParameters {
		result.UserSecretParameters[k] = userSecretParams[k]
	}
	d.Set("user_secret_parameters", result.UserSecretParameters)
	var storageClassList []map[string]string
	for _, v := range result.StorageClassParameters {
		if getDefinedStorageClasses(scDefinedList, v) {
			storageClassList = append(storageClassList, v)
		}
	}
	// Set terraform defined storage classes
	d.Set("storage_class_parameters", storageClassList)
	d.Set("uuid", *result.UUID)
	d.Set("update_assignments", false)
	delete_assignments := d.Get("delete_assignments").(bool)
	d.Set("delete_assignments", delete_assignments)

	return nil
}

func resourceIBMContainerStorageConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return err
	}
	updateStorageConfigurationOptions := &kubernetesserviceapiv1.UpdateStorageConfigurationOptions{}
	updateAssignments := d.Get("update_assignments").(bool)
	configName := d.Get("config_name").(string)
	satLocation := d.Get("location").(string)
	updateStorageConfigurationOptions.Controller = &satLocation

	err = validateStorageConfig(d, meta)
	if err != nil {
		return err
	}

	if d.HasChange("user_config_parameters") || d.HasChange("user_secret_parameters") || d.HasChange("storage_class_parameters") && !d.IsNewResource() {

		if v, ok := d.GetOk("config_name"); ok {
			updateStorageConfigurationOptions.SetConfigName(v.(string))
		}

		if v, ok := d.GetOk("storage_template_name"); ok {
			updateStorageConfigurationOptions.SetStorageTemplateName(v.(string))
		}

		if v, ok := d.GetOk("storage_template_version"); ok {
			updateStorageConfigurationOptions.SetStorageTemplateVersion(v.(string))
		}

		if v, ok := d.GetOk("user_config_parameters"); ok {
			userConfigParameters := convertToMapStringString(v.(map[string]interface{}))
			updateStorageConfigurationOptions.SetUserConfigParameters(userConfigParameters)
		}

		if v, ok := d.GetOk("user_secret_parameters"); ok {
			userSecretParams := convertToMapStringString(v.(map[string]interface{}))
			updateStorageConfigurationOptions.SetUserSecretParameters(userSecretParams)
		}

		if storageClassParamsList, ok := d.GetOk("storage_class_parameters"); ok {
			var scpList []map[string]string
			for _, value := range storageClassParamsList.([]interface{}) {
				storageclassParams := convertToMapStringString(value.(map[string]interface{}))
				scpList = append(scpList, storageclassParams)
			}
			updateStorageConfigurationOptions.SetStorageClassParameters(scpList)
		}

		_, _, err := satClient.UpdateStorageConfiguration(updateStorageConfigurationOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Unable to Update Storage Configuration %s - %v", *updateStorageConfigurationOptions.ConfigName, err)
		}

		getStorageConfigurationOptions := &kubernetesserviceapiv1.GetStorageConfigurationOptions{
			Name: updateStorageConfigurationOptions.ConfigName,
		}
		// If we are able to successful get the configuration, then update is assumed to be a success
		_, err = waitForStorageConfigurationStatus(getStorageConfigurationOptions, meta, d)
		if err != nil {
			return err
		}

		// if the user has set the updateAssignments parameter to true, then all the assignments are auto updated with the latest configuration revision
		if updateAssignments {
			getAssignmentsByConfigOptions := &kubernetesserviceapiv1.GetAssignmentsByConfigOptions{
				Config: updateStorageConfigurationOptions.ConfigName,
			}

			result, _, err := satClient.GetAssignmentsByConfig(getAssignmentsByConfigOptions)
			if err != nil {
				return err
			}

			for _, v := range result {
				updateStorageAssignment := &kubernetesserviceapiv1.UpdateAssignmentOptions{
					Name:                v.Name,
					UpdateConfigVersion: &updateAssignments,
					UUID:                v.UUID,
				}
				_, _, err := satClient.UpdateAssignment(updateStorageAssignment)
				if err != nil {
					return err
				}
			}
		}
	}

	// Version Upgrade Scenario, the existing configuration is deleted along with its assignments,
	// a new configuration with the new template version is created and assigned back to the previously assigned clusters and groups
	if d.HasChange("storage_template_version") {
		getAssignmentsByConfigOptions := &kubernetesserviceapiv1.GetAssignmentsByConfigOptions{
			Config: &configName,
		}
		result, _, err := satClient.GetAssignmentsByConfig(getAssignmentsByConfigOptions)
		if err != nil {
			return err
		}
		err = resourceIBMContainerStorageConfigurationDelete(d, meta)
		if err != nil {
			return err
		}
		err = resourceIBMContainerStorageConfigurationCreate(d, meta)
		if err != nil {
			return err
		}
		for _, v := range result {
			if len(v.Groups) != 0 {
				createAssignmentGroupOptions := &kubernetesserviceapiv1.CreateAssignmentOptions{
					Name:       v.Name,
					Groups:     v.Groups,
					Config:     v.ChannelName,
					Controller: &satLocation,
				}
				_, _, err = satClient.CreateAssignment(createAssignmentGroupOptions)
				if err != nil {
					return fmt.Errorf("[ERROR] Creating Assignment during Storage Configuration Upgrade - %v", err)
				}
			} else {
				createAssignmentConfigOptions := &kubernetesserviceapiv1.CreateAssignmentOptions{
					Name:       v.Name,
					Config:     v.ChannelName,
					Cluster:    v.Cluster,
					Controller: &satLocation,
				}
				_, _, err = satClient.CreateAssignmentByCluster(createAssignmentConfigOptions)
				if err != nil {
					return fmt.Errorf("[ERROR] Creating Assignment during Storage Configuration Upgrade - %v", err)
				}
			}
		}
	}

	return resourceIBMContainerStorageConfigurationRead(d, meta)
}

func resourceIBMContainerStorageConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	uuid := d.Get("uuid").(string)
	name := d.Get("config_name").(string)
	controller := d.Get("location").(string)
	delete_assignments := d.Get("delete_assignments").(bool)
	removeStorageConfigurationOptions := &kubernetesserviceapiv1.RemoveStorageConfigurationOptions{}
	removeStorageConfigurationOptions.UUID = &uuid
	removeStorageConfigurationOptions.Controller = &controller
	if delete_assignments || d.HasChange("storage_template_version") {
		removeStorageConfigurationOptions.RemoveAssignments = &delete_assignments
	}

	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return err
	}
	_, _, err = satClient.RemoveStorageConfiguration(removeStorageConfigurationOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Deleting Storage Configuration %s - %v", name, err)
	}
	getStorageConfigurationOptions := &kubernetesserviceapiv1.GetStorageConfigurationOptions{
		Name: &name,
	}

	// If we cannot Get the storage configuration, it is assumed to be successfully deleted.
	_, err = waitForStorageConfigurationDeletionStatus(getStorageConfigurationOptions, meta, d)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

// Helper function to extract only the terraform defined storage classes from the server equivalent
func getDefinedStorageClasses(definedMaps []map[string]string, getMaps map[string]string) bool {
	for _, v := range definedMaps {
		eq := reflect.DeepEqual(v, getMaps)
		if eq {
			return true
		}
	}
	return false
}

func waitForStorageConfigurationStatus(getStorageConfigurationOptions *kubernetesserviceapiv1.GetStorageConfigurationOptions, meta interface{}, d *schema.ResourceData) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:        []string{"NotReady"},
		Target:         []string{"Ready"},
		Refresh:        storageConfigurationStatusRefreshFunc(getStorageConfigurationOptions, meta),
		Timeout:        d.Timeout(schema.TimeoutCreate),
		Delay:          10 * time.Second,
		MinTimeout:     10 * time.Second,
		NotFoundChecks: 100,
	}
	return stateConf.WaitForState()
}

func storageConfigurationStatusRefreshFunc(getStorageConfigurationOptions *kubernetesserviceapiv1.GetStorageConfigurationOptions, meta interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		satClient, err := meta.(conns.ClientSession).SatelliteClientSession()

		if err != nil {
			return nil, "NotReady", err
		}
		_, response, err := satClient.GetStorageConfiguration(getStorageConfigurationOptions)

		if response.GetStatusCode() == 200 {
			return true, "Ready", nil
		}

		return nil, "NotReady", nil
	}
}

func waitForStorageConfigurationDeletionStatus(getStorageConfigurationOptions *kubernetesserviceapiv1.GetStorageConfigurationOptions, meta interface{}, d *schema.ResourceData) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:        []string{"NotReady"},
		Target:         []string{"Ready"},
		Refresh:        storageConfigurationDeletionStatusRefreshFunc(getStorageConfigurationOptions, meta),
		Timeout:        d.Timeout(schema.TimeoutDelete),
		Delay:          10 * time.Second,
		MinTimeout:     10 * time.Second,
		NotFoundChecks: 100,
	}
	return stateConf.WaitForState()
}

func storageConfigurationDeletionStatusRefreshFunc(getStorageConfigurationOptions *kubernetesserviceapiv1.GetStorageConfigurationOptions, meta interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		satClient, err := meta.(conns.ClientSession).SatelliteClientSession()

		if err != nil {
			return nil, "NotReady", err
		}
		_, response, err := satClient.GetStorageConfiguration(getStorageConfigurationOptions)
		if response.GetStatusCode() == 404 {
			return true, "Ready", nil
		}

		return nil, "NotReady", nil
	}
}
