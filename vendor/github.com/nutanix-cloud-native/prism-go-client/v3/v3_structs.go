package v3

import (
	"time"

	"github.com/nutanix-cloud-native/prism-go-client"
)

// Reference ...
type Reference struct {
	Kind *string `json:"kind" mapstructure:"kind"`
	Name *string `json:"name,omitempty" mapstructure:"name,omitempty"`
	UUID *string `json:"uuid" mapstructure:"uuid"`
}

// VMVnumaConfig Indicates how VM vNUMA should be configured
type VMVnumaConfig struct {
	// Number of vNUMA nodes. 0 means vNUMA is disabled.
	NumVnumaNodes *int64 `json:"num_vnuma_nodes,omitempty" mapstructure:"num_vnuma_nodes,omitempty"`
}

type VMSerialPort struct {
	Index       *int64 `json:"index,omitempty" mapstructure:"index,omitempty"`
	IsConnected *bool  `json:"is_connected,omitempty" mapstructure:"is_connected,omitempty"`
}

// IPAddress An IP address.
type IPAddress struct {
	// Address *string.
	IP *string `json:"ip,omitempty" mapstructure:"ip,omitempty"`

	// Address type. It can only be \"ASSIGNED\" in the spec. If no type is specified in the spec, the default type is
	// set to \"ASSIGNED\".
	Type *string `json:"type,omitempty" mapstructure:"type,omitempty"`
}

// VMNic Virtual Machine NIC.
type VMNic struct {
	// IP endpoints for the adapter. Currently, IPv4 addresses are supported.
	IPEndpointList []*IPAddress `json:"ip_endpoint_list,omitempty" mapstructure:"ip_endpoint_list,omitempty"`

	// The MAC address for the adapter.
	MacAddress *string `json:"mac_address,omitempty" mapstructure:"mac_address,omitempty"`

	// The model of this NIC.
	Model *string `json:"model,omitempty" mapstructure:"model,omitempty"`

	NetworkFunctionChainReference *Reference `json:"network_function_chain_reference,omitempty" mapstructure:"network_function_chain_reference,omitempty"`

	// The number of queues for this NIC
	NumQueues *int64 `json:"num_queues,omitempty" mapstructure:"num_queues,omitempty"`

	// The type of this Network function NIC. Defaults to INGRESS.
	NetworkFunctionNicType *string `json:"network_function_nic_type,omitempty" mapstructure:"network_function_nic_type,omitempty"`

	// The type of this NIC. Defaults to NORMAL_NIC.
	NicType *string `json:"nic_type,omitempty" mapstructure:"nic_type,omitempty"`

	SubnetReference *Reference `json:"subnet_reference,omitempty" mapstructure:"subnet_reference,omitempty"`

	// The NIC's UUID, which is used to uniquely identify this particular NIC. This UUID may be used to refer to the NIC
	// outside the context of the particular VM it is attached to.
	UUID *string `json:"uuid,omitempty" mapstructure:"uuid,omitempty"`

	IsConnected *bool `json:"is_connected,omitempty" mapstructure:"is_connected,omitempty"`
}

// DiskAddress Disk Address.
type DiskAddress struct {
	AdapterType *string `json:"adapter_type,omitempty" mapstructure:"adapter_type,omitempty"`
	DeviceIndex *int64  `json:"device_index,omitempty" mapstructure:"device_index,omitempty"`
}

// VMBootDevice Indicates which device a VM should boot from. One of disk_address or mac_address should be provided.
type VMBootDevice struct {
	// Address of disk to boot from.
	DiskAddress *DiskAddress `json:"disk_address,omitempty" mapstructure:"disk_address,omitempty"`

	// MAC address of nic to boot from.
	MacAddress *string `json:"mac_address,omitempty" mapstructure:"mac_address,omitempty"`
}

// VMBootConfig Indicates which device a VM should boot from.
type VMBootConfig struct {
	// Indicates which device a VM should boot from. Boot device takes precdence over boot device order. If both are
	// given then specified boot device will be primary boot device and remaining devices will be assigned boot order
	// according to boot device order field.
	BootDevice *VMBootDevice `json:"boot_device,omitempty" mapstructure:"boot_device,omitempty"`
	BootType   *string       `json:"boot_type,omitempty" mapstructure:"boot_type,omitempty"`

	// Indicates the order of device types in which VM should try to boot from. If boot device order is not provided the
	// system will decide appropriate boot device order.
	BootDeviceOrderList []*string `json:"boot_device_order_list,omitempty" mapstructure:"boot_device_order_list,omitempty"`
}

// NutanixGuestToolsSpec Information regarding Nutanix Guest Tools.
type NutanixGuestToolsSpec struct {
	State                 *string           `json:"state,omitempty" mapstructure:"state,omitempty"`                                     // Nutanix Guest Tools is enabled or not.
	Version               *string           `json:"version,omitempty" mapstructure:"version,omitempty"`                                 // Version of Nutanix Guest Tools installed on the VM.
	NgtState              *string           `json:"ngt_state,omitempty" mapstructure:"ngt_state,omitempty"`                             // Nutanix Guest Tools installed or not.
	Credentials           map[string]string `json:"credentials,omitempty" mapstructure:"credentials,omitempty"`                         // Credentials to login server
	IsoMountState         *string           `json:"iso_mount_state,omitempty" mapstructure:"iso_mount_state,omitempty"`                 // Desired mount state of Nutanix Guest Tools ISO.
	EnabledCapabilityList []*string         `json:"enabled_capability_list,omitempty" mapstructure:"enabled_capability_list,omitempty"` // Application names that are enabled.
}

// GuestToolsSpec Information regarding guest tools.
type GuestToolsSpec struct {
	// Nutanix Guest Tools information
	NutanixGuestTools *NutanixGuestToolsSpec `json:"nutanix_guest_tools,omitempty" mapstructure:"nutanix_guest_tools,omitempty"`
}

// VMGpu Graphics resource information for the Virtual Machine.
type VMGpu struct {
	// The device ID of the GPU.
	DeviceID *int64 `json:"device_id,omitempty" mapstructure:"device_id,omitempty"`

	// The mode of this GPU.
	Mode *string `json:"mode,omitempty" mapstructure:"mode,omitempty"`

	// The vendor of the GPU.
	Vendor *string `json:"vendor,omitempty" mapstructure:"vendor,omitempty"`
}

// GuestCustomizationCloudInit If this field is set, the guest will be customized using cloud-init. Either user_data or
// custom_key_values should be provided. If custom_key_ves are provided then the user data will be generated using these
// key-value pairs.
type GuestCustomizationCloudInit struct {
	// Generic key value pair used for custom attributes
	CustomKeyValues map[string]string `json:"custom_key_values,omitempty" mapstructure:"custom_key_values,omitempty"`

	// The contents of the meta_data configuration for cloud-init. This can be formatted as YAML or JSON. The value must
	// be base64 encoded.
	MetaData *string `json:"meta_data,omitempty" mapstructure:"meta_data,omitempty"`

	// The contents of the user_data configuration for cloud-init. This can be formatted as YAML, JSON, or could be a
	// shell script. The value must be base64 encoded.
	UserData *string `json:"user_data,omitempty" mapstructure:"user_data,omitempty"`
}

// GuestCustomizationSysprep If this field is set, the guest will be customized using Sysprep. Either unattend_xml or
// custom_key_values should be provided. If custom_key_values are provided then the unattended answer file will be
// generated using these key-value pairs.
type GuestCustomizationSysprep struct {
	// Generic key value pair used for custom attributes
	CustomKeyValues map[string]string `json:"custom_key_values,omitempty" mapstructure:"custom_key_values,omitempty"`

	// Whether the guest will be freshly installed using this unattend configuration, or whether this unattend
	// configuration will be applied to a pre-prepared image. Default is \"PREPARED\".
	InstallType *string `json:"install_type,omitempty" mapstructure:"install_type,omitempty"`

	// This field contains a Sysprep unattend xml definition, as a *string. The value must be base64 encoded.
	UnattendXML *string `json:"unattend_xml,omitempty" mapstructure:"unattend_xml,omitempty"`
}

// GuestCustomization VM guests may be customized at boot time using one of several different methods. Currently,
// cloud-init w/ ConfigDriveV2 (for Linux VMs) and Sysprep (for Windows VMs) are supported. Only ONE OF sysprep or
// cloud_init should be provided. Note that guest customization can currently only be set during VM creation. Attempting
// to change it after creation will result in an error. Additional properties can be specified. For example - in the
// context of VM template creation if \"override_script\" is set to \"True\" then the deployer can upload their own
// custom script.
type GuestCustomization struct {
	CloudInit *GuestCustomizationCloudInit `json:"cloud_init,omitempty" mapstructure:"cloud_init,omitempty"`

	// Flag to allow override of customization by deployer.
	IsOverridable *bool `json:"is_overridable,omitempty" mapstructure:"is_overridable,omitempty"`

	Sysprep *GuestCustomizationSysprep `json:"sysprep,omitempty" mapstructure:"sysprep,omitempty"`
}

// VMGuestPowerStateTransitionConfig Extra configs related to power state transition.
type VMGuestPowerStateTransitionConfig struct {
	// Indicates whether to execute set script before ngt shutdown/reboot.
	EnableScriptExec *bool `json:"enable_script_exec,omitempty" mapstructure:"enable_script_exec,omitempty"`

	// Indicates whether to abort ngt shutdown/reboot if script fails.
	ShouldFailOnScriptFailure *bool `json:"should_fail_on_script_failure,omitempty" mapstructure:"should_fail_on_script_failure,omitempty"`
}

// VMPowerStateMechanism Indicates the mechanism guiding the VM power state transition. Currently used for the transition
// to \"OFF\" state.
type VMPowerStateMechanism struct {
	GuestTransitionConfig *VMGuestPowerStateTransitionConfig `json:"guest_transition_config,omitempty" mapstructure:"guest_transition_config,omitempty"`

	// Power state mechanism (ACPI/GUEST/HARD).
	Mechanism *string `json:"mechanism,omitempty" mapstructure:"mechanism,omitempty"`
}

// VMDiskDeviceProperties ...
type VMDiskDeviceProperties struct {
	DeviceType  *string      `json:"device_type,omitempty" mapstructure:"device_type,omitempty"`
	DiskAddress *DiskAddress `json:"disk_address,omitempty" mapstructure:"disk_address,omitempty"`
}

// StorageContainerReference references to a kind. Either one of (kind, uuid) or url needs to be specified.
type StorageContainerReference struct {
	URL  string `json:"url,omitempty"`
	Kind string `json:"kind,omitempty"`
	UUID string `json:"uuid,omitempty"`
	Name string `json:"name,omitempty"`
}

// VMStorageConfig specifies the storage configuration parameters for VM disks.
type VMStorageConfig struct {
	FlashMode                 string                     `json:"flash_mode,omitempty"`
	StorageContainerReference *StorageContainerReference `json:"storage_container_reference,omitempty"`
}

// VMDisk VirtualMachine Disk (VM Disk).
type VMDisk struct {
	DataSourceReference *Reference `json:"data_source_reference,omitempty" mapstructure:"data_source_reference,omitempty"`

	DeviceProperties *VMDiskDeviceProperties `json:"device_properties,omitempty" mapstructure:"device_properties,omitempty"`

	// Size of the disk in Bytes.
	DiskSizeBytes *int64 `json:"disk_size_bytes,omitempty" mapstructure:"disk_size_bytes,omitempty"`

	// Size of the disk in MiB. Must match the size specified in 'disk_size_bytes' - rounded up to the nearest MiB -
	// when that field is present.
	DiskSizeMib *int64 `json:"disk_size_mib,omitempty" mapstructure:"disk_size_mib,omitempty"`

	// The device ID which is used to uniquely identify this particular disk.
	UUID *string `json:"uuid,omitempty" mapstructure:"uuid,omitempty"`

	VolumeGroupReference *Reference `json:"volume_group_reference,omitempty" mapstructure:"volume_group_reference,omitempty"`

	// This preference specifies the storage configuration parameters for VM disks.
	StorageConfig *VMStorageConfig `json:"storage_config,omitempty" mapstructure:"storage_config,omitempty"`
}

// VMResources VM Resources Definition.
type VMResources struct {
	// Indicates which device the VM should boot from.
	BootConfig *VMBootConfig `json:"boot_config,omitempty" mapstructure:"boot_config,omitempty"`

	// Disks attached to the VM.
	DiskList []*VMDisk `json:"disk_list,omitempty" mapstructure:"disk_list,omitempty"`

	// GPUs attached to the VM.
	GpuList []*VMGpu `json:"gpu_list,omitempty" mapstructure:"gpu_list,omitempty"`

	GuestCustomization *GuestCustomization `json:"guest_customization,omitempty" mapstructure:"guest_customization,omitempty"`

	// Guest OS Identifier. For ESX, refer to VMware documentation link
	// https://www.vmware.com/support/orchestrator/doc/vro-vsphere65-api/html/VcVirtualMachineGuestOsIdentifier.html
	// for the list of guest OS identifiers.
	GuestOsID *string `json:"guest_os_id,omitempty" mapstructure:"guest_os_id,omitempty"`

	// Information regarding guest tools.
	GuestTools *GuestToolsSpec `json:"guest_tools,omitempty" mapstructure:"guest_tools,omitempty"`

	// VM's hardware clock timezone in IANA TZDB format (America/Los_Angeles).
	HardwareClockTimezone *string `json:"hardware_clock_timezone,omitempty" mapstructure:"hardware_clock_timezone,omitempty"`

	// Memory size in MiB.
	MemorySizeMib *int64 `json:"memory_size_mib,omitempty" mapstructure:"memory_size_mib,omitempty"`

	// NICs attached to the VM.
	NicList []*VMNic `json:"nic_list,omitempty" mapstructure:"nic_list,omitempty"`

	// Number of threads per core
	NumThreads *int64 `json:"num_threads_per_core,omitempty" mapstructure:"num_threads_per_core,omitempty"`

	// Number of vCPU sockets.
	NumSockets *int64 `json:"num_sockets,omitempty" mapstructure:"num_sockets,omitempty"`

	// Number of vCPUs per socket.
	NumVcpusPerSocket *int64 `json:"num_vcpus_per_socket,omitempty" mapstructure:"num_vcpus_per_socket,omitempty"`

	// *Reference to an entity that the VM should be cloned from.
	ParentReference *Reference `json:"parent_reference,omitempty" mapstructure:"parent_reference,omitempty"`

	// The current or desired power state of the VM.
	PowerState *string `json:"power_state,omitempty" mapstructure:"power_state,omitempty"`

	PowerStateMechanism *VMPowerStateMechanism `json:"power_state_mechanism,omitempty" mapstructure:"power_state_mechanism,omitempty"`

	// Indicates whether VGA console should be enabled or not.
	VgaConsoleEnabled *bool `json:"vga_console_enabled,omitempty" mapstructure:"vga_console_enabled,omitempty"`

	// Indicates whether to passthrough the host’s CPU features to the guest. Enabling this will disable live migration of the VM.
	EnableCPUPassthrough *bool `json:"enable_cpu_passthrough,omitempty" mapstructure:"enable_cpu_passthrough,omitempty"`

	// Indicates whether a VMs cores are getting pinned on either CPU1, CPU2, CPU3 or CPU4. By default, the Linux Scheduler in AHV will pin all cores wherever they are best available.
	EnableCPUPinning *bool `json:"is_vcpu_hard_pinned,omitempty" mapstructure:"is_vcpu_hard_pinned,omitempty"`

	// Information regarding vNUMA configuration.
	VMVnumaConfig *VMVnumaConfig `json:"vnuma_config,omitempty" mapstructure:"vnuma_config,omitempty"`

	SerialPortList []*VMSerialPort `json:"serial_port_list,omitempty" mapstructure:"serial_port_list,omitempty"`

	MachineType *string `json:"machine_type,omitempty" mapstructure:"machine_type,omitempty"`
}

// VM An intentful representation of a vm spec
type VM struct {
	AvailabilityZoneReference *Reference `json:"availability_zone_reference,omitempty" mapstructure:"availability_zone_reference,omitempty"`

	ClusterReference *Reference `json:"cluster_reference,omitempty" mapstructure:"cluster_reference,omitempty"`

	// A description for vm.
	Description *string `json:"description,omitempty" mapstructure:"description,omitempty"`

	// vm Name.
	Name *string `json:"name" mapstructure:"name"`

	Resources *VMResources `json:"resources,omitempty" mapstructure:"resources,omitempty"`
}

// VMIntentInput ...
type VMIntentInput struct {
	APIVersion *string `json:"api_version,omitempty" mapstructure:"api_version,omitempty"`

	Metadata *Metadata `json:"metadata" mapstructure:"metadata"`

	Spec *VM `json:"spec" mapstructure:"spec"`
}

// MessageResource ...
type MessageResource struct {
	// Custom key-value details relevant to the status.
	Details interface{} `json:"details,omitempty" mapstructure:"details,omitempty"`

	// If state is ERROR, a message describing the error.
	Message *string `json:"message" mapstructure:"message"`

	// If state is ERROR, a machine-readable snake-cased *string.
	Reason *string `json:"reason" mapstructure:"reason"`
}

// VMStatus The status of a REST API call. Only used when there is a failure to report.
type VMStatus struct {
	APIVersion *string `json:"api_version,omitempty" mapstructure:"api_version,omitempty"`

	// The HTTP error code.
	Code *int64 `json:"code,omitempty" mapstructure:"code,omitempty"`

	// The kind name
	Kind *string `json:"kind,omitempty" mapstructure:"kind,omitempty"`

	MessageList []*MessageResource `json:"message_list,omitempty" mapstructure:"message_list,omitempty"`

	State *string `json:"state,omitempty" mapstructure:"state,omitempty"`
}

