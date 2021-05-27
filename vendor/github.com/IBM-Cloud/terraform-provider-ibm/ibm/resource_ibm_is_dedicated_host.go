// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

const (
	isDedicatedHostStable     = "stable"
	isDedicatedHostDeleting   = "deleting"
	isDedicatedHostDeleteDone = "done"
	isDedicatedHostFailed     = "failed"

	isDedicatedHostUpdating             = "updating"
	isDedicatedHostProvisioningDone     = "done"
	isDedicatedHostWaiting              = "waiting"
	isDedicatedHostSuspended            = "suspended"
	isDedicatedHostActionStatusStopping = "stopping"
	isDedicatedHostActionStatusStopped  = "stopped"
	isDedicatedHostStatusPending        = "pending"
	isDedicatedHostStatusRunning        = "running"
	isDedicatedHostStatusFailed         = "failed"
)

func resourceIbmIsDedicatedHost() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIbmIsDedicatedHostCreate,
		Read:     resourceIbmIsDedicatedHostRead,
		Update:   resourceIbmIsDedicatedHostUpdate,
		Delete:   resourceIbmIsDedicatedHostDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"instance_placement_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If set to true, instances can be placed on this dedicated host.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: InvokeValidator("ibm_is_dedicated_host", "name"),
				Description:  "The unique user-defined name for this dedicated host. If unspecified, the name will be a hyphenated list of randomly-selected words.",
			},
			"profile": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Globally unique name of the dedicated host profile to use for this dedicated host.",
			},
			"resource_group": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The unique identifier for the resource group to use. If unspecified, the account's [default resourcegroup](https://cloud.ibm.com/apidocs/resource-manager#introduction) is used.",
			},
			"host_group": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The unique identifier of the dedicated host group for this dedicated host.",
			},
			"available_memory": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The amount of memory in gibibytes that is currently available for instances.",
			},
			"available_vcpu": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The available VCPU for the dedicated host.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"architecture": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VCPU architecture.",
						},
						"count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of VCPUs assigned.",
						},
					},
				},
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the dedicated host was created.",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this dedicated host.",
			},
			"disks": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of the dedicated host's disks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"available": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The remaining space left for instance placement in GB (gigabytes).",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the disk was created.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this disk.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this disk.",
						},
						"instance_disks": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Instance disks that are on this dedicated host disk.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"deleted": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this instance disk.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this instance disk.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this disk.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"interface_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The disk interface used for attaching the diskThe enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered.",
						},
						"lifecycle_state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of this dedicated host disk.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined or system-provided name for this disk.",
						},
						"provisionable": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether this dedicated host disk is available for instance disk creation.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of resource referenced.",
						},
						"size": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The size of the disk in GB (gigabytes).",
						},
						"supported_instance_interface_types": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The instance disk interfaces supported for this dedicated host disk.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this dedicated host.",
			},
			"instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Array of instances that are allocated to this dedicated host.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this virtual server instance.",
						},
						"deleted": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this virtual server instance.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this virtual server instance.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this virtual server instance (and default system hostname).",
						},
					},
				},
			},
			"lifecycle_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the dedicated host resource.",
			},
			"memory": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total amount of memory in gibibytes for this host.",
			},
			"provisionable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this dedicated host is available for instance creation.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of resource referenced.",
			},
			"socket_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of sockets for this host.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The administrative state of the dedicated host.The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the dedicated host on which the unexpected property value was encountered.",
			},
			"supported_instance_profiles": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Array of instance profiles that can be used by instances placed on this dedicated host.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this virtual server instance profile.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique name for this virtual server instance profile.",
						},
					},
				},
			},
			"vcpu": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The total VCPU of the dedicated host.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"architecture": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VCPU architecture.",
						},
						"count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of VCPUs assigned.",
						},
					},
				},
			},
			"zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The globally unique name of the zone this dedicated host resides in.",
			},
		},
	}
}

func resourceIbmIsDedicatedHostValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Optional:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		})

	resourceValidator := ResourceValidator{ResourceName: "ibm_is_dedicated_host", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmIsDedicatedHostCreate(d *schema.ResourceData, meta interface{}) error {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return err
	}
	createDedicatedHostOptions := &vpcv1.CreateDedicatedHostOptions{}
	dedicatedHostPrototype := vpcv1.DedicatedHostPrototype{}

	if dhname, ok := d.GetOk("name"); ok {

		namestr := dhname.(string)
		dedicatedHostPrototype.Name = &namestr
	}
	if insplacementenabled, ok := d.GetOk("instance_placement_enabled"); ok {
		insplacementenabledbool := insplacementenabled.(bool)
		dedicatedHostPrototype.InstancePlacementEnabled = &insplacementenabledbool
	}

	if dhprofile, ok := d.GetOk("profile"); ok {
		dhprofilename := dhprofile.(string)
		dedicatedHostProfileIdentity := vpcv1.DedicatedHostProfileIdentity{
			Name: &dhprofilename,
		}
		dedicatedHostPrototype.Profile = &dedicatedHostProfileIdentity
	}

	if dhgroup, ok := d.GetOk("host_group"); ok {
		dhgroupid := dhgroup.(string)
		dedicatedHostGroupIdentity := vpcv1.DedicatedHostGroupIdentity{
			ID: &dhgroupid,
		}
		dedicatedHostPrototype.Group = &dedicatedHostGroupIdentity
	}

	if resgroup, ok := d.GetOk("resource_group"); ok {
		resgroupid := resgroup.(string)
		resourceGroupIdentity := vpcv1.ResourceGroupIdentity{
			ID: &resgroupid,
		}
		dedicatedHostPrototype.ResourceGroup = &resourceGroupIdentity
	}

	createDedicatedHostOptions.SetDedicatedHostPrototype(&dedicatedHostPrototype)

	dedicatedHost, response, err := vpcClient.CreateDedicatedHostWithContext(context.TODO(), createDedicatedHostOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateDedicatedHostWithContext failed %s\n%s", err, response)
		return err
	}

	d.SetId(*dedicatedHost.ID)

	_, err = isWaitForDedicatedHostAvailable(vpcClient, d.Id(), d.Timeout(schema.TimeoutCreate), d)
	if err != nil {
		return err
	}

	return resourceIbmIsDedicatedHostRead(d, meta)
}

func resourceIbmIsDedicatedHostRead(d *schema.ResourceData, meta interface{}) error {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return err
	}

	getDedicatedHostOptions := &vpcv1.GetDedicatedHostOptions{}

	getDedicatedHostOptions.SetID(d.Id())

	dedicatedHost, response, err := vpcClient.GetDedicatedHostWithContext(context.TODO(), getDedicatedHostOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetDedicatedHostWithContext failed %s\n%s", err, response)
		return err
	}

	if err = d.Set("available_memory", intValue(dedicatedHost.AvailableMemory)); err != nil {
		return fmt.Errorf("Error setting available_memory: %s", err)
	}
	availableVcpuMap := resourceIbmIsDedicatedHostVCPUToMap(*dedicatedHost.AvailableVcpu)
	if err = d.Set("available_vcpu", []map[string]interface{}{availableVcpuMap}); err != nil {
		return fmt.Errorf("Error setting available_vcpu: %s", err)
	}
	if err = d.Set("created_at", dedicatedHost.CreatedAt.String()); err != nil {
		return fmt.Errorf("Error setting created_at: %s", err)
	}
	if err = d.Set("crn", dedicatedHost.CRN); err != nil {
		return fmt.Errorf("Error setting crn: %s", err)
	}
	disks := []map[string]interface{}{}
	for _, disksItem := range dedicatedHost.Disks {
		disksItemMap := resourceIbmIsDedicatedHostDedicatedHostDiskToMap(disksItem)
		disks = append(disks, disksItemMap)
	}
	if err = d.Set("disks", disks); err != nil {
		return fmt.Errorf("Error setting disks: %s", err)
	}
	d.Set("host_group", *dedicatedHost.Group.ID)

	if err = d.Set("href", dedicatedHost.Href); err != nil {
		return fmt.Errorf("Error setting href: %s", err)
	}
	if err = d.Set("instance_placement_enabled", dedicatedHost.InstancePlacementEnabled); err != nil {
		return fmt.Errorf("Error setting instance_placement_enabled: %s", err)
	}
	instances := []map[string]interface{}{}
	for _, instancesItem := range dedicatedHost.Instances {
		instancesItemMap := resourceIbmIsDedicatedHostInstanceReferenceToMap(instancesItem)
		instances = append(instances, instancesItemMap)
	}
	if err = d.Set("instances", instances); err != nil {
		return fmt.Errorf("Error setting instances: %s", err)
	}
	if err = d.Set("lifecycle_state", dedicatedHost.LifecycleState); err != nil {
		return fmt.Errorf("Error setting lifecycle_state: %s", err)
	}
	if err = d.Set("memory", intValue(dedicatedHost.Memory)); err != nil {
		return fmt.Errorf("Error setting memory: %s", err)
	}
	if err = d.Set("name", dedicatedHost.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}

	if err = d.Set("profile", *dedicatedHost.Profile.Name); err != nil {
		return fmt.Errorf("Error setting profile: %s", err)
	}
	if err = d.Set("provisionable", dedicatedHost.Provisionable); err != nil {
		return fmt.Errorf("Error setting provisionable: %s", err)
	}
	if err = d.Set("resource_group", *dedicatedHost.ResourceGroup.ID); err != nil {
		return fmt.Errorf("Error setting resource_group: %s", err)
	}
	if err = d.Set("resource_type", dedicatedHost.ResourceType); err != nil {
		return fmt.Errorf("Error setting resource_type: %s", err)
	}
	if err = d.Set("socket_count", intValue(dedicatedHost.SocketCount)); err != nil {
		return fmt.Errorf("Error setting socket_count: %s", err)
	}
	if err = d.Set("state", dedicatedHost.State); err != nil {
		return fmt.Errorf("Error setting state: %s", err)
	}
	supportedInstanceProfiles := []map[string]interface{}{}
	for _, supportedInstanceProfilesItem := range dedicatedHost.SupportedInstanceProfiles {
		supportedInstanceProfilesItemMap := resourceIbmIsDedicatedHostInstanceProfileReferenceToMap(supportedInstanceProfilesItem)
		supportedInstanceProfiles = append(supportedInstanceProfiles, supportedInstanceProfilesItemMap)
	}
	if err = d.Set("supported_instance_profiles", supportedInstanceProfiles); err != nil {
		return fmt.Errorf("Error setting supported_instance_profiles: %s", err)
	}
	vcpuMap := resourceIbmIsDedicatedHostVCPUToMap(*dedicatedHost.Vcpu)
	if err = d.Set("vcpu", []map[string]interface{}{vcpuMap}); err != nil {
		return fmt.Errorf("Error setting vcpu: %s", err)
	}

	if err = d.Set("zone", *dedicatedHost.Zone.Name); err != nil {
		return fmt.Errorf("Error setting zone: %s", err)
	}

	return nil
}

