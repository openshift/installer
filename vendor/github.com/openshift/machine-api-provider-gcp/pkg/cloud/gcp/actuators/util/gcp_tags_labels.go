/*
Copyright The Kubernetes Authors.
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

package util

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	machinecontroller "github.com/openshift/machine-api-operator/pkg/controller/machine"
	tagservice "github.com/openshift/machine-api-provider-gcp/pkg/cloud/gcp/actuators/services/tags"

	configv1 "github.com/openshift/api/config/v1"
	machinev1 "github.com/openshift/api/machine/v1beta1"

	"github.com/googleapis/gax-go/v2/apierror"

	"k8s.io/klog/v2"
	controllerclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// maxUserTagLimit is the maximum userTags that can be configured as defined in openshift/api.
	// https://github.com/openshift/api/blob/master/machine/v1beta1/types_gcpprovider.go#L153-L160
	maxUserTagLimit = 50

	// globalInfrastructureName is the default name of the Infrastructure object
	globalInfrastructureName = "cluster"

	// ocpDefaultLabelFmt is the format string for the default label
	// added to the OpenShift created GCP resources.
	ocpDefaultLabelFmt = "kubernetes-io-cluster-%s"
)

// GetInfrastructure returns the Infrastructure object infrastructure/cluster or empty
// on encountering any error.
func GetInfrastructure(client controllerclient.Client) (*configv1.Infrastructure, error) {
	infra := &configv1.Infrastructure{}
	infraName := controllerclient.ObjectKey{Name: globalInfrastructureName}

	if err := client.Get(context.Background(), infraName, infra); err != nil {
		return nil, fmt.Errorf("failed to get infrastructure: %w", err)
	}
	return infra, nil
}

// getInfraResourceLabels returns the user-defined labels present in the
// status sub-resource of Infrastructure.
func getInfraResourceLabels(platformStatus *configv1.PlatformStatus) (labels map[string]string) {
	if platformStatus != nil && platformStatus.GCP != nil && platformStatus.GCP.ResourceLabels != nil {
		labels = make(map[string]string, len(platformStatus.GCP.ResourceLabels))
		for _, label := range platformStatus.GCP.ResourceLabels {
			labels[label.Key] = label.Value
		}
	}
	return
}

// getOCPLabels returns the OCP specific labels to be added to the resources.
func getOCPLabels(clusterID string) map[string]string {
	return map[string]string{
		fmt.Sprintf(ocpDefaultLabelFmt, clusterID): "owned",
	}
}

// mergeLabels is for merging OCP specific labels, labels defined in Infrastructure.Status and
// GCPMachineProviderSpec with OCP, GCPMachineProviderSpec, Infrastructure labels precedence order.
func mergeLabels(ocpLabels, providerSpecLabels, infraLabels map[string]string) map[string]string {
	labels := make(map[string]string)

	if infraLabels != nil {
		// copy user defined labels in platform.Infrastructure.Status.
		for k, v := range infraLabels {
			labels[k] = v
		}
	}

	if providerSpecLabels != nil {
		// merge labels present in Infrastructure.Status with
		// the labels configured in GCPMachineProviderSpec, with
		// precedence given to those in GCPMachineProviderSpec
		// for new or updated labels.
		for k, v := range providerSpecLabels {
			labels[k] = v
		}
	}

	if ocpLabels != nil {
		// copy OCP labels, overwrite any OCP reserved labels found in
		// the user defined label list.
		for k, v := range ocpLabels {
			labels[k] = v
		}
	}

	return labels
}

// GetLabelsList returns the merged list of user-defined labels in Infrastructure.Status
// and GCPMachineProviderSpec to apply on the resources.
func GetLabelsList(client controllerclient.Client, machineClusterID string, providerSpecLabels map[string]string) (map[string]string, error) {
	ocpLabels := getOCPLabels(machineClusterID)

	infra, err := GetInfrastructure(client)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster infrastructure: %w", err)
	}

	infraLabels := getInfraResourceLabels(infra.Status.PlatformStatus)
	labels := mergeLabels(ocpLabels, providerSpecLabels, infraLabels)

	if len(ocpLabels) > 32 || (len(labels)-len(ocpLabels)) > 32 {
		return nil, fmt.Errorf("ocp can define upto 32 labels and user can define upto 32 labels,"+
			"infrstructure.status.resourceLabels and Machine.Spec.ProviderSpec.Labels put together configured label count is %d", len(labels))
	}

	return labels, nil
}

// getInfraResourceTagsList returns the user-defined tags present in the
// status sub-resource of Infrastructure.
func getInfraResourceTagsList(platformStatus *configv1.PlatformStatus) []machinev1.ResourceManagerTag {
	if platformStatus == nil || platformStatus.GCP == nil || platformStatus.GCP.ResourceTags == nil {
		return nil
	}

	tags := make([]machinev1.ResourceManagerTag, len(platformStatus.GCP.ResourceTags))
	for i, tag := range platformStatus.GCP.ResourceTags {
		tags[i] = machinev1.ResourceManagerTag{
			ParentID: tag.ParentID,
			Key:      tag.Key,
			Value:    tag.Value,
		}
	}

	return tags
}

// getTagValuesNames returns the list of tags in Compute APIs required format, which
// is a map containing keys of the form Key(`tagKeys/{tag_key_id}`) and values in the
// form (`tagValues/{tag_value_id}`).
func getTagValuesNames(ctx context.Context, tagService tagservice.TagService, tagList []machinev1.ResourceManagerTag) (map[string]string, error) {
	// identify tags which are inaccessible due to permissions issues
	// or does not exist and report back to user to fix in one go.
	inaccessibleTags := make([]string, 0)
	tagValueList := make(map[string]string, len(tagList))

	for _, tag := range tagList {
		name := fmt.Sprintf("%s/%s/%s", tag.ParentID, tag.Key, tag.Value)
		value, err := tagService.GetNamespacedName(ctx, name)
		if err != nil {
			var gErr *apierror.APIError
			// google API returns StatusForbidden or StatusNotFound when the tag
			// does not exist, since it could be because of permission issues
			// or genuinely tag does not exist.
			if errors.As(err, &gErr) && (gErr.HTTPCode() == http.StatusNotFound ||
				gErr.HTTPCode() == http.StatusForbidden) {
				klog.Errorf("does not have permission to access %s tag or tag does not exist", name)
				inaccessibleTags = append(inaccessibleTags, name)
				continue
			}
			// fetching tag's metadata could fail due to errors like timeout, server
			// internal errors, permission issues among others. Since tag's key and
			// value names are required for binding tag to compute resource, will
			// return error and retry during next reconciliation.
			return nil, fmt.Errorf("failed to fetch %s tag details: %w", name, err)
		}
		tagValueList[value.Parent] = value.Name
	}

	if len(inaccessibleTags) != 0 {
		return nil, machinecontroller.InvalidMachineConfiguration("%v tag(s) do not exist or does not have required permission to access", inaccessibleTags)
	}

	return tagValueList, nil
}

// mergeInfraProviderSpecTags merges user-defined tags in Infrastructure.Status and
// GCPMachineProviderSpec, with precedence given to those in GCPMachineProviderSpec
// for new or updated tags.
func mergeInfraProviderSpecTags(infraTags []machinev1.ResourceManagerTag, providerSpecTags []machinev1.ResourceManagerTag) []machinev1.ResourceManagerTag {
	mergedTags := make([]machinev1.ResourceManagerTag, 0, len(infraTags))

	for _, tag := range providerSpecTags {
		mergedTags = append(mergedTags, tag)
	}

	for _, iTag := range infraTags {
		appendTag := true
		for _, pTag := range providerSpecTags {
			if iTag.ParentID == pTag.ParentID && iTag.Key == pTag.Key {
				appendTag = false
				break
			}
		}
		if appendTag {
			mergedTags = append(mergedTags, iTag)
		}
	}

	return mergedTags
}

// GetResourceManagerTags returns the merged list of user-defined tags in Infrastructure.Status
// and GCPMachineProviderSpec to apply on the resources.
func GetResourceManagerTags(ctx context.Context,
	client controllerclient.Client,
	tagService tagservice.TagService,
	providerSpecTags []machinev1.ResourceManagerTag) (map[string]string, error) {
	infra, err := GetInfrastructure(client)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster infrastructure: %w", err)
	}
	userTags := getInfraResourceTagsList(infra.Status.PlatformStatus)

	if len(userTags) == 0 && len(providerSpecTags) == 0 {
		klog.V(3).Infof("user-defined tags in infrastructure and machineProviderSpec is empty")
		return nil, nil
	}

	mergedTags := mergeInfraProviderSpecTags(userTags, providerSpecTags)

	if len(mergedTags) > maxUserTagLimit {
		return nil, fmt.Errorf("maximum of %d tags can be added to a compute instance, "+
			"infrastructure.status.resourceTags and machines.spec.providerSpec.resourceManagerTags "+
			"combined tag count is %d", maxUserTagLimit, len(mergedTags))
	}

	return getTagValuesNames(ctx, tagService, mergedTags)
}
