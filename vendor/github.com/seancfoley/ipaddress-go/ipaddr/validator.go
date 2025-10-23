//
// Copyright 2020-2022 Sean C Foley
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
//

package ipaddr

import (
	"github.com/seancfoley/ipaddress-go/ipaddr/addrerr"
	"github.com/seancfoley/ipaddress-go/ipaddr/addrstrparam"
)

const (
	SmtpIPv6Identifier = "IPv6:"
	IPvFuture          = 'v'
)

// Interface for validation and parsing of host identifier strings
type hostIdentifierStringValidator interface {
	validateHostName(fromHost *HostName, validationOptions addrstrparam.HostNameParams) (*parsedHost, addrerr.HostNameError)

	validateIPAddressStr(fromString *IPAddressString, validationOptions addrstrparam.IPAddressStringParams) (ipAddressProvider, addrerr.AddressStringError)

	validateMACAddressStr(fromString *MACAddressString, validationOptions addrstrparam.MACAddressStringParams) (macAddressProvider, addrerr.AddressStringError)

	validatePrefixLenStr(fullAddr string, version IPVersion) (PrefixLen, addrerr.AddressStringError)
}

var _ hostIdentifierStringValidator = strValidator{}
