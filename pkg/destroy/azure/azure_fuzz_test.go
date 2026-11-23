package azure

import (
	"strings"
	"testing"
)

func FuzzGroupFromID(f *testing.F) {
	f.Add("my-rg")
	f.Add("")
	f.Add("a/resource-group")

	f.Fuzz(func(t *testing.T, resourceGroup string) {
		// Arbitrary malformed IDs must not panic.
		_ = groupFromID(resourceGroup)

		component := strings.ReplaceAll(resourceGroup, "/", "_")
		id := "/subscriptions/sub/resourceGroups/" + component + "/providers/Microsoft.Network/dnszones/example.com"
		if got := groupFromID(id); got != component {
			t.Fatalf("groupFromID(valid resource ID) = %q, want %q", got, component)
		}
		for _, malformedID := range []string{"", "/subscriptions/sub", "/subscriptions/sub/resourceGroups"} {
			if got := groupFromID(malformedID); got != "" {
				t.Fatalf("groupFromID(%q) = %q, want empty for a short ID", malformedID, got)
			}
		}
	})
}
