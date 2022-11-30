package baremetal

import (
	"testing"

	"github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestMakeHintMap(t *testing.T) {
	addressableTrue := true
	addressableFalse := false

	for _, tc := range []struct {
		Scenario string
		Hints    RootDeviceHints
		Expected map[string]string
	}{
		{
			Scenario: "device-name",
			Hints: RootDeviceHints{
				DeviceName: "userd_devicename",
			},
			Expected: map[string]string{
				"name": "s== userd_devicename",
			},
		},
		{
			Scenario: "hctl",
			Hints: RootDeviceHints{
				HCTL: "1:2:3:4",
			},
			Expected: map[string]string{
				"hctl": "s== 1:2:3:4",
			},
		},
		{
			Scenario: "model",
			Hints: RootDeviceHints{
				Model: "userd_model",
			},
			Expected: map[string]string{
				"model": "<in> userd_model",
			},
		},
		{
			Scenario: "vendor",
			Hints: RootDeviceHints{
				Vendor: "userd_vendor",
			},
			Expected: map[string]string{
				"vendor": "<in> userd_vendor",
			},
		},
		{
			Scenario: "serial-number",
			Hints: RootDeviceHints{
				SerialNumber: "userd_serial",
			},
			Expected: map[string]string{
				"serial": "s== userd_serial",
			},
		},
		{
			Scenario: "min-size",
			Hints: RootDeviceHints{
				MinSizeGigabytes: 40,
			},
			Expected: map[string]string{
				"size": ">= 40",
			},
		},
		{
			Scenario: "wwn",
			Hints: RootDeviceHints{
				WWN: "userd_wwn",
			},
			Expected: map[string]string{
				"wwn": "s== userd_wwn",
			},
		},
		{
			Scenario: "wwn-with-extension",
			Hints: RootDeviceHints{
				WWNWithExtension: "userd_with_extension",
			},
			Expected: map[string]string{
				"wwn_with_extension": "s== userd_with_extension",
			},
		},
		{
			Scenario: "wwn-extension",
			Hints: RootDeviceHints{
				WWNVendorExtension: "userd_vendor_extension",
			},
			Expected: map[string]string{
				"wwn_vendor_extension": "s== userd_vendor_extension",
			},
		},
		{
			Scenario: "rotational-true",
			Hints: RootDeviceHints{
				Rotational: &addressableTrue,
			},
			Expected: map[string]string{
				"rotational": "true",
			},
		},
		{
			Scenario: "rotational-false",
			Hints: RootDeviceHints{
				Rotational: &addressableFalse,
			},
			Expected: map[string]string{
				"rotational": "false",
			},
		},
		{
			Scenario: "everything-bagel",
			Hints: RootDeviceHints{
				DeviceName:         "userd_devicename",
				HCTL:               "1:2:3:4",
				Model:              "userd_model",
				Vendor:             "userd_vendor",
				SerialNumber:       "userd_serial",
				MinSizeGigabytes:   40,
				WWN:                "userd_wwn",
				WWNWithExtension:   "userd_with_extension",
				WWNVendorExtension: "userd_vendor_extension",
				Rotational:         &addressableTrue,
			},
			Expected: map[string]string{
				"name":                 "s== userd_devicename",
				"hctl":                 "s== 1:2:3:4",
				"model":                "<in> userd_model",
				"vendor":               "<in> userd_vendor",
				"serial":               "s== userd_serial",
				"size":                 ">= 40",
				"wwn":                  "s== userd_wwn",
				"wwn_with_extension":   "s== userd_with_extension",
				"wwn_vendor_extension": "s== userd_vendor_extension",
				"rotational":           "true",
			},
		},
		{
			Scenario: "empty",
			Hints:    RootDeviceHints{},
			Expected: map[string]string{},
		},
	} {
		t.Run(tc.Scenario, func(t *testing.T) {
			actual := tc.Hints.MakeHintMap()
			assert.Equal(t, tc.Expected, actual, "hint map does not match")
		})
	}
}

