package azure

import (
	"strings"
	"testing"
)

func FuzzGroupFromID(f *testing.F) {
	// A DNS zone ID shorter than five slash-separated parts used to panic on [4].
	f.Add("")
	f.Add("/subscriptions/sub")
	f.Add("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/dnszones/example.com")

	f.Fuzz(func(t *testing.T, id string) {
		got := groupFromID(id)
		parts := strings.Split(id, "/")
		want := ""
		if len(parts) >= 5 {
			want = parts[4]
		}
		if got != want {
			t.Fatalf("groupFromID(%q) = %q, want %q", id, got, want)
		}
	})
}
