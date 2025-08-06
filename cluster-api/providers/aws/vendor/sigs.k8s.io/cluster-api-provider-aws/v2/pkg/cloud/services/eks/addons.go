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

package eks

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/converters"
	eksaddons "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/eks/addons"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
)

func (s *Service) reconcileAddons(ctx context.Context) error {
	s.scope.Info("Reconciling EKS addons")

	eksClusterName := s.scope.KubernetesClusterName()

	// Get available addon names for the cluster
	addonNames, err := s.listAddons(ctx, eksClusterName)
	if err != nil {
		s.Error(err, "failed listing addons")
		return fmt.Errorf("listing eks addons: %w", err)
	}

	// Get installed addons for the cluster
	s.scope.Debug("getting installed eks addons", "cluster", eksClusterName)
	installed, err := s.getClusterAddonsInstalled(ctx, eksClusterName, addonNames)
	if err != nil {
		return fmt.Errorf("getting installed eks addons: %w", err)
	}

	// Get the addons from the spec we want for the cluster
	desiredAddons := s.translateAPIToAddon(s.scope.Addons())

	// If there are no addons desired or installed then do nothing
	if len(installed) == 0 && len(desiredAddons) == 0 {
		s.scope.Info("no addons installed and no addons to install, no action needed")
		return nil
	}

	//  Compute operations to move installed to desired
	s.scope.Debug("creating eks addons plan", "cluster", eksClusterName, "numdesired", len(desiredAddons), "numinstalled", len(installed))
	addonsPlan := eksaddons.NewPlan(eksClusterName, desiredAddons, installed, s.EKSClient, s.scope.MaxWaitActiveUpdateDelete)
	procedures, err := addonsPlan.Create(ctx)
	if err != nil {
		s.scope.Error(err, "failed creating eks addons plane")
		return fmt.Errorf("creating eks addons plan: %w", err)
	}
	s.scope.Debug("computed EKS addons plan", "numprocs", len(procedures))

	// Perform required operations
	for _, procedure := range procedures {
		s.scope.Debug("Executing addon procedure", "name", procedure.Name())
		if err := procedure.Do(ctx); err != nil {
			s.scope.Error(err, "failed executing addon procedure", "name", procedure.Name())
			return fmt.Errorf("%s: %w", procedure.Name(), err)
		}
	}

	// Update status with addons installed details
	// Note: we are not relying on the computed state from the operations as we still want
	// to update the state even if there are no operations to capture things like status changes
	s.scope.Debug("getting installed eks addons to update status", "cluster", eksClusterName)
	addonState, err := s.getInstalledState(ctx, eksClusterName, addonNames)
	if err != nil {
		return fmt.Errorf("getting installed state of eks addons: %w", err)
	}
	s.scope.ControlPlane.Status.Addons = addonState

	// Persist status and record event
	if err := s.scope.PatchObject(); err != nil {
		return fmt.Errorf("failed to update control plane: %w", err)
	}
	record.Eventf(s.scope.ControlPlane, "SuccessfulReconcileEKSClusterAddons", "Reconciled addons for EKS Cluster %s", s.scope.KubernetesClusterName())
	s.scope.Debug("Reconcile EKS addons completed successfully")

	return nil
}

func (s *Service) getClusterAddonsInstalled(ctx context.Context, eksClusterName string, addonNames []string) ([]*eksaddons.EKSAddon, error) {
	s.Debug("getting eks addons installed")

	addonsInstalled := []*eksaddons.EKSAddon{}
	if len(addonNames) == 0 {
		s.scope.Info("no eks addons installed in cluster", "cluster", eksClusterName)
		return addonsInstalled, nil
	}

	for _, addon := range addonNames {
		describeInput := &eks.DescribeAddonInput{
			AddonName:   aws.String(addon),
			ClusterName: &eksClusterName,
		}
		describeOutput, err := s.EKSClient.DescribeAddon(ctx, describeInput)
		if err != nil {
			return addonsInstalled, fmt.Errorf("describing eks addon %s: %w", addon, err)
		}

		if describeOutput.Addon == nil {
			continue
		}
		s.scope.Debug("describe output", "output", describeOutput.Addon)
		status := string(describeOutput.Addon.Status)
		installedAddon := &eksaddons.EKSAddon{
			Name:                  describeOutput.Addon.AddonName,
			Version:               describeOutput.Addon.AddonVersion,
			ARN:                   describeOutput.Addon.AddonArn,
			Configuration:         describeOutput.Addon.ConfigurationValues,
			Tags:                  infrav1.Tags{},
			Status:                aws.String(status),
			ServiceAccountRoleARN: describeOutput.Addon.ServiceAccountRoleArn,
		}
		for k, v := range describeOutput.Addon.Tags {
			installedAddon.Tags[k] = v
		}

		addonsInstalled = append(addonsInstalled, installedAddon)
	}

	return addonsInstalled, nil
}

