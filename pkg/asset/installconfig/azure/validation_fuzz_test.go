package azure

import (
	"net/url"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/util/validation/field"

	aztypes "github.com/openshift/installer/pkg/types/azure"
)

func FuzzAzureStackClusterOSImage(f *testing.F) {
	f.Add("https://storage.test-endpoint.com/rhcos-image", "storage.test-endpoint.com")
	f.Add("https://storage.not-in-the-cluster.com/rhcos-image", "storage.test-endpoint.com")
	f.Add("https://storage.test-endpoint.com/rhcos-image", "endpoint.com")
	f.Add("", "storage.test-endpoint.com")
	// "%" is not a valid URL escape and used to panic on a nil parse result.
	f.Add("%", "storage.test-endpoint.com")

	f.Fuzz(func(t *testing.T, image, suffix string) {
		errs := validateAzureStackClusterOSImage(suffix, image, field.NewPath("platform").Child("azure"))
		if len(errs) > 1 {
			t.Fatalf("got %d errors, want at most 1: %v", len(errs), errs)
		}
		parsed, err := url.Parse(image)
		if err != nil {
			if len(errs) != 1 || errs[0] == nil || errs[0].Type != field.ErrorTypeInvalid || !strings.Contains(errs[0].Detail, "clusterOSImage URL is invalid") {
				t.Fatalf("invalid URL %q = %v", image, errs)
			}
			return
		}
		if strings.HasSuffix(parsed.Host, suffix) {
			if len(errs) != 0 {
				t.Fatalf("image %q host %q suffix %q = %v, want accept", image, parsed.Host, suffix, errs)
			}
			return
		}
		if len(errs) != 1 || errs[0] == nil || errs[0].Detail != "clusterOSImage must be in the Azure Stack environment" {
			t.Fatalf("image %q host %q suffix %q = %v, want environment error", image, parsed.Host, suffix, errs)
		}
	})
}

func FuzzAzureStackClusterOSImageSet(f *testing.F) {
	f.Add("https://storage.test-endpoint.com/rhcos-image")
	f.Add("")

	f.Fuzz(func(t *testing.T, image string) {
		errs := checkAzureStackClusterOSImageSet(image, field.NewPath("platform").Child("azure"))
		if image == "" {
			if len(errs) != 1 || errs[0] == nil || errs[0].Type != field.ErrorTypeRequired {
				t.Fatalf("empty clusterOSImage = %v, want one Required error", errs)
			}
			return
		}
		if len(errs) != 0 {
			t.Fatalf("clusterOSImage %q = %v, want accept", image, errs)
		}
	})
}

func FuzzPremiumDiskCapability(f *testing.F) {
	f.Add("Premium_LRS", "Standard_D8s_v3", "True", true)
	f.Add("Premium_LRS", "Standard_D4_v4", "False", true)
	f.Add("Premium_LRS", "Standard_D4_v4", "false", true)
	f.Add("Premium_LRS", "Standard_D8s_v3", "", false)

	f.Fuzz(func(t *testing.T, diskType, instanceType, premiumIO string, hasPremiumIO bool) {
		capabilities := map[string]string{}
		if hasPremiumIO {
			capabilities["PremiumIO"] = premiumIO
		}
		errs := validatePremiumDisk(field.NewPath("pool"), diskType, instanceType, capabilities)
		if len(errs) > 1 {
			t.Fatalf("got %d errors, want at most 1: %v", len(errs), errs)
		}
		switch {
		case !hasPremiumIO:
			if len(errs) != 1 || errs[0] == nil || errs[0].Detail != "capability not found: PremiumIO" {
				t.Fatalf("missing PremiumIO = %v", errs)
			}
		case strings.EqualFold(premiumIO, "False"):
			if len(errs) != 1 || errs[0] == nil || !strings.Contains(errs[0].Detail, "PremiumIO not supported for instance type "+instanceType) {
				t.Fatalf("PremiumIO %q instance %q = %v", premiumIO, instanceType, errs)
			}
		default:
			if len(errs) != 0 {
				t.Fatalf("PremiumIO %q = %v, want accept", premiumIO, errs)
			}
		}
	})
}

