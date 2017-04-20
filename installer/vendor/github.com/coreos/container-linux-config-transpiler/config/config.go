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

package config

import (
	"reflect"

	yaml "github.com/ajeddeloh/yaml"
	"github.com/coreos/container-linux-config-transpiler/config/types"
	ignTypes "github.com/coreos/ignition/config/v2_0/types"
	"github.com/coreos/ignition/config/validate"
	"github.com/coreos/ignition/config/validate/report"
)

func Parse(data []byte) (types.Config, report.Report) {
	var cfg types.Config
	var r report.Report

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return types.Config{}, report.ReportFromError(err, report.EntryError)
	}

	nodes := yaml.UnmarshalToNode(data)
	if nodes == nil {
		r.Add(report.Entry{
			Kind:    report.EntryWarning,
			Message: "Configuration is empty",
		})
		r.Merge(validate.ValidateWithoutSource(reflect.ValueOf(cfg)))
	} else {
		root, err := fromYamlDocumentNode(*nodes)
		if err != nil {
			return types.Config{}, report.ReportFromError(err, report.EntryError)
		}

		r.Merge(validate.Validate(reflect.ValueOf(cfg), root, nil))
	}

	if r.IsFatal() {
		return types.Config{}, r
	}
	return cfg, r
}

func ConvertAs2_0(in types.Config, platform string) (ignTypes.Config, report.Report) {
	return types.ConvertAs2_0(in, platform)
}
