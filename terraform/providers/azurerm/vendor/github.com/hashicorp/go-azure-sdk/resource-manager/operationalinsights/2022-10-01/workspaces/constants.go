package workspaces

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CapacityReservationLevel int64

const (
	CapacityReservationLevelFiveHundred  CapacityReservationLevel = 500
	CapacityReservationLevelFiveThousand CapacityReservationLevel = 5000
	CapacityReservationLevelFourHundred  CapacityReservationLevel = 400
	CapacityReservationLevelOneHundred   CapacityReservationLevel = 100
	CapacityReservationLevelOneThousand  CapacityReservationLevel = 1000
	CapacityReservationLevelThreeHundred CapacityReservationLevel = 300
	CapacityReservationLevelTwoHundred   CapacityReservationLevel = 200
	CapacityReservationLevelTwoThousand  CapacityReservationLevel = 2000
)

func PossibleValuesForCapacityReservationLevel() []int64 {
	return []int64{
		int64(CapacityReservationLevelFiveHundred),
		int64(CapacityReservationLevelFiveThousand),
		int64(CapacityReservationLevelFourHundred),
		int64(CapacityReservationLevelOneHundred),
		int64(CapacityReservationLevelOneThousand),
		int64(CapacityReservationLevelThreeHundred),
		int64(CapacityReservationLevelTwoHundred),
		int64(CapacityReservationLevelTwoThousand),
	}
}

func parseCapacityReservationLevel(input int64) (*CapacityReservationLevel, error) {
	vals := map[int64]CapacityReservationLevel{
		500:  CapacityReservationLevelFiveHundred,
		5000: CapacityReservationLevelFiveThousand,
		400:  CapacityReservationLevelFourHundred,
		100:  CapacityReservationLevelOneHundred,
		1000: CapacityReservationLevelOneThousand,
		300:  CapacityReservationLevelThreeHundred,
		200:  CapacityReservationLevelTwoHundred,
		2000: CapacityReservationLevelTwoThousand,
	}
	if v, ok := vals[input]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CapacityReservationLevel(input)
	return &out, nil
}

type DataIngestionStatus string

const (
	DataIngestionStatusApproachingQuota      DataIngestionStatus = "ApproachingQuota"
	DataIngestionStatusForceOff              DataIngestionStatus = "ForceOff"
	DataIngestionStatusForceOn               DataIngestionStatus = "ForceOn"
	DataIngestionStatusOverQuota             DataIngestionStatus = "OverQuota"
	DataIngestionStatusRespectQuota          DataIngestionStatus = "RespectQuota"
	DataIngestionStatusSubscriptionSuspended DataIngestionStatus = "SubscriptionSuspended"
)

func PossibleValuesForDataIngestionStatus() []string {
	return []string{
		string(DataIngestionStatusApproachingQuota),
		string(DataIngestionStatusForceOff),
		string(DataIngestionStatusForceOn),
		string(DataIngestionStatusOverQuota),
		string(DataIngestionStatusRespectQuota),
		string(DataIngestionStatusSubscriptionSuspended),
	}
}

