package azure

import (
	"strings"
	"testing"
)

func FuzzSetBaseDomain(f *testing.F) {
	f.Add("/subscriptions/sub/resourceGroups/rg-name/providers/Microsoft.Network/dnszones/example.com")
	f.Add("")
	f.Add("no-slashes")
	f.Add("/subscriptions/sub")
	f.Add("/subscriptions/sub/resourceGroups/")

	f.Fuzz(func(t *testing.T, zoneID string) {
		platform := &Platform{BaseDomainResourceGroupName: "unchanged"}
		err := platform.SetBaseDomain(zoneID)
		parts := strings.Split(zoneID, "/")
		if len(parts) < 5 {
			if err == nil || platform.BaseDomainResourceGroupName != "unchanged" {
				t.Fatalf("SetBaseDomain(%q) = %v, group %q", zoneID, err, platform.BaseDomainResourceGroupName)
			}
			return
		}
		if err != nil || platform.BaseDomainResourceGroupName != parts[4] {
			t.Fatalf("SetBaseDomain(%q) = %v, group %q, want %q", zoneID, err, platform.BaseDomainResourceGroupName, parts[4])
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
