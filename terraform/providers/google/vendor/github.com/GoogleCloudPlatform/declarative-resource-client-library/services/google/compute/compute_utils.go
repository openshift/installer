// Copyright 2024 Google LLC. All Rights Reserved.
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
// Package compute contains handwritten support code for the compute service.
package compute

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/url"
	"path"
	"strings"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl/operations"
)

// EncodeImageDeprecateRequest properly encodes the image deprecation request for an image.
func EncodeImageDeprecateRequest(m map[string]interface{}) map[string]interface{} {
	req := make(map[string]interface{})
	// Deprecate requests involve removing the "deprecated" key.
	if deprecatedVal, ok := m["deprecated"]; ok {
		deprecatedMap := deprecatedVal.(map[string]interface{})
		for key, value := range deprecatedMap {
			req[key] = value
		}
	}

	return req
}

// WrapTargetPoolInstance wraps the instances provided by AddInstance and RemoveInstance
// in their required format.
func WrapTargetPoolInstance(m map[string]interface{}) map[string]interface{} {
	is, ok := m["instances"].([]string)
	if !ok {
		return nil
	}
	wrapped := make([]interface{}, len(is))
	for idx, i := range is {
		wrapped[idx] = map[string]interface{}{
			"instance": i,
		}
	}
	return map[string]interface{}{
		"instances": wrapped,
	}
}

// WrapTargetPoolHealthCheck wraps the instances provided by AddHC and RemoveHC
// in their required format.
func WrapTargetPoolHealthCheck(m map[string]interface{}) map[string]interface{} {
	hcs, ok := m["healthChecks"].([]string)
	if !ok {
		return nil
	}
	wrapped := make([]interface{}, len(hcs))
	for idx, hc := range hcs {
		wrapped[idx] = map[string]interface{}{
			"healthCheck": hc,
		}
	}
	return map[string]interface{}{
		"healthChecks": wrapped,
	}
}

// forwardingRuleEncodeCreateRequest removes the labels parameter - it cannot be supplied on create.
func forwardingRuleEncodeCreateRequest(m map[string]any) map[string]any {
	// labels cannot be specified on create
	delete(m, "labels")
	return m
}

// forwardingRuleSetLabelsPostCreate adds a 'setLabels' operation after
// a create operation, because creation cannot set labels due to a
// long-standing bug in the API for most compute networking resources.
// createPubsubConfigs adds a patch to apply PubsubConfigs after create (if applicable).
func forwardingRuleSetLabelsPostCreate(inOps []forwardingRuleApiOperation) ([]forwardingRuleApiOperation, error) {
	for _, op := range inOps {
		if _, ok := op.(*createForwardingRuleOperation); ok {
			return append(inOps, &updateForwardingRuleSetLabelsOperation{FieldDiffs: []*dcl.FieldDiff{{FieldName: "labels"}}}), nil
		}
	}
	return inOps, nil
}

func canonicalizeReservationCPUPlatform(o, n interface{}) bool {
	oVal, _ := o.(*string)
	nVal, _ := n.(*string)
	return equalReservationCPUPlatform(oVal, nVal)
}

func equalReservationCPUPlatform(o, n *string) bool {
	if o == nil && n == nil {
		return true
	}
	if o == nil || n == nil {
		return false
	}
	if *o == "automatic" {
		return true
	}
	if *n == "automatic" {
		return true
	}

	return *o == *n
}

func canonicalizeIPAddressToReference(o, n interface{}) bool {
	oVal, _ := o.(*string)
	nVal, _ := n.(*string)
	if oVal == nil && nVal == nil {
		return true
	}
	if oVal == nil || nVal == nil {
		return false
	}
	if isIPV4Address(*oVal) && !isIPV4Address(*nVal) {
		return true
	}
	if isIPV4Address(*nVal) && !isIPV4Address(*oVal) {
		return true
	}
	return dcl.NameToSelfLink(oVal, nVal)
}

func isIPV4Address(addr string) bool {
	return net.ParseIP(addr) != nil
}

func canonicalizePortRange(o, n interface{}) bool {
	oVal, _ := o.(*string)
	nVal, _ := n.(*string)
	return equalPortRanges(oVal, nVal)
}

