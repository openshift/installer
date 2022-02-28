// Copyright 2021 Google LLC. All Rights Reserved.
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
package dataproc

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/api/googleapi"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

type Cluster struct {
	Project       *string                `json:"project"`
	Name          *string                `json:"name"`
	Config        *ClusterClusterConfig  `json:"config"`
	Labels        map[string]string      `json:"labels"`
	Status        *ClusterStatus         `json:"status"`
	StatusHistory []ClusterStatusHistory `json:"statusHistory"`
	ClusterUuid   *string                `json:"clusterUuid"`
	Metrics       *ClusterMetrics        `json:"metrics"`
	Location      *string                `json:"location"`
}

func (r *Cluster) String() string {
	return dcl.SprintResource(r)
}

// The enum ClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum.
type ClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum string

// ClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnumRef returns a *ClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum with the value of string s
// If the empty string is provided, nil is returned.
func ClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnumRef(s string) *ClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum {
	v := ClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum(s)
	return &v
}

func (v ClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"PRIVATE_IPV6_GOOGLE_ACCESS_UNSPECIFIED", "INHERIT_FROM_SUBNETWORK", "OUTBOUND", "BIDIRECTIONAL"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "ClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum ClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum.
type ClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum string

// ClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnumRef returns a *ClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func ClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnumRef(s string) *ClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum {
	v := ClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum(s)
	return &v
}

