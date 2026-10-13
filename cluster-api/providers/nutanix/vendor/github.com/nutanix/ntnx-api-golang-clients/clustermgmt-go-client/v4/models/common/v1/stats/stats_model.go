/*
 * Generated file models/common/v1/stats/stats_model.go.
 *
 * Product version: 4.0.1-beta-2
 *
 * Part of the Nutanix Clustermgmt Versioned APIs
 *
 * (c) 2024 Nutanix Inc.  All rights reserved
 *
 */

/*
  Module common.v1.stats of Nutanix Clustermgmt Versioned APIs
*/
package stats

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

/*
The operator to use while performing down-sampling on stats data. Allowed values are SUM, MIN, MAX, AVG, COUNT and LAST.
*/
type DownSamplingOperator int

const (
	DOWNSAMPLINGOPERATOR_UNKNOWN  DownSamplingOperator = 0
	DOWNSAMPLINGOPERATOR_REDACTED DownSamplingOperator = 1
	DOWNSAMPLINGOPERATOR_SUM      DownSamplingOperator = 2
	DOWNSAMPLINGOPERATOR_MIN      DownSamplingOperator = 3
	DOWNSAMPLINGOPERATOR_MAX      DownSamplingOperator = 4
	DOWNSAMPLINGOPERATOR_AVG      DownSamplingOperator = 5
	DOWNSAMPLINGOPERATOR_COUNT    DownSamplingOperator = 6
	DOWNSAMPLINGOPERATOR_LAST     DownSamplingOperator = 7
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *DownSamplingOperator) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"SUM",
		"MIN",
		"MAX",
		"AVG",
		"COUNT",
		"LAST",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e DownSamplingOperator) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"SUM",
		"MIN",
		"MAX",
		"AVG",
		"COUNT",
		"LAST",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *DownSamplingOperator) index(name string) DownSamplingOperator {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"SUM",
		"MIN",
		"MAX",
		"AVG",
		"COUNT",
		"LAST",
	}
	for idx := range names {
		if names[idx] == name {
			return DownSamplingOperator(idx)
		}
	}
	return DOWNSAMPLINGOPERATOR_UNKNOWN
}

func (e *DownSamplingOperator) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for DownSamplingOperator:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *DownSamplingOperator) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e DownSamplingOperator) Ref() *DownSamplingOperator {
	return &e
}

/*
A time value pair representing a stat associated with a given entity at a given point of date and time represented in extended ISO-8601 format."
*/
type TimeIntValuePair struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  The date and time at which the stat was recorded.The value should be in extended ISO-8601 format. For example, start time of 2022-04-23T01:23:45.678+09:00 would consider all stats starting at 1:23:45.678 on the 23rd of April 2022. Details around ISO-8601 format can be found at https://www.iso.org/standard/70907.html
	*/
	Timestamp *time.Time `json:"timestamp,omitempty"`
	/*
	  Value of the stat at the recorded date and time in extended ISO-8601 format."
	*/
	Value *int64 `json:"value,omitempty"`
}

func NewTimeIntValuePair() *TimeIntValuePair {
	p := new(TimeIntValuePair)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "common.v1.stats.TimeIntValuePair"
	p.Reserved_ = map[string]interface{}{"$fv": "v1.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

type FileDetail struct {
	Path        *string `json:"-"`
	ObjectType_ *string `json:"-"`
}

func NewFileDetail() *FileDetail {
	p := new(FileDetail)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "FileDetail"

	return p
}
