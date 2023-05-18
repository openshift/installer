// Copyright 2023 Google LLC. All Rights Reserved.
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package dcl

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// DiffInfo is a struct that contains all information about the diff that's about to occur.
type DiffInfo struct {
	// Ignore + OutputOnly cause the diff checker to always return no-diff.
	Ignore           bool
	OutputOnly       bool
	ServerDefault    bool
	MergeNestedDiffs bool
	IgnoredPrefixes  []string
	Type             string

	// ObjectFunction is the function used to diff a Nested Object.
	ObjectFunction func(desired, actual interface{}, fn FieldName) ([]*FieldDiff, error)

	// CustomDiff is used to handle diffing a field when normal diff functions will not suffice.
	// It should return false if there is any diff between 'desired' and 'actual'.
	CustomDiff func(desired, actual interface{}) bool

	// OperationSelector takes in the field's diff and returns the name of the operation (or Recreate) that should be triggered.
	OperationSelector func(d *FieldDiff) []string

	EmptyObject interface{}
}

// FieldName is used to add information about a field's name for logging purposes.
type FieldName struct {
	FieldName string
}

// AddIndex adds an index to a FieldName and returns the same item.
// Info is always pass-by-value, so the original field name still exists.
func (i FieldName) AddIndex(index int) FieldName {
	newInfo := i
	newInfo.FieldName = newInfo.FieldName + fmt.Sprintf("[%v]", index)
	return newInfo
}

// AddNest adds an index to a FieldName and returns the same item.
// Info is always pass-by-value, so the original field name still exists.
func (i FieldName) AddNest(field string) FieldName {
	newInfo := i
	if i.FieldName == "" {
		newInfo.FieldName = field
	} else {
		newInfo.FieldName = newInfo.FieldName + fmt.Sprintf(".%s", field)
	}
	return newInfo
}

// FieldDiff contains all information about a diff that exists in the resource.
type FieldDiff struct {
	FieldName string
	Message   string
	Desired   interface{}
	Actual    interface{}

	ToAdd    []interface{}
	ToRemove []interface{}

	// The name of the operation that should result (may be Recreate)
	// In the case of sets, more than one operation may be returned.
	ResultingOperation []string
}

func (d *FieldDiff) String() string {
	if d.Message != "" {
		return fmt.Sprintf("Field %s diff: %s", d.FieldName, d.Message)
	} else if len(d.ToAdd) != 0 || len(d.ToRemove) != 0 {
		return fmt.Sprintf("Field %s: add %v, remove %v", d.FieldName, d.ToAdd, d.ToRemove)
	}
	return fmt.Sprintf("Field %s: got %s, want %s", d.FieldName, SprintResourceCompact(d.Actual), SprintResourceCompact(d.Desired))
}

func stringValue(i interface{}) string {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return "nil"
		}
		return fmt.Sprintf("%v", reflect.Indirect(v))
	}
	return fmt.Sprintf("%v", i)
}

