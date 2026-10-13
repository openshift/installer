/*
Copyright 2024 The ORC Authors.

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

package v1alpha1

// +kubebuilder:validation:MinLength:=1
// +kubebuilder:validation:MaxLength:=80
type ServerTag string

type FilterByServerTags struct {
	// tags is a list of tags to filter by. If specified, the resource must
	// have all of the tags specified to be included in the result.
	// +listType=set
	// +optional
	// +kubebuilder:validation:MaxItems:=50
	Tags []ServerTag `json:"tags,omitempty"`

	// tagsAny is a list of tags to filter by. If specified, the resource
	// must have at least one of the tags specified to be included in the
	// result.
	// +listType=set
	// +optional
	// +kubebuilder:validation:MaxItems:=50
	TagsAny []ServerTag `json:"tagsAny,omitempty"`

	// notTags is a list of tags to filter by. If specified, resources which
	// contain all of the given tags will be excluded from the result.
	// +listType=set
	// +optional
	// +kubebuilder:validation:MaxItems:=50
	NotTags []ServerTag `json:"notTags,omitempty"`

	// notTagsAny is a list of tags to filter by. If specified, resources
	// which contain any of the given tags will be excluded from the result.
	// +listType=set
	// +optional
	// +kubebuilder:validation:MaxItems:=50
	NotTagsAny []ServerTag `json:"notTagsAny,omitempty"`
}

// +kubebuilder:validation:MinProperties:=1
// +kubebuilder:validation:MaxProperties:=1
type ServerPortSpec struct {
	// portRef is a reference to a Port object. Server creation will wait for
	// this port to be created and available.
	// +optional
	PortRef *KubernetesNameRef `json:"portRef,omitempty"`
}

// +kubebuilder:validation:MinProperties:=1
type ServerVolumeSpec struct {
	// volumeRef is a reference to a Volume object. Server creation will wait for
	// this volume to be created and available.
	// +required
	VolumeRef KubernetesNameRef `json:"volumeRef,omitempty"`

	// device is the name of the device, such as `/dev/vdb`.
	// Omit for auto-assignment
	// +kubebuilder:validation:MaxLength:=255
	// +optional
	Device *string `json:"device,omitempty"`
}

type ServerVolumeStatus struct {
	// id is the ID of a volume attached to the server.
	// +kubebuilder:validation:MaxLength:=1024
	// +optional
	ID string `json:"id,omitempty"`
}

type ServerInterfaceFixedIP struct {
	// ipAddress is the IP address assigned to the port.
	// +kubebuilder:validation:MaxLength:=1024
	// +optional
	IPAddress string `json:"ipAddress,omitempty"`

	// subnetID is the ID of the subnet from which the IP address is allocated.
	// +kubebuilder:validation:MaxLength:=1024
	// +optional
	SubnetID string `json:"subnetID,omitempty"`
}

type ServerInterfaceStatus struct {
	// portID is the ID of a port attached to the server.
	// +kubebuilder:validation:MaxLength:=1024
	// +optional
	PortID string `json:"portID,omitempty"`

	// netID is the ID of the network to which the interface is attached.
	// +kubebuilder:validation:MaxLength:=1024
	// +optional
	NetID string `json:"netID,omitempty"`

	// macAddr is the MAC address of the interface.
	// +kubebuilder:validation:MaxLength:=1024
	// +optional
	MACAddr string `json:"macAddr,omitempty"`

	// portState is the state of the port (e.g., ACTIVE, DOWN).
	// +kubebuilder:validation:MaxLength:=1024
	// +optional
	PortState string `json:"portState,omitempty"`

	// fixedIPs is the list of fixed IP addresses assigned to the interface.
	// +kubebuilder:validation:MaxItems:=32
	// +listType=atomic
	// +optional
	FixedIPs []ServerInterfaceFixedIP `json:"fixedIPs,omitempty"`
}

// ServerResourceSpec contains the desired state of a server
type ServerResourceSpec struct {
	// name will be the name of the created resource. If not specified, the
	// name of the ORC object will be used.
	// +optional
	Name *OpenStackName `json:"name,omitempty"`

	// imageRef references the image to use for the server instance.
	// NOTE: This is not required in case of boot from volume.
	// +required
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="imageRef is immutable"
	ImageRef KubernetesNameRef `json:"imageRef,omitempty"`

	// flavorRef references the flavor to use for the server instance.
	// +required
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="flavorRef is immutable"
	FlavorRef KubernetesNameRef `json:"flavorRef,omitempty"`

	// userData specifies data which will be made available to the server at
	// boot time, either via the metadata service or a config drive. It is
	// typically read by a configuration service such as cloud-init or ignition.
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="userData is immutable"
	UserData *UserDataSpec `json:"userData,omitempty"`

	// ports defines a list of ports which will be attached to the server.
	// +kubebuilder:validation:MaxItems:=64
	// +listType=atomic
	// +required
	Ports []ServerPortSpec `json:"ports,omitempty"`

	// volumes is a list of volumes attached to the server.
	// +kubebuilder:validation:MaxItems:=64
	// +listType=atomic
	// +optional
	Volumes []ServerVolumeSpec `json:"volumes,omitempty"`

	// serverGroupRef is a reference to a ServerGroup object. The server
	// will be created in the server group.
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="serverGroupRef is immutable"
	ServerGroupRef *KubernetesNameRef `json:"serverGroupRef,omitempty"`

	// availabilityZone is the availability zone in which to create the server.
	// +kubebuilder:validation:MaxLength=255
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="availabilityZone is immutable"
	AvailabilityZone string `json:"availabilityZone,omitempty"`

	// tags is a list of tags which will be applied to the server.
	// +kubebuilder:validation:MaxItems:=50
	// +listType=set
	// +optional
	Tags []ServerTag `json:"tags,omitempty"`
}

// +kubebuilder:validation:MinProperties:=1
// +kubebuilder:validation:MaxProperties:=1
type UserDataSpec struct {
	// secretRef is a reference to a Secret containing the user data for this server.
	// +optional
	SecretRef *KubernetesNameRef `json:"secretRef,omitempty"`
}

// ServerFilter defines an existing resource by its properties
// +kubebuilder:validation:MinProperties:=1
type ServerFilter struct {
	// name of the existing resource
	// +optional
	Name *OpenStackName `json:"name,omitempty"`

	// availabilityZone is the availability zone of the existing resource
	// +kubebuilder:validation:MaxLength=255
	// +optional
	AvailabilityZone string `json:"availabilityZone,omitempty"`

	FilterByServerTags `json:",inline"`
}

// ServerResourceStatus represents the observed state of the resource.
type ServerResourceStatus struct {
	// name is the human-readable name of the resource. Might not be unique.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Name string `json:"name,omitempty"`

	// hostID is the host where the server is located in the cloud.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	HostID string `json:"hostID,omitempty"`

	// status contains the current operational status of the server,
	// such as IN_PROGRESS or ACTIVE.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Status string `json:"status,omitempty"`

	// imageID indicates the OS image used to deploy the server.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	ImageID string `json:"imageID,omitempty"`

	// availabilityZone is the availability zone where the server is located.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	AvailabilityZone string `json:"availabilityZone,omitempty"`

	// serverGroups is a slice of strings containing the UUIDs of the
	// server groups to which the server belongs. Currently this can
	// contain at most one entry.
	// +kubebuilder:validation:MaxItems:=32
	// +kubebuilder:validation:items:MaxLength=1024
	// +listType=atomic
	// +optional
	ServerGroups []string `json:"serverGroups,omitempty"`

	// volumes contains the volumes attached to the server.
	// +kubebuilder:validation:MaxItems:=64
	// +listType=atomic
	// +optional
	Volumes []ServerVolumeStatus `json:"volumes,omitempty"`

	// interfaces contains the list of interfaces attached to the server.
	// +kubebuilder:validation:MaxItems:=64
	// +listType=atomic
	// +optional
	Interfaces []ServerInterfaceStatus `json:"interfaces,omitempty"`

	// tags is the list of tags on the resource.
	// +kubebuilder:validation:MaxItems:=50
	// +kubebuilder:validation:items:MaxLength=1024
	// +listType=atomic
	// +optional
	Tags []string `json:"tags,omitempty"`
}
