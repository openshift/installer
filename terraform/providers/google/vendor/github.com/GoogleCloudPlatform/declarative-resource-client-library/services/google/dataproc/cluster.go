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
	Project              *string                      `json:"project"`
	Name                 *string                      `json:"name"`
	Config               *ClusterConfig               `json:"config"`
	Labels               map[string]string            `json:"labels"`
	Status               *ClusterStatus               `json:"status"`
	StatusHistory        []ClusterStatusHistory       `json:"statusHistory"`
	ClusterUuid          *string                      `json:"clusterUuid"`
	Metrics              *ClusterMetrics              `json:"metrics"`
	Location             *string                      `json:"location"`
	VirtualClusterConfig *ClusterVirtualClusterConfig `json:"virtualClusterConfig"`
}

func (r *Cluster) String() string {
	return dcl.SprintResource(r)
}

// The enum ClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum.
type ClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum string

// ClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnumRef returns a *ClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum with the value of string s
// If the empty string is provided, nil is returned.
func ClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnumRef(s string) *ClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum {
	v := ClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum(s)
	return &v
}

func (v ClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum) Validate() error {
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
		Enum:  "ClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum ClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum.
type ClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum string

// ClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnumRef returns a *ClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func ClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnumRef(s string) *ClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum {
	v := ClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum(s)
	return &v
}

func (v ClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum) Validate() error {
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
		Enum:  "ClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum ClusterConfigMasterConfigPreemptibilityEnum.
type ClusterConfigMasterConfigPreemptibilityEnum string

// ClusterConfigMasterConfigPreemptibilityEnumRef returns a *ClusterConfigMasterConfigPreemptibilityEnum with the value of string s
// If the empty string is provided, nil is returned.
func ClusterConfigMasterConfigPreemptibilityEnumRef(s string) *ClusterConfigMasterConfigPreemptibilityEnum {
	v := ClusterConfigMasterConfigPreemptibilityEnum(s)
	return &v
}

func (v ClusterConfigMasterConfigPreemptibilityEnum) Validate() error {
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
		Enum:  "ClusterConfigMasterConfigPreemptibilityEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum ClusterConfigWorkerConfigPreemptibilityEnum.
type ClusterConfigWorkerConfigPreemptibilityEnum string

// ClusterConfigWorkerConfigPreemptibilityEnumRef returns a *ClusterConfigWorkerConfigPreemptibilityEnum with the value of string s
// If the empty string is provided, nil is returned.
func ClusterConfigWorkerConfigPreemptibilityEnumRef(s string) *ClusterConfigWorkerConfigPreemptibilityEnum {
	v := ClusterConfigWorkerConfigPreemptibilityEnum(s)
	return &v
}

func (v ClusterConfigWorkerConfigPreemptibilityEnum) Validate() error {
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
		Enum:  "ClusterConfigWorkerConfigPreemptibilityEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum ClusterConfigSecondaryWorkerConfigPreemptibilityEnum.
type ClusterConfigSecondaryWorkerConfigPreemptibilityEnum string

// ClusterConfigSecondaryWorkerConfigPreemptibilityEnumRef returns a *ClusterConfigSecondaryWorkerConfigPreemptibilityEnum with the value of string s
// If the empty string is provided, nil is returned.
func ClusterConfigSecondaryWorkerConfigPreemptibilityEnumRef(s string) *ClusterConfigSecondaryWorkerConfigPreemptibilityEnum {
	v := ClusterConfigSecondaryWorkerConfigPreemptibilityEnum(s)
	return &v
}

func (v ClusterConfigSecondaryWorkerConfigPreemptibilityEnum) Validate() error {
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
		Enum:  "ClusterConfigSecondaryWorkerConfigPreemptibilityEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum ClusterConfigSoftwareConfigOptionalComponentsEnum.
type ClusterConfigSoftwareConfigOptionalComponentsEnum string

// ClusterConfigSoftwareConfigOptionalComponentsEnumRef returns a *ClusterConfigSoftwareConfigOptionalComponentsEnum with the value of string s
// If the empty string is provided, nil is returned.
func ClusterConfigSoftwareConfigOptionalComponentsEnumRef(s string) *ClusterConfigSoftwareConfigOptionalComponentsEnum {
	v := ClusterConfigSoftwareConfigOptionalComponentsEnum(s)
	return &v
}

func (v ClusterConfigSoftwareConfigOptionalComponentsEnum) Validate() error {
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
		Enum:  "ClusterConfigSoftwareConfigOptionalComponentsEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum ClusterConfigDataprocMetricConfigMetricsMetricSourceEnum.
type ClusterConfigDataprocMetricConfigMetricsMetricSourceEnum string

// ClusterConfigDataprocMetricConfigMetricsMetricSourceEnumRef returns a *ClusterConfigDataprocMetricConfigMetricsMetricSourceEnum with the value of string s
// If the empty string is provided, nil is returned.
func ClusterConfigDataprocMetricConfigMetricsMetricSourceEnumRef(s string) *ClusterConfigDataprocMetricConfigMetricsMetricSourceEnum {
	v := ClusterConfigDataprocMetricConfigMetricsMetricSourceEnum(s)
	return &v
}

func (v ClusterConfigDataprocMetricConfigMetricsMetricSourceEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"METRIC_SOURCE_UNSPECIFIED", "MONITORING_AGENT_DEFAULTS", "HDFS", "SPARK", "YARN", "SPARK_HISTORY_SERVER", "HIVESERVER2"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "ClusterConfigDataprocMetricConfigMetricsMetricSourceEnum",
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

// The enum ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetRolesEnum.
type ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetRolesEnum string

// ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetRolesEnumRef returns a *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetRolesEnum with the value of string s
// If the empty string is provided, nil is returned.
func ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetRolesEnumRef(s string) *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetRolesEnum {
	v := ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetRolesEnum(s)
	return &v
}

func (v ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetRolesEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ROLE_UNSPECIFIED", "DEFAULT", "CONTROLLER", "SPARK_DRIVER", "SPARK_EXECUTOR"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetRolesEnum",
		Value: string(v),
		Valid: []string{},
	}
}

type ClusterConfig struct {
	empty                 bool                                 `json:"-"`
	StagingBucket         *string                              `json:"stagingBucket"`
	TempBucket            *string                              `json:"tempBucket"`
	GceClusterConfig      *ClusterConfigGceClusterConfig       `json:"gceClusterConfig"`
	MasterConfig          *ClusterConfigMasterConfig           `json:"masterConfig"`
	WorkerConfig          *ClusterConfigWorkerConfig           `json:"workerConfig"`
	SecondaryWorkerConfig *ClusterConfigSecondaryWorkerConfig  `json:"secondaryWorkerConfig"`
	SoftwareConfig        *ClusterConfigSoftwareConfig         `json:"softwareConfig"`
	InitializationActions []ClusterConfigInitializationActions `json:"initializationActions"`
	EncryptionConfig      *ClusterConfigEncryptionConfig       `json:"encryptionConfig"`
	AutoscalingConfig     *ClusterConfigAutoscalingConfig      `json:"autoscalingConfig"`
	SecurityConfig        *ClusterConfigSecurityConfig         `json:"securityConfig"`
	LifecycleConfig       *ClusterConfigLifecycleConfig        `json:"lifecycleConfig"`
	EndpointConfig        *ClusterConfigEndpointConfig         `json:"endpointConfig"`
	MetastoreConfig       *ClusterConfigMetastoreConfig        `json:"metastoreConfig"`
	DataprocMetricConfig  *ClusterConfigDataprocMetricConfig   `json:"dataprocMetricConfig"`
}

type jsonClusterConfig ClusterConfig

func (r *ClusterConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfig
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

		r.MetastoreConfig = res.MetastoreConfig

		r.DataprocMetricConfig = res.DataprocMetricConfig

	}
	return nil
}

// This object is used to assert a desired state where this ClusterConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfig *ClusterConfig = &ClusterConfig{empty: true}

func (r *ClusterConfig) Empty() bool {
	return r.empty
}

func (r *ClusterConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterConfigGceClusterConfig struct {
	empty                      bool                                                      `json:"-"`
	Zone                       *string                                                   `json:"zone"`
	Network                    *string                                                   `json:"network"`
	Subnetwork                 *string                                                   `json:"subnetwork"`
	InternalIPOnly             *bool                                                     `json:"internalIPOnly"`
	PrivateIPv6GoogleAccess    *ClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum `json:"privateIPv6GoogleAccess"`
	ServiceAccount             *string                                                   `json:"serviceAccount"`
	ServiceAccountScopes       []string                                                  `json:"serviceAccountScopes"`
	Tags                       []string                                                  `json:"tags"`
	Metadata                   map[string]string                                         `json:"metadata"`
	ReservationAffinity        *ClusterConfigGceClusterConfigReservationAffinity         `json:"reservationAffinity"`
	NodeGroupAffinity          *ClusterConfigGceClusterConfigNodeGroupAffinity           `json:"nodeGroupAffinity"`
	ShieldedInstanceConfig     *ClusterConfigGceClusterConfigShieldedInstanceConfig      `json:"shieldedInstanceConfig"`
	ConfidentialInstanceConfig *ClusterConfigGceClusterConfigConfidentialInstanceConfig  `json:"confidentialInstanceConfig"`
}

type jsonClusterConfigGceClusterConfig ClusterConfigGceClusterConfig

func (r *ClusterConfigGceClusterConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfigGceClusterConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfigGceClusterConfig
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

		r.ShieldedInstanceConfig = res.ShieldedInstanceConfig

		r.ConfidentialInstanceConfig = res.ConfidentialInstanceConfig

	}
	return nil
}

// This object is used to assert a desired state where this ClusterConfigGceClusterConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfigGceClusterConfig *ClusterConfigGceClusterConfig = &ClusterConfigGceClusterConfig{empty: true}

func (r *ClusterConfigGceClusterConfig) Empty() bool {
	return r.empty
}

func (r *ClusterConfigGceClusterConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfigGceClusterConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterConfigGceClusterConfigReservationAffinity struct {
	empty                  bool                                                                        `json:"-"`
	ConsumeReservationType *ClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum `json:"consumeReservationType"`
	Key                    *string                                                                     `json:"key"`
	Values                 []string                                                                    `json:"values"`
}

type jsonClusterConfigGceClusterConfigReservationAffinity ClusterConfigGceClusterConfigReservationAffinity

func (r *ClusterConfigGceClusterConfigReservationAffinity) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfigGceClusterConfigReservationAffinity
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfigGceClusterConfigReservationAffinity
	} else {

		r.ConsumeReservationType = res.ConsumeReservationType

		r.Key = res.Key

		r.Values = res.Values

	}
	return nil
}

// This object is used to assert a desired state where this ClusterConfigGceClusterConfigReservationAffinity is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfigGceClusterConfigReservationAffinity *ClusterConfigGceClusterConfigReservationAffinity = &ClusterConfigGceClusterConfigReservationAffinity{empty: true}

func (r *ClusterConfigGceClusterConfigReservationAffinity) Empty() bool {
	return r.empty
}

func (r *ClusterConfigGceClusterConfigReservationAffinity) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfigGceClusterConfigReservationAffinity) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterConfigGceClusterConfigNodeGroupAffinity struct {
	empty     bool    `json:"-"`
	NodeGroup *string `json:"nodeGroup"`
}

type jsonClusterConfigGceClusterConfigNodeGroupAffinity ClusterConfigGceClusterConfigNodeGroupAffinity

func (r *ClusterConfigGceClusterConfigNodeGroupAffinity) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfigGceClusterConfigNodeGroupAffinity
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfigGceClusterConfigNodeGroupAffinity
	} else {

		r.NodeGroup = res.NodeGroup

	}
	return nil
}

