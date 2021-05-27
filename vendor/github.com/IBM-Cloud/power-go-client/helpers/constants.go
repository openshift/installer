package helpers

import "time"

const (
	// IBM PI Instance

	PIInstanceName            = "pi_instance_name"
	PIInstanceDate            = "pi_creation_date"
	PIInstanceSSHKeyName      = "pi_key_pair_name"
	PIInstanceImageName       = "pi_image_id"
	PIInstanceProcessors      = "pi_processors"
	PIInstanceProcType        = "pi_proc_type"
	PIInstanceMemory          = "pi_memory"
	PIInstanceSystemType      = "pi_sys_type"
	PIInstanceId              = "pi_instance_id"
	PIInstanceDiskSize        = "pi_disk_size"
	PIInstanceStatus          = "pi_instance_status"
	PIInstanceMinProc         = "pi_minproc"
	PIInstanceVolumeIds       = "pi_volume_ids"
	PIInstanceNetworkIds      = "pi_network_ids"
	PIInstancePublicNetwork   = "pi_public_network"
	PIInstanceMigratable      = "pi_migratable"
	PICloudInstanceId         = "pi_cloud_instance_id"
	PICloudInstanceSubnetName = "pi_cloud_instance_subnet_name"
	PIInstanceMimMem          = "pi_minmem"
	PIInstanceMaxProc         = "pi_maxproc"
	PIInstanceMaxMem          = "pi_maxmem"
	PIInstanceReboot          = "pi_reboot"
	PITenantId                = "pi_tenant_id"
	PIVirtualCoresAssigned    = "pi_virtual_cores_assigned"
	PIVirtualCoresMax         = "pi_virtual_cores_max"
	PIVirtualCoresMin         = "pi_virutal_cores_min"

	PIInstanceHealthStatus      = "pi_health_status"
	PIInstanceReplicants        = "pi_replicants"
	PIInstanceReplicationPolicy = "pi_replication_policy"
	PIInstanceReplicationScheme = "pi_replication_scheme"
	PIInstanceProgress          = "pi_progress"
	PIInstanceUserData          = "pi_user_data"
	PIInstancePinPolicy         = "pi_pin_policy"

	// IBM PI Volume
	PIVolumeName              = "pi_volume_name"
	PIVolumeSize              = "pi_volume_size"
	PIVolumeType              = "pi_volume_type"
	PIVolumeShareable         = "pi_volume_shareable"
	PIVolumeId                = "pi_volume_id"
	PIVolumeStatus            = "pi_volume_status"
	PIVolumeWWN               = "pi_volume_wwn"
	PIVolumeDeleteOnTerminate = "pi_volume_delete_on_terminate"
	PIVolumeCreateDate        = "pi_volume_create_date"
	PIVolumeLastUpdate        = "pi_last_updated_date"
	PIVolumePool              = "pi_volume_pool"
	PIAffinityPolicy          = "pi_volume_affinity_policy"
	PIAffinityVolume          = "pi_volume_affinity"

	// IBM PI Snapshots

	PISnapshot         = "pi_snap_shot_id"
	PISnapshotName     = "pi_snap_shot_name"
	PISnapshotStatus   = "pi_snap_shot_status"
	PISnapshotAction   = "pi_snap_shot_action"
	PISnapshotComplete = "pi_snap_shot_complete"

	// IBM PI Image

	PIImageName       = "pi_image_name"
	PIImageAccessKey  = "pi_image_access_key"
	PIImageSecretKey  = "pi_image_secret_key"
	PIImageSource     = "pi_image_source"
	PIImageBucketName = "pi_image_bucket_name"
	PIImageFileName   = "pi_image_file_name"
	PIImageRegion     = "pi_image_region"
	PIImageDisk       = "pi_image_disk"
	PIImageCopyID     = "pi_image_copy_id"
	PIImagePath       = "pi_image_path"
	PIImageOsType     = "pi_image_os_type"

	// IBM PI Key

	PIKeyName = "pi_key_name"
	PIKey     = "pi_ssh_key"
	PIKeyDate = "pi_creation_date"
	PIKeyId   = "pi_key_id"

	// IBM PI Network

	PINetworkReady           = "ready"
	PINetworkID              = "pi_networkid"
	PINetworkName            = "pi_network_name"
	PINetworkCidr            = "pi_cidr"
	PINetworkDNS             = "pi_dns"
	PINetworkType            = "pi_network_type"
	PINetworkGateway         = "pi_gateway"
	PINetworkIPAddressRange  = "pi_ipaddress_range"
	PINetworkVlanId          = "pi_vlan_id"
	PINetworkProvisioning    = "build"
	PINetworkPortDescription = "pi_network_port_description"
	PINetworkPortIPAddress   = "pi_network_port_ipaddress"
	PINetworkPortMacAddress  = "pi_network_port_macaddress"
	PINetworkPortStatus      = "pi_network_port_status"
	PINetworkPortPortID      = "pi_network_port_portid"

	// IBM PI Operations
	PIInstanceOperationType       = "pi_operation"
	PIInstanceOperationProgress   = "pi_progress"
	PIInstanceOperationStatus     = "pi_status"
	PIInstanceOperationServerName = "pi_instance_name"

	// IBM PI Volume Attach
	PIVolumeAttachName                = "pi_volume_attach_name"
	PIVolumeAllowableAttachStatus     = "in-use"
	PIVolumeAttachStatus              = "status"
	PowerVolumeAttachDeleting         = "deleting"
	PowerVolumeAttachProvisioning     = "creating"
	PowerVolumeAttachProvisioningDone = "available"

	// IBM PI Instance Capture
	PIInstanceCaptureName                  = "pi_capture_name"
	PIInstanceCaptureDestination           = "pi_capture_destination"
	PIInstanceCaptureVolumeIds             = "pi_capture_volume_ids"
	PIInstanceCaptureCloudStorageImagePath = "pi_capture_storage_image_path"
	PIInstanceCaptureCloudStorageRegion    = "pi_capture_cloud_storage_region"
	PIInstanceCaptureCloudStorageAccessKey = "pi_capture_cloud_storage_access_key"
	PIInstanceCaptureCloudStorageSecretKey = "pi_capture_cloud_storage_secret_key"

	// Status For all the resources

	PIVolumeDeleting         = "deleting"
	PIVolumeDeleted          = "done"
	PIVolumeProvisioning     = "creating"
	PIVolumeProvisioningDone = "available"
	PIInstanceAvailable      = "ACTIVE"
	PIInstanceHealthOk       = "OK"
	PIInstanceHealthWarning  = "WARNING"
	PIInstanceBuilding       = "BUILD"
	PIInstanceDeleting       = "DELETING"
	PIInstanceNotFound       = "Not Found"
	PIImageQueStatus         = "queued"
	PIImageActiveStatus      = "active"

	//Timeout values for Power VS -

	PICreateTimeOut = 5 * time.Minute
	PIDeleteTimeOut = 3 * time.Minute
	PIGetTimeOut    = 2 * time.Minute
)
