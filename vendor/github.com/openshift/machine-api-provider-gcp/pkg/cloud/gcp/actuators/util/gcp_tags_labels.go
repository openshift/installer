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
	"fmt"

	configv1 "github.com/openshift/api/config/v1"
	controllerclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
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
func GetLabelsList(userLabelsAllowed bool, client controllerclient.Client, machineClusterID string, providerSpecLabels map[string]string) (map[string]string, error) {
	ocpLabels := getOCPLabels(machineClusterID)

	if !userLabelsAllowed {
		return mergeLabels(ocpLabels, providerSpecLabels, nil), nil
	}

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
