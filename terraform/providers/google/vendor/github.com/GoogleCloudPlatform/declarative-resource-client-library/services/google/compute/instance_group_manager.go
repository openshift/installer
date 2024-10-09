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
package compute

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/api/googleapi"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

type InstanceGroupManager struct {
	Id                  *int64                                    `json:"id"`
	CreationTimestamp   *string                                   `json:"creationTimestamp"`
	Name                *string                                   `json:"name"`
	Description         *string                                   `json:"description"`
	Zone                *string                                   `json:"zone"`
	Region              *string                                   `json:"region"`
	DistributionPolicy  *InstanceGroupManagerDistributionPolicy   `json:"distributionPolicy"`
	InstanceTemplate    *string                                   `json:"instanceTemplate"`
	Versions            []InstanceGroupManagerVersions            `json:"versions"`
	InstanceGroup       *string                                   `json:"instanceGroup"`
	TargetPools         []string                                  `json:"targetPools"`
	BaseInstanceName    *string                                   `json:"baseInstanceName"`
	Fingerprint         *string                                   `json:"fingerprint"`
	CurrentActions      *InstanceGroupManagerCurrentActions       `json:"currentActions"`
	Status              *InstanceGroupManagerStatus               `json:"status"`
	TargetSize          *int64                                    `json:"targetSize"`
	SelfLink            *string                                   `json:"selfLink"`
	AutoHealingPolicies []InstanceGroupManagerAutoHealingPolicies `json:"autoHealingPolicies"`
	UpdatePolicy        *InstanceGroupManagerUpdatePolicy         `json:"updatePolicy"`
	NamedPorts          []InstanceGroupManagerNamedPorts          `json:"namedPorts"`
	StatefulPolicy      *InstanceGroupManagerStatefulPolicy       `json:"statefulPolicy"`
	Project             *string                                   `json:"project"`
	Location            *string                                   `json:"location"`
}

func (r *InstanceGroupManager) String() string {
	return dcl.SprintResource(r)
}

// The enum InstanceGroupManagerDistributionPolicyTargetShapeEnum.
type InstanceGroupManagerDistributionPolicyTargetShapeEnum string

// InstanceGroupManagerDistributionPolicyTargetShapeEnumRef returns a *InstanceGroupManagerDistributionPolicyTargetShapeEnum with the value of string s
// If the empty string is provided, nil is returned.
func InstanceGroupManagerDistributionPolicyTargetShapeEnumRef(s string) *InstanceGroupManagerDistributionPolicyTargetShapeEnum {
	v := InstanceGroupManagerDistributionPolicyTargetShapeEnum(s)
	return &v
}

func (v InstanceGroupManagerDistributionPolicyTargetShapeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"TARGET_SHAPE_UNSPECIFIED", "ANY", "BALANCED", "ANY_SINGLE_ZONE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "InstanceGroupManagerDistributionPolicyTargetShapeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum InstanceGroupManagerUpdatePolicyTypeEnum.
type InstanceGroupManagerUpdatePolicyTypeEnum string

// InstanceGroupManagerUpdatePolicyTypeEnumRef returns a *InstanceGroupManagerUpdatePolicyTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func InstanceGroupManagerUpdatePolicyTypeEnumRef(s string) *InstanceGroupManagerUpdatePolicyTypeEnum {
	v := InstanceGroupManagerUpdatePolicyTypeEnum(s)
	return &v
}

func (v InstanceGroupManagerUpdatePolicyTypeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"OPPORTUNISTIC", "PROACTIVE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "InstanceGroupManagerUpdatePolicyTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum InstanceGroupManagerUpdatePolicyInstanceRedistributionTypeEnum.
type InstanceGroupManagerUpdatePolicyInstanceRedistributionTypeEnum string

// InstanceGroupManagerUpdatePolicyInstanceRedistributionTypeEnumRef returns a *InstanceGroupManagerUpdatePolicyInstanceRedistributionTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func InstanceGroupManagerUpdatePolicyInstanceRedistributionTypeEnumRef(s string) *InstanceGroupManagerUpdatePolicyInstanceRedistributionTypeEnum {
	v := InstanceGroupManagerUpdatePolicyInstanceRedistributionTypeEnum(s)
	return &v
}

