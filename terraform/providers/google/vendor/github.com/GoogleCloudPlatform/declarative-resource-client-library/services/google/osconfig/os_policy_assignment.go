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
package osconfig

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/api/googleapi"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

type OSPolicyAssignment struct {
	Name               *string                             `json:"name"`
	Description        *string                             `json:"description"`
	OSPolicies         []OSPolicyAssignmentOSPolicies      `json:"osPolicies"`
	InstanceFilter     *OSPolicyAssignmentInstanceFilter   `json:"instanceFilter"`
	Rollout            *OSPolicyAssignmentRollout          `json:"rollout"`
	RevisionId         *string                             `json:"revisionId"`
	RevisionCreateTime *string                             `json:"revisionCreateTime"`
	Etag               *string                             `json:"etag"`
	RolloutState       *OSPolicyAssignmentRolloutStateEnum `json:"rolloutState"`
	Baseline           *bool                               `json:"baseline"`
	Deleted            *bool                               `json:"deleted"`
	Reconciling        *bool                               `json:"reconciling"`
	Uid                *string                             `json:"uid"`
	Project            *string                             `json:"project"`
	Location           *string                             `json:"location"`
	SkipAwaitRollout   *bool                               `json:"skipAwaitRollout"`
}

func (r *OSPolicyAssignment) String() string {
	return dcl.SprintResource(r)
}

// The enum OSPolicyAssignmentOSPoliciesModeEnum.
type OSPolicyAssignmentOSPoliciesModeEnum string

// OSPolicyAssignmentOSPoliciesModeEnumRef returns a *OSPolicyAssignmentOSPoliciesModeEnum with the value of string s
// If the empty string is provided, nil is returned.
func OSPolicyAssignmentOSPoliciesModeEnumRef(s string) *OSPolicyAssignmentOSPoliciesModeEnum {
	v := OSPolicyAssignmentOSPoliciesModeEnum(s)
	return &v
}

