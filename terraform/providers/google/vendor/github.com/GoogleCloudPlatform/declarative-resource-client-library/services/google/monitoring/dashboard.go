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

type Dashboard struct {
	Name         *string                `json:"name"`
	DisplayName  *string                `json:"displayName"`
	GridLayout   *DashboardGridLayout   `json:"gridLayout"`
	MosaicLayout *DashboardMosaicLayout `json:"mosaicLayout"`
	RowLayout    *DashboardRowLayout    `json:"rowLayout"`
	ColumnLayout *DashboardColumnLayout `json:"columnLayout"`
	Project      *string                `json:"project"`
	Etag         *string                `json:"etag"`
}

func (r *Dashboard) String() string {
	return dcl.SprintResource(r)
}

// The enum DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum.
type DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum string

// DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnumRef returns a *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnumRef(s string) *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum {
	v := DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum.
type DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum string

// DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnumRef returns a *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnumRef(s string) *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum {
	v := DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum.
type DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum string

// DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnumRef returns a *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnumRef(s string) *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum {
	v := DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum.
type DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum string

// DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnumRef returns a *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnumRef(s string) *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum {
	v := DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum.
type DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum string

// DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnumRef returns a *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnumRef(s string) *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum {
	v := DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"METHOD_UNSPECIFIED", "METHOD_MEAN", "METHOD_MAX", "METHOD_MIN", "METHOD_SUM", "METHOD_LATEST"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum.
type DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum string

// DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnumRef returns a *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnumRef(s string) *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum {
	v := DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"DIRECTION_UNSPECIFIED", "TOP", "BOTTOM"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum.
type DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum string

// DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnumRef returns a *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnumRef(s string) *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum {
	v := DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum.
type DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum string

// DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnumRef returns a *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnumRef(s string) *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum {
	v := DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum.
type DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum string

// DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnumRef returns a *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnumRef(s string) *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum {
	v := DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum.
type DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum string

// DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnumRef returns a *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnumRef(s string) *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum {
	v := DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum.
type DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum string

// DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnumRef returns a *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnumRef(s string) *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum {
	v := DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum.
type DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum string

// DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnumRef returns a *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnumRef(s string) *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum {
	v := DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum.
type DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum string

// DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnumRef returns a *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnumRef(s string) *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum {
	v := DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"METHOD_UNSPECIFIED", "METHOD_MEAN", "METHOD_MAX", "METHOD_MIN", "METHOD_SUM", "METHOD_LATEST"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum.
type DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum string

// DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnumRef returns a *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnumRef(s string) *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum {
	v := DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"DIRECTION_UNSPECIFIED", "TOP", "BOTTOM"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsXyChartDataSetsPlotTypeEnum.
type DashboardGridLayoutWidgetsXyChartDataSetsPlotTypeEnum string

// DashboardGridLayoutWidgetsXyChartDataSetsPlotTypeEnumRef returns a *DashboardGridLayoutWidgetsXyChartDataSetsPlotTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsXyChartDataSetsPlotTypeEnumRef(s string) *DashboardGridLayoutWidgetsXyChartDataSetsPlotTypeEnum {
	v := DashboardGridLayoutWidgetsXyChartDataSetsPlotTypeEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsXyChartDataSetsPlotTypeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"PLOT_TYPE_UNSPECIFIED", "LINE", "STACKED_AREA", "STACKED_BAR", "HEATMAP"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsXyChartDataSetsPlotTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsXyChartThresholdsColorEnum.
type DashboardGridLayoutWidgetsXyChartThresholdsColorEnum string

// DashboardGridLayoutWidgetsXyChartThresholdsColorEnumRef returns a *DashboardGridLayoutWidgetsXyChartThresholdsColorEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsXyChartThresholdsColorEnumRef(s string) *DashboardGridLayoutWidgetsXyChartThresholdsColorEnum {
	v := DashboardGridLayoutWidgetsXyChartThresholdsColorEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsXyChartThresholdsColorEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"COLOR_UNSPECIFIED", "GREY", "BLUE", "GREEN", "YELLOW", "ORANGE", "RED"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsXyChartThresholdsColorEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsXyChartThresholdsDirectionEnum.
type DashboardGridLayoutWidgetsXyChartThresholdsDirectionEnum string

// DashboardGridLayoutWidgetsXyChartThresholdsDirectionEnumRef returns a *DashboardGridLayoutWidgetsXyChartThresholdsDirectionEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsXyChartThresholdsDirectionEnumRef(s string) *DashboardGridLayoutWidgetsXyChartThresholdsDirectionEnum {
	v := DashboardGridLayoutWidgetsXyChartThresholdsDirectionEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsXyChartThresholdsDirectionEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"DIRECTION_UNSPECIFIED", "ABOVE", "BELOW"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsXyChartThresholdsDirectionEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsXyChartXAxisScaleEnum.
type DashboardGridLayoutWidgetsXyChartXAxisScaleEnum string

// DashboardGridLayoutWidgetsXyChartXAxisScaleEnumRef returns a *DashboardGridLayoutWidgetsXyChartXAxisScaleEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsXyChartXAxisScaleEnumRef(s string) *DashboardGridLayoutWidgetsXyChartXAxisScaleEnum {
	v := DashboardGridLayoutWidgetsXyChartXAxisScaleEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsXyChartXAxisScaleEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"SCALE_UNSPECIFIED", "LINEAR", "LOG10"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsXyChartXAxisScaleEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsXyChartYAxisScaleEnum.
type DashboardGridLayoutWidgetsXyChartYAxisScaleEnum string

// DashboardGridLayoutWidgetsXyChartYAxisScaleEnumRef returns a *DashboardGridLayoutWidgetsXyChartYAxisScaleEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsXyChartYAxisScaleEnumRef(s string) *DashboardGridLayoutWidgetsXyChartYAxisScaleEnum {
	v := DashboardGridLayoutWidgetsXyChartYAxisScaleEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsXyChartYAxisScaleEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"SCALE_UNSPECIFIED", "LINEAR", "LOG10"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsXyChartYAxisScaleEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsXyChartChartOptionsModeEnum.
type DashboardGridLayoutWidgetsXyChartChartOptionsModeEnum string

// DashboardGridLayoutWidgetsXyChartChartOptionsModeEnumRef returns a *DashboardGridLayoutWidgetsXyChartChartOptionsModeEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsXyChartChartOptionsModeEnumRef(s string) *DashboardGridLayoutWidgetsXyChartChartOptionsModeEnum {
	v := DashboardGridLayoutWidgetsXyChartChartOptionsModeEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsXyChartChartOptionsModeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"MODE_UNSPECIFIED", "COLOR", "X_RAY", "STATS"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsXyChartChartOptionsModeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum.
type DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum string

// DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnumRef returns a *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnumRef(s string) *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum {
	v := DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum.
type DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum string

// DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnumRef returns a *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnumRef(s string) *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum {
	v := DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum.
type DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum string

// DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnumRef returns a *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnumRef(s string) *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum {
	v := DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum.
type DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum string

// DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnumRef returns a *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnumRef(s string) *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum {
	v := DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum.
type DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum string

// DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnumRef returns a *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnumRef(s string) *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum {
	v := DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"METHOD_UNSPECIFIED", "METHOD_MEAN", "METHOD_MAX", "METHOD_MIN", "METHOD_SUM", "METHOD_LATEST"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum.
type DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum string

// DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnumRef returns a *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnumRef(s string) *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum {
	v := DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"DIRECTION_UNSPECIFIED", "TOP", "BOTTOM"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum.
type DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum string

// DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnumRef returns a *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnumRef(s string) *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum {
	v := DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum.
type DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum string

// DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnumRef returns a *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnumRef(s string) *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum {
	v := DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum.
type DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum string

// DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnumRef returns a *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnumRef(s string) *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum {
	v := DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum.
type DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum string

// DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnumRef returns a *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnumRef(s string) *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum {
	v := DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum.
type DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum string

// DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnumRef returns a *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnumRef(s string) *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum {
	v := DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum.
type DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum string

// DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnumRef returns a *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnumRef(s string) *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum {
	v := DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum.
type DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum string

// DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnumRef returns a *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnumRef(s string) *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum {
	v := DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"METHOD_UNSPECIFIED", "METHOD_MEAN", "METHOD_MAX", "METHOD_MIN", "METHOD_SUM", "METHOD_LATEST"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum.
type DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum string

// DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnumRef returns a *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnumRef(s string) *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum {
	v := DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"DIRECTION_UNSPECIFIED", "TOP", "BOTTOM"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsScorecardSparkChartViewSparkChartTypeEnum.
type DashboardGridLayoutWidgetsScorecardSparkChartViewSparkChartTypeEnum string

// DashboardGridLayoutWidgetsScorecardSparkChartViewSparkChartTypeEnumRef returns a *DashboardGridLayoutWidgetsScorecardSparkChartViewSparkChartTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsScorecardSparkChartViewSparkChartTypeEnumRef(s string) *DashboardGridLayoutWidgetsScorecardSparkChartViewSparkChartTypeEnum {
	v := DashboardGridLayoutWidgetsScorecardSparkChartViewSparkChartTypeEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsScorecardSparkChartViewSparkChartTypeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"SPARK_CHART_TYPE_UNSPECIFIED", "SPARK_LINE", "SPARK_BAR"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsScorecardSparkChartViewSparkChartTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsScorecardThresholdsColorEnum.
type DashboardGridLayoutWidgetsScorecardThresholdsColorEnum string

// DashboardGridLayoutWidgetsScorecardThresholdsColorEnumRef returns a *DashboardGridLayoutWidgetsScorecardThresholdsColorEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsScorecardThresholdsColorEnumRef(s string) *DashboardGridLayoutWidgetsScorecardThresholdsColorEnum {
	v := DashboardGridLayoutWidgetsScorecardThresholdsColorEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsScorecardThresholdsColorEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"COLOR_UNSPECIFIED", "GREY", "BLUE", "GREEN", "YELLOW", "ORANGE", "RED"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsScorecardThresholdsColorEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsScorecardThresholdsDirectionEnum.
type DashboardGridLayoutWidgetsScorecardThresholdsDirectionEnum string

// DashboardGridLayoutWidgetsScorecardThresholdsDirectionEnumRef returns a *DashboardGridLayoutWidgetsScorecardThresholdsDirectionEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsScorecardThresholdsDirectionEnumRef(s string) *DashboardGridLayoutWidgetsScorecardThresholdsDirectionEnum {
	v := DashboardGridLayoutWidgetsScorecardThresholdsDirectionEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsScorecardThresholdsDirectionEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"DIRECTION_UNSPECIFIED", "ABOVE", "BELOW"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsScorecardThresholdsDirectionEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardGridLayoutWidgetsTextFormatEnum.
type DashboardGridLayoutWidgetsTextFormatEnum string

// DashboardGridLayoutWidgetsTextFormatEnumRef returns a *DashboardGridLayoutWidgetsTextFormatEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardGridLayoutWidgetsTextFormatEnumRef(s string) *DashboardGridLayoutWidgetsTextFormatEnum {
	v := DashboardGridLayoutWidgetsTextFormatEnum(s)
	return &v
}

func (v DashboardGridLayoutWidgetsTextFormatEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"FORMAT_UNSPECIFIED", "MARKDOWN", "RAW"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardGridLayoutWidgetsTextFormatEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum.
type DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum string

// DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnumRef returns a *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnumRef(s string) *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum {
	v := DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum.
type DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum string

// DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnumRef returns a *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnumRef(s string) *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum {
	v := DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum.
type DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum string

// DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnumRef returns a *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnumRef(s string) *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum {
	v := DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum.
type DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum string

// DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnumRef returns a *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnumRef(s string) *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum {
	v := DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum.
type DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum string

// DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnumRef returns a *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnumRef(s string) *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum {
	v := DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"METHOD_UNSPECIFIED", "METHOD_MEAN", "METHOD_MAX", "METHOD_MIN", "METHOD_SUM", "METHOD_LATEST"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum.
type DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum string

// DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnumRef returns a *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnumRef(s string) *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum {
	v := DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"DIRECTION_UNSPECIFIED", "TOP", "BOTTOM"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum.
type DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum string

// DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnumRef returns a *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnumRef(s string) *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum {
	v := DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum.
type DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum string

// DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnumRef returns a *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnumRef(s string) *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum {
	v := DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum.
type DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum string

// DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnumRef returns a *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnumRef(s string) *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum {
	v := DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum.
type DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum string

// DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnumRef returns a *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnumRef(s string) *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum {
	v := DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum.
type DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum string

// DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnumRef returns a *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnumRef(s string) *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum {
	v := DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum.
type DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum string

// DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnumRef returns a *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnumRef(s string) *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum {
	v := DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum.
type DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum string

// DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnumRef returns a *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnumRef(s string) *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum {
	v := DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"METHOD_UNSPECIFIED", "METHOD_MEAN", "METHOD_MAX", "METHOD_MIN", "METHOD_SUM", "METHOD_LATEST"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum.
type DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum string

// DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnumRef returns a *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnumRef(s string) *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum {
	v := DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"DIRECTION_UNSPECIFIED", "TOP", "BOTTOM"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetXyChartDataSetsPlotTypeEnum.
type DashboardMosaicLayoutTilesWidgetXyChartDataSetsPlotTypeEnum string

// DashboardMosaicLayoutTilesWidgetXyChartDataSetsPlotTypeEnumRef returns a *DashboardMosaicLayoutTilesWidgetXyChartDataSetsPlotTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetXyChartDataSetsPlotTypeEnumRef(s string) *DashboardMosaicLayoutTilesWidgetXyChartDataSetsPlotTypeEnum {
	v := DashboardMosaicLayoutTilesWidgetXyChartDataSetsPlotTypeEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetXyChartDataSetsPlotTypeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"PLOT_TYPE_UNSPECIFIED", "LINE", "STACKED_AREA", "STACKED_BAR", "HEATMAP"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetXyChartDataSetsPlotTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetXyChartThresholdsColorEnum.
type DashboardMosaicLayoutTilesWidgetXyChartThresholdsColorEnum string

// DashboardMosaicLayoutTilesWidgetXyChartThresholdsColorEnumRef returns a *DashboardMosaicLayoutTilesWidgetXyChartThresholdsColorEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetXyChartThresholdsColorEnumRef(s string) *DashboardMosaicLayoutTilesWidgetXyChartThresholdsColorEnum {
	v := DashboardMosaicLayoutTilesWidgetXyChartThresholdsColorEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetXyChartThresholdsColorEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"COLOR_UNSPECIFIED", "GREY", "BLUE", "GREEN", "YELLOW", "ORANGE", "RED"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetXyChartThresholdsColorEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetXyChartThresholdsDirectionEnum.
type DashboardMosaicLayoutTilesWidgetXyChartThresholdsDirectionEnum string

// DashboardMosaicLayoutTilesWidgetXyChartThresholdsDirectionEnumRef returns a *DashboardMosaicLayoutTilesWidgetXyChartThresholdsDirectionEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetXyChartThresholdsDirectionEnumRef(s string) *DashboardMosaicLayoutTilesWidgetXyChartThresholdsDirectionEnum {
	v := DashboardMosaicLayoutTilesWidgetXyChartThresholdsDirectionEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetXyChartThresholdsDirectionEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"DIRECTION_UNSPECIFIED", "ABOVE", "BELOW"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetXyChartThresholdsDirectionEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetXyChartXAxisScaleEnum.
type DashboardMosaicLayoutTilesWidgetXyChartXAxisScaleEnum string

// DashboardMosaicLayoutTilesWidgetXyChartXAxisScaleEnumRef returns a *DashboardMosaicLayoutTilesWidgetXyChartXAxisScaleEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetXyChartXAxisScaleEnumRef(s string) *DashboardMosaicLayoutTilesWidgetXyChartXAxisScaleEnum {
	v := DashboardMosaicLayoutTilesWidgetXyChartXAxisScaleEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetXyChartXAxisScaleEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"SCALE_UNSPECIFIED", "LINEAR", "LOG10"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetXyChartXAxisScaleEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetXyChartYAxisScaleEnum.
type DashboardMosaicLayoutTilesWidgetXyChartYAxisScaleEnum string

// DashboardMosaicLayoutTilesWidgetXyChartYAxisScaleEnumRef returns a *DashboardMosaicLayoutTilesWidgetXyChartYAxisScaleEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetXyChartYAxisScaleEnumRef(s string) *DashboardMosaicLayoutTilesWidgetXyChartYAxisScaleEnum {
	v := DashboardMosaicLayoutTilesWidgetXyChartYAxisScaleEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetXyChartYAxisScaleEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"SCALE_UNSPECIFIED", "LINEAR", "LOG10"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetXyChartYAxisScaleEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetXyChartChartOptionsModeEnum.
type DashboardMosaicLayoutTilesWidgetXyChartChartOptionsModeEnum string

// DashboardMosaicLayoutTilesWidgetXyChartChartOptionsModeEnumRef returns a *DashboardMosaicLayoutTilesWidgetXyChartChartOptionsModeEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetXyChartChartOptionsModeEnumRef(s string) *DashboardMosaicLayoutTilesWidgetXyChartChartOptionsModeEnum {
	v := DashboardMosaicLayoutTilesWidgetXyChartChartOptionsModeEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetXyChartChartOptionsModeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"MODE_UNSPECIFIED", "COLOR", "X_RAY", "STATS"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetXyChartChartOptionsModeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum.
type DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum string

// DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnumRef returns a *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnumRef(s string) *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum {
	v := DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum.
type DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum string

// DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnumRef returns a *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnumRef(s string) *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum {
	v := DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum.
type DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum string

// DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnumRef returns a *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnumRef(s string) *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum {
	v := DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum.
type DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum string

// DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnumRef returns a *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnumRef(s string) *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum {
	v := DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum.
type DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum string

// DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnumRef returns a *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnumRef(s string) *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum {
	v := DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"METHOD_UNSPECIFIED", "METHOD_MEAN", "METHOD_MAX", "METHOD_MIN", "METHOD_SUM", "METHOD_LATEST"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum.
type DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum string

// DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnumRef returns a *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnumRef(s string) *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum {
	v := DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"DIRECTION_UNSPECIFIED", "TOP", "BOTTOM"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum.
type DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum string

// DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnumRef returns a *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnumRef(s string) *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum {
	v := DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum.
type DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum string

// DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnumRef returns a *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnumRef(s string) *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum {
	v := DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum.
type DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum string

// DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnumRef returns a *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnumRef(s string) *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum {
	v := DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum.
type DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum string

// DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnumRef returns a *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnumRef(s string) *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum {
	v := DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum.
type DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum string

// DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnumRef returns a *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnumRef(s string) *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum {
	v := DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum.
type DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum string

// DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnumRef returns a *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnumRef(s string) *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum {
	v := DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum.
type DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum string

// DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnumRef returns a *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnumRef(s string) *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum {
	v := DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"METHOD_UNSPECIFIED", "METHOD_MEAN", "METHOD_MAX", "METHOD_MIN", "METHOD_SUM", "METHOD_LATEST"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum.
type DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum string

// DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnumRef returns a *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnumRef(s string) *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum {
	v := DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"DIRECTION_UNSPECIFIED", "TOP", "BOTTOM"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetScorecardSparkChartViewSparkChartTypeEnum.
type DashboardMosaicLayoutTilesWidgetScorecardSparkChartViewSparkChartTypeEnum string

// DashboardMosaicLayoutTilesWidgetScorecardSparkChartViewSparkChartTypeEnumRef returns a *DashboardMosaicLayoutTilesWidgetScorecardSparkChartViewSparkChartTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetScorecardSparkChartViewSparkChartTypeEnumRef(s string) *DashboardMosaicLayoutTilesWidgetScorecardSparkChartViewSparkChartTypeEnum {
	v := DashboardMosaicLayoutTilesWidgetScorecardSparkChartViewSparkChartTypeEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetScorecardSparkChartViewSparkChartTypeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"SPARK_CHART_TYPE_UNSPECIFIED", "SPARK_LINE", "SPARK_BAR"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetScorecardSparkChartViewSparkChartTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetScorecardThresholdsColorEnum.
type DashboardMosaicLayoutTilesWidgetScorecardThresholdsColorEnum string

// DashboardMosaicLayoutTilesWidgetScorecardThresholdsColorEnumRef returns a *DashboardMosaicLayoutTilesWidgetScorecardThresholdsColorEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetScorecardThresholdsColorEnumRef(s string) *DashboardMosaicLayoutTilesWidgetScorecardThresholdsColorEnum {
	v := DashboardMosaicLayoutTilesWidgetScorecardThresholdsColorEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetScorecardThresholdsColorEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"COLOR_UNSPECIFIED", "GREY", "BLUE", "GREEN", "YELLOW", "ORANGE", "RED"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetScorecardThresholdsColorEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetScorecardThresholdsDirectionEnum.
type DashboardMosaicLayoutTilesWidgetScorecardThresholdsDirectionEnum string

// DashboardMosaicLayoutTilesWidgetScorecardThresholdsDirectionEnumRef returns a *DashboardMosaicLayoutTilesWidgetScorecardThresholdsDirectionEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetScorecardThresholdsDirectionEnumRef(s string) *DashboardMosaicLayoutTilesWidgetScorecardThresholdsDirectionEnum {
	v := DashboardMosaicLayoutTilesWidgetScorecardThresholdsDirectionEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetScorecardThresholdsDirectionEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"DIRECTION_UNSPECIFIED", "ABOVE", "BELOW"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetScorecardThresholdsDirectionEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardMosaicLayoutTilesWidgetTextFormatEnum.
type DashboardMosaicLayoutTilesWidgetTextFormatEnum string

// DashboardMosaicLayoutTilesWidgetTextFormatEnumRef returns a *DashboardMosaicLayoutTilesWidgetTextFormatEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardMosaicLayoutTilesWidgetTextFormatEnumRef(s string) *DashboardMosaicLayoutTilesWidgetTextFormatEnum {
	v := DashboardMosaicLayoutTilesWidgetTextFormatEnum(s)
	return &v
}

