// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	homedir "github.com/mitchellh/go-homedir"
	gouuid "github.com/satori/go.uuid"

	"github.com/IBM-Cloud/bluemix-go/helpers"
)

var (
	validHRef *regexp.Regexp
)

func init() {
	validHRef = regexp.MustCompile(`^http(s)?:\/\/([^\/?#]*)([^?#]*)(\?([^#]*))?(#(.*))?$`)
}

func validateSecondaryIPCount(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value != 4 && value != 8 {
		errors = append(errors, fmt.Errorf(
			"%q must be either 4 or 8", k))
	}
	return
}

func validateServiceTags(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 2048 {
		errors = append(errors, fmt.Errorf(
			"%q must contain tags whose maximum length is 2048 characters", k))
	}
	return
}

func validateAllowedStringValue(validValues []string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		input := v.(string)
		existed := false
		for _, s := range validValues {
			if s == input {
				existed = true
				break
			}
		}
		if !existed {
			errors = append(errors, fmt.Errorf(
				"%q must contain a value from %#v, got %q",
				k, validValues, input))
		}
		return

	}
}

func validateRegexpLen(min, max int, regex string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(string)

		acceptedcharacters, _ := regexp.MatchString(regex, value)

		if acceptedcharacters {
			if (len(value) < min) || (len(value) > max) && (min > 0 && max > 0) {
				errors = append(errors, fmt.Errorf(
					"%q (%q) must contain from %d to %d characters ", k, value, min, max))
			}
		} else {
			errors = append(errors, fmt.Errorf(
				"%q (%q) should match regexp %s ", k, v, regex))
		}

		return

	}
}

func validateAllowedIntValue(is []int) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(int)
		existed := false
		for _, i := range is {
			if i == value {
				existed = true
				break
			}
		}
		if !existed {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid int value should in array %#v, got %q",
				k, is, value))
		}
		return

	}
}

func validateAllowedEnterpriseNameValue() schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(string)

		if len(value) < 3 || len(value) > 60 {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid string value with length between 3 and 60", value))
		}
		return

	}
}
func validateRoutePath(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	//Somehow API allows this
	if value == "" {
		return
	}

	if (len(value) < 2) || (len(value) > 128) {
		errors = append(errors, fmt.Errorf(
			"%q (%q) must contain from 2 to 128 characters ", k, value))
	}
	if !(strings.HasPrefix(value, "/")) {
		errors = append(errors, fmt.Errorf(
			"%q (%q) must start with a forward slash '/'", k, value))

	}
	if strings.Contains(value, "?") {
		errors = append(errors, fmt.Errorf(
			"%q (%q) must not contain a '?'", k, value))
	}

	return
}

func validateRoutePort(v interface{}, k string) (ws []string, errors []error) {
	return validatePortRange(1024, 65535)(v, k)
}

func validateAppPort(v interface{}, k string) (ws []string, errors []error) {
	return validatePortRange(1024, 65535)(v, k)
}
func validateLBListenerPolicyPriority(v interface{}, k string) (ws []string, errors []error) {
	interval := v.(int)
	if interval < 1 || interval > 10 {
		errors = append(errors, fmt.Errorf(
			"%q must be between 1 and 10",
			k))
	}
	return
}

func validateStringLength(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if (len(value) < 1) || (len(value) > 128) {
		errors = append(errors, fmt.Errorf(
			"%q (%q) must contain from 1 to 128 characters ", k, value))
	}
	return
}

func validatePortRange(start, end int) func(v interface{}, k string) (ws []string, errors []error) {
	f := func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(int)
		if (value < start) || (value > end) {
			errors = append(errors, fmt.Errorf(
				"%q (%d) must be in the range of %d to %d", k, value, start, end))
		}
		return
	}
	return f
}

func validateDomainName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if !(strings.Contains(value, ".")) {
		errors = append(errors, fmt.Errorf(
			"%q (%q) must contain a '.',example.com,foo.example.com", k, value))
	}

	return
}

func validateAppInstance(v interface{}, k string) (ws []string, errors []error) {
	instances := v.(int)
	if instances < 0 {
		errors = append(errors, fmt.Errorf(
			"%q (%q) must be greater than 0", k, instances))
	}
	return

}

func validateWorkerNum(v interface{}, k string) (ws []string, errors []error) {
	workerNum := v.(int)
	if workerNum <= 0 {
		errors = append(errors, fmt.Errorf(
			"%q  must be greater than 0", k))
	}
	return

}

func validateAppZipPath(v interface{}, k string) (ws []string, errors []error) {
	path := v.(string)
	applicationZip, err := homedir.Expand(path)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q (%q) home directory in the given path couldn't be expanded", k, path))
	}
	if !helpers.FileExists(applicationZip) {
		errors = append(errors, fmt.Errorf(
			"%q (%q) doesn't exist", k, path))
	}

	return

}

