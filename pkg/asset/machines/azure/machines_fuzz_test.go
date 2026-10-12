package azure

import (
	"fmt"
	"strings"
	"testing"

	machineapi "github.com/openshift/api/machine/v1beta1"
	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"

	aztypes "github.com/openshift/installer/pkg/types/azure"
)

func FuzzTrimSubscriptionPrefix(f *testing.F) {
	f.Add("/subscriptions/12345678-1234-1234-1234-123456789abc/resourceGroups/my-rg/providers/Microsoft.Compute/images/my-image")
	f.Add("/resourceGroups/my-rg/providers/Microsoft.Compute/images/my-image")
	f.Add("/subscriptions/12345678-1234-1234-1234-123456789abc/providers/Microsoft.Compute/images/my-image")
	f.Add("")
	f.Add("RedHat:RHEL:8-lvm:latest")

	f.Fuzz(func(t *testing.T, image string) {
		got := trimSubscriptionPrefix(image)
		marker := "/resourceGroups/"
		index := strings.Index(image, marker)
		want := image
		if index >= 0 {
			want = image[index:]
		}
		if got != want {
			t.Fatalf("trimSubscriptionPrefix(%q) = %q, want %q", image, got, want)
		}
		if len(got) > len(image) {
			t.Fatalf("trimmed image %q is longer than %q", got, image)
		}
	})
}

func FuzzMachineImage(f *testing.F) {
	f.Add("RedHat", "RHEL", "8-lvm-gen2", "latest", string(aztypes.ImageWithPurchasePlan), string(aztypes.PublicCloud), false, "V2", "test-rg", "sub", "infra-id", "")
	f.Add("", "", "", "", "", string(aztypes.StackCloud), false, "V1", "rg", "sub", "infra", "")
	f.Add("", "", "", "", "", string(aztypes.PublicCloud), false, "V2", "rg", "sub", "infra-id", "RedHat:RHEL:8-lvm:latest")
	f.Add("", "", "", "", "", string(aztypes.PublicCloud), false, "V1", "rg", "sub", "infra", "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Compute/images/img")
	f.Add("", "", "", "", "", string(aztypes.PublicCloud), false, "V2", "rg", "sub", "infra", "")
	f.Add("", "", "", "", "", string(aztypes.PublicCloud), true, "V2", "rg", "sub", "infra-id", "")

	f.Fuzz(func(t *testing.T, publisher, offer, sku, version, plan, cloud string, confidential bool, gen, rg, sub, infraID, rhcos string) {
		osImage := aztypes.OSImage{
			Publisher: publisher,
			Offer:     offer,
			SKU:       sku,
			Version:   version,
			Plan:      aztypes.ImagePurchasePlan(plan),
		}
		image := capzImage(osImage, aztypes.CloudEnvironment(cloud), confidential, gen, rg, sub, infraID, rhcos)
		if image == nil {
			t.Fatal("capzImage returned nil")
		}
		assertCapzImage(t, image, osImage, cloud, confidential, gen, rg, sub, infraID, rhcos)

		mapi := mapiImage(osImage, aztypes.CloudEnvironment(cloud), confidential, gen, rg, sub, infraID, rhcos)
		if image.ID != nil && mapi.ResourceID != trimSubscriptionPrefix(*image.ID) {
			t.Fatalf("mapi resource ID %q, want trimmed %q", mapi.ResourceID, *image.ID)
		}
		if image.ID == nil && image.Marketplace != nil {
			if mapi.Publisher != image.Marketplace.Publisher || mapi.Offer != image.Marketplace.Offer || mapi.SKU != image.Marketplace.SKU || mapi.Version != image.Marketplace.Version {
				t.Fatalf("mapi marketplace = %+v, capz = %+v", mapi, image.Marketplace)
			}
		}
	})
}

