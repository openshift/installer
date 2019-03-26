// Copyright 2015 CoreOS, Inc.
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
	"github.com/coreos/go-semver/semver"

	"github.com/coreos/ignition/config/shared/errors"
	"github.com/coreos/ignition/config/validate/report"
)

func (c ConfigReference) Key() string {
	if c.Source == nil {
		return ""
	}
	return *c.Source
}

func (c ConfigReference) ValidateSource() (r report.Report) {
	if c.Source == nil {
		return
	}
	r.AddOnError(validateURL(*c.Source))
	return
}

func (v Ignition) Semver() (*semver.Version, error) {
	return semver.NewVersion(v.Version)
}

func (v Ignition) Validate() report.Report {
	tv, err := v.Semver()
	if err != nil {
		return report.ReportFromError(errors.ErrInvalidVersion, report.EntryError)
	}

	if MaxVersion != *tv {
		return report.ReportFromError(errors.ErrUnknownVersion, report.EntryError)
	}
	return report.Report{}
}
