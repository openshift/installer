// Copyright 2023 Google LLC. All Rights Reserved.
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
package containeraws

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/api/googleapi"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

type NodePool struct {
	Name              *string                    `json:"name"`
	Version           *string                    `json:"version"`
	Config            *NodePoolConfig            `json:"config"`
	Autoscaling       *NodePoolAutoscaling       `json:"autoscaling"`
	SubnetId          *string                    `json:"subnetId"`
	State             *NodePoolStateEnum         `json:"state"`
	Uid               *string                    `json:"uid"`
	Reconciling       *bool                      `json:"reconciling"`
	CreateTime        *string                    `json:"createTime"`
	UpdateTime        *string                    `json:"updateTime"`
	Etag              *string                    `json:"etag"`
	Annotations       map[string]string          `json:"annotations"`
	MaxPodsConstraint *NodePoolMaxPodsConstraint `json:"maxPodsConstraint"`
	Project           *string                    `json:"project"`
	Location          *string                    `json:"location"`
	Cluster           *string                    `json:"cluster"`
}

func (r *NodePool) String() string {
	return dcl.SprintResource(r)
}

// The enum NodePoolConfigRootVolumeVolumeTypeEnum.
type NodePoolConfigRootVolumeVolumeTypeEnum string

// NodePoolConfigRootVolumeVolumeTypeEnumRef returns a *NodePoolConfigRootVolumeVolumeTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func NodePoolConfigRootVolumeVolumeTypeEnumRef(s string) *NodePoolConfigRootVolumeVolumeTypeEnum {
	v := NodePoolConfigRootVolumeVolumeTypeEnum(s)
	return &v
}

func (v NodePoolConfigRootVolumeVolumeTypeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"VOLUME_TYPE_UNSPECIFIED", "GP2", "GP3"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "NodePoolConfigRootVolumeVolumeTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum NodePoolConfigTaintsEffectEnum.
type NodePoolConfigTaintsEffectEnum string

// NodePoolConfigTaintsEffectEnumRef returns a *NodePoolConfigTaintsEffectEnum with the value of string s
// If the empty string is provided, nil is returned.
func NodePoolConfigTaintsEffectEnumRef(s string) *NodePoolConfigTaintsEffectEnum {
	v := NodePoolConfigTaintsEffectEnum(s)
	return &v
}

func (v NodePoolConfigTaintsEffectEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"EFFECT_UNSPECIFIED", "NO_SCHEDULE", "PREFER_NO_SCHEDULE", "NO_EXECUTE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "NodePoolConfigTaintsEffectEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum NodePoolStateEnum.
type NodePoolStateEnum string

// NodePoolStateEnumRef returns a *NodePoolStateEnum with the value of string s
// If the empty string is provided, nil is returned.
func NodePoolStateEnumRef(s string) *NodePoolStateEnum {
	v := NodePoolStateEnum(s)
	return &v
}

func (v NodePoolStateEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"STATE_UNSPECIFIED", "PROVISIONING", "RUNNING", "RECONCILING", "STOPPING", "ERROR", "DEGRADED"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "NodePoolStateEnum",
		Value: string(v),
		Valid: []string{},
	}
}

type NodePoolConfig struct {
	empty                        bool                                        `json:"-"`
	InstanceType                 *string                                     `json:"instanceType"`
	RootVolume                   *NodePoolConfigRootVolume                   `json:"rootVolume"`
	Taints                       []NodePoolConfigTaints                      `json:"taints"`
	Labels                       map[string]string                           `json:"labels"`
	Tags                         map[string]string                           `json:"tags"`
	IamInstanceProfile           *string                                     `json:"iamInstanceProfile"`
	ConfigEncryption             *NodePoolConfigConfigEncryption             `json:"configEncryption"`
	SshConfig                    *NodePoolConfigSshConfig                    `json:"sshConfig"`
	SecurityGroupIds             []string                                    `json:"securityGroupIds"`
	ProxyConfig                  *NodePoolConfigProxyConfig                  `json:"proxyConfig"`
	AutoscalingMetricsCollection *NodePoolConfigAutoscalingMetricsCollection `json:"autoscalingMetricsCollection"`
}

type jsonNodePoolConfig NodePoolConfig

func (r *NodePoolConfig) UnmarshalJSON(data []byte) error {
	var res jsonNodePoolConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyNodePoolConfig
	} else {

		r.InstanceType = res.InstanceType

		r.RootVolume = res.RootVolume

		r.Taints = res.Taints

		r.Labels = res.Labels

		r.Tags = res.Tags

		r.IamInstanceProfile = res.IamInstanceProfile

		r.ConfigEncryption = res.ConfigEncryption

		r.SshConfig = res.SshConfig

		r.SecurityGroupIds = res.SecurityGroupIds

		r.ProxyConfig = res.ProxyConfig

		r.AutoscalingMetricsCollection = res.AutoscalingMetricsCollection

	}
	return nil
}

