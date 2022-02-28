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

type PatchDeployment struct {
	Name              *string                           `json:"name"`
	Description       *string                           `json:"description"`
	InstanceFilter    *PatchDeploymentInstanceFilter    `json:"instanceFilter"`
	PatchConfig       *PatchDeploymentPatchConfig       `json:"patchConfig"`
	Duration          *string                           `json:"duration"`
	OneTimeSchedule   *PatchDeploymentOneTimeSchedule   `json:"oneTimeSchedule"`
	RecurringSchedule *PatchDeploymentRecurringSchedule `json:"recurringSchedule"`
	CreateTime        *string                           `json:"createTime"`
	UpdateTime        *string                           `json:"updateTime"`
	LastExecuteTime   *string                           `json:"lastExecuteTime"`
	Rollout           *PatchDeploymentRollout           `json:"rollout"`
	Project           *string                           `json:"project"`
}

func (r *PatchDeployment) String() string {
	return dcl.SprintResource(r)
}

// The enum PatchDeploymentPatchConfigRebootConfigEnum.
type PatchDeploymentPatchConfigRebootConfigEnum string

// PatchDeploymentPatchConfigRebootConfigEnumRef returns a *PatchDeploymentPatchConfigRebootConfigEnum with the value of string s
// If the empty string is provided, nil is returned.
func PatchDeploymentPatchConfigRebootConfigEnumRef(s string) *PatchDeploymentPatchConfigRebootConfigEnum {
	v := PatchDeploymentPatchConfigRebootConfigEnum(s)
	return &v
}

func (v PatchDeploymentPatchConfigRebootConfigEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REBOOT_CONFIG_UNSPECIFIED", "DEFAULT", "ALWAYS", "NEVER"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "PatchDeploymentPatchConfigRebootConfigEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum PatchDeploymentPatchConfigAptTypeEnum.
type PatchDeploymentPatchConfigAptTypeEnum string

// PatchDeploymentPatchConfigAptTypeEnumRef returns a *PatchDeploymentPatchConfigAptTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func PatchDeploymentPatchConfigAptTypeEnumRef(s string) *PatchDeploymentPatchConfigAptTypeEnum {
	v := PatchDeploymentPatchConfigAptTypeEnum(s)
	return &v
}