func (s *Service) getInstalledState(ctx context.Context, eksClusterName string, addonNames []string) ([]ekscontrolplanev1.AddonState, error) {
	s.Debug("getting eks addons installed to create state")

	addonState := []ekscontrolplanev1.AddonState{}
	if len(addonNames) == 0 {
		s.scope.Info("no eks addons installed in cluster", "cluster", eksClusterName)
		return addonState, nil
	}

	for _, addonName := range addonNames {
		describeInput := &eks.DescribeAddonInput{
			AddonName:   aws.String(addonName),
			ClusterName: &eksClusterName,
		}
		describeOutput, err := s.EKSClient.DescribeAddon(ctx, describeInput)
		if err != nil {
			return addonState, fmt.Errorf("describing eks addon %s: %w", addonName, err)
		}

		if describeOutput.Addon == nil {
			continue
		}
		s.scope.Debug("describe output", "output", describeOutput.Addon)

		installedAddonState := converters.AddonSDKToAddonState(describeOutput.Addon)
		addonState = append(addonState, *installedAddonState)
	}

	return addonState, nil
}

func (s *Service) listAddons(ctx context.Context, eksClusterName string) ([]string, error) {
	s.Debug("getting list of eks addons")

	input := &eks.ListAddonsInput{
		ClusterName: &eksClusterName,
	}

	addons := []string{}
	output, err := s.EKSClient.ListAddons(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("listing eks addons: %w", err)
	}

	addons = append(addons, output.Addons...)

	return addons, nil
}

func (s *Service) translateAPIToAddon(addons []ekscontrolplanev1.Addon) []*eksaddons.EKSAddon {
	converted := []*eksaddons.EKSAddon{}

	for i := range addons {
		addon := addons[i]
		conflict, err := convertConflictResolution(*addon.ConflictResolution)
		if err != nil {
			s.scope.Error(err, err.Error())
		}
		convertedAddon := &eksaddons.EKSAddon{
			Name:                  &addon.Name,
			Version:               &addon.Version,
			Configuration:         convertConfiguration(addon.Configuration),
			Tags:                  ngTags(s.scope.Cluster.Name, s.scope.AdditionalTags()),
			ResolveConflict:       conflict,
			ServiceAccountRoleARN: addon.ServiceAccountRoleArn,
		}

		converted = append(converted, convertedAddon)
	}

	return converted
}

// WaitUntilAddonDeleted is blocking function to wait until EKS Addon is Deleted.
func (k *EKSClient) WaitUntilAddonDeleted(ctx context.Context, input *eks.DescribeAddonInput, maxWait time.Duration) error {
	waiter := eks.NewAddonDeletedWaiter(k, func(o *eks.AddonDeletedWaiterOptions) {
		o.LogWaitAttempts = true
	})

	return waiter.Wait(ctx, input, maxWait)
}

func convertConflictResolution(conflict ekscontrolplanev1.AddonResolution) (*string, error) {
	switch conflict {
	case ekscontrolplanev1.AddonResolutionNone:
		return aws.String(string(ekstypes.ResolveConflictsNone)), nil

	case ekscontrolplanev1.AddonResolutionOverwrite:
		return aws.String(string(ekstypes.ResolveConflictsOverwrite)), nil

	case ekscontrolplanev1.AddonResolutionPreserve:
		return aws.String(string(ekstypes.ResolveConflictsPreserve)), nil

	// Defaulting to behavior "Overwrite" as documented
	default:
		return aws.String(string(ekstypes.ResolveConflictsOverwrite)), fmt.Errorf("failed to determine adddonResolution; defaulting to Overwrite")
	}
}

func convertConfiguration(configuration string) *string {
	if configuration == "" {
		return nil
	}
	return &configuration
}
