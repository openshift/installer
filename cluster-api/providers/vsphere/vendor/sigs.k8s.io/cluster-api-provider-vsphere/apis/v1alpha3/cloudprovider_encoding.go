/*
Copyright 2019 The Kubernetes Authors.

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

package v1alpha3

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"
	gcfg "gopkg.in/gcfg.v1"
)

const gcfgTag = "gcfg"

var iniEscapeChars = regexp.MustCompile(`([\\"])`)

// MarshalINI marshals the cloud provider configuration to INI-style
// configuration data.
func (c *CPIConfig) MarshalINI() ([]byte, error) {
	if c == nil {
		return nil, errors.New("config is nil")
	}

	buf := &bytes.Buffer{}

	// Get the reflected type and value of the CPIConfig object.
	configValue := reflect.ValueOf(*c)
	configType := reflect.TypeOf(*c)

	for sectionIndex := 0; sectionIndex < configValue.NumField(); sectionIndex++ {
		sectionType := configType.Field(sectionIndex)
		sectionValue := configValue.Field(sectionIndex)

		// Get the value of the gcfg tag to help determine the section
		// name and whether to omit an empty value. Also ignore fields without the gcfg tag
		sectionName, omitEmpty, hasTag := parseGcfgTag(sectionType)
		if !hasTag {
			continue
		}

		// Do not marshal a section if it is empty.
		if omitEmpty && isEmpty(sectionValue) {
			continue
		}

		switch sectionValue.Kind() {
		case reflect.Map:
			keys := sectionValue.MapKeys()
			sort.Slice(keys, func(i, j int) bool {
				return keys[i].String() < keys[j].String()
			})

			for _, key := range keys {
				sectionNameKey, sectionValue := key, sectionValue.MapIndex(key)
				sectionName := fmt.Sprintf(`%s "%v"`, sectionName, sectionNameKey.String())
				if err := c.marshalINISectionProperties(buf, sectionValue, sectionName); err != nil {
					return nil, err
				}
			}
		default:
			if err := c.marshalINISectionProperties(buf, sectionValue, sectionName); err != nil {
				return nil, err
			}
		}
	}

	return buf.Bytes(), nil
}

func (c *CPIConfig) marshalINISectionProperties(out io.Writer, sectionValue reflect.Value, sectionName string) error {
	switch sectionValue.Kind() {
	case reflect.Interface, reflect.Ptr:
		return c.marshalINISectionProperties(out, sectionValue.Elem(), sectionName)
	}

	fmt.Fprintf(out, "[%s]\n", sectionName)

	sectionType := sectionValue.Type()
	for propertyIndex := 0; propertyIndex < sectionType.NumField(); propertyIndex++ {
		propertyType := sectionType.Field(propertyIndex)
		propertyValue := sectionValue.Field(propertyIndex)

		// Get the value of the gcfg tag to help determine the property
		// name and whether to omit an empty value.
		propertyName, omitEmpty, hasTag := parseGcfgTag(propertyType)
		if !hasTag {
			continue
		}

		// Do not marshal a property if it is empty.
		if omitEmpty && isEmpty(propertyValue) {
			continue
		}

		switch propertyValue.Kind() {
		case reflect.Interface, reflect.Ptr:
			propertyValue = propertyValue.Elem()
		}

		fmt.Fprintf(out, "%s", propertyName)
		if propertyValue.IsValid() {
			rawVal := fmt.Sprintf("%v", propertyValue.Interface())
			val := iniEscapeChars.ReplaceAllString(rawVal, "\\$1")
			val = strings.ReplaceAll(val, "\t", "\\t")
			if propertyValue.Kind() == reflect.String {
				val = "\"" + val + "\""
			}
			fmt.Fprintf(out, " = %s\n", val)
		}
	}

	fmt.Fprintf(out, "\n")

	return nil
}

func parseGcfgTag(field reflect.StructField) (string, bool, bool) {
	name := field.Name
	omitEmpty := false
	hasTag := false

	if tagVal, ok := field.Tag.Lookup(gcfgTag); ok {
		hasTag = true
		tagParts := strings.Split(tagVal, ",")
		lenTagParts := len(tagParts)
		if lenTagParts > 0 {
			tagName := tagParts[0]
			if len(tagName) > 0 && tagName != "-" {
				name = tagName
			}
		}
		if lenTagParts > 1 {
			omitEmpty = tagParts[1] == "omitempty"
		}
	}

	return name, omitEmpty, hasTag
}

// UnmarshalINIOptions defines the options used to influence how INI data is
// unmarshalled.
//
// +kubebuilder:object:generate=false
type UnmarshalINIOptions struct {
	// WarnAsFatal indicates that warnings that occur when unmarshalling INI
	// data should be treated as fatal errors.
	WarnAsFatal bool
}

// UnmarshalINIOptionFunc is used to set unmarshal options.
//
// +kubebuilder:object:generate=false
type UnmarshalINIOptionFunc func(*UnmarshalINIOptions)

// WarnAsFatal sets the option to treat warnings as fatal errors when
// unmarshalling INI data.
func WarnAsFatal(opts *UnmarshalINIOptions) {
	opts.WarnAsFatal = true
}

// UnmarshalINI unmarshals the cloud provider configuration from INI-style
// configuration data.
func (c *CPIConfig) UnmarshalINI(data []byte, optFuncs ...UnmarshalINIOptionFunc) error {
	opts := &UnmarshalINIOptions{}
	for _, setOpts := range optFuncs {
		setOpts(opts)
	}
	var config unmarshallableConfig
	if err := gcfg.ReadStringInto(&config, string(data)); err != nil {
		if opts.WarnAsFatal {
			return err
		}
		if err := gcfg.FatalOnly(err); err != nil {
			return err
		}
	}
	c.Global = config.Global
	c.Network = config.Network
	c.Disk = config.Disk
	c.Workspace = config.Workspace
	c.Labels = config.Labels
	c.VCenter = map[string]CPIVCenterConfig{}
	for k, v := range config.VCenter {
		c.VCenter[k] = *v
	}
	return nil
}

// IsEmpty returns true if an object is its empty value or if a struct, all of
// its fields are their empty values.
func IsEmpty(obj interface{}) bool {
	return isEmpty(reflect.ValueOf(obj))
}

// IsNotEmpty returns true when IsEmpty returns false.
func IsNotEmpty(obj interface{}) bool {
	return !IsEmpty(obj)
}

// isEmpty returns true if an object's fields are all set to their empty values.
func isEmpty(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.Interface, reflect.Ptr:
		return val.IsNil() || isEmpty(val.Elem())

	case reflect.Struct:
		structIsEmpty := true
		for fieldIndex := 0; fieldIndex < val.NumField(); fieldIndex++ {
			if structIsEmpty = isEmpty(val.Field(fieldIndex)); !structIsEmpty {
				break
			}
		}
		return structIsEmpty

	case reflect.Array, reflect.String:
		return val.Len() == 0

	case reflect.Bool:
		return !val.Bool()

	case reflect.Map, reflect.Slice:
		return val.IsNil() || val.Len() == 0

	case reflect.Float32, reflect.Float64:
		return val.Float() == 0

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int() == 0

	default:
		panic(errors.Errorf("invalid kind: %s", val.Kind()))
	}
}

// MarshalCloudProviderArgs marshals the cloud provider arguments for passing
// into a pod spec.
func (cpic *CPICloudConfig) MarshalCloudProviderArgs() []string {
	args := []string{
		"--v=2",
		"--cloud-provider=vsphere",
		"--cloud-config=/etc/cloud/vsphere.conf",
	}
	if cpic.ExtraArgs != nil {
		for k, v := range cpic.ExtraArgs {
			args = append(args, fmt.Sprintf("--%s=%s", k, v))
		}
	}
	return args
}