// Diff takes in two interfaces and diffs them according to Info.
func Diff(desired, actual interface{}, info DiffInfo, fn FieldName) ([]*FieldDiff, error) {
	var diffs []*FieldDiff
	// All Output-only fields should not be diffed.
	if info.OutputOnly || info.Ignore {
		return nil, nil
	}

	// If desired is a zero value, we do not care about the field.
	if IsZeroValue(desired) {
		return nil, nil
	}

	if info.OperationSelector == nil {
		return nil, fmt.Errorf("an operation selector function must exist")
	}

	desiredType := ValueType(desired)

	if desiredType == "invalid" {
		return nil, nil
	}

	if info.CustomDiff != nil {
		if !info.CustomDiff(desired, actual) {
			diffs = append(diffs, &FieldDiff{FieldName: fn.FieldName, Desired: desired, Actual: actual})
		}
		addOperationToDiffs(diffs, info)
		return diffs, nil
	}

	if desiredType == "slice" {
		dSlice, iSlice, err := slices(desired, actual)
		if err != nil {
			return nil, err
		}
		var arrDiffs []*FieldDiff
		if info.Type == "Set" {
			arrDiffs, err = setDiff(dSlice, iSlice, info, fn)
		} else {
			arrDiffs, err = arrayDiff(dSlice, iSlice, info, fn)
		}
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, arrDiffs...)
		addOperationToDiffs(diffs, info)
		return diffs, nil
	}

	if info.Type == "EnumType" {
		if !reflect.DeepEqual(desired, actual) {
			diffs = append(diffs, &FieldDiff{FieldName: fn.FieldName, Desired: desired, Actual: actual})
			addOperationToDiffs(diffs, info)
			return diffs, nil
		}
		return nil, nil
	}

	switch desiredType {
	case "string":
		dStr, err := str(desired)
		if err != nil {
			return nil, err
		}

		// Protobufs cannot differentiate between empty primitive values + null.
		// If the API returns null or does not return a value for the field and we have set the empty string, those are equivalent.
		if IsZeroValue(actual) && *dStr == "" {
			return diffs, nil
		}

		aStr, err := str(actual)
		if err != nil {
			return nil, err
		}

		if info.Type == "ReferenceType" {
			if !StringEqualsWithSelfLink(dStr, aStr) {
				diffs = append(diffs, &FieldDiff{FieldName: fn.FieldName, Desired: dStr, Actual: aStr})
			}
		} else {
			if !StringCanonicalize(dStr, aStr) {
				diffs = append(diffs, &FieldDiff{FieldName: fn.FieldName, Desired: dStr, Actual: aStr})
			}
		}

	case "map":
		dMap, aMap, err := maps(desired, actual)
		if err != nil {
			return nil, err
		}
		mapDiffs, err := mapCompare(dMap, aMap, info.IgnoredPrefixes, info, fn)
		if err != nil {
			return nil, err
		}
		if len(mapDiffs) > 0 {
			diffs = append(diffs, mapDiffs...)
		}

	case "int64":
		dInt, err := makeint64(desired)
		if err != nil {
			return nil, err
		}

		// 0 is the empty value for integers.
		if IsZeroValue(actual) && *dInt == 0 {
			return diffs, nil
		}

		if !reflect.DeepEqual(desired, actual) {
			diffs = append(diffs, &FieldDiff{FieldName: fn.FieldName, Desired: desired, Actual: actual})
		}

	case "int":
		dInt, err := makeint(desired)
		if err != nil {
			return nil, err
		}

		// 0 is the empty value for integers.
		if IsZeroValue(actual) && *dInt == 0 {
			return diffs, nil
		}

		if !reflect.DeepEqual(desired, actual) {
			diffs = append(diffs, &FieldDiff{FieldName: fn.FieldName, Desired: desired, Actual: actual})
		}

	case "float64":
		dFloat, err := makefloat64(desired)
		if err != nil {
			return nil, err
		}

		// 0 is the empty value for integers.
		if IsZeroValue(actual) && *dFloat == 0.0 {
			return diffs, nil
		}

		if !reflect.DeepEqual(desired, actual) {
			diffs = append(diffs, &FieldDiff{FieldName: fn.FieldName, Desired: desired, Actual: actual})
		}

	case "bool":
		dBool, aBool, err := bools(desired, actual)
		if err != nil {
			return nil, err
		}
		if !BoolCanonicalize(dBool, aBool) {
			diffs = append(diffs, &FieldDiff{FieldName: fn.FieldName, Desired: dBool, Actual: aBool})
		}

	case "struct":
		// If API returns nil (which means field is unset) && we have the empty-struct, no diff occurs.
		if IsZeroValue(actual) && IsEmptyValueIndirect(desired) {
			return nil, nil
		}

		// Want empty value, but non-empty value currrently exists.
		// Only consider *explicitly* empty values, rather than "some combination
		// of nils and falses" (as IEVI would do), because of the case comparing
		// a non-explicitly empty struct with a struct containing only computed fields.
		// See compute's `validate_test.go` for example.
		if hasEmptyStructField(desired) && !IsEmptyValueIndirect(actual) {
			if info.ServerDefault {
				// The API can return values where none are in the desired state.
				return nil, nil
			}
			diffs = append(diffs, &FieldDiff{FieldName: fn.FieldName, Desired: desired, Actual: actual})
			addOperationToDiffs(diffs, info)
			return diffs, nil
		}

		if info.ObjectFunction == nil {
			return nil, fmt.Errorf("struct %v given without an object function", desired)
		}

		if info.EmptyObject == nil {
			return nil, fmt.Errorf("struct %v given without an empty object type", desired)
		}

		// If the API returns nil, we can't diff against a nil. We should use the empty object instead.
		// This is because the user could write out a config that is functionally equivalent to the empty object (contains all 0s and ""),
		// but is not technically the empty object.
		if actual == nil || ValueType(actual) == "invalid" {
			actual = info.EmptyObject
		}

		ds, err := info.ObjectFunction(desired, actual, fn)
		if err != nil {
			return nil, err
		}
		if info.MergeNestedDiffs {
			// Replace any nested diffs with a recreate operation with a diff in this field.
			nonRecreateCount := 0
			for _, d := range ds {
				if len(d.ResultingOperation) == 0 {
					return nil, fmt.Errorf("diff found in field %q with no operation", d.FieldName)
				}
				if d.ResultingOperation[0] != "Recreate" {
					ds[nonRecreateCount] = d
					nonRecreateCount++
				}
			}
			if nonRecreateCount < len(ds) {
				// At least one nested diff requires a recreate.
				ds[nonRecreateCount] = &FieldDiff{FieldName: fn.FieldName, Desired: desired, Actual: actual}
				nonRecreateCount++
			}
			ds = ds[:nonRecreateCount]
		}
		diffs = append(diffs, ds...)
	default:
		return nil, fmt.Errorf("no diffing logic exists for type: %q", desiredType)
	}

	addOperationToDiffs(diffs, info)
	return diffs, nil
}

