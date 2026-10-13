package baremetal

// This package replicates the code from
// https://github.com/metal3-io/baremetal-operator/pkg/provisioner/ironic/devicehints

import (
	"fmt"
	"strings"

	"github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
)

// RootDeviceHints holds the hints for specifying the storage location
// for the root filesystem for the image.
type RootDeviceHints struct {
	// A Linux device name like "/dev/vda". The hint must match the
	// actual value exactly.
	DeviceName string `json:"deviceName,omitempty"`

	// A SCSI bus address like 0:0:0:0. The hint must match the actual
	// value exactly.
	HCTL string `json:"hctl,omitempty"`

	// A vendor-specific device identifier. The hint can be a
	// substring of the actual value.
	Model string `json:"model,omitempty"`

	// The name of the vendor or manufacturer of the device. The hint
	// can be a substring of the actual value.
	Vendor string `json:"vendor,omitempty"`

	// Device serial number. The hint must match the actual value
	// exactly.
	SerialNumber string `json:"serialNumber,omitempty"`

	// The minimum size of the device in Gigabytes.
	// +kubebuilder:validation:Minimum=0
	MinSizeGigabytes int `json:"minSizeGigabytes,omitempty"`

	// Unique storage identifier. The hint must match the actual value
	// exactly.
	WWN string `json:"wwn,omitempty"`

	// Unique storage identifier with the vendor extension
	// appended. The hint must match the actual value exactly.
	WWNWithExtension string `json:"wwnWithExtension,omitempty"`

	// Unique vendor storage identifier. The hint must match the
	// actual value exactly.
	WWNVendorExtension string `json:"wwnVendorExtension,omitempty"`

	// True if the device should use spinning media, false otherwise.
	Rotational *bool `json:"rotational,omitempty"`
}

// MakeHintMap converts a RootDeviceHints instance into a string map
// suitable to pass to ironic.
func (source *RootDeviceHints) MakeHintMap() map[string]string {
	hints := map[string]string{}

	if source == nil {
		return hints
	}

	if source.DeviceName != "" {
		if strings.HasPrefix(source.DeviceName, "/dev/disk/by-path/") {
			hints["by_path"] = fmt.Sprintf("s== %s", source.DeviceName)
		} else {
			hints["name"] = fmt.Sprintf("s== %s", source.DeviceName)
		}
	}
	if source.HCTL != "" {
		hints["hctl"] = fmt.Sprintf("s== %s", source.HCTL)
	}
	if source.Model != "" {
		hints["model"] = fmt.Sprintf("<in> %s", source.Model)
	}
	if source.Vendor != "" {
		hints["vendor"] = fmt.Sprintf("<in> %s", source.Vendor)
	}
	if source.SerialNumber != "" {
		hints["serial"] = fmt.Sprintf("s== %s", source.SerialNumber)
	}
	if source.MinSizeGigabytes != 0 {
		hints["size"] = fmt.Sprintf(">= %d", source.MinSizeGigabytes)
	}
	if source.WWN != "" {
		hints["wwn"] = fmt.Sprintf("s== %s", source.WWN)
	}
	if source.WWNWithExtension != "" {
		hints["wwn_with_extension"] = fmt.Sprintf("s== %s", source.WWNWithExtension)
	}
	if source.WWNVendorExtension != "" {
		hints["wwn_vendor_extension"] = fmt.Sprintf("s== %s", source.WWNVendorExtension)
	}
	switch {
	case source.Rotational == nil:
	case *source.Rotational:
		hints["rotational"] = "true"
	case !*source.Rotational:
		hints["rotational"] = "false"
	}

	return hints
}

// MakeCRDHints returns the hints in the format needed to pass to
// create a BareMetalHost resource.
func (source *RootDeviceHints) MakeCRDHints() *v1alpha1.RootDeviceHints {
	if source == nil {
		return nil
	}
	return &v1alpha1.RootDeviceHints{
		DeviceName:         source.DeviceName,
		HCTL:               source.HCTL,
		Model:              source.Model,
		Vendor:             source.Vendor,
		SerialNumber:       source.SerialNumber,
		MinSizeGigabytes:   source.MinSizeGigabytes,
		WWN:                source.WWN,
		WWNWithExtension:   source.WWNWithExtension,
		WWNVendorExtension: source.WWNVendorExtension,
		Rotational:         source.Rotational,
	}
}
