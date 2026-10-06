/*
Copyright 2026 The Kubernetes Authors.

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

package userdata

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
)

var (
	defaultTemplateFuncMap = template.FuncMap{
		"Indent": templateYAMLIndent,
		"toYaml": templateToYAML,
	}
)

func templateYAMLIndent(i int, input string) string {
	split := strings.Split(input, "\n")
	ident := "\n" + strings.Repeat(" ", i)
	return strings.Repeat(" ", i) + strings.Join(split, ident)
}

func templateToYAML(r *runtime.RawExtension) (string, error) {
	if r == nil {
		return "", nil
	}
	if r.Object != nil {
		b, err := yaml.Marshal(r.Object)
		if err != nil {
			return "", errors.Wrap(err, "failed to convert to yaml")
		}
		return string(b), nil
	}
	if len(r.Raw) == 0 {
		return "", nil
	}
	if yb, err := yaml.JSONToYAML(r.Raw); err == nil {
		return string(yb), nil
	}
	var temp interface{}
	err := yaml.Unmarshal(r.Raw, &temp)
	if err == nil {
		return string(r.Raw), nil
	}
	return "", fmt.Errorf("runtime object raw is neither json nor yaml %s", string(r.Raw))
}