// This object is used to assert a desired state where this ClusterConfigGceClusterConfigNodeGroupAffinity is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfigGceClusterConfigNodeGroupAffinity *ClusterConfigGceClusterConfigNodeGroupAffinity = &ClusterConfigGceClusterConfigNodeGroupAffinity{empty: true}

func (r *ClusterConfigGceClusterConfigNodeGroupAffinity) Empty() bool {
	return r.empty
}

func (r *ClusterConfigGceClusterConfigNodeGroupAffinity) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfigGceClusterConfigNodeGroupAffinity) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterConfigGceClusterConfigShieldedInstanceConfig struct {
	empty                     bool  `json:"-"`
	EnableSecureBoot          *bool `json:"enableSecureBoot"`
	EnableVtpm                *bool `json:"enableVtpm"`
	EnableIntegrityMonitoring *bool `json:"enableIntegrityMonitoring"`
}

type jsonClusterConfigGceClusterConfigShieldedInstanceConfig ClusterConfigGceClusterConfigShieldedInstanceConfig

func (r *ClusterConfigGceClusterConfigShieldedInstanceConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfigGceClusterConfigShieldedInstanceConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfigGceClusterConfigShieldedInstanceConfig
	} else {

		r.EnableSecureBoot = res.EnableSecureBoot

		r.EnableVtpm = res.EnableVtpm

		r.EnableIntegrityMonitoring = res.EnableIntegrityMonitoring

	}
	return nil
}

// This object is used to assert a desired state where this ClusterConfigGceClusterConfigShieldedInstanceConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfigGceClusterConfigShieldedInstanceConfig *ClusterConfigGceClusterConfigShieldedInstanceConfig = &ClusterConfigGceClusterConfigShieldedInstanceConfig{empty: true}

func (r *ClusterConfigGceClusterConfigShieldedInstanceConfig) Empty() bool {
	return r.empty
}

func (r *ClusterConfigGceClusterConfigShieldedInstanceConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfigGceClusterConfigShieldedInstanceConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterConfigGceClusterConfigConfidentialInstanceConfig struct {
	empty                     bool  `json:"-"`
	EnableConfidentialCompute *bool `json:"enableConfidentialCompute"`
}

type jsonClusterConfigGceClusterConfigConfidentialInstanceConfig ClusterConfigGceClusterConfigConfidentialInstanceConfig

func (r *ClusterConfigGceClusterConfigConfidentialInstanceConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfigGceClusterConfigConfidentialInstanceConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfigGceClusterConfigConfidentialInstanceConfig
	} else {

		r.EnableConfidentialCompute = res.EnableConfidentialCompute

	}
	return nil
}

