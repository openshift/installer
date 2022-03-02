package vsphere

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/provider"

	"github.com/vmware/govmomi/vim25/mo"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/clustercomputeresource"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/customattribute"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/folder"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/hostsystem"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/viapi"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/vsansystem"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

const resourceVSphereComputeClusterName = "vsphere_compute_cluster"

const (
	clusterAdmissionControlTypeResourcePercentage = "resourcePercentage"
	clusterAdmissionControlTypeSlotPolicy         = "slotPolicy"
	clusterAdmissionControlTypeFailoverHosts      = "failoverHosts"
	clusterAdmissionControlTypeDisabled           = "disabled"
)

var clusterAdmissionControlTypeAllowedValues = []string{
	clusterAdmissionControlTypeResourcePercentage,
	clusterAdmissionControlTypeSlotPolicy,
	clusterAdmissionControlTypeFailoverHosts,
	clusterAdmissionControlTypeDisabled,
}

var drsBehaviorAllowedValues = []string{
	string(types.DrsBehaviorManual),
	string(types.DrsBehaviorPartiallyAutomated),
	string(types.DrsBehaviorFullyAutomated),
}

var dpmBehaviorAllowedValues = []string{
	string(types.DpmBehaviorManual),
	string(types.DpmBehaviorAutomated),
}

var clusterDasConfigInfoServiceStateAllowedValues = []string{
	string(types.ClusterDasConfigInfoServiceStateEnabled),
	string(types.ClusterDasConfigInfoServiceStateDisabled),
}

var computeClusterDasConfigInfoServiceStateAllowedValues = []string{
	string(types.ClusterDasVmSettingsRestartPriorityLowest),
	string(types.ClusterDasVmSettingsRestartPriorityLow),
	string(types.ClusterDasVmSettingsRestartPriorityMedium),
	string(types.ClusterDasVmSettingsRestartPriorityHigh),
	string(types.ClusterDasVmSettingsRestartPriorityHighest),
}

var computeClusterVMReadinessReadyConditionAllowedValues = []string{
	string(types.ClusterVmReadinessReadyConditionNone),
	string(types.ClusterVmReadinessReadyConditionPoweredOn),
	string(types.ClusterVmReadinessReadyConditionGuestHbStatusGreen),
	string(types.ClusterVmReadinessReadyConditionAppHbStatusGreen),
}

var computeClusterDasVMSettingsIsolationResponseAllowedValues = []string{
	string(types.ClusterDasVmSettingsIsolationResponseNone),
	string(types.ClusterDasVmSettingsIsolationResponsePowerOff),
	string(types.ClusterDasVmSettingsIsolationResponseShutdown),
}

var computeClusterVMStorageProtectionForPDLAllowedValues = []string{
	string(types.ClusterVmComponentProtectionSettingsStorageVmReactionDisabled),
	string(types.ClusterVmComponentProtectionSettingsStorageVmReactionWarning),
	string(types.ClusterVmComponentProtectionSettingsStorageVmReactionRestartAggressive),
}

var computeClusterVMStorageProtectionForAPDAllowedValues = []string{
	string(types.ClusterVmComponentProtectionSettingsStorageVmReactionDisabled),
	string(types.ClusterVmComponentProtectionSettingsStorageVmReactionWarning),
	string(types.ClusterVmComponentProtectionSettingsStorageVmReactionRestartConservative),
	string(types.ClusterVmComponentProtectionSettingsStorageVmReactionRestartAggressive),
}

var computeClusterVMReactionOnAPDClearedAllowedValues = []string{
	string(types.ClusterVmComponentProtectionSettingsVmReactionOnAPDClearedNone),
	string(types.ClusterVmComponentProtectionSettingsVmReactionOnAPDClearedReset),
}

var clusterDasConfigInfoVMMonitoringStateAllowedValues = []string{
	string(types.ClusterDasConfigInfoVmMonitoringStateVmMonitoringDisabled),
	string(types.ClusterDasConfigInfoVmMonitoringStateVmMonitoringOnly),
	string(types.ClusterDasConfigInfoVmMonitoringStateVmAndAppMonitoring),
}

var clusterDasConfigInfoHBDatastoreCandidateAllowedValues = []string{
	string(types.ClusterDasConfigInfoHBDatastoreCandidateUserSelectedDs),
	string(types.ClusterDasConfigInfoHBDatastoreCandidateAllFeasibleDs),
	string(types.ClusterDasConfigInfoHBDatastoreCandidateAllFeasibleDsWithUserPreference),
}

var clusterInfraUpdateHaConfigInfoBehaviorTypeAllowedValues = []string{
	string(types.ClusterInfraUpdateHaConfigInfoBehaviorTypeManual),
	string(types.ClusterInfraUpdateHaConfigInfoBehaviorTypeAutomated),
}

var clusterInfraUpdateHaConfigInfoRemediationTypeAllowedValues = []string{
	string(types.ClusterInfraUpdateHaConfigInfoRemediationTypeMaintenanceMode),
	string(types.ClusterInfraUpdateHaConfigInfoRemediationTypeQuarantineMode),
}

func resourceVSphereComputeCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceVSphereComputeClusterCreate,
		Read:   resourceVSphereComputeClusterRead,
		Update: resourceVSphereComputeClusterUpdate,
		Delete: resourceVSphereComputeClusterDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVSphereComputeClusterImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name for the new cluster.",
			},
			"datacenter_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The managed object ID of the datacenter to put the cluster in.",
			},
			"host_system_ids": {
				Type:          schema.TypeSet,
				Optional:      true,
				MaxItems:      64,
				Description:   "The managed object IDs of the hosts to put in the cluster.",
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"host_managed"},
			},
			"folder": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the folder to locate the cluster in.",
				StateFunc:   folder.NormalizePath,
			},
			"host_cluster_exit_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3600,
				Description:  "The timeout for each host maintenance mode operation when removing hosts from a cluster.",
				ValidateFunc: validation.IntBetween(0, 604800),
			},
			"force_evacuate_on_destroy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Force removal of all hosts in the cluster during destroy and make them standalone hosts. Use of this flag mainly exists for testing and is not recommended in normal use.",
			},
			// DRS - General/automation
			"drs_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable DRS for this cluster.",
			},
			"drs_automation_level": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      string(types.DrsBehaviorManual),
				Description:  "The default automation level for all virtual machines in this cluster. Can be one of manual, partiallyAutomated, or fullyAutomated.",
				ValidateFunc: validation.StringInSlice(drsBehaviorAllowedValues, false),
			},
			"drs_migration_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				Description:  "A value between 1 and 5 indicating the threshold of imbalance tolerated between hosts. A lower setting will tolerate more imbalance while a higher setting will tolerate less.",
				ValidateFunc: validation.IntBetween(1, 5),
			},
			"drs_enable_vm_overrides": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "When true, allows individual VM overrides within this cluster to be set.",
			},
			"drs_enable_predictive_drs": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "When true, enables DRS to use data from vRealize Operations Manager to make proactive DRS recommendations.",
			},
			// DRS - DPM
			"dpm_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable DPM support for DRS. This allows you to dynamically control the power of hosts depending on the needs of virtual machines in the cluster. Requires that DRS be enabled.",
			},
			"dpm_automation_level": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      string(types.DpmBehaviorManual),
				Description:  "The automation level for host power operations in this cluster. Can be one of manual or automated.",
				ValidateFunc: validation.StringInSlice(dpmBehaviorAllowedValues, false),
			},
			"dpm_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				Description:  "A value between 1 and 5 indicating the threshold of load within the cluster that influences host power operations. This affects both power on and power off operations - a lower setting will tolerate more of a surplus/deficit than a higher setting.",
				ValidateFunc: validation.IntBetween(1, 5),
			},
			// DRS - Advanced options
			"drs_advanced_options": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Advanced configuration options for DRS and DPM.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			// HA - General
			"ha_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable vSphere HA for this cluster.",
			},
			// HA - Host monitoring settings
			"ha_host_monitoring": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      string(types.ClusterDasConfigInfoServiceStateEnabled),
				Description:  "Global setting that controls whether vSphere HA remediates VMs on host failure. Can be one of enabled or disabled.",
				ValidateFunc: validation.StringInSlice(clusterDasConfigInfoServiceStateAllowedValues, false),
			},
			// Host monitoring - VM restarts
			"ha_vm_restart_priority": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      string(types.ClusterDasVmSettingsRestartPriorityMedium),
				Description:  "The default restart priority for affected VMs when vSphere detects a host failure. Can be one of lowest, low, medium, high, or highest.",
				ValidateFunc: validation.StringInSlice(computeClusterDasConfigInfoServiceStateAllowedValues, false),
			},
			"ha_vm_dependency_restart_condition": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      string(types.ClusterVmReadinessReadyConditionNone),
				Description:  "The condition used to determine whether or not VMs in a certain restart priority class are online, allowing HA to move on to restarting VMs on the next priority. Can be one of none, poweredOn, guestHbStatusGreen, or appHbStatusGreen.",
				ValidateFunc: validation.StringInSlice(computeClusterVMReadinessReadyConditionAllowedValues, false),
			},
			"ha_vm_restart_additional_delay": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Additional delay in seconds after ready condition is met. A VM is considered ready at this point.",
			},
			"ha_vm_restart_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     600,
				Description: "The maximum time, in seconds, that vSphere HA will wait for virtual machines in one priority to be ready before proceeding with the next priority.",
			},
			// Host monitoring - host isolation
			"ha_host_isolation_response": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      string(types.ClusterDasVmSettingsIsolationResponseNone),
				Description:  "The action to take on virtual machines when a host has detected that it has been isolated from the rest of the cluster. Can be one of none, powerOff, or shutdown.",
				ValidateFunc: validation.StringInSlice(computeClusterDasVMSettingsIsolationResponseAllowedValues, false),
			},
			// VM component protection
			"ha_vm_component_protection": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      string(types.ClusterDasConfigInfoServiceStateEnabled),
				Description:  "Controls vSphere VM component protection for virtual machines in this cluster. This allows vSphere HA to react to failures between hosts and specific virtual machine components, such as datastores. Can be one of enabled or disabled.",
				ValidateFunc: validation.StringInSlice(clusterDasConfigInfoServiceStateAllowedValues, false),
			},
			// VM component protection - datastore monitoring - Permanent Device Loss
			"ha_datastore_pdl_response": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      string(types.ClusterVmComponentProtectionSettingsStorageVmReactionDisabled),
				Description:  "When ha_vm_component_protection is enabled, controls the action to take on virtual machines when the cluster has detected a permanent device loss to a relevant datastore. Can be one of disabled, warning, or restartAggressive.",
				ValidateFunc: validation.StringInSlice(computeClusterVMStorageProtectionForPDLAllowedValues, false),
			},
			// VM component protection - datastore monitoring - All Paths Down
			"ha_datastore_apd_response": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      string(types.ClusterVmComponentProtectionSettingsStorageVmReactionDisabled),
				Description:  "When ha_vm_component_protection is enabled, controls the action to take on virtual machines when the cluster has detected loss to all paths to a relevant datastore. Can be one of disabled, warning, restartConservative, or restartAggressive.",
				ValidateFunc: validation.StringInSlice(computeClusterVMStorageProtectionForAPDAllowedValues, false),
			},
			"ha_datastore_apd_recovery_action": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      string(types.ClusterVmComponentProtectionSettingsVmReactionOnAPDClearedNone),
				Description:  "When ha_vm_component_protection is enabled, controls the action to take on virtual machines if an APD status on an affected datastore clears in the middle of an APD event. Can be one of none or reset.",
				ValidateFunc: validation.StringInSlice(computeClusterVMReactionOnAPDClearedAllowedValues, false),
			},
			"ha_datastore_apd_response_delay": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3,
				Description: "When ha_vm_component_protection is enabled, controls the delay in minutes to wait after an APD timeout event to execute the response action defined in ha_datastore_apd_response.",
			},
			// VM monitoring
			"ha_vm_monitoring": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      string(types.ClusterDasConfigInfoVmMonitoringStateVmMonitoringDisabled),
				Description:  "The type of virtual machine monitoring to use when HA is enabled in the cluster. Can be one of vmMonitoringDisabled, vmMonitoringOnly, or vmAndAppMonitoring.",
				ValidateFunc: validation.StringInSlice(clusterDasConfigInfoVMMonitoringStateAllowedValues, false),
			},
			"ha_vm_failure_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     30,
				Description: "If a heartbeat from a virtual machine is not received within this configured interval, the virtual machine is marked as failed. The value is in seconds.",
			},
			"ha_vm_minimum_uptime": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     120,
				Description: "The time, in seconds, that HA waits after powering on a virtual machine before monitoring for heartbeats.",
			},
			"ha_vm_maximum_resets": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3,
				Description: "The maximum number of resets that HA will perform to a virtual machine when responding to a failure event.",
			},
			"ha_vm_maximum_failure_window": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     -1,
				Description: "The length of the reset window in which ha_vm_maximum_resets can operate. When this window expires, no more resets are attempted regardless of the setting configured in ha_vm_maximum_resets. -1 means no window, meaning an unlimited reset time is allotted.",
			},
			// Admission control
			"ha_admission_control_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      clusterAdmissionControlTypeResourcePercentage,
				Description:  "The type of admission control policy to use with vSphere HA, which controls whether or not specific VM operations are permitted in the cluster in order to protect the reliability of the cluster. Can be one of resourcePercentage, slotPolicy, failoverHosts, or disabled. Note that disabling admission control is not recommended and can lead to service issues.",
				ValidateFunc: validation.StringInSlice(clusterAdmissionControlTypeAllowedValues, false),
			},
			"ha_admission_control_host_failure_tolerance": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "The maximum number of failed hosts that admission control tolerates when making decisions on whether to permit virtual machine operations. The maximum is one less than the number of hosts in the cluster.",
			},
			"ha_admission_control_performance_tolerance": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      100,
				Description:  "The percentage of resource reduction that a cluster of VMs can tolerate in case of a failover. A value of 0 produces warnings only, whereas a value of 100 disables the setting.",
				ValidateFunc: validation.IntBetween(0, 100),
			},
			"ha_admission_control_resource_percentage_auto_compute": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "When ha_admission_control_policy is resourcePercentage, automatically determine available resource percentages by subtracting the average number of host resources represented by the ha_admission_control_host_failure_tolerance setting from the total amount of resources in the cluster. Disable to supply user-defined values.",
			},
			"ha_admission_control_resource_percentage_cpu": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      100,
				Description:  "When ha_admission_control_policy is resourcePercentage, this controls the user-defined percentage of CPU resources in the cluster to reserve for failover.",
				ValidateFunc: validation.IntBetween(1, 100),
			},
			"ha_admission_control_resource_percentage_memory": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      100,
				Description:  "When ha_admission_control_policy is resourcePercentage, this controls the user-defined percentage of memory resources in the cluster to reserve for failover.",
				ValidateFunc: validation.IntBetween(1, 100),
			},
			"ha_admission_control_slot_policy_use_explicit_size": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "When ha_admission_control_policy is slotPolicy, this setting controls whether or not you wish to supply explicit values to CPU and memory slot sizes. The default is to gather a automatic average based on all powered-on virtual machines currently in the cluster.",
			},
			"ha_admission_control_slot_policy_explicit_cpu": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     32,
				Description: "When ha_admission_control_policy is slotPolicy, this controls the user-defined CPU slot size, in MHz.",
			},
			"ha_admission_control_slot_policy_explicit_memory": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
				Description: "When ha_admission_control_policy is slotPolicy, this controls the user-defined memory slot size, in MB.",
			},
			"ha_admission_control_failover_host_system_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "When ha_admission_control_policy is failoverHosts, this defines the managed object IDs of hosts to use as dedicated failover hosts. These hosts are kept as available as possible - admission control will block access to the host, and DRS will ignore the host when making recommendations.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			// HA - datastores
			"ha_heartbeat_datastore_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      string(types.ClusterDasConfigInfoHBDatastoreCandidateAllFeasibleDsWithUserPreference),
				Description:  "The selection policy for HA heartbeat datastores. Can be one of allFeasibleDs, userSelectedDs, or allFeasibleDsWithUserPreference.",
				ValidateFunc: validation.StringInSlice(clusterDasConfigInfoHBDatastoreCandidateAllowedValues, false),
			},
			"ha_heartbeat_datastore_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The list of managed object IDs for preferred datastores to use for HA heartbeating. This setting is only useful when ha_heartbeat_datastore_policy is set to either userSelectedDs or allFeasibleDsWithUserPreference.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			// HA - Advanced options
			"ha_advanced_options": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Advanced configuration options for vSphere HA.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			// Proactive HA
			"proactive_ha_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enables proactive HA, allowing for vSphere to get HA data from external providers and use DRS to perform remediation.",
			},
			"proactive_ha_automation_level": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      string(types.ClusterInfraUpdateHaConfigInfoBehaviorTypeManual),
				Description:  "The DRS behavior for proactive HA recommendations. Can be one of Automated or Manual.",
				ValidateFunc: validation.StringInSlice(clusterInfraUpdateHaConfigInfoBehaviorTypeAllowedValues, false),
			},
			"proactive_ha_moderate_remediation": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      string(types.ClusterInfraUpdateHaConfigInfoRemediationTypeQuarantineMode),
				Description:  "The configured remediation for moderately degraded hosts. Can be one of MaintenanceMode or QuarantineMode. Note that this cannot be set to MaintenanceMode when proactive_ha_severe_remediation is set to QuarantineMode.",
				ValidateFunc: validation.StringInSlice(clusterInfraUpdateHaConfigInfoRemediationTypeAllowedValues, false),
			},
			"proactive_ha_severe_remediation": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      string(types.ClusterInfraUpdateHaConfigInfoRemediationTypeQuarantineMode),
				Description:  "The configured remediation for severely degraded hosts. Can be one of MaintenanceMode or QuarantineMode. Note that this cannot be set to QuarantineMode when proactive_ha_moderate_remediation is set to MaintenanceMode.",
				ValidateFunc: validation.StringInSlice(clusterInfraUpdateHaConfigInfoRemediationTypeAllowedValues, false),
			},
			"proactive_ha_provider_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The list of IDs for health update providers configured for this cluster.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"resource_pool_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The managed object ID of the cluster's root resource pool.",
			},
			"host_managed": {
				Type:          schema.TypeBool,
				Optional:      true,
				Description:   "Must be set if cluster enrollment is managed from host resource.",
				ConflictsWith: []string{"host_system_ids"},
			},
			"vsan_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether the VSAN service is enabled for the cluster.",
			},
			"vsan_disk_group": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "A list of disk UUIDs to add to the vSAN cluster.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cache": {
							Type:        schema.TypeString,
							Description: "Cache disk.",
							Optional:    true,
						},
						"storage": {
							Type:        schema.TypeSet,
							Description: "List of storage disks.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			vSphereTagAttributeKey:    tagsSchema(),
			customattribute.ConfigKey: customattribute.ConfigSchema(),
		},
	}
}

func resourceVSphereComputeClusterCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning create", resourceVSphereComputeClusterIDString(d))

	// We create the cluster here. This function creates a cluster with no
	// configuration, as we want to add the hosts before applying the full
	// configuration.
	cluster, err := resourceVSphereComputeClusterApplyCreate(d, meta)
	if err != nil {
		return err
	}

	// The cluster can be tagged here now.
	if err := resourceVSphereComputeClusterApplyTags(d, meta, cluster); err != nil {
		return err
	}
	if err := resourceVSphereComputeClusterApplyCustomAttributes(d, meta, cluster); err != nil {
		return err
	}

	// Move the hosts in now.
	if err := resourceVSphereComputeClusterProcessHostUpdate(d, meta, cluster); err != nil {
		return err
	}

	// Now that all the hosts that will be in the cluster have been added, apply
	// the cluster configuration.
	if err := resourceVSphereComputeClusterApplyClusterConfiguration(d, meta, cluster); err != nil {
		return err
	}

	// All done!
	log.Printf("[DEBUG] %s: Create finished successfully", resourceVSphereComputeClusterIDString(d))
	return resourceVSphereComputeClusterRead(d, meta)
}

func resourceVSphereComputeClusterRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning read", resourceVSphereComputeClusterIDString(d))

	cluster, err := resourceVSphereComputeClusterGetCluster(d, meta)
	if err != nil {
		if viapi.IsManagedObjectNotFoundError(err) {
			log.Printf("[DEBUG] %s: Resource has been deleted", resourceVSphereComputeClusterIDString(d))
			d.SetId("")
			return nil
		}
		return err
	}

	if err := resourceVSphereComputeClusterSaveDatacenter(d, meta, cluster); err != nil {
		return err
	}

	if err := resourceVSphereComputeClusterSaveNameAndPath(d, cluster); err != nil {
		return err
	}

	if err := resourceVSphereComputeClusterFlattenData(d, meta, cluster); err != nil {
		return err
	}

	if err := resourceVSphereComputeClusterReadTags(d, meta, cluster); err != nil {
		return err
	}

	if err := resourceVSphereComputeClusterReadCustomAttributes(d, meta, cluster); err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Read completed successfully", resourceVSphereComputeClusterIDString(d))
	return nil
}

func resourceVSphereComputeClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning update", resourceVSphereComputeClusterIDString(d))

	cluster, err := resourceVSphereComputeClusterGetCluster(d, meta)
	if err != nil {
		return err
	}

	cluster, err = resourceVSphereComputeClusterApplyNameChange(d, meta, cluster)
	if err != nil {
		return err
	}
	cluster, err = resourceVSphereComputeClusterApplyFolderChange(d, meta, cluster)
	if err != nil {
		return err
	}

	if err := resourceVSphereComputeClusterProcessHostUpdate(d, meta, cluster); err != nil {
		return err
	}

	if err := resourceVSphereComputeClusterApplyClusterConfiguration(d, meta, cluster); err != nil {
		return err
	}

	if err := resourceVSphereComputeClusterApplyTags(d, meta, cluster); err != nil {
		return err
	}

	if err := resourceVSphereComputeClusterApplyCustomAttributes(d, meta, cluster); err != nil {
		return err
	}

	if err = updateVsanDisks(d, cluster, meta); err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Update finished successfully", resourceVSphereComputeClusterIDString(d))
	return resourceVSphereComputeClusterRead(d, meta)
}

func resourceVSphereComputeClusterDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning delete", resourceVSphereComputeClusterIDString(d))
	cluster, err := resourceVSphereComputeClusterGetCluster(d, meta)
	if err != nil {
		return err
	}

	client, err := resourceVSphereComputeClusterClient(meta)
	if err != nil {
		return err
	}

	version := viapi.ParseVersionFromClient(client)
	spec := expandClusterConfigSpecEx(d, version)

	if *spec.DasConfig.Enabled && *spec.DasConfig.AdmissionControlEnabled {
		switch v := spec.DasConfig.AdmissionControlPolicy.(type) {
		case *types.ClusterFailoverHostAdmissionControlPolicy:
			_ = v
			log.Printf("[DEBUG] if Admission Control Policy set to Failover Host than turn HA OFF before removing hosts")
			spec.DasConfig.Enabled = structure.BoolPtr(false)
			if err := clustercomputeresource.Reconfigure(cluster, spec); err != nil {
				return err
			}
		}
	}

	if err := resourceVSphereComputeClusterDeleteProcessForceRemoveHosts(d, meta, cluster); err != nil {
		return err
	}

	if err := resourceVSphereComputeClusterValidateEmptyCluster(d, cluster); err != nil {
		return err
	}

	if err := resourceVSphereComputeClusterApplyDelete(d, cluster); err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Deleted successfully", resourceVSphereComputeClusterIDString(d))
	return nil
}

func resourceVSphereComputeClusterImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	p := d.Id()
	cluster, err := resourceVSphereComputeClusterGetClusterFromPath(meta, p, "")
	if err != nil {
		return nil, fmt.Errorf("error loading cluster: %s", err)
	}

	d.SetId(cluster.Reference().Value)
	if err := resourceVSphereComputeClusterImportSetDefaults(d); err != nil {
		return nil, err
	}

	err = resourceVSphereComputeClusterRead(d, meta)
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

// resourceVSphereComputeClusterImportSetDefaults sets configuration defaults
// for resource attributes that are not managed by read. This gives
// completeness for import, in addition to reducing the number of post-import
// diffs that may need to happen.
func resourceVSphereComputeClusterImportSetDefaults(d *schema.ResourceData) error {
	s := resourceVSphereComputeCluster().Schema
	return structure.SetBatch(d, map[string]interface{}{
		"ha_admission_control_performance_tolerance":       s["ha_admission_control_performance_tolerance"].Default,
		"ha_admission_control_resource_percentage_cpu":     s["ha_admission_control_resource_percentage_cpu"].Default,
		"ha_admission_control_resource_percentage_memory":  s["ha_admission_control_resource_percentage_memory"].Default,
		"ha_admission_control_slot_policy_explicit_cpu":    s["ha_admission_control_slot_policy_explicit_cpu"].Default,
		"ha_admission_control_slot_policy_explicit_memory": s["ha_admission_control_slot_policy_explicit_memory"].Default,
		"host_cluster_exit_timeout":                        s["host_cluster_exit_timeout"].Default,
	})
}

// resourceVSphereComputeClusterApplyCreate processes the creation part of
// resourceVSphereComputeClusterCreate.
func resourceVSphereComputeClusterApplyCreate(d *schema.ResourceData, meta interface{}) (*object.ClusterComputeResource, error) {
	log.Printf("[DEBUG] %s: Processing compute cluster creation", resourceVSphereComputeClusterIDString(d))
	client, err := resourceVSphereComputeClusterClient(meta)
	if err != nil {
		return nil, err
	}

	dc, err := datacenterFromID(client, d.Get("datacenter_id").(string))
	if err != nil {
		return nil, fmt.Errorf("cannot locate datacenter: %s", err)
	}

	// Find the folder based off the path to the datacenter. This is where we
	// create the datastore cluster.
	f, err := folder.FromPath(client, d.Get("folder").(string), folder.VSphereFolderTypeHost, dc)
	if err != nil {
		return nil, fmt.Errorf("cannot locate folder: %s", err)
	}

	// Create the cluster. We use an empty config spec so that we can move the
	// necessary hosts into the cluster *before* we send the full configuration,
	// ensuring that any host-dependent configuration does not break.
	cluster, err := clustercomputeresource.Create(f, d.Get("name").(string), types.ClusterConfigSpecEx{})
	if err != nil {
		return nil, fmt.Errorf("error creating cluster: %s", err)
	}
	if err = updateVsanDisks(d, cluster, meta); err != nil {
		return nil, err
	}
	// Set the ID now before proceeding any further. Any other operation past
	// this point is recoverable.
	d.SetId(cluster.Reference().Value)

	return cluster, nil
}

// resourceVSphereComputeClusterProcessHostUpdate processes any changes in host
// membership in the cluster.
//
// Note that this has implications for interoperability with any future host
// resources that we may set up in Terraform. This namely exists to support the
// fact that some cluster configuration settings depend on hosts actually
// existing in the cluster before they can be carried out, in addition to the
// fact that we don't have any actual host resources at this point in time, and
// may actually not in the future as the addition of hosts will require
// passwords to be supplied to Terraform, which will propagate to state and
// have security implications.
//
// Currently, this process expects the hosts supplied to host_system_ids to be
// already added to vSphere - the recommended method would be to add these
// hosts OOB as standalone hosts to the datacenter that the cluster is being
// deployed to, and then use the vsphere_host data source to get the necessary
// ID to pass into the vsphere_compute_cluster resource.
//
// Hosts moved *out* of the cluster will be moved to the root host folder of
// the datacenter the cluster belongs to. This will create a ComputeResource MO
// for this host OOB from Terraform. Conversely, moving a host into a cluster
// removes the ComputeResource MO for that host, in addition to moving any VMs
// into the cluster at the root cluster resource pool, removing any resource
// pools that exist on the standalone host.
//
// Hosts being removed are placed into maintenance mode. It is up to the
// operator to determine what the implications of this are - if DRS is set up
// correctly and sufficient resources exist, placing a host into maintenance
// mode *should* migrate powered on VMs off the cluster. Powered off VMs will
// be migrated as well, leaving the host as empty as possible after it leaves
// the cluster. The host will be taken out of maintenance mode after being
// removed.
func resourceVSphereComputeClusterProcessHostUpdate(
	d *schema.ResourceData,
	meta interface{},
	cluster *object.ClusterComputeResource,
) error {
	log.Printf("[DEBUG] %s: Processing any necessary host addition/removal operations", resourceVSphereComputeClusterIDString(d))
	client, err := resourceVSphereComputeClusterClient(meta)
	if err != nil {
		return err
	}

	o, n := d.GetChange("host_system_ids")

	newHosts, err := resourceVSphereComputeClusterGetHostSystemObjects(
		client,
		structure.SliceInterfacesToStrings(n.(*schema.Set).Difference(o.(*schema.Set)).List()),
	)
	if err != nil {
		return err
	}

	oldHosts, err := resourceVSphereComputeClusterGetHostSystemObjects(
		client,
		structure.SliceInterfacesToStrings(o.(*schema.Set).Difference(n.(*schema.Set)).List()),
	)
	if err != nil {
		return err
	}

	// Add new hosts first
	if len(newHosts) > 0 {
		if err := clustercomputeresource.MoveHostsInto(client, cluster, newHosts); err != nil {
			return fmt.Errorf("error moving new hosts into cluster: %s", err)
		}

		for _, hs := range newHosts {
			hsProps, err := hostsystem.Properties(hs)
			if err != nil {
				return fmt.Errorf("while fetching properties for host %q: %s", hs.Reference().Value, err)
			}
			if hsProps.Runtime.InMaintenanceMode {
				err := hostsystem.ExitMaintenanceMode(hs, provider.DefaultAPITimeout)
				if err != nil {
					return fmt.Errorf("while getting host %q out of maintenance mode: %s", hs.Reference().Value, err)
				}
			}
		}
	}

	// Remove hosts next
	if err := clustercomputeresource.MoveHostsOutOf(cluster, oldHosts, d.Get("host_cluster_exit_timeout").(int)); err != nil {
		return fmt.Errorf("error moving old hosts out of cluster: %s", err)
	}

	return nil
}

