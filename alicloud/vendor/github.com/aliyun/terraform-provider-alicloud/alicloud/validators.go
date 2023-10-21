package alicloud

import (
	"encoding/json"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"gopkg.in/yaml.v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// validateCIDRNetworkAddress ensures that the string value is a valid CIDR that
// represents a network address - it adds an error otherwise
func validateCIDRNetworkAddress(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid CIDR, got error parsing: %s", k, err))
		return
	}

	if ipnet == nil || value != ipnet.String() {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid network CIDR, expected %q, got %q",
			k, ipnet, value))
	}

	return
}

func validateVpnCIDRNetworkAddress(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	cidrs := strings.Split(value, ",")
	for _, cidr := range cidrs {
		_, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid CIDR, got error parsing: %s", k, err))
			return
		}

		if ipnet == nil || cidr != ipnet.String() {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid network CIDR, expected %q, got %q",
				k, ipnet, cidr))
			return
		}
	}

	return
}

func validateSwitchCIDRNetworkAddress(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid CIDR, got error parsing: %s", k, err))
		return
	}

	if ipnet == nil || value != ipnet.String() {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid network CIDR, expected %q, got %q",
			k, ipnet, value))
		return
	}

	mark, _ := strconv.Atoi(strings.Split(ipnet.String(), "/")[1])
	if mark < 16 || mark > 29 {
		errors = append(errors, fmt.Errorf(
			"%q must contain a network CIDR which mark between 16 and 29",
			k))
	}

	return
}

func validateAllowedSplitStringValue(ss []string, splitStr string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(string)
		existed := false
		tsList := strings.Split(value, splitStr)

		for _, ts := range tsList {
			existed = false
			for _, s := range ss {
				if ts == s {
					existed = true
					break
				}
			}
		}
		if !existed {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid string value should in %#v, got %q",
				k, ss, value))
		}
		return

	}
}

func validateStringConvertInt64() schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		if value, ok := v.(string); ok {
			_, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				errors = append(errors, fmt.Errorf(
					"%q should be convert to int64, got %q", k, value))
			}
		} else {
			errors = append(errors, fmt.Errorf(
				"%q should be convert to string, got %q", k, value))
		}

		return
	}
}

func validateForwardPort(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "any" {
		valueConv, err := strconv.Atoi(value)
		if err != nil || valueConv < 1 || valueConv > 65535 {
			errors = append(errors, fmt.Errorf("%q must be a valid port between 1 and 65535 or any ", k))
		}
	}
	return
}

func validateOssBucketDateTimestamp(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, err := time.Parse("2006-01-02", value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q cannot be parsed as date YYYY-MM-DD Format", value))
	}
	return
}

func validateOnsGroupId(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !(strings.HasPrefix(value, "GID-") || strings.HasPrefix(value, "GID_")) {
		errors = append(errors, fmt.Errorf("%q is invalid, it must start with 'GID-' or 'GID_'", k))
	}
	if reg := regexp.MustCompile(`^[\w\-]{7,64}$`); !reg.MatchString(value) {
		errors = append(errors, fmt.Errorf("%q length is limited to 7-64 and only characters such as letters, digits, '_' and '-' are allowed", k))
	}
	return
}

func validateRR(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if strings.HasPrefix(value, "-") || strings.HasSuffix(value, "-") {
		errors = append(errors, fmt.Errorf("RR is invalid, it can not starts or ends with '-'"))
	}

	if len(value) > 253 {
		errors = append(errors, fmt.Errorf("RR can not longer than 253 characters."))
	}

	for _, part := range strings.Split(value, ".") {
		if len(part) > 63 {
			errors = append(errors, fmt.Errorf("Each part of RR split with . can not longer than 63 characters."))
			return
		}
	}
	return
}

// Takes a value containing JSON string and passes it through
// the JSON parser to normalize it, returns either a parsing
// error or normalized JSON string.
func normalizeYamlString(yamlString interface{}) (string, error) {
	var j interface{}

	if yamlString == nil || yamlString.(string) == "" {
		return "", nil
	}

	s := yamlString.(string)

	err := yaml.Unmarshal([]byte(s), &j)
	if err != nil {
		return s, err
	}

	// The error is intentionally ignored here to allow empty policies to passthrough validation.
	// This covers any interpolated values
	bytes, _ := yaml.Marshal(j)

	return string(bytes[:]), nil
}

// Takes a value containing JSON string and passes it through
// the JSON parser to normalize it, returns either a parsing
// error or normalized JSON string.
func normalizeJsonString(jsonString interface{}) (string, error) {
	var j interface{}

	if jsonString == nil || jsonString.(string) == "" {
		return "", nil
	}

	s := jsonString.(string)

	err := json.Unmarshal([]byte(s), &j)
	if err != nil {
		return s, err
	}

	// The error is intentionally ignored here to allow empty policies to passthrough validation.
	// This covers any interpolated values
	bytes, _ := json.Marshal(j)

	return string(bytes[:]), nil
}

func validateYamlString(v interface{}, k string) (ws []string, errors []error) {
	if _, err := normalizeYamlString(v); err != nil {
		errors = append(errors, fmt.Errorf("%q contains an invalid YAML: %s", k, err))
	}

	return
}

func validateDBConnectionPort(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		port, err := strconv.Atoi(value)
		if err != nil {
			errors = append(errors, err)
		}
		if port < 1000 || port > 5999 {
			errors = append(errors, fmt.Errorf("%q cannot be less than 3001 and larger than 3999.", k))
		}
	}
	return
}

func validateSslVpnPortValue(is []int) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		ws, errors = validation.IntBetween(1, 65535)(v, k)
		if errors != nil {
			return
		}

		value := v.(int)
		for _, i := range is {
			if i == value {
				errors = append(errors, fmt.Errorf(
					"%q must contain a valid int value should not be in array %#v, got %q",
					k, is, value))
				return
			}
		}
		return

	}
}

// below copy/pasta from https://github.com/hashicorp/terraform-plugin-sdk/blob/master/helper/validation/validation.go
// alicloud vendor contains very old version of Terraform which lacks this functions

// IntBetween returns a SchemaValidateFunc which tests if the provided value
// is of type int and is between min and max (inclusive)
func intBetween(min, max int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(int)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be int", k))
			return
		}

		if v < min || v > max {
			es = append(es, fmt.Errorf("expected %s to be in the range (%d - %d), got %d", k, min, max, v))
			return
		}

		return
	}
}

// Validate length(2~128) and prefix of the name.
func validateNormalName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) < 2 || len(value) > 128 {
		errors = append(errors, fmt.Errorf("%s cannot be longer than 128 characters", k))
	}
	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		errors = append(errors, fmt.Errorf("%s cannot starts with http:// or https://", k))
	}
	return
}
