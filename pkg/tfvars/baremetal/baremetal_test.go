// Package baremetal contains bare metal specific Terraform-variable logic.
package baremetal

import (
	"encoding/json"
	"testing"

	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/stretchr/testify/assert"
)

func TestMastersSelectionByRole(t *testing.T) {

	cases := []struct {
		scenario string

		numControlPlaneReplicas int64
		libvirtURI              string
		apiVIP                  string
		imageCacheIP            string
		bootstrapOSImage        string
		externalBridge          string
		externalMAC             string
		provisioningBridge      string
		provisioningMAC         string
		platformHosts           []*baremetal.Host
		image                   string
		ironicUsername          string
		ironicPassword          string
		ignition                string

		expectedError string
		expectedHosts []string
	}{
		{
			scenario:                "filter-workers",
			numControlPlaneReplicas: 2,
			platformHosts: platformHosts(
				host("master-0", "master"),
				host("master-1", "master"),
				host("worker-0", "worker"),
			),
			expectedHosts: []string{"master-0", "master-1"},
		},
		{
			scenario:                "filter-all-workers",
			numControlPlaneReplicas: 1,
			platformHosts: platformHosts(
				host("master-0", "master"),
				host("worker-0", "worker"),
				host("worker-1", "worker"),
			),
			expectedHosts: []string{"master-0"},
		},
		{
			scenario:                "hosts-norole",
			numControlPlaneReplicas: 2,
			platformHosts: platformHosts(
				host("worker-0", "worker"),
				host("master-0", ""),
				host("master-1", ""),
			),
			expectedHosts: []string{"master-0", "master-1"},
		},
		{
			scenario:                "no-role",
			numControlPlaneReplicas: 3,
			platformHosts: platformHosts(
				host("master-0", ""),
				host("master-1", ""),
				host("master-2", ""),
			),
			expectedHosts: []string{"master-0", "master-1", "master-2"},
		},
		{
			scenario:                "more-masters-than-required",
			numControlPlaneReplicas: 2,
			platformHosts: platformHosts(
				host("master-0", "master"),
				host("master-1", "master"),
				host("master-2", "master"),
			),
			expectedHosts: []string{"master-0", "master-1"},
		},
		{
			scenario:                "more-hosts-than-required-mixed",
			numControlPlaneReplicas: 2,
			platformHosts: platformHosts(
				host("master-0", "master"),
				host("master-1", "master"),
				host("master-2", ""),
			),
			expectedHosts: []string{"master-0", "master-1"},
		},
		{
			scenario:                "more-hosts-than-required-norole",
			numControlPlaneReplicas: 2,
			platformHosts: platformHosts(
				host("master-0", ""),
				host("master-1", ""),
				host("master-2", ""),
			),
			expectedHosts: []string{"master-0", "master-1"},
		},
		{
			scenario:                "more-hosts-than-required-mixed-again",
			numControlPlaneReplicas: 2,
			platformHosts: platformHosts(
				host("worker-0", "worker"),
				host("worker-1", "worker"),
				host("worker-2", "worker"),
				host("master-0", ""),
				host("master-1", ""),
				host("master-2", ""),
			),
			expectedHosts: []string{"master-0", "master-1"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.scenario, func(t *testing.T) {

			imageDownloader = func(baseURL string) (string, error) {
				return "", nil
			}

			data, err := TFVars(
				tc.numControlPlaneReplicas,
				tc.libvirtURI,
				tc.apiVIP,
				tc.imageCacheIP,
				tc.bootstrapOSImage,
				tc.externalBridge,
				tc.externalMAC,
				tc.provisioningBridge,
				tc.provisioningMAC,
				tc.platformHosts,
				tc.image,
				tc.ironicUsername,
				tc.ironicPassword,
				tc.ignition)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expectedError, err)
			}

			var cfg config
			err = json.Unmarshal(data, &cfg)
			assert.Nil(t, err)

			assert.Equal(t, len(tc.expectedHosts), len(cfg.Masters))
			for i, hostName := range tc.expectedHosts {
				assert.Equal(t, hostName, cfg.Masters[i]["name"])
			}
		})
	}
}

func host(name, tag string) *baremetal.Host {
	return &baremetal.Host{
		Name:            name,
		Role:            tag,
		HardwareProfile: "default",
		BMC: baremetal.BMC{
			Address: "redfish+http://192.168.111.1:8000/redfish/v1/Systems/e4427260-6250-4df9-9e8a-120f78a46aa6",
		},
	}
}

func platformHosts(hosts ...*baremetal.Host) []*baremetal.Host {
	return hosts
}
