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
	"github.com/vincent-petithory/dataurl"
)

type File struct {
	Filesystem string       `yaml:"filesystem"`
	Path       string       `yaml:"path"`
	Mode       int          `yaml:"mode"`
	Contents   FileContents `yaml:"contents"`
	User       FileUser     `yaml:"user"`
	Group      FileGroup    `yaml:"group"`
}

type FileContents struct {
	Remote Remote `yaml:"remote"`
	Inline string `yaml:"inline"`
}

type Remote struct {
	Url          string       `yaml:"url"`
	Compression  string       `yaml:"compression"`
	Verification Verification `yaml:"verification"`
}

type FileUser struct {
	Id int `yaml:"id"`
}

type FileGroup struct {
	Id int `yaml:"id"`
}

func init() {
	register2_0(func(in Config, out ignTypes.Config, platform string) (ignTypes.Config, report.Report) {
		for _, file := range in.Storage.Files {
			newFile := ignTypes.File{
				Filesystem: file.Filesystem,
				Path:       ignTypes.Path(file.Path),
				Mode:       ignTypes.FileMode(file.Mode),
				User:       ignTypes.FileUser{Id: file.User.Id},
				Group:      ignTypes.FileGroup{Id: file.Group.Id},
			}

			if file.Contents.Inline != "" {
				newFile.Contents = ignTypes.FileContents{
					Source: ignTypes.Url{
						Scheme: "data",
						Opaque: "," + dataurl.EscapeString(file.Contents.Inline),
					},
				}
			}

			if file.Contents.Remote.Url != "" {
				source, err := url.Parse(file.Contents.Remote.Url)
				if err != nil {
					return out, report.ReportFromError(err, report.EntryError)
				}

				newFile.Contents = ignTypes.FileContents{Source: ignTypes.Url(*source)}
			}

			if newFile.Contents == (ignTypes.FileContents{}) {
				newFile.Contents = ignTypes.FileContents{
					Source: ignTypes.Url{
						Scheme: "data",
						Opaque: ",",
					},
				}
			}

			newFile.Contents.Compression = ignTypes.Compression(file.Contents.Remote.Compression)
			newFile.Contents.Verification = convertVerification(file.Contents.Remote.Verification)

			out.Storage.Files = append(out.Storage.Files, newFile)
		}
		return out, report.Report{}
	})
}
