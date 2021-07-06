// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
)

const HEALTH_CHECK_TYPE_HTTP_CUSTOM = "HTTP-CUSTOM"

var IBMComputeAutoScaleGroupObjectMask = []string{
	"id",
	"name",
	"minimumMemberCount",
	"maximumMemberCount",
	"cooldown",
	"status[keyName]",
	"regionalGroup[id,name]",
	"terminationPolicy[keyName]",
	"virtualGuestMemberTemplate[blockDeviceTemplateGroup,primaryNetworkComponent[networkVlan[id]],primaryBackendNetworkComponent[networkVlan[id]]]",
	"loadBalancers[id,port,virtualServerId,healthCheck[id]]",
	"networkVlans[id,networkVlanId,networkVlan[vlanNumber,primaryRouter[hostname]]]",
	"loadBalancers[healthCheck[healthCheckTypeId,type[keyname],attributes[value,type[id,keyname]]]]",
}

func resourceIBMComputeAutoScaleGroup() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMComputeAutoScaleGroupCreate,
		Read:     resourceIBMComputeAutoScaleGroupRead,
		Update:   resourceIBMComputeAutoScaleGroupUpdate,
		Delete:   resourceIBMComputeAutoScaleGroupDelete,
		Exists:   resourceIBMComputeAutoScaleGroupExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name",
			},

			"regional_group": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "regional group",
			},

			"minimum_member_count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Minimum member count",
			},

			"maximum_member_count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Maximum member count",
			},

			"cooldown": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Cooldown value",
			},

			"termination_policy": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Termination policy",
			},

			"virtual_server_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "virtual server ID",
			},

			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Port number",
			},

			"health_check": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			// This has to be a TypeList, because TypeMap does not handle non-primitive
			// members properly.
			"virtual_guest_member_template": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        getModifiedVirtualGuestResource(),
				Description: "Virtual guest member template",
			},

			"network_vlan_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Set: func(v interface{}) int {
					return v.(int)
				},
				Description: "List of network VLAN ids",
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

// Returns a modified version of the virtual guest resource, with all members set to ForceNew = false.
// Otherwise a modified template parameter unnecessarily forces scale group drop/create
func getModifiedVirtualGuestResource() *schema.Resource {

	r := resourceIBMComputeVmInstance()
	// wait_time_minutes is only used in virtual_guest resource.
	delete(r.Schema, "wait_time_minutes")

	for _, elem := range r.Schema {
		elem.ForceNew = false
		elem.ConflictsWith = []string{}
	}

	return r
}

// Helper method to parse healthcheck data in the resource schema format to the SoftLayer datatypes
func buildHealthCheckFromResourceData(d map[string]interface{}) (datatypes.Network_Application_Delivery_Controller_LoadBalancer_Health_Check, error) {
	healthCheckOpts := datatypes.Network_Application_Delivery_Controller_LoadBalancer_Health_Check{
		Type: &datatypes.Network_Application_Delivery_Controller_LoadBalancer_Health_Check_Type{
			Keyname: sl.String(d["type"].(string)),
		},
	}

	if *healthCheckOpts.Type.Keyname == HEALTH_CHECK_TYPE_HTTP_CUSTOM {
		// Validate and apply type-specific fields
		healthCheckMethod, ok := d["custom_method"]
		if !ok {
			return datatypes.Network_Application_Delivery_Controller_LoadBalancer_Health_Check{}, errors.New("\"custom_method\" is required when HTTP-CUSTOM healthcheck is specified")
		}

		healthCheckRequest, ok := d["custom_request"]
		if !ok {
			return datatypes.Network_Application_Delivery_Controller_LoadBalancer_Health_Check{}, errors.New("\"custom_request\" is required when HTTP-CUSTOM healthcheck is specified")
		}

		healthCheckResponse, ok := d["custom_response"]
		if !ok {
			return datatypes.Network_Application_Delivery_Controller_LoadBalancer_Health_Check{}, errors.New("\"custom_response\" is required when HTTP-CUSTOM healthcheck is specified")
		}

		// HTTP-CUSTOM values are represented as an array of SoftLayer_Health_Check_Attributes
		healthCheckOpts.Attributes = []datatypes.Network_Application_Delivery_Controller_LoadBalancer_Health_Attribute{
			{
				Type: &datatypes.Network_Application_Delivery_Controller_LoadBalancer_Health_Attribute_Type{
					Keyname: sl.String("HTTP_CUSTOM_TYPE"),
				},
				Value: sl.String(healthCheckMethod.(string)),
			},
			{
				Type: &datatypes.Network_Application_Delivery_Controller_LoadBalancer_Health_Attribute_Type{
					Keyname: sl.String("LOCATION"),
				},
				Value: sl.String(healthCheckRequest.(string)),
			},
			{
				Type: &datatypes.Network_Application_Delivery_Controller_LoadBalancer_Health_Attribute_Type{
					Keyname: sl.String("EXPECTED_RESPONSE"),
				},
				Value: sl.String(healthCheckResponse.(string)),
			},
		}
	}

	return healthCheckOpts, nil
}