func (v DashboardMosaicLayoutTilesWidgetTextFormatEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"FORMAT_UNSPECIFIED", "MARKDOWN", "RAW"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardMosaicLayoutTilesWidgetTextFormatEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum.
type DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum string

// DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnumRef returns a *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnumRef(s string) *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum {
	v := DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum.
type DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum string

// DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnumRef returns a *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnumRef(s string) *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum {
	v := DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum.
type DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum string

// DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnumRef returns a *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnumRef(s string) *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum {
	v := DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum.
type DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum string

// DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnumRef returns a *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnumRef(s string) *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum {
	v := DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum.
type DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum string

// DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnumRef returns a *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnumRef(s string) *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum {
	v := DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"METHOD_UNSPECIFIED", "METHOD_MEAN", "METHOD_MAX", "METHOD_MIN", "METHOD_SUM", "METHOD_LATEST"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum.
type DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum string

// DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnumRef returns a *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnumRef(s string) *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum {
	v := DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"DIRECTION_UNSPECIFIED", "TOP", "BOTTOM"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum.
type DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum string

// DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnumRef returns a *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnumRef(s string) *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum {
	v := DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum.
type DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum string

// DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnumRef returns a *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnumRef(s string) *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum {
	v := DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum.
type DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum string

// DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnumRef returns a *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnumRef(s string) *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum {
	v := DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum.
type DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum string

// DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnumRef returns a *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnumRef(s string) *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum {
	v := DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum.
type DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum string

// DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnumRef returns a *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnumRef(s string) *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum {
	v := DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum.
type DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum string

// DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnumRef returns a *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnumRef(s string) *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum {
	v := DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum.
type DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum string

// DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnumRef returns a *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnumRef(s string) *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum {
	v := DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"METHOD_UNSPECIFIED", "METHOD_MEAN", "METHOD_MAX", "METHOD_MIN", "METHOD_SUM", "METHOD_LATEST"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum.
type DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum string

// DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnumRef returns a *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnumRef(s string) *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum {
	v := DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"DIRECTION_UNSPECIFIED", "TOP", "BOTTOM"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsXyChartDataSetsPlotTypeEnum.
type DashboardRowLayoutRowsWidgetsXyChartDataSetsPlotTypeEnum string

// DashboardRowLayoutRowsWidgetsXyChartDataSetsPlotTypeEnumRef returns a *DashboardRowLayoutRowsWidgetsXyChartDataSetsPlotTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsXyChartDataSetsPlotTypeEnumRef(s string) *DashboardRowLayoutRowsWidgetsXyChartDataSetsPlotTypeEnum {
	v := DashboardRowLayoutRowsWidgetsXyChartDataSetsPlotTypeEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsXyChartDataSetsPlotTypeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"PLOT_TYPE_UNSPECIFIED", "LINE", "STACKED_AREA", "STACKED_BAR", "HEATMAP"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsXyChartDataSetsPlotTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsXyChartThresholdsColorEnum.
type DashboardRowLayoutRowsWidgetsXyChartThresholdsColorEnum string

// DashboardRowLayoutRowsWidgetsXyChartThresholdsColorEnumRef returns a *DashboardRowLayoutRowsWidgetsXyChartThresholdsColorEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsXyChartThresholdsColorEnumRef(s string) *DashboardRowLayoutRowsWidgetsXyChartThresholdsColorEnum {
	v := DashboardRowLayoutRowsWidgetsXyChartThresholdsColorEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsXyChartThresholdsColorEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"COLOR_UNSPECIFIED", "GREY", "BLUE", "GREEN", "YELLOW", "ORANGE", "RED"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsXyChartThresholdsColorEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsXyChartThresholdsDirectionEnum.
type DashboardRowLayoutRowsWidgetsXyChartThresholdsDirectionEnum string

// DashboardRowLayoutRowsWidgetsXyChartThresholdsDirectionEnumRef returns a *DashboardRowLayoutRowsWidgetsXyChartThresholdsDirectionEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsXyChartThresholdsDirectionEnumRef(s string) *DashboardRowLayoutRowsWidgetsXyChartThresholdsDirectionEnum {
	v := DashboardRowLayoutRowsWidgetsXyChartThresholdsDirectionEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsXyChartThresholdsDirectionEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"DIRECTION_UNSPECIFIED", "ABOVE", "BELOW"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsXyChartThresholdsDirectionEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsXyChartXAxisScaleEnum.
type DashboardRowLayoutRowsWidgetsXyChartXAxisScaleEnum string

// DashboardRowLayoutRowsWidgetsXyChartXAxisScaleEnumRef returns a *DashboardRowLayoutRowsWidgetsXyChartXAxisScaleEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsXyChartXAxisScaleEnumRef(s string) *DashboardRowLayoutRowsWidgetsXyChartXAxisScaleEnum {
	v := DashboardRowLayoutRowsWidgetsXyChartXAxisScaleEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsXyChartXAxisScaleEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"SCALE_UNSPECIFIED", "LINEAR", "LOG10"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsXyChartXAxisScaleEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsXyChartYAxisScaleEnum.
type DashboardRowLayoutRowsWidgetsXyChartYAxisScaleEnum string

// DashboardRowLayoutRowsWidgetsXyChartYAxisScaleEnumRef returns a *DashboardRowLayoutRowsWidgetsXyChartYAxisScaleEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsXyChartYAxisScaleEnumRef(s string) *DashboardRowLayoutRowsWidgetsXyChartYAxisScaleEnum {
	v := DashboardRowLayoutRowsWidgetsXyChartYAxisScaleEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsXyChartYAxisScaleEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"SCALE_UNSPECIFIED", "LINEAR", "LOG10"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsXyChartYAxisScaleEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsXyChartChartOptionsModeEnum.
type DashboardRowLayoutRowsWidgetsXyChartChartOptionsModeEnum string

// DashboardRowLayoutRowsWidgetsXyChartChartOptionsModeEnumRef returns a *DashboardRowLayoutRowsWidgetsXyChartChartOptionsModeEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsXyChartChartOptionsModeEnumRef(s string) *DashboardRowLayoutRowsWidgetsXyChartChartOptionsModeEnum {
	v := DashboardRowLayoutRowsWidgetsXyChartChartOptionsModeEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsXyChartChartOptionsModeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"MODE_UNSPECIFIED", "COLOR", "X_RAY", "STATS"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsXyChartChartOptionsModeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum.
type DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum string

// DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnumRef returns a *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnumRef(s string) *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum {
	v := DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum.
type DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum string

// DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnumRef returns a *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnumRef(s string) *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum {
	v := DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum.
type DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum string

// DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnumRef returns a *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnumRef(s string) *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum {
	v := DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum.
type DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum string

// DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnumRef returns a *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnumRef(s string) *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum {
	v := DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum.
type DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum string

// DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnumRef returns a *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnumRef(s string) *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum {
	v := DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"METHOD_UNSPECIFIED", "METHOD_MEAN", "METHOD_MAX", "METHOD_MIN", "METHOD_SUM", "METHOD_LATEST"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum.
type DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum string

// DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnumRef returns a *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnumRef(s string) *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum {
	v := DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"DIRECTION_UNSPECIFIED", "TOP", "BOTTOM"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum.
type DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum string

// DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnumRef returns a *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnumRef(s string) *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum {
	v := DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum.
type DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum string

// DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnumRef returns a *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnumRef(s string) *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum {
	v := DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum.
type DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum string

// DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnumRef returns a *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnumRef(s string) *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum {
	v := DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum.
type DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum string

// DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnumRef returns a *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnumRef(s string) *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum {
	v := DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum.
type DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum string

// DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnumRef returns a *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnumRef(s string) *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum {
	v := DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum.
type DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum string

// DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnumRef returns a *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnumRef(s string) *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum {
	v := DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum.
type DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum string

// DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnumRef returns a *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnumRef(s string) *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum {
	v := DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"METHOD_UNSPECIFIED", "METHOD_MEAN", "METHOD_MAX", "METHOD_MIN", "METHOD_SUM", "METHOD_LATEST"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum.
type DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum string

// DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnumRef returns a *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnumRef(s string) *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum {
	v := DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"DIRECTION_UNSPECIFIED", "TOP", "BOTTOM"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsScorecardSparkChartViewSparkChartTypeEnum.
type DashboardRowLayoutRowsWidgetsScorecardSparkChartViewSparkChartTypeEnum string

// DashboardRowLayoutRowsWidgetsScorecardSparkChartViewSparkChartTypeEnumRef returns a *DashboardRowLayoutRowsWidgetsScorecardSparkChartViewSparkChartTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsScorecardSparkChartViewSparkChartTypeEnumRef(s string) *DashboardRowLayoutRowsWidgetsScorecardSparkChartViewSparkChartTypeEnum {
	v := DashboardRowLayoutRowsWidgetsScorecardSparkChartViewSparkChartTypeEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsScorecardSparkChartViewSparkChartTypeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"SPARK_CHART_TYPE_UNSPECIFIED", "SPARK_LINE", "SPARK_BAR"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsScorecardSparkChartViewSparkChartTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsScorecardThresholdsColorEnum.
type DashboardRowLayoutRowsWidgetsScorecardThresholdsColorEnum string

// DashboardRowLayoutRowsWidgetsScorecardThresholdsColorEnumRef returns a *DashboardRowLayoutRowsWidgetsScorecardThresholdsColorEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsScorecardThresholdsColorEnumRef(s string) *DashboardRowLayoutRowsWidgetsScorecardThresholdsColorEnum {
	v := DashboardRowLayoutRowsWidgetsScorecardThresholdsColorEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsScorecardThresholdsColorEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"COLOR_UNSPECIFIED", "GREY", "BLUE", "GREEN", "YELLOW", "ORANGE", "RED"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsScorecardThresholdsColorEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsScorecardThresholdsDirectionEnum.
type DashboardRowLayoutRowsWidgetsScorecardThresholdsDirectionEnum string

// DashboardRowLayoutRowsWidgetsScorecardThresholdsDirectionEnumRef returns a *DashboardRowLayoutRowsWidgetsScorecardThresholdsDirectionEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsScorecardThresholdsDirectionEnumRef(s string) *DashboardRowLayoutRowsWidgetsScorecardThresholdsDirectionEnum {
	v := DashboardRowLayoutRowsWidgetsScorecardThresholdsDirectionEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsScorecardThresholdsDirectionEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"DIRECTION_UNSPECIFIED", "ABOVE", "BELOW"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsScorecardThresholdsDirectionEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardRowLayoutRowsWidgetsTextFormatEnum.
type DashboardRowLayoutRowsWidgetsTextFormatEnum string

// DashboardRowLayoutRowsWidgetsTextFormatEnumRef returns a *DashboardRowLayoutRowsWidgetsTextFormatEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardRowLayoutRowsWidgetsTextFormatEnumRef(s string) *DashboardRowLayoutRowsWidgetsTextFormatEnum {
	v := DashboardRowLayoutRowsWidgetsTextFormatEnum(s)
	return &v
}

