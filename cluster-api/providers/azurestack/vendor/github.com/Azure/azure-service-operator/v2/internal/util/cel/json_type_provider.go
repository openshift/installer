/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package cel

import (
	"encoding/json"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/pkg/errors"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

const (
	typeName = "v1.JSON"
)

var jsonMapType = types.NewMapType(types.StringType, types.NewObjectType(typeName))

type jsonProvider struct {
	baseAdapter  types.Adapter
	baseProvider types.Provider
}

// newJSONProvider creates a JSON type provider that defers all type lookups to the inner provider EXCEPT for lookups
// involving map[string]v1.JSON, which is the type used by ASO when a property doesn't have a concrete type.
// For map[string]v1.JSON, this provider transforms it into map[string]any.
// This means that (for example) a structure in Go looking like
//
//	type MyResource struct {
//		Spec              MySpec   `json:"spec,omitempty"`
//	}
//
//	type MySpec struct {
//		Untyped  map[string]v1.JSON                 `json:"untyped,omitempty"`
//	}
//
// When accessed in CEL with an expression like self.spec.untyped.nestedStruct.value, the "nestedStruct.value" part
// of that expression will reach into the v1.JSON raw payload and return the value there (if it exists). If it's not found
// a missingKey error will be returned.
// You can think of this provider as "pretend map[string]v1.JSON is actually map[string]any, and that the v1.JSON type
// doesn't exist".
func newJSONProvider() cel.EnvOption {
	return func(env *cel.Env) (*cel.Env, error) {
		provider := &jsonProvider{
			baseAdapter:  env.CELTypeAdapter(),
			baseProvider: env.CELTypeProvider(),
		}

		env, err := cel.CustomTypeAdapter(provider)(env)
		if err != nil {
			return nil, err
		}
		return cel.CustomTypeProvider(provider)(env)
	}
}

var (
	_ types.Adapter  = &jsonProvider{}
	_ types.Provider = &jsonProvider{}
)

// types.Adapter impl
func (j *jsonProvider) NativeToValue(value any) ref.Val {
	if value == nil {
		return types.NullValue
	}

	// TODO: NO idea
	switch v := value.(type) {
	case *v1.JSON:
		return nativeToValue(j, v)
	case v1.JSON:
		return nativeToValue(j, &v)
	default:
		return j.baseAdapter.NativeToValue(value)
	}
}

func nativeToValue(a types.Adapter, v *v1.JSON) ref.Val {
	if v == nil {
		return types.NullValue
	}
	var data map[string]any
	err := json.Unmarshal(v.Raw, &data)
	if err != nil {
		return types.NewErr("failed to unmarshal JSON: %s", err)
	}

	return types.NewStringInterfaceMap(a, data)
}

// types.Provider impl

func (j *jsonProvider) EnumValue(enumName string) ref.Val {
	return j.baseProvider.EnumValue(enumName)
}

func (j *jsonProvider) FindIdent(identName string) (ref.Val, bool) {
	if identName == typeName {
		return types.MapType, true
	}
	return j.baseProvider.FindIdent(identName)
}

func (j *jsonProvider) FindStructType(structType string) (*types.Type, bool) {
	if structType == typeName {
		return types.NewTypeTypeWithParam(types.MapType), true // OK so we're a map?
	}
	return j.baseProvider.FindStructType(structType)
}

func (j *jsonProvider) FindStructFieldNames(structType string) ([]string, bool) {
	if structType == typeName {
		return nil, true
	}
	return j.baseProvider.FindStructFieldNames(typeName)
}

func (j *jsonProvider) FindStructFieldType(structType, fieldName string) (*types.FieldType, bool) {
	ft, found := j.baseProvider.FindStructFieldType(structType, fieldName)
	if !found {
		return ft, found
	}

	if ft.Type.IsExactType(jsonMapType) {
		return &types.FieldType{
			Type: types.NewMapType(types.StringType, types.DynType),
			IsSet: func(target any) bool {
				return ft.IsSet(target) // For now we just proxy this, but that may not be correct
			},
			GetFrom: func(target any) (any, error) {
				inner, err := ft.GetFrom(target)
				if err != nil {
					return inner, err
				}

				typed, ok := inner.(map[string]v1.JSON)
				if !ok {
					return nil, errors.Errorf("unexpected actual type for map[string]v1.JSON")
				}

				var raw []byte
				raw, err = json.Marshal(typed)
				if err != nil {
					return nil, err
				}

				var result map[string]any
				err = json.Unmarshal(raw, &result)
				if err != nil {
					return nil, err
				}

				return result, nil
			},
		}, true
	}

	return ft, found
}

func (j *jsonProvider) NewValue(structType string, fields map[string]ref.Val) ref.Val {
	// Just proxy to baseProvider NewValue here because we don't support creating a new value of type
	// v1.JSON - as far as CEL is concerned no such type exists.
	// TODO: If we want to support construction of map[string]dyn -> map[string]v1.JSON, we need a mapping of
	// TODO: structs that contain map[string]v1.JSON, so that we can perform special handling for those.
	// TODO: See for example the commented out "construct simple resource with map[string]v1.JSON" test
	// TODO: in cel_test.go. This would boil down to performing our own conversion for that single field
	// TODO: that just serializes the map to a JSON string.
	// TODO: This isn't currently implemented because it doesn't seem likely to ever come up for our use-case (where
	// TODO: return types must be string or map[string]string
	return j.baseProvider.NewValue(structType, fields)
}