// This object is used to assert a desired state where this ClusterConfigGceClusterConfigConfidentialInstanceConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfigGceClusterConfigConfidentialInstanceConfig *ClusterConfigGceClusterConfigConfidentialInstanceConfig = &ClusterConfigGceClusterConfigConfidentialInstanceConfig{empty: true}

func (r *ClusterConfigGceClusterConfigConfidentialInstanceConfig) Empty() bool {
	return r.empty
}

func (r *ClusterConfigGceClusterConfigConfidentialInstanceConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfigGceClusterConfigConfidentialInstanceConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterConfigMasterConfig struct {
	empty              bool                                          `json:"-"`
	NumInstances       *int64                                        `json:"numInstances"`
	InstanceNames      []string                                      `json:"instanceNames"`
	Image              *string                                       `json:"image"`
	MachineType        *string                                       `json:"machineType"`
	DiskConfig         *ClusterConfigMasterConfigDiskConfig          `json:"diskConfig"`
	IsPreemptible      *bool                                         `json:"isPreemptible"`
	Preemptibility     *ClusterConfigMasterConfigPreemptibilityEnum  `json:"preemptibility"`
	ManagedGroupConfig *ClusterConfigMasterConfigManagedGroupConfig  `json:"managedGroupConfig"`
	Accelerators       []ClusterConfigMasterConfigAccelerators       `json:"accelerators"`
	MinCpuPlatform     *string                                       `json:"minCpuPlatform"`
	InstanceReferences []ClusterConfigMasterConfigInstanceReferences `json:"instanceReferences"`
}

type jsonClusterConfigMasterConfig ClusterConfigMasterConfig

func (r *ClusterConfigMasterConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfigMasterConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfigMasterConfig
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

		r.InstanceReferences = res.InstanceReferences

	}
	return nil
}

// This object is used to assert a desired state where this ClusterConfigMasterConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfigMasterConfig *ClusterConfigMasterConfig = &ClusterConfigMasterConfig{empty: true}

func (r *ClusterConfigMasterConfig) Empty() bool {
	return r.empty
}

func (r *ClusterConfigMasterConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfigMasterConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterConfigMasterConfigDiskConfig struct {
	empty             bool    `json:"-"`
	BootDiskType      *string `json:"bootDiskType"`
	BootDiskSizeGb    *int64  `json:"bootDiskSizeGb"`
	NumLocalSsds      *int64  `json:"numLocalSsds"`
	LocalSsdInterface *string `json:"localSsdInterface"`
}

type jsonClusterConfigMasterConfigDiskConfig ClusterConfigMasterConfigDiskConfig

func (r *ClusterConfigMasterConfigDiskConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfigMasterConfigDiskConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfigMasterConfigDiskConfig
	} else {

		r.BootDiskType = res.BootDiskType

		r.BootDiskSizeGb = res.BootDiskSizeGb

		r.NumLocalSsds = res.NumLocalSsds

		r.LocalSsdInterface = res.LocalSsdInterface

	}
	return nil
}

// This object is used to assert a desired state where this ClusterConfigMasterConfigDiskConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfigMasterConfigDiskConfig *ClusterConfigMasterConfigDiskConfig = &ClusterConfigMasterConfigDiskConfig{empty: true}

func (r *ClusterConfigMasterConfigDiskConfig) Empty() bool {
	return r.empty
}

func (r *ClusterConfigMasterConfigDiskConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfigMasterConfigDiskConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterConfigMasterConfigManagedGroupConfig struct {
	empty                    bool    `json:"-"`
	InstanceTemplateName     *string `json:"instanceTemplateName"`
	InstanceGroupManagerName *string `json:"instanceGroupManagerName"`
}

type jsonClusterConfigMasterConfigManagedGroupConfig ClusterConfigMasterConfigManagedGroupConfig

func (r *ClusterConfigMasterConfigManagedGroupConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfigMasterConfigManagedGroupConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfigMasterConfigManagedGroupConfig
	} else {

		r.InstanceTemplateName = res.InstanceTemplateName

		r.InstanceGroupManagerName = res.InstanceGroupManagerName

	}
	return nil
}

// This object is used to assert a desired state where this ClusterConfigMasterConfigManagedGroupConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfigMasterConfigManagedGroupConfig *ClusterConfigMasterConfigManagedGroupConfig = &ClusterConfigMasterConfigManagedGroupConfig{empty: true}

func (r *ClusterConfigMasterConfigManagedGroupConfig) Empty() bool {
	return r.empty
}

func (r *ClusterConfigMasterConfigManagedGroupConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfigMasterConfigManagedGroupConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterConfigMasterConfigAccelerators struct {
	empty            bool    `json:"-"`
	AcceleratorType  *string `json:"acceleratorType"`
	AcceleratorCount *int64  `json:"acceleratorCount"`
}

type jsonClusterConfigMasterConfigAccelerators ClusterConfigMasterConfigAccelerators

func (r *ClusterConfigMasterConfigAccelerators) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfigMasterConfigAccelerators
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfigMasterConfigAccelerators
	} else {

		r.AcceleratorType = res.AcceleratorType

		r.AcceleratorCount = res.AcceleratorCount

	}
	return nil
}

// This object is used to assert a desired state where this ClusterConfigMasterConfigAccelerators is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfigMasterConfigAccelerators *ClusterConfigMasterConfigAccelerators = &ClusterConfigMasterConfigAccelerators{empty: true}

func (r *ClusterConfigMasterConfigAccelerators) Empty() bool {
	return r.empty
}

func (r *ClusterConfigMasterConfigAccelerators) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfigMasterConfigAccelerators) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterConfigMasterConfigInstanceReferences struct {
	empty          bool    `json:"-"`
	InstanceName   *string `json:"instanceName"`
	InstanceId     *string `json:"instanceId"`
	PublicKey      *string `json:"publicKey"`
	PublicEciesKey *string `json:"publicEciesKey"`
}

type jsonClusterConfigMasterConfigInstanceReferences ClusterConfigMasterConfigInstanceReferences

func (r *ClusterConfigMasterConfigInstanceReferences) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfigMasterConfigInstanceReferences
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfigMasterConfigInstanceReferences
	} else {

		r.InstanceName = res.InstanceName

		r.InstanceId = res.InstanceId

		r.PublicKey = res.PublicKey

		r.PublicEciesKey = res.PublicEciesKey

	}
	return nil
}

// This object is used to assert a desired state where this ClusterConfigMasterConfigInstanceReferences is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfigMasterConfigInstanceReferences *ClusterConfigMasterConfigInstanceReferences = &ClusterConfigMasterConfigInstanceReferences{empty: true}

func (r *ClusterConfigMasterConfigInstanceReferences) Empty() bool {
	return r.empty
}

