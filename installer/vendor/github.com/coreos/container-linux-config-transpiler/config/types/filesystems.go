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

type Filesystem struct {
	Name  string `yaml:"name"`
	Mount *Mount `yaml:"mount"`
	Path  string `yaml:"path"`
}

type Mount struct {
	Device string  `yaml:"device"`
	Format string  `yaml:"format"`
	Create *Create `yaml:"create"`
}

type Create struct {
	Force   bool     `yaml:"force"`
	Options []string `yaml:"options"`
}

func init() {
	register2_0(func(in Config, out ignTypes.Config, platform string) (ignTypes.Config, report.Report) {
		for _, filesystem := range in.Storage.Filesystems {
			newFilesystem := ignTypes.Filesystem{
				Name: filesystem.Name,
				Path: func(p ignTypes.Path) *ignTypes.Path {
					if p == "" {
						return nil
					}

					return &p
				}(ignTypes.Path(filesystem.Path)),
			}

			if filesystem.Mount != nil {
				newFilesystem.Mount = &ignTypes.FilesystemMount{
					Device: ignTypes.Path(filesystem.Mount.Device),
					Format: ignTypes.FilesystemFormat(filesystem.Mount.Format),
				}

				if filesystem.Mount.Create != nil {
					newFilesystem.Mount.Create = &ignTypes.FilesystemCreate{
						Force:   filesystem.Mount.Create.Force,
						Options: ignTypes.MkfsOptions(filesystem.Mount.Create.Options),
					}
				}
			}

			out.Storage.Filesystems = append(out.Storage.Filesystems, newFilesystem)
		}
		return out, report.Report{}
	})
}
