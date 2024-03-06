/*
Copyright 2022 The Kubernetes Authors.

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

package resourcecontroller

import (
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
)

// ResourceController interface defines a method that a IBMCLOUD service object should implement in order to
// use the resourcecontrollerv2 package for listing resource instances.
type ResourceController interface {
	ListResourceInstances(listResourceInstancesOptions *resourcecontrollerv2.ListResourceInstancesOptions) (result *resourcecontrollerv2.ResourceInstancesList, response *core.DetailedResponse, err error)
	GetResourceInstance(*resourcecontrollerv2.GetResourceInstanceOptions) (*resourcecontrollerv2.ResourceInstance, *core.DetailedResponse, error)
	CreateResourceInstance(*resourcecontrollerv2.CreateResourceInstanceOptions) (*resourcecontrollerv2.ResourceInstance, *core.DetailedResponse, error)
	GetServiceInstance(string, string) (*resourcecontrollerv2.ResourceInstance, error)
	DeleteResourceInstance(*resourcecontrollerv2.DeleteResourceInstanceOptions) (*core.DetailedResponse, error)

	GetInstanceByName(string, string, string) (*resourcecontrollerv2.ResourceInstance, error)
	CreateResourceKey(*resourcecontrollerv2.CreateResourceKeyOptions) (*resourcecontrollerv2.ResourceKey, *core.DetailedResponse, error)

	SetServiceURL(string) error
	GetServiceURL() string
}