func parseDataIngestionStatus(input string) (*DataIngestionStatus, error) {
	vals := map[string]DataIngestionStatus{
		"approachingquota":      DataIngestionStatusApproachingQuota,
		"forceoff":              DataIngestionStatusForceOff,
		"forceon":               DataIngestionStatusForceOn,
		"overquota":             DataIngestionStatusOverQuota,
		"respectquota":          DataIngestionStatusRespectQuota,
		"subscriptionsuspended": DataIngestionStatusSubscriptionSuspended,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataIngestionStatus(input)
	return &out, nil
}

type PublicNetworkAccessType string

const (
	PublicNetworkAccessTypeDisabled PublicNetworkAccessType = "Disabled"
	PublicNetworkAccessTypeEnabled  PublicNetworkAccessType = "Enabled"
)

func PossibleValuesForPublicNetworkAccessType() []string {
	return []string{
		string(PublicNetworkAccessTypeDisabled),
		string(PublicNetworkAccessTypeEnabled),
	}
}

func parsePublicNetworkAccessType(input string) (*PublicNetworkAccessType, error) {
	vals := map[string]PublicNetworkAccessType{
		"disabled": PublicNetworkAccessTypeDisabled,
		"enabled":  PublicNetworkAccessTypeEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicNetworkAccessType(input)
	return &out, nil
}

type WorkspaceEntityStatus string

const (
	WorkspaceEntityStatusCanceled            WorkspaceEntityStatus = "Canceled"
	WorkspaceEntityStatusCreating            WorkspaceEntityStatus = "Creating"
	WorkspaceEntityStatusDeleting            WorkspaceEntityStatus = "Deleting"
	WorkspaceEntityStatusFailed              WorkspaceEntityStatus = "Failed"
	WorkspaceEntityStatusProvisioningAccount WorkspaceEntityStatus = "ProvisioningAccount"
	WorkspaceEntityStatusSucceeded           WorkspaceEntityStatus = "Succeeded"
	WorkspaceEntityStatusUpdating            WorkspaceEntityStatus = "Updating"
)

func PossibleValuesForWorkspaceEntityStatus() []string {
	return []string{
		string(WorkspaceEntityStatusCanceled),
		string(WorkspaceEntityStatusCreating),
		string(WorkspaceEntityStatusDeleting),
		string(WorkspaceEntityStatusFailed),
		string(WorkspaceEntityStatusProvisioningAccount),
		string(WorkspaceEntityStatusSucceeded),
		string(WorkspaceEntityStatusUpdating),
	}
}

func parseWorkspaceEntityStatus(input string) (*WorkspaceEntityStatus, error) {
	vals := map[string]WorkspaceEntityStatus{
		"canceled":            WorkspaceEntityStatusCanceled,
		"creating":            WorkspaceEntityStatusCreating,
		"deleting":            WorkspaceEntityStatusDeleting,
		"failed":              WorkspaceEntityStatusFailed,
		"provisioningaccount": WorkspaceEntityStatusProvisioningAccount,
		"succeeded":           WorkspaceEntityStatusSucceeded,
		"updating":            WorkspaceEntityStatusUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WorkspaceEntityStatus(input)
	return &out, nil
}

type WorkspaceSkuNameEnum string

const (
	WorkspaceSkuNameEnumCapacityReservation  WorkspaceSkuNameEnum = "CapacityReservation"
	WorkspaceSkuNameEnumFree                 WorkspaceSkuNameEnum = "Free"
	WorkspaceSkuNameEnumLACluster            WorkspaceSkuNameEnum = "LACluster"
	WorkspaceSkuNameEnumPerGBTwoZeroOneEight WorkspaceSkuNameEnum = "PerGB2018"
	WorkspaceSkuNameEnumPerNode              WorkspaceSkuNameEnum = "PerNode"
	WorkspaceSkuNameEnumPremium              WorkspaceSkuNameEnum = "Premium"
	WorkspaceSkuNameEnumStandalone           WorkspaceSkuNameEnum = "Standalone"
	WorkspaceSkuNameEnumStandard             WorkspaceSkuNameEnum = "Standard"
)

func PossibleValuesForWorkspaceSkuNameEnum() []string {
	return []string{
		string(WorkspaceSkuNameEnumCapacityReservation),
		string(WorkspaceSkuNameEnumFree),
		string(WorkspaceSkuNameEnumLACluster),
		string(WorkspaceSkuNameEnumPerGBTwoZeroOneEight),
		string(WorkspaceSkuNameEnumPerNode),
		string(WorkspaceSkuNameEnumPremium),
		string(WorkspaceSkuNameEnumStandalone),
		string(WorkspaceSkuNameEnumStandard),
	}
}

func parseWorkspaceSkuNameEnum(input string) (*WorkspaceSkuNameEnum, error) {
	vals := map[string]WorkspaceSkuNameEnum{
		"capacityreservation": WorkspaceSkuNameEnumCapacityReservation,
		"free":                WorkspaceSkuNameEnumFree,
		"lacluster":           WorkspaceSkuNameEnumLACluster,
		"pergb2018":           WorkspaceSkuNameEnumPerGBTwoZeroOneEight,
		"pernode":             WorkspaceSkuNameEnumPerNode,
		"premium":             WorkspaceSkuNameEnumPremium,
		"standalone":          WorkspaceSkuNameEnumStandalone,
		"standard":            WorkspaceSkuNameEnumStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WorkspaceSkuNameEnum(input)
	return &out, nil
}