func (v DashboardRowLayoutRowsWidgetsTextFormatEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"FORMAT_UNSPECIFIED", "MARKDOWN", "RAW"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardRowLayoutRowsWidgetsTextFormatEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum.
type DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum string

// DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnumRef returns a *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum {
	v := DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum.
type DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum string

// DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnumRef returns a *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum {
	v := DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum.
type DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum string

// DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnumRef returns a *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum {
	v := DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum.
type DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum string

// DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnumRef returns a *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum {
	v := DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum.
type DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum string

// DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnumRef returns a *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum {
	v := DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"METHOD_UNSPECIFIED", "METHOD_MEAN", "METHOD_MAX", "METHOD_MIN", "METHOD_SUM", "METHOD_LATEST"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum.
type DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum string

// DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnumRef returns a *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum {
	v := DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"DIRECTION_UNSPECIFIED", "TOP", "BOTTOM"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum.
type DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum string

// DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnumRef returns a *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum {
	v := DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum.
type DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum string

// DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnumRef returns a *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum {
	v := DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum.
type DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum string

// DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnumRef returns a *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum {
	v := DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum.
type DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum string

// DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnumRef returns a *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum {
	v := DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum.
type DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum string

// DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnumRef returns a *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum {
	v := DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum.
type DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum string

// DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnumRef returns a *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum {
	v := DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum.
type DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum string

// DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnumRef returns a *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum {
	v := DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"METHOD_UNSPECIFIED", "METHOD_MEAN", "METHOD_MAX", "METHOD_MIN", "METHOD_SUM", "METHOD_LATEST"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum.
type DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum string

// DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnumRef returns a *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum {
	v := DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"DIRECTION_UNSPECIFIED", "TOP", "BOTTOM"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsXyChartDataSetsPlotTypeEnum.
type DashboardColumnLayoutColumnsWidgetsXyChartDataSetsPlotTypeEnum string

// DashboardColumnLayoutColumnsWidgetsXyChartDataSetsPlotTypeEnumRef returns a *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsPlotTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsXyChartDataSetsPlotTypeEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsPlotTypeEnum {
	v := DashboardColumnLayoutColumnsWidgetsXyChartDataSetsPlotTypeEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsXyChartDataSetsPlotTypeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"PLOT_TYPE_UNSPECIFIED", "LINE", "STACKED_AREA", "STACKED_BAR", "HEATMAP"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsXyChartDataSetsPlotTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsXyChartThresholdsColorEnum.
type DashboardColumnLayoutColumnsWidgetsXyChartThresholdsColorEnum string

// DashboardColumnLayoutColumnsWidgetsXyChartThresholdsColorEnumRef returns a *DashboardColumnLayoutColumnsWidgetsXyChartThresholdsColorEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsXyChartThresholdsColorEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsXyChartThresholdsColorEnum {
	v := DashboardColumnLayoutColumnsWidgetsXyChartThresholdsColorEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsXyChartThresholdsColorEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"COLOR_UNSPECIFIED", "GREY", "BLUE", "GREEN", "YELLOW", "ORANGE", "RED"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsXyChartThresholdsColorEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsXyChartThresholdsDirectionEnum.
type DashboardColumnLayoutColumnsWidgetsXyChartThresholdsDirectionEnum string

// DashboardColumnLayoutColumnsWidgetsXyChartThresholdsDirectionEnumRef returns a *DashboardColumnLayoutColumnsWidgetsXyChartThresholdsDirectionEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsXyChartThresholdsDirectionEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsXyChartThresholdsDirectionEnum {
	v := DashboardColumnLayoutColumnsWidgetsXyChartThresholdsDirectionEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsXyChartThresholdsDirectionEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"DIRECTION_UNSPECIFIED", "ABOVE", "BELOW"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsXyChartThresholdsDirectionEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsXyChartXAxisScaleEnum.
type DashboardColumnLayoutColumnsWidgetsXyChartXAxisScaleEnum string

// DashboardColumnLayoutColumnsWidgetsXyChartXAxisScaleEnumRef returns a *DashboardColumnLayoutColumnsWidgetsXyChartXAxisScaleEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsXyChartXAxisScaleEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsXyChartXAxisScaleEnum {
	v := DashboardColumnLayoutColumnsWidgetsXyChartXAxisScaleEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsXyChartXAxisScaleEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"SCALE_UNSPECIFIED", "LINEAR", "LOG10"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsXyChartXAxisScaleEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsXyChartYAxisScaleEnum.
type DashboardColumnLayoutColumnsWidgetsXyChartYAxisScaleEnum string

// DashboardColumnLayoutColumnsWidgetsXyChartYAxisScaleEnumRef returns a *DashboardColumnLayoutColumnsWidgetsXyChartYAxisScaleEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsXyChartYAxisScaleEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsXyChartYAxisScaleEnum {
	v := DashboardColumnLayoutColumnsWidgetsXyChartYAxisScaleEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsXyChartYAxisScaleEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"SCALE_UNSPECIFIED", "LINEAR", "LOG10"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsXyChartYAxisScaleEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsXyChartChartOptionsModeEnum.
type DashboardColumnLayoutColumnsWidgetsXyChartChartOptionsModeEnum string

// DashboardColumnLayoutColumnsWidgetsXyChartChartOptionsModeEnumRef returns a *DashboardColumnLayoutColumnsWidgetsXyChartChartOptionsModeEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsXyChartChartOptionsModeEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsXyChartChartOptionsModeEnum {
	v := DashboardColumnLayoutColumnsWidgetsXyChartChartOptionsModeEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsXyChartChartOptionsModeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"MODE_UNSPECIFIED", "COLOR", "X_RAY", "STATS"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsXyChartChartOptionsModeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum.
type DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum string

// DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnumRef returns a *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum {
	v := DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum.
type DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum string

// DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnumRef returns a *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum {
	v := DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum.
type DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum string

// DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnumRef returns a *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum {
	v := DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum.
type DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum string

// DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnumRef returns a *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum {
	v := DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum.
type DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum string

// DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnumRef returns a *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum {
	v := DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"METHOD_UNSPECIFIED", "METHOD_MEAN", "METHOD_MAX", "METHOD_MIN", "METHOD_SUM", "METHOD_LATEST"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum.
type DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum string

// DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnumRef returns a *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum {
	v := DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"DIRECTION_UNSPECIFIED", "TOP", "BOTTOM"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum.
type DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum string

// DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnumRef returns a *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum {
	v := DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum.
type DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum string

// DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnumRef returns a *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum {
	v := DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum.
type DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum string

// DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnumRef returns a *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum {
	v := DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum.
type DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum string

// DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnumRef returns a *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum {
	v := DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum.
type DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum string

// DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnumRef returns a *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum {
	v := DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ALIGN_NONE", "ALIGN_DELTA", "ALIGN_RATE", "ALIGN_INTERPOLATE", "ALIGN_NEXT_OLDER", "ALIGN_MIN", "ALIGN_MAX", "ALIGN_MEAN", "ALIGN_COUNT", "ALIGN_SUM", "ALIGN_STDDEV", "ALIGN_COUNT_TRUE", "ALIGN_COUNT_FALSE", "ALIGN_FRACTION_TRUE", "ALIGN_PERCENTILE_99", "ALIGN_PERCENTILE_95", "ALIGN_PERCENTILE_50", "ALIGN_PERCENTILE_05", "ALIGN_MAKE_DISTRIBUTION", "ALIGN_PERCENT_CHANGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum.
type DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum string

// DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnumRef returns a *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum {
	v := DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REDUCE_NONE", "REDUCE_MEAN", "REDUCE_MIN", "REDUCE_MAX", "REDUCE_SUM", "REDUCE_STDDEV", "REDUCE_COUNT", "REDUCE_COUNT_TRUE", "REDUCE_COUNT_FALSE", "REDUCE_FRACTION_TRUE", "REDUCE_PERCENTILE_99", "REDUCE_PERCENTILE_95", "REDUCE_PERCENTILE_50", "REDUCE_PERCENTILE_05", "REDUCE_FRACTION_LESS_THAN", "REDUCE_MAKE_DISTRIBUTION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum.
type DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum string

// DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnumRef returns a *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum {
	v := DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"METHOD_UNSPECIFIED", "METHOD_MEAN", "METHOD_MAX", "METHOD_MIN", "METHOD_SUM", "METHOD_LATEST"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum.
type DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum string

// DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnumRef returns a *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum {
	v := DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"DIRECTION_UNSPECIFIED", "TOP", "BOTTOM"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsScorecardSparkChartViewSparkChartTypeEnum.
type DashboardColumnLayoutColumnsWidgetsScorecardSparkChartViewSparkChartTypeEnum string

// DashboardColumnLayoutColumnsWidgetsScorecardSparkChartViewSparkChartTypeEnumRef returns a *DashboardColumnLayoutColumnsWidgetsScorecardSparkChartViewSparkChartTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsScorecardSparkChartViewSparkChartTypeEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsScorecardSparkChartViewSparkChartTypeEnum {
	v := DashboardColumnLayoutColumnsWidgetsScorecardSparkChartViewSparkChartTypeEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsScorecardSparkChartViewSparkChartTypeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"SPARK_CHART_TYPE_UNSPECIFIED", "SPARK_LINE", "SPARK_BAR"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsScorecardSparkChartViewSparkChartTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsScorecardThresholdsColorEnum.
type DashboardColumnLayoutColumnsWidgetsScorecardThresholdsColorEnum string

// DashboardColumnLayoutColumnsWidgetsScorecardThresholdsColorEnumRef returns a *DashboardColumnLayoutColumnsWidgetsScorecardThresholdsColorEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsScorecardThresholdsColorEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsScorecardThresholdsColorEnum {
	v := DashboardColumnLayoutColumnsWidgetsScorecardThresholdsColorEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsScorecardThresholdsColorEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"COLOR_UNSPECIFIED", "GREY", "BLUE", "GREEN", "YELLOW", "ORANGE", "RED"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsScorecardThresholdsColorEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsScorecardThresholdsDirectionEnum.
type DashboardColumnLayoutColumnsWidgetsScorecardThresholdsDirectionEnum string

// DashboardColumnLayoutColumnsWidgetsScorecardThresholdsDirectionEnumRef returns a *DashboardColumnLayoutColumnsWidgetsScorecardThresholdsDirectionEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsScorecardThresholdsDirectionEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsScorecardThresholdsDirectionEnum {
	v := DashboardColumnLayoutColumnsWidgetsScorecardThresholdsDirectionEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsScorecardThresholdsDirectionEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"DIRECTION_UNSPECIFIED", "ABOVE", "BELOW"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsScorecardThresholdsDirectionEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum DashboardColumnLayoutColumnsWidgetsTextFormatEnum.
type DashboardColumnLayoutColumnsWidgetsTextFormatEnum string

// DashboardColumnLayoutColumnsWidgetsTextFormatEnumRef returns a *DashboardColumnLayoutColumnsWidgetsTextFormatEnum with the value of string s
// If the empty string is provided, nil is returned.
func DashboardColumnLayoutColumnsWidgetsTextFormatEnumRef(s string) *DashboardColumnLayoutColumnsWidgetsTextFormatEnum {
	v := DashboardColumnLayoutColumnsWidgetsTextFormatEnum(s)
	return &v
}

func (v DashboardColumnLayoutColumnsWidgetsTextFormatEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"FORMAT_UNSPECIFIED", "MARKDOWN", "RAW"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "DashboardColumnLayoutColumnsWidgetsTextFormatEnum",
		Value: string(v),
		Valid: []string{},
	}
}

type DashboardGridLayout struct {
	empty   bool                         `json:"-"`
	Columns *int64                       `json:"columns"`
	Widgets []DashboardGridLayoutWidgets `json:"widgets"`
}

type jsonDashboardGridLayout DashboardGridLayout

func (r *DashboardGridLayout) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayout
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayout
	} else {

		r.Columns = res.Columns

		r.Widgets = res.Widgets

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayout is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayout *DashboardGridLayout = &DashboardGridLayout{empty: true}

func (r *DashboardGridLayout) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayout) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayout) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgets struct {
	empty     bool                                 `json:"-"`
	Title     *string                              `json:"title"`
	XyChart   *DashboardGridLayoutWidgetsXyChart   `json:"xyChart"`
	Scorecard *DashboardGridLayoutWidgetsScorecard `json:"scorecard"`
	Text      *DashboardGridLayoutWidgetsText      `json:"text"`
	Blank     *DashboardGridLayoutWidgetsBlank     `json:"blank"`
	LogsPanel *DashboardGridLayoutWidgetsLogsPanel `json:"logsPanel"`
}

type jsonDashboardGridLayoutWidgets DashboardGridLayoutWidgets

func (r *DashboardGridLayoutWidgets) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgets
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgets
	} else {

		r.Title = res.Title

		r.XyChart = res.XyChart

		r.Scorecard = res.Scorecard

		r.Text = res.Text

		r.Blank = res.Blank

		r.LogsPanel = res.LogsPanel

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgets is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgets *DashboardGridLayoutWidgets = &DashboardGridLayoutWidgets{empty: true}

func (r *DashboardGridLayoutWidgets) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgets) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgets) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsXyChart struct {
	empty             bool                                           `json:"-"`
	DataSets          []DashboardGridLayoutWidgetsXyChartDataSets    `json:"dataSets"`
	TimeshiftDuration *string                                        `json:"timeshiftDuration"`
	Thresholds        []DashboardGridLayoutWidgetsXyChartThresholds  `json:"thresholds"`
	XAxis             *DashboardGridLayoutWidgetsXyChartXAxis        `json:"xAxis"`
	YAxis             *DashboardGridLayoutWidgetsXyChartYAxis        `json:"yAxis"`
	ChartOptions      *DashboardGridLayoutWidgetsXyChartChartOptions `json:"chartOptions"`
}

type jsonDashboardGridLayoutWidgetsXyChart DashboardGridLayoutWidgetsXyChart

func (r *DashboardGridLayoutWidgetsXyChart) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsXyChart
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsXyChart
	} else {

		r.DataSets = res.DataSets

		r.TimeshiftDuration = res.TimeshiftDuration

		r.Thresholds = res.Thresholds

		r.XAxis = res.XAxis

		r.YAxis = res.YAxis

		r.ChartOptions = res.ChartOptions

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsXyChart is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsXyChart *DashboardGridLayoutWidgetsXyChart = &DashboardGridLayoutWidgetsXyChart{empty: true}

func (r *DashboardGridLayoutWidgetsXyChart) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsXyChart) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsXyChart) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsXyChartDataSets struct {
	empty              bool                                                      `json:"-"`
	TimeSeriesQuery    *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQuery `json:"timeSeriesQuery"`
	PlotType           *DashboardGridLayoutWidgetsXyChartDataSetsPlotTypeEnum    `json:"plotType"`
	LegendTemplate     *string                                                   `json:"legendTemplate"`
	MinAlignmentPeriod *string                                                   `json:"minAlignmentPeriod"`
}

type jsonDashboardGridLayoutWidgetsXyChartDataSets DashboardGridLayoutWidgetsXyChartDataSets

func (r *DashboardGridLayoutWidgetsXyChartDataSets) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsXyChartDataSets
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsXyChartDataSets
	} else {

		r.TimeSeriesQuery = res.TimeSeriesQuery

		r.PlotType = res.PlotType

		r.LegendTemplate = res.LegendTemplate

		r.MinAlignmentPeriod = res.MinAlignmentPeriod

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsXyChartDataSets is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsXyChartDataSets *DashboardGridLayoutWidgetsXyChartDataSets = &DashboardGridLayoutWidgetsXyChartDataSets{empty: true}

func (r *DashboardGridLayoutWidgetsXyChartDataSets) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsXyChartDataSets) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsXyChartDataSets) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQuery struct {
	empty                   bool                                                                           `json:"-"`
	TimeSeriesFilter        *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter      `json:"timeSeriesFilter"`
	TimeSeriesFilterRatio   *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio `json:"timeSeriesFilterRatio"`
	TimeSeriesQueryLanguage *string                                                                        `json:"timeSeriesQueryLanguage"`
	UnitOverride            *string                                                                        `json:"unitOverride"`
}

type jsonDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQuery DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQuery

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQuery) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQuery
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQuery
	} else {

		r.TimeSeriesFilter = res.TimeSeriesFilter

		r.TimeSeriesFilterRatio = res.TimeSeriesFilterRatio

		r.TimeSeriesQueryLanguage = res.TimeSeriesQueryLanguage

		r.UnitOverride = res.UnitOverride

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQuery is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQuery *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQuery = &DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQuery{empty: true}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQuery) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQuery) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQuery) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter struct {
	empty                bool                                                                                          `json:"-"`
	Filter               *string                                                                                       `json:"filter"`
	Aggregation          *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation          `json:"aggregation"`
	SecondaryAggregation *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation `json:"secondaryAggregation"`
	PickTimeSeriesFilter *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter `json:"pickTimeSeriesFilter"`
}

type jsonDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter
	} else {

		r.Filter = res.Filter

		r.Aggregation = res.Aggregation

		r.SecondaryAggregation = res.SecondaryAggregation

		r.PickTimeSeriesFilter = res.PickTimeSeriesFilter

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter = &DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter{empty: true}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation struct {
	empty              bool                                                                                                       `json:"-"`
	AlignmentPeriod    *string                                                                                                    `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                   `json:"groupByFields"`
}

type jsonDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation = &DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation{empty: true}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation struct {
	empty              bool                                                                                                                `json:"-"`
	AlignmentPeriod    *string                                                                                                             `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                            `json:"groupByFields"`
}

type jsonDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation = &DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation{empty: true}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter struct {
	empty         bool                                                                                                           `json:"-"`
	RankingMethod *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum `json:"rankingMethod"`
	NumTimeSeries *int64                                                                                                         `json:"numTimeSeries"`
	Direction     *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum     `json:"direction"`
}

type jsonDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter
	} else {

		r.RankingMethod = res.RankingMethod

		r.NumTimeSeries = res.NumTimeSeries

		r.Direction = res.Direction

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter = &DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter{empty: true}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio struct {
	empty                bool                                                                                               `json:"-"`
	Numerator            *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator            `json:"numerator"`
	Denominator          *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator          `json:"denominator"`
	SecondaryAggregation *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation `json:"secondaryAggregation"`
	PickTimeSeriesFilter *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter `json:"pickTimeSeriesFilter"`
}

type jsonDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio
	} else {

		r.Numerator = res.Numerator

		r.Denominator = res.Denominator

		r.SecondaryAggregation = res.SecondaryAggregation

		r.PickTimeSeriesFilter = res.PickTimeSeriesFilter

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio = &DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio{empty: true}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator struct {
	empty       bool                                                                                               `json:"-"`
	Filter      *string                                                                                            `json:"filter"`
	Aggregation *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation `json:"aggregation"`
}

type jsonDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator
	} else {

		r.Filter = res.Filter

		r.Aggregation = res.Aggregation

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator = &DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator{empty: true}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation struct {
	empty              bool                                                                                                                     `json:"-"`
	AlignmentPeriod    *string                                                                                                                  `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                                 `json:"groupByFields"`
}

type jsonDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation = &DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation{empty: true}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator struct {
	empty       bool                                                                                                 `json:"-"`
	Filter      *string                                                                                              `json:"filter"`
	Aggregation *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation `json:"aggregation"`
}

type jsonDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator
	} else {

		r.Filter = res.Filter

		r.Aggregation = res.Aggregation

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator = &DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator{empty: true}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation struct {
	empty              bool                                                                                                                       `json:"-"`
	AlignmentPeriod    *string                                                                                                                    `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                                   `json:"groupByFields"`
}

type jsonDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation = &DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation{empty: true}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation struct {
	empty              bool                                                                                                                     `json:"-"`
	AlignmentPeriod    *string                                                                                                                  `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                                 `json:"groupByFields"`
}

type jsonDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation = &DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation{empty: true}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter struct {
	empty         bool                                                                                                                `json:"-"`
	RankingMethod *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum `json:"rankingMethod"`
	NumTimeSeries *int64                                                                                                              `json:"numTimeSeries"`
	Direction     *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum     `json:"direction"`
}

type jsonDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter
	} else {

		r.RankingMethod = res.RankingMethod

		r.NumTimeSeries = res.NumTimeSeries

		r.Direction = res.Direction

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter = &DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter{empty: true}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsXyChartThresholds struct {
	empty     bool                                                      `json:"-"`
	Label     *string                                                   `json:"label"`
	Value     *float64                                                  `json:"value"`
	Color     *DashboardGridLayoutWidgetsXyChartThresholdsColorEnum     `json:"color"`
	Direction *DashboardGridLayoutWidgetsXyChartThresholdsDirectionEnum `json:"direction"`
}

type jsonDashboardGridLayoutWidgetsXyChartThresholds DashboardGridLayoutWidgetsXyChartThresholds

func (r *DashboardGridLayoutWidgetsXyChartThresholds) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsXyChartThresholds
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsXyChartThresholds
	} else {

		r.Label = res.Label

		r.Value = res.Value

		r.Color = res.Color

		r.Direction = res.Direction

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsXyChartThresholds is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsXyChartThresholds *DashboardGridLayoutWidgetsXyChartThresholds = &DashboardGridLayoutWidgetsXyChartThresholds{empty: true}

func (r *DashboardGridLayoutWidgetsXyChartThresholds) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsXyChartThresholds) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsXyChartThresholds) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsXyChartXAxis struct {
	empty bool                                             `json:"-"`
	Label *string                                          `json:"label"`
	Scale *DashboardGridLayoutWidgetsXyChartXAxisScaleEnum `json:"scale"`
}

type jsonDashboardGridLayoutWidgetsXyChartXAxis DashboardGridLayoutWidgetsXyChartXAxis

func (r *DashboardGridLayoutWidgetsXyChartXAxis) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsXyChartXAxis
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsXyChartXAxis
	} else {

		r.Label = res.Label

		r.Scale = res.Scale

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsXyChartXAxis is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsXyChartXAxis *DashboardGridLayoutWidgetsXyChartXAxis = &DashboardGridLayoutWidgetsXyChartXAxis{empty: true}

func (r *DashboardGridLayoutWidgetsXyChartXAxis) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsXyChartXAxis) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsXyChartXAxis) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsXyChartYAxis struct {
	empty bool                                             `json:"-"`
	Label *string                                          `json:"label"`
	Scale *DashboardGridLayoutWidgetsXyChartYAxisScaleEnum `json:"scale"`
}

type jsonDashboardGridLayoutWidgetsXyChartYAxis DashboardGridLayoutWidgetsXyChartYAxis

func (r *DashboardGridLayoutWidgetsXyChartYAxis) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsXyChartYAxis
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsXyChartYAxis
	} else {

		r.Label = res.Label

		r.Scale = res.Scale

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsXyChartYAxis is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsXyChartYAxis *DashboardGridLayoutWidgetsXyChartYAxis = &DashboardGridLayoutWidgetsXyChartYAxis{empty: true}

func (r *DashboardGridLayoutWidgetsXyChartYAxis) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsXyChartYAxis) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsXyChartYAxis) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsXyChartChartOptions struct {
	empty bool                                                   `json:"-"`
	Mode  *DashboardGridLayoutWidgetsXyChartChartOptionsModeEnum `json:"mode"`
}

type jsonDashboardGridLayoutWidgetsXyChartChartOptions DashboardGridLayoutWidgetsXyChartChartOptions

func (r *DashboardGridLayoutWidgetsXyChartChartOptions) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsXyChartChartOptions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsXyChartChartOptions
	} else {

		r.Mode = res.Mode

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsXyChartChartOptions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsXyChartChartOptions *DashboardGridLayoutWidgetsXyChartChartOptions = &DashboardGridLayoutWidgetsXyChartChartOptions{empty: true}

func (r *DashboardGridLayoutWidgetsXyChartChartOptions) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsXyChartChartOptions) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsXyChartChartOptions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsScorecard struct {
	empty           bool                                                `json:"-"`
	TimeSeriesQuery *DashboardGridLayoutWidgetsScorecardTimeSeriesQuery `json:"timeSeriesQuery"`
	GaugeView       *DashboardGridLayoutWidgetsScorecardGaugeView       `json:"gaugeView"`
	SparkChartView  *DashboardGridLayoutWidgetsScorecardSparkChartView  `json:"sparkChartView"`
	Thresholds      []DashboardGridLayoutWidgetsScorecardThresholds     `json:"thresholds"`
}

type jsonDashboardGridLayoutWidgetsScorecard DashboardGridLayoutWidgetsScorecard

func (r *DashboardGridLayoutWidgetsScorecard) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsScorecard
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsScorecard
	} else {

		r.TimeSeriesQuery = res.TimeSeriesQuery

		r.GaugeView = res.GaugeView

		r.SparkChartView = res.SparkChartView

		r.Thresholds = res.Thresholds

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsScorecard is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsScorecard *DashboardGridLayoutWidgetsScorecard = &DashboardGridLayoutWidgetsScorecard{empty: true}

func (r *DashboardGridLayoutWidgetsScorecard) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsScorecard) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsScorecard) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsScorecardTimeSeriesQuery struct {
	empty                   bool                                                                     `json:"-"`
	TimeSeriesFilter        *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilter      `json:"timeSeriesFilter"`
	TimeSeriesFilterRatio   *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio `json:"timeSeriesFilterRatio"`
	TimeSeriesQueryLanguage *string                                                                  `json:"timeSeriesQueryLanguage"`
	UnitOverride            *string                                                                  `json:"unitOverride"`
}

type jsonDashboardGridLayoutWidgetsScorecardTimeSeriesQuery DashboardGridLayoutWidgetsScorecardTimeSeriesQuery

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQuery) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsScorecardTimeSeriesQuery
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsScorecardTimeSeriesQuery
	} else {

		r.TimeSeriesFilter = res.TimeSeriesFilter

		r.TimeSeriesFilterRatio = res.TimeSeriesFilterRatio

		r.TimeSeriesQueryLanguage = res.TimeSeriesQueryLanguage

		r.UnitOverride = res.UnitOverride

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsScorecardTimeSeriesQuery is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsScorecardTimeSeriesQuery *DashboardGridLayoutWidgetsScorecardTimeSeriesQuery = &DashboardGridLayoutWidgetsScorecardTimeSeriesQuery{empty: true}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQuery) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQuery) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQuery) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilter struct {
	empty                bool                                                                                    `json:"-"`
	Filter               *string                                                                                 `json:"filter"`
	Aggregation          *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation          `json:"aggregation"`
	SecondaryAggregation *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation `json:"secondaryAggregation"`
	PickTimeSeriesFilter *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter `json:"pickTimeSeriesFilter"`
}

type jsonDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilter DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilter

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilter) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilter
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilter
	} else {

		r.Filter = res.Filter

		r.Aggregation = res.Aggregation

		r.SecondaryAggregation = res.SecondaryAggregation

		r.PickTimeSeriesFilter = res.PickTimeSeriesFilter

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilter is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilter *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilter = &DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilter{empty: true}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilter) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilter) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilter) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation struct {
	empty              bool                                                                                                 `json:"-"`
	AlignmentPeriod    *string                                                                                              `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                             `json:"groupByFields"`
}

type jsonDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation = &DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation{empty: true}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation struct {
	empty              bool                                                                                                          `json:"-"`
	AlignmentPeriod    *string                                                                                                       `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                      `json:"groupByFields"`
}

type jsonDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation = &DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation{empty: true}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter struct {
	empty         bool                                                                                                     `json:"-"`
	RankingMethod *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum `json:"rankingMethod"`
	NumTimeSeries *int64                                                                                                   `json:"numTimeSeries"`
	Direction     *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum     `json:"direction"`
}

type jsonDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter
	} else {

		r.RankingMethod = res.RankingMethod

		r.NumTimeSeries = res.NumTimeSeries

		r.Direction = res.Direction

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter = &DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter{empty: true}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio struct {
	empty                bool                                                                                         `json:"-"`
	Numerator            *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator            `json:"numerator"`
	Denominator          *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator          `json:"denominator"`
	SecondaryAggregation *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation `json:"secondaryAggregation"`
	PickTimeSeriesFilter *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter `json:"pickTimeSeriesFilter"`
}

type jsonDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio
	} else {

		r.Numerator = res.Numerator

		r.Denominator = res.Denominator

		r.SecondaryAggregation = res.SecondaryAggregation

		r.PickTimeSeriesFilter = res.PickTimeSeriesFilter

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio = &DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio{empty: true}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator struct {
	empty       bool                                                                                         `json:"-"`
	Filter      *string                                                                                      `json:"filter"`
	Aggregation *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation `json:"aggregation"`
}

type jsonDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator
	} else {

		r.Filter = res.Filter

		r.Aggregation = res.Aggregation

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator = &DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator{empty: true}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation struct {
	empty              bool                                                                                                               `json:"-"`
	AlignmentPeriod    *string                                                                                                            `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                           `json:"groupByFields"`
}

type jsonDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation = &DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation{empty: true}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator struct {
	empty       bool                                                                                           `json:"-"`
	Filter      *string                                                                                        `json:"filter"`
	Aggregation *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation `json:"aggregation"`
}

type jsonDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator
	} else {

		r.Filter = res.Filter

		r.Aggregation = res.Aggregation

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator = &DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator{empty: true}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation struct {
	empty              bool                                                                                                                 `json:"-"`
	AlignmentPeriod    *string                                                                                                              `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                             `json:"groupByFields"`
}

type jsonDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation = &DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation{empty: true}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation struct {
	empty              bool                                                                                                               `json:"-"`
	AlignmentPeriod    *string                                                                                                            `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                           `json:"groupByFields"`
}

type jsonDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation = &DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation{empty: true}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter struct {
	empty         bool                                                                                                          `json:"-"`
	RankingMethod *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum `json:"rankingMethod"`
	NumTimeSeries *int64                                                                                                        `json:"numTimeSeries"`
	Direction     *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum     `json:"direction"`
}

type jsonDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter
	} else {

		r.RankingMethod = res.RankingMethod

		r.NumTimeSeries = res.NumTimeSeries

		r.Direction = res.Direction

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter = &DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter{empty: true}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsScorecardGaugeView struct {
	empty      bool     `json:"-"`
	LowerBound *float64 `json:"lowerBound"`
	UpperBound *float64 `json:"upperBound"`
}

type jsonDashboardGridLayoutWidgetsScorecardGaugeView DashboardGridLayoutWidgetsScorecardGaugeView

func (r *DashboardGridLayoutWidgetsScorecardGaugeView) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsScorecardGaugeView
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsScorecardGaugeView
	} else {

		r.LowerBound = res.LowerBound

		r.UpperBound = res.UpperBound

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsScorecardGaugeView is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsScorecardGaugeView *DashboardGridLayoutWidgetsScorecardGaugeView = &DashboardGridLayoutWidgetsScorecardGaugeView{empty: true}

func (r *DashboardGridLayoutWidgetsScorecardGaugeView) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsScorecardGaugeView) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsScorecardGaugeView) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsScorecardSparkChartView struct {
	empty              bool                                                                 `json:"-"`
	SparkChartType     *DashboardGridLayoutWidgetsScorecardSparkChartViewSparkChartTypeEnum `json:"sparkChartType"`
	MinAlignmentPeriod *string                                                              `json:"minAlignmentPeriod"`
}

type jsonDashboardGridLayoutWidgetsScorecardSparkChartView DashboardGridLayoutWidgetsScorecardSparkChartView

func (r *DashboardGridLayoutWidgetsScorecardSparkChartView) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsScorecardSparkChartView
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsScorecardSparkChartView
	} else {

		r.SparkChartType = res.SparkChartType

		r.MinAlignmentPeriod = res.MinAlignmentPeriod

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsScorecardSparkChartView is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsScorecardSparkChartView *DashboardGridLayoutWidgetsScorecardSparkChartView = &DashboardGridLayoutWidgetsScorecardSparkChartView{empty: true}

func (r *DashboardGridLayoutWidgetsScorecardSparkChartView) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsScorecardSparkChartView) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsScorecardSparkChartView) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsScorecardThresholds struct {
	empty     bool                                                        `json:"-"`
	Label     *string                                                     `json:"label"`
	Value     *float64                                                    `json:"value"`
	Color     *DashboardGridLayoutWidgetsScorecardThresholdsColorEnum     `json:"color"`
	Direction *DashboardGridLayoutWidgetsScorecardThresholdsDirectionEnum `json:"direction"`
}

type jsonDashboardGridLayoutWidgetsScorecardThresholds DashboardGridLayoutWidgetsScorecardThresholds

func (r *DashboardGridLayoutWidgetsScorecardThresholds) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsScorecardThresholds
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsScorecardThresholds
	} else {

		r.Label = res.Label

		r.Value = res.Value

		r.Color = res.Color

		r.Direction = res.Direction

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsScorecardThresholds is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsScorecardThresholds *DashboardGridLayoutWidgetsScorecardThresholds = &DashboardGridLayoutWidgetsScorecardThresholds{empty: true}

func (r *DashboardGridLayoutWidgetsScorecardThresholds) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsScorecardThresholds) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsScorecardThresholds) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsText struct {
	empty   bool                                      `json:"-"`
	Content *string                                   `json:"content"`
	Format  *DashboardGridLayoutWidgetsTextFormatEnum `json:"format"`
}

type jsonDashboardGridLayoutWidgetsText DashboardGridLayoutWidgetsText

func (r *DashboardGridLayoutWidgetsText) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsText
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsText
	} else {

		r.Content = res.Content

		r.Format = res.Format

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsText is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsText *DashboardGridLayoutWidgetsText = &DashboardGridLayoutWidgetsText{empty: true}

func (r *DashboardGridLayoutWidgetsText) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsText) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsText) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsBlank struct {
	empty bool `json:"-"`
}

type jsonDashboardGridLayoutWidgetsBlank DashboardGridLayoutWidgetsBlank

func (r *DashboardGridLayoutWidgetsBlank) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsBlank
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsBlank
	} else {

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsBlank is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsBlank *DashboardGridLayoutWidgetsBlank = &DashboardGridLayoutWidgetsBlank{empty: true}

func (r *DashboardGridLayoutWidgetsBlank) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsBlank) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsBlank) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardGridLayoutWidgetsLogsPanel struct {
	empty         bool     `json:"-"`
	Filter        *string  `json:"filter"`
	ResourceNames []string `json:"resourceNames"`
}

type jsonDashboardGridLayoutWidgetsLogsPanel DashboardGridLayoutWidgetsLogsPanel

func (r *DashboardGridLayoutWidgetsLogsPanel) UnmarshalJSON(data []byte) error {
	var res jsonDashboardGridLayoutWidgetsLogsPanel
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardGridLayoutWidgetsLogsPanel
	} else {

		r.Filter = res.Filter

		r.ResourceNames = res.ResourceNames

	}
	return nil
}

// This object is used to assert a desired state where this DashboardGridLayoutWidgetsLogsPanel is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardGridLayoutWidgetsLogsPanel *DashboardGridLayoutWidgetsLogsPanel = &DashboardGridLayoutWidgetsLogsPanel{empty: true}

func (r *DashboardGridLayoutWidgetsLogsPanel) Empty() bool {
	return r.empty
}

func (r *DashboardGridLayoutWidgetsLogsPanel) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardGridLayoutWidgetsLogsPanel) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayout struct {
	empty   bool                         `json:"-"`
	Columns *int64                       `json:"columns"`
	Tiles   []DashboardMosaicLayoutTiles `json:"tiles"`
}

type jsonDashboardMosaicLayout DashboardMosaicLayout

func (r *DashboardMosaicLayout) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayout
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayout
	} else {

		r.Columns = res.Columns

		r.Tiles = res.Tiles

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayout is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayout *DashboardMosaicLayout = &DashboardMosaicLayout{empty: true}

func (r *DashboardMosaicLayout) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayout) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayout) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTiles struct {
	empty  bool                              `json:"-"`
	XPos   *int64                            `json:"xPos"`
	YPos   *int64                            `json:"yPos"`
	Width  *int64                            `json:"width"`
	Height *int64                            `json:"height"`
	Widget *DashboardMosaicLayoutTilesWidget `json:"widget"`
}

type jsonDashboardMosaicLayoutTiles DashboardMosaicLayoutTiles

func (r *DashboardMosaicLayoutTiles) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTiles
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTiles
	} else {

		r.XPos = res.XPos

		r.YPos = res.YPos

		r.Width = res.Width

		r.Height = res.Height

		r.Widget = res.Widget

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTiles is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTiles *DashboardMosaicLayoutTiles = &DashboardMosaicLayoutTiles{empty: true}

func (r *DashboardMosaicLayoutTiles) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTiles) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTiles) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidget struct {
	empty     bool                                       `json:"-"`
	Title     *string                                    `json:"title"`
	XyChart   *DashboardMosaicLayoutTilesWidgetXyChart   `json:"xyChart"`
	Scorecard *DashboardMosaicLayoutTilesWidgetScorecard `json:"scorecard"`
	Text      *DashboardMosaicLayoutTilesWidgetText      `json:"text"`
	Blank     *DashboardMosaicLayoutTilesWidgetBlank     `json:"blank"`
	LogsPanel *DashboardMosaicLayoutTilesWidgetLogsPanel `json:"logsPanel"`
}

type jsonDashboardMosaicLayoutTilesWidget DashboardMosaicLayoutTilesWidget

func (r *DashboardMosaicLayoutTilesWidget) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidget
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidget
	} else {

		r.Title = res.Title

		r.XyChart = res.XyChart

		r.Scorecard = res.Scorecard

		r.Text = res.Text

		r.Blank = res.Blank

		r.LogsPanel = res.LogsPanel

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidget is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidget *DashboardMosaicLayoutTilesWidget = &DashboardMosaicLayoutTilesWidget{empty: true}

func (r *DashboardMosaicLayoutTilesWidget) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidget) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidget) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetXyChart struct {
	empty             bool                                                 `json:"-"`
	DataSets          []DashboardMosaicLayoutTilesWidgetXyChartDataSets    `json:"dataSets"`
	TimeshiftDuration *string                                              `json:"timeshiftDuration"`
	Thresholds        []DashboardMosaicLayoutTilesWidgetXyChartThresholds  `json:"thresholds"`
	XAxis             *DashboardMosaicLayoutTilesWidgetXyChartXAxis        `json:"xAxis"`
	YAxis             *DashboardMosaicLayoutTilesWidgetXyChartYAxis        `json:"yAxis"`
	ChartOptions      *DashboardMosaicLayoutTilesWidgetXyChartChartOptions `json:"chartOptions"`
}

type jsonDashboardMosaicLayoutTilesWidgetXyChart DashboardMosaicLayoutTilesWidgetXyChart

func (r *DashboardMosaicLayoutTilesWidgetXyChart) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetXyChart
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetXyChart
	} else {

		r.DataSets = res.DataSets

		r.TimeshiftDuration = res.TimeshiftDuration

		r.Thresholds = res.Thresholds

		r.XAxis = res.XAxis

		r.YAxis = res.YAxis

		r.ChartOptions = res.ChartOptions

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetXyChart is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetXyChart *DashboardMosaicLayoutTilesWidgetXyChart = &DashboardMosaicLayoutTilesWidgetXyChart{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetXyChart) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetXyChart) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetXyChart) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetXyChartDataSets struct {
	empty              bool                                                            `json:"-"`
	TimeSeriesQuery    *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQuery `json:"timeSeriesQuery"`
	PlotType           *DashboardMosaicLayoutTilesWidgetXyChartDataSetsPlotTypeEnum    `json:"plotType"`
	LegendTemplate     *string                                                         `json:"legendTemplate"`
	MinAlignmentPeriod *string                                                         `json:"minAlignmentPeriod"`
}