// VMNicOutputStatus Virtual Machine NIC Status.
type VMNicOutputStatus struct {
	// The Floating IP associated with the vnic.
	FloatingIP *string `json:"floating_ip,omitempty" mapstructure:"floating_ip,omitempty"`

	// IP endpoints for the adapter. Currently, IPv4 addresses are supported.
	IPEndpointList []*IPAddress `json:"ip_endpoint_list,omitempty" mapstructure:"ip_endpoint_list,omitempty"`

	// The MAC address for the adapter.
	MacAddress *string `json:"mac_address,omitempty" mapstructure:"mac_address,omitempty"`

	// The model of this NIC.
	Model *string `json:"model,omitempty" mapstructure:"model,omitempty"`

	NetworkFunctionChainReference *Reference `json:"network_function_chain_reference,omitempty" mapstructure:"network_function_chain_reference,omitempty"`

	// The number of queues for this NIC
	NumQueues *int64 `json:"num_queues,omitempty" mapstructure:"num_queues,omitempty"`

	// The type of this Network function NIC. Defaults to INGRESS.
	NetworkFunctionNicType *string `json:"network_function_nic_type,omitempty" mapstructure:"network_function_nic_type,omitempty"`

	// The type of this NIC. Defaults to NORMAL_NIC.
	NicType *string `json:"nic_type,omitempty" mapstructure:"nic_type,omitempty"`

	SubnetReference *Reference `json:"subnet_reference,omitempty" mapstructure:"subnet_reference,omitempty"`

	// The NIC's UUID, which is used to uniquely identify this particular NIC. This UUID may be used to refer to the NIC
	// outside the context of the particular VM it is attached to.
	UUID *string `json:"uuid,omitempty" mapstructure:"uuid,omitempty"`

	IsConnected *bool `json:"is_connected,omitempty" mapstructure:"is_connected,omitempty"`
}

// NutanixGuestToolsStatus Information regarding Nutanix Guest Tools.
type NutanixGuestToolsStatus struct {
	// Version of Nutanix Guest Tools available on the cluster.
	AvailableVersion *string `json:"available_version,omitempty" mapstructure:"available_version,omitempty"`
	// Nutanix Guest Tools installed or not.
	NgtState *string `json:"ngt_state,omitempty" mapstructure:"ngt_state,omitempty"`
	// Desired mount state of Nutanix Guest Tools ISO.
	IsoMountState *string `json:"iso_mount_state,omitempty" mapstructure:"iso_mount_state,omitempty"`
	// Nutanix Guest Tools is enabled or not.
	State *string `json:"state,omitempty" mapstructure:"state,omitempty"`
	// Version of Nutanix Guest Tools installed on the VM.
	Version *string `json:"version,omitempty" mapstructure:"version,omitempty"`
	// Application names that are enabled.
	EnabledCapabilityList []*string `json:"enabled_capability_list,omitempty" mapstructure:"enabled_capability_list,omitempty"`
	// Credentials to login server
	Credentials map[string]string `json:"credentials,omitempty" mapstructure:"credentials,omitempty"`
	// Version of the operating system on the VM.
	GuestOsVersion *string `json:"guest_os_version,omitempty" mapstructure:"guest_os_version,omitempty"`
	// Whether the VM is configured to take VSS snapshots through NGT.
	VSSSnapshotCapable *bool `json:"vss_snapshot_capable,omitempty" mapstructure:"vss_snapshot_capable,omitempty"`
	// Communication from VM to CVM is active or not.
	IsReachable *bool `json:"is_reachable,omitempty" mapstructure:"is_reachable,omitempty"`
	// Whether VM mobility drivers are installed in the VM.
	VMMobilityDriversInstalled *bool `json:"vm_mobility_drivers_installed,omitempty" mapstructure:"vm_mobility_drivers_installed,omitempty"`
}

// GuestToolsStatus Information regarding guest tools.
type GuestToolsStatus struct {
	// Nutanix Guest Tools information
	NutanixGuestTools *NutanixGuestToolsStatus `json:"nutanix_guest_tools,omitempty" mapstructure:"nutanix_guest_tools,omitempty"`
}

// VMGpuOutputStatus Graphics resource status information for the Virtual Machine.
type VMGpuOutputStatus struct {
	// The device ID of the GPU.
	DeviceID *int64 `json:"device_id,omitempty" mapstructure:"device_id,omitempty"`

	// Fraction of the physical GPU assigned.
	Fraction *int64 `json:"fraction,omitempty" mapstructure:"fraction,omitempty"`

	// GPU frame buffer size in MiB.
	FrameBufferSizeMib *int64 `json:"frame_buffer_size_mib,omitempty" mapstructure:"frame_buffer_size_mib,omitempty"`

	// Last determined guest driver version.
	GuestDriverVersion *string `json:"guest_driver_version,omitempty" mapstructure:"guest_driver_version,omitempty"`

	// The mode of this GPU
	Mode *string `json:"mode,omitempty" mapstructure:"mode,omitempty"`

	// Name of the GPU resource.
	Name *string `json:"name,omitempty" mapstructure:"name,omitempty"`

	// Number of supported virtual display heads.
	NumVirtualDisplayHeads *int64 `json:"num_virtual_display_heads,omitempty" mapstructure:"num_virtual_display_heads,omitempty"`

	// GPU {segment:bus:device:function} (sbdf) address if assigned.
	PCIAddress *string `json:"pci_address,omitempty" mapstructure:"pci_address,omitempty"`

	// UUID of the GPU.
	UUID *string `json:"uuid,omitempty" mapstructure:"uuid,omitempty"`

	// The vendor of the GPU.
	Vendor *string `json:"vendor,omitempty" mapstructure:"vendor,omitempty"`
}

// GuestCustomizationStatus VM guests may be customized at boot time using one of several different methods. Currently,
// cloud-init w/ ConfigDriveV2 (for Linux VMs) and Sysprep (for Windows VMs) are supported. Only ONE OF sysprep or
// cloud_init should be provided. Note that guest customization can currently only be set during VM creation. Attempting
// to change it after creation will result in an error. Additional properties can be specified. For example - in the
// context of VM template creation if \"override_script\" is set to \"True\" then the deployer can upload their own
// custom script.
type GuestCustomizationStatus struct {
	CloudInit *GuestCustomizationCloudInit `json:"cloud_init,omitempty" mapstructure:"cloud_init,omitempty"`

	// Flag to allow override of customization by deployer.
	IsOverridable *bool `json:"is_overridable,omitempty" mapstructure:"is_overridable,omitempty"`

	Sysprep *GuestCustomizationSysprep `json:"sysprep,omitempty" mapstructure:"sysprep,omitempty"`
}

// VMResourcesDefStatus VM Resources Status Definition.
type VMResourcesDefStatus struct {
	// Indicates which device the VM should boot from.
	BootConfig *VMBootConfig `json:"boot_config,omitempty" mapstructure:"boot_config,omitempty"`

	// Disks attached to the VM.
	DiskList []*VMDisk `json:"disk_list,omitempty" mapstructure:"disk_list,omitempty"`

	// GPUs attached to the VM.
	GpuList []*VMGpuOutputStatus `json:"gpu_list,omitempty" mapstructure:"gpu_list,omitempty"`

	GuestCustomization *GuestCustomizationStatus `json:"guest_customization,omitempty" mapstructure:"guest_customization,omitempty"`

	// Guest OS Identifier. For ESX, refer to VMware documentation link
	// https://www.vmware.com/support/orchestrator/doc/vro-vsphere65-api/html/VcVirtualMachineGuestOsIdentifier.html
	// for the list of guest OS identifiers.
	GuestOsID *string `json:"guest_os_id,omitempty" mapstructure:"guest_os_id,omitempty"`

	// Information regarding guest tools.
	GuestTools *GuestToolsStatus `json:"guest_tools,omitempty" mapstructure:"guest_tools,omitempty"`

	// VM's hardware clock timezone in IANA TZDB format (America/Los_Angeles).
	HardwareClockTimezone *string `json:"hardware_clock_timezone,omitempty" mapstructure:"hardware_clock_timezone,omitempty"`

	HostReference *Reference `json:"host_reference,omitempty" mapstructure:"host_reference,omitempty"`

	// The hypervisor type for the hypervisor the VM is hosted on.
	HypervisorType *string `json:"hypervisor_type,omitempty" mapstructure:"hypervisor_type,omitempty"`

	// Memory size in MiB.
	MemorySizeMib *int64 `json:"memory_size_mib,omitempty" mapstructure:"memory_size_mib,omitempty"`

	// NICs attached to the VM.
	NicList []*VMNicOutputStatus `json:"nic_list,omitempty" mapstructure:"nic_list,omitempty"`

	// Number of vCPU sockets.
	NumSockets *int64 `json:"num_sockets,omitempty" mapstructure:"num_sockets,omitempty"`

	// Number of vCPUs per socket.
	NumVcpusPerSocket *int64 `json:"num_vcpus_per_socket,omitempty" mapstructure:"num_vcpus_per_socket,omitempty"`

	// *Reference to an entity that the VM cloned from.
	ParentReference *Reference `json:"parent_reference,omitempty" mapstructure:"parent_reference,omitempty"`

	// Current power state of the VM.
	PowerState *string `json:"power_state,omitempty" mapstructure:"power_state,omitempty"`

	PowerStateMechanism *VMPowerStateMechanism `json:"power_state_mechanism,omitempty" mapstructure:"power_state_mechanism,omitempty"`

	// Indicates whether VGA console has been enabled or not.
	VgaConsoleEnabled *bool `json:"vga_console_enabled,omitempty" mapstructure:"vga_console_enabled,omitempty"`

	// Indicates whether to passthrough the host’s CPU features to the guest. Enabling this will disable live migration of the VM.
	EnableCPUPassthrough *bool `json:"enable_cpu_passthrough,omitempty" mapstructure:"enable_cpu_passthrough,omitempty"`

	// Indicates whether a VMs cores are getting pinned on either CPU1, CPU2, CPU3 or CPU4. By default, the Linux Scheduler in AHV will pin all cores wherever they are best available.
	EnableCPUPinning *bool `json:"is_vcpu_hard_pinned,omitempty" mapstructure:"is_vcpu_hard_pinned,omitempty"`

	// Information regarding vNUMA configuration.
	VnumaConfig *VMVnumaConfig `json:"vnuma_config,omitempty" mapstructure:"vnuma_config,omitempty"`

	SerialPortList []*VMSerialPort `json:"serial_port_list,omitempty" mapstructure:"serial_port_list,omitempty"`

	MachineType *string `json:"machine_type,omitempty" mapstructure:"machine_type,omitempty"`
}

// VMDefStatus An intentful representation of a vm status
type VMDefStatus struct {
	AvailabilityZoneReference *Reference `json:"availability_zone_reference,omitempty" mapstructure:"availability_zone_reference,omitempty"`

	ClusterReference *Reference `json:"cluster_reference,omitempty" mapstructure:"cluster_reference,omitempty"`

	// A description for vm.
	Description *string `json:"description,omitempty" mapstructure:"description,omitempty"`

	// Any error messages for the vm, if in an error state.
	MessageList []*MessageResource `json:"message_list,omitempty" mapstructure:"message_list,omitempty"`

	// vm Name.
	Name *string `json:"name,omitempty" mapstructure:"name,omitempty"`

	Resources *VMResourcesDefStatus `json:"resources,omitempty" mapstructure:"resources,omitempty"`

	// The state of the vm.
	State *string `json:"state,omitempty" mapstructure:"state,omitempty"`

	ExecutionContext *ExecutionContext `json:"execution_context,omitempty" mapstructure:"execution_context,omitempty"`
}

// ExecutionContext ...
type ExecutionContext struct {
	TaskUUID interface{} `json:"task_uuid,omitempty" mapstructure:"task_uuid,omitempty"`
}

// VMIntentResponse Response object for intentful operations on a vm
type VMIntentResponse struct {
	APIVersion *string `json:"api_version" mapstructure:"api_version"`

	Metadata *Metadata `json:"metadata,omitempty" mapstructure:"metadata,omitempty"`

	Spec *VM `json:"spec,omitempty" mapstructure:"spec,omitempty"`

	Status *VMDefStatus `json:"status,omitempty" mapstructure:"status,omitempty"`
}

// DSMetadata All api calls that return a list will have this metadata block as input
type DSMetadata struct {
	// The filter in FIQL syntax used for the results.
	Filter *string `json:"filter,omitempty" mapstructure:"filter,omitempty"`

	// The kind name
	Kind *string `json:"kind,omitempty" mapstructure:"kind,omitempty"`

	// The number of records to retrieve relative to the offset
	Length *int64 `json:"length,omitempty" mapstructure:"length,omitempty"`

	// Offset from the start of the entity list
	Offset *int64 `json:"offset,omitempty" mapstructure:"offset,omitempty"`

	// The attribute to perform sort on
	SortAttribute *string `json:"sort_attribute,omitempty" mapstructure:"sort_attribute,omitempty"`

	// The sort order in which results are returned
	SortOrder *string `json:"sort_order,omitempty" mapstructure:"sort_order,omitempty"`

	// Additional filters for client side filtering api response
	ClientSideFilters []*prismgoclient.AdditionalFilter `json:"-"`
}

// VMIntentResource Response object for intentful operations on a vm
type VMIntentResource struct {
	APIVersion *string `json:"api_version,omitempty" mapstructure:"api_version,omitempty"`

	Metadata *Metadata `json:"metadata" mapstructure:"metadata"`

	Spec *VM `json:"spec,omitempty" mapstructure:"spec,omitempty"`

	Status *VMDefStatus `json:"status,omitempty" mapstructure:"status,omitempty"`
}

// VMListIntentResponse Response object for intentful operation of vms
type VMListIntentResponse struct {
	APIVersion *string `json:"api_version" mapstructure:"api_version"`

	Entities []*VMIntentResource `json:"entities,omitempty" mapstructure:"entities,omitempty"`

	Metadata *ListMetadataOutput `json:"metadata" mapstructure:"metadata"`
}

// SubnetMetadata The subnet kind metadata
type SubnetMetadata struct {
	// Categories for the subnet
	Categories map[string]string `json:"categories,omitempty" mapstructure:"categories,omitempty"`

	// UTC date and time in RFC-3339 format when subnet was created
	CreationTime *time.Time `json:"creation_time,omitempty" mapstructure:"creation_time,omitempty"`

	// The kind name
	Kind *string `json:"kind" mapstructure:"kind"`

	// UTC date and time in RFC-3339 format when subnet was last updated
	LastUpdateTime *time.Time `json:"last_update_time,omitempty" mapstructure:"last_update_time,omitempty"`

	// subnet name
	Name *string `json:"name,omitempty" mapstructure:"name,omitempty"`

	OwnerReference *Reference `json:"owner_reference,omitempty" mapstructure:"owner_reference,omitempty"`

	// project reference
	ProjectReference *Reference `json:"project_reference,omitempty" mapstructure:"project_reference,omitempty"`

	// Hash of the spec. This will be returned from server.
	SpecHash *string `json:"spec_hash,omitempty" mapstructure:"spec_hash,omitempty"`

	// Version number of the latest spec.
	SpecVersion *int64 `json:"spec_version,omitempty" mapstructure:"spec_version,omitempty"`

	// subnet uuid
	UUID *string `json:"uuid,omitempty" mapstructure:"uuid,omitempty"`
}

// Address represents the Host address.
type Address struct {
	// Fully qualified domain name.
	FQDN *string `json:"fqdn,omitempty" mapstructure:"fqdn,omitempty"`

	// IPV4 address.
	IP *string `json:"ip,omitempty" mapstructure:"ip,omitempty"`

	// IPV6 address.
	IPV6 *string `json:"ipv6,omitempty" mapstructure:"ipv6,omitempty"`

	// Port Number
	Port *int64 `json:"port,omitempty" mapstructure:"port,omitempty"`
}

// IPPool represents IP pool.
type IPPool struct {
	// Range of IPs (example: 10.0.0.9 10.0.0.19).
	Range *string `json:"range,omitempty" mapstructure:"range,omitempty"`
}

// DHCPOptions Spec for defining DHCP options.
type DHCPOptions struct {
	BootFileName *string `json:"boot_file_name,omitempty" mapstructure:"boot_file_name,omitempty"`

	DomainName *string `json:"domain_name,omitempty" mapstructure:"domain_name,omitempty"`

	DomainNameServerList []*string `json:"domain_name_server_list,omitempty" mapstructure:"domain_name_server_list,omitempty"`

	DomainSearchList []*string `json:"domain_search_list,omitempty" mapstructure:"domain_search_list,omitempty"`

	TFTPServerName *string `json:"tftp_server_name,omitempty" mapstructure:"tftp_server_name,omitempty"`
}

// IPConfig represents the configurtion of IP.
type IPConfig struct {
	// Default gateway IP address.
	DefaultGatewayIP *string `json:"default_gateway_ip,omitempty" mapstructure:"default_gateway_ip,omitempty"`

	DHCPOptions *DHCPOptions `json:"dhcp_options,omitempty" mapstructure:"dhcp_options,omitempty"`

	DHCPServerAddress *Address `json:"dhcp_server_address,omitempty" mapstructure:"dhcp_server_address,omitempty"`

	PoolList []*IPPool `json:"pool_list,omitempty" mapstructure:"pool_list,omitempty"`

	PrefixLength *int64 `json:"prefix_length,omitempty" mapstructure:"prefix_length,omitempty"`

	// Subnet IP address.
	SubnetIP *string `json:"subnet_ip,omitempty" mapstructure:"subnet_ip,omitempty"`
}