func equalPortRanges(o, n *string) bool {
	if o == nil && n == nil {
		return true
	}
	if o == nil || n == nil {
		return false
	}
	if strings.Contains(*o, "-") && !strings.Contains(*n, "-") {
		// If one of them contains a dash but not the other, ensure that the one with a dash is 'n'.
		swap := n
		n = o
		o = swap
	}
	if !strings.Contains(*o, "-") && strings.Contains(*n, "-") {
		o = dcl.String(fmt.Sprintf("%s-%s", *o, *o))
	}
	return *o == *n
}

// Custom create method for firewall policy which waits on a ComputeGlobalOrganizationOperation.
func (op *createFirewallPolicyOperation) do(ctx context.Context, r *FirewallPolicy, c *Client) error {
	c.Config.Logger.Infof("Attempting to create %v", r)

	u, err := r.createURL(c.Config.BasePath)

	if err != nil {
		return err
	}

	req, err := r.marshal(c)
	if err != nil {
		return err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "POST", u, bytes.NewBuffer(req), c.Config.RetryProvider)
	if err != nil {
		return err
	}
	// Wait for object to be created.
	var o operations.ComputeGlobalOrganizationOperation
	if err := dcl.ParseResponse(resp.Response, &o.BaseOperation); err != nil {
		return err
	}
	if err := o.Wait(context.WithValue(ctx, dcl.DoNotLogRequestsKey, true), c.Config, r.Parent); err != nil {
		c.Config.Logger.Warningf("Creation failed after waiting for operation: %v", err)
		return err
	}
	c.Config.Logger.Infof("Successfully waited for operation")

	r.Name = &o.BaseOperation.TargetID

	if _, err := c.GetFirewallPolicy(ctx, r); err != nil {
		return err
	}

	return nil
}

// Custom update method for network which updates mtu field before updating other fields.
func (op *updateNetworkUpdateOperation) do(ctx context.Context, r *Network, c *Client) error {
	_, err := c.GetNetwork(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "update")
	if err != nil {
		return err
	}

	req, err := newUpdateNetworkUpdateRequest(ctx, r, c)
	if err != nil {
		return err
	}

	if mtu, ok := req["mtu"]; ok {
		// Update mtu field first.
		if err := performNetworkUpdate(ctx, r, c, u, map[string]interface{}{"mtu": mtu}); err != nil {
			return err
		}
		delete(req, "mtu")
	}

	if err := performNetworkUpdate(ctx, r, c, u, req); err != nil {
		return err
	}

	return nil
}

// Send the given update request to the given url on the given network with the given client in the given context and wait for the resulting operation.
func performNetworkUpdate(ctx context.Context, r *Network, c *Client, u string, req map[string]interface{}) error {
	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateNetworkUpdateRequest(c, req)
	if err != nil {
		return err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "PATCH", u, bytes.NewBuffer(body), c.Config.RetryProvider)
	if err != nil {
		return err
	}

	var o operations.ComputeOperation
	if err := dcl.ParseResponse(resp.Response, &o); err != nil {
		return err
	}
	err = o.Wait(context.WithValue(ctx, dcl.DoNotLogRequestsKey, true), c.Config, r.basePath(), "GET")

	if err != nil {
		return err
	}
	return nil
}

func expandFirewallPolicyRuleTLSInspect(_ *Client, tlsInspect *bool, res *FirewallPolicyRule) (*bool, error) {
	if *res.Action == "apply_security_profile_group" || dcl.ValueOrEmptyBool(tlsInspect) {
		return tlsInspect, nil
	}
	return nil, nil
}

func expandNetworkFirewallPolicyRuleTLSInspect(_ *Client, tlsInspect *bool, res *NetworkFirewallPolicyRule) (*bool, error) {
	if *res.Action == "apply_security_profile_group" || dcl.ValueOrEmptyBool(tlsInspect) {
		return tlsInspect, nil
	}
	return nil, nil
}

// Because the server will return both versions and instance template and expects only one to
// be set in our requests, instance template will flatten to nil if versions is non-empty.
func flattenInstanceGroupManagerInstanceTemplateWithConflict(c *Client, instanceTemplate interface{}, resource *InstanceGroupManager) *string {
	if len(resource.Versions) > 0 {
		c.Config.Logger.Info("flattening instance_template field to nil because versions were present")
		return nil
	}
	return dcl.FlattenString(instanceTemplate)
}