func (v InstanceGroupManagerUpdatePolicyInstanceRedistributionTypeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"NONE", "PROACTIVE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "InstanceGroupManagerUpdatePolicyInstanceRedistributionTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum InstanceGroupManagerUpdatePolicyMinimalActionEnum.
type InstanceGroupManagerUpdatePolicyMinimalActionEnum string

// InstanceGroupManagerUpdatePolicyMinimalActionEnumRef returns a *InstanceGroupManagerUpdatePolicyMinimalActionEnum with the value of string s
// If the empty string is provided, nil is returned.
func InstanceGroupManagerUpdatePolicyMinimalActionEnumRef(s string) *InstanceGroupManagerUpdatePolicyMinimalActionEnum {
	v := InstanceGroupManagerUpdatePolicyMinimalActionEnum(s)
	return &v
}

func (v InstanceGroupManagerUpdatePolicyMinimalActionEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REPLACE", "RESTART", "REFRESH", "NONE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "InstanceGroupManagerUpdatePolicyMinimalActionEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum InstanceGroupManagerUpdatePolicyReplacementMethodEnum.
type InstanceGroupManagerUpdatePolicyReplacementMethodEnum string

// InstanceGroupManagerUpdatePolicyReplacementMethodEnumRef returns a *InstanceGroupManagerUpdatePolicyReplacementMethodEnum with the value of string s
// If the empty string is provided, nil is returned.
func InstanceGroupManagerUpdatePolicyReplacementMethodEnumRef(s string) *InstanceGroupManagerUpdatePolicyReplacementMethodEnum {
	v := InstanceGroupManagerUpdatePolicyReplacementMethodEnum(s)
	return &v
}

func (v InstanceGroupManagerUpdatePolicyReplacementMethodEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"SUBSTITUTE", "RECREATE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "InstanceGroupManagerUpdatePolicyReplacementMethodEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum InstanceGroupManagerStatefulPolicyPreservedStateDisksAutoDeleteEnum.
type InstanceGroupManagerStatefulPolicyPreservedStateDisksAutoDeleteEnum string

// InstanceGroupManagerStatefulPolicyPreservedStateDisksAutoDeleteEnumRef returns a *InstanceGroupManagerStatefulPolicyPreservedStateDisksAutoDeleteEnum with the value of string s
// If the empty string is provided, nil is returned.
func InstanceGroupManagerStatefulPolicyPreservedStateDisksAutoDeleteEnumRef(s string) *InstanceGroupManagerStatefulPolicyPreservedStateDisksAutoDeleteEnum {
	v := InstanceGroupManagerStatefulPolicyPreservedStateDisksAutoDeleteEnum(s)
	return &v
}

func (v InstanceGroupManagerStatefulPolicyPreservedStateDisksAutoDeleteEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"NEVER", "ON_PERMANENT_INSTANCE_DELETION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "InstanceGroupManagerStatefulPolicyPreservedStateDisksAutoDeleteEnum",
		Value: string(v),
		Valid: []string{},
	}
}

type InstanceGroupManagerDistributionPolicy struct {
	empty       bool                                                   `json:"-"`
	Zones       []InstanceGroupManagerDistributionPolicyZones          `json:"zones"`
	TargetShape *InstanceGroupManagerDistributionPolicyTargetShapeEnum `json:"targetShape"`
}

type jsonInstanceGroupManagerDistributionPolicy InstanceGroupManagerDistributionPolicy

func (r *InstanceGroupManagerDistributionPolicy) UnmarshalJSON(data []byte) error {
	var res jsonInstanceGroupManagerDistributionPolicy
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInstanceGroupManagerDistributionPolicy
	} else {

		r.Zones = res.Zones

		r.TargetShape = res.TargetShape

	}
	return nil
}

// This object is used to assert a desired state where this InstanceGroupManagerDistributionPolicy is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyInstanceGroupManagerDistributionPolicy *InstanceGroupManagerDistributionPolicy = &InstanceGroupManagerDistributionPolicy{empty: true}

func (r *InstanceGroupManagerDistributionPolicy) Empty() bool {
	return r.empty
}

func (r *InstanceGroupManagerDistributionPolicy) String() string {
	return dcl.SprintResource(r)
}

func (r *InstanceGroupManagerDistributionPolicy) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type InstanceGroupManagerDistributionPolicyZones struct {
	empty bool    `json:"-"`
	Zone  *string `json:"zone"`
}

type jsonInstanceGroupManagerDistributionPolicyZones InstanceGroupManagerDistributionPolicyZones

func (r *InstanceGroupManagerDistributionPolicyZones) UnmarshalJSON(data []byte) error {
	var res jsonInstanceGroupManagerDistributionPolicyZones
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInstanceGroupManagerDistributionPolicyZones
	} else {

		r.Zone = res.Zone

	}
	return nil
}

// This object is used to assert a desired state where this InstanceGroupManagerDistributionPolicyZones is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyInstanceGroupManagerDistributionPolicyZones *InstanceGroupManagerDistributionPolicyZones = &InstanceGroupManagerDistributionPolicyZones{empty: true}