func arrayDiff(desired, actual []interface{}, info DiffInfo, fn FieldName) ([]*FieldDiff, error) {
	var diffs []*FieldDiff

	// Nothing to diff against.
	if actual == nil {
		return diffs, nil
	}

	if len(desired) != len(actual) && !IsZeroValue(desired) {
		diffs = append(diffs, &FieldDiff{FieldName: fn.FieldName, Message: fmt.Sprintf("different lengths: desired %d, actual %d", len(desired), len(actual))})
		return diffs, nil
	}

	for i, dItem := range desired {
		aItem := actual[i]
		diff, err := Diff(dItem, aItem, info, fn.AddIndex(i))
		if err != nil {
			return nil, err
		}
		if diff != nil {
			diffs = append(diffs, diff...)
		}
	}
	return diffs, nil
}

func setDiff(desired, actual []interface{}, info DiffInfo, fn FieldName) ([]*FieldDiff, error) {
	var diffs []*FieldDiff

	// Everything should be added.
	if actual == nil {
		diffs = append(diffs, &FieldDiff{FieldName: fn.FieldName, ToAdd: desired})
		return diffs, nil
	}

	var toAdd, toRemove []interface{}

	for i, aItem := range actual {
		found := false
		for _, desItem := range desired {
			if ds, _ := Diff(desItem, aItem, info, fn.AddIndex(i)); len(ds) == 0 {
				found = true
				break
			}
		}
		if !found {
			toRemove = append(toRemove, aItem)
		}
	}

	for i, dItem := range desired {
		found := false
		for _, actItem := range actual {
			if ds, _ := Diff(dItem, actItem, info, fn.AddIndex(i)); len(ds) == 0 {
				found = true
				break
			}
		}
		if !found {
			toAdd = append(toAdd, dItem)
		}
	}

	if len(toAdd) > 0 || len(toRemove) > 0 {
		return []*FieldDiff{&FieldDiff{FieldName: fn.FieldName, ToAdd: toAdd, ToRemove: toRemove}}, nil
	}
	return nil, nil
}

// ValueType returns the reflect-style kind of an interface or the underlying type of a pointer.
func ValueType(i interface{}) string {
	if reflect.ValueOf(i).Kind() == reflect.Ptr {
		return reflect.Indirect(reflect.ValueOf(i)).Kind().String()
	}
	return reflect.ValueOf(i).Kind().String()
}