func (v ClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"TYPE_UNSPECIFIED", "NO_RESERVATION", "ANY_RESERVATION", "SPECIFIC_RESERVATION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "ClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum ClusterInstanceGroupConfigPreemptibilityEnum.
type ClusterInstanceGroupConfigPreemptibilityEnum string

// ClusterInstanceGroupConfigPreemptibilityEnumRef returns a *ClusterInstanceGroupConfigPreemptibilityEnum with the value of string s
// If the empty string is provided, nil is returned.
func ClusterInstanceGroupConfigPreemptibilityEnumRef(s string) *ClusterInstanceGroupConfigPreemptibilityEnum {
	v := ClusterInstanceGroupConfigPreemptibilityEnum(s)
	return &v
}

func (v ClusterInstanceGroupConfigPreemptibilityEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"PREEMPTIBILITY_UNSPECIFIED", "NON_PREEMPTIBLE", "PREEMPTIBLE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "ClusterInstanceGroupConfigPreemptibilityEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum ClusterClusterConfigSoftwareConfigOptionalComponentsEnum.
type ClusterClusterConfigSoftwareConfigOptionalComponentsEnum string

// ClusterClusterConfigSoftwareConfigOptionalComponentsEnumRef returns a *ClusterClusterConfigSoftwareConfigOptionalComponentsEnum with the value of string s
// If the empty string is provided, nil is returned.
func ClusterClusterConfigSoftwareConfigOptionalComponentsEnumRef(s string) *ClusterClusterConfigSoftwareConfigOptionalComponentsEnum {
	v := ClusterClusterConfigSoftwareConfigOptionalComponentsEnum(s)
	return &v
}

func (v ClusterClusterConfigSoftwareConfigOptionalComponentsEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"COMPONENT_UNSPECIFIED", "ANACONDA", "DOCKER", "DRUID", "FLINK", "HBASE", "HIVE_WEBHCAT", "JUPYTER", "KERBEROS", "PRESTO", "RANGER", "SOLR", "ZEPPELIN", "ZOOKEEPER"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "ClusterClusterConfigSoftwareConfigOptionalComponentsEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum ClusterStatusStateEnum.
type ClusterStatusStateEnum string

// ClusterStatusStateEnumRef returns a *ClusterStatusStateEnum with the value of string s
// If the empty string is provided, nil is returned.
func ClusterStatusStateEnumRef(s string) *ClusterStatusStateEnum {
	v := ClusterStatusStateEnum(s)
	return &v
}

func (v ClusterStatusStateEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"UNKNOWN", "CREATING", "RUNNING", "ERROR", "DELETING", "UPDATING", "STOPPING", "STOPPED", "STARTING"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "ClusterStatusStateEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum ClusterStatusSubstateEnum.
type ClusterStatusSubstateEnum string

// ClusterStatusSubstateEnumRef returns a *ClusterStatusSubstateEnum with the value of string s
// If the empty string is provided, nil is returned.
func ClusterStatusSubstateEnumRef(s string) *ClusterStatusSubstateEnum {
	v := ClusterStatusSubstateEnum(s)
	return &v
}

func (v ClusterStatusSubstateEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"UNSPECIFIED", "UNHEALTHY", "STALE_STATUS"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "ClusterStatusSubstateEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum ClusterStatusHistoryStateEnum.
type ClusterStatusHistoryStateEnum string

// ClusterStatusHistoryStateEnumRef returns a *ClusterStatusHistoryStateEnum with the value of string s
// If the empty string is provided, nil is returned.
func ClusterStatusHistoryStateEnumRef(s string) *ClusterStatusHistoryStateEnum {
	v := ClusterStatusHistoryStateEnum(s)
	return &v
}

func (v ClusterStatusHistoryStateEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"UNKNOWN", "CREATING", "RUNNING", "ERROR", "DELETING", "UPDATING", "STOPPING", "STOPPED", "STARTING"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "ClusterStatusHistoryStateEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum ClusterStatusHistorySubstateEnum.
type ClusterStatusHistorySubstateEnum string

// ClusterStatusHistorySubstateEnumRef returns a *ClusterStatusHistorySubstateEnum with the value of string s
// If the empty string is provided, nil is returned.
func ClusterStatusHistorySubstateEnumRef(s string) *ClusterStatusHistorySubstateEnum {
	v := ClusterStatusHistorySubstateEnum(s)
	return &v
}

func (v ClusterStatusHistorySubstateEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"UNSPECIFIED", "UNHEALTHY", "STALE_STATUS"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "ClusterStatusHistorySubstateEnum",
		Value: string(v),
		Valid: []string{},
	}
}

type ClusterClusterConfig struct {
	empty                 bool                                        `json:"-"`
	StagingBucket         *string                                     `json:"stagingBucket"`
	TempBucket            *string                                     `json:"tempBucket"`
	GceClusterConfig      *ClusterClusterConfigGceClusterConfig       `json:"gceClusterConfig"`
	MasterConfig          *ClusterInstanceGroupConfig                 `json:"masterConfig"`
	WorkerConfig          *ClusterInstanceGroupConfig                 `json:"workerConfig"`
	SecondaryWorkerConfig *ClusterInstanceGroupConfig                 `json:"secondaryWorkerConfig"`
	SoftwareConfig        *ClusterClusterConfigSoftwareConfig         `json:"softwareConfig"`
	InitializationActions []ClusterClusterConfigInitializationActions `json:"initializationActions"`
	EncryptionConfig      *ClusterClusterConfigEncryptionConfig       `json:"encryptionConfig"`
	AutoscalingConfig     *ClusterClusterConfigAutoscalingConfig      `json:"autoscalingConfig"`
	SecurityConfig        *ClusterClusterConfigSecurityConfig         `json:"securityConfig"`
	LifecycleConfig       *ClusterClusterConfigLifecycleConfig        `json:"lifecycleConfig"`
	EndpointConfig        *ClusterClusterConfigEndpointConfig         `json:"endpointConfig"`
}

type jsonClusterClusterConfig ClusterClusterConfig

func (r *ClusterClusterConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterClusterConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterClusterConfig
	} else {

		r.StagingBucket = res.StagingBucket

		r.TempBucket = res.TempBucket

		r.GceClusterConfig = res.GceClusterConfig

		r.MasterConfig = res.MasterConfig

		r.WorkerConfig = res.WorkerConfig

		r.SecondaryWorkerConfig = res.SecondaryWorkerConfig

		r.SoftwareConfig = res.SoftwareConfig

		r.InitializationActions = res.InitializationActions

		r.EncryptionConfig = res.EncryptionConfig

		r.AutoscalingConfig = res.AutoscalingConfig

		r.SecurityConfig = res.SecurityConfig

		r.LifecycleConfig = res.LifecycleConfig

		r.EndpointConfig = res.EndpointConfig

	}
	return nil
}

// This object is used to assert a desired state where this ClusterClusterConfig is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyClusterClusterConfig *ClusterClusterConfig = &ClusterClusterConfig{empty: true}

func (r *ClusterClusterConfig) Empty() bool {
	return r.empty
}

func (r *ClusterClusterConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterClusterConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterClusterConfigGceClusterConfig struct {
	empty                   bool                                                             `json:"-"`
	Zone                    *string                                                          `json:"zone"`
	Network                 *string                                                          `json:"network"`
	Subnetwork              *string                                                          `json:"subnetwork"`
	InternalIPOnly          *bool                                                            `json:"internalIPOnly"`
	PrivateIPv6GoogleAccess *ClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum `json:"privateIPv6GoogleAccess"`
	ServiceAccount          *string                                                          `json:"serviceAccount"`
	ServiceAccountScopes    []string                                                         `json:"serviceAccountScopes"`
	Tags                    []string                                                         `json:"tags"`
	Metadata                map[string]string                                                `json:"metadata"`
	ReservationAffinity     *ClusterClusterConfigGceClusterConfigReservationAffinity         `json:"reservationAffinity"`
	NodeGroupAffinity       *ClusterClusterConfigGceClusterConfigNodeGroupAffinity           `json:"nodeGroupAffinity"`
}

type jsonClusterClusterConfigGceClusterConfig ClusterClusterConfigGceClusterConfig

func (r *ClusterClusterConfigGceClusterConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterClusterConfigGceClusterConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterClusterConfigGceClusterConfig
	} else {

		r.Zone = res.Zone

		r.Network = res.Network

		r.Subnetwork = res.Subnetwork

		r.InternalIPOnly = res.InternalIPOnly

		r.PrivateIPv6GoogleAccess = res.PrivateIPv6GoogleAccess

		r.ServiceAccount = res.ServiceAccount

		r.ServiceAccountScopes = res.ServiceAccountScopes

		r.Tags = res.Tags

		r.Metadata = res.Metadata

		r.ReservationAffinity = res.ReservationAffinity

		r.NodeGroupAffinity = res.NodeGroupAffinity

	}
	return nil
}

// This object is used to assert a desired state where this ClusterClusterConfigGceClusterConfig is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyClusterClusterConfigGceClusterConfig *ClusterClusterConfigGceClusterConfig = &ClusterClusterConfigGceClusterConfig{empty: true}

func (r *ClusterClusterConfigGceClusterConfig) Empty() bool {
	return r.empty
}

func (r *ClusterClusterConfigGceClusterConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterClusterConfigGceClusterConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterClusterConfigGceClusterConfigReservationAffinity struct {
	empty                  bool                                                                               `json:"-"`
	ConsumeReservationType *ClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum `json:"consumeReservationType"`
	Key                    *string                                                                            `json:"key"`
	Values                 []string                                                                           `json:"values"`
}

type jsonClusterClusterConfigGceClusterConfigReservationAffinity ClusterClusterConfigGceClusterConfigReservationAffinity

func (r *ClusterClusterConfigGceClusterConfigReservationAffinity) UnmarshalJSON(data []byte) error {
	var res jsonClusterClusterConfigGceClusterConfigReservationAffinity
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterClusterConfigGceClusterConfigReservationAffinity
	} else {

		r.ConsumeReservationType = res.ConsumeReservationType

		r.Key = res.Key

		r.Values = res.Values

	}
	return nil
}

// This object is used to assert a desired state where this ClusterClusterConfigGceClusterConfigReservationAffinity is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyClusterClusterConfigGceClusterConfigReservationAffinity *ClusterClusterConfigGceClusterConfigReservationAffinity = &ClusterClusterConfigGceClusterConfigReservationAffinity{empty: true}

func (r *ClusterClusterConfigGceClusterConfigReservationAffinity) Empty() bool {
	return r.empty
}

func (r *ClusterClusterConfigGceClusterConfigReservationAffinity) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterClusterConfigGceClusterConfigReservationAffinity) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterClusterConfigGceClusterConfigNodeGroupAffinity struct {
	empty     bool    `json:"-"`
	NodeGroup *string `json:"nodeGroup"`
}

type jsonClusterClusterConfigGceClusterConfigNodeGroupAffinity ClusterClusterConfigGceClusterConfigNodeGroupAffinity

func (r *ClusterClusterConfigGceClusterConfigNodeGroupAffinity) UnmarshalJSON(data []byte) error {
	var res jsonClusterClusterConfigGceClusterConfigNodeGroupAffinity
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterClusterConfigGceClusterConfigNodeGroupAffinity
	} else {

		r.NodeGroup = res.NodeGroup

	}
	return nil
}

// This object is used to assert a desired state where this ClusterClusterConfigGceClusterConfigNodeGroupAffinity is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyClusterClusterConfigGceClusterConfigNodeGroupAffinity *ClusterClusterConfigGceClusterConfigNodeGroupAffinity = &ClusterClusterConfigGceClusterConfigNodeGroupAffinity{empty: true}

func (r *ClusterClusterConfigGceClusterConfigNodeGroupAffinity) Empty() bool {
	return r.empty
}

func (r *ClusterClusterConfigGceClusterConfigNodeGroupAffinity) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterClusterConfigGceClusterConfigNodeGroupAffinity) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterInstanceGroupConfig struct {
	empty              bool                                          `json:"-"`
	NumInstances       *int64                                        `json:"numInstances"`
	InstanceNames      []string                                      `json:"instanceNames"`
	Image              *string                                       `json:"image"`
	MachineType        *string                                       `json:"machineType"`
	DiskConfig         *ClusterInstanceGroupConfigDiskConfig         `json:"diskConfig"`
	IsPreemptible      *bool                                         `json:"isPreemptible"`
	Preemptibility     *ClusterInstanceGroupConfigPreemptibilityEnum `json:"preemptibility"`
	ManagedGroupConfig *ClusterInstanceGroupConfigManagedGroupConfig `json:"managedGroupConfig"`
	Accelerators       []ClusterInstanceGroupConfigAccelerators      `json:"accelerators"`
	MinCpuPlatform     *string                                       `json:"minCpuPlatform"`
}

type jsonClusterInstanceGroupConfig ClusterInstanceGroupConfig

func (r *ClusterInstanceGroupConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterInstanceGroupConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterInstanceGroupConfig
	} else {

		r.NumInstances = res.NumInstances

		r.InstanceNames = res.InstanceNames

		r.Image = res.Image

		r.MachineType = res.MachineType

		r.DiskConfig = res.DiskConfig

		r.IsPreemptible = res.IsPreemptible

		r.Preemptibility = res.Preemptibility

		r.ManagedGroupConfig = res.ManagedGroupConfig

		r.Accelerators = res.Accelerators

		r.MinCpuPlatform = res.MinCpuPlatform

	}
	return nil
}