func (r *ClusterConfigMasterConfigInstanceReferences) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfigMasterConfigInstanceReferences) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterConfigWorkerConfig struct {
	empty              bool                                          `json:"-"`
	NumInstances       *int64                                        `json:"numInstances"`
	InstanceNames      []string                                      `json:"instanceNames"`
	Image              *string                                       `json:"image"`
	MachineType        *string                                       `json:"machineType"`
	DiskConfig         *ClusterConfigWorkerConfigDiskConfig          `json:"diskConfig"`
	IsPreemptible      *bool                                         `json:"isPreemptible"`
	Preemptibility     *ClusterConfigWorkerConfigPreemptibilityEnum  `json:"preemptibility"`
	ManagedGroupConfig *ClusterConfigWorkerConfigManagedGroupConfig  `json:"managedGroupConfig"`
	Accelerators       []ClusterConfigWorkerConfigAccelerators       `json:"accelerators"`
	MinCpuPlatform     *string                                       `json:"minCpuPlatform"`
	InstanceReferences []ClusterConfigWorkerConfigInstanceReferences `json:"instanceReferences"`
}

type jsonClusterConfigWorkerConfig ClusterConfigWorkerConfig

func (r *ClusterConfigWorkerConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfigWorkerConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfigWorkerConfig
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

		r.InstanceReferences = res.InstanceReferences

	}
	return nil
}

// This object is used to assert a desired state where this ClusterConfigWorkerConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfigWorkerConfig *ClusterConfigWorkerConfig = &ClusterConfigWorkerConfig{empty: true}

func (r *ClusterConfigWorkerConfig) Empty() bool {
	return r.empty
}

func (r *ClusterConfigWorkerConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfigWorkerConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterConfigWorkerConfigDiskConfig struct {
	empty             bool    `json:"-"`
	BootDiskType      *string `json:"bootDiskType"`
	BootDiskSizeGb    *int64  `json:"bootDiskSizeGb"`
	NumLocalSsds      *int64  `json:"numLocalSsds"`
	LocalSsdInterface *string `json:"localSsdInterface"`
}

type jsonClusterConfigWorkerConfigDiskConfig ClusterConfigWorkerConfigDiskConfig

func (r *ClusterConfigWorkerConfigDiskConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfigWorkerConfigDiskConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfigWorkerConfigDiskConfig
	} else {

		r.BootDiskType = res.BootDiskType

		r.BootDiskSizeGb = res.BootDiskSizeGb

		r.NumLocalSsds = res.NumLocalSsds

		r.LocalSsdInterface = res.LocalSsdInterface

	}
	return nil
}

// This object is used to assert a desired state where this ClusterConfigWorkerConfigDiskConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfigWorkerConfigDiskConfig *ClusterConfigWorkerConfigDiskConfig = &ClusterConfigWorkerConfigDiskConfig{empty: true}

func (r *ClusterConfigWorkerConfigDiskConfig) Empty() bool {
	return r.empty
}

func (r *ClusterConfigWorkerConfigDiskConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfigWorkerConfigDiskConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterConfigWorkerConfigManagedGroupConfig struct {
	empty                    bool    `json:"-"`
	InstanceTemplateName     *string `json:"instanceTemplateName"`
	InstanceGroupManagerName *string `json:"instanceGroupManagerName"`
}

type jsonClusterConfigWorkerConfigManagedGroupConfig ClusterConfigWorkerConfigManagedGroupConfig

func (r *ClusterConfigWorkerConfigManagedGroupConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfigWorkerConfigManagedGroupConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfigWorkerConfigManagedGroupConfig
	} else {

		r.InstanceTemplateName = res.InstanceTemplateName

		r.InstanceGroupManagerName = res.InstanceGroupManagerName

	}
	return nil
}

// This object is used to assert a desired state where this ClusterConfigWorkerConfigManagedGroupConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfigWorkerConfigManagedGroupConfig *ClusterConfigWorkerConfigManagedGroupConfig = &ClusterConfigWorkerConfigManagedGroupConfig{empty: true}

func (r *ClusterConfigWorkerConfigManagedGroupConfig) Empty() bool {
	return r.empty
}

func (r *ClusterConfigWorkerConfigManagedGroupConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfigWorkerConfigManagedGroupConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterConfigWorkerConfigAccelerators struct {
	empty            bool    `json:"-"`
	AcceleratorType  *string `json:"acceleratorType"`
	AcceleratorCount *int64  `json:"acceleratorCount"`
}

type jsonClusterConfigWorkerConfigAccelerators ClusterConfigWorkerConfigAccelerators

func (r *ClusterConfigWorkerConfigAccelerators) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfigWorkerConfigAccelerators
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfigWorkerConfigAccelerators
	} else {

		r.AcceleratorType = res.AcceleratorType

		r.AcceleratorCount = res.AcceleratorCount

	}
	return nil
}

// This object is used to assert a desired state where this ClusterConfigWorkerConfigAccelerators is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfigWorkerConfigAccelerators *ClusterConfigWorkerConfigAccelerators = &ClusterConfigWorkerConfigAccelerators{empty: true}

func (r *ClusterConfigWorkerConfigAccelerators) Empty() bool {
	return r.empty
}

func (r *ClusterConfigWorkerConfigAccelerators) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfigWorkerConfigAccelerators) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterConfigWorkerConfigInstanceReferences struct {
	empty          bool    `json:"-"`
	InstanceName   *string `json:"instanceName"`
	InstanceId     *string `json:"instanceId"`
	PublicKey      *string `json:"publicKey"`
	PublicEciesKey *string `json:"publicEciesKey"`
}

type jsonClusterConfigWorkerConfigInstanceReferences ClusterConfigWorkerConfigInstanceReferences

func (r *ClusterConfigWorkerConfigInstanceReferences) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfigWorkerConfigInstanceReferences
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfigWorkerConfigInstanceReferences
	} else {

		r.InstanceName = res.InstanceName

		r.InstanceId = res.InstanceId

		r.PublicKey = res.PublicKey

		r.PublicEciesKey = res.PublicEciesKey

	}
	return nil
}

// This object is used to assert a desired state where this ClusterConfigWorkerConfigInstanceReferences is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfigWorkerConfigInstanceReferences *ClusterConfigWorkerConfigInstanceReferences = &ClusterConfigWorkerConfigInstanceReferences{empty: true}

func (r *ClusterConfigWorkerConfigInstanceReferences) Empty() bool {
	return r.empty
}

func (r *ClusterConfigWorkerConfigInstanceReferences) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfigWorkerConfigInstanceReferences) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterConfigSecondaryWorkerConfig struct {
	empty              bool                                                   `json:"-"`
	NumInstances       *int64                                                 `json:"numInstances"`
	InstanceNames      []string                                               `json:"instanceNames"`
	Image              *string                                                `json:"image"`
	MachineType        *string                                                `json:"machineType"`
	DiskConfig         *ClusterConfigSecondaryWorkerConfigDiskConfig          `json:"diskConfig"`
	IsPreemptible      *bool                                                  `json:"isPreemptible"`
	Preemptibility     *ClusterConfigSecondaryWorkerConfigPreemptibilityEnum  `json:"preemptibility"`
	ManagedGroupConfig *ClusterConfigSecondaryWorkerConfigManagedGroupConfig  `json:"managedGroupConfig"`
	Accelerators       []ClusterConfigSecondaryWorkerConfigAccelerators       `json:"accelerators"`
	MinCpuPlatform     *string                                                `json:"minCpuPlatform"`
	InstanceReferences []ClusterConfigSecondaryWorkerConfigInstanceReferences `json:"instanceReferences"`
}

