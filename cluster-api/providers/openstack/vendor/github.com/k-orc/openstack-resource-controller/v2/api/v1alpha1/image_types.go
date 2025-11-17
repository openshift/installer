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

// GlanceTag is the name of the go field tag in properties structs used to specify the Glance property name.
const GlanceTag = "glance"

// +kubebuilder:validation:MinLength:=1
// +kubebuilder:validation:MaxLength:=255
type ImageTag string

// +kubebuilder:validation:Enum:=ami;ari;aki;bare;ovf;ova;docker;compressed
type ImageContainerFormat string

const (
	ImageContainerFormatAKI        ImageContainerFormat = "aki"
	ImageContainerFormatAMI        ImageContainerFormat = "ami"
	ImageContainerFormatARI        ImageContainerFormat = "ari"
	ImageContainerFormatBare       ImageContainerFormat = "bare"
	ImageContainerFormatCompressed ImageContainerFormat = "compressed"
	ImageContainerFormatDocker     ImageContainerFormat = "docker"
	ImageContainerFormatOVA        ImageContainerFormat = "ova"
	ImageContainerFormatOVF        ImageContainerFormat = "ovf"
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
	// cpuSockets is the preferred number of sockets to expose to the guest
	// +kubebuilder:validation:Minimum:=1
	// +optional
	CPUSockets *int32 `json:"cpuSockets,omitempty" glance:"hw_cpu_sockets"`

	// cpuCores is the preferred number of cores to expose to the guest
	// +kubebuilder:validation:Minimum:=1
	// +optional
	CPUCores *int32 `json:"cpuCores,omitempty" glance:"hw_cpu_cores"`

	// cpuThreads is the preferred number of threads to expose to the guest
	// +kubebuilder:validation:Minimum:=1
	// +optional
	CPUThreads *int32 `json:"cpuThreads,omitempty" glance:"hw_cpu_threads"`

	// cpuPolicy is used to pin the virtual CPUs (vCPUs) of instances to the
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

	// cpuThreadPolicy further refines a CPUPolicy of 'dedicated' by stating
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

	// cdromBus specifies the type of disk controller to attach CD-ROM devices to.
	// +optional
	CDROMBus *ImageHWBus `json:"cdromBus,omitempty" glance:"hw_cdrom_bus"`

	// diskBus specifies the type of disk controller to attach disk devices to.
	// +optional
	DiskBus *ImageHWBus `json:"diskBus,omitempty" glance:"hw_disk_bus"`

	// TODO: hw_machine_type seems important to support early, but how to
	// select a supported set?

	// scsiModel enables the use of VirtIO SCSI (virtio-scsi) to provide
	// block device access for compute instances; by default, instances use
	// VirtIO Block (virtio-blk). VirtIO SCSI is a para-virtualized SCSI
	// controller device that provides improved scalability and performance,
	// and supports advanced SCSI hardware.
	//
	// The only permitted value is virtio-scsi.
	// +kubebuilder:validation:Enum:=virtio-scsi
	// +optional
	SCSIModel *string `json:"scsiModel,omitempty" glance:"hw_scsi_model"`

	// vifModel specifies the model of virtual network interface device to use.
	//
	// Permitted values are e1000, e1000e, ne2k_pci, pcnet, rtl8139, virtio,
	// and vmxnet3.
	// +kubebuilder:validation:Enum:=e1000;e1000e;ne2k_pci;pcnet;rtl8139;virtio;vmxnet3
	// +optional
	VIFModel *string `json:"vifModel,omitempty" glance:"hw_vif_model"`

	// rngModel adds a random-number generator device to the image’s instances.
	// This image property by itself does not guarantee that a hardware RNG will be used;
	// it expresses a preference that may or may not be satisfied depending upon Nova configuration.
	// +kubebuilder:validation:MaxLength:=255
	// +optional
	RngModel *string `json:"rngModel,omitempty" glance:"hw_rng_model"`

	// qemuGuestAgent enables QEMU guest agent.
	// +optional
	QemuGuestAgent *bool `json:"qemuGuestAgent,omitempty" glance:"hw_qemu_guest_agent"`
}