func TestMakeCRDHints(t *testing.T) {
	addressableTrue := true
	addressableFalse := false

	for _, tc := range []struct {
		Scenario string
		Hints    *RootDeviceHints
		Expected *v1alpha1.RootDeviceHints
	}{
		{
			Scenario: "nil",
			Hints:    nil,
			Expected: nil,
		},
		{
			Scenario: "device-name",
			Hints: &RootDeviceHints{
				DeviceName: "userd_devicename",
			},
			Expected: &v1alpha1.RootDeviceHints{
				DeviceName: "userd_devicename",
			},
		},
		{
			Scenario: "hctl",
			Hints: &RootDeviceHints{
				HCTL: "1:2:3:4",
			},
			Expected: &v1alpha1.RootDeviceHints{
				HCTL: "1:2:3:4",
			},
		},
		{
			Scenario: "model",
			Hints: &RootDeviceHints{
				Model: "userd_model",
			},
			Expected: &v1alpha1.RootDeviceHints{
				Model: "userd_model",
			},
		},
		{
			Scenario: "vendor",
			Hints: &RootDeviceHints{
				Vendor: "userd_vendor",
			},
			Expected: &v1alpha1.RootDeviceHints{
				Vendor: "userd_vendor",
			},
		},
		{
			Scenario: "serial-number",
			Hints: &RootDeviceHints{
				SerialNumber: "userd_serial",
			},
			Expected: &v1alpha1.RootDeviceHints{
				SerialNumber: "userd_serial",
			},
		},
		{
			Scenario: "min-size",
			Hints: &RootDeviceHints{
				MinSizeGigabytes: 40,
			},
			Expected: &v1alpha1.RootDeviceHints{
				MinSizeGigabytes: 40,
			},
		},
		{
			Scenario: "wwn",
			Hints: &RootDeviceHints{
				WWN: "userd_wwn",
			},
			Expected: &v1alpha1.RootDeviceHints{
				WWN: "userd_wwn",
			},
		},
		{
			Scenario: "wwn-with-extension",
			Hints: &RootDeviceHints{
				WWNWithExtension: "userd_with_extension",
			},
			Expected: &v1alpha1.RootDeviceHints{
				WWNWithExtension: "userd_with_extension",
			},
		},
		{
			Scenario: "wwn-extension",
			Hints: &RootDeviceHints{
				WWNVendorExtension: "userd_vendor_extension",
			},
			Expected: &v1alpha1.RootDeviceHints{
				WWNVendorExtension: "userd_vendor_extension",
			},
		},
		{
			Scenario: "rotational-true",
			Hints: &RootDeviceHints{
				Rotational: &addressableTrue,
			},
			Expected: &v1alpha1.RootDeviceHints{
				Rotational: &addressableTrue,
			},
		},
		{
			Scenario: "rotational-false",
			Hints: &RootDeviceHints{
				Rotational: &addressableFalse,
			},
			Expected: &v1alpha1.RootDeviceHints{
				Rotational: &addressableFalse,
			},
		},
		{
			Scenario: "everything-bagel",
			Hints: &RootDeviceHints{
				DeviceName:         "userd_devicename",
				HCTL:               "1:2:3:4",
				Model:              "userd_model",
				Vendor:             "userd_vendor",
				SerialNumber:       "userd_serial",
				MinSizeGigabytes:   40,
				WWN:                "userd_wwn",
				WWNWithExtension:   "userd_with_extension",
				WWNVendorExtension: "userd_vendor_extension",
				Rotational:         &addressableTrue,
			},
			Expected: &v1alpha1.RootDeviceHints{
				DeviceName:         "userd_devicename",
				HCTL:               "1:2:3:4",
				Model:              "userd_model",
				Vendor:             "userd_vendor",
				SerialNumber:       "userd_serial",
				MinSizeGigabytes:   40,
				WWN:                "userd_wwn",
				WWNWithExtension:   "userd_with_extension",
				WWNVendorExtension: "userd_vendor_extension",
				Rotational:         &addressableTrue,
			},
		},
		{
			Scenario: "empty",
			Hints:    &RootDeviceHints{},
			Expected: &v1alpha1.RootDeviceHints{},
		},
	} {
		t.Run(tc.Scenario, func(t *testing.T) {
			actual := tc.Hints.MakeCRDHints()
			assert.Equal(t, tc.Expected, actual, "hint map does not match")
		})
	}
}
