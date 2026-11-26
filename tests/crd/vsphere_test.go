package crd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVSphereSchema(t *testing.T) {
	crd := loadCRD(t)
	root := rootSchema(t, crd)

	platformSchema := schemaAtPath(t, root, "platform", "vsphere")

	t.Run("Platform", func(t *testing.T) {
		t.Run("apiVIPs", func(t *testing.T) {
			s := schemaAtPath(t, platformSchema, "apiVIPs")
			requireArrayItems(t, s, nil, int64Ptr(2))
			requireFormat(t, s, "ip")
		})

		t.Run("ingressVIPs", func(t *testing.T) {
			s := schemaAtPath(t, platformSchema, "ingressVIPs")
			requireArrayItems(t, s, nil, int64Ptr(2))
			requireFormat(t, s, "ip")
		})

		t.Run("vcenters", func(t *testing.T) {
			s := schemaAtPath(t, platformSchema, "vcenters")
			requireArrayItems(t, s, int64Ptr(1), int64Ptr(3))
		})
	})

	t.Run("VCenter", func(t *testing.T) {
		vcenterSchema := schemaAtPath(t, platformSchema, "vcenters", "[]")

		t.Run("server", func(t *testing.T) {
			s := schemaAtPath(t, vcenterSchema, "server")
			requireStringLength(t, s, nil, int64Ptr(255))
		})
		t.Run("port", func(t *testing.T) {
			s := schemaAtPath(t, vcenterSchema, "port")
			requireNumericRange(t, s, float64Ptr(1), float64Ptr(32767))
		})
		t.Run("datacenters", func(t *testing.T) {
			s := schemaAtPath(t, vcenterSchema, "datacenters")
			requireArrayItems(t, s, int64Ptr(1), nil)
		})
		t.Run("required fields", func(t *testing.T) {
			requireRequired(t, vcenterSchema, "server", "user", "password", "datacenters")
		})
	})

	t.Run("FailureDomain", func(t *testing.T) {
		fdSchema := schemaAtPath(t, platformSchema, "failureDomains", "[]")

		t.Run("name", func(t *testing.T) {
			s := schemaAtPath(t, fdSchema, "name")
			requireStringLength(t, s, int64Ptr(1), int64Ptr(256))
		})
		t.Run("server", func(t *testing.T) {
			s := schemaAtPath(t, fdSchema, "server")
			requireStringLength(t, s, int64Ptr(1), int64Ptr(255))
		})
		t.Run("regionType", func(t *testing.T) {
			s := schemaAtPath(t, fdSchema, "regionType")
			requireEnum(t, s, "Datacenter", "ComputeCluster")
		})
		t.Run("zoneType", func(t *testing.T) {
			s := schemaAtPath(t, fdSchema, "zoneType")
			requireEnum(t, s, "ComputeCluster", "HostGroup")
		})
		t.Run("required fields", func(t *testing.T) {
			requireRequired(t, fdSchema, "name", "region", "zone", "server", "topology")
		})
	})

	t.Run("Topology", func(t *testing.T) {
		topoSchema := schemaAtPath(t, platformSchema, "failureDomains", "[]", "topology")

		t.Run("datacenter", func(t *testing.T) {
			s := schemaAtPath(t, topoSchema, "datacenter")
			requireStringLength(t, s, int64Ptr(1), int64Ptr(80))
		})
		t.Run("computeCluster", func(t *testing.T) {
			s := schemaAtPath(t, topoSchema, "computeCluster")
			requireStringLength(t, s, int64Ptr(1), int64Ptr(2048))
		})
		t.Run("networks", func(t *testing.T) {
			s := schemaAtPath(t, topoSchema, "networks")
			requireArrayItems(t, s, int64Ptr(1), int64Ptr(10))
		})
		t.Run("datastore", func(t *testing.T) {
			s := schemaAtPath(t, topoSchema, "datastore")
			requireStringLength(t, s, int64Ptr(1), int64Ptr(2048))
		})
		t.Run("resourcePool", func(t *testing.T) {
			s := schemaAtPath(t, topoSchema, "resourcePool")
			requireStringLength(t, s, int64Ptr(1), int64Ptr(2048))
			requirePattern(t, s, `^/.*?/host/.*?/Resources.*`)
		})
		t.Run("folder", func(t *testing.T) {
			s := schemaAtPath(t, topoSchema, "folder")
			requireStringLength(t, s, int64Ptr(1), int64Ptr(2048))
			requirePattern(t, s, `^/.*?/vm/.*?`)
		})
		t.Run("template", func(t *testing.T) {
			s := schemaAtPath(t, topoSchema, "template")
			requireStringLength(t, s, int64Ptr(1), int64Ptr(2048))
			requirePattern(t, s, `^/.*?/vm/.*?`)
		})
		t.Run("hostGroup", func(t *testing.T) {
			s := schemaAtPath(t, topoSchema, "hostGroup")
			requireStringLength(t, s, nil, int64Ptr(80))
		})
		t.Run("required fields", func(t *testing.T) {
			requireRequired(t, topoSchema, "datacenter", "computeCluster", "networks", "datastore")
		})
	})

	t.Run("Host", func(t *testing.T) {
		hostSchema := schemaAtPath(t, platformSchema, "hosts", "[]")

		t.Run("role", func(t *testing.T) {
			s := schemaAtPath(t, hostSchema, "role")
			requireEnum(t, s, "", "bootstrap", "control-plane", "compute")
		})
		t.Run("required fields", func(t *testing.T) {
			requireRequired(t, hostSchema, "networkDevice", "role")
		})
	})

	t.Run("NetworkDeviceSpec", func(t *testing.T) {
		netDevSchema := schemaAtPath(t, platformSchema, "hosts", "[]", "networkDevice")

		t.Run("gateway", func(t *testing.T) {
			s := schemaAtPath(t, netDevSchema, "gateway")
			requireFormat(t, s, "")
			requireStringLength(t, s, nil, int64Ptr(45))
			requireCELRule(t, s, "self == '' || isIP(self)")
		})
		t.Run("ipAddrs", func(t *testing.T) {
			s := schemaAtPath(t, netDevSchema, "ipAddrs")
			requireFormat(t, s, "")
			requireArrayItems(t, s, nil, int64Ptr(10))
			requireCELRule(t, s, "self.all(x, isCIDR(x))")
		})
		t.Run("nameservers", func(t *testing.T) {
			s := schemaAtPath(t, netDevSchema, "nameservers")
			requireFormat(t, s, "")
			requireArrayItems(t, s, nil, int64Ptr(3))
			requireCELRule(t, s, "self.all(x, isIP(x))")
		})
		t.Run("ipAddrs is required", func(t *testing.T) {
			requireRequired(t, netDevSchema, "ipAddrs")
		})
	})

	t.Run("MachinePool", func(t *testing.T) {
		mpSchema := schemaAtPath(t, platformSchema, "defaultMachinePlatform")

		t.Run("dataDisks maxItems", func(t *testing.T) {
			s := schemaAtPath(t, mpSchema, "dataDisks")
			requireArrayItems(t, s, nil, int64Ptr(29))
		})
	})

	t.Run("DataDisk", func(t *testing.T) {
		ddSchema := schemaAtPath(t, platformSchema, "defaultMachinePlatform", "dataDisks", "[]")

		t.Run("name", func(t *testing.T) {
			s := schemaAtPath(t, ddSchema, "name")
			requireStringLength(t, s, nil, int64Ptr(80))
			requirePattern(t, s, `^[a-zA-Z0-9]([-_a-zA-Z0-9]*[a-zA-Z0-9])?$`)
		})
		t.Run("sizeGiB", func(t *testing.T) {
			s := schemaAtPath(t, ddSchema, "sizeGiB")
			requireNumericRange(t, s, float64Ptr(1), float64Ptr(16384))
		})
		t.Run("provisioningMode", func(t *testing.T) {
			s := schemaAtPath(t, ddSchema, "provisioningMode")
			requireEnum(t, s, "Thin", "Thick", "EagerlyZeroed")
		})
	})
}

