/*
Copyright 2025 The Kubernetes Authors.

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

package util

import (
	"bytes"
	"fmt"
	"text/template"

	sprig "github.com/go-task/slim-sprig/v3"
	"github.com/pkg/errors"
)

const (
	maxNameLength = 63
)

var nameTemplateFuncs = map[string]any{
	"trimSuffix": sprig.GenericFuncMap()["trimSuffix"],
	"trunc":      sprig.GenericFuncMap()["trunc"],
}

var nameTpl = template.New("name generator").Funcs(nameTemplateFuncs).Option("missingkey=error")

// GenerateMachineNameFromTemplate generate a name from machine name and a naming strategy template.
// the template supports only `trimSuffix` and `trunc` functions.
func GenerateMachineNameFromTemplate(machineName string, nameTemplate *string) (string, error) {
	if machineName == "" {
		return "", fmt.Errorf("machine name can not be emmpty")
	}

	if nameTemplate == nil {
		return machineName, nil
	}

	data := map[string]interface{}{
		"machine": map[string]interface{}{
			"name": machineName,
		},
	}

	tpl, err := nameTpl.Parse(*nameTemplate)
	if err != nil {
		return "", errors.Wrapf(err, "unable to parse template %q", *nameTemplate)
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}

	name := buf.String()

	// If the name exceeds the maxNameLength, trim to maxNameLength.
	// Note: we're not adding a random suffix as the name has to be deterministic.
	if len(name) > maxNameLength {
		name = name[:maxNameLength]
	}

	return name, nil
}
