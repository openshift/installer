/*
Copyright 2020 The Kubernetes Authors.

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
	"strings"

	yaml "gopkg.in/yaml.v2"
	klog "k8s.io/klog/v2"
)

/*
	TODO:
	When the INI based cloud-config is deprecated, this file should be merged into config.go
	and this file should be deleted.
*/

// CreateConfig generates a common Config object based on what other structs and funcs
// are already dependent upon in other packages.
func (ccy *CommonConfigYAML) CreateConfig() *Config {
	cfg := &Config{
		VirtualCenter: make(map[string]*VirtualCenterConfig),
	}

	cfg.Global.User = ccy.Global.User
	cfg.Global.Password = ccy.Global.Password
	cfg.Global.VCenterIP = ccy.Global.VCenterIP
	cfg.Global.VCenterPort = fmt.Sprint(ccy.Global.VCenterPort)
	cfg.Global.InsecureFlag = ccy.Global.InsecureFlag
	cfg.Global.Datacenters = strings.Join(ccy.Global.Datacenters, ",")
	cfg.Global.RoundTripperCount = ccy.Global.RoundTripperCount
	cfg.Global.CAFile = ccy.Global.CAFile
	cfg.Global.Thumbprint = ccy.Global.Thumbprint
	cfg.Global.SecretName = ccy.Global.SecretName
	cfg.Global.SecretNamespace = ccy.Global.SecretNamespace
	cfg.Global.SecretsDirectory = ccy.Global.SecretsDirectory

	for keyVcConfig, valVcConfig := range ccy.Vcenter {
		cfg.VirtualCenter[keyVcConfig] = &VirtualCenterConfig{
			User:              valVcConfig.User,
			Password:          valVcConfig.Password,
			TenantRef:         valVcConfig.TenantRef,
			VCenterIP:         valVcConfig.VCenterIP,
			VCenterPort:       fmt.Sprint(valVcConfig.VCenterPort),
			InsecureFlag:      valVcConfig.InsecureFlag,
			Datacenters:       strings.Join(valVcConfig.Datacenters, ","),
			RoundTripperCount: valVcConfig.RoundTripperCount,
			CAFile:            valVcConfig.CAFile,
			Thumbprint:        valVcConfig.Thumbprint,
			SecretRef:         valVcConfig.SecretRef,
			SecretName:        valVcConfig.SecretName,
			SecretNamespace:   valVcConfig.SecretNamespace,
			IPFamilyPriority:  valVcConfig.IPFamilyPriority,
		}
	}

	cfg.Labels.Region = ccy.Labels.Region
	cfg.Labels.Zone = ccy.Labels.Zone

	return cfg
}

// isSecretInfoProvided returns true if k8s secret is set or using generic CO secret method.
// If both k8s secret and generic CO both are true, we don't know which to use, so return false.
func (ccy *CommonConfigYAML) isSecretInfoProvided() bool {
	return (ccy.Global.SecretName != "" && ccy.Global.SecretNamespace != "" && ccy.Global.SecretsDirectory == "") ||
		(ccy.Global.SecretName == "" && ccy.Global.SecretNamespace == "" && ccy.Global.SecretsDirectory != "")
}

// isSecretInfoProvided returns true if the secret per VC has been configured
func (vccy *VirtualCenterConfigYAML) isSecretInfoProvided() bool {
	return vccy.SecretName != "" && vccy.SecretNamespace != ""
}