// Helper method to parse network vlan information in the resource schema format to the SoftLayer datatypes
func buildScaleVlansFromResourceData(v interface{}, meta interface{}) ([]datatypes.Scale_Network_Vlan, error) {
	vlanIds := v.([]interface{})
	scaleNetworkVlans := make([]datatypes.Scale_Network_Vlan, 0, len(vlanIds))

	for _, iVlanId := range vlanIds {
		vlanId := iVlanId.(int)
		scaleNetworkVlans = append(
			scaleNetworkVlans,
			datatypes.Scale_Network_Vlan{NetworkVlanId: &vlanId},
		)
	}

	return scaleNetworkVlans, nil
}

func getVirtualGuestTemplate(vGuestTemplateList []interface{}, meta interface{}) (datatypes.Virtual_Guest, error) {
	if len(vGuestTemplateList) != 1 {
		return datatypes.Virtual_Guest{},
			errors.New("Only one virtual_guest_member_template can be provided")
	}

	// Retrieve the map of virtual_guest_member_template attributes
	vGuestMap := vGuestTemplateList[0].(map[string]interface{})

	// Create an empty ResourceData instance for a IBM_Compute_VM_Instance resource
	vGuestResourceData := resourceIBMComputeVmInstance().Data(nil)

	// For each item in the map, call Set on the ResourceData.  This handles
	// validation and yields a completed ResourceData object
	for k, v := range vGuestMap {
		log.Printf("****** %s: %#v", k, v)
		err := vGuestResourceData.Set(k, v)
		if err != nil {
			return datatypes.Virtual_Guest{},
				fmt.Errorf("Error while parsing virtual_guest_member_template values: %s", err)
		}
	}
	dc := vGuestResourceData.Get("datacenter").(string)
	publicVlan := vGuestResourceData.Get("public_vlan_id").(int)
	privateVlan := vGuestResourceData.Get("private_vlan_id").(int)
	quote_id := 0
	// Get the virtual guest creation template from the completed resource data object
	vgs, err := getVirtualGuestTemplateFromResourceData(vGuestResourceData, meta, dc, publicVlan, privateVlan, quote_id)
	return vgs[0], err
}