func resourceVSphereComputeClusterGetHostSystemObjects(client *govmomi.Client, hsIDs []string) ([]*object.HostSystem, error) {
	var hosts []*object.HostSystem

	for _, hsID := range hsIDs {
		hs, err := hostsystem.FromID(client, hsID)
		if err != nil {
			return nil, fmt.Errorf("error locating host system ID %q: %s", hsID, err)
		}
		hosts = append(hosts, hs)
	}

	return hosts, nil
}

func resourceVSphereComputeClusterApplyClusterConfiguration(
	d *schema.ResourceData,
	meta interface{},
	cluster *object.ClusterComputeResource,
) error {
	// This is a no-op if there is no config changed
	if !resourceVSphereComputeClusterHasClusterConfigChange(d) {
		log.Printf("[DEBUG] %s: No cluster-specific configuration attributes have changed", resourceVSphereComputeClusterIDString(d))
		return nil
	}

	log.Printf("[DEBUG] %s: Applying cluster configuration", resourceVSphereComputeClusterIDString(d))

	client, err := resourceVSphereComputeClusterClient(meta)
	if err != nil {
		return err
	}

	// Get the version of the vSphere connection to help determine what
	// attributes we need to set
	version := viapi.ParseVersionFromClient(client)

	// Expand the cluster configuration.
	spec := expandClusterConfigSpecEx(d, version)

	// Note that the reconfigure for a cluster is the same as a standalone host,
	// hence we send this to the computeresource helper's Reconfigure function.
	return clustercomputeresource.Reconfigure(cluster, spec)
}

// resourceVSphereComputeClusterApplyTags processes the tags step for both
// create and update for vsphere_compute_cluster.
func resourceVSphereComputeClusterApplyTags(d *schema.ResourceData, meta interface{}, cluster *object.ClusterComputeResource) error {
	tagsClient, err := tagsManagerIfDefined(d, meta)
	if err != nil {
		return err
	}

	// Apply any pending tags now
	if tagsClient == nil {
		log.Printf("[DEBUG] %s: Tags unsupported on this connection, skipping", resourceVSphereComputeClusterIDString(d))
		return nil
	}

	log.Printf("[DEBUG] %s: Applying any pending tags", resourceVSphereComputeClusterIDString(d))
	return processTagDiff(tagsClient, d, cluster)
}

// resourceVSphereComputeClusterReadTags reads the tags for
// vsphere_compute_cluster.
func resourceVSphereComputeClusterReadTags(d *schema.ResourceData, meta interface{}, cluster *object.ClusterComputeResource) error {
	if tagsClient, _ := meta.(*VSphereClient).TagsManager(); tagsClient != nil {
		log.Printf("[DEBUG] %s: Reading tags", resourceVSphereComputeClusterIDString(d))
		if err := readTagsForResource(tagsClient, cluster, d); err != nil {
			return err
		}
	} else {
		log.Printf("[DEBUG] %s: Tags unsupported on this connection, skipping tag read", resourceVSphereComputeClusterIDString(d))
	}
	return nil
}

// resourceVSphereComputeClusterApplyCustomAttributes processes the custom
// attributes step for both create and update for vsphere_compute_cluster.
func resourceVSphereComputeClusterApplyCustomAttributes(
	d *schema.ResourceData,
	meta interface{},
	cluster *object.ClusterComputeResource,
) error {
	client := meta.(*VSphereClient).vimClient
	// Verify a proper vCenter before proceeding if custom attributes are defined
	attrsProcessor, err := customattribute.GetDiffProcessorIfAttributesDefined(client, d)
	if err != nil {
		return err
	}

	if attrsProcessor == nil {
		log.Printf("[DEBUG] %s: Custom attributes unsupported on this connection, skipping", resourceVSphereComputeClusterIDString(d))
		return nil
	}

	log.Printf("[DEBUG] %s: Applying any pending custom attributes", resourceVSphereComputeClusterIDString(d))
	return attrsProcessor.ProcessDiff(cluster)
}

// resourceVSphereComputeClusterReadCustomAttributes reads the custom
// attributes for vsphere_compute_cluster.
func resourceVSphereComputeClusterReadCustomAttributes(
	d *schema.ResourceData,
	meta interface{},
	cluster *object.ClusterComputeResource,
) error {
	client := meta.(*VSphereClient).vimClient
	// Read custom attributes
	if customattribute.IsSupported(client) {
		log.Printf("[DEBUG] %s: Reading custom attributes", resourceVSphereComputeClusterIDString(d))
		props, err := clustercomputeresource.Properties(cluster)
		if err != nil {
			return err
		}
		customattribute.ReadFromResource(client, props.Entity(), d)
	} else {
		log.Printf("[DEBUG] %s: Custom attributes unsupported on this connection, skipping", resourceVSphereComputeClusterIDString(d))
	}

	return nil
}

// resourceVSphereComputeClusterGetCluster gets the ComputeClusterResource from the ID
// in the supplied ResourceData.
func resourceVSphereComputeClusterGetCluster(
	d structure.ResourceIDStringer,
	meta interface{},
) (*object.ClusterComputeResource, error) {
	log.Printf("[DEBUG] %s: Fetching ComputeClusterResource object from resource ID", resourceVSphereComputeClusterIDString(d))
	client, err := resourceVSphereComputeClusterClient(meta)
	if err != nil {
		return nil, err
	}

	return clustercomputeresource.FromID(client, d.Id())
}

// resourceVSphereComputeClusterGetClusterFromPath gets the ComputeClusterResource from a
// supplied path. If no datacenter is supplied, the path must be a full path.
func resourceVSphereComputeClusterGetClusterFromPath(
	meta interface{},
	path string,
	dcID string,
) (*object.ClusterComputeResource, error) {
	client, err := resourceVSphereComputeClusterClient(meta)
	if err != nil {
		return nil, err
	}
	var dc *object.Datacenter
	if dcID != "" {
		var err error
		dc, err = datacenterFromID(client, dcID)
		if err != nil {
			return nil, fmt.Errorf("cannot locate datacenter: %s", err)
		}
		log.Printf("[DEBUG] Looking for cluster %q in datacenter %q", path, dc.InventoryPath)
	} else {
		log.Printf("[DEBUG] Fetching cluster at path %q", path)
	}

	return clustercomputeresource.FromPath(client, path, dc)
}

// resourceVSphereComputeClusterSaveDatacenter saves the datacenter that the
// cluster is a member of for correctness on imports.
func resourceVSphereComputeClusterSaveDatacenter(
	d *schema.ResourceData,
	meta interface{},
	cluster *object.ClusterComputeResource,
) error {
	log.Printf("[DEBUG] %s: Saving datacenter", resourceVSphereComputeClusterIDString(d))
	client, err := resourceVSphereComputeClusterClient(meta)
	if err != nil {
		return err
	}

	p, err := folder.RootPathParticleHost.SplitDatacenter(cluster.InventoryPath)
	if err != nil {
		return fmt.Errorf("error parsing datacenter path from cluster: %s", err)
	}

	dc, err := getDatacenter(client, p)
	if err != nil {
		return fmt.Errorf("error fetching datacenter for cluster: %s", err)
	}

	return d.Set("datacenter_id", dc.Reference().Value)
}

// resourceVSphereComputeClusterSaveNameAndPath saves the name and path of a
// StoragePod into the supplied ResourceData.
func resourceVSphereComputeClusterSaveNameAndPath(d *schema.ResourceData, cluster *object.ClusterComputeResource) error {
	log.Printf(
		"[DEBUG] %s: Saving name and path data for cluster %q",
		resourceVSphereComputeClusterIDString(d),
		cluster.InventoryPath,
	)

	if err := d.Set("name", cluster.Name()); err != nil {
		return fmt.Errorf("error saving name: %s", err)
	}

	f, err := folder.RootPathParticleHost.SplitRelativeFolder(cluster.InventoryPath)
	if err != nil {
		return fmt.Errorf("error parsing cluster path %q: %s", cluster.InventoryPath, err)
	}
	if err := d.Set("folder", folder.NormalizePath(f)); err != nil {
		return fmt.Errorf("error saving folder: %s", err)
	}
	return nil
}

// resourceVSphereComputeClusterApplyNameChange applies any changes to a
// ClusterComputeResource's name.
func resourceVSphereComputeClusterApplyNameChange(
	d *schema.ResourceData,
	meta interface{},
	cluster *object.ClusterComputeResource,
) (*object.ClusterComputeResource, error) {
	log.Printf(
		"[DEBUG] %s: Applying any name changes (old path = %q)",
		resourceVSphereComputeClusterIDString(d),
		cluster.InventoryPath,
	)

	var changed bool
	var err error

	if d.HasChange("name") {
		if err = clustercomputeresource.Rename(cluster, d.Get("name").(string)); err != nil {
			return nil, fmt.Errorf("error renaming cluster: %s", err)
		}
		changed = true
	}

	if changed {
		// Update the cluster so that we have the new inventory path for logging and
		// other things
		cluster, err = resourceVSphereComputeClusterGetCluster(d, meta)
		if err != nil {
			return nil, fmt.Errorf("error refreshing cluster after name change: %s", err)
		}
		log.Printf(
			"[DEBUG] %s: Name changed, new path = %q",
			resourceVSphereComputeClusterIDString(d),
			cluster.InventoryPath,
		)
	}

	return cluster, nil
}

// resourceVSphereComputeClusterApplyFolderChange applies any changes to a
// ClusterComputeResource's folder location.
func resourceVSphereComputeClusterApplyFolderChange(
	d *schema.ResourceData,
	meta interface{},
	cluster *object.ClusterComputeResource,
) (*object.ClusterComputeResource, error) {
	log.Printf(
		"[DEBUG] %s: Applying any folder changes (old path = %q)",
		resourceVSphereComputeClusterIDString(d),
		cluster.InventoryPath,
	)

	var changed bool
	var err error

	if d.HasChange("folder") {
		f := d.Get("folder").(string)
		client := meta.(*VSphereClient).vimClient
		if err = clustercomputeresource.MoveToFolder(client, cluster, f); err != nil {
			return nil, fmt.Errorf("could not move cluster to folder %q: %s", f, err)
		}
		changed = true
	}

	if changed {
		// Update the cluster so that we have the new inventory path for logging and
		// other things
		cluster, err = resourceVSphereComputeClusterGetCluster(d, meta)
		if err != nil {
			return nil, fmt.Errorf("error refreshing cluster after folder change: %s", err)
		}
		log.Printf(
			"[DEBUG] %s: Folder changed, new path = %q",
			resourceVSphereComputeClusterIDString(d),
			cluster.InventoryPath,
		)
	}

	return cluster, nil
}

