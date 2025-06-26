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

// Package identityprovider provides a plan to manage EKS OIDC identity provider association.
package identityprovider

import (
	"context"

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/eks"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/planner"
)

// NewPlan creates plan to manage EKS OIDC identity provider association.
func NewPlan(clusterName string, currentIdentityProvider, desiredIdentityProvider *OidcIdentityProviderConfig, client eks.Client, log logger.Wrapper) planner.Plan {
	return &plan{
		currentIdentityProvider: currentIdentityProvider,
		desiredIdentityProvider: desiredIdentityProvider,
		eksClient:               client,
		clusterName:             clusterName,
		log:                     log,
	}
}

// Plan is a plan that will manage EKS OIDC identity provider association.
type plan struct {
	currentIdentityProvider *OidcIdentityProviderConfig
	desiredIdentityProvider *OidcIdentityProviderConfig
	eksClient               eks.Client
	log                     logger.Wrapper
	clusterName             string
}

// Create will create the plan (i.e. list of procedures) for managing EKS OIDC identity provider association.
func (p *plan) Create(_ context.Context) ([]planner.Procedure, error) {
	procedures := []planner.Procedure{}

	if p.desiredIdentityProvider == nil && p.currentIdentityProvider == nil {
		return procedures, nil
	}

	// no config is mentioned deleted provider if we have one
	if p.desiredIdentityProvider == nil {
		// disassociation will also trigger deletion hence
		// we do nothing in case of ConfigStatusDeleting as it will happen eventually
		if p.currentIdentityProvider.Status == string(ekstypes.ConfigStatusActive) {
			procedures = append(procedures, &DisassociateIdentityProviderConfig{plan: p})
		}

		return procedures, nil
	}

	// create case
	if p.currentIdentityProvider == nil {
		procedures = append(procedures, &AssociateIdentityProviderProcedure{plan: p})
		return procedures, nil
	}

	if p.currentIdentityProvider.IsEqual(p.desiredIdentityProvider) {
		tagsDiff := p.desiredIdentityProvider.Tags.Difference(p.currentIdentityProvider.Tags)
		if len(tagsDiff) > 0 {
			procedures = append(procedures, &UpdatedIdentityProviderTagsProcedure{plan: p})
		}

		if len(p.desiredIdentityProvider.Tags) == 0 && len(p.currentIdentityProvider.Tags) != 0 {
			procedures = append(procedures, &RemoveIdentityProviderTagsProcedure{plan: p})
		}
		switch p.currentIdentityProvider.Status {
		case string(ekstypes.ConfigStatusActive):
			// config active no work to be done
			return procedures, nil
		case string(ekstypes.ConfigStatusCreating):
			// no change need wait for association to complete
			procedures = append(procedures, &WaitIdentityProviderAssociatedProcedure{plan: p})
		}
	} else {
		procedures = append(procedures, &DisassociateIdentityProviderConfig{plan: p})
	}

	return procedures, nil
}
