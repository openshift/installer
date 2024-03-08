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

package identityprovider

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/pkg/errors"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/wait"
)

var oidcType = aws.String("oidc")

// WaitIdentityProviderAssociatedProcedure waits for the identity provider to be associated.
type WaitIdentityProviderAssociatedProcedure struct {
	plan *plan
}

// Name returns the name of the procedure.
func (w *WaitIdentityProviderAssociatedProcedure) Name() string {
	return "wait_identity_provider_association"
}

// Do waits for the identity provider to be associated.
func (w *WaitIdentityProviderAssociatedProcedure) Do(ctx context.Context) error {
	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
		out, err := w.plan.eksClient.DescribeIdentityProviderConfigWithContext(ctx, &eks.DescribeIdentityProviderConfigInput{
			ClusterName: aws.String(w.plan.clusterName),
			IdentityProviderConfig: &eks.IdentityProviderConfig{
				Name: aws.String(w.plan.currentIdentityProvider.IdentityProviderConfigName),
				Type: oidcType,
			},
		})

		if err != nil {
			return false, err
		}

		if aws.StringValue(out.IdentityProviderConfig.Oidc.Status) == eks.ConfigStatusActive {
			return true, nil
		}

		return false, nil
	}); err != nil {
		return errors.Wrap(err, "failed waiting for identity provider association to be ready")
	}

	return nil
}

// DisassociateIdentityProviderConfig disassociates the identity provider.
type DisassociateIdentityProviderConfig struct {
	plan *plan
}

// Name returns the name of the procedure.
func (d *DisassociateIdentityProviderConfig) Name() string {
	return "dissociate_identity_provider"
}

// Do disassociates the identity provider.
func (d *DisassociateIdentityProviderConfig) Do(ctx context.Context) error {
	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
		_, err := d.plan.eksClient.DisassociateIdentityProviderConfigWithContext(ctx, &eks.DisassociateIdentityProviderConfigInput{
			ClusterName: aws.String(d.plan.clusterName),
			IdentityProviderConfig: &eks.IdentityProviderConfig{
				Name: aws.String(d.plan.currentIdentityProvider.IdentityProviderConfigName),
				Type: oidcType,
			},
		})

		if err != nil {
			return false, err
		}

		return true, nil
	}); err != nil {
		return errors.Wrap(err, "failing disassociating identity provider config")
	}

	return nil
}

// AssociateIdentityProviderProcedure associates the identity provider.
type AssociateIdentityProviderProcedure struct {
	plan *plan
}

// Name returns the name of the procedure.
func (a *AssociateIdentityProviderProcedure) Name() string {
	return "associate_identity_provider"
}

// Do associates the identity provider.
func (a *AssociateIdentityProviderProcedure) Do(ctx context.Context) error {
	oidc := a.plan.desiredIdentityProvider
	input := &eks.AssociateIdentityProviderConfigInput{
		ClusterName: aws.String(a.plan.clusterName),
		Oidc: &eks.OidcIdentityProviderConfigRequest{
			ClientId:                   aws.String(oidc.ClientID),
			GroupsClaim:                aws.String(oidc.GroupsClaim),
			GroupsPrefix:               aws.String(oidc.GroupsPrefix),
			IdentityProviderConfigName: aws.String(oidc.IdentityProviderConfigName),
			IssuerUrl:                  aws.String(oidc.IssuerURL),
			RequiredClaims:             aws.StringMap(oidc.RequiredClaims),
			UsernameClaim:              aws.String(oidc.UsernameClaim),
			UsernamePrefix:             aws.String(oidc.UsernamePrefix),
		},
	}

	if len(oidc.Tags) > 0 {
		input.Tags = aws.StringMap(oidc.Tags)
	}

	_, err := a.plan.eksClient.AssociateIdentityProviderConfigWithContext(ctx, input)
	if err != nil {
		return errors.Wrap(err, "failed associating identity provider")
	}

	return nil
}

// UpdatedIdentityProviderTagsProcedure updates the tags for the identity provider.
type UpdatedIdentityProviderTagsProcedure struct {
	plan *plan
}

// Name returns the name of the procedure.
func (u *UpdatedIdentityProviderTagsProcedure) Name() string {
	return "update_identity_provider_tags"
}

// Do updates the tags for the identity provider.
func (u *UpdatedIdentityProviderTagsProcedure) Do(_ context.Context) error {
	arn := u.plan.currentIdentityProvider.IdentityProviderConfigArn
	_, err := u.plan.eksClient.TagResource(&eks.TagResourceInput{
		ResourceArn: &arn,
		Tags:        aws.StringMap(u.plan.desiredIdentityProvider.Tags),
	})

	if err != nil {
		return errors.Wrap(err, "updating identity provider tags")
	}

	return nil
}

// RemoveIdentityProviderTagsProcedure removes the tags from the identity provider.
type RemoveIdentityProviderTagsProcedure struct {
	plan *plan
}

// Name returns the name of the procedure.
func (r *RemoveIdentityProviderTagsProcedure) Name() string {
	return "remove_identity_provider_tags"
}

// Do removes the tags from the identity provider.
func (r *RemoveIdentityProviderTagsProcedure) Do(_ context.Context) error {
	keys := make([]*string, 0, len(r.plan.currentIdentityProvider.Tags))

	for key := range r.plan.currentIdentityProvider.Tags {
		keys = append(keys, aws.String(key))
	}

	arn := r.plan.currentIdentityProvider.IdentityProviderConfigArn
	_, err := r.plan.eksClient.UntagResource(&eks.UntagResourceInput{
		ResourceArn: &arn,
		TagKeys:     keys,
	})

	if err != nil {
		return errors.Wrap(err, "untagging identity provider")
	}
	return nil
}
