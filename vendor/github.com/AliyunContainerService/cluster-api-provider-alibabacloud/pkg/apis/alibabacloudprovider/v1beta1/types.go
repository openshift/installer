/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

// InstanceTenancy Specifies whether to create the instance on a dedicated host
type InstanceTenancy string

const (
	// DefaultTenancy creates the instance on a non-dedicated host.
	DefaultTenancy InstanceTenancy = "default"

	// HostTenancy creates the instance on a dedicated host. If you do not specify the DedicatedHostID parameter, Alibaba Cloud automatically selects a dedicated host for the instance.
	HostTenancy InstanceTenancy = "host"
)

// DataDisk The datadisk of Instance
type DataDisk struct {
	//The name of data disk N. The name must be 2 to 128 characters in length. It must start with a letter and cannot start with http:// or https://. It can contain letters, digits, colons (:), underscores (_), and hyphens (-).
	//
	//This parameter is empty by default.
	DiskName string `name:"diskName,omitempty"`

	//The ID of the snapshot used to create data disk N. Valid values of N: 1 to 16.
	//
	//When the DataDisk.N.SnapshotID parameter is specified, the DataDisk.N.Size parameter is ignored. The data disk is created based on the size of the specified snapshot.
	//Use snapshots created after July 15, 2013. Otherwise, an error is returned and your request is rejected.
	SnapshotID string `name:"snapshotId,omitempty"`

	//The size of data disk N. Valid values of N: 1 to 16. Unit: GiB. Valid values:
	//
	//Valid values when DataDisk.N.Category is set to cloud_efficiency: 20 to 32768
	//Valid values when DataDisk.N.Category is set to cloud_ssd: 20 to 32768
	//Valid values when DataDisk.N.Category is set to cloud_essd: 20 to 32768
	//Valid values when DataDisk.N.Category is set to cloud: 5 to 2000
	//The value of this parameter must be greater than or equal to the size of the snapshot specified by the SnapshotID parameter.
	Size int `name:"size,omitempty"`

	//Specifies whether to encrypt data disk N.
	//
	//Default value: false.
	Encrypted bool `name:"encrypted,omitempty"`

	//
	//The performance level of the ESSD used as data disk N. The N value must be the same as that in DataDisk.N.Category when DataDisk.N.Category is set to cloud_essd. Default value: PL1. Valid values:
	//
	//PL0: A single ESSD can deliver up to 10,000 random read/write IOPS.
	//PL1: A single ESSD can deliver up to 50,000 random read/write IOPS.
	//PL2: A single ESSD can deliver up to 100,000 random read/write IOPS.
	//PL3: A single ESSD can deliver up to 1,000,000 random read/write IOPS.
	//For more information about ESSD performance levels, see ESSDs.
	PerformanceLevel string `name:"performanceLevel,omitempty"`

	//TODO
	//EncryptAlgorithm string `name:"EncryptAlgorithm"`

	//The description of data disk N. The description must be 2 to 256 characters in length and cannot start with http:// or https://.
	//
	//This parameter is empty by default.
	Description string `name:"description,omitempty"`

	//The category of data disk N. Valid values:
	//
	//cloud_efficiency: ultra disk
	//cloud_ssd: standard SSD
	//cloud_essd: ESSD
	//cloud: basic disk
	//For I/O optimized instances, the default value is cloud_efficiency. For non-I/O optimized instances, the default value is cloud.
	Category string `name:"category,omitempty"`

	//The ID of the Key Management Service (KMS) key to be used by data disk N.
	KMSKeyID string `name:"kmsKeyId,omitempty"`

	//The mount point of data disk N.
	Device string `name:"device,omitempty"`

	//Specifies whether to release data disk N along with the instance.
	//
	//Default value: true.
	DeleteWithInstance *bool `name:"deleteWithInstance,omitempty"`
}

// Tag  The tags of ECS Instance
type Tag struct {
	Value string `name:"value"`
	Key   string `name:"key"`
}
