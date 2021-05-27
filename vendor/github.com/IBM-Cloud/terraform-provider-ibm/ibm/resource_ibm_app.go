// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"time"

	v2 "github.com/IBM-Cloud/bluemix-go/api/mccp/mccpv2"
	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/bluemix-go/helpers"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	homedir "github.com/mitchellh/go-homedir"
)

func resourceIBMApp() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMAppCreate,
		Read:     resourceIBMAppRead,
		Update:   resourceIBMAppUpdate,
		Delete:   resourceIBMAppDelete,
		Exists:   resourceIBMAppExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name for the app",
			},
			"memory": {
				Description: "The amount of memory each instance should have. In megabytes.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"instances": {
				Description: "The number of instances",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
			},
			"disk_quota": {
				Description: "The maximum amount of disk available to an instance of an app. In megabytes.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"space_guid": {
				Description: "Define space guid to which app belongs",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"buildpack": {
				Description: "Buildpack to build the app. 3 options: a) Blank means autodetection; b) A Git Url pointing to a buildpack; c) Name of an installed buildpack.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"environment_json": {
				Description: "Key/value pairs of all the environment variables to run in your app. Does not include any system or service variables.",
				Type:        schema.TypeMap,
				Optional:    true,
			},
			"route_guid": {
				Description: "Define the route guids which should be bound to the application.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
			},
			"service_instance_guid": {
				Description: "Define the service instance guids that should be bound to this application.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
			},
			"app_path": {
				Description: "Define the  path of the zip file of the application.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"app_version": {
				Description: "Version of the application",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"command": {
				Description: "The initial command for the app",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"wait_time_minutes": {
				Description: "Define timeout to wait for the app instances to start/update/restage etc.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     20,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"health_check_http_endpoint": {
				Description: "Endpoint called to determine if the app is healthy.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"health_check_type": {
				Description:  "Type of health check to perform.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "port",
				ValidateFunc: validateAllowedStringValue([]string{"port", "process"}),
			},
			"health_check_timeout": {
				Description: "Timeout in seconds for health checking of an staged app when starting up.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
		},
	}
}

func resourceIBMAppCreate(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(ClientSession).MccpAPI()
	if err != nil {
		return err
	}
	appAPI := cfClient.Apps()
	name := d.Get("name").(string)
	spaceGUID := d.Get("space_guid").(string)
	healthChekcType := d.Get("health_check_type").(string)

	appCreatePayload := v2.AppRequest{
		Name:            helpers.String(name),
		SpaceGUID:       helpers.String(spaceGUID),
		HealthCheckType: helpers.String(healthChekcType),
	}

	if memory, ok := d.GetOk("memory"); ok {
		appCreatePayload.Memory = memory.(int)
	}

	if instances, ok := d.GetOk("instances"); ok {
		appCreatePayload.Instances = instances.(int)
	}

	if diskQuota, ok := d.GetOk("disk_quota"); ok {
		appCreatePayload.DiskQuota = diskQuota.(int)
	}

	if buildpack, ok := d.GetOk("buildpack"); ok {
		appCreatePayload.BuildPack = helpers.String(buildpack.(string))
	}

	if environmentJSON, ok := d.GetOk("environment_json"); ok {
		appCreatePayload.EnvironmentJSON = helpers.Map(environmentJSON.(map[string]interface{}))

	}

	if command, ok := d.GetOk("command"); ok {
		appCreatePayload.Command = helpers.String(command.(string))
	}

	if healtChkEndpoint, ok := d.GetOk("health_check_http_endpoint"); ok {
		appCreatePayload.HealthCheckHTTPEndpoint = helpers.String(healtChkEndpoint.(string))
	}

	if healtChkTimeout, ok := d.GetOk("health_check_timeout"); ok {
		appCreatePayload.HealthCheckTimeout = healtChkTimeout.(int)
	}

	_, err = appAPI.FindByName(spaceGUID, name)
	if err == nil {
		return fmt.Errorf("%s already exists in the given space %s", name, spaceGUID)
	}

	log.Println("[INFO] Creating Cloud Foundary Application")
	app, err := appAPI.Create(appCreatePayload)
	if err != nil {
		return fmt.Errorf("Error creating app: %s", err)
	}

	appGUID := app.Metadata.GUID
	log.Println("[INFO] Cloud Foundary Application is created successfully")

	d.SetId(appGUID)

	if v, ok := d.Get("route_guid").(*schema.Set); ok && v.Len() > 0 {
		log.Println("[INFO] Bind the route with cloud foundary application")
		for _, routeID := range v.List() {
			_, err := appAPI.BindRoute(appGUID, routeID.(string))
			if err != nil {
				return fmt.Errorf("Error binding route %s to app: %s", routeID.(string), err)
			}
		}
	}
	if v, ok := d.Get("service_instance_guid").(*schema.Set); ok && v.Len() > 0 {
		sbAPI := cfClient.ServiceBindings()
		for _, svcID := range v.List() {
			req := v2.ServiceBindingRequest{
				ServiceInstanceGUID: svcID.(string),
				AppGUID:             appGUID,
			}
			_, err := sbAPI.Create(req)
			if err != nil {
				return fmt.Errorf("Error binding service instance %s to  app: %s", svcID.(string), err)
			}
		}
	}
	log.Println("[INFO] Upload the app bits to the cloud foundary application")
	applicationZip, err := processAppZipPath(d.Get("app_path").(string))
	if err != nil {
		return err
	}

	_, err = appAPI.Upload(appGUID, applicationZip)
	if err != nil {
		return fmt.Errorf("Error uploading app bits: %s", err)
	}

	err = restartApp(appGUID, d, meta)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Application: %s has started successfully", name)
	return resourceIBMAppRead(d, meta)
}

func resourceIBMAppRead(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(ClientSession).MccpAPI()
	if err != nil {
		return err
	}
	appAPI := cfClient.Apps()
	appGUID := d.Id()

	appData, err := appAPI.Get(appGUID)
	if err != nil {
		return fmt.Errorf("Error retrieving app details %s : %s", appGUID, err)
	}

	d.SetId(appData.Metadata.GUID)
	d.Set("name", appData.Entity.Name)
	d.Set("memory", appData.Entity.Memory)
	d.Set("instances", appData.Entity.Instances)
	d.Set("space_guid", appData.Entity.SpaceGUID)
	d.Set("disk_quota", appData.Entity.DiskQuota)
	d.Set("buildpack", appData.Entity.BuildPack)
	d.Set("environment_json", Flatten(appData.Entity.EnvironmentJSON))
	d.Set("command", appData.Entity.Command)
	d.Set("health_check_type", appData.Entity.HealthCheckType)
	d.Set("health_check_http_endpoint", appData.Entity.HealthCheckHTTPEndpoint)
	d.Set("health_check_timeout", appData.Entity.HealthCheckTimeout)

	route, err := appAPI.ListRoutes(appGUID)
	if err != nil {
		return err
	}
	if len(route) > 0 {
		d.Set("route_guid", flattenRoute(route))
	}

	svcBindings, err := appAPI.ListServiceBindings(appGUID)
	if err != nil {
		return err
	}
	if len(svcBindings) > 0 {
		d.Set("service_instance_guid", flattenServiceBindings(svcBindings))
	}

	return nil

}

func resourceIBMAppUpdate(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(ClientSession).MccpAPI()
	if err != nil {
		return err
	}
	appAPI := cfClient.Apps()
	appGUID := d.Id()

	appUpdatePayload := v2.AppRequest{}
	restartRequired := false
	restageRequired := false

	waitTimeout := time.Duration(d.Get("wait_time_minutes").(int)) * time.Minute

	if d.HasChange("name") {
		appUpdatePayload.Name = helpers.String(d.Get("name").(string))
	}

	if d.HasChange("memory") {
		appUpdatePayload.Memory = d.Get("memory").(int)
	}

	if d.HasChange("instances") {
		appUpdatePayload.Instances = d.Get("instances").(int)
	}

	if d.HasChange("disk_quota") {
		appUpdatePayload.DiskQuota = d.Get("disk_quota").(int)
	}

	if d.HasChange("buildpack") {
		appUpdatePayload.BuildPack = helpers.String(d.Get("buildpack").(string))
		restageRequired = true
	}

	if d.HasChange("command") {
		appUpdatePayload.Command = helpers.String(d.Get("command").(string))
		restartRequired = true
	}

	if d.HasChange("environment_json") {
		appUpdatePayload.EnvironmentJSON = helpers.Map(d.Get("environment_json").(map[string]interface{}))
		restageRequired = true
	}

	if d.HasChange("health_check_type") {
		appUpdatePayload.HealthCheckType = helpers.String(d.Get("health_check_type").(string))
		restartRequired = true
	}

	if d.HasChange("health_check_http_endpoint") {
		appUpdatePayload.HealthCheckHTTPEndpoint = helpers.String(d.Get("health_check_http_endpoint").(string))
		restartRequired = true
	}

	if d.HasChange("health_check_timeout") {
		appUpdatePayload.HealthCheckTimeout = d.Get("health_check_timeout").(int)
		restartRequired = true
	}

	if d.HasChange("command") {
		appUpdatePayload.Command = helpers.String(d.Get("command").(string))
		restartRequired = true
	}
	log.Println("[INFO] Update cloud foundary application")

	_, err = appAPI.Update(appGUID, appUpdatePayload)
	if err != nil {
		return fmt.Errorf("Error updating application: %s", err)
	}
	//TODO find the digest of the zip and avoid upload if it is same
	if d.HasChange("app_path") || d.HasChange("app_version") {
		appZipLoc, err := processAppZipPath(d.Get("app_path").(string))
		if err != nil {
			return err
		}
		log.Println("[DEBUG] Uploading application bits")
		_, err = appAPI.Upload(appGUID, appZipLoc)
		if err != nil {
			return fmt.Errorf("Error uploading  app: %s", err)
		}
		restartRequired = true
	}

	err = updateRouteGUID(appGUID, appAPI, d)
	if err != nil {
		return err
	}

	restage, err := updateServiceInstanceGUID(appGUID, d, meta)
	if err != nil {
		return err
	}
	if restage {
		restageRequired = true
	}

	/*Wait if any previous staging is going on
	log.Println("[INFO] Waiting to see any previous staging is on or not")
	state, err := appAPI.WaitForAppStatus(v2.AppStagedState, appGUID, waitTimeout)
	if waitTimeout != 0 && (err != nil || state == v2.AppPendingState) {
		return fmt.Errorf("The application is still in %s from last operations.Please try again after sometime by increasing timeout value %q", state, err)
	}*/

	//If restage and restart both are required then we only need restage as that starts over everything
	if restageRequired {
		log.Println("[INFO] Restage since buildpack has changed")
		err := restageApp(appGUID, d, meta)
		if err != nil {
			return err
		}
	} else if restartRequired {
		err := restartApp(appGUID, d, meta)
		if err != nil {
			return err
		}
	} else {
		//In case only memory/disk etc are updated then cloud controller would destroy the current instances
		//and spin new ones, so we are waiting till they come up again
		state, err := appAPI.WaitForInstanceStatus(v2.AppRunningState, appGUID, waitTimeout)
		if waitTimeout != 0 && (err != nil || state != v2.AppRunningState) {
			return fmt.Errorf("All applications instances aren't %s, Current status is %s, %q", v2.AppRunningState, state, err)
		}
	}

	return resourceIBMAppRead(d, meta)
}

func resourceIBMAppDelete(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(ClientSession).MccpAPI()
	if err != nil {
		return err
	}
	appAPI := cfClient.Apps()
	id := d.Id()

	err = appAPI.Delete(id, false, true)
	if err != nil {
		return fmt.Errorf("Error deleting app: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceIBMAppExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	cfClient, err := meta.(ClientSession).MccpAPI()
	if err != nil {
		return false, err
	}
	appAPI := cfClient.Apps()
	id := d.Id()

	app, err := appAPI.Get(id)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}

	return app.Metadata.GUID == id, nil
}

func updateRouteGUID(appGUID string, appAPI v2.Apps, d *schema.ResourceData) (err error) {
	if d.HasChange("route_guid") {
		ors, nrs := d.GetChange("route_guid")
		or := ors.(*schema.Set)
		nr := nrs.(*schema.Set)

		remove := expandStringList(or.Difference(nr).List())
		add := expandStringList(nr.Difference(or).List())

		if len(add) > 0 {
			for i := range add {
				_, err = appAPI.BindRoute(appGUID, add[i])
				if err != nil {
					return fmt.Errorf("Error while binding route %q to application %s: %q", add[i], appGUID, err)
				}
			}
		}
		if len(remove) > 0 {
			for i := range remove {
				err = appAPI.UnBindRoute(appGUID, remove[i])
				if err != nil {
					return fmt.Errorf("Error while un-binding route %q from application %s: %q", add[i], appGUID, err)
				}
			}
		}
	}
	return
}

func updateServiceInstanceGUID(appGUID string, d *schema.ResourceData, meta interface{}) (restageRequired bool, err error) {
	cfClient, err := meta.(ClientSession).MccpAPI()
	if err != nil {
		return false, err
	}
	appAPI := cfClient.Apps()
	sbAPI := cfClient.ServiceBindings()
	if d.HasChange("service_instance_guid") {
		oss, nss := d.GetChange("service_instance_guid")
		os := oss.(*schema.Set)
		ns := nss.(*schema.Set)
		remove := expandStringList(os.Difference(ns).List())
		add := expandStringList(ns.Difference(os).List())

		if len(add) > 0 {
			for i := range add {
				sbPayload := v2.ServiceBindingRequest{
					ServiceInstanceGUID: add[i],
					AppGUID:             appGUID,
				}
				_, err = sbAPI.Create(sbPayload)
				if err != nil {
					err = fmt.Errorf("Error while binding service instance %s to application %s: %q", add[i], appGUID, err)
					return
				}
				restageRequired = true
			}
		}
		if len(remove) > 0 {
			var appFilters, svcFilters string
			var bindings []v2.ServiceBinding
			appFilters, err = new(v2.Filter).Name("app_guid").Eq(appGUID).Build()
			if err != nil {
				return
			}
			svcFilters, err = new(v2.Filter).Name("service_instance_guid").In(remove...).Build()
			if err != nil {
				return
			}
			bindings, err = sbAPI.List(appFilters, svcFilters)
			if err != nil {
				return
			}
			sbIds := make([]string, len(bindings))
			for i, sb := range bindings {
				sbIds[i] = sb.GUID
			}
			err = appAPI.DeleteServiceBindings(appGUID, sbIds...)
			if err != nil {
				err = fmt.Errorf("Error while un-binding service instances %s to application %s: %q", remove, appGUID, err)
				return
			}
		}
	}
	return
}
func restartApp(appGUID string, d *schema.ResourceData, meta interface{}) error {
	cfClient, _ := meta.(ClientSession).MccpAPI()
	appAPI := cfClient.Apps()

	appUpdatePayload := v2.AppRequest{
		State: helpers.String(v2.AppStoppedState),
	}
	log.Println("[INFO] Stopping Application")
	_, err := appAPI.Update(appGUID, appUpdatePayload)
	if err != nil {
		return fmt.Errorf("Error updating application status to %s %s", v2.AppStoppedState, err)
	}
	waitTimeout := time.Duration(d.Get("wait_time_minutes").(int)) * time.Minute
	log.Println("[INFO] Starting Application")
	status, err := appAPI.Start(appGUID, waitTimeout)
	if err != nil {
		return fmt.Errorf("Error while starting application : %s", err)
	}
	if waitTimeout != 0 {
		return checkAppStatus(status)
	}
	return nil
}

func restageApp(appGUID string, d *schema.ResourceData, meta interface{}) error {
	cfClient, _ := meta.(ClientSession).MccpAPI()
	appAPI := cfClient.Apps()

	log.Println("[INFO] Restage Application")
	waitTimeout := time.Duration(d.Get("wait_time_minutes").(int)) * time.Minute
	status, err := appAPI.Restage(appGUID, waitTimeout)
	if err != nil {
		return fmt.Errorf("Error while restaging application : %s", err)
	}
	if waitTimeout != 0 {
		return checkAppStatus(status)
	}
	return nil
}

func checkAppStatus(status *v2.AppState) error {
	if status.PackageState != v2.AppStagedState {
		return fmt.Errorf("Applications couldn't be staged, current status is  %s", status.PackageState)
	}
	if status.InstanceState != v2.AppRunningState {
		return fmt.Errorf("All applications instances aren't %s, Current status is %s", v2.AppRunningState, status.InstanceState)
	}
	return nil
}

func processAppZipPath(path string) (string, error) {
	applicationZip, err := homedir.Expand(path)
	if err != nil {
		return path, fmt.Errorf("home directory in the given path %s couldn't be expanded", path)
	}
	if !helpers.FileExists(applicationZip) {
		return path, fmt.Errorf("The given app path: %s doesn't exist", path)
	}
	return applicationZip, nil
}
