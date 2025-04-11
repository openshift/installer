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

// IMPORTANT: This file has been generated automatically, refrain from modifying it manually as all
// your changes will be lost when the file is generated again.

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// ClusterClient is the client of the 'cluster' resource.
//
// Manages a specific cluster.
type ClusterClient struct {
	transport http.RoundTripper
	path      string
}

// NewClusterClient creates a new client for the 'cluster'
// resource using the given transport to send the requests and receive the
// responses.
func NewClusterClient(transport http.RoundTripper, path string) *ClusterClient {
	return &ClusterClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
//
// Deletes the cluster.
func (c *ClusterClient) Delete() *ClusterDeleteRequest {
	return &ClusterDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the cluster.
func (c *ClusterClient) Get() *ClusterGetRequest {
	return &ClusterGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Hibernate creates a request for the 'hibernate' method.
//
// Initiates cluster hibernation. While hibernating a cluster will not consume any cloud provider infrastructure
// but will be counted for quota.
func (c *ClusterClient) Hibernate() *ClusterHibernateRequest {
	return &ClusterHibernateRequest{
		transport: c.transport,
		path:      path.Join(c.path, "hibernate"),
	}
}

// Resume creates a request for the 'resume' method.
//
// Resumes from Hibernation.
func (c *ClusterClient) Resume() *ClusterResumeRequest {
	return &ClusterResumeRequest{
		transport: c.transport,
		path:      path.Join(c.path, "resume"),
	}
}

// Update creates a request for the 'update' method.
//
// Updates the cluster.
func (c *ClusterClient) Update() *ClusterUpdateRequest {
	return &ClusterUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// AWS returns the target 'AWS' resource.
func (c *ClusterClient) AWS() *AWSClient {
	return NewAWSClient(
		c.transport,
		path.Join(c.path, "aws"),
	)
}

// AWSInfrastructureAccessRoleGrants returns the target 'AWS_infrastructure_access_role_grants' resource.
//
// Reference to the resource that manages the collection of AWS infrastructure
// access role grants on this cluster.
func (c *ClusterClient) AWSInfrastructureAccessRoleGrants() *AWSInfrastructureAccessRoleGrantsClient {
	return NewAWSInfrastructureAccessRoleGrantsClient(
		c.transport,
		path.Join(c.path, "aws_infrastructure_access_role_grants"),
	)
}

// STSOperatorRoles returns the target 'operator_IAM_roles' resource.
func (c *ClusterClient) STSOperatorRoles() *OperatorIAMRolesClient {
	return NewOperatorIAMRolesClient(
		c.transport,
		path.Join(c.path, "sts_operator_roles"),
	)
}

// AddonInquiries returns the target 'addon_inquiries' resource.
//
// Reference to the resource that manages the collection of the add-on inquiries on this cluster.
func (c *ClusterClient) AddonInquiries() *AddonInquiriesClient {
	return NewAddonInquiriesClient(
		c.transport,
		path.Join(c.path, "addon_inquiries"),
	)
}

// AddonUpgradePolicies returns the target 'addon_upgrade_policies' resource.
//
// Reference to the resource that manages the collection of addon upgrade policies defined for this cluster.
func (c *ClusterClient) AddonUpgradePolicies() *AddonUpgradePoliciesClient {
	return NewAddonUpgradePoliciesClient(
		c.transport,
		path.Join(c.path, "addon_upgrade_policies"),
	)
}

// Addons returns the target 'add_on_installations' resource.
//
// Reference to the resource that manages the collection of add-ons installed on this cluster.
func (c *ClusterClient) Addons() *AddOnInstallationsClient {
	return NewAddOnInstallationsClient(
		c.transport,
		path.Join(c.path, "addons"),
	)
}

// Autoscaler returns the target 'autoscaler' resource.
//
// Reference to the resource that manages the cluster autoscaler.
func (c *ClusterClient) Autoscaler() *AutoscalerClient {
	return NewAutoscalerClient(
		c.transport,
		path.Join(c.path, "autoscaler"),
	)
}

// BreakGlassCredentials returns the target 'break_glass_credentials' resource.
//
// Reference to the resource that manages the break glass credentials.
func (c *ClusterClient) BreakGlassCredentials() *BreakGlassCredentialsClient {
	return NewBreakGlassCredentialsClient(
		c.transport,
		path.Join(c.path, "break_glass_credentials"),
	)
}

// Clusterdeployment returns the target 'clusterdeployment' resource.
//
// Reference to the resource that manages the cluster deployment.
func (c *ClusterClient) Clusterdeployment() *ClusterdeploymentClient {
	return NewClusterdeploymentClient(
		c.transport,
		path.Join(c.path, "clusterdeployment"),
	)
}

// ControlPlane returns the target 'control_plane' resource.
//
// Reference to the resource that manages the collection of upgrade policies defined for control plane for this cluster.
func (c *ClusterClient) ControlPlane() *ControlPlaneClient {
	return NewControlPlaneClient(
		c.transport,
		path.Join(c.path, "control_plane"),
	)
}

// Credentials returns the target 'credentials' resource.
//
// Reference to the resource that manages the credentials of the cluster.
func (c *ClusterClient) Credentials() *CredentialsClient {
	return NewCredentialsClient(
		c.transport,
		path.Join(c.path, "credentials"),
	)
}

// DeleteProtection returns the target 'delete_protection' resource.
func (c *ClusterClient) DeleteProtection() *DeleteProtectionClient {
	return NewDeleteProtectionClient(
		c.transport,
		path.Join(c.path, "delete_protection"),
	)
}

// ExternalAuthConfig returns the target 'external_auth_config' resource.
//
// Reference to the resource that manages the external authentication configuration.
func (c *ClusterClient) ExternalAuthConfig() *ExternalAuthConfigClient {
	return NewExternalAuthConfigClient(
		c.transport,
		path.Join(c.path, "external_auth_config"),
	)
}

// ExternalConfiguration returns the target 'external_configuration' resource.
//
// Reference to the resource that manages the external configuration.
func (c *ClusterClient) ExternalConfiguration() *ExternalConfigurationClient {
	return NewExternalConfigurationClient(
		c.transport,
		path.Join(c.path, "external_configuration"),
	)
}

// GateAgreements returns the target 'version_gate_agreements' resource.
//
// Reference to cluster's agreed version gate.
func (c *ClusterClient) GateAgreements() *VersionGateAgreementsClient {
	return NewVersionGateAgreementsClient(
		c.transport,
		path.Join(c.path, "gate_agreements"),
	)
}

// Groups returns the target 'groups' resource.
//
// Reference to the resource that manages the collection of groups.
func (c *ClusterClient) Groups() *GroupsClient {
	return NewGroupsClient(
		c.transport,
		path.Join(c.path, "groups"),
	)
}

// Hypershift returns the target 'hypershift' resource.
func (c *ClusterClient) Hypershift() *HypershiftClient {
	return NewHypershiftClient(
		c.transport,
		path.Join(c.path, "hypershift"),
	)
}

// IdentityProviders returns the target 'identity_providers' resource.
//
// Reference to the resource that manages the collection of identity providers.
func (c *ClusterClient) IdentityProviders() *IdentityProvidersClient {
	return NewIdentityProvidersClient(
		c.transport,
		path.Join(c.path, "identity_providers"),
	)
}

// InflightChecks returns the target 'inflight_checks' resource.
//
// Reference to the resource that manages the collection of inflight checks.
func (c *ClusterClient) InflightChecks() *InflightChecksClient {
	return NewInflightChecksClient(
		c.transport,
		path.Join(c.path, "inflight_checks"),
	)
}

// Ingresses returns the target 'ingresses' resource.
//
// Reference to the resource that manages the collection of ingress resources.
func (c *ClusterClient) Ingresses() *IngressesClient {
	return NewIngressesClient(
		c.transport,
		path.Join(c.path, "ingresses"),
	)
}

// KubeletConfig returns the target 'kubelet_config' resource.
func (c *ClusterClient) KubeletConfig() *KubeletConfigClient {
	return NewKubeletConfigClient(
		c.transport,
		path.Join(c.path, "kubelet_config"),
	)
}

// KubeletConfigs returns the target 'kubelet_configs' resource.
func (c *ClusterClient) KubeletConfigs() *KubeletConfigsClient {
	return NewKubeletConfigsClient(
		c.transport,
		path.Join(c.path, "kubelet_configs"),
	)
}

// LimitedSupportReasons returns the target 'limited_support_reasons' resource.
//
// Reference to cluster limited support reasons.
func (c *ClusterClient) LimitedSupportReasons() *LimitedSupportReasonsClient {
	return NewLimitedSupportReasonsClient(
		c.transport,
		path.Join(c.path, "limited_support_reasons"),
	)
}

// Logs returns the target 'logs' resource.
//
// Reference to the resource that manages the collection of logs of the cluster.
func (c *ClusterClient) Logs() *LogsClient {
	return NewLogsClient(
		c.transport,
		path.Join(c.path, "logs"),
	)
}

// MachinePools returns the target 'machine_pools' resource.
//
// Reference to the resource that manages the collection of machine pool resources.
func (c *ClusterClient) MachinePools() *MachinePoolsClient {
	return NewMachinePoolsClient(
		c.transport,
		path.Join(c.path, "machine_pools"),
	)
}

// MetricQueries returns the target 'metric_queries' resource.
//
// Reference to the resource that manages metrics queries for the cluster.
func (c *ClusterClient) MetricQueries() *MetricQueriesClient {
	return NewMetricQueriesClient(
		c.transport,
		path.Join(c.path, "metric_queries"),
	)
}

// NodePools returns the target 'node_pools' resource.
//
// Reference to the resource that manages the collection of node pool resources.
func (c *ClusterClient) NodePools() *NodePoolsClient {
	return NewNodePoolsClient(
		c.transport,
		path.Join(c.path, "node_pools"),
	)
}

// ProvisionShard returns the target 'provision_shard' resource.
//
// Reference to the resource that manages the cluster's provision shard.
func (c *ClusterClient) ProvisionShard() *ProvisionShardClient {
	return NewProvisionShardClient(
		c.transport,
		path.Join(c.path, "provision_shard"),
	)
}

// Resources returns the target 'resources' resource.
//
// Reference to cluster resources.
func (c *ClusterClient) Resources() *ResourcesClient {
	return NewResourcesClient(
		c.transport,
		path.Join(c.path, "resources"),
	)
}

// Status returns the target 'cluster_status' resource.
//
// Reference to the resource that manages the detailed status of the cluster.
func (c *ClusterClient) Status() *ClusterStatusClient {
	return NewClusterStatusClient(
		c.transport,
		path.Join(c.path, "status"),
	)
}

// StsSupportJumpRole returns the target 'sts_support_jump_role' resource.
func (c *ClusterClient) StsSupportJumpRole() *StsSupportJumpRoleClient {
	return NewStsSupportJumpRoleClient(
		c.transport,
		path.Join(c.path, "sts_support_jump_role"),
	)
}

// TuningConfigs returns the target 'tuning_configs' resource.
//
// Reference to the resource that manages the collection of tuning configs for this cluster.
func (c *ClusterClient) TuningConfigs() *TuningConfigsClient {
	return NewTuningConfigsClient(
		c.transport,
		path.Join(c.path, "tuning_configs"),
	)
}

// UpgradePolicies returns the target 'upgrade_policies' resource.
//
// Reference to the resource that manages the collection of upgrade policies defined for this cluster.
func (c *ClusterClient) UpgradePolicies() *UpgradePoliciesClient {
	return NewUpgradePoliciesClient(
		c.transport,
		path.Join(c.path, "upgrade_policies"),
	)
}

// Vpc returns the target 'vpc' resource.
//
// Reference to the resource that manages the vpc resource.
func (c *ClusterClient) Vpc() *VpcClient {
	return NewVpcClient(
		c.transport,
		path.Join(c.path, "vpc"),
	)
}

// ClusterPollRequest is the request for the Poll method.
type ClusterPollRequest struct {
	request    *ClusterGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *ClusterPollRequest) Parameter(name string, value interface{}) *ClusterPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *ClusterPollRequest) Header(name string, value interface{}) *ClusterPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *ClusterPollRequest) Interval(value time.Duration) *ClusterPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *ClusterPollRequest) Status(value int) *ClusterPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *ClusterPollRequest) Predicate(value func(*ClusterGetResponse) bool) *ClusterPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*ClusterGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *ClusterPollRequest) StartContext(ctx context.Context) (response *ClusterPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &ClusterPollResponse{
			response: result.(*ClusterGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *ClusterPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// ClusterPollResponse is the response for the Poll method.
type ClusterPollResponse struct {
	response *ClusterGetResponse
}

// Status returns the response status code.
func (r *ClusterPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *ClusterPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *ClusterPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *ClusterPollResponse) Body() *Cluster {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ClusterPollResponse) GetBody() (value *Cluster, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *ClusterClient) Poll() *ClusterPollRequest {
	return &ClusterPollRequest{
		request: c.Get(),
	}
}

// ClusterDeleteRequest is the request for the 'delete' method.
type ClusterDeleteRequest struct {
	transport   http.RoundTripper
	path        string
	query       url.Values
	header      http.Header
	bestEffort  *bool
	deprovision *bool
	dryRun      *bool
}

// Parameter adds a query parameter.
func (r *ClusterDeleteRequest) Parameter(name string, value interface{}) *ClusterDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ClusterDeleteRequest) Header(name string, value interface{}) *ClusterDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ClusterDeleteRequest) Impersonate(user string) *ClusterDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// BestEffort sets the value of the 'best_effort' parameter.
//
// BestEffort flag is used to check if the cluster deletion should be best-effort mode or not.
func (r *ClusterDeleteRequest) BestEffort(value bool) *ClusterDeleteRequest {
	r.bestEffort = &value
	return r
}

// Deprovision sets the value of the 'deprovision' parameter.
//
// If false it will only delete from OCM but not the actual cluster resources.
// false is only allowed for OCP clusters. true by default.
func (r *ClusterDeleteRequest) Deprovision(value bool) *ClusterDeleteRequest {
	r.deprovision = &value
	return r
}

// DryRun sets the value of the 'dry_run' parameter.
//
// Dry run flag is used to check if the operation can be completed, but won't delete.
func (r *ClusterDeleteRequest) DryRun(value bool) *ClusterDeleteRequest {
	r.dryRun = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ClusterDeleteRequest) Send() (result *ClusterDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ClusterDeleteRequest) SendContext(ctx context.Context) (result *ClusterDeleteResponse, err error) {
	query := helpers.CopyQuery(r.query)
	if r.bestEffort != nil {
		helpers.AddValue(&query, "best_effort", *r.bestEffort)
	}
	if r.deprovision != nil {
		helpers.AddValue(&query, "deprovision", *r.deprovision)
	}
	if r.dryRun != nil {
		helpers.AddValue(&query, "dry_run", *r.dryRun)
	}
	header := helpers.CopyHeader(r.header)
	uri := &url.URL{
		Path:     r.path,
		RawQuery: query.Encode(),
	}
	request := &http.Request{
		Method: "DELETE",
		URL:    uri,
		Header: header,
	}
	if ctx != nil {
		request = request.WithContext(ctx)
	}
	response, err := r.transport.RoundTrip(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	result = &ClusterDeleteResponse{}
	result.status = response.StatusCode
	result.header = response.Header
	reader := bufio.NewReader(response.Body)
	_, err = reader.Peek(1)
	if err == io.EOF {
		err = nil
		return
	}
	if result.status >= 400 {
		result.err, err = errors.UnmarshalErrorStatus(reader, result.status)
		if err != nil {
			return
		}
		err = result.err
		return
	}
	return
}

// ClusterDeleteResponse is the response for the 'delete' method.
type ClusterDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *ClusterDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ClusterDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ClusterDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// ClusterGetRequest is the request for the 'get' method.
type ClusterGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *ClusterGetRequest) Parameter(name string, value interface{}) *ClusterGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ClusterGetRequest) Header(name string, value interface{}) *ClusterGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ClusterGetRequest) Impersonate(user string) *ClusterGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ClusterGetRequest) Send() (result *ClusterGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ClusterGetRequest) SendContext(ctx context.Context) (result *ClusterGetResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	uri := &url.URL{
		Path:     r.path,
		RawQuery: query.Encode(),
	}
	request := &http.Request{
		Method: "GET",
		URL:    uri,
		Header: header,
	}
	if ctx != nil {
		request = request.WithContext(ctx)
	}
	response, err := r.transport.RoundTrip(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	result = &ClusterGetResponse{}
	result.status = response.StatusCode
	result.header = response.Header
	reader := bufio.NewReader(response.Body)
	_, err = reader.Peek(1)
	if err == io.EOF {
		err = nil
		return
	}
	if result.status >= 400 {
		result.err, err = errors.UnmarshalErrorStatus(reader, result.status)
		if err != nil {
			return
		}
		err = result.err
		return
	}
	err = readClusterGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ClusterGetResponse is the response for the 'get' method.
type ClusterGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *Cluster
}

// Status returns the response status code.
func (r *ClusterGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ClusterGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ClusterGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *ClusterGetResponse) Body() *Cluster {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ClusterGetResponse) GetBody() (value *Cluster, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// ClusterHibernateRequest is the request for the 'hibernate' method.
type ClusterHibernateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *ClusterHibernateRequest) Parameter(name string, value interface{}) *ClusterHibernateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ClusterHibernateRequest) Header(name string, value interface{}) *ClusterHibernateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ClusterHibernateRequest) Impersonate(user string) *ClusterHibernateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ClusterHibernateRequest) Send() (result *ClusterHibernateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ClusterHibernateRequest) SendContext(ctx context.Context) (result *ClusterHibernateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	uri := &url.URL{
		Path:     r.path,
		RawQuery: query.Encode(),
	}
	request := &http.Request{
		Method: "POST",
		URL:    uri,
		Header: header,
	}
	if ctx != nil {
		request = request.WithContext(ctx)
	}
	response, err := r.transport.RoundTrip(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	result = &ClusterHibernateResponse{}
	result.status = response.StatusCode
	result.header = response.Header
	reader := bufio.NewReader(response.Body)
	_, err = reader.Peek(1)
	if err == io.EOF {
		err = nil
		return
	}
	if result.status >= 400 {
		result.err, err = errors.UnmarshalErrorStatus(reader, result.status)
		if err != nil {
			return
		}
		err = result.err
		return
	}
	return
}

// ClusterHibernateResponse is the response for the 'hibernate' method.
type ClusterHibernateResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *ClusterHibernateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ClusterHibernateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ClusterHibernateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// ClusterResumeRequest is the request for the 'resume' method.
type ClusterResumeRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *ClusterResumeRequest) Parameter(name string, value interface{}) *ClusterResumeRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ClusterResumeRequest) Header(name string, value interface{}) *ClusterResumeRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ClusterResumeRequest) Impersonate(user string) *ClusterResumeRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ClusterResumeRequest) Send() (result *ClusterResumeResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ClusterResumeRequest) SendContext(ctx context.Context) (result *ClusterResumeResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	uri := &url.URL{
		Path:     r.path,
		RawQuery: query.Encode(),
	}
	request := &http.Request{
		Method: "POST",
		URL:    uri,
		Header: header,
	}
	if ctx != nil {
		request = request.WithContext(ctx)
	}
	response, err := r.transport.RoundTrip(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	result = &ClusterResumeResponse{}
	result.status = response.StatusCode
	result.header = response.Header
	reader := bufio.NewReader(response.Body)
	_, err = reader.Peek(1)
	if err == io.EOF {
		err = nil
		return
	}
	if result.status >= 400 {
		result.err, err = errors.UnmarshalErrorStatus(reader, result.status)
		if err != nil {
			return
		}
		err = result.err
		return
	}
	return
}

// ClusterResumeResponse is the response for the 'resume' method.
type ClusterResumeResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *ClusterResumeResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ClusterResumeResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ClusterResumeResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// ClusterUpdateRequest is the request for the 'update' method.
type ClusterUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *Cluster
}