type jsonClusterConfigSecondaryWorkerConfig ClusterConfigSecondaryWorkerConfig

func (r *ClusterConfigSecondaryWorkerConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfigSecondaryWorkerConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfigSecondaryWorkerConfig
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

		r.InstanceReferences = res.InstanceReferences

	}
	return nil
}

// This object is used to assert a desired state where this ClusterConfigSecondaryWorkerConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfigSecondaryWorkerConfig *ClusterConfigSecondaryWorkerConfig = &ClusterConfigSecondaryWorkerConfig{empty: true}

func (r *ClusterConfigSecondaryWorkerConfig) Empty() bool {
	return r.empty
}

func (r *ClusterConfigSecondaryWorkerConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfigSecondaryWorkerConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterConfigSecondaryWorkerConfigDiskConfig struct {
	empty             bool    `json:"-"`
	BootDiskType      *string `json:"bootDiskType"`
	BootDiskSizeGb    *int64  `json:"bootDiskSizeGb"`
	NumLocalSsds      *int64  `json:"numLocalSsds"`
	LocalSsdInterface *string `json:"localSsdInterface"`
}

type jsonClusterConfigSecondaryWorkerConfigDiskConfig ClusterConfigSecondaryWorkerConfigDiskConfig

func (r *ClusterConfigSecondaryWorkerConfigDiskConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfigSecondaryWorkerConfigDiskConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfigSecondaryWorkerConfigDiskConfig
	} else {

		r.BootDiskType = res.BootDiskType

		r.BootDiskSizeGb = res.BootDiskSizeGb

		r.NumLocalSsds = res.NumLocalSsds

		r.LocalSsdInterface = res.LocalSsdInterface

	}
	return nil
}

// This object is used to assert a desired state where this ClusterConfigSecondaryWorkerConfigDiskConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfigSecondaryWorkerConfigDiskConfig *ClusterConfigSecondaryWorkerConfigDiskConfig = &ClusterConfigSecondaryWorkerConfigDiskConfig{empty: true}

func (r *ClusterConfigSecondaryWorkerConfigDiskConfig) Empty() bool {
	return r.empty
}

func (r *ClusterConfigSecondaryWorkerConfigDiskConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfigSecondaryWorkerConfigDiskConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterConfigSecondaryWorkerConfigManagedGroupConfig struct {
	empty                    bool    `json:"-"`
	InstanceTemplateName     *string `json:"instanceTemplateName"`
	InstanceGroupManagerName *string `json:"instanceGroupManagerName"`
}

type jsonClusterConfigSecondaryWorkerConfigManagedGroupConfig ClusterConfigSecondaryWorkerConfigManagedGroupConfig

func (r *ClusterConfigSecondaryWorkerConfigManagedGroupConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfigSecondaryWorkerConfigManagedGroupConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfigSecondaryWorkerConfigManagedGroupConfig
	} else {

		r.InstanceTemplateName = res.InstanceTemplateName

		r.InstanceGroupManagerName = res.InstanceGroupManagerName

	}
	return nil
}

// This object is used to assert a desired state where this ClusterConfigSecondaryWorkerConfigManagedGroupConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfigSecondaryWorkerConfigManagedGroupConfig *ClusterConfigSecondaryWorkerConfigManagedGroupConfig = &ClusterConfigSecondaryWorkerConfigManagedGroupConfig{empty: true}

func (r *ClusterConfigSecondaryWorkerConfigManagedGroupConfig) Empty() bool {
	return r.empty
}

func (r *ClusterConfigSecondaryWorkerConfigManagedGroupConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfigSecondaryWorkerConfigManagedGroupConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterConfigSecondaryWorkerConfigAccelerators struct {
	empty            bool    `json:"-"`
	AcceleratorType  *string `json:"acceleratorType"`
	AcceleratorCount *int64  `json:"acceleratorCount"`
}

type jsonClusterConfigSecondaryWorkerConfigAccelerators ClusterConfigSecondaryWorkerConfigAccelerators

func (r *ClusterConfigSecondaryWorkerConfigAccelerators) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfigSecondaryWorkerConfigAccelerators
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfigSecondaryWorkerConfigAccelerators
	} else {

		r.AcceleratorType = res.AcceleratorType

		r.AcceleratorCount = res.AcceleratorCount

	}
	return nil
}

// This object is used to assert a desired state where this ClusterConfigSecondaryWorkerConfigAccelerators is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfigSecondaryWorkerConfigAccelerators *ClusterConfigSecondaryWorkerConfigAccelerators = &ClusterConfigSecondaryWorkerConfigAccelerators{empty: true}

func (r *ClusterConfigSecondaryWorkerConfigAccelerators) Empty() bool {
	return r.empty
}

func (r *ClusterConfigSecondaryWorkerConfigAccelerators) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfigSecondaryWorkerConfigAccelerators) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterConfigSecondaryWorkerConfigInstanceReferences struct {
	empty          bool    `json:"-"`
	InstanceName   *string `json:"instanceName"`
	InstanceId     *string `json:"instanceId"`
	PublicKey      *string `json:"publicKey"`
	PublicEciesKey *string `json:"publicEciesKey"`
}

type jsonClusterConfigSecondaryWorkerConfigInstanceReferences ClusterConfigSecondaryWorkerConfigInstanceReferences

func (r *ClusterConfigSecondaryWorkerConfigInstanceReferences) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfigSecondaryWorkerConfigInstanceReferences
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfigSecondaryWorkerConfigInstanceReferences
	} else {

		r.InstanceName = res.InstanceName

		r.InstanceId = res.InstanceId

		r.PublicKey = res.PublicKey

		r.PublicEciesKey = res.PublicEciesKey

	}
	return nil
}

// This object is used to assert a desired state where this ClusterConfigSecondaryWorkerConfigInstanceReferences is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfigSecondaryWorkerConfigInstanceReferences *ClusterConfigSecondaryWorkerConfigInstanceReferences = &ClusterConfigSecondaryWorkerConfigInstanceReferences{empty: true}

func (r *ClusterConfigSecondaryWorkerConfigInstanceReferences) Empty() bool {
	return r.empty
}

func (r *ClusterConfigSecondaryWorkerConfigInstanceReferences) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfigSecondaryWorkerConfigInstanceReferences) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterConfigSoftwareConfig struct {
	empty              bool                                                `json:"-"`
	ImageVersion       *string                                             `json:"imageVersion"`
	Properties         map[string]string                                   `json:"properties"`
	OptionalComponents []ClusterConfigSoftwareConfigOptionalComponentsEnum `json:"optionalComponents"`
}

type jsonClusterConfigSoftwareConfig ClusterConfigSoftwareConfig

func (r *ClusterConfigSoftwareConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfigSoftwareConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfigSoftwareConfig
	} else {

		r.ImageVersion = res.ImageVersion

		r.Properties = res.Properties

		r.OptionalComponents = res.OptionalComponents

	}
	return nil
}

// This object is used to assert a desired state where this ClusterConfigSoftwareConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfigSoftwareConfig *ClusterConfigSoftwareConfig = &ClusterConfigSoftwareConfig{empty: true}

func (r *ClusterConfigSoftwareConfig) Empty() bool {
	return r.empty
}

func (r *ClusterConfigSoftwareConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfigSoftwareConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterConfigInitializationActions struct {
	empty            bool    `json:"-"`
	ExecutableFile   *string `json:"executableFile"`
	ExecutionTimeout *string `json:"executionTimeout"`
}

type jsonClusterConfigInitializationActions ClusterConfigInitializationActions

func (r *ClusterConfigInitializationActions) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfigInitializationActions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfigInitializationActions
	} else {

		r.ExecutableFile = res.ExecutableFile

		r.ExecutionTimeout = res.ExecutionTimeout

	}
	return nil
}

// This object is used to assert a desired state where this ClusterConfigInitializationActions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfigInitializationActions *ClusterConfigInitializationActions = &ClusterConfigInitializationActions{empty: true}

func (r *ClusterConfigInitializationActions) Empty() bool {
	return r.empty
}

func (r *ClusterConfigInitializationActions) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfigInitializationActions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterConfigEncryptionConfig struct {
	empty           bool    `json:"-"`
	GcePdKmsKeyName *string `json:"gcePdKmsKeyName"`
}

type jsonClusterConfigEncryptionConfig ClusterConfigEncryptionConfig

func (r *ClusterConfigEncryptionConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfigEncryptionConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfigEncryptionConfig
	} else {

		r.GcePdKmsKeyName = res.GcePdKmsKeyName

	}
	return nil
}

// This object is used to assert a desired state where this ClusterConfigEncryptionConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfigEncryptionConfig *ClusterConfigEncryptionConfig = &ClusterConfigEncryptionConfig{empty: true}

func (r *ClusterConfigEncryptionConfig) Empty() bool {
	return r.empty
}

func (r *ClusterConfigEncryptionConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfigEncryptionConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterConfigAutoscalingConfig struct {
	empty  bool    `json:"-"`
	Policy *string `json:"policy"`
}

type jsonClusterConfigAutoscalingConfig ClusterConfigAutoscalingConfig

func (r *ClusterConfigAutoscalingConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfigAutoscalingConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfigAutoscalingConfig
	} else {

		r.Policy = res.Policy

	}
	return nil
}

// This object is used to assert a desired state where this ClusterConfigAutoscalingConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfigAutoscalingConfig *ClusterConfigAutoscalingConfig = &ClusterConfigAutoscalingConfig{empty: true}

func (r *ClusterConfigAutoscalingConfig) Empty() bool {
	return r.empty
}

func (r *ClusterConfigAutoscalingConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfigAutoscalingConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterConfigSecurityConfig struct {
	empty          bool                                       `json:"-"`
	KerberosConfig *ClusterConfigSecurityConfigKerberosConfig `json:"kerberosConfig"`
	IdentityConfig *ClusterConfigSecurityConfigIdentityConfig `json:"identityConfig"`
}

type jsonClusterConfigSecurityConfig ClusterConfigSecurityConfig

func (r *ClusterConfigSecurityConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfigSecurityConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfigSecurityConfig
	} else {

		r.KerberosConfig = res.KerberosConfig

		r.IdentityConfig = res.IdentityConfig

	}
	return nil
}

// This object is used to assert a desired state where this ClusterConfigSecurityConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfigSecurityConfig *ClusterConfigSecurityConfig = &ClusterConfigSecurityConfig{empty: true}

func (r *ClusterConfigSecurityConfig) Empty() bool {
	return r.empty
}