func validateNotes(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 1000 {
		errors = append(errors, fmt.Errorf(
			"%q should not exceed 1000 characters", k))
	}
	return
}

func validatePublicBandwidth(v interface{}, k string) (ws []string, errors []error) {
	bandwidth := v.(int)
	if bandwidth < 0 {
		errors = append(errors, fmt.Errorf(
			"%q (%q) must be greater than 0", k, bandwidth))
		return
	}
	validBandwidths := []int{250, 1000, 5000, 10000, 20000}
	for _, b := range validBandwidths {
		if b == bandwidth {
			return
		}
	}
	errors = append(errors, fmt.Errorf(
		"%q (%d) must be one of the value from %d", k, bandwidth, validBandwidths))
	return

}

func validateMaxConn(v interface{}, k string) (ws []string, errors []error) {
	maxConn := v.(int)
	if maxConn < 1 || maxConn > 64000 {
		errors = append(errors, fmt.Errorf(
			"%q must be between 1 and 64000",
			k))
		return
	}
	return
}

func validateKeyLifeTime(v interface{}, k string) (ws []string, errors []error) {
	secs := v.(int)
	if secs < 1800 || secs > 86400 {
		errors = append(errors, fmt.Errorf(
			"%q must be between 1800 and 86400",
			k))
		return
	}
	return
}

func validateWeight(v interface{}, k string) (ws []string, errors []error) {
	weight := v.(int)
	if weight < 0 || weight > 100 {
		errors = append(errors, fmt.Errorf(
			"%q must be between 1 and 100",
			k))
	}
	return
}

func validateSizePerZone(v interface{}, k string) (ws []string, errors []error) {
	sizePerZone := v.(int)
	if sizePerZone <= 0 {
		errors = append(errors, fmt.Errorf(
			"%q must be greater than 0",
			k))
	}
	return
}

func validateInterval(v interface{}, k string) (ws []string, errors []error) {
	interval := v.(int)
	if interval < 2 || interval > 60 {
		errors = append(errors, fmt.Errorf(
			"%q must be between 2 and 60",
			k))
	}
	return
}

func validateMaxRetries(v interface{}, k string) (ws []string, errors []error) {
	maxRetries := v.(int)
	if maxRetries < 1 || maxRetries > 10 {
		errors = append(errors, fmt.Errorf(
			"%q must be between 1 and 10",
			k))
	}
	return
}

func validateTimeout(v interface{}, k string) (ws []string, errors []error) {
	timeout := v.(int)
	if timeout < 1 || timeout > 59 {
		errors = append(errors, fmt.Errorf(
			"%q must be between 1 and 59",
			k))
	}
	return
}

func validateURLPath(v interface{}, k string) (ws []string, errors []error) {
	urlPath := v.(string)
	if len(urlPath) > 250 || !strings.HasPrefix(urlPath, "/") {
		errors = append(errors, fmt.Errorf(
			"%q should start with ‘/‘ and has a max length of 250 characters.",
			k))
	}
	return
}

func validateSecurityRuleDirection(v interface{}, k string) (ws []string, errors []error) {
	validDirections := map[string]bool{
		"ingress": true,
		"egress":  true,
	}

	value := v.(string)
	_, found := validDirections[value]
	if !found {
		strarray := make([]string, 0, len(validDirections))
		for key := range validDirections {
			strarray = append(strarray, key)
		}
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid security group rule direction %q. Valid types are %q.",
			k, value, strings.Join(strarray, ",")))
	}
	return
}

func validateSecurityRuleEtherType(v interface{}, k string) (ws []string, errors []error) {
	validEtherTypes := map[string]bool{
		"IPv4": true,
		"IPv6": true,
	}

	value := v.(string)
	_, found := validEtherTypes[value]
	if !found {
		strarray := make([]string, 0, len(validEtherTypes))
		for key := range validEtherTypes {
			strarray = append(strarray, key)
		}
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid security group rule ethernet type %q. Valid types are %q.",
			k, value, strings.Join(strarray, ",")))
	}
	return
}

//validateIP...
func validateIP(v interface{}, k string) (ws []string, errors []error) {
	address := v.(string)
	if net.ParseIP(address) == nil {
		errors = append(errors, fmt.Errorf(
			"%q must be a valid ip address",
			k))
	}
	return
}

//validateCIDR...
func validateCIDR(v interface{}, k string) (ws []string, errors []error) {
	address := v.(string)
	_, _, err := net.ParseCIDR(address)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q must be a valid cidr address",
			k))
	}
	return
}