func FuzzCapzUserTags(f *testing.F) {
	f.Add("cluster", "key", "value", "other", "value2")
	f.Add("cluster", "kubernetes.io_cluster.cluster", "owned", "key", "value")
	f.Add("", "", "", "", "")

	f.Fuzz(func(t *testing.T, clusterID, k1, v1, k2, v2 string) {
		clusterKey := fmt.Sprintf("kubernetes.io_cluster.%s", clusterID)
		owned, err := CapzTagsFromUserTags(clusterID, nil)
		if err != nil || owned[clusterKey] != "owned" || len(owned) != 1 {
			t.Fatalf("nil user tags cluster %q = %v, %v", clusterID, owned, err)
		}

		userTags := map[string]string{k1: v1, k2: v2}
		tags, err := CapzTagsFromUserTags(clusterID, userTags)
		blocked := false
		for key := range userTags {
			if key == clusterKey {
				blocked = true
			}
		}
		if blocked {
			if err == nil || tags != nil {
				t.Fatalf("cluster key %q in user tags = %v, %v", clusterKey, tags, err)
			}
			return
		}
		if err != nil {
			t.Fatalf("user tags %v: %v", userTags, err)
		}
		if tags[clusterKey] != "owned" {
			t.Fatalf("missing owned tag %q in %v", clusterKey, tags)
		}
		if len(tags) != 1+len(userTags) {
			t.Fatalf("tag count %d, want %d (%v)", len(tags), 1+len(userTags), tags)
		}
		for key, value := range userTags {
			if tags[key] != value {
				t.Fatalf("tag %q = %q, want %q", key, tags[key], value)
			}
		}
	})
}

func FuzzBootDiagnosticsObject(f *testing.F) {
	f.Add(false, "", "", "", "master")
	f.Add(false, "", "", "", "worker")
	f.Add(true, string(capz.DisabledDiagnosticsStorage), "", "", "master")
	f.Add(true, string(capz.ManagedDiagnosticsStorage), "", "", "worker")
	f.Add(true, string(capz.UserManagedDiagnosticsStorage), "account", "core.windows.net", "master")
	f.Add(true, string(capz.UserManagedDiagnosticsStorage), "", "core.windows.net", "worker")

	f.Fuzz(func(t *testing.T, hasDiag bool, diagType, account, suffix, role string) {
		var diag *aztypes.BootDiagnostics
		if hasDiag {
			diag = &aztypes.BootDiagnostics{
				Type:               capz.BootDiagnosticsStorageAccountType(diagType),
				StorageAccountName: account,
			}
		}
		got := getBootDiagnosticObject(diag, suffix, role)
		assertBootDiagnostics(t, got, hasDiag, diagType, account, suffix, role)
	})
}