func resourceIBMComputeAutoScaleGroupCreate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	accountServiceNoRetry := services.GetScaleGroupService(sess.SetRetries(0))

	virtualGuestTemplateOpts, err := getVirtualGuestTemplate(d.Get("virtual_guest_member_template").([]interface{}), meta)
	if err != nil {
		return fmt.Errorf("Error while parsing virtual_guest_member_template values: %s", err)
	}

	scaleNetworkVlans, err := buildScaleVlansFromResourceData(d.Get("network_vlan_ids").(*schema.Set).List(), meta)
	if err != nil {
		return fmt.Errorf("Error while parsing network vlan values: %s", err)
	}

	locationGroupRegionalId, err := getLocationGroupRegionalId(sess, d.Get("regional_group").(string))
	if err != nil {
		return err
	}

	// Build up our creation options
	opts := datatypes.Scale_Group{
		Name:                       sl.String(d.Get("name").(string)),
		Cooldown:                   sl.Int(d.Get("cooldown").(int)),
		MinimumMemberCount:         sl.Int(d.Get("minimum_member_count").(int)),
		MaximumMemberCount:         sl.Int(d.Get("maximum_member_count").(int)),
		SuspendedFlag:              sl.Bool(false),
		VirtualGuestMemberTemplate: &virtualGuestTemplateOpts,
		NetworkVlans:               scaleNetworkVlans,
		RegionalGroupId:            &locationGroupRegionalId,
	}

	opts.TerminationPolicy = &datatypes.Scale_Termination_Policy{
		KeyName: sl.String(d.Get("termination_policy").(string)),
	}

	opts.LoadBalancers, err = buildLoadBalancers(d)
	if err != nil {
		return fmt.Errorf("Error creating Scale Group: %s", err)
	}

	res, err := accountServiceNoRetry.CreateObject(&opts)
	if err != nil {
		return fmt.Errorf("Error creating Scale Group: %s", err)
	}

	d.SetId(strconv.Itoa(*res.Id))
	log.Printf("[INFO] Scale Group ID: %d", *res.Id)

	time.Sleep(60)

	// wait for scale group to become active
	_, err = waitForActiveStatus(d, meta)

	if err != nil {
		return fmt.Errorf("Error waiting for scale group (%s) to become active: %s", d.Id(), err)
	}

	return resourceIBMComputeAutoScaleGroupRead(d, meta)
}

func buildLoadBalancers(d *schema.ResourceData, ids ...int) ([]datatypes.Scale_LoadBalancer, error) {
	isLoadBalancerEmpty := true
	loadBalancers := []datatypes.Scale_LoadBalancer{{}}

	if virtualServerId, ok := d.GetOk("virtual_server_id"); ok {
		isLoadBalancerEmpty = false
		loadBalancers[0].VirtualServerId = sl.Int(virtualServerId.(int))
		if len(ids) > 0 {
			loadBalancers[0].Id = sl.Int(ids[0])
		}
	}

	if healthCheck, ok := d.GetOk("health_check"); ok {
		isLoadBalancerEmpty = false
		healthCheckOpts, err := buildHealthCheckFromResourceData(healthCheck.(map[string]interface{}))
		if err != nil {
			return []datatypes.Scale_LoadBalancer{}, fmt.Errorf("Error while parsing health check options: %s", err)
		}
		loadBalancers[0].HealthCheck = &healthCheckOpts
	}

	if port, ok := d.GetOk("port"); ok {
		isLoadBalancerEmpty = false
		loadBalancers[0].Port = sl.Int(port.(int))
	}

	if isLoadBalancerEmpty {
		return []datatypes.Scale_LoadBalancer{}, nil
	} else {
		return loadBalancers, nil
	}
}

func resourceIBMComputeAutoScaleGroupRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetScaleGroupService(sess)

	groupId, _ := strconv.Atoi(d.Id())

	slGroupObj, err := service.Id(groupId).Mask(strings.Join(IBMComputeAutoScaleGroupObjectMask, ",")).GetObject()
	if err != nil {
		// If the scale group is somehow already destroyed, mark as successfully gone
		if apiErr, ok := err.(sl.Error); ok && apiErr.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving autoscale Group: %s", err)
	}

	d.Set("name", slGroupObj.Name)
	if slGroupObj.RegionalGroup != nil && slGroupObj.RegionalGroup.Name != nil {
		d.Set("regional_group", slGroupObj.RegionalGroup.Name)
	}
	d.Set("minimum_member_count", slGroupObj.MinimumMemberCount)
	d.Set("maximum_member_count", slGroupObj.MaximumMemberCount)
	d.Set("cooldown", slGroupObj.Cooldown)
	d.Set("status", slGroupObj.Status.KeyName)
	d.Set("termination_policy", slGroupObj.TerminationPolicy.KeyName)
	if len(slGroupObj.LoadBalancers) > 0 {
		d.Set("virtual_server_id", slGroupObj.LoadBalancers[0].VirtualServerId)
		d.Set("port", slGroupObj.LoadBalancers[0].Port)

		// Health Check
		healthCheckObj := slGroupObj.LoadBalancers[0].HealthCheck
		currentHealthCheck := d.Get("health_check").(map[string]interface{})

		currentHealthCheck["type"] = *healthCheckObj.Type.Keyname

		if *healthCheckObj.Type.Keyname == HEALTH_CHECK_TYPE_HTTP_CUSTOM {
			for _, elem := range healthCheckObj.Attributes {
				switch *elem.Type.Keyname {
				case "HTTP_CUSTOM_TYPE":
					currentHealthCheck["custom_method"] = *elem.Value
				case "LOCATION":
					currentHealthCheck["custom_request"] = *elem.Value
				case "EXPECTED_RESPONSE":
					currentHealthCheck["custom_response"] = *elem.Value
				}
			}
		}

		d.Set("health_check", currentHealthCheck)
	}

	// Network Vlans
	vlanIds := make([]int, len(slGroupObj.NetworkVlans))
	for i, vlan := range slGroupObj.NetworkVlans {
		vlanIds[i] = *vlan.NetworkVlanId
	}
	d.Set("network_vlan_ids", vlanIds)

	virtualGuestTemplate := populateMemberTemplateResourceData(*slGroupObj.VirtualGuestMemberTemplate)
	d.Set("virtual_guest_member_template", virtualGuestTemplate)

	return nil
}

func populateMemberTemplateResourceData(template datatypes.Virtual_Guest) []map[string]interface{} {

	d := make(map[string]interface{})

	d["hostname"] = *template.Hostname
	d["domain"] = *template.Domain
	d["datacenter"] = *template.Datacenter.Name
	d["network_speed"] = *template.NetworkComponents[0].MaxSpeed
	d["cores"] = *template.StartCpus
	d["memory"] = *template.MaxMemory
	d["private_network_only"] = *template.PrivateNetworkOnlyFlag
	d["hourly_billing"] = *template.HourlyBillingFlag
	d["local_disk"] = *template.LocalDiskFlag

	// Guard against nil values for optional fields in virtual_guest resource
	d["dedicated_acct_host_only"] = sl.Get(template.DedicatedAccountHostOnlyFlag)
	d["os_reference_code"] = sl.Get(template.OperatingSystemReferenceCode)
	d["post_install_script_uri"] = sl.Get(template.PostInstallScriptUri)

	if template.PrimaryNetworkComponent != nil && template.PrimaryNetworkComponent.NetworkVlan != nil {
		d["public_vlan_id"] = sl.Get(template.PrimaryNetworkComponent.NetworkVlan.Id)
	}

	if template.PrimaryBackendNetworkComponent != nil && template.PrimaryBackendNetworkComponent.NetworkVlan != nil {
		d["private_vlan_id"] = sl.Get(template.PrimaryBackendNetworkComponent.NetworkVlan.Id)
	}
	if template.BlockDeviceTemplateGroup != nil {
		d["image_id"] = sl.Get(template.BlockDeviceTemplateGroup.GlobalIdentifier)
	}

	if len(template.UserData) > 0 {
		d["user_metadata"] = *template.UserData[0].Value
	}

	sshKeys := make([]interface{}, 0, len(template.SshKeys))
	for _, elem := range template.SshKeys {
		sshKeys = append(sshKeys, *elem.Id)
	}
	d["ssh_key_ids"] = sshKeys

	disks := make([]interface{}, 0, len(template.BlockDevices))
	for _, elem := range template.BlockDevices {
		disks = append(disks, *elem.DiskImage.Capacity)
	}
	d["disks"] = disks

	return []map[string]interface{}{d}
}

