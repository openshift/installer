package redis

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DayOfWeek string

const (
	DayOfWeekEveryday  DayOfWeek = "Everyday"
	DayOfWeekFriday    DayOfWeek = "Friday"
	DayOfWeekMonday    DayOfWeek = "Monday"
	DayOfWeekSaturday  DayOfWeek = "Saturday"
	DayOfWeekSunday    DayOfWeek = "Sunday"
	DayOfWeekThursday  DayOfWeek = "Thursday"
	DayOfWeekTuesday   DayOfWeek = "Tuesday"
	DayOfWeekWednesday DayOfWeek = "Wednesday"
	DayOfWeekWeekend   DayOfWeek = "Weekend"
)

func PossibleValuesForDayOfWeek() []string {
	return []string{
		string(DayOfWeekEveryday),
		string(DayOfWeekFriday),
		string(DayOfWeekMonday),
		string(DayOfWeekSaturday),
		string(DayOfWeekSunday),
		string(DayOfWeekThursday),
		string(DayOfWeekTuesday),
		string(DayOfWeekWednesday),
		string(DayOfWeekWeekend),
	}
}

func parseDayOfWeek(input string) (*DayOfWeek, error) {
	vals := map[string]DayOfWeek{
		"everyday":  DayOfWeekEveryday,
		"friday":    DayOfWeekFriday,
		"monday":    DayOfWeekMonday,
		"saturday":  DayOfWeekSaturday,
		"sunday":    DayOfWeekSunday,
		"thursday":  DayOfWeekThursday,
		"tuesday":   DayOfWeekTuesday,
		"wednesday": DayOfWeekWednesday,
		"weekend":   DayOfWeekWeekend,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DayOfWeek(input)
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
	ProvisioningStateCreating               ProvisioningState = "Creating"
	ProvisioningStateDeleting               ProvisioningState = "Deleting"
	ProvisioningStateDisabled               ProvisioningState = "Disabled"
	ProvisioningStateFailed                 ProvisioningState = "Failed"
	ProvisioningStateLinking                ProvisioningState = "Linking"
	ProvisioningStateProvisioning           ProvisioningState = "Provisioning"
	ProvisioningStateRecoveringScaleFailure ProvisioningState = "RecoveringScaleFailure"
	ProvisioningStateScaling                ProvisioningState = "Scaling"
	ProvisioningStateSucceeded              ProvisioningState = "Succeeded"
	ProvisioningStateUnlinking              ProvisioningState = "Unlinking"
	ProvisioningStateUnprovisioning         ProvisioningState = "Unprovisioning"
	ProvisioningStateUpdating               ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateDisabled),
		string(ProvisioningStateFailed),
		string(ProvisioningStateLinking),
		string(ProvisioningStateProvisioning),
		string(ProvisioningStateRecoveringScaleFailure),
		string(ProvisioningStateScaling),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUnlinking),
		string(ProvisioningStateUnprovisioning),
		string(ProvisioningStateUpdating),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"creating":               ProvisioningStateCreating,
		"deleting":               ProvisioningStateDeleting,
		"disabled":               ProvisioningStateDisabled,
		"failed":                 ProvisioningStateFailed,
		"linking":                ProvisioningStateLinking,
		"provisioning":           ProvisioningStateProvisioning,
		"recoveringscalefailure": ProvisioningStateRecoveringScaleFailure,
		"scaling":                ProvisioningStateScaling,
		"succeeded":              ProvisioningStateSucceeded,
		"unlinking":              ProvisioningStateUnlinking,
		"unprovisioning":         ProvisioningStateUnprovisioning,
		"updating":               ProvisioningStateUpdating,
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

type RebootType string

const (
	RebootTypeAllNodes      RebootType = "AllNodes"
	RebootTypePrimaryNode   RebootType = "PrimaryNode"
	RebootTypeSecondaryNode RebootType = "SecondaryNode"
)

func PossibleValuesForRebootType() []string {
	return []string{
		string(RebootTypeAllNodes),
		string(RebootTypePrimaryNode),
		string(RebootTypeSecondaryNode),
	}
}

func parseRebootType(input string) (*RebootType, error) {
	vals := map[string]RebootType{
		"allnodes":      RebootTypeAllNodes,
		"primarynode":   RebootTypePrimaryNode,
		"secondarynode": RebootTypeSecondaryNode,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RebootType(input)
	return &out, nil
}

type RedisKeyType string

const (
	RedisKeyTypePrimary   RedisKeyType = "Primary"
	RedisKeyTypeSecondary RedisKeyType = "Secondary"
)

func PossibleValuesForRedisKeyType() []string {
	return []string{
		string(RedisKeyTypePrimary),
		string(RedisKeyTypeSecondary),
	}
}

func parseRedisKeyType(input string) (*RedisKeyType, error) {
	vals := map[string]RedisKeyType{
		"primary":   RedisKeyTypePrimary,
		"secondary": RedisKeyTypeSecondary,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RedisKeyType(input)
	return &out, nil
}

type ReplicationRole string

const (
	ReplicationRolePrimary   ReplicationRole = "Primary"
	ReplicationRoleSecondary ReplicationRole = "Secondary"
)

func PossibleValuesForReplicationRole() []string {
	return []string{
		string(ReplicationRolePrimary),
		string(ReplicationRoleSecondary),
	}
}

func parseReplicationRole(input string) (*ReplicationRole, error) {
	vals := map[string]ReplicationRole{
		"primary":   ReplicationRolePrimary,
		"secondary": ReplicationRoleSecondary,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReplicationRole(input)
	return &out, nil
}

type SkuFamily string

const (
	SkuFamilyC SkuFamily = "C"
	SkuFamilyP SkuFamily = "P"
)

func PossibleValuesForSkuFamily() []string {
	return []string{
		string(SkuFamilyC),
		string(SkuFamilyP),
	}
}

func parseSkuFamily(input string) (*SkuFamily, error) {
	vals := map[string]SkuFamily{
		"c": SkuFamilyC,
		"p": SkuFamilyP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuFamily(input)
	return &out, nil
}

type SkuName string

const (
	SkuNameBasic    SkuName = "Basic"
	SkuNamePremium  SkuName = "Premium"
	SkuNameStandard SkuName = "Standard"
)

func PossibleValuesForSkuName() []string {
	return []string{
		string(SkuNameBasic),
		string(SkuNamePremium),
		string(SkuNameStandard),
	}
}

func parseSkuName(input string) (*SkuName, error) {
	vals := map[string]SkuName{
		"basic":    SkuNameBasic,
		"premium":  SkuNamePremium,
		"standard": SkuNameStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuName(input)
	return &out, nil
}

type TlsVersion string

const (
	TlsVersionOnePointOne  TlsVersion = "1.1"
	TlsVersionOnePointTwo  TlsVersion = "1.2"
	TlsVersionOnePointZero TlsVersion = "1.0"
)

func PossibleValuesForTlsVersion() []string {
	return []string{
		string(TlsVersionOnePointOne),
		string(TlsVersionOnePointTwo),
		string(TlsVersionOnePointZero),
	}
}

func parseTlsVersion(input string) (*TlsVersion, error) {
	vals := map[string]TlsVersion{
		"1.1": TlsVersionOnePointOne,
		"1.2": TlsVersionOnePointTwo,
		"1.0": TlsVersionOnePointZero,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TlsVersion(input)
	return &out, nil
}
