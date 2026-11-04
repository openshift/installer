/*
Copyright 2026 The Kubernetes Authors.

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

package firewalls

import (
	"slices"
	"strings"

	"google.golang.org/api/compute/v1"
)

// allIPv4Range is the range GCP fills a source or destination in with when a rule leaves
// it out, since an omitted range means every address.
const allIPv4Range = "0.0.0.0/0"

// driftedFields returns the names of the user-configurable fields that differ between the
// desired spec and the firewall rule that currently exists in GCP.
//
// Server-populated fields (SelfLink, Id, CreationTimestamp, ...) are not compared, and neither
// is the network, which cannot be changed without recreating the rule.
func driftedFields(spec, actual *compute.Firewall) []string {
	var drifted []string

	if spec.Description != actual.Description {
		drifted = append(drifted, "description")
	}
	// GCP cannot update the direction of an existing rule, and the webhooks reject a spec
	// that asks for it, so the only rule that drifts here is one that already existed under
	// a name the spec now claims. It is still reported, because the update GCP refuses is
	// what tells the user that the rule they adopted is not the rule they described.
	if !strings.EqualFold(spec.Direction, actual.Direction) {
		drifted = append(drifted, "direction")
	}
	if spec.Disabled != actual.Disabled {
		drifted = append(drifted, "disabled")
	}
	// The default rules do not set a priority, leaving GCP to assign one, so only compare the
	// priority when the spec carries one. Otherwise every default rule looks drifted. A user
	// specified rule always has one, since the field has a minimum of 1 and defaults to 1000.
	if spec.Priority != 0 && spec.Priority != actual.Priority {
		drifted = append(drifted, "priority")
	}
	if !equalSourceRanges(spec, actual) {
		drifted = append(drifted, "sourceRanges")
	}
	// A rule applies to every destination unless it names one, so GCP always fills the
	// destination of a rule that leaves it out.
	if !equalRanges(spec.DestinationRanges, actual.DestinationRanges) {
		drifted = append(drifted, "destinationRanges")
	}
	if !equalUnordered(spec.SourceTags, actual.SourceTags) {
		drifted = append(drifted, "sourceTags")
	}
	if !equalUnordered(spec.TargetTags, actual.TargetTags) {
		drifted = append(drifted, "targetTags")
	}
	if !equalUnordered(allowedDescriptors(spec.Allowed), allowedDescriptors(actual.Allowed)) {
		drifted = append(drifted, "allowed")
	}
	if !equalUnordered(deniedDescriptors(spec.Denied), deniedDescriptors(actual.Denied)) {
		drifted = append(drifted, "denied")
	}

	return drifted
}

func allowedDescriptors(allowed []*compute.FirewallAllowed) []string {
	descriptors := make([]string, 0, len(allowed))
	for _, a := range allowed {
		descriptors = append(descriptors, describeProtocolPorts(a.IPProtocol, a.Ports))
	}

	return descriptors
}

func deniedDescriptors(denied []*compute.FirewallDenied) []string {
	descriptors := make([]string, 0, len(denied))
	for _, d := range denied {
		descriptors = append(descriptors, describeProtocolPorts(d.IPProtocol, d.Ports))
	}

	return descriptors
}

// describeProtocolPorts flattens a protocol and its ports into a single comparable string so
// that allow/deny lists can be compared without caring about ordering.
func describeProtocolPorts(protocol string, ports []string) string {
	sorted := slices.Clone(ports)
	slices.Sort(sorted)

	return strings.ToLower(protocol) + ":" + strings.Join(sorted, ",")
}

// equalSourceRanges reports whether the source ranges the spec asks for match the ones the rule
// carries in GCP.
//
// GCP only fills the source in when the rule names no source at all. A rule that selects its
// source by tag or by service account keeps an empty list, so an empty list means "none of them"
// rather than "all of them", and a range that appears on such a rule widens it to the whole
// internet. That is a change made out of band, not a default, and reverting it is the point of
// comparing at all.
func equalSourceRanges(spec, actual *compute.Firewall) bool {
	if len(spec.SourceTags) > 0 || len(spec.SourceServiceAccounts) > 0 {
		return equalUnordered(spec.SourceRanges, actual.SourceRanges)
	}

	return equalRanges(spec.SourceRanges, actual.SourceRanges)
}

// equalRanges reports whether a source or destination range the spec asks for matches the one
// the rule carries in GCP.
//
// A rule that leaves the field out applies to every address, and GCP stores that by filling the
// field in with the range that matches everything rather than leaving it empty. Comparing the
// two literally would report an egress rule without a destination, or an ingress rule without a
// source, as drifted on every reconcile and update it in a loop. An empty list therefore matches
// the range that matches everything, which is what the empty list already means.
func equalRanges(spec, actual []string) bool {
	if len(spec) > 0 {
		return equalUnordered(spec, actual)
	}

	return len(actual) == 0 || slices.Equal(actual, []string{allIPv4Range})
}

// equalUnordered reports whether two string slices hold the same elements, ignoring order. GCP
// does not guarantee that it returns list fields in the order they were submitted.
func equalUnordered(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	sortedA, sortedB := slices.Clone(a), slices.Clone(b)
	slices.Sort(sortedA)
	slices.Sort(sortedB)

	return slices.Equal(sortedA, sortedB)
}