func (r *InstanceGroupManagerDistributionPolicyZones) Empty() bool {
	return r.empty
}

func (r *InstanceGroupManagerDistributionPolicyZones) String() string {
	return dcl.SprintResource(r)
}

func (r *InstanceGroupManagerDistributionPolicyZones) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type InstanceGroupManagerVersions struct {
	empty            bool                                    `json:"-"`
	Name             *string                                 `json:"name"`
	InstanceTemplate *string                                 `json:"instanceTemplate"`
	TargetSize       *InstanceGroupManagerVersionsTargetSize `json:"targetSize"`
}

type jsonInstanceGroupManagerVersions InstanceGroupManagerVersions

func (r *InstanceGroupManagerVersions) UnmarshalJSON(data []byte) error {
	var res jsonInstanceGroupManagerVersions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInstanceGroupManagerVersions
	} else {

		r.Name = res.Name

		r.InstanceTemplate = res.InstanceTemplate

		r.TargetSize = res.TargetSize

	}
	return nil
}

// This object is used to assert a desired state where this InstanceGroupManagerVersions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyInstanceGroupManagerVersions *InstanceGroupManagerVersions = &InstanceGroupManagerVersions{empty: true}

func (r *InstanceGroupManagerVersions) Empty() bool {
	return r.empty
}

func (r *InstanceGroupManagerVersions) String() string {
	return dcl.SprintResource(r)
}

func (r *InstanceGroupManagerVersions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type InstanceGroupManagerVersionsTargetSize struct {
	empty      bool   `json:"-"`
	Fixed      *int64 `json:"fixed"`
	Percent    *int64 `json:"percent"`
	Calculated *int64 `json:"calculated"`
}

type jsonInstanceGroupManagerVersionsTargetSize InstanceGroupManagerVersionsTargetSize

func (r *InstanceGroupManagerVersionsTargetSize) UnmarshalJSON(data []byte) error {
	var res jsonInstanceGroupManagerVersionsTargetSize
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInstanceGroupManagerVersionsTargetSize
	} else {

		r.Fixed = res.Fixed

		r.Percent = res.Percent

		r.Calculated = res.Calculated

	}
	return nil
}

// This object is used to assert a desired state where this InstanceGroupManagerVersionsTargetSize is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyInstanceGroupManagerVersionsTargetSize *InstanceGroupManagerVersionsTargetSize = &InstanceGroupManagerVersionsTargetSize{empty: true}

func (r *InstanceGroupManagerVersionsTargetSize) Empty() bool {
	return r.empty
}

func (r *InstanceGroupManagerVersionsTargetSize) String() string {
	return dcl.SprintResource(r)
}

func (r *InstanceGroupManagerVersionsTargetSize) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type InstanceGroupManagerCurrentActions struct {
	empty                  bool   `json:"-"`
	None                   *int64 `json:"none"`
	Creating               *int64 `json:"creating"`
	CreatingWithoutRetries *int64 `json:"creatingWithoutRetries"`
	Verifying              *int64 `json:"verifying"`
	Recreating             *int64 `json:"recreating"`
	Deleting               *int64 `json:"deleting"`
	Abandoning             *int64 `json:"abandoning"`
	Restarting             *int64 `json:"restarting"`
	Refreshing             *int64 `json:"refreshing"`
}

type jsonInstanceGroupManagerCurrentActions InstanceGroupManagerCurrentActions

func (r *InstanceGroupManagerCurrentActions) UnmarshalJSON(data []byte) error {
	var res jsonInstanceGroupManagerCurrentActions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInstanceGroupManagerCurrentActions
	} else {

		r.None = res.None

		r.Creating = res.Creating

		r.CreatingWithoutRetries = res.CreatingWithoutRetries

		r.Verifying = res.Verifying

		r.Recreating = res.Recreating

		r.Deleting = res.Deleting

		r.Abandoning = res.Abandoning

		r.Restarting = res.Restarting

		r.Refreshing = res.Refreshing

	}
	return nil
}

// This object is used to assert a desired state where this InstanceGroupManagerCurrentActions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyInstanceGroupManagerCurrentActions *InstanceGroupManagerCurrentActions = &InstanceGroupManagerCurrentActions{empty: true}

func (r *InstanceGroupManagerCurrentActions) Empty() bool {
	return r.empty
}

func (r *InstanceGroupManagerCurrentActions) String() string {
	return dcl.SprintResource(r)
}

