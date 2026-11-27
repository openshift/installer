package validation

import (
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/azure"
)

const (
	diskEncryptionUnsupportedMsg = "disk encryption sets are not supported on this platform"
	subscriptionIDInvalidMsg     = "invalid subscription ID format"
	resourceGroupInvalidMsg      = "invalid resource group format"
	diskEncryptionNameInvalidMsg = "invalid name format"
	encryptionAtHostUnsupported  = "encryption at host is not supported on this platform"
)

func FuzzValidateDiskEncryption(f *testing.F) {
	f.Add(subscriptionID, resourceGroup, diskEncryptionSetName, string(azure.PublicCloud))
	f.Add(subscriptionID, resourceGroup, diskEncryptionSetName, string(azure.StackCloud))
	f.Add("invalid", resourceGroup, diskEncryptionSetName, string(azure.PublicCloud))
	f.Add(subscriptionID, "", diskEncryptionSetName, string(azure.PublicCloud))
	f.Add(subscriptionID, resourceGroup, "", string(azure.PublicCloud))
	f.Add(strings.ToUpper(subscriptionID), strings.Repeat("a", 90), strings.Repeat("b", 80), string(azure.USGovernmentCloud))
	f.Add(subscriptionID, strings.Repeat("a", 91), strings.Repeat("b", 81), "")
	f.Add(subscriptionID, "group.", "name.", string(azure.ChinaCloud))
	f.Add("", "", "", string(azure.StackCloud))

	f.Fuzz(func(t *testing.T, sub, rg, name, cloud string) {
		pool := &azure.MachinePool{
			OSDisk: azure.OSDisk{
				DiskEncryptionSet: &azure.DiskEncryptionSet{
					SubscriptionID: sub,
					ResourceGroup:  rg,
					Name:           name,
				},
			},
		}
		errs := ValidateDiskEncryption(pool, azure.CloudEnvironment(cloud), field.NewPath("test-path"))
		if len(errs) > 1 {
			t.Fatalf("got %d errors, want at most 1: %v", len(errs), errs)
		}

		want := expectedDiskEncryptionError(sub, rg, name, cloud)
		if want == "" {
			if len(errs) != 0 {
				t.Fatalf("ValidateDiskEncryption(%q, %q, %q, %q) = %v, want accept", sub, rg, name, cloud, errs)
			}
			return
		}
		if len(errs) != 1 {
			t.Fatalf("ValidateDiskEncryption(%q, %q, %q, %q) accepted, want %q", sub, rg, name, cloud, want)
		}
		if errs[0] == nil {
			t.Fatal("ValidateDiskEncryption returned a nil error")
		}
		if errs[0].Type != field.ErrorTypeInvalid || errs[0].Detail != want {
			t.Fatalf("error = type %q detail %q, want Invalid %q", errs[0].Type, errs[0].Detail, want)
		}
	})
}

func FuzzDiskEncryptionSetID(f *testing.F) {
	f.Add(subscriptionID, resourceGroup, diskEncryptionSetName)
	f.Add("invalid", resourceGroup, diskEncryptionSetName)
	f.Add(subscriptionID, "", diskEncryptionSetName)
	f.Add(subscriptionID, resourceGroup, "")
	f.Add(strings.Repeat("a", 8)+"-bbbb-cccc-dddd-eeeeeeeeeeee", strings.Repeat("r", 90), strings.Repeat("n", 80))
	f.Add(subscriptionID, "a.", "-name")
	f.Add(subscriptionID, "group/extra", "name/extra")

	f.Fuzz(func(t *testing.T, sub, rg, name string) {
		if RxSubscriptionID.MatchString(sub) != subscriptionIDAccepted(sub) {
			t.Fatalf("RxSubscriptionID(%q) disagrees with subscription ID rules", sub)
		}
		if RxResourceGroup.MatchString(rg) != resourceGroupAccepted(rg) {
			t.Fatalf("RxResourceGroup(%q) disagrees with resource group rules", rg)
		}
		if RxDiskEncryptionSetName.MatchString(name) != diskEncryptionSetNameAccepted(name) {
			t.Fatalf("RxDiskEncryptionSetName(%q) disagrees with disk encryption set name rules", name)
		}

		id := (&azure.DiskEncryptionSet{SubscriptionID: sub, ResourceGroup: rg, Name: name}).ToID()
		partsOK := subscriptionIDAccepted(sub) && resourceGroupAccepted(rg) && diskEncryptionSetNameAccepted(name)
		if matched := RxDiskEncryptionSetID.MatchString(id); matched != partsOK {
			t.Fatalf("ToID %q match=%t, parts accepted=%t", id, matched, partsOK)
		}
	})
}

func FuzzValidateEncryptionAtHost(f *testing.F) {
	f.Add(true, string(azure.StackCloud))
	f.Add(true, string(azure.PublicCloud))
	f.Add(false, string(azure.StackCloud))
	f.Add(false, string(azure.PublicCloud))
	f.Add(true, "")
	f.Add(true, string(azure.USGovernmentCloud))

	f.Fuzz(func(t *testing.T, enabled bool, cloud string) {
		pool := &azure.MachinePool{EncryptionAtHost: enabled}
		errs := ValidateEncryptionAtHost(pool, azure.CloudEnvironment(cloud), field.NewPath("test-path"))
		wantUnsupported := enabled && cloud == string(azure.StackCloud)
		if wantUnsupported {
			if len(errs) != 1 || errs[0] == nil || errs[0].Type != field.ErrorTypeInvalid || errs[0].Detail != encryptionAtHostUnsupported {
				t.Fatalf("encryption at host on %q = %v, want one Invalid error", cloud, errs)
			}
			return
		}
		if len(errs) != 0 {
			t.Fatalf("encryption at host enabled=%t cloud=%q = %v, want accept", enabled, cloud, errs)
		}
	})
}

func expectedDiskEncryptionError(sub, rg, name, cloud string) string {
	if cloud == string(azure.StackCloud) {
		return diskEncryptionUnsupportedMsg
	}
	if !subscriptionIDAccepted(sub) {
		return subscriptionIDInvalidMsg
	}
	if !resourceGroupAccepted(rg) {
		return resourceGroupInvalidMsg
	}
	if !diskEncryptionSetNameAccepted(name) {
		return diskEncryptionNameInvalidMsg
	}
	return ""
}

func subscriptionIDAccepted(s string) bool {
	// 8-4-4-4-12 hexadecimal characters.
	if len(s) != 36 {
		return false
	}
	widths := []int{8, 4, 4, 4, 12}
	i := 0
	for group, width := range widths {
		if group > 0 {
			if s[i] != '-' {
				return false
			}
			i++
		}
		for k := 0; k < width; k++ {
			if !isHex(s[i]) {
				return false
			}
			i++
		}
	}
	return i == len(s)
}

func resourceGroupAccepted(s string) bool {
	if len(s) < 1 || len(s) > 90 {
		return false
	}
	for i := 0; i < len(s); i++ {
		if !isResourceGroupChar(s[i]) {
			return false
		}
	}
	return isResourceGroupTail(s[len(s)-1])
}

func diskEncryptionSetNameAccepted(s string) bool {
	if len(s) < 1 || len(s) > 80 {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !isASCIILetter(c) && !isDigit(c) && c != '_' && c != '-' {
			return false
		}
	}
	return true
}

func isHex(c byte) bool {
	return isDigit(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func isResourceGroupTail(c byte) bool {
	return isASCIILetter(c) || isDigit(c) || c == '_' || c == '(' || c == ')' || c == '-'
}

func isResourceGroupChar(c byte) bool {
	return isResourceGroupTail(c) || c == '.'
}