// This object is used to assert a desired state where this ClusterInstanceGroupConfig is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyClusterInstanceGroupConfig *ClusterInstanceGroupConfig = &ClusterInstanceGroupConfig{empty: true}

func (r *ClusterInstanceGroupConfig) Empty() bool {
	return r.empty
}

func (r *ClusterInstanceGroupConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterInstanceGroupConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterInstanceGroupConfigDiskConfig struct {
	empty          bool    `json:"-"`
	BootDiskType   *string `json:"bootDiskType"`
	BootDiskSizeGb *int64  `json:"bootDiskSizeGb"`
	NumLocalSsds   *int64  `json:"numLocalSsds"`
}

type jsonClusterInstanceGroupConfigDiskConfig ClusterInstanceGroupConfigDiskConfig

func (r *ClusterInstanceGroupConfigDiskConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterInstanceGroupConfigDiskConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterInstanceGroupConfigDiskConfig
	} else {

		r.BootDiskType = res.BootDiskType

		r.BootDiskSizeGb = res.BootDiskSizeGb

		r.NumLocalSsds = res.NumLocalSsds

	}
	return nil
}

// This object is used to assert a desired state where this ClusterInstanceGroupConfigDiskConfig is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyClusterInstanceGroupConfigDiskConfig *ClusterInstanceGroupConfigDiskConfig = &ClusterInstanceGroupConfigDiskConfig{empty: true}

