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

	ini "gopkg.in/gcfg.v1"
	klog "k8s.io/klog/v2"
)

/*
	TODO:
	When the INI based cloud-config is deprecated. This file should be deleted.
*/

// CreateConfig generates a common Config object based on what other structs and funcs
// are already dependent upon in other packages.
func (cci *CommonConfigINI) CreateConfig() *Config {
	cfg := &Config{
		VirtualCenter: make(map[string]*VirtualCenterConfig),
	}

	cfg.Global.User = cci.Global.User
	cfg.Global.Password = cci.Global.Password
	cfg.Global.VCenterIP = cci.Global.VCenterIP
	cfg.Global.VCenterPort = cci.Global.VCenterPort
	cfg.Global.InsecureFlag = cci.Global.InsecureFlag
	cfg.Global.Datacenters = cci.Global.Datacenters
	cfg.Global.RoundTripperCount = cci.Global.RoundTripperCount
	cfg.Global.CAFile = cci.Global.CAFile
	cfg.Global.Thumbprint = cci.Global.Thumbprint
	cfg.Global.SecretName = cci.Global.SecretName
	cfg.Global.SecretNamespace = cci.Global.SecretNamespace
	cfg.Global.SecretsDirectory = cci.Global.SecretsDirectory
	cfg.Global.APIDisable = cci.Global.APIDisable
	cfg.Global.APIBinding = cci.Global.APIBinding

	for keyVcConfig, valVcConfig := range cci.VirtualCenter {
		cfg.VirtualCenter[keyVcConfig] = &VirtualCenterConfig{
			User:              valVcConfig.User,
			Password:          valVcConfig.Password,
			TenantRef:         valVcConfig.TenantRef,
			VCenterIP:         valVcConfig.VCenterIP,
			VCenterPort:       valVcConfig.VCenterPort,
			InsecureFlag:      valVcConfig.InsecureFlag,
			Datacenters:       valVcConfig.Datacenters,
			RoundTripperCount: valVcConfig.RoundTripperCount,
			CAFile:            valVcConfig.CAFile,
			Thumbprint:        valVcConfig.Thumbprint,
			SecretRef:         valVcConfig.SecretRef,
			SecretName:        valVcConfig.SecretName,
			SecretNamespace:   valVcConfig.SecretNamespace,
			IPFamilyPriority:  valVcConfig.IPFamilyPriority,
		}
	}

	cfg.Labels.Region = cci.Labels.Region
	cfg.Labels.Zone = cci.Labels.Zone

	return cfg
}

// validateIPFamily takes the possible values of IPFamily and initializes the
// slice as determined bby priority
func (vcci *VirtualCenterConfigINI) validateIPFamily() error {
	if len(vcci.IPFamily) == 0 {
		vcci.IPFamily = DefaultIPFamily
	}

	ipFamilies := strings.Split(vcci.IPFamily, ",")
	for i, ipFamily := range ipFamilies {
		ipFamily = strings.TrimSpace(ipFamily)
		if len(ipFamily) == 0 {
			copy(ipFamilies[i:], ipFamilies[i+1:])      // Shift a[i+1:] left one index.
			ipFamilies[len(ipFamilies)-1] = ""          // Erase last element (write zero value).
			ipFamilies = ipFamilies[:len(ipFamilies)-1] // Truncate slice.
			continue
		}
		if !strings.EqualFold(ipFamily, IPv4Family) && !strings.EqualFold(ipFamily, IPv6Family) {
			return ErrInvalidIPFamilyType
		}
	}

	vcci.IPFamilyPriority = ipFamilies
	return nil
}

// isSecretInfoProvided returns true if k8s secret is set or using generic CO secret method.
// If both k8s secret and generic CO both are true, we don't know which to use, so return false.
func (cci *CommonConfigINI) isSecretInfoProvided() bool {
	return (cci.Global.SecretName != "" && cci.Global.SecretNamespace != "" && cci.Global.SecretsDirectory == "") ||
		(cci.Global.SecretName == "" && cci.Global.SecretNamespace == "" && cci.Global.SecretsDirectory != "")
}

// isSecretInfoProvided returns true if the secret per VC has been configured
func (vcci *VirtualCenterConfigINI) isSecretInfoProvided() bool {
	return vcci.SecretName != "" && vcci.SecretNamespace != ""
}

