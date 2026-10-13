/*
 * Generated file models/common/v1/stats/stats_model.go.
 *
 * Product version: 4.0.2-beta-1
 *
 * Part of the Nutanix Networking Versioned APIs
 *
 * (c) 2024 Nutanix Inc.  All rights reserved
 *
 */

/*
  Nutanix Stats Configuration
*/
package stats

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
