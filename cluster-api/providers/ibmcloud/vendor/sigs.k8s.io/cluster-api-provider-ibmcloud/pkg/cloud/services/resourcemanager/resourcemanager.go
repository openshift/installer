/*
Copyright 2024 The Kubernetes Authors.

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

package resourcemanager

import (
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/resourcemanagerv2"
)

//go:generate ../../../../hack/tools/bin/mockgen -source=./resourcemanager.go -destination=./mock/resourcemanager_generated.go -package=mock
//go:generate /usr/bin/env bash -c "cat ../../../../hack/boilerplate/boilerplate.generatego.txt ./mock/resourcemanager_generated.go > ./mock/_resourcemanager_generated.go && mv ./mock/_resourcemanager_generated.go ./mock/resourcemanager_generated.go"

// ResourceManager interface defines a method that a IBMCLOUD service object should implement in order to
// use the manage lifecycle of cloud resource groups using Resource Manager APIs.
type ResourceManager interface {
	GetResourceGroup(*resourcemanagerv2.GetResourceGroupOptions) (*resourcemanagerv2.ResourceGroup, *core.DetailedResponse, error)
	ListResourceGroups(*resourcemanagerv2.ListResourceGroupsOptions) (*resourcemanagerv2.ResourceGroupList, *core.DetailedResponse, error)

	GetResourceGroupByName(string) (*resourcemanagerv2.ResourceGroup, error)
}
