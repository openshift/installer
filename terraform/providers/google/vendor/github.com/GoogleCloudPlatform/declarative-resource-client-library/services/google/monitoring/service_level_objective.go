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
package monitoring

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/api/googleapi"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

type ServiceLevelObjective struct {
	Name                   *string                                     `json:"name"`
	DisplayName            *string                                     `json:"displayName"`
	ServiceLevelIndicator  *ServiceLevelObjectiveServiceLevelIndicator `json:"serviceLevelIndicator"`
	Goal                   *float64                                    `json:"goal"`
	RollingPeriod          *string                                     `json:"rollingPeriod"`
	CalendarPeriod         *ServiceLevelObjectiveCalendarPeriodEnum    `json:"calendarPeriod"`
	CreateTime             *string                                     `json:"createTime"`
	DeleteTime             *string                                     `json:"deleteTime"`
	ServiceManagementOwned *bool                                       `json:"serviceManagementOwned"`
	UserLabels             map[string]string                           `json:"userLabels"`
	Project                *string                                     `json:"project"`
	Service                *string                                     `json:"service"`
}

func (r *ServiceLevelObjective) String() string {
	return dcl.SprintResource(r)
}

// The enum ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnum.
type ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnum string

// ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnumRef returns a *ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnum with the value of string s
// If the empty string is provided, nil is returned.
func ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnumRef(s string) *ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnum {
	v := ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnum(s)
	return &v
}

func (v ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"LATENCY_EXPERIENCE_UNSPECIFIED", "DELIGHTING", "SATISFYING", "ANNOYING"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnum.
type ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnum string

// ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnumRef returns a *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnum with the value of string s
// If the empty string is provided, nil is returned.
func ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnumRef(s string) *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnum {
	v := ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnum(s)
	return &v
}

func (v ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"LATENCY_EXPERIENCE_UNSPECIFIED", "DELIGHTING", "SATISFYING", "ANNOYING"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnum.
type ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnum string

// ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnumRef returns a *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnum with the value of string s
// If the empty string is provided, nil is returned.
func ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnumRef(s string) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnum {
	v := ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnum(s)
	return &v
}

func (v ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"LATENCY_EXPERIENCE_UNSPECIFIED", "DELIGHTING", "SATISFYING", "ANNOYING"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnum.
type ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnum string

// ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnumRef returns a *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnum with the value of string s
// If the empty string is provided, nil is returned.
func ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnumRef(s string) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnum {
	v := ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnum(s)
	return &v
}

func (v ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"LATENCY_EXPERIENCE_UNSPECIFIED", "DELIGHTING", "SATISFYING", "ANNOYING"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum ServiceLevelObjectiveCalendarPeriodEnum.
type ServiceLevelObjectiveCalendarPeriodEnum string

// ServiceLevelObjectiveCalendarPeriodEnumRef returns a *ServiceLevelObjectiveCalendarPeriodEnum with the value of string s
// If the empty string is provided, nil is returned.
func ServiceLevelObjectiveCalendarPeriodEnumRef(s string) *ServiceLevelObjectiveCalendarPeriodEnum {
	v := ServiceLevelObjectiveCalendarPeriodEnum(s)
	return &v
}

func (v ServiceLevelObjectiveCalendarPeriodEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"CALENDAR_PERIOD_UNSPECIFIED", "DAY", "WEEK", "FORTNIGHT", "MONTH", "QUARTER", "HALF", "YEAR"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "ServiceLevelObjectiveCalendarPeriodEnum",
		Value: string(v),
		Valid: []string{},
	}
}

type ServiceLevelObjectiveServiceLevelIndicator struct {
	empty        bool                                                    `json:"-"`
	BasicSli     *ServiceLevelObjectiveServiceLevelIndicatorBasicSli     `json:"basicSli"`
	RequestBased *ServiceLevelObjectiveServiceLevelIndicatorRequestBased `json:"requestBased"`
	WindowsBased *ServiceLevelObjectiveServiceLevelIndicatorWindowsBased `json:"windowsBased"`
}

type jsonServiceLevelObjectiveServiceLevelIndicator ServiceLevelObjectiveServiceLevelIndicator

func (r *ServiceLevelObjectiveServiceLevelIndicator) UnmarshalJSON(data []byte) error {
	var res jsonServiceLevelObjectiveServiceLevelIndicator
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyServiceLevelObjectiveServiceLevelIndicator
	} else {

		r.BasicSli = res.BasicSli

		r.RequestBased = res.RequestBased

		r.WindowsBased = res.WindowsBased

	}
	return nil
}

// This object is used to assert a desired state where this ServiceLevelObjectiveServiceLevelIndicator is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyServiceLevelObjectiveServiceLevelIndicator *ServiceLevelObjectiveServiceLevelIndicator = &ServiceLevelObjectiveServiceLevelIndicator{empty: true}

func (r *ServiceLevelObjectiveServiceLevelIndicator) Empty() bool {
	return r.empty
}

func (r *ServiceLevelObjectiveServiceLevelIndicator) String() string {
	return dcl.SprintResource(r)
}

func (r *ServiceLevelObjectiveServiceLevelIndicator) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ServiceLevelObjectiveServiceLevelIndicatorBasicSli struct {
	empty                 bool                                                                     `json:"-"`
	Method                []string                                                                 `json:"method"`
	Location              []string                                                                 `json:"location"`
	Version               []string                                                                 `json:"version"`
	Availability          *ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability          `json:"availability"`
	Latency               *ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency               `json:"latency"`
	OperationAvailability *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability `json:"operationAvailability"`
	OperationLatency      *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency      `json:"operationLatency"`
}

type jsonServiceLevelObjectiveServiceLevelIndicatorBasicSli ServiceLevelObjectiveServiceLevelIndicatorBasicSli

func (r *ServiceLevelObjectiveServiceLevelIndicatorBasicSli) UnmarshalJSON(data []byte) error {
	var res jsonServiceLevelObjectiveServiceLevelIndicatorBasicSli
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyServiceLevelObjectiveServiceLevelIndicatorBasicSli
	} else {

		r.Method = res.Method

		r.Location = res.Location

		r.Version = res.Version

		r.Availability = res.Availability

		r.Latency = res.Latency

		r.OperationAvailability = res.OperationAvailability

		r.OperationLatency = res.OperationLatency

	}
	return nil
}

// This object is used to assert a desired state where this ServiceLevelObjectiveServiceLevelIndicatorBasicSli is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyServiceLevelObjectiveServiceLevelIndicatorBasicSli *ServiceLevelObjectiveServiceLevelIndicatorBasicSli = &ServiceLevelObjectiveServiceLevelIndicatorBasicSli{empty: true}

func (r *ServiceLevelObjectiveServiceLevelIndicatorBasicSli) Empty() bool {
	return r.empty
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorBasicSli) String() string {
	return dcl.SprintResource(r)
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorBasicSli) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability struct {
	empty bool `json:"-"`
}

type jsonServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability

func (r *ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability) UnmarshalJSON(data []byte) error {
	var res jsonServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability
	} else {

	}
	return nil
}