func (ccy *CommonConfigYAML) validateConfig() error {
	//Fix default global values
	if ccy.Global.RoundTripperCount == 0 {
		ccy.Global.RoundTripperCount = DefaultRoundTripperCount
	}
	if ccy.Global.VCenterPort == 0 {
		ccy.Global.VCenterPort = DefaultVCenterPort
	}
	if ccy.Global.APIBinding == "" {
		ccy.Global.APIBinding = DefaultAPIBinding
	}
	if len(ccy.Global.IPFamilyPriority) == 0 {
		ccy.Global.IPFamilyPriority = []string{DefaultIPFamily}
	}

	// Create a single instance of VSphereInstance for the Global VCenterIP if the
	// VirtualCenter does not already exist in the map
	if ccy.Global.VCenterIP != "" && ccy.Vcenter[ccy.Global.VCenterIP] == nil {
		ccy.Vcenter[ccy.Global.VCenterIP] = &VirtualCenterConfigYAML{
			User:              ccy.Global.User,
			Password:          ccy.Global.Password,
			TenantRef:         ccy.Global.VCenterIP,
			VCenterIP:         ccy.Global.VCenterIP,
			VCenterPort:       ccy.Global.VCenterPort,
			InsecureFlag:      ccy.Global.InsecureFlag,
			Datacenters:       ccy.Global.Datacenters,
			RoundTripperCount: ccy.Global.RoundTripperCount,
			CAFile:            ccy.Global.CAFile,
			Thumbprint:        ccy.Global.Thumbprint,
			SecretRef:         DefaultCredentialManager,
			SecretName:        ccy.Global.SecretName,
			SecretNamespace:   ccy.Global.SecretNamespace,
			IPFamilyPriority:  ccy.Global.IPFamilyPriority,
		}
	}

	// Must have at least one vCenter defined
	if len(ccy.Vcenter) == 0 {
		klog.Error(ErrMissingVCenter)
		return ErrMissingVCenter
	}

	// vsphere.conf is no longer supported in the old format.
	for tenantRef, vcConfig := range ccy.Vcenter {
		klog.V(4).Infof("Initializing vc server %s", tenantRef)
		if vcConfig.VCenterIP == "" {
			klog.Error(ErrInvalidVCenterIP)
			return ErrInvalidVCenterIP
		}

		// in the YAML-based config, the tenant ref is required in the config
		vcConfig.TenantRef = tenantRef

		if !ccy.isSecretInfoProvided() && !vcConfig.isSecretInfoProvided() {
			if vcConfig.User == "" {
				vcConfig.User = ccy.Global.User
				if vcConfig.User == "" {
					klog.Errorf("vcConfig.User is empty for vc %s!", tenantRef)
					return ErrUsernameMissing
				}
			}
			if vcConfig.Password == "" {
				vcConfig.Password = ccy.Global.Password
				if vcConfig.Password == "" {
					klog.Errorf("vcConfig.Password is empty for vc %s!", tenantRef)
					return ErrPasswordMissing
				}
			}
		} else if ccy.isSecretInfoProvided() && !vcConfig.isSecretInfoProvided() {
			vcConfig.SecretRef = DefaultCredentialManager
		} else if vcConfig.isSecretInfoProvided() {
			vcConfig.SecretRef = vcConfig.SecretNamespace + "/" + vcConfig.SecretName
		}

		if vcConfig.VCenterPort == 0 {
			vcConfig.VCenterPort = ccy.Global.VCenterPort
		}

		if len(vcConfig.Datacenters) == 0 {
			if len(ccy.Global.Datacenters) != 0 {
				vcConfig.Datacenters = ccy.Global.Datacenters
			}
		}
		if vcConfig.RoundTripperCount == 0 {
			vcConfig.RoundTripperCount = ccy.Global.RoundTripperCount
		}
		if vcConfig.CAFile == "" {
			vcConfig.CAFile = ccy.Global.CAFile
		}
		if vcConfig.Thumbprint == "" {
			vcConfig.Thumbprint = ccy.Global.Thumbprint
		}

		if len(vcConfig.IPFamilyPriority) == 0 {
			vcConfig.IPFamilyPriority = ccy.Global.IPFamilyPriority
		}

		insecure := vcConfig.InsecureFlag
		if !insecure {
			vcConfig.InsecureFlag = ccy.Global.InsecureFlag
		}
	}

	return nil
}

// ReadRawConfigYAML parses vSphere cloud config file and stores it into ConfigYAML
func ReadRawConfigYAML(byConfig []byte) (*CommonConfigYAML, error) {
	if len(byConfig) == 0 {
		klog.Errorf("Invalid YAML file")
		return nil, fmt.Errorf("Invalid YAML file")
	}

	cfg := CommonConfigYAML{
		Vcenter: make(map[string]*VirtualCenterConfigYAML),
	}

	if err := yaml.Unmarshal(byConfig, &cfg); err != nil {
		klog.Errorf("Unmarshal failed: %s", err)
		return nil, err
	}

	err := cfg.validateConfig()
	if err != nil {
		klog.Errorf("validateConfig failed: %s", err)
		return nil, err
	}

	return &cfg, nil
}

// ReadConfigYAML parses vSphere cloud config file and stores it into Config
func ReadConfigYAML(byConfig []byte) (*Config, error) {
	cfg, err := ReadRawConfigYAML(byConfig)
	if err != nil {
		return nil, err
	}

	return cfg.CreateConfig(), nil
}
