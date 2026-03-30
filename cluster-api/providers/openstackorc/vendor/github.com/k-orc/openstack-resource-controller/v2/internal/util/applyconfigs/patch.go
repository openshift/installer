/*
Copyright 2024 The ORC Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package applyconfigs

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// applyConfigPatch is a slightly more ergonomic version of client.RawPatch which json marshals its argument
type applyConfigPatch struct {
	patchType   types.PatchType
	applyConfig interface{}
}

// Type implements Patch.
func (p applyConfigPatch) Type() types.PatchType {
	return p.patchType
}

// Data implements Patch.
func (p applyConfigPatch) Data(_ client.Object) ([]byte, error) {
	return json.Marshal(p.applyConfig)
}

func Patch(patchType types.PatchType, applyConfig interface{}) client.Patch {
	return &applyConfigPatch{
		patchType:   patchType,
		applyConfig: applyConfig,
	}
}