// This object is used to assert a desired state where this ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability *ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability = &ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability{empty: true}

func (r *ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability) Empty() bool {
	return r.empty
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability) String() string {
	return dcl.SprintResource(r)
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency struct {
	empty      bool                                                                     `json:"-"`
	Threshold  *string                                                                  `json:"threshold"`
	Experience *ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnum `json:"experience"`
}

type jsonServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency

func (r *ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency) UnmarshalJSON(data []byte) error {
	var res jsonServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency
	} else {

		r.Threshold = res.Threshold

		r.Experience = res.Experience

	}
	return nil
}

// This object is used to assert a desired state where this ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency *ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency = &ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency{empty: true}

func (r *ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency) Empty() bool {
	return r.empty
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency) String() string {
	return dcl.SprintResource(r)
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability struct {
	empty bool `json:"-"`
}

type jsonServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability

func (r *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability) UnmarshalJSON(data []byte) error {
	var res jsonServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability
	} else {

	}
	return nil
}

// This object is used to assert a desired state where this ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability = &ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability{empty: true}

func (r *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability) Empty() bool {
	return r.empty
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability) String() string {
	return dcl.SprintResource(r)
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency struct {
	empty      bool                                                                              `json:"-"`
	Threshold  *string                                                                           `json:"threshold"`
	Experience *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnum `json:"experience"`
}

type jsonServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency

func (r *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency) UnmarshalJSON(data []byte) error {
	var res jsonServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency
	} else {

		r.Threshold = res.Threshold

		r.Experience = res.Experience

	}
	return nil
}

// This object is used to assert a desired state where this ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency = &ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency{empty: true}