type ImagePropertiesOperatingSystem struct {
	// distro is the common name of the operating system distribution in lowercase.
	// +kubebuilder:validation:Enum:=arch;centos;debian;fedora;freebsd;gentoo;mandrake;mandriva;mes;msdos;netbsd;netware;openbsd;opensolaris;opensuse;rocky;rhel;sled;ubuntu;windows
	// +optional
	Distro *string `json:"distro,omitempty" glance:"os_distro"`
	// version is the operating system version as specified by the distributor.
	// +kubebuilder:validation:MaxLength:=255
	// +optional
	Version *string `json:"version,omitempty" glance:"os_version"`
}

type ImageProperties struct {
	// architecture is the CPU architecture that must be supported by the hypervisor.
	// +kubebuilder:validation:Enum:=aarch64;alpha;armv7l;cris;i686;ia64;lm32;m68k;microblaze;microblazeel;mips;mipsel;mips64;mips64el;openrisc;parisc;parisc64;ppc;ppc64;ppcemb;s390;s390x;sh4;sh4eb;sparc;sparc64;unicore32;x86_64;xtensa;xtensaeb
	// +optional
	Architecture *string `json:"architecture,omitempty" glance:"architecture"`

	// hypervisorType is the hypervisor type
	// +kubebuilder:validation:Enum:=hyperv;ironic;lxc;qemu;uml;vmware;xen
	// +optional
	HypervisorType *string `json:"hypervisorType,omitempty" glance:"hypervisor_type"`

	// minDiskGB is the minimum amount of disk space in GB that is required to boot the image
	// +kubebuilder:validation:Minimum:=1
	// +optional
	MinDiskGB *int32 `json:"minDiskGB,omitempty"`

	// minMemoryMB is the minimum amount of RAM in MB that is required to boot the image.
	// +kubebuilder:validation:Minimum:=1
	// +optional
	MinMemoryMB *int32 `json:"minMemoryMB,omitempty"`

	// hardware is a set of properties which control the virtual hardware
	// created by Nova.
	// +optional
	Hardware *ImagePropertiesHardware `json:"hardware,omitempty"`

	// operatingSystem is a set of properties that specify and influence the behavior
	// of the operating system within the virtual machine.
	// +optional
	OperatingSystem *ImagePropertiesOperatingSystem `json:"operatingSystem,omitempty"`
}

// +kubebuilder:validation:Enum:=xz;gz;bz2
type ImageCompression string

const (
	ImageCompressionXZ  ImageCompression = "xz"
	ImageCompressionGZ  ImageCompression = "gz"
	ImageCompressionBZ2 ImageCompression = "bz2"
)

type ImageContent struct {
	// containerFormat is the format of the image container.
	// qcow2 and raw images do not usually have a container. This is specified as "bare", which is also the default.
	// Permitted values are ami, ari, aki, bare, compressed, ovf, ova, and docker.
	// +kubebuilder:default:=bare
	// +optional
	ContainerFormat ImageContainerFormat `json:"containerFormat,omitempty"`

	// diskFormat is the format of the disk image.
	// Normal values are "qcow2", or "raw". Glance may be configured to support others.
	// +required
	DiskFormat ImageDiskFormat `json:"diskFormat,omitempty"`

	// download describes how to obtain image data by downloading it from a URL.
	// Must be set when creating a managed image.
	// +required
	//nolint:kubeapilinter
	Download *ImageContentSourceDownload `json:"download"`
}

type ImageContentSourceDownload struct {
	// url containing image data
	// +kubebuilder:validation:Format=uri
	// +kubebuilder:validation:MaxLength=2048
	// +required
	URL string `json:"url"`

	// decompress specifies that the source data must be decompressed with the
	// given compression algorithm before being stored. Specifying Decompress
	// will disable the use of Glance's web-download, as web-download cannot
	// currently deterministically decompress downloaded content.
	// +optional
	Decompress *ImageCompression `json:"decompress,omitempty"`

	// hash is a hash which will be used to verify downloaded data, i.e.
	// before any decompression. If not specified, no hash verification will be
	// performed. Specifying a Hash will disable the use of Glance's
	// web-download, as web-download cannot currently deterministically verify
	// the hash of downloaded content.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="hash is immutable"
	// +optional
	Hash *ImageHash `json:"hash,omitempty"`
}

