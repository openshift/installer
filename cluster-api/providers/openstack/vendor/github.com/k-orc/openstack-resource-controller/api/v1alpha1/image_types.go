/*
Copyright 2024 The Kubernetes Authors.

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

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GlanceTag is the name of the go field tag in properties structs used to specify the Glance property name.
const GlanceTag = "glance"

// +kubebuilder:validation:MinLength:=1
// +kubebuilder:validation:MaxLength:=255
type ImageTag string

// +kubebuilder:validation:Enum:=ami;ari;aki;bare;ovf;ova;docker
type ImageContainerFormat string

const (
	ImageContainerFormatAKI    ImageContainerFormat = "aki"
	ImageContainerFormatAMI    ImageContainerFormat = "ami"
	ImageContainerFormatARI    ImageContainerFormat = "ari"
	ImageContainerFormatBare   ImageContainerFormat = "bare"
	ImageContainerFormatDocker ImageContainerFormat = "docker"
	ImageContainerFormatOVA    ImageContainerFormat = "ova"
	ImageContainerFormatOVF    ImageContainerFormat = "ovf"
)

// +kubebuilder:validation:Enum:=ami;ari;aki;vhd;vhdx;vmdk;raw;qcow2;vdi;ploop;iso
type ImageDiskFormat string

const (
	ImageDiskFormatAMI   ImageDiskFormat = "ami"
	ImageDiskFormatARI   ImageDiskFormat = "ari"
	ImageDiskFormatAKI   ImageDiskFormat = "aki"
	ImageDiskFormatVHD   ImageDiskFormat = "vhd"
	ImageDiskFormatVHDX  ImageDiskFormat = "vhdx"
	ImageDiskFormatVMDK  ImageDiskFormat = "vmdk"
	ImageDiskFormatRaw   ImageDiskFormat = "raw"
	ImageDiskFormatQCOW2 ImageDiskFormat = "qcow2"
	ImageDiskFormatVDI   ImageDiskFormat = "vdi"
	ImageDiskFormatPLoop ImageDiskFormat = "ploop"
	ImageDiskFormatISO   ImageDiskFormat = "iso"
)

// +kubebuilder:validation:Enum:=public;private;shared;community
type ImageVisibility string

const (
	ImageVisibilityPublic    ImageVisibility = "public"
	ImageVisibilityPrivate   ImageVisibility = "private"
	ImageVisibilityShared    ImageVisibility = "shared"
	ImageVisibilityCommunity ImageVisibility = "community"
)

// +kubebuilder:validation:Enum:=md5;sha1;sha256;sha512
type ImageHashAlgorithm string

const (
	ImageHashAlgorithmMD5    ImageHashAlgorithm = "md5"
	ImageHashAlgorithmSHA1   ImageHashAlgorithm = "sha1"
	ImageHashAlgorithmSHA256 ImageHashAlgorithm = "sha256"
	ImageHashAlgorithmSHA512 ImageHashAlgorithm = "sha512"
)

// See https://docs.openstack.org/glance/latest/admin/useful-image-properties.html
// for a list of 'well known' image properties we might consider supporting explicitly.
//
// The set of supported properties is currently arbitrarily selective. We should
// add supported options here freely.

// ImageHWBus is a type of hardware bus.
//
// Permitted values are scsi, virtio, uml, xen, ide, usb, and lxc.
// +kubebuilder:validation:Enum:=scsi;virtio;uml;xen;ide;usb;lxc
type ImageHWBus string

type ImagePropertiesHardware struct {
	// CPUSockets is the preferred number of sockets to expose to the guest
	// +optional
	CPUSockets *int `json:"cpuSockets,omitempty" glance:"hw_cpu_sockets"`

	// CPUCores is the preferred number of cores to expose to the guest
	// +optional
	CPUCores *int `json:"cpuCores,omitempty" glance:"hw_cpu_cores"`

	// CPUThreads is the preferred number of threads to expose to the guest
	// +optional
	CPUThreads *int `json:"cpuThreads,omitempty" glance:"hw_cpu_threads"`

	// CPUPolicy is used to pin the virtual CPUs (vCPUs) of instances to the
	// host's physical CPU cores (pCPUs). Host aggregates should be used to
	// separate these pinned instances from unpinned instances as the latter
	// will not respect the resourcing requirements of the former.
	//
	// Permitted values are shared (the default), and dedicated.
	//
	// shared: The guest vCPUs will be allowed to freely float across host
	// pCPUs, albeit potentially constrained by NUMA policy.
	//
	// dedicated: The guest vCPUs will be strictly pinned to a set of host
	// pCPUs. In the absence of an explicit vCPU topology request, the
	// drivers typically expose all vCPUs as sockets with one core and one
	// thread. When strict CPU pinning is in effect the guest CPU topology
	// will be setup to match the topology of the CPUs to which it is
	// pinned. This option implies an overcommit ratio of 1.0. For example,
	// if a two vCPU guest is pinned to a single host core with two threads,
	// then the guest will get a topology of one socket, one core, two
	// threads.
	// +kubebuilder:validation:Enum:=shared;dedicated
	// +optional
	CPUPolicy *string `json:"cpuPolicy,omitempty" glance:"hw_cpu_policy"`

	// CPUThreadPolicy further refines a CPUPolicy of 'dedicated' by stating
	// how hardware CPU threads in a simultaneous multithreading-based (SMT)
	// architecture be used. SMT-based architectures include Intel
	// processors with Hyper-Threading technology. In these architectures,
	// processor cores share a number of components with one or more other
	// cores. Cores in such architectures are commonly referred to as
	// hardware threads, while the cores that a given core share components
	// with are known as thread siblings.
	//
	// Permitted values are prefer (the default), isolate, and require.
	//
	// prefer: The host may or may not have an SMT architecture. Where an
	// SMT architecture is present, thread siblings are preferred.
	//
	// isolate: The host must not have an SMT architecture or must emulate a
	// non-SMT architecture. If the host does not have an SMT architecture,
	// each vCPU is placed on a different core as expected. If the host does
	// have an SMT architecture - that is, one or more cores have thread
	// siblings - then each vCPU is placed on a different physical core. No
	// vCPUs from other guests are placed on the same core. All but one
	// thread sibling on each utilized core is therefore guaranteed to be
	// unusable.
	//
	// require: The host must have an SMT architecture. Each vCPU is
	// allocated on thread siblings. If the host does not have an SMT
	// architecture, then it is not used. If the host has an SMT
	// architecture, but not enough cores with free thread siblings are
	// available, then scheduling fails.
	// +kubebuilder:validation:Enum:=prefer;isolate;require
	// +optional
	CPUThreadPolicy *string `json:"cpuThreadPolicy,omitempty" glance:"hw_cpu_thread_policy"`

	// CDROMBus specifies the type of disk controller to attach CD-ROM devices to.
	// +optional
	CDROMBus *ImageHWBus `json:"cdromBus,omitempty" glance:"hw_cdrom_bus"`

	// DiskBus specifies the type of disk controller to attach disk devices to.
	// +optional
	DiskBus *ImageHWBus `json:"diskBus,omitempty" glance:"hw_disk_bus"`

	// TODO: hw_machine_type seems important to support early, but how to
	// select a supported set?

	// SCSIModel enables the use of VirtIO SCSI (virtio-scsi) to provide
	// block device access for compute instances; by default, instances use
	// VirtIO Block (virtio-blk). VirtIO SCSI is a para-virtualized SCSI
	// controller device that provides improved scalability and performance,
	// and supports advanced SCSI hardware.
	//
	// The only permitted value is virtio-scsi.
	// +kubebuilder:validation:Enum:=virtio-scsi
	// +optional
	SCSIModel *string `json:"scsiModel,omitempty" glance:"hw_scsi_model"`

	// VIFModel specifies the model of virtual network interface device to use.
	//
	// Permitted values are e1000, e1000e, ne2k_pci, pcnet, rtl8139, virtio,
	// and vmxnet3.
	// +kubebuilder:validation:Enum:=e1000;e1000e;ne2k_pci;pcnet;rtl8139;virtio;vmxnet3
	// +optional
	VIFModel *string `json:"vifModel,omitempty" glance:"hw_vif_model"`
}

type ImageProperties struct {
	// MinDisk is the minimum amount of disk space in GB that is required to boot the image
	// +kubebuilder:validation:Minimum:=1
	// +optional
	MinDiskGB *int `json:"minDiskGB,omitempty"`

	// MinMemoryMB is the minimum amount of RAM in MB that is required to boot the image.
	// +kubebuilder:validation:Minimum:=1
	// +optional
	MinMemoryMB *int `json:"minMemoryMB,omitempty"`

	// Hardware is a set of properties which control the virtual hardware
	// created by Nova.
	// +optional
	Hardware *ImagePropertiesHardware `json:"hardware,omitempty"`
}

// +kubebuilder:validation:Enum:=xz;gz;bz2
type ImageCompression string

const (
	ImageCompressionXZ  ImageCompression = "xz"
	ImageCompressionGZ  ImageCompression = "gz"
	ImageCompressionBZ2 ImageCompression = "bz2"
)

type ImageContent struct {
	// ContainerFormat is the format of the image container.
	// qcow2 and raw images do not usually have a container. This is specified as "bare", which is also the default.
	// Permitted values are ami, ari, aki, bare, ovf, ova, and docker.
	// +kubebuilder:default:=bare
	// +optional
	ContainerFormat ImageContainerFormat `json:"containerFormat,omitempty"`

	// DiskFormat is the format of the disk image.
	// Normal values are "qcow2", or "raw". Glance may be configured to support others.
	// +kubebuilder:validation:Required
	DiskFormat ImageDiskFormat `json:"diskFormat"`

	// Download describes how to obtain image data by downloading it from a URL.
	// Must be set when creating a managed image.
	// +kubebuilder:validation:Required
	Download *ImageContentSourceDownload `json:"download,omitempty"`
}

type ImageContentSourceDownload struct {
	// URL containing image data
	// +kubebuilder:validation:Format=uri
	// +kubebuilder:validation:Required
	URL string `json:"url"`

	// Decompress specifies that the source data must be decompressed with the
	// given compression algorithm before being stored. Specifying Decompress
	// will disable the use of Glance's web-download, as web-download cannot
	// currently deterministically decompress downloaded content.
	// +optional
	Decompress *ImageCompression `json:"decompress,omitempty"`

	// Hash is a hash which will be used to verify downloaded data, i.e.
	// before any decompression. If not specified, no hash verification will be
	// performed. Specifying a Hash will disable the use of Glance's
	// web-download, as web-download cannot currently deterministically verify
	// the hash of downloaded content.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="hash is immutable"
	// +optional
	Hash *ImageHash `json:"hash,omitempty"`
}

type ImageHash struct {
	// Algorithm is the hash algorithm used to generate value.
	// +kubebuilder:validation:Required
	Algorithm ImageHashAlgorithm `json:"algorithm"`

	// Value is the hash of the image data using Algorithm. It must be hex encoded using lowercase letters.
	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:MaxLength:=1024
	// +kubebuilder:validation:Pattern:=`^[0-9a-f]+$`
	// +kubebuilder:validation:Required
	Value string `json:"value"`
}

// ImageResourceSpec contains the desired state of a Glance image
// +kubebuilder:validation:XValidation:rule="has(self.name) ? self.name == oldSelf.name : !has(oldSelf.name)",message="name is immutable"
// +kubebuilder:validation:XValidation:rule="has(self.protected) ? self.protected == oldSelf.protected : !has(oldSelf.protected)",message="name is immutable"
// +kubebuilder:validation:XValidation:rule="has(self.tags) ? self.tags == oldSelf.tags : !has(oldSelf.tags)",message="tags is immutable"
// +kubebuilder:validation:XValidation:rule="has(self.visibility) ? self.visibility == oldSelf.visibility : !has(oldSelf.visibility)",message="visibility is immutable"
// +kubebuilder:validation:XValidation:rule="has(self.properties) ? self.properties == oldSelf.properties : !has(oldSelf.properties)",message="properties is immutable"
type ImageResourceSpec struct {
	// Name will be the name of the created Glance image. If not specified, the
	// name of the Image object will be used.
	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:MaxLength:=1000
	// +optional
	Name string `json:"name,omitempty"`

	// Protected specifies that the image is protected from deletion.
	// If not specified, the default is false.
	// +optional
	Protected *bool `json:"protected,omitempty"`

	// Tags is a list of tags which will be applied to the image. A tag has a maximum length of 255 characters.
	// +listType=set
	// +optional
	Tags []ImageTag `json:"tags,omitempty"`

	// Visibility of the image
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="visibility is immutable"
	// +optional
	Visibility *ImageVisibility `json:"visibility,omitempty"`

	// Properties is metadata available to consumers of the image
	// +optional
	Properties *ImageProperties `json:"properties,omitempty"`

	// Content specifies how to obtain the image content.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="content is immutable"
	// +optional
	Content *ImageContent `json:"content,omitempty"`
}

// ImageFilter defines a Glance query
// +kubebuilder:validation:MinProperties:=1
type ImageFilter struct {
	// Name specifies the name of a Glance image
	// +optional
	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:MaxLength:=1000
	Name *string `json:"name,omitempty"`
}

// ImageImport specifies an existing image which will be imported instead of
// creating a new image
// +kubebuilder:validation:MinProperties:=1
// +kubebuilder:validation:MaxProperties:=1
type ImageImport struct {
	// ID contains the unique identifier of an existing Glance image. Note that
	// when specifying an image import by ID, the image MUST already exist. The
	// Image will enter an error state if the image does not exist.
	// +optional
	// +kubebuilder:validation:Format:=uuid
	ID *string `json:"id,omitempty"`

	// Filter contains an image query which is expected to return a single
	// result. The controller will continue to retry if filter returns no
	// results. If filter returns multiple results the controller will set an
	// error state and will not continue to retry.
	// +optional
	Filter *ImageFilter `json:"filter,omitempty"`
}

// ImageSpec defines the desired state of an Image.
// +kubebuilder:validation:XValidation:rule="self.managementPolicy == 'managed' ? has(self.resource) : true",message="resource must be specified when policy is managed"
// +kubebuilder:validation:XValidation:rule="self.managementPolicy == 'managed' ? !has(self.__import__) : true",message="import may not be specified when policy is managed"
// +kubebuilder:validation:XValidation:rule="self.managementPolicy == 'unmanaged' ? !has(self.resource) : true",message="resource may not be specified when policy is unmanaged"
// +kubebuilder:validation:XValidation:rule="self.managementPolicy == 'unmanaged' ? has(self.__import__) : true",message="import must be specified when policy is unmanaged"
// +kubebuilder:validation:XValidation:rule="has(self.managedOptions) ? self.managementPolicy == 'managed' : true",message="managedOptions may only be provided when policy is managed"
// +kubebuilder:validation:XValidation:rule="!has(self.__import__) ? has(self.resource.content) : true",message="resource content must be specified when not importing"
type ImageSpec struct {
	// Import refers to an existing image which will be imported instead of
	// creating a new image.
	// +optional
	Import *ImageImport `json:"import,omitempty"`

	// Resource specifies the desired state of the Glance image.
	//
	// Resource may not be specified if the management policy is `unmanaged`.
	//
	// Resource must be specified when the management policy is `managed`.
	// +optional
	Resource *ImageResourceSpec `json:"resource,omitempty"`

	// ManagementPolicy defines how ORC will treat the object. Valid values are
	// `managed`: ORC will create, update, and delete the resource; `unmanaged`:
	// ORC will import an existing image, and will not apply updates to it or
	// delete it.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="managementPolicy is immutable"
	// +kubebuilder:default:=managed
	// +optional
	ManagementPolicy ManagementPolicy `json:"managementPolicy,omitempty"`

	// ManagedOptions specifies options which may be applied to managed objects.
	// +optional
	ManagedOptions *ManagedOptions `json:"managedOptions,omitempty"`

	// CloudCredentialsRef points to a secret containing OpenStack credentials
	// +kubebuilder:validation:Required
	CloudCredentialsRef CloudCredentialsReference `json:"cloudCredentialsRef"`
}

// ImageResourceStatus represents the observed state of a Glance image
type ImageResourceStatus struct {
	// Status is the image status as reported by Glance
	// +optional
	Status *string `json:"status,omitempty"`

	// Hash is the hash of the image data published by Glance. Note that this is
	// a hash of the data stored internally by Glance, which will have been
	// decompressed and potentially format converted depending on server-side
	// configuration which is not visible to clients. It is expected that this
	// hash will usually differ from the download hash.
	// +optional
	Hash *ImageHash `json:"hash,omitempty"`

	// SizeB is the size of the image data, in bytes
	// +optional
	SizeB *int64 `json:"sizeB,omitempty"`

	// VirtualSizeB is the size of the disk the image data represents, in bytes
	// +optional
	VirtualSizeB *int64 `json:"virtualSizeB,omitempty"`
}

// ImageStatus defines the observed state of an Image.
type ImageStatus struct {
	// Conditions represents the observed status of the object.
	// Known .status.conditions.type are: "Available", "Progressing"
	//
	// Available represents the availability of the Glance image. If it is
	// true then the image is ready for use in Glance, and its hash has been
	// verified.
	//
	// Progressing indicates the state of the Glance image does not currently
	// reflect the desired state, but that reconciliation is progressing.
	// Progressing will be False either because the desired state has been
	// achieved, or some terminal error prevents it from being achieved and the
	// controller is no longer attempting to reconcile.
	//
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// ID is the unique identifier of the Glance image
	// +optional
	ID *string `json:"id,omitempty"`

	// Resource contains the observed state of the Glance image
	// +optional
	Resource *ImageResourceStatus `json:"resource,omitempty"`

	// DownloadAttempts is the number of times the controller has attempted to download the image contents
	// +optional
	DownloadAttempts *int `json:"downloadAttempts,omitempty"`
}

var _ ObjectWithConditions = &Image{}

func (i *Image) GetConditions() []metav1.Condition {
	return i.Status.Conditions
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="ID",type="string",JSONPath=".status.id",description="Glance image ID"
// +kubebuilder:printcolumn:name="Available",type="string",JSONPath=".status.conditions[?(@.type=='Available')].status",description="Availability status of image"
// +kubebuilder:printcolumn:name="Message",type="string",JSONPath=".status.conditions[?(@.type=='Available')].message",description="Message describing current availability status"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation"

// Image is the Schema for the ORC images API.
type Image struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ImageSpec   `json:"spec,omitempty"`
	Status ImageStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ImageList contains a list of Image.
type ImageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Image `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Image{}, &ImageList{})
}

func (i *Image) GetCloudCredentialsRef() (*string, *CloudCredentialsReference) {
	if i == nil {
		return nil, nil
	}

	return &i.Namespace, &i.Spec.CloudCredentialsRef
}

var _ CloudCredentialsRefProvider = &Image{}
