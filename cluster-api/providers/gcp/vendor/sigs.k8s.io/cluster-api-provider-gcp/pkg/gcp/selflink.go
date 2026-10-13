/*
Copyright 2025 The Kubernetes Authors.

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

// Package gcp implements helper functions for working with GCP.
package gcp

import "github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/meta"

// FormatKey builds a selfLink style string from a meta.Key for logging / human-facing error messages
func FormatKey(resourceType string, key *meta.Key) string {
	return SelfLink(resourceType, key)
}

// SelfLink builds a selfLink for passing to GCP APIs
func SelfLink(resourceType string, key *meta.Key) string {
	if key.Region != "" {
		return "regions/" + key.Region + "/" + resourceType + "/" + key.Name
	}
	if key.Zone != "" {
		return "zones/" + key.Zone + "/" + resourceType + "/" + key.Name
	}
	return "global/" + resourceType + "/" + key.Name
}
