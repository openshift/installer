/*
Copyright (c) 2020 Red Hat, Inc.

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

// Package configuration provides a mechanism to load configuration from JSON or YAML files. The
// typical use will be to create a configuration object and then load one or more configuration
// sources:
//
//	// Load the configuration from a file:
//	cfg, err := configuration.New().
//		Load("myconfig.yaml").
//		Build()
//	if err != nil {
//		...
//	}
//
// Once the configuration is loaded it can be copied into an object containing the same tags
// used by the YAML library:
//
//	// Copy the configuration into our object:
//	type MyConfig struct {
//		MyKey string `yaml:"mykey"`
//		YouKey int `yaml:"yourkey"`
//	}
//	var myCfg MyConfig
//	err = cfg.Populate(&myCfg)
//	if err != nil {
//		...
//	}
//
// The advantage of using this configuration instead of using plain YAML is that configuration
// sources can use the tags to reference environment variables, files and scripts. For example:
//
//	mykey: !variable MYVARIABLE
//	yourkey: !file /my/file.txt
//
// The following tags are supported:
//
//	!file /my/file.txt - Is replaced by the content of the file `/my/file.txt`.
//	!script myscript - Is replaced by the result of executing the `myscript` script.
//	!trim mytext - Is replaced by the result of trimming white space from `mytext`.
//	!variable MYVARIABLE - Is replaced by the content of the environment variable `MYVARIABLE`.
//	!yaml mytext - Is replaced by the result of parsing `mytext` as YAML.
//
// Tag names can be abbreviated. For example these are all valid tags:
//
//	!f /my/file.txt - Replaced by the content of the `/my.file.txt` file.
//	!s myscript - Replaced by the result of execution the `myscript` script.
//	!v MYVARIABLE - Replaced by the content of the environment variablel `MYVARIABLE`.
//
// By default the tags replace the value of node they are applied to with a string. This will not
// work for fields that are declared of other types in the configuration struct. In those cases it
// is possible to add a suffix to the tag to indicate the type of the replacmenet. For example:
//
//	# A configuration with an integer loaded from an environment variable and a boolean
//	# loaded from a file:
//	myid: !variable/integer MYID
//	myenabled: !file/boolean /my/enabled.txt
//
// This can be used with the following Go code:
//
//	type MyConfig struct {
//		MyId      int  `yaml:"myid"`
//		MyEnabled bool `yaml:"myenabled"`
//	}
//	var myCfg MyConfig
//	err = cfg.Populate(&myCfg)
//	if err != nil {
//		...
//	}
//
// Tags can be chained. For example to read a value from a file and trim white space from
// the content:
//
//	mypassword: !file/trim mypassword.txt
//
// When multiple sources are configured (calling the Load method multiple times) they will all
// be merged, and sources loaded later sources will override sources loaded earlier.
package configuration