type jsonDashboardMosaicLayoutTilesWidgetXyChartDataSets DashboardMosaicLayoutTilesWidgetXyChartDataSets

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSets) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetXyChartDataSets
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetXyChartDataSets
	} else {

		r.TimeSeriesQuery = res.TimeSeriesQuery

		r.PlotType = res.PlotType

		r.LegendTemplate = res.LegendTemplate

		r.MinAlignmentPeriod = res.MinAlignmentPeriod

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetXyChartDataSets is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetXyChartDataSets *DashboardMosaicLayoutTilesWidgetXyChartDataSets = &DashboardMosaicLayoutTilesWidgetXyChartDataSets{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSets) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSets) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSets) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQuery struct {
	empty                   bool                                                                                 `json:"-"`
	TimeSeriesFilter        *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilter      `json:"timeSeriesFilter"`
	TimeSeriesFilterRatio   *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio `json:"timeSeriesFilterRatio"`
	TimeSeriesQueryLanguage *string                                                                              `json:"timeSeriesQueryLanguage"`
	UnitOverride            *string                                                                              `json:"unitOverride"`
}

type jsonDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQuery DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQuery

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQuery) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQuery
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQuery
	} else {

		r.TimeSeriesFilter = res.TimeSeriesFilter

		r.TimeSeriesFilterRatio = res.TimeSeriesFilterRatio

		r.TimeSeriesQueryLanguage = res.TimeSeriesQueryLanguage

		r.UnitOverride = res.UnitOverride

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQuery is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQuery *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQuery = &DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQuery{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQuery) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQuery) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQuery) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilter struct {
	empty                bool                                                                                                `json:"-"`
	Filter               *string                                                                                             `json:"filter"`
	Aggregation          *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation          `json:"aggregation"`
	SecondaryAggregation *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation `json:"secondaryAggregation"`
	PickTimeSeriesFilter *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter `json:"pickTimeSeriesFilter"`
}

type jsonDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilter DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilter

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilter) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilter
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilter
	} else {

		r.Filter = res.Filter

		r.Aggregation = res.Aggregation

		r.SecondaryAggregation = res.SecondaryAggregation

		r.PickTimeSeriesFilter = res.PickTimeSeriesFilter

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilter is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilter *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilter = &DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilter{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilter) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilter) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilter) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation struct {
	empty              bool                                                                                                             `json:"-"`
	AlignmentPeriod    *string                                                                                                          `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                         `json:"groupByFields"`
}

type jsonDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation = &DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation struct {
	empty              bool                                                                                                                      `json:"-"`
	AlignmentPeriod    *string                                                                                                                   `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                                  `json:"groupByFields"`
}

type jsonDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation = &DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter struct {
	empty         bool                                                                                                                 `json:"-"`
	RankingMethod *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum `json:"rankingMethod"`
	NumTimeSeries *int64                                                                                                               `json:"numTimeSeries"`
	Direction     *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum     `json:"direction"`
}

type jsonDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter
	} else {

		r.RankingMethod = res.RankingMethod

		r.NumTimeSeries = res.NumTimeSeries

		r.Direction = res.Direction

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter = &DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio struct {
	empty                bool                                                                                                     `json:"-"`
	Numerator            *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator            `json:"numerator"`
	Denominator          *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator          `json:"denominator"`
	SecondaryAggregation *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation `json:"secondaryAggregation"`
	PickTimeSeriesFilter *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter `json:"pickTimeSeriesFilter"`
}

type jsonDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio
	} else {

		r.Numerator = res.Numerator

		r.Denominator = res.Denominator

		r.SecondaryAggregation = res.SecondaryAggregation

		r.PickTimeSeriesFilter = res.PickTimeSeriesFilter

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio = &DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator struct {
	empty       bool                                                                                                     `json:"-"`
	Filter      *string                                                                                                  `json:"filter"`
	Aggregation *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation `json:"aggregation"`
}

type jsonDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator
	} else {

		r.Filter = res.Filter

		r.Aggregation = res.Aggregation

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator = &DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation struct {
	empty              bool                                                                                                                           `json:"-"`
	AlignmentPeriod    *string                                                                                                                        `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                                       `json:"groupByFields"`
}

type jsonDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation = &DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator struct {
	empty       bool                                                                                                       `json:"-"`
	Filter      *string                                                                                                    `json:"filter"`
	Aggregation *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation `json:"aggregation"`
}

type jsonDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator
	} else {

		r.Filter = res.Filter

		r.Aggregation = res.Aggregation

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator = &DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation struct {
	empty              bool                                                                                                                             `json:"-"`
	AlignmentPeriod    *string                                                                                                                          `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                                         `json:"groupByFields"`
}

type jsonDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation = &DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation struct {
	empty              bool                                                                                                                           `json:"-"`
	AlignmentPeriod    *string                                                                                                                        `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                                       `json:"groupByFields"`
}

type jsonDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation = &DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter struct {
	empty         bool                                                                                                                      `json:"-"`
	RankingMethod *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum `json:"rankingMethod"`
	NumTimeSeries *int64                                                                                                                    `json:"numTimeSeries"`
	Direction     *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum     `json:"direction"`
}

type jsonDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter
	} else {

		r.RankingMethod = res.RankingMethod

		r.NumTimeSeries = res.NumTimeSeries

		r.Direction = res.Direction

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter = &DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetXyChartThresholds struct {
	empty     bool                                                            `json:"-"`
	Label     *string                                                         `json:"label"`
	Value     *float64                                                        `json:"value"`
	Color     *DashboardMosaicLayoutTilesWidgetXyChartThresholdsColorEnum     `json:"color"`
	Direction *DashboardMosaicLayoutTilesWidgetXyChartThresholdsDirectionEnum `json:"direction"`
}

type jsonDashboardMosaicLayoutTilesWidgetXyChartThresholds DashboardMosaicLayoutTilesWidgetXyChartThresholds

func (r *DashboardMosaicLayoutTilesWidgetXyChartThresholds) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetXyChartThresholds
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetXyChartThresholds
	} else {

		r.Label = res.Label

		r.Value = res.Value

		r.Color = res.Color

		r.Direction = res.Direction

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetXyChartThresholds is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetXyChartThresholds *DashboardMosaicLayoutTilesWidgetXyChartThresholds = &DashboardMosaicLayoutTilesWidgetXyChartThresholds{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetXyChartThresholds) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartThresholds) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartThresholds) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetXyChartXAxis struct {
	empty bool                                                   `json:"-"`
	Label *string                                                `json:"label"`
	Scale *DashboardMosaicLayoutTilesWidgetXyChartXAxisScaleEnum `json:"scale"`
}

type jsonDashboardMosaicLayoutTilesWidgetXyChartXAxis DashboardMosaicLayoutTilesWidgetXyChartXAxis

func (r *DashboardMosaicLayoutTilesWidgetXyChartXAxis) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetXyChartXAxis
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetXyChartXAxis
	} else {

		r.Label = res.Label

		r.Scale = res.Scale

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetXyChartXAxis is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetXyChartXAxis *DashboardMosaicLayoutTilesWidgetXyChartXAxis = &DashboardMosaicLayoutTilesWidgetXyChartXAxis{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetXyChartXAxis) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartXAxis) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartXAxis) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetXyChartYAxis struct {
	empty bool                                                   `json:"-"`
	Label *string                                                `json:"label"`
	Scale *DashboardMosaicLayoutTilesWidgetXyChartYAxisScaleEnum `json:"scale"`
}

type jsonDashboardMosaicLayoutTilesWidgetXyChartYAxis DashboardMosaicLayoutTilesWidgetXyChartYAxis

func (r *DashboardMosaicLayoutTilesWidgetXyChartYAxis) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetXyChartYAxis
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetXyChartYAxis
	} else {

		r.Label = res.Label

		r.Scale = res.Scale

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetXyChartYAxis is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetXyChartYAxis *DashboardMosaicLayoutTilesWidgetXyChartYAxis = &DashboardMosaicLayoutTilesWidgetXyChartYAxis{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetXyChartYAxis) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartYAxis) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartYAxis) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetXyChartChartOptions struct {
	empty bool                                                         `json:"-"`
	Mode  *DashboardMosaicLayoutTilesWidgetXyChartChartOptionsModeEnum `json:"mode"`
}

type jsonDashboardMosaicLayoutTilesWidgetXyChartChartOptions DashboardMosaicLayoutTilesWidgetXyChartChartOptions

func (r *DashboardMosaicLayoutTilesWidgetXyChartChartOptions) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetXyChartChartOptions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetXyChartChartOptions
	} else {

		r.Mode = res.Mode

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetXyChartChartOptions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetXyChartChartOptions *DashboardMosaicLayoutTilesWidgetXyChartChartOptions = &DashboardMosaicLayoutTilesWidgetXyChartChartOptions{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetXyChartChartOptions) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartChartOptions) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetXyChartChartOptions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetScorecard struct {
	empty           bool                                                      `json:"-"`
	TimeSeriesQuery *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQuery `json:"timeSeriesQuery"`
	GaugeView       *DashboardMosaicLayoutTilesWidgetScorecardGaugeView       `json:"gaugeView"`
	SparkChartView  *DashboardMosaicLayoutTilesWidgetScorecardSparkChartView  `json:"sparkChartView"`
	Thresholds      []DashboardMosaicLayoutTilesWidgetScorecardThresholds     `json:"thresholds"`
}

type jsonDashboardMosaicLayoutTilesWidgetScorecard DashboardMosaicLayoutTilesWidgetScorecard

func (r *DashboardMosaicLayoutTilesWidgetScorecard) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetScorecard
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetScorecard
	} else {

		r.TimeSeriesQuery = res.TimeSeriesQuery

		r.GaugeView = res.GaugeView

		r.SparkChartView = res.SparkChartView

		r.Thresholds = res.Thresholds

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetScorecard is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetScorecard *DashboardMosaicLayoutTilesWidgetScorecard = &DashboardMosaicLayoutTilesWidgetScorecard{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetScorecard) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetScorecard) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetScorecard) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQuery struct {
	empty                   bool                                                                           `json:"-"`
	TimeSeriesFilter        *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilter      `json:"timeSeriesFilter"`
	TimeSeriesFilterRatio   *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatio `json:"timeSeriesFilterRatio"`
	TimeSeriesQueryLanguage *string                                                                        `json:"timeSeriesQueryLanguage"`
	UnitOverride            *string                                                                        `json:"unitOverride"`
}

type jsonDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQuery DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQuery

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQuery) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQuery
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQuery
	} else {

		r.TimeSeriesFilter = res.TimeSeriesFilter

		r.TimeSeriesFilterRatio = res.TimeSeriesFilterRatio

		r.TimeSeriesQueryLanguage = res.TimeSeriesQueryLanguage

		r.UnitOverride = res.UnitOverride

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQuery is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQuery *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQuery = &DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQuery{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQuery) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQuery) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQuery) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilter struct {
	empty                bool                                                                                          `json:"-"`
	Filter               *string                                                                                       `json:"filter"`
	Aggregation          *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregation          `json:"aggregation"`
	SecondaryAggregation *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation `json:"secondaryAggregation"`
	PickTimeSeriesFilter *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter `json:"pickTimeSeriesFilter"`
}

type jsonDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilter DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilter

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilter) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilter
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilter
	} else {

		r.Filter = res.Filter

		r.Aggregation = res.Aggregation

		r.SecondaryAggregation = res.SecondaryAggregation

		r.PickTimeSeriesFilter = res.PickTimeSeriesFilter

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilter is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilter *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilter = &DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilter{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilter) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilter) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilter) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregation struct {
	empty              bool                                                                                                       `json:"-"`
	AlignmentPeriod    *string                                                                                                    `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                   `json:"groupByFields"`
}

type jsonDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregation DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregation

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregation *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregation = &DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregation{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation struct {
	empty              bool                                                                                                                `json:"-"`
	AlignmentPeriod    *string                                                                                                             `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                            `json:"groupByFields"`
}

type jsonDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation = &DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter struct {
	empty         bool                                                                                                           `json:"-"`
	RankingMethod *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum `json:"rankingMethod"`
	NumTimeSeries *int64                                                                                                         `json:"numTimeSeries"`
	Direction     *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum     `json:"direction"`
}

type jsonDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter
	} else {

		r.RankingMethod = res.RankingMethod

		r.NumTimeSeries = res.NumTimeSeries

		r.Direction = res.Direction

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter = &DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatio struct {
	empty                bool                                                                                               `json:"-"`
	Numerator            *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator            `json:"numerator"`
	Denominator          *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator          `json:"denominator"`
	SecondaryAggregation *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation `json:"secondaryAggregation"`
	PickTimeSeriesFilter *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter `json:"pickTimeSeriesFilter"`
}

type jsonDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatio DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatio

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatio) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatio
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatio
	} else {

		r.Numerator = res.Numerator

		r.Denominator = res.Denominator

		r.SecondaryAggregation = res.SecondaryAggregation

		r.PickTimeSeriesFilter = res.PickTimeSeriesFilter

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatio is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatio *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatio = &DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatio{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatio) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatio) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatio) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator struct {
	empty       bool                                                                                               `json:"-"`
	Filter      *string                                                                                            `json:"filter"`
	Aggregation *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation `json:"aggregation"`
}

type jsonDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator
	} else {

		r.Filter = res.Filter

		r.Aggregation = res.Aggregation

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator = &DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation struct {
	empty              bool                                                                                                                     `json:"-"`
	AlignmentPeriod    *string                                                                                                                  `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                                 `json:"groupByFields"`
}

type jsonDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation = &DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator struct {
	empty       bool                                                                                                 `json:"-"`
	Filter      *string                                                                                              `json:"filter"`
	Aggregation *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation `json:"aggregation"`
}

type jsonDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator
	} else {

		r.Filter = res.Filter

		r.Aggregation = res.Aggregation

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator = &DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation struct {
	empty              bool                                                                                                                       `json:"-"`
	AlignmentPeriod    *string                                                                                                                    `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                                   `json:"groupByFields"`
}

type jsonDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation = &DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation struct {
	empty              bool                                                                                                                     `json:"-"`
	AlignmentPeriod    *string                                                                                                                  `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                                 `json:"groupByFields"`
}

type jsonDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation = &DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter struct {
	empty         bool                                                                                                                `json:"-"`
	RankingMethod *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum `json:"rankingMethod"`
	NumTimeSeries *int64                                                                                                              `json:"numTimeSeries"`
	Direction     *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum     `json:"direction"`
}

type jsonDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter
	} else {

		r.RankingMethod = res.RankingMethod

		r.NumTimeSeries = res.NumTimeSeries

		r.Direction = res.Direction

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter = &DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetScorecardGaugeView struct {
	empty      bool     `json:"-"`
	LowerBound *float64 `json:"lowerBound"`
	UpperBound *float64 `json:"upperBound"`
}

type jsonDashboardMosaicLayoutTilesWidgetScorecardGaugeView DashboardMosaicLayoutTilesWidgetScorecardGaugeView

func (r *DashboardMosaicLayoutTilesWidgetScorecardGaugeView) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetScorecardGaugeView
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetScorecardGaugeView
	} else {

		r.LowerBound = res.LowerBound

		r.UpperBound = res.UpperBound

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetScorecardGaugeView is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetScorecardGaugeView *DashboardMosaicLayoutTilesWidgetScorecardGaugeView = &DashboardMosaicLayoutTilesWidgetScorecardGaugeView{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetScorecardGaugeView) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetScorecardGaugeView) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetScorecardGaugeView) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetScorecardSparkChartView struct {
	empty              bool                                                                       `json:"-"`
	SparkChartType     *DashboardMosaicLayoutTilesWidgetScorecardSparkChartViewSparkChartTypeEnum `json:"sparkChartType"`
	MinAlignmentPeriod *string                                                                    `json:"minAlignmentPeriod"`
}

type jsonDashboardMosaicLayoutTilesWidgetScorecardSparkChartView DashboardMosaicLayoutTilesWidgetScorecardSparkChartView

func (r *DashboardMosaicLayoutTilesWidgetScorecardSparkChartView) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetScorecardSparkChartView
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetScorecardSparkChartView
	} else {

		r.SparkChartType = res.SparkChartType

		r.MinAlignmentPeriod = res.MinAlignmentPeriod

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetScorecardSparkChartView is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetScorecardSparkChartView *DashboardMosaicLayoutTilesWidgetScorecardSparkChartView = &DashboardMosaicLayoutTilesWidgetScorecardSparkChartView{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetScorecardSparkChartView) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetScorecardSparkChartView) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetScorecardSparkChartView) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetScorecardThresholds struct {
	empty     bool                                                              `json:"-"`
	Label     *string                                                           `json:"label"`
	Value     *float64                                                          `json:"value"`
	Color     *DashboardMosaicLayoutTilesWidgetScorecardThresholdsColorEnum     `json:"color"`
	Direction *DashboardMosaicLayoutTilesWidgetScorecardThresholdsDirectionEnum `json:"direction"`
}

type jsonDashboardMosaicLayoutTilesWidgetScorecardThresholds DashboardMosaicLayoutTilesWidgetScorecardThresholds

func (r *DashboardMosaicLayoutTilesWidgetScorecardThresholds) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetScorecardThresholds
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetScorecardThresholds
	} else {

		r.Label = res.Label

		r.Value = res.Value

		r.Color = res.Color

		r.Direction = res.Direction

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetScorecardThresholds is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetScorecardThresholds *DashboardMosaicLayoutTilesWidgetScorecardThresholds = &DashboardMosaicLayoutTilesWidgetScorecardThresholds{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetScorecardThresholds) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetScorecardThresholds) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetScorecardThresholds) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetText struct {
	empty   bool                                            `json:"-"`
	Content *string                                         `json:"content"`
	Format  *DashboardMosaicLayoutTilesWidgetTextFormatEnum `json:"format"`
}

type jsonDashboardMosaicLayoutTilesWidgetText DashboardMosaicLayoutTilesWidgetText

func (r *DashboardMosaicLayoutTilesWidgetText) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetText
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetText
	} else {

		r.Content = res.Content

		r.Format = res.Format

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetText is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetText *DashboardMosaicLayoutTilesWidgetText = &DashboardMosaicLayoutTilesWidgetText{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetText) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetText) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetText) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetBlank struct {
	empty bool `json:"-"`
}

type jsonDashboardMosaicLayoutTilesWidgetBlank DashboardMosaicLayoutTilesWidgetBlank