// SubnetResources represents Subnet creation/modification spec.
type SubnetResources struct {
	IPConfig *IPConfig `json:"ip_config,omitempty" mapstructure:"ip_config,omitempty"`

	NetworkFunctionChainReference *Reference `json:"network_function_chain_reference,omitempty" mapstructure:"network_function_chain_reference,omitempty"`

	SubnetType *string `json:"subnet_type" mapstructure:"subnet_type"`

	VlanID *int64 `json:"vlan_id,omitempty" mapstructure:"vlan_id,omitempty"`

	VswitchName *string `json:"vswitch_name,omitempty" mapstructure:"vswitch_name,omitempty"`
}

// Subnet An intentful representation of a subnet spec
type Subnet struct {
	AvailabilityZoneReference *Reference `json:"availability_zone_reference,omitempty" mapstructure:"availability_zone_reference,omitempty"`

	ClusterReference *Reference `json:"cluster_reference,omitempty" mapstructure:"cluster_reference,omitempty"`

	// A description for subnet.
	Description *string `json:"description,omitempty" mapstructure:"description,omitempty"`

	// subnet Name.
	Name *string `json:"name" mapstructure:"name"`

	Resources *SubnetResources `json:"resources,omitempty" mapstructure:"resources,omitempty"`
}

// SubnetIntentInput An intentful representation of a subnet
type SubnetIntentInput struct {
	APIVersion *string `json:"api_version,omitempty" mapstructure:"api_version,omitempty"`

	Metadata *Metadata `json:"metadata" mapstructure:"metadata"`

	Spec *Subnet `json:"spec" mapstructure:"spec"`
}

// SubnetStatus represents The status of a REST API call. Only used when there is a failure to report.
type SubnetStatus struct {
	APIVersion *string `json:"api_version,omitempty" mapstructure:"api_version,omitempty"`

	// The HTTP error code.
	Code *int64 `json:"code,omitempty" mapstructure:"code,omitempty"`

	// The kind name
	Kind *string `json:"kind,omitempty" mapstructure:"kind,omitempty"`

	MessageList []*MessageResource `json:"message_list,omitempty" mapstructure:"message_list,omitempty"`

	State *string `json:"state,omitempty" mapstructure:"state,omitempty"`
}

// SubnetResourcesDefStatus represents a Subnet creation/modification status.
type SubnetResourcesDefStatus struct {
	IPConfig *IPConfig `json:"ip_config,omitempty" mapstructure:"ip_config,omitempty"`

	NetworkFunctionChainReference *Reference `json:"network_function_chain_reference,omitempty" mapstructure:"network_function_chain_reference,omitempty"`

	SubnetType *string `json:"subnet_type" mapstructure:"subnet_type"`

	VlanID *int64 `json:"vlan_id,omitempty" mapstructure:"vlan_id,omitempty"`

	VswitchName *string `json:"vswitch_name,omitempty" mapstructure:"vswitch_name,omitempty"`
}

// SubnetDefStatus An intentful representation of a subnet status
type SubnetDefStatus struct {
	AvailabilityZoneReference *Reference `json:"availability_zone_reference,omitempty" mapstructure:"availability_zone_reference,omitempty"`

	ClusterReference *Reference `json:"cluster_reference,omitempty" mapstructure:"cluster_reference,omitempty"`

	// A description for subnet.
	Description *string `json:"description" mapstructure:"description"`

	// Any error messages for the subnet, if in an error state.
	MessageList []*MessageResource `json:"message_list,omitempty" mapstructure:"message_list,omitempty"`

	// subnet Name.
	Name *string `json:"name" mapstructure:"name"`

	Resources *SubnetResourcesDefStatus `json:"resources,omitempty" mapstructure:"resources,omitempty"`

	// The state of the subnet.
	State *string `json:"state,omitempty" mapstructure:"state,omitempty"`

	ExecutionContext *ExecutionContext `json:"execution_context,omitempty" mapstructure:"execution_context,omitempty"`
}

// SubnetIntentResponse represents the response object for intentful operations on a subnet
type SubnetIntentResponse struct {
	APIVersion *string `json:"api_version" mapstructure:"api_version"`

	Metadata *Metadata `json:"metadata,omitempty" mapstructure:"metadata,omitempty"`

	Spec *Subnet `json:"spec,omitempty" mapstructure:"spec,omitempty"`

	Status *SubnetDefStatus `json:"status,omitempty" mapstructure:"status,omitempty"`
}

// SubnetIntentResource represents Response object for intentful operations on a subnet
type SubnetIntentResource struct {
	APIVersion *string `json:"api_version,omitempty" mapstructure:"api_version,omitempty"`

	Metadata *Metadata `json:"metadata" mapstructure:"metadata"`

	Spec *Subnet `json:"spec,omitempty" mapstructure:"spec,omitempty"`

	Status *SubnetDefStatus `json:"status,omitempty" mapstructure:"status,omitempty"`
}

// SubnetListIntentResponse represents the response object for intentful operation of subnets
type SubnetListIntentResponse struct {
	APIVersion *string `json:"api_version" mapstructure:"api_version"`

	Entities []*SubnetIntentResponse `json:"entities,omitempty" mapstructure:"entities,omitempty"`

	Metadata *ListMetadataOutput `json:"metadata" mapstructure:"metadata"`
}

// SubnetListMetadata ...
type SubnetListMetadata struct {
	// The filter in FIQL syntax used for the results.
	Filter *string `json:"filter,omitempty" mapstructure:"filter,omitempty"`

	// The kind name
	Kind *string `json:"kind,omitempty" mapstructure:"kind,omitempty"`

	// The number of records to retrieve relative to the offset
	Length *int64 `json:"length,omitempty" mapstructure:"length,omitempty"`

	// Offset from the start of the entity list
	Offset *int64 `json:"offset,omitempty" mapstructure:"offset,omitempty"`

	// The attribute to perform sort on
	SortAttribute *string `json:"sort_attribute,omitempty" mapstructure:"sort_attribute,omitempty"`

	// The sort order in which results are returned
	SortOrder *string `json:"sort_order,omitempty" mapstructure:"sort_order,omitempty"`
}

// Checksum represents the image checksum
type Checksum struct {
	ChecksumAlgorithm *string `json:"checksum_algorithm" mapstructure:"checksum_algorithm"`
	ChecksumValue     *string `json:"checksum_value" mapstructure:"checksum_value"`
}

// ImageVersionResources The image version, which is composed of a product name and product version.
type ImageVersionResources struct {
	// Name of the producer/distribution of the image. For example windows or red hat.
	ProductName *string `json:"product_name" mapstructure:"product_name"`

	// Version *string for the disk image.
	ProductVersion *string `json:"product_version" mapstructure:"product_version"`
}

// ImageResources describes the image spec resources object.
type ImageResources struct {
	// The supported CPU architecture for a disk image.
	Architecture *string `json:"architecture,omitempty" mapstructure:"architecture,omitempty"`

	// Checksum of the image. The checksum is used for image validation if the image has a source specified. For images
	// that do not have their source specified the checksum is generated by the image service.
	Checksum *Checksum `json:"checksum,omitempty" mapstructure:"checksum,omitempty"`

	// The type of image.
	ImageType *string `json:"image_type,omitempty" mapstructure:"image_type,omitempty"`

	// The source URI points at the location of a the source image which is used to create/update image.
	SourceURI *string `json:"source_uri,omitempty" mapstructure:"source_uri,omitempty"`

	// The image version
	Version *ImageVersionResources `json:"version,omitempty" mapstructure:"version,omitempty"`

	// Cluster reference lists
	InitialPlacementRefList []*ReferenceValues `json:"initial_placement_ref_list,omitempty" mapstructure:"initial_placement_ref_list, omitempty"`

	// Reference to the source image such as 'vm_disk
	DataSourceReference *Reference `json:"data_source_reference,omitempty" mapstructure:"data_source_reference,omitempty"`
}

// Image An intentful representation of a image spec
type Image struct {
	// A description for image.
	Description *string `json:"description,omitempty" mapstructure:"description,omitempty"`

	// image Name.
	Name *string `json:"name,omitempty" mapstructure:"name,omitempty"`

	Resources *ImageResources `json:"resources" mapstructure:"resources"`
}

// ImageMetadata Metadata The image kind metadata
type ImageMetadata struct {
	// Categories for the image
	Categories map[string]string `json:"categories,omitempty" mapstructure:"categories,omitempty"`

	// UTC date and time in RFC-3339 format when vm was created
	CreationTime *time.Time `json:"creation_time,omitempty" mapstructure:"creation_time,omitempty"`

	// The kind name
	Kind *string `json:"kind" mapstructure:"kind"`

	// UTC date and time in RFC-3339 format when image was last updated
	LastUpdateTime *time.Time `json:"last_update_time,omitempty" mapstructure:"last_update_time,omitempty"`

	// image name
	Name *string `json:"name,omitempty" mapstructure:"name,omitempty"`

	// project reference
	ProjectReference *Reference `json:"project_reference,omitempty" mapstructure:"project_reference,omitempty"`

	OwnerReference *Reference `json:"owner_reference,omitempty" mapstructure:"owner_reference,omitempty"`

	// Hash of the spec. This will be returned from server.
	SpecHash *string `json:"spec_hash,omitempty" mapstructure:"spec_hash,omitempty"`

	// Version number of the latest spec.
	SpecVersion *int64 `json:"spec_version,omitempty" mapstructure:"spec_version,omitempty"`

	// image uuid
	UUID *string `json:"uuid,omitempty" mapstructure:"uuid,omitempty"`
}

// ImageIntentInput An intentful representation of a image
type ImageIntentInput struct {
	APIVersion *string `json:"api_version,omitempty" mapstructure:"api_version,omitempty"`

	Metadata *Metadata `json:"metadata,omitempty" mapstructure:"metadata,omitempty"`

	Spec *Image `json:"spec,omitempty" mapstructure:"spec,omitempty"`
}

// ImageStatus represents the status of a REST API call. Only used when there is a failure to report.
type ImageStatus struct {
	APIVersion *string `json:"api_version,omitempty" mapstructure:"api_version,omitempty"`

	// The HTTP error code.
	Code *int64 `json:"code,omitempty" mapstructure:"code,omitempty"`

	// The kind name
	Kind *string `json:"kind,omitempty" mapstructure:"kind,omitempty"`

	MessageList []*MessageResource `json:"message_list,omitempty" mapstructure:"message_list,omitempty"`

	State *string `json:"state,omitempty" mapstructure:"state,omitempty"`
}

// ImageVersionStatus represents the image version, which is composed of a product name and product version.
type ImageVersionStatus struct {
	// Name of the producer/distribution of the image. For example windows or red hat.
	ProductName *string `json:"product_name" mapstructure:"product_name"`

	// Version *string for the disk image.
	ProductVersion *string `json:"product_version" mapstructure:"product_version"`
}

// ImageResourcesDefStatus describes the image status resources object.
type ImageResourcesDefStatus struct {
	// The supported CPU architecture for a disk image.
	Architecture *string `json:"architecture,omitempty" mapstructure:"architecture,omitempty"`

	// Checksum of the image. The checksum is used for image validation if the image has a source specified. For images
	// that do not have their source specified the checksum is generated by the image service.
	Checksum *Checksum `json:"checksum,omitempty" mapstructure:"checksum,omitempty"`

	// The type of image.
	ImageType *string `json:"image_type,omitempty" mapstructure:"image_type,omitempty"`

	// List of URIs where the raw image data can be accessed.
	RetrievalURIList []*string `json:"retrieval_uri_list,omitempty" mapstructure:"retrieval_uri_list,omitempty"`

	// The size of the image in bytes.
	SizeBytes *int64 `json:"size_bytes,omitempty" mapstructure:"size_bytes,omitempty"`

	// The source URI points at the location of a the source image which is used to create/update image.
	SourceURI *string `json:"source_uri,omitempty" mapstructure:"source_uri,omitempty"`

	// Cluster reference lists
	InitialPlacementRefList []*ReferenceValues `json:"initial_placement_ref_list,omitempty" mapstructure:"initial_placement_ref_list, omitempty"`

	// cluster reference list when request was made without refs
	CurrentClusterReferenceList []*ReferenceValues `json:"current_cluster_reference_list,omitempty" mapstructure:"current_cluster_reference_list, omitempty"`

	// The image version
	Version *ImageVersionStatus `json:"version,omitempty" mapstructure:"version,omitempty"`
}

// ImageDefStatus represents an intentful representation of a image status
type ImageDefStatus struct {
	AvailabilityZoneReference *Reference `json:"availability_zone_reference,omitempty" mapstructure:"availability_zone_reference,omitempty"`

	ClusterReference *Reference `json:"cluster_reference,omitempty" mapstructure:"cluster_reference,omitempty"`

	// A description for image.
	Description *string `json:"description,omitempty" mapstructure:"description,omitempty"`

	// Any error messages for the image, if in an error state.
	MessageList []*MessageResource `json:"message_list,omitempty" mapstructure:"message_list,omitempty"`

	// image Name.
	Name *string `json:"name" mapstructure:"name"`

	Resources ImageResourcesDefStatus `json:"resources" mapstructure:"resources"`

	// The state of the image.
	State *string `json:"state,omitempty" mapstructure:"state,omitempty"`

	ExecutionContext *ExecutionContext `json:"execution_context,omitempty" mapstructure:"execution_context,omitempty"`
}

// ImageIntentResponse represents the response object for intentful operations on a image
type ImageIntentResponse struct {
	APIVersion *string `json:"api_version" mapstructure:"api_version"`

	Metadata *Metadata `json:"metadata" mapstructure:"metadata"`

	Spec *Image `json:"spec,omitempty" mapstructure:"spec,omitempty"`

	Status *ImageDefStatus `json:"status,omitempty" mapstructure:"status,omitempty"`
}

// ImageListMetadata represents metadata input
type ImageListMetadata struct {
	// The filter in FIQL syntax used for the results.
	Filter *string `json:"filter,omitempty" mapstructure:"filter,omitempty"`

	// The kind name
	Kind *string `json:"kind,omitempty" mapstructure:"kind,omitempty"`

	// The number of records to retrieve relative to the offset
	Length *int64 `json:"length,omitempty" mapstructure:"length,omitempty"`

	// Offset from the start of the entity list
	Offset *int64 `json:"offset,omitempty" mapstructure:"offset,omitempty"`

	// The attribute to perform sort on
	SortAttribute *string `json:"sort_attribute,omitempty" mapstructure:"sort_attribute,omitempty"`

	// The sort order in which results are returned
	SortOrder *string `json:"sort_order,omitempty" mapstructure:"sort_order,omitempty"`
}

// ImageIntentResource represents the response object for intentful operations on a image
type ImageIntentResource struct {
	APIVersion *string `json:"api_version,omitempty" mapstructure:"api_version,omitempty"`

	Metadata *Metadata `json:"metadata" mapstructure:"metadata"`

	Spec *Image `json:"spec,omitempty" mapstructure:"spec,omitempty"`

	Status *ImageDefStatus `json:"status,omitempty" mapstructure:"status,omitempty"`
}

// ImageListIntentResponse represents the response object for intentful operation of images
type ImageListIntentResponse struct {
	APIVersion *string `json:"api_version" mapstructure:"api_version"`

	Entities []*ImageIntentResponse `json:"entities,omitempty" mapstructure:"entities,omitempty"`

	Metadata *ListMetadataOutput `json:"metadata" mapstructure:"metadata"`
}

// ClusterListIntentResponse ...
type ClusterListIntentResponse struct {
	APIVersion *string                  `json:"api_version" mapstructure:"api_version"`
	Entities   []*ClusterIntentResponse `json:"entities,omitempty" mapstructure:"entities,omitempty"`
	Metadata   *ListMetadataOutput      `json:"metadata" mapstructure:"metadata"`
}

// ClusterIntentResponse ...
type ClusterIntentResponse struct {
	APIVersion *string `json:"api_version,omitempty" mapstructure:"api_version,omitempty"`

	Metadata *Metadata `json:"metadata" mapstructure:"metadata"`

	Spec *Cluster `json:"spec,omitempty" mapstructure:"spec,omitempty"`

	Status *ClusterDefStatus `json:"status,omitempty" mapstructure:"status,omitempty"`
}

// Cluster ...
type Cluster struct {
	Name      *string          `json:"name,omitempty" mapstructure:"name,omitempty"`
	Resources *ClusterResource `json:"resources,omitempty" mapstructure:"resources,omitempty"`
}

// ClusterDefStatus ...
type ClusterDefStatus struct {
	State       *string            `json:"state,omitempty" mapstructure:"state,omitempty"`
	MessageList []*MessageResource `json:"message_list,omitempty" mapstructure:"message_list,omitempty"`
	Name        *string            `json:"name,omitempty" mapstructure:"name,omitempty"`
	Resources   *ClusterObj        `json:"resources,omitempty" mapstructure:"resources,omitempty"`
}

// ClusterObj ...
type ClusterObj struct {
	Nodes             *ClusterNodes    `json:"nodes,omitempty" mapstructure:"nodes,omitempty"`
	Config            *ClusterConfig   `json:"config,omitempty" mapstructure:"config,omitempty"`
	Network           *ClusterNetwork  `json:"network,omitempty" mapstructure:"network,omitempty"`
	Analysis          *ClusterAnalysis `json:"analysis,omitempty" mapstructure:"analysis,omitempty"`
	RuntimeStatusList []*string        `json:"runtime_status_list,omitempty" mapstructure:"runtime_status_list,omitempty"`
}