//validateCIDRAddress...
func validateCIDRAddress() schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		address := v.(string)
		_, _, err := net.ParseCIDR(address)
		if err != nil {
			errors = append(errors, fmt.Errorf(
				"%q must be a valid cidr address",
				k))
		}
		return
	}
}

//validateOverlappingAddress...
func validateOverlappingAddress() schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		nonOverlappingCIDR := map[string]bool{
			"127.0.0.0/8":    true,
			"161.26.0.0/16":  true,
			"166.8.0.0/14":   true,
			"169.254.0.0/16": true,
			"224.0.0.0/4":    true,
		}

		address := v.(string)
		_, found := nonOverlappingCIDR[address]
		if found {
			errors = append(errors, fmt.Errorf(
				"%q the request is overlapping with reserved address ranges",
				k))
		}
		return
	}
}

//validateRemoteIP...
func validateRemoteIP(v interface{}, k string) (ws []string, errors []error) {
	_, err1 := validateCIDR(v, k)
	_, err2 := validateIP(v, k)

	if len(err1) != 0 && len(err2) != 0 {
		errors = append(errors, fmt.Errorf(
			"%q must be a valid remote ip address (cidr or ip)",
			k))
	}
	return
}

//validateIPorCIDR...
func validateIPorCIDR() schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		_, err1 := validateCIDR(v, k)
		_, err2 := validateIP(v, k)

		if len(err1) != 0 && len(err2) != 0 {
			errors = append(errors, fmt.Errorf(
				"%q must be a valid remote ip address (cidr or ip)",
				k))
		}
		return
	}
}

func validateSecurityRuleProtocol(v interface{}, k string) (ws []string, errors []error) {
	validProtocols := map[string]bool{
		"icmp": true,
		"tcp":  true,
		"udp":  true,
	}

	value := v.(string)
	_, found := validProtocols[value]
	if !found {
		strarray := make([]string, 0, len(validProtocols))
		for key := range validProtocols {
			strarray = append(strarray, key)
		}
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid security group rule ethernet type %q. Valid types are %q.",
			k, value, strings.Join(strarray, ",")))
	}
	return
}

func validateNamespace(ns string) error {
	os := strings.Split(ns, "_")
	if len(os) < 2 || (len(os) == 2 && (len(os[0]) == 0 || len(os[1]) == 0)) {
		return fmt.Errorf(
			"Namespace is (%s), it must be of the form <org>_<space>, provider can't find the auth key if you use _ as well", ns)
	}
	return nil
}

//func validateJSONString(v interface{}, k string) (ws []string, errors []error) {
//	if _, err := normalizeJSONString(v); err != nil {
//		errors = append(errors, fmt.Errorf("%q contains an invalid JSON: %s", k, err))
//	}
//	if err := validateKeyValue(v); err != nil {
//		errors = append(errors, fmt.Errorf("%q contains an invalid JSON: %s", k, err))
//	}
//	return
//}

func validateJSONString() schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		if _, err := normalizeJSONString(v); err != nil {
			errors = append(errors, fmt.Errorf("%q contains an invalid JSON: %s", k, err))
		}
		if err := validateKeyValue(v); err != nil {
			errors = append(errors, fmt.Errorf("%q contains an invalid JSON: %s", k, err))
		}
		return
	}
}

func validateRegexp(regex string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(string)

		acceptedcharacters, _ := regexp.MatchString(regex, value)

		if !acceptedcharacters {
			errors = append(errors, fmt.Errorf(
				"%q (%q) should match regexp %s ", k, v, regex))
		}

		return

	}
}

// NoZeroValues is a SchemaValidateFunc which tests if the provided value is
// not a zero value. It's useful in situations where you want to catch
// explicit zero values on things like required fields during validation.
func validateNoZeroValues() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (ws []string, errors []error) {

		if reflect.ValueOf(i).Interface() == reflect.Zero(reflect.TypeOf(i)).Interface() {
			switch reflect.TypeOf(i).Kind() {
			case reflect.String:
				errors = append(errors, fmt.Errorf("%s value must not be empty.", k))
			case reflect.Int, reflect.Float64:
				errors = append(errors, fmt.Errorf("%s value must not be zero.", k))
			default:
				// this validator should only ever be applied to TypeString, TypeInt and TypeFloat
				errors = append(errors, fmt.Errorf("can't use NoZeroValues with %T attribute %s", k, i))
			}
		}
		return
	}
}

