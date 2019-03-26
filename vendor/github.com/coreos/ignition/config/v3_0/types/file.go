// Copyright 2016 CoreOS, Inc.
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

package types

import (
	"fmt"

	"github.com/coreos/ignition/config/shared/errors"
	"github.com/coreos/ignition/config/validate/report"
)

func (f File) ValidateMode() (r report.Report) {
	r.AddOnError(validateMode(f.Mode))
	if f.Mode == nil {
		r.AddOnWarning(errors.ErrFilePermissionsUnset)
	}
	return r
}

func (f FileEmbedded1) IgnoreDuplicates() map[string]struct{} {
	return map[string]struct{}{
		"Append": {},
	}
}

func (fc FileContents) ValidateCompression() (r report.Report) {
	if fc.Compression == nil {
		return
	}
	switch *fc.Compression {
	case "", "gzip":
	default:
		r.AddOnError(errors.ErrCompressionInvalid)
	}
	return
}

func (fc FileContents) ValidateSource() (r report.Report) {
	if fc.Source == nil {
		return
	}
	err := validateURL(*fc.Source)
	if err != nil {
		r.AddOnError(fmt.Errorf("invalid url %q: %v", *fc.Source, err))
	}
	return
}