// ClusterNodes ...
type ClusterNodes struct {
	HypervisorServerList []*HypervisorServer `json:"hypervisor_server_list,omitempty" mapstructure:"hypervisor_server_list,omitempty"`
}

// SoftwareMapValues ...
type SoftwareMapValues struct {
	SoftwareType *string `json:"software_type,omitempty" mapstructure:"software_type,omitempty"`
	Status       *string `json:"status,omitempty" mapstructure:"status,omitempty"`
	Version      *string `json:"version,omitempty" mapstructure:"version,omitempty"`
}

// SoftwareMap ...
type SoftwareMap struct {
	NCC *SoftwareMapValues `json:"ncc,omitempty" mapstructure:"ncc,omitempty"`
	NOS *SoftwareMapValues `json:"nos,omitempty" mapstructure:"nos,omitempty"`
}

// ClusterConfig ...
type ClusterConfig struct {
	GpuDriverVersion              *string                    `json:"gpu_driver_version,omitempty" mapstructure:"gpu_driver_version,omitempty"`
	ClientAuth                    *ClientAuth                `json:"client_auth,omitempty" mapstructure:"client_auth,omitempty"`
	AuthorizedPublicKeyList       []*PublicKey               `json:"authorized_public_key_list,omitempty" mapstructure:"authorized_public_key_list,omitempty"`
	SoftwareMap                   *SoftwareMap               `json:"software_map,omitempty" mapstructure:"software_map,omitempty"`
	EncryptionStatus              *string                    `json:"encryption_status,omitempty" mapstructure:"encryption_status,omitempty"`
	SslKey                        *SslKey                    `json:"ssl_key,omitempty" mapstructure:"ssl_key,omitempty"`
	ServiceList                   []*string                  `json:"service_list,omitempty" mapstructure:"service_list,omitempty"`
	SupportedInformationVerbosity *string                    `json:"supported_information_verbosity,omitempty" mapstructure:"supported_information_verbosity,omitempty"`
	CertificationSigningInfo      *CertificationSigningInfo  `json:"certification_signing_info,omitempty" mapstructure:"certification_signing_info,omitempty"`
	RedundancyFactor              *int64                     `json:"redundancy_factor,omitempty" mapstructure:"redundancy_factor,omitempty"`
	ExternalConfigurations        *ExternalConfigurations    `json:"external_configurations,omitempty" mapstructure:"external_configurations,omitempty"`
	OperationMode                 *string                    `json:"operation_mode,omitempty" mapstructure:"operation_mode,omitempty"`
	CaCertificateList             []*CaCert                  `json:"ca_certificate_list,omitempty" mapstructure:"ca_certificate_list,omitempty"`
	EnabledFeatureList            []*string                  `json:"enabled_feature_list,omitempty" mapstructure:"enabled_feature_list,omitempty"`
	IsAvailable                   *bool                      `json:"is_available,omitempty" mapstructure:"is_available,omitempty"`
	Build                         *BuildInfo                 `json:"build,omitempty" mapstructure:"build,omitempty"`
	Timezone                      *string                    `json:"timezone,omitempty" mapstructure:"timezone,omitempty"`
	ClusterArch                   *string                    `json:"cluster_arch,omitempty" mapstructure:"cluster_arch,omitempty"`
	ManagementServerList          []*ClusterManagementServer `json:"management_server_list,omitempty" mapstructure:"management_server_list,omitempty"`
}

// ClusterManagementServer ...
type ClusterManagementServer struct {
	IP         *string   `json:"ip,omitempty" mapstructure:"ip,omitempty"`
	DrsEnabled *bool     `json:"drs_enabled,omitempty" mapstructure:"drs_enabled,omitempty"`
	StatusList []*string `json:"status_list,omitempty" mapstructure:"status_list,omitempty"`
	Type       *string   `json:"type,omitempty" mapstructure:"type,omitempty"`
}

// BuildInfo ...
type BuildInfo struct {
	CommitID      *string `json:"commit_id,omitempty" mapstructure:"commit_id,omitempty"`
	FullVersion   *string `json:"full_version,omitempty" mapstructure:"full_version,omitempty"`
	CommitDate    *string `json:"commit_date,omitempty" mapstructure:"commit_date,omitempty"`
	Version       *string `json:"version,omitempty" mapstructure:"version,omitempty"`
	ShortCommitID *string `json:"short_commit_id,omitempty" mapstructure:"short_commit_id,omitempty"`
	BuildType     *string `json:"build_type,omitempty" mapstructure:"build_type,omitempty"`
}

// CaCert ...
type CaCert struct {
	CaName      *string `json:"ca_name,omitempty" mapstructure:"ca_name,omitempty"`
	Certificate *string `json:"certificate,omitempty" mapstructure:"certificate,omitempty"`
}

// ExternalConfigurations ...
type ExternalConfigurations struct {
	CitrixConnectorConfig *CitrixConnectorConfigDetails `json:"citrix_connector_config,omitempty" mapstructure:"citrix_connector_config,omitempty"`
}

// CitrixConnectorConfigDetails ...
type CitrixConnectorConfigDetails struct {
	CitrixVMReferenceList *[]Reference            `json:"citrix_vm_reference_list,omitempty" mapstructure:"citrix_vm_reference_list,omitempty"`
	ClientSecret          *string                 `json:"client_secret,omitempty" mapstructure:"client_secret,omitempty"`
	CustomerID            *string                 `json:"customer_id,omitempty" mapstructure:"customer_id,omitempty"`
	ClientID              *string                 `json:"client_id,omitempty" mapstructure:"client_id,omitempty"`
	ResourceLocation      *CitrixResourceLocation `json:"resource_location,omitempty" mapstructure:"resource_location,omitempty"`
}

// CitrixResourceLocation ...
type CitrixResourceLocation struct {
	ID   *string `json:"id,omitempty" mapstructure:"id,omitempty"`
	Name *string `json:"name,omitempty" mapstructure:"name,omitempty"`
}

// SslKey ...
type SslKey struct {
	KeyType        *string                   `json:"key_type,omitempty" mapstructure:"key_type,omitempty"`
	KeyName        *string                   `json:"key_name,omitempty" mapstructure:"key_name,omitempty"`
	SigningInfo    *CertificationSigningInfo `json:"signing_info,omitempty" mapstructure:"signing_info,omitempty"`
	ExpireDatetime *string                   `json:"expire_datetime,omitempty" mapstructure:"expire_datetime,omitempty"`
}

// CertificationSigningInfo ...
type CertificationSigningInfo struct {
	City             *string `json:"city,omitempty" mapstructure:"city,omitempty"`
	CommonNameSuffix *string `json:"common_name_suffix,omitempty" mapstructure:"common_name_suffix,omitempty"`
	State            *string `json:"state,omitempty" mapstructure:"state,omitempty"`
	CountryCode      *string `json:"country_code,omitempty" mapstructure:"country_code,omitempty"`
	CommonName       *string `json:"common_name,omitempty" mapstructure:"common_name,omitempty"`
	Organization     *string `json:"organization,omitempty" mapstructure:"organization,omitempty"`
	EmailAddress     *string `json:"email_address,omitempty" mapstructure:"email_address,omitempty"`
}

// PublicKey ...
type PublicKey struct {
	Key  *string `json:"key,omitempty" mapstructure:"key,omitempty"`
	Name *string `json:"name,omitempty" mapstructure:"name,omitempty"`
}

// ClientAuth ...
type ClientAuth struct {
	Status  *string `json:"status,omitempty" mapstructure:"status,omitempty"`
	CaChain *string `json:"ca_chain,omitempty" mapstructure:"ca_chain,omitempty"`
	Name    *string `json:"name,omitempty" mapstructure:"name,omitempty"`
}

// HypervisorServer ...
type HypervisorServer struct {
	IP      *string `json:"ip,omitempty" mapstructure:"ip,omitempty"`
	Version *string `json:"version,omitempty" mapstructure:"version,omitempty"`
	Type    *string `json:"type,omitempty" mapstructure:"type,omitempty"`
}

// ClusterResource ...
type ClusterResource struct {
	Config            *ConfigClusterSpec `json:"config,omitempty" mapstructure:"config,omitempty"`
	Network           *ClusterNetwork    `json:"network,omitempty" mapstructure:"network,omitempty"`
	RunTimeStatusList []*string          `json:"runtime_status_list,omitempty" mapstructure:"runtime_status_list,omitempty"`
}

// ConfigClusterSpec ...
type ConfigClusterSpec struct {
	GpuDriverVersion              *string                     `json:"gpu_driver_version,omitempty" mapstructure:"gpu_driver_version,omitempty"`
	ClientAuth                    *ClientAuth                 `json:"client_auth,omitempty" mapstructure:"client_auth,omitempty"`
	AuthorizedPublicKeyList       []*PublicKey                `json:"authorized_public_key_list,omitempty" mapstructure:"authorized_public_key_list,omitempty"`
	SoftwareMap                   map[string]interface{}      `json:"software_map,omitempty" mapstructure:"software_map,omitempty"`
	EncryptionStatus              string                      `json:"encryption_status,omitempty" mapstructure:"encryption_status,omitempty"`
	RedundancyFactor              *int64                      `json:"redundancy_factor,omitempty" mapstructure:"redundancy_factor,omitempty"`
	CertificationSigningInfo      *CertificationSigningInfo   `json:"certification_signing_info,omitempty" mapstructure:"certification_signing_info,omitempty"`
	SupportedInformationVerbosity *string                     `json:"supported_information_verbosity,omitempty" mapstructure:"supported_information_verbosity,omitempty"`
	ExternalConfigurations        *ExternalConfigurationsSpec `json:"external_configurations,omitempty" mapstructure:"external_configurations,omitempty"`
	EnabledFeatureList            []*string                   `json:"enabled_feature_list,omitempty" mapstructure:"enabled_feature_list,omitempty"`
	Timezone                      *string                     `json:"timezone,omitempty" mapstructure:"timezone,omitempty"`
	OperationMode                 *string                     `json:"operation_mode,omitempty" mapstructure:"operation_mode,omitempty"`
}

// ExternalConfigurationsSpec ...
type ExternalConfigurationsSpec struct {
	CitrixConnectorConfig *CitrixConnectorConfigDetailsSpec `json:"citrix_connector_config,omitempty" mapstructure:"citrix_connector_config,omitempty"`
}

// CitrixConnectorConfigDetailsSpec ...
type CitrixConnectorConfigDetailsSpec struct {
	CitrixVMReferenceList []*Reference                `json:"citrix_connector_config,omitempty" mapstructure:"citrix_connector_config,omitempty"`
	ClientSecret          *string                     `json:"client_secret,omitempty" mapstructure:"client_secret,omitempty"`
	CustomerID            *string                     `json:"customer_id,omitempty" mapstructure:"customer_id,omitempty"`
	ClientID              *string                     `json:"client_id,omitempty" mapstructure:"client_id,omitempty"`
	ResourceLocation      *CitrixResourceLocationSpec `json:"resource_location,omitempty" mapstructure:"resource_location,omitempty"`
}

// CitrixResourceLocationSpec ...
type CitrixResourceLocationSpec struct {
	ID   *string `json:"id,omitempty" mapstructure:"id,omitempty"`
	Name *string `json:"name,omitempty" mapstructure:"name,omitempty"`
}

// ClusterNetwork ...
type ClusterNetwork struct {
	MasqueradingPort       *int64                  `json:"masquerading_port,omitempty" mapstructure:"masquerading_port,omitempty"`
	MasqueradingIP         *string                 `json:"masquerading_ip,omitempty" mapstructure:"masquerading_ip,omitempty"`
	ExternalIP             *string                 `json:"external_ip,omitempty" mapstructure:"external_ip,omitempty"`
	HTTPProxyList          []*ClusterNetworkEntity `json:"http_proxy_list,omitempty" mapstructure:"http_proxy_list,omitempty"`
	SMTPServer             *SMTPServer             `json:"smtp_server,omitempty" mapstructure:"smtp_server,omitempty"`
	NTPServerIPList        []*string               `json:"ntp_server_ip_list,omitempty" mapstructure:"ntp_server_ip_list,omitempty"`
	ExternalSubnet         *string                 `json:"external_subnet,omitempty" mapstructure:"external_subnet,omitempty"`
	NFSSubnetWhitelist     []*string               `json:"nfs_subnet_whitelist,omitempty" mapstructure:"nfs_subnet_whitelist,omitempty"`
	ExternalDataServicesIP *string                 `json:"external_data_services_ip,omitempty" mapstructure:"external_data_services_ip,omitempty"`
	DomainServer           *ClusterDomainServer    `json:"domain_server,omitempty" mapstructure:"domain_server,omitempty"`
	NameServerIPList       []*string               `json:"name_server_ip_list,omitempty" mapstructure:"name_server_ip_list,omitempty"`
	HTTPProxyWhitelist     []*HTTPProxyWhitelist   `json:"http_proxy_whitelist,omitempty" mapstructure:"http_proxy_whitelist,omitempty"`
	InternalSubnet         *string                 `json:"internal_subnet,omitempty" mapstructure:"internal_subnet,omitempty"`
}

// HTTPProxyWhitelist ...
type HTTPProxyWhitelist struct {
	Target     *string `json:"target,omitempty" mapstructure:"target,omitempty"`
	TargetType *string `json:"target_type,omitempty" mapstructure:"target_type,omitempty"`
}

// ClusterDomainServer ...
type ClusterDomainServer struct {
	Nameserver        *string      `json:"nameserver,omitempty" mapstructure:"nameserver,omitempty"`
	Name              *string      `json:"name,omitempty" mapstructure:"name,omitempty"`
	DomainCredentials *Credentials `json:"external_data_services_ip,omitempty" mapstructure:"external_data_services_ip,omitempty"`
}

// SMTPServer ...
type SMTPServer struct {
	Type         *string               `json:"type,omitempty" mapstructure:"type,omitempty"`
	EmailAddress *string               `json:"email_address,omitempty" mapstructure:"email_address,omitempty"`
	Server       *ClusterNetworkEntity `json:"server,omitempty" mapstructure:"server,omitempty"`
}

// ClusterNetworkEntity ...
type ClusterNetworkEntity struct {
	Credentials   *Credentials `json:"credentials,omitempty" mapstructure:"credentials,omitempty"`
	ProxyTypeList []*string    `json:"proxy_type_list,omitempty" mapstructure:"proxy_type_list,omitempty"`
	Address       *Address     `json:"address,omitempty" mapstructure:"address,omitempty"`
}

// Credentials ...
type Credentials struct {
	Username *string `json:"username,omitempty" mapstructure:"username,omitempty"`
	Password *string `json:"password,omitempty" mapstructure:"password,omitempty"`
}

// VMEfficiencyMap ...
type VMEfficiencyMap struct {
	BullyVMNum           *string `json:"bully_vm_num,omitempty" mapstructure:"bully_vm_num,omitempty"`
	ConstrainedVMNum     *string `json:"constrained_vm_num,omitempty" mapstructure:"constrained_vm_num,omitempty"`
	DeadVMNum            *string `json:"dead_vm_num,omitempty" mapstructure:"dead_vm_num,omitempty"`
	InefficientVMNum     *string `json:"inefficient_vm_num,omitempty" mapstructure:"inefficient_vm_num,omitempty"`
	OverprovisionedVMNum *string `json:"overprovisioned_vm_num,omitempty" mapstructure:"overprovisioned_vm_num,omitempty"`
}

// ClusterAnalysis ...
type ClusterAnalysis struct {
	VMEfficiencyMap *VMEfficiencyMap `json:"vm_efficiency_map,omitempty" mapstructure:"vm_efficiency_map,omitempty"`
}

// CategoryListMetadata All api calls that return a list will have this metadata block as input
type CategoryListMetadata struct {
	// The filter in FIQL syntax used for the results.
	Filter *string `json:"filter,omitempty" mapstructure:"filter,omitempty"`

	// The kind name
	Kind *string `json:"kind,omitempty" mapstructure:"kind,omitempty"`

	// The number of records to retrieve relative to the offset
	Length *int64 `json:"length,omitempty" mapstructure:"length,omitempty"`

	// Offset from the start of the entity list
	Offset *int64 `json:"offset,omitempty" mapstructure:"offset,omitempty"`

	// The attribute to perform sort on
	SortAttribute *string `json:"sort_attribute,omitempty" mapstructure:"sort_attribute,omitempty"`

	// The sort order in which results are returned
	SortOrder *string `json:"sort_order,omitempty" mapstructure:"sort_order,omitempty"`

	// Total number of matched results.
	TotalMatches *int64 `json:"total_matches,omitempty" mapstructure:"total_matches,omitempty"`
}

// CategoryKeyStatus represents Category Key Definition.
type CategoryKeyStatus struct {
	// API version.
	APIVersion *string `json:"api_version,omitempty" mapstructure:"api_version,omitempty"`

	// Description of the category.
	Description *string `json:"description,omitempty" mapstructure:"description,omitempty"`

	// Name of the category.
	Name *string `json:"name" mapstructure:"name"`

	// Specifying whether its a system defined category.
	SystemDefined *bool `json:"system_defined,omitempty" mapstructure:"system_defined,omitempty"`
}

// CategoryKeyListResponse represents the category key list response.
type CategoryKeyListResponse struct {
	// API Version.
	APIVersion *string `json:"api_version,omitempty" mapstructure:"api_version,omitempty"`

	Entities []*CategoryKeyStatus `json:"entities,omitempty" mapstructure:"entities,omitempty"`

	Metadata *CategoryListMetadata `json:"metadata,omitempty" mapstructure:"metadata,omitempty"`
}

