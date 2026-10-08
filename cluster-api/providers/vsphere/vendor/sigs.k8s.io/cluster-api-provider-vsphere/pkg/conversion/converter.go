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
	"fmt"
	"reflect"
	"strings"

	pkgerrors "github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	conversionmeta "sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion/api/meta"
)

// ConvertFunc defines a func that perform conversions.
type ConvertFunc func(context.Context, runtime.Object, runtime.Object) error

// Converter defines methods for converting API objects.
type Converter struct {
	// gvkToType allows to figure out the go type for an object with the given gkv.
	gvkToType map[schema.GroupVersionKind]reflect.Type

	// typeToGVK allows to find the gkv for a given go object.
	typeToGVK map[reflect.Type]schema.GroupVersionKind

	// gvkHubTypes allows to figure out if an object is a hub type or not.
	gvkHubTypes map[schema.GroupVersionKind]bool

	// conversionFuncs stores func to convert objects with a given gvk to another.
	conversionFuncs map[schema.GroupVersionKind]map[schema.GroupVersionKind]ConvertFunc

	// targetVersionSelector stores func that selects the target version for conversions.
	targetVersionSelector func(gk schema.GroupKind) (string, error)
}

// NewConverter returns a Converter.
func NewConverter(targetVersionSelector func(_ schema.GroupKind) (string, error)) *Converter {
	s := &Converter{
		gvkToType:             map[schema.GroupVersionKind]reflect.Type{},
		typeToGVK:             map[reflect.Type]schema.GroupVersionKind{},
		gvkHubTypes:           map[schema.GroupVersionKind]bool{},
		conversionFuncs:       map[schema.GroupVersionKind]map[schema.GroupVersionKind]ConvertFunc{},
		targetVersionSelector: targetVersionSelector,
	}
	return s
}

// AddHubTypes adds to the converter hub types that require conversion.
func (s *Converter) AddHubTypes(gv schema.GroupVersion, types ...runtime.Object) error {
	if gv.Group == "" {
		return pkgerrors.Errorf("invalid group, group cannot be empty")
	}

	if gv.Version == "" {
		return pkgerrors.Errorf("invalid version, version cannot be empty")
	}

	for _, obj := range types {
		t, err := objType(obj)
		if err != nil {
			return err
		}

		_, isConvertible := obj.(conversionmeta.Convertible)
		if !strings.HasSuffix(t.Name(), "List") && !isConvertible {
			return pkgerrors.Errorf("all objects must implement sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion/api/meta.Convertible, %v.%v does not", t.PkgPath(), t.Name())
		}

		gvk := gv.WithKind(t.Name())
		if oldT, found := s.gvkToType[gvk]; found && oldT != t {
			return pkgerrors.Errorf("double registration of different types for %v: old=%v.%v, new=%v.%v", gvk, oldT.PkgPath(), oldT.Name(), t.PkgPath(), t.Name())
		}

		if oldGvk, found := s.typeToGVK[t]; found && oldGvk != gvk {
			return pkgerrors.Errorf("double registration of different gvk for %v.%v: old=%s, new=%s", t.PkgPath(), t.Name(), oldGvk, gvk)
		}

		s.gvkToType[gvk] = t
		s.gvkHubTypes[gvk] = true
		s.typeToGVK[t] = gvk
	}
	return nil
}

// AllKnownHubTypes returns the all known hub types.
func (s *Converter) AllKnownHubTypes() map[schema.GroupVersionKind]reflect.Type {
	r := map[schema.GroupVersionKind]reflect.Type{}
	for gvk, isHub := range s.gvkHubTypes {
		if isHub {
			r[gvk] = s.gvkToType[gvk]
		}
	}
	return r
}

// Recognizes returns true if the converter is able to handle the provided group,version,kind
// of an object.
func (s *Converter) Recognizes(gvk schema.GroupVersionKind) bool {
	_, exists := s.gvkToType[gvk]
	return exists
}

