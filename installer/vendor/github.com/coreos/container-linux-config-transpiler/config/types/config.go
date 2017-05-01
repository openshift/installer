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
	"net/url"

	ignTypes "github.com/coreos/ignition/config/v2_0/types"
	"github.com/coreos/ignition/config/validate/report"
)

type Config struct {
	Ignition  Ignition   `yaml:"ignition"`
	Storage   Storage    `yaml:"storage"`
	Systemd   Systemd    `yaml:"systemd"`
	Networkd  Networkd   `yaml:"networkd"`
	Passwd    Passwd     `yaml:"passwd"`
	Etcd      *Etcd      `yaml:"etcd"`
	Flannel   *Flannel   `yaml:"flannel"`
	Update    *Update    `yaml:"update"`
	Docker    *Docker    `yaml:"docker"`
	Locksmith *Locksmith `yaml:"locksmith"`
}

type Ignition struct {
	Config IgnitionConfig `yaml:"config"`
}

type IgnitionConfig struct {
	Append  []ConfigReference `yaml:"append"`
	Replace *ConfigReference  `yaml:"replace"`
}

type ConfigReference struct {
	Source       string       `yaml:"source"`
	Verification Verification `yaml:"verification"`
}

func init() {
	register2_0(func(in Config, out ignTypes.Config, platform string) (ignTypes.Config, report.Report) {
		for _, ref := range in.Ignition.Config.Append {
			newRef, err := convertConfigReference(ref)
			if err != nil {
				return out, report.ReportFromError(err, report.EntryError)
			}
			out.Ignition.Config.Append = append(out.Ignition.Config.Append, newRef)
		}

		if in.Ignition.Config.Replace != nil {
			newRef, err := convertConfigReference(*in.Ignition.Config.Replace)
			if err != nil {
				return out, report.ReportFromError(err, report.EntryError)
			}
			out.Ignition.Config.Replace = &newRef
		}
		return out, report.Report{}
	})
}

func convertConfigReference(in ConfigReference) (ignTypes.ConfigReference, error) {
	source, err := url.Parse(in.Source)
	if err != nil {
		return ignTypes.ConfigReference{}, err
	}

	return ignTypes.ConfigReference{
		Source:       ignTypes.Url(*source),
		Verification: convertVerification(in.Verification),
	}, nil
}

func convertVerification(in Verification) ignTypes.Verification {
	if in.Hash.Function == "" || in.Hash.Sum == "" {
		return ignTypes.Verification{}
	}

	return ignTypes.Verification{
		&ignTypes.Hash{
			Function: in.Hash.Function,
			Sum:      in.Hash.Sum,
		},
	}
}