func FuzzAcceleratedNetworkingCapability(f *testing.F) {
	f.Add("Accelerated", "Standard_D8s_v3", "True", true)
	f.Add("Accelerated", "Standard_A1_v2", "False", true)
	f.Add("Basic", "Standard_Dc4_v4", "", false)

	f.Fuzz(func(t *testing.T, networkingType, instanceType, enabled string, hasCapability bool) {
		capabilities := map[string]string{}
		if hasCapability {
			capabilities[string(aztypes.AcceleratedNetworkingEnabled)] = enabled
		}
		errs := validateAcceleratedNetworking(field.NewPath("pool"), networkingType, instanceType, capabilities)
		if len(errs) > 1 {
			t.Fatalf("got %d errors, want at most 1: %v", len(errs), errs)
		}
		switch {
		case !hasCapability:
			if len(errs) != 1 || errs[0] == nil || errs[0].Detail != "capability not found: AcceleratedNetworkingEnabled" {
				t.Fatalf("missing capability = %v", errs)
			}
		case !strings.EqualFold(enabled, "True"):
			if len(errs) != 1 || errs[0] == nil || !strings.Contains(errs[0].Detail, "vm networking type is not supported for instance type "+instanceType) {
				t.Fatalf("accelerated %q instance %q = %v", enabled, instanceType, errs)
			}
		default:
			if len(errs) != 0 {
				t.Fatalf("accelerated %q = %v, want accept", enabled, errs)
			}
		}
	})
}

func FuzzInstanceSecurityType(f *testing.F) {
	f.Add(string(aztypes.SecurityTypesConfidentialVM), "Standard_DC8ads_v5", "SNP", true, false)
	f.Add(string(aztypes.SecurityTypesConfidentialVM), "Standard_DC8eds_v5", "TDX", true, false)
	f.Add(string(aztypes.SecurityTypesConfidentialVM), "Standard_DC8s_v3", "SGX", true, false)
	f.Add(string(aztypes.SecurityTypesConfidentialVM), "Standard_D8s_v3", "", false, false)
	f.Add(string(aztypes.SecurityTypesTrustedLaunch), "Standard_D8s_v3", "", false, false)
	f.Add(string(aztypes.SecurityTypesTrustedLaunch), "Standard_D8ps_v5", "", false, true)
	f.Add(string(aztypes.SecurityTypesTrustedLaunch), "Standard_DC8ads_v5", "SNP", true, false)
	f.Add("", "Standard_D8s_v3", "", false, false)

	f.Fuzz(func(t *testing.T, securityType, instanceType, confidentialType string, hasConfidential, trustedLaunchDisabled bool) {
		capabilities := map[string]string{}
		if hasConfidential {
			capabilities["ConfidentialComputingType"] = confidentialType
		}
		if trustedLaunchDisabled {
			capabilities["TrustedLaunchDisabled"] = "True"
		}
		errs := validateSecurityType(field.NewPath("settings").Child("securityType"), aztypes.SecurityTypes(securityType), instanceType, capabilities)
		if len(errs) > 1 {
			t.Fatalf("got %d errors, want at most 1: %v", len(errs), errs)
		}
		want := expectedSecurityTypeDetail(securityType, instanceType, confidentialType, hasConfidential, trustedLaunchDisabled)
		if want == "" {
			if len(errs) != 0 {
				t.Fatalf("security type %q instance %q = %v, want accept", securityType, instanceType, errs)
			}
			return
		}
		if len(errs) != 1 || errs[0] == nil || errs[0].Detail != want {
			t.Fatalf("security type %q = %v, want %q", securityType, errs, want)
		}
	})
}

func FuzzWindowsVMFamily(f *testing.F) {
	f.Add("Standard_NV6ads_A10_v5", "standardNVSv4Family")
	f.Add("Standard_D8s_v3", "standardDSv3Family")
	f.Add("Standard_D8s_v3", "")

	f.Fuzz(func(t *testing.T, instanceType, family string) {
		errs := validateFamily(field.NewPath("type"), instanceType, family)
		if family == "standardNVSv4Family" {
			if len(errs) != 1 || errs[0] == nil || errs[0].Detail != "standardNVSv4Family is currently only supported on Windows" {
				t.Fatalf("family %q = %v", family, errs)
			}
			return
		}
		if len(errs) != 0 {
			t.Fatalf("family %q = %v, want accept", family, errs)
		}
	})
}

func expectedSecurityTypeDetail(securityType, instanceType, confidentialType string, hasConfidential, trustedLaunchDisabled bool) string {
	var reason string
	supported := true
	switch aztypes.SecurityTypes(securityType) {
	case aztypes.SecurityTypesConfidentialVM:
		supportedType := confidentialType == confidentialComputingTypeSNP || confidentialType == confidentialComputingTypeTDX
		supported = hasConfidential && supportedType
		if !hasConfidential {
			reason = "no support for Confidential Computing"
		} else if !supportedType {
			reason = "no support for required confidential computing type (found: " + confidentialType + ", supported: " + strings.Join(supportedConfidentialComputingTypes, ", ") + ")"
		}
	case aztypes.SecurityTypesTrustedLaunch:
		supported = !(trustedLaunchDisabled || hasConfidential)
		if trustedLaunchDisabled {
			reason = "no support for Trusted Launch"
		} else if hasConfidential {
			reason = "confidential VMs do not support Trusted Launch for VMs"
		}
	}
	if supported {
		return ""
	}
	return "this security type is not supported for instance type " + instanceType + ", " + reason
}
