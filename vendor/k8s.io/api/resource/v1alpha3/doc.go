/*
<<<<<<<< HEAD:vendor/k8s.io/api/resource/v1alpha3/doc.go
Copyright 2022 The Kubernetes Authors.
========
Copyright 2016 The Kubernetes Authors.
>>>>>>>> 4e32654a50 (Changed logic to use CAPI for data disks):cluster-api/providers/vsphere/vendor/k8s.io/apimachinery/pkg/util/portforward/constants.go

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

<<<<<<<< HEAD:vendor/k8s.io/api/resource/v1alpha3/doc.go
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=package
// +k8s:protobuf-gen=package

// +groupName=resource.k8s.io

// Package v1alpha3 is the v1alpha3 version of the resource API.
package v1alpha3 // import "k8s.io/api/resource/v1alpha3"
========
package portforward

const (
	PortForwardV1Name                    = "portforward.k8s.io"
	WebsocketsSPDYTunnelingPrefix        = "SPDY/3.1+"
	KubernetesSuffix                     = ".k8s.io"
	WebsocketsSPDYTunnelingPortForwardV1 = WebsocketsSPDYTunnelingPrefix + PortForwardV1Name
)
>>>>>>>> 4e32654a50 (Changed logic to use CAPI for data disks):cluster-api/providers/vsphere/vendor/k8s.io/apimachinery/pkg/util/portforward/constants.go