func (v PatchDeploymentPatchConfigAptTypeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"TYPE_UNSPECIFIED", "DIST", "UPGRADE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "PatchDeploymentPatchConfigAptTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum PatchDeploymentPatchConfigWindowsUpdateClassificationsEnum.
type PatchDeploymentPatchConfigWindowsUpdateClassificationsEnum string

// PatchDeploymentPatchConfigWindowsUpdateClassificationsEnumRef returns a *PatchDeploymentPatchConfigWindowsUpdateClassificationsEnum with the value of string s
// If the empty string is provided, nil is returned.
func PatchDeploymentPatchConfigWindowsUpdateClassificationsEnumRef(s string) *PatchDeploymentPatchConfigWindowsUpdateClassificationsEnum {
	v := PatchDeploymentPatchConfigWindowsUpdateClassificationsEnum(s)
	return &v
}

func (v PatchDeploymentPatchConfigWindowsUpdateClassificationsEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"CLASSIFICATION_UNSPECIFIED", "CRITICAL", "SECURITY", "DEFINITION", "DRIVER", "FEATURE_PACK", "SERVICE_PACK", "TOOL", "UPDATE_ROLLUP", "UPDATE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "PatchDeploymentPatchConfigWindowsUpdateClassificationsEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum PatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnum.
type PatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnum string

// PatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnumRef returns a *PatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnum with the value of string s
// If the empty string is provided, nil is returned.
func PatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnumRef(s string) *PatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnum {
	v := PatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnum(s)
	return &v
}

func (v PatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnum) Validate() error {
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
		Enum:  "PatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum PatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnum.
type PatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnum string

// PatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnumRef returns a *PatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnum with the value of string s
// If the empty string is provided, nil is returned.
func PatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnumRef(s string) *PatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnum {
	v := PatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnum(s)
	return &v
}

func (v PatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnum) Validate() error {
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
		Enum:  "PatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum PatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnum.
type PatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnum string

// PatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnumRef returns a *PatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnum with the value of string s
// If the empty string is provided, nil is returned.
func PatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnumRef(s string) *PatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnum {
	v := PatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnum(s)
	return &v
}

func (v PatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnum) Validate() error {
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
		Enum:  "PatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum PatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnum.
type PatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnum string

// PatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnumRef returns a *PatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnum with the value of string s
// If the empty string is provided, nil is returned.
func PatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnumRef(s string) *PatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnum {
	v := PatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnum(s)
	return &v
}

func (v PatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnum) Validate() error {
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
		Enum:  "PatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum PatchDeploymentRecurringScheduleFrequencyEnum.
type PatchDeploymentRecurringScheduleFrequencyEnum string

// PatchDeploymentRecurringScheduleFrequencyEnumRef returns a *PatchDeploymentRecurringScheduleFrequencyEnum with the value of string s
// If the empty string is provided, nil is returned.
func PatchDeploymentRecurringScheduleFrequencyEnumRef(s string) *PatchDeploymentRecurringScheduleFrequencyEnum {
	v := PatchDeploymentRecurringScheduleFrequencyEnum(s)
	return &v
}

func (v PatchDeploymentRecurringScheduleFrequencyEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"FREQUENCY_UNSPECIFIED", "WEEKLY", "MONTHLY", "DAILY"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "PatchDeploymentRecurringScheduleFrequencyEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum PatchDeploymentRecurringScheduleWeeklyDayOfWeekEnum.
type PatchDeploymentRecurringScheduleWeeklyDayOfWeekEnum string

// PatchDeploymentRecurringScheduleWeeklyDayOfWeekEnumRef returns a *PatchDeploymentRecurringScheduleWeeklyDayOfWeekEnum with the value of string s
// If the empty string is provided, nil is returned.
func PatchDeploymentRecurringScheduleWeeklyDayOfWeekEnumRef(s string) *PatchDeploymentRecurringScheduleWeeklyDayOfWeekEnum {
	v := PatchDeploymentRecurringScheduleWeeklyDayOfWeekEnum(s)
	return &v
}

func (v PatchDeploymentRecurringScheduleWeeklyDayOfWeekEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"DAY_OF_WEEK_UNSPECIFIED", "MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY", "SATURDAY", "SUNDAY"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "PatchDeploymentRecurringScheduleWeeklyDayOfWeekEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnum.
type PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnum string

// PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnumRef returns a *PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnum with the value of string s
// If the empty string is provided, nil is returned.
func PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnumRef(s string) *PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnum {
	v := PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnum(s)
	return &v
}

func (v PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"DAY_OF_WEEK_UNSPECIFIED", "MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY", "SATURDAY", "SUNDAY"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum PatchDeploymentRolloutModeEnum.
type PatchDeploymentRolloutModeEnum string

// PatchDeploymentRolloutModeEnumRef returns a *PatchDeploymentRolloutModeEnum with the value of string s
// If the empty string is provided, nil is returned.
func PatchDeploymentRolloutModeEnumRef(s string) *PatchDeploymentRolloutModeEnum {
	v := PatchDeploymentRolloutModeEnum(s)
	return &v
}

func (v PatchDeploymentRolloutModeEnum) Validate() error {
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
		Enum:  "PatchDeploymentRolloutModeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

type PatchDeploymentInstanceFilter struct {
	empty                bool                                       `json:"-"`
	All                  *bool                                      `json:"all"`
	GroupLabels          []PatchDeploymentInstanceFilterGroupLabels `json:"groupLabels"`
	Zones                []string                                   `json:"zones"`
	Instances            []string                                   `json:"instances"`
	InstanceNamePrefixes []string                                   `json:"instanceNamePrefixes"`
}

type jsonPatchDeploymentInstanceFilter PatchDeploymentInstanceFilter

func (r *PatchDeploymentInstanceFilter) UnmarshalJSON(data []byte) error {
	var res jsonPatchDeploymentInstanceFilter
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPatchDeploymentInstanceFilter
	} else {

		r.All = res.All

		r.GroupLabels = res.GroupLabels

		r.Zones = res.Zones

		r.Instances = res.Instances

		r.InstanceNamePrefixes = res.InstanceNamePrefixes

	}
	return nil
}

// This object is used to assert a desired state where this PatchDeploymentInstanceFilter is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyPatchDeploymentInstanceFilter *PatchDeploymentInstanceFilter = &PatchDeploymentInstanceFilter{empty: true}

func (r *PatchDeploymentInstanceFilter) Empty() bool {
	return r.empty
}

func (r *PatchDeploymentInstanceFilter) String() string {
	return dcl.SprintResource(r)
}

func (r *PatchDeploymentInstanceFilter) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PatchDeploymentInstanceFilterGroupLabels struct {
	empty  bool              `json:"-"`
	Labels map[string]string `json:"labels"`
}

type jsonPatchDeploymentInstanceFilterGroupLabels PatchDeploymentInstanceFilterGroupLabels

func (r *PatchDeploymentInstanceFilterGroupLabels) UnmarshalJSON(data []byte) error {
	var res jsonPatchDeploymentInstanceFilterGroupLabels
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPatchDeploymentInstanceFilterGroupLabels
	} else {

		r.Labels = res.Labels

	}
	return nil
}

// This object is used to assert a desired state where this PatchDeploymentInstanceFilterGroupLabels is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyPatchDeploymentInstanceFilterGroupLabels *PatchDeploymentInstanceFilterGroupLabels = &PatchDeploymentInstanceFilterGroupLabels{empty: true}

func (r *PatchDeploymentInstanceFilterGroupLabels) Empty() bool {
	return r.empty
}

func (r *PatchDeploymentInstanceFilterGroupLabels) String() string {
	return dcl.SprintResource(r)
}

func (r *PatchDeploymentInstanceFilterGroupLabels) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PatchDeploymentPatchConfig struct {
	empty         bool                                        `json:"-"`
	RebootConfig  *PatchDeploymentPatchConfigRebootConfigEnum `json:"rebootConfig"`
	Apt           *PatchDeploymentPatchConfigApt              `json:"apt"`
	Yum           *PatchDeploymentPatchConfigYum              `json:"yum"`
	Goo           *PatchDeploymentPatchConfigGoo              `json:"goo"`
	Zypper        *PatchDeploymentPatchConfigZypper           `json:"zypper"`
	WindowsUpdate *PatchDeploymentPatchConfigWindowsUpdate    `json:"windowsUpdate"`
	PreStep       *PatchDeploymentPatchConfigPreStep          `json:"preStep"`
	PostStep      *PatchDeploymentPatchConfigPostStep         `json:"postStep"`
}

type jsonPatchDeploymentPatchConfig PatchDeploymentPatchConfig

func (r *PatchDeploymentPatchConfig) UnmarshalJSON(data []byte) error {
	var res jsonPatchDeploymentPatchConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPatchDeploymentPatchConfig
	} else {

		r.RebootConfig = res.RebootConfig

		r.Apt = res.Apt

		r.Yum = res.Yum

		r.Goo = res.Goo

		r.Zypper = res.Zypper

		r.WindowsUpdate = res.WindowsUpdate

		r.PreStep = res.PreStep

		r.PostStep = res.PostStep

	}
	return nil
}

// This object is used to assert a desired state where this PatchDeploymentPatchConfig is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyPatchDeploymentPatchConfig *PatchDeploymentPatchConfig = &PatchDeploymentPatchConfig{empty: true}

func (r *PatchDeploymentPatchConfig) Empty() bool {
	return r.empty
}

func (r *PatchDeploymentPatchConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *PatchDeploymentPatchConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PatchDeploymentPatchConfigApt struct {
	empty             bool                                   `json:"-"`
	Type              *PatchDeploymentPatchConfigAptTypeEnum `json:"type"`
	Excludes          []string                               `json:"excludes"`
	ExclusivePackages []string                               `json:"exclusivePackages"`
}

type jsonPatchDeploymentPatchConfigApt PatchDeploymentPatchConfigApt

func (r *PatchDeploymentPatchConfigApt) UnmarshalJSON(data []byte) error {
	var res jsonPatchDeploymentPatchConfigApt
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPatchDeploymentPatchConfigApt
	} else {

		r.Type = res.Type

		r.Excludes = res.Excludes

		r.ExclusivePackages = res.ExclusivePackages

	}
	return nil
}

// This object is used to assert a desired state where this PatchDeploymentPatchConfigApt is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyPatchDeploymentPatchConfigApt *PatchDeploymentPatchConfigApt = &PatchDeploymentPatchConfigApt{empty: true}

func (r *PatchDeploymentPatchConfigApt) Empty() bool {
	return r.empty
}

func (r *PatchDeploymentPatchConfigApt) String() string {
	return dcl.SprintResource(r)
}

func (r *PatchDeploymentPatchConfigApt) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PatchDeploymentPatchConfigYum struct {
	empty             bool     `json:"-"`
	Security          *bool    `json:"security"`
	Minimal           *bool    `json:"minimal"`
	Excludes          []string `json:"excludes"`
	ExclusivePackages []string `json:"exclusivePackages"`
}

type jsonPatchDeploymentPatchConfigYum PatchDeploymentPatchConfigYum

func (r *PatchDeploymentPatchConfigYum) UnmarshalJSON(data []byte) error {
	var res jsonPatchDeploymentPatchConfigYum
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPatchDeploymentPatchConfigYum
	} else {

		r.Security = res.Security

		r.Minimal = res.Minimal

		r.Excludes = res.Excludes

		r.ExclusivePackages = res.ExclusivePackages

	}
	return nil
}

// This object is used to assert a desired state where this PatchDeploymentPatchConfigYum is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyPatchDeploymentPatchConfigYum *PatchDeploymentPatchConfigYum = &PatchDeploymentPatchConfigYum{empty: true}

func (r *PatchDeploymentPatchConfigYum) Empty() bool {
	return r.empty
}

func (r *PatchDeploymentPatchConfigYum) String() string {
	return dcl.SprintResource(r)
}

func (r *PatchDeploymentPatchConfigYum) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PatchDeploymentPatchConfigGoo struct {
	empty bool `json:"-"`
}

type jsonPatchDeploymentPatchConfigGoo PatchDeploymentPatchConfigGoo

func (r *PatchDeploymentPatchConfigGoo) UnmarshalJSON(data []byte) error {
	var res jsonPatchDeploymentPatchConfigGoo
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPatchDeploymentPatchConfigGoo
	} else {

	}
	return nil
}

// This object is used to assert a desired state where this PatchDeploymentPatchConfigGoo is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyPatchDeploymentPatchConfigGoo *PatchDeploymentPatchConfigGoo = &PatchDeploymentPatchConfigGoo{empty: true}

func (r *PatchDeploymentPatchConfigGoo) Empty() bool {
	return r.empty
}

func (r *PatchDeploymentPatchConfigGoo) String() string {
	return dcl.SprintResource(r)
}

func (r *PatchDeploymentPatchConfigGoo) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PatchDeploymentPatchConfigZypper struct {
	empty            bool     `json:"-"`
	WithOptional     *bool    `json:"withOptional"`
	WithUpdate       *bool    `json:"withUpdate"`
	Categories       []string `json:"categories"`
	Severities       []string `json:"severities"`
	Excludes         []string `json:"excludes"`
	ExclusivePatches []string `json:"exclusivePatches"`
}

type jsonPatchDeploymentPatchConfigZypper PatchDeploymentPatchConfigZypper

func (r *PatchDeploymentPatchConfigZypper) UnmarshalJSON(data []byte) error {
	var res jsonPatchDeploymentPatchConfigZypper
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPatchDeploymentPatchConfigZypper
	} else {

		r.WithOptional = res.WithOptional

		r.WithUpdate = res.WithUpdate

		r.Categories = res.Categories

		r.Severities = res.Severities

		r.Excludes = res.Excludes

		r.ExclusivePatches = res.ExclusivePatches

	}
	return nil
}

// This object is used to assert a desired state where this PatchDeploymentPatchConfigZypper is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyPatchDeploymentPatchConfigZypper *PatchDeploymentPatchConfigZypper = &PatchDeploymentPatchConfigZypper{empty: true}

func (r *PatchDeploymentPatchConfigZypper) Empty() bool {
	return r.empty
}

func (r *PatchDeploymentPatchConfigZypper) String() string {
	return dcl.SprintResource(r)
}

func (r *PatchDeploymentPatchConfigZypper) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PatchDeploymentPatchConfigWindowsUpdate struct {
	empty            bool                                                         `json:"-"`
	Classifications  []PatchDeploymentPatchConfigWindowsUpdateClassificationsEnum `json:"classifications"`
	Excludes         []string                                                     `json:"excludes"`
	ExclusivePatches []string                                                     `json:"exclusivePatches"`
}

type jsonPatchDeploymentPatchConfigWindowsUpdate PatchDeploymentPatchConfigWindowsUpdate

func (r *PatchDeploymentPatchConfigWindowsUpdate) UnmarshalJSON(data []byte) error {
	var res jsonPatchDeploymentPatchConfigWindowsUpdate
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPatchDeploymentPatchConfigWindowsUpdate
	} else {

		r.Classifications = res.Classifications

		r.Excludes = res.Excludes

		r.ExclusivePatches = res.ExclusivePatches

	}
	return nil
}

// This object is used to assert a desired state where this PatchDeploymentPatchConfigWindowsUpdate is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyPatchDeploymentPatchConfigWindowsUpdate *PatchDeploymentPatchConfigWindowsUpdate = &PatchDeploymentPatchConfigWindowsUpdate{empty: true}

func (r *PatchDeploymentPatchConfigWindowsUpdate) Empty() bool {
	return r.empty
}

func (r *PatchDeploymentPatchConfigWindowsUpdate) String() string {
	return dcl.SprintResource(r)
}

func (r *PatchDeploymentPatchConfigWindowsUpdate) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PatchDeploymentPatchConfigPreStep struct {
	empty                 bool                                                    `json:"-"`
	LinuxExecStepConfig   *PatchDeploymentPatchConfigPreStepLinuxExecStepConfig   `json:"linuxExecStepConfig"`
	WindowsExecStepConfig *PatchDeploymentPatchConfigPreStepWindowsExecStepConfig `json:"windowsExecStepConfig"`
}

type jsonPatchDeploymentPatchConfigPreStep PatchDeploymentPatchConfigPreStep

func (r *PatchDeploymentPatchConfigPreStep) UnmarshalJSON(data []byte) error {
	var res jsonPatchDeploymentPatchConfigPreStep
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPatchDeploymentPatchConfigPreStep
	} else {

		r.LinuxExecStepConfig = res.LinuxExecStepConfig

		r.WindowsExecStepConfig = res.WindowsExecStepConfig

	}
	return nil
}

// This object is used to assert a desired state where this PatchDeploymentPatchConfigPreStep is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyPatchDeploymentPatchConfigPreStep *PatchDeploymentPatchConfigPreStep = &PatchDeploymentPatchConfigPreStep{empty: true}

func (r *PatchDeploymentPatchConfigPreStep) Empty() bool {
	return r.empty
}

func (r *PatchDeploymentPatchConfigPreStep) String() string {
	return dcl.SprintResource(r)
}

func (r *PatchDeploymentPatchConfigPreStep) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PatchDeploymentPatchConfigPreStepLinuxExecStepConfig struct {
	empty               bool                                                                 `json:"-"`
	LocalPath           *string                                                              `json:"localPath"`
	GcsObject           *PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject       `json:"gcsObject"`
	AllowedSuccessCodes []int64                                                              `json:"allowedSuccessCodes"`
	Interpreter         *PatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnum `json:"interpreter"`
}

type jsonPatchDeploymentPatchConfigPreStepLinuxExecStepConfig PatchDeploymentPatchConfigPreStepLinuxExecStepConfig

func (r *PatchDeploymentPatchConfigPreStepLinuxExecStepConfig) UnmarshalJSON(data []byte) error {
	var res jsonPatchDeploymentPatchConfigPreStepLinuxExecStepConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPatchDeploymentPatchConfigPreStepLinuxExecStepConfig
	} else {

		r.LocalPath = res.LocalPath

		r.GcsObject = res.GcsObject

		r.AllowedSuccessCodes = res.AllowedSuccessCodes

		r.Interpreter = res.Interpreter

	}
	return nil
}

// This object is used to assert a desired state where this PatchDeploymentPatchConfigPreStepLinuxExecStepConfig is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyPatchDeploymentPatchConfigPreStepLinuxExecStepConfig *PatchDeploymentPatchConfigPreStepLinuxExecStepConfig = &PatchDeploymentPatchConfigPreStepLinuxExecStepConfig{empty: true}

func (r *PatchDeploymentPatchConfigPreStepLinuxExecStepConfig) Empty() bool {
	return r.empty
}

func (r *PatchDeploymentPatchConfigPreStepLinuxExecStepConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *PatchDeploymentPatchConfigPreStepLinuxExecStepConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject struct {
	empty            bool    `json:"-"`
	Bucket           *string `json:"bucket"`
	Object           *string `json:"object"`
	GenerationNumber *int64  `json:"generationNumber"`
}

type jsonPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject

func (r *PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject) UnmarshalJSON(data []byte) error {
	var res jsonPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject
	} else {

		r.Bucket = res.Bucket

		r.Object = res.Object

		r.GenerationNumber = res.GenerationNumber

	}
	return nil
}

// This object is used to assert a desired state where this PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject *PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject = &PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject{empty: true}

func (r *PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject) Empty() bool {
	return r.empty
}

func (r *PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject) String() string {
	return dcl.SprintResource(r)
}

func (r *PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PatchDeploymentPatchConfigPreStepWindowsExecStepConfig struct {
	empty               bool                                                                   `json:"-"`
	LocalPath           *string                                                                `json:"localPath"`
	GcsObject           *PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject       `json:"gcsObject"`
	AllowedSuccessCodes []int64                                                                `json:"allowedSuccessCodes"`
	Interpreter         *PatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnum `json:"interpreter"`
}

type jsonPatchDeploymentPatchConfigPreStepWindowsExecStepConfig PatchDeploymentPatchConfigPreStepWindowsExecStepConfig

func (r *PatchDeploymentPatchConfigPreStepWindowsExecStepConfig) UnmarshalJSON(data []byte) error {
	var res jsonPatchDeploymentPatchConfigPreStepWindowsExecStepConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPatchDeploymentPatchConfigPreStepWindowsExecStepConfig
	} else {

		r.LocalPath = res.LocalPath

		r.GcsObject = res.GcsObject

		r.AllowedSuccessCodes = res.AllowedSuccessCodes

		r.Interpreter = res.Interpreter

	}
	return nil
}

// This object is used to assert a desired state where this PatchDeploymentPatchConfigPreStepWindowsExecStepConfig is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyPatchDeploymentPatchConfigPreStepWindowsExecStepConfig *PatchDeploymentPatchConfigPreStepWindowsExecStepConfig = &PatchDeploymentPatchConfigPreStepWindowsExecStepConfig{empty: true}

func (r *PatchDeploymentPatchConfigPreStepWindowsExecStepConfig) Empty() bool {
	return r.empty
}

func (r *PatchDeploymentPatchConfigPreStepWindowsExecStepConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *PatchDeploymentPatchConfigPreStepWindowsExecStepConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject struct {
	empty            bool    `json:"-"`
	Bucket           *string `json:"bucket"`
	Object           *string `json:"object"`
	GenerationNumber *int64  `json:"generationNumber"`
}

type jsonPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject

func (r *PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject) UnmarshalJSON(data []byte) error {
	var res jsonPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject
	} else {

		r.Bucket = res.Bucket

		r.Object = res.Object

		r.GenerationNumber = res.GenerationNumber

	}
	return nil
}

// This object is used to assert a desired state where this PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject *PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject = &PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject{empty: true}

func (r *PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject) Empty() bool {
	return r.empty
}

func (r *PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject) String() string {
	return dcl.SprintResource(r)
}

func (r *PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PatchDeploymentPatchConfigPostStep struct {
	empty                 bool                                                     `json:"-"`
	LinuxExecStepConfig   *PatchDeploymentPatchConfigPostStepLinuxExecStepConfig   `json:"linuxExecStepConfig"`
	WindowsExecStepConfig *PatchDeploymentPatchConfigPostStepWindowsExecStepConfig `json:"windowsExecStepConfig"`
}

type jsonPatchDeploymentPatchConfigPostStep PatchDeploymentPatchConfigPostStep

func (r *PatchDeploymentPatchConfigPostStep) UnmarshalJSON(data []byte) error {
	var res jsonPatchDeploymentPatchConfigPostStep
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPatchDeploymentPatchConfigPostStep
	} else {

		r.LinuxExecStepConfig = res.LinuxExecStepConfig

		r.WindowsExecStepConfig = res.WindowsExecStepConfig

	}
	return nil
}

// This object is used to assert a desired state where this PatchDeploymentPatchConfigPostStep is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyPatchDeploymentPatchConfigPostStep *PatchDeploymentPatchConfigPostStep = &PatchDeploymentPatchConfigPostStep{empty: true}

func (r *PatchDeploymentPatchConfigPostStep) Empty() bool {
	return r.empty
}

func (r *PatchDeploymentPatchConfigPostStep) String() string {
	return dcl.SprintResource(r)
}

func (r *PatchDeploymentPatchConfigPostStep) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PatchDeploymentPatchConfigPostStepLinuxExecStepConfig struct {
	empty               bool                                                                  `json:"-"`
	LocalPath           *string                                                               `json:"localPath"`
	GcsObject           *PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject       `json:"gcsObject"`
	AllowedSuccessCodes []int64                                                               `json:"allowedSuccessCodes"`
	Interpreter         *PatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnum `json:"interpreter"`
}

type jsonPatchDeploymentPatchConfigPostStepLinuxExecStepConfig PatchDeploymentPatchConfigPostStepLinuxExecStepConfig

func (r *PatchDeploymentPatchConfigPostStepLinuxExecStepConfig) UnmarshalJSON(data []byte) error {
	var res jsonPatchDeploymentPatchConfigPostStepLinuxExecStepConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPatchDeploymentPatchConfigPostStepLinuxExecStepConfig
	} else {

		r.LocalPath = res.LocalPath

		r.GcsObject = res.GcsObject

		r.AllowedSuccessCodes = res.AllowedSuccessCodes

		r.Interpreter = res.Interpreter

	}
	return nil
}

