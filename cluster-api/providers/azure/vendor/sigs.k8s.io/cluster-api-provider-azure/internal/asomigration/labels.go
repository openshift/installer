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

package asomigration

import (
	"context"
	"errors"
	"fmt"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	asoCRDAppLabel              = "app.kubernetes.io/name"
	asoCRDAppValue              = "azure-service-operator"
	infrastructureProviderLabel = "infrastructure-azure"
)

// LabelCRDsForClusterctlUpgrade labels ASO-managed CRDs as CAPZ provider resources.
func LabelCRDsForClusterctlUpgrade(ctx context.Context, c client.Client) error {
	crds := &apiextensionsv1.CustomResourceDefinitionList{}
	if err := c.List(ctx, crds, client.MatchingLabels{asoCRDAppLabel: asoCRDAppValue}); err != nil {
		return fmt.Errorf("listing ASO-managed CRDs: %w", err)
	}
	if len(crds.Items) == 0 {
		return errors.New("no ASO-managed CRDs found")
	}

	for i := range crds.Items {
		crd := &crds.Items[i]
		if crd.Labels[clusterv1.ProviderNameLabel] == infrastructureProviderLabel {
			continue
		}

		base := crd.DeepCopy()
		if crd.Labels == nil {
			crd.Labels = map[string]string{}
		}
		crd.Labels[clusterv1.ProviderNameLabel] = infrastructureProviderLabel
		if err := c.Patch(ctx, crd, client.MergeFrom(base)); err != nil {
			return fmt.Errorf("labeling ASO-managed CRD %s: %w", crd.Name, err)
		}
	}

	verified := &apiextensionsv1.CustomResourceDefinitionList{}
	if err := c.List(ctx, verified, client.MatchingLabels{asoCRDAppLabel: asoCRDAppValue}); err != nil {
		return fmt.Errorf("verifying ASO-managed CRD labels: %w", err)
	}
	for i := range verified.Items {
		crd := &verified.Items[i]
		if crd.Labels[clusterv1.ProviderNameLabel] != infrastructureProviderLabel {
			return fmt.Errorf("ASO-managed CRD %s is missing provider label", crd.Name)
		}
	}

	return nil
}
