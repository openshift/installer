/*
Copyright 2021.

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

// Package v1beta1 contains API Schema definitions for the ibmcloudproviderconfig v1beta1 API group
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=package,register
// +k8s:conversion-gen=github.com/openshift/cluster-api-provider-ibmcloud/pkg/apis/ibmcloudprovider
// +k8s:defaulter-gen=TypeMeta
// +groupName=ibmcloudproviderconfig.openshift.io
package v1beta1

import (
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	klog "k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
	yaml "sigs.k8s.io/yaml"
)

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: "ibmcloudproviderconfig.openshift.io", Version: "v1beta1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}
)

// RawExtensionFromProviderSpec marshals the machine provider spec.
func RawExtensionFromProviderSpec(spec *IBMCloudMachineProviderSpec) (*runtime.RawExtension, error) {
	if spec == nil {
		return &runtime.RawExtension{}, nil
	}

	var rawBytes []byte
	var err error
	if rawBytes, err = json.Marshal(spec); err != nil {
		return nil, fmt.Errorf("error marshalling providerSpec: %v", err)
	}

	return &runtime.RawExtension{
		Raw: rawBytes,
	}, nil
}

// RawExtensionFromProviderStatus marshals the provider status
func RawExtensionFromProviderStatus(status *IBMCloudMachineProviderStatus) (*runtime.RawExtension, error) {
	if status == nil {
		return &runtime.RawExtension{}, nil
	}

	var rawBytes []byte
	var err error
	if rawBytes, err = json.Marshal(status); err != nil {
		return nil, fmt.Errorf("error marshalling providerStatus: %v", err)
	}

	return &runtime.RawExtension{
		Raw: rawBytes,
	}, nil
}

// ProviderSpecFromRawExtension unmarshals the JSON-encoded spec
func ProviderSpecFromRawExtension(rawExtension *runtime.RawExtension) (*IBMCloudMachineProviderSpec, error) {
	if rawExtension == nil {
		return &IBMCloudMachineProviderSpec{}, nil
	}

	spec := new(IBMCloudMachineProviderSpec)
	if err := yaml.Unmarshal(rawExtension.Raw, &spec); err != nil {
		return nil, fmt.Errorf("error unmarshalling providerSpec: %v", err)
	}

	klog.V(5).Infof("Got provider Spec from raw extension: %+v", spec)
	return spec, nil
}

// ProviderStatusFromRawExtension unmarshals the JSON-encoded into a IBMCloudMachineProviderStatus type
func ProviderStatusFromRawExtension(rawExtension *runtime.RawExtension) (*IBMCloudMachineProviderStatus, error) {
	if rawExtension == nil {
		return &IBMCloudMachineProviderStatus{}, nil
	}

	providerStatus := new(IBMCloudMachineProviderStatus)
	if err := yaml.Unmarshal(rawExtension.Raw, providerStatus); err != nil {
		return nil, fmt.Errorf("error unmarshalling providerStatus: %v", err)
	}

	klog.V(5).Infof("Got provider Status from raw extension: %+v", providerStatus)
	return providerStatus, nil
}