func strs(d, i interface{}) (*string, *string, error) {
	dStr, err := str(d)
	if err != nil {
		return nil, nil, err
	}

	iStr, err := str(i)
	if err != nil {
		return nil, nil, err
	}

	return dStr, iStr, nil
}

func str(d interface{}) (*string, error) {
	dPtr, dOk := d.(*string)
	if !dOk {
		dStr, dOk2 := d.(string)
		if !dOk2 {
			return nil, fmt.Errorf("was given non string %v", d)
		}
		dPtr = String(dStr)
	}
	return dPtr, nil
}

func makeint64(d interface{}) (*int64, error) {
	dPtr, dOk := d.(*int64)
	if !dOk {
		dInt, dOk2 := d.(int64)
		if !dOk2 {
			return nil, fmt.Errorf("was given non int64 %v", d)
		}
		dPtr = Int64(dInt)
	}
	return dPtr, nil
}

func makeint(d interface{}) (*int, error) {
	dPtr, dOk := d.(*int)
	if !dOk {
		dInt, dOk2 := d.(int)
		if !dOk2 {
			return nil, fmt.Errorf("was given non int %v", d)
		}
		dPtr = &dInt
	}
	return dPtr, nil
}

func makefloat64(d interface{}) (*float64, error) {
	dPtr, dOk := d.(*float64)
	if !dOk {
		dFloat, dOk2 := d.(float64)
		if !dOk2 {
			return nil, fmt.Errorf("was given non float64 %v", d)
		}
		dPtr = &dFloat
	}
	return dPtr, nil
}

func bools(d, i interface{}) (*bool, *bool, error) {
	dBool, err := boolean(d)
	if err != nil {
		return nil, nil, err
	}

	iBool, err := boolean(i)
	if err != nil {
		return nil, nil, err
	}

	return dBool, iBool, nil
}

func boolean(d interface{}) (*bool, error) {
	dPtr, dOk := d.(*bool)
	if !dOk {
		return nil, nil
	}
	return dPtr, nil
}

func maps(d, a interface{}) (map[string]interface{}, map[string]interface{}, error) {
	dMap, _ := mapCast(d)
	aMap, _ := mapCast(a)
	return dMap, aMap, nil
}

func mapCast(m interface{}) (map[string]interface{}, error) {
	j, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	var mi map[string]interface{}
	json.Unmarshal(j, &mi)
	return mi, nil
}

func slices(d, i interface{}) ([]interface{}, []interface{}, error) {
	dSlice, err := slice(d)
	if err != nil {
		return nil, nil, err
	}

	iSlice, err := slice(i)
	if err != nil {
		return nil, nil, err
	}

	return dSlice, iSlice, nil
}

func slice(slice interface{}) ([]interface{}, error) {
	// Keep the distinction between nil and empty slice input
	// This isn't going to be an error though.
	if slice == nil {
		return nil, nil
	}

	s := reflect.ValueOf(slice)

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret, nil
}

func addOperationToDiffs(fds []*FieldDiff, i DiffInfo) {
	for _, fd := range fds {
		// Do not overwrite update operations on nested fields with parent field operations.
		if len(fd.ResultingOperation) == 0 {
			fd.ResultingOperation = i.OperationSelector(fd)
		}
	}
}

func mapCompare(d, a map[string]interface{}, ignorePrefixes []string, info DiffInfo, fn FieldName) ([]*FieldDiff, error) {
	var diffs []*FieldDiff
	for k, v := range d {
		if isIgnored(k, ignorePrefixes) {
			continue
		}

		av, ok := a[k]
		if !ok {
			diffs = append(diffs, &FieldDiff{FieldName: fn.FieldName, Message: fmt.Sprintf("%v is missing from actual", k)})
			continue
		}

		objDiffs, err := Diff(v, av, info, fn)
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, objDiffs...)
	}

	for k, v := range a {
		if isIgnored(k, ignorePrefixes) {
			continue
		}

		dv, ok := d[k]
		if !ok {
			diffs = append(diffs, &FieldDiff{FieldName: fn.FieldName, Message: fmt.Sprintf("%v is missing from desired", k)})
			continue
		}

		objDiffs, err := Diff(dv, v, info, fn)
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, objDiffs...)
	}

	return diffs, nil
}
