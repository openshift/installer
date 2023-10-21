package pool

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AllocationState string

const (
	AllocationStateResizing AllocationState = "Resizing"
	AllocationStateSteady   AllocationState = "Steady"
	AllocationStateStopping AllocationState = "Stopping"
)

func PossibleValuesForAllocationState() []string {
	return []string{
		string(AllocationStateResizing),
		string(AllocationStateSteady),
		string(AllocationStateStopping),
	}
}

func parseAllocationState(input string) (*AllocationState, error) {
	vals := map[string]AllocationState{
		"resizing": AllocationStateResizing,
		"steady":   AllocationStateSteady,
		"stopping": AllocationStateStopping,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AllocationState(input)
	return &out, nil
}

type AutoUserScope string

const (
	AutoUserScopePool AutoUserScope = "Pool"
	AutoUserScopeTask AutoUserScope = "Task"
)

func PossibleValuesForAutoUserScope() []string {
	return []string{
		string(AutoUserScopePool),
		string(AutoUserScopeTask),
	}
}

func parseAutoUserScope(input string) (*AutoUserScope, error) {
	vals := map[string]AutoUserScope{
		"pool": AutoUserScopePool,
		"task": AutoUserScopeTask,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutoUserScope(input)
	return &out, nil
}

type CachingType string

const (
	CachingTypeNone      CachingType = "None"
	CachingTypeReadOnly  CachingType = "ReadOnly"
	CachingTypeReadWrite CachingType = "ReadWrite"
)

func PossibleValuesForCachingType() []string {
	return []string{
		string(CachingTypeNone),
		string(CachingTypeReadOnly),
		string(CachingTypeReadWrite),
	}
}

func parseCachingType(input string) (*CachingType, error) {
	vals := map[string]CachingType{
		"none":      CachingTypeNone,
		"readonly":  CachingTypeReadOnly,
		"readwrite": CachingTypeReadWrite,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CachingType(input)
	return &out, nil
}

type CertificateStoreLocation string

const (
	CertificateStoreLocationCurrentUser  CertificateStoreLocation = "CurrentUser"
	CertificateStoreLocationLocalMachine CertificateStoreLocation = "LocalMachine"
)

func PossibleValuesForCertificateStoreLocation() []string {
	return []string{
		string(CertificateStoreLocationCurrentUser),
		string(CertificateStoreLocationLocalMachine),
	}
}

func parseCertificateStoreLocation(input string) (*CertificateStoreLocation, error) {
	vals := map[string]CertificateStoreLocation{
		"currentuser":  CertificateStoreLocationCurrentUser,
		"localmachine": CertificateStoreLocationLocalMachine,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CertificateStoreLocation(input)
	return &out, nil
}

type CertificateVisibility string

const (
	CertificateVisibilityRemoteUser CertificateVisibility = "RemoteUser"
	CertificateVisibilityStartTask  CertificateVisibility = "StartTask"
	CertificateVisibilityTask       CertificateVisibility = "Task"
)

func PossibleValuesForCertificateVisibility() []string {
	return []string{
		string(CertificateVisibilityRemoteUser),
		string(CertificateVisibilityStartTask),
		string(CertificateVisibilityTask),
	}
}

func parseCertificateVisibility(input string) (*CertificateVisibility, error) {
	vals := map[string]CertificateVisibility{
		"remoteuser": CertificateVisibilityRemoteUser,
		"starttask":  CertificateVisibilityStartTask,
		"task":       CertificateVisibilityTask,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CertificateVisibility(input)
	return &out, nil
}

type ComputeNodeDeallocationOption string

const (
	ComputeNodeDeallocationOptionRequeue        ComputeNodeDeallocationOption = "Requeue"
	ComputeNodeDeallocationOptionRetainedData   ComputeNodeDeallocationOption = "RetainedData"
	ComputeNodeDeallocationOptionTaskCompletion ComputeNodeDeallocationOption = "TaskCompletion"
	ComputeNodeDeallocationOptionTerminate      ComputeNodeDeallocationOption = "Terminate"
)

func PossibleValuesForComputeNodeDeallocationOption() []string {
	return []string{
		string(ComputeNodeDeallocationOptionRequeue),
		string(ComputeNodeDeallocationOptionRetainedData),
		string(ComputeNodeDeallocationOptionTaskCompletion),
		string(ComputeNodeDeallocationOptionTerminate),
	}
}

func parseComputeNodeDeallocationOption(input string) (*ComputeNodeDeallocationOption, error) {
	vals := map[string]ComputeNodeDeallocationOption{
		"requeue":        ComputeNodeDeallocationOptionRequeue,
		"retaineddata":   ComputeNodeDeallocationOptionRetainedData,
		"taskcompletion": ComputeNodeDeallocationOptionTaskCompletion,
		"terminate":      ComputeNodeDeallocationOptionTerminate,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ComputeNodeDeallocationOption(input)
	return &out, nil
}

type ComputeNodeFillType string

const (
	ComputeNodeFillTypePack   ComputeNodeFillType = "Pack"
	ComputeNodeFillTypeSpread ComputeNodeFillType = "Spread"
)

func PossibleValuesForComputeNodeFillType() []string {
	return []string{
		string(ComputeNodeFillTypePack),
		string(ComputeNodeFillTypeSpread),
	}
}

func parseComputeNodeFillType(input string) (*ComputeNodeFillType, error) {
	vals := map[string]ComputeNodeFillType{
		"pack":   ComputeNodeFillTypePack,
		"spread": ComputeNodeFillTypeSpread,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ComputeNodeFillType(input)
	return &out, nil
}

type ContainerType string

const (
	ContainerTypeDockerCompatible ContainerType = "DockerCompatible"
)

func PossibleValuesForContainerType() []string {
	return []string{
		string(ContainerTypeDockerCompatible),
	}
}

func parseContainerType(input string) (*ContainerType, error) {
	vals := map[string]ContainerType{
		"dockercompatible": ContainerTypeDockerCompatible,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContainerType(input)
	return &out, nil
}

type ContainerWorkingDirectory string

const (
	ContainerWorkingDirectoryContainerImageDefault ContainerWorkingDirectory = "ContainerImageDefault"
	ContainerWorkingDirectoryTaskWorkingDirectory  ContainerWorkingDirectory = "TaskWorkingDirectory"
)

func PossibleValuesForContainerWorkingDirectory() []string {
	return []string{
		string(ContainerWorkingDirectoryContainerImageDefault),
		string(ContainerWorkingDirectoryTaskWorkingDirectory),
	}
}

func parseContainerWorkingDirectory(input string) (*ContainerWorkingDirectory, error) {
	vals := map[string]ContainerWorkingDirectory{
		"containerimagedefault": ContainerWorkingDirectoryContainerImageDefault,
		"taskworkingdirectory":  ContainerWorkingDirectoryTaskWorkingDirectory,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContainerWorkingDirectory(input)
	return &out, nil
}

type DiffDiskPlacement string

const (
	DiffDiskPlacementCacheDisk DiffDiskPlacement = "CacheDisk"
)

func PossibleValuesForDiffDiskPlacement() []string {
	return []string{
		string(DiffDiskPlacementCacheDisk),
	}
}

func parseDiffDiskPlacement(input string) (*DiffDiskPlacement, error) {
	vals := map[string]DiffDiskPlacement{
		"cachedisk": DiffDiskPlacementCacheDisk,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DiffDiskPlacement(input)
	return &out, nil
}

type DiskEncryptionTarget string

const (
	DiskEncryptionTargetOsDisk        DiskEncryptionTarget = "OsDisk"
	DiskEncryptionTargetTemporaryDisk DiskEncryptionTarget = "TemporaryDisk"
)

func PossibleValuesForDiskEncryptionTarget() []string {
	return []string{
		string(DiskEncryptionTargetOsDisk),
		string(DiskEncryptionTargetTemporaryDisk),
	}
}

func parseDiskEncryptionTarget(input string) (*DiskEncryptionTarget, error) {
	vals := map[string]DiskEncryptionTarget{
		"osdisk":        DiskEncryptionTargetOsDisk,
		"temporarydisk": DiskEncryptionTargetTemporaryDisk,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DiskEncryptionTarget(input)
	return &out, nil
}

type DynamicVNetAssignmentScope string

const (
	DynamicVNetAssignmentScopeJob  DynamicVNetAssignmentScope = "job"
	DynamicVNetAssignmentScopeNone DynamicVNetAssignmentScope = "none"
)

func PossibleValuesForDynamicVNetAssignmentScope() []string {
	return []string{
		string(DynamicVNetAssignmentScopeJob),
		string(DynamicVNetAssignmentScopeNone),
	}
}

func parseDynamicVNetAssignmentScope(input string) (*DynamicVNetAssignmentScope, error) {
	vals := map[string]DynamicVNetAssignmentScope{
		"job":  DynamicVNetAssignmentScopeJob,
		"none": DynamicVNetAssignmentScopeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DynamicVNetAssignmentScope(input)
	return &out, nil
}

type ElevationLevel string

const (
	ElevationLevelAdmin    ElevationLevel = "Admin"
	ElevationLevelNonAdmin ElevationLevel = "NonAdmin"
)

func PossibleValuesForElevationLevel() []string {
	return []string{
		string(ElevationLevelAdmin),
		string(ElevationLevelNonAdmin),
	}
}

func parseElevationLevel(input string) (*ElevationLevel, error) {
	vals := map[string]ElevationLevel{
		"admin":    ElevationLevelAdmin,
		"nonadmin": ElevationLevelNonAdmin,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ElevationLevel(input)
	return &out, nil
}

type IPAddressProvisioningType string

const (
	IPAddressProvisioningTypeBatchManaged        IPAddressProvisioningType = "BatchManaged"
	IPAddressProvisioningTypeNoPublicIPAddresses IPAddressProvisioningType = "NoPublicIPAddresses"
	IPAddressProvisioningTypeUserManaged         IPAddressProvisioningType = "UserManaged"
)

func PossibleValuesForIPAddressProvisioningType() []string {
	return []string{
		string(IPAddressProvisioningTypeBatchManaged),
		string(IPAddressProvisioningTypeNoPublicIPAddresses),
		string(IPAddressProvisioningTypeUserManaged),
	}
}

func parseIPAddressProvisioningType(input string) (*IPAddressProvisioningType, error) {
	vals := map[string]IPAddressProvisioningType{
		"batchmanaged":        IPAddressProvisioningTypeBatchManaged,
		"nopublicipaddresses": IPAddressProvisioningTypeNoPublicIPAddresses,
		"usermanaged":         IPAddressProvisioningTypeUserManaged,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IPAddressProvisioningType(input)
	return &out, nil
}

type InboundEndpointProtocol string

const (
	InboundEndpointProtocolTCP InboundEndpointProtocol = "TCP"
	InboundEndpointProtocolUDP InboundEndpointProtocol = "UDP"
)

func PossibleValuesForInboundEndpointProtocol() []string {
	return []string{
		string(InboundEndpointProtocolTCP),
		string(InboundEndpointProtocolUDP),
	}
}

func parseInboundEndpointProtocol(input string) (*InboundEndpointProtocol, error) {
	vals := map[string]InboundEndpointProtocol{
		"tcp": InboundEndpointProtocolTCP,
		"udp": InboundEndpointProtocolUDP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InboundEndpointProtocol(input)
	return &out, nil
}

type InterNodeCommunicationState string

const (
	InterNodeCommunicationStateDisabled InterNodeCommunicationState = "Disabled"
	InterNodeCommunicationStateEnabled  InterNodeCommunicationState = "Enabled"
)

func PossibleValuesForInterNodeCommunicationState() []string {
	return []string{
		string(InterNodeCommunicationStateDisabled),
		string(InterNodeCommunicationStateEnabled),
	}
}

func parseInterNodeCommunicationState(input string) (*InterNodeCommunicationState, error) {
	vals := map[string]InterNodeCommunicationState{
		"disabled": InterNodeCommunicationStateDisabled,
		"enabled":  InterNodeCommunicationStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InterNodeCommunicationState(input)
	return &out, nil
}

type LoginMode string

const (
	LoginModeBatch       LoginMode = "Batch"
	LoginModeInteractive LoginMode = "Interactive"
)

func PossibleValuesForLoginMode() []string {
	return []string{
		string(LoginModeBatch),
		string(LoginModeInteractive),
	}
}

func parseLoginMode(input string) (*LoginMode, error) {
	vals := map[string]LoginMode{
		"batch":       LoginModeBatch,
		"interactive": LoginModeInteractive,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LoginMode(input)
	return &out, nil
}

type NetworkSecurityGroupRuleAccess string

const (
	NetworkSecurityGroupRuleAccessAllow NetworkSecurityGroupRuleAccess = "Allow"
	NetworkSecurityGroupRuleAccessDeny  NetworkSecurityGroupRuleAccess = "Deny"
)

func PossibleValuesForNetworkSecurityGroupRuleAccess() []string {
	return []string{
		string(NetworkSecurityGroupRuleAccessAllow),
		string(NetworkSecurityGroupRuleAccessDeny),
	}
}

func parseNetworkSecurityGroupRuleAccess(input string) (*NetworkSecurityGroupRuleAccess, error) {
	vals := map[string]NetworkSecurityGroupRuleAccess{
		"allow": NetworkSecurityGroupRuleAccessAllow,
		"deny":  NetworkSecurityGroupRuleAccessDeny,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkSecurityGroupRuleAccess(input)
	return &out, nil
}

type NodePlacementPolicyType string

const (
	NodePlacementPolicyTypeRegional NodePlacementPolicyType = "Regional"
	NodePlacementPolicyTypeZonal    NodePlacementPolicyType = "Zonal"
)

func PossibleValuesForNodePlacementPolicyType() []string {
	return []string{
		string(NodePlacementPolicyTypeRegional),
		string(NodePlacementPolicyTypeZonal),
	}
}

func parseNodePlacementPolicyType(input string) (*NodePlacementPolicyType, error) {
	vals := map[string]NodePlacementPolicyType{
		"regional": NodePlacementPolicyTypeRegional,
		"zonal":    NodePlacementPolicyTypeZonal,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NodePlacementPolicyType(input)
	return &out, nil
}

type PoolProvisioningState string

const (
	PoolProvisioningStateDeleting  PoolProvisioningState = "Deleting"
	PoolProvisioningStateSucceeded PoolProvisioningState = "Succeeded"
)

func PossibleValuesForPoolProvisioningState() []string {
	return []string{
		string(PoolProvisioningStateDeleting),
		string(PoolProvisioningStateSucceeded),
	}
}

func parsePoolProvisioningState(input string) (*PoolProvisioningState, error) {
	vals := map[string]PoolProvisioningState{
		"deleting":  PoolProvisioningStateDeleting,
		"succeeded": PoolProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PoolProvisioningState(input)
	return &out, nil
}

type StorageAccountType string

const (
	StorageAccountTypePremiumLRS  StorageAccountType = "Premium_LRS"
	StorageAccountTypeStandardLRS StorageAccountType = "Standard_LRS"
)

func PossibleValuesForStorageAccountType() []string {
	return []string{
		string(StorageAccountTypePremiumLRS),
		string(StorageAccountTypeStandardLRS),
	}
}

func parseStorageAccountType(input string) (*StorageAccountType, error) {
	vals := map[string]StorageAccountType{
		"premium_lrs":  StorageAccountTypePremiumLRS,
		"standard_lrs": StorageAccountTypeStandardLRS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageAccountType(input)
	return &out, nil
}
