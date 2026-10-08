/*
Copyright 2026 The Kubernetes Authors.

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

package conversion

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// ConverterBuilder collects functions that add things to a converter. It's to allow
// code to compile without explicitly referencing generated types.
type ConverterBuilder struct {
	funcs []func(converter *Converter) error
	gv    schema.GroupVersion
}

// AddToConverter applies all the stored functions to the converter. A non-nil error
// indicates that one function failed and the attempt was abandoned.
func (sb *ConverterBuilder) AddToConverter(s *Converter) error {
	for _, f := range sb.funcs {
		if err := f(s); err != nil {
			return err
		}
	}
	return nil
}

// AddTypes adds to the Converter types that require conversion.
func (sb *ConverterBuilder) AddTypes(types ...runtime.Object) {
	sb.funcs = append(sb.funcs, func(s *Converter) error {
		return s.AddHubTypes(sb.gv, types...)
	})
}

// AddConversion adds to the Converter functions to be used when converting objects from one version to another.
// For instance, adding conversion from vmoprhub.VirtualMachine to vmoprv1alpha2.VirtualMachine will look like
//
// converterBuilder.AddConversion(
//
//	conversion.NewAddConversionBuilder(convert_hub_VirtualMachine_To_v1alpha2_VirtualMachine, convert_v1alpha2_VirtualMachine_To_hub_VirtualMachine),
//
// )
//
// More examples can be found in pkg/conversion/api/vmoperator.
func (sb *ConverterBuilder) AddConversion(builder AddConversionBuilder) {
	sb.funcs = append(sb.funcs, builder.Build(sb.gv.Version))
}

// NewConverterBuilder returns a ConverterBuilder.
func NewConverterBuilder(gv schema.GroupVersion, fs ...func(s *Converter) error) ConverterBuilder {
	cb := ConverterBuilder{
		gv:    gv,
		funcs: fs,
	}
	return cb
}

// AddConversionBuilder build a func that adds a conversion to a Converter.
type AddConversionBuilder interface {
	// Build a func that adds a conversion to a Converter.
	Build(version string) func(converter *Converter) error
}

// NewAddConversionBuilder return a AddConversionBuilder.
func NewAddConversionBuilder[hubT, spokeT runtime.Object](
	hubToSpokeFunc func(ctx context.Context, hub hubT, spoke spokeT) error,
	spokeToHubFunc func(ctx context.Context, spoke spokeT, hub hubT) error,
) AddConversionBuilder {
	return &conversionBuilder[hubT, spokeT]{
		hubToSpokeFunc: hubToSpokeFunc,
		spokeToHubFunc: spokeToHubFunc,
	}
}

type conversionBuilder[hubT, spokeT runtime.Object] struct {
	hubToSpokeFunc func(ctx context.Context, hub hubT, spoke spokeT) error
	spokeToHubFunc func(ctx context.Context, spoke spokeT, hub hubT) error
}

// Build a func that adds a conversion to a Converter.
func (c conversionBuilder[hubObject, spokeObject]) Build(version string) func(converter *Converter) error {
	return func(converter *Converter) error {
		return converter.AddConversion(
			createZero[hubObject](),
			version, createZero[spokeObject](),
			func(ctx context.Context, hub runtime.Object, spoke runtime.Object) error {
				return c.hubToSpokeFunc(ctx, hub.(hubObject), spoke.(spokeObject))
			},
			func(ctx context.Context, spoke runtime.Object, hub runtime.Object) error {
				return c.spokeToHubFunc(ctx, spoke.(spokeObject), hub.(hubObject))
			},
		)
	}
}

func createZero[T runtime.Object]() T {
	var val T
	return val
}