func (cci *CommonConfigINI) validateConfig() error {
	//Fix default global values
	if cci.Global.RoundTripperCount == 0 {
		cci.Global.RoundTripperCount = DefaultRoundTripperCount
	}
	if cci.Global.VCenterPort == "" {
		cci.Global.VCenterPort = DefaultVCenterPortStr
	}
	if cci.Global.APIBinding == "" {
		cci.Global.APIBinding = DefaultAPIBinding
	}
	if cci.Global.IPFamily == "" {
		cci.Global.IPFamily = DefaultIPFamily
	}

	// Create a single instance of VSphereInstance for the Global VCenterIP if the
	// VirtualCenter does not already exist in the map
	if cci.Global.VCenterIP != "" && cci.VirtualCenter[cci.Global.VCenterIP] == nil {
		cci.VirtualCenter[cci.Global.VCenterIP] = &VirtualCenterConfigINI{
			User:              cci.Global.User,
			Password:          cci.Global.Password,
			TenantRef:         cci.Global.VCenterIP,
			VCenterIP:         cci.Global.VCenterIP,
			VCenterPort:       cci.Global.VCenterPort,
			InsecureFlag:      cci.Global.InsecureFlag,
			Datacenters:       cci.Global.Datacenters,
			RoundTripperCount: cci.Global.RoundTripperCount,
			CAFile:            cci.Global.CAFile,
			Thumbprint:        cci.Global.Thumbprint,
			SecretRef:         DefaultCredentialManager,
			SecretName:        cci.Global.SecretName,
			SecretNamespace:   cci.Global.SecretNamespace,
			IPFamily:          cci.Global.IPFamily,
		}
	}

	// Must have at least one vCenter defined
	if len(cci.VirtualCenter) == 0 {
		klog.Error(ErrMissingVCenter)
		return ErrMissingVCenter
	}

	// vsphere.conf is no longer supported in the old format.
	for vcServer, vcConfig := range cci.VirtualCenter {
		klog.V(4).Infof("Initializing vc server %s", vcServer)
		if vcServer == "" {
			klog.Error(ErrInvalidVCenterIP)
			return ErrInvalidVCenterIP
		}

		// If vcConfig.VCenterIP is explicitly set, that means the vcServer
		// above is the TenantRef
		if vcConfig.VCenterIP != "" {
			//vcConfig.VCenterIP is already set
			vcConfig.TenantRef = vcServer
		} else {
			vcConfig.VCenterIP = vcServer
			vcConfig.TenantRef = vcServer
		}

		if !cci.isSecretInfoProvided() && !vcConfig.isSecretInfoProvided() {
			if vcConfig.User == "" {
				vcConfig.User = cci.Global.User
				if vcConfig.User == "" {
					klog.Errorf("vcConfig.User is empty for vc %s!", vcServer)
					return ErrUsernameMissing
				}
			}
			if vcConfig.Password == "" {
				vcConfig.Password = cci.Global.Password
				if vcConfig.Password == "" {
					klog.Errorf("vcConfig.Password is empty for vc %s!", vcServer)
					return ErrPasswordMissing
				}
			}
		} else if cci.isSecretInfoProvided() && !vcConfig.isSecretInfoProvided() {
			vcConfig.SecretRef = DefaultCredentialManager
		} else if vcConfig.isSecretInfoProvided() {
			vcConfig.SecretRef = vcConfig.SecretNamespace + "/" + vcConfig.SecretName
		}

		if vcConfig.VCenterPort == "" {
			vcConfig.VCenterPort = cci.Global.VCenterPort
		}

		if vcConfig.Datacenters == "" {
			if cci.Global.Datacenters != "" {
				vcConfig.Datacenters = cci.Global.Datacenters
			}
		}
		if vcConfig.RoundTripperCount == 0 {
			vcConfig.RoundTripperCount = cci.Global.RoundTripperCount
		}
		if vcConfig.CAFile == "" {
			vcConfig.CAFile = cci.Global.CAFile
		}
		if vcConfig.Thumbprint == "" {
			vcConfig.Thumbprint = cci.Global.Thumbprint
		}

		if vcConfig.IPFamily == "" {
			vcConfig.IPFamily = cci.Global.IPFamily
		}

		err := vcConfig.validateIPFamily()
		if err != nil {
			klog.Errorf("Invalid vcConfig IPFamily: %s, err=%s", vcConfig.IPFamily, err)
			return err
		}

		insecure := vcConfig.InsecureFlag
		if !insecure {
			vcConfig.InsecureFlag = cci.Global.InsecureFlag
		}
	}

	return nil
}

// ReadRawConfigINI parses vSphere cloud config file and stores it into ConfigINI
func ReadRawConfigINI(byConfig []byte) (*CommonConfigINI, error) {
	if len(byConfig) == 0 {
		return nil, fmt.Errorf("Invalid INI file")
	}

	strConfig := string(byConfig[:])

	cfg := &CommonConfigINI{
		VirtualCenter: make(map[string]*VirtualCenterConfigINI),
	}

	if err := ini.FatalOnly(ini.ReadStringInto(cfg, strConfig)); err != nil {
		return nil, err
	}

	err := cfg.validateConfig()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// ReadConfigINI parses vSphere cloud config file and stores it into Config
func ReadConfigINI(byConfig []byte) (*Config, error) {
	cfg, err := ReadRawConfigINI(byConfig)
	if err != nil {
		return nil, err
	}

	return cfg.CreateConfig(), nil
}
