// Copyright (c) 2023 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha2

import (
	"github.com/vmware-tanzu/vm-operator/api/v1alpha2/cloudinit"
	"github.com/vmware-tanzu/vm-operator/api/v1alpha2/common"
	"github.com/vmware-tanzu/vm-operator/api/v1alpha2/sysprep"
)

// VirtualMachineBootstrapSpec defines the desired state of a VM's bootstrap
// configuration.
type VirtualMachineBootstrapSpec struct {

	// CloudInit may be used to bootstrap Linux guests with Cloud-Init or
	// Windows guests that support Cloudbase-Init.
	//
	// The guest's networking stack is configured by Cloud-Init on Linux guests
	// and Cloudbase-Init on Windows guests.
	//
	// Please note this bootstrap provider may not be used in conjunction with
	// the other bootstrap providers.
	//
	// +optional
	CloudInit *VirtualMachineBootstrapCloudInitSpec `json:"cloudInit,omitempty"`

	// LinuxPrep may be used to bootstrap Linux guests.
	//
	// The guest's networking stack is configured by Guest OS Customization
	// (GOSC).
	//
	// Please note this bootstrap provider may be used in conjunction with the
	// VAppConfig bootstrap provider when wanting to configure the guest's
	// network with GOSC but also send vApp/OVF properties into the guest.
	//
	// This bootstrap provider may not be used in conjunction with the CloudInit
	// or Sysprep bootstrap providers.
	//
	// +optional
	LinuxPrep *VirtualMachineBootstrapLinuxPrepSpec `json:"linuxPrep,omitempty"`

	// Sysprep may be used to bootstrap Windows guests.
	//
	// The guest's networking stack is configured by Guest OS Customization
	// (GOSC).
	//
	// Please note this bootstrap provider may be used in conjunction with the
	// VAppConfig bootstrap provider when wanting to configure the guest's
	// network with GOSC but also send vApp/OVF properties into the guest.
	//
	// This bootstrap provider may not be used in conjunction with the CloudInit
	// or LinuxPrep bootstrap providers.
	//
	// +optional
	Sysprep *VirtualMachineBootstrapSysprepSpec `json:"sysprep,omitempty"`

	// VAppConfig may be used to bootstrap guests that rely on vApp properties
	// (how VMware surfaces OVF properties on guests) to transport data into the
	// guest.
	//
	// The guest's networking stack may be configured using either vApp
	// properties or GOSC.
	//
	// Many OVFs define one or more properties that are used by the guest to
	// bootstrap its networking stack. If the VirtualMachineImage defines one or
	// more properties like this, then they can be configured to use the network
	// data provided for this VM at runtime by setting these properties to Go
	// template strings.
	//
	// It is also possible to use GOSC to bootstrap this VM's network stack by
	// configuring either the LinuxPrep or Sysprep bootstrap providers.
	//
	// Please note the VAppConfig bootstrap provider in conjunction with the
	// LinuxPrep bootstrap provider is the equivalent of setting the v1alpha1
	// VM metadata transport to "OvfEnv".
	//
	// This bootstrap provider may not be used in conjunction with the CloudInit
	// bootstrap provider.
	//
	// +optional
	VAppConfig *VirtualMachineBootstrapVAppConfigSpec `json:"vAppConfig,omitempty"`
}

// VirtualMachineBootstrapCloudInitSpec describes the CloudInit configuration
// used to bootstrap the VM.
type VirtualMachineBootstrapCloudInitSpec struct {
	// CloudConfig describes a subset of a Cloud-Init CloudConfig, used to
	// bootstrap the VM.
	//
	// Please note this field and RawCloudConfig are mutually exclusive.
	//
	// +optional
	CloudConfig *cloudinit.CloudConfig `json:"cloudConfig,omitempty"`

	// RawCloudConfig describes a key in a Secret resource that contains the
	// CloudConfig data used to bootstrap the VM.
	//
	// The CloudConfig data specified by the key may be plain-text,
	// base64-encoded, or gzipped and base64-encoded.
	//
	// Please note this field and CloudConfig are mutually exclusive.
	//
	// +optional
	RawCloudConfig *common.SecretKeySelector `json:"rawCloudConfig,omitempty"`

	// SSHAuthorizedKeys is a list of public keys that CloudInit will apply to
	// the guest's default user.
	//
	// +optional
	SSHAuthorizedKeys []string `json:"sshAuthorizedKeys,omitempty"`
}

// VirtualMachineBootstrapLinuxPrepSpec describes the LinuxPrep configuration
// used to bootstrap the VM.
type VirtualMachineBootstrapLinuxPrepSpec struct {
	// HardwareClockIsUTC specifies whether the hardware clock is in UTC or
	// local time.
	//
	// +optional
	HardwareClockIsUTC bool `json:"hardwareClockIsUTC,omitempty"`

	// TimeZone is a case-sensitive timezone, such as Europe/Sofia.
	//
	// Valid values are based on the tz (timezone) database used by Linux and
	// other Unix systems. The values are strings in the form of
	// "Area/Location," in which Area is a continent or ocean name, and
	// Location is the city, island, or other regional designation.
	//
	// Please see https://kb.vmware.com/s/article/2145518 for a list of valid
	// time zones for Linux systems.
	//
	// +optional
	TimeZone string `json:"timeZone,omitempty"`
}

// VirtualMachineBootstrapSysprepSpec describes the Sysprep configuration used
// to bootstrap the VM.
type VirtualMachineBootstrapSysprepSpec struct {
	// Sysprep is an object representation of a Windows sysprep.xml answer file.
	//
	// This field encloses all the individual keys listed in a sysprep.xml file.
	//
	// For more detailed information please see
	// https://technet.microsoft.com/en-us/library/cc771830(v=ws.10).aspx.
	//
	// Please note this field and RawSysprep are mutually exclusive.
	//
	// +optional
	Sysprep *sysprep.Sysprep `json:"sysprep,omitempty"`

	// RawSysprep describes a key in a Secret resource that contains an XML
	// string of the Sysprep text used to bootstrap the VM.
	//
	// The data specified by the Secret key may be plain-text, base64-encoded,
	// or gzipped and base64-encoded.
	//
	// Please note this field and Sysprep are mutually exclusive.
	//
	// +optional
	RawSysprep *common.SecretKeySelector `json:"rawSysprep,omitempty"`
}

// VirtualMachineBootstrapVAppConfigSpec describes the vApp configuration
// used to bootstrap the VM.
type VirtualMachineBootstrapVAppConfigSpec struct {
	// Properties is a list of vApp/OVF property key/value pairs.
	//
	// Please note this field and RawProperties are mutually exclusive.
	//
	// +optional
	// +listType=map
	// +listMapKey=key
	Properties []common.KeyValueOrSecretKeySelectorPair `json:"properties,omitempty"`

	// RawProperties is the name of a Secret resource in the same Namespace as
	// this VM where each key/value pair from the Secret is used as a vApp
	// key/value pair.
	//
	// Please note this field and Properties are mutually exclusive.
	//
	// +optional
	RawProperties string `json:"rawProperties,omitempty"`
}
