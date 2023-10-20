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

package iamauth

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	iamauthv1 "sigs.k8s.io/aws-iam-authenticator/pkg/mapper/crd/apis/iamauthenticator/v1alpha1"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
)

type crdBackend struct {
	client crclient.Client
}

func (b *crdBackend) MapRole(mapping ekscontrolplanev1.RoleMapping) error {
	ctx := context.TODO()

	if errs := mapping.Validate(); errs != nil {
		return kerrors.NewAggregate(errs)
	}

	mappingList := iamauthv1.IAMIdentityMappingList{}

	if err := b.client.List(ctx, &mappingList); err != nil {
		return fmt.Errorf("getting list of mappings: %w", err)
	}

	for _, existingMapping := range mappingList.Items {
		existing := existingMapping
		if roleMappingMatchesIAMMap(mapping, &existing) {
			// We already have a mapping so do nothing
			return nil
		}
	}

	iamMapping := &iamauthv1.IAMIdentityMapping{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    metav1.NamespaceSystem,
			GenerateName: "capa-iamauth-",
		},
		Spec: iamauthv1.IAMIdentityMappingSpec{
			ARN:      mapping.RoleARN,
			Username: mapping.UserName,
			Groups:   mapping.Groups,
		},
	}

	return b.client.Create(ctx, iamMapping)
}

func (b *crdBackend) MapUser(mapping ekscontrolplanev1.UserMapping) error {
	ctx := context.TODO()

	if errs := mapping.Validate(); errs != nil {
		return kerrors.NewAggregate(errs)
	}

	mappingList := iamauthv1.IAMIdentityMappingList{}

	if err := b.client.List(ctx, &mappingList); err != nil {
		return fmt.Errorf("getting list of mappings: %w", err)
	}

	for _, existingMapping := range mappingList.Items {
		existing := existingMapping
		if userMappingMatchesIAMMap(mapping, &existing) {
			// We already have a mapping so do nothing
			return nil
		}
	}

	iamMapping := &iamauthv1.IAMIdentityMapping{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    metav1.NamespaceSystem,
			GenerateName: "capa-iamauth-",
		},
		Spec: iamauthv1.IAMIdentityMappingSpec{
			ARN:      mapping.UserARN,
			Username: mapping.UserName,
			Groups:   mapping.Groups,
		},
	}

	return b.client.Create(ctx, iamMapping)
}

func roleMappingMatchesIAMMap(mapping ekscontrolplanev1.RoleMapping, iamMapping *iamauthv1.IAMIdentityMapping) bool {
	if mapping.RoleARN != iamMapping.Spec.ARN {
		return false
	}

	if mapping.UserName != iamMapping.Spec.Username {
		return false
	}

	if len(mapping.Groups) != len(iamMapping.Spec.Groups) {
		return false
	}

	for _, mappingGroup := range mapping.Groups {
		found := false
		for _, iamGroup := range iamMapping.Spec.Groups {
			if iamGroup == mappingGroup {
				found = true
			}
		}
		if !found {
			return false
		}
	}

	return true
}

func userMappingMatchesIAMMap(mapping ekscontrolplanev1.UserMapping, iamMapping *iamauthv1.IAMIdentityMapping) bool {
	if mapping.UserARN != iamMapping.Spec.ARN {
		return false
	}

	if mapping.UserName != iamMapping.Spec.Username {
		return false
	}

	if len(mapping.Groups) != len(iamMapping.Spec.Groups) {
		return false
	}

	for _, mappingGroup := range mapping.Groups {
		found := false
		for _, iamGroup := range iamMapping.Spec.Groups {
			if iamGroup == mappingGroup {
				found = true
			}
		}
		if !found {
			return false
		}
	}

	return true
}