func (r *InstanceGroupManagerCurrentActions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type InstanceGroupManagerStatus struct {
	empty         bool                                     `json:"-"`
	IsStable      *bool                                    `json:"isStable"`
	VersionTarget *InstanceGroupManagerStatusVersionTarget `json:"versionTarget"`
	Stateful      *InstanceGroupManagerStatusStateful      `json:"stateful"`
	Autoscaler    *string                                  `json:"autoscaler"`
}

type jsonInstanceGroupManagerStatus InstanceGroupManagerStatus

func (r *InstanceGroupManagerStatus) UnmarshalJSON(data []byte) error {
	var res jsonInstanceGroupManagerStatus
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInstanceGroupManagerStatus
	} else {

		r.IsStable = res.IsStable

		r.VersionTarget = res.VersionTarget

		r.Stateful = res.Stateful

		r.Autoscaler = res.Autoscaler

	}
	return nil
}

// This object is used to assert a desired state where this InstanceGroupManagerStatus is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyInstanceGroupManagerStatus *InstanceGroupManagerStatus = &InstanceGroupManagerStatus{empty: true}

func (r *InstanceGroupManagerStatus) Empty() bool {
	return r.empty
}

func (r *InstanceGroupManagerStatus) String() string {
	return dcl.SprintResource(r)
}

func (r *InstanceGroupManagerStatus) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type InstanceGroupManagerStatusVersionTarget struct {
	empty     bool  `json:"-"`
	IsReached *bool `json:"isReached"`
}

type jsonInstanceGroupManagerStatusVersionTarget InstanceGroupManagerStatusVersionTarget

func (r *InstanceGroupManagerStatusVersionTarget) UnmarshalJSON(data []byte) error {
	var res jsonInstanceGroupManagerStatusVersionTarget
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInstanceGroupManagerStatusVersionTarget
	} else {

		r.IsReached = res.IsReached

	}
	return nil
}

// This object is used to assert a desired state where this InstanceGroupManagerStatusVersionTarget is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyInstanceGroupManagerStatusVersionTarget *InstanceGroupManagerStatusVersionTarget = &InstanceGroupManagerStatusVersionTarget{empty: true}

func (r *InstanceGroupManagerStatusVersionTarget) Empty() bool {
	return r.empty
}

func (r *InstanceGroupManagerStatusVersionTarget) String() string {
	return dcl.SprintResource(r)
}

func (r *InstanceGroupManagerStatusVersionTarget) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type InstanceGroupManagerStatusStateful struct {
	empty              bool                                                  `json:"-"`
	HasStatefulConfig  *bool                                                 `json:"hasStatefulConfig"`
	PerInstanceConfigs *InstanceGroupManagerStatusStatefulPerInstanceConfigs `json:"perInstanceConfigs"`
}

type jsonInstanceGroupManagerStatusStateful InstanceGroupManagerStatusStateful

func (r *InstanceGroupManagerStatusStateful) UnmarshalJSON(data []byte) error {
	var res jsonInstanceGroupManagerStatusStateful
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInstanceGroupManagerStatusStateful
	} else {

		r.HasStatefulConfig = res.HasStatefulConfig

		r.PerInstanceConfigs = res.PerInstanceConfigs

	}
	return nil
}

// This object is used to assert a desired state where this InstanceGroupManagerStatusStateful is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyInstanceGroupManagerStatusStateful *InstanceGroupManagerStatusStateful = &InstanceGroupManagerStatusStateful{empty: true}

func (r *InstanceGroupManagerStatusStateful) Empty() bool {
	return r.empty
}

func (r *InstanceGroupManagerStatusStateful) String() string {
	return dcl.SprintResource(r)
}

func (r *InstanceGroupManagerStatusStateful) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type InstanceGroupManagerStatusStatefulPerInstanceConfigs struct {
	empty        bool  `json:"-"`
	AllEffective *bool `json:"allEffective"`
}

type jsonInstanceGroupManagerStatusStatefulPerInstanceConfigs InstanceGroupManagerStatusStatefulPerInstanceConfigs

func (r *InstanceGroupManagerStatusStatefulPerInstanceConfigs) UnmarshalJSON(data []byte) error {
	var res jsonInstanceGroupManagerStatusStatefulPerInstanceConfigs
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInstanceGroupManagerStatusStatefulPerInstanceConfigs
	} else {

		r.AllEffective = res.AllEffective

	}
	return nil
}

// This object is used to assert a desired state where this InstanceGroupManagerStatusStatefulPerInstanceConfigs is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyInstanceGroupManagerStatusStatefulPerInstanceConfigs *InstanceGroupManagerStatusStatefulPerInstanceConfigs = &InstanceGroupManagerStatusStatefulPerInstanceConfigs{empty: true}

