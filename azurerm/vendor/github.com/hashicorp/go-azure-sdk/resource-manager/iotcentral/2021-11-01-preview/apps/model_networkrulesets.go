package apps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkRuleSets struct {
	ApplyToDevices    *bool                   `json:"applyToDevices,omitempty"`
	ApplyToIoTCentral *bool                   `json:"applyToIoTCentral,omitempty"`
	DefaultAction     *NetworkAction          `json:"defaultAction,omitempty"`
	IPRules           *[]NetworkRuleSetIPRule `json:"ipRules,omitempty"`
}
