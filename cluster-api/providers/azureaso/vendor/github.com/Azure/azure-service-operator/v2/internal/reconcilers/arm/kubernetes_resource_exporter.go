/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package arm

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/Azure/azure-service-operator/v2/internal/set"
	asocel "github.com/Azure/azure-service-operator/v2/internal/util/cel"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/conditions"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/configmaps"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/core"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/secrets"
)

type kubernetesResourceExporter interface {
	Export(ctx context.Context) ([]client.Object, error)
}

var _ kubernetesResourceExporter = &configMapExpressionExporter{}

type configMapExpressionExporter struct {
	obj                 genruntime.ARMMetaObject
	versionedObj        genruntime.ARMMetaObject
	expressionEvaluator asocel.ExpressionEvaluator
}

func (c *configMapExpressionExporter) Export(ctx context.Context) ([]client.Object, error) {
	cmExporter, ok := c.obj.(configmaps.Exporter)
	if !ok {
		return nil, nil
	}

	cmExpressions := cmExporter.ConfigMapDestinationExpressions()
	if len(cmExpressions) == 0 {
		return nil, nil
	}

	resources, err := c.parseConfigMaps(cmExpressions)
	if err != nil {
		err = errors.Wrap(err, "failed to parse configmap expressions for export")
		return nil, conditions.NewReadyConditionImpactingError(err, conditions.ConditionSeverityError, conditions.ReasonAdditionalKubernetesObjWriteFailure)
	}
	return resources, err
}

func (c *configMapExpressionExporter) parseConfigMaps(expressions []*core.DestinationExpression) ([]client.Object, error) {
	collector := configmaps.NewCollector(c.obj.GetNamespace())

	for _, expression := range expressions {
		value, err := c.expressionEvaluator.CompileAndRun(expression.Value, c.versionedObj, nil)
		if err != nil {
			return nil, err
		}

		if value.Value != "" {
			collector.AddValue(
				&genruntime.ConfigMapDestination{
					Name: expression.Name,
					Key:  expression.Key,
				}, value.Value)
		} else if len(value.Values) > 0 {
			for k, v := range value.Values {
				collector.AddValue(
					&genruntime.ConfigMapDestination{
						Name: expression.Name,
						Key:  k,
					}, v)
			}
		} else {
			return nil, errors.Errorf("unexpected expression output")
		}
	}

	result, err := collector.Values()
	if err != nil {
		return nil, err
	}
	return configmaps.SliceToClientObjectSlice(result), nil
}

var _ kubernetesResourceExporter = &secretExpressionExporter{}

type secretExpressionExporter struct {
	obj                 genruntime.ARMMetaObject
	versionedObj        genruntime.ARMMetaObject
	rawSecrets          map[string]string
	expressionEvaluator asocel.ExpressionEvaluator
}

func (s *secretExpressionExporter) Export(ctx context.Context) ([]client.Object, error) {
	secretExporter, ok := s.obj.(secrets.Exporter)
	if !ok {
		return nil, nil
	}

	secretExpressions := secretExporter.SecretDestinationExpressions()
	if len(secretExpressions) == 0 {
		return nil, nil
	}
	resources, err := s.parseSecrets(secretExpressions, s.versionedObj, s.rawSecrets)
	if err != nil {
		err = errors.Wrap(err, "failed to parse secret expressions for export")
		return nil, conditions.NewReadyConditionImpactingError(err, conditions.ConditionSeverityError, conditions.ReasonAdditionalKubernetesObjWriteFailure)
	}

	return resources, nil
}

func (s *secretExpressionExporter) parseSecrets(
	expressions []*core.DestinationExpression,
	versionedObj genruntime.ARMMetaObject,
	rawSecrets map[string]string,
) ([]client.Object, error) {
	collector := secrets.NewCollector(versionedObj.GetNamespace())

	for _, expression := range expressions {
		value, err := s.expressionEvaluator.CompileAndRun(expression.Value, versionedObj, rawSecrets)
		if err != nil {
			return nil, err
		}

		if value.Value != "" {
			collector.AddValue(
				&genruntime.SecretDestination{
					Name: expression.Name,
					Key:  expression.Key,
				}, value.Value)
		} else if len(value.Values) > 0 {
			for k, v := range value.Values {
				collector.AddValue(
					&genruntime.SecretDestination{
						Name: expression.Name,
						Key:  k,
					}, v)
			}
		} else {
			return nil, errors.Errorf("unexpected expression output")
		}
	}

	result, err := collector.Values()
	if err != nil {
		return nil, err
	}
	return secrets.SliceToClientObjectSlice(result), nil
}

func findRequiredSecrets(
	expressionEvaluator asocel.ExpressionEvaluator,
	obj genruntime.ARMMetaObject,
	versionedObj genruntime.ARMMetaObject,
) (set.Set[string], error) {
	secretExporter, ok := obj.(secrets.Exporter)
	if !ok {
		return nil, nil
	}
	expressions := secretExporter.SecretDestinationExpressions()
	if len(expressions) == 0 {
		return nil, nil
	}

	result := make(set.Set[string])
	for _, expression := range expressions {
		secretsNeeded, err := expressionEvaluator.FindSecretUsage(expression.Value, versionedObj)
		if err != nil {
			return nil, err
		}

		result.AddAll(secretsNeeded)
	}

	return result, nil
}

var _ kubernetesResourceExporter = &kubernetesSecretExporter{}

type kubernetesSecretExporter struct {
	obj               genruntime.ARMMetaObject
	log               logr.Logger
	extension         genruntime.ResourceExtension
	additionalSecrets set.Set[string]
	connection        Connection

	// This is a bit hacky but we stash the RawSecrets here for use by other exporters later
	rawSecrets map[string]string
}

func (e *kubernetesSecretExporter) Export(ctx context.Context) ([]client.Object, error) {
	retriever := extensions.CreateKubernetesSecretExporter(ctx, e.extension, e.connection.Client(), e.log)
	result, err := retriever(e.obj, e.additionalSecrets)
	if err != nil {
		return nil, errors.Wrap(err, "extension failed to produce resources for export")
	}
	if result == nil {
		return nil, nil
	}

	e.rawSecrets = result.RawSecrets

	return result.Objs, nil
}

var _ kubernetesResourceExporter = &autoGeneratedConfigExporter{}

type autoGeneratedConfigExporter struct {
	versionedObj genruntime.ARMMetaObject
	log          logr.Logger
	connection   Connection
}

func (a *autoGeneratedConfigExporter) Export(ctx context.Context) ([]client.Object, error) {
	exporter, ok := a.versionedObj.(genruntime.KubernetesConfigExporter)
	if !ok {
		return nil, nil
	}

	resources, err := exporter.ExportKubernetesConfigMaps(ctx, a.versionedObj, a.connection.Client(), a.log)
	if err != nil {
		return nil, errors.Wrap(err, "failed to produce resources for export")
	}

	return resources, nil
}

var _ kubernetesResourceExporter = &manualConfigExporter{}

type manualConfigExporter struct {
	obj        genruntime.ARMMetaObject
	extension  genruntime.ResourceExtension
	log        logr.Logger
	connection Connection
}

func (a *manualConfigExporter) Export(ctx context.Context) ([]client.Object, error) {
	exporter, ok := a.extension.(genruntime.KubernetesConfigExporter)
	if !ok {
		return nil, nil
	}

	resources, err := exporter.ExportKubernetesConfigMaps(ctx, a.obj, a.connection.Client(), a.log)
	if err != nil {
		return nil, errors.Wrap(err, "failed to produce resources for export")
	}

	return resources, nil
}
