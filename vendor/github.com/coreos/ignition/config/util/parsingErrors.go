// Copyright 2018 CoreOS, Inc.
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

package util

import (
	"github.com/coreos/ignition/config/shared/errors"
	"github.com/coreos/ignition/config/validate/report"
	"github.com/coreos/ignition/config/validate/util"

	json "github.com/ajeddeloh/go-json"
)

// HandleParseErrors will attempt to unmarshal an invalid rawConfig into "to".
// If it fails to unmarsh it will generate a report.Report from the errors.
func HandleParseErrors(rawConfig []byte, to interface{}) (report.Report, error) {
	err := json.Unmarshal(rawConfig, to)
	if err == nil {
		return report.Report{}, nil
	}

	// Handle json syntax and type errors first, since they are fatal but have offset info
	if serr, ok := err.(*json.SyntaxError); ok {
		line, col, highlight := util.Highlight(rawConfig, serr.Offset)
		return report.Report{
				Entries: []report.Entry{{
					Kind:      report.EntryError,
					Message:   serr.Error(),
					Line:      line,
					Column:    col,
					Highlight: highlight,
				}},
			},
			errors.ErrInvalid
	}

	if terr, ok := err.(*json.UnmarshalTypeError); ok {
		line, col, highlight := util.Highlight(rawConfig, terr.Offset)
		return report.Report{
				Entries: []report.Entry{{
					Kind:      report.EntryError,
					Message:   terr.Error(),
					Line:      line,
					Column:    col,
					Highlight: highlight,
				}},
			},
			errors.ErrInvalid
	}

	return report.ReportFromError(err, report.EntryError), errors.ErrInvalid
}
