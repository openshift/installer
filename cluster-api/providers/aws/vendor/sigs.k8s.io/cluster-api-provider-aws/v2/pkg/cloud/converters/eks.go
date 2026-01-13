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

// Package converters provides conversion functions for AWS SDK types to CAPA types.
package converters

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/eks/identityprovider"
)

var (
	// ErrUnknowTaintEffect is an error when a unknown TaintEffect is used.
	ErrUnknowTaintEffect = errors.New("uknown taint effect")

	// ErrUnknownCapacityType is an error when a unknown CapacityType is used.
	ErrUnknownCapacityType = errors.New("unknown capacity type")
)

// AddonSDKToAddonState is used to convert an AWS SDK Addon to a control plane AddonState.
func AddonSDKToAddonState(eksAddon *ekstypes.Addon) *ekscontrolplanev1.AddonState {
	addonState := &ekscontrolplanev1.AddonState{
		Name:                  aws.ToString(eksAddon.AddonName),
		Version:               aws.ToString(eksAddon.AddonVersion),
		ARN:                   aws.ToString(eksAddon.AddonArn),
		CreatedAt:             metav1.NewTime(*eksAddon.CreatedAt),
		ModifiedAt:            metav1.NewTime(*eksAddon.ModifiedAt),
		Status:                aws.String(string(eksAddon.Status)),
		ServiceAccountRoleArn: eksAddon.ServiceAccountRoleArn,
		Issues:                []ekscontrolplanev1.AddonIssue{},
	}
	if eksAddon.Health != nil {
		for _, issue := range eksAddon.Health.Issues {
			addonState.Issues = append(addonState.Issues, ekscontrolplanev1.AddonIssue{
				Code:        aws.String(string(issue.Code)),
				Message:     issue.Message,
				ResourceIDs: issue.ResourceIds,
			})
		}
	}

	return addonState
}

// FromAWSStringSlice will converts an AWS string pointer slice.
func FromAWSStringSlice(from []*string) []string {
	converted := []string{}
	for _, s := range from {
		converted = append(converted, *s)
	}

	return converted
}

// TaintToSDK is used to a CAPA Taint to AWS SDK taint.
func TaintToSDK(taint expinfrav1.Taint) (ekstypes.Taint, error) {
	convertedEffect, err := TaintEffectToSDK(taint.Effect)
	if err != nil {
		return ekstypes.Taint{}, fmt.Errorf("converting taint effect %s: %w", taint.Effect, err)
	}
	return ekstypes.Taint{
		Effect: convertedEffect,
		Key:    aws.String(taint.Key),
		Value:  aws.String(taint.Value),
	}, nil
}

// TaintsToSDK is used to convert an array of CAPA Taints to AWS SDK taints.
func TaintsToSDK(taints expinfrav1.Taints) ([]ekstypes.Taint, error) {
	converted := []ekstypes.Taint{}

	for _, taint := range taints {
		convertedTaint, err := TaintToSDK(taint)
		if err != nil {
			return nil, fmt.Errorf("converting taint: %w", err)
		}
		converted = append(converted, convertedTaint)
	}

	return converted, nil
}

// TaintsFromSDK is used to convert an array of AWS SDK taints to CAPA Taints.
func TaintsFromSDK(taints []ekstypes.Taint) (expinfrav1.Taints, error) {
	converted := expinfrav1.Taints{}
	for _, taint := range taints {
		convertedEffect, err := TaintEffectFromSDK(taint.Effect)
		if err != nil {
			return nil, fmt.Errorf("converting taint effect %s: %w", taint.Effect, err)
		}
		converted = append(converted, expinfrav1.Taint{
			Effect: convertedEffect,
			Key:    *taint.Key,
			Value:  *taint.Value,
		})
	}

	return converted, nil
}

// TaintEffectToSDK is used to convert a TaintEffect to the AWS SDK taint effect value.
func TaintEffectToSDK(effect expinfrav1.TaintEffect) (ekstypes.TaintEffect, error) {
	switch effect {
	case expinfrav1.TaintEffectNoExecute:
		return ekstypes.TaintEffectNoExecute, nil
	case expinfrav1.TaintEffectPreferNoSchedule:
		return ekstypes.TaintEffectPreferNoSchedule, nil
	case expinfrav1.TaintEffectNoSchedule:
		return ekstypes.TaintEffectNoSchedule, nil
	default:
		return "", ErrUnknowTaintEffect
	}
}