func resourceIbmIsDedicatedHostVCPUToMap(vCPU vpcv1.Vcpu) map[string]interface{} {
	vCPUMap := map[string]interface{}{}

	vCPUMap["architecture"] = vCPU.Architecture
	vCPUMap["count"] = intValue(vCPU.Count)

	return vCPUMap
}

func resourceIbmIsDedicatedHostInstanceReferenceToMap(instanceReference vpcv1.InstanceReference) map[string]interface{} {
	instanceReferenceMap := map[string]interface{}{}

	instanceReferenceMap["crn"] = instanceReference.CRN
	if instanceReference.Deleted != nil {
		DeletedMap := resourceIbmIsDedicatedHostInstanceReferenceDeletedToMap(*instanceReference.Deleted)
		instanceReferenceMap["deleted"] = []map[string]interface{}{DeletedMap}
	}
	instanceReferenceMap["href"] = instanceReference.Href
	instanceReferenceMap["id"] = instanceReference.ID
	instanceReferenceMap["name"] = instanceReference.Name

	return instanceReferenceMap
}

func resourceIbmIsDedicatedHostInstanceReferenceDeletedToMap(instanceReferenceDeleted vpcv1.InstanceReferenceDeleted) map[string]interface{} {
	instanceReferenceDeletedMap := map[string]interface{}{}

	instanceReferenceDeletedMap["more_info"] = instanceReferenceDeleted.MoreInfo

	return instanceReferenceDeletedMap
}

func resourceIbmIsDedicatedHostInstanceProfileReferenceToMap(instanceProfileReference vpcv1.InstanceProfileReference) map[string]interface{} {
	instanceProfileReferenceMap := map[string]interface{}{}

	instanceProfileReferenceMap["href"] = instanceProfileReference.Href
	instanceProfileReferenceMap["name"] = instanceProfileReference.Name

	return instanceProfileReferenceMap
}