// This object is used to assert a desired state where this PatchDeploymentPatchConfigPostStepLinuxExecStepConfig is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyPatchDeploymentPatchConfigPostStepLinuxExecStepConfig *PatchDeploymentPatchConfigPostStepLinuxExecStepConfig = &PatchDeploymentPatchConfigPostStepLinuxExecStepConfig{empty: true}

func (r *PatchDeploymentPatchConfigPostStepLinuxExecStepConfig) Empty() bool {
	return r.empty
}

func (r *PatchDeploymentPatchConfigPostStepLinuxExecStepConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *PatchDeploymentPatchConfigPostStepLinuxExecStepConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject struct {
	empty            bool    `json:"-"`
	Bucket           *string `json:"bucket"`
	Object           *string `json:"object"`
	GenerationNumber *int64  `json:"generationNumber"`
}

type jsonPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject

func (r *PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject) UnmarshalJSON(data []byte) error {
	var res jsonPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject
	} else {

		r.Bucket = res.Bucket

		r.Object = res.Object

		r.GenerationNumber = res.GenerationNumber

	}
	return nil
}

// This object is used to assert a desired state where this PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject *PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject = &PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject{empty: true}

func (r *PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject) Empty() bool {
	return r.empty
}

func (r *PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject) String() string {
	return dcl.SprintResource(r)
}

