/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package configmaps

import (
	"sort"

	"github.com/pkg/errors"
	"golang.org/x/exp/maps"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kerrors "k8s.io/apimachinery/pkg/util/errors"

	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

// TODO: This is similar to secrets.Collector. If this is updated, that should likely be as well

// Collector collects configmap values and their associated genruntime.ConfigMapDestination's
// and produces a merged set of v1.ConfigMap's that can be written.
type Collector struct {
	configs   map[string]*v1.ConfigMap
	namespace string
	errors    []error
}

// NewCollector creates a new Collector for collecting multiple config map writes together
func NewCollector(namespace string) *Collector {
	return &Collector{
		configs:   make(map[string]*v1.ConfigMap),
		namespace: namespace,
	}
}

func (c *Collector) get(dest *genruntime.ConfigMapDestination) *v1.ConfigMap {
	existing, ok := c.configs[dest.Name]
	if !ok {
		existing = &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      dest.Name,
				Namespace: c.namespace,
			},
			Data:       make(map[string]string),
			BinaryData: make(map[string][]byte),
		}
		c.configs[dest.Name] = existing
	}
	return existing
}

func (c *Collector) errIfKeyExists(val *v1.ConfigMap, key string) error {
	if _, ok := val.Data[key]; ok {
		return errors.Errorf("key collision, entry exists for key %s in Data", key)
	}

	if _, ok := val.BinaryData[key]; ok {
		return errors.Errorf("key collision, entry exists for key %s in BinaryData", key)
	}

	return nil
}

// AddValue adds the dest and ConfigMapDestination pair to the collector. If another value has already
// been added going to the same config map (but with a different key) the new key is merged into the
// existing map.
func (c *Collector) AddValue(dest *genruntime.ConfigMapDestination, value string) {
	if dest == nil || value == "" {
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

// AddBinaryValue adds the dest and ConfigMapDestination pair to the collector. If another value has already
// been added going to the same config map (but with a different key) the new key is merged into the
// existing map.
func (c *Collector) AddBinaryValue(dest *genruntime.ConfigMapDestination, value []byte) {
	if dest == nil || value == nil {
		return
	}

	existing := c.get(dest)
	err := c.errIfKeyExists(existing, dest.Key)
	if err != nil {
		c.errors = append(c.errors, err)
		return
	}

	existing.BinaryData[dest.Key] = value
}

// Values returns the set of Values that have been collected.
func (c *Collector) Values() ([]*v1.ConfigMap, error) {
	err := kerrors.NewAggregate(c.errors)
	if err != nil {
		return nil, err
	}

	result := maps.Values(c.configs)

	// Force a deterministic ordering
	sort.Slice(result, func(i, j int) bool {
		left := result[i]
		right := result[j]

		return left.Namespace < right.Namespace || (left.Namespace == right.Namespace && left.Name < right.Name)
	})

	return result, nil
}
