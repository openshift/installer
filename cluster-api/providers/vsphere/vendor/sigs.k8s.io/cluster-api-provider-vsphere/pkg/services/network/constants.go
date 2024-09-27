/*
Copyright 2021 The Kubernetes Authors.

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

// Package network contains code for configuring network services.
package network

import (
	netopv1 "github.com/vmware-tanzu/net-operator-api/api/v1alpha1"
	nsxvpcv1 "github.com/vmware-tanzu/nsx-operator/pkg/apis/vpc/v1alpha1"
	ncpv1 "github.com/vmware-tanzu/vm-operator/external/ncp/api/v1alpha1"
)

const (
	// NSXTVNetSelectorKey is also defined in VM Operator.
	NSXTVNetSelectorKey = "ncp.vmware.com/virtual-network-name"

	// CAPVDefaultNetworkLabel is a label used to identify the default network.
	CAPVDefaultNetworkLabel = "capv.vmware.com/is-default-network"
	// NetOpNetworkNameAnnotation is the key used in an annotation to define the NetOp network. The expected value is the network name.
	NetOpNetworkNameAnnotation = "netoperator.vmware.com/network-name"

	// SystemNamespace is the namespace where supervisor control plane VMs reside.
	SystemNamespace = "kube-system"

	// legacyDefaultNetworkLabel was the label used for default networks.
	//
	// Deprecated: legacyDefaultNetworkLabel will be removed in a future release.
	legacyDefaultNetworkLabel = "capw.vmware.com/is-default-network"

	// AnnotationEnableEndpointHealthCheckKey is the key of the annotation that is used to enable health check on the
	// Service endpoint port. vm-operator propagates annotations in VMService to Service and LB providers like NSX-T
	// will enable health check on the endpoint target port when this annotation is present on the Service.
	AnnotationEnableEndpointHealthCheckKey = "lb.iaas.vmware.com/enable-endpoint-health-check"
)

var (
	// NetworkGVKNetOperator is the GVK used for networks in net-operator mode.
	NetworkGVKNetOperator = netopv1.SchemeGroupVersion.WithKind("Network")

	// NetworkGVKNSXT is the GVK used for networks in NSX-T mode.
	NetworkGVKNSXT = ncpv1.SchemeGroupVersion.WithKind("VirtualNetwork")

	// NetworkGVKNSXTVPC is the GVK used for networks in NSX-T VPC mode.
	NetworkGVKNSXTVPC = nsxvpcv1.SchemeGroupVersion.WithKind("SubnetSet")
)