// AddConversion adds to the Converter functions to be used when converting objects from one version to another.
// For instance, adding conversion from vmoprhub.VirtualMachine to vmoprv1alpha2.VirtualMachine will look like
//
// converter.AddConversion(
//
//	&vmoprvhub.VirtualMachine{},
//	vmoprv1alpha2.GroupVersion.Version, &vmoprv1alpha2.VirtualMachine{},
//	convert_hub_VirtualMachine_To_v1alpha2_VirtualMachine, convert_v1alpha2_VirtualMachine_To_hub_VirtualMachine,
//
// )
//
// AddConversionBuilder provides a convenient and typed way to perform this operation. More examples can be found in pkg/conversion/api/vmoperator.
func (s *Converter) AddConversion(hub runtime.Object, spokeVersion string, spoke runtime.Object, hubToSpokeFunc, spokeToHubFunc ConvertFunc) error {
	hubGVK, err := s.GroupVersionKindFor(hub)
	if err != nil {
		return err
	}

	if strings.HasSuffix(hubGVK.Kind, "List") {
		return pkgerrors.New("invalid source type, source type for a conversion cannot have the List suffix")
	}

	spokeType, err := objType(spoke)
	if err != nil {
		return err
	}

	spokeGVK := schema.GroupVersionKind{
		Group:   hubGVK.Group,
		Version: spokeVersion,
		Kind:    spokeType.Name(),
	}
	if spokeGVK.Group == "" {
		return pkgerrors.Errorf("invalid group, group cannot be empty")
	}
	if spokeGVK.Version == "" {
		return pkgerrors.New("invalid version, version cannot be empty")
	}

	if spokeGVK.Version == hubGVK.Version {
		return pkgerrors.New("invalid version, spokeVersion for a conversion cannot be the equal to the version registered for the hub object")
	}

	if hubGVK.Kind != spokeGVK.Kind {
		return pkgerrors.New("invalid spoke type, spoke type for a conversion must be of the same kind as the hub object")
	}

	if oldT, found := s.gvkToType[spokeGVK]; found && oldT != spokeType {
		return pkgerrors.Errorf("double registration of different types for %v: old=%v.%v, new=%v.%v", spokeGVK, oldT.PkgPath(), oldT.Name(), spokeType.PkgPath(), spokeType.Name())
	}

	if oldGVK, found := s.typeToGVK[spokeType]; found && oldGVK != spokeGVK {
		return pkgerrors.Errorf("double registration of different gvk for %v.%v: old=%s, new=%s", spokeType.PkgPath(), spokeType.Name(), oldGVK, spokeGVK)
	}

	s.gvkToType[spokeGVK] = spokeType
	s.gvkHubTypes[spokeGVK] = false
	s.typeToGVK[spokeType] = spokeGVK

	if s.conversionFuncs[hubGVK] == nil {
		s.conversionFuncs[hubGVK] = map[schema.GroupVersionKind]ConvertFunc{}
	}

	if oldC, found := s.conversionFuncs[hubGVK][spokeGVK]; found && reflect.ValueOf(oldC) != reflect.ValueOf(hubToSpokeFunc) {
		return pkgerrors.Errorf("double registration of conversion function from %v to %v: old function is different from the new function", hubGVK, spokeGVK.Version)
	}
	s.conversionFuncs[hubGVK][spokeGVK] = hubToSpokeFunc

	if s.conversionFuncs[spokeGVK] == nil {
		s.conversionFuncs[spokeGVK] = map[schema.GroupVersionKind]ConvertFunc{}
	}
	if oldC, found := s.conversionFuncs[spokeGVK][hubGVK]; found && reflect.ValueOf(oldC) != reflect.ValueOf(spokeToHubFunc) {
		return pkgerrors.Errorf("double registration of conversion function from %v to %v: old function is different from the new function", spokeGVK, hubGVK.Version)
	}
	s.conversionFuncs[spokeGVK][hubGVK] = spokeToHubFunc

	return nil
}

