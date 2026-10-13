/*
Copyright 2020 The Kubernetes Authors.

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

// Package iam provides a service for managing IAM roles and policies.
package iam

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/converters"
	iamv1 "sigs.k8s.io/cluster-api-provider-aws/v2/iam/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/iamauth"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
)

const (
	// EKSFargateService is the service to trust for fargate pod execution roles.
	EKSFargateService = "eks-fargate-pods.amazonaws.com"
)

// IAMService defines the specs for an IAM service.
type IAMService struct {
	logger.Wrapper
	IAMClient iamauth.IAMAPI
	Client    *http.Client
}

// GetIAMRole will return the IAM role for the IAMService.
func (s *IAMService) GetIAMRole(ctx context.Context, name string) (*iamtypes.Role, error) {
	input := &iam.GetRoleInput{
		RoleName: aws.String(name),
	}

	out, err := s.IAMClient.GetRole(ctx, input)
	if err != nil {
		return nil, err
	}

	return out.Role, nil
}

func (s *IAMService) getIAMPolicy(ctx context.Context, policyArn string) (*iamtypes.Policy, error) {
	input := &iam.GetPolicyInput{
		PolicyArn: &policyArn,
	}

	out, err := s.IAMClient.GetPolicy(ctx, input)
	if err != nil {
		return nil, err
	}

	return out.Policy, nil
}

func (s *IAMService) getIAMRolePolicies(ctx context.Context, roleName string) ([]string, error) {
	input := &iam.ListAttachedRolePoliciesInput{
		RoleName: &roleName,
	}

	out, err := s.IAMClient.ListAttachedRolePolicies(ctx, input)
	if err != nil {
		return nil, errors.Wrapf(err, "error listing role polices for %s", roleName)
	}

	policies := []string{}
	for _, policy := range out.AttachedPolicies {
		if policy.PolicyArn != nil {
			policies = append(policies, *policy.PolicyArn)
		}
	}

	return policies, nil
}

func (s *IAMService) detachIAMRolePolicy(ctx context.Context, roleName string, policyARN string) error {
	input := &iam.DetachRolePolicyInput{
		RoleName:  aws.String(roleName),
		PolicyArn: aws.String(policyARN),
	}

	if _, err := s.IAMClient.DetachRolePolicy(ctx, input); err != nil {
		return errors.Wrapf(err, "error detaching policy %s from role %s", policyARN, roleName)
	}

	return nil
}

func (s *IAMService) attachIAMRolePolicy(ctx context.Context, roleName string, policyARN string) error {
	input := &iam.AttachRolePolicyInput{
		RoleName:  aws.String(roleName),
		PolicyArn: aws.String(policyARN),
	}

	if _, err := s.IAMClient.AttachRolePolicy(ctx, input); err != nil {
		return errors.Wrapf(err, "error attaching policy %s to role %s", policyARN, roleName)
	}

	return nil
}

// EnsurePoliciesAttached will ensure the IAMService has policies attached.
func (s *IAMService) EnsurePoliciesAttached(ctx context.Context, role *iamtypes.Role, policies []string) (bool, error) {
	s.Debug("Ensuring Polices are attached to role")
	existingPolices, err := s.getIAMRolePolicies(ctx, *role.RoleName)
	if err != nil {
		return false, err
	}

	var updatedPolicies bool
	// Remove polices that aren't in the list
	for _, existingPolicy := range existingPolices {
		found := findStringInSlice(policies, existingPolicy)
		if !found {
			updatedPolicies = true
			err = s.detachIAMRolePolicy(ctx, *role.RoleName, existingPolicy)
			if err != nil {
				return false, err
			}
			s.Debug("Detached policy from role", "role", role.RoleName, "policy", existingPolicy)
		}
	}

	// Add any policies that aren't currently attached
	for _, policy := range policies {
		found := findStringInSlice(existingPolices, policy)
		if !found {
			// Make sure policy exists before attaching
			_, err := s.getIAMPolicy(ctx, policy)
			if err != nil {
				return false, errors.Wrapf(err, "error getting policy %s", policy)
			}

			updatedPolicies = true
			err = s.attachIAMRolePolicy(ctx, *role.RoleName, policy)
			if err != nil {
				return false, err
			}
			s.Debug("Attached policy to role", "role", role.RoleName, "policy", policy)
		}
	}

	return updatedPolicies, nil
}

// RoleTags returns the tags for the given role.
func RoleTags(key string, additionalTags infrav1.Tags) []iamtypes.Tag {
	additionalTags[infrav1.ClusterAWSCloudProviderTagKey(key)] = string(infrav1.ResourceLifecycleOwned)
	tags := []iamtypes.Tag{}
	for k, v := range additionalTags {
		tags = append(tags, iamtypes.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}

	// Sort so that unit tests can expect a stable order
	sort.Slice(tags, func(i, j int) bool { return *tags[i].Key < *tags[j].Key })

	return tags
}

// CreateRole will create a role from the IAMService.
func (s *IAMService) CreateRole(
	ctx context.Context,
	roleName string,
	key string,
	trustRelationship *iamv1.PolicyDocument,
	additionalTags infrav1.Tags,
	path string,
	permissionsBoundary string,
) (*iamtypes.Role, error) {
	tags := RoleTags(key, additionalTags)

	trustRelationshipJSON, err := converters.IAMPolicyDocumentToJSON(*trustRelationship)
	if err != nil {
		return nil, errors.Wrap(err, "error converting trust relationship to json")
	}

	input := &iam.CreateRoleInput{
		RoleName:                 aws.String(roleName),
		Tags:                     tags,
		AssumeRolePolicyDocument: aws.String(trustRelationshipJSON),
	}

	if len(path) > 0 {
		input.Path = aws.String(path)
	}

	if len(permissionsBoundary) > 0 {
		input.PermissionsBoundary = aws.String(permissionsBoundary)
	}

	out, err := s.IAMClient.CreateRole(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call CreateRole")
	}

	return out.Role, nil
}

// EnsureTagsAndPolicy will ensure any tags and policies against the IAMService.
func (s *IAMService) EnsureTagsAndPolicy(
	ctx context.Context,
	role *iamtypes.Role,
	key string,
	trustRelationship *iamv1.PolicyDocument,
	additionalTags infrav1.Tags,
) (bool, error) {
	s.Debug("Ensuring tags and AssumeRolePolicyDocument are set on role")

	rolePolicyDocumentRaw, err := url.PathUnescape(*role.AssumeRolePolicyDocument)
	if err != nil {
		return false, errors.Wrap(err, "couldn't decode AssumeRolePolicyDocument")
	}

	var rolePolicyDocument iamv1.PolicyDocument
	err = json.Unmarshal([]byte(rolePolicyDocumentRaw), &rolePolicyDocument)
	if err != nil {
		return false, errors.Wrap(err, "couldn't unmarshal AssumeRolePolicyDocument")
	}

	var updated bool
	if !cmp.Equal(*trustRelationship, rolePolicyDocument) {
		trustRelationshipJSON, err := converters.IAMPolicyDocumentToJSON(*trustRelationship)
		if err != nil {
			return false, errors.Wrap(err, "error converting trust relationship to json")
		}
		policyInput := &iam.UpdateAssumeRolePolicyInput{
			RoleName:       role.RoleName,
			PolicyDocument: aws.String(trustRelationshipJSON),
		}
		updated = true
		if _, err := s.IAMClient.UpdateAssumeRolePolicy(ctx, policyInput); err != nil {
			return updated, err
		}
	}

	tagInput := &iam.TagRoleInput{
		RoleName: role.RoleName,
	}
	untagInput := &iam.UntagRoleInput{
		RoleName: role.RoleName,
	}
	currentTags := make(map[string]string)
	for _, tag := range role.Tags {
		currentTags[*tag.Key] = *tag.Value
		if *tag.Key == infrav1.ClusterAWSCloudProviderTagKey(key) {
			continue
		}
		if _, ok := additionalTags[*tag.Key]; !ok {
			untagInput.TagKeys = append(untagInput.TagKeys, *tag.Key)
		}
	}
	for key, value := range additionalTags {
		if currentV, ok := currentTags[key]; !ok || value != currentV {
			tagInput.Tags = append(tagInput.Tags, iamtypes.Tag{
				Key:   aws.String(key),
				Value: aws.String(value),
			})
		}
	}

	if len(tagInput.Tags) > 0 {
		updated = true
		_, err = s.IAMClient.TagRole(ctx, tagInput)
		if err != nil {
			return updated, err
		}
	}

	if len(untagInput.TagKeys) > 0 {
		updated = true
		_, err = s.IAMClient.UntagRole(ctx, untagInput)
		if err != nil {
			return updated, err
		}
	}

	return updated, nil
}

func (s *IAMService) detachAllPoliciesForRole(ctx context.Context, name string) error {
	s.Debug("Detaching all policies for role", "role", name)
	input := &iam.ListAttachedRolePoliciesInput{
		RoleName: &name,
	}
	policies, err := s.IAMClient.ListAttachedRolePolicies(ctx, input)
	if err != nil {
		return errors.Wrapf(err, "error fetching policies for role %s", name)
	}
	for _, p := range policies.AttachedPolicies {
		s.Debug("Detaching policy", "policy", *p.PolicyArn)
		if err := s.detachIAMRolePolicy(ctx, name, *p.PolicyArn); err != nil {
			return err
		}
	}
	return nil
}

// DeleteRole will delete a role from the IAMService.
func (s *IAMService) DeleteRole(ctx context.Context, name string) error {
	if err := s.detachAllPoliciesForRole(ctx, name); err != nil {
		return errors.Wrapf(err, "error detaching policies for role %s", name)
	}

	input := &iam.DeleteRoleInput{
		RoleName: aws.String(name),
	}

	if _, err := s.IAMClient.DeleteRole(ctx, input); err != nil {
		return errors.Wrapf(err, "error deleting role %s", name)
	}

	return nil
}

// IsUnmanaged will check if a given role and tag are unmanaged against the IAMService.
func (s *IAMService) IsUnmanaged(role *iamtypes.Role, key string) bool {
	keyToFind := infrav1.ClusterAWSCloudProviderTagKey(key)
	for _, tag := range role.Tags {
		if *tag.Key == keyToFind && *tag.Value == string(infrav1.ResourceLifecycleOwned) {
			return false
		}
	}

	return true
}

// ControlPlaneTrustRelationship will generate a ControlPlane PolicyDocument.
func ControlPlaneTrustRelationship(enableFargate bool) *iamv1.PolicyDocument {
	identity := make(iamv1.Principals)
	identity["Service"] = []string{"eks.amazonaws.com"}
	if enableFargate {
		identity["Service"] = append(identity["Service"], EKSFargateService)
	}

	policy := &iamv1.PolicyDocument{
		Version: "2012-10-17",
		Statement: []iamv1.StatementEntry{
			{
				Effect: "Allow",
				Action: []string{
					"sts:AssumeRole",
				},
				Principal: identity,
			},
		},
	}

	return policy
}

// FargateTrustRelationship will generate a Fargate PolicyDocument.
func FargateTrustRelationship() *iamv1.PolicyDocument {
	identity := make(iamv1.Principals)
	identity["Service"] = []string{EKSFargateService}

	policy := &iamv1.PolicyDocument{
		Version: "2012-10-17",
		Statement: []iamv1.StatementEntry{
			{
				Effect: "Allow",
				Action: []string{
					"sts:AssumeRole",
				},
				Principal: identity,
			},
		},
	}

	return policy
}

// NodegroupTrustRelationship will generate a Nodegroup PolicyDocument.
func NodegroupTrustRelationship() *iamv1.PolicyDocument {
	identity := make(iamv1.Principals)
	identity["Service"] = []string{"ec2.amazonaws.com"}

	policy := &iamv1.PolicyDocument{
		Version: "2012-10-17",
		Statement: []iamv1.StatementEntry{
			{
				Effect: "Allow",
				Action: []string{
					"sts:AssumeRole",
				},
				Principal: identity,
			},
		},
	}

	return policy
}

func findStringInSlice(slice []string, toFind string) bool {
	for _, item := range slice {
		if item == toFind {
			return true
		}
	}

	return false
}

const stsAWSAudience = "sts.amazonaws.com"

// CreateOIDCProvider will create an OIDC provider.
func (s *IAMService) CreateOIDCProvider(ctx context.Context, cluster *ekstypes.Cluster) (string, error) {
	issuerURL, err := url.Parse(*cluster.Identity.Oidc.Issuer)
	if err != nil {
		return "", err
	}
	if issuerURL.Scheme != "https" {
		return "", errors.Errorf("invalid scheme for issuer URL %s", issuerURL.String())
	}

	thumbprint, err := fetchRootCAThumbprint(ctx, issuerURL.String(), s.Client)
	if err != nil {
		return "", err
	}
	input := iam.CreateOpenIDConnectProviderInput{
		ClientIDList:   []string{stsAWSAudience},
		ThumbprintList: []string{thumbprint},
		Url:            aws.String(issuerURL.String()),
	}
	provider, err := s.IAMClient.CreateOpenIDConnectProvider(ctx, &input)
	if err != nil {
		return "", errors.Wrap(err, "error creating provider")
	}
	return *provider.OpenIDConnectProviderArn, nil
}

// FindAndVerifyOIDCProvider will try to find an OIDC provider. It will return an error if the found provider does not
// match the cluster spec.
func (s *IAMService) FindAndVerifyOIDCProvider(ctx context.Context, cluster *ekstypes.Cluster) (string, error) {
	issuerURL, err := url.Parse(*cluster.Identity.Oidc.Issuer)
	if err != nil {
		return "", err
	}
	if issuerURL.Scheme != "https" {
		return "", errors.Errorf("invalid scheme for issuer URL %s", issuerURL.String())
	}

	thumbprint, err := fetchRootCAThumbprint(ctx, issuerURL.String(), s.Client)
	if err != nil {
		return "", err
	}
	output, err := s.IAMClient.ListOpenIDConnectProviders(ctx, &iam.ListOpenIDConnectProvidersInput{})
	if err != nil {
		return "", errors.Wrap(err, "error listing providers")
	}
	for _, r := range output.OpenIDConnectProviderList {
		provider, err := s.IAMClient.GetOpenIDConnectProvider(ctx, &iam.GetOpenIDConnectProviderInput{OpenIDConnectProviderArn: r.Arn})
		if err != nil {
			return "", errors.Wrap(err, "error getting provider")
		}
		// URL should always contain `https`.
		if *provider.Url != issuerURL.String() && *provider.Url != strings.Replace(issuerURL.String(), "https://", "", 1) {
			continue
		}
		if len(provider.ThumbprintList) != 1 || provider.ThumbprintList[0] != thumbprint {
			return "", errors.Wrap(err, "found provider with matching issuerURL but with non-matching thumbprint")
		}
		if len(provider.ClientIDList) != 1 || provider.ClientIDList[0] != stsAWSAudience {
			return "", errors.Wrap(err, "found provider with matching issuerURL but with non-matching clientID")
		}
		return *r.Arn, nil
	}
	return "", nil
}

func fetchRootCAThumbprint(ctx context.Context, issuerURL string, client *http.Client) (string, error) {
	// needed to appease noctx.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, issuerURL, http.NoBody)
	if err != nil {
		return "", err
	}

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	rootCA := response.TLS.PeerCertificates[len(response.TLS.PeerCertificates)-1]
	sha1Sum := sha1.Sum(rootCA.Raw) //nolint:gosec
	return hex.EncodeToString(sha1Sum[:]), nil
}

// DeleteOIDCProvider will delete an OIDC provider.
func (s *IAMService) DeleteOIDCProvider(ctx context.Context, arn *string) error {
	input := iam.DeleteOpenIDConnectProviderInput{
		OpenIDConnectProviderArn: arn,
	}

	_, err := s.IAMClient.DeleteOpenIDConnectProvider(ctx, &input)
	if err != nil {
		return errors.Wrap(err, "error deleting provider")
	}
	return nil
}