func (r *InstanceGroupManagerStatusStatefulPerInstanceConfigs) Empty() bool {
	return r.empty
}

func (r *InstanceGroupManagerStatusStatefulPerInstanceConfigs) String() string {
	return dcl.SprintResource(r)
}

func (r *InstanceGroupManagerStatusStatefulPerInstanceConfigs) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type InstanceGroupManagerAutoHealingPolicies struct {
	empty           bool    `json:"-"`
	HealthCheck     *string `json:"healthCheck"`
	InitialDelaySec *int64  `json:"initialDelaySec"`
}

type jsonInstanceGroupManagerAutoHealingPolicies InstanceGroupManagerAutoHealingPolicies

func (r *InstanceGroupManagerAutoHealingPolicies) UnmarshalJSON(data []byte) error {
	var res jsonInstanceGroupManagerAutoHealingPolicies
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInstanceGroupManagerAutoHealingPolicies
	} else {

		r.HealthCheck = res.HealthCheck

		r.InitialDelaySec = res.InitialDelaySec

	}
	return nil
}

// This object is used to assert a desired state where this InstanceGroupManagerAutoHealingPolicies is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyInstanceGroupManagerAutoHealingPolicies *InstanceGroupManagerAutoHealingPolicies = &InstanceGroupManagerAutoHealingPolicies{empty: true}

func (r *InstanceGroupManagerAutoHealingPolicies) Empty() bool {
	return r.empty
}

func (r *InstanceGroupManagerAutoHealingPolicies) String() string {
	return dcl.SprintResource(r)
}

func (r *InstanceGroupManagerAutoHealingPolicies) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type InstanceGroupManagerUpdatePolicy struct {
	empty                      bool                                                            `json:"-"`
	Type                       *InstanceGroupManagerUpdatePolicyTypeEnum                       `json:"type"`
	InstanceRedistributionType *InstanceGroupManagerUpdatePolicyInstanceRedistributionTypeEnum `json:"instanceRedistributionType"`
	MinimalAction              *InstanceGroupManagerUpdatePolicyMinimalActionEnum              `json:"minimalAction"`
	MaxSurge                   *InstanceGroupManagerUpdatePolicyMaxSurge                       `json:"maxSurge"`
	MaxUnavailable             *InstanceGroupManagerUpdatePolicyMaxUnavailable                 `json:"maxUnavailable"`
	ReplacementMethod          *InstanceGroupManagerUpdatePolicyReplacementMethodEnum          `json:"replacementMethod"`
}

type jsonInstanceGroupManagerUpdatePolicy InstanceGroupManagerUpdatePolicy

func (r *InstanceGroupManagerUpdatePolicy) UnmarshalJSON(data []byte) error {
	var res jsonInstanceGroupManagerUpdatePolicy
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInstanceGroupManagerUpdatePolicy
	} else {

		r.Type = res.Type

		r.InstanceRedistributionType = res.InstanceRedistributionType

		r.MinimalAction = res.MinimalAction

		r.MaxSurge = res.MaxSurge

		r.MaxUnavailable = res.MaxUnavailable

		r.ReplacementMethod = res.ReplacementMethod

	}
	return nil
}

// This object is used to assert a desired state where this InstanceGroupManagerUpdatePolicy is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyInstanceGroupManagerUpdatePolicy *InstanceGroupManagerUpdatePolicy = &InstanceGroupManagerUpdatePolicy{empty: true}

func (r *InstanceGroupManagerUpdatePolicy) Empty() bool {
	return r.empty
}

func (r *InstanceGroupManagerUpdatePolicy) String() string {
	return dcl.SprintResource(r)
}

func (r *InstanceGroupManagerUpdatePolicy) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type InstanceGroupManagerUpdatePolicyMaxSurge struct {
	empty      bool   `json:"-"`
	Fixed      *int64 `json:"fixed"`
	Percent    *int64 `json:"percent"`
	Calculated *int64 `json:"calculated"`
}

type jsonInstanceGroupManagerUpdatePolicyMaxSurge InstanceGroupManagerUpdatePolicyMaxSurge

func (r *InstanceGroupManagerUpdatePolicyMaxSurge) UnmarshalJSON(data []byte) error {
	var res jsonInstanceGroupManagerUpdatePolicyMaxSurge
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInstanceGroupManagerUpdatePolicyMaxSurge
	} else {

		r.Fixed = res.Fixed

		r.Percent = res.Percent

		r.Calculated = res.Calculated

	}
	return nil
}