func validateBindedPackageName() schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(string)

		if !(strings.HasPrefix(value, "/")) {
			errors = append(errors, fmt.Errorf(
				"%q (%q) must start with a forward slash '/'.The package name should be '/whisk.system/cloudant', '/test@in.ibm.com_new/utils' or '/_/utils'", k, value))

		}

		index := strings.LastIndex(value, "/")

		if index < 2 || index == len(value)-1 {
			errors = append(errors, fmt.Errorf(
				"%q (%q) is not a valid bind package name.The package name should be '/whisk.system/cloudant','/test@in.ibm.com_new/utils' or '/_/utils'", k, value))

		}

		return
	}
}

func validateKeyValue(jsonString interface{}) error {
	var j [](map[string]interface{})
	if jsonString == nil || jsonString.(string) == "" {
		return nil
	}
	s := jsonString.(string)
	err := json.Unmarshal([]byte(s), &j)
	if err != nil {
		return err
	}
	for _, v := range j {
		_, exists := v["key"]
		if !exists {
			return errors.New("'key' is missing from json")
		}
		_, exists = v["value"]
		if !exists {
			return errors.New("'value' is missing from json")
		}
	}
	return nil
}

func validateActionName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if strings.HasPrefix(value, "/") {
		errors = append(errors, fmt.Errorf(
			"%q (%q) must not start with a forward slash '/'.The action name should be like 'myaction' or utils/cloudant'", k, value))

	}

	const alphaNumeric = "abcdefghijklmnopqrstuvwxyz0123456789/_@.-"

	for _, char := range value {
		if !strings.Contains(alphaNumeric, strings.ToLower(string(char))) {
			errors = append(errors, fmt.Errorf(
				"%q (%q) The name of the package contains illegal characters", k, value))
		}
	}

	return
}

func validateActionKind(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	kindList := []string{"php:7.3", "nodejs:8", "swift:3", "nodejs", "blackbox", "java", "sequence", "nodejs:10", "python:3", "python", "python:2", "swift", "swift:4.2"}
	if !stringInSlice(value, kindList) {
		errors = append(errors, fmt.Errorf(
			"%q (%q) Invalid kind is provided.Supported list of kinds of actions are (%q)", k, value, kindList))
	}
	return
}

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func validateFunctionName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	var validName = regexp.MustCompile(`\A([\w]|[\w][\w@ .-]*[\w@.-]+)\z`)
	if !validName.MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q (%q) The name contains illegal characters", k, value))

	}
	return
}

func validateStorageType(v interface{}, k string) (ws []string, errors []error) {
	validEtherTypes := map[string]bool{
		"Endurance":   true,
		"Performance": true,
	}

	value := v.(string)
	_, found := validEtherTypes[value]
	if !found {
		strarray := make([]string, 0, len(validEtherTypes))
		for key := range validEtherTypes {
			strarray = append(strarray, key)
		}
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid storage type %q. Valid types are %q.",
			k, value, strings.Join(strarray, ",")))
	}
	return
}

func validateRole(v interface{}, k string) (ws []string, errors []error) {
	validRolesTypes := map[string]bool{
		"Writer":        true,
		"Reader":        true,
		"Manager":       true,
		"Administrator": true,
		"Operator":      true,
		"Viewer":        true,
		"Editor":        true,
	}

	value := v.(string)
	_, found := validRolesTypes[value]
	if !found {
		strarray := make([]string, 0, len(validRolesTypes))
		for key := range validRolesTypes {
			strarray = append(strarray, key)
		}
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid role %q. Valid roles are %q.",
			k, value, strings.Join(strarray, ",")))
	}
	return
}

func validateDayOfWeek(v interface{}, k string) (ws []string, errors []error) {
	validDayTypes := map[string]bool{
		"SUNDAY":    true,
		"MONDAY":    true,
		"TUESDAY":   true,
		"WEDNESDAY": true,
		"THURSDAY":  true,
		"FRIDAY":    true,
		"SATURDAY":  true,
	}

	value := v.(string)
	_, found := validDayTypes[value]
	if !found {
		strarray := make([]string, 0, len(validDayTypes))
		for key := range validDayTypes {
			strarray = append(strarray, key)
		}
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid day %q. Valid days are %q.",
			k, value, strings.Join(strarray, ",")))
	}
	return
}

func validateScheduleType(v interface{}, k string) (ws []string, errors []error) {
	validSchdTypes := map[string]bool{
		"HOURLY": true,
		"DAILY":  true,
		"WEEKLY": true,
	}

	value := v.(string)
	_, found := validSchdTypes[value]
	if !found {
		strarray := make([]string, 0, len(validSchdTypes))
		for key := range validSchdTypes {
			strarray = append(strarray, key)
		}
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid schedule type %q. Valid schedules are %q.",
			k, value, strings.Join(strarray, ",")))
	}
	return
}

func validateHour(start, end int) func(v interface{}, k string) (ws []string, errors []error) {
	f := func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(int)
		if (value < start) || (value > end) {
			errors = append(errors, fmt.Errorf(
				"%q (%d) must be in the range of %d to %d", k, value, start, end))
		}
		return
	}
	return f
}

