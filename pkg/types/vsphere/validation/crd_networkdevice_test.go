package validation

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/yaml"
)

func loadInstallConfigCRD(t *testing.T) *apiextensionsv1.CustomResourceDefinition {
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

func networkDeviceSchemaProps(t *testing.T, crd *apiextensionsv1.CustomResourceDefinition) map[string]apiextensionsv1.JSONSchemaProps {
	t.Helper()

	if len(crd.Spec.Versions) == 0 || crd.Spec.Versions[0].Schema == nil || crd.Spec.Versions[0].Schema.OpenAPIV3Schema == nil {
		t.Fatal("CRD OpenAPI schema is missing")
	}
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

// TestCRDNetworkDeviceSpecFormats guards against controller-gen collapsing dual
// Format=ipv4/Format=ipv6 markers into format: ipv6 (openshift/installer#10377).
// gateway and nameservers must use format: ip; ipAddrs must not use an IP format
// because entries are CIDRs.
func TestCRDNetworkDeviceSpecFormats(t *testing.T) {
	props := networkDeviceSchemaProps(t, loadInstallConfigCRD(t))

	t.Run("gateway", func(t *testing.T) {
		gw, ok := props["gateway"]
		if !assert.True(t, ok, "gateway property must exist") {
			return
		}
		assert.Equal(t, "ip", gw.Format, "gateway must use format: ip to accept IPv4 and IPv6")
	})

	t.Run("ipAddrs", func(t *testing.T) {
		ip, ok := props["ipAddrs"]
		if !assert.True(t, ok, "ipAddrs property must exist") {
			return
		}
		assert.Empty(t, ip.Format, "ipAddrs must not set format; values are CIDRs not bare IPs")
	})

	t.Run("nameservers", func(t *testing.T) {
		ns, ok := props["nameservers"]
		if !assert.True(t, ok, "nameservers property must exist") {
			return
		}
		assert.Equal(t, "ip", ns.Format, "nameservers must use format: ip to accept IPv4 and IPv6")
	})
}