func (r *PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PatchDeploymentPatchConfigPostStepWindowsExecStepConfig struct {
	empty               bool                                                                    `json:"-"`
	LocalPath           *string                                                                 `json:"localPath"`
	GcsObject           *PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject       `json:"gcsObject"`
	AllowedSuccessCodes []int64                                                                 `json:"allowedSuccessCodes"`
	Interpreter         *PatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnum `json:"interpreter"`
}

type jsonPatchDeploymentPatchConfigPostStepWindowsExecStepConfig PatchDeploymentPatchConfigPostStepWindowsExecStepConfig

func (r *PatchDeploymentPatchConfigPostStepWindowsExecStepConfig) UnmarshalJSON(data []byte) error {
	var res jsonPatchDeploymentPatchConfigPostStepWindowsExecStepConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPatchDeploymentPatchConfigPostStepWindowsExecStepConfig
	} else {

		r.LocalPath = res.LocalPath

		r.GcsObject = res.GcsObject

		r.AllowedSuccessCodes = res.AllowedSuccessCodes

		r.Interpreter = res.Interpreter

	}
	return nil
}

// This object is used to assert a desired state where this PatchDeploymentPatchConfigPostStepWindowsExecStepConfig is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyPatchDeploymentPatchConfigPostStepWindowsExecStepConfig *PatchDeploymentPatchConfigPostStepWindowsExecStepConfig = &PatchDeploymentPatchConfigPostStepWindowsExecStepConfig{empty: true}

