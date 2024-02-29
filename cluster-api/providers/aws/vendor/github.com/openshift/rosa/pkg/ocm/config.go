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
	"strings"

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

func GetEnv() (string, error) {
	cfg, err := config.Load()
	if err != nil {
		return "", err
	}

	urlAliases := URLAliases
	if cfg.FedRAMP {
		urlAliases = fedramp.URLAliases
	}

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