func resourceIbmIsDedicatedHostUpdate(d *schema.ResourceData, meta interface{}) error {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return err
	}

	updateDedicatedHostOptions := &vpcv1.UpdateDedicatedHostOptions{}

	updateDedicatedHostOptions.SetID(d.Id())

	hasChange := false

	dedicatedHostPrototypemap := map[string]interface{}{}

	if d.HasChange("name") {

		dedicatedHostPrototypemap["name"] = d.Get("name").(interface{})
		hasChange = true
	}
	if d.HasChange("instance_placement_enabled") {

		dedicatedHostPrototypemap["instance_placement_enabled"] = d.Get("instance_placement_enabled").(interface{})
		hasChange = true
	}
	if d.HasChange("profile") {
		dedicatedHostPrototypemap["profile"] = d.Get("profile").(interface{})
		hasChange = true
	}
	if d.HasChange("resource_group") {
		dedicatedHostPrototypemap["resource_group"] = d.Get("resource_group").(interface{})
		hasChange = true
	}
	if d.HasChange("host_group") {
		dedicatedHostPrototypemap["group"] = d.Get("host_group").(interface{})
		hasChange = true
	}

	if hasChange {
		updateDedicatedHostOptions.SetDedicatedHostPatch(dedicatedHostPrototypemap)
		_, response, err := vpcClient.UpdateDedicatedHostWithContext(context.TODO(), updateDedicatedHostOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateDedicatedHostWithContext fails %s\n%s", err, response)
			return err
		}
	}

	return resourceIbmIsDedicatedHostRead(d, meta)
}