// resourceVSphereComputeClusterValidateEmptyCluster validates that the cluster
// is empty. This is used to ensure a safe deletion of the cluster - we do not
// allow deletion of clusters that still virtual machines or hosts in them.
func resourceVSphereComputeClusterValidateEmptyCluster(
	d structure.ResourceIDStringer,
	cluster *object.ClusterComputeResource,
) error {
	log.Printf("[DEBUG] %s: Checking to ensure that cluster is empty", resourceVSphereComputeClusterIDString(d))
	ne, err := clustercomputeresource.HasChildren(cluster)
	if err != nil {
		return fmt.Errorf("error checking for cluster contents: %s", err)
	}
	if ne {
		return fmt.Errorf(
			"cluster %q still has hosts or virtual machines. Please move or remove all items before deleting",
			cluster.InventoryPath,
		)
	}
	return nil
}

// resourceVSphereComputeClusterDeleteProcessForceRemoveHosts process
// force-evacuation if the resource has been configured to do so.
//
// NOTE: As documented, this should only be used in testing. Improper use
// of this option can lead to service disruptions and/or may fail to
// actually succeed depending on the resources actually in use in the
// cluster, and specific constraints that exist in the cluster.
func resourceVSphereComputeClusterDeleteProcessForceRemoveHosts(
	d *schema.ResourceData,
	meta interface{},
	cluster *object.ClusterComputeResource,
) error {
	if !d.Get("force_evacuate_on_destroy").(bool) {
		return nil
	}

	client, err := resourceVSphereComputeClusterClient(meta)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Force-evacuating hosts in cluster before removal", resourceVSphereComputeClusterIDString(d))
	hosts, err := resourceVSphereComputeClusterGetHostSystemObjects(
		client,
		structure.SliceInterfacesToStrings(d.Get("host_system_ids").(*schema.Set).List()),
	)
	if err != nil {
		return err
	}

	if err := clustercomputeresource.MoveHostsOutOf(cluster, hosts, d.Get("host_cluster_exit_timeout").(int)); err != nil {
		return fmt.Errorf("error force-removing old hosts out of cluster: %s", err)
	}

	return nil
}

// resourceVSphereComputeClusterApplyDelete process the removal of a
// cluster.
func resourceVSphereComputeClusterApplyDelete(d structure.ResourceIDStringer, cluster *object.ClusterComputeResource) error {
	log.Printf("[DEBUG] %s: Proceeding with cluster deletion", resourceVSphereComputeClusterIDString(d))
	if err := clustercomputeresource.Delete(cluster); err != nil {
		return err
	}
	return nil
}

// resourceVSphereComputeClusterFlattenData saves the configuration attributes
// from a ClusterComputeResource into the supplied ResourceData. It also saves
// the root resource pool for the cluster in resource_pool_id.
//
// Note that other functions handle other non-configuration related items, such
// as path, name, tags, and custom attributes.
func resourceVSphereComputeClusterFlattenData(
	d *schema.ResourceData,
	meta interface{},
	cluster *object.ClusterComputeResource,
) error {
	log.Printf("[DEBUG] %s: Saving cluster attributes", resourceVSphereComputeClusterIDString(d))
	client, err := resourceVSphereComputeClusterClient(meta)
	if err != nil {
		return err
	}

	// Get the version of the vSphere connection to help determine what
	// attributes we need to set
	version := viapi.ParseVersionFromClient(client)

	props, err := clustercomputeresource.Properties(cluster)
	if err != nil {
		return err
	}

	// Save the root resource pool ID so that it can be passed on to other
	// resources without having to resort to the data source.
	if err := d.Set("resource_pool_id", props.ResourcePool.Value); err != nil {
		return err
	}

	if !d.Get("host_managed").(bool) {
		hostList := []string{}
		for _, host := range props.Host {
			hostList = append(hostList, host.Value)
		}
		d.Set("host_system_ids", hostList)
	}

	err = flattenVsanDisks(d, cluster, meta.(*VSphereClient).vimClient)
	if err != nil {
		return err
	}

	return flattenClusterConfigSpecEx(d, props.ConfigurationEx.(*types.ClusterConfigInfoEx), version)
}

// expandClusterConfigSpecEx reads certain ResourceData keys and returns a
// ClusterConfigSpecEx.
func expandClusterConfigSpecEx(d *schema.ResourceData, version viapi.VSphereVersion) *types.ClusterConfigSpecEx {
	obj := &types.ClusterConfigSpecEx{
		DasConfig: expandClusterDasConfigInfo(d, version),
		DpmConfig: expandClusterDpmConfigInfo(d),
		DrsConfig: expandClusterDrsConfigInfo(d),
	}

	if version.Newer(viapi.VSphereVersion{Product: version.Product, Major: 6, Minor: 5}) {
		obj.InfraUpdateHaConfig = expandClusterInfraUpdateHaConfigInfo(d)
		obj.Orchestration = expandClusterOrchestrationInfo(d)
		obj.ProactiveDrsConfig = expandClusterProactiveDrsConfigInfo(d)
	}

	obj.VsanConfig = expandVsanConfig(d)

	return obj
}

func expandVsanConfig(d *schema.ResourceData) *types.VsanClusterConfigInfo {
	conf := &types.VsanClusterConfigInfo{}
	enabled := d.Get("vsan_enabled").(bool)

	conf.Enabled = &enabled
	conf.DefaultConfig = &types.VsanClusterConfigInfoHostDefaultInfo{}
	return conf
}

func updateVsanDisks(d *schema.ResourceData, cluster *object.ClusterComputeResource, meta interface{}) error {
	client := meta.(*VSphereClient).vimClient
	od, nd := d.GetChange("vsan_disk_group")
	delSet := structure.DiffSlice(od.([]interface{}), nd.([]interface{}))
	addSet := structure.DiffSlice(nd.([]interface{}), od.([]interface{}))
	delSetI := delSet
	addSetI := addSet

	for i, del := range delSetI {
		r := del.(map[string]interface{})
		for n, add := range addSetI {
			a := add.(map[string]interface{})
			if r["cache"].(string) != a["cache"].(string) {
				continue
			}

			ds := r["storage"].(*schema.Set)
			as := a["storage"].(*schema.Set)
			switch {
			case ds.Len() > as.Len():
				addSet = structure.DropSliceItem(addSet, n)
				delSet = structure.DropSliceItem(delSet, i)
				r["storage"] = ds.Difference(as)
				if r["storage"].(*schema.Set).Len() >= 1 {
					r["cache"] = ""
				}
				delSet = append(delSet, r)
			case ds.Len() < as.Len():
				addSet = structure.DropSliceItem(addSet, n)
				delSet = structure.DropSliceItem(delSet, i)
				a["storage"] = as.Difference(ds)
				addSet = append(addSet, a)
			default:
				addSet = structure.DropSliceItem(addSet, n)
				delSet = structure.DropSliceItem(delSet, i)
			}
		}
	}

	hosts, err := clustercomputeresource.Hosts(cluster)
	if err != nil {
		return err
	}

	for _, host := range hosts {
		if err = deleteVsanDisks(host, delSet, client); err != nil {
			return err
		}
		if err = addVsanDisks(host, addSet, client); err != nil {
			return err
		}
	}
	return nil
}

func hostStorageSystemPropertiesFromHostSystemID(client *govmomi.Client, hostID string) (*mo.HostStorageSystem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	hss, err := hostStorageSystemFromHostSystemID(client, hostID)
	if err != nil {
		return nil, err
	}
	var hssProps mo.HostStorageSystem
	err = hss.Properties(ctx, hss.Reference(), nil, &hssProps)
	return &hssProps, err
}

func generateDiskMap(client *govmomi.Client, host *object.HostSystem, list []interface{}) (*types.VsanHostDiskMapping, error) {
	diskMap := types.VsanHostDiskMapping{
		NonSsd: []types.HostScsiDisk{},
	}
	hssProps, err := hostStorageSystemPropertiesFromHostSystemID(client, host.Reference().Value)
	if err != nil {
		return nil, err
	}
	if hssProps.StorageDeviceInfo == nil {
		return &diskMap, nil
	}
	for _, scsiLun := range hssProps.StorageDeviceInfo.ScsiLun {
		for _, diskGroup := range list {
			if hostDisk, ok := scsiLun.(*types.HostScsiDisk); ok {
				if err != nil {
					return nil, err
				}
				for _, storageDisk := range diskGroup.(map[string]interface{})["storage"].(*schema.Set).List() {
					if hostDisk.CanonicalName == storageDisk.(string) {
						diskMap.NonSsd = append(diskMap.NonSsd, *hostDisk)
					}
				}
				if diskGroup.(map[string]interface{})["cache"].(string) == hostDisk.CanonicalName {
					diskMap.Ssd = *hostDisk
				}
			}
		}
	}
	return &diskMap, nil
}

func deleteVsanDisks(host *object.HostSystem, list []interface{}, client *govmomi.Client) error {
	log.Printf("deleteVsanDisks: Starting removal of vSAN disks on %s.", host.Name())
	hvs, err := vsansystem.FromHost(client, host, defaultAPITimeout)
	if err != nil {
		return nil
	}
	diskMap, err := generateDiskMap(client, host, list)
	if err != nil {
		return err
	}
	if diskMap.Ssd.CanonicalName != "" || len(diskMap.NonSsd) > 0 {
		log.Printf("deleteVsanDisks: Scheduled disks are being removed.")
		if err = vsansystem.RemoveDiskMapping(client, host, hvs, diskMap, defaultAPITimeout); err != nil {
			return err
		}
		log.Printf("deleteVsanDisks: vSAN disks successfully removed.")
	} else {
		log.Printf("deleteVsanDisks: No vSAN disks to remove on %s.", host.Name())
	}
	return nil
}