func (r *PatchDeploymentPatchConfigPostStepWindowsExecStepConfig) Empty() bool {
	return r.empty
}

func (r *PatchDeploymentPatchConfigPostStepWindowsExecStepConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *PatchDeploymentPatchConfigPostStepWindowsExecStepConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject struct {
	empty            bool    `json:"-"`
	Bucket           *string `json:"bucket"`
	Object           *string `json:"object"`
	GenerationNumber *int64  `json:"generationNumber"`
}

type jsonPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject

func (r *PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject) UnmarshalJSON(data []byte) error {
	var res jsonPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject
	} else {

		r.Bucket = res.Bucket

		r.Object = res.Object

		r.GenerationNumber = res.GenerationNumber

	}
	return nil
}

// This object is used to assert a desired state where this PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject *PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject = &PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject{empty: true}

func (r *PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject) Empty() bool {
	return r.empty
}

func (r *PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject) String() string {
	return dcl.SprintResource(r)
}

func (r *PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PatchDeploymentOneTimeSchedule struct {
	empty       bool    `json:"-"`
	ExecuteTime *string `json:"executeTime"`
}

type jsonPatchDeploymentOneTimeSchedule PatchDeploymentOneTimeSchedule

func (r *PatchDeploymentOneTimeSchedule) UnmarshalJSON(data []byte) error {
	var res jsonPatchDeploymentOneTimeSchedule
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPatchDeploymentOneTimeSchedule
	} else {

		r.ExecuteTime = res.ExecuteTime

	}
	return nil
}

