package deviceupdates

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthenticationType string

const (
	AuthenticationTypeKeyBased AuthenticationType = "KeyBased"
)

func PossibleValuesForAuthenticationType() []string {
	return []string{
		string(AuthenticationTypeKeyBased),
	}
}

func parseAuthenticationType(input string) (*AuthenticationType, error) {
	vals := map[string]AuthenticationType{
		"keybased": AuthenticationTypeKeyBased,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AuthenticationType(input)
	return &out, nil
}

type CheckNameAvailabilityReason string

const (
	CheckNameAvailabilityReasonAlreadyExists CheckNameAvailabilityReason = "AlreadyExists"
	CheckNameAvailabilityReasonInvalid       CheckNameAvailabilityReason = "Invalid"
)

func PossibleValuesForCheckNameAvailabilityReason() []string {
	return []string{
		string(CheckNameAvailabilityReasonAlreadyExists),
		string(CheckNameAvailabilityReasonInvalid),
	}
}

func parseCheckNameAvailabilityReason(input string) (*CheckNameAvailabilityReason, error) {
	vals := map[string]CheckNameAvailabilityReason{
		"alreadyexists": CheckNameAvailabilityReasonAlreadyExists,
		"invalid":       CheckNameAvailabilityReasonInvalid,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CheckNameAvailabilityReason(input)
	return &out, nil
}

type PrivateEndpointConnectionProvisioningState string

const (
	PrivateEndpointConnectionProvisioningStateCreating  PrivateEndpointConnectionProvisioningState = "Creating"
	PrivateEndpointConnectionProvisioningStateDeleting  PrivateEndpointConnectionProvisioningState = "Deleting"
	PrivateEndpointConnectionProvisioningStateFailed    PrivateEndpointConnectionProvisioningState = "Failed"
	PrivateEndpointConnectionProvisioningStateSucceeded PrivateEndpointConnectionProvisioningState = "Succeeded"
)

func PossibleValuesForPrivateEndpointConnectionProvisioningState() []string {
	return []string{
		string(PrivateEndpointConnectionProvisioningStateCreating),
		string(PrivateEndpointConnectionProvisioningStateDeleting),
		string(PrivateEndpointConnectionProvisioningStateFailed),
		string(PrivateEndpointConnectionProvisioningStateSucceeded),
	}
}

func parsePrivateEndpointConnectionProvisioningState(input string) (*PrivateEndpointConnectionProvisioningState, error) {
	vals := map[string]PrivateEndpointConnectionProvisioningState{
		"creating":  PrivateEndpointConnectionProvisioningStateCreating,
		"deleting":  PrivateEndpointConnectionProvisioningStateDeleting,
		"failed":    PrivateEndpointConnectionProvisioningStateFailed,
		"succeeded": PrivateEndpointConnectionProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateEndpointConnectionProvisioningState(input)
	return &out, nil
}

type PrivateEndpointConnectionProxyProvisioningState string

const (
	PrivateEndpointConnectionProxyProvisioningStateCreating  PrivateEndpointConnectionProxyProvisioningState = "Creating"
	PrivateEndpointConnectionProxyProvisioningStateDeleting  PrivateEndpointConnectionProxyProvisioningState = "Deleting"
	PrivateEndpointConnectionProxyProvisioningStateFailed    PrivateEndpointConnectionProxyProvisioningState = "Failed"
	PrivateEndpointConnectionProxyProvisioningStateSucceeded PrivateEndpointConnectionProxyProvisioningState = "Succeeded"
)

func PossibleValuesForPrivateEndpointConnectionProxyProvisioningState() []string {
	return []string{
		string(PrivateEndpointConnectionProxyProvisioningStateCreating),
		string(PrivateEndpointConnectionProxyProvisioningStateDeleting),
		string(PrivateEndpointConnectionProxyProvisioningStateFailed),
		string(PrivateEndpointConnectionProxyProvisioningStateSucceeded),
	}
}

func parsePrivateEndpointConnectionProxyProvisioningState(input string) (*PrivateEndpointConnectionProxyProvisioningState, error) {
	vals := map[string]PrivateEndpointConnectionProxyProvisioningState{
		"creating":  PrivateEndpointConnectionProxyProvisioningStateCreating,
		"deleting":  PrivateEndpointConnectionProxyProvisioningStateDeleting,
		"failed":    PrivateEndpointConnectionProxyProvisioningStateFailed,
		"succeeded": PrivateEndpointConnectionProxyProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateEndpointConnectionProxyProvisioningState(input)
	return &out, nil
}

type PrivateEndpointServiceConnectionStatus string

const (
	PrivateEndpointServiceConnectionStatusApproved PrivateEndpointServiceConnectionStatus = "Approved"
	PrivateEndpointServiceConnectionStatusPending  PrivateEndpointServiceConnectionStatus = "Pending"
	PrivateEndpointServiceConnectionStatusRejected PrivateEndpointServiceConnectionStatus = "Rejected"
)

func PossibleValuesForPrivateEndpointServiceConnectionStatus() []string {
	return []string{
		string(PrivateEndpointServiceConnectionStatusApproved),
		string(PrivateEndpointServiceConnectionStatusPending),
		string(PrivateEndpointServiceConnectionStatusRejected),
	}
}

func parsePrivateEndpointServiceConnectionStatus(input string) (*PrivateEndpointServiceConnectionStatus, error) {
	vals := map[string]PrivateEndpointServiceConnectionStatus{
		"approved": PrivateEndpointServiceConnectionStatusApproved,
		"pending":  PrivateEndpointServiceConnectionStatusPending,
		"rejected": PrivateEndpointServiceConnectionStatusRejected,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateEndpointServiceConnectionStatus(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted  ProvisioningState = "Accepted"
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleted   ProvisioningState = "Deleted"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleted),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"accepted":  ProvisioningStateAccepted,
		"canceled":  ProvisioningStateCanceled,
		"creating":  ProvisioningStateCreating,
		"deleted":   ProvisioningStateDeleted,
		"failed":    ProvisioningStateFailed,
		"succeeded": ProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type PublicNetworkAccess string

const (
	PublicNetworkAccessDisabled PublicNetworkAccess = "Disabled"
	PublicNetworkAccessEnabled  PublicNetworkAccess = "Enabled"
)

func PossibleValuesForPublicNetworkAccess() []string {
	return []string{
		string(PublicNetworkAccessDisabled),
		string(PublicNetworkAccessEnabled),
	}
}

func parsePublicNetworkAccess(input string) (*PublicNetworkAccess, error) {
	vals := map[string]PublicNetworkAccess{
		"disabled": PublicNetworkAccessDisabled,
		"enabled":  PublicNetworkAccessEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicNetworkAccess(input)
	return &out, nil
}

type Role string

const (
	RoleFailover Role = "Failover"
	RolePrimary  Role = "Primary"
)

func PossibleValuesForRole() []string {
	return []string{
		string(RoleFailover),
		string(RolePrimary),
	}
}

func parseRole(input string) (*Role, error) {
	vals := map[string]Role{
		"failover": RoleFailover,
		"primary":  RolePrimary,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Role(input)
	return &out, nil
}

type SKU string

const (
	SKUFree     SKU = "Free"
	SKUStandard SKU = "Standard"
)

func PossibleValuesForSKU() []string {
	return []string{
		string(SKUFree),
		string(SKUStandard),
	}
}

func parseSKU(input string) (*SKU, error) {
	vals := map[string]SKU{
		"free":     SKUFree,
		"standard": SKUStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SKU(input)
	return &out, nil
}
