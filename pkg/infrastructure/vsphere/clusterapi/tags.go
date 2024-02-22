package clusterapi

import (
	"context"
	"fmt"

	"github.com/vmware/govmomi/vapi/tags"
	"github.com/vmware/govmomi/vim25/types"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/session"
)

func attachTag(ctx context.Context, session *session.Session, vmMoRefValue, tagID string) error {
	tagManager := session.TagManager

	moRef := types.ManagedObjectReference{
		Value: vmMoRefValue,
		Type:  "VirtualMachine",
	}

	err := tagManager.AttachTag(ctx, tagID, moRef)

	if err != nil {
		return fmt.Errorf("unable to attach tag: %w", err)
	}
	return nil
}

func createClusterTagID(ctx context.Context, session *session.Session, clusterID string) (string, error) {
	tagManager := session.TagManager
	categories, err := tagManager.GetCategories(ctx)
	if err != nil {
		return "", fmt.Errorf("unable to get tag categories: %w", err)
	}

	var clusterTagCategory *tags.Category
	clusterTagCategoryName := fmt.Sprintf("openshift-%s", clusterID)
	tagCategoryID := ""

	for i, category := range categories {
		if category.Name == clusterTagCategoryName {
			clusterTagCategory = &categories[i]
			tagCategoryID = category.ID
			break
		}
	}

	if clusterTagCategory == nil {
		clusterTagCategory = &tags.Category{
			Name:        clusterTagCategoryName,
			Description: "Added by openshift-install do not remove",
			Cardinality: "SINGLE",
			AssociableTypes: []string{
				"urn:vim25:VirtualMachine",
				"urn:vim25:ResourcePool",
				"urn:vim25:Folder",
				"urn:vim25:Datastore",
				"urn:vim25:StoragePod",
			},
		}
		tagCategoryID, err = tagManager.CreateCategory(ctx, clusterTagCategory)
		if err != nil {
			return "", fmt.Errorf("unable to create tag category: %w", err)
		}
	}

	var categoryTag *tags.Tag
	tagID := ""

	categoryTags, err := tagManager.GetTagsForCategory(ctx, tagCategoryID)
	if err != nil {
		return "", fmt.Errorf("unable to get tags for category: %w", err)
	}
	for i, tag := range categoryTags {
		if tag.Name == clusterID {
			categoryTag = &categoryTags[i]
			tagID = tag.ID
			break
		}
	}

	if categoryTag == nil {
		categoryTag = &tags.Tag{
			Description: "Added by openshift-install do not remove",
			Name:        clusterID,
			CategoryID:  tagCategoryID,
		}
		tagID, err = tagManager.CreateTag(ctx, categoryTag)
		if err != nil {
			return "", fmt.Errorf("unable to create tag: %w", err)
		}
	}

	return tagID, nil
}
