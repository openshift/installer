/*
Copyright (c) 2020 Red Hat, Inc.

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

// This file contains the types and functions used to manage the configuration of the command line
// client.

package ocm

import (
	"fmt"
	"net/url"
	"strings"

	sdk "github.com/openshift-online/ocm-sdk-go"

	"github.com/openshift/rosa/pkg/config"
	"github.com/openshift/rosa/pkg/fedramp"
)

const Production = "production"

// URLAliases allows the value of the `--env` option to map to the various API URLs.
var URLAliases = map[string]string{
	"production":  "https://api.openshift.com",
	"staging":     "https://api.stage.openshift.com",
	"integration": "https://api.integration.openshift.com",
	"local":       "http://localhost:8000",
	"local-proxy": "http://localhost:9000",
	"crc":         "https://clusters-service.apps-crc.testing",
}

func ValidOCMUrlAliases() []string {
	keys := make([]string, 0, len(URLAliases))
	for k := range URLAliases {
		keys = append(keys, k)
	}
	return keys
}

// URL Precedent (from highest priority to lowest priority):
//  1. runtime `--env` cli arg (key found in `urlAliases`)
//  2. runtime `--env` cli arg (non-empty string)
//  3. config file `URL` value (non-empty string)
//  4. sdk.DefaultURL
//
// Finally, it will try to url.ParseRequestURI the resolved URL to make sure it's a valid URL.
func ResolveGatewayUrl(optionalParsedCliFlagValue string, optionalParsedConfig *config.Config) (string, error) {
	gatewayURL := sdk.DefaultURL
	source := "default"
	if optionalParsedCliFlagValue != "" {
		gatewayURL = optionalParsedCliFlagValue
		source = "flag"
		if _, ok := URLAliases[optionalParsedCliFlagValue]; ok {
			gatewayURL = URLAliases[optionalParsedCliFlagValue]
		}
	} else if optionalParsedConfig != nil && optionalParsedConfig.URL != "" {
		// re-use the URL from the config file
		gatewayURL = optionalParsedConfig.URL
		source = "config"
	}

	url, err := url.ParseRequestURI(gatewayURL)
	if err != nil {
		return "", fmt.Errorf(
			"%w\n\nURL Source: %s\nExpected an absolute URI/path (e.g. %s) or a case-sensitive alias, one of: [%s]",
			err, source, sdk.DefaultURL, strings.Join(ValidOCMUrlAliases(), ", "))
	}

	return url.String(), nil
}

func GetEnv() (string, error) {
	cfg, err := config.Load()
	if err != nil {
		return "", err
	}

	urlAliases := URLAliases
	if cfg.FedRAMP {
		urlAliases = fedramp.URLAliases
	}

	// Check for OCM environments (including regionalized URLs)
	if strings.HasSuffix(strings.TrimSuffix(cfg.URL, "/"), "openshift.com") {
		regionDiscoveryUrl, err := sdk.DetermineRegionDiscoveryUrl(cfg.URL)
		if err == nil {
			discoveryGatewayUrl, _ := url.Parse(regionDiscoveryUrl)
			// Check for URL aliases
			for env, api := range urlAliases {
				if api == fmt.Sprintf("%s://%s", discoveryGatewayUrl.Scheme, discoveryGatewayUrl.Host) {
					return env, nil
				}
			}
		}
	}

	// URL check as a fallback mechanism (in case of other URLs like local envs, fedRAMP envs, etc.)
	for env, api := range urlAliases {
		if api == strings.TrimSuffix(cfg.URL, "/") {
			return env, nil
		}
	}

	// Special use case for Admin users in the GovCloud environment
	for env, api := range fedramp.AdminURLAliases {
		if api == strings.TrimSuffix(cfg.URL, "/") {
			return env, nil
		}
	}

	return "", fmt.Errorf("Invalid OCM API")
}