func (r *ClusterInstanceGroupConfigDiskConfig) Empty() bool {
	return r.empty
}

func (r *ClusterInstanceGroupConfigDiskConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterInstanceGroupConfigDiskConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterInstanceGroupConfigManagedGroupConfig struct {
	empty                    bool    `json:"-"`
	InstanceTemplateName     *string `json:"instanceTemplateName"`
	InstanceGroupManagerName *string `json:"instanceGroupManagerName"`
}

type jsonClusterInstanceGroupConfigManagedGroupConfig ClusterInstanceGroupConfigManagedGroupConfig

func (r *ClusterInstanceGroupConfigManagedGroupConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterInstanceGroupConfigManagedGroupConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterInstanceGroupConfigManagedGroupConfig
	} else {

		r.InstanceTemplateName = res.InstanceTemplateName

		r.InstanceGroupManagerName = res.InstanceGroupManagerName

	}
	return nil
}

// This object is used to assert a desired state where this ClusterInstanceGroupConfigManagedGroupConfig is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyClusterInstanceGroupConfigManagedGroupConfig *ClusterInstanceGroupConfigManagedGroupConfig = &ClusterInstanceGroupConfigManagedGroupConfig{empty: true}

func (r *ClusterInstanceGroupConfigManagedGroupConfig) Empty() bool {
	return r.empty
}

