/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
   SPDX-License-Identifier: Apache-2.0 */

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
//+kubebuilder:object:root=true
//+kubebuilder:resource:scope="Cluster",path=ipblocksinfos

// IPBlocksInfo is the Schema for the ipblocksinfo API
type IPBlocksInfo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// ExternalIPCIDRs is a list of CIDR strings. Each CIDR is a contiguous IP address
	// spaces represented by network address and prefix length. The visibility of the
	// IPBlocks is External.
	ExternalIPCIDRs []string `json:"externalIPCIDRs,omitempty"`
	// PrivateTGWIPCIDRs is a list of CIDR strings. Each CIDR is a contiguous IP address
	// spaces represented by network address and prefix length. The visibility of the
	// IPBlocks is PrivateTWG. Only IPBlocks in default project will be included.
	PrivateTGWIPCIDRs []string `json:"privateTGWIPCIDRs,omitempty"`
}

//+kubebuilder:object:root=true

// IPBlocksInfoList contains a list of IPBlocksInfo
type IPBlocksInfoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IPBlocksInfo `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IPBlocksInfo{}, &IPBlocksInfoList{})
}
