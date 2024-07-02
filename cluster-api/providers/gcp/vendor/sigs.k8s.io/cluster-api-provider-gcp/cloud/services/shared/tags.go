/*
Copyright 2023 The Kubernetes Authors.

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

package shared

import (
	"context"
	"fmt"

	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	rmpb "cloud.google.com/go/resourcemanager/apiv3/resourcemanagerpb"
	infrav1 "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	infrav1exp "sigs.k8s.io/cluster-api-provider-gcp/exp/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ResourceTagBinding creates a TagBinding between a TagValue and a Google Cloud resource.
// If any of the SDK calls fail, the error is logged and no action is taken.
func ResourceTagBinding(ctx context.Context, client *resourcemanager.TagBindingsClient, spec infrav1exp.GCPManagedClusterSpec, name string) error {
	for _, tag := range spec.ResourceManagerTags {
		tagValue, err := getTagValues(ctx, tag)
		if err != nil {
			return fmt.Errorf("failed to retrieve tag value: %w", err)
		}
		req := &rmpb.CreateTagBindingRequest{
			TagBinding: &rmpb.TagBinding{
				Parent:   fmt.Sprintf("//container.googleapis.com/projects/%s/locations/%s/clusters/%s", tag.ParentID, spec.Region, name),
				TagValue: tagValue.GetName(),
			},
		}
		op, err := client.CreateTagBinding(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to create tag binding: %w", err)
		}

		_, err = op.Wait(ctx)
		if err != nil {
			return fmt.Errorf("tag binding operation failed: %w", err)
		}
	}

	return nil
}

// ResourceTagConvert converts the passed resource-manager tags to a GCP API valid format.
// Tag keys and Tag Values will be created by the user and only the Tag bindings to the Compute Instance will be
// handled by CAPG. If the Tag Key/Tag Value cannot be retrieved or no tags are provided, this will be empty and no tags will be added.
func ResourceTagConvert(ctx context.Context, t infrav1.ResourceManagerTags) infrav1.ResourceManagerTagsMap {
	tagValueList := make(infrav1.ResourceManagerTagsMap, len(t))
	log := log.FromContext(ctx)
	if len(t) == 0 {
		return tagValueList
	}

	for _, tag := range t {
		tagValue, err := getTagValues(ctx, tag)
		if err != nil {
			log.Error(err, "failed to retrieve tag value")
			continue
		}
		tagValueList[tagValue.GetParent()] = tagValue.GetName()
	}

	return tagValueList
}

func getTagValues(ctx context.Context, tag infrav1.ResourceManagerTag) (*rmpb.TagValue, error) {
	log := log.FromContext(ctx)
	client, err := resourcemanager.NewTagValuesClient(ctx)
	if err != nil {
		log.Error(err, "failed to create tag values client")
		return &rmpb.TagValue{}, err
	}
	defer client.Close()

	req := &rmpb.GetNamespacedTagValueRequest{
		Name: fmt.Sprintf("%s/%s/%s", tag.ParentID, tag.Key, tag.Value),
	}
	tagValue, err := client.GetNamespacedTagValue(ctx, req)
	if err != nil {
		return &rmpb.TagValue{}, err
	}

	return tagValue, nil
}