// This object is used to assert a desired state where this PatchDeploymentOneTimeSchedule is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyPatchDeploymentOneTimeSchedule *PatchDeploymentOneTimeSchedule = &PatchDeploymentOneTimeSchedule{empty: true}

func (r *PatchDeploymentOneTimeSchedule) Empty() bool {
	return r.empty
}

func (r *PatchDeploymentOneTimeSchedule) String() string {
	return dcl.SprintResource(r)
}

func (r *PatchDeploymentOneTimeSchedule) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PatchDeploymentRecurringSchedule struct {
	empty           bool                                           `json:"-"`
	TimeZone        *PatchDeploymentRecurringScheduleTimeZone      `json:"timeZone"`
	StartTime       *string                                        `json:"startTime"`
	EndTime         *string                                        `json:"endTime"`
	TimeOfDay       *PatchDeploymentRecurringScheduleTimeOfDay     `json:"timeOfDay"`
	Frequency       *PatchDeploymentRecurringScheduleFrequencyEnum `json:"frequency"`
	Weekly          *PatchDeploymentRecurringScheduleWeekly        `json:"weekly"`
	Monthly         *PatchDeploymentRecurringScheduleMonthly       `json:"monthly"`
	LastExecuteTime *string                                        `json:"lastExecuteTime"`
	NextExecuteTime *string                                        `json:"nextExecuteTime"`
}

type jsonPatchDeploymentRecurringSchedule PatchDeploymentRecurringSchedule

func (r *PatchDeploymentRecurringSchedule) UnmarshalJSON(data []byte) error {
	var res jsonPatchDeploymentRecurringSchedule
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPatchDeploymentRecurringSchedule
	} else {

		r.TimeZone = res.TimeZone

		r.StartTime = res.StartTime

		r.EndTime = res.EndTime

		r.TimeOfDay = res.TimeOfDay

		r.Frequency = res.Frequency

		r.Weekly = res.Weekly

		r.Monthly = res.Monthly

		r.LastExecuteTime = res.LastExecuteTime

		r.NextExecuteTime = res.NextExecuteTime

	}
	return nil
}

// This object is used to assert a desired state where this PatchDeploymentRecurringSchedule is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyPatchDeploymentRecurringSchedule *PatchDeploymentRecurringSchedule = &PatchDeploymentRecurringSchedule{empty: true}

func (r *PatchDeploymentRecurringSchedule) Empty() bool {
	return r.empty
}

func (r *PatchDeploymentRecurringSchedule) String() string {
	return dcl.SprintResource(r)
}

func (r *PatchDeploymentRecurringSchedule) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PatchDeploymentRecurringScheduleTimeZone struct {
	empty   bool    `json:"-"`
	Id      *string `json:"id"`
	Version *string `json:"version"`
}

type jsonPatchDeploymentRecurringScheduleTimeZone PatchDeploymentRecurringScheduleTimeZone

func (r *PatchDeploymentRecurringScheduleTimeZone) UnmarshalJSON(data []byte) error {
	var res jsonPatchDeploymentRecurringScheduleTimeZone
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPatchDeploymentRecurringScheduleTimeZone
	} else {

		r.Id = res.Id

		r.Version = res.Version

	}
	return nil
}