func resourceIbmIsDedicatedHostDelete(d *schema.ResourceData, meta interface{}) error {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return err
	}

	getDedicatedHostOptions := &vpcv1.GetDedicatedHostOptions{}

	getDedicatedHostOptions.SetID(d.Id())

	dedicatedHost, response, err := vpcClient.GetDedicatedHostWithContext(context.TODO(), getDedicatedHostOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetDedicatedHostWithContext failed %s\n%s", err, response)
		return err
	}
	if dedicatedHost != nil && dedicatedHost.LifecycleState != nil && *dedicatedHost.LifecycleState != isDedicatedHostSuspended && *dedicatedHost.LifecycleState != isDedicatedHostFailed {

		updateDedicatedHostOptions := &vpcv1.UpdateDedicatedHostOptions{}
		dedicatedHostPrototypeMap := map[string]interface{}{}
		dedicatedHostPrototypeMap["instance_placement_enabled"] = core.BoolPtr(false)
		updateDedicatedHostOptions.SetID(d.Id())
		updateDedicatedHostOptions.SetDedicatedHostPatch(dedicatedHostPrototypeMap)
		_, updateresponse, err := vpcClient.UpdateDedicatedHostWithContext(context.TODO(), updateDedicatedHostOptions)
		if err != nil {
			log.Printf("[DEBUG] Failed disabling instance placement %s\n%s", err, updateresponse)
			return err
		}
	}
	deleteDedicatedHostOptions := &vpcv1.DeleteDedicatedHostOptions{}

	deleteDedicatedHostOptions.SetID(d.Id())

	response, err = vpcClient.DeleteDedicatedHostWithContext(context.TODO(), deleteDedicatedHostOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteDedicatedHostWithContext failed %s\n%s", err, response)
		return err
	}
	_, err = isWaitForDedicatedHostDelete(vpcClient, d, d.Id())
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func isWaitForDedicatedHostDelete(instanceC *vpcv1.VpcV1, d *schema.ResourceData, id string) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending: []string{isDedicatedHostDeleting, isDedicatedHostStable},
		Target:  []string{isDedicatedHostDeleteDone, ""},
		Refresh: func() (interface{}, string, error) {
			getdhoptions := &vpcv1.GetDedicatedHostOptions{
				ID: &id,
			}
			dedicatedhost, response, err := instanceC.GetDedicatedHost(getdhoptions)
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					return dedicatedhost, isDedicatedHostDeleteDone, nil
				}
				return nil, "", fmt.Errorf("Error getting dedicated Host: %s\n%s", err, response)
			}
			if *dedicatedhost.State == isDedicatedHostFailed {
				return dedicatedhost, *dedicatedhost.State, fmt.Errorf("The  Dedicated host %s failed to delete: %v", d.Id(), err)
			}
			return dedicatedhost, isDedicatedHostDeleting, nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isWaitForDedicatedHostAvailable(instanceC *vpcv1.VpcV1, id string, timeout time.Duration, d *schema.ResourceData) (interface{}, error) {
	log.Printf("Waiting for dedicated host (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{isDedicatedHostStatusPending, isDedicatedHostUpdating, isDedicatedHostWaiting},
		Target:     []string{isDedicatedHostFailed, isDedicatedHostStable, isDedicatedHostSuspended},
		Refresh:    isDedicatedHostRefreshFunc(instanceC, id, d),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isDedicatedHostRefreshFunc(instanceC *vpcv1.VpcV1, id string, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getinsOptions := &vpcv1.GetDedicatedHostOptions{
			ID: &id,
		}
		dhost, response, err := instanceC.GetDedicatedHost(getinsOptions)
		if dhost == nil || err != nil {
			return nil, "", fmt.Errorf("Error getting dedicated host : %s\n%s", err, response)
		}
		d.Set("state", *dhost.State)
		d.Set("lifecycle_state", *dhost.LifecycleState)

		if *dhost.LifecycleState == isDedicatedHostSuspended || *dhost.LifecycleState == isDedicatedHostFailed {

			return dhost, *dhost.LifecycleState, fmt.Errorf("status of dedicated host is %s : \n%s", *dhost.LifecycleState, response)

		}
		return dhost, *dhost.LifecycleState, nil
	}
}

func resourceIbmIsDedicatedHostDedicatedHostDiskToMap(dedicatedHostDisk vpcv1.DedicatedHostDisk) map[string]interface{} {
	dedicatedHostDiskMap := map[string]interface{}{}

	dedicatedHostDiskMap["available"] = intValue(dedicatedHostDisk.Available)
	dedicatedHostDiskMap["created_at"] = dedicatedHostDisk.CreatedAt.String()
	dedicatedHostDiskMap["href"] = dedicatedHostDisk.Href
	dedicatedHostDiskMap["id"] = dedicatedHostDisk.ID
	instanceDisks := []map[string]interface{}{}
	for _, instanceDisksItem := range dedicatedHostDisk.InstanceDisks {
		instanceDisksItemMap := resourceIbmIsDedicatedHostInstanceDiskReferenceToMap(instanceDisksItem)
		instanceDisks = append(instanceDisks, instanceDisksItemMap)
		// TODO: handle InstanceDisks of type TypeList -- list of non-primitive, not model items
	}
	dedicatedHostDiskMap["instance_disks"] = instanceDisks
	dedicatedHostDiskMap["interface_type"] = dedicatedHostDisk.InterfaceType
	if dedicatedHostDisk.LifecycleState != nil {
		dedicatedHostDiskMap["lifecycle_state"] = dedicatedHostDisk.LifecycleState
	}
	dedicatedHostDiskMap["name"] = dedicatedHostDisk.Name
	dedicatedHostDiskMap["provisionable"] = dedicatedHostDisk.Provisionable
	dedicatedHostDiskMap["resource_type"] = dedicatedHostDisk.ResourceType
	dedicatedHostDiskMap["size"] = intValue(dedicatedHostDisk.Size)
	dedicatedHostDiskMap["supported_instance_interface_types"] = dedicatedHostDisk.SupportedInstanceInterfaceTypes

	return dedicatedHostDiskMap
}

func resourceIbmIsDedicatedHostInstanceDiskReferenceToMap(instanceDiskReference vpcv1.InstanceDiskReference) map[string]interface{} {
	instanceDiskReferenceMap := map[string]interface{}{}

	if instanceDiskReference.Deleted != nil {
		DeletedMap := resourceIbmIsDedicatedHostInstanceDiskReferenceDeletedToMap(*instanceDiskReference.Deleted)
		instanceDiskReferenceMap["deleted"] = []map[string]interface{}{DeletedMap}
	}
	instanceDiskReferenceMap["href"] = instanceDiskReference.Href
	instanceDiskReferenceMap["id"] = instanceDiskReference.ID
	instanceDiskReferenceMap["name"] = instanceDiskReference.Name
	instanceDiskReferenceMap["resource_type"] = instanceDiskReference.ResourceType

	return instanceDiskReferenceMap
}

func resourceIbmIsDedicatedHostInstanceDiskReferenceDeletedToMap(instanceDiskReferenceDeleted vpcv1.InstanceDiskReferenceDeleted) map[string]interface{} {
	instanceDiskReferenceDeletedMap := map[string]interface{}{}

	instanceDiskReferenceDeletedMap["more_info"] = instanceDiskReferenceDeleted.MoreInfo

	return instanceDiskReferenceDeletedMap
}