func addVsanDisks(host *object.HostSystem, list []interface{}, client *govmomi.Client) error {
	log.Printf("addVsanDisks: Starting initialization of vSAN disks on %s.", host.Name())
	hvs, err := vsansystem.FromHost(client, host, defaultAPITimeout)
	if err != nil {
		return nil
	}
	diskMap, err := generateDiskMap(client, host, list)
	if diskMap.Ssd.CanonicalName != "" {
		log.Printf("addVsanDisks: Scheduled disks are being initialized.")
		if err = vsansystem.InitializeDisks(client, host, hvs, diskMap, defaultAPITimeout); err != nil {
			return err
		}
		log.Printf("addVsanDisks: vSAN disks successfully initialized.")
	} else {
		log.Printf("addVsanDisks: No vSAN disks to initialize on %s.", host.Name())
	}
	return nil
}

func flattenVsanDisks(d *schema.ResourceData, cluster *object.ClusterComputeResource, client *govmomi.Client) error {
	diskMap := []interface{}{}

	hosts, err := clustercomputeresource.Hosts(cluster)
	if err != nil {
		return err
	}
	for _, host := range hosts {
		hvs, err := vsansystem.FromHost(client, host, defaultAPITimeout)
		if err != nil {
			return err
		}
		hvsProps, err := vsansystem.Properties(client, hvs, defaultAPITimeout)
		if err != nil {
			return err
		}
		if hvsProps.Config.StorageInfo == nil {
			return nil
		}
		for _, diskGroup := range hvsProps.Config.StorageInfo.DiskMapping {
			var vsanStorage []string
			for _, disk := range diskGroup.NonSsd {
				vsanStorage = append(vsanStorage, disk.CanonicalName)
			}
			diskMap = append(diskMap, map[string]interface{}{
				"cache":   diskGroup.Ssd.CanonicalName,
				"storage": vsanStorage,
			})
		}
	}
	return d.Set("vsan_disk_group", diskMap)
}

// flattenClusterConfigSpecEx saves a ClusterConfigSpecEx into the supplied
// ResourceData.
func flattenClusterConfigSpecEx(d *schema.ResourceData, obj *types.ClusterConfigInfoEx, version viapi.VSphereVersion) error {
	if err := flattenClusterDasConfigInfo(d, obj.DasConfig, version); err != nil {
		return err
	}
	if err := flattenClusterDpmConfigInfo(d, obj.DpmConfigInfo); err != nil {
		return err
	}
	if err := flattenClusterDrsConfigInfo(d, obj.DrsConfig); err != nil {
		return err
	}
	d.Set("vsan_enabled", obj.VsanConfigInfo.Enabled)

	if version.Newer(viapi.VSphereVersion{Product: version.Product, Major: 6, Minor: 5}) {
		if err := flattenClusterInfraUpdateHaConfigInfo(d, obj.InfraUpdateHaConfig); err != nil {
			return err
		}
	}
	if version.Newer(viapi.VSphereVersion{Product: version.Product, Major: 6, Minor: 5}) {
		if err := flattenClusterOrchestrationInfo(d, obj.Orchestration); err != nil {
			return err
		}
	}
	if version.Newer(viapi.VSphereVersion{Product: version.Product, Major: 6, Minor: 5}) {
		return flattenClusterProactiveDrsConfigInfo(d, obj.ProactiveDrsConfig)
	}

	return nil
}

// expandClusterDasConfigInfo reads certain ResourceData keys and returns a
// ClusterDasConfigInfo.
func expandClusterDasConfigInfo(d *schema.ResourceData, version viapi.VSphereVersion) *types.ClusterDasConfigInfo {
	obj := &types.ClusterDasConfigInfo{
		DefaultVmSettings:          expandClusterDasVMSettings(d, version),
		Enabled:                    structure.GetBool(d, "ha_enabled"),
		HBDatastoreCandidatePolicy: d.Get("ha_heartbeat_datastore_policy").(string),
		HostMonitoring:             d.Get("ha_host_monitoring").(string),
		Option:                     expandResourceVSphereComputeClusterDasAdvancedOptions(d),
		VmMonitoring:               d.Get("ha_vm_monitoring").(string),
		HeartbeatDatastore: structure.SliceInterfacesToManagedObjectReferences(
			d.Get("ha_heartbeat_datastore_ids").(*schema.Set).List(),
			"Datastore",
		),
	}

	policy := d.Get("ha_admission_control_policy").(string)
	if policy != clusterAdmissionControlTypeDisabled {
		obj.AdmissionControlEnabled = structure.BoolPtr(true)
	} else {
		obj.AdmissionControlEnabled = structure.BoolPtr(false)
	}
	obj.AdmissionControlPolicy = expandBaseClusterDasAdmissionControlPolicy(d, policy, version)

	if version.Newer(viapi.VSphereVersion{Product: version.Product, Major: 6}) {
		obj.VmComponentProtecting = d.Get("ha_vm_component_protection").(string)
	}

	return obj
}

// flattenClusterDasConfigInfo saves a ClusterDasConfigInfo into the supplied
// ResourceData.
func flattenClusterDasConfigInfo(d *schema.ResourceData, obj types.ClusterDasConfigInfo, version viapi.VSphereVersion) error {
	var dsIDs []string
	for _, v := range obj.HeartbeatDatastore {
		dsIDs = append(dsIDs, v.Value)
	}

	err := structure.SetBatch(d, map[string]interface{}{
		"ha_enabled":                    obj.Enabled,
		"ha_heartbeat_datastore_policy": obj.HBDatastoreCandidatePolicy,
		"ha_host_monitoring":            obj.HostMonitoring,
		"ha_vm_monitoring":              obj.VmMonitoring,
		"ha_heartbeat_datastore_ids":    dsIDs,
	})
	if err != nil {
		return err
	}

	if err := flattenClusterDasVMSettings(d, obj.DefaultVmSettings, version); err != nil {
		return err
	}
	if err := flattenResourceVSphereComputeClusterDasAdvancedOptions(d, obj.Option); err != nil {
		return err
	}

	if version.Newer(viapi.VSphereVersion{Product: version.Product, Major: 6}) {
		if err := d.Set("ha_vm_component_protection", obj.VmComponentProtecting); err != nil {
			return err
		}
	}

	// If AdmissionControlEnabled is false, AdmissionControlPolicy is still
	// configured. Set ha_admission_control_policy to disabled before
	// flattenBaseClusterDasAdmissionControlPolicy, so AdmissionControlEnabled
	// can still be checked.
	if *obj.AdmissionControlEnabled == false {
		return d.Set("ha_admission_control_policy", clusterAdmissionControlTypeDisabled)
	}
	return flattenBaseClusterDasAdmissionControlPolicy(d, obj.AdmissionControlPolicy, version)
}

// expandBaseClusterDasAdmissionControlPolicy reads certain ResourceData keys
// and returns a BaseClusterDasAdmissionControlPolicy.
func expandBaseClusterDasAdmissionControlPolicy(
	d *schema.ResourceData,
	policy string,
	version viapi.VSphereVersion,
) types.BaseClusterDasAdmissionControlPolicy {
	var obj types.BaseClusterDasAdmissionControlPolicy

	switch policy {
	case clusterAdmissionControlTypeResourcePercentage:
		obj = expandClusterFailoverResourcesAdmissionControlPolicy(d, version)
	case clusterAdmissionControlTypeSlotPolicy:
		obj = expandClusterFailoverLevelAdmissionControlPolicy(d)
	case clusterAdmissionControlTypeFailoverHosts:
		obj = expandClusterFailoverHostAdmissionControlPolicy(d, version)
	case clusterAdmissionControlTypeDisabled:
		return nil
	}

	if version.Newer(viapi.VSphereVersion{Product: version.Product, Major: 6, Minor: 5}) {
		obj.GetClusterDasAdmissionControlPolicy().ResourceReductionToToleratePercent = structure.Int32Ptr(int32(d.Get("ha_admission_control_performance_tolerance").(int)))
	}

	return obj
}

// flattenBaseClusterDasAdmissionControlPolicy saves a
// BaseClusterDasAdmissionControlPolicy into the supplied ResourceData.
func flattenBaseClusterDasAdmissionControlPolicy(
	d *schema.ResourceData,
	obj types.BaseClusterDasAdmissionControlPolicy,
	version viapi.VSphereVersion,
) error {
	var policy string

	switch t := obj.(type) {
	case *types.ClusterFailoverResourcesAdmissionControlPolicy:
		if err := flattenClusterFailoverResourcesAdmissionControlPolicy(d, t, version); err != nil {
			return err
		}
		policy = clusterAdmissionControlTypeResourcePercentage
	case *types.ClusterFailoverLevelAdmissionControlPolicy:
		if err := flattenClusterFailoverLevelAdmissionControlPolicy(d, t); err != nil {
			return err
		}
		policy = clusterAdmissionControlTypeSlotPolicy
	case *types.ClusterFailoverHostAdmissionControlPolicy:
		if err := flattenClusterFailoverHostAdmissionControlPolicy(d, t, version); err != nil {
			return err
		}
		policy = clusterAdmissionControlTypeFailoverHosts
	default:
		policy = clusterAdmissionControlTypeDisabled
	}

	return d.Set("ha_admission_control_policy", policy)
}

// expandClusterFailoverResourcesAdmissionControlPolicy reads certain
// ResourceData keys and returns a
// ClusterFailoverResourcesAdmissionControlPolicy.
func expandClusterFailoverResourcesAdmissionControlPolicy(
	d *schema.ResourceData,
	version viapi.VSphereVersion,
) *types.ClusterFailoverResourcesAdmissionControlPolicy {
	obj := &types.ClusterFailoverResourcesAdmissionControlPolicy{
		CpuFailoverResourcesPercent:    int32(d.Get("ha_admission_control_resource_percentage_cpu").(int)),
		MemoryFailoverResourcesPercent: int32(d.Get("ha_admission_control_resource_percentage_memory").(int)),
	}

	if version.Newer(viapi.VSphereVersion{Product: version.Product, Major: 6, Minor: 5}) {
		obj.AutoComputePercentages = structure.GetBool(d, "ha_admission_control_resource_percentage_auto_compute")
		obj.FailoverLevel = int32(d.Get("ha_admission_control_host_failure_tolerance").(int))
	}

	return obj
}