func resourceIBMComputeAutoScaleGroupUpdate(d *schema.ResourceData, meta interface{}) error {

	sess := meta.(ClientSession).SoftLayerSession()
	scaleGroupService := services.GetScaleGroupService(sess)
	scaleNetworkVlanService := services.GetScaleNetworkVlanService(sess)
	scaleLoadBalancerService := services.GetScaleLoadBalancerService(sess)
	scaleGroupServiceNoRetry := services.GetScaleGroupService(sess.SetRetries(0))

	groupId, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid ID. Must be an integer: %s", err)
	}

	// Fetch the complete object from SoftLayer, update with current values from the configuration, and send the
	// whole thing back to SoftLayer (effectively, a PUT)
	groupObj, err := scaleGroupService.Id(groupId).Mask(strings.Join(IBMComputeAutoScaleGroupObjectMask, ",")).GetObject()
	if err != nil {
		return fmt.Errorf("Error retrieving autoscale_group resource: %s", err)
	}

	groupObj.Name = sl.String(d.Get("name").(string))
	groupObj.MinimumMemberCount = sl.Int(d.Get("minimum_member_count").(int))
	groupObj.MaximumMemberCount = sl.Int(d.Get("maximum_member_count").(int))
	groupObj.Cooldown = sl.Int(d.Get("cooldown").(int))
	groupObj.TerminationPolicy.KeyName = sl.String(d.Get("termination_policy").(string))

	currentLoadBalancers := groupObj.LoadBalancers
	if len(currentLoadBalancers) > 0 {
		groupObj.LoadBalancers, err = buildLoadBalancers(d, *currentLoadBalancers[0].Id)
	} else {
		groupObj.LoadBalancers, err = buildLoadBalancers(d)
	}
	if err != nil {
		return fmt.Errorf("Error creating Scale Group: %s", err)
	}

	if d.HasChange("network_vlan_ids") {
		// Vlans require special handling:
		//
		// 1. Delete any scale_network_vlans which no longer appear in the updated configuration
		// 2. Pass the updated list of vlans to the Scale_Group.editObject function.  SoftLayer determines
		// which Vlans are new, and which already exist.

		_, newValue := d.GetChange("network_vlan_ids")
		newIds := newValue.(*schema.Set).List()

		// Delete all Vlans
		oldScaleVlans, err := scaleGroupService.
			Id(groupId).
			GetNetworkVlans()
		if err != nil {
			return fmt.Errorf("Could not retrieve current vlans for scale group (%d): %s", groupId, err)
		}

		for _, oldScaleVlan := range oldScaleVlans {
			_, err := scaleNetworkVlanService.Id(*oldScaleVlan.Id).DeleteObject()
			if err != nil {
				return fmt.Errorf("Error deleting scale network vlan %d: %s", *oldScaleVlan.Id, err)
			}
		}

		// Parse the new list of vlans into the appropriate input structure
		scaleVlans, err := buildScaleVlansFromResourceData(newIds, meta)

		if err != nil {
			return fmt.Errorf("Unable to parse network vlan options: %s", err)
		}

		groupObj.NetworkVlans = scaleVlans
	}

	if d.HasChange("virtual_guest_member_template") {
		virtualGuestTemplateOpts, err := getVirtualGuestTemplate(d.Get("virtual_guest_member_template").([]interface{}), meta)
		if err != nil {
			return fmt.Errorf("Unable to parse virtual guest member template options: %s", err)
		}

		groupObj.VirtualGuestMemberTemplate = &virtualGuestTemplateOpts

	}
	_, err = scaleGroupServiceNoRetry.Id(groupId).EditObject(&groupObj)
	if err != nil {
		return fmt.Errorf("Error received while editing autoscale_group: %s", err)
	}

	// wait for scale group to become active
	_, err = waitForActiveStatus(d, meta)

	if err != nil {
		return fmt.Errorf("Error waiting for scale group (%s) to become active: %s", d.Id(), err)
	}

	// Delete a load balancer if there is the load balancer in a scale group
	// and a request doesn't have virtual_server_id, port, and health_check.
	if len(currentLoadBalancers) > 0 && len(groupObj.LoadBalancers) <= 0 {
		_, err = scaleLoadBalancerService.Id(*currentLoadBalancers[0].Id).DeleteObject()
		if err != nil {
			return fmt.Errorf("Error received while deleting loadbalancers: %s", err)
		}
	}

	return nil
}

