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

// This file contains information about the tool.

package info

const DefaultVersion = "1.2.46"

// Build contains the short Git SHA of the CLI at the point it was build. Set via `-ldflags` at build time
var Build = "local"

const DefaultUserAgent = "ROSACLI"
