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

package metadata

import (
	"context"

	"github.com/vmware/govmomi/vapi/tags"
	ctrl "sigs.k8s.io/controller-runtime"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/session"
)

type metadataContext interface {
	context.Context

	GetSession() *session.Session
}

func getCategoryAssociableType(domainType infrav1.FailureDomainType) string {
	switch domainType {
	case infrav1.DatacenterFailureDomain:
		return "Datacenter"
	case infrav1.ComputeClusterFailureDomain:
		return "ClusterComputeResource"
	case infrav1.HostGroupFailureDomain:
		return "HostSystem"
	default:
		return ""
	}
}

// CreateCategory either creates a new vSphere category or updates the associable type for an existing category.
func CreateCategory(ctx metadataContext, name string, failureDomainType infrav1.FailureDomainType) (string, error) {
	logger := ctrl.LoggerFrom(ctx, "category", name)
	manager := ctx.GetSession().TagManager
	category, err := manager.GetCategory(ctx, name)
	if err != nil {
		logger.V(4).Info("failed to find existing category, creating a new category")
		id, err := manager.CreateCategory(ctx, getCategoryObject(name, failureDomainType))
		if err != nil {
			return "", err
		}
		return id, nil
	}
	category.Patch(getCategoryObject(name, failureDomainType))
	if err := manager.UpdateCategory(ctx, category); err != nil {
		logger.V(4).Error(err, "failed to update existing category")
		return "", err
	}
	return category.ID, nil
}

func getCategoryObject(name string, failureDomainType infrav1.FailureDomainType) *tags.Category {
	return &tags.Category{
		Name:            name,
		Description:     "CAPV generated category for Failure Domain support",
		AssociableTypes: []string{getCategoryAssociableType(failureDomainType)},
		Cardinality:     "MULTIPLE",
	}
}

func CreateTag(ctx metadataContext, name, categoryID string) error {
	logger := ctrl.LoggerFrom(ctx, "tag", name, "category", categoryID)
	manager := ctx.GetSession().TagManager
	_, err := manager.GetTag(ctx, name)
	if err != nil {
		logger.V(4).Info("failed to find existing tag, creating a new tag")
		_, err = manager.CreateTag(ctx, &tags.Tag{
			Description: "CAPV generated tag for Failure Domain support",
			Name:        name,
			CategoryID:  categoryID,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
