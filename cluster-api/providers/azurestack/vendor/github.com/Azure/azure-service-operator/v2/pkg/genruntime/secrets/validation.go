/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package secrets

import (
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/Azure/azure-service-operator/v2/internal/set"
	asocel "github.com/Azure/azure-service-operator/v2/internal/util/cel"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/core"
)

// For unit testing purposes
var expressionEvaluator = asocel.Evaluator

type keyPair struct {
	name string
	key  string
}

// ValidateDestinations checks that no two destinations are writing to the same secret/key, as that could cause
// those secrets to overwrite one another.
func ValidateDestinations(
	self any,
	destinations []*genruntime.SecretDestination,
	destinationExpressions []*core.DestinationExpression,
) (admission.Warnings, error) {
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

	celEvaluator := expressionEvaluator()

	for _, dest := range destinationExpressions {
		if dest == nil {
			continue
		}

		if celEvaluator != nil {
			// Validate that the expression is a valid CEL expression
			outputType, err := celEvaluator.Check(dest.Value, self)
			if err != nil {
				// Error should already have enough context
				return nil, err
			}

			if outputType.IsExactType(asocel.StringType) && dest.Key == "" {
				return nil, errors.Errorf("CEL expression with output type string must specify destination 'key', %s", dest.String())
			}
			if outputType.IsExactType(asocel.MapType) && dest.Key != "" {
				return nil, errors.Errorf("CEL expression with output type map[string]string must not specify destination 'key', %s", dest.String())
			}
		}

		if dest.Key == "" {
			continue
		}

		// Validate no destination collisions. Note that we can't check collisions for map keys because to do that
		// we'd have to actually evaluate the CEL expression which we can't do yet as it may be dependent on
		// things like .status fields which aren't yet filled out in the resource.
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
