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
	"fmt"
	"strings"

	ignTypes "github.com/coreos/ignition/config/v2_0/types"
	"github.com/coreos/ignition/config/validate/report"
)

type Docker struct {
	Flags []string `yaml:"flags"`
}

func init() {
	register2_0(func(in Config, out ignTypes.Config, platform string) (ignTypes.Config, report.Report) {
		if in.Docker != nil {
			contents := fmt.Sprintf("[Service]\nEnvironment=\"DOCKER_OPTS=%s\"", strings.Join(in.Docker.Flags, " "))
			out.Systemd.Units = append(out.Systemd.Units, ignTypes.SystemdUnit{
				Name:   "docker.service",
				Enable: true,
				DropIns: []ignTypes.SystemdUnitDropIn{{
					Name:     "20-clct-docker.conf",
					Contents: contents,
				}},
			})
		}
		return out, report.Report{}
	})
}
