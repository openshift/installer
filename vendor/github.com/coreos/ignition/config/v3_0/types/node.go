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
	"path/filepath"

	"github.com/coreos/ignition/config/shared/errors"
	"github.com/coreos/ignition/config/validate/report"
)

func (n Node) Key() string {
	return n.Path
}

func (n Node) ValidatePath() (r report.Report) {
	r.AddOnError(validatePath(n.Path))
	return
}

func (n Node) Depth() int {
	count := 0
	for p := filepath.Clean(string(n.Path)); p != "/"; count++ {
		p = filepath.Dir(p)
	}
	return count
}

func (nu NodeUser) Validate() (r report.Report) {
	if nu.ID != nil && (nu.Name != nil && *nu.Name != "") {
		r.AddOnError(errors.ErrBothIDAndNameSet)
	}
	return
}
func (ng NodeGroup) Validate() (r report.Report) {
	if ng.ID != nil && (ng.Name != nil && *ng.Name != "") {
		r.AddOnError(errors.ErrBothIDAndNameSet)
	}
	return
}