func (r *ClusterInstanceGroupConfigManagedGroupConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterInstanceGroupConfigManagedGroupConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterInstanceGroupConfigAccelerators struct {
	empty            bool    `json:"-"`
	AcceleratorType  *string `json:"acceleratorType"`
	AcceleratorCount *int64  `json:"acceleratorCount"`
}

type jsonClusterInstanceGroupConfigAccelerators ClusterInstanceGroupConfigAccelerators

func (r *ClusterInstanceGroupConfigAccelerators) UnmarshalJSON(data []byte) error {
	var res jsonClusterInstanceGroupConfigAccelerators
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterInstanceGroupConfigAccelerators
	} else {

		r.AcceleratorType = res.AcceleratorType

		r.AcceleratorCount = res.AcceleratorCount

	}
	return nil
}

// This object is used to assert a desired state where this ClusterInstanceGroupConfigAccelerators is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyClusterInstanceGroupConfigAccelerators *ClusterInstanceGroupConfigAccelerators = &ClusterInstanceGroupConfigAccelerators{empty: true}

func (r *ClusterInstanceGroupConfigAccelerators) Empty() bool {
	return r.empty
}

func (r *ClusterInstanceGroupConfigAccelerators) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterInstanceGroupConfigAccelerators) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterClusterConfigSoftwareConfig struct {
	empty              bool                                                       `json:"-"`
	ImageVersion       *string                                                    `json:"imageVersion"`
	Properties         map[string]string                                          `json:"properties"`
	OptionalComponents []ClusterClusterConfigSoftwareConfigOptionalComponentsEnum `json:"optionalComponents"`
}

type jsonClusterClusterConfigSoftwareConfig ClusterClusterConfigSoftwareConfig

func (r *ClusterClusterConfigSoftwareConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterClusterConfigSoftwareConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterClusterConfigSoftwareConfig
	} else {

		r.ImageVersion = res.ImageVersion

		r.Properties = res.Properties

		r.OptionalComponents = res.OptionalComponents

	}
	return nil
}

// This object is used to assert a desired state where this ClusterClusterConfigSoftwareConfig is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyClusterClusterConfigSoftwareConfig *ClusterClusterConfigSoftwareConfig = &ClusterClusterConfigSoftwareConfig{empty: true}

func (r *ClusterClusterConfigSoftwareConfig) Empty() bool {
	return r.empty
}

func (r *ClusterClusterConfigSoftwareConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterClusterConfigSoftwareConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterClusterConfigInitializationActions struct {
	empty            bool    `json:"-"`
	ExecutableFile   *string `json:"executableFile"`
	ExecutionTimeout *string `json:"executionTimeout"`
}

type jsonClusterClusterConfigInitializationActions ClusterClusterConfigInitializationActions

func (r *ClusterClusterConfigInitializationActions) UnmarshalJSON(data []byte) error {
	var res jsonClusterClusterConfigInitializationActions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterClusterConfigInitializationActions
	} else {

		r.ExecutableFile = res.ExecutableFile

		r.ExecutionTimeout = res.ExecutionTimeout

	}
	return nil
}

// This object is used to assert a desired state where this ClusterClusterConfigInitializationActions is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyClusterClusterConfigInitializationActions *ClusterClusterConfigInitializationActions = &ClusterClusterConfigInitializationActions{empty: true}

func (r *ClusterClusterConfigInitializationActions) Empty() bool {
	return r.empty
}

func (r *ClusterClusterConfigInitializationActions) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterClusterConfigInitializationActions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterClusterConfigEncryptionConfig struct {
	empty           bool    `json:"-"`
	GcePdKmsKeyName *string `json:"gcePdKmsKeyName"`
}

type jsonClusterClusterConfigEncryptionConfig ClusterClusterConfigEncryptionConfig

func (r *ClusterClusterConfigEncryptionConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterClusterConfigEncryptionConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterClusterConfigEncryptionConfig
	} else {

		r.GcePdKmsKeyName = res.GcePdKmsKeyName

	}
	return nil
}

// This object is used to assert a desired state where this ClusterClusterConfigEncryptionConfig is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyClusterClusterConfigEncryptionConfig *ClusterClusterConfigEncryptionConfig = &ClusterClusterConfigEncryptionConfig{empty: true}

func (r *ClusterClusterConfigEncryptionConfig) Empty() bool {
	return r.empty
}

func (r *ClusterClusterConfigEncryptionConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterClusterConfigEncryptionConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterClusterConfigAutoscalingConfig struct {
	empty  bool    `json:"-"`
	Policy *string `json:"policy"`
}

type jsonClusterClusterConfigAutoscalingConfig ClusterClusterConfigAutoscalingConfig

func (r *ClusterClusterConfigAutoscalingConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterClusterConfigAutoscalingConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterClusterConfigAutoscalingConfig
	} else {

		r.Policy = res.Policy

	}
	return nil
}

// This object is used to assert a desired state where this ClusterClusterConfigAutoscalingConfig is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyClusterClusterConfigAutoscalingConfig *ClusterClusterConfigAutoscalingConfig = &ClusterClusterConfigAutoscalingConfig{empty: true}

func (r *ClusterClusterConfigAutoscalingConfig) Empty() bool {
	return r.empty
}

func (r *ClusterClusterConfigAutoscalingConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterClusterConfigAutoscalingConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterClusterConfigSecurityConfig struct {
	empty          bool                                              `json:"-"`
	KerberosConfig *ClusterClusterConfigSecurityConfigKerberosConfig `json:"kerberosConfig"`
}

type jsonClusterClusterConfigSecurityConfig ClusterClusterConfigSecurityConfig

func (r *ClusterClusterConfigSecurityConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterClusterConfigSecurityConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterClusterConfigSecurityConfig
	} else {

		r.KerberosConfig = res.KerberosConfig

	}
	return nil
}

// This object is used to assert a desired state where this ClusterClusterConfigSecurityConfig is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyClusterClusterConfigSecurityConfig *ClusterClusterConfigSecurityConfig = &ClusterClusterConfigSecurityConfig{empty: true}

func (r *ClusterClusterConfigSecurityConfig) Empty() bool {
	return r.empty
}

func (r *ClusterClusterConfigSecurityConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterClusterConfigSecurityConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterClusterConfigSecurityConfigKerberosConfig struct {
	empty                         bool    `json:"-"`
	EnableKerberos                *bool   `json:"enableKerberos"`
	RootPrincipalPassword         *string `json:"rootPrincipalPassword"`
	KmsKey                        *string `json:"kmsKey"`
	Keystore                      *string `json:"keystore"`
	Truststore                    *string `json:"truststore"`
	KeystorePassword              *string `json:"keystorePassword"`
	KeyPassword                   *string `json:"keyPassword"`
	TruststorePassword            *string `json:"truststorePassword"`
	CrossRealmTrustRealm          *string `json:"crossRealmTrustRealm"`
	CrossRealmTrustKdc            *string `json:"crossRealmTrustKdc"`
	CrossRealmTrustAdminServer    *string `json:"crossRealmTrustAdminServer"`
	CrossRealmTrustSharedPassword *string `json:"crossRealmTrustSharedPassword"`
	KdcDbKey                      *string `json:"kdcDbKey"`
	TgtLifetimeHours              *int64  `json:"tgtLifetimeHours"`
	Realm                         *string `json:"realm"`
}

type jsonClusterClusterConfigSecurityConfigKerberosConfig ClusterClusterConfigSecurityConfigKerberosConfig

func (r *ClusterClusterConfigSecurityConfigKerberosConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterClusterConfigSecurityConfigKerberosConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterClusterConfigSecurityConfigKerberosConfig
	} else {

		r.EnableKerberos = res.EnableKerberos

		r.RootPrincipalPassword = res.RootPrincipalPassword

		r.KmsKey = res.KmsKey

		r.Keystore = res.Keystore

		r.Truststore = res.Truststore

		r.KeystorePassword = res.KeystorePassword

		r.KeyPassword = res.KeyPassword

		r.TruststorePassword = res.TruststorePassword

		r.CrossRealmTrustRealm = res.CrossRealmTrustRealm

		r.CrossRealmTrustKdc = res.CrossRealmTrustKdc

		r.CrossRealmTrustAdminServer = res.CrossRealmTrustAdminServer

		r.CrossRealmTrustSharedPassword = res.CrossRealmTrustSharedPassword

		r.KdcDbKey = res.KdcDbKey

		r.TgtLifetimeHours = res.TgtLifetimeHours

		r.Realm = res.Realm

	}
	return nil
}

// This object is used to assert a desired state where this ClusterClusterConfigSecurityConfigKerberosConfig is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyClusterClusterConfigSecurityConfigKerberosConfig *ClusterClusterConfigSecurityConfigKerberosConfig = &ClusterClusterConfigSecurityConfigKerberosConfig{empty: true}