// This object is used to assert a desired state where this NodePoolConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyNodePoolConfig *NodePoolConfig = &NodePoolConfig{empty: true}

func (r *NodePoolConfig) Empty() bool {
	return r.empty
}

func (r *NodePoolConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *NodePoolConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type NodePoolConfigRootVolume struct {
	empty      bool                                    `json:"-"`
	SizeGib    *int64                                  `json:"sizeGib"`
	VolumeType *NodePoolConfigRootVolumeVolumeTypeEnum `json:"volumeType"`
	Iops       *int64                                  `json:"iops"`
	KmsKeyArn  *string                                 `json:"kmsKeyArn"`
}

type jsonNodePoolConfigRootVolume NodePoolConfigRootVolume

func (r *NodePoolConfigRootVolume) UnmarshalJSON(data []byte) error {
	var res jsonNodePoolConfigRootVolume
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyNodePoolConfigRootVolume
	} else {

		r.SizeGib = res.SizeGib

		r.VolumeType = res.VolumeType

		r.Iops = res.Iops

		r.KmsKeyArn = res.KmsKeyArn

	}
	return nil
}

// This object is used to assert a desired state where this NodePoolConfigRootVolume is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyNodePoolConfigRootVolume *NodePoolConfigRootVolume = &NodePoolConfigRootVolume{empty: true}

func (r *NodePoolConfigRootVolume) Empty() bool {
	return r.empty
}

func (r *NodePoolConfigRootVolume) String() string {
	return dcl.SprintResource(r)
}

func (r *NodePoolConfigRootVolume) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type NodePoolConfigTaints struct {
	empty  bool                            `json:"-"`
	Key    *string                         `json:"key"`
	Value  *string                         `json:"value"`
	Effect *NodePoolConfigTaintsEffectEnum `json:"effect"`
}

type jsonNodePoolConfigTaints NodePoolConfigTaints

func (r *NodePoolConfigTaints) UnmarshalJSON(data []byte) error {
	var res jsonNodePoolConfigTaints
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyNodePoolConfigTaints
	} else {

		r.Key = res.Key

		r.Value = res.Value

		r.Effect = res.Effect

	}
	return nil
}

// This object is used to assert a desired state where this NodePoolConfigTaints is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyNodePoolConfigTaints *NodePoolConfigTaints = &NodePoolConfigTaints{empty: true}

func (r *NodePoolConfigTaints) Empty() bool {
	return r.empty
}

func (r *NodePoolConfigTaints) String() string {
	return dcl.SprintResource(r)
}

func (r *NodePoolConfigTaints) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type NodePoolConfigConfigEncryption struct {
	empty     bool    `json:"-"`
	KmsKeyArn *string `json:"kmsKeyArn"`
}

type jsonNodePoolConfigConfigEncryption NodePoolConfigConfigEncryption

func (r *NodePoolConfigConfigEncryption) UnmarshalJSON(data []byte) error {
	var res jsonNodePoolConfigConfigEncryption
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyNodePoolConfigConfigEncryption
	} else {

		r.KmsKeyArn = res.KmsKeyArn

	}
	return nil
}

// This object is used to assert a desired state where this NodePoolConfigConfigEncryption is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyNodePoolConfigConfigEncryption *NodePoolConfigConfigEncryption = &NodePoolConfigConfigEncryption{empty: true}