func validateMinute(start, end int) func(v interface{}, k string) (ws []string, errors []error) {
	f := func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(int)
		if (value < start) || (value > end) {
			errors = append(errors, fmt.Errorf(
				"%q (%d) must be in the range of %d to %d", k, value, start, end))
		}
		return
	}
	return f
}

func validateDatacenterOption(v []interface{}, allowedValues []string) error {
	for _, option := range v {
		if option == nil {
			return fmt.Errorf("Provide a valid `datacenter_choice`")
		}
		values := option.(map[string]interface{})
		for k := range values {
			if !stringInSlice(k, allowedValues) {
				return fmt.Errorf(
					"%q Invalid values are provided in `datacenter_choice`. Supported list of keys are (%q)", k, allowedValues)
			}

		}
	}
	return nil
}

func validateLBTimeout(v interface{}, k string) (ws []string, errors []error) {
	timeout := v.(int)
	if timeout <= 0 || timeout > 3600 {
		errors = append(errors, fmt.Errorf(
			"%q must be between 1 and 3600",
			k))
	}
	return
}

// validateRecordType ensures that the dns record type is valid
func validateRecordType(t string, proxied bool) error {
	switch t {
	case "A", "AAAA", "CNAME":
		return nil
	case "TXT", "SRV", "LOC", "MX", "NS", "SPF", "CAA", "CERT", "DNSKEY", "DS", "NAPTR", "SMIMEA", "SSHFP", "TLSA", "URI":
		if !proxied {
			return nil
		}
	default:
		return fmt.Errorf(
			`Invalid type %q. Valid types are "A", "AAAA", "CNAME", "TXT", "SRV", "LOC", "MX", "NS", "SPF", "CAA", "CERT", "DNSKEY", "DS", "NAPTR", "SMIMEA", "SSHFP", "TLSA" or "URI".`, t)
	}

	return fmt.Errorf("Type %q cannot be proxied", t)
}

// validateRecordName ensures that based on supplied record type, the name content matches
// Currently only validates A and AAAA types
func validateRecordName(t string, value string) error {
	switch t {
	case "A":
		// Must be ipv4 addr
		addr := net.ParseIP(value)
		if addr == nil || !strings.Contains(value, ".") {
			return fmt.Errorf("A record must be a valid IPv4 address, got: %q", value)
		}
	case "AAAA":
		// Must be ipv6 addr
		addr := net.ParseIP(value)
		if addr == nil || !strings.Contains(value, ":") {
			return fmt.Errorf("AAAA record must be a valid IPv6 address, got: %q", value)
		}
	case "TXT":
		// Must be printable ASCII
		for i := 0; i < len(value); i++ {
			char := value[i]
			if (char < 0x20) || (0x7F < char) {
				return fmt.Errorf("TXT record must contain printable ASCII, found: %q", char)
			}
		}
	}

	return nil
}

func validateVLANName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 20 {
		errors = append(errors, fmt.Errorf(
			"Length provided for '%q' is too long. Maximum length is 20 characters", k))
	}
	return
}

func validateAuthProtocol(v interface{}, k string) (ws []string, errors []error) {
	authProtocol := v.(string)
	if authProtocol != "MD5" && authProtocol != "SHA1" && authProtocol != "SHA256" {
		errors = append(errors, fmt.Errorf(
			"%q auth protocol can be MD5 or SHA1 or SHA256", k))
	}
	return
}

//ValidateIPVersion
func validateIPVersion(v interface{}, k string) (ws []string, errors []error) {
	validVersions := map[string]bool{
		"ipv4": true,
		"ipv6": true,
	}

	value := v.(string)
	_, found := validVersions[value]
	if !found {
		strarray := make([]string, 0, len(validVersions))
		for key := range validVersions {
			strarray = append(strarray, key)
		}
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid ip version type %q. Valid types are %q.",
			k, value, strings.Join(strarray, ",")))
	}
	return
}

func validateVPCIdentity(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	// We do not currently accept CRN or HRef
	validators := []func(string) bool{isSecurityGroupAddress, isSecurityGroupCIDR,
		isVPCIdentityByID}

	for _, validator := range validators {
		if validator(value) {
			return
		}
	}
	errors = append(errors, fmt.Errorf("%q (%s) invalid vpc identity", k, value))
	return
}

func validateResourceGroupId(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, err := gouuid.FromString(value)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q contains an invalid resource group id, %q.", k, value))
	}
	return
}

func validateSecurityGroupId(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, err := gouuid.FromString(value)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q contains an invalid security group id, %q.", k, value))
	}
	return
}

