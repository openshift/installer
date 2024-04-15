/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package reflecthelpers

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/Azure/azure-service-operator/v2/internal/set"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

// ValueOfPtr dereferences a pointer and returns the value the pointer points to.
// Use this as carefully as you would the * operator
// TODO: Can we delete this helper later when we have some better code generated functions?
func ValueOfPtr(ptr interface{}) interface{} {
	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("Can't get value of pointer for non-pointer type %T", ptr))
	}
	val := reflect.Indirect(v)

	return val.Interface()
}

// DeepCopyInto calls in.DeepCopyInto(out)
func DeepCopyInto(in client.Object, out client.Object) {
	inVal := reflect.ValueOf(in)

	method := inVal.MethodByName("DeepCopyInto")
	method.Call([]reflect.Value{reflect.ValueOf(out)})
}

// FindReferences finds references of the given type on the provided object
func FindReferences(obj interface{}, t reflect.Type) (map[interface{}]struct{}, error) {
	result := make(map[interface{}]struct{})

	visitor := NewReflectVisitor()
	visitor.VisitStruct = func(this *ReflectVisitor, it reflect.Value, ctx interface{}) error {
		if it.Type() == t {
			if it.CanInterface() {
				result[it.Interface()] = struct{}{}
			}
			return nil
		}

		return IdentityVisitStruct(this, it, ctx)
	}

	err := visitor.Visit(obj, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "scanning for references of type %s", t.String())
	}

	return result, nil
}

// FindPropertiesWithTag finds all the properties with the given tag on the specified object and
// returns a map of the property name to the property value
func FindPropertiesWithTag(obj interface{}, tag string) (map[string][]interface{}, error) {
	result := make(map[string][]interface{})

	visitor := NewReflectVisitor()
	visitor.VisitStruct = func(this *ReflectVisitor, it reflect.Value, ctx interface{}) error {
		// This was adapted from IdentityVisitStruct
		for i := 0; i < it.NumField(); i++ {
			fieldVal := it.Field(i)
			if !fieldVal.CanInterface() {
				// Bypass unexported fields
				continue
			}

			structField := it.Type().Field(i)
			path := ctx.(string)
			if path == "" {
				path = structField.Name
			} else {
				path += "." + structField.Name
			}
			_, ok := structField.Tag.Lookup(tag)
			field := it.Field(i)
			if ok && field.CanInterface() {
				if len(result[path]) == 0 {
					result[path] = []interface{}{}
				}
				result[path] = append(result[path], field.Interface())
			}

			err := this.visit(field, path)
			if err != nil {
				return err
			}
		}

		return nil
	}

	err := visitor.Visit(obj, "")
	if err != nil {
		return nil, errors.Wrapf(err, "scanning for references to tag %s", tag)
	}

	return result, nil
}

// FindResourceReferences finds all the genruntime.ResourceReference's on the provided object
func FindResourceReferences(obj interface{}) (set.Set[genruntime.ResourceReference], error) {
	return Find[genruntime.ResourceReference](obj)
}

// FindSecretReferences finds all the genruntime.SecretReference's on the provided object
func FindSecretReferences(obj interface{}) (set.Set[genruntime.SecretReference], error) {
	return Find[genruntime.SecretReference](obj)
}

// FindConfigMapReferences finds all the genruntime.ConfigMapReference's on the provided object
func FindConfigMapReferences(obj interface{}) (set.Set[genruntime.ConfigMapReference], error) {
	return Find[genruntime.ConfigMapReference](obj)
}

// Find finds all the references of the given type on the provided object
func Find[T comparable](obj interface{}) (set.Set[T], error) {
	var t T
	untypedResult, err := FindReferences(obj, reflect.TypeOf(t))
	if err != nil {
		return nil, err
	}

	result := set.Make[T]()
	for k := range untypedResult {
		result.Add(k.(T))
	}

	return result, nil
}

// FindOptionalConfigMapReferences finds all the genruntime.ConfigMapReference's on the provided object
func FindOptionalConfigMapReferences(obj interface{}) ([]*genruntime.OptionalConfigMapReferencePair, error) {
	untypedResult, err := FindPropertiesWithTag(obj, "optionalConfigMapPair") // TODO: This is astmodel.OptionalConfigMapPairTag
	if err != nil {
		return nil, err
	}

	collector := make(map[string][]*genruntime.OptionalConfigMapReferencePair)
	suffix := "FromConfig" // TODO This is astmodel.OptionalConfigMapReferenceSuffix

	// This could probably be more efficient, but this avoids code duplication, and we're not dealing
	// with huge collections here.
	for key, values := range untypedResult {
		if strings.HasSuffix(key, suffix) {
			continue
		}

		collector[key] = make([]*genruntime.OptionalConfigMapReferencePair, 0, len(values))
		for _, val := range values {
			typedValue, ok := val.(*string)
			if !ok {
				return nil, errors.Errorf("value of property %s was not a *string like expected", key)
			}
			collector[key] = append(collector[key], &genruntime.OptionalConfigMapReferencePair{
				Name:  key,
				Value: typedValue,
			})
		}
	}

	for key, values := range untypedResult {
		if !strings.HasSuffix(key, suffix) {
			continue
		}
		idx := strings.TrimSuffix(key, suffix)
		if len(values) != len(collector[idx]) {
			return nil, errors.Errorf("number of Ref's didn't match number of Values for %s", idx)
		}

		for i, val := range values {
			typedValue, ok := val.(*genruntime.ConfigMapReference)
			if !ok {
				return nil, errors.Errorf("value of property %s was not a genruntime.ConfigMapReference like expected", key)
			}
			collector[idx][i].RefName = key
			collector[idx][i].Ref = typedValue
		}
	}

	// Translate our collector into a simple list
	var result []*genruntime.OptionalConfigMapReferencePair
	for _, values := range collector {
		for _, val := range values {
			result = append(result, val)
		}
	}

	return result, nil
}

