// Copyright (C) 2018 Joey Ma <majunjiev@gmail.com>
// All rights reserved.
//
// This software may be modified and distributed under the terms
// of the BSD-2 license.  See the LICENSE file for details.

package ovirt

import (
	"bytes"
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func extractSemanticVerion(version string) (major, minor, patch string) {
	vs := strings.Split(version, ".")
	switch len(vs) {
	case 1:
		return vs[0], "", ""
	case 2:
		return vs[0], vs[1], ""
	case 3:
		return vs[0], vs[1], vs[2]
	default:
		return "", "", ""
	}
}

// macRange returns a SchemaValidateFunc which tests if the provided value
// is of type string, and in valid MAC range notation
func macRange() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}
		macs := strings.Split(v, ",")
		if len(macs) != 2 {
			es = append(es, fmt.Errorf(
				"expected %s to contain a valid MAC range, got: %s", k, v))
			return
		}

		mac1, err1 := net.ParseMAC(macs[0])
		mac2, err2 := net.ParseMAC(macs[1])
		if err1 != nil || err2 != nil || bytes.Compare(mac1, mac2) > 0 {
			es = append(es, fmt.Errorf(
				"expected %s to contain a valid MAC range, got: %s", k, v))
		}
		return

	}
}

func parseResourceID(id string, count int) ([]string, error) {
	parts := strings.Split(id, ":")

	if len(parts) != count {
		return nil, fmt.Errorf("Invalid Resource ID %s, expected %d parts, got %d", id, count, len(parts))
	}
	return parts, nil
}

//Converts an array of type []interface{} to []string, notice that all the interface elements needs to be strings
func convInterfaceArrToStringArr(arr []interface{}) ([]string, error) {
	var newArr []string
	for _, val := range arr {
		if s, ok := val.(string); ok {
			newArr = append(newArr, s)
		} else {
			return nil, fmt.Errorf(
				"error converting []interface{} to []string, "+
					"provided interface array %v contains non string elements %v",
				arr, val)
		}
	}
	return newArr, nil
}