func (r *NodePoolConfigConfigEncryption) Empty() bool {
	return r.empty
}

func (r *NodePoolConfigConfigEncryption) String() string {
	return dcl.SprintResource(r)
}

func (r *NodePoolConfigConfigEncryption) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type NodePoolConfigSshConfig struct {
	empty      bool    `json:"-"`
	Ec2KeyPair *string `json:"ec2KeyPair"`
}

type jsonNodePoolConfigSshConfig NodePoolConfigSshConfig

func (r *NodePoolConfigSshConfig) UnmarshalJSON(data []byte) error {
	var res jsonNodePoolConfigSshConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyNodePoolConfigSshConfig
	} else {

		r.Ec2KeyPair = res.Ec2KeyPair

	}
	return nil
}

// This object is used to assert a desired state where this NodePoolConfigSshConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyNodePoolConfigSshConfig *NodePoolConfigSshConfig = &NodePoolConfigSshConfig{empty: true}

func (r *NodePoolConfigSshConfig) Empty() bool {
	return r.empty
}

func (r *NodePoolConfigSshConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *NodePoolConfigSshConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type NodePoolConfigProxyConfig struct {
	empty         bool    `json:"-"`
	SecretArn     *string `json:"secretArn"`
	SecretVersion *string `json:"secretVersion"`
}

type jsonNodePoolConfigProxyConfig NodePoolConfigProxyConfig

func (r *NodePoolConfigProxyConfig) UnmarshalJSON(data []byte) error {
	var res jsonNodePoolConfigProxyConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyNodePoolConfigProxyConfig
	} else {

		r.SecretArn = res.SecretArn

		r.SecretVersion = res.SecretVersion

	}
	return nil
}

// This object is used to assert a desired state where this NodePoolConfigProxyConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyNodePoolConfigProxyConfig *NodePoolConfigProxyConfig = &NodePoolConfigProxyConfig{empty: true}

func (r *NodePoolConfigProxyConfig) Empty() bool {
	return r.empty
}

func (r *NodePoolConfigProxyConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *NodePoolConfigProxyConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type NodePoolConfigAutoscalingMetricsCollection struct {
	empty       bool     `json:"-"`
	Granularity *string  `json:"granularity"`
	Metrics     []string `json:"metrics"`
}

type jsonNodePoolConfigAutoscalingMetricsCollection NodePoolConfigAutoscalingMetricsCollection

func (r *NodePoolConfigAutoscalingMetricsCollection) UnmarshalJSON(data []byte) error {
	var res jsonNodePoolConfigAutoscalingMetricsCollection
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyNodePoolConfigAutoscalingMetricsCollection
	} else {

		r.Granularity = res.Granularity

		r.Metrics = res.Metrics

	}
	return nil
}

// This object is used to assert a desired state where this NodePoolConfigAutoscalingMetricsCollection is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyNodePoolConfigAutoscalingMetricsCollection *NodePoolConfigAutoscalingMetricsCollection = &NodePoolConfigAutoscalingMetricsCollection{empty: true}

func (r *NodePoolConfigAutoscalingMetricsCollection) Empty() bool {
	return r.empty
}

func (r *NodePoolConfigAutoscalingMetricsCollection) String() string {
	return dcl.SprintResource(r)
}

func (r *NodePoolConfigAutoscalingMetricsCollection) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type NodePoolAutoscaling struct {
	empty        bool   `json:"-"`
	MinNodeCount *int64 `json:"minNodeCount"`
	MaxNodeCount *int64 `json:"maxNodeCount"`
}

type jsonNodePoolAutoscaling NodePoolAutoscaling

func (r *NodePoolAutoscaling) UnmarshalJSON(data []byte) error {
	var res jsonNodePoolAutoscaling
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyNodePoolAutoscaling
	} else {

		r.MinNodeCount = res.MinNodeCount

		r.MaxNodeCount = res.MaxNodeCount

	}
	return nil
}

// This object is used to assert a desired state where this NodePoolAutoscaling is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyNodePoolAutoscaling *NodePoolAutoscaling = &NodePoolAutoscaling{empty: true}