func isSecurityGroupAddress(s string) bool {
	return net.ParseIP(s) != nil
}

func isSecurityGroupCIDR(s string) bool {
	_, _, err := net.ParseCIDR(s)
	return err == nil
}

func isSecurityGroupIdentityByID(s string) bool {
	_, err := gouuid.FromString(s)
	return err == nil
}

func isSecurityGroupIdentityByCRN(s string) bool {
	segments := strings.Split(s, ":")
	return len(segments) == 10 && segments[0] == "crn"
}

func isSecurityGroupIdentityByHRef(s string) bool {
	return validHRef.MatchString(s)
}

func isVPCIdentityByID(s string) bool {
	_, err := gouuid.FromString(s)
	return err == nil
}

func validateSecurityGroupRemote(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	validators := []func(string) bool{isSecurityGroupAddress, isSecurityGroupCIDR,
		isSecurityGroupIdentityByID /*, isSecurityGroupIdentityByCRN, isSecurityGroupIdentityByHRef*/}

	for _, validator := range validators {
		if validator(value) {
			return
		}
	}
	errors = append(errors, fmt.Errorf("%q (%s) invalid security group remote", k, value))
	return
}

func validateGeneration(v interface{}, k string) (ws []string, errors []error) {
	validVersions := map[string]bool{
		"gc": true,
		"gt": true,
	}

	value := v.(string)
	_, found := validVersions[value]
	if !found {
		strarray := make([]string, 0, len(validVersions))
		for key := range validVersions {
			strarray = append(strarray, key)
		}
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid generation type %q. Valid types are %q.",
			k, value, strings.Join(strarray, ",")))
	}
	return
}

func validateEncyptionProtocol(v interface{}, k string) (ws []string, errors []error) {
	encyptionProtocol := v.(string)
	if encyptionProtocol != "DES" && encyptionProtocol != "3DES" && encyptionProtocol != "AES128" && encyptionProtocol != "AES192" && encyptionProtocol != "AES256" {
		errors = append(errors, fmt.Errorf(
			"%q encryption protocol can be DES or 3DES or AES128 or AES192 or AES256", k))
	}
	return
}

func validateDeadPeerDetectionInterval(v interface{}, k string) (ws []string, errors []error) {
	secs := v.(int)
	if secs < 15 || secs > 86399 {
		errors = append(errors, fmt.Errorf(
			"%q must be between 15 and 86399",
			k))
		return
	}
	return
}

func validateDiffieHellmanGroup(v interface{}, k string) (ws []string, errors []error) {
	diffieHellmanGroup := v.(int)
	if diffieHellmanGroup != 0 && diffieHellmanGroup != 1 && diffieHellmanGroup != 2 && diffieHellmanGroup != 5 {
		errors = append(errors, fmt.Errorf(
			"%q Diffie Hellman Group can be 0 or 1 or 2 or 5", k))
	}
	return
}

func validateAllowedRangeInt(start, end int) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(int)
		if value < start || value > end {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid int value should be in range(%d, %d), got %d",
				k, start, end, value))
		}
		return
	}
}

func validateDeadPeerDetectionTimeout(v interface{}, k string) (ws []string, errors []error) {
	secs := v.(int)
	if secs < 15 || secs > 86399 {
		errors = append(errors, fmt.Errorf(
			"%q must be between 15 and 86399",
			k))
		return
	}
	return
}

func validatekeylife(v interface{}, k string) (ws []string, errors []error) {
	keylife := v.(int)
	if keylife < 120 || keylife > 172800 {
		errors = append(errors, fmt.Errorf(
			"%q keylife value can be between 120 and 172800", k))
	}
	return
}

func validateLBListenerPort(v interface{}, k string) (ws []string, errors []error) {
	return validatePortRange(1, 65535)(v, k)
}

func validateLBListenerConnectionLimit(v interface{}, k string) (ws []string, errors []error) {
	conns := v.(int)
	if conns < 1 || conns > 15000 {
		errors = append(errors, fmt.Errorf(
			"%q must be between 1 and 15000",
			k))
		return
	}
	return
}

//ValidateISName
func validateISName(v interface{}, k string) (ws []string, errors []error) {
	name := v.(string)
	acceptedcharacters, _ := regexp.MatchString(`^[a-z][-a-z0-9]*$`, name)
	endwithalphanumeric, _ := regexp.MatchString(`.*[a-z0-9]$`, name)
	length := len(name)
	if acceptedcharacters == true {
		if length <= 40 {
			if endwithalphanumeric == true {
				if strings.Contains(name, "--") != true {
					return
				} else {
					errors = append(errors, fmt.Errorf(
						"%q (%q) should not contain consecutive dash(-)", k, v))
				}
			} else {
				errors = append(errors, fmt.Errorf(
					"%q (%q) should not end with dash(-) ", k, v))
			}
		} else {
			errors = append(errors, fmt.Errorf(
				"%q (%q) should not exceed 40 characters", k, v))
		}

	} else {
		errors = append(errors, fmt.Errorf(
			"%q (%q) should contain only lowercase alphanumeric,dash and should begin with lowercase character", k, v))
	}
	return
}