// CategoryKey represents category key definition.
type CategoryKey struct {
	// API version.
	APIVersion *string `json:"api_version,omitempty" mapstructure:"api_version,omitempty"`

	// Description of the category.
	Description *string `json:"description,omitempty" mapstructure:"description,omitempty"`

	// Name of the category.
	Name *string `json:"name" mapstructure:"name"`
}

// CategoryStatus represents The status of a REST API call. Only used when there is a failure to report.
type CategoryStatus struct {
	APIVersion *string `json:"api_version,omitempty" mapstructure:"api_version,omitempty"`

	// The HTTP error code.
	Code *int64 `json:"code,omitempty" mapstructure:"code,omitempty"`

	// The kind name
	Kind *string `json:"kind,omitempty" mapstructure:"kind,omitempty"`

	MessageList []*MessageResource `json:"message_list,omitempty" mapstructure:"message_list,omitempty"`

	State *string `json:"state,omitempty" mapstructure:"state,omitempty"`
}

// CategoryValueListResponse represents Category Value list response.
type CategoryValueListResponse struct {
	APIVersion *string `json:"api_version,omitempty" mapstructure:"api_version,omitempty"`

	Entities []*CategoryValueStatus `json:"entities,omitempty" mapstructure:"entities,omitempty"`

	Metadata *CategoryListMetadata `json:"metadata,omitempty" mapstructure:"metadata,omitempty"`
}

// CategoryValueStatus represents Category value definition.
type CategoryValueStatus struct {
	// API version.
	APIVersion *string `json:"api_version,omitempty" mapstructure:"api_version,omitempty"`

	// Description of the category value.
	Description *string `json:"description,omitempty" mapstructure:"description,omitempty"`

	// The name of the category.
	Name *string `json:"name,omitempty" mapstructure:"name,omitempty"`

	// Specifying whether its a system defined category.
	SystemDefined *bool `json:"system_defined,omitempty" mapstructure:"system_defined,omitempty"`

	// The value of the category.
	Value *string `json:"value,omitempty" mapstructure:"value,omitempty"`
}

// CategoryFilter represents A category filter.
type CategoryFilter struct {
	// List of kinds associated with this filter.
	KindList []*string `json:"kind_list,omitempty" mapstructure:"kind_list,omitempty"`

	// A list of category key and list of values.
	Params map[string][]string `json:"params,omitempty" mapstructure:"params,omitempty"`

	// The type of the filter being used.
	Type *string `json:"type,omitempty" mapstructure:"type,omitempty"`
}

// CategoryQueryInput represents Categories query input object.
type CategoryQueryInput struct {
	// API version.
	APIVersion *string `json:"api_version,omitempty" mapstructure:"api_version,omitempty"`

	CategoryFilter *CategoryFilter `json:"category_filter,omitempty" mapstructure:"category_filter,omitempty"`

	// The maximum number of members to return per group.
	GroupMemberCount *int64 `json:"group_member_count,omitempty" mapstructure:"group_member_count,omitempty"`

	// The offset into the total member set to return per group.
	GroupMemberOffset *int64 `json:"group_member_offset,omitempty" mapstructure:"group_member_offset,omitempty"`

	// TBD: USED_IN - to get policies in which specified categories are used. APPLIED_TO - to get entities attached to
	// specified categories.
	UsageType *string `json:"usage_type,omitempty" mapstructure:"usage_type,omitempty"`
}

// CategoryQueryResponseMetadata represents Response metadata.
type CategoryQueryResponseMetadata struct {
	// The maximum number of records to return per group.
	GroupMemberCount *int64 `json:"group_member_count,omitempty" mapstructure:"group_member_count,omitempty"`

	// The offset into the total records set to return per group.
	GroupMemberOffset *int64 `json:"group_member_offset,omitempty" mapstructure:"group_member_offset,omitempty"`

	// Total number of matched results.
	TotalMatches *int64 `json:"total_matches,omitempty" mapstructure:"total_matches,omitempty"`

	// TBD: USED_IN - to get policies in which specified categories are used. APPLIED_TO - to get entities attached to specified categories.
	UsageType *string `json:"usage_type,omitempty" mapstructure:"usage_type,omitempty"`
}

// EntityReference Reference to an entity.
type EntityReference struct {
	// Categories for the entity.
	Categories map[string]string `json:"categories,omitempty" mapstructure:"categories,omitempty"`

	// Kind of the reference.
	Kind *string `json:"kind,omitempty" mapstructure:"kind,omitempty"`

	// Name of the entity.
	Name *string `json:"name,omitempty" mapstructure:"name,omitempty"`

	// The type of filter being used. (Options : CATEGORIES_MATCH_ALL , CATEGORIES_MATCH_ANY)
	Type *string `json:"type,omitempty" mapstructure:"type,omitempty"`

	// UUID of the entity.
	UUID *string `json:"uuid,omitempty" mapstructure:"uuid,omitempty"`
}

// CategoryQueryResponseResults ...
type CategoryQueryResponseResults struct {
	// List of entity references.
	EntityAnyReferenceList []*EntityReference `json:"entity_any_reference_list,omitempty" mapstructure:"entity_any_reference_list,omitempty"`

	// Total number of filtered results.
	FilteredEntityCount *int64 `json:"filtered_entity_count,omitempty" mapstructure:"filtered_entity_count,omitempty"`

	// The entity kind.
	Kind *string `json:"kind,omitempty" mapstructure:"kind,omitempty"`

	// Total number of the matched results.
	TotalEntityCount *int64 `json:"total_entity_count,omitempty" mapstructure:"total_entity_count,omitempty"`
}

// CategoryQueryResponse represents Categories query response object.
type CategoryQueryResponse struct {
	// API version.
	APIVersion *string `json:"api_version,omitempty" mapstructure:"api_version,omitempty"`

	Metadata *CategoryQueryResponseMetadata `json:"metadata,omitempty" mapstructure:"metadata,omitempty"`

	Results []*CategoryQueryResponseResults `json:"results,omitempty" mapstructure:"results,omitempty"`
}

// CategoryValue represents Category value definition.
type CategoryValue struct {
	// API version.
	APIVersion *string `json:"api_version,omitempty" mapstructure:"api_version,omitempty"`

	// Description of the category value.
	Description *string `json:"description,omitempty" `

	// Value for the category.
	Value *string `json:"value,omitempty" mapstructure:"value,omitempty"`
}

// PortRange represents Range of TCP/UDP ports.
type PortRange struct {
	EndPort *int64 `json:"end_port,omitempty" mapstructure:"end_port,omitempty"`

	StartPort *int64 `json:"start_port,omitempty" mapstructure:"start_port,omitempty"`
}

// IPSubnet IP subnet provided as an address and prefix length.
type IPSubnet struct {
	// IPV4 address.
	IP *string `json:"ip,omitempty" mapstructure:"ip,omitempty"`

	PrefixLength *int64 `json:"prefix_length,omitempty" mapstructure:"prefix_length,omitempty"`
}

// NetworkRuleIcmpTypeCodeList ..
type NetworkRuleIcmpTypeCodeList struct {
	Code *int64 `json:"code,omitempty" mapstructure:"code,omitempty"`

	Type *int64 `json:"type,omitempty" mapstructure:"type,omitempty"`
}

// NetworkRule ...
type NetworkRule struct {
	ExpirationTime                *string                        `json:"expiration_time,omitempty" mapstructure:"expiration_time,omitempty"`
	Filter                        *CategoryFilter                `json:"filter,omitempty" mapstructure:"filter,omitempty"`
	IcmpTypeCodeList              []*NetworkRuleIcmpTypeCodeList `json:"icmp_type_code_list,omitempty" mapstructure:"icmp_type_code_list,omitempty"`
	IPSubnet                      *IPSubnet                      `json:"ip_subnet,omitempty" mapstructure:"ip_subnet,omitempty"`
	NetworkFunctionChainReference *Reference                     `json:"network_function_chain_reference,omitempty" mapstructure:"network_function_chain_reference,omitempty"`
	PeerSpecificationType         *string                        `json:"peer_specification_type,omitempty" mapstructure:"peer_specification_type,omitempty"`
	Protocol                      *string                        `json:"protocol,omitempty" mapstructure:"protocol,omitempty"`
	TCPPortRangeList              []*PortRange                   `json:"tcp_port_range_list,omitempty" mapstructure:"tcp_port_range_list,omitempty"`
	UDPPortRangeList              []*PortRange                   `json:"udp_port_range_list,omitempty" mapstructure:"udp_port_range_list,omitempty"`
	AddressGroupInclusionList     []*Reference                   `json:"address_group_inclusion_list,omitempty" mapstructure:"address_group_inclusion_list,omitempty"`
	Description                   *string                        `json:"description,omitempty" mapstructure:"description,omitempty"`
	ServiceGroupList              []*Reference                   `json:"service_group_list,omitempty" mapstructure:"service_group_list,omitempty"`
}

// TargetGroup ...
type TargetGroup struct {
	// Default policy for communication within target group.
	DefaultInternalPolicy *string `json:"default_internal_policy,omitempty" mapstructure:"default_internal_policy,omitempty"`

	// The set of categories that matching VMs need to have.
	Filter *CategoryFilter `json:"filter,omitempty" mapstructure:"filter,omitempty"`

	// Way to identify the object for which rule is applied.
	PeerSpecificationType *string `json:"peer_specification_type,omitempty" mapstructure:"peer_specification_type,omitempty"`
}

// NetworkSecurityRuleResourcesRule These rules are used for quarantining suspected VMs. Target group is a required
// attribute.  Empty inbound_allow_list will not allow anything into target group. Empty outbound_allow_list will allow
// everything from target group.
type NetworkSecurityRuleResourcesRule struct {
	Action            *string        `json:"action,omitempty" mapstructure:"action,omitempty"`                         // Type of action.
	InboundAllowList  []*NetworkRule `json:"inbound_allow_list,omitempty" mapstructure:"inbound_allow_list,omitempty"` //
	OutboundAllowList []*NetworkRule `json:"outbound_allow_list,omitempty" mapstructure:"outbound_allow_list,omitempty"`
	TargetGroup       *TargetGroup   `json:"target_group,omitempty" mapstructure:"target_group,omitempty"`
}

// NetworkSecurityRuleIsolationRule These rules are used for environmental isolation.
type NetworkSecurityRuleIsolationRule struct {
	Action             *string         `json:"action,omitempty" mapstructure:"action,omitempty"`                             // Type of action.
	FirstEntityFilter  *CategoryFilter `json:"first_entity_filter,omitempty" mapstructure:"first_entity_filter,omitempty"`   // The set of categories that matching VMs need to have.
	SecondEntityFilter *CategoryFilter `json:"second_entity_filter,omitempty" mapstructure:"second_entity_filter,omitempty"` // The set of categories that matching VMs need to have.
}

// NetworkSecurityRuleResources ...
type NetworkSecurityRuleResources struct {
	AllowIpv6Traffic      *bool                             `json:"allow_ipv6_traffic,omitempty" mapstructure:"allow_ipv6_traffic,omitempty"`
	IsPolicyHitlogEnabled *bool                             `json:"is_policy_hitlog_enabled,omitempty" mapstructure:"is_policy_hitlog_enabled,omitempty"`
	AdRule                *NetworkSecurityRuleResourcesRule `json:"ad_rule,omitempty" mapstructure:"ad_rule,omitempty"`
	AppRule               *NetworkSecurityRuleResourcesRule `json:"app_rule,omitempty" mapstructure:"app_rule,omitempty"`
	IsolationRule         *NetworkSecurityRuleIsolationRule `json:"isolation_rule,omitempty" mapstructure:"isolation_rule,omitempty"`
	QuarantineRule        *NetworkSecurityRuleResourcesRule `json:"quarantine_rule,omitempty" mapstructure:"quarantine_rule,omitempty"`
}

// NetworkSecurityRule ...
type NetworkSecurityRule struct {
	Description *string                       `json:"description" mapstructure:"description"`
	Name        *string                       `json:"name,omitempty" mapstructure:"name,omitempty"`
	Resources   *NetworkSecurityRuleResources `json:"resources,omitempty" `
}

// Metadata Metadata The kind metadata
type Metadata struct {
	LastUpdateTime   *time.Time        `json:"last_update_time,omitempty" mapstructure:"last_update_time,omitempty"`   //
	Kind             *string           `json:"kind" mapstructure:"kind"`                                               //
	UUID             *string           `json:"uuid,omitempty" mapstructure:"uuid,omitempty"`                           //
	ProjectReference *Reference        `json:"project_reference,omitempty" mapstructure:"project_reference,omitempty"` // project reference
	CreationTime     *time.Time        `json:"creation_time,omitempty" mapstructure:"creation_time,omitempty"`
	SpecVersion      *int64            `json:"spec_version,omitempty" mapstructure:"spec_version,omitempty"`
	SpecHash         *string           `json:"spec_hash,omitempty" mapstructure:"spec_hash,omitempty"`
	OwnerReference   *Reference        `json:"owner_reference,omitempty" mapstructure:"owner_reference,omitempty"`
	Categories       map[string]string `json:"categories,omitempty" mapstructure:"categories,omitempty"`
	Name             *string           `json:"name,omitempty" mapstructure:"name,omitempty"`

	// Applied on Prism Central only. Indicate whether force to translate the spec of the fanout request to fit the target cluster API schema.
	ShouldForceTranslate *bool `json:"should_force_translate,omitempty" mapstructure:"should_force_translate,omitempty"`

	CategoriesMapping    map[string][]string `json:"categories_mapping,omitempty" mapstructure:"categories_mapping,omitempty"`
	EntityVersion        *string             `json:"entity_version,omitempty" mapstructure:"entity_version,omitempty"`
	UseCategoriesMapping *bool               `json:"use_categories_mapping,omitempty" mapstructure:"use_categories_mapping,omitempty"`
}

// NetworkSecurityRuleIntentInput An intentful representation of a network_security_rule
type NetworkSecurityRuleIntentInput struct {
	APIVersion *string              `json:"api_version,omitempty" mapstructure:"api_version,omitempty"`
	Metadata   *Metadata            `json:"metadata" mapstructure:"metadata"`
	Spec       *NetworkSecurityRule `json:"spec" mapstructure:"spec"`
}

// NetworkSecurityRuleDefStatus ... Network security rule status
type NetworkSecurityRuleDefStatus struct {
	Resources        *NetworkSecurityRuleResources `json:"resources,omitempty" mapstructure:"resources,omitempty"`
	State            *string                       `json:"state,omitempty" mapstructure:"state,omitempty"`
	ExecutionContext *ExecutionContext             `json:"execution_context,omitempty" mapstructure:"execution_context,omitempty"`
	Name             *string                       `json:"name,omitempty" mapstructure:"name,omitempty"`
	Description      *string                       `json:"description,omitempty" mapstructure:"description,omitempty"`
}

// NetworkSecurityRuleIntentResponse Response object for intentful operations on a network_security_rule
type NetworkSecurityRuleIntentResponse struct {
	APIVersion *string                       `json:"api_version,omitempty" mapstructure:"api_version,omitempty"`
	Metadata   *Metadata                     `json:"metadata" mapstructure:"metadata"`
	Spec       *NetworkSecurityRule          `json:"spec,omitempty" mapstructure:"spec,omitempty"`
	Status     *NetworkSecurityRuleDefStatus `json:"status,omitempty" mapstructure:"status,omitempty"`
}

// NetworkSecurityRuleStatus The status of a REST API call. Only used when there is a failure to report.
type NetworkSecurityRuleStatus struct {
	APIVersion  *string            `json:"api_version,omitempty" mapstructure:"api_version,omitempty"` //
	Code        *int64             `json:"code,omitempty" mapstructure:"code,omitempty"`               // The HTTP error code.
	Kind        *string            `json:"kind,omitempty" mapstructure:"kind,omitempty"`               // The kind name
	MessageList []*MessageResource `json:"message_list,omitempty" mapstructure:"message_list,omitempty"`
	State       *string            `json:"state,omitempty" mapstructure:"state,omitempty"`
}

// ListMetadata All api calls that return a list will have this metadata block as input
type ListMetadata struct {
	Filter        *string `json:"filter,omitempty" mapstructure:"filter,omitempty"`                 // The filter in FIQL syntax used for the results.
	Kind          *string `json:"kind,omitempty" mapstructure:"kind,omitempty"`                     // The kind name
	Length        *int64  `json:"length,omitempty" mapstructure:"length,omitempty"`                 // The number of records to retrieve relative to the offset
	Offset        *int64  `json:"offset,omitempty" mapstructure:"offset,omitempty"`                 // Offset from the start of the entity list
	SortAttribute *string `json:"sort_attribute,omitempty" mapstructure:"sort_attribute,omitempty"` // The attribute to perform sort on
	SortOrder     *string `json:"sort_order,omitempty" mapstructure:"sort_order,omitempty"`         // The sort order in which results are returned
}

// ListMetadataOutput All api calls that return a list will have this metadata block
type ListMetadataOutput struct {
	Filter        *string `json:"filter,omitempty" mapstructure:"filter,omitempty"`                 // The filter used for the results
	Kind          *string `json:"kind,omitempty" mapstructure:"kind,omitempty"`                     // The kind name
	Length        *int64  `json:"length,omitempty" mapstructure:"length,omitempty"`                 // The number of records retrieved relative to the offset
	Offset        *int64  `json:"offset,omitempty" mapstructure:"offset,omitempty"`                 // Offset from the start of the entity list
	SortAttribute *string `json:"sort_attribute,omitempty" mapstructure:"sort_attribute,omitempty"` // The attribute to perform sort on
	SortOrder     *string `json:"sort_order,omitempty" mapstructure:"sort_order,omitempty"`         // The sort order in which results are returned
	TotalMatches  *int64  `json:"total_matches,omitempty" mapstructure:"total_matches,omitempty"`   // Total matches found
}