func (r *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency) Empty() bool {
	return r.empty
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency) String() string {
	return dcl.SprintResource(r)
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ServiceLevelObjectiveServiceLevelIndicatorRequestBased struct {
	empty           bool                                                                   `json:"-"`
	GoodTotalRatio  *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio  `json:"goodTotalRatio"`
	DistributionCut *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut `json:"distributionCut"`
}

type jsonServiceLevelObjectiveServiceLevelIndicatorRequestBased ServiceLevelObjectiveServiceLevelIndicatorRequestBased

func (r *ServiceLevelObjectiveServiceLevelIndicatorRequestBased) UnmarshalJSON(data []byte) error {
	var res jsonServiceLevelObjectiveServiceLevelIndicatorRequestBased
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyServiceLevelObjectiveServiceLevelIndicatorRequestBased
	} else {

		r.GoodTotalRatio = res.GoodTotalRatio

		r.DistributionCut = res.DistributionCut

	}
	return nil
}

// This object is used to assert a desired state where this ServiceLevelObjectiveServiceLevelIndicatorRequestBased is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyServiceLevelObjectiveServiceLevelIndicatorRequestBased *ServiceLevelObjectiveServiceLevelIndicatorRequestBased = &ServiceLevelObjectiveServiceLevelIndicatorRequestBased{empty: true}

func (r *ServiceLevelObjectiveServiceLevelIndicatorRequestBased) Empty() bool {
	return r.empty
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorRequestBased) String() string {
	return dcl.SprintResource(r)
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorRequestBased) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio struct {
	empty              bool    `json:"-"`
	GoodServiceFilter  *string `json:"goodServiceFilter"`
	BadServiceFilter   *string `json:"badServiceFilter"`
	TotalServiceFilter *string `json:"totalServiceFilter"`
}

type jsonServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio

func (r *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio) UnmarshalJSON(data []byte) error {
	var res jsonServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio
	} else {

		r.GoodServiceFilter = res.GoodServiceFilter

		r.BadServiceFilter = res.BadServiceFilter

		r.TotalServiceFilter = res.TotalServiceFilter

	}
	return nil
}

// This object is used to assert a desired state where this ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio = &ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio{empty: true}

func (r *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio) Empty() bool {
	return r.empty
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio) String() string {
	return dcl.SprintResource(r)
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut struct {
	empty              bool                                                                        `json:"-"`
	DistributionFilter *string                                                                     `json:"distributionFilter"`
	Range              *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange `json:"range"`
}

type jsonServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut

func (r *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut) UnmarshalJSON(data []byte) error {
	var res jsonServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut
	} else {

		r.DistributionFilter = res.DistributionFilter

		r.Range = res.Range

	}
	return nil
}

// This object is used to assert a desired state where this ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut = &ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut{empty: true}

func (r *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut) Empty() bool {
	return r.empty
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut) String() string {
	return dcl.SprintResource(r)
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange struct {
	empty bool     `json:"-"`
	Min   *float64 `json:"min"`
	Max   *float64 `json:"max"`
}

type jsonServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange

func (r *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange) UnmarshalJSON(data []byte) error {
	var res jsonServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange
	} else {

		r.Min = res.Min

		r.Max = res.Max

	}
	return nil
}

// This object is used to assert a desired state where this ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange = &ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange{empty: true}

func (r *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange) Empty() bool {
	return r.empty
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange) String() string {
	return dcl.SprintResource(r)
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ServiceLevelObjectiveServiceLevelIndicatorWindowsBased struct {
	empty                   bool                                                                           `json:"-"`
	GoodBadMetricFilter     *string                                                                        `json:"goodBadMetricFilter"`
	GoodTotalRatioThreshold *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold `json:"goodTotalRatioThreshold"`
	MetricMeanInRange       *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange       `json:"metricMeanInRange"`
	MetricSumInRange        *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange        `json:"metricSumInRange"`
	WindowPeriod            *string                                                                        `json:"windowPeriod"`
}

type jsonServiceLevelObjectiveServiceLevelIndicatorWindowsBased ServiceLevelObjectiveServiceLevelIndicatorWindowsBased

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBased) UnmarshalJSON(data []byte) error {
	var res jsonServiceLevelObjectiveServiceLevelIndicatorWindowsBased
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBased
	} else {

		r.GoodBadMetricFilter = res.GoodBadMetricFilter

		r.GoodTotalRatioThreshold = res.GoodTotalRatioThreshold

		r.MetricMeanInRange = res.MetricMeanInRange

		r.MetricSumInRange = res.MetricSumInRange

		r.WindowPeriod = res.WindowPeriod

	}
	return nil
}

