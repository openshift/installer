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

package v3_0

import (
	"reflect"

	"github.com/coreos/ignition/v2/config/merge"
	"github.com/coreos/ignition/v2/config/shared/errors"
	"github.com/coreos/ignition/v2/config/util"
	"github.com/coreos/ignition/v2/config/v3_0/types"
	"github.com/coreos/ignition/v2/config/validate"

	"github.com/coreos/go-semver/semver"
	"github.com/coreos/vcontext/report"
)

func Merge(parent, child types.Config) types.Config {
	vParent := reflect.ValueOf(parent)
	vChild := reflect.ValueOf(child)

	vRes := merge.MergeStruct(vParent, vChild)
	res := vRes.Interface().(types.Config)
	return res
}

// Parse parses the raw config into a types.Config struct and generates a report of any
// errors, warnings, info, and deprecations it encountered
func Parse(rawConfig []byte) (types.Config, report.Report, error) {
	if isEmpty(rawConfig) {
		return types.Config{}, report.Report{}, errors.ErrEmpty
	}

	var config types.Config
	if rpt, err := util.HandleParseErrors(rawConfig, &config); err != nil {
		return types.Config{}, rpt, err
	}

	version, err := semver.NewVersion(config.Ignition.Version)

	if err != nil || *version != types.MaxVersion {
		return types.Config{}, report.Report{}, errors.ErrUnknownVersion
	}

	rpt := validate.ValidateWithContext(config, rawConfig)
	if rpt.IsFatal() {
		return types.Config{}, rpt, errors.ErrInvalid
	}

	return config, rpt, nil
}

func isEmpty(userdata []byte) bool {
	return len(userdata) == 0
}
