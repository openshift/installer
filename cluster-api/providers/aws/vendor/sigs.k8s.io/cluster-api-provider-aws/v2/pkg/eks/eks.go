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

// Package eks contains the EKS API implementation.
package eks

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/pkg/errors"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/hash"
)

const (
	resourcePrefix = "capa_"
)

// Client implements EKSAPI as it can not be imported from pkg/cloud/services/eks/service.go due to import cycle.
type Client interface {
	AssociateIdentityProviderConfig(ctx context.Context, params *eks.AssociateIdentityProviderConfigInput, optFns ...func(*eks.Options)) (*eks.AssociateIdentityProviderConfigOutput, error)
	DisassociateIdentityProviderConfig(ctx context.Context, params *eks.DisassociateIdentityProviderConfigInput, optFns ...func(*eks.Options)) (*eks.DisassociateIdentityProviderConfigOutput, error)
	DescribeIdentityProviderConfig(ctx context.Context, params *eks.DescribeIdentityProviderConfigInput, optFns ...func(*eks.Options)) (*eks.DescribeIdentityProviderConfigOutput, error)
	TagResource(ctx context.Context, params *eks.TagResourceInput, optFns ...func(*eks.Options)) (*eks.TagResourceOutput, error)
	UntagResource(ctx context.Context, params *eks.UntagResourceInput, optFns ...func(*eks.Options)) (*eks.UntagResourceOutput, error)
	CreateAddon(ctx context.Context, params *eks.CreateAddonInput, optFns ...func(*eks.Options)) (*eks.CreateAddonOutput, error)
	DeleteAddon(ctx context.Context, params *eks.DeleteAddonInput, optFns ...func(*eks.Options)) (*eks.DeleteAddonOutput, error)
	UpdateAddon(ctx context.Context, params *eks.UpdateAddonInput, optFns ...func(*eks.Options)) (*eks.UpdateAddonOutput, error)
	DescribeAddon(ctx context.Context, params *eks.DescribeAddonInput, optFns ...func(*eks.Options)) (*eks.DescribeAddonOutput, error)
	WaitUntilAddonDeleted(ctx context.Context, params *eks.DescribeAddonInput, maxWait time.Duration) error
}

// GenerateEKSName generates a name of an EKS resources.
func GenerateEKSName(resourceName, namespace string, maxLength int) (string, error) {
	escapedName := strings.ReplaceAll(resourceName, ".", "_")
	eksName := fmt.Sprintf("%s_%s", namespace, escapedName)

	if len(eksName) < maxLength {
		return eksName, nil
	}

	hashLength := 32 - len(resourcePrefix)
	hashedName, err := hash.Base36TruncatedHash(eksName, hashLength)
	if err != nil {
		return "", errors.Wrap(err, "creating hash from name")
	}

	return fmt.Sprintf("%s%s", resourcePrefix, hashedName), nil
}