func (r *ClusterConfigSecurityConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfigSecurityConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterConfigSecurityConfigKerberosConfig struct {
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

type jsonClusterConfigSecurityConfigKerberosConfig ClusterConfigSecurityConfigKerberosConfig

func (r *ClusterConfigSecurityConfigKerberosConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfigSecurityConfigKerberosConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfigSecurityConfigKerberosConfig
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

// This object is used to assert a desired state where this ClusterConfigSecurityConfigKerberosConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfigSecurityConfigKerberosConfig *ClusterConfigSecurityConfigKerberosConfig = &ClusterConfigSecurityConfigKerberosConfig{empty: true}

func (r *ClusterConfigSecurityConfigKerberosConfig) Empty() bool {
	return r.empty
}

func (r *ClusterConfigSecurityConfigKerberosConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfigSecurityConfigKerberosConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterConfigSecurityConfigIdentityConfig struct {
	empty                     bool              `json:"-"`
	UserServiceAccountMapping map[string]string `json:"userServiceAccountMapping"`
}

type jsonClusterConfigSecurityConfigIdentityConfig ClusterConfigSecurityConfigIdentityConfig

func (r *ClusterConfigSecurityConfigIdentityConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfigSecurityConfigIdentityConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfigSecurityConfigIdentityConfig
	} else {

		r.UserServiceAccountMapping = res.UserServiceAccountMapping

	}
	return nil
}

// This object is used to assert a desired state where this ClusterConfigSecurityConfigIdentityConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfigSecurityConfigIdentityConfig *ClusterConfigSecurityConfigIdentityConfig = &ClusterConfigSecurityConfigIdentityConfig{empty: true}

func (r *ClusterConfigSecurityConfigIdentityConfig) Empty() bool {
	return r.empty
}

func (r *ClusterConfigSecurityConfigIdentityConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfigSecurityConfigIdentityConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterConfigLifecycleConfig struct {
	empty          bool    `json:"-"`
	IdleDeleteTtl  *string `json:"idleDeleteTtl"`
	AutoDeleteTime *string `json:"autoDeleteTime"`
	AutoDeleteTtl  *string `json:"autoDeleteTtl"`
	IdleStartTime  *string `json:"idleStartTime"`
}

type jsonClusterConfigLifecycleConfig ClusterConfigLifecycleConfig

func (r *ClusterConfigLifecycleConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfigLifecycleConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfigLifecycleConfig
	} else {

		r.IdleDeleteTtl = res.IdleDeleteTtl

		r.AutoDeleteTime = res.AutoDeleteTime

		r.AutoDeleteTtl = res.AutoDeleteTtl

		r.IdleStartTime = res.IdleStartTime

	}
	return nil
}

// This object is used to assert a desired state where this ClusterConfigLifecycleConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfigLifecycleConfig *ClusterConfigLifecycleConfig = &ClusterConfigLifecycleConfig{empty: true}

func (r *ClusterConfigLifecycleConfig) Empty() bool {
	return r.empty
}

func (r *ClusterConfigLifecycleConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfigLifecycleConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterConfigEndpointConfig struct {
	empty                bool              `json:"-"`
	HttpPorts            map[string]string `json:"httpPorts"`
	EnableHttpPortAccess *bool             `json:"enableHttpPortAccess"`
}

type jsonClusterConfigEndpointConfig ClusterConfigEndpointConfig

func (r *ClusterConfigEndpointConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfigEndpointConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfigEndpointConfig
	} else {

		r.HttpPorts = res.HttpPorts

		r.EnableHttpPortAccess = res.EnableHttpPortAccess

	}
	return nil
}

// This object is used to assert a desired state where this ClusterConfigEndpointConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfigEndpointConfig *ClusterConfigEndpointConfig = &ClusterConfigEndpointConfig{empty: true}

func (r *ClusterConfigEndpointConfig) Empty() bool {
	return r.empty
}

func (r *ClusterConfigEndpointConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfigEndpointConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterConfigMetastoreConfig struct {
	empty                    bool    `json:"-"`
	DataprocMetastoreService *string `json:"dataprocMetastoreService"`
}

type jsonClusterConfigMetastoreConfig ClusterConfigMetastoreConfig

func (r *ClusterConfigMetastoreConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfigMetastoreConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfigMetastoreConfig
	} else {

		r.DataprocMetastoreService = res.DataprocMetastoreService

	}
	return nil
}

// This object is used to assert a desired state where this ClusterConfigMetastoreConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfigMetastoreConfig *ClusterConfigMetastoreConfig = &ClusterConfigMetastoreConfig{empty: true}

func (r *ClusterConfigMetastoreConfig) Empty() bool {
	return r.empty
}

func (r *ClusterConfigMetastoreConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfigMetastoreConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterConfigDataprocMetricConfig struct {
	empty   bool                                       `json:"-"`
	Metrics []ClusterConfigDataprocMetricConfigMetrics `json:"metrics"`
}

type jsonClusterConfigDataprocMetricConfig ClusterConfigDataprocMetricConfig

func (r *ClusterConfigDataprocMetricConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfigDataprocMetricConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfigDataprocMetricConfig
	} else {

		r.Metrics = res.Metrics

	}
	return nil
}

// This object is used to assert a desired state where this ClusterConfigDataprocMetricConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfigDataprocMetricConfig *ClusterConfigDataprocMetricConfig = &ClusterConfigDataprocMetricConfig{empty: true}

func (r *ClusterConfigDataprocMetricConfig) Empty() bool {
	return r.empty
}

func (r *ClusterConfigDataprocMetricConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfigDataprocMetricConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterConfigDataprocMetricConfigMetrics struct {
	empty           bool                                                      `json:"-"`
	MetricSource    *ClusterConfigDataprocMetricConfigMetricsMetricSourceEnum `json:"metricSource"`
	MetricOverrides []string                                                  `json:"metricOverrides"`
}

type jsonClusterConfigDataprocMetricConfigMetrics ClusterConfigDataprocMetricConfigMetrics

func (r *ClusterConfigDataprocMetricConfigMetrics) UnmarshalJSON(data []byte) error {
	var res jsonClusterConfigDataprocMetricConfigMetrics
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterConfigDataprocMetricConfigMetrics
	} else {

		r.MetricSource = res.MetricSource

		r.MetricOverrides = res.MetricOverrides

	}
	return nil
}

// This object is used to assert a desired state where this ClusterConfigDataprocMetricConfigMetrics is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterConfigDataprocMetricConfigMetrics *ClusterConfigDataprocMetricConfigMetrics = &ClusterConfigDataprocMetricConfigMetrics{empty: true}

func (r *ClusterConfigDataprocMetricConfigMetrics) Empty() bool {
	return r.empty
}

func (r *ClusterConfigDataprocMetricConfigMetrics) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterConfigDataprocMetricConfigMetrics) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
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
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
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
	hash := sha256.Sum256([]byte(r.String()))
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
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
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
	hash := sha256.Sum256([]byte(r.String()))
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
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
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
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterVirtualClusterConfig struct {
	empty                   bool                                                `json:"-"`
	StagingBucket           *string                                             `json:"stagingBucket"`
	KubernetesClusterConfig *ClusterVirtualClusterConfigKubernetesClusterConfig `json:"kubernetesClusterConfig"`
	AuxiliaryServicesConfig *ClusterVirtualClusterConfigAuxiliaryServicesConfig `json:"auxiliaryServicesConfig"`
}

type jsonClusterVirtualClusterConfig ClusterVirtualClusterConfig

func (r *ClusterVirtualClusterConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterVirtualClusterConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterVirtualClusterConfig
	} else {

		r.StagingBucket = res.StagingBucket

		r.KubernetesClusterConfig = res.KubernetesClusterConfig

		r.AuxiliaryServicesConfig = res.AuxiliaryServicesConfig

	}
	return nil
}

// This object is used to assert a desired state where this ClusterVirtualClusterConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterVirtualClusterConfig *ClusterVirtualClusterConfig = &ClusterVirtualClusterConfig{empty: true}

func (r *ClusterVirtualClusterConfig) Empty() bool {
	return r.empty
}

func (r *ClusterVirtualClusterConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterVirtualClusterConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterVirtualClusterConfigKubernetesClusterConfig struct {
	empty                    bool                                                                        `json:"-"`
	KubernetesNamespace      *string                                                                     `json:"kubernetesNamespace"`
	GkeClusterConfig         *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig         `json:"gkeClusterConfig"`
	KubernetesSoftwareConfig *ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig `json:"kubernetesSoftwareConfig"`
}

type jsonClusterVirtualClusterConfigKubernetesClusterConfig ClusterVirtualClusterConfigKubernetesClusterConfig

func (r *ClusterVirtualClusterConfigKubernetesClusterConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterVirtualClusterConfigKubernetesClusterConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterVirtualClusterConfigKubernetesClusterConfig
	} else {

		r.KubernetesNamespace = res.KubernetesNamespace

		r.GkeClusterConfig = res.GkeClusterConfig

		r.KubernetesSoftwareConfig = res.KubernetesSoftwareConfig

	}
	return nil
}

// This object is used to assert a desired state where this ClusterVirtualClusterConfigKubernetesClusterConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterVirtualClusterConfigKubernetesClusterConfig *ClusterVirtualClusterConfigKubernetesClusterConfig = &ClusterVirtualClusterConfigKubernetesClusterConfig{empty: true}

func (r *ClusterVirtualClusterConfigKubernetesClusterConfig) Empty() bool {
	return r.empty
}

func (r *ClusterVirtualClusterConfigKubernetesClusterConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterVirtualClusterConfigKubernetesClusterConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig struct {
	empty            bool                                                                               `json:"-"`
	GkeClusterTarget *string                                                                            `json:"gkeClusterTarget"`
	NodePoolTarget   []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget `json:"nodePoolTarget"`
}

type jsonClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig

func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig
	} else {

		r.GkeClusterTarget = res.GkeClusterTarget

		r.NodePoolTarget = res.NodePoolTarget

	}
	return nil
}

// This object is used to assert a desired state where this ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig = &ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig{empty: true}

func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig) Empty() bool {
	return r.empty
}

func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget struct {
	empty          bool                                                                                            `json:"-"`
	NodePool       *string                                                                                         `json:"nodePool"`
	Roles          []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetRolesEnum     `json:"roles"`
	NodePoolConfig *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig `json:"nodePoolConfig"`
}

type jsonClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget

func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget) UnmarshalJSON(data []byte) error {
	var res jsonClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget
	} else {

		r.NodePool = res.NodePool

		r.Roles = res.Roles

		r.NodePoolConfig = res.NodePoolConfig

	}
	return nil
}

// This object is used to assert a desired state where this ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget = &ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget{empty: true}

func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget) Empty() bool {
	return r.empty
}

