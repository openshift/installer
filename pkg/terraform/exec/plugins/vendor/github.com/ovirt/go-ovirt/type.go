//
// Copyright (c) 2017 Joey <majunjiev@gmail.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package ovirtsdk

type Href interface {
	Href() (string, bool)
}

// Link represents struct of href and rel attributes
type Link struct {
	href *string
	rel  *string
}

// Struct represents the base for all struts defined in types.go
type Struct struct {
	href *string
}

func (p *Struct) Href() (string, bool) {
	if p.href != nil {
		return *p.href, true
	}
	return "", false
}

func (p *Struct) MustHref() string {
	if p.href == nil {
		panic("href attribute must exist")
	}
	return *p.href
}

func (p *Struct) SetHref(attr string) {
	p.href = &attr
}
