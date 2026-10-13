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

package v1beta2

// Hub marks AWSCluster as a conversion hub.
func (*AWSCluster) Hub() {}

// Hub marks AWSClusterList as a conversion hub.
func (*AWSClusterList) Hub() {}

// Hub marks AWSMachine as a conversion hub.
func (*AWSMachine) Hub() {}

// Hub marks AWSMachineList as a conversion hub.
func (*AWSMachineList) Hub() {}

// Hub marks AWSMachineTemplate as a conversion hub.
func (*AWSMachineTemplate) Hub() {}

// Hub marks AWSMachineTemplateList as a conversion hub.
func (*AWSMachineTemplateList) Hub() {}

// Hub marks AWSClusterStaticIdentity as a conversion hub.
func (*AWSClusterStaticIdentity) Hub() {}

// Hub marks AWSClusterStaticIdentityList as a conversion hub.
func (*AWSClusterStaticIdentityList) Hub() {}

// Hub marks AWSClusterRoleIdentity as a conversion hub.
func (*AWSClusterRoleIdentity) Hub() {}

// Hub marks AWSClusterRoleIdentityList as a conversion hub.
func (*AWSClusterRoleIdentityList) Hub() {}

// Hub marks AWSClusterControllerIdentity as a conversion hub.
func (*AWSClusterControllerIdentity) Hub() {}

// Hub marks AWSClusterControllerIdentityList as a conversion hub.
func (*AWSClusterControllerIdentityList) Hub() {}

// Hub marks AWSClusterTemplate as a conversion hub.
func (*AWSClusterTemplate) Hub() {}

// Hub marks AWSClusterTemplateList as a conversion hub.
func (*AWSClusterTemplateList) Hub() {}