// This object is used to assert a desired state where this PatchDeploymentRecurringScheduleTimeZone is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyPatchDeploymentRecurringScheduleTimeZone *PatchDeploymentRecurringScheduleTimeZone = &PatchDeploymentRecurringScheduleTimeZone{empty: true}

func (r *PatchDeploymentRecurringScheduleTimeZone) Empty() bool {
	return r.empty
}

func (r *PatchDeploymentRecurringScheduleTimeZone) String() string {
	return dcl.SprintResource(r)
}

func (r *PatchDeploymentRecurringScheduleTimeZone) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PatchDeploymentRecurringScheduleTimeOfDay struct {
	empty   bool   `json:"-"`
	Hours   *int64 `json:"hours"`
	Minutes *int64 `json:"minutes"`
	Seconds *int64 `json:"seconds"`
	Nanos   *int64 `json:"nanos"`
}

type jsonPatchDeploymentRecurringScheduleTimeOfDay PatchDeploymentRecurringScheduleTimeOfDay

func (r *PatchDeploymentRecurringScheduleTimeOfDay) UnmarshalJSON(data []byte) error {
	var res jsonPatchDeploymentRecurringScheduleTimeOfDay
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPatchDeploymentRecurringScheduleTimeOfDay
	} else {

		r.Hours = res.Hours

		r.Minutes = res.Minutes

		r.Seconds = res.Seconds

		r.Nanos = res.Nanos

	}
	return nil
}

// This object is used to assert a desired state where this PatchDeploymentRecurringScheduleTimeOfDay is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyPatchDeploymentRecurringScheduleTimeOfDay *PatchDeploymentRecurringScheduleTimeOfDay = &PatchDeploymentRecurringScheduleTimeOfDay{empty: true}

func (r *PatchDeploymentRecurringScheduleTimeOfDay) Empty() bool {
	return r.empty
}

func (r *PatchDeploymentRecurringScheduleTimeOfDay) String() string {
	return dcl.SprintResource(r)
}

func (r *PatchDeploymentRecurringScheduleTimeOfDay) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PatchDeploymentRecurringScheduleWeekly struct {
	empty     bool                                                 `json:"-"`
	DayOfWeek *PatchDeploymentRecurringScheduleWeeklyDayOfWeekEnum `json:"dayOfWeek"`
}

type jsonPatchDeploymentRecurringScheduleWeekly PatchDeploymentRecurringScheduleWeekly

func (r *PatchDeploymentRecurringScheduleWeekly) UnmarshalJSON(data []byte) error {
	var res jsonPatchDeploymentRecurringScheduleWeekly
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPatchDeploymentRecurringScheduleWeekly
	} else {

		r.DayOfWeek = res.DayOfWeek

	}
	return nil
}

// This object is used to assert a desired state where this PatchDeploymentRecurringScheduleWeekly is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyPatchDeploymentRecurringScheduleWeekly *PatchDeploymentRecurringScheduleWeekly = &PatchDeploymentRecurringScheduleWeekly{empty: true}

func (r *PatchDeploymentRecurringScheduleWeekly) Empty() bool {
	return r.empty
}

func (r *PatchDeploymentRecurringScheduleWeekly) String() string {
	return dcl.SprintResource(r)
}

func (r *PatchDeploymentRecurringScheduleWeekly) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PatchDeploymentRecurringScheduleMonthly struct {
	empty          bool                                                   `json:"-"`
	WeekDayOfMonth *PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth `json:"weekDayOfMonth"`
	MonthDay       *int64                                                 `json:"monthDay"`
}

type jsonPatchDeploymentRecurringScheduleMonthly PatchDeploymentRecurringScheduleMonthly

func (r *PatchDeploymentRecurringScheduleMonthly) UnmarshalJSON(data []byte) error {
	var res jsonPatchDeploymentRecurringScheduleMonthly
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPatchDeploymentRecurringScheduleMonthly
	} else {

		r.WeekDayOfMonth = res.WeekDayOfMonth

		r.MonthDay = res.MonthDay

	}
	return nil
}

// This object is used to assert a desired state where this PatchDeploymentRecurringScheduleMonthly is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyPatchDeploymentRecurringScheduleMonthly *PatchDeploymentRecurringScheduleMonthly = &PatchDeploymentRecurringScheduleMonthly{empty: true}

func (r *PatchDeploymentRecurringScheduleMonthly) Empty() bool {
	return r.empty
}

func (r *PatchDeploymentRecurringScheduleMonthly) String() string {
	return dcl.SprintResource(r)
}

func (r *PatchDeploymentRecurringScheduleMonthly) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth struct {
	empty       bool                                                                `json:"-"`
	WeekOrdinal *int64                                                              `json:"weekOrdinal"`
	DayOfWeek   *PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnum `json:"dayOfWeek"`
}

type jsonPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth

func (r *PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth) UnmarshalJSON(data []byte) error {
	var res jsonPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth
	} else {

		r.WeekOrdinal = res.WeekOrdinal

		r.DayOfWeek = res.DayOfWeek

	}
	return nil
}

// This object is used to assert a desired state where this PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth *PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth = &PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth{empty: true}

func (r *PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth) Empty() bool {
	return r.empty
}

func (r *PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth) String() string {
	return dcl.SprintResource(r)
}

func (r *PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PatchDeploymentRollout struct {
	empty            bool                                    `json:"-"`
	Mode             *PatchDeploymentRolloutModeEnum         `json:"mode"`
	DisruptionBudget *PatchDeploymentRolloutDisruptionBudget `json:"disruptionBudget"`
}

type jsonPatchDeploymentRollout PatchDeploymentRollout

func (r *PatchDeploymentRollout) UnmarshalJSON(data []byte) error {
	var res jsonPatchDeploymentRollout
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPatchDeploymentRollout
	} else {

		r.Mode = res.Mode

		r.DisruptionBudget = res.DisruptionBudget

	}
	return nil
}

