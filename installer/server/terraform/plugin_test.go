package terraform

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
)

func TestProvider(t *testing.T) {
	for _, plugin := range plugins {
		if err := plugin.ProviderFunc().(*schema.Provider).InternalValidate(); err != nil {
			t.Fatalf("err: %s", err)
		}
	}
}
