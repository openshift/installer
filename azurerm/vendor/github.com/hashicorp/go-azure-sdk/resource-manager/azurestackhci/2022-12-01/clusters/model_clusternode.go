package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterNode struct {
	CoreCount                 *float64                   `json:"coreCount,omitempty"`
	EhcResourceId             *string                    `json:"ehcResourceId,omitempty"`
	Id                        *float64                   `json:"id,omitempty"`
	Manufacturer              *string                    `json:"manufacturer,omitempty"`
	MemoryInGiB               *float64                   `json:"memoryInGiB,omitempty"`
	Model                     *string                    `json:"model,omitempty"`
	Name                      *string                    `json:"name,omitempty"`
	NodeType                  *ClusterNodeType           `json:"nodeType,omitempty"`
	OsDisplayVersion          *string                    `json:"osDisplayVersion,omitempty"`
	OsName                    *string                    `json:"osName,omitempty"`
	OsVersion                 *string                    `json:"osVersion,omitempty"`
	SerialNumber              *string                    `json:"serialNumber,omitempty"`
	WindowsServerSubscription *WindowsServerSubscription `json:"windowsServerSubscription,omitempty"`
}