// ValidateFunc is honored only when the schema's Type is set to TypeInt,
// TypeFloat, TypeString, TypeBool, or TypeMap. It is ignored for all other types.
// enum to list all the validator functions supported by this tool.
type FunctionIdentifier int

const (
	IntBetween FunctionIdentifier = iota
	IntAtLeast
	IntAtMost
	ValidateAllowedStringValue
	StringLenBetween
	ValidateIPorCIDR
	ValidateCIDRAddress
	ValidateAllowedIntValue
	ValidateRegexpLen
	ValidateRegexp
	ValidateNoZeroValues
	ValidateJSONString
	ValidateJSONParam
	ValidateBindedPackageName
	ValidateOverlappingAddress
)

// ValueType -- Copied from Terraform for now. You can refer to Terraform ValueType directly.
// ValueType is an enum of the type that can be represented by a schema.
type ValueType int

const (
	TypeInvalid ValueType = iota
	TypeBool
	TypeInt
	TypeFloat
	TypeString
)

// Type of constraints required for validation
type ValueConstraintType int

const (
	MinValue ValueConstraintType = iota
	MaxValue
	MinValueLength
	MaxValueLength
	AllowedValues
	MatchesValue
)

// Schema is used to describe the validation schema.
type ValidateSchema struct {

	//This is the parameter name.
	//Ex: private_subnet in ibm_compute_bare_metal resource
	Identifier string

	// this is similar to schema.ValueType
	Type ValueType

	// The actual validation function that needs to be invoked.
	// Ex: IntBetween, validateAllowedIntValue, validateAllowedStringValue
	ValidateFunctionIdentifier FunctionIdentifier

	MinValue       string
	MaxValue       string
	AllowedValues  string //Comma separated list of strings.
	Matches        string
	Regexp         string
	MinValueLength int
	MaxValueLength int

	// Is this nullable
	Nullable bool

	Optional bool
	Required bool
	Default  interface{}
	ForceNew bool
}

type ResourceValidator struct {
	// This is the resource name - Found in provider.go of IBM Terraform provider.
	// Ex: ibm_compute_monitor, ibm_compute_bare_metal, ibm_compute_dedicated_host, ibm_cis_global_load_balancer etc.,
	ResourceName string

	// Array of validator objects. Each object refers to one parameter in the resource provider.
	Schema []ValidateSchema
}

type ValidatorDict struct {
	ResourceValidatorDictionary   map[string]*ResourceValidator
	DataSourceValidatorDictionary map[string]*ResourceValidator
}

// Resource Validator Dictionary -- For all terraform IBM Resource Providers.
// This is of type - Array of ResourceValidators.
// Each object in this array is a type of map, where key == ResourceName and value == array of ValidateSchema objects. Each of these
// ValidateSchema corresponds to a parameter in the resourceProvider.

var validatorDict = Validator()

// This is the main validation function. This function will be used in all the provider code.
func InvokeValidator(resourceName, identifier string) schema.SchemaValidateFunc {
	// Loop through dictionary and identify the resource and then the parameter configuration.
	var schemaToInvoke ValidateSchema
	found := false
	resourceItem := validatorDict.ResourceValidatorDictionary[resourceName]
	if resourceItem.ResourceName == resourceName {
		parameterValidateSchema := resourceItem.Schema
		for _, validateSchema := range parameterValidateSchema {
			if validateSchema.Identifier == identifier {
				schemaToInvoke = validateSchema
				found = true
				break
			}
		}
	}

	if found {
		return invokeValidatorInternal(schemaToInvoke)
	} else {
		// Add error code later. TODO
		return nil
	}
}

func InvokeDataSourceValidator(resourceName, identifier string) schema.SchemaValidateFunc {
	// Loop through dictionary and identify the resource and then the parameter configuration.
	var schemaToInvoke ValidateSchema
	found := false

	dataSourceItem := validatorDict.DataSourceValidatorDictionary[resourceName]
	if dataSourceItem.ResourceName == resourceName {
		parameterValidateSchema := dataSourceItem.Schema
		for _, validateSchema := range parameterValidateSchema {
			if validateSchema.Identifier == identifier {
				schemaToInvoke = validateSchema
				found = true
				break
			}
		}
	}

	if found {
		return invokeValidatorInternal(schemaToInvoke)
	} else {
		// Add error code later. TODO
		return nil
	}
}