func (v OSPolicyAssignmentOSPoliciesModeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"MODE_UNSPECIFIED", "VALIDATION", "ENFORCEMENT"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "OSPolicyAssignmentOSPoliciesModeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnum.
type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnum string

// OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnumRef returns a *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnum with the value of string s
// If the empty string is provided, nil is returned.
func OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnumRef(s string) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnum {
	v := OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnum(s)
	return &v
}

func (v OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"DESIRED_STATE_UNSPECIFIED", "INSTALLED", "REMOVED"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnum.
type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnum string

// OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnumRef returns a *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnumRef(s string) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnum {
	v := OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnum(s)
	return &v
}

func (v OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ARCHIVE_TYPE_UNSPECIFIED", "DEB", "DEB_SRC"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnum.
type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnum string

// OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnumRef returns a *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnum with the value of string s
// If the empty string is provided, nil is returned.
func OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnumRef(s string) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnum {
	v := OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnum(s)
	return &v
}

func (v OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"INTERPRETER_UNSPECIFIED", "NONE", "SHELL", "POWERSHELL"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnum.
type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnum string

// OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnumRef returns a *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnum with the value of string s
// If the empty string is provided, nil is returned.
func OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnumRef(s string) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnum {
	v := OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnum(s)
	return &v
}

func (v OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"INTERPRETER_UNSPECIFIED", "NONE", "SHELL", "POWERSHELL"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnum.
type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnum string

// OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnumRef returns a *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnum with the value of string s
// If the empty string is provided, nil is returned.
func OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnumRef(s string) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnum {
	v := OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnum(s)
	return &v
}

func (v OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"OS_POLICY_COMPLIANCE_STATE_UNSPECIFIED", "COMPLIANT", "NON_COMPLIANT", "UNKNOWN", "NO_OS_POLICIES_APPLICABLE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum OSPolicyAssignmentRolloutStateEnum.
type OSPolicyAssignmentRolloutStateEnum string

// OSPolicyAssignmentRolloutStateEnumRef returns a *OSPolicyAssignmentRolloutStateEnum with the value of string s
// If the empty string is provided, nil is returned.
func OSPolicyAssignmentRolloutStateEnumRef(s string) *OSPolicyAssignmentRolloutStateEnum {
	v := OSPolicyAssignmentRolloutStateEnum(s)
	return &v
}

func (v OSPolicyAssignmentRolloutStateEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ROLLOUT_STATE_UNSPECIFIED", "IN_PROGRESS", "CANCELLING", "CANCELLED", "SUCCEEDED"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "OSPolicyAssignmentRolloutStateEnum",
		Value: string(v),
		Valid: []string{},
	}
}

type OSPolicyAssignmentOSPolicies struct {
	empty                     bool                                         `json:"-"`
	Id                        *string                                      `json:"id"`
	Description               *string                                      `json:"description"`
	Mode                      *OSPolicyAssignmentOSPoliciesModeEnum        `json:"mode"`
	ResourceGroups            []OSPolicyAssignmentOSPoliciesResourceGroups `json:"resourceGroups"`
	AllowNoResourceGroupMatch *bool                                        `json:"allowNoResourceGroupMatch"`
}

type jsonOSPolicyAssignmentOSPolicies OSPolicyAssignmentOSPolicies

func (r *OSPolicyAssignmentOSPolicies) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPolicies
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPolicies
	} else {

		r.Id = res.Id

		r.Description = res.Description

		r.Mode = res.Mode

		r.ResourceGroups = res.ResourceGroups

		r.AllowNoResourceGroupMatch = res.AllowNoResourceGroupMatch

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPolicies is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPolicies *OSPolicyAssignmentOSPolicies = &OSPolicyAssignmentOSPolicies{empty: true}

func (r *OSPolicyAssignmentOSPolicies) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPolicies) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPolicies) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroups struct {
	empty            bool                                                         `json:"-"`
	InventoryFilters []OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters `json:"inventoryFilters"`
	Resources        []OSPolicyAssignmentOSPoliciesResourceGroupsResources        `json:"resources"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroups OSPolicyAssignmentOSPoliciesResourceGroups

func (r *OSPolicyAssignmentOSPoliciesResourceGroups) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroups
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroups
	} else {

		r.InventoryFilters = res.InventoryFilters

		r.Resources = res.Resources

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroups is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroups *OSPolicyAssignmentOSPoliciesResourceGroups = &OSPolicyAssignmentOSPoliciesResourceGroups{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroups) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroups) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroups) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters struct {
	empty       bool    `json:"-"`
	OSShortName *string `json:"osShortName"`
	OSVersion   *string `json:"osVersion"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters
	} else {

		r.OSShortName = res.OSShortName

		r.OSVersion = res.OSVersion

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters *OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters = &OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResources struct {
	empty      bool                                                           `json:"-"`
	Id         *string                                                        `json:"id"`
	Pkg        *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg        `json:"pkg"`
	Repository *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository `json:"repository"`
	Exec       *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec       `json:"exec"`
	File       *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile       `json:"file"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResources OSPolicyAssignmentOSPoliciesResourceGroupsResources

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResources) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResources
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResources
	} else {

		r.Id = res.Id

		r.Pkg = res.Pkg

		r.Repository = res.Repository

		r.Exec = res.Exec

		r.File = res.File

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResources is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResources *OSPolicyAssignmentOSPoliciesResourceGroupsResources = &OSPolicyAssignmentOSPoliciesResourceGroupsResources{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResources) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResources) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResources) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg struct {
	empty        bool                                                                    `json:"-"`
	DesiredState *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnum `json:"desiredState"`
	Apt          *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt              `json:"apt"`
	Deb          *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb              `json:"deb"`
	Yum          *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum              `json:"yum"`
	Zypper       *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper           `json:"zypper"`
	Rpm          *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm              `json:"rpm"`
	Googet       *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget           `json:"googet"`
	Msi          *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi              `json:"msi"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg
	} else {

		r.DesiredState = res.DesiredState

		r.Apt = res.Apt

		r.Deb = res.Deb

		r.Yum = res.Yum

		r.Zypper = res.Zypper

		r.Rpm = res.Rpm

		r.Googet = res.Googet

		r.Msi = res.Msi

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt struct {
	empty bool    `json:"-"`
	Name  *string `json:"name"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt
	} else {

		r.Name = res.Name

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb struct {
	empty    bool                                                             `json:"-"`
	Source   *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource `json:"source"`
	PullDeps *bool                                                            `json:"pullDeps"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb
	} else {

		r.Source = res.Source

		r.PullDeps = res.PullDeps

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource struct {
	empty         bool                                                                   `json:"-"`
	Remote        *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote `json:"remote"`
	Gcs           *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs    `json:"gcs"`
	LocalPath     *string                                                                `json:"localPath"`
	AllowInsecure *bool                                                                  `json:"allowInsecure"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource
	} else {

		r.Remote = res.Remote

		r.Gcs = res.Gcs

		r.LocalPath = res.LocalPath

		r.AllowInsecure = res.AllowInsecure

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote struct {
	empty          bool    `json:"-"`
	Uri            *string `json:"uri"`
	Sha256Checksum *string `json:"sha256Checksum"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote
	} else {

		r.Uri = res.Uri

		r.Sha256Checksum = res.Sha256Checksum

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs struct {
	empty      bool    `json:"-"`
	Bucket     *string `json:"bucket"`
	Object     *string `json:"object"`
	Generation *int64  `json:"generation"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs
	} else {

		r.Bucket = res.Bucket

		r.Object = res.Object

		r.Generation = res.Generation

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum struct {
	empty bool    `json:"-"`
	Name  *string `json:"name"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum
	} else {

		r.Name = res.Name

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper struct {
	empty bool    `json:"-"`
	Name  *string `json:"name"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper
	} else {

		r.Name = res.Name

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm struct {
	empty    bool                                                             `json:"-"`
	Source   *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource `json:"source"`
	PullDeps *bool                                                            `json:"pullDeps"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm
	} else {

		r.Source = res.Source

		r.PullDeps = res.PullDeps

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource struct {
	empty         bool                                                                   `json:"-"`
	Remote        *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote `json:"remote"`
	Gcs           *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs    `json:"gcs"`
	LocalPath     *string                                                                `json:"localPath"`
	AllowInsecure *bool                                                                  `json:"allowInsecure"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource
	} else {

		r.Remote = res.Remote

		r.Gcs = res.Gcs

		r.LocalPath = res.LocalPath

		r.AllowInsecure = res.AllowInsecure

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote struct {
	empty          bool    `json:"-"`
	Uri            *string `json:"uri"`
	Sha256Checksum *string `json:"sha256Checksum"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote
	} else {

		r.Uri = res.Uri

		r.Sha256Checksum = res.Sha256Checksum

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs struct {
	empty      bool    `json:"-"`
	Bucket     *string `json:"bucket"`
	Object     *string `json:"object"`
	Generation *int64  `json:"generation"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs
	} else {

		r.Bucket = res.Bucket

		r.Object = res.Object

		r.Generation = res.Generation

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget struct {
	empty bool    `json:"-"`
	Name  *string `json:"name"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget
	} else {

		r.Name = res.Name

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi struct {
	empty      bool                                                             `json:"-"`
	Source     *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource `json:"source"`
	Properties []string                                                         `json:"properties"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi
	} else {

		r.Source = res.Source

		r.Properties = res.Properties

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource struct {
	empty         bool                                                                   `json:"-"`
	Remote        *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote `json:"remote"`
	Gcs           *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs    `json:"gcs"`
	LocalPath     *string                                                                `json:"localPath"`
	AllowInsecure *bool                                                                  `json:"allowInsecure"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource
	} else {

		r.Remote = res.Remote

		r.Gcs = res.Gcs

		r.LocalPath = res.LocalPath

		r.AllowInsecure = res.AllowInsecure

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote struct {
	empty          bool    `json:"-"`
	Uri            *string `json:"uri"`
	Sha256Checksum *string `json:"sha256Checksum"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote
	} else {

		r.Uri = res.Uri

		r.Sha256Checksum = res.Sha256Checksum

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs struct {
	empty      bool    `json:"-"`
	Bucket     *string `json:"bucket"`
	Object     *string `json:"object"`
	Generation *int64  `json:"generation"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs
	} else {

		r.Bucket = res.Bucket

		r.Object = res.Object

		r.Generation = res.Generation

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository struct {
	empty  bool                                                                 `json:"-"`
	Apt    *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt    `json:"apt"`
	Yum    *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum    `json:"yum"`
	Zypper *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper `json:"zypper"`
	Goo    *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo    `json:"goo"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository
	} else {

		r.Apt = res.Apt

		r.Yum = res.Yum

		r.Zypper = res.Zypper

		r.Goo = res.Goo

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt struct {
	empty        bool                                                                             `json:"-"`
	ArchiveType  *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnum `json:"archiveType"`
	Uri          *string                                                                          `json:"uri"`
	Distribution *string                                                                          `json:"distribution"`
	Components   []string                                                                         `json:"components"`
	GpgKey       *string                                                                          `json:"gpgKey"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt
	} else {

		r.ArchiveType = res.ArchiveType

		r.Uri = res.Uri

		r.Distribution = res.Distribution

		r.Components = res.Components

		r.GpgKey = res.GpgKey

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum struct {
	empty       bool     `json:"-"`
	Id          *string  `json:"id"`
	DisplayName *string  `json:"displayName"`
	BaseUrl     *string  `json:"baseUrl"`
	GpgKeys     []string `json:"gpgKeys"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum
	} else {

		r.Id = res.Id

		r.DisplayName = res.DisplayName

		r.BaseUrl = res.BaseUrl

		r.GpgKeys = res.GpgKeys

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper struct {
	empty       bool     `json:"-"`
	Id          *string  `json:"id"`
	DisplayName *string  `json:"displayName"`
	BaseUrl     *string  `json:"baseUrl"`
	GpgKeys     []string `json:"gpgKeys"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper
	} else {

		r.Id = res.Id

		r.DisplayName = res.DisplayName

		r.BaseUrl = res.BaseUrl

		r.GpgKeys = res.GpgKeys

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo struct {
	empty bool    `json:"-"`
	Name  *string `json:"name"`
	Url   *string `json:"url"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo
	} else {

		r.Name = res.Name

		r.Url = res.Url

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec struct {
	empty    bool                                                             `json:"-"`
	Validate *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate `json:"validate"`
	Enforce  *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce  `json:"enforce"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec
	} else {

		r.Validate = res.Validate

		r.Enforce = res.Enforce

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate struct {
	empty          bool                                                                            `json:"-"`
	File           *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile            `json:"file"`
	Script         *string                                                                         `json:"script"`
	Args           []string                                                                        `json:"args"`
	Interpreter    *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnum `json:"interpreter"`
	OutputFilePath *string                                                                         `json:"outputFilePath"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate
	} else {

		r.File = res.File

		r.Script = res.Script

		r.Args = res.Args

		r.Interpreter = res.Interpreter

		r.OutputFilePath = res.OutputFilePath

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile struct {
	empty         bool                                                                       `json:"-"`
	Remote        *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote `json:"remote"`
	Gcs           *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs    `json:"gcs"`
	LocalPath     *string                                                                    `json:"localPath"`
	AllowInsecure *bool                                                                      `json:"allowInsecure"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile
	} else {

		r.Remote = res.Remote

		r.Gcs = res.Gcs

		r.LocalPath = res.LocalPath

		r.AllowInsecure = res.AllowInsecure

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote struct {
	empty          bool    `json:"-"`
	Uri            *string `json:"uri"`
	Sha256Checksum *string `json:"sha256Checksum"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote
	} else {

		r.Uri = res.Uri

		r.Sha256Checksum = res.Sha256Checksum

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs struct {
	empty      bool    `json:"-"`
	Bucket     *string `json:"bucket"`
	Object     *string `json:"object"`
	Generation *int64  `json:"generation"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs
	} else {

		r.Bucket = res.Bucket

		r.Object = res.Object

		r.Generation = res.Generation

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce struct {
	empty          bool                                                                           `json:"-"`
	File           *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile            `json:"file"`
	Script         *string                                                                        `json:"script"`
	Args           []string                                                                       `json:"args"`
	Interpreter    *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnum `json:"interpreter"`
	OutputFilePath *string                                                                        `json:"outputFilePath"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce
	} else {

		r.File = res.File

		r.Script = res.Script

		r.Args = res.Args

		r.Interpreter = res.Interpreter

		r.OutputFilePath = res.OutputFilePath

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile struct {
	empty         bool                                                                      `json:"-"`
	Remote        *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote `json:"remote"`
	Gcs           *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs    `json:"gcs"`
	LocalPath     *string                                                                   `json:"localPath"`
	AllowInsecure *bool                                                                     `json:"allowInsecure"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile
	} else {

		r.Remote = res.Remote

		r.Gcs = res.Gcs

		r.LocalPath = res.LocalPath

		r.AllowInsecure = res.AllowInsecure

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote struct {
	empty          bool    `json:"-"`
	Uri            *string `json:"uri"`
	Sha256Checksum *string `json:"sha256Checksum"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote
	} else {

		r.Uri = res.Uri

		r.Sha256Checksum = res.Sha256Checksum

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs struct {
	empty      bool    `json:"-"`
	Bucket     *string `json:"bucket"`
	Object     *string `json:"object"`
	Generation *int64  `json:"generation"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs
	} else {

		r.Bucket = res.Bucket

		r.Object = res.Object

		r.Generation = res.Generation

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile struct {
	empty       bool                                                              `json:"-"`
	File        *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile      `json:"file"`
	Content     *string                                                           `json:"content"`
	Path        *string                                                           `json:"path"`
	State       *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnum `json:"state"`
	Permissions *string                                                           `json:"permissions"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile
	} else {

		r.File = res.File

		r.Content = res.Content

		r.Path = res.Path

		r.State = res.State

		r.Permissions = res.Permissions

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile struct {
	empty         bool                                                               `json:"-"`
	Remote        *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote `json:"remote"`
	Gcs           *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs    `json:"gcs"`
	LocalPath     *string                                                            `json:"localPath"`
	AllowInsecure *bool                                                              `json:"allowInsecure"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile
	} else {

		r.Remote = res.Remote

		r.Gcs = res.Gcs

		r.LocalPath = res.LocalPath

		r.AllowInsecure = res.AllowInsecure

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote struct {
	empty          bool    `json:"-"`
	Uri            *string `json:"uri"`
	Sha256Checksum *string `json:"sha256Checksum"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote
	} else {

		r.Uri = res.Uri

		r.Sha256Checksum = res.Sha256Checksum

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs struct {
	empty      bool    `json:"-"`
	Bucket     *string `json:"bucket"`
	Object     *string `json:"object"`
	Generation *int64  `json:"generation"`
}

type jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs
	} else {

		r.Bucket = res.Bucket

		r.Object = res.Object

		r.Generation = res.Generation

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs{empty: true}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentInstanceFilter struct {
	empty           bool                                              `json:"-"`
	All             *bool                                             `json:"all"`
	InclusionLabels []OSPolicyAssignmentInstanceFilterInclusionLabels `json:"inclusionLabels"`
	ExclusionLabels []OSPolicyAssignmentInstanceFilterExclusionLabels `json:"exclusionLabels"`
	Inventories     []OSPolicyAssignmentInstanceFilterInventories     `json:"inventories"`
}

type jsonOSPolicyAssignmentInstanceFilter OSPolicyAssignmentInstanceFilter

func (r *OSPolicyAssignmentInstanceFilter) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentInstanceFilter
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentInstanceFilter
	} else {

		r.All = res.All

		r.InclusionLabels = res.InclusionLabels

		r.ExclusionLabels = res.ExclusionLabels

		r.Inventories = res.Inventories

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentInstanceFilter is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentInstanceFilter *OSPolicyAssignmentInstanceFilter = &OSPolicyAssignmentInstanceFilter{empty: true}

func (r *OSPolicyAssignmentInstanceFilter) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentInstanceFilter) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentInstanceFilter) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentInstanceFilterInclusionLabels struct {
	empty  bool              `json:"-"`
	Labels map[string]string `json:"labels"`
}

type jsonOSPolicyAssignmentInstanceFilterInclusionLabels OSPolicyAssignmentInstanceFilterInclusionLabels

func (r *OSPolicyAssignmentInstanceFilterInclusionLabels) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentInstanceFilterInclusionLabels
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentInstanceFilterInclusionLabels
	} else {

		r.Labels = res.Labels

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentInstanceFilterInclusionLabels is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentInstanceFilterInclusionLabels *OSPolicyAssignmentInstanceFilterInclusionLabels = &OSPolicyAssignmentInstanceFilterInclusionLabels{empty: true}

func (r *OSPolicyAssignmentInstanceFilterInclusionLabels) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentInstanceFilterInclusionLabels) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentInstanceFilterInclusionLabels) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentInstanceFilterExclusionLabels struct {
	empty  bool              `json:"-"`
	Labels map[string]string `json:"labels"`
}

type jsonOSPolicyAssignmentInstanceFilterExclusionLabels OSPolicyAssignmentInstanceFilterExclusionLabels

func (r *OSPolicyAssignmentInstanceFilterExclusionLabels) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentInstanceFilterExclusionLabels
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentInstanceFilterExclusionLabels
	} else {

		r.Labels = res.Labels

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentInstanceFilterExclusionLabels is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentInstanceFilterExclusionLabels *OSPolicyAssignmentInstanceFilterExclusionLabels = &OSPolicyAssignmentInstanceFilterExclusionLabels{empty: true}

func (r *OSPolicyAssignmentInstanceFilterExclusionLabels) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentInstanceFilterExclusionLabels) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentInstanceFilterExclusionLabels) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentInstanceFilterInventories struct {
	empty       bool    `json:"-"`
	OSShortName *string `json:"osShortName"`
	OSVersion   *string `json:"osVersion"`
}

type jsonOSPolicyAssignmentInstanceFilterInventories OSPolicyAssignmentInstanceFilterInventories

func (r *OSPolicyAssignmentInstanceFilterInventories) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentInstanceFilterInventories
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentInstanceFilterInventories
	} else {

		r.OSShortName = res.OSShortName

		r.OSVersion = res.OSVersion

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentInstanceFilterInventories is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentInstanceFilterInventories *OSPolicyAssignmentInstanceFilterInventories = &OSPolicyAssignmentInstanceFilterInventories{empty: true}

func (r *OSPolicyAssignmentInstanceFilterInventories) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentInstanceFilterInventories) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentInstanceFilterInventories) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentRollout struct {
	empty            bool                                       `json:"-"`
	DisruptionBudget *OSPolicyAssignmentRolloutDisruptionBudget `json:"disruptionBudget"`
	MinWaitDuration  *string                                    `json:"minWaitDuration"`
}

type jsonOSPolicyAssignmentRollout OSPolicyAssignmentRollout

func (r *OSPolicyAssignmentRollout) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentRollout
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentRollout
	} else {

		r.DisruptionBudget = res.DisruptionBudget

		r.MinWaitDuration = res.MinWaitDuration

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentRollout is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentRollout *OSPolicyAssignmentRollout = &OSPolicyAssignmentRollout{empty: true}

func (r *OSPolicyAssignmentRollout) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentRollout) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentRollout) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type OSPolicyAssignmentRolloutDisruptionBudget struct {
	empty   bool   `json:"-"`
	Fixed   *int64 `json:"fixed"`
	Percent *int64 `json:"percent"`
}

type jsonOSPolicyAssignmentRolloutDisruptionBudget OSPolicyAssignmentRolloutDisruptionBudget

func (r *OSPolicyAssignmentRolloutDisruptionBudget) UnmarshalJSON(data []byte) error {
	var res jsonOSPolicyAssignmentRolloutDisruptionBudget
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyOSPolicyAssignmentRolloutDisruptionBudget
	} else {

		r.Fixed = res.Fixed

		r.Percent = res.Percent

	}
	return nil
}

// This object is used to assert a desired state where this OSPolicyAssignmentRolloutDisruptionBudget is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyOSPolicyAssignmentRolloutDisruptionBudget *OSPolicyAssignmentRolloutDisruptionBudget = &OSPolicyAssignmentRolloutDisruptionBudget{empty: true}

func (r *OSPolicyAssignmentRolloutDisruptionBudget) Empty() bool {
	return r.empty
}

func (r *OSPolicyAssignmentRolloutDisruptionBudget) String() string {
	return dcl.SprintResource(r)
}

func (r *OSPolicyAssignmentRolloutDisruptionBudget) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *OSPolicyAssignment) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "os_config",
		Type:    "OSPolicyAssignment",
		Version: "osconfig",
	}
}

func (r *OSPolicyAssignment) ID() (string, error) {
	if err := extractOSPolicyAssignmentFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name":                 dcl.ValueOrEmptyString(nr.Name),
		"description":          dcl.ValueOrEmptyString(nr.Description),
		"os_policies":          dcl.ValueOrEmptyString(nr.OSPolicies),
		"instance_filter":      dcl.ValueOrEmptyString(nr.InstanceFilter),
		"rollout":              dcl.ValueOrEmptyString(nr.Rollout),
		"revision_id":          dcl.ValueOrEmptyString(nr.RevisionId),
		"revision_create_time": dcl.ValueOrEmptyString(nr.RevisionCreateTime),
		"etag":                 dcl.ValueOrEmptyString(nr.Etag),
		"rollout_state":        dcl.ValueOrEmptyString(nr.RolloutState),
		"baseline":             dcl.ValueOrEmptyString(nr.Baseline),
		"deleted":              dcl.ValueOrEmptyString(nr.Deleted),
		"reconciling":          dcl.ValueOrEmptyString(nr.Reconciling),
		"uid":                  dcl.ValueOrEmptyString(nr.Uid),
		"project":              dcl.ValueOrEmptyString(nr.Project),
		"location":             dcl.ValueOrEmptyString(nr.Location),
		"skip_await_rollout":   dcl.ValueOrEmptyString(nr.SkipAwaitRollout),
	}
	return dcl.Nprintf("projects/{{project}}/locations/{{location}}/osPolicyAssignments/{{name}}", params), nil
}

const OSPolicyAssignmentMaxPage = -1

type OSPolicyAssignmentList struct {
	Items []*OSPolicyAssignment

	nextToken string

	pageSize int32

	resource *OSPolicyAssignment
}

func (l *OSPolicyAssignmentList) HasNext() bool {
	return l.nextToken != ""
}

func (l *OSPolicyAssignmentList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listOSPolicyAssignment(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListOSPolicyAssignment(ctx context.Context, project, location string) (*OSPolicyAssignmentList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListOSPolicyAssignmentWithMaxResults(ctx, project, location, OSPolicyAssignmentMaxPage)

}

func (c *Client) ListOSPolicyAssignmentWithMaxResults(ctx context.Context, project, location string, pageSize int32) (*OSPolicyAssignmentList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &OSPolicyAssignment{
		Project:  &project,
		Location: &location,
	}
	items, token, err := c.listOSPolicyAssignment(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &OSPolicyAssignmentList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetOSPolicyAssignment(ctx context.Context, r *OSPolicyAssignment) (*OSPolicyAssignment, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractOSPolicyAssignmentFields(r)

	b, err := c.getOSPolicyAssignmentRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalOSPolicyAssignment(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Location = r.Location
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeOSPolicyAssignmentNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractOSPolicyAssignmentFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteOSPolicyAssignment(ctx context.Context, r *OSPolicyAssignment) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("OSPolicyAssignment resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting OSPolicyAssignment...")
	deleteOp := deleteOSPolicyAssignmentOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllOSPolicyAssignment deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllOSPolicyAssignment(ctx context.Context, project, location string, filter func(*OSPolicyAssignment) bool) error {
	listObj, err := c.ListOSPolicyAssignment(ctx, project, location)
	if err != nil {
		return err
	}

	err = c.deleteAllOSPolicyAssignment(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllOSPolicyAssignment(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyOSPolicyAssignment(ctx context.Context, rawDesired *OSPolicyAssignment, opts ...dcl.ApplyOption) (*OSPolicyAssignment, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *OSPolicyAssignment
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyOSPolicyAssignmentHelper(c, ctx, rawDesired, opts...)
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

func applyOSPolicyAssignmentHelper(c *Client, ctx context.Context, rawDesired *OSPolicyAssignment, opts ...dcl.ApplyOption) (*OSPolicyAssignment, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyOSPolicyAssignment...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractOSPolicyAssignmentFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.oSPolicyAssignmentDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToOSPolicyAssignmentDiffs(c.Config, fieldDiffs, opts)
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
	var ops []oSPolicyAssignmentApiOperation
	if create {
		ops = append(ops, &createOSPolicyAssignmentOperation{})
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
	return applyOSPolicyAssignmentDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyOSPolicyAssignmentDiff(c *Client, ctx context.Context, desired *OSPolicyAssignment, rawDesired *OSPolicyAssignment, ops []oSPolicyAssignmentApiOperation, opts ...dcl.ApplyOption) (*OSPolicyAssignment, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetOSPolicyAssignment(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createOSPolicyAssignmentOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapOSPolicyAssignment(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeOSPolicyAssignmentNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeOSPolicyAssignmentNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeOSPolicyAssignmentDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractOSPolicyAssignmentFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractOSPolicyAssignmentFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffOSPolicyAssignment(c, newDesired, newState)
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