// Convert converts an object into another with the same kind, but a different version.
// Note:
// - If src is a hub type, dst must be one of the corresponding spoke types.
// - If src is a spoke type, dst must be the corresponding hub type.
func (s *Converter) Convert(ctx context.Context, src runtime.Object, dst runtime.Object) error {
	srcGVK, err := s.GroupVersionKindFor(src)
	if err != nil {
		return err
	}

	dstGVK, err := s.GroupVersionKindFor(dst)
	if err != nil {
		return err
	}

	if s.gvkHubTypes[srcGVK] {
		hub, _ := src.(conversionmeta.Convertible)
		source := hub.GetSource()

		if source.APIVersion != "" && source.APIVersion != dstGVK.GroupVersion().String() {
			return pkgerrors.Errorf("objects with kind %s and source.APIVersion %s cannot be converted to %s", srcGVK.Kind, source.APIVersion, dstGVK.Version)
		}
	}

	conversionFunc, ok := s.conversionFuncs[srcGVK][dstGVK]
	if !ok {
		return pkgerrors.Errorf("no conversion registered from %s to %s", srcGVK, dstGVK)
	}

	if err := conversionFunc(ctx, src, dst); err != nil {
		return pkgerrors.Wrapf(err, "error converting from %s to %s", srcGVK, dstGVK)
	}

	if s.gvkHubTypes[dstGVK] {
		hub, _ := dst.(conversionmeta.Convertible)
		hub.SetSource(conversionmeta.SourceTypeMeta{APIVersion: srcGVK.GroupVersion().String()})
	}
	return nil
}

// IsHub return true if an object is a hub type that requires conversion before write and after read.
func (s *Converter) IsHub(obj runtime.Object) bool {
	gvk, err := s.GroupVersionKindFor(obj)
	if err != nil {
		return false
	}
	return s.gvkHubTypes[gvk]
}

// SpokeGroupVersionKindFor returns the GroupVersionKind of the target spoke object for the given object.
// Note: This func should only be called for hub types.
func (s *Converter) SpokeGroupVersionKindFor(obj runtime.Object) (schema.GroupVersionKind, error) {
	gvk, err := s.GroupVersionKindFor(obj)
	if err != nil {
		return schema.GroupVersionKind{}, err
	}

	isList := false
	if strings.HasSuffix(gvk.Kind, "List") {
		isList = true
		gvk.Kind = strings.TrimSuffix(gvk.Kind, "List")
	}

	if !s.gvkHubTypes[gvk] {
		return schema.GroupVersionKind{}, pkgerrors.Errorf("no type registered for %s", gvk)
	}

	spokeVersion, err := s.targetVersionSelector(gvk.GroupKind())
	if err != nil {
		return schema.GroupVersionKind{}, pkgerrors.Wrapf(err, "no target version registered for %s", gvk.GroupKind())
	}

	spokeGVK := schema.GroupVersionKind{
		Group:   gvk.Group,
		Version: spokeVersion,
		Kind:    gvk.Kind,
	}

	if _, ok := s.conversionFuncs[gvk][spokeGVK]; !ok {
		return schema.GroupVersionKind{}, pkgerrors.Errorf("no conversion registered from %s to %s", gvk, spokeGVK.Version)
	}

	if isList {
		spokeGVK.Kind = fmt.Sprintf("%sList", spokeGVK.Kind)
	}

	return spokeGVK, nil
}

// GroupVersionKindFor returns the GroupVersionKind for the given object.
// Note: obj can be either a hub type or a spoke type.
func (s *Converter) GroupVersionKindFor(obj runtime.Object) (schema.GroupVersionKind, error) {
	t, err := objType(obj)
	if err != nil {
		return schema.GroupVersionKind{}, err
	}

	gvk, ok := s.typeToGVK[t]
	if !ok {
		return schema.GroupVersionKind{}, pkgerrors.Errorf("no type registered for %s.%s", t.PkgPath(), t.Name())
	}
	return gvk, nil
}

func objType(obj runtime.Object) (reflect.Type, error) {
	if obj == nil {
		return nil, pkgerrors.New("all objects must be pointers to structs, got nil")
	}

	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Pointer {
		return nil, pkgerrors.Errorf("all objects must be pointers to structs, got %s", t.Kind())
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return nil, pkgerrors.Errorf("all objects must be pointers to structs, got *%s", t.Kind())
	}
	return t, nil
}