// the function is currently modified to invoke SchemaValidateFunc directly.
// But in terraform, we will just return SchemaValidateFunc as shown below.. So terraform will invoke this func
func invokeValidatorInternal(schema ValidateSchema) schema.SchemaValidateFunc {

	funcIdentifier := schema.ValidateFunctionIdentifier
	switch funcIdentifier {
	case IntBetween:
		minValue := schema.GetValue(MinValue)
		maxValue := schema.GetValue(MaxValue)
		return validation.IntBetween(minValue.(int), maxValue.(int))
	case IntAtLeast:
		minValue := schema.GetValue(MinValue)
		return validation.IntAtLeast(minValue.(int))
	case IntAtMost:
		maxValue := schema.GetValue(MaxValue)
		return validation.IntAtMost(maxValue.(int))
	case ValidateAllowedStringValue:
		allowedValues := schema.GetValue(AllowedValues)
		return validateAllowedStringValue(allowedValues.([]string))
	case StringLenBetween:
		return validation.StringLenBetween(schema.MinValueLength, schema.MaxValueLength)
	case ValidateIPorCIDR:
		return validateIPorCIDR()
	case ValidateCIDRAddress:
		return validateCIDRAddress()
	case ValidateAllowedIntValue:
		allowedValues := schema.GetValue(AllowedValues)
		return validateAllowedIntValue(allowedValues.([]int))
	case ValidateRegexpLen:
		return validateRegexpLen(schema.MinValueLength, schema.MaxValueLength, schema.Regexp)
	case ValidateRegexp:
		return validateRegexp(schema.Regexp)
	case ValidateNoZeroValues:
		return validateNoZeroValues()
	case ValidateJSONString:
		return validateJSONString()
	case ValidateBindedPackageName:
		return validateBindedPackageName()
	case ValidateOverlappingAddress:
		return validateOverlappingAddress()

	default:
		return nil
	}
}

// utility functions - Move to different package
func (vs ValidateSchema) GetValue(valueConstraint ValueConstraintType) interface{} {

	var valueToConvert string
	switch valueConstraint {
	case MinValue:
		valueToConvert = vs.MinValue
	case MaxValue:
		valueToConvert = vs.MaxValue
	case AllowedValues:
		valueToConvert = vs.AllowedValues
	case MatchesValue:
		valueToConvert = vs.Matches
	}

	switch vs.Type {
	case TypeInvalid:
		return nil
	case TypeBool:
		b, err := strconv.ParseBool(valueToConvert)
		if err != nil {
			return vs.Zero()
		}
		return b
	case TypeInt:
		// Convert comma separated string to array
		if strings.Contains(valueToConvert, ",") {
			var arr2 []int
			arr1 := strings.Split(valueToConvert, ",")
			for _, ele := range arr1 {
				e, err := strconv.Atoi(strings.TrimSpace(ele))
				if err != nil {
					return vs.Zero()
				}
				arr2 = append(arr2, e)
			}
			return arr2
		} else {
			num, err := strconv.Atoi(valueToConvert)
			if err != nil {
				return vs.Zero()
			}
			return num
		}

	case TypeFloat:
		f, err := strconv.ParseFloat(valueToConvert, 32)
		if err != nil {
			return vs.Zero()
		}
		return f
	case TypeString:
		//return valueToConvert
		// Convert comma separated string to array
		arr := strings.Split(valueToConvert, ",")
		for i, ele := range arr {
			arr[i] = strings.TrimSpace(ele)
		}
		return arr
	default:
		panic(fmt.Sprintf("unknown type %s", vs.Type))
	}
}

// Use stringer tool to generate this later.
func (i FunctionIdentifier) String() string {
	return [...]string{"IntBetween", "IntAtLeast", "IntAtMost"}[i]
}

// Use Stringer tool to generate this later.
func (i ValueType) String() string {
	return [...]string{"TypeInvalid", "TypeBool", "TypeInt", "TypeFloat", "TypeString"}[i]
}

// Use Stringer tool to generate this later.
func (i ValueConstraintType) String() string {
	return [...]string{"MinValue", "MaxValue", "MinValueLength", "MaxValueLength", "AllowedValues", "MatchesValue"}[i]
}

// Zero returns the zero value for a type.
func (vs ValidateSchema) Zero() interface{} {
	switch vs.Type {
	case TypeInvalid:
		return nil
	case TypeBool:
		return false
	case TypeInt:
		return make([]string, 0)
	case TypeFloat:
		return 0.0
	case TypeString:
		return make([]int, 0)
	default:
		panic(fmt.Sprintf("unknown type %s", vs.Type))
	}
}
