/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package genruntime

import (
	"bytes"
	"encoding/json"

	"github.com/pkg/errors"
)

// PropertyBag is an unordered set of stashed information that used for properties not directly supported by storage
// resources, allowing for full fidelity round trip conversions
type PropertyBag map[string]string

// Background:
// We store items in the bag as serialized JSON, which we can then deserialize in a just-in-time fashion once we know
// the type of the instance we're going to populate. Unlike other platforms, Go doesn't embed type information as it
// serializes to JSON or YAML, which means that deserialization requires a type hint that's not available when our
// containing resource is hydrated. We only have the required type available when we are doing the conversion to a
// related type.
// This comment kept separate from the definition above so that it doesn't get copied into the generated YAML files.

// PropertyBag returns a new property bag
// originals is a (potentially empty) sequence of existing property bags who's content will be copied into the new
// property bag. In the case of key overlaps, values from bags later in the parameter list overwrite the earlier value.
func NewPropertyBag(originals ...PropertyBag) PropertyBag {
	result := make(PropertyBag)

	for _, orig := range originals {
		for k, v := range orig {
			result[k] = v
		}
	}

	return result
}

// Contains returns true if the specified name is present in the bag; false otherwise
func (bag PropertyBag) Contains(name string) bool {
	_, found := bag[name]
	return found
}

// Add is used to add a value into the bag; exact formatting depends on the type.
// Any existing value will be overwritten.
// property is the name of the item to put into the bag
// value is the instance to be stashed away for later
func (bag PropertyBag) Add(property string, value interface{}) error {
	switch v := value.(type) {
	case string:
		bag[property] = v
	default:
		// Default to treating as a JSON blob
		j, err := json.Marshal(v)
		if err != nil {
			return errors.Wrapf(err, "adding %s as JSON", property)
		}
		bag[property] = string(j)
	}

	return nil
}

// Pull removes a value from the bag, using it to populate the destination
// property is the name of the item to remove and return
// destination should be a pointer to where the value is to be placed
// If the item is present and successfully deserialized, returns no error (nil); otherwise returns an error.
// If an error happens deserializing an item from the bag, it is still removed from the bag.
func (bag PropertyBag) Pull(property string, destination interface{}) error {
	value, found := bag[property]
	if !found {
		// Property not found in the bag
		return errors.Errorf("property bag does not contain %q", property)
	}

	// Property found, remove the value
	delete(bag, property)

	switch d := destination.(type) {
	case *string:
		*d = value
	default:
		data := []byte(value)
		decoder := json.NewDecoder(bytes.NewReader(data))
		decoder.DisallowUnknownFields()
		err := decoder.Decode(destination)
		if err != nil {
			return errors.Wrapf(err, "pulling %q from PropertyBag", property)
		}
	}

	return nil
}

// Remove ensures the property bag doesn't contain a value for the specified name
// property is the name of the item to remove
// It is not an error to try and remove an item that's not present
func (bag PropertyBag) Remove(property string) {
	delete(bag, property)
}
