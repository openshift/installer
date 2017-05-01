// Copyright 2017 CoreOS, Inc.
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
	"reflect"

	ignTypes "github.com/coreos/ignition/config/v2_0/types"
	"github.com/coreos/ignition/config/validate"
	"github.com/coreos/ignition/config/validate/report"
)

type converterFor2_0 func(in Config, out ignTypes.Config, platform string) (ignTypes.Config, report.Report)

var convertersFor2_0 []converterFor2_0

func register2_0(f converterFor2_0) {
	convertersFor2_0 = append(convertersFor2_0, f)
}

func ConvertAs2_0(in Config, platform string) (ignTypes.Config, report.Report) {
	out := ignTypes.Config{
		Ignition: ignTypes.Ignition{
			Version: ignTypes.IgnitionVersion{Major: 2, Minor: 0},
		},
	}

	r := report.Report{}

	for _, convert := range convertersFor2_0 {
		var subReport report.Report
		out, subReport = convert(in, out, platform)
		r.Merge(subReport)
	}
	if r.IsFatal() {
		return ignTypes.Config{}, r
	}

	validationReport := validate.ValidateWithoutSource(reflect.ValueOf(out))
	r.Merge(validationReport)
	if r.IsFatal() {
		return ignTypes.Config{}, r
	}

	return out, r
}