// This object is used to assert a desired state where this PatchDeploymentRollout is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyPatchDeploymentRollout *PatchDeploymentRollout = &PatchDeploymentRollout{empty: true}

func (r *PatchDeploymentRollout) Empty() bool {
	return r.empty
}

func (r *PatchDeploymentRollout) String() string {
	return dcl.SprintResource(r)
}

func (r *PatchDeploymentRollout) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PatchDeploymentRolloutDisruptionBudget struct {
	empty   bool   `json:"-"`
	Fixed   *int64 `json:"fixed"`
	Percent *int64 `json:"percent"`
}

type jsonPatchDeploymentRolloutDisruptionBudget PatchDeploymentRolloutDisruptionBudget

func (r *PatchDeploymentRolloutDisruptionBudget) UnmarshalJSON(data []byte) error {
	var res jsonPatchDeploymentRolloutDisruptionBudget
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPatchDeploymentRolloutDisruptionBudget
	} else {

		r.Fixed = res.Fixed

		r.Percent = res.Percent

	}
	return nil
}

// This object is used to assert a desired state where this PatchDeploymentRolloutDisruptionBudget is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyPatchDeploymentRolloutDisruptionBudget *PatchDeploymentRolloutDisruptionBudget = &PatchDeploymentRolloutDisruptionBudget{empty: true}

func (r *PatchDeploymentRolloutDisruptionBudget) Empty() bool {
	return r.empty
}

func (r *PatchDeploymentRolloutDisruptionBudget) String() string {
	return dcl.SprintResource(r)
}

func (r *PatchDeploymentRolloutDisruptionBudget) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *PatchDeployment) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "os_config",
		Type:    "PatchDeployment",
		Version: "osconfig",
	}
}

func (r *PatchDeployment) ID() (string, error) {
	if err := extractPatchDeploymentFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name":              dcl.ValueOrEmptyString(nr.Name),
		"description":       dcl.ValueOrEmptyString(nr.Description),
		"instanceFilter":    dcl.ValueOrEmptyString(nr.InstanceFilter),
		"patchConfig":       dcl.ValueOrEmptyString(nr.PatchConfig),
		"duration":          dcl.ValueOrEmptyString(nr.Duration),
		"oneTimeSchedule":   dcl.ValueOrEmptyString(nr.OneTimeSchedule),
		"recurringSchedule": dcl.ValueOrEmptyString(nr.RecurringSchedule),
		"createTime":        dcl.ValueOrEmptyString(nr.CreateTime),
		"updateTime":        dcl.ValueOrEmptyString(nr.UpdateTime),
		"lastExecuteTime":   dcl.ValueOrEmptyString(nr.LastExecuteTime),
		"rollout":           dcl.ValueOrEmptyString(nr.Rollout),
		"project":           dcl.ValueOrEmptyString(nr.Project),
	}
	return dcl.Nprintf("projects/{{project}}/patchDeployments/{{name}}", params), nil
}

const PatchDeploymentMaxPage = -1

type PatchDeploymentList struct {
	Items []*PatchDeployment

	nextToken string

	pageSize int32

	resource *PatchDeployment
}

func (l *PatchDeploymentList) HasNext() bool {
	return l.nextToken != ""
}

func (l *PatchDeploymentList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listPatchDeployment(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListPatchDeployment(ctx context.Context, project string) (*PatchDeploymentList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListPatchDeploymentWithMaxResults(ctx, project, PatchDeploymentMaxPage)

}

func (c *Client) ListPatchDeploymentWithMaxResults(ctx context.Context, project string, pageSize int32) (*PatchDeploymentList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &PatchDeployment{
		Project: &project,
	}
	items, token, err := c.listPatchDeployment(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &PatchDeploymentList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetPatchDeployment(ctx context.Context, r *PatchDeployment) (*PatchDeployment, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractPatchDeploymentFields(r)

	b, err := c.getPatchDeploymentRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalPatchDeployment(b, c)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizePatchDeploymentNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractPatchDeploymentFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeletePatchDeployment(ctx context.Context, r *PatchDeployment) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("PatchDeployment resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting PatchDeployment...")
	deleteOp := deletePatchDeploymentOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllPatchDeployment deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllPatchDeployment(ctx context.Context, project string, filter func(*PatchDeployment) bool) error {
	listObj, err := c.ListPatchDeployment(ctx, project)
	if err != nil {
		return err
	}

	err = c.deleteAllPatchDeployment(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllPatchDeployment(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyPatchDeployment(ctx context.Context, rawDesired *PatchDeployment, opts ...dcl.ApplyOption) (*PatchDeployment, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *PatchDeployment
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyPatchDeploymentHelper(c, ctx, rawDesired, opts...)
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

func applyPatchDeploymentHelper(c *Client, ctx context.Context, rawDesired *PatchDeployment, opts ...dcl.ApplyOption) (*PatchDeployment, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyPatchDeployment...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractPatchDeploymentFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.patchDeploymentDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToPatchDeploymentDiffs(c.Config, fieldDiffs, opts)
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
	var ops []patchDeploymentApiOperation
	if create {
		ops = append(ops, &createPatchDeploymentOperation{})
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
	return applyPatchDeploymentDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyPatchDeploymentDiff(c *Client, ctx context.Context, desired *PatchDeployment, rawDesired *PatchDeployment, ops []patchDeploymentApiOperation, opts ...dcl.ApplyOption) (*PatchDeployment, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetPatchDeployment(ctx, desired.urlNormalized())
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createPatchDeploymentOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapPatchDeployment(r, c)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizePatchDeploymentNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizePatchDeploymentNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizePatchDeploymentDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractPatchDeploymentFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractPatchDeploymentFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffPatchDeployment(c, newDesired, newState)
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