// This object is used to assert a desired state where this InstanceGroupManagerUpdatePolicyMaxSurge is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyInstanceGroupManagerUpdatePolicyMaxSurge *InstanceGroupManagerUpdatePolicyMaxSurge = &InstanceGroupManagerUpdatePolicyMaxSurge{empty: true}

func (r *InstanceGroupManagerUpdatePolicyMaxSurge) Empty() bool {
	return r.empty
}

func (r *InstanceGroupManagerUpdatePolicyMaxSurge) String() string {
	return dcl.SprintResource(r)
}

func (r *InstanceGroupManagerUpdatePolicyMaxSurge) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type InstanceGroupManagerUpdatePolicyMaxUnavailable struct {
	empty      bool   `json:"-"`
	Fixed      *int64 `json:"fixed"`
	Percent    *int64 `json:"percent"`
	Calculated *int64 `json:"calculated"`
}

type jsonInstanceGroupManagerUpdatePolicyMaxUnavailable InstanceGroupManagerUpdatePolicyMaxUnavailable

func (r *InstanceGroupManagerUpdatePolicyMaxUnavailable) UnmarshalJSON(data []byte) error {
	var res jsonInstanceGroupManagerUpdatePolicyMaxUnavailable
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInstanceGroupManagerUpdatePolicyMaxUnavailable
	} else {

		r.Fixed = res.Fixed

		r.Percent = res.Percent

		r.Calculated = res.Calculated

	}
	return nil
}

// This object is used to assert a desired state where this InstanceGroupManagerUpdatePolicyMaxUnavailable is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyInstanceGroupManagerUpdatePolicyMaxUnavailable *InstanceGroupManagerUpdatePolicyMaxUnavailable = &InstanceGroupManagerUpdatePolicyMaxUnavailable{empty: true}

func (r *InstanceGroupManagerUpdatePolicyMaxUnavailable) Empty() bool {
	return r.empty
}

func (r *InstanceGroupManagerUpdatePolicyMaxUnavailable) String() string {
	return dcl.SprintResource(r)
}

func (r *InstanceGroupManagerUpdatePolicyMaxUnavailable) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type InstanceGroupManagerNamedPorts struct {
	empty bool    `json:"-"`
	Name  *string `json:"name"`
	Port  *int64  `json:"port"`
}

type jsonInstanceGroupManagerNamedPorts InstanceGroupManagerNamedPorts

func (r *InstanceGroupManagerNamedPorts) UnmarshalJSON(data []byte) error {
	var res jsonInstanceGroupManagerNamedPorts
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInstanceGroupManagerNamedPorts
	} else {

		r.Name = res.Name

		r.Port = res.Port

	}
	return nil
}

// This object is used to assert a desired state where this InstanceGroupManagerNamedPorts is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyInstanceGroupManagerNamedPorts *InstanceGroupManagerNamedPorts = &InstanceGroupManagerNamedPorts{empty: true}

func (r *InstanceGroupManagerNamedPorts) Empty() bool {
	return r.empty
}

func (r *InstanceGroupManagerNamedPorts) String() string {
	return dcl.SprintResource(r)
}

func (r *InstanceGroupManagerNamedPorts) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type InstanceGroupManagerStatefulPolicy struct {
	empty          bool                                              `json:"-"`
	PreservedState *InstanceGroupManagerStatefulPolicyPreservedState `json:"preservedState"`
}

type jsonInstanceGroupManagerStatefulPolicy InstanceGroupManagerStatefulPolicy

func (r *InstanceGroupManagerStatefulPolicy) UnmarshalJSON(data []byte) error {
	var res jsonInstanceGroupManagerStatefulPolicy
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInstanceGroupManagerStatefulPolicy
	} else {

		r.PreservedState = res.PreservedState

	}
	return nil
}

// This object is used to assert a desired state where this InstanceGroupManagerStatefulPolicy is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyInstanceGroupManagerStatefulPolicy *InstanceGroupManagerStatefulPolicy = &InstanceGroupManagerStatefulPolicy{empty: true}

func (r *InstanceGroupManagerStatefulPolicy) Empty() bool {
	return r.empty
}

func (r *InstanceGroupManagerStatefulPolicy) String() string {
	return dcl.SprintResource(r)
}

func (r *InstanceGroupManagerStatefulPolicy) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type InstanceGroupManagerStatefulPolicyPreservedState struct {
	empty bool                                                             `json:"-"`
	Disks map[string]InstanceGroupManagerStatefulPolicyPreservedStateDisks `json:"disks"`
}

type jsonInstanceGroupManagerStatefulPolicyPreservedState InstanceGroupManagerStatefulPolicyPreservedState

