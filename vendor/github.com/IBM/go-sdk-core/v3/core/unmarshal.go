/**
 * (C) Copyright IBM Corp. 2020.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package core

import (
	"encoding/base64"
	"fmt"
	"github.com/go-openapi/strfmt"
	"reflect"
)

//
// This file contains a set of methods that are used by generated code to unmarshal common data types.
// Specifically, these methods are designed to fit within the unmarshalling framework implemented by the Go generator.
// For user-defined models, the Go generator will emit code that first unmarshals the operation response body into a 
// generic map (or a slice of maps for a response which is a JSON array), then it invokes the appropriate generated
// "Unmarshal<model-name>()" method to transform the map into an instance of the model (struct).  The generated
// "Unmarshal<model-name>()" methods use the methods in this file to handle primitive and other common data types.
//
// In the methods below that deal with numeric values, they first cast the data value to a "float64" value, then
// eventually transform the float64 value into the appropriate numeric type (int64, float32, etc.).  
// This is done because the json unmarshaller unmarshals all JSON numbers into float64 values within the generic
// map.
//
// For each data type, there are two methods - Unmarshal<type>() and Unmarshal<type>Slice().
//
// Each Unmarshal<type>() method will retrieve the specified property from the generic map, and then return 
// an appropriate value that can then be assigned to a field within a generated struct of that type.
// For example, UnmarshalString() will return a pointer to the string value from the generic map.   UnmarshalInt64()
// will return a pointer to the int64 value, etc.
//
// Each Unmarshal<type>Slice() method will retrieve the specified property from the generic map and interpret it as 
// a slice of <type>.  Within these methods, we first cast the value retrieved from the map as a slice of interface{}
// (i.e. a slice of anything), then we walk through the slice and cast each element to the appropriate type.
// Some of the Unmarshal<type>() and Unmarshal<type>Slice() methods perform an additional transformation on the 
// data values retrieved from the generic map.  For example, the UnmarshalByteArray() method will interpret the 
// JSON value as a string, then perform a base64-decoding of the string to produce the resulting byte array ([]byte).
// 

// Error messages constants.
const (
	errorValueType        = "the '%s' property value should be a %s but was %s"
	errorNotAnArray       = "the '%s' property value should be an array but was %s"
	errorElementValueType = "the '%s' property array element should be a %s but was %s"
	errorDecodeBase64     = "error decoding the '%s' property's base64-encoded string value '%s': %s"
	errorDecodeDate       = "error decoding the '%s' property's Date value '%s': %s"
	errorDecodeDateTime   = "error decoding the '%s' property's DateTime value '%s': %s"
)

// CopyMap returns a shallow copy of the specified map.
// The returned map contains a copy of the mappings found in the original map,
// but uses values from the original map.
func CopyMap(m map[string]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{})
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}

// UnmarshalString retrieves the specified property from the map and returns it as a string
// or nil if the property wasn't found in the map.
func UnmarshalString(m map[string]interface{}, propertyName string) (result *string, err error) {
	var v interface{}
	v, foundIt := m[propertyName]
	if foundIt && v != nil {
		// Interpret the map value as a string.
		s, ok := v.(string)
		if ok {
			result = StringPtr(s)
		} else {
			err = fmt.Errorf(errorValueType, propertyName, "string", reflect.TypeOf(v).String())
		}
	}
	return
}

// UnmarshalStringSlice retrieves the specified property from the map and returns it as a string slice
// or nil if the property wasn't found in the map.
func UnmarshalStringSlice(m map[string]interface{}, propertyName string) (slice []string, err error) {
	var v interface{}
	v, foundIt := m[propertyName]
	if foundIt && v != nil {
		// Interpret the map value as a slice of anything.
		vSlice, ok := v.([]interface{})
		if !ok {
			err = fmt.Errorf(errorNotAnArray, propertyName, reflect.TypeOf(v).String())
			return
		}
		for _, element := range vSlice {
			// Interpret each slice element as a string and add it to the result slice.
			s, ok := element.(string)
			if !ok {
				err = fmt.Errorf(errorElementValueType, propertyName, "string", reflect.TypeOf(element).String())
				return
			}
			slice = append(slice, s)
		}
	}
	return
}

// UnmarshalByteArray retrieves the specified property from the map and returns it as a byte array
// or nil if the property wasn't found in the map.
func UnmarshalByteArray(m map[string]interface{}, propertyName string) (result *[]byte, err error) {
	var v interface{}
	v, foundIt := m[propertyName]
	if foundIt && v != nil {
		// Interpret the map value as a string.
		s, ok := v.(string)
		if ok {
			// Decode the string into a byte array.
			ba, localErr := base64.StdEncoding.DecodeString(s)
			if localErr != nil {
				err = fmt.Errorf(errorDecodeBase64, propertyName, truncateString(s, 16), localErr.Error())
				return
			} else {
				result = &ba
			}
		} else {
			err = fmt.Errorf(errorValueType, propertyName, "base64-encoded string", reflect.TypeOf(v).String())
		}
	}
	return
}

// UnmarshalByteArraySlice retrieves the specified property from the map and returns it as a byte array slice
// or nil if the property wasn't found in the map.
func UnmarshalByteArraySlice(m map[string]interface{}, propertyName string) (slice [][]byte, err error) {
	var v interface{}
	v, foundIt := m[propertyName]
	if foundIt && v != nil {
		// Interpret the map value as a slice of anything.
		vSlice, ok := v.([]interface{})
		if ok {
			for _, element := range vSlice {
				// Interpret each slice element as a string.
				s, ok := element.(string)
				if !ok {
					err = fmt.Errorf(errorElementValueType, propertyName, "base64-encoded string",
						reflect.TypeOf(element).String())
					return
				}
				// Decode the string into a byte array and add it to the result slice.
				ba, localErr := base64.StdEncoding.DecodeString(s)
				if localErr != nil {
					err = fmt.Errorf(errorDecodeBase64, propertyName, truncateString(s, 16), localErr.Error())
					return
				} else {
					slice = append(slice, ba)
				}
			}
		} else {
			err = fmt.Errorf(errorNotAnArray, propertyName, reflect.TypeOf(v).String())
		}
	}
	return
}

// UnmarshalBool retrieves the specified property from the map and returns it as a bool
// or nil if the property wasn't found in the map.
func UnmarshalBool(m map[string]interface{}, propertyName string) (result *bool, err error) {
	var v interface{}
	v, foundIt := m[propertyName]
	if foundIt && v != nil {
		// Interpret the map value as a bool.
		b, ok := v.(bool)
		if ok {
			result = BoolPtr(b)
		} else {
			err = fmt.Errorf(errorValueType, propertyName, "boolean", reflect.TypeOf(v).String())
		}
	}
	return
}

// UnmarshalBoolSlice retrieves the specified property from the map and returns it as a bool slice
// or nil if the property wasn't found in the map.
func UnmarshalBoolSlice(m map[string]interface{}, propertyName string) (slice []bool, err error) {
	var v interface{}
	v, foundIt := m[propertyName]
	if foundIt && v != nil {
		// Interpret the map value as a slice of anything.
		vSlice, ok := v.([]interface{})
		if !ok {
			err = fmt.Errorf(errorNotAnArray, propertyName, reflect.TypeOf(v).String())
			return
		}
		for _, element := range vSlice {
			// Interpret each slice element as a bool and add it to the result slice.
			b, ok := element.(bool)
			if !ok {
				err = fmt.Errorf(errorElementValueType, propertyName, "boolean", reflect.TypeOf(element).String())
				return
			}
			slice = append(slice, b)
		}
	}
	return
}

// UnmarshalInt64 retrieves the specified property from the map and returns it as an int64
// or nil if the property wasn't found in the map.
func UnmarshalInt64(m map[string]interface{}, propertyName string) (result *int64, err error) {
	var v interface{}
	v, foundIt := m[propertyName]
	if foundIt && v != nil {
		// Interpret the map value as a float64 to match the behavior of the JSON unmarshaller.
		f, ok := v.(float64)
		if ok {
			// Convert the value to an int64 and return a pointer to it.
			result = Int64Ptr(int64(f))
		} else {
			err = fmt.Errorf(errorValueType, propertyName, "integer", reflect.TypeOf(v).String())
		}
	}
	return
}

// UnmarshalInt64Slice retrieves the specified property from the map and returns it as an int64 slice
// or nil if the property wasn't found in the map.
func UnmarshalInt64Slice(m map[string]interface{}, propertyName string) (slice []int64, err error) {
	var v interface{}
	v, foundIt := m[propertyName]
	if foundIt && v != nil {
		// Interpret the map value as a slice of anything.
		vSlice, ok := v.([]interface{})
		if !ok {
			err = fmt.Errorf(errorNotAnArray, propertyName, reflect.TypeOf(v).String())
			return
		}
		for _, element := range vSlice {
			// Interpret each slice element as a float64, then convert it to an int64 and add it to the result slice.
			f, ok := element.(float64)
			if !ok {
				err = fmt.Errorf(errorElementValueType, propertyName, "integer", reflect.TypeOf(element).String())
				return
			}
			slice = append(slice, int64(f))
		}
	}
	return
}

// UnmarshalFloat32 retrieves the specified property from the map and returns it as a float32
// or nil if the property wasn't found in the map.
func UnmarshalFloat32(m map[string]interface{}, propertyName string) (result *float32, err error) {
	var v interface{}
	v, foundIt := m[propertyName]
	if foundIt && v != nil {
		// Interpret the map value as a float64 to match the behavior of the JSON unmarshaller.
		f, ok := v.(float64)
		if ok {
			// Convert the value to a float32 and return a pointer to it.
			result = Float32Ptr(float32(f))
		} else {
			err = fmt.Errorf(errorValueType, propertyName, "float32", reflect.TypeOf(v).String())
		}
	}
	return
}

// UnmarshalFloat32Slice retrieves the specified property from the map and returns it as a float32 slice
// or nil if the property wasn't found in the map.
func UnmarshalFloat32Slice(m map[string]interface{}, propertyName string) (slice []float32, err error) {
	var v interface{}
	v, foundIt := m[propertyName]
	if foundIt && v != nil {
		// Interpret the map value as a slice of anything.
		vSlice, ok := v.([]interface{})
		if !ok {
			err = fmt.Errorf(errorNotAnArray, propertyName, reflect.TypeOf(v).String())
			return
		}
		for _, element := range vSlice {
			// Interpret each slice element as a float64, then convert it to a float32 and add it to the result slice.
			f, ok := element.(float64)
			if !ok {
				err = fmt.Errorf(errorElementValueType, propertyName, "float32", reflect.TypeOf(element).String())
				return
			}
			slice = append(slice, float32(f))
		}
	}
	return
}

// UnmarshalFloat64 retrieves the specified property from the map and returns it as a float64
// or nil if the property wasn't found in the map.
func UnmarshalFloat64(m map[string]interface{}, propertyName string) (result *float64, err error) {
	var v interface{}
	v, foundIt := m[propertyName]
	if foundIt && v != nil {
		// Interpret the map value as a float64.
		f, ok := v.(float64)
		if ok {
			result = Float64Ptr(f)
		} else {
			err = fmt.Errorf(errorValueType, propertyName, "float64", reflect.TypeOf(v).String())
		}
	}
	return
}

// UnmarshalFloat64Slice retrieves the specified property from the map and returns it as a float64 slice
// or nil if the property wasn't found in the map.
func UnmarshalFloat64Slice(m map[string]interface{}, propertyName string) (slice []float64, err error) {
	var v interface{}
	v, foundIt := m[propertyName]
	if foundIt && v != nil {
		// Interpret the map value as a slice of anything.
		vSlice, ok := v.([]interface{})
		if !ok {
			err = fmt.Errorf(errorNotAnArray, propertyName, reflect.TypeOf(v).String())
			return
		}
		for _, element := range vSlice {
			// Interpret each slice element as a float64 and add it to the result slice.
			f, ok := element.(float64)
			if !ok {
				err = fmt.Errorf(errorElementValueType, propertyName, "float64", reflect.TypeOf(element).String())
				return
			}
			slice = append(slice, f)
		}
	}
	return
}

// UnmarshalUUID retrieves the specified property from the map and returns it as a UUID
// or nil if the property wasn't found in the map.
func UnmarshalUUID(m map[string]interface{}, propertyName string) (result *strfmt.UUID, err error) {
	var v interface{}
	v, foundIt := m[propertyName]
	if foundIt && v != nil {
		// Interpret the map value as a string.
		s, ok := v.(string)
		if ok {
			// Convert the string to a UUID.
			uuid := strfmt.UUID(s)
			result = &uuid
		} else {
			err = fmt.Errorf(errorValueType, propertyName, "UUID", reflect.TypeOf(v).String())
		}
	}
	return
}

// UnmarshalUUIDSlice retrieves the specified property from the map and returns it as a UUID slice
// or nil if the property wasn't found in the map.
func UnmarshalUUIDSlice(m map[string]interface{}, propertyName string) (slice []strfmt.UUID, err error) {
	var v interface{}
	v, foundIt := m[propertyName]
	if foundIt && v != nil {
		// Interpret the map value as a slice of anything.
		vSlice, ok := v.([]interface{})
		if !ok {
			err = fmt.Errorf(errorNotAnArray, propertyName, reflect.TypeOf(v).String())
			return
		}
		for _, element := range vSlice {
			// Interpret each slice element as a string, then convert it to a UUID value and add it to the result slice.
			s, ok := element.(string)
			if !ok {
				err = fmt.Errorf(errorElementValueType, propertyName, "UUID", reflect.TypeOf(element).String())
				return
			}
			slice = append(slice, strfmt.UUID(s))
		}
	}
	return
}

// UnmarshalDate retrieves the specified property from the map and returns it as a Date
// or nil if the property wasn't found in the map.
func UnmarshalDate(m map[string]interface{}, propertyName string) (result *strfmt.Date, err error) {
	var v interface{}
	v, foundIt := m[propertyName]
	if foundIt && v != nil {
		// Interpret the map value as a string.
		s, ok := v.(string)
		if ok {
			// Convert the string to a Date value.
			var date strfmt.Date
			localErr := date.UnmarshalText([]byte(s))
			if localErr != nil {
				err = fmt.Errorf(errorDecodeDate, propertyName, truncateString(s, 16), localErr.Error())
			} else {
				result = &date
			}
		} else {
			err = fmt.Errorf(errorValueType, propertyName, "Date", reflect.TypeOf(v).String())
		}
	}
	return
}

// UnmarshalDateSlice retrieves the specified property from the map and returns it as a Date slice
// or nil if the property wasn't found in the map.
func UnmarshalDateSlice(m map[string]interface{}, propertyName string) (slice []strfmt.Date, err error) {
	var v interface{}
	v, foundIt := m[propertyName]
	if foundIt && v != nil {
		// Interpret the map value as a slice of anything.
		vSlice, ok := v.([]interface{})
		if !ok {
			err = fmt.Errorf(errorNotAnArray, propertyName, reflect.TypeOf(v).String())
			return
		}
		for _, element := range vSlice {
			// Interpret each slice element to a string, convert it to a Date and add it to the result slice.
			s, ok := element.(string)
			if !ok {
				err = fmt.Errorf(errorElementValueType, propertyName, "Date", reflect.TypeOf(element).String())
				return
			}
			var date strfmt.Date
			localErr := date.UnmarshalText([]byte(s))
			if localErr != nil {
				err = fmt.Errorf(errorDecodeDate, propertyName, truncateString(s, 16), localErr.Error())
				return
			} else {
				slice = append(slice, date)
			}
		}
	}
	return
}

// UnmarshalDateTime retrieves the specified property from the map and returns it as a DateTime
// or nil if the property wasn't found in the map.
func UnmarshalDateTime(m map[string]interface{}, propertyName string) (result *strfmt.DateTime, err error) {
	var v interface{}
	v, foundIt := m[propertyName]
	if foundIt && v != nil {
		// Interpret the map value as a string.
		s, ok := v.(string)
		if ok {
			// Convert the string to a DateTime value.
			var datetime strfmt.DateTime
			localErr := datetime.UnmarshalText([]byte(s))
			if localErr != nil {
				err = fmt.Errorf(errorDecodeDateTime, propertyName, truncateString(s, 16), localErr.Error())
			} else {
				result = &datetime
			}
		} else {
			err = fmt.Errorf(errorValueType, propertyName, "DateTime", reflect.TypeOf(v).String())
		}
	}
	return
}

// UnmarshalDateTimeSlice retrieves the specified property from the map and returns it as a DateTime slice
// or nil if the property wasn't found in the map.
func UnmarshalDateTimeSlice(m map[string]interface{}, propertyName string) (slice []strfmt.DateTime, err error) {
	var v interface{}
	v, foundIt := m[propertyName]
	if foundIt && v != nil {
		// Interpret the map value as a slice of anything.
		vSlice, ok := v.([]interface{})
		if !ok {
			err = fmt.Errorf(errorNotAnArray, propertyName, reflect.TypeOf(v).String())
			return
		}
		for _, element := range vSlice {
			// Interpret each slice element as a string, convert to a DateTime and add it to the result slice.
			s, ok := element.(string)
			if !ok {
				err = fmt.Errorf(errorElementValueType, propertyName, "DateTime", reflect.TypeOf(element).String())
				return
			}
			var datetime strfmt.DateTime
			localErr := datetime.UnmarshalText([]byte(s))
			if localErr != nil {
				err = fmt.Errorf(errorDecodeDateTime, propertyName, truncateString(s, 16), localErr.Error())
				return
			} else {
				slice = append(slice, datetime)
			}
		}
	}
	return
}

// UnmarshalObject retrieves the specified property from the map and returns it as a generic
// object (i.e. map[string]interface{}), or nil if the property wasn't found in the map.
func UnmarshalObject(m map[string]interface{}, propertyName string) (result map[string]interface{}, err error) {
	var v interface{}
	v, foundIt := m[propertyName]
	if foundIt && v != nil {
		// Interpret the map value as a generic map containing a JSON object.
		jsonMap, ok := v.(map[string]interface{})
		if ok {
			result = jsonMap
		} else {
			err = fmt.Errorf(errorValueType, propertyName, "JSON object", reflect.TypeOf(v).String())
		}
	}
	return
}

// UnmarshalObjectSlice retrieves the specified property from the map and returns it as a slice of
// generic objects (i.e. []map[string]interface{}), or nil if the property wasn't found in the map.
func UnmarshalObjectSlice(m map[string]interface{}, propertyName string) (slice []map[string]interface{}, err error) {
	v, foundIt := m[propertyName]
	if foundIt && v != nil {
		// Interpret the map value as a slice of anything.
		vSlice, ok := v.([]interface{})
		if !ok {
			err = fmt.Errorf(errorNotAnArray, propertyName, reflect.TypeOf(v).String())
			return
		}
		for _, element := range vSlice {
			// Interpret each slice element as a map and then add it to the result slice.
			jsonMap, ok := element.(map[string]interface{})
			if !ok {
				err = fmt.Errorf(errorElementValueType, propertyName, "JSON object", reflect.TypeOf(element).String())
				return
			}
			slice = append(slice, jsonMap)
		}
	}
	return
}

// UnmarshalAny retrieves the specified property from the map and returns it as a generic
// value (i.e. interface{}), or nil if the property wasn't found in the map.
func UnmarshalAny(m map[string]interface{}, propertyName string) (result interface{}, err error) {
	var v interface{}
	v, foundIt := m[propertyName]
	if foundIt && v != nil {
		result = v
	}
	return
}

// UnmarshalAnySlice retrieves the specified property from the map and returns it as a slice of
// generic values (i.e. []interface{}), or nil if the property wasn't found in the map.
func UnmarshalAnySlice(m map[string]interface{}, propertyName string) (slice []interface{}, err error) {
	v, foundIt := m[propertyName]
	if foundIt && v != nil {
		// Interpret the map value as a slice of anything.
		vSlice, ok := v.([]interface{})
		if !ok {
			err = fmt.Errorf(errorNotAnArray, propertyName, reflect.TypeOf(v).String())
			return
		}
		slice = vSlice
	}
	return
}

// truncateString returns a string suitable for inclusion in an error message.
// If the input string is longer than the specified length, we'll just return the first <length> 
// bytes followed by "...".
func truncateString(s string, length int) string {
	if len(s) > length {
		runes := []rune(s)
		return string(runes[0:length]) + "..."
	}
	return s
}