// GetObjectListItems gets the list of items from an ObjectList
func GetObjectListItems(listPtr client.ObjectList) ([]client.Object, error) {
	itemsField, err := getItemsField(listPtr)
	if err != nil {
		return nil, err
	}

	var result []client.Object
	for i := 0; i < itemsField.Len(); i++ {
		item := itemsField.Index(i)

		if item.Kind() == reflect.Struct {
			if !item.CanAddr() {
				return nil, errors.Errorf("provided list elements were not pointers, but cannot be addressed")
			}
			item = item.Addr()
		}

		typedItem, ok := item.Interface().(client.Object)
		if !ok {
			return nil, errors.Errorf("provided list elements did not implement client.Object interface")
		}

		result = append(result, typedItem)
	}

	return result, nil
}

// SetObjectListItems gets the list of items from an ObjectList
func SetObjectListItems(listPtr client.ObjectList, items []client.Object) (returnErr error) {
	itemsField, err := getItemsField(listPtr)
	if err != nil {
		return err
	}

	if !itemsField.CanSet() {
		return errors.Errorf("cannot set items field of %T", listPtr)
	}
	defer func() {
		if recovered := recover(); recovered != nil {
			returnErr = errors.Errorf("failed to set items field of %T: %s", listPtr, recovered)
		}
	}()

	slice := reflect.MakeSlice(itemsField.Type(), 0, 0)
	for _, item := range items {
		val := reflect.ValueOf(item)

		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		slice = reflect.Append(slice, val)
	}

	itemsField.Set(slice)
	return nil
}

func getItemsField(listPtr client.ObjectList) (reflect.Value, error) {
	val := reflect.ValueOf(listPtr)
	if val.Kind() != reflect.Ptr {
		return reflect.Value{}, errors.Errorf("provided list was not a pointer, was %s", val.Kind())
	}

	list := val.Elem()

	if list.Kind() != reflect.Struct {
		return reflect.Value{}, errors.Errorf("provided list was not a struct, was %s", val.Kind())
	}

	itemsField := list.FieldByName("Items")
	if (itemsField == reflect.Value{}) {
		return reflect.Value{}, errors.Errorf("provided list has no field \"Items\"")
	}
	if itemsField.Kind() != reflect.Slice {
		return reflect.Value{}, errors.Errorf("provided list \"Items\" field was not of type slice")
	}

	return itemsField, nil
}

// SetProperty sets the property on the provided object to the provided value.
// obj is the object to modify.
// propertyPath is a dot-separated path to the property to set.
// value is the value to set the property to.
// Returns an error if any of the properties in the path do not exist, if the property is not settable,
// or if the value provided is incompatible.
func SetProperty(obj any, propertyPath string, value any) error {
	if obj == nil {
		return errors.Errorf("provided object was nil")
	}

	if propertyPath == "" {
		return errors.Errorf("property path was empty")
	}

	steps := strings.Split(propertyPath, ".")
	return setPropertyCore(obj, steps, value)
}

func setPropertyCore(obj any, propertyPath []string, value any) (err error) {
	// Catch any panic that occurs when setting the field and turn it into an error return
	defer func() {
		if recovered := recover(); recovered != nil {
			err = errors.Errorf("failed to set property %s: %s", propertyPath[0], recovered)
		}
	}()

	// Get the underlying object we need to modify
	subject := reflect.ValueOf(obj)

	// Dereference pointers
	if subject.Kind() == reflect.Ptr {
		subject = subject.Elem()
	}

	// Check we have a struct
	if subject.Kind() != reflect.Struct {
		return errors.Errorf("provided object was not a struct, was %s", subject.Kind())
	}

	// Get the field we need to modify
	field := subject.FieldByName(propertyPath[0])

	// Check the field exists
	if field == (reflect.Value{}) {
		return errors.Errorf("provided object did not have a field named %s", propertyPath[0])
	}

	// If this is not the last property in the path, we need to recurse
	if len(propertyPath) > 1 {
		if field.Kind() == reflect.Ptr {
			// Field is a pointer; initialize it if needed, then pass the pointer recursively
			if field.IsNil() {
				newValue := reflect.New(field.Type().Elem())
				field.Set(newValue)
			}

			err = setPropertyCore(field.Interface(), propertyPath[1:], value)
			if err != nil {
				return errors.Wrapf(err, "failed to set property %s",
					propertyPath[0])
			}

			return nil
		}

		// Field is not a pointer, so we need to pass the address of the field recursively
		err = setPropertyCore(field.Addr().Interface(), propertyPath[1:], value)
		if err != nil {
			return errors.Wrapf(err, "failed to set property %s",
				propertyPath[0])
		}

		return nil
	}

	// If this is the last property in the path, we need to set the value, if we can
	if !field.CanSet() {
		return errors.Errorf("field %s was not settable", propertyPath[0])
	}

	// Cast value to the type required by the field
	valueKind := reflect.ValueOf(value)
	if !valueKind.CanConvert(field.Type()) {
		return errors.Errorf("value of kind %s was not compatible with field %s", valueKind, propertyPath[0])
	}

	value = valueKind.Convert(field.Type()).Interface()

	field.Set(reflect.ValueOf(value))
	return nil
}