func (r *DashboardMosaicLayoutTilesWidgetBlank) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetBlank
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetBlank
	} else {

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetBlank is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetBlank *DashboardMosaicLayoutTilesWidgetBlank = &DashboardMosaicLayoutTilesWidgetBlank{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetBlank) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetBlank) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetBlank) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardMosaicLayoutTilesWidgetLogsPanel struct {
	empty         bool     `json:"-"`
	Filter        *string  `json:"filter"`
	ResourceNames []string `json:"resourceNames"`
}

type jsonDashboardMosaicLayoutTilesWidgetLogsPanel DashboardMosaicLayoutTilesWidgetLogsPanel

func (r *DashboardMosaicLayoutTilesWidgetLogsPanel) UnmarshalJSON(data []byte) error {
	var res jsonDashboardMosaicLayoutTilesWidgetLogsPanel
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardMosaicLayoutTilesWidgetLogsPanel
	} else {

		r.Filter = res.Filter

		r.ResourceNames = res.ResourceNames

	}
	return nil
}

// This object is used to assert a desired state where this DashboardMosaicLayoutTilesWidgetLogsPanel is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardMosaicLayoutTilesWidgetLogsPanel *DashboardMosaicLayoutTilesWidgetLogsPanel = &DashboardMosaicLayoutTilesWidgetLogsPanel{empty: true}

func (r *DashboardMosaicLayoutTilesWidgetLogsPanel) Empty() bool {
	return r.empty
}

func (r *DashboardMosaicLayoutTilesWidgetLogsPanel) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardMosaicLayoutTilesWidgetLogsPanel) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayout struct {
	empty bool                     `json:"-"`
	Rows  []DashboardRowLayoutRows `json:"rows"`
}

type jsonDashboardRowLayout DashboardRowLayout

func (r *DashboardRowLayout) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayout
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayout
	} else {

		r.Rows = res.Rows

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayout is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayout *DashboardRowLayout = &DashboardRowLayout{empty: true}

func (r *DashboardRowLayout) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayout) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayout) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRows struct {
	empty   bool                            `json:"-"`
	Weight  *int64                          `json:"weight"`
	Widgets []DashboardRowLayoutRowsWidgets `json:"widgets"`
}

type jsonDashboardRowLayoutRows DashboardRowLayoutRows

func (r *DashboardRowLayoutRows) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRows
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRows
	} else {

		r.Weight = res.Weight

		r.Widgets = res.Widgets

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRows is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRows *DashboardRowLayoutRows = &DashboardRowLayoutRows{empty: true}

func (r *DashboardRowLayoutRows) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRows) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRows) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgets struct {
	empty     bool                                    `json:"-"`
	Title     *string                                 `json:"title"`
	XyChart   *DashboardRowLayoutRowsWidgetsXyChart   `json:"xyChart"`
	Scorecard *DashboardRowLayoutRowsWidgetsScorecard `json:"scorecard"`
	Text      *DashboardRowLayoutRowsWidgetsText      `json:"text"`
	Blank     *DashboardRowLayoutRowsWidgetsBlank     `json:"blank"`
	LogsPanel *DashboardRowLayoutRowsWidgetsLogsPanel `json:"logsPanel"`
}

type jsonDashboardRowLayoutRowsWidgets DashboardRowLayoutRowsWidgets

func (r *DashboardRowLayoutRowsWidgets) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgets
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgets
	} else {

		r.Title = res.Title

		r.XyChart = res.XyChart

		r.Scorecard = res.Scorecard

		r.Text = res.Text

		r.Blank = res.Blank

		r.LogsPanel = res.LogsPanel

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgets is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgets *DashboardRowLayoutRowsWidgets = &DashboardRowLayoutRowsWidgets{empty: true}

func (r *DashboardRowLayoutRowsWidgets) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgets) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgets) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsXyChart struct {
	empty             bool                                              `json:"-"`
	DataSets          []DashboardRowLayoutRowsWidgetsXyChartDataSets    `json:"dataSets"`
	TimeshiftDuration *string                                           `json:"timeshiftDuration"`
	Thresholds        []DashboardRowLayoutRowsWidgetsXyChartThresholds  `json:"thresholds"`
	XAxis             *DashboardRowLayoutRowsWidgetsXyChartXAxis        `json:"xAxis"`
	YAxis             *DashboardRowLayoutRowsWidgetsXyChartYAxis        `json:"yAxis"`
	ChartOptions      *DashboardRowLayoutRowsWidgetsXyChartChartOptions `json:"chartOptions"`
}

type jsonDashboardRowLayoutRowsWidgetsXyChart DashboardRowLayoutRowsWidgetsXyChart

func (r *DashboardRowLayoutRowsWidgetsXyChart) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsXyChart
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsXyChart
	} else {

		r.DataSets = res.DataSets

		r.TimeshiftDuration = res.TimeshiftDuration

		r.Thresholds = res.Thresholds

		r.XAxis = res.XAxis

		r.YAxis = res.YAxis

		r.ChartOptions = res.ChartOptions

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsXyChart is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsXyChart *DashboardRowLayoutRowsWidgetsXyChart = &DashboardRowLayoutRowsWidgetsXyChart{empty: true}

func (r *DashboardRowLayoutRowsWidgetsXyChart) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsXyChart) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsXyChart) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsXyChartDataSets struct {
	empty              bool                                                         `json:"-"`
	TimeSeriesQuery    *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQuery `json:"timeSeriesQuery"`
	PlotType           *DashboardRowLayoutRowsWidgetsXyChartDataSetsPlotTypeEnum    `json:"plotType"`
	LegendTemplate     *string                                                      `json:"legendTemplate"`
	MinAlignmentPeriod *string                                                      `json:"minAlignmentPeriod"`
}

type jsonDashboardRowLayoutRowsWidgetsXyChartDataSets DashboardRowLayoutRowsWidgetsXyChartDataSets

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSets) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsXyChartDataSets
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsXyChartDataSets
	} else {

		r.TimeSeriesQuery = res.TimeSeriesQuery

		r.PlotType = res.PlotType

		r.LegendTemplate = res.LegendTemplate

		r.MinAlignmentPeriod = res.MinAlignmentPeriod

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsXyChartDataSets is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsXyChartDataSets *DashboardRowLayoutRowsWidgetsXyChartDataSets = &DashboardRowLayoutRowsWidgetsXyChartDataSets{empty: true}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSets) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSets) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSets) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQuery struct {
	empty                   bool                                                                              `json:"-"`
	TimeSeriesFilter        *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter      `json:"timeSeriesFilter"`
	TimeSeriesFilterRatio   *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio `json:"timeSeriesFilterRatio"`
	TimeSeriesQueryLanguage *string                                                                           `json:"timeSeriesQueryLanguage"`
	UnitOverride            *string                                                                           `json:"unitOverride"`
}

type jsonDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQuery DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQuery

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQuery) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQuery
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQuery
	} else {

		r.TimeSeriesFilter = res.TimeSeriesFilter

		r.TimeSeriesFilterRatio = res.TimeSeriesFilterRatio

		r.TimeSeriesQueryLanguage = res.TimeSeriesQueryLanguage

		r.UnitOverride = res.UnitOverride

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQuery is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQuery *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQuery = &DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQuery{empty: true}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQuery) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQuery) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQuery) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter struct {
	empty                bool                                                                                             `json:"-"`
	Filter               *string                                                                                          `json:"filter"`
	Aggregation          *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation          `json:"aggregation"`
	SecondaryAggregation *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation `json:"secondaryAggregation"`
	PickTimeSeriesFilter *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter `json:"pickTimeSeriesFilter"`
}

type jsonDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter
	} else {

		r.Filter = res.Filter

		r.Aggregation = res.Aggregation

		r.SecondaryAggregation = res.SecondaryAggregation

		r.PickTimeSeriesFilter = res.PickTimeSeriesFilter

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter = &DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter{empty: true}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation struct {
	empty              bool                                                                                                          `json:"-"`
	AlignmentPeriod    *string                                                                                                       `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                      `json:"groupByFields"`
}

type jsonDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation = &DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation{empty: true}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation struct {
	empty              bool                                                                                                                   `json:"-"`
	AlignmentPeriod    *string                                                                                                                `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                               `json:"groupByFields"`
}

type jsonDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation = &DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation{empty: true}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter struct {
	empty         bool                                                                                                              `json:"-"`
	RankingMethod *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum `json:"rankingMethod"`
	NumTimeSeries *int64                                                                                                            `json:"numTimeSeries"`
	Direction     *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum     `json:"direction"`
}

type jsonDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter
	} else {

		r.RankingMethod = res.RankingMethod

		r.NumTimeSeries = res.NumTimeSeries

		r.Direction = res.Direction

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter = &DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter{empty: true}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio struct {
	empty                bool                                                                                                  `json:"-"`
	Numerator            *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator            `json:"numerator"`
	Denominator          *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator          `json:"denominator"`
	SecondaryAggregation *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation `json:"secondaryAggregation"`
	PickTimeSeriesFilter *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter `json:"pickTimeSeriesFilter"`
}

type jsonDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio
	} else {

		r.Numerator = res.Numerator

		r.Denominator = res.Denominator

		r.SecondaryAggregation = res.SecondaryAggregation

		r.PickTimeSeriesFilter = res.PickTimeSeriesFilter

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio = &DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio{empty: true}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator struct {
	empty       bool                                                                                                  `json:"-"`
	Filter      *string                                                                                               `json:"filter"`
	Aggregation *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation `json:"aggregation"`
}

type jsonDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator
	} else {

		r.Filter = res.Filter

		r.Aggregation = res.Aggregation

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator = &DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator{empty: true}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation struct {
	empty              bool                                                                                                                        `json:"-"`
	AlignmentPeriod    *string                                                                                                                     `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                                    `json:"groupByFields"`
}

type jsonDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation = &DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation{empty: true}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator struct {
	empty       bool                                                                                                    `json:"-"`
	Filter      *string                                                                                                 `json:"filter"`
	Aggregation *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation `json:"aggregation"`
}

type jsonDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator
	} else {

		r.Filter = res.Filter

		r.Aggregation = res.Aggregation

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator = &DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator{empty: true}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation struct {
	empty              bool                                                                                                                          `json:"-"`
	AlignmentPeriod    *string                                                                                                                       `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                                      `json:"groupByFields"`
}

type jsonDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation = &DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation{empty: true}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation struct {
	empty              bool                                                                                                                        `json:"-"`
	AlignmentPeriod    *string                                                                                                                     `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                                    `json:"groupByFields"`
}

type jsonDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation = &DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation{empty: true}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter struct {
	empty         bool                                                                                                                   `json:"-"`
	RankingMethod *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum `json:"rankingMethod"`
	NumTimeSeries *int64                                                                                                                 `json:"numTimeSeries"`
	Direction     *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum     `json:"direction"`
}

type jsonDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter
	} else {

		r.RankingMethod = res.RankingMethod

		r.NumTimeSeries = res.NumTimeSeries

		r.Direction = res.Direction

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter = &DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter{empty: true}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsXyChartThresholds struct {
	empty     bool                                                         `json:"-"`
	Label     *string                                                      `json:"label"`
	Value     *float64                                                     `json:"value"`
	Color     *DashboardRowLayoutRowsWidgetsXyChartThresholdsColorEnum     `json:"color"`
	Direction *DashboardRowLayoutRowsWidgetsXyChartThresholdsDirectionEnum `json:"direction"`
}

type jsonDashboardRowLayoutRowsWidgetsXyChartThresholds DashboardRowLayoutRowsWidgetsXyChartThresholds

func (r *DashboardRowLayoutRowsWidgetsXyChartThresholds) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsXyChartThresholds
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsXyChartThresholds
	} else {

		r.Label = res.Label

		r.Value = res.Value

		r.Color = res.Color

		r.Direction = res.Direction

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsXyChartThresholds is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsXyChartThresholds *DashboardRowLayoutRowsWidgetsXyChartThresholds = &DashboardRowLayoutRowsWidgetsXyChartThresholds{empty: true}

func (r *DashboardRowLayoutRowsWidgetsXyChartThresholds) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsXyChartThresholds) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsXyChartThresholds) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsXyChartXAxis struct {
	empty bool                                                `json:"-"`
	Label *string                                             `json:"label"`
	Scale *DashboardRowLayoutRowsWidgetsXyChartXAxisScaleEnum `json:"scale"`
}

type jsonDashboardRowLayoutRowsWidgetsXyChartXAxis DashboardRowLayoutRowsWidgetsXyChartXAxis

func (r *DashboardRowLayoutRowsWidgetsXyChartXAxis) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsXyChartXAxis
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsXyChartXAxis
	} else {

		r.Label = res.Label

		r.Scale = res.Scale

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsXyChartXAxis is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsXyChartXAxis *DashboardRowLayoutRowsWidgetsXyChartXAxis = &DashboardRowLayoutRowsWidgetsXyChartXAxis{empty: true}

func (r *DashboardRowLayoutRowsWidgetsXyChartXAxis) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsXyChartXAxis) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsXyChartXAxis) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsXyChartYAxis struct {
	empty bool                                                `json:"-"`
	Label *string                                             `json:"label"`
	Scale *DashboardRowLayoutRowsWidgetsXyChartYAxisScaleEnum `json:"scale"`
}

type jsonDashboardRowLayoutRowsWidgetsXyChartYAxis DashboardRowLayoutRowsWidgetsXyChartYAxis

func (r *DashboardRowLayoutRowsWidgetsXyChartYAxis) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsXyChartYAxis
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsXyChartYAxis
	} else {

		r.Label = res.Label

		r.Scale = res.Scale

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsXyChartYAxis is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsXyChartYAxis *DashboardRowLayoutRowsWidgetsXyChartYAxis = &DashboardRowLayoutRowsWidgetsXyChartYAxis{empty: true}

func (r *DashboardRowLayoutRowsWidgetsXyChartYAxis) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsXyChartYAxis) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsXyChartYAxis) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsXyChartChartOptions struct {
	empty bool                                                      `json:"-"`
	Mode  *DashboardRowLayoutRowsWidgetsXyChartChartOptionsModeEnum `json:"mode"`
}

type jsonDashboardRowLayoutRowsWidgetsXyChartChartOptions DashboardRowLayoutRowsWidgetsXyChartChartOptions

func (r *DashboardRowLayoutRowsWidgetsXyChartChartOptions) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsXyChartChartOptions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsXyChartChartOptions
	} else {

		r.Mode = res.Mode

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsXyChartChartOptions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsXyChartChartOptions *DashboardRowLayoutRowsWidgetsXyChartChartOptions = &DashboardRowLayoutRowsWidgetsXyChartChartOptions{empty: true}

func (r *DashboardRowLayoutRowsWidgetsXyChartChartOptions) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsXyChartChartOptions) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsXyChartChartOptions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsScorecard struct {
	empty           bool                                                   `json:"-"`
	TimeSeriesQuery *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQuery `json:"timeSeriesQuery"`
	GaugeView       *DashboardRowLayoutRowsWidgetsScorecardGaugeView       `json:"gaugeView"`
	SparkChartView  *DashboardRowLayoutRowsWidgetsScorecardSparkChartView  `json:"sparkChartView"`
	Thresholds      []DashboardRowLayoutRowsWidgetsScorecardThresholds     `json:"thresholds"`
}

type jsonDashboardRowLayoutRowsWidgetsScorecard DashboardRowLayoutRowsWidgetsScorecard

func (r *DashboardRowLayoutRowsWidgetsScorecard) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsScorecard
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsScorecard
	} else {

		r.TimeSeriesQuery = res.TimeSeriesQuery

		r.GaugeView = res.GaugeView

		r.SparkChartView = res.SparkChartView

		r.Thresholds = res.Thresholds

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsScorecard is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsScorecard *DashboardRowLayoutRowsWidgetsScorecard = &DashboardRowLayoutRowsWidgetsScorecard{empty: true}

func (r *DashboardRowLayoutRowsWidgetsScorecard) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsScorecard) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsScorecard) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQuery struct {
	empty                   bool                                                                        `json:"-"`
	TimeSeriesFilter        *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilter      `json:"timeSeriesFilter"`
	TimeSeriesFilterRatio   *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio `json:"timeSeriesFilterRatio"`
	TimeSeriesQueryLanguage *string                                                                     `json:"timeSeriesQueryLanguage"`
	UnitOverride            *string                                                                     `json:"unitOverride"`
}

type jsonDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQuery DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQuery

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQuery) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQuery
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQuery
	} else {

		r.TimeSeriesFilter = res.TimeSeriesFilter

		r.TimeSeriesFilterRatio = res.TimeSeriesFilterRatio

		r.TimeSeriesQueryLanguage = res.TimeSeriesQueryLanguage

		r.UnitOverride = res.UnitOverride

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQuery is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQuery *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQuery = &DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQuery{empty: true}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQuery) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQuery) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQuery) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilter struct {
	empty                bool                                                                                       `json:"-"`
	Filter               *string                                                                                    `json:"filter"`
	Aggregation          *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation          `json:"aggregation"`
	SecondaryAggregation *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation `json:"secondaryAggregation"`
	PickTimeSeriesFilter *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter `json:"pickTimeSeriesFilter"`
}

type jsonDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilter DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilter

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilter) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilter
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilter
	} else {

		r.Filter = res.Filter

		r.Aggregation = res.Aggregation

		r.SecondaryAggregation = res.SecondaryAggregation

		r.PickTimeSeriesFilter = res.PickTimeSeriesFilter

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilter is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilter *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilter = &DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilter{empty: true}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilter) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilter) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilter) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation struct {
	empty              bool                                                                                                    `json:"-"`
	AlignmentPeriod    *string                                                                                                 `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                `json:"groupByFields"`
}

type jsonDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation = &DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation{empty: true}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation struct {
	empty              bool                                                                                                             `json:"-"`
	AlignmentPeriod    *string                                                                                                          `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                         `json:"groupByFields"`
}

type jsonDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation = &DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation{empty: true}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter struct {
	empty         bool                                                                                                        `json:"-"`
	RankingMethod *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum `json:"rankingMethod"`
	NumTimeSeries *int64                                                                                                      `json:"numTimeSeries"`
	Direction     *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum     `json:"direction"`
}

type jsonDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter
	} else {

		r.RankingMethod = res.RankingMethod

		r.NumTimeSeries = res.NumTimeSeries

		r.Direction = res.Direction

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter = &DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter{empty: true}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio struct {
	empty                bool                                                                                            `json:"-"`
	Numerator            *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator            `json:"numerator"`
	Denominator          *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator          `json:"denominator"`
	SecondaryAggregation *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation `json:"secondaryAggregation"`
	PickTimeSeriesFilter *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter `json:"pickTimeSeriesFilter"`
}

type jsonDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio
	} else {

		r.Numerator = res.Numerator

		r.Denominator = res.Denominator

		r.SecondaryAggregation = res.SecondaryAggregation

		r.PickTimeSeriesFilter = res.PickTimeSeriesFilter

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio = &DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio{empty: true}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator struct {
	empty       bool                                                                                            `json:"-"`
	Filter      *string                                                                                         `json:"filter"`
	Aggregation *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation `json:"aggregation"`
}

type jsonDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator
	} else {

		r.Filter = res.Filter

		r.Aggregation = res.Aggregation

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator = &DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator{empty: true}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation struct {
	empty              bool                                                                                                                  `json:"-"`
	AlignmentPeriod    *string                                                                                                               `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                              `json:"groupByFields"`
}

type jsonDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation = &DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation{empty: true}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator struct {
	empty       bool                                                                                              `json:"-"`
	Filter      *string                                                                                           `json:"filter"`
	Aggregation *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation `json:"aggregation"`
}

type jsonDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator
	} else {

		r.Filter = res.Filter

		r.Aggregation = res.Aggregation

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator = &DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator{empty: true}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation struct {
	empty              bool                                                                                                                    `json:"-"`
	AlignmentPeriod    *string                                                                                                                 `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                                `json:"groupByFields"`
}

type jsonDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation = &DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation{empty: true}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation struct {
	empty              bool                                                                                                                  `json:"-"`
	AlignmentPeriod    *string                                                                                                               `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                              `json:"groupByFields"`
}

type jsonDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation = &DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation{empty: true}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter struct {
	empty         bool                                                                                                             `json:"-"`
	RankingMethod *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum `json:"rankingMethod"`
	NumTimeSeries *int64                                                                                                           `json:"numTimeSeries"`
	Direction     *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum     `json:"direction"`
}

type jsonDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter
	} else {

		r.RankingMethod = res.RankingMethod

		r.NumTimeSeries = res.NumTimeSeries

		r.Direction = res.Direction

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter = &DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter{empty: true}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsScorecardGaugeView struct {
	empty      bool     `json:"-"`
	LowerBound *float64 `json:"lowerBound"`
	UpperBound *float64 `json:"upperBound"`
}

type jsonDashboardRowLayoutRowsWidgetsScorecardGaugeView DashboardRowLayoutRowsWidgetsScorecardGaugeView

func (r *DashboardRowLayoutRowsWidgetsScorecardGaugeView) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsScorecardGaugeView
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsScorecardGaugeView
	} else {

		r.LowerBound = res.LowerBound

		r.UpperBound = res.UpperBound

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsScorecardGaugeView is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsScorecardGaugeView *DashboardRowLayoutRowsWidgetsScorecardGaugeView = &DashboardRowLayoutRowsWidgetsScorecardGaugeView{empty: true}

func (r *DashboardRowLayoutRowsWidgetsScorecardGaugeView) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsScorecardGaugeView) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsScorecardGaugeView) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsScorecardSparkChartView struct {
	empty              bool                                                                    `json:"-"`
	SparkChartType     *DashboardRowLayoutRowsWidgetsScorecardSparkChartViewSparkChartTypeEnum `json:"sparkChartType"`
	MinAlignmentPeriod *string                                                                 `json:"minAlignmentPeriod"`
}