// This object is used to assert a desired state where this ServiceLevelObjectiveServiceLevelIndicatorWindowsBased is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBased *ServiceLevelObjectiveServiceLevelIndicatorWindowsBased = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBased{empty: true}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBased) Empty() bool {
	return r.empty
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBased) String() string {
	return dcl.SprintResource(r)
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBased) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold struct {
	empty               bool                                                                                              `json:"-"`
	Performance         *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance         `json:"performance"`
	BasicSliPerformance *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance `json:"basicSliPerformance"`
	Threshold           *float64                                                                                          `json:"threshold"`
}

type jsonServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold) UnmarshalJSON(data []byte) error {
	var res jsonServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold
	} else {

		r.Performance = res.Performance

		r.BasicSliPerformance = res.BasicSliPerformance

		r.Threshold = res.Threshold

	}
	return nil
}

// This object is used to assert a desired state where this ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold{empty: true}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold) Empty() bool {
	return r.empty
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold) String() string {
	return dcl.SprintResource(r)
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance struct {
	empty           bool                                                                                                     `json:"-"`
	GoodTotalRatio  *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio  `json:"goodTotalRatio"`
	DistributionCut *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut `json:"distributionCut"`
}

type jsonServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance) UnmarshalJSON(data []byte) error {
	var res jsonServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance
	} else {

		r.GoodTotalRatio = res.GoodTotalRatio

		r.DistributionCut = res.DistributionCut

	}
	return nil
}

// This object is used to assert a desired state where this ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance{empty: true}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance) Empty() bool {
	return r.empty
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance) String() string {
	return dcl.SprintResource(r)
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio struct {
	empty              bool    `json:"-"`
	GoodServiceFilter  *string `json:"goodServiceFilter"`
	BadServiceFilter   *string `json:"badServiceFilter"`
	TotalServiceFilter *string `json:"totalServiceFilter"`
}

type jsonServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio) UnmarshalJSON(data []byte) error {
	var res jsonServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio
	} else {

		r.GoodServiceFilter = res.GoodServiceFilter

		r.BadServiceFilter = res.BadServiceFilter

		r.TotalServiceFilter = res.TotalServiceFilter

	}
	return nil
}

// This object is used to assert a desired state where this ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio{empty: true}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio) Empty() bool {
	return r.empty
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio) String() string {
	return dcl.SprintResource(r)
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut struct {
	empty              bool                                                                                                          `json:"-"`
	DistributionFilter *string                                                                                                       `json:"distributionFilter"`
	Range              *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange `json:"range"`
}

type jsonServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut) UnmarshalJSON(data []byte) error {
	var res jsonServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut
	} else {

		r.DistributionFilter = res.DistributionFilter

		r.Range = res.Range

	}
	return nil
}

// This object is used to assert a desired state where this ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut{empty: true}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut) Empty() bool {
	return r.empty
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut) String() string {
	return dcl.SprintResource(r)
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange struct {
	empty bool     `json:"-"`
	Min   *float64 `json:"min"`
	Max   *float64 `json:"max"`
}

type jsonServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange) UnmarshalJSON(data []byte) error {
	var res jsonServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange
	} else {

		r.Min = res.Min

		r.Max = res.Max

	}
	return nil
}

// This object is used to assert a desired state where this ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange{empty: true}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange) Empty() bool {
	return r.empty
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange) String() string {
	return dcl.SprintResource(r)
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance struct {
	empty                 bool                                                                                                                   `json:"-"`
	Method                []string                                                                                                               `json:"method"`
	Location              []string                                                                                                               `json:"location"`
	Version               []string                                                                                                               `json:"version"`
	Availability          *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability          `json:"availability"`
	Latency               *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency               `json:"latency"`
	OperationAvailability *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability `json:"operationAvailability"`
	OperationLatency      *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency      `json:"operationLatency"`
}

type jsonServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance) UnmarshalJSON(data []byte) error {
	var res jsonServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance
	} else {

		r.Method = res.Method

		r.Location = res.Location

		r.Version = res.Version

		r.Availability = res.Availability

		r.Latency = res.Latency

		r.OperationAvailability = res.OperationAvailability

		r.OperationLatency = res.OperationLatency

	}
	return nil
}