func (r *NodePoolAutoscaling) Empty() bool {
	return r.empty
}

func (r *NodePoolAutoscaling) String() string {
	return dcl.SprintResource(r)
}

func (r *NodePoolAutoscaling) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type NodePoolMaxPodsConstraint struct {
	empty          bool   `json:"-"`
	MaxPodsPerNode *int64 `json:"maxPodsPerNode"`
}

type jsonNodePoolMaxPodsConstraint NodePoolMaxPodsConstraint

func (r *NodePoolMaxPodsConstraint) UnmarshalJSON(data []byte) error {
	var res jsonNodePoolMaxPodsConstraint
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyNodePoolMaxPodsConstraint
	} else {

		r.MaxPodsPerNode = res.MaxPodsPerNode

	}
	return nil
}

// This object is used to assert a desired state where this NodePoolMaxPodsConstraint is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyNodePoolMaxPodsConstraint *NodePoolMaxPodsConstraint = &NodePoolMaxPodsConstraint{empty: true}

func (r *NodePoolMaxPodsConstraint) Empty() bool {
	return r.empty
}

func (r *NodePoolMaxPodsConstraint) String() string {
	return dcl.SprintResource(r)
}

func (r *NodePoolMaxPodsConstraint) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *NodePool) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "container_aws",
		Type:    "NodePool",
		Version: "containeraws",
	}
}

func (r *NodePool) ID() (string, error) {
	if err := extractNodePoolFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name":                dcl.ValueOrEmptyString(nr.Name),
		"version":             dcl.ValueOrEmptyString(nr.Version),
		"config":              dcl.ValueOrEmptyString(nr.Config),
		"autoscaling":         dcl.ValueOrEmptyString(nr.Autoscaling),
		"subnet_id":           dcl.ValueOrEmptyString(nr.SubnetId),
		"state":               dcl.ValueOrEmptyString(nr.State),
		"uid":                 dcl.ValueOrEmptyString(nr.Uid),
		"reconciling":         dcl.ValueOrEmptyString(nr.Reconciling),
		"create_time":         dcl.ValueOrEmptyString(nr.CreateTime),
		"update_time":         dcl.ValueOrEmptyString(nr.UpdateTime),
		"etag":                dcl.ValueOrEmptyString(nr.Etag),
		"annotations":         dcl.ValueOrEmptyString(nr.Annotations),
		"max_pods_constraint": dcl.ValueOrEmptyString(nr.MaxPodsConstraint),
		"project":             dcl.ValueOrEmptyString(nr.Project),
		"location":            dcl.ValueOrEmptyString(nr.Location),
		"cluster":             dcl.ValueOrEmptyString(nr.Cluster),
	}
	return dcl.Nprintf("projects/{{project}}/locations/{{location}}/awsClusters/{{cluster}}/awsNodePools/{{name}}", params), nil
}

const NodePoolMaxPage = -1

type NodePoolList struct {
	Items []*NodePool

	nextToken string

	pageSize int32

	resource *NodePool
}

func (l *NodePoolList) HasNext() bool {
	return l.nextToken != ""
}

