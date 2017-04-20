// Copyright 2017 CoreOS, Inc.
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

package templating

import (
	"fmt"
	"strings"
)

var (
	ErrUnknownPlatform = fmt.Errorf("unsupported platform")
	ErrUnknownField    = fmt.Errorf("unknown field")
)

const (
	PlatformAzure  = "azure"
	PlatformDO     = "digitalocean"
	PlatformEC2    = "ec2"
	PlatformGCE    = "gce"
	PlatformPacket = "packet"
)

var Platforms = []string{
	PlatformAzure,
	PlatformDO,
	PlatformEC2,
	PlatformGCE,
	PlatformPacket,
}

const (
	fieldHostname  = "HOSTNAME"
	fieldV4Private = "PRIVATE_IPV4"
	fieldV4Public  = "PUBLIC_IPV4"
	fieldV6Private = "PRIVATE_IPV6"
	fieldV6Public  = "PUBLIC_IPV6"
)

var platformTemplatingMap = map[string]map[string]string{
	PlatformAzure: {
		// TODO: is this right?
		fieldV4Private: "COREOS_AZURE_IPV4_DYNAMIC",
		fieldV4Public:  "COREOS_AZURE_IPV4_VIRTUAL",
	},
	PlatformDO: {
		// TODO: unused: COREOS_DIGITALOCEAN_IPV4_ANCHOR_0
		fieldHostname:  "COREOS_DIGITALOCEAN_HOSTNAME",
		fieldV4Private: "COREOS_DIGITALOCEAN_IPV4_PRIVATE_0",
		fieldV4Public:  "COREOS_DIGITALOCEAN_IPV4_PUBLIC_0",
		fieldV6Private: "COREOS_DIGITALOCEAN_IPV6_PRIVATE_0",
		fieldV6Public:  "COREOS_DIGITALOCEAN_IPV6_PUBLIC_0",
	},
	PlatformEC2: {
		fieldHostname:  "COREOS_EC2_HOSTNAME",
		fieldV4Private: "COREOS_EC2_IPV4_LOCAL",
		fieldV4Public:  "COREOS_EC2_IPV4_PUBLIC",
	},
	PlatformGCE: {
		fieldHostname:  "COREOS_GCE_HOSTNAME",
		fieldV4Private: "COREOS_GCE_IP_EXTERNAL_0",
		fieldV4Public:  "COREOS_GCE_IP_LOCAL_0",
	},
	PlatformPacket: {
		fieldHostname:  "COREOS_PACKET_HOSTNAME",
		fieldV4Private: "COREOS_PACKET_IPV4_PRIVATE_0",
		fieldV4Public:  "COREOS_PACKET_IPV4_PUBLIC_0",
		fieldV6Public:  "COREOS_PACKET_IPV6_PUBLIC_0",
	},
}

// HasTemplating returns whether or not any of the environment variables present
// in the passed in list use ct templating
func HasTemplating(vars []string) bool {
	for _, v := range vars {
		if strings.ContainsRune(v, '{') || strings.ContainsRune(v, '}') {
			return true
		}
	}
	return false
}

func PerformTemplating(platform string, vars []string) ([]string, error) {
	if _, ok := platformTemplatingMap[platform]; !ok {
		return nil, ErrUnknownPlatform
	}

	for i := range vars {
		startIndex := strings.IndexRune(vars[i], '{')
		endIndex := strings.IndexRune(vars[i], '}')
		for startIndex != -1 && endIndex != -1 && startIndex < endIndex {
			fieldName := vars[i][startIndex+1 : endIndex]
			fieldVal, ok := platformTemplatingMap[platform][fieldName]
			if !ok {
				return nil, ErrUnknownField
			}
			vars[i] = strings.Replace(vars[i], "{"+fieldName+"}", "${"+fieldVal+"}", 1)

			// start the search for a new start index from the old end index, or
			// we'll just find the curly braces we just substituted in
			startIndex = strings.IndexRune(vars[i][endIndex:], '{')
			if startIndex != -1 {
				startIndex += endIndex

				// and start the search for a new end index from the new start
				// index, or as before we'll just find the curly braces we just
				// substituted in
				endIndex = strings.IndexRune(vars[i][startIndex:], '}')
				if endIndex != -1 {
					endIndex += startIndex
				}
			}

		}
	}
	return vars, nil
}
