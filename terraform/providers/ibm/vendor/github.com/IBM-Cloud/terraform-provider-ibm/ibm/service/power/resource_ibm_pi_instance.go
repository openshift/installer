// Copyright IBM Corp. 2017, 2021, 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
)

func ResourceIBMPIInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPIInstanceCreate,
		ReadContext:   resourceIBMPIInstanceRead,
		UpdateContext: resourceIBMPIInstanceUpdate,
		DeleteContext: resourceIBMPIInstanceDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
		CustomizeDiff: customdiff.Sequence(
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return flex.ResourcePowerUserTagsCustomizeDiff(diff)
			},
		),

		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_AffinityInstance: {
				ConflictsWith: []string{Arg_AffinityVolume},
				Description:   "PVM Instance (ID or Name) to base storage affinity policy against; required if requesting storage affinity and pi_affinity_volume is not provided",
				Optional:      true,
				Type:          schema.TypeString,
			},
			Arg_AffinityPolicy: {
				Description:  "Affinity policy for pvm instance being created; ignored if pi_storage_pool provided; for policy affinity requires one of pi_affinity_instance or pi_affinity_volume to be specified; for policy anti-affinity requires one of pi_anti_affinity_instances or pi_anti_affinity_volumes to be specified",
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{Affinity, AntiAffinity}),
			},
			Arg_AffinityVolume: {
				ConflictsWith: []string{Arg_AffinityInstance},
				Description:   "Volume (ID or Name) to base storage affinity policy against; required if requesting affinity and pi_affinity_instance is not provided",
				Optional:      true,
				Type:          schema.TypeString,
			},
			Arg_AntiAffinityInstances: {
				ConflictsWith: []string{Arg_AntiAffinityVolumes},
				Description:   "List of pvmInstances to base storage anti-affinity policy against; required if requesting anti-affinity and pi_anti_affinity_volumes is not provided",
				Elem:          &schema.Schema{Type: schema.TypeString},
				Optional:      true,
				Type:          schema.TypeList,
			},
			Arg_AntiAffinityVolumes: {
				ConflictsWith: []string{Arg_AntiAffinityInstances},
				Description:   "List of volumes to base storage anti-affinity policy against; required if requesting anti-affinity and pi_anti_affinity_instances is not provided",
				Elem:          &schema.Schema{Type: schema.TypeString},
				Optional:      true,
				Type:          schema.TypeList,
			},
			Arg_BootVolumeReplicationEnabled: {
				Description: "Indicates if the boot volume should be replication enabled or not.",
				ForceNew:    true,
				Optional:    true,
				Type:        schema.TypeBool,
			},
			Arg_CloudInstanceID: {
				Description: "This is the Power Instance id that is assigned to the account",
				ForceNew:    true,
				Required:    true,
				Type:        schema.TypeString,
			},
			Arg_DeploymentTarget: {
				Description: "The deployment of a dedicated host.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_ID: {
							Description: "The uuid of the host group or host.",
							Required:    true,
							Type:        schema.TypeString,
						},
						Attr_Type: {
							Description:  "The deployment target type. Supported values are `host` and `hostGroup`.",
							Required:     true,
							Type:         schema.TypeString,
							ValidateFunc: validate.ValidateAllowedStringValues([]string{Host, HostGroup}),
						},
					},
				},
				MaxItems:     1,
				Optional:     true,
				RequiredWith: []string{Arg_SysType},
				Type:         schema.TypeSet,
			},
			Arg_DeploymentType: {
				Description:  "Custom Deployment Type Information",
				ForceNew:     true,
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{DeploymentTypeEpic, DeploymentTypeVMNoStorage}),
			},
			Arg_HealthStatus: {
				Default:      OK,
				Description:  "Allow the user to set the status of the lpar so that they can connect to it faster",
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{OK, Warning}),
			},
			Arg_IBMiCSS: {
				Description: "IBM i Cloud Storage Solution",
				Optional:    true,
				Type:        schema.TypeBool,
			},
			Arg_IBMiPHA: {
				Description: "IBM i Power High Availability",
				Optional:    true,
				Type:        schema.TypeBool,
			},
			Arg_IBMiRDSUsers: {
				Description: "IBM i Rational Dev Studio Number of User Licenses",
				Optional:    true,
				Type:        schema.TypeInt,
			},
			Arg_ImageID: {
				Description:      "PI instance image id",
				DiffSuppressFunc: flex.ApplyOnce,
				Required:         true,
				Type:             schema.TypeString,
			},
			Arg_InstanceName: {
				Description: "PI Instance name",
				Required:    true,
				Type:        schema.TypeString,
			},
			Arg_KeyPairName: {
				Description: "SSH key name",
				ForceNew:    true,
				Optional:    true,
				Type:        schema.TypeString,
			},
			Arg_LicenseRepositoryCapacity: {
				Computed:    true,
				Deprecated:  "This field is deprecated.",
				Description: "The VTL license repository capacity TB value",
				Optional:    true,
				Type:        schema.TypeInt,
			},
			Arg_Memory: {
				Computed:      true,
				ConflictsWith: []string{Arg_SAPProfileID},
				Description:   "Memory size",
				Optional:      true,
				Type:          schema.TypeFloat,
			},
			Arg_Network: {
				Description:      "List of one or more networks to attach to the instance",
				DiffSuppressFunc: flex.ApplyOnce,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_IPAddress: {
							Computed: true,
							Optional: true,
							Type:     schema.TypeString,
						},
						Attr_MacAddress: {
							Computed: true,
							Type:     schema.TypeString,
						},
						Attr_NetworkID: {
							Required: true,
							Type:     schema.TypeString,
						},
						Attr_NetworkInterfaceID: {
							Computed:    true,
							Description: "ID of the network interface.",
							Type:        schema.TypeString,
						},
						Attr_NetworkName: {
							Computed: true,
							Type:     schema.TypeString,
						},
						Attr_NetworkSecurityGroupIDs: {
							Computed:    true,
							Description: "Network security groups that the network interface is a member of. There is a limit of 1 network security group in the array. If not specified, default network security group is used.",
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Type:        schema.TypeSet,
						},
						Attr_NetworkSecurityGroupsHref: {
							Computed:    true,
							Description: "Links to the network security groups that the network interface is a member of.",
							Elem:        &schema.Schema{Type: schema.TypeString},
							Type:        schema.TypeList,
						},
						Attr_Type: {
							Computed: true,
							Type:     schema.TypeString,
						},
						Attr_ExternalIP: {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
				Required: true,
				Type:     schema.TypeList,
			},
			Arg_PinPolicy: {
				Default:      None,
				Description:  "Pin Policy of the instance",
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{None, Soft, Hard}),
			},
			Arg_PlacementGroupID: {
				Description: "Placement group ID",
				Computed:    true,
				Optional:    true,
				Type:        schema.TypeString,
			},
			Arg_ProcType: {
				Computed:      true,
				ConflictsWith: []string{Arg_SAPProfileID},
				Description:   "Instance processor type",
				Optional:      true,
				Type:          schema.TypeString,
				ValidateFunc:  validate.ValidateAllowedStringValues([]string{Dedicated, Shared, Capped}),
			},
			Arg_Processors: {
				Computed:      true,
				ConflictsWith: []string{Arg_SAPProfileID},
				Description:   "Processors count",
				Optional:      true,
				Type:          schema.TypeFloat,
			},
			Arg_Replicants: {
				Default:     1,
				Description: "PI Instance replicas count",
				ForceNew:    true,
				Optional:    true,
				Type:        schema.TypeInt,
			},
			Arg_ReplicationPolicy: {
				Default:      None,
				Description:  "Replication policy for the PI Instance",
				ForceNew:     true,
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{Affinity, AntiAffinity, None}),
			},
			Arg_ReplicationScheme: {
				Default:      Suffix,
				Description:  "Replication scheme",
				ForceNew:     true,
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{Prefix, Suffix}),
			},
			Arg_ReplicationSites: {
				Description: "Indicates the replication sites of the boot volume.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				ForceNew:    true,
				Optional:    true,
				Set:         schema.HashString,
				Type:        schema.TypeSet,
			},
			Arg_RetainVirtualSerialNumber: {
				Default:     false,
				Description: "Indicates whether to retain virtual serial number when changed or deleted.",
				Optional:    true,
				Type:        schema.TypeBool,
			},
			Arg_SAPProfileID: {
				ConflictsWith: []string{Arg_Processors, Arg_Memory, Arg_ProcType},
				Description:   "SAP Profile ID for the amount of cores and memory",
				Optional:      true,
				Type:          schema.TypeString,
			},
			Arg_SAPDeploymentType: {
				Description: "Custom SAP Deployment Type Information",
				ForceNew:    true,
				Optional:    true,
				Type:        schema.TypeString,
			},
			Arg_SharedProcessorPool: {
				ConflictsWith: []string{Arg_SAPProfileID},
				Description:   "Shared Processor Pool the instance is deployed on",
				ForceNew:      true,
				Optional:      true,
				Type:          schema.TypeString,
			},
			Arg_StoragePool: {
				Computed:    true,
				Description: "Storage Pool for server deployment; if provided then pi_storage_pool_affinity will be ignored; Only valid when you deploy one of the IBM supplied stock images. Storage pool for a custom image (an imported image or an image that is created from a VM capture) defaults to the storage pool the image was created in",
				Optional:    true,
				Type:        schema.TypeString,
			},
			Arg_StoragePoolAffinity: {
				Default:     true,
				Description: "Indicates if all volumes attached to the server must reside in the same storage pool",
				Optional:    true,
				Type:        schema.TypeBool,
			},
			Arg_StorageType: {
				Computed:    true,
				Description: "Storage type for server deployment; if pi_storage_type is not provided the storage type will default to tier3",
				Optional:    true,
				Type:        schema.TypeString,
			},
			Arg_StorageConnection: {
				Description:  "Storage Connectivity Group for server deployment",
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{vSCSI, MaxVolumeSupport}),
			},
			Arg_SysType: {
				Computed:    true,
				Description: "PI Instance system type",
				ForceNew:    true,
				Optional:    true,
				Type:        schema.TypeString,
			},
			Arg_UserData: {
				Description: "Base64 encoded data to be passed in for invoking a cloud init script",
				ForceNew:    true,
				Optional:    true,
				Type:        schema.TypeString,
			},
			Arg_UserTags: {
				Computed:    true,
				Description: "The user tags attached to this resource.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Set:         schema.HashString,
				Type:        schema.TypeSet,
			},
			Arg_VirtualCoresAssigned: {
				Computed:    true,
				Description: "Virtual Cores Assigned to the PVMInstance",
				Optional:    true,
				Type:        schema.TypeInt,
			},
			Arg_VirtualOpticalDevice: {
				Description:  "Virtual Machine's Cloud Initialization Virtual Optical Device",
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{Attach}),
			},
			Arg_VirtualSerialNumber: {
				ConflictsWith: []string{Arg_SAPProfileID},
				Description:   "Virtual Serial Number information",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_Description: {
							Description: "Description of the Virtual Serial Number",
							Optional:    true,
							Type:        schema.TypeString,
						},
						Attr_Serial: {
							Description:      "Provide an existing reserved Virtual Serial Number or specify 'auto-assign' for auto generated Virtual Serial Number.",
							Required:         true,
							DiffSuppressFunc: supressVSNDiffAutoAssign,
							Type:             schema.TypeString,
						},
					},
				},
				MaxItems: 1,
				Optional: true,
				Type:     schema.TypeList,
			},
			Arg_VolumeIDs: {
				Description:      "List of PI volumes",
				DiffSuppressFunc: flex.ApplyOnce,
				Elem:             &schema.Schema{Type: schema.TypeString},
				Optional:         true,
				Set:              schema.HashString,
				Type:             schema.TypeSet,
			},

			// Attributes
			Attr_CRN: {
				Computed:    true,
				Description: "The CRN of this resource.",
				Type:        schema.TypeString,
			},
			Attr_DedicatedHostID: {
				Computed:    true,
				Description: "The dedicated host ID where the shared processor pool resides.",
				Type:        schema.TypeString,
			},
			Attr_Fault: {
				Computed:    true,
				Description: "Fault information.",
				Type:        schema.TypeMap,
			},
			Attr_HealthStatus: {
				Computed:    true,
				Description: "PI Instance health status",
				Type:        schema.TypeString,
			},
			Attr_IBMiRDS: {
				Computed:    true,
				Description: "IBM i Rational Dev Studio",
				Optional:    false,
				Required:    false,
				Type:        schema.TypeBool,
			},
			Attr_InstanceID: {
				Computed:    true,
				Description: "Instance ID",
				Type:        schema.TypeString,
			},
			Attr_MaxMemory: {
				Computed:    true,
				Description: "Maximum memory size",
				Type:        schema.TypeFloat,
			},
			Attr_MaxProcessors: {
				Computed:    true,
				Description: "Maximum number of processors",
				Type:        schema.TypeFloat,
			},
			Attr_MaxVirtualCores: {
				Computed:    true,
				Description: "Maximum Virtual Cores Assigned to the PVMInstance",
				Type:        schema.TypeInt,
			},
			Attr_MinMemory: {
				Computed:    true,
				Description: "Minimum memory",
				Type:        schema.TypeFloat,
			},
			Attr_MinProcessors: {
				Computed:    true,
				Description: "Minimum number of the CPUs",
				Type:        schema.TypeFloat,
			},
			Attr_MinVirtualCores: {
				Computed:    true,
				Description: "Minimum Virtual Cores Assigned to the PVMInstance",
				Type:        schema.TypeInt,
			},
			Attr_OperatingSystem: {
				Computed:    true,
				Description: "Operating System",
				Type:        schema.TypeString,
			},
			Attr_OSType: {
				Computed:    true,
				Description: "OS Type",
				Type:        schema.TypeString,
			},
			Attr_PinPolicy: {
				Computed:    true,
				Description: "PIN Policy of the Instance",
				Type:        schema.TypeString,
			},
			Attr_Progress: {
				Computed:    true,
				Description: "Progress of the operation",
				Type:        schema.TypeFloat,
			},
			Attr_SharedProcessorPoolID: {
				Computed:    true,
				Description: "Shared Processor Pool ID the instance is deployed on",
				Type:        schema.TypeString,
			},
			Attr_Status: {
				Computed:    true,
				Description: "PI instance status",
				Type:        schema.TypeString,
			},
		},
	}
}

func resourceIBMPIInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("Now in the PowerVMCreate")
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	client := instance.NewIBMPIInstanceClient(ctx, sess, cloudInstanceID)
	sapClient := instance.NewIBMPISAPInstanceClient(ctx, sess, cloudInstanceID)
	imageClient := instance.NewIBMPIImageClient(ctx, sess, cloudInstanceID)

	var pvmList *models.PVMInstanceList
	if _, ok := d.GetOk(Arg_SAPProfileID); ok {
		pvmList, err = createSAPInstance(d, sapClient)
	} else {
		pvmList, err = createPVMInstance(d, client, imageClient)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	var instanceReadyStatus string
	if r, ok := d.GetOk(Arg_HealthStatus); ok {
		instanceReadyStatus = r.(string)
	}

	// id is a combination of the cloud instance id and all of the pvm instance ids
	id := cloudInstanceID
	for _, pvm := range *pvmList {
		id += "/" + *pvm.PvmInstanceID
	}

	d.SetId(id)

	for _, s := range *pvmList {
		if dt, ok := d.GetOk(Arg_DeploymentType); ok && dt.(string) == DeploymentTypeVMNoStorage {
			_, err = isWaitForPIInstanceShutoff(ctx, client, *s.PvmInstanceID, instanceReadyStatus, d.Timeout(schema.TimeoutCreate))
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			_, err = isWaitForPIInstanceAvailable(ctx, client, *s.PvmInstanceID, instanceReadyStatus, d.Timeout(schema.TimeoutCreate))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	// If Storage Pool Affinity is given as false we need to update the vm instance.
	// Default value is true which indicates that all volumes attached to the server
	// must reside in the same storage pool.
	storagePoolAffinity := d.Get(Arg_StoragePoolAffinity).(bool)
	if !storagePoolAffinity {
		for _, s := range *pvmList {
			body := &models.PVMInstanceUpdate{
				StoragePoolAffinity: &storagePoolAffinity,
			}
			// This is a synchronous process hence no need to check for health status
			_, err = client.Update(*s.PvmInstanceID, body)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	// If user tags are set, make sure tags are set correctly before moving on
	if _, ok := d.GetOk(Arg_UserTags); ok {
		oldList, newList := d.GetChange(Arg_UserTags)
		for _, s := range *pvmList {
			if s.Crn != "" {
				err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, string(s.Crn), "", UserTagType)
				if err != nil {
					log.Printf("Error on update of pi instance (%s) pi_user_tags during creation: %s", *s.PvmInstanceID, err)
				}
			}
		}
	}

	// If virtual optical device provided then update cloud initialization
	if vod, ok := d.GetOk(Arg_VirtualOpticalDevice); ok {
		for _, s := range *pvmList {
			body := &models.PVMInstanceUpdate{
				CloudInitialization: &models.CloudInitialization{
					VirtualOpticalDevice: vod.(string),
				},
			}
			_, err = client.Update(*s.PvmInstanceID, body)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceIBMPIInstanceRead(ctx, d, meta)
}

func resourceIBMPIInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	idArr, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := idArr[0]
	instanceID := idArr[1]

	client := instance.NewIBMPIInstanceClient(ctx, sess, cloudInstanceID)
	powervmdata, err := client.Get(instanceID)
	if err != nil {
		return diag.FromErr(err)
	}

	if powervmdata.Crn != "" {
		d.Set(Attr_CRN, powervmdata.Crn)
		tags, err := flex.GetTagsUsingCRN(meta, string(powervmdata.Crn))
		if err != nil {
			log.Printf("Error on get of ibm pi instance (%s) pi_user_tags: %s", *powervmdata.PvmInstanceID, err)
		}
		d.Set(Arg_UserTags, tags)
	}
	d.Set(Arg_Memory, powervmdata.Memory)
	d.Set(Arg_Processors, powervmdata.Processors)
	if powervmdata.Status != nil {
		d.Set(Attr_Status, powervmdata.Status)
	}
	d.Set(Arg_CloudInstanceID, cloudInstanceID)
	d.Set(Arg_ImageID, powervmdata.ImageID)
	d.Set(Arg_InstanceName, powervmdata.ServerName)
	d.Set(Arg_ProcType, powervmdata.ProcType)
	d.Set(Arg_StoragePool, powervmdata.StoragePool)
	d.Set(Arg_StoragePoolAffinity, powervmdata.StoragePoolAffinity)
	d.Set(Attr_InstanceID, powervmdata.PvmInstanceID)
	d.Set(Attr_MinProcessors, powervmdata.Minproc)
	d.Set(Attr_Progress, powervmdata.Progress)
	if *powervmdata.PlacementGroup != None {
		d.Set(Arg_PlacementGroupID, powervmdata.PlacementGroup)
	}
	d.Set(Arg_SharedProcessorPool, powervmdata.SharedProcessorPool)
	d.Set(Attr_SharedProcessorPoolID, powervmdata.SharedProcessorPoolID)

	networksMap := []map[string]interface{}{}
	if powervmdata.Networks != nil {
		for _, n := range powervmdata.Networks {
			if n != nil {
				v := map[string]interface{}{
					Attr_ExternalIP:         n.ExternalIP,
					Attr_IPAddress:          n.IPAddress,
					Attr_MacAddress:         n.MacAddress,
					Attr_NetworkID:          n.NetworkID,
					Attr_NetworkInterfaceID: n.NetworkInterfaceID,
					Attr_NetworkName:        n.NetworkName,
					Attr_Type:               n.Type,
				}
				if len(n.NetworkSecurityGroupIDs) > 0 {
					v[Attr_NetworkSecurityGroupIDs] = n.NetworkSecurityGroupIDs
				}
				if len(n.NetworkSecurityGroupsHref) > 0 {
					v[Attr_NetworkSecurityGroupsHref] = n.NetworkSecurityGroupsHref
				}
				networksMap = append(networksMap, v)
			}
		}
	}
	d.Set(Arg_Network, networksMap)

	if powervmdata.SapProfile != nil && powervmdata.SapProfile.ProfileID != nil {
		d.Set(Arg_SAPProfileID, powervmdata.SapProfile.ProfileID)
	}
	d.Set(Arg_SysType, powervmdata.SysType)
	d.Set(Attr_DedicatedHostID, powervmdata.DedicatedHostID)
	d.Set(Attr_MinMemory, powervmdata.Minmem)
	d.Set(Attr_MaxProcessors, powervmdata.Maxproc)
	d.Set(Attr_MaxMemory, powervmdata.Maxmem)
	d.Set(Attr_PinPolicy, powervmdata.PinPolicy)
	d.Set(Attr_OperatingSystem, powervmdata.OperatingSystem)
	d.Set(Attr_OSType, powervmdata.OsType)

	if powervmdata.Health != nil {
		d.Set(Attr_HealthStatus, powervmdata.Health.Status)
	}
	if powervmdata.VirtualCores != nil {
		d.Set(Arg_VirtualCoresAssigned, powervmdata.VirtualCores.Assigned)
		d.Set(Attr_MaxVirtualCores, powervmdata.VirtualCores.Max)
		d.Set(Attr_MinVirtualCores, powervmdata.VirtualCores.Min)
	}
	d.Set(Arg_DeploymentType, powervmdata.DeploymentType)
	d.Set(Arg_LicenseRepositoryCapacity, powervmdata.LicenseRepositoryCapacity)
	if powervmdata.SoftwareLicenses != nil {
		d.Set(Arg_IBMiCSS, powervmdata.SoftwareLicenses.IbmiCSS)
		d.Set(Arg_IBMiPHA, powervmdata.SoftwareLicenses.IbmiPHA)
		d.Set(Attr_IBMiRDS, powervmdata.SoftwareLicenses.IbmiRDS)
		if *powervmdata.SoftwareLicenses.IbmiRDS {
			d.Set(Arg_IBMiRDSUsers, powervmdata.SoftwareLicenses.IbmiRDSUsers)
		} else {
			d.Set(Arg_IBMiRDSUsers, 0)
		}
	}
	if powervmdata.Fault != nil {
		d.Set(Attr_Fault, flattenPvmInstanceFault(powervmdata.Fault))
	} else {
		d.Set(Attr_Fault, nil)
	}

	if powervmdata.VirtualSerialNumber != nil {
		d.Set(Arg_VirtualSerialNumber, flattenVirtualSerialNumberToList(powervmdata.VirtualSerialNumber))
	} else {
		d.Set(Arg_VirtualSerialNumber, nil)
	}

	return nil
}

func resourceIBMPIInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get(Arg_InstanceName).(string)
	mem := d.Get(Arg_Memory).(float64)
	procs := d.Get(Arg_Processors).(float64)
	processortype := d.Get(Arg_ProcType).(string)
	assignedVirtualCores := int64(d.Get(Arg_VirtualCoresAssigned).(int))

	if d.Get(Attr_HealthStatus) == Warning {
		return diag.Errorf("the operation cannot be performed when the lpar health in the WARNING State")
	}

	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.Errorf("failed to get the session from the IBM Cloud Service")
	}

	cloudInstanceID, instanceID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := instance.NewIBMPIInstanceClient(ctx, sess, cloudInstanceID)

	// Check if cloud instance is capable of changing virtual cores
	cloudInstanceClient := instance.NewIBMPICloudInstanceClient(ctx, sess, cloudInstanceID)
	cloudInstance, err := cloudInstanceClient.Get(cloudInstanceID)
	if err != nil {
		return diag.FromErr(err)
	}
	cores_enabled := checkCloudInstanceCapability(cloudInstance, CUSTOM_VIRTUAL_CORES)

	if d.HasChanges(Arg_InstanceName, Arg_VirtualOpticalDevice) {
		if d.HasChange(Arg_InstanceName) && d.HasChange(Arg_VirtualOpticalDevice) {
			oldVOD, _ := d.GetChange(Arg_VirtualOpticalDevice)
			d.Set(Arg_VirtualOpticalDevice, oldVOD)
			return diag.Errorf("updates to %s and %s are mutually exclusive", Arg_InstanceName, Arg_VirtualOpticalDevice)
		}
		body := &models.PVMInstanceUpdate{}
		if d.HasChange(Arg_InstanceName) {
			body.ServerName = name
		}
		if d.HasChange(Arg_VirtualOpticalDevice) {
			body.CloudInitialization = &models.CloudInitialization{}
			if vod, ok := d.GetOk(Arg_VirtualOpticalDevice); ok {
				body.CloudInitialization.VirtualOpticalDevice = vod.(string)
			} else {
				body.CloudInitialization.VirtualOpticalDevice = Detach
			}
		}
		_, err = client.Update(instanceID, body)
		if err != nil {
			return diag.Errorf("failed to update the lpar: %v", err)
		}
		_, err = isWaitForPIInstanceAvailableOrShutoffAfterUpdate(ctx, client, instanceID, OK, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange(Arg_ProcType) {
		// Stop the lpar
		status := d.Get(Attr_Status).(string)
		if strings.ToLower(status) == State_Shutoff {
			log.Printf("the lpar is in the shutoff state. Nothing to do . Moving on ")
		} else {
			err := stopLparForResourceChange(ctx, client, instanceID, d)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		// Modify
		log.Printf("At this point the lpar should be off. Executing the Processor Update Change")
		updatebody := &models.PVMInstanceUpdate{ProcType: processortype}
		if cores_enabled {
			log.Printf("support for %s is enabled", CUSTOM_VIRTUAL_CORES)
			updatebody.VirtualCores = &models.VirtualCores{Assigned: &assignedVirtualCores}
		} else {
			log.Printf("no virtual cores support enabled for this customer..")
		}
		_, err = client.Update(instanceID, updatebody)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = isWaitForPIInstanceStopped(ctx, client, instanceID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}

		// Start the lpar
		err := startLparAfterResourceChange(ctx, client, instanceID, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Virtual core will be updated only if service instance capability is enabled
	if d.HasChange(Arg_VirtualCoresAssigned) {
		body := &models.PVMInstanceUpdate{
			VirtualCores: &models.VirtualCores{Assigned: &assignedVirtualCores},
		}
		_, err = client.Update(instanceID, body)
		if err != nil {
			return diag.Errorf("failed to update the lpar with the change for virtual cores: %v", err)
		}
		_, err = isWaitForPIInstanceAvailable(ctx, client, instanceID, OK, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Start of the change for Memory and Processors
	if d.HasChange(Arg_Memory) || d.HasChange(Arg_Processors) {

		maxMemLpar := d.Get(Attr_MaxMemory).(float64)
		maxCPULpar := d.Get(Attr_MaxProcessors).(float64)

		if mem > maxMemLpar || procs > maxCPULpar {
			log.Printf("Will require a shutdown to perform the change")
		} else {
			log.Printf("maxMemLpar is set to %f", maxMemLpar)
			log.Printf("maxCPULpar is set to %f", maxCPULpar)
		}

		instanceState := d.Get(Attr_Status).(string)
		log.Printf("the instance state is %s", instanceState)

		if (mem > maxMemLpar || procs > maxCPULpar) && strings.ToLower(instanceState) != State_Shutoff {
			err = performChangeAndReboot(ctx, client, d, instanceID, mem, procs)
			if err != nil {
				return diag.FromErr(err)
			}

		} else {
			body := &models.PVMInstanceUpdate{
				Memory:     mem,
				Processors: procs,
			}
			if cores_enabled {
				log.Printf("support for %s is enabled", CUSTOM_VIRTUAL_CORES)
				body.VirtualCores = &models.VirtualCores{Assigned: &assignedVirtualCores}
			} else {
				log.Printf("no virtual cores support enabled for this customer..")
			}

			_, err = client.Update(instanceID, body)
			if err != nil {
				return diag.Errorf("failed to update the lpar with the change %v", err)
			}
			instanceReadyStatus := d.Get(Arg_HealthStatus).(string)
			_, err = isWaitForPIInstanceAvailableOrShutoffAfterUpdate(ctx, client, instanceID, instanceReadyStatus, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	// License repository capacity will be updated only if service instance is a vtl instance
	// might need to check if lrc was set
	if d.HasChange(Arg_LicenseRepositoryCapacity) {
		lrc := d.Get(Arg_LicenseRepositoryCapacity).(int64)
		body := &models.PVMInstanceUpdate{
			LicenseRepositoryCapacity: lrc,
		}
		_, err = client.Update(instanceID, body)
		if err != nil {
			return diag.Errorf("failed to update the lpar with the change for license repository capacity %s", err)
		}
		_, err = isWaitForPIInstanceAvailable(ctx, client, instanceID, OK, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			diag.FromErr(err)
		}
	}

	if d.HasChange(Arg_SAPProfileID) {
		// Stop the lpar
		status := d.Get(Attr_Status).(string)
		if strings.ToLower(status) == State_Shutoff {
			log.Printf("the lpar is in the shutoff state. Nothing to do... Moving on ")
		} else {
			err := stopLparForResourceChange(ctx, client, instanceID, d)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		// Update the profile id
		profileID := d.Get(Arg_SAPProfileID).(string)
		body := &models.PVMInstanceUpdate{
			SapProfileID: profileID,
		}
		_, err = client.Update(instanceID, body)
		if err != nil {
			return diag.Errorf("failed to update the lpar with the change for sap profile: %v", err)
		}

		// Wait for the resize to complete and status to reset
		_, err = isWaitForPIInstanceStopped(ctx, client, instanceID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}

		// Start the lpar
		err := startLparAfterResourceChange(ctx, client, instanceID, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange(Arg_StoragePoolAffinity) {
		storagePoolAffinity := d.Get(Arg_StoragePoolAffinity).(bool)
		body := &models.PVMInstanceUpdate{
			StoragePoolAffinity: &storagePoolAffinity,
		}
		// This is a synchronous process hence no need to check for health status
		_, err = client.Update(instanceID, body)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange(Arg_PlacementGroupID) {
		pgClient := instance.NewIBMPIPlacementGroupClient(ctx, sess, cloudInstanceID)

		oldRaw, newRaw := d.GetChange(Arg_PlacementGroupID)
		old := oldRaw.(string)
		new := newRaw.(string)

		if len(strings.TrimSpace(old)) > 0 {
			placementGroupID := old
			//remove server from old placement group
			body := &models.PlacementGroupServer{
				ID: &instanceID,
			}
			pgID, err := pgClient.DeleteMember(placementGroupID, body)
			if err != nil {
				// ignore delete member error where the server is already not in the PG
				if !strings.Contains(err.Error(), "is not part of placement-group") {
					return diag.FromErr(err)
				}
			} else {
				_, err = isWaitForPIInstancePlacementGroupDelete(ctx, pgClient, *pgID.ID, instanceID, d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}

		if len(strings.TrimSpace(new)) > 0 {
			placementGroupID := new
			// add server to a new placement group
			body := &models.PlacementGroupServer{
				ID: &instanceID,
			}
			pgID, err := pgClient.AddMember(placementGroupID, body)
			if err != nil {
				return diag.FromErr(err)
			} else {
				_, err = isWaitForPIInstancePlacementGroupAdd(ctx, pgClient, *pgID.ID, instanceID, d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}
	if d.HasChanges(Arg_IBMiCSS, Arg_IBMiPHA, Arg_IBMiRDSUsers) {
		status := d.Get(Attr_Status).(string)
		if strings.ToLower(status) == State_Active {
			log.Printf("the lpar is in the Active state, continuing with update")
		} else {
			_, err = isWaitForPIInstanceAvailable(ctx, client, instanceID, OK, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.FromErr(err)
			}
		}

		sl := &models.SoftwareLicenses{}
		sl.IbmiCSS = flex.PtrToBool(d.Get(Arg_IBMiCSS).(bool))
		sl.IbmiPHA = flex.PtrToBool(d.Get(Arg_IBMiPHA).(bool))
		ibmrdsUsers := d.Get(Arg_IBMiRDSUsers).(int)
		if ibmrdsUsers < 0 {
			return diag.Errorf("request with  IBM i Rational Dev Studio property requires IBM i Rational Dev Studio number of users")
		}
		sl.IbmiRDS = flex.PtrToBool(ibmrdsUsers > 0)
		sl.IbmiRDSUsers = int64(ibmrdsUsers)

		updatebody := &models.PVMInstanceUpdate{SoftwareLicenses: sl}
		_, err = client.Update(instanceID, updatebody)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = isWaitForPIInstanceSoftwareLicenses(ctx, client, instanceID, sl, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange(Arg_UserTags) {
		if crn, ok := d.GetOk(Attr_CRN); ok {
			oldList, newList := d.GetChange(Arg_UserTags)
			err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, crn.(string), "", UserTagType)
			if err != nil {
				log.Printf("Error on update of pi instance (%s) pi_user_tags: %s", instanceID, err)
			}
		}
	}

	if d.HasChange(Arg_VirtualSerialNumber) {
		vsnClient := instance.NewIBMPIVSNClient(ctx, sess, cloudInstanceID)

		if d.HasChange(Arg_VirtualSerialNumber + ".0." + Attr_Serial) {
			instanceRestart := false
			status := d.Get(Attr_Status).(string)
			if strings.ToLower(status) != State_Shutoff {
				err := stopLparForResourceChange(ctx, client, instanceID, d)
				if err != nil {
					return diag.FromErr(err)
				}
				instanceRestart = true
			}

			oldVSN, newVSN := d.GetChange(Arg_VirtualSerialNumber)
			if len(oldVSN.([]interface{})) > 0 {
				retainVSN := d.Get(Arg_RetainVirtualSerialNumber).(bool)
				deleteBody := &models.DeleteServerVirtualSerialNumber{
					RetainVSN: retainVSN,
				}
				err := vsnClient.PVMInstanceDeleteVSN(instanceID, deleteBody)
				if err != nil {
					return diag.FromErr(err)
				}

				_, err = isWaitForPIInstanceStopped(ctx, client, instanceID, d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return diag.FromErr(err)
				}
			}

			if len(newVSN.([]interface{})) > 0 {
				newVSNMap := newVSN.([]interface{})[0].(map[string]interface{})
				description := newVSNMap[Attr_Description].(string)
				serial := newVSNMap[Attr_Serial].(string)
				addBody := &models.AddServerVirtualSerialNumber{
					Description: description,
					Serial:      &serial,
				}
				_, err := vsnClient.PVMInstanceAttachVSN(instanceID, addBody)
				if err != nil {
					return diag.FromErr(err)
				}

				_, err = isWaitForPIInstanceStopped(ctx, client, instanceID, d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return diag.FromErr(err)
				}
			}

			if instanceRestart {
				err = startLparAfterResourceChange(ctx, client, instanceID, d)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}

		if !d.HasChange(Arg_VirtualSerialNumber+".0."+Attr_Serial) && d.HasChange(Arg_VirtualSerialNumber+".0."+Attr_Description) {
			newDescriptionString := d.Get(Arg_VirtualSerialNumber + ".0." + Attr_Description).(string)
			updateBody := &models.UpdateServerVirtualSerialNumber{
				Description: &newDescriptionString,
			}
			_, err = vsnClient.PVMInstanceUpdateVSN(instanceID, updateBody)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceIBMPIInstanceRead(ctx, d, meta)
}

func resourceIBMPIInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	idArr, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := idArr[0]
	client := instance.NewIBMPIInstanceClient(ctx, sess, cloudInstanceID)
	for _, instanceID := range idArr[1:] {
		retainVSNBool := d.Get(Arg_RetainVirtualSerialNumber).(bool)
		if _, ok := d.GetOk(Arg_VirtualSerialNumber); ok && retainVSNBool {
			body := &models.PVMInstanceDelete{
				RetainVSN: &retainVSNBool,
			}
			err = client.DeleteWithBody(instanceID, body)
		} else {
			err = client.Delete(instanceID)
		}
		if err != nil {
			return diag.FromErr(err)
		}
	}

	for _, instanceID := range idArr[1:] {
		_, err = isWaitForPIInstanceDeleted(ctx, client, instanceID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId("")
	return nil
}

func isWaitForPIInstanceDeleted(ctx context.Context, client *instance.IBMPIInstanceClient, id string, timeout time.Duration) (interface{}, error) {

	log.Printf("Waiting for  (%s) to be deleted.", id)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Retry, State_Deleting},
		Target:     []string{State_NotFound},
		Refresh:    isPIInstanceDeleteRefreshFunc(client, id),
		Delay:      Timeout_Delay,
		MinTimeout: Timeout_Active,
		Timeout:    timeout,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isPIInstanceDeleteRefreshFunc(client *instance.IBMPIInstanceClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		pvm, err := client.Get(id)
		if err != nil {
			log.Printf("The power vm does not exist")
			return pvm, State_NotFound, nil
		}
		return pvm, State_Deleting, nil
	}
}

func isWaitForPIInstanceAvailable(ctx context.Context, client *instance.IBMPIInstanceClient, id string, instanceReadyStatus string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for PIInstance (%s) to be available and active ", id)

	queryTimeOut := Timeout_Active
	if instanceReadyStatus == Warning {
		queryTimeOut = Timeout_Warning
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Pending, State_Build, Warning},
		Target:     []string{State_Active, OK, State_Error, "", State_Shutoff},
		Refresh:    isPIInstanceRefreshFunc(client, id, instanceReadyStatus),
		Delay:      Timeout_Delay,
		MinTimeout: queryTimeOut,
		Timeout:    timeout,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isPIInstanceRefreshFunc(client *instance.IBMPIInstanceClient, id, instanceReadyStatus string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {

		pvm, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}
		// Check for `instanceReadyStatus` health status and also the final health status "OK"
		if strings.ToLower(*pvm.Status) == State_Active && (pvm.Health.Status == instanceReadyStatus || pvm.Health.Status == OK) {
			return pvm, State_Active, nil
		}
		if strings.ToLower(*pvm.Status) == State_Error {
			if pvm.Fault != nil {
				err = fmt.Errorf("failed to create the lpar: %s", pvm.Fault.Message)
			} else {
				err = fmt.Errorf("failed to create the lpar")
			}
			return pvm, *pvm.Status, err
		}

		return pvm, State_Build, nil
	}
}

func isWaitForPIInstanceAvailableOrShutoffAfterUpdate(ctx context.Context, client *instance.IBMPIInstanceClient, id string, instanceReadyStatus string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for PIInstance (%s) to be available and active or shutoff ", id)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Updating, Warning},
		Target:     []string{State_Active, OK, State_Shutoff},
		Refresh:    isPIInstanceShutoffOrActiveAfterResourceChange(client, id, instanceReadyStatus),
		Delay:      Timeout_Delay,
		MinTimeout: 5 * time.Minute,
		Timeout:    timeout,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isPIInstanceShutoffOrActiveAfterResourceChange(client *instance.IBMPIInstanceClient, id string, instanceReadyStatus string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {

		pvm, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}
		if strings.ToLower(*pvm.Status) == State_Active && (pvm.Health.Status == instanceReadyStatus || pvm.Health.Status == OK) {
			log.Printf("The lpar is now active after the resource change...")
			return pvm, State_Active, nil
		}
		if strings.ToLower(*pvm.Status) == State_Shutoff && pvm.Health.Status == OK {
			log.Printf("The lpar is now off after the resource change...")
			return pvm, State_Shutoff, nil
		}

		return pvm, State_Updating, nil
	}
}

func isWaitForPIInstancePlacementGroupAdd(ctx context.Context, client *instance.IBMPIPlacementGroupClient, pgID string, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for PIInstance Placement Group (%s) to be updated ", id)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Adding},
		Target:     []string{State_Added},
		Refresh:    isPIInstancePlacementGroupAddRefreshFunc(client, pgID, id),
		Delay:      Timeout_Delay,
		MinTimeout: Timeout_Active,
		Timeout:    timeout,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isPIInstancePlacementGroupAddRefreshFunc(client *instance.IBMPIPlacementGroupClient, pgID string, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		pg, err := client.Get(pgID)
		if err != nil {
			return nil, "", err
		}
		for _, x := range pg.Members {
			if x == id {
				return pg, State_Added, nil
			}
		}
		return pg, State_Adding, nil
	}
}

func isWaitForPIInstancePlacementGroupDelete(ctx context.Context, client *instance.IBMPIPlacementGroupClient, pgID string, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for PIInstance Placement Group (%s) to be updated ", id)

	queryTimeOut := Timeout_Active

	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Deleting},
		Target:     []string{State_Deleted},
		Refresh:    isPIInstancePlacementGroupDeleteRefreshFunc(client, pgID, id),
		Delay:      Timeout_Delay,
		MinTimeout: queryTimeOut,
		Timeout:    timeout,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isPIInstancePlacementGroupDeleteRefreshFunc(client *instance.IBMPIPlacementGroupClient, pgID string, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		pg, err := client.Get(pgID)
		if err != nil {
			return nil, "", err
		}
		for _, x := range pg.Members {
			if x == id {
				return pg, State_Deleting, nil
			}
		}
		return pg, State_Deleted, nil
	}
}

func isWaitForPIInstanceSoftwareLicenses(ctx context.Context, client *instance.IBMPIInstanceClient, id string, softwareLicenses *models.SoftwareLicenses, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for PIInstance Software Licenses (%s) to be updated ", id)

	queryTimeOut := Timeout_Active

	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_InProgress},
		Target:     []string{State_Available},
		Refresh:    isPIInstanceSoftwareLicensesRefreshFunc(client, id, softwareLicenses),
		Delay:      Timeout_Delay,
		MinTimeout: queryTimeOut,
		Timeout:    timeout,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isPIInstanceSoftwareLicensesRefreshFunc(client *instance.IBMPIInstanceClient, id string, softwareLicenses *models.SoftwareLicenses) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {

		pvm, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}

		// Check that each software license we modified has been updated
		if softwareLicenses.IbmiCSS != nil {
			if *softwareLicenses.IbmiCSS != *pvm.SoftwareLicenses.IbmiCSS {
				return pvm, State_InProgress, nil
			}
		}

		if softwareLicenses.IbmiPHA != nil {
			if *softwareLicenses.IbmiPHA != *pvm.SoftwareLicenses.IbmiPHA {
				return pvm, State_InProgress, nil
			}
		}

		if softwareLicenses.IbmiRDS != nil {
			// If the update set IBMiRDS to false, don't check IBMiRDSUsers as it will be updated on the terraform side on the read
			if !*softwareLicenses.IbmiRDS {
				if *softwareLicenses.IbmiRDS != *pvm.SoftwareLicenses.IbmiRDS {
					return pvm, State_InProgress, nil
				}
			} else if (*softwareLicenses.IbmiRDS != *pvm.SoftwareLicenses.IbmiRDS) || (softwareLicenses.IbmiRDSUsers != pvm.SoftwareLicenses.IbmiRDSUsers) {
				return pvm, State_InProgress, nil
			}
		}

		return pvm, State_Available, nil
	}
}

func isWaitForPIInstanceShutoff(ctx context.Context, client *instance.IBMPIInstanceClient, id string, instanceReadyStatus string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for PIInstance (%s) to be shutoff and health active ", id)

	queryTimeOut := Timeout_Active
	if instanceReadyStatus == Warning {
		queryTimeOut = Timeout_Warning
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Pending, State_Build, Warning},
		Target:     []string{OK, State_Error, "", State_Shutoff},
		Refresh:    isPIInstanceShutoffRefreshFunc(client, id, instanceReadyStatus),
		Delay:      Timeout_Delay,
		MinTimeout: queryTimeOut,
		Timeout:    timeout,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isPIInstanceShutoffRefreshFunc(client *instance.IBMPIInstanceClient, id, instanceReadyStatus string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {

		pvm, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}
		if strings.ToLower(*pvm.Status) == State_Shutoff && (pvm.Health.Status == instanceReadyStatus || pvm.Health.Status == OK) {
			return pvm, State_Shutoff, nil
		}
		if strings.ToLower(*pvm.Status) == State_Error {
			if pvm.Fault != nil {
				err = fmt.Errorf("failed to create the lpar: %s", pvm.Fault.Message)
			} else {
				err = fmt.Errorf("failed to create the lpar")
			}
			return pvm, *pvm.Status, err
		}

		return pvm, State_Build, nil
	}
}

// This function takes the input string and encodes into base64 if isn't already encoded
func encodeBase64(userData string) string {
	_, err := base64.StdEncoding.DecodeString(userData)
	if err != nil {
		return base64.StdEncoding.EncodeToString([]byte(userData))
	}
	return userData
}

func isWaitForPIInstanceStopped(ctx context.Context, client *instance.IBMPIInstanceClient, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for PIInstance (%s) to be stopped and powered off ", id)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Stopping, State_Resize, State_VerifyResize, Warning},
		Target:     []string{OK, State_Shutoff},
		Refresh:    isPIInstanceRefreshFuncOff(client, id),
		Delay:      Timeout_Delay,
		MinTimeout: Timeout_Active, // This is the time that the client will execute to check the status of the request
		Timeout:    timeout,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isPIInstanceRefreshFuncOff(client *instance.IBMPIInstanceClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {

		log.Printf("Calling the check Refresh status of the pvm instance %s", id)
		pvm, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}
		if strings.ToLower(*pvm.Status) == State_Shutoff && pvm.Health.Status == OK {
			return pvm, State_Shutoff, nil
		}
		return pvm, State_Stopping, nil
	}
}

func stopLparForResourceChange(ctx context.Context, client *instance.IBMPIInstanceClient, id string, d *schema.ResourceData) error {
	body := &models.PVMInstanceAction{
		//Action: flex.PtrToString("stop"),
		Action: flex.PtrToString(Action_ImmediateShutdown),
	}
	err := client.Action(id, body)
	if err != nil {
		return fmt.Errorf("failed to perform the stop action on the pvm instance %v", err)
	}

	_, err = isWaitForPIInstanceStopped(ctx, client, id, d.Timeout(schema.TimeoutUpdate))

	return err
}

// Start the lpar
func startLparAfterResourceChange(ctx context.Context, client *instance.IBMPIInstanceClient, id string, d *schema.ResourceData) error {
	body := &models.PVMInstanceAction{
		Action: flex.PtrToString(Action_Start),
	}
	err := client.Action(id, body)
	if err != nil {
		return fmt.Errorf("failed to perform the start action on the pvm instance %v", err)
	}

	_, err = isWaitForPIInstanceAvailable(ctx, client, id, OK, d.Timeout(schema.TimeoutUpdate))

	return err
}

// Stop / Modify / Start only when the lpar is off limits
func performChangeAndReboot(ctx context.Context, client *instance.IBMPIInstanceClient, d *schema.ResourceData, id string, mem, procs float64) error {
	/*
		These are the steps
		1. Stop the lpar - Check if the lpar is SHUTOFF
		2. Once the lpar is SHUTOFF - Make the cpu / memory change - DUring this time , you can check for RESIZE and VERIFY_RESIZE as the transition states
		3. If the change is successful , the lpar state will be back in SHUTOFF
		4. Once the LPAR state is SHUTOFF , initiate the start again and check for ACTIVE + OK
	*/
	//Execute the stop

	log.Printf("Calling the stop lpar for Resource Change code ..")
	err := stopLparForResourceChange(ctx, client, id, d)
	if err != nil {
		return err
	}

	body := &models.PVMInstanceUpdate{
		Memory:     mem,
		Processors: procs,
	}

	_, updateErr := client.Update(id, body)
	if updateErr != nil {
		return fmt.Errorf("failed to update the lpar with the change, %s", updateErr)
	}

	_, err = isWaitForPIInstanceShutoffAfterUpdate(ctx, client, id, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return fmt.Errorf("failed to get an update from the Service after the resource change, %s", err)
	}

	// Now we can start the lpar
	log.Printf("Calling the start lpar After the  Resource Change code ..")
	err = startLparAfterResourceChange(ctx, client, id, d)
	if err != nil {
		return err
	}

	return nil

}

func isWaitForPIInstanceShutoffAfterUpdate(ctx context.Context, client *instance.IBMPIInstanceClient, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for PIInstance (%s) to be ACTIVE or SHUTOFF AFTER THE RESIZE Due to DLPAR Operation ", id)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Resize, State_VerifyResize},
		Target:     []string{State_Active, State_Shutoff, OK},
		Refresh:    isPIInstanceShutAfterResourceChange(client, id),
		Delay:      Timeout_Delay,
		MinTimeout: 5 * time.Minute,
		Timeout:    timeout,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isPIInstanceShutAfterResourceChange(client *instance.IBMPIInstanceClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {

		pvm, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}

		if strings.ToLower(*pvm.Status) == State_Shutoff && pvm.Health.Status == OK {
			log.Printf("The lpar is now off after the resource change...")
			return pvm, State_Shutoff, nil
		}

		return pvm, State_Resize, nil
	}
}

func expandPVMNetworks(networks []interface{}) []*models.PVMInstanceAddNetwork {
	pvmNetworks := make([]*models.PVMInstanceAddNetwork, 0, len(networks))
	for _, v := range networks {
		network := v.(map[string]interface{})
		pvmInstanceNetwork := &models.PVMInstanceAddNetwork{
			IPAddress:               network[Attr_IPAddress].(string),
			NetworkID:               flex.PtrToString(network[Attr_NetworkID].(string)),
			NetworkSecurityGroupIDs: flex.ExpandStringList((network[Attr_NetworkSecurityGroupIDs].(*schema.Set)).List()),
		}
		pvmNetworks = append(pvmNetworks, pvmInstanceNetwork)
	}
	return pvmNetworks
}

func checkCloudInstanceCapability(cloudInstance *models.CloudInstance, custom_capability string) bool {
	log.Printf("Checking for the following capability %s", custom_capability)
	log.Printf("the instance features are %s", cloudInstance.Capabilities)
	for _, v := range cloudInstance.Capabilities {
		if v == custom_capability {
			return true
		}
	}
	return false
}

func createSAPInstance(d *schema.ResourceData, sapClient *instance.IBMPISAPInstanceClient) (*models.PVMInstanceList, error) {
	name := d.Get(Arg_InstanceName).(string)
	profileID := d.Get(Arg_SAPProfileID).(string)
	imageid := d.Get(Arg_ImageID).(string)

	pvmNetworks := expandPVMNetworks(d.Get(Arg_Network).([]interface{}))

	var replicants int64
	if r, ok := d.GetOk(Arg_Replicants); ok {
		replicants = int64(r.(int))
	}
	var replicationpolicy string
	if r, ok := d.GetOk(Arg_ReplicationPolicy); ok {
		replicationpolicy = r.(string)
	}
	var replicationNamingScheme string
	if r, ok := d.GetOk(Arg_ReplicationScheme); ok {
		replicationNamingScheme = r.(string)
	}
	instances := &models.PVMInstanceMultiCreate{
		AffinityPolicy: &replicationpolicy,
		Count:          replicants,
		Numerical:      &replicationNamingScheme,
	}

	body := &models.SAPCreate{
		ImageID:   &imageid,
		Instances: instances,
		Name:      &name,
		Networks:  pvmNetworks,
		ProfileID: &profileID,
	}

	if v, ok := d.GetOk(Arg_SAPDeploymentType); ok {
		body.DeploymentType = v.(string)
	}
	if v, ok := d.GetOk(Arg_VolumeIDs); ok {
		volids := flex.ExpandStringList((v.(*schema.Set)).List())
		if len(volids) > 0 {
			body.VolumeIDs = volids
		}
	}
	if p, ok := d.GetOk(Arg_PinPolicy); ok {
		pinpolicy := p.(string)
		if d.Get(Arg_PinPolicy) == Soft || d.Get(Arg_PinPolicy) == Hard {
			body.PinPolicy = models.PinPolicy(pinpolicy)
		}
	}

	if v, ok := d.GetOk(Arg_KeyPairName); ok {
		sshkey := v.(string)
		body.SSHKeyName = sshkey
	}
	if u, ok := d.GetOk(Arg_UserData); ok {
		userData := u.(string)
		body.UserData = encodeBase64(userData)
	}
	if sys, ok := d.GetOk(Arg_SysType); ok {
		body.SysType = sys.(string)
	}

	if st, ok := d.GetOk(Arg_StorageType); ok {
		body.StorageType = st.(string)
	}
	var bootVolumeReplicationEnabled bool
	if bootVolumeReplicationBoolean, ok := d.GetOk(Arg_BootVolumeReplicationEnabled); ok {
		bootVolumeReplicationEnabled = bootVolumeReplicationBoolean.(bool)
		body.BootVolumeReplicationEnabled = &bootVolumeReplicationEnabled
	}
	var replicationSites []string
	if sites, ok := d.GetOk(Arg_ReplicationSites); ok {
		if !bootVolumeReplicationEnabled {
			return nil, fmt.Errorf("must set %s to true in order to specify replication sites", Arg_BootVolumeReplicationEnabled)
		} else {
			replicationSites = flex.FlattenSet(sites.(*schema.Set))
			body.ReplicationSites = replicationSites
		}
	}
	if sp, ok := d.GetOk(Arg_StoragePool); ok {
		body.StoragePool = sp.(string)
	}

	if ap, ok := d.GetOk(Arg_AffinityPolicy); ok {
		policy := ap.(string)
		affinity := &models.StorageAffinity{
			AffinityPolicy: &policy,
		}

		if policy == Affinity {
			if av, ok := d.GetOk(Arg_AffinityVolume); ok {
				afvol := av.(string)
				affinity.AffinityVolume = &afvol
			}
			if ai, ok := d.GetOk(Arg_AffinityInstance); ok {
				afins := ai.(string)
				affinity.AffinityPVMInstance = &afins
			}
		} else {
			if avs, ok := d.GetOk(Arg_AntiAffinityVolumes); ok {
				afvols := flex.ExpandStringList(avs.([]interface{}))
				affinity.AntiAffinityVolumes = afvols
			}
			if ais, ok := d.GetOk(Arg_AntiAffinityInstances); ok {
				afinss := flex.ExpandStringList(ais.([]interface{}))
				affinity.AntiAffinityPVMInstances = afinss
			}
		}
		body.StorageAffinity = affinity
	}

	if pg, ok := d.GetOk(Arg_PlacementGroupID); ok {
		body.PlacementGroup = pg.(string)
	}
	if deploymentTarget, ok := d.GetOk(Arg_DeploymentTarget); ok {
		body.DeploymentTarget = expandDeploymentTarget(deploymentTarget.(*schema.Set).List())
	}
	if tags, ok := d.GetOk(Arg_UserTags); ok {
		body.UserTags = flex.FlattenSet(tags.(*schema.Set))
	}
	pvmList, err := sapClient.Create(body)
	if err != nil {
		return nil, fmt.Errorf("failed to provision: %v", err)
	}
	if pvmList == nil {
		return nil, fmt.Errorf("failed to provision")
	}

	return pvmList, nil
}

func createPVMInstance(d *schema.ResourceData, client *instance.IBMPIInstanceClient, imageClient *instance.IBMPIImageClient) (*models.PVMInstanceList, error) {

	name := d.Get(Arg_InstanceName).(string)
	imageid := d.Get(Arg_ImageID).(string)

	var mem, procs float64
	var systype, processortype string
	if v, ok := d.GetOk(Arg_Memory); ok {
		mem = v.(float64)
	} else {
		return nil, fmt.Errorf("%s is required for creating pvm instances", Arg_Memory)
	}
	if v, ok := d.GetOk(Arg_Processors); ok {
		procs = v.(float64)
	} else {
		return nil, fmt.Errorf("%s is required for creating pvm instances", Arg_Processors)
	}
	if v, ok := d.GetOk(Arg_SysType); ok {
		systype = v.(string)
	} else {
		return nil, fmt.Errorf("%s is required for creating pvm instances", Arg_SysType)
	}
	if v, ok := d.GetOk(Arg_ProcType); ok {
		processortype = v.(string)
	} else {
		return nil, fmt.Errorf("%s is required for creating pvm instances", Arg_ProcType)
	}

	pvmNetworks := expandPVMNetworks(d.Get(Arg_Network).([]interface{}))

	var volids []string
	if v, ok := d.GetOk(Arg_VolumeIDs); ok {
		volids = flex.ExpandStringList((v.(*schema.Set)).List())
	}
	var replicants float64
	if r, ok := d.GetOk(Arg_Replicants); ok {
		replicants = float64(r.(int))
	}
	var replicationpolicy string
	if r, ok := d.GetOk(Arg_ReplicationPolicy); ok {
		replicationpolicy = r.(string)
	}
	var replicationNamingScheme string
	if r, ok := d.GetOk(Arg_ReplicationScheme); ok {
		replicationNamingScheme = r.(string)
	}
	var pinpolicy string
	if p, ok := d.GetOk(Arg_PinPolicy); ok {
		pinpolicy = p.(string)
		if pinpolicy == "" {
			pinpolicy = None
		}
	}

	var userData string
	if u, ok := d.GetOk(Arg_UserData); ok {
		userData = u.(string)
	}

	body := &models.PVMInstanceCreate{
		Processors:              &procs,
		Memory:                  &mem,
		ServerName:              flex.PtrToString(name),
		SysType:                 systype,
		ImageID:                 flex.PtrToString(imageid),
		ProcType:                flex.PtrToString(processortype),
		Replicants:              &replicants,
		UserData:                encodeBase64(userData),
		ReplicantNamingScheme:   flex.PtrToString(replicationNamingScheme),
		ReplicantAffinityPolicy: flex.PtrToString(replicationpolicy),
		Networks:                pvmNetworks,
	}
	if s, ok := d.GetOk(Arg_KeyPairName); ok {
		sshkey := s.(string)
		body.KeyPairName = sshkey
	}
	if len(volids) > 0 {
		body.VolumeIDs = volids
	}
	if d.Get(Arg_PinPolicy) == Soft || d.Get(Arg_PinPolicy) == Hard {
		body.PinPolicy = models.PinPolicy(pinpolicy)
	}

	var assignedVirtualCores int64
	if a, ok := d.GetOk(Arg_VirtualCoresAssigned); ok {
		assignedVirtualCores = int64(a.(int))
		body.VirtualCores = &models.VirtualCores{Assigned: &assignedVirtualCores}
	}

	if st, ok := d.GetOk(Arg_StorageType); ok {
		body.StorageType = st.(string)
	}
	if sp, ok := d.GetOk(Arg_StoragePool); ok {
		body.StoragePool = sp.(string)
	}

	if dt, ok := d.GetOk(Arg_DeploymentType); ok {
		body.DeploymentType = dt.(string)
	}

	if ap, ok := d.GetOk(Arg_AffinityPolicy); ok {
		policy := ap.(string)
		affinity := &models.StorageAffinity{
			AffinityPolicy: &policy,
		}

		if policy == Affinity {
			if av, ok := d.GetOk(Arg_AffinityVolume); ok {
				afvol := av.(string)
				affinity.AffinityVolume = &afvol
			}
			if ai, ok := d.GetOk(Arg_AffinityInstance); ok {
				afins := ai.(string)
				affinity.AffinityPVMInstance = &afins
			}
		} else {
			if avs, ok := d.GetOk(Arg_AntiAffinityVolumes); ok {
				afvols := flex.ExpandStringList(avs.([]interface{}))
				affinity.AntiAffinityVolumes = afvols
			}
			if ais, ok := d.GetOk(Arg_AntiAffinityInstances); ok {
				afinss := flex.ExpandStringList(ais.([]interface{}))
				affinity.AntiAffinityPVMInstances = afinss
			}
		}
		body.StorageAffinity = affinity
	}

	if sc, ok := d.GetOk(Arg_StorageConnection); ok {
		body.StorageConnection = sc.(string)
	}

	if pg, ok := d.GetOk(Arg_PlacementGroupID); ok {
		body.PlacementGroup = pg.(string)
	}

	if spp, ok := d.GetOk(Arg_SharedProcessorPool); ok {
		body.SharedProcessorPool = spp.(string)
	}
	imageData, err := imageClient.GetStockImage(imageid)
	if err != nil {
		// check if vtl image is cloud instance image
		imageData, err = imageClient.Get(imageid)
		if err != nil {
			return nil, fmt.Errorf("image doesn't exist. %e", err)
		}
	}
	if lrc, ok := d.GetOk(Arg_LicenseRepositoryCapacity); ok {

		if imageData.Specifications.ImageType == StockVTL {
			body.LicenseRepositoryCapacity = int64(lrc.(int))
		} else {
			return nil, fmt.Errorf("pi_license_repository_capacity should only be used when creating VTL instances. %e", err)
		}
	}

	if imageData.Specifications.OperatingSystem == OS_IBMI {
		// Default value
		falseBool := false
		sl := &models.SoftwareLicenses{
			IbmiCSS:      &falseBool,
			IbmiPHA:      &falseBool,
			IbmiRDS:      &falseBool,
			IbmiRDSUsers: 0,
		}
		if ibmiCSS, ok := d.GetOk(Arg_IBMiCSS); ok {
			sl.IbmiCSS = flex.PtrToBool(ibmiCSS.(bool))
		}
		if ibmiPHA, ok := d.GetOk(Arg_IBMiPHA); ok {
			sl.IbmiPHA = flex.PtrToBool(ibmiPHA.(bool))
		}
		if ibmrdsUsers, ok := d.GetOk(Arg_IBMiRDSUsers); ok {
			if ibmrdsUsers.(int) < 0 {
				return nil, fmt.Errorf("request with IBM i Rational Dev Studio property requires IBM i Rational Dev Studio number of users")
			}
			sl.IbmiRDS = flex.PtrToBool(ibmrdsUsers.(int) > 0)
			sl.IbmiRDSUsers = int64(ibmrdsUsers.(int))
		}
		body.SoftwareLicenses = sl
	}
	if deploymentTarget, ok := d.GetOk(Arg_DeploymentTarget); ok {
		body.DeploymentTarget = expandDeploymentTarget(deploymentTarget.(*schema.Set).List())
	}
	var bootVolumeReplicationEnabled bool
	if bootVolumeReplicationBoolean, ok := d.GetOk(Arg_BootVolumeReplicationEnabled); ok {
		bootVolumeReplicationEnabled = bootVolumeReplicationBoolean.(bool)
		body.BootVolumeReplicationEnabled = &bootVolumeReplicationEnabled
	}
	var replicationSites []string
	if sites, ok := d.GetOk(Arg_ReplicationSites); ok {
		if !bootVolumeReplicationEnabled {
			return nil, fmt.Errorf("must set %s to true in order to specify replication sites", Arg_BootVolumeReplicationEnabled)
		} else {
			replicationSites = flex.FlattenSet(sites.(*schema.Set))
			body.ReplicationSites = replicationSites
		}
	}

	if tags, ok := d.GetOk(Arg_UserTags); ok {
		body.UserTags = flex.FlattenSet(tags.(*schema.Set))
	}
	if vsn, ok := d.GetOk(Arg_VirtualSerialNumber); ok {
		vsnListType := vsn.([]interface{})
		vsnCreateModel := vsnSetToCreateModel(vsnListType)
		body.VirtualSerialNumber = vsnCreateModel
	}

	pvmList, err := client.Create(body)

	if err != nil {
		return nil, fmt.Errorf("failed to provision: %v", err)
	}
	if pvmList == nil {
		return nil, fmt.Errorf("failed to provision")
	}

	return pvmList, nil
}

func expandDeploymentTarget(dt []interface{}) *models.DeploymentTarget {
	dtexpanded := &models.DeploymentTarget{}
	for _, v := range dt {
		dtarget := v.(map[string]interface{})
		dtexpanded.ID = core.StringPtr(dtarget[Attr_ID].(string))
		dtexpanded.Type = core.StringPtr(dtarget[Attr_Type].(string))
	}
	return dtexpanded
}

func splitID(id string) (id1, id2 string, err error) {
	parts, err := flex.IdParts(id)
	if err != nil {
		return
	}
	id1 = parts[0]
	id2 = parts[1]
	return
}

func vsnSetToCreateModel(vsnSetList []interface{}) *models.CreateServerVirtualSerialNumber {
	vsnItemMap := vsnSetList[0].(map[string]interface{})
	serialString := vsnItemMap[Attr_Serial].(string)
	model := &models.CreateServerVirtualSerialNumber{
		Serial: &serialString,
	}
	description := vsnItemMap[Attr_Description].(string)
	if description != "" {
		model.Description = description
	}

	return model
}

func flattenVirtualSerialNumberToList(vsn *models.GetServerVirtualSerialNumber) []map[string]interface{} {
	v := make([]map[string]interface{}, 1)
	v[0] = map[string]interface{}{
		Attr_Description: vsn.Description,
		Attr_Serial:      vsn.Serial,
	}
	return v
}

// Do not show a diff if VSN is changed to existing assigned VSN
func supressVSNDiffAutoAssign(k, old, new string, d *schema.ResourceData) bool {
	return new == old || (new == AutoAssign && old != "")
}