func (r *InstanceGroupManagerStatefulPolicyPreservedState) UnmarshalJSON(data []byte) error {
	var res jsonInstanceGroupManagerStatefulPolicyPreservedState
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInstanceGroupManagerStatefulPolicyPreservedState
	} else {

		r.Disks = res.Disks

	}
	return nil
}

// This object is used to assert a desired state where this InstanceGroupManagerStatefulPolicyPreservedState is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyInstanceGroupManagerStatefulPolicyPreservedState *InstanceGroupManagerStatefulPolicyPreservedState = &InstanceGroupManagerStatefulPolicyPreservedState{empty: true}

func (r *InstanceGroupManagerStatefulPolicyPreservedState) Empty() bool {
	return r.empty
}

func (r *InstanceGroupManagerStatefulPolicyPreservedState) String() string {
	return dcl.SprintResource(r)
}

func (r *InstanceGroupManagerStatefulPolicyPreservedState) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type InstanceGroupManagerStatefulPolicyPreservedStateDisks struct {
	empty      bool                                                                 `json:"-"`
	AutoDelete *InstanceGroupManagerStatefulPolicyPreservedStateDisksAutoDeleteEnum `json:"autoDelete"`
}

type jsonInstanceGroupManagerStatefulPolicyPreservedStateDisks InstanceGroupManagerStatefulPolicyPreservedStateDisks

func (r *InstanceGroupManagerStatefulPolicyPreservedStateDisks) UnmarshalJSON(data []byte) error {
	var res jsonInstanceGroupManagerStatefulPolicyPreservedStateDisks
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInstanceGroupManagerStatefulPolicyPreservedStateDisks
	} else {

		r.AutoDelete = res.AutoDelete

	}
	return nil
}

// This object is used to assert a desired state where this InstanceGroupManagerStatefulPolicyPreservedStateDisks is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyInstanceGroupManagerStatefulPolicyPreservedStateDisks *InstanceGroupManagerStatefulPolicyPreservedStateDisks = &InstanceGroupManagerStatefulPolicyPreservedStateDisks{empty: true}

func (r *InstanceGroupManagerStatefulPolicyPreservedStateDisks) Empty() bool {
	return r.empty
}

func (r *InstanceGroupManagerStatefulPolicyPreservedStateDisks) String() string {
	return dcl.SprintResource(r)
}

func (r *InstanceGroupManagerStatefulPolicyPreservedStateDisks) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *InstanceGroupManager) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "compute",
		Type:    "InstanceGroupManager",
		Version: "compute",
	}
}

func (r *InstanceGroupManager) ID() (string, error) {
	if err := extractInstanceGroupManagerFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"id":                    dcl.ValueOrEmptyString(nr.Id),
		"creation_timestamp":    dcl.ValueOrEmptyString(nr.CreationTimestamp),
		"name":                  dcl.ValueOrEmptyString(nr.Name),
		"description":           dcl.ValueOrEmptyString(nr.Description),
		"zone":                  dcl.ValueOrEmptyString(nr.Zone),
		"region":                dcl.ValueOrEmptyString(nr.Region),
		"distribution_policy":   dcl.ValueOrEmptyString(nr.DistributionPolicy),
		"instance_template":     dcl.ValueOrEmptyString(nr.InstanceTemplate),
		"versions":              dcl.ValueOrEmptyString(nr.Versions),
		"instance_group":        dcl.ValueOrEmptyString(nr.InstanceGroup),
		"target_pools":          dcl.ValueOrEmptyString(nr.TargetPools),
		"base_instance_name":    dcl.ValueOrEmptyString(nr.BaseInstanceName),
		"fingerprint":           dcl.ValueOrEmptyString(nr.Fingerprint),
		"current_actions":       dcl.ValueOrEmptyString(nr.CurrentActions),
		"status":                dcl.ValueOrEmptyString(nr.Status),
		"target_size":           dcl.ValueOrEmptyString(nr.TargetSize),
		"self_link":             dcl.ValueOrEmptyString(nr.SelfLink),
		"auto_healing_policies": dcl.ValueOrEmptyString(nr.AutoHealingPolicies),
		"update_policy":         dcl.ValueOrEmptyString(nr.UpdatePolicy),
		"named_ports":           dcl.ValueOrEmptyString(nr.NamedPorts),
		"stateful_policy":       dcl.ValueOrEmptyString(nr.StatefulPolicy),
		"project":               dcl.ValueOrEmptyString(nr.Project),
		"location":              dcl.ValueOrEmptyString(nr.Location),
	}
	if dcl.IsRegion(nr.Location) {
		return dcl.Nprintf("projects/{{project}}/regions/{{location}}/instanceGroupManagers/{{name}}", params), nil
	}

	if dcl.IsZone(nr.Location) {
		return dcl.Nprintf("projects/{{project}}/zones/{{location}}/instanceGroupManagers/{{name}}", params), nil
	}

	return dcl.Nprintf("", params), nil
}