func TestVSphereCELValidation(t *testing.T) {
	skipIfNoEnvtest(t)
	validIPAddrs := []any{"192.168.1.100/24"}
	validNameservers := []any{"8.8.8.8"}

	buildCR := func(gateway string, ipAddrs, nameservers []any) map[string]any {
		netDev := map[string]any{
			"ipAddrs": ipAddrs,
		}
		if gateway != "" {
			netDev["gateway"] = gateway
		}
		if len(nameservers) > 0 {
			netDev["nameservers"] = nameservers
		}

		return map[string]any{
			"apiVersion": "install.openshift.io/v1",
			"kind":       "InstallConfig",
			"metadata": map[string]any{
				"generateName": "cel-test-",
				"namespace":    "default",
			},
			"baseDomain": "example.com",
			"platform": map[string]any{
				"vsphere": map[string]any{
					"hosts": []any{
						map[string]any{
							"role":          "bootstrap",
							"networkDevice": netDev,
						},
					},
				},
			},
			"pullSecret": "{\"auths\":{}}",
		}
	}

	t.Run("gateway", func(t *testing.T) {
		cases := []struct {
			name      string
			gateway   string
			expectErr bool
		}{
			{name: "empty string", gateway: "", expectErr: false},
			{name: "valid IPv4", gateway: "192.168.1.1", expectErr: false},
			{name: "valid IPv6", gateway: "2001:db8::1", expectErr: false},
			{name: "invalid string", gateway: "not-an-ip", expectErr: true},
			{name: "CIDR not bare IP", gateway: "192.168.1.0/24", expectErr: true},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				cr := buildCR(tc.gateway, validIPAddrs, validNameservers)
				err := createCR(t, cr)
				if tc.expectErr {
					assert.Error(t, err, "expected validation error for gateway=%q", tc.gateway)
				} else {
					assert.NoError(t, err, "expected no error for gateway=%q", tc.gateway)
				}
			})
		}
	})

	t.Run("ipAddrs", func(t *testing.T) {
		cases := []struct {
			name      string
			ipAddrs   []any
			expectErr bool
		}{
			{name: "valid IPv4 CIDR", ipAddrs: []any{"192.168.1.100/24"}, expectErr: false},
			{name: "valid IPv6 CIDR", ipAddrs: []any{"2001:db8::1/64"}, expectErr: false},
			{name: "mixed families", ipAddrs: []any{"192.168.1.100/24", "2001:db8::1/64"}, expectErr: false},
			{name: "bare IP without prefix", ipAddrs: []any{"192.168.1.1"}, expectErr: true},
			{name: "invalid string", ipAddrs: []any{"not-a-cidr"}, expectErr: true},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				cr := buildCR("", tc.ipAddrs, validNameservers)
				err := createCR(t, cr)
				if tc.expectErr {
					assert.Error(t, err, "expected validation error for ipAddrs=%v", tc.ipAddrs)
				} else {
					assert.NoError(t, err, "expected no error for ipAddrs=%v", tc.ipAddrs)
				}
			})
		}
	})

	t.Run("nameservers", func(t *testing.T) {
		cases := []struct {
			name        string
			nameservers []any
			expectErr   bool
		}{
			{name: "valid IPv4", nameservers: []any{"8.8.8.8"}, expectErr: false},
			{name: "valid IPv6", nameservers: []any{"2001:4860:4860::8888"}, expectErr: false},
			{name: "mixed families", nameservers: []any{"8.8.8.8", "2001:4860:4860::8888"}, expectErr: false},
			{name: "invalid string", nameservers: []any{"not-an-ip"}, expectErr: true},
			{name: "CIDR not bare IP", nameservers: []any{"8.8.8.8/32"}, expectErr: true},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				cr := buildCR("", validIPAddrs, tc.nameservers)
				err := createCR(t, cr)
				if tc.expectErr {
					assert.Error(t, err, "expected validation error for nameservers=%v", tc.nameservers)
				} else {
					assert.NoError(t, err, "expected no error for nameservers=%v", tc.nameservers)
				}
			})
		}
	})
}