func (r *ClusterClusterConfigSecurityConfigKerberosConfig) Empty() bool {
	return r.empty
}

func (r *ClusterClusterConfigSecurityConfigKerberosConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterClusterConfigSecurityConfigKerberosConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterClusterConfigLifecycleConfig struct {
	empty          bool    `json:"-"`
	IdleDeleteTtl  *string `json:"idleDeleteTtl"`
	AutoDeleteTime *string `json:"autoDeleteTime"`
	AutoDeleteTtl  *string `json:"autoDeleteTtl"`
	IdleStartTime  *string `json:"idleStartTime"`
}

type jsonClusterClusterConfigLifecycleConfig ClusterClusterConfigLifecycleConfig

func (r *ClusterClusterConfigLifecycleConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterClusterConfigLifecycleConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterClusterConfigLifecycleConfig
	} else {

		r.IdleDeleteTtl = res.IdleDeleteTtl

		r.AutoDeleteTime = res.AutoDeleteTime

		r.AutoDeleteTtl = res.AutoDeleteTtl

		r.IdleStartTime = res.IdleStartTime

	}
	return nil
}

// This object is used to assert a desired state where this ClusterClusterConfigLifecycleConfig is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyClusterClusterConfigLifecycleConfig *ClusterClusterConfigLifecycleConfig = &ClusterClusterConfigLifecycleConfig{empty: true}

func (r *ClusterClusterConfigLifecycleConfig) Empty() bool {
	return r.empty
}

func (r *ClusterClusterConfigLifecycleConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterClusterConfigLifecycleConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterClusterConfigEndpointConfig struct {
	empty                bool              `json:"-"`
	HttpPorts            map[string]string `json:"httpPorts"`
	EnableHttpPortAccess *bool             `json:"enableHttpPortAccess"`
}

type jsonClusterClusterConfigEndpointConfig ClusterClusterConfigEndpointConfig

func (r *ClusterClusterConfigEndpointConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterClusterConfigEndpointConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterClusterConfigEndpointConfig
	} else {

		r.HttpPorts = res.HttpPorts

		r.EnableHttpPortAccess = res.EnableHttpPortAccess

	}
	return nil
}

// This object is used to assert a desired state where this ClusterClusterConfigEndpointConfig is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyClusterClusterConfigEndpointConfig *ClusterClusterConfigEndpointConfig = &ClusterClusterConfigEndpointConfig{empty: true}

func (r *ClusterClusterConfigEndpointConfig) Empty() bool {
	return r.empty
}

func (r *ClusterClusterConfigEndpointConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterClusterConfigEndpointConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterStatus struct {
	empty          bool                       `json:"-"`
	State          *ClusterStatusStateEnum    `json:"state"`
	Detail         *string                    `json:"detail"`
	StateStartTime *string                    `json:"stateStartTime"`
	Substate       *ClusterStatusSubstateEnum `json:"substate"`
}

type jsonClusterStatus ClusterStatus

func (r *ClusterStatus) UnmarshalJSON(data []byte) error {
	var res jsonClusterStatus
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterStatus
	} else {

		r.State = res.State

		r.Detail = res.Detail

		r.StateStartTime = res.StateStartTime

		r.Substate = res.Substate

	}
	return nil
}

// This object is used to assert a desired state where this ClusterStatus is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyClusterStatus *ClusterStatus = &ClusterStatus{empty: true}

func (r *ClusterStatus) Empty() bool {
	return r.empty
}

func (r *ClusterStatus) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterStatus) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterStatusHistory struct {
	empty          bool                              `json:"-"`
	State          *ClusterStatusHistoryStateEnum    `json:"state"`
	Detail         *string                           `json:"detail"`
	StateStartTime *string                           `json:"stateStartTime"`
	Substate       *ClusterStatusHistorySubstateEnum `json:"substate"`
}

type jsonClusterStatusHistory ClusterStatusHistory

func (r *ClusterStatusHistory) UnmarshalJSON(data []byte) error {
	var res jsonClusterStatusHistory
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterStatusHistory
	} else {

		r.State = res.State

		r.Detail = res.Detail

		r.StateStartTime = res.StateStartTime

		r.Substate = res.Substate

	}
	return nil
}

// This object is used to assert a desired state where this ClusterStatusHistory is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyClusterStatusHistory *ClusterStatusHistory = &ClusterStatusHistory{empty: true}

