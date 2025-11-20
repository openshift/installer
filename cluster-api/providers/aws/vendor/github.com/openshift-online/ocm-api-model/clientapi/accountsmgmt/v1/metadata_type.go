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

// IMPORTANT: This file has been generated automatically, refrain from modifying it manually as all
// your changes will be lost when the file is generated again.

package v1 // github.com/openshift-online/ocm-api-model/clientapi/accountsmgmt/v1

// Metadata contains the version metadata.
type Metadata struct {
	fieldSet_     []bool
	serverVersion string
}

// ServerVersion returns the version of the server.
func (m *Metadata) ServerVersion() string {
	if m != nil && len(m.fieldSet_) > 0 && m.fieldSet_[0] {
		return m.serverVersion
	}
	return ""
}

// GetServerVersion returns the value of the server version and a flag indicating if
// the attribute has a value.
func (m *Metadata) GetServerVersion() (value string, ok bool) {
	ok = m != nil && len(m.fieldSet_) > 0 && m.fieldSet_[0]
	if ok {
		value = m.serverVersion
	}
	return
}