func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig struct {
	empty       bool                                                                                                       `json:"-"`
	Config      *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig      `json:"config"`
	Locations   []string                                                                                                   `json:"locations"`
	Autoscaling *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling `json:"autoscaling"`
}

type jsonClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig

func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig
	} else {

		r.Config = res.Config

		r.Locations = res.Locations

		r.Autoscaling = res.Autoscaling

	}
	return nil
}

// This object is used to assert a desired state where this ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig = &ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig{empty: true}

func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig) Empty() bool {
	return r.empty
}

func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig struct {
	empty                  bool                                                                                                                        `json:"-"`
	MachineType            *string                                                                                                                     `json:"machineType"`
	LocalSsdCount          *int64                                                                                                                      `json:"localSsdCount"`
	Preemptible            *bool                                                                                                                       `json:"preemptible"`
	Accelerators           []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators          `json:"accelerators"`
	MinCpuPlatform         *string                                                                                                                     `json:"minCpuPlatform"`
	BootDiskKmsKey         *string                                                                                                                     `json:"bootDiskKmsKey"`
	EphemeralStorageConfig *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig `json:"ephemeralStorageConfig"`
	Spot                   *bool                                                                                                                       `json:"spot"`
}

type jsonClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig

func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig
	} else {

		r.MachineType = res.MachineType

		r.LocalSsdCount = res.LocalSsdCount

		r.Preemptible = res.Preemptible

		r.Accelerators = res.Accelerators

		r.MinCpuPlatform = res.MinCpuPlatform

		r.BootDiskKmsKey = res.BootDiskKmsKey

		r.EphemeralStorageConfig = res.EphemeralStorageConfig

		r.Spot = res.Spot

	}
	return nil
}

// This object is used to assert a desired state where this ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig = &ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig{empty: true}

func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig) Empty() bool {
	return r.empty
}

func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators struct {
	empty            bool    `json:"-"`
	AcceleratorCount *int64  `json:"acceleratorCount"`
	AcceleratorType  *string `json:"acceleratorType"`
	GpuPartitionSize *string `json:"gpuPartitionSize"`
}

type jsonClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators

func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators) UnmarshalJSON(data []byte) error {
	var res jsonClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators
	} else {

		r.AcceleratorCount = res.AcceleratorCount

		r.AcceleratorType = res.AcceleratorType

		r.GpuPartitionSize = res.GpuPartitionSize

	}
	return nil
}

// This object is used to assert a desired state where this ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators = &ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators{empty: true}

func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators) Empty() bool {
	return r.empty
}

func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig struct {
	empty         bool   `json:"-"`
	LocalSsdCount *int64 `json:"localSsdCount"`
}

type jsonClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig

func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig
	} else {

		r.LocalSsdCount = res.LocalSsdCount

	}
	return nil
}

// This object is used to assert a desired state where this ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig = &ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig{empty: true}

func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig) Empty() bool {
	return r.empty
}

func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling struct {
	empty        bool   `json:"-"`
	MinNodeCount *int64 `json:"minNodeCount"`
	MaxNodeCount *int64 `json:"maxNodeCount"`
}

type jsonClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling

func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling) UnmarshalJSON(data []byte) error {
	var res jsonClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling
	} else {

		r.MinNodeCount = res.MinNodeCount

		r.MaxNodeCount = res.MaxNodeCount

	}
	return nil
}

// This object is used to assert a desired state where this ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling = &ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling{empty: true}

func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling) Empty() bool {
	return r.empty
}

func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig struct {
	empty            bool              `json:"-"`
	ComponentVersion map[string]string `json:"componentVersion"`
	Properties       map[string]string `json:"properties"`
}

type jsonClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig

func (r *ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig
	} else {

		r.ComponentVersion = res.ComponentVersion

		r.Properties = res.Properties

	}
	return nil
}

// This object is used to assert a desired state where this ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig *ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig = &ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig{empty: true}

func (r *ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig) Empty() bool {
	return r.empty
}

func (r *ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterVirtualClusterConfigAuxiliaryServicesConfig struct {
	empty                    bool                                                                        `json:"-"`
	MetastoreConfig          *ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig          `json:"metastoreConfig"`
	SparkHistoryServerConfig *ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig `json:"sparkHistoryServerConfig"`
}

type jsonClusterVirtualClusterConfigAuxiliaryServicesConfig ClusterVirtualClusterConfigAuxiliaryServicesConfig

func (r *ClusterVirtualClusterConfigAuxiliaryServicesConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterVirtualClusterConfigAuxiliaryServicesConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterVirtualClusterConfigAuxiliaryServicesConfig
	} else {

		r.MetastoreConfig = res.MetastoreConfig

		r.SparkHistoryServerConfig = res.SparkHistoryServerConfig

	}
	return nil
}

// This object is used to assert a desired state where this ClusterVirtualClusterConfigAuxiliaryServicesConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterVirtualClusterConfigAuxiliaryServicesConfig *ClusterVirtualClusterConfigAuxiliaryServicesConfig = &ClusterVirtualClusterConfigAuxiliaryServicesConfig{empty: true}

func (r *ClusterVirtualClusterConfigAuxiliaryServicesConfig) Empty() bool {
	return r.empty
}

func (r *ClusterVirtualClusterConfigAuxiliaryServicesConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterVirtualClusterConfigAuxiliaryServicesConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig struct {
	empty                    bool    `json:"-"`
	DataprocMetastoreService *string `json:"dataprocMetastoreService"`
}

type jsonClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig

func (r *ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig
	} else {

		r.DataprocMetastoreService = res.DataprocMetastoreService

	}
	return nil
}

// This object is used to assert a desired state where this ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig *ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig = &ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig{empty: true}

func (r *ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig) Empty() bool {
	return r.empty
}

func (r *ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig struct {
	empty           bool    `json:"-"`
	DataprocCluster *string `json:"dataprocCluster"`
}

type jsonClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig

func (r *ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig
	} else {

		r.DataprocCluster = res.DataprocCluster

	}
	return nil
}

// This object is used to assert a desired state where this ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig *ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig = &ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig{empty: true}

func (r *ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig) Empty() bool {
	return r.empty
}

func (r *ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
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
		"project":                dcl.ValueOrEmptyString(nr.Project),
		"name":                   dcl.ValueOrEmptyString(nr.Name),
		"config":                 dcl.ValueOrEmptyString(nr.Config),
		"labels":                 dcl.ValueOrEmptyString(nr.Labels),
		"status":                 dcl.ValueOrEmptyString(nr.Status),
		"status_history":         dcl.ValueOrEmptyString(nr.StatusHistory),
		"cluster_uuid":           dcl.ValueOrEmptyString(nr.ClusterUuid),
		"metrics":                dcl.ValueOrEmptyString(nr.Metrics),
		"location":               dcl.ValueOrEmptyString(nr.Location),
		"virtual_cluster_config": dcl.ValueOrEmptyString(nr.VirtualClusterConfig),
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
	result, err := unmarshalCluster(b, c, r)
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
	rawNew, err := c.GetCluster(ctx, desired)
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

				fullResp, err := unmarshalMapCluster(r, c, rawDesired)
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