// flattenClusterFailoverResourcesAdmissionControlPolicy saves a
// ClusterFailoverResourcesAdmissionControlPolicy into the supplied
// ResourceData.
func flattenClusterFailoverResourcesAdmissionControlPolicy(
	d *schema.ResourceData,
	obj *types.ClusterFailoverResourcesAdmissionControlPolicy,
	version viapi.VSphereVersion,
) error {
	// AutoComputePercentages is a vSphere >= 6.5 feature, but when it's
	// enabled the admission control CPU/memory values will be auto-set and
	// caused spurious diffs, so do a nil check to see if we have the value or
	// if it's disabled before we set the values.
	if obj.AutoComputePercentages == nil || !*obj.AutoComputePercentages {
		err := structure.SetBatch(d, map[string]interface{}{
			"ha_admission_control_resource_percentage_cpu":    obj.CpuFailoverResourcesPercent,
			"ha_admission_control_resource_percentage_memory": obj.MemoryFailoverResourcesPercent,
		})
		if err != nil {
			return err
		}
	}

	if version.Newer(viapi.VSphereVersion{Product: version.Product, Major: 6, Minor: 5}) {
		return structure.SetBatch(d, map[string]interface{}{
			"ha_admission_control_resource_percentage_auto_compute": obj.AutoComputePercentages,
			"ha_admission_control_host_failure_tolerance":           obj.FailoverLevel,
		})
	}

	return nil
}

// expandClusterFailoverLevelAdmissionControlPolicy reads certain ResourceData
// keys and returns a ClusterFailoverLevelAdmissionControlPolicy.
func expandClusterFailoverLevelAdmissionControlPolicy(d *schema.ResourceData) *types.ClusterFailoverLevelAdmissionControlPolicy {
	obj := &types.ClusterFailoverLevelAdmissionControlPolicy{
		FailoverLevel: int32(d.Get("ha_admission_control_host_failure_tolerance").(int)),
	}

	if d.Get("ha_admission_control_slot_policy_use_explicit_size").(bool) {
		obj.SlotPolicy = &types.ClusterFixedSizeSlotPolicy{
			Cpu:    int32(d.Get("ha_admission_control_resource_percentage_cpu").(int)),
			Memory: int32(d.Get("ha_admission_control_resource_percentage_memory").(int)),
		}
	}

	return obj
}

// flattenClusterFailoverLevelAdmissionControlPolicy saves a
// ClusterFailoverLevelAdmissionControlPolicy into the supplied ResourceData.
func flattenClusterFailoverLevelAdmissionControlPolicy(
	d *schema.ResourceData,
	obj *types.ClusterFailoverLevelAdmissionControlPolicy,
) error {
	if err := d.Set("ha_admission_control_host_failure_tolerance", obj.FailoverLevel); err != nil {
		return err
	}

	if obj.SlotPolicy != nil {
		return structure.SetBatch(d, map[string]interface{}{
			"ha_admission_control_resource_percentage_cpu":    obj.SlotPolicy.(*types.ClusterFixedSizeSlotPolicy).Cpu,
			"ha_admission_control_resource_percentage_memory": obj.SlotPolicy.(*types.ClusterFixedSizeSlotPolicy).Memory,
		})
	}

	return nil
}

// expandClusterFailoverHostAdmissionControlPolicy reads certain ResourceData
// keys and returns a ClusterFailoverHostAdmissionControlPolicy.
func expandClusterFailoverHostAdmissionControlPolicy(
	d *schema.ResourceData,
	version viapi.VSphereVersion,
) *types.ClusterFailoverHostAdmissionControlPolicy {
	obj := &types.ClusterFailoverHostAdmissionControlPolicy{
		FailoverHosts: structure.SliceInterfacesToManagedObjectReferences(
			d.Get("ha_admission_control_failover_host_system_ids").(*schema.Set).List(),
			"HostSystem",
		),
	}

	if version.Newer(viapi.VSphereVersion{Product: version.Product, Major: 6, Minor: 5}) {
		obj.FailoverLevel = int32(d.Get("ha_admission_control_host_failure_tolerance").(int))
	}

	return obj
}

// flattenClusterFailoverHostAdmissionControlPolicy saves a
// ClusterFailoverHostAdmissionControlPolicy into the supplied ResourceData.
func flattenClusterFailoverHostAdmissionControlPolicy(
	d *schema.ResourceData,
	obj *types.ClusterFailoverHostAdmissionControlPolicy,
	version viapi.VSphereVersion,
) error {
	var hsIDs []string
	for _, v := range obj.FailoverHosts {
		hsIDs = append(hsIDs, v.Value)
	}

	if err := d.Set("ha_admission_control_failover_host_system_ids", hsIDs); err != nil {
		return err
	}

	if version.Newer(viapi.VSphereVersion{Product: version.Product, Major: 6, Minor: 5}) {
		return d.Set("ha_admission_control_host_failure_tolerance", obj.FailoverLevel)
	}

	return nil
}

// expandClusterDasVMSettings reads certain ResourceData keys and returns a
// ClusterDasVmSettings.
func expandClusterDasVMSettings(d *schema.ResourceData, version viapi.VSphereVersion) *types.ClusterDasVmSettings {
	obj := &types.ClusterDasVmSettings{
		IsolationResponse:         d.Get("ha_host_isolation_response").(string),
		RestartPriority:           d.Get("ha_vm_restart_priority").(string),
		VmToolsMonitoringSettings: expandClusterVMToolsMonitoringSettings(d),
	}

	if version.Newer(viapi.VSphereVersion{Product: version.Product, Major: 6}) {
		obj.VmComponentProtectionSettings = expandClusterVMComponentProtectionSettings(d)
	}
	if version.Newer(viapi.VSphereVersion{Product: version.Product, Major: 6, Minor: 5}) {
		obj.RestartPriorityTimeout = int32(d.Get("ha_vm_restart_timeout").(int))
	}

	return obj
}

// flattenClusterDasVMSettings saves a ClusterDasVmSettings into the supplied
// ResourceData.
func flattenClusterDasVMSettings(d *schema.ResourceData, obj *types.ClusterDasVmSettings, version viapi.VSphereVersion) error {
	err := structure.SetBatch(d, map[string]interface{}{
		"ha_host_isolation_response": obj.IsolationResponse,
		"ha_vm_restart_priority":     obj.RestartPriority,
	})
	if err != nil {
		return err
	}

	if err := flattenClusterVMToolsMonitoringSettings(d, obj.VmToolsMonitoringSettings); err != nil {
		return err
	}

	if version.Newer(viapi.VSphereVersion{Product: version.Product, Major: 6}) {
		if err := flattenClusterVMComponentProtectionSettings(d, obj.VmComponentProtectionSettings); err != nil {
			return err
		}
	}
	if version.Newer(viapi.VSphereVersion{Product: version.Product, Major: 6, Minor: 5}) {
		return d.Set("ha_vm_restart_timeout", obj.RestartPriorityTimeout)
	}

	return nil
}

// expandClusterVMComponentProtectionSettings reads certain ResourceData keys and returns a
// ClusterVmComponentProtectionSettings.
func expandClusterVMComponentProtectionSettings(d *schema.ResourceData) *types.ClusterVmComponentProtectionSettings {
	obj := &types.ClusterVmComponentProtectionSettings{
		VmReactionOnAPDCleared:    d.Get("ha_datastore_apd_recovery_action").(string),
		VmStorageProtectionForAPD: d.Get("ha_datastore_apd_response").(string),
		VmStorageProtectionForPDL: d.Get("ha_datastore_pdl_response").(string),
		VmTerminateDelayForAPDSec: int32(d.Get("ha_datastore_apd_response_delay").(int)),
	}

	if d.Get("ha_datastore_apd_response").(string) != string(types.ClusterVmComponentProtectionSettingsStorageVmReactionDisabled) {
		// Flag EnableAPDTimeoutForHosts to ensure that APD is enabled for all
		// hosts in the cluster and our other settings here will be effective. Note
		// that this setting is not persisted to state or the vSphere backend and
		// is actually a host operation, not a cluster operation. It's here to
		// ensure that the settings specified here are otherwise effective. We may
		// need to revisit this if we introduce more robust host management
		// capabilities in the provider.
		obj.EnableAPDTimeoutForHosts = structure.BoolPtr(true)
	}

	return obj
}

// flattenClusterVmComponentProtectionSettings saves a
// ClusterVmComponentProtectionSettings into the supplied ResourceData.
func flattenClusterVMComponentProtectionSettings(d *schema.ResourceData, obj *types.ClusterVmComponentProtectionSettings) error {
	return structure.SetBatch(d, map[string]interface{}{
		"ha_datastore_apd_recovery_action": obj.VmReactionOnAPDCleared,
		"ha_datastore_apd_response":        obj.VmStorageProtectionForAPD,
		"ha_datastore_pdl_response":        obj.VmStorageProtectionForPDL,
		"ha_datastore_apd_response_delay":  obj.VmTerminateDelayForAPDSec,
	})
}

// expandClusterVMToolsMonitoringSettings reads certain ResourceData keys and returns a
// ClusterVmToolsMonitoringSettings.
func expandClusterVMToolsMonitoringSettings(d *schema.ResourceData) *types.ClusterVmToolsMonitoringSettings {
	obj := &types.ClusterVmToolsMonitoringSettings{
		FailureInterval:  int32(d.Get("ha_vm_failure_interval").(int)),
		MaxFailures:      int32(d.Get("ha_vm_maximum_resets").(int)),
		MaxFailureWindow: int32(d.Get("ha_vm_maximum_failure_window").(int)),
		MinUpTime:        int32(d.Get("ha_vm_minimum_uptime").(int)),
		VmMonitoring:     d.Get("ha_vm_monitoring").(string),
	}

	return obj
}