// Parameter adds a query parameter.
func (r *ClusterUpdateRequest) Parameter(name string, value interface{}) *ClusterUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ClusterUpdateRequest) Header(name string, value interface{}) *ClusterUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ClusterUpdateRequest) Impersonate(user string) *ClusterUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *ClusterUpdateRequest) Body(value *Cluster) *ClusterUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ClusterUpdateRequest) Send() (result *ClusterUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ClusterUpdateRequest) SendContext(ctx context.Context) (result *ClusterUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeClusterUpdateRequest(r, buffer)
	if err != nil {
		return
	}
	uri := &url.URL{
		Path:     r.path,
		RawQuery: query.Encode(),
	}
	request := &http.Request{
		Method: "PATCH",
		URL:    uri,
		Header: header,
		Body:   io.NopCloser(buffer),
	}
	if ctx != nil {
		request = request.WithContext(ctx)
	}
	response, err := r.transport.RoundTrip(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	result = &ClusterUpdateResponse{}
	result.status = response.StatusCode
	result.header = response.Header
	reader := bufio.NewReader(response.Body)
	_, err = reader.Peek(1)
	if err == io.EOF {
		err = nil
		return
	}
	if result.status >= 400 {
		result.err, err = errors.UnmarshalErrorStatus(reader, result.status)
		if err != nil {
			return
		}
		err = result.err
		return
	}
	err = readClusterUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ClusterUpdateResponse is the response for the 'update' method.
type ClusterUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *Cluster
}

// Status returns the response status code.
func (r *ClusterUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ClusterUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ClusterUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *ClusterUpdateResponse) Body() *Cluster {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ClusterUpdateResponse) GetBody() (value *Cluster, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