const InstanceGroupManagerMaxPage = -1

type InstanceGroupManagerList struct {
	Items []*InstanceGroupManager

	nextToken string

	pageSize int32

	resource *InstanceGroupManager
}

func (l *InstanceGroupManagerList) HasNext() bool {
	return l.nextToken != ""
}

func (l *InstanceGroupManagerList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listInstanceGroupManager(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListInstanceGroupManager(ctx context.Context, project, location string) (*InstanceGroupManagerList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	c = NewClient(c.Config.Clone(dcl.WithCodeRetryability(map[int]dcl.Retryability{
		412: dcl.Retryability{
			Retryable: false,
			Pattern:   "",
			Timeout:   0,
		},
	})))
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListInstanceGroupManagerWithMaxResults(ctx, project, location, InstanceGroupManagerMaxPage)

}

func (c *Client) ListInstanceGroupManagerWithMaxResults(ctx context.Context, project, location string, pageSize int32) (*InstanceGroupManagerList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &InstanceGroupManager{
		Project:  &project,
		Location: &location,
	}
	items, token, err := c.listInstanceGroupManager(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &InstanceGroupManagerList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetInstanceGroupManager(ctx context.Context, r *InstanceGroupManager) (*InstanceGroupManager, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	c = NewClient(c.Config.Clone(dcl.WithCodeRetryability(map[int]dcl.Retryability{
		412: dcl.Retryability{
			Retryable: false,
			Pattern:   "",
			Timeout:   0,
		},
	})))
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractInstanceGroupManagerFields(r)

	b, err := c.getInstanceGroupManagerRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalInstanceGroupManager(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Location = r.Location
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeInstanceGroupManagerNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractInstanceGroupManagerFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteInstanceGroupManager(ctx context.Context, r *InstanceGroupManager) error {
	ctx = dcl.ContextWithRequestID(ctx)
	c = NewClient(c.Config.Clone(dcl.WithCodeRetryability(map[int]dcl.Retryability{
		412: dcl.Retryability{
			Retryable: false,
			Pattern:   "",
			Timeout:   0,
		},
	})))
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("InstanceGroupManager resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting InstanceGroupManager...")
	deleteOp := deleteInstanceGroupManagerOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllInstanceGroupManager deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllInstanceGroupManager(ctx context.Context, project, location string, filter func(*InstanceGroupManager) bool) error {
	listObj, err := c.ListInstanceGroupManager(ctx, project, location)
	if err != nil {
		return err
	}

	err = c.deleteAllInstanceGroupManager(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllInstanceGroupManager(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyInstanceGroupManager(ctx context.Context, rawDesired *InstanceGroupManager, opts ...dcl.ApplyOption) (*InstanceGroupManager, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	c = NewClient(c.Config.Clone(dcl.WithCodeRetryability(map[int]dcl.Retryability{
		412: dcl.Retryability{
			Retryable: false,
			Pattern:   "",
			Timeout:   0,
		},
	})))
	var resultNewState *InstanceGroupManager
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyInstanceGroupManagerHelper(c, ctx, rawDesired, opts...)
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

func applyInstanceGroupManagerHelper(c *Client, ctx context.Context, rawDesired *InstanceGroupManager, opts ...dcl.ApplyOption) (*InstanceGroupManager, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyInstanceGroupManager...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractInstanceGroupManagerFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.instanceGroupManagerDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToInstanceGroupManagerDiffs(c.Config, fieldDiffs, opts)
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
	var ops []instanceGroupManagerApiOperation
	if create {
		ops = append(ops, &createInstanceGroupManagerOperation{})
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
	return applyInstanceGroupManagerDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyInstanceGroupManagerDiff(c *Client, ctx context.Context, desired *InstanceGroupManager, rawDesired *InstanceGroupManager, ops []instanceGroupManagerApiOperation, opts ...dcl.ApplyOption) (*InstanceGroupManager, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetInstanceGroupManager(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createInstanceGroupManagerOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapInstanceGroupManager(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeInstanceGroupManagerNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeInstanceGroupManagerNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeInstanceGroupManagerDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractInstanceGroupManagerFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractInstanceGroupManagerFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffInstanceGroupManager(c, newDesired, newState)
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