func resourceIBMComputeAutoScaleGroupDelete(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	scaleGroupService := services.GetScaleGroupService(sess)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting scale group: %s", err)
	}

	log.Printf("[INFO] Deleting scale group: %d", id)
	_, err = scaleGroupService.Id(id).ForceDeleteObject()
	if err != nil {
		return fmt.Errorf("Error deleting scale group: %s", err)
	}

	d.SetId("")

	return nil
}

func waitForActiveStatus(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	sess := meta.(ClientSession).SoftLayerSession()
	scaleGroupService := services.GetScaleGroupService(sess)

	log.Printf("Waiting for scale group (%s) to become active", d.Id())
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return nil, fmt.Errorf("The scale group ID %s must be numeric", d.Id())
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"BUSY", "SCALING", "SUSPENDED"},
		Target:  []string{"ACTIVE"},
		Refresh: func() (interface{}, string, error) {
			// get the status of the scale group
			result, err := scaleGroupService.Id(id).Mask("status.keyName,minimumMemberCount," +
				"virtualGuestMembers[virtualGuest[primaryBackendIpAddress,primaryIpAddress,privateNetworkOnlyFlag,fullyQualifiedDomainName]]").
				GetObject()
			if err != nil {
				if apiErr, ok := err.(sl.Error); ok && apiErr.StatusCode == 404 {
					return nil, "", fmt.Errorf("The scale group %d does not exist anymore: %s", id, err)
				}

				return result, "BUSY", nil // Retry
			}

			status := "BUSY"

			// Return "BUSY" if member VMs don't have ip addresses.
			for _, scaleMemberVirtualGuest := range result.VirtualGuestMembers {
				// Checking primary backend IP address.
				if scaleMemberVirtualGuest.VirtualGuest.PrimaryBackendIpAddress == nil {
					log.Printf("The member vm of scale group does not have private IP yet. Hostname : %s",
						*scaleMemberVirtualGuest.VirtualGuest.FullyQualifiedDomainName)
					return result, status, nil
				}

				// Checking primary IP address.
				if !(*scaleMemberVirtualGuest.VirtualGuest.PrivateNetworkOnlyFlag) &&
					scaleMemberVirtualGuest.VirtualGuest.PrimaryIpAddress == nil {
					log.Printf("The member vm of scale group does not have IP yet. Hostname : %s",
						*scaleMemberVirtualGuest.VirtualGuest.FullyQualifiedDomainName)
					return result, status, nil
				}
			}
			if result.Status.KeyName != nil {
				status = *result.Status.KeyName
				log.Printf("The status of scale group with id (%d) is (%s)", id, *result.Status.KeyName)
			} else {
				log.Printf("Could not get the status of scale group with id (%d). Retrying...", id)
			}

			return result, status, nil
		},
		Timeout:    120 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func resourceIBMComputeAutoScaleGroupExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	sess := meta.(ClientSession).SoftLayerSession()
	scaleGroupService := services.GetScaleGroupService(sess)

	groupId, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	result, err := scaleGroupService.Id(groupId).Mask("id").GetObject()
	if err != nil {
		if apiErr, ok := err.(sl.Error); ok {
			if apiErr.StatusCode == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}
	return result.Id != nil && *result.Id == groupId, nil
}

func getLocationGroupRegionalId(sess *session.Session, locationGroupRegionalName string) (int, error) {
	locationGroupRegionals, err := services.GetLocationGroupRegionalService(sess).
		Mask("id,name").
		// FIXME: Someday, filters may actually work in SoftLayer
		//Filter(filter.Build(
		//	filter.Path("name").Eq(locationGroupRegionalName))).
		//Limit(1).
		GetAllObjects()

	if err != nil {
		return -1, err
	}

	if len(locationGroupRegionals) < 1 {
		return -1, fmt.Errorf("Invalid location group regional: %s", locationGroupRegionalName)
	}

	for _, locationGroupRegional := range locationGroupRegionals {
		if *locationGroupRegional.Name == locationGroupRegionalName {
			return *locationGroupRegional.Id, nil
		}
	}

	return -1, fmt.Errorf("Invalid regional_group_id: %s", locationGroupRegionalName)
}