// NetworkSecurityRuleIntentResource ... Response object for intentful operations on a network_security_rule
type NetworkSecurityRuleIntentResource struct {
	APIVersion *string                       `json:"api_version,omitempty" mapstructure:"api_version,omitempty"`
	Metadata   *Metadata                     `json:"metadata,omitempty" mapstructure:"metadata,omitempty"`
	Spec       *NetworkSecurityRule          `json:"spec,omitempty" mapstructure:"spec,omitempty"`
	Status     *NetworkSecurityRuleDefStatus `json:"status,omitempty" mapstructure:"status,omitempty"`
}

// NetworkSecurityRuleListIntentResponse Response object for intentful operation of network_security_rules
type NetworkSecurityRuleListIntentResponse struct {
	APIVersion string                               `json:"api_version" mapstructure:"api_version"`
	Entities   []*NetworkSecurityRuleIntentResource `json:"entities,omitempty" bson:"entities,omitempty" mapstructure:"entities,omitempty"`
	Metadata   *ListMetadataOutput                  `json:"metadata" mapstructure:"metadata"`
}

// VolumeGroupInput Represents the request body for create volume_grop request
type VolumeGroupInput struct {
	APIVersion *string      `json:"api_version,omitempty" mapstructure:"api_version,omitempty"` // default 3.1.0
	Metadata   *Metadata    `json:"metadata,omitempty" mapstructure:"metadata,omitempty"`       // The volume_group kind metadata.
	Spec       *VolumeGroup `json:"spec,omitempty" mapstructure:"spec,omitempty"`               // Volume group input spec.
}

// VolumeGroup Represents volume group input spec.
type VolumeGroup struct {
	Name        *string               `json:"name" mapstructure:"name"`                                   // Volume Group name (required)
	Description *string               `json:"description,omitempty" mapstructure:"description,omitempty"` // Volume Group description.
	Resources   *VolumeGroupResources `json:"resources" mapstructure:"resources"`                         // Volume Group resources.
}

// VolumeGroupResources Represents the volume group resources
type VolumeGroupResources struct {
	FlashMode         *string         `json:"flash_mode,omitempty" mapstructure:"flash_mode,omitempty"`                   // Flash Mode, if enabled all disks of the VG are pinned to SSD
	FileSystemType    *string         `json:"file_system_type,omitempty" mapstructure:"file_system_type,omitempty"`       // File system to be used for volume
	SharingStatus     *string         `json:"sharing_status,omitempty" mapstructure:"sharing_status,omitempty"`           // Whether the VG can be shared across multiple iSCSI initiators
	AttachmentList    []*VMAttachment `json:"attachment_list,omitempty" mapstructure:"attachment_list,omitempty"`         // VMs attached to volume group.
	DiskList          []*VGDisk       `json:"disk_list,omitempty" mapstructure:"disk_list,omitempty"`                     // VGDisk Volume group disk specification.
	IscsiTargetPrefix *string         `json:"iscsi_target_prefix,omitempty" mapstructure:"iscsi_target_prefix,omitempty"` // iSCSI target prefix-name.
}

// VMAttachment VMs attached to volume group.
type VMAttachment struct {
	VMReference        *Reference `json:"vm_reference" mapstructure:"vm_reference"`                 // Reference to a kind
	IscsiInitiatorName *string    `json:"iscsi_initiator_name" mapstructure:"iscsi_initiator_name"` // Name of the iSCSI initiator of the workload outside Nutanix cluster.
}

// VGDisk Volume group disk specification.
type VGDisk struct {
	VmdiskUUID           *string    `json:"vmdisk_uuid" mapstructure:"vmdisk_uuid"`                       // The UUID of this volume disk
	Index                *int64     `json:"index" mapstructure:"index"`                                   // Index of the volume disk in the group.
	DataSourceReference  *Reference `json:"data_source_reference" mapstructure:"data_source_reference"`   // Reference to a kind
	DiskSizeMib          *int64     `json:"disk_size_mib" mapstructure:"disk_size_mib"`                   // Size of the disk in MiB.
	StorageContainerUUID *string    `json:"storage_container_uuid" mapstructure:"storage_container_uuid"` // Container UUID on which to create the disk.
}

// VolumeGroupResponse Response object for intentful operations on a volume_group
type VolumeGroupResponse struct {
	APIVersion *string               `json:"api_version" mapstructure:"api_version"`           //
	Metadata   *Metadata             `json:"metadata" mapstructure:"metadata"`                 // The volume_group kind metadata
	Spec       *VolumeGroup          `json:"spec,omitempty" mapstructure:"spec,omitempty"`     // Volume group input spec.
	Status     *VolumeGroupDefStatus `json:"status,omitempty" mapstructure:"status,omitempty"` // Volume group configuration.
}

// VolumeGroupDefStatus  Volume group configuration.
type VolumeGroupDefStatus struct {
	State       *string               `json:"state" mapstructure:"state"`               // The state of the volume group entity.
	MessageList []*MessageResource    `json:"message_list" mapstructure:"message_list"` // Volume group message list.
	Name        *string               `json:"name" mapstructure:"name"`                 // Volume group name.
	Resources   *VolumeGroupResources `json:"resources" mapstructure:"resources"`       // Volume group resources.
	Description *string               `json:"description" mapstructure:"description"`   // Volume group description.
}

// VolumeGroupListResponse Response object for intentful operation of volume_groups
type VolumeGroupListResponse struct {
	APIVersion *string                `json:"api_version" mapstructure:"api_version"`
	Entities   []*VolumeGroupResponse `json:"entities,omitempty" mapstructure:"entities,omitempty"`
	Metadata   *ListMetadataOutput    `json:"metadata" mapstructure:"metadata"`
}

// TasksResponse ...
type TasksResponse struct {
	Status               *string      `json:"status,omitempty" mapstructure:"status,omitempty"`
	LastUpdateTime       *time.Time   `json:"last_update_time,omitempty" mapstructure:"last_update_time,omitempty"`
	LogicalTimestamp     *int64       `json:"logical_timestamp,omitempty" mapstructure:"logical_timestamp,omitempty"`
	EntityReferenceList  []*Reference `json:"entity_reference_list,omitempty" mapstructure:"entity_reference_list,omitempty"`
	StartTime            *time.Time   `json:"start_time,omitempty" mapstructure:"start_time,omitempty"`
	CreationTime         *time.Time   `json:"creation_time,omitempty" mapstructure:"creation_time,omitempty"`
	ClusterReference     *Reference   `json:"cluster_reference,omitempty" mapstructure:"cluster_reference,omitempty"`
	SubtaskReferenceList []*Reference `json:"subtask_reference_list,omitempty" mapstructure:"subtask_reference_list,omitempty"`
	CompletionTime       *time.Time   `json:"completion_timev" mapstructure:"completion_timev"`
	ProgressMessage      *string      `json:"progress_message,omitempty" mapstructure:"progress_message,omitempty"`
	OperationType        *string      `json:"operation_type,omitempty" mapstructure:"operation_type,omitempty"`
	PercentageComplete   *int64       `json:"percentage_complete,omitempty" mapstructure:"percentage_complete,omitempty"`
	APIVersion           *string      `json:"api_version,omitempty" mapstructure:"api_version,omitempty"`
	UUID                 *string      `json:"uuid,omitempty" mapstructure:"uuid,omitempty"`
	ErrorDetail          *string      `json:"error_detail,omitempty" mapstructure:"error_detail,omitempty"`
}

// DeleteResponse ...
type DeleteResponse struct {
	Status     *DeleteStatus `json:"status" mapstructure:"status"`
	Spec       string        `json:"spec" mapstructure:"spec"`
	APIVersion string        `json:"api_version" mapstructure:"api_version"`
	Metadata   *Metadata     `json:"metadata" mapstructure:"metadata"`
}

// DeleteStatus ...
type DeleteStatus struct {
	State            string            `json:"state" mapstructure:"state"`
	ExecutionContext *ExecutionContext `json:"execution_context" mapstructure:"execution_context"`
}

/* Host Resource */

// DomainCredencial represents the way to login server
type DomainCredencial struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// WindowsDomain means Hyper-V node domain
type WindowsDomain struct {
	Name                 string            `json:"name,omitempty"`
	NameServerIP         string            `json:"name_server_ip,omitempty"`
	OrganizationUnitPath string            `json:"organization_unit_path,omitempty"`
	NamePrefix           string            `json:"name_prefix,omitempty"`
	DomainName           string            `json:"domain_name,omitempty"`
	DomainCredencial     *DomainCredencial `json:"domain_credencial,omitempty"`
}

// OplogUsage represents oplog disk usage
type OplogUsage struct {
	OplogDiskPct  *float64 `json:"oplog_disk_pct,omitempty"`
	OplogDiskSize *int64   `json:"oplog_disk_size,omitempty"`
}

// ControllerVM means Hyper-V node domain
type ControllerVM struct {
	IP         string      `json:"ip,omitempty"`
	NatIP      string      `json:"nat_ip,omitempty"`
	NatPort    *int64      `json:"nat_port,omitempty"`
	OplogUsage *OplogUsage `json:"oplog_usage,omitempty"`
}

// FailoverCluster means Hiper-V failover cluster
type FailoverCluster struct {
	IP               string            `json:"ip,omitempty"`
	Name             string            `json:"name,omitempty"`
	DomainCredencial *DomainCredencial `json:"domain_credencial,omitempty"`
}

// IPMI means Host IPMI Information
type IPMI struct {
	IP string `json:"ip,omitempty"`
}

// ReferenceValues references to a kind
type ReferenceValues struct {
	Kind string `json:"kind,omitempty"`
	UUID string `json:"uuid,omitempty"`
	Name string `json:"name,omitempty"`
}

// GPU represnts list of GPUs on the host
type GPU struct {
	Status                 string           `json:"status,omitempty"`
	Vendor                 string           `json:"vendor,omitempty"`
	NumVirtualDisplayHeads *int64           `json:"num_virtual_display_heads,omitempty"`
	Assignable             bool             `json:"assignable,omitempty"`
	LicenseList            []*string        `json:"license_list,omitempty"`
	NumVgpusAllocated      *int64           `json:"num_vgpus_allocated,omitempty"`
	PciAddress             string           `json:"pci_address,omitempty"`
	Name                   string           `json:"name,omitempty"`
	FrameBufferSizeMib     *int64           `json:"frame_buffer_size_mib,omitempty"`
	Index                  *int64           `json:"index,omitempty"`
	UUID                   string           `json:"uuid,omitempty"`
	NumaNode               *int64           `json:"numa_node,omitempty"`
	MaxResoution           string           `json:"max_resolution,omitempty"`
	ConsumerReference      *ReferenceValues `json:"consumer_reference,omitempty"`
	Mode                   string           `json:"mode,omitempty"`
	Fraction               *int64           `json:"fraction,omitempty"`
	GuestDriverVersion     string           `json:"guest_driver_version,omitempty"`
	DeviceID               *int64           `json:"device_id,omitempty"`
}

// Hypervisor Full name of hypervisor running on Host
type Hypervisor struct {
	NumVms             *int64 `json:"num_vms,omitempty"`
	IP                 string `json:"ip,omitempty"`
	HypervisorFullName string `json:"hypervisor_full_name,omitempty"`
}

// Block represents Host block config info.
type Block struct {
	BlockSerialNumber string `json:"block_serial_number,omitempty"`
	BlockModel        string `json:"block_model,omitempty"`
}

// HostResources represents the host resources
type HostResources struct {
	GPUDriverVersion       string             `json:"gpu_driver_version,omitempty"`
	FailoverCluster        *FailoverCluster   `json:"failover_cluster,omitempty"`
	IPMI                   *IPMI              `json:"ipmi,omitempty"`
	CPUModel               string             `json:"cpu_model,omitempty"`
	HostNicsIDList         []*string          `json:"host_nics_id_list,omitempty"`
	NumCPUSockets          *int64             `json:"num_cpu_sockets,omitempty"`
	WindowsDomain          *WindowsDomain     `json:"windows_domain,omitempty"`
	GPUList                []*GPU             `json:"gpu_list,omitempty"`
	SerialNumber           string             `json:"serial_number,omitempty"`
	CPUCapacityHZ          *int64             `json:"cpu_capacity_hz,omitempty"`
	MemoryVapacityMib      *int64             `json:"memory_capacity_mib,omitempty"`
	HostDisksReferenceList []*ReferenceValues `json:"host_disks_reference_list,omitempty"`
	MonitoringState        string             `json:"monitoring_state,omitempty"`
	Hypervisor             *Hypervisor        `json:"hypervisor,omitempty"`
	HostType               string             `json:"host_type,omitempty"`
	NumCPUCores            *int64             `json:"num_cpu_cores,omitempty"`
	RackableUnitReference  *ReferenceValues   `json:"rackable_unit_reference,omitempty"`
	ControllerVM           *ControllerVM      `json:"controller_vm,omitempty"`
	Block                  *Block             `json:"block,omitempty"`
}

// HostSpec Represents volume group input spec.
type HostSpec struct {
	Name      string         `json:"name,omitempty"`
	Resources *HostResources `json:"resources,omitempty"`
}

// HostStatus  Volume group configuration.
type HostStatus struct {
	State            string             `json:"state,omitempty"`
	MessageList      []*MessageResource `json:"message_list,omitempty"`
	Name             string             `json:"name,omitempty"`
	Resources        *HostResources     `json:"resources,omitempty"`
	ClusterReference *ReferenceValues   `json:"cluster_reference,omitempty"`
}

// HostResponse Response object for intentful operations on a Host
type HostResponse struct {
	APIVersion string      `json:"api_version,omitempty"`
	Metadata   *Metadata   `json:"metadata,omitempty"`
	Spec       *HostSpec   `json:"spec,omitempty"`
	Status     *HostStatus `json:"status,omitempty"`
}

// HostListResponse Response object for intentful operation of Host
type HostListResponse struct {
	APIVersion string              `json:"api_version,omitempty"`
	Entities   []*HostResponse     `json:"entities,omitempty"`
	Metadata   *ListMetadataOutput `json:"metadata,omitempty"`
}

/* Project Resource */

// Resources represents the utilization limits for resource types
type Resources struct {
	Units        string `json:"units,omitempty"`
	Limit        *int64 `json:"limit,omitempty"`
	ResourceType string `json:"resource_type,omitempty"`
	Value        *int64 `json:"value,omitempty"`
}

// ResourceDomain specification (limits)
type ResourceDomain struct {
	Resources []*Resources `json:"resources,omitempty"`
}

// ProjectResources ...
type ProjectResources struct {
	ResourceDomain                 *ResourceDomain    `json:"resource_domain,omitempty"`
	AccountReferenceList           []*ReferenceValues `json:"account_reference_list,omitempty"`
	EnvironmentReferenceList       []*ReferenceValues `json:"environment_reference_list,omitempty"`
	DefaultSubnetReference         *ReferenceValues   `json:"default_subnet_reference,omitempty"`
	UserReferenceList              []*ReferenceValues `json:"user_reference_list,omitempty"`
	IsDefault                      bool               `json:"is_default,omitempty"`
	ExternalUserGroupReferenceList []*ReferenceValues `json:"external_user_group_reference_list,omitempty"`
	SubnetReferenceList            []*ReferenceValues `json:"subnet_reference_list,omitempty"`
	ExternalNetworkList            []*ReferenceValues `json:"external_network_list,omitempty"`
}

// ProjectStatus ...
type ProjectStatus struct {
	State            string             `json:"state,omitempty"`
	MessageList      []*MessageResource `json:"message_list,omitempty"`
	Name             string             `json:"name,omitempty"`
	Resources        *ProjectResources  `json:"resources,omitempty"`
	Descripion       string             `json:"description,omitempty"`
	ExecutionContext *ExecutionContext  `json:"execution_context,omitempty"`
}

// ProjectSpec ...
type ProjectSpec struct {
	Name       string            `json:"name,omitempty"`
	Resources  *ProjectResources `json:"resources,omitempty"`
	Descripion string            `json:"description,omitempty"`
}

// Project Response object for intentful operations on a Host
type Project struct {
	Status     *ProjectStatus `json:"status,omitempty"`
	Spec       *ProjectSpec   `json:"spec,omitempty"`
	APIVersion string         `json:"api_version,omitempty"`
	Metadata   *Metadata      `json:"metadata,omitempty"`
}

// ProjectListResponse Response object for intentful operation of Host
type ProjectListResponse struct {
	APIVersion string              `json:"api_version,omitempty"`
	Entities   []*Project          `json:"entities,omitempty"`
	Metadata   *ListMetadataOutput `json:"metadata,omitempty"`
}

// AccessControlPolicyResources ...
type AccessControlPolicyResources struct {
	UserReferenceList      []*Reference `json:"user_reference_list,omitempty"`
	UserGroupReferenceList []*Reference `json:"user_group_reference_list,omitempty"`
	RoleReference          *Reference   `json:"role_reference,omitempty"`
	FilterList             *FilterList  `json:"filter_list,omitempty"`
}

// FilterList ...
type FilterList struct {
	ContextList []*ContextList `json:"context_list,omitempty"`
}

// ContextList ...
type ContextList struct {
	ScopeFilterExpressionList  []*ScopeFilterExpressionList `json:"scope_filter_expression_list,omitempty"`
	EntityFilterExpressionList []EntityFilterExpressionList `json:"entity_filter_expression_list,omitempty"`
}