func assertCapzImage(t *testing.T, image *capz.Image, osImage aztypes.OSImage, cloud string, confidential bool, gen, rg, sub, infraID, rhcos string) {
	t.Helper()
	switch {
	case osImage.Publisher != "":
		if image.Marketplace == nil || image.ID != nil {
			t.Fatalf("publisher image = %+v", image)
		}
		if image.Marketplace.Publisher != osImage.Publisher || image.Marketplace.Offer != osImage.Offer || image.Marketplace.SKU != osImage.SKU || image.Marketplace.Version != osImage.Version {
			t.Fatalf("marketplace fields = %+v", image.Marketplace)
		}
		if image.Marketplace.ThirdPartyImage != (osImage.Plan != aztypes.ImageNoPurchasePlan) {
			t.Fatalf("third party = %t, plan %q", image.Marketplace.ThirdPartyImage, osImage.Plan)
		}
	case cloud == string(aztypes.StackCloud):
		want := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/images/%s", sub, rg, infraID)
		if image.ID == nil || *image.ID != want || image.Marketplace != nil {
			t.Fatalf("stack image = %+v, want ID %q", image, want)
		}
	case strings.Count(rhcos, ":") == 3:
		parts := strings.Split(rhcos, ":")
		if len(parts) != 4 || image.Marketplace == nil || image.Marketplace.ThirdPartyImage {
			t.Fatalf("URN image %q = %+v", rhcos, image)
		}
		if image.Marketplace.Publisher != parts[0] || image.Marketplace.Offer != parts[1] || image.Marketplace.SKU != parts[2] || image.Marketplace.Version != parts[3] {
			t.Fatalf("URN parts %v, image %+v", parts, image.Marketplace)
		}
	case strings.HasPrefix(rhcos, "/subscriptions/"):
		if image.ID == nil || *image.ID != rhcos {
			t.Fatalf("resource image = %+v, want %q", image, rhcos)
		}
	case rhcos == "" && !confidential:
		if image.ID != nil || image.Marketplace != nil {
			t.Fatalf("empty RHCOS image = %+v, want empty", image)
		}
	default:
		gallery := strings.ReplaceAll(infraID, "-", "_")
		want := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/galleries/gallery_%s/images/%s", sub, rg, gallery, infraID)
		if gen == "V2" {
			want += "-gen2"
		}
		if image.ID == nil || *image.ID != want {
			t.Fatalf("gallery image = %+v, want %q", image, want)
		}
	}
}

func FuzzBootstrapFailureDomain(f *testing.F) {
	// No zones used to panic on Zones[0] while building the bootstrap machine.
	f.Add(false, "")
	f.Add(true, "")
	f.Add(true, "1")
	f.Add(true, "eastus-1")

	f.Fuzz(func(t *testing.T, hasZone bool, zone string) {
		var zones []string
		if hasZone {
			zones = []string{zone}
		}
		got := bootstrapFailureDomain(zones)
		if !hasZone || zone == "" {
			if got != nil {
				t.Fatalf("bootstrapFailureDomain(%q) = %q, want nil", zones, *got)
			}
			return
		}
		if got == nil || *got != zone {
			t.Fatalf("bootstrapFailureDomain(%q) = %v, want %q", zones, got, zone)
		}
	})
}

func assertBootDiagnostics(t *testing.T, got *machineapi.AzureDiagnostics, hasDiag bool, diagType, account, suffix, role string) {
	t.Helper()
	if !hasDiag {
		if role == controlPlaneRoleName {
			if got == nil || got.Boot == nil || got.Boot.StorageAccountType != machineapi.AzureManagedAzureDiagnosticsStorage || got.Boot.CustomerManaged != nil {
				t.Fatalf("default master diagnostics = %+v", got)
			}
			return
		}
		if got != nil {
			t.Fatalf("default %q diagnostics = %+v, want nil", role, got)
		}
		return
	}
	switch capz.BootDiagnosticsStorageAccountType(diagType) {
	case capz.DisabledDiagnosticsStorage:
		if got != nil {
			t.Fatalf("disabled diagnostics = %+v", got)
		}
	case capz.ManagedDiagnosticsStorage:
		if got == nil || got.Boot == nil || got.Boot.StorageAccountType != machineapi.AzureManagedAzureDiagnosticsStorage || got.Boot.CustomerManaged != nil {
			t.Fatalf("managed diagnostics = %+v", got)
		}
	default:
		if got == nil || got.Boot == nil || got.Boot.StorageAccountType != machineapi.CustomerManagedAzureDiagnosticsStorage || got.Boot.CustomerManaged == nil {
			t.Fatalf("customer managed diagnostics = %+v", got)
		}
		wantURI := ""
		if diagType == string(capz.UserManagedDiagnosticsStorage) && account != "" {
			wantURI = fmt.Sprintf("https://%s.blob.%s", account, suffix)
		}
		if got.Boot.CustomerManaged.StorageAccountURI != wantURI {
			t.Fatalf("storage URI %q, want %q", got.Boot.CustomerManaged.StorageAccountURI, wantURI)
		}
	}
}