// This object is used to assert a desired state where this ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance{empty: true}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance) Empty() bool {
	return r.empty
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance) String() string {
	return dcl.SprintResource(r)
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability struct {
	empty bool `json:"-"`
}

type jsonServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability) UnmarshalJSON(data []byte) error {
	var res jsonServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability
	} else {

	}
	return nil
}

// This object is used to assert a desired state where this ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability{empty: true}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability) Empty() bool {
	return r.empty
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability) String() string {
	return dcl.SprintResource(r)
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency struct {
	empty      bool                                                                                                                   `json:"-"`
	Threshold  *string                                                                                                                `json:"threshold"`
	Experience *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnum `json:"experience"`
}

type jsonServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency) UnmarshalJSON(data []byte) error {
	var res jsonServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency
	} else {

		r.Threshold = res.Threshold

		r.Experience = res.Experience

	}
	return nil
}

// This object is used to assert a desired state where this ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency{empty: true}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency) Empty() bool {
	return r.empty
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency) String() string {
	return dcl.SprintResource(r)
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability struct {
	empty bool `json:"-"`
}

type jsonServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability) UnmarshalJSON(data []byte) error {
	var res jsonServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability
	} else {

	}
	return nil
}

// This object is used to assert a desired state where this ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability{empty: true}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability) Empty() bool {
	return r.empty
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability) String() string {
	return dcl.SprintResource(r)
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency struct {
	empty      bool                                                                                                                            `json:"-"`
	Threshold  *string                                                                                                                         `json:"threshold"`
	Experience *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnum `json:"experience"`
}

type jsonServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency) UnmarshalJSON(data []byte) error {
	var res jsonServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency
	} else {

		r.Threshold = res.Threshold

		r.Experience = res.Experience

	}
	return nil
}

// This object is used to assert a desired state where this ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency{empty: true}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency) Empty() bool {
	return r.empty
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency) String() string {
	return dcl.SprintResource(r)
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange struct {
	empty      bool                                                                          `json:"-"`
	TimeSeries *string                                                                       `json:"timeSeries"`
	Range      *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange `json:"range"`
}

type jsonServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange) UnmarshalJSON(data []byte) error {
	var res jsonServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange
	} else {

		r.TimeSeries = res.TimeSeries

		r.Range = res.Range

	}
	return nil
}

// This object is used to assert a desired state where this ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange{empty: true}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange) Empty() bool {
	return r.empty
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange) String() string {
	return dcl.SprintResource(r)
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange struct {
	empty bool     `json:"-"`
	Min   *float64 `json:"min"`
	Max   *float64 `json:"max"`
}

type jsonServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange) UnmarshalJSON(data []byte) error {
	var res jsonServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange
	} else {

		r.Min = res.Min

		r.Max = res.Max

	}
	return nil
}

// This object is used to assert a desired state where this ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange{empty: true}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange) Empty() bool {
	return r.empty
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange) String() string {
	return dcl.SprintResource(r)
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange struct {
	empty      bool                                                                         `json:"-"`
	TimeSeries *string                                                                      `json:"timeSeries"`
	Range      *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange `json:"range"`
}

type jsonServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange) UnmarshalJSON(data []byte) error {
	var res jsonServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange
	} else {

		r.TimeSeries = res.TimeSeries

		r.Range = res.Range

	}
	return nil
}

// This object is used to assert a desired state where this ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange{empty: true}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange) Empty() bool {
	return r.empty
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange) String() string {
	return dcl.SprintResource(r)
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange struct {
	empty bool     `json:"-"`
	Min   *float64 `json:"min"`
	Max   *float64 `json:"max"`
}

type jsonServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange) UnmarshalJSON(data []byte) error {
	var res jsonServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange
	} else {

		r.Min = res.Min

		r.Max = res.Max

	}
	return nil
}

// This object is used to assert a desired state where this ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange{empty: true}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange) Empty() bool {
	return r.empty
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange) String() string {
	return dcl.SprintResource(r)
}

func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *ServiceLevelObjective) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "monitoring",
		Type:    "ServiceLevelObjective",
		Version: "monitoring",
	}
}

