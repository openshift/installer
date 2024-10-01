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

// Package constants contains constants used in CAPV.
package constants

import (
	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
)

const (
	// CloudProviderSecretName is the name of the Secret that stores the
	// cloud provider credentials.
	CloudProviderSecretName = "cloud-provider-vsphere-credentials"

	// CloudProviderSecretNamespace is the namespace in which the cloud provider
	// credentials secret is located.
	CloudProviderSecretNamespace = "kube-system"

	// DefaultBindPort is the default API port used to generate the kubeadm
	// configurations.
	DefaultBindPort = 6443

	// VSphereCredentialSecretUserKey is the key used to store/retrieve the
	// vSphere username from a Kubernetes secret.
	VSphereCredentialSecretUserKey = "username"

	// VSphereCredentialSecretPassKey is the key used to store/retrieve the
	// vSphere password from a Kubernetes secret.
	VSphereCredentialSecretPassKey = "password"

	// MachineReadyAnnotationLabel is the annotation used to indicate that a
	// machine is ready.
	MachineReadyAnnotationLabel = "capv." + infrav1.GroupName + "/machine-ready"

	// MaintenanceAnnotationLabel is the annotation used to indicate a machine and/or
	// cluster are in maintenance mode.
	MaintenanceAnnotationLabel = "capv." + infrav1.GroupName + "/maintenance"

	// NodeLabelPrefix is the prefix for node labels.
	NodeLabelPrefix = "node.cluster.x-k8s.io"

	// ESXiHostInfoLabel is the label for esxi host info.
	ESXiHostInfoLabel = NodeLabelPrefix + "/esxi-host"
)
