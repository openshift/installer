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
	"path"
	"strings"

	"github.com/coreos/go-systemd/unit"

	"github.com/coreos/ignition/config/shared/errors"
	"github.com/coreos/ignition/config/shared/validations"
	"github.com/coreos/ignition/config/validate/report"
)

func (u Unit) Key() string {
	return u.Name
}

func (d Dropin) Key() string {
	return d.Name
}

func (u Unit) ValidateContents() (r report.Report) {
	opts, err := validateUnitContent(u.Contents)
	r.AddOnError(err)

	isEnabled := u.Enabled != nil && *u.Enabled
	r.Merge(validations.ValidateInstallSection(u.Name, isEnabled, (u.Contents == nil || *u.Contents == ""), opts))

	return r
}

func (u Unit) ValidateName() (r report.Report) {
	switch path.Ext(u.Name) {
	case ".service", ".socket", ".device", ".mount", ".automount", ".swap", ".target", ".path", ".timer", ".snapshot", ".slice", ".scope":
	default:
		r.AddOnError(errors.ErrInvalidSystemdExt)
	}
	return
}

func (d Dropin) Validate() report.Report {
	r := report.Report{}

	_, err := validateUnitContent(d.Contents)
	r.AddOnError(err)

	switch path.Ext(d.Name) {
	case ".conf":
	default:
		r.AddOnError(errors.ErrInvalidSystemdDropinExt)
	}

	return r
}

func validateUnitContent(content *string) ([]*unit.UnitOption, error) {
	if content == nil {
		return []*unit.UnitOption{}, nil
	}
	c := strings.NewReader(*content)
	opts, err := unit.Deserialize(c)
	if err != nil {
		return nil, fmt.Errorf("invalid unit content: %s", err)
	}
	return opts, nil
}
