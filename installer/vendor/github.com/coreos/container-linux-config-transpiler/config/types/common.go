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
	"reflect"

	"github.com/coreos/container-linux-config-transpiler/config/templating"
)

var (
	ErrPlatformUnspecified = fmt.Errorf("platform must be specified to use templating")
)

func isZero(v interface{}) bool {
	if v == nil {
		return true
	}
	zv := reflect.Zero(reflect.TypeOf(v))
	return reflect.DeepEqual(v, zv.Interface())
}

// assembleUnit will assemble the contents of a systemd unit dropin that will
// have the given environment variables, and call the given exec line with the
// provided args prepended to it
func assembleUnit(exec string, args, vars []string, platform string) (string, error) {
	hasTemplating := templating.HasTemplating(args)

	var out string
	if hasTemplating {
		if platform == "" {
			return "", ErrPlatformUnspecified
		}
		out = "[Unit]\nRequires=coreos-metadata.service\nAfter=coreos-metadata.service\n\n[Service]\nEnvironmentFile=/run/metadata/coreos\n"
		var err error
		args, err = templating.PerformTemplating(platform, args)
		if err != nil {
			return "", err
		}
	} else {
		out = "[Service]\n"
	}

	for _, v := range vars {
		out += fmt.Sprintf("Environment=\"%s\"\n", v)
	}
	for _, a := range args {
		exec += fmt.Sprintf(" \\\n  %s", a)
	}
	out += "ExecStart=\n"
	out += fmt.Sprintf("ExecStart=%s", exec)
	return out, nil
}

// getCliArgs builds a list of --ARG=VAL from a struct with cli: tags on its members.
func getCliArgs(e interface{}) []string {
	if e == nil {
		return nil
	}
	et := reflect.TypeOf(e)
	ev := reflect.ValueOf(e)

	vars := []string{}
	for i := 0; i < et.NumField(); i++ {
		if val := ev.Field(i).Interface(); !isZero(val) {
			if et.Field(i).Anonymous {
				vars = append(vars, getCliArgs(val)...)
			} else {
				key := et.Field(i).Tag.Get("cli")
				vars = append(vars, fmt.Sprintf("--%s=%q", key, val))
			}
		}
	}

	return vars
}
