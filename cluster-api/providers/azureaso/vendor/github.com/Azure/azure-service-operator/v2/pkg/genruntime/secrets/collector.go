/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package secrets

import (
	"sort"

	"github.com/pkg/errors"
	"golang.org/x/exp/maps"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kerrors "k8s.io/apimachinery/pkg/util/errors"

	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

// TODO: This is similar to configmaps.Collector. If this is updated, that should likely be as well

// Collector collects secret values and their associated genruntime.SecretDestination's
// and produces a merged set of v1.Secret's that can be written.
type Collector struct {
	secrets   map[string]*v1.Secret
	namespace string
	errors    []error
}

// NewCollector creates a new Collector
func NewCollector(namespace string) *Collector {
	return &Collector{
		secrets:   make(map[string]*v1.Secret),
		namespace: namespace,
	}
}

func (c *Collector) get(dest *genruntime.SecretDestination) *v1.Secret {
	existing, ok := c.secrets[dest.Name]
	if !ok {
		existing = &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      dest.Name,
				Namespace: c.namespace,
			},
			StringData: make(map[string]string),
			Data:       make(map[string][]byte),
		}
		c.secrets[dest.Name] = existing
	}
	return existing
}

func (c *Collector) errIfKeyExists(val *v1.Secret, key string) error {
	if _, ok := val.StringData[key]; ok {
		return errors.Errorf("key collision, entry exists for key %s in StringData", key)
	}

	if _, ok := val.Data[key]; ok {
		return errors.Errorf("key collision, entry exists for key %s in Data", key)
	}

	return nil
}

// AddValue adds the dest and secretValue pair to the collector. If another value has already
// been added going to the same secret (but with a different key) the new key is merged into the
// existing secret.
func (c *Collector) AddValue(dest *genruntime.SecretDestination, value string) {
	if dest == nil || value == "" {
		return
	}

	existing := c.get(dest)
	err := c.errIfKeyExists(existing, dest.Key)
	if err != nil {
		c.errors = append(c.errors, err)
		return
	}

	existing.StringData[dest.Key] = value
}

// AddBinaryValue adds the dest and secretValue pair to the collector. If another value has already
// been added going to the same secret (but with a different key) the new key is merged into the
// existing secret.
func (c *Collector) AddBinaryValue(dest *genruntime.SecretDestination, value []byte) {
	if dest == nil || value == nil {
		return
	}

	existing := c.get(dest)
	err := c.errIfKeyExists(existing, dest.Key)
	if err != nil {
		c.errors = append(c.errors, err)
		return
	}

	existing.Data[dest.Key] = value
}

// Values returns the set of secrets that have been collected.
func (c *Collector) Values() ([]*v1.Secret, error) {
	err := kerrors.NewAggregate(c.errors)
	if err != nil {
		return nil, err
	}

	result := maps.Values(c.secrets)

	// Force a deterministic ordering
	sort.Slice(result, func(i, j int) bool {
		left := result[i]
		right := result[j]

		return left.Namespace < right.Namespace || (left.Namespace == right.Namespace && left.Name < right.Name)
	})

	return result, nil
}