type ImageHash struct {
	// algorithm is the hash algorithm used to generate value.
	// +required
	Algorithm ImageHashAlgorithm `json:"algorithm,omitempty"`

	// value is the hash of the image data using Algorithm. It must be hex encoded using lowercase letters.
	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:MaxLength:=1024
	// +kubebuilder:validation:Pattern:=`^[0-9a-f]+$`
	// +required
	Value string `json:"value,omitempty"`
}

// ImageResourceSpec contains the desired state of a Glance image
type ImageResourceSpec struct {
	// name will be the name of the created Glance image. If not specified, the
	// name of the Image object will be used.
	// +optional
	Name *OpenStackName `json:"name,omitempty"`

	// protected specifies that the image is protected from deletion.
	// If not specified, the default is false.
	// +optional
	Protected *bool `json:"protected,omitempty"`

	// tags is a list of tags which will be applied to the image. A tag has a maximum length of 255 characters.
	// +kubebuilder:validation:MaxItems:=64
	// +listType=set
	// +optional
	Tags []ImageTag `json:"tags,omitempty"`

	// visibility of the image
	// +optional
	Visibility *ImageVisibility `json:"visibility,omitempty"`

	// properties is metadata available to consumers of the image
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="properties is immutable"
	// +optional
	Properties *ImageProperties `json:"properties,omitempty"`

	// content specifies how to obtain the image content.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="content is immutable"
	// +optional
	Content *ImageContent `json:"content,omitempty"`
}

// ImageFilter defines a Glance query
// +kubebuilder:validation:MinProperties:=1
type ImageFilter struct {
	// name specifies the name of a Glance image
	// +optional
	Name *OpenStackName `json:"name,omitempty"`

	// visibility specifies the visibility of a Glance image.
	// +optional
	Visibility *ImageVisibility `json:"visibility,omitempty"`

	// tags is the list of tags on the resource.
	// +kubebuilder:validation:MaxItems:=64
	// +listType=set
	// +optional
	Tags []ImageTag `json:"tags,omitempty"`
}

// ImageResourceStatus represents the observed state of a Glance image
type ImageResourceStatus struct {
	// name is a Human-readable name for the image. Might not be unique.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Name string `json:"name,omitempty"`

	// status is the image status as reported by Glance
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Status string `json:"status,omitempty"`

	// protected specifies that the image is protected from deletion.
	// +optional
	Protected bool `json:"protected,omitempty"`

	// visibility of the image
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Visibility string `json:"visibility,omitempty"`

	// hash is the hash of the image data published by Glance. Note that this is
	// a hash of the data stored internally by Glance, which will have been
	// decompressed and potentially format converted depending on server-side
	// configuration which is not visible to clients. It is expected that this
	// hash will usually differ from the download hash.
	// +optional
	Hash *ImageHash `json:"hash,omitempty"`

	// sizeB is the size of the image data, in bytes
	// +optional
	SizeB *int64 `json:"sizeB,omitempty"`

	// virtualSizeB is the size of the disk the image data represents, in bytes
	// +optional
	VirtualSizeB *int64 `json:"virtualSizeB,omitempty"`

	// tags is the list of tags on the resource.
	// +kubebuilder:validation:MaxItems:=64
	// +kubebuilder:validation:items:MaxLength=1024
	// +listType=atomic
	// +optional
	Tags []string `json:"tags,omitempty"`
}

type ImageStatusExtra struct {
	// downloadAttempts is the number of times the controller has attempted to download the image contents
	// +optional
	DownloadAttempts *int32 `json:"downloadAttempts,omitempty"`
}
