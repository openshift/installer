package azure

import (
	"strings"
	"testing"
)

func FuzzSetBaseDomain(f *testing.F) {
	f.Add("rg-name")
	f.Add("")
	f.Add("rg/name")

	f.Fuzz(func(t *testing.T, resourceGroup string) {
		// Arbitrary malformed IDs must not panic.
		malformedPlatform := &Platform{BaseDomainResourceGroupName: "unchanged"}
		if err := malformedPlatform.SetBaseDomain(resourceGroup); err != nil && malformedPlatform.BaseDomainResourceGroupName != "unchanged" {
			t.Fatalf("SetBaseDomain(%q) changed the platform on error: %q", resourceGroup, malformedPlatform.BaseDomainResourceGroupName)
		}

		component := strings.ReplaceAll(resourceGroup, "/", "_")
		zoneID := "/subscriptions/sub/resourceGroups/" + component + "/providers/Microsoft.Network/dnszones/example.com"
		platform := &Platform{BaseDomainResourceGroupName: "unchanged"}
		err := platform.SetBaseDomain(zoneID)
		if component == "" {
			if err == nil || platform.BaseDomainResourceGroupName != "unchanged" {
				t.Fatalf("SetBaseDomain(%q) = %v, group %q", zoneID, err, platform.BaseDomainResourceGroupName)
			}
			return
		}
		if err != nil || platform.BaseDomainResourceGroupName != component {
			t.Fatalf("SetBaseDomain(%q) = %v, group %q, want %q", zoneID, err, platform.BaseDomainResourceGroupName, component)
		}
		for _, malformedID := range []string{"", "no-slashes", "/subscriptions/sub", "/subscriptions/sub/resourceGroups/"} {
			invalidPlatform := &Platform{BaseDomainResourceGroupName: "unchanged"}
			if err := invalidPlatform.SetBaseDomain(malformedID); err == nil || invalidPlatform.BaseDomainResourceGroupName != "unchanged" {
				t.Fatalf("SetBaseDomain(%q) = %v, group %q, want error without mutation", malformedID, err, invalidPlatform.BaseDomainResourceGroupName)
			}
		}
	})
}

func FuzzAzureResourceNames(f *testing.F) {
	f.Add("infra-id", "", "")
	f.Add("infra-id", "custom-rg", "custom-vnet")
	f.Add("", "", "")

	f.Fuzz(func(t *testing.T, infraID, resourceGroup, virtualNetwork string) {
		platform := &Platform{ResourceGroupName: resourceGroup, VirtualNetwork: virtualNetwork}
		if got := platform.ClusterResourceGroupName(infraID); resourceGroup != "" && got != resourceGroup || resourceGroup == "" && got != infraID+"-rg" {
			t.Fatalf("resource group %q from %q/%q", got, resourceGroup, infraID)
		}
		if got := platform.VirtualNetworkName(infraID); virtualNetwork != "" && got != virtualNetwork || virtualNetwork == "" && got != infraID+"-vnet" {
			t.Fatalf("virtual network %q from %q/%q", got, virtualNetwork, infraID)
		}
		if got := platform.ControlPlaneSubnetName(infraID); got != infraID+"-master-subnet" {
			t.Fatalf("control plane subnet %q", got)
		}
		if got := platform.ComputeSubnetName(infraID); got != infraID+"-worker-subnet" {
			t.Fatalf("compute subnet %q", got)
		}
		if got := platform.NetworkSecurityGroupName(infraID); got != infraID+"-nsg" {
			t.Fatalf("network security group %q", got)
		}
	})
}

func FuzzStorageAccountName(f *testing.F) {
	f.Add("infra-id")
	f.Add("")
	f.Add(strings.Repeat("a", 30))
	f.Add(strings.Repeat("ab-", 20))

	f.Fuzz(func(t *testing.T, infraID string) {
		got := GetStorageAccountName(infraID)
		if len(got) > 24 {
			t.Fatalf("storage account %q length %d", got, len(got))
		}
		if strings.Contains(got, "-") || !strings.HasSuffix(got, "sa") {
			t.Fatalf("storage account %q", got)
		}
		compact := strings.ReplaceAll(infraID, "-", "")
		if len(compact) > 22 {
			compact = compact[:22]
		}
		if got != compact+"sa" {
			t.Fatalf("GetStorageAccountName(%q) = %q, want %q", infraID, got, compact+"sa")
		}
	})
}