// TaintEffectFromSDK is used to convert a AWS SDK taint effect value to a TaintEffect.
func TaintEffectFromSDK(effect ekstypes.TaintEffect) (expinfrav1.TaintEffect, error) {
	switch effect {
	case ekstypes.TaintEffectNoExecute:
		return expinfrav1.TaintEffectNoExecute, nil
	case ekstypes.TaintEffectPreferNoSchedule:
		return expinfrav1.TaintEffectPreferNoSchedule, nil
	case ekstypes.TaintEffectNoSchedule:
		return expinfrav1.TaintEffectNoSchedule, nil
	default:
		return "", ErrUnknowTaintEffect
	}
}

// ConvertSDKToIdentityProvider is used to convert an AWS SDK OIDCIdentityProviderConfig to a CAPA OidcIdentityProviderConfig.
func ConvertSDKToIdentityProvider(in *ekscontrolplanev1.OIDCIdentityProviderConfig) *identityprovider.OidcIdentityProviderConfig {
	if in != nil {
		if in.RequiredClaims == nil {
			in.RequiredClaims = make(map[string]string)
		}
		return &identityprovider.OidcIdentityProviderConfig{
			ClientID:                   in.ClientID,
			GroupsClaim:                aws.ToString(in.GroupsClaim),
			GroupsPrefix:               aws.ToString(in.GroupsPrefix),
			IdentityProviderConfigName: in.IdentityProviderConfigName,
			IssuerURL:                  in.IssuerURL,
			RequiredClaims:             in.RequiredClaims,
			Tags:                       in.Tags,
			UsernameClaim:              aws.ToString(in.UsernameClaim),
			UsernamePrefix:             aws.ToString(in.UsernamePrefix),
		}
	}

	return nil
}

// CapacityTypeToSDK is used to convert a CapacityType to the AWS SDK capacity type value.
func CapacityTypeToSDK(capacityType expinfrav1.ManagedMachinePoolCapacityType) (ekstypes.CapacityTypes, error) {
	switch capacityType {
	case expinfrav1.ManagedMachinePoolCapacityTypeOnDemand:
		return ekstypes.CapacityTypesOnDemand, nil
	case expinfrav1.ManagedMachinePoolCapacityTypeSpot:
		return ekstypes.CapacityTypesSpot, nil
	default:
		return "", ErrUnknownCapacityType
	}
}

// NodegroupUpdateconfigToSDK is used to convert a CAPA UpdateConfig to AWS SDK NodegroupUpdateConfig.
func NodegroupUpdateconfigToSDK(updateConfig *expinfrav1.UpdateConfig) (*ekstypes.NodegroupUpdateConfig, error) {
	if updateConfig == nil {
		return nil, nil
	}

	converted := &ekstypes.NodegroupUpdateConfig{}
	if updateConfig.MaxUnavailable != nil {
		//nolint:G115 // Added golint exception as there is a kubebuilder validation configured
		converted.MaxUnavailable = aws.Int32(int32(*updateConfig.MaxUnavailable))
	}
	if updateConfig.MaxUnavailablePercentage != nil {
		//nolint:G115 // Added golint exception as there is a kubebuilder validation configured
		converted.MaxUnavailablePercentage = aws.Int32(int32(*updateConfig.MaxUnavailablePercentage))
	}

	return converted, nil
}

// NodegroupUpdateconfigFromSDK is used to convert a AWS SDK NodegroupUpdateConfig to a CAPA UpdateConfig.
func NodegroupUpdateconfigFromSDK(ngUpdateConfig *ekstypes.NodegroupUpdateConfig) *expinfrav1.UpdateConfig {
	if ngUpdateConfig == nil {
		return nil
	}

	converted := &expinfrav1.UpdateConfig{}
	if ngUpdateConfig.MaxUnavailable != nil {
		converted.MaxUnavailable = aws.Int(int(*ngUpdateConfig.MaxUnavailable))
	}
	if ngUpdateConfig.MaxUnavailablePercentage != nil {
		converted.MaxUnavailablePercentage = aws.Int(int(*ngUpdateConfig.MaxUnavailablePercentage))
	}

	return converted
}