func (r *ServiceLevelObjective) ID() (string, error) {
	if err := extractServiceLevelObjectiveFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name":                     dcl.ValueOrEmptyString(nr.Name),
		"display_name":             dcl.ValueOrEmptyString(nr.DisplayName),
		"service_level_indicator":  dcl.ValueOrEmptyString(nr.ServiceLevelIndicator),
		"goal":                     dcl.ValueOrEmptyString(nr.Goal),
		"rolling_period":           dcl.ValueOrEmptyString(nr.RollingPeriod),
		"calendar_period":          dcl.ValueOrEmptyString(nr.CalendarPeriod),
		"create_time":              dcl.ValueOrEmptyString(nr.CreateTime),
		"delete_time":              dcl.ValueOrEmptyString(nr.DeleteTime),
		"service_management_owned": dcl.ValueOrEmptyString(nr.ServiceManagementOwned),
		"user_labels":              dcl.ValueOrEmptyString(nr.UserLabels),
		"project":                  dcl.ValueOrEmptyString(nr.Project),
		"service":                  dcl.ValueOrEmptyString(nr.Service),
	}
	return dcl.Nprintf("projects/{{project}}/services/{{service}}/serviceLevelObjectives/{{name}}", params), nil
}

const ServiceLevelObjectiveMaxPage = -1

type ServiceLevelObjectiveList struct {
	Items []*ServiceLevelObjective

	nextToken string

	pageSize int32

	resource *ServiceLevelObjective
}

func (l *ServiceLevelObjectiveList) HasNext() bool {
	return l.nextToken != ""
}

func (l *ServiceLevelObjectiveList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listServiceLevelObjective(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListServiceLevelObjective(ctx context.Context, project, service string) (*ServiceLevelObjectiveList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListServiceLevelObjectiveWithMaxResults(ctx, project, service, ServiceLevelObjectiveMaxPage)

}

func (c *Client) ListServiceLevelObjectiveWithMaxResults(ctx context.Context, project, service string, pageSize int32) (*ServiceLevelObjectiveList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &ServiceLevelObjective{
		Project: &project,
		Service: &service,
	}
	items, token, err := c.listServiceLevelObjective(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &ServiceLevelObjectiveList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetServiceLevelObjective(ctx context.Context, r *ServiceLevelObjective) (*ServiceLevelObjective, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractServiceLevelObjectiveFields(r)

	b, err := c.getServiceLevelObjectiveRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalServiceLevelObjective(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Service = r.Service
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeServiceLevelObjectiveNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractServiceLevelObjectiveFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteServiceLevelObjective(ctx context.Context, r *ServiceLevelObjective) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("ServiceLevelObjective resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting ServiceLevelObjective...")
	deleteOp := deleteServiceLevelObjectiveOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllServiceLevelObjective deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllServiceLevelObjective(ctx context.Context, project, service string, filter func(*ServiceLevelObjective) bool) error {
	listObj, err := c.ListServiceLevelObjective(ctx, project, service)
	if err != nil {
		return err
	}

	err = c.deleteAllServiceLevelObjective(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllServiceLevelObjective(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyServiceLevelObjective(ctx context.Context, rawDesired *ServiceLevelObjective, opts ...dcl.ApplyOption) (*ServiceLevelObjective, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *ServiceLevelObjective
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyServiceLevelObjectiveHelper(c, ctx, rawDesired, opts...)
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

func applyServiceLevelObjectiveHelper(c *Client, ctx context.Context, rawDesired *ServiceLevelObjective, opts ...dcl.ApplyOption) (*ServiceLevelObjective, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyServiceLevelObjective...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractServiceLevelObjectiveFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.serviceLevelObjectiveDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToServiceLevelObjectiveDiffs(c.Config, fieldDiffs, opts)
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
	var ops []serviceLevelObjectiveApiOperation
	if create {
		ops = append(ops, &createServiceLevelObjectiveOperation{})
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
	return applyServiceLevelObjectiveDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyServiceLevelObjectiveDiff(c *Client, ctx context.Context, desired *ServiceLevelObjective, rawDesired *ServiceLevelObjective, ops []serviceLevelObjectiveApiOperation, opts ...dcl.ApplyOption) (*ServiceLevelObjective, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetServiceLevelObjective(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createServiceLevelObjectiveOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapServiceLevelObjective(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeServiceLevelObjectiveNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeServiceLevelObjectiveNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeServiceLevelObjectiveDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractServiceLevelObjectiveFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractServiceLevelObjectiveFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffServiceLevelObjective(c, newDesired, newState)
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
