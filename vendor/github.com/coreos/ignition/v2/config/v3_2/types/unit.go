// Copyright 2020 Red Hat, Inc.
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

	"github.com/coreos/ignition/v2/config/shared/errors"
	"github.com/coreos/ignition/v2/config/shared/validations"

	"github.com/coreos/go-systemd/v22/unit"
	cpath "github.com/coreos/vcontext/path"
	"github.com/coreos/vcontext/report"
)

func (u Unit) Key() string {
	return u.Name
}

func (d Dropin) Key() string {
	return d.Name
}

func (u Unit) Validate(c cpath.ContextPath) (r report.Report) {
	r.AddOnError(c.Append("name"), validateName(u.Name))
	c = c.Append("contents")
	opts, err := validateUnitContent(u.Contents)
	r.AddOnError(c, err)

	isEnabled := u.Enabled != nil && *u.Enabled
	r.AddOnWarn(c, validations.ValidateInstallSection(u.Name, isEnabled, (u.Contents == nil || *u.Contents == ""), opts))

	return
}

func validateName(name string) error {
	switch path.Ext(name) {
	case ".service", ".socket", ".device", ".mount", ".automount", ".swap", ".target", ".path", ".timer", ".snapshot", ".slice", ".scope":
	default:
		return errors.ErrInvalidSystemdExt
	}
	return nil
}

func (d Dropin) Validate(c cpath.ContextPath) (r report.Report) {
	_, err := validateUnitContent(d.Contents)
	r.AddOnError(c.Append("contents"), err)

	switch path.Ext(d.Name) {
	case ".conf":
	default:
		r.AddOnError(c.Append("name"), errors.ErrInvalidSystemdDropinExt)
	}

	return
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