// Because the server will return both instance_template and instance template and expects only one to
// be set in our requests, instance template will flatten to nil if instance_template is non-nil.
func flattenInstanceGroupManagerVersionsWithConflict(c *Client, Versions interface{}, resource *InstanceGroupManager) []InstanceGroupManagerVersions {
	if resource.InstanceTemplate != nil {
		c.Config.Logger.Info("flattening versions field to nil because instance_template was present")
		return nil
	}
	return flattenInstanceGroupManagerVersionsSlice(c, Versions, resource)
}

func machineTypeOperations() func(fd *dcl.FieldDiff) []string {
	return func(fd *dcl.FieldDiff) []string {
		// We're assuming that the instance is currently running. If it isn't, this will lead to a no-op stop operation.
		return []string{"updateInstanceStopOperation", "updateInstanceMachineTypeOperation", "updateInstanceStartOperation"}
	}
}

func flattenPacketMirroringRegion(_ *Client, region interface{}) *string {
	regionString, ok := region.(string)
	if !ok {
		return nil
	}
	return dcl.SelfLinkToName(&regionString)
}

func targetPoolHealthCheck() func(fd *dcl.FieldDiff) []string {
	return func(fd *dcl.FieldDiff) []string {
		var ops []string
		if !dcl.IsZeroValue(fd.ToAdd) {
			ops = append(ops, "updateTargetPoolAddHCOperation")
		}
		if !dcl.IsZeroValue(fd.ToRemove) {
			ops = append(ops, "updateTargetPoolRemoveHCOperation")
		}
		return ops
	}
}

func targetPoolInstances() func(fd *dcl.FieldDiff) []string {
	return func(fd *dcl.FieldDiff) []string {
		var ops []string
		if !dcl.IsZeroValue(fd.ToAdd) {
			ops = append(ops, "updateTargetPoolAddInstanceOperation")
		}
		if !dcl.IsZeroValue(fd.ToRemove) {
			ops = append(ops, "updateTargetPoolRemoveInstanceOperation")
		}
		return ops
	}
}

func flattenNetworkSelfLinkWithID(_ *Client, _ interface{}, _ *Network, r map[string]interface{}) *string {
	selfLink, ok := r["selfLink"].(string)
	if !ok {
		return nil
	}
	id, ok := r["id"].(string)
	if !ok {
		return nil
	}
	u, err := url.Parse(selfLink)
	if err != nil {
		return nil
	}
	u.Path = fmt.Sprintf("%s/%s", path.Dir(u.Path), id)
	selfLinkWithID := u.String()
	return &selfLinkWithID
}

// Subnetwork's update operation has a custom method because a separate request must be performed for each field.
func (op *updateSubnetworkUpdateOperation) do(ctx context.Context, r *Subnetwork, c *Client) error {
	_, err := c.GetSubnetwork(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "update")
	if err != nil {
		return err
	}

	req, err := newUpdateSubnetworkUpdateRequest(ctx, r, c)
	if err != nil {
		return err
	}

	fingerprint := req["fingerprint"]
	for field, value := range req {
		if field == "fingerprint" {
			continue
		}
		sr := map[string]interface{}{
			field:         value,
			"fingerprint": fingerprint,
		}
		c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", sr)
		body, err := marshalUpdateSubnetworkUpdateRequest(c, sr)
		if err != nil {
			return err
		}
		resp, err := dcl.SendRequest(ctx, c.Config, "PATCH", u, bytes.NewBuffer(body), c.Config.RetryProvider)
		if err != nil {
			return err
		}

		var o operations.ComputeOperation
		if err := dcl.ParseResponse(resp.Response, &o); err != nil {
			return err
		}
		err = o.Wait(context.WithValue(ctx, dcl.DoNotLogRequestsKey, true), c.Config, r.basePath(), "GET")

		if err != nil {
			return err
		}
		// Perform a get request to pick up the new fingerprint for the resource.
		ur, err := c.GetSubnetwork(ctx, r)
		if err != nil {
			return err
		}
		fingerprint = ur.Fingerprint
	}

	return nil
}
