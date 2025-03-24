/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package secrets

import (
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/Azure/azure-service-operator/v2/internal/set"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

type keyPair struct {
	name string
	key  string
}

// ValidateDestinations checks that no two destinations are writing to the same secret/key, as that could cause
// those secrets to overwrite one another.
func ValidateDestinations(destinations []*genruntime.SecretDestination) (admission.Warnings, error) {
	// Map of secret -> keys
	locations := set.Make[keyPair]()

	for _, dest := range destinations {
		if dest == nil {
			continue
		}

		pair := keyPair{
			name: dest.Name,
			key:  dest.Key,
		}
		if locations.Contains(pair) {
			return nil, errors.Errorf("cannot write more than one secret to destination %s", dest.String())
		}

		locations.Add(pair)
	}

	return nil, nil
}
