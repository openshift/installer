/*
Copyright 2025 The Kubernetes Authors.

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

package instancetemplates

import (
	"context"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/filter"
	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/meta"
	"google.golang.org/api/compute/v1"

	"sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud/gcperrors"
	"sigs.k8s.io/cluster-api-provider-gcp/pkg/gcp"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// Reconcile reconciles GCP instanceTemplates.
func (s *Service) Reconcile(ctx context.Context) (*meta.Key, error) {
	log := log.FromContext(ctx)
	log.Info("Reconciling instanceTemplate resources")
	instanceTemplate, instanceTemplateKey, err := s.createOrGetInstanceTemplate(ctx)
	if err != nil {
		return nil, err
	}
	log.V(2).Info("binding to instanceTemplate", "selfLink", instanceTemplate.SelfLink)

	return instanceTemplateKey, nil
}

// Delete deletes the GCP instanceTemplate.
func (s *Service) Delete(ctx context.Context) error {
	log := log.FromContext(ctx)

	baseKey, err := s.scope.BaseInstanceTemplateResourceName()
	if err != nil {
		return err
	}

	selfLink := gcp.FormatKey("instanceTemplates", baseKey)
	log = log.WithValues("instanceTemplatesPrefix", selfLink)

	log.Info("Deleting instanceTemplate resources")

	log.V(2).Info("Looking for instanceTemplates for deletion")
	var predicate *filter.F
	instanceTemplates, err := s.instanceTemplates.List(ctx, predicate)
	if err != nil {
		log.Error(err, "looking for instanceTemplates for deletion")
		return err
	}

	var errs []error
	for _, instanceTemplate := range instanceTemplates {
		instanceName := instanceTemplate.Name

		if instanceTemplate.Properties == nil || instanceTemplate.Properties.Labels == nil {
			continue
		}
		labels := v1beta1.Labels(instanceTemplate.Properties.Labels)
		if !labels.HasOwned(s.scope.ClusterName()) {
			continue
		}

		log.V(2).Info("Deleting instanceTemplate", "selfLink", instanceTemplate.SelfLink)
		key := meta.GlobalKey(instanceName)
		if err := s.instanceTemplates.Delete(ctx, key); err != nil {
			if gcperrors.IsNotFound(err) {
				log.V(2).Info("instanceTemplate not found for deletion", "instanceTemplate", instanceTemplate.SelfLink)
			} else {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) == 0 {
		return nil
	}

	joined := errors.Join(errs...)
	log.Error(joined, "failed to delete instanceTemplates")
	return joined
}

func (s *Service) createOrGetInstanceTemplate(ctx context.Context) (*compute.InstanceTemplate, *meta.Key, error) {
	log := log.FromContext(ctx)

	baseKey, err := s.scope.BaseInstanceTemplateResourceName()
	if err != nil {
		return nil, nil, err
	}

	desired, err := s.scope.InstanceTemplateResource(ctx)
	if err != nil {
		return nil, nil, err
	}

	desiredJSON, err := json.Marshal(desired)
	if err != nil {
		return nil, nil, fmt.Errorf("marshalling instance template to json: %w", err)
	}
	encoded := append([]byte(baseKey.Name), desiredJSON...)
	hash := sha256.Sum256(encoded)
	hashHex := hex.EncodeToString(hash[:])

	// Store the full hash in the labels to detect potential collisions
	// GCP labels can only be 63 characters, but we fit with base32
	configHashValue := strings.ToLower(base32.HexEncoding.WithPadding('-').EncodeToString(hash[:]))
	desired.Properties.Labels[v1beta1.ConfigHashKey] = configHashValue

	namePrefix := baseKey.Name
	suffix := hashHex[:16]
	name := namePrefix + suffix

	instanceTemplateKey := meta.GlobalKey(name)

	selfLink := gcp.FormatKey("instanceTemplates", baseKey)
	log = log.WithValues("instanceTemplate", selfLink)

	log.V(2).Info("Looking for instanceTemplate")
	instanceTemplate, err := s.instanceTemplates.Get(ctx, instanceTemplateKey)
	if err != nil {
		if !gcperrors.IsNotFound(err) {
			log.Error(err, "Error looking for instanceTemplate")
			return nil, nil, err
		}

		log.V(2).Info("Creating instanceTemplate")
		if err := s.instanceTemplates.Insert(ctx, instanceTemplateKey, desired); err != nil {
			log.Error(err, "creating instanceTemplate")
			return nil, nil, err
		}

		instanceTemplate, err = s.instanceTemplates.Get(ctx, instanceTemplateKey)
		if err != nil {
			return nil, nil, err
		}
	}

	// Verify that the full hash is correct, in case of collisions on the short hash
	if v := getConfigHashKey(instanceTemplate); v != "" && v != configHashValue {
		return nil, nil, fmt.Errorf("instance template name collision detected for name %q", name)
	}

	return instanceTemplate, instanceTemplateKey, nil
}

func getConfigHashKey(it *compute.InstanceTemplate) string {
	if it.Properties != nil && it.Properties.Labels != nil {
		return it.Properties.Labels[v1beta1.ConfigHashKey]
	}
	return ""
}