func (r *ClusterStatusHistory) Empty() bool {
	return r.empty
}

func (r *ClusterStatusHistory) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterStatusHistory) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterMetrics struct {
	empty       bool              `json:"-"`
	HdfsMetrics map[string]string `json:"hdfsMetrics"`
	YarnMetrics map[string]string `json:"yarnMetrics"`
}

type jsonClusterMetrics ClusterMetrics

func (r *ClusterMetrics) UnmarshalJSON(data []byte) error {
	var res jsonClusterMetrics
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterMetrics
	} else {

		r.HdfsMetrics = res.HdfsMetrics

		r.YarnMetrics = res.YarnMetrics

	}
	return nil
}

// This object is used to assert a desired state where this ClusterMetrics is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyClusterMetrics *ClusterMetrics = &ClusterMetrics{empty: true}

func (r *ClusterMetrics) Empty() bool {
	return r.empty
}

func (r *ClusterMetrics) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterMetrics) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *Cluster) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "dataproc",
		Type:    "Cluster",
		Version: "dataproc",
	}
}

func (r *Cluster) ID() (string, error) {
	if err := extractClusterFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":       dcl.ValueOrEmptyString(nr.Project),
		"name":          dcl.ValueOrEmptyString(nr.Name),
		"config":        dcl.ValueOrEmptyString(nr.Config),
		"labels":        dcl.ValueOrEmptyString(nr.Labels),
		"status":        dcl.ValueOrEmptyString(nr.Status),
		"statusHistory": dcl.ValueOrEmptyString(nr.StatusHistory),
		"clusterUuid":   dcl.ValueOrEmptyString(nr.ClusterUuid),
		"metrics":       dcl.ValueOrEmptyString(nr.Metrics),
		"location":      dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.Nprintf("projects/{{project}}/regions/{{location}}/clusters/{{name}}", params), nil
}

const ClusterMaxPage = -1

type ClusterList struct {
	Items []*Cluster

	nextToken string

	pageSize int32

	resource *Cluster
}

func (l *ClusterList) HasNext() bool {
	return l.nextToken != ""
}

func (l *ClusterList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listCluster(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListCluster(ctx context.Context, project, location string) (*ClusterList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListClusterWithMaxResults(ctx, project, location, ClusterMaxPage)

}

func (c *Client) ListClusterWithMaxResults(ctx context.Context, project, location string, pageSize int32) (*ClusterList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &Cluster{
		Project:  &project,
		Location: &location,
	}
	items, token, err := c.listCluster(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &ClusterList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetCluster(ctx context.Context, r *Cluster) (*Cluster, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractClusterFields(r)

	b, err := c.getClusterRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalCluster(b, c)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Location = r.Location
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeClusterNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractClusterFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteCluster(ctx context.Context, r *Cluster) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("Cluster resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting Cluster...")
	deleteOp := deleteClusterOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllCluster deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllCluster(ctx context.Context, project, location string, filter func(*Cluster) bool) error {
	listObj, err := c.ListCluster(ctx, project, location)
	if err != nil {
		return err
	}

	err = c.deleteAllCluster(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllCluster(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyCluster(ctx context.Context, rawDesired *Cluster, opts ...dcl.ApplyOption) (*Cluster, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *Cluster
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyClusterHelper(c, ctx, rawDesired, opts...)
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

func applyClusterHelper(c *Client, ctx context.Context, rawDesired *Cluster, opts ...dcl.ApplyOption) (*Cluster, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyCluster...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractClusterFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.clusterDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToClusterDiffs(c.Config, fieldDiffs, opts)
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
	var ops []clusterApiOperation
	if create {
		ops = append(ops, &createClusterOperation{})
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
	return applyClusterDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyClusterDiff(c *Client, ctx context.Context, desired *Cluster, rawDesired *Cluster, ops []clusterApiOperation, opts ...dcl.ApplyOption) (*Cluster, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetCluster(ctx, desired.urlNormalized())
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createClusterOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapCluster(r, c)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeClusterNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeClusterNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeClusterDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractClusterFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractClusterFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffCluster(c, newDesired, newState)
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

func (r *Cluster) GetPolicy(basePath string) (string, string, *bytes.Buffer, error) {
	u := r.getPolicyURL(basePath)
	body := &bytes.Buffer{}
	body.WriteString(fmt.Sprintf(`{"options":{"requestedPolicyVersion": %d}}`, r.IAMPolicyVersion()))
	return u, "POST", body, nil
}