func (l *NodePoolList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listNodePool(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListNodePool(ctx context.Context, project, location, cluster string) (*NodePoolList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListNodePoolWithMaxResults(ctx, project, location, cluster, NodePoolMaxPage)

}

func (c *Client) ListNodePoolWithMaxResults(ctx context.Context, project, location, cluster string, pageSize int32) (*NodePoolList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &NodePool{
		Project:  &project,
		Location: &location,
		Cluster:  &cluster,
	}
	items, token, err := c.listNodePool(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &NodePoolList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetNodePool(ctx context.Context, r *NodePool) (*NodePool, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractNodePoolFields(r)

	b, err := c.getNodePoolRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalNodePool(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Location = r.Location
	result.Cluster = r.Cluster
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeNodePoolNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractNodePoolFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteNodePool(ctx context.Context, r *NodePool) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("NodePool resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting NodePool...")
	deleteOp := deleteNodePoolOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllNodePool deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllNodePool(ctx context.Context, project, location, cluster string, filter func(*NodePool) bool) error {
	listObj, err := c.ListNodePool(ctx, project, location, cluster)
	if err != nil {
		return err
	}

	err = c.deleteAllNodePool(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllNodePool(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyNodePool(ctx context.Context, rawDesired *NodePool, opts ...dcl.ApplyOption) (*NodePool, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *NodePool
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyNodePoolHelper(c, ctx, rawDesired, opts...)
		resultNewState = newState
		if err != nil {
			// If the error is 409, there is conflict in resource update.
			// Here we want to apply changes based on latest state.
			if dcl.IsConflictError(err) {
				return &dcl.RetryDetails{}, dcl.OperationNotDone{Err: err}
			}
			return nil, err
		}
		return nil, nil
	}, c.Config.RetryProvider)
	return resultNewState, err
}

func applyNodePoolHelper(c *Client, ctx context.Context, rawDesired *NodePool, opts ...dcl.ApplyOption) (*NodePool, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyNodePool...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractNodePoolFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.nodePoolDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToNodePoolDiffs(c.Config, fieldDiffs, opts)
	if err != nil {
		return nil, err
	}

	// TODO(magic-modules-eng): 2.2 Feasibility check (all updates are feasible so far).

	// 2.3: Lifecycle Directive Check
	var create bool
	lp := dcl.FetchLifecycleParams(opts)
	if initial == nil {
		if dcl.HasLifecycleParam(lp, dcl.BlockCreation) {
			return nil, dcl.ApplyInfeasibleError{Message: fmt.Sprintf("Creation blocked by lifecycle params: %#v.", desired)}
		}
		create = true
	} else if dcl.HasLifecycleParam(lp, dcl.BlockAcquire) {
		return nil, dcl.ApplyInfeasibleError{
			Message: fmt.Sprintf("Resource already exists - apply blocked by lifecycle params: %#v.", initial),
		}
	} else {
		for _, d := range diffs {
			if d.RequiresRecreate {
				return nil, dcl.ApplyInfeasibleError{
					Message: fmt.Sprintf("infeasible update: (%v) would require recreation", d),
				}
			}
			if dcl.HasLifecycleParam(lp, dcl.BlockModification) {
				return nil, dcl.ApplyInfeasibleError{Message: fmt.Sprintf("Modification blocked, diff (%v) unresolvable.", d)}
			}
		}
	}

	// 2.4 Imperative Request Planning
	var ops []nodePoolApiOperation
	if create {
		ops = append(ops, &createNodePoolOperation{})
	} else {
		for _, d := range diffs {
			ops = append(ops, d.UpdateOp)
		}
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created plan: %#v", ops)

	// 2.5 Request Actuation
	for _, op := range ops {
		c.Config.Logger.InfoWithContextf(ctx, "Performing operation %T %+v", op, op)
		if err := op.do(ctx, desired, c); err != nil {
			c.Config.Logger.InfoWithContextf(ctx, "Failed operation %T %+v: %v", op, op, err)
			return nil, err
		}
		c.Config.Logger.InfoWithContextf(ctx, "Finished operation %T %+v", op, op)
	}
	return applyNodePoolDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyNodePoolDiff(c *Client, ctx context.Context, desired *NodePool, rawDesired *NodePool, ops []nodePoolApiOperation, opts ...dcl.ApplyOption) (*NodePool, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetNodePool(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createNodePoolOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapNodePool(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeNodePoolNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeNodePoolNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeNodePoolDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractNodePoolFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractNodePoolFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffNodePool(c, newDesired, newState)
	if err != nil {
		return newState, err
	}

	if len(newDiffs) == 0 {
		c.Config.Logger.InfoWithContext(ctx, "No diffs found. Apply was successful.")
	} else {
		c.Config.Logger.InfoWithContextf(ctx, "Found diffs: %v", newDiffs)
		diffMessages := make([]string, len(newDiffs))
		for i, d := range newDiffs {
			diffMessages[i] = fmt.Sprintf("%v", d)
		}
		return newState, dcl.DiffAfterApplyError{Diffs: diffMessages}
	}
	c.Config.Logger.InfoWithContext(ctx, "Done Apply.")
	return newState, nil
}