type jsonDashboardRowLayoutRowsWidgetsScorecardSparkChartView DashboardRowLayoutRowsWidgetsScorecardSparkChartView

func (r *DashboardRowLayoutRowsWidgetsScorecardSparkChartView) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsScorecardSparkChartView
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsScorecardSparkChartView
	} else {

		r.SparkChartType = res.SparkChartType

		r.MinAlignmentPeriod = res.MinAlignmentPeriod

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsScorecardSparkChartView is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsScorecardSparkChartView *DashboardRowLayoutRowsWidgetsScorecardSparkChartView = &DashboardRowLayoutRowsWidgetsScorecardSparkChartView{empty: true}

func (r *DashboardRowLayoutRowsWidgetsScorecardSparkChartView) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsScorecardSparkChartView) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsScorecardSparkChartView) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsScorecardThresholds struct {
	empty     bool                                                           `json:"-"`
	Label     *string                                                        `json:"label"`
	Value     *float64                                                       `json:"value"`
	Color     *DashboardRowLayoutRowsWidgetsScorecardThresholdsColorEnum     `json:"color"`
	Direction *DashboardRowLayoutRowsWidgetsScorecardThresholdsDirectionEnum `json:"direction"`
}

type jsonDashboardRowLayoutRowsWidgetsScorecardThresholds DashboardRowLayoutRowsWidgetsScorecardThresholds

func (r *DashboardRowLayoutRowsWidgetsScorecardThresholds) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsScorecardThresholds
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsScorecardThresholds
	} else {

		r.Label = res.Label

		r.Value = res.Value

		r.Color = res.Color

		r.Direction = res.Direction

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsScorecardThresholds is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsScorecardThresholds *DashboardRowLayoutRowsWidgetsScorecardThresholds = &DashboardRowLayoutRowsWidgetsScorecardThresholds{empty: true}

func (r *DashboardRowLayoutRowsWidgetsScorecardThresholds) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsScorecardThresholds) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsScorecardThresholds) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsText struct {
	empty   bool                                         `json:"-"`
	Content *string                                      `json:"content"`
	Format  *DashboardRowLayoutRowsWidgetsTextFormatEnum `json:"format"`
}

type jsonDashboardRowLayoutRowsWidgetsText DashboardRowLayoutRowsWidgetsText

func (r *DashboardRowLayoutRowsWidgetsText) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsText
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsText
	} else {

		r.Content = res.Content

		r.Format = res.Format

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsText is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsText *DashboardRowLayoutRowsWidgetsText = &DashboardRowLayoutRowsWidgetsText{empty: true}

func (r *DashboardRowLayoutRowsWidgetsText) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsText) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsText) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsBlank struct {
	empty bool `json:"-"`
}

type jsonDashboardRowLayoutRowsWidgetsBlank DashboardRowLayoutRowsWidgetsBlank

func (r *DashboardRowLayoutRowsWidgetsBlank) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsBlank
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsBlank
	} else {

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsBlank is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsBlank *DashboardRowLayoutRowsWidgetsBlank = &DashboardRowLayoutRowsWidgetsBlank{empty: true}

func (r *DashboardRowLayoutRowsWidgetsBlank) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsBlank) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsBlank) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardRowLayoutRowsWidgetsLogsPanel struct {
	empty         bool     `json:"-"`
	Filter        *string  `json:"filter"`
	ResourceNames []string `json:"resourceNames"`
}

type jsonDashboardRowLayoutRowsWidgetsLogsPanel DashboardRowLayoutRowsWidgetsLogsPanel

func (r *DashboardRowLayoutRowsWidgetsLogsPanel) UnmarshalJSON(data []byte) error {
	var res jsonDashboardRowLayoutRowsWidgetsLogsPanel
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardRowLayoutRowsWidgetsLogsPanel
	} else {

		r.Filter = res.Filter

		r.ResourceNames = res.ResourceNames

	}
	return nil
}

// This object is used to assert a desired state where this DashboardRowLayoutRowsWidgetsLogsPanel is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardRowLayoutRowsWidgetsLogsPanel *DashboardRowLayoutRowsWidgetsLogsPanel = &DashboardRowLayoutRowsWidgetsLogsPanel{empty: true}

func (r *DashboardRowLayoutRowsWidgetsLogsPanel) Empty() bool {
	return r.empty
}

func (r *DashboardRowLayoutRowsWidgetsLogsPanel) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardRowLayoutRowsWidgetsLogsPanel) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayout struct {
	empty   bool                           `json:"-"`
	Columns []DashboardColumnLayoutColumns `json:"columns"`
}

type jsonDashboardColumnLayout DashboardColumnLayout

func (r *DashboardColumnLayout) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayout
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayout
	} else {

		r.Columns = res.Columns

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayout is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayout *DashboardColumnLayout = &DashboardColumnLayout{empty: true}

func (r *DashboardColumnLayout) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayout) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayout) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumns struct {
	empty   bool                                  `json:"-"`
	Weight  *int64                                `json:"weight"`
	Widgets []DashboardColumnLayoutColumnsWidgets `json:"widgets"`
}

type jsonDashboardColumnLayoutColumns DashboardColumnLayoutColumns

func (r *DashboardColumnLayoutColumns) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumns
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumns
	} else {

		r.Weight = res.Weight

		r.Widgets = res.Widgets

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumns is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumns *DashboardColumnLayoutColumns = &DashboardColumnLayoutColumns{empty: true}

func (r *DashboardColumnLayoutColumns) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumns) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumns) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgets struct {
	empty     bool                                          `json:"-"`
	Title     *string                                       `json:"title"`
	XyChart   *DashboardColumnLayoutColumnsWidgetsXyChart   `json:"xyChart"`
	Scorecard *DashboardColumnLayoutColumnsWidgetsScorecard `json:"scorecard"`
	Text      *DashboardColumnLayoutColumnsWidgetsText      `json:"text"`
	Blank     *DashboardColumnLayoutColumnsWidgetsBlank     `json:"blank"`
	LogsPanel *DashboardColumnLayoutColumnsWidgetsLogsPanel `json:"logsPanel"`
}

type jsonDashboardColumnLayoutColumnsWidgets DashboardColumnLayoutColumnsWidgets

func (r *DashboardColumnLayoutColumnsWidgets) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgets
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgets
	} else {

		r.Title = res.Title

		r.XyChart = res.XyChart

		r.Scorecard = res.Scorecard

		r.Text = res.Text

		r.Blank = res.Blank

		r.LogsPanel = res.LogsPanel

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgets is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgets *DashboardColumnLayoutColumnsWidgets = &DashboardColumnLayoutColumnsWidgets{empty: true}

func (r *DashboardColumnLayoutColumnsWidgets) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgets) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgets) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsXyChart struct {
	empty             bool                                                    `json:"-"`
	DataSets          []DashboardColumnLayoutColumnsWidgetsXyChartDataSets    `json:"dataSets"`
	TimeshiftDuration *string                                                 `json:"timeshiftDuration"`
	Thresholds        []DashboardColumnLayoutColumnsWidgetsXyChartThresholds  `json:"thresholds"`
	XAxis             *DashboardColumnLayoutColumnsWidgetsXyChartXAxis        `json:"xAxis"`
	YAxis             *DashboardColumnLayoutColumnsWidgetsXyChartYAxis        `json:"yAxis"`
	ChartOptions      *DashboardColumnLayoutColumnsWidgetsXyChartChartOptions `json:"chartOptions"`
}

type jsonDashboardColumnLayoutColumnsWidgetsXyChart DashboardColumnLayoutColumnsWidgetsXyChart

func (r *DashboardColumnLayoutColumnsWidgetsXyChart) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsXyChart
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsXyChart
	} else {

		r.DataSets = res.DataSets

		r.TimeshiftDuration = res.TimeshiftDuration

		r.Thresholds = res.Thresholds

		r.XAxis = res.XAxis

		r.YAxis = res.YAxis

		r.ChartOptions = res.ChartOptions

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsXyChart is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsXyChart *DashboardColumnLayoutColumnsWidgetsXyChart = &DashboardColumnLayoutColumnsWidgetsXyChart{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsXyChart) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChart) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChart) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsXyChartDataSets struct {
	empty              bool                                                               `json:"-"`
	TimeSeriesQuery    *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQuery `json:"timeSeriesQuery"`
	PlotType           *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsPlotTypeEnum    `json:"plotType"`
	LegendTemplate     *string                                                            `json:"legendTemplate"`
	MinAlignmentPeriod *string                                                            `json:"minAlignmentPeriod"`
}

type jsonDashboardColumnLayoutColumnsWidgetsXyChartDataSets DashboardColumnLayoutColumnsWidgetsXyChartDataSets

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSets) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsXyChartDataSets
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsXyChartDataSets
	} else {

		r.TimeSeriesQuery = res.TimeSeriesQuery

		r.PlotType = res.PlotType

		r.LegendTemplate = res.LegendTemplate

		r.MinAlignmentPeriod = res.MinAlignmentPeriod

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsXyChartDataSets is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsXyChartDataSets *DashboardColumnLayoutColumnsWidgetsXyChartDataSets = &DashboardColumnLayoutColumnsWidgetsXyChartDataSets{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSets) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSets) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSets) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQuery struct {
	empty                   bool                                                                                    `json:"-"`
	TimeSeriesFilter        *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter      `json:"timeSeriesFilter"`
	TimeSeriesFilterRatio   *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio `json:"timeSeriesFilterRatio"`
	TimeSeriesQueryLanguage *string                                                                                 `json:"timeSeriesQueryLanguage"`
	UnitOverride            *string                                                                                 `json:"unitOverride"`
}

type jsonDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQuery DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQuery

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQuery) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQuery
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQuery
	} else {

		r.TimeSeriesFilter = res.TimeSeriesFilter

		r.TimeSeriesFilterRatio = res.TimeSeriesFilterRatio

		r.TimeSeriesQueryLanguage = res.TimeSeriesQueryLanguage

		r.UnitOverride = res.UnitOverride

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQuery is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQuery *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQuery = &DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQuery{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQuery) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQuery) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQuery) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter struct {
	empty                bool                                                                                                   `json:"-"`
	Filter               *string                                                                                                `json:"filter"`
	Aggregation          *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation          `json:"aggregation"`
	SecondaryAggregation *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation `json:"secondaryAggregation"`
	PickTimeSeriesFilter *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter `json:"pickTimeSeriesFilter"`
}

type jsonDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter
	} else {

		r.Filter = res.Filter

		r.Aggregation = res.Aggregation

		r.SecondaryAggregation = res.SecondaryAggregation

		r.PickTimeSeriesFilter = res.PickTimeSeriesFilter

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter = &DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilter) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation struct {
	empty              bool                                                                                                                `json:"-"`
	AlignmentPeriod    *string                                                                                                             `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                            `json:"groupByFields"`
}

type jsonDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation = &DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation struct {
	empty              bool                                                                                                                         `json:"-"`
	AlignmentPeriod    *string                                                                                                                      `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                                     `json:"groupByFields"`
}

type jsonDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation = &DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterSecondaryAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter struct {
	empty         bool                                                                                                                    `json:"-"`
	RankingMethod *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum `json:"rankingMethod"`
	NumTimeSeries *int64                                                                                                                  `json:"numTimeSeries"`
	Direction     *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum     `json:"direction"`
}

type jsonDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter
	} else {

		r.RankingMethod = res.RankingMethod

		r.NumTimeSeries = res.NumTimeSeries

		r.Direction = res.Direction

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter = &DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio struct {
	empty                bool                                                                                                        `json:"-"`
	Numerator            *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator            `json:"numerator"`
	Denominator          *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator          `json:"denominator"`
	SecondaryAggregation *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation `json:"secondaryAggregation"`
	PickTimeSeriesFilter *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter `json:"pickTimeSeriesFilter"`
}

type jsonDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio
	} else {

		r.Numerator = res.Numerator

		r.Denominator = res.Denominator

		r.SecondaryAggregation = res.SecondaryAggregation

		r.PickTimeSeriesFilter = res.PickTimeSeriesFilter

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio = &DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatio) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator struct {
	empty       bool                                                                                                        `json:"-"`
	Filter      *string                                                                                                     `json:"filter"`
	Aggregation *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation `json:"aggregation"`
}

type jsonDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator
	} else {

		r.Filter = res.Filter

		r.Aggregation = res.Aggregation

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator = &DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumerator) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation struct {
	empty              bool                                                                                                                              `json:"-"`
	AlignmentPeriod    *string                                                                                                                           `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                                          `json:"groupByFields"`
}

type jsonDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation = &DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator struct {
	empty       bool                                                                                                          `json:"-"`
	Filter      *string                                                                                                       `json:"filter"`
	Aggregation *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation `json:"aggregation"`
}

type jsonDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator
	} else {

		r.Filter = res.Filter

		r.Aggregation = res.Aggregation

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator = &DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominator) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation struct {
	empty              bool                                                                                                                                `json:"-"`
	AlignmentPeriod    *string                                                                                                                             `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                                            `json:"groupByFields"`
}

type jsonDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation = &DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation struct {
	empty              bool                                                                                                                              `json:"-"`
	AlignmentPeriod    *string                                                                                                                           `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                                          `json:"groupByFields"`
}

type jsonDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation = &DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter struct {
	empty         bool                                                                                                                         `json:"-"`
	RankingMethod *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum `json:"rankingMethod"`
	NumTimeSeries *int64                                                                                                                       `json:"numTimeSeries"`
	Direction     *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum     `json:"direction"`
}

type jsonDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter
	} else {

		r.RankingMethod = res.RankingMethod

		r.NumTimeSeries = res.NumTimeSeries

		r.Direction = res.Direction

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter = &DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartDataSetsTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsXyChartThresholds struct {
	empty     bool                                                               `json:"-"`
	Label     *string                                                            `json:"label"`
	Value     *float64                                                           `json:"value"`
	Color     *DashboardColumnLayoutColumnsWidgetsXyChartThresholdsColorEnum     `json:"color"`
	Direction *DashboardColumnLayoutColumnsWidgetsXyChartThresholdsDirectionEnum `json:"direction"`
}

type jsonDashboardColumnLayoutColumnsWidgetsXyChartThresholds DashboardColumnLayoutColumnsWidgetsXyChartThresholds

func (r *DashboardColumnLayoutColumnsWidgetsXyChartThresholds) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsXyChartThresholds
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsXyChartThresholds
	} else {

		r.Label = res.Label

		r.Value = res.Value

		r.Color = res.Color

		r.Direction = res.Direction

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsXyChartThresholds is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsXyChartThresholds *DashboardColumnLayoutColumnsWidgetsXyChartThresholds = &DashboardColumnLayoutColumnsWidgetsXyChartThresholds{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartThresholds) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartThresholds) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartThresholds) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsXyChartXAxis struct {
	empty bool                                                      `json:"-"`
	Label *string                                                   `json:"label"`
	Scale *DashboardColumnLayoutColumnsWidgetsXyChartXAxisScaleEnum `json:"scale"`
}

type jsonDashboardColumnLayoutColumnsWidgetsXyChartXAxis DashboardColumnLayoutColumnsWidgetsXyChartXAxis

func (r *DashboardColumnLayoutColumnsWidgetsXyChartXAxis) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsXyChartXAxis
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsXyChartXAxis
	} else {

		r.Label = res.Label

		r.Scale = res.Scale

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsXyChartXAxis is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsXyChartXAxis *DashboardColumnLayoutColumnsWidgetsXyChartXAxis = &DashboardColumnLayoutColumnsWidgetsXyChartXAxis{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartXAxis) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartXAxis) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartXAxis) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsXyChartYAxis struct {
	empty bool                                                      `json:"-"`
	Label *string                                                   `json:"label"`
	Scale *DashboardColumnLayoutColumnsWidgetsXyChartYAxisScaleEnum `json:"scale"`
}

type jsonDashboardColumnLayoutColumnsWidgetsXyChartYAxis DashboardColumnLayoutColumnsWidgetsXyChartYAxis

func (r *DashboardColumnLayoutColumnsWidgetsXyChartYAxis) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsXyChartYAxis
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsXyChartYAxis
	} else {

		r.Label = res.Label

		r.Scale = res.Scale

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsXyChartYAxis is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsXyChartYAxis *DashboardColumnLayoutColumnsWidgetsXyChartYAxis = &DashboardColumnLayoutColumnsWidgetsXyChartYAxis{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartYAxis) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartYAxis) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartYAxis) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsXyChartChartOptions struct {
	empty bool                                                            `json:"-"`
	Mode  *DashboardColumnLayoutColumnsWidgetsXyChartChartOptionsModeEnum `json:"mode"`
}

type jsonDashboardColumnLayoutColumnsWidgetsXyChartChartOptions DashboardColumnLayoutColumnsWidgetsXyChartChartOptions

func (r *DashboardColumnLayoutColumnsWidgetsXyChartChartOptions) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsXyChartChartOptions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsXyChartChartOptions
	} else {

		r.Mode = res.Mode

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsXyChartChartOptions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsXyChartChartOptions *DashboardColumnLayoutColumnsWidgetsXyChartChartOptions = &DashboardColumnLayoutColumnsWidgetsXyChartChartOptions{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartChartOptions) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartChartOptions) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsXyChartChartOptions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsScorecard struct {
	empty           bool                                                         `json:"-"`
	TimeSeriesQuery *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQuery `json:"timeSeriesQuery"`
	GaugeView       *DashboardColumnLayoutColumnsWidgetsScorecardGaugeView       `json:"gaugeView"`
	SparkChartView  *DashboardColumnLayoutColumnsWidgetsScorecardSparkChartView  `json:"sparkChartView"`
	Thresholds      []DashboardColumnLayoutColumnsWidgetsScorecardThresholds     `json:"thresholds"`
}

type jsonDashboardColumnLayoutColumnsWidgetsScorecard DashboardColumnLayoutColumnsWidgetsScorecard

func (r *DashboardColumnLayoutColumnsWidgetsScorecard) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsScorecard
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsScorecard
	} else {

		r.TimeSeriesQuery = res.TimeSeriesQuery

		r.GaugeView = res.GaugeView

		r.SparkChartView = res.SparkChartView

		r.Thresholds = res.Thresholds

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsScorecard is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsScorecard *DashboardColumnLayoutColumnsWidgetsScorecard = &DashboardColumnLayoutColumnsWidgetsScorecard{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsScorecard) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsScorecard) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsScorecard) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQuery struct {
	empty                   bool                                                                              `json:"-"`
	TimeSeriesFilter        *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilter      `json:"timeSeriesFilter"`
	TimeSeriesFilterRatio   *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio `json:"timeSeriesFilterRatio"`
	TimeSeriesQueryLanguage *string                                                                           `json:"timeSeriesQueryLanguage"`
	UnitOverride            *string                                                                           `json:"unitOverride"`
}

type jsonDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQuery DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQuery

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQuery) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQuery
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQuery
	} else {

		r.TimeSeriesFilter = res.TimeSeriesFilter

		r.TimeSeriesFilterRatio = res.TimeSeriesFilterRatio

		r.TimeSeriesQueryLanguage = res.TimeSeriesQueryLanguage

		r.UnitOverride = res.UnitOverride

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQuery is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQuery *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQuery = &DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQuery{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQuery) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQuery) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQuery) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilter struct {
	empty                bool                                                                                             `json:"-"`
	Filter               *string                                                                                          `json:"filter"`
	Aggregation          *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation          `json:"aggregation"`
	SecondaryAggregation *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation `json:"secondaryAggregation"`
	PickTimeSeriesFilter *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter `json:"pickTimeSeriesFilter"`
}

type jsonDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilter DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilter

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilter) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilter
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilter
	} else {

		r.Filter = res.Filter

		r.Aggregation = res.Aggregation

		r.SecondaryAggregation = res.SecondaryAggregation

		r.PickTimeSeriesFilter = res.PickTimeSeriesFilter

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilter is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilter *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilter = &DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilter{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilter) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilter) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilter) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation struct {
	empty              bool                                                                                                          `json:"-"`
	AlignmentPeriod    *string                                                                                                       `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                      `json:"groupByFields"`
}

type jsonDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation = &DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation struct {
	empty              bool                                                                                                                   `json:"-"`
	AlignmentPeriod    *string                                                                                                                `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                               `json:"groupByFields"`
}

type jsonDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation = &DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterSecondaryAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter struct {
	empty         bool                                                                                                              `json:"-"`
	RankingMethod *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterRankingMethodEnum `json:"rankingMethod"`
	NumTimeSeries *int64                                                                                                            `json:"numTimeSeries"`
	Direction     *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilterDirectionEnum     `json:"direction"`
}

type jsonDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter
	} else {

		r.RankingMethod = res.RankingMethod

		r.NumTimeSeries = res.NumTimeSeries

		r.Direction = res.Direction

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter = &DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterPickTimeSeriesFilter) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio struct {
	empty                bool                                                                                                  `json:"-"`
	Numerator            *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator            `json:"numerator"`
	Denominator          *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator          `json:"denominator"`
	SecondaryAggregation *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation `json:"secondaryAggregation"`
	PickTimeSeriesFilter *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter `json:"pickTimeSeriesFilter"`
}

type jsonDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio
	} else {

		r.Numerator = res.Numerator

		r.Denominator = res.Denominator

		r.SecondaryAggregation = res.SecondaryAggregation

		r.PickTimeSeriesFilter = res.PickTimeSeriesFilter

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio = &DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatio) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator struct {
	empty       bool                                                                                                  `json:"-"`
	Filter      *string                                                                                               `json:"filter"`
	Aggregation *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation `json:"aggregation"`
}

type jsonDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator
	} else {

		r.Filter = res.Filter

		r.Aggregation = res.Aggregation

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator = &DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumerator) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation struct {
	empty              bool                                                                                                                        `json:"-"`
	AlignmentPeriod    *string                                                                                                                     `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                                    `json:"groupByFields"`
}

type jsonDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation = &DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioNumeratorAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator struct {
	empty       bool                                                                                                    `json:"-"`
	Filter      *string                                                                                                 `json:"filter"`
	Aggregation *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation `json:"aggregation"`
}

type jsonDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator
	} else {

		r.Filter = res.Filter

		r.Aggregation = res.Aggregation

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator = &DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominator) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation struct {
	empty              bool                                                                                                                          `json:"-"`
	AlignmentPeriod    *string                                                                                                                       `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                                      `json:"groupByFields"`
}

type jsonDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation = &DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioDenominatorAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation struct {
	empty              bool                                                                                                                        `json:"-"`
	AlignmentPeriod    *string                                                                                                                     `json:"alignmentPeriod"`
	PerSeriesAligner   *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationPerSeriesAlignerEnum   `json:"perSeriesAligner"`
	CrossSeriesReducer *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregationCrossSeriesReducerEnum `json:"crossSeriesReducer"`
	GroupByFields      []string                                                                                                                    `json:"groupByFields"`
}

type jsonDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation
	} else {

		r.AlignmentPeriod = res.AlignmentPeriod

		r.PerSeriesAligner = res.PerSeriesAligner

		r.CrossSeriesReducer = res.CrossSeriesReducer

		r.GroupByFields = res.GroupByFields

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation = &DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioSecondaryAggregation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter struct {
	empty         bool                                                                                                                   `json:"-"`
	RankingMethod *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterRankingMethodEnum `json:"rankingMethod"`
	NumTimeSeries *int64                                                                                                                 `json:"numTimeSeries"`
	Direction     *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilterDirectionEnum     `json:"direction"`
}

type jsonDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter
	} else {

		r.RankingMethod = res.RankingMethod

		r.NumTimeSeries = res.NumTimeSeries

		r.Direction = res.Direction

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter = &DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardTimeSeriesQueryTimeSeriesFilterRatioPickTimeSeriesFilter) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsScorecardGaugeView struct {
	empty      bool     `json:"-"`
	LowerBound *float64 `json:"lowerBound"`
	UpperBound *float64 `json:"upperBound"`
}

type jsonDashboardColumnLayoutColumnsWidgetsScorecardGaugeView DashboardColumnLayoutColumnsWidgetsScorecardGaugeView

func (r *DashboardColumnLayoutColumnsWidgetsScorecardGaugeView) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsScorecardGaugeView
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsScorecardGaugeView
	} else {

		r.LowerBound = res.LowerBound

		r.UpperBound = res.UpperBound

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsScorecardGaugeView is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsScorecardGaugeView *DashboardColumnLayoutColumnsWidgetsScorecardGaugeView = &DashboardColumnLayoutColumnsWidgetsScorecardGaugeView{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardGaugeView) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardGaugeView) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardGaugeView) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsScorecardSparkChartView struct {
	empty              bool                                                                          `json:"-"`
	SparkChartType     *DashboardColumnLayoutColumnsWidgetsScorecardSparkChartViewSparkChartTypeEnum `json:"sparkChartType"`
	MinAlignmentPeriod *string                                                                       `json:"minAlignmentPeriod"`
}

type jsonDashboardColumnLayoutColumnsWidgetsScorecardSparkChartView DashboardColumnLayoutColumnsWidgetsScorecardSparkChartView

func (r *DashboardColumnLayoutColumnsWidgetsScorecardSparkChartView) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsScorecardSparkChartView
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsScorecardSparkChartView
	} else {

		r.SparkChartType = res.SparkChartType

		r.MinAlignmentPeriod = res.MinAlignmentPeriod

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsScorecardSparkChartView is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsScorecardSparkChartView *DashboardColumnLayoutColumnsWidgetsScorecardSparkChartView = &DashboardColumnLayoutColumnsWidgetsScorecardSparkChartView{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardSparkChartView) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardSparkChartView) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardSparkChartView) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsScorecardThresholds struct {
	empty     bool                                                                 `json:"-"`
	Label     *string                                                              `json:"label"`
	Value     *float64                                                             `json:"value"`
	Color     *DashboardColumnLayoutColumnsWidgetsScorecardThresholdsColorEnum     `json:"color"`
	Direction *DashboardColumnLayoutColumnsWidgetsScorecardThresholdsDirectionEnum `json:"direction"`
}

type jsonDashboardColumnLayoutColumnsWidgetsScorecardThresholds DashboardColumnLayoutColumnsWidgetsScorecardThresholds

func (r *DashboardColumnLayoutColumnsWidgetsScorecardThresholds) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsScorecardThresholds
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsScorecardThresholds
	} else {

		r.Label = res.Label

		r.Value = res.Value

		r.Color = res.Color

		r.Direction = res.Direction

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsScorecardThresholds is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsScorecardThresholds *DashboardColumnLayoutColumnsWidgetsScorecardThresholds = &DashboardColumnLayoutColumnsWidgetsScorecardThresholds{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardThresholds) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardThresholds) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsScorecardThresholds) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsText struct {
	empty   bool                                               `json:"-"`
	Content *string                                            `json:"content"`
	Format  *DashboardColumnLayoutColumnsWidgetsTextFormatEnum `json:"format"`
}

type jsonDashboardColumnLayoutColumnsWidgetsText DashboardColumnLayoutColumnsWidgetsText

func (r *DashboardColumnLayoutColumnsWidgetsText) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsText
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsText
	} else {

		r.Content = res.Content

		r.Format = res.Format

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsText is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsText *DashboardColumnLayoutColumnsWidgetsText = &DashboardColumnLayoutColumnsWidgetsText{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsText) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsText) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsText) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsBlank struct {
	empty bool `json:"-"`
}

type jsonDashboardColumnLayoutColumnsWidgetsBlank DashboardColumnLayoutColumnsWidgetsBlank

func (r *DashboardColumnLayoutColumnsWidgetsBlank) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsBlank
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsBlank
	} else {

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsBlank is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsBlank *DashboardColumnLayoutColumnsWidgetsBlank = &DashboardColumnLayoutColumnsWidgetsBlank{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsBlank) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsBlank) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsBlank) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DashboardColumnLayoutColumnsWidgetsLogsPanel struct {
	empty         bool     `json:"-"`
	Filter        *string  `json:"filter"`
	ResourceNames []string `json:"resourceNames"`
}

type jsonDashboardColumnLayoutColumnsWidgetsLogsPanel DashboardColumnLayoutColumnsWidgetsLogsPanel

func (r *DashboardColumnLayoutColumnsWidgetsLogsPanel) UnmarshalJSON(data []byte) error {
	var res jsonDashboardColumnLayoutColumnsWidgetsLogsPanel
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDashboardColumnLayoutColumnsWidgetsLogsPanel
	} else {

		r.Filter = res.Filter

		r.ResourceNames = res.ResourceNames

	}
	return nil
}

// This object is used to assert a desired state where this DashboardColumnLayoutColumnsWidgetsLogsPanel is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDashboardColumnLayoutColumnsWidgetsLogsPanel *DashboardColumnLayoutColumnsWidgetsLogsPanel = &DashboardColumnLayoutColumnsWidgetsLogsPanel{empty: true}

func (r *DashboardColumnLayoutColumnsWidgetsLogsPanel) Empty() bool {
	return r.empty
}

func (r *DashboardColumnLayoutColumnsWidgetsLogsPanel) String() string {
	return dcl.SprintResource(r)
}

func (r *DashboardColumnLayoutColumnsWidgetsLogsPanel) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *Dashboard) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "monitoring",
		Type:    "Dashboard",
		Version: "monitoring",
	}
}

func (r *Dashboard) ID() (string, error) {
	if err := extractDashboardFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name":          dcl.ValueOrEmptyString(nr.Name),
		"display_name":  dcl.ValueOrEmptyString(nr.DisplayName),
		"grid_layout":   dcl.ValueOrEmptyString(nr.GridLayout),
		"mosaic_layout": dcl.ValueOrEmptyString(nr.MosaicLayout),
		"row_layout":    dcl.ValueOrEmptyString(nr.RowLayout),
		"column_layout": dcl.ValueOrEmptyString(nr.ColumnLayout),
		"project":       dcl.ValueOrEmptyString(nr.Project),
		"etag":          dcl.ValueOrEmptyString(nr.Etag),
	}
	return dcl.Nprintf("projects/{{project}}/dashboards/{{name}}", params), nil
}

const DashboardMaxPage = -1

type DashboardList struct {
	Items []*Dashboard

	nextToken string

	pageSize int32

	resource *Dashboard
}

func (l *DashboardList) HasNext() bool {
	return l.nextToken != ""
}

func (l *DashboardList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listDashboard(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListDashboard(ctx context.Context, project string) (*DashboardList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListDashboardWithMaxResults(ctx, project, DashboardMaxPage)

}

func (c *Client) ListDashboardWithMaxResults(ctx context.Context, project string, pageSize int32) (*DashboardList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &Dashboard{
		Project: &project,
	}
	items, token, err := c.listDashboard(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &DashboardList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetDashboard(ctx context.Context, r *Dashboard) (*Dashboard, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractDashboardFields(r)

	b, err := c.getDashboardRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalDashboard(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeDashboardNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractDashboardFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteDashboard(ctx context.Context, r *Dashboard) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("Dashboard resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting Dashboard...")
	deleteOp := deleteDashboardOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllDashboard deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllDashboard(ctx context.Context, project string, filter func(*Dashboard) bool) error {
	listObj, err := c.ListDashboard(ctx, project)
	if err != nil {
		return err
	}

	err = c.deleteAllDashboard(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllDashboard(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyDashboard(ctx context.Context, rawDesired *Dashboard, opts ...dcl.ApplyOption) (*Dashboard, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *Dashboard
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyDashboardHelper(c, ctx, rawDesired, opts...)
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

func applyDashboardHelper(c *Client, ctx context.Context, rawDesired *Dashboard, opts ...dcl.ApplyOption) (*Dashboard, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyDashboard...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractDashboardFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.dashboardDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToDashboardDiffs(c.Config, fieldDiffs, opts)
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
	var ops []dashboardApiOperation
	if create {
		ops = append(ops, &createDashboardOperation{})
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
	return applyDashboardDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyDashboardDiff(c *Client, ctx context.Context, desired *Dashboard, rawDesired *Dashboard, ops []dashboardApiOperation, opts ...dcl.ApplyOption) (*Dashboard, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetDashboard(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createDashboardOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapDashboard(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeDashboardNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeDashboardNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeDashboardDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractDashboardFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractDashboardFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffDashboard(c, newDesired, newState)
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
