/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	klog "k8s.io/klog/v2"
)

/*
	TODO:
	When the INI based cloud-config is deprecated, this functions below should be preserved
*/

func getEnvKeyValue(match string, partial bool) (string, string, error) {
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		if len(pair) != 2 {
			continue
		}

		key := pair[0]
		value := pair[1]

		if partial && strings.Contains(key, match) {
			return key, value, nil
		}

		if strings.Compare(key, match) == 0 {
			return key, value, nil
		}
	}

	matchType := "match"
	if partial {
		matchType = "partial match"
	}

	return "", "", fmt.Errorf("Failed to find %s with %s", matchType, match)
}

// FromEnv initializes the provided configuratoin object with values
// obtained from environment variables. If an environment variable is set
// for a property that's already initialized, the environment variable's value
// takes precedence.
func (cfg *Config) FromEnv() error {

	//Init
	if cfg.VirtualCenter == nil {
		cfg.VirtualCenter = make(map[string]*VirtualCenterConfig)
	}

	//Globals
	if v := os.Getenv("VSPHERE_VCENTER"); v != "" {
		cfg.Global.VCenterIP = v
	}
	if v := os.Getenv("VSPHERE_VCENTER_PORT"); v != "" {
		cfg.Global.VCenterPort = v
	}
	if v := os.Getenv("VSPHERE_USER"); v != "" {
		cfg.Global.User = v
	}
	if v := os.Getenv("VSPHERE_PASSWORD"); v != "" {
		cfg.Global.Password = v
	}
	if v := os.Getenv("VSPHERE_DATACENTER"); v != "" {
		cfg.Global.Datacenters = v
	}
	if v := os.Getenv("VSPHERE_SECRET_NAME"); v != "" {
		cfg.Global.SecretName = v
	}
	if v := os.Getenv("VSPHERE_SECRET_NAMESPACE"); v != "" {
		cfg.Global.SecretNamespace = v
	}

	if v := os.Getenv("VSPHERE_ROUNDTRIP_COUNT"); v != "" {
		tmp, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			klog.Errorf("Failed to parse VSPHERE_ROUNDTRIP_COUNT: %s", err)
		} else {
			cfg.Global.RoundTripperCount = uint(tmp)
		}
	}

	if v := os.Getenv("VSPHERE_INSECURE"); v != "" {
		InsecureFlag, err := strconv.ParseBool(v)
		if err != nil {
			klog.Errorf("Failed to parse VSPHERE_INSECURE: %s", err)
		} else {
			cfg.Global.InsecureFlag = InsecureFlag
		}
	}

	if v := os.Getenv("VSPHERE_API_DISABLE"); v != "" {
		APIDisable, err := strconv.ParseBool(v)
		if err != nil {
			klog.Errorf("Failed to parse VSPHERE_API_DISABLE: %s", err)
		} else {
			cfg.Global.APIDisable = APIDisable
		}
	}

	if v := os.Getenv("VSPHERE_API_BINDING"); v != "" {
		cfg.Global.APIBinding = v
	}

	if v := os.Getenv("VSPHERE_SECRETS_DIRECTORY"); v != "" {
		cfg.Global.SecretsDirectory = v
	}
	if cfg.Global.SecretsDirectory == "" {
		cfg.Global.SecretsDirectory = DefaultSecretDirectory
	}
	if _, err := os.Stat(cfg.Global.SecretsDirectory); os.IsNotExist(err) {
		cfg.Global.SecretsDirectory = "" //Dir does not exist, set to empty string
	}

	if v := os.Getenv("VSPHERE_CAFILE"); v != "" {
		cfg.Global.CAFile = v
	}
	if v := os.Getenv("VSPHERE_THUMBPRINT"); v != "" {
		cfg.Global.Thumbprint = v
	}
	if v := os.Getenv("VSPHERE_LABEL_REGION"); v != "" {
		cfg.Labels.Region = v
	}
	if v := os.Getenv("VSPHERE_LABEL_ZONE"); v != "" {
		cfg.Labels.Zone = v
	}

	//Build VirtualCenter from ENVs
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")

		if len(pair) != 2 {
			continue
		}

		key := pair[0]
		value := pair[1]

		if strings.HasPrefix(key, "VSPHERE_VCENTER_") && len(value) > 0 {
			id := strings.TrimPrefix(key, "VSPHERE_VCENTER_")
			vcenter := value

			_, username, errUsername := getEnvKeyValue("VCENTER_"+id+"_USERNAME", false)
			if errUsername != nil {
				username = cfg.Global.User
			}
			_, password, errPassword := getEnvKeyValue("VCENTER_"+id+"_PASSWORD", false)
			if errPassword != nil {
				password = cfg.Global.Password
			}
			_, server, errServer := getEnvKeyValue("VCENTER_"+id+"_SERVER", false)
			if errServer != nil {
				server = ""
			}
			_, port, errPort := getEnvKeyValue("VCENTER_"+id+"_PORT", false)
			if errPort != nil {
				port = cfg.Global.VCenterPort
			}
			insecureFlag := false
			_, insecureTmp, errInsecure := getEnvKeyValue("VCENTER_"+id+"_INSECURE", false)
			if errInsecure != nil {
				insecureFlagTmp, errTmp := strconv.ParseBool(insecureTmp)
				if errTmp == nil {
					insecureFlag = insecureFlagTmp
				}
			}
			_, datacenters, errDatacenters := getEnvKeyValue("VCENTER_"+id+"_DATACENTERS", false)
			if errDatacenters != nil {
				datacenters = cfg.Global.Datacenters
			}
			roundtrip := DefaultRoundTripperCount
			_, roundtripTmp, errRoundtrip := getEnvKeyValue("VCENTER_"+id+"_ROUNDTRIP", false)
			if errRoundtrip != nil {
				roundtripFlagTmp, errTmp := strconv.ParseUint(roundtripTmp, 10, 32)
				if errTmp == nil {
					roundtrip = uint(roundtripFlagTmp)
				}
			}
			_, caFile, errCaFile := getEnvKeyValue("VCENTER_"+id+"_CAFILE", false)
			if errCaFile != nil {
				caFile = cfg.Global.CAFile
			}
			_, thumbprint, errThumbprint := getEnvKeyValue("VCENTER_"+id+"_THUMBPRINT", false)
			if errThumbprint != nil {
				thumbprint = cfg.Global.Thumbprint
			}

			_, secretName, secretNameErr := getEnvKeyValue("VCENTER_"+id+"_SECRET_NAME", false)
			_, secretNamespace, secretNamespaceErr := getEnvKeyValue("VCENTER_"+id+"_SECRET_NAMESPACE", false)

			if secretNameErr != nil || secretNamespaceErr != nil {
				secretName = ""
				secretNamespace = ""
			}
			secretRef := DefaultCredentialManager
			if secretName != "" && secretNamespace != "" {
				secretRef = vcenter
			}

			iPFamilyPriority := []string{DefaultIPFamily}
			_, ipFamily, errIPFamily := getEnvKeyValue("VCENTER_"+id+"_IP_FAMILY", false)
			if errIPFamily != nil {
				iPFamilyPriority = []string{ipFamily}
			}

			// If server is explicitly set, that means the vcenter value above is the TenantRef
			vcenterIP := vcenter
			tenantRef := vcenter
			if server != "" {
				vcenterIP = server
				tenantRef = vcenter
			}

			var vcc *VirtualCenterConfig
			if cfg.VirtualCenter[tenantRef] != nil {
				vcc = cfg.VirtualCenter[tenantRef]
			} else {
				vcc = &VirtualCenterConfig{}
				cfg.VirtualCenter[tenantRef] = vcc
			}

			vcc.User = username
			vcc.Password = password
			vcc.TenantRef = tenantRef
			vcc.VCenterIP = vcenterIP
			vcc.VCenterPort = port
			vcc.InsecureFlag = insecureFlag
			vcc.Datacenters = datacenters
			vcc.RoundTripperCount = roundtrip
			vcc.CAFile = caFile
			vcc.Thumbprint = thumbprint
			vcc.SecretRef = secretRef
			vcc.SecretName = secretName
			vcc.SecretNamespace = secretNamespace
			vcc.IPFamilyPriority = iPFamilyPriority
		}
	}

	return nil
}

/*
	TODO:
	When the INI based cloud-config is deprecated, the references to the
	INI based code (ie the call to ReadConfigINI) below should be deleted.
*/

// ReadConfig parses vSphere cloud config file and stores it into VSphereConfig.
// Environment variables are also checked
func ReadConfig(byConfig []byte) (*Config, error) {
	if len(byConfig) == 0 {
		return nil, fmt.Errorf("Invalid YAML/INI file")
	}

	cfg, err := ReadConfigYAML(byConfig)
	if err != nil {
		klog.Warningf("ReadConfigYAML failed: %s", err)

		cfg, err = ReadConfigINI(byConfig)
		if err != nil {
			klog.Errorf("ReadConfigINI failed: %s", err)
			return nil, err
		}

		klog.Info("ReadConfig INI succeeded. INI-based cloud-config is deprecated and will be removed in 2.0. Please use YAML based cloud-config.")
	} else {
		klog.Info("ReadConfig YAML succeeded")
	}

	// Env Vars should override config file entries if present
	if err := cfg.FromEnv(); err != nil {
		klog.Errorf("FromEnv failed: %s", err)
		return nil, err
	}

	klog.Info("Config initialized")
	return cfg, nil
}
