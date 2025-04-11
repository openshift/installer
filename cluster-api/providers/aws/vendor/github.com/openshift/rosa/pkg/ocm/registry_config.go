/**
Copyright (c) 2024 Red Hat, Inc.

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

package ocm

import (
	"fmt"
	"strconv"
	"strings"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"

	"github.com/openshift/rosa/pkg/helper"
)

func BuildRegistryConfig(spec Spec) (*cmv1.ClusterRegistryConfigBuilder, error) {
	isClusterRegistryConfigured := spec.AllowedRegistries != nil ||
		spec.BlockedRegistries != nil || spec.InsecureRegistries != nil ||
		spec.AllowedRegistriesForImport != "" || spec.PlatformAllowlist != "" ||
		spec.AdditionalTrustedCa != nil

	if isClusterRegistryConfigured {
		registryResources := cmv1.NewRegistrySources()

		if spec.AllowedRegistries != nil {
			registryResources.AllowedRegistries(spec.AllowedRegistries...)
		}
		if spec.BlockedRegistries != nil {
			registryResources.BlockedRegistries(spec.BlockedRegistries...)
		}
		if spec.InsecureRegistries != nil {
			registryResources.InsecureRegistries(spec.InsecureRegistries...)
		}

		clusterRegistryConfig := cmv1.NewClusterRegistryConfig().
			RegistrySources(registryResources)

		if spec.AdditionalTrustedCa != nil {
			clusterRegistryConfig.AdditionalTrustedCa(spec.AdditionalTrustedCa)
		}
		if spec.PlatformAllowlist != "" {
			clusterRegistryConfig.PlatformAllowlist(
				cmv1.NewRegistryAllowlist().ID(spec.PlatformAllowlist))
		}
		if spec.AllowedRegistriesForImport != "" {
			obj, err := BuildAllowedRegistriesForImport(spec.AllowedRegistriesForImport)
			if err != nil {
				return nil, fmt.Errorf("Failed to build allowed registries for import, received error: %s", err)
			}

			//construct the location list
			var locationList []*cmv1.RegistryLocationBuilder
			for registryName, insecure := range obj {
				locationList = append(locationList, cmv1.NewRegistryLocation().
					Insecure(insecure).DomainName(registryName))
			}
			clusterRegistryConfig.AllowedRegistriesForImport(locationList...)
		}
		return clusterRegistryConfig, nil
	}
	return nil, nil
}

func BuildAllowedRegistriesForImport(allowedRegistriesForImport string) (map[string]bool, error) {
	obj := map[string]bool{}
	err := ValidateAllowedRegistriesForImport(allowedRegistriesForImport)
	if err != nil {
		return nil, err
	}
	list := helper.HandleEmptyStringOnSlice(strings.Split(allowedRegistriesForImport, ","))
	for _, registry := range list {
		registryObj := helper.HandleEmptyStringOnSlice(strings.Split(registry, ":"))
		objLen := len(registryObj)
		if objLen >= 2 {
			// insecure will always be the last item
			insecure := registryObj[objLen-1]
			// the registryName will remain the same
			registryName := strings.Join(registryObj[:objLen-1], ":")
			boolValue, err := strconv.ParseBool(insecure)
			if err != nil {
				return nil, fmt.Errorf("failed to convert the value '%s' to bool: %s", insecure, err)
			}
			obj[registryName] = boolValue
		}
	}
	return obj, nil
}

func (c *Client) GetAllowlist(id string) (*cmv1.RegistryAllowlist, error) {
	response, err := c.ocm.ClustersMgmt().V1().RegistryAllowlists().
		RegistryAllowlist(id).Get().
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}

	return response.Body(), nil
}
