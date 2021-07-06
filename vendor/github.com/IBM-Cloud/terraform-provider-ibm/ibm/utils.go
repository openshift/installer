// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Used for retry logic on resource timeout.
func isResourceTimeoutError(err error) bool {
	timeoutErr, ok := err.(*resource.TimeoutError)
	return ok && timeoutErr.LastError == nil
}
func GetPrivateServiceURLForRegion(region string) (string, error) {
	var endpoints = map[string]string{
		"us-south":   "https://private.us.icr.io",  // us-south
		"uk-south":   "https://private.uk.icr.io",  // uk-south
		"eu-gb":      "https://private.uk.icr.io",  // eu-gb
		"eu-central": "https://private.de.icr.io",  // eu-central
		"eu-de":      "https://private.de.icr.io",  // eu-de
		"ap-north":   "https://private.jp.icr.io",  // ap-north
		"jp-tok":     "https://private.jp.icr.io",  // jp-tok
		"ap-south":   "https://private.au.icr.io",  // ap-south
		"au-syd":     "https://private.au.icr.io",  // au-syd
		"global":     "https://private.icr.io",     // global
		"jp-osa":     "https://private.jp2.icr.io", // jp-osa
	}

	if url, ok := endpoints[region]; ok {
		return url, nil
	}
	return "", fmt.Errorf("service URL for region '%s' not found", region)
}