// ScopeFilterExpressionList ...
type ScopeFilterExpressionList struct {
	LeftHandSide  string        `json:"left_hand_side,omitempty"`
	Operator      string        `json:"operator,omitempty"`
	RightHandSide RightHandSide `json:"right_hand_side,omitempty"`
}

// EntityFilterExpressionList ...
type EntityFilterExpressionList struct {
	LeftHandSide  LeftHandSide  `json:"left_hand_side,omitempty"`
	Operator      string        `json:"operator,omitempty"`
	RightHandSide RightHandSide `json:"right_hand_side,omitempty"`
}

// LeftHandSide ...
type LeftHandSide struct {
	EntityType *string `json:"entity_type,omitempty"`
}

// RightHandSide ...
type RightHandSide struct {
	Collection *string             `json:"collection,omitempty"`
	Categories map[string][]string `json:"categories,omitempty"`
	UUIDList   []string            `json:"uuid_list,omitempty"`
}

// AccessControlPolicyStatus ...
type AccessControlPolicyStatus struct {
	State            *string                       `json:"state,omitempty"`
	MessageList      []*MessageResource            `json:"message_list,omitempty"`
	Name             *string                       `json:"name,omitempty"`
	Resources        *AccessControlPolicyResources `json:"resources,omitempty"`
	Description      *string                       `json:"description,omitempty"`
	ExecutionContext *ExecutionContext             `json:"execution_context,omitempty"`
}

// AccessControlPolicySpec ...
type AccessControlPolicySpec struct {
	Name        *string                       `json:"name,omitempty"`
	Resources   *AccessControlPolicyResources `json:"resources,omitempty"`
	Description *string                       `json:"description,omitempty"`
}

// AccessControlPolicy Response object for intentful operations on a access policy
type AccessControlPolicy struct {
	Status     *AccessControlPolicyStatus `json:"status,omitempty"`
	Spec       *AccessControlPolicySpec   `json:"spec,omitempty"`
	APIVersion string                     `json:"api_version,omitempty"`
	Metadata   *Metadata                  `json:"metadata,omitempty"`
}

// AccessControlPolicyListResponse Response object for intentful operation of access policy
type AccessControlPolicyListResponse struct {
	APIVersion string                 `json:"api_version,omitempty"`
	Entities   []*AccessControlPolicy `json:"entities,omitempty"`
	Metadata   *ListMetadataOutput    `json:"metadata,omitempty"`
}

// RoleResources ...
type RoleResources struct {
	PermissionReferenceList []*Reference `json:"permission_reference_list,omitempty"`
}

// RoleStatus ...
type RoleStatus struct {
	State            *string            `json:"state,omitempty"`
	MessageList      []*MessageResource `json:"message_list,omitempty"`
	Name             *string            `json:"name,omitempty"`
	Resources        *RoleResources     `json:"resources,omitempty"`
	Description      *string            `json:"description,omitempty"`
	ExecutionContext *ExecutionContext  `json:"execution_context,omitempty"`
}

// RoleSpec ...
type RoleSpec struct {
	Name        *string        `json:"name,omitempty"`
	Resources   *RoleResources `json:"resources,omitempty"`
	Description *string        `json:"description,omitempty"`
}

// Role Response object for intentful operations on a access policy
type Role struct {
	Status     *RoleStatus `json:"status,omitempty"`
	Spec       *RoleSpec   `json:"spec,omitempty"`
	APIVersion string      `json:"api_version,omitempty"`
	Metadata   *Metadata   `json:"metadata,omitempty"`
}

// RoleListResponse Response object for intentful operation of access policy
type RoleListResponse struct {
	APIVersion string              `json:"api_version,omitempty"`
	Entities   []*Role             `json:"entities,omitempty"`
	Metadata   *ListMetadataOutput `json:"metadata,omitempty"`
}

type ResourceUsageSummary struct {
	ResourceDomain *ResourceDomainStatus `json:"resource_domain"` // The status for a resource domain (limits and values)
}

type ResourceDomainStatus struct {
	Resources []ResourceUtilizationStatus `json:"resources,omitempty"` // The utilization/limit for resource types
}

type ResourceUtilizationStatus struct {
	Limit        *int64  `json:"limit,omitempty"`         // The resource consumption limit (unspecified is unlimited)
	ResourceType *string `json:"resource_type,omitempty"` // The type of resource (for example storage, CPUs)
	Units        *string `json:"units,omitempty"`         // The units of the resource type
	Value        *int64  `json:"value,omitempty"`         // The amount of resource consumed
}

// An intentful representation of a user
type UserIntentInput struct {
	APIVersion *string   `json:"api_version,omitempty"` // API Version of the Nutanix v3 API framework.
	Metadata   *Metadata `json:"metadata,omitempty"`    // The user kind metadata
	Spec       *UserSpec `json:"spec,omitempty"`        // User Input Definition.
}

// Response object for intentful operations on a user
type UserIntentResponse struct {
	APIVersion *string     `json:"api_version,omitempty"` // API Version of the Nutanix v3 API framework.
	Metadata   *Metadata   `json:"metadata,omitempty"`    // The user kind metadata
	Spec       *UserSpec   `json:"spec,omitempty"`        // User Input Definition.
	Status     *UserStatus `json:"status,omitempty"`      // User status definition.
}

// User Input Definition.
type UserSpec struct {
	Resources *UserResources `json:"resources,omitempty"` // User Resource Definition.
}

// User Resource Definition.
type UserResources struct {
	DirectoryServiceUser *DirectoryServiceUser `json:"directory_service_user,omitempty"` // A Directory Service user.
	IdentityProviderUser *IdentityProvider     `json:"identity_provider_user,omitempty"` // An Identity Provider user.
}

// A Directory Service user.
type DirectoryServiceUser struct {
	DefaultUserPrincipalName  *string    `json:"default_user_principal_name,omitempty"` // The Default UserPrincipalName of the user from the directory service.
	DirectoryServiceReference *Reference `json:"directory_service_reference,omitempty"` // The reference to a directory_service
	UserPrincipalName         *string    `json:"user_principal_name,omitempty"`         // The UserPrincipalName of the user from the directory service.
}

// An Identity Provider user.
type IdentityProvider struct {
	IdentityProviderReference *Reference `json:"identity_provider_reference,omitempty"` // The reference to a identity_provider
	Username                  *string    `json:"username,omitempty"`                    // The username from the identity provider. Name Id for SAML Identity Provider.
}

// User status definition.
type UserStatus struct {
	MessageList      []MessageResource    `json:"message_list,omitempty"`
	Name             *string              `json:"name,omitempty"`      // Name of the User.
	Resources        *UserStatusResources `json:"resources,omitempty"` // User Resource Definition.
	State            *string              `json:"state,omitempty"`     // The state of the entity.
	ExecutionContext *ExecutionContext    `json:"execution_context,omitempty"`
}

// User Resource Definition.
type UserStatusResources struct {
	AccessControlPolicyReferenceList []*Reference          `json:"access_control_policy_reference_list,omitempty"` // List of ACP references.
	DirectoryServiceUser             *DirectoryServiceUser `json:"directory_service_user,omitempty"`               // A Directory Service user.
	DisplayName                      *string               `json:"display_name,omitempty"`                         // The display name of the user (common name) provided by the directory service.
	IdentityProviderUser             *IdentityProvider     `json:"identity_provider_user,omitempty"`               // An Identity Provider user.
	ProjectsReferenceList            []*Reference          `json:"projects_reference_list,omitempty"`              // A list of projects the user is part of.
	ResourceUsageSummary             *ResourceUsageSummary `json:"resource_usage_summary,omitempty"`
	UserType                         *string               `json:"user_type,omitempty"`
}

type UserListResponse struct {
	APIVersion *string               `json:"api_version,omitempty"` // API Version of the Nutanix v3 API framework.
	Entities   []*UserIntentResponse `json:"entities,omitempty"`
	Metadata   *ListMetadataOutput   `json:"metadata,omitempty"` // All api calls that return a list will have this metadata block
}

// Response object for intentful operations on a user_group
type UserGroupIntentResponse struct {
	APIVersion *string          `json:"api_version,omitempty"` // API Version of the Nutanix v3 API framework.
	Metadata   *Metadata        `json:"metadata,omitempty"`    // The user_group kind metadata
	Spec       *UserGroupSpec   `json:"spec,omitempty"`        // User Group Input Definition.
	Status     *UserGroupStatus `json:"status,omitempty"`      // User group status definition.
}

// User Group Input Definition.
type UserGroupSpec struct {
	Resources *UserGroupResources `json:"resources,omitempty"` // User Group Resource Definition
}

// User Group Resource Definition
type UserGroupResources struct {
	AccessControlPolicyReferenceList []*Reference               `json:"access_control_policy_reference_list,omitempty"` // List of ACP references.
	DirectoryServiceUserGroup        *DirectoryServiceUserGroup `json:"directory_service_user_group,omitempty"`         // A Directory Service user group.
	DisplayName                      *string                    `json:"display_name,omitempty"`                         // The display name for the user group.
	ProjectsReferenceList            []*Reference               `json:"projects_reference_list,omitempty"`              // A list of projects the user group is part of.
	UserGroupType                    *string                    `json:"user_group_type,omitempty"`
}

// User group status definition.
type UserGroupStatus struct {
	MessageList []MessageResource   `json:"message_list,omitempty"`
	Resources   *UserGroupResources `json:"resources,omitempty"` // User Group Resource Definition.
	State       *string             `json:"state,omitempty"`     // The state of the entity.
}

// A Directory Service user group.
type DirectoryServiceUserGroup struct {
	DirectoryServiceReference *Reference `json:"directory_service_reference,omitempty"` // The reference to a directory_service
	DistinguishedName         *string    `json:"distinguished_name,omitempty"`          // The Distinguished name for the user group.
}

// Response object for intentful operation of user_groups
type UserGroupListResponse struct {
	APIVersion *string                    `json:"api_version,omitempty"` // API Version of the Nutanix v3 API framework.
	Entities   []*UserGroupIntentResponse `json:"entities,omitempty"`
	Metadata   *ListMetadataOutput        `json:"metadata,omitempty"` // All api calls that return a list will have this metadata block
}

// Response object for intentful operations on a user_group
type PermissionIntentResponse struct {
	APIVersion *string           `json:"api_version,omitempty"` // API Version of the Nutanix v3 API framework.
	Metadata   *Metadata         `json:"metadata,omitempty"`    // The user_group kind metadata
	Spec       *PermissionSpec   `json:"spec,omitempty"`        // Permission Input Definition.
	Status     *PermissionStatus `json:"status,omitempty"`      // User group status definition.
}

// Permission Input Definition.
type PermissionSpec struct {
	Name        *string              `json:"name,omitempty"`        // The name for the permission.
	Description *string              `json:"description,omitempty"` // The display name for the permission.
	Resources   *PermissionResources `json:"resources,omitempty"`   // Permission Resource Definition
}

// Permission Resource Definition
type PermissionResources struct {
	Operation *string           `json:"operation,omitempty"`
	Kind      *string           `json:"kind,omitempty"`
	Fields    *FieldsPermission `json:"fields,omitempty"`
}

type FieldsPermission struct {
	FieldMode     *string   `json:"field_mode,omitempty"`
	FieldNameList []*string `json:"field_name_list,omitempty"`
}

// Permission status definition.
type PermissionStatus struct {
	Name        *string              `json:"name,omitempty"`        // The name for the permission.
	Description *string              `json:"description,omitempty"` // The display name for the permission.
	Resources   *PermissionResources `json:"resources,omitempty"`   // Permission Resource Definition
	MessageList []MessageResource    `json:"message_list,omitempty"`
	State       *string              `json:"state,omitempty"` // The state of the entity.
}

// Response object for intentful operation of Permissions
type PermissionListResponse struct {
	APIVersion *string                     `json:"api_version,omitempty"` // API Version of the Nutanix v3 API framework.
	Entities   []*PermissionIntentResponse `json:"entities,omitempty"`
	Metadata   *ListMetadataOutput         `json:"metadata,omitempty"` // All api calls that return a list will have this metadata block
}

// ProtectionRuleResources represents the resources of protection rules
type ProtectionRuleResources struct {
	StartTime                        string                              `json:"start_time,omitempty"`
	AvailabilityZoneConnectivityList []*AvailabilityZoneConnectivityList `json:"availability_zone_connectivity_list,omitempty"`
	OrderedAvailabilityZoneList      []*OrderedAvailabilityZoneList      `json:"ordered_availability_zone_list,omitempty"`
	CategoryFilter                   *CategoryFilter                     `json:"category_filter,omitempty"`
}

// AvailabilityZoneConnectivityList represents a object for resource of protection rule
type AvailabilityZoneConnectivityList struct {
	DestinationAvailabilityZoneIndex *int64                  `json:"destination_availability_zone_index,omitempty"`
	SourceAvailabilityZoneIndex      *int64                  `json:"source_availability_zone_index,omitempty"`
	SnapshotScheduleList             []*SnapshotScheduleList `json:"snapshot_schedule_list,omitempty"`
}

// SnapshotScheduleList represents a object for resource of protection rule
type SnapshotScheduleList struct {
	RecoveryPointObjectiveSecs    *int64                   `json:"recovery_point_objective_secs,omitempty"`
	LocalSnapshotRetentionPolicy  *SnapshotRetentionPolicy `json:"local_snapshot_retention_policy,omitempty"`
	AutoSuspendTimeoutSecs        *int64                   `json:"auto_suspend_timeout_secs,omitempty"`
	SnapshotType                  string                   `json:"snapshot_type,omitempty"`
	RemoteSnapshotRetentionPolicy *SnapshotRetentionPolicy `json:"remote_snapshot_retention_policy,omitempty"`
}

// SnapshotRetentionPolicy represents a object for resource of protection rule
type SnapshotRetentionPolicy struct {
	NumSnapshots          *int64                 `json:"num_snapshots,omitempty"`
	RollupRetentionPolicy *RollupRetentionPolicy `json:"rollup_retention_policy,omitempty"`
}

// RollupRetentionPolicy represents a object for resource of protection rule
type RollupRetentionPolicy struct {
	Multiple             *int64 `json:"multiple,omitempty"`
	SnapshotIntervalType string `json:"snapshot_interval_type,omitempty"`
}

// OrderedAvailabilityZoneList represents a object for resource of protection rule
type OrderedAvailabilityZoneList struct {
	ClusterUUID         string `json:"cluster_uuid,omitempty"`
	AvailabilityZoneURL string `json:"availability_zone_url,omitempty"`
}

// ProtectionRuleStatus represents a status of a protection rule
type ProtectionRuleStatus struct {
	State            string                   `json:"state,omitempty"`
	MessageList      []*MessageResource       `json:"message_list,omitempty"`
	Name             string                   `json:"name,omitempty"`
	Resources        *ProtectionRuleResources `json:"resources,omitempty"`
	ExecutionContext *ExecutionContext        `json:"execution_context,omitempty"`
}

// ProtectionRuleSpec represents a spec of protection rules
type ProtectionRuleSpec struct {
	Name        string                   `json:"name,omitempty"`
	Description string                   `json:"description,omitempty"`
	Resources   *ProtectionRuleResources `json:"resources,omitempty"`
}

// ProtectionRuleResponse represents a response object of a protection rule
type ProtectionRuleResponse struct {
	APIVersion string                `json:"api_version,omitempty"`
	Metadata   *Metadata             `json:"metadata,omitempty"`
	Spec       *ProtectionRuleSpec   `json:"spec,omitempty"`
	Status     *ProtectionRuleStatus `json:"status,omitempty"`
}

// ProtectionRulesListResponse represents the response of a list of protection rules
type ProtectionRulesListResponse struct {
	APIVersion string                    `json:"api_version,omitempty"`
	Entities   []*ProtectionRuleResponse `json:"entities,omitempty"`
	Metadata   *ListMetadataOutput       `json:"metadata,omitempty"`
}

// ProtectionRuleInput Represents the request of create protection rule
type ProtectionRuleInput struct {
	APIVersion string              `json:"api_version,omitempty"`
	Metadata   *Metadata           `json:"metadata,omitempty"`
	Spec       *ProtectionRuleSpec `json:"spec,omitempty"`
}

// RecoveryPlanResources represents the resources of recovery plan
type RecoveryPlanResources struct {
	VolumeGroupRecoveryInfoList []*VolumeGroupRecoveryInfoList `json:"volume_group_recovery_info_list,omitempty"`
	StageList                   []*StageList                   `json:"stage_list,omitempty"`
	Parameters                  *Parameters                    `json:"parameters,omitempty"`
}

// VolumeGroupRecoveryInfoList encapsulates the specification of volume
// group instances to be included in the recovery plan.
type VolumeGroupRecoveryInfoList struct {
	CategoryFilter *CategoryFilter `json:"category_filter,omitempty"`
}

// Parameters represents a object for resource of recovery plan
type Parameters struct {
	FloatingIPAssignmentList []*FloatingIPAssignmentList `json:"floating_ip_assignment_list,omitempty"`
	NetworkMappingList       []*NetworkMappingList       `json:"network_mapping_list,omitempty"`
	AvailabilityZoneList     []*AvailabilityZoneList     `json:"availability_zone_list,omitempty"`
	PrimaryLocationIndex     *int64                      `json:"primary_location_index,omitempty"`
}

