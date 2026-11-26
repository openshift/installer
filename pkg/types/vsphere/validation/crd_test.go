package validation

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/yaml"
)

func loadCRD(t *testing.T) *apiextensionsv1.CustomResourceDefinition {
	t.Helper()
	data, err := os.ReadFile("../../../../data/data/install.openshift.io_installconfigs.yaml")
	if err != nil {
		t.Fatalf("failed to load CRD: %v", err)
	}
	var crd apiextensionsv1.CustomResourceDefinition
	if err := yaml.Unmarshal(data, &crd); err != nil {
		t.Fatalf("failed to unmarshal CRD: %v", err)
	}
	return &crd
}

func networkDeviceSchema(t *testing.T, crd *apiextensionsv1.CustomResourceDefinition) map[string]apiextensionsv1.JSONSchemaProps {
	t.Helper()

	schema := crd.Spec.Versions[0].Schema.OpenAPIV3Schema
	platform, ok := schema.Properties["platform"]
	if !ok {
		t.Fatal("missing platform property in CRD schema")
	}
	vsphere, ok := platform.Properties["vsphere"]
	if !ok {
		t.Fatal("missing vsphere property in CRD schema")
	}
	hosts, ok := vsphere.Properties["hosts"]
	if !ok {
		t.Fatal("missing hosts property in CRD schema")
	}
	if hosts.Items == nil || hosts.Items.Schema == nil {
		t.Fatal("hosts items schema is nil")
	}
	netDev, ok := hosts.Items.Schema.Properties["networkDevice"]
	if !ok {
		t.Fatal("missing networkDevice property in CRD schema")
	}
	return netDev.Properties
}

func TestCRDNetworkDeviceSpec(t *testing.T) {
	crd := loadCRD(t)
	props := networkDeviceSchema(t, crd)

	t.Run("gateway", func(t *testing.T) {
		gw, ok := props["gateway"]
		if !assert.True(t, ok, "gateway property must exist") {
			return
		}
		assert.Empty(t, gw.Format, "gateway must not have a format constraint")
		if assert.Len(t, gw.XValidations, 1, "gateway must have exactly 1 CEL validation rule") {
			assert.Equal(t, "self == '' || isIP(self)", gw.XValidations[0].Rule)
			assert.NotEmpty(t, gw.XValidations[0].Message)
		}
	})

	t.Run("ipAddrs", func(t *testing.T) {
		ip, ok := props["ipAddrs"]
		if !assert.True(t, ok, "ipAddrs property must exist") {
			return
		}
		assert.Empty(t, ip.Format, "ipAddrs must not have a format constraint")
		if assert.Len(t, ip.XValidations, 1, "ipAddrs must have exactly 1 CEL validation rule") {
			assert.Equal(t, "self.all(x, isCIDR(x))", ip.XValidations[0].Rule)
		}
		if assert.NotNil(t, ip.MaxItems, "ipAddrs must have maxItems set") {
			assert.Equal(t, int64(10), *ip.MaxItems)
		}
	})

	t.Run("nameservers", func(t *testing.T) {
		ns, ok := props["nameservers"]
		if !assert.True(t, ok, "nameservers property must exist") {
			return
		}
		assert.Empty(t, ns.Format, "nameservers must not have a format constraint")
		if assert.Len(t, ns.XValidations, 1, "nameservers must have exactly 1 CEL validation rule") {
			assert.Equal(t, "self.all(x, isIP(x))", ns.XValidations[0].Rule)
		}
		if assert.NotNil(t, ns.MaxItems, "nameservers must have maxItems set") {
			assert.Equal(t, int64(3), *ns.MaxItems)
		}
	})
}
