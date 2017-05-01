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
	ignTypes "github.com/coreos/ignition/config/v2_0/types"
	"github.com/coreos/ignition/config/validate/report"
)

type Networkd struct {
	Units []NetworkdUnit `yaml:"units"`
}

type NetworkdUnit struct {
	Name     string `yaml:"name"`
	Contents string `yaml:"contents"`
}

func init() {
	register2_0(func(in Config, out ignTypes.Config, platform string) (ignTypes.Config, report.Report) {
		for _, unit := range in.Networkd.Units {
			out.Networkd.Units = append(out.Networkd.Units, ignTypes.NetworkdUnit{
				Name:     ignTypes.NetworkdUnitName(unit.Name),
				Contents: unit.Contents,
			})
		}
		return out, report.Report{}
	})
}
