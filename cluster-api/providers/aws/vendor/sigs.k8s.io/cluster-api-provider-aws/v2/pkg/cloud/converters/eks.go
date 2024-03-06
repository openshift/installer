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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
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
func AddonSDKToAddonState(eksAddon *eks.Addon) *ekscontrolplanev1.AddonState {
	addonState := &ekscontrolplanev1.AddonState{
		Name:                  aws.StringValue(eksAddon.AddonName),
		Version:               aws.StringValue(eksAddon.AddonVersion),
		ARN:                   aws.StringValue(eksAddon.AddonArn),
		CreatedAt:             metav1.NewTime(*eksAddon.CreatedAt),
		ModifiedAt:            metav1.NewTime(*eksAddon.ModifiedAt),
		Status:                eksAddon.Status,
		ServiceAccountRoleArn: eksAddon.ServiceAccountRoleArn,
		Issues:                []ekscontrolplanev1.AddonIssue{},
	}
	if eksAddon.Health != nil {
		for _, issue := range eksAddon.Health.Issues {
			addonState.Issues = append(addonState.Issues, ekscontrolplanev1.AddonIssue{
				Code:        issue.Code,
				Message:     issue.Message,
				ResourceIDs: FromAWSStringSlice(issue.ResourceIds),
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
func TaintToSDK(taint expinfrav1.Taint) (*eks.Taint, error) {
	convertedEffect, err := TaintEffectToSDK(taint.Effect)
	if err != nil {
		return nil, fmt.Errorf("converting taint effect %s: %w", taint.Effect, err)
	}
	return &eks.Taint{
		Effect: aws.String(convertedEffect),
		Key:    aws.String(taint.Key),
		Value:  aws.String(taint.Value),
	}, nil
}

// TaintsToSDK is used to convert an array of CAPA Taints to AWS SDK taints.
func TaintsToSDK(taints expinfrav1.Taints) ([]*eks.Taint, error) {
	converted := []*eks.Taint{}

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
func TaintsFromSDK(taints []*eks.Taint) (expinfrav1.Taints, error) {
	converted := expinfrav1.Taints{}
	for _, taint := range taints {
		convertedEffect, err := TaintEffectFromSDK(*taint.Effect)
		if err != nil {
			return nil, fmt.Errorf("converting taint effect %s: %w", *taint.Effect, err)
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
func TaintEffectToSDK(effect expinfrav1.TaintEffect) (string, error) {
	switch effect {
	case expinfrav1.TaintEffectNoExecute:
		return eks.TaintEffectNoExecute, nil
	case expinfrav1.TaintEffectPreferNoSchedule:
		return eks.TaintEffectPreferNoSchedule, nil
	case expinfrav1.TaintEffectNoSchedule:
		return eks.TaintEffectNoSchedule, nil
	default:
		return "", ErrUnknowTaintEffect
	}
}

// TaintEffectFromSDK is used to convert a AWS SDK taint effect value to a TaintEffect.
func TaintEffectFromSDK(effect string) (expinfrav1.TaintEffect, error) {
	switch effect {
	case eks.TaintEffectNoExecute:
		return expinfrav1.TaintEffectNoExecute, nil
	case eks.TaintEffectPreferNoSchedule:
		return expinfrav1.TaintEffectPreferNoSchedule, nil
	case eks.TaintEffectNoSchedule:
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
			GroupsClaim:                aws.StringValue(in.GroupsClaim),
			GroupsPrefix:               aws.StringValue(in.GroupsPrefix),
			IdentityProviderConfigName: in.IdentityProviderConfigName,
			IssuerURL:                  in.IssuerURL,
			RequiredClaims:             in.RequiredClaims,
			Tags:                       in.Tags,
			UsernameClaim:              aws.StringValue(in.UsernameClaim),
			UsernamePrefix:             aws.StringValue(in.UsernamePrefix),
		}
	}

	return nil
}

// CapacityTypeToSDK is used to convert a CapacityType to the AWS SDK capacity type value.
func CapacityTypeToSDK(capacityType expinfrav1.ManagedMachinePoolCapacityType) (string, error) {
	switch capacityType {
	case expinfrav1.ManagedMachinePoolCapacityTypeOnDemand:
		return eks.CapacityTypesOnDemand, nil
	case expinfrav1.ManagedMachinePoolCapacityTypeSpot:
		return eks.CapacityTypesSpot, nil
	default:
		return "", ErrUnknownCapacityType
	}
}

// NodegroupUpdateconfigToSDK is used to convert a CAPA UpdateConfig to AWS SDK NodegroupUpdateConfig.
func NodegroupUpdateconfigToSDK(updateConfig *expinfrav1.UpdateConfig) *eks.NodegroupUpdateConfig {
	if updateConfig == nil {
		return nil
	}

	converted := &eks.NodegroupUpdateConfig{}
	if updateConfig.MaxUnavailable != nil {
		converted.MaxUnavailable = aws.Int64(int64(*updateConfig.MaxUnavailable))
	}
	if updateConfig.MaxUnavailablePercentage != nil {
		converted.MaxUnavailablePercentage = aws.Int64(int64(*updateConfig.MaxUnavailablePercentage))
	}

	return converted
}

// NodegroupUpdateconfigFromSDK is used to convert a AWS SDK NodegroupUpdateConfig to a CAPA UpdateConfig.
func NodegroupUpdateconfigFromSDK(ngUpdateConfig *eks.NodegroupUpdateConfig) *expinfrav1.UpdateConfig {
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
