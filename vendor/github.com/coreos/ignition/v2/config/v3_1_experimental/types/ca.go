// Copyright 2018 CoreOS, Inc.
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

	"github.com/coreos/vcontext/path"
	"github.com/coreos/vcontext/report"

	"github.com/coreos/ignition/v2/config/shared/errors"
)

func (c CaReference) Key() string {
	return c.Source
}

func (ca CaReference) Validate(c path.ContextPath) (r report.Report) {
	r.AddOnError(c.Append("source"), validateURL(ca.Source))
	r.AddOnError(c.Append("httpHeaders"), ca.validateSchemeForHTTPHeaders())
	return
}

func (ca CaReference) validateSchemeForHTTPHeaders() error {
	if len(ca.HTTPHeaders) < 1 {
		return nil
	}

	if ca.Source == "" {
		return errors.ErrInvalidUrl
	}

	u, err := url.Parse(ca.Source)
	if err != nil {
		return errors.ErrInvalidUrl
	}

	switch u.Scheme {
	case "http", "https":
		return nil
	default:
		return errors.ErrUnsupportedSchemeForHTTPHeaders
	}
}
