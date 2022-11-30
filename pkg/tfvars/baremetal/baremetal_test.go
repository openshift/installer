// Package baremetal contains bare metal specific Terraform-variable logic.
package baremetal

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/types/baremetal"
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
		hostFiles               []*asset.File
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
			hostFiles: hostFiles(
				files("master-0", nil),
				files("master-2", nil),
				files("worker-0", nil),
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
			hostFiles: hostFiles(
				files("master-0", nil),
				files("worker-0", nil),
				files("worker-1", nil),
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
			hostFiles: hostFiles(
				files("worker-0", nil),
				files("master-0", nil),
				files("master-1", nil),
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
			hostFiles: hostFiles(
				files("master-0", nil),
				files("master-1", nil),
				files("master-2", nil),
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
			hostFiles: hostFiles(
				files("master-0", nil),
				files("master-1", nil),
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
			hostFiles: hostFiles(
				files("master-0", nil),
				files("master-1", nil),
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
			hostFiles: hostFiles(
				files("master-0", nil),
				files("master-1", nil),
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
			hostFiles: hostFiles(
				files("worker-0", nil),
				files("worker-1", nil),
				files("worker-2", nil),
				files("master-0", nil),
				files("master-1", nil),
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
				tc.hostFiles,
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

func TestRAIDBIOSConfig(t *testing.T) {

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
		hostFiles               []*asset.File
		image                   string
		ironicUsername          string
		ironicPassword          string
		ignition                string

		expectedError        string
		expectedRAIDConfig   []string
		expectedBIOSSettings []string
	}{
		{
			scenario:                "raid",
			numControlPlaneReplicas: 2,
			platformHosts: platformHosts(
				iRMChost("master-0", "master"),
				iRMChost("master-1", "master"),
			),
			hostFiles: hostFiles(
				files("master-0", configuration("raid")),
				files("master-1", configuration("raid")),
			),
			expectedRAIDConfig: []string{"{\"hardwareRAIDVolumes\":[{\"level\":\"0\",\"name\":\"raid0\"}],\"softwareRAIDVolumes\":null}",
				"{\"hardwareRAIDVolumes\":[{\"level\":\"0\",\"name\":\"raid0\"}],\"softwareRAIDVolumes\":null}"},
			expectedBIOSSettings: []string{"", ""},
		},
		{
			scenario:                "bios",
			numControlPlaneReplicas: 2,
			platformHosts: platformHosts(
				iRMChost("master-0", "master"),
				iRMChost("master-1", "master"),
			),
			hostFiles: hostFiles(
				files("master-0", configuration("bios")),
				files("master-1", configuration("bios")),
			),
			expectedRAIDConfig: []string{"", ""},
			expectedBIOSSettings: []string{"[{\"name\":\"cpu_vt_enabled\",\"value\":\"True\"},{\"name\":\"hyper_threading_enabled\",\"value\":\"True\"}]",
				"[{\"name\":\"cpu_vt_enabled\",\"value\":\"True\"},{\"name\":\"hyper_threading_enabled\",\"value\":\"True\"}]"},
		},
		{
			scenario:                "raid and bios",
			numControlPlaneReplicas: 2,
			platformHosts: platformHosts(
				iRMChost("master-0", "master"),
				iRMChost("master-1", "master"),
			),
			hostFiles: hostFiles(
				files("master-0", configuration("raidbios")),
				files("master-1", configuration("raidbios")),
			),
			expectedRAIDConfig: []string{"{\"hardwareRAIDVolumes\":[{\"level\":\"0\",\"name\":\"raid0\"}],\"softwareRAIDVolumes\":null}",
				"{\"hardwareRAIDVolumes\":[{\"level\":\"0\",\"name\":\"raid0\"}],\"softwareRAIDVolumes\":null}"},
			expectedBIOSSettings: []string{"[{\"name\":\"cpu_vt_enabled\",\"value\":\"True\"},{\"name\":\"hyper_threading_enabled\",\"value\":\"True\"}]",
				"[{\"name\":\"cpu_vt_enabled\",\"value\":\"True\"},{\"name\":\"hyper_threading_enabled\",\"value\":\"True\"}]"},
		},
		{
			scenario:                "no raid and no bios",
			numControlPlaneReplicas: 2,
			platformHosts: platformHosts(
				iRMChost("master-0", "master"),
				iRMChost("master-1", "master"),
			),
			hostFiles: hostFiles(
				files("master-0", configuration("")),
				files("master-1", configuration("")),
			),
			expectedRAIDConfig:   []string{"", ""},
			expectedBIOSSettings: []string{"", ""},
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
				tc.hostFiles,
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

			for i, RAIDConfig := range tc.expectedRAIDConfig {
				assert.Equal(t, RAIDConfig, cfg.Masters[i]["raid_config"])
			}

			for i, BIOSSettins := range tc.expectedBIOSSettings {
				assert.Equal(t, BIOSSettins, cfg.Masters[i]["bios_settings"])
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

func iRMChost(name, tag string) *baremetal.Host {
	return &baremetal.Host{
		Name:            name,
		Role:            tag,
		HardwareProfile: "default",
		BMC: baremetal.BMC{
			Address: "irmc://127.0.0.1",
		},
	}
}

func platformHosts(hosts ...*baremetal.Host) []*baremetal.Host {
	return hosts
}

func files(name string, data []byte) *asset.File {
	return &asset.File{
		Filename: name,
		Data:     data,
	}
}

func hostFiles(files ...*asset.File) []*asset.File {
	return files
}

func configuration(cfg string) []byte {
	switch cfg {
	case "raid":
		return []byte(`
spec:
  raid:
    hardwareRAIDVolumes:
    - level: "0"
      name: "raid0"
`)
	case "bios":
		return []byte(`
spec:
  firmware:
    virtualizationEnabled: true
    simultaneousMultithreadingEnabled: true
`)
	case "raidbios":
		return []byte(`
spec:
  raid:
    hardwareRAIDVolumes:
    - level: "0"
      name: "raid0"
  firmware:
    virtualizationEnabled: true
    simultaneousMultithreadingEnabled: true
`)
	default:
		return nil
	}
}