// NodeRepairConfigToSDK is used to convert a CAPA NodeRepairConfig to AWS SDK NodeRepairConfig.
func NodeRepairConfigToSDK(repairConfig *expinfrav1.NodeRepairConfig) *ekstypes.NodeRepairConfig {
	if repairConfig == nil {
		// Default to disabled if not specified to avoid behavior changes
		return &ekstypes.NodeRepairConfig{
			Enabled: aws.Bool(false),
		}
	}

	return &ekstypes.NodeRepairConfig{
		Enabled: repairConfig.Enabled,
	}
}

// AMITypeToSDK converts a CAPA ManagedMachineAMIType to AWS SDK AMIType.
func AMITypeToSDK(amiType expinfrav1.ManagedMachineAMIType) ekstypes.AMITypes {
	switch amiType {
	case expinfrav1.Al2x86_64:
		return ekstypes.AMITypesAl2X8664
	case expinfrav1.Al2x86_64GPU:
		return ekstypes.AMITypesAl2X8664Gpu
	case expinfrav1.Al2Arm64:
		return ekstypes.AMITypesAl2Arm64
	case expinfrav1.Custom:
		return ekstypes.AMITypesCustom
	case expinfrav1.BottleRocketArm64:
		return ekstypes.AMITypesBottlerocketArm64
	case expinfrav1.BottleRocketx86_64:
		return ekstypes.AMITypesBottlerocketX8664
	case expinfrav1.BottleRocketArm64Fips:
		return ekstypes.AMITypesBottlerocketArm64Fips
	case expinfrav1.BottleRocketx86_64Fips:
		return ekstypes.AMITypesBottlerocketX8664Fips
	case expinfrav1.BottleRocketArm64Nvidia:
		return ekstypes.AMITypesBottlerocketArm64Nvidia
	case expinfrav1.BottleRocketx86_64Nvidia:
		return ekstypes.AMITypesBottlerocketX8664Nvidia
	case expinfrav1.WindowsCore2019x86_64:
		return ekstypes.AMITypesWindowsCore2019X8664
	case expinfrav1.WindowsFull2019x86_64:
		return ekstypes.AMITypesWindowsFull2019X8664
	case expinfrav1.WindowsCore2022x86_64:
		return ekstypes.AMITypesWindowsCore2022X8664
	case expinfrav1.WindowsFull2022x86_64:
		return ekstypes.AMITypesWindowsFull2022X8664
	case expinfrav1.Al2023Arm64:
		return ekstypes.AMITypesAl2023Arm64Standard
	case expinfrav1.Al2023x86_64:
		return ekstypes.AMITypesAl2023X8664Standard
	case expinfrav1.Al2023x86_64Neuron:
		return ekstypes.AMITypesAl2023X8664Neuron
	case expinfrav1.Al2023x86_64Nvidia:
		return ekstypes.AMITypesAl2023X8664Nvidia
	case expinfrav1.Al2023Arm64Nvidia:
		return ekstypes.AMITypesAl2023Arm64Nvidia
	default:
		return ekstypes.AMITypesCustom
	}
}

// AddonConflictResolutionToSDK converts CAPA conflict resolution types to SDK types.
func AddonConflictResolutionToSDK(conflict *string) ekstypes.ResolveConflicts {
	if *conflict == string(ekscontrolplanev1.AddonResolutionNone) {
		return ekstypes.ResolveConflictsNone
	}
	return ekstypes.ResolveConflictsOverwrite
}

// AddonConflictResolutionFromSDK converts SDK conflict resolution types to CAPA types.
func AddonConflictResolutionFromSDK(conflict ekstypes.ResolveConflicts) *string {
	if conflict == ekstypes.ResolveConflictsNone {
		return aws.String(string(ekscontrolplanev1.AddonResolutionNone))
	}
	return aws.String(string(ekscontrolplanev1.AddonResolutionOverwrite))
}

// SupportTypeToSDK converts CAPA upgrade support policy types to SDK types.
func SupportTypeToSDK(input ekscontrolplanev1.UpgradePolicy) ekstypes.SupportType {
	if input == ekscontrolplanev1.UpgradePolicyStandard {
		return ekstypes.SupportTypeStandard
	}
	return ekstypes.SupportTypeExtended
}