// flattenClusterVmToolsMonitoringSettings saves a
// ClusterVmToolsMonitoringSettings into the supplied ResourceData.
func flattenClusterVMToolsMonitoringSettings(d *schema.ResourceData, obj *types.ClusterVmToolsMonitoringSettings) error {
	return structure.SetBatch(d, map[string]interface{}{
		"ha_vm_failure_interval":       obj.FailureInterval,
		"ha_vm_maximum_resets":         obj.MaxFailures,
		"ha_vm_maximum_failure_window": obj.MaxFailureWindow,
		"ha_vm_minimum_uptime":         obj.MinUpTime,
		"ha_vm_monitoring":             obj.VmMonitoring,
	})
}

// expandResourceVSphereComputeClusterDasAdvancedOptions reads certain
// ResourceData keys and returns a BaseOptionValue list designed for use as DAS
// (vSphere HA) advanced options.
func expandResourceVSphereComputeClusterDasAdvancedOptions(d *schema.ResourceData) []types.BaseOptionValue {
	var opts []types.BaseOptionValue

	m := d.Get("ha_advanced_options").(map[string]interface{})
	for k, v := range m {
		opts = append(opts, &types.OptionValue{
			Key:   k,
			Value: types.AnyType(v),
		})
	}
	return opts
}

// flattenResourceVSphereComputeClusterDasAdvancedOptions saves a
// BaseOptionValue into the supplied ResourceData for DAS (vSphere HA) advanced
// options.
func flattenResourceVSphereComputeClusterDasAdvancedOptions(d *schema.ResourceData, opts []types.BaseOptionValue) error {
	m := make(map[string]interface{})
	for _, opt := range opts {
		m[opt.GetOptionValue().Key] = opt.GetOptionValue().Value
	}

	return d.Set("ha_advanced_options", m)
}

// expandClusterDpmConfigInfo reads certain ResourceData keys and returns a
// ClusterDpmConfigInfo.
func expandClusterDpmConfigInfo(d *schema.ResourceData) *types.ClusterDpmConfigInfo {
	obj := &types.ClusterDpmConfigInfo{
		DefaultDpmBehavior:  types.DpmBehavior(d.Get("dpm_automation_level").(string)),
		Enabled:             structure.GetBool(d, "dpm_enabled"),
		HostPowerActionRate: int32(d.Get("dpm_threshold").(int)),
	}

	return obj
}

// flattenClusterDpmConfigInfo saves a ClusterDpmConfigInfo into the supplied
// ResourceData.
func flattenClusterDpmConfigInfo(d *schema.ResourceData, obj *types.ClusterDpmConfigInfo) error {
	return structure.SetBatch(d, map[string]interface{}{
		"dpm_automation_level": obj.DefaultDpmBehavior,
		"dpm_enabled":          obj.Enabled,
		"dpm_threshold":        obj.HostPowerActionRate,
	})
}

// expandClusterDrsConfigInfo reads certain ResourceData keys and returns a
// ClusterDrsConfigInfo.
func expandClusterDrsConfigInfo(d *schema.ResourceData) *types.ClusterDrsConfigInfo {
	obj := &types.ClusterDrsConfigInfo{
		DefaultVmBehavior:         types.DrsBehavior(d.Get("drs_automation_level").(string)),
		Enabled:                   structure.GetBool(d, "drs_enabled"),
		EnableVmBehaviorOverrides: structure.GetBool(d, "drs_enable_vm_overrides"),
		VmotionRate:               int32(d.Get("drs_migration_threshold").(int)),
		Option:                    expandResourceVSphereComputeClusterDrsAdvancedOptions(d),
	}

	return obj
}

// flattenClusterDrsConfigInfo saves a ClusterDrsConfigInfo into the supplied
// ResourceData.
func flattenClusterDrsConfigInfo(d *schema.ResourceData, obj types.ClusterDrsConfigInfo) error {
	err := structure.SetBatch(d, map[string]interface{}{
		"drs_automation_level":    obj.DefaultVmBehavior,
		"drs_enabled":             obj.Enabled,
		"drs_enable_vm_overrides": obj.EnableVmBehaviorOverrides,
		"drs_migration_threshold": obj.VmotionRate,
	})
	if err != nil {
		return err
	}

	return flattenResourceVSphereComputeClusterDrsAdvancedOptions(d, obj.Option)
}

// expandResourceVSphereComputeClusterDrsAdvancedOptions reads certain
// ResourceData keys and returns a BaseOptionValue list designed for use as DRS
// advanced options.
func expandResourceVSphereComputeClusterDrsAdvancedOptions(d *schema.ResourceData) []types.BaseOptionValue {
	var opts []types.BaseOptionValue

	m := d.Get("drs_advanced_options").(map[string]interface{})
	for k, v := range m {
		opts = append(opts, &types.OptionValue{
			Key:   k,
			Value: types.AnyType(v),
		})
	}
	return opts
}

// flattenResourceVSphereComputeClusterDrsAdvancedOptions saves a
// BaseOptionValue into the supplied ResourceData for DRS and DPM advanced
// options.
func flattenResourceVSphereComputeClusterDrsAdvancedOptions(d *schema.ResourceData, opts []types.BaseOptionValue) error {
	m := make(map[string]interface{})
	for _, opt := range opts {
		m[opt.GetOptionValue().Key] = opt.GetOptionValue().Value
	}

	return d.Set("drs_advanced_options", m)
}

// expandClusterInfraUpdateHaConfigInfo reads certain ResourceData keys and returns a
// ClusterInfraUpdateHaConfigInfo.
func expandClusterInfraUpdateHaConfigInfo(d *schema.ResourceData) *types.ClusterInfraUpdateHaConfigInfo {
	obj := &types.ClusterInfraUpdateHaConfigInfo{
		Behavior:            d.Get("proactive_ha_automation_level").(string),
		Enabled:             structure.GetBool(d, "proactive_ha_enabled"),
		ModerateRemediation: d.Get("proactive_ha_moderate_remediation").(string),
		Providers:           structure.SliceInterfacesToStrings(d.Get("proactive_ha_provider_ids").(*schema.Set).List()),
		SevereRemediation:   d.Get("proactive_ha_severe_remediation").(string),
	}

	return obj
}

// flattenClusterInfraUpdateHaConfigInfo saves a ClusterInfraUpdateHaConfigInfo into the
// supplied ResourceData.
func flattenClusterInfraUpdateHaConfigInfo(d *schema.ResourceData, obj *types.ClusterInfraUpdateHaConfigInfo) error {
	return structure.SetBatch(d, map[string]interface{}{
		"proactive_ha_automation_level":     obj.Behavior,
		"proactive_ha_enabled":              obj.Enabled,
		"proactive_ha_moderate_remediation": obj.ModerateRemediation,
		"proactive_ha_provider_ids":         obj.Providers,
		"proactive_ha_severe_remediation":   obj.SevereRemediation,
	})
}

// expandClusterOrchestrationInfo reads certain ResourceData keys and returns a
// ClusterOrchestrationInfo.
func expandClusterOrchestrationInfo(d *schema.ResourceData) *types.ClusterOrchestrationInfo {
	obj := &types.ClusterOrchestrationInfo{
		DefaultVmReadiness: &types.ClusterVmReadiness{
			PostReadyDelay: int32(d.Get("ha_vm_restart_additional_delay").(int)),
			ReadyCondition: d.Get("ha_vm_dependency_restart_condition").(string),
		},
	}

	return obj
}

// flattenClusterOrchestrationInfo saves a ClusterOrchestrationInfo into the
// supplied ResourceData.
func flattenClusterOrchestrationInfo(d *schema.ResourceData, obj *types.ClusterOrchestrationInfo) error {
	return structure.SetBatch(d, map[string]interface{}{
		"ha_vm_restart_additional_delay":     obj.DefaultVmReadiness.PostReadyDelay,
		"ha_vm_dependency_restart_condition": obj.DefaultVmReadiness.ReadyCondition,
	})
}

// expandClusterProactiveDrsConfigInfo reads certain ResourceData keys and returns a
// ClusterProactiveDrsConfigInfo.
func expandClusterProactiveDrsConfigInfo(d *schema.ResourceData) *types.ClusterProactiveDrsConfigInfo {
	obj := &types.ClusterProactiveDrsConfigInfo{
		Enabled: structure.GetBool(d, "drs_enable_predictive_drs"),
	}

	return obj
}

// flattenClusterProactiveDrsConfigInfo saves a ClusterProactiveDrsConfigInfo into the
// supplied ResourceData.
func flattenClusterProactiveDrsConfigInfo(d *schema.ResourceData, obj *types.ClusterProactiveDrsConfigInfo) error {
	return structure.SetBatch(d, map[string]interface{}{
		"drs_enable_predictive_drs": obj.Enabled,
	})
}

// resourceVSphereComputeClusterIDString prints a friendly string for the
// vsphere_compute_cluster resource.
func resourceVSphereComputeClusterIDString(d structure.ResourceIDStringer) string {
	return structure.ResourceIDString(d, resourceVSphereComputeClusterName)
}

func resourceVSphereComputeClusterClient(meta interface{}) (*govmomi.Client, error) {
	client := meta.(*VSphereClient).vimClient
	if err := viapi.ValidateVirtualCenter(client); err != nil {
		return nil, err
	}
	return client, nil
}

// resourceVSphereComputeClusterHasClusterConfigChange checks all resource keys
// associated with cluster configuration (and not, for example, member hosts,
// folder, tags, etc) to see if there has been a change in the configuration of
// those keys. This helper is designed to detect no-ops in a cluster
// configuration to see if we really need to send a configure API call to
// vSphere.
func resourceVSphereComputeClusterHasClusterConfigChange(d *schema.ResourceData) bool {
	for k := range resourceVSphereComputeCluster().Schema {
		switch {
		case resourceVSphereComputeClusterHasClusterConfigChangeExcluded(k):
			continue
		case d.HasChange(k):
			return true
		}
	}

	return false
}

func resourceVSphereComputeClusterHasClusterConfigChangeExcluded(k string) bool {
	// It's easier to track which keys don't belong to storage DRS versus the
	// ones that do.
	excludeKeys := []string{
		"name",
		"datacenter_id",
		"host_system_ids",
		"folder",
		"host_cluster_exit_timeout",
		"force_evacuate_on_destroy",
		vSphereTagAttributeKey,
		customattribute.ConfigKey,
	}

	for _, exclude := range excludeKeys {
		if k == exclude {
			return true
		}
	}

	return false
}