// AvailabilityZoneList represents objects that encapsulate the list of AOS
// clusters in an AZ.
type AvailabilityZoneList struct {
	ClusterReferenceList []*Reference `json:"cluster_reference_list,omitempty" mapstructure:"cluster_reference_list,omitempty"`
	AvailabilityZoneURL  *string      `json:"availability_zone_url"`
}

// FloatingIPAssignmentList represents a object for resource of recovery plan
type FloatingIPAssignmentList struct {
	AvailabilityZoneURL string                `json:"availability_zone_url,omitempty"`
	VMIPAssignmentList  []*VMIPAssignmentList `json:"vm_ip_assignment_list,omitempty"`
}

// VMIPAssignmentList represents a object for resource of recovery plan
type VMIPAssignmentList struct {
	TestFloatingIPConfig     *FloatingIPConfig `json:"test_floating_ip_config,omitempty"`
	RecoveryFloatingIPConfig *FloatingIPConfig `json:"recovery_floating_ip_config,omitempty"`
	VMReference              *Reference        `json:"vm_reference,omitempty"`
	VMNICInformation         *VMNICInformation `json:"vm_nic_information,omitempty"`
}

// FloatingIPConfig represents a object for resource of recovery plan
type FloatingIPConfig struct {
	IP                        string `json:"ip,omitempty"`
	ShouldAllocateDynamically *bool  `json:"should_allocate_dynamically,omitempty"`
}

// VMNICInformation represents a object for resource of recovery plan
type VMNICInformation struct {
	IP   string `json:"ip,omitempty"`
	UUID string `json:"uuid,omitempty"`
}

// represents a object for resource of recovery plan
type NetworkMappingList struct {
	AvailabilityZoneNetworkMappingList []*AvailabilityZoneNetworkMappingList `json:"availability_zone_network_mapping_list,omitempty"`
	AreNetworksStretched               *bool                                 `json:"are_networks_stretched,omitempty"`
}

// AvailabilityZoneNetworkMappingList represents a object for resource of recovery plan
type AvailabilityZoneNetworkMappingList struct {
	RecoveryNetwork          *Network            `json:"recovery_network,omitempty"`
	AvailabilityZoneURL      string              `json:"availability_zone_url,omitempty"`
	TestNetwork              *Network            `json:"test_network,omitempty"`
	RecoveryIPAssignmentList []*IPAssignmentList `json:"recovery_ip_assignment_list,omitempty"`
	TestIPAssignmentList     []*IPAssignmentList `json:"test_ip_assignment_list,omitempty"`
	ClusterReferenceList     []*Reference        `json:"cluster_reference_list,omitempty"`
}

type IPAssignmentList struct {
	VMReference  *Reference      `json:"vm_reference,omitempty"`
	IPConfigList []*IPConfigList `json:"ip_config_list,omitempty"`
}

type IPConfigList struct {
	IPAddress string `json:"ip_address,omitempty"`
}

// Network represents a object for resource of recovery plan
type Network struct {
	VirtualNetworkReference *Reference    `json:"virtual_network_reference,omitempty"`
	SubnetList              []*SubnetList `json:"subnet_list,omitempty"`
	Name                    string        `json:"name,omitempty"`
	VPCReference            *Reference    `json:"vpc_reference,omitempty"`
	UseVPCReference         *bool         `json:"use_vpc_reference,omitempty"`
}

// SubnetList represents a object for resource of recovery plan
type SubnetList struct {
	GatewayIP                 string `json:"gateway_ip,omitempty"`
	ExternalConnectivityState string `json:"external_connectivity_state,omitempty"`
	PrefixLength              *int64 `json:"prefix_length,omitempty"`
}

// StageList represents a object for resource of recovery plan
type StageList struct {
	StageWork     *StageWork `json:"stage_work,omitempty"`
	StageUUID     string     `json:"stage_uuid,omitempty"`
	DelayTimeSecs *int64     `json:"delay_time_secs,omitempty"`
}

// StageWork represents a object for resource of recovery plan
type StageWork struct {
	RecoverEntities *RecoverEntities `json:"recover_entities,omitempty"`
}

// RecoverEntities represents a object for resource of recovery plan
type RecoverEntities struct {
	EntityInfoList []*EntityInfoList `json:"entity_info_list,omitempty"`
}

// EntityInfoList represents a object for resource of recovery plan
type EntityInfoList struct {
	AnyEntityReference *Reference        `json:"any_entity_reference,omitempty"`
	Categories         map[string]string `json:"categories,omitempty"`
	ScriptList         []*ScriptList     `json:"script_list,omitempty"`
}

type ScriptList struct {
	EnableScriptExec *bool  `json:"enable_script_exec,omitempty"`
	Timeout          *int64 `json:"timeout,omitempty"`
}

// RecoveryPlanStatus represents a status of a recovery plan
type RecoveryPlanStatus struct {
	State            string                 `json:"state,omitempty"`
	MessageList      []*MessageResource     `json:"message_list,omitempty"`
	Name             string                 `json:"name,omitempty"`
	Resources        *RecoveryPlanResources `json:"resources,omitempty"`
	ExecutionContext *ExecutionContext      `json:"execution_context,omitempty"`
}

// RecoveryPlanSpec represents a spec of recovery plans
type RecoveryPlanSpec struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Resources   *RecoveryPlanResources `json:"resources,omitempty"`
}

// RecoveryPlanResponse represents a response object of a recovery plan
type RecoveryPlanResponse struct {
	APIVersion string              `json:"api_version,omitempty"`
	Metadata   *Metadata           `json:"metadata,omitempty"`
	Spec       *RecoveryPlanSpec   `json:"spec,omitempty"`
	Status     *RecoveryPlanStatus `json:"status,omitempty"`
}

// RecoveryPlanListResponse represents the response of a list of recovery plans
type RecoveryPlanListResponse struct {
	APIVersion string                  `json:"api_version,omitempty"`
	Entities   []*RecoveryPlanResponse `json:"entities,omitempty"`
	Metadata   *ListMetadataOutput     `json:"metadata,omitempty"`
}

// RecoveryPlanInput Represents the request of create recovery plan
type RecoveryPlanInput struct {
	APIVersion string            `json:"api_version,omitempty"`
	Metadata   *Metadata         `json:"metadata,omitempty"`
	Spec       *RecoveryPlanSpec `json:"spec,omitempty"`
}

type ServiceListEntry struct {
	Protocol         *string                        `json:"protocol,omitempty"`
	TCPPortRangeList []*PortRange                   `json:"tcp_port_range_list,omitempty"`
	UDPPortRangeList []*PortRange                   `json:"udp_port_range_list,omitempty"`
	IcmpTypeCodeList []*NetworkRuleIcmpTypeCodeList `json:"icmp_type_code_list,omitempty"`
}

type ServiceGroupListEntry struct {
	UUID                   *string            `json:"uuid,omitempty"`
	ServiceGroup           *ServiceGroupInput `json:"service_group,omitempty"`
	AssociatedPoliciesList []*ReferenceValues `json:"associated_policies_list,omitempty"`
}

type ServiceGroupInput struct {
	Name          *string             `json:"name,omitempty"`
	Description   *string             `json:"description,omitempty"`
	ServiceList   []*ServiceListEntry `json:"service_list,omitempty"`
	SystemDefined *bool               `json:"is_system_defined,omitempty"`
}

type ServiceGroupListResponse struct {
	Metadata *ListMetadataOutput      `json:"metadata,omitempty"`
	Entities []*ServiceGroupListEntry `json:"entities,omitempty"`
}

type ServiceGroupResponse struct {
	ServiceGroup *ServiceGroupInput `json:"service_group,omitempty"`
	UUID         *string            `json:"uuid,omitempty"`
}

type IPAddressBlock struct {
	IPAddress    *string `json:"ip,omitempty"`
	PrefixLength *int64  `json:"prefix_length,omitempty"`
}

type AddressGroupInput struct {
	Name               *string           `json:"name,omitempty"`
	Description        *string           `json:"description,omitempty"`
	BlockList          []*IPAddressBlock `json:"ip_address_block_list,omitempty"`
	AddressGroupString *string           `json:"address_group_string,omitempty"`
}

type AddressGroupResponse struct {
	UUID         *string            `json:"uuid,omitempty"`
	AddressGroup *AddressGroupInput `json:"address_group,omitempty"`
}

type AddressGroupListEntry struct {
	AddressGroup           *AddressGroupInput `json:"address_group,omitempty"`
	AssociatedPoliciesList []*ReferenceValues `json:"associated_policies_list,omitempty"`
}

type AddressGroupListResponse struct {
	Metadata *ListMetadataOutput      `json:"metadata,omitempty"`
	Entities []*AddressGroupListEntry `json:"entities,omitempty"`
}

type RecoveryPlanJobIntentInput struct {
	APIVersion *string          `json:"api_version,omitempty"`
	Metadata   *Metadata        `json:"metadata"`
	Spec       *RecoveryPlanJob `json:"spec"`
}

type RecoveryPlanJob struct {
	Name      *string                   `json:"name"`
	Resources *RecoveryPlanJobResources `json:"resources"`
}

type RecoveryPlanJobResources struct {
	ExecutionParameters   *RecoveryPlanJobResourcesExecutionParameters `json:"execution_parameters"`
	RecoveryPlanReference *Reference                                   `json:"recovery_plan_reference"`
}

type RecoveryPlanJobResourcesExecutionParameters struct {
	ActionType                        *string                 `json:"action_type"`
	FailedAvailabilityZoneList        []*AvailabilityZoneList `json:"failed_availability_zone_list"`
	RecoveryAvailabilityZoneList      []*AvailabilityZoneList `json:"recovery_availability_zone_list"`
	RecoveryReferenceTime             *time.Time              `json:"recovery_reference_time,omitempty"`
	ShouldContinueOnValidationFailure *bool                   `json:"should_continue_on_validation_failure,omitempty"`
}

type RecoveryPlanJobIntentResponse struct {
	APIVersion *string                   `json:"api_version"`
	Metadata   *Metadata                 `json:"metadata"`
	Spec       *RecoveryPlanJob          `json:"spec,omitempty"`
	Status     *RecoveryPlanJobDefStatus `json:"status,omitempty"`
}

type RecoveryPlanJobDefStatus struct {
	CleanupStatus                  *RecoveryPlanJobExecutionPhasesStatus `json:"cleanup_status,omitempty"`
	EndTime                        *time.Time                            `json:"end_time,omitempty"`
	ExecutionStatus                *RecoveryPlanJobExecutionPhasesStatus `json:"execution_status,omitempty"`
	Name                           *string                               `json:"name"`
	ParentRecoveryPlanJobReference *Reference                            `json:"parent_recovery_plan_job_reference,omitempty"`
	RecoveryPlanSpecification      *RecoveryPlanSpec                     `json:"recovery_plan_specification,omitempty"`
	Resources                      *RecoveryPlanJobResources             `json:"resources"`
	RootRecoveryPlanJobReference   *Reference                            `json:"root_recovery_plan_job_reference,omitempty"`
	StartTime                      *time.Time                            `json:"start_time,omitempty"`
	ValidationInformation          *RecoveryPlanJobValidationInformation `json:"validation_information,omitempty"`
	WitnessAddress                 string                                `json:"witness_address,omitempty"`
}

type RecoveryPlanJobExecutionPhasesStatus struct {
	OperationStatus      *RecoveryPlanJobStepStatus `json:"operation_status,omitempty"`
	PercentageComplete   *int32                     `json:"percentage_complete"`
	PostprocessingStatus *RecoveryPlanJobStepStatus `json:"postprocessing_status,omitempty"`
	PreprocessingStatus  *RecoveryPlanJobStepStatus `json:"preprocessing_status,omitempty"`
	Status               *string                    `json:"status"`
}

type RecoveryPlanJobStepStatus struct {
	PercentageComplete *int32  `json:"percentage_complete"`
	Status             *string `json:"status"`
}

type RecoveryPlanJobValidationInformation struct {
	ErrorsList   []*RecoveryPlanValidationMessage `json:"errors_list"`
	WarningsList []*RecoveryPlanValidationMessage `json:"warnings_list"`
}

type RecoveryPlanValidationMessage struct {
	AffectedAnyReferenceList      []*Reference                 `json:"affected_any_reference_list"`
	CauseAndResolutionMessageList []*CauseAndResolutionMessage `json:"cause_and_resolution_message_list"`
	ImpactMessageList             []string                     `json:"impact_message_list"`
	Message                       *string                      `json:"message"`
	ValidationType                *string                      `json:"validation_type"`
}

type CauseAndResolutionMessage struct {
	Cause          *string  `json:"cause"`
	ResolutionList []string `json:"resolution_list"`
}

type RecoveryPlanJobResponse struct {
	TaskUUID string `json:"task_uuid,omitempty"`
}

type RecoveryPlanJobExecutionStatus struct {
	OperationStatus      *RecoveryPlanJobPhaseExecutionStatus `json:"operation_status,omitempty"`
	PostprocessingStatus *RecoveryPlanJobPhaseExecutionStatus `json:"postprocessing_status,omitempty"`
	PreprocessingStatus  *RecoveryPlanJobPhaseExecutionStatus `json:"preprocessing_status,omitempty"`
}

type RecoveryPlanJobPhaseExecutionStatus struct {
	PercentageComplete      int32                                 `json:"percentage_complete,omitempty"`
	Status                  string                                `json:"status,omitempty"`
	StepExecutionStatusList []*RecoveryPlanJobStepExecutionStatus `json:"step_execution_status_list"`
}

type RecoveryPlanJobStepExecutionStatus struct {
	AnyEntityReferenceList  []*Reference                  `json:"any_entity_reference_list"`
	EndTime                 *time.Time                    `json:"end_time,omitempty"`
	ErrorCode               string                        `json:"error_code,omitempty"`
	ErrorDetail             string                        `json:"error_detail,omitempty"`
	Message                 string                        `json:"message,omitempty"`
	OperationType           *string                       `json:"operation_type"`
	ParentStepUUID          string                        `json:"parent_step_uuid,omitempty"`
	PercentageComplete      int32                         `json:"percentage_complete,omitempty"`
	RecoveredEntityInfoList []*RecoveredEntityInformation `json:"recovered_entity_info_list"`
	StartTime               *time.Time                    `json:"start_time,omitempty"`
	Status                  *string                       `json:"status"`
	StepSequenceNumber      int64                         `json:"step_sequence_number,omitempty"`
	StepUUID                *string                       `json:"step_uuid"`
}

type RecoveredEntityInformation struct {
	ErrorDetail           string           `json:"error_detail,omitempty"`
	RecoveredEntityInfo   *RecoveredEntity `json:"recovered_entity_info,omitempty"`
	SourceEntityReference *Reference       `json:"source_entity_reference,omitempty"`
}

type RecoveredEntity struct {
	EntityName string `json:"entity_name,omitempty"`
	EntityUUID string `json:"entity_uuid,omitempty"`
}

type RecoveryPlanJobListResponse struct {
	APIVersion *string                          `json:"api_version"`
	Entities   []*RecoveryPlanJobIntentResponse `json:"entities"`
	Metadata   *ListMetadataOutput              `json:"metadata"`
}

type RecoveryPlanJobActionRequest struct {
	RerunRecoveryPlanJobUUID               string `json:"rerun_recovery_plan_job_uuid,omitempty"`
	ShouldContinueRerunOnValidationFailure *bool  `json:"should_continue_rerun_on_validation_failure,omitempty"`
}

type GroupsRequestedAttribute struct {
	Attribute *string `json:"attribute"`
}

type GroupsGetEntitiesRequest struct {
	EntityType            *string                     `json:"entity_type"`
	FilterCriteria        string                      `json:"filter_criteria,omitempty"`
	GroupMemberAttributes []*GroupsRequestedAttribute `json:"group_member_attributes"`
}

type GroupsGetEntitiesResponse struct {
	FilteredGroupCount int64                `json:"filtered_group_count,omitempty"`
	GroupResults       []*GroupsGroupResult `json:"group_results"`
}

type GroupsGroupResult struct {
	EntityResults []*GroupsEntity `json:"entity_results"`
}

type GroupsEntity struct {
	Data     []*GroupsFieldData `json:"data"`
	EntityID string             `json:"entity_id,omitempty"`
}

type GroupsFieldData struct {
	Name   string                 `json:"name,omitempty"`
	Values []*GroupsTimevaluePair `json:"values"`
}

type GroupsTimevaluePair struct {
	Time   int64    `json:"time,omitempty"`
	Values []string `json:"values"`
}

type AvailabilityZoneIntentResponse struct {
	APIVersion *string                 `json:"api_version"`
	Metadata   *Metadata               `json:"metadata"`
	Spec       *AvailabilityZoneSpec   `json:"spec,omitempty"`
	Status     *AvailabilityZoneStatus `json:"status,omitempty"`
}

// AvailabilityZone Input Definition.
type AvailabilityZoneSpec struct {
	Name      *string                    `json:"name,omitempty"`      // The name of the AZ
	Resources *AvailabilityZoneResources `json:"resources,omitempty"` // AvailabilityZone Resource Definition
}

// AvailabilityZone Resource Definition
type AvailabilityZoneResources struct {
	ManagementUrl       *string `json:"management_url,omitempty"` // The URL of the management server
	ManagementPlaneType *string `json:"management_plane_type"`    // The type of the management plane
}

// AvailabilityZone status definition.
type AvailabilityZoneStatus struct {
	Name        *string                    `json:"name,omitempty"`      // The name of the AZ
	Resources   *AvailabilityZoneResources `json:"resources,omitempty"` // AvailabilityZone Resource Definition
	MessageList []MessageResource          `json:"message_list,omitempty"`
	State       *string                    `json:"state,omitempty"` // The state of the entity
}
