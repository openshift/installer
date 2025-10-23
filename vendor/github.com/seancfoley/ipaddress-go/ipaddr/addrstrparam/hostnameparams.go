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

package addrstrparam

// CopyHostNameParams produces an immutable copy of the original HostNameParams.
// Copying a HostNameParams created by a HostNameParamsBuilder is unnecessary since it is already immutable.
func CopyHostNameParams(orig HostNameParams) HostNameParams {
	if p, ok := orig.(*hostNameParameters); ok {
		return p
	}
	return new(HostNameParamsBuilder).Set(orig).ToParams()
}

// HostNameParams provides parameters for parsing host name strings.
//
// This allows you to control the validation performed by HostName.
//
// HostName uses a default permissive HostNameParams object when you do not specify one.
//
// If you wish to use parameters different from the default, then use this interface.  Immutable instances can be constructed with HostNameParamsBuilder.
type HostNameParams interface {
	// AllowsEmpty determines if an empty host string is considered valid.
	// The parser will first parse as an empty address, if allowed by the nested IPAddressStringParams.
	// Otherwise, it will be considered an empty host if this returns true, or an invalid host if it returns false.
	AllowsEmpty() bool

	// GetPreferredVersion indicates the version to prefer when resolving host names.
	GetPreferredVersion() IPVersion

	// AllowsBracketedIPv4 allows bracketed IPv4 addresses like "[1.2.3.4]".
	AllowsBracketedIPv4() bool

	// AllowsBracketedIPv6 allows bracketed IPv6 addresses like "[1::2]".
	AllowsBracketedIPv6() bool

	// NormalizesToLowercase indicates whether to normalize the host name to lowercase characters when parsing.
	NormalizesToLowercase() bool

	// AllowsIPAddress allows a host name to specify an IP address or subnet.
	AllowsIPAddress() bool

	// AllowsPort allows a host name to specify a port.
	AllowsPort() bool

	// AllowsService allows a host name to specify a service, which typically maps to a port.
	AllowsService() bool

	// ExpectsPort indicates whether a port should be inferred from a host like 1:2:3:4::80 that is ambiguous if a port might have been appended.
	// The final segment would normally be considered part of the address, but can be interpreted as a port instead.
	ExpectsPort() bool

	// GetIPAddressParams returns the parameters that apply specifically to IP addresses and subnets, whenever a host name specifies an IP addresses or subnet.
	GetIPAddressParams() IPAddressStringParams
}

// hostNameParameters has parameters for parsing host name strings.
// They are immutable and can be constructed using an HostNameParamsBuilder.
type hostNameParameters struct {
	ipParams ipAddressStringParameters

	preferredVersion IPVersion

	noEmpty, noBracketedIPv4, noBracketedIPv6,
	noNormalizeToLower, noIPAddress, noPort, noService, expectPort bool
}

// AllowsEmpty determines if an empty host string is considered valid.
// The parser will first parse as an empty address, if allowed by the nested IPAddressStringParams.
// Otherwise, it will be considered an empty host if this returns true, or an invalid host if it returns false.
func (params *hostNameParameters) AllowsEmpty() bool {
	return !params.noEmpty
}

// GetPreferredVersion indicates the version to prefer when resolving host names.
func (params *hostNameParameters) GetPreferredVersion() IPVersion {
	return params.preferredVersion
}

// AllowsBracketedIPv4 allows bracketed IPv4 addresses like "[1.2.3.4]".
func (params *hostNameParameters) AllowsBracketedIPv4() bool {
	return !params.noBracketedIPv4
}

// AllowsBracketedIPv6 allows bracketed IPv6 addresses like "[1::2]".
func (params *hostNameParameters) AllowsBracketedIPv6() bool {
	return !params.noBracketedIPv6
}

// NormalizesToLowercase indicates whether to normalize the host name to lowercase characters when parsing.
func (params *hostNameParameters) NormalizesToLowercase() bool {
	return !params.noNormalizeToLower
}

// AllowsIPAddress allows a host name to specify an IP address or subnet.
func (params *hostNameParameters) AllowsIPAddress() bool {
	return !params.noIPAddress
}

// AllowsPort allows a host name to specify a port.
func (params *hostNameParameters) AllowsPort() bool {
	return !params.noPort
}

// AllowsService allows a host name to specify a service, which typically maps to a port.
func (params *hostNameParameters) AllowsService() bool {
	return !params.noService
}

// ExpectsPort indicates whether a port should be inferred from a host like 1:2:3:4::80 that is ambiguous if a port might have been appended.
// The final segment would normally be considered part of the address, but can be interpreted as a port instead.
func (params *hostNameParameters) ExpectsPort() bool {
	return params.expectPort
}

// GetIPAddressParams returns the parameters that apply specifically to IP addresses and subnets, whenever a host name specifies an IP addresses or subnet.
func (params *hostNameParameters) GetIPAddressParams() IPAddressStringParams {
	return &params.ipParams
}

// HostNameParamsBuilder builds an immutable HostNameParams for controlling parsing of host names.
type HostNameParamsBuilder struct {
	hostNameParameters

	ipAddressBuilder IPAddressStringParamsBuilder
}

// ToParams returns an immutable HostNameParams instance built by this builder.
func (builder *HostNameParamsBuilder) ToParams() HostNameParams {
	// We do not return a pointer to builder.hostNameParameters because that would make it possible to change params
	// by continuing to use the same builder,
	// and we want immutable objects for concurrency-safety,
	// so we cannot allow it
	result := builder.hostNameParameters
	result.ipParams = *builder.ipAddressBuilder.ToParams().(*ipAddressStringParameters)
	return &result
}

// GetIPAddressParamsBuilder returns a builder that builds the IPAddressStringParams for the HostNameParams being built by this builder.
func (builder *HostNameParamsBuilder) GetIPAddressParamsBuilder() (result *IPAddressStringParamsBuilder) {
	result = &builder.ipAddressBuilder
	result.parent = builder
	return
}

// Set populates this builder with the values from the given HostNameParams.
func (builder *HostNameParamsBuilder) Set(params HostNameParams) *HostNameParamsBuilder {
	if p, ok := params.(*hostNameParameters); ok {
		builder.hostNameParameters = *p
	} else {
		builder.hostNameParameters = hostNameParameters{
			preferredVersion:   params.GetPreferredVersion(),
			noEmpty:            !params.AllowsEmpty(),
			noBracketedIPv4:    !params.AllowsBracketedIPv4(),
			noBracketedIPv6:    !params.AllowsBracketedIPv6(),
			noNormalizeToLower: !params.NormalizesToLowercase(),
			noIPAddress:        !params.AllowsIPAddress(),
			noPort:             !params.AllowsPort(),
			noService:          !params.AllowsService(),
			expectPort:         params.ExpectsPort(),
		}
	}
	builder.SetIPAddressParams(params.GetIPAddressParams())
	return builder
}

// SetIPAddressParams populates this builder with the values from the given IPAddressStringParams.
func (builder *HostNameParamsBuilder) SetIPAddressParams(params IPAddressStringParams) *HostNameParamsBuilder {
	//builder.ipAddressBuilder = *ToIPAddressStringParamsBuilder(params)
	builder.ipAddressBuilder.Set(params)
	return builder
}

// AllowEmpty dictates whether an empty host string is considered valid.
// The parser will first parse as an empty address, if allowed by the nested IPAddressStringParams.
// Otherwise, this setting dictates whether it will be considered an invalid host.
func (builder *HostNameParamsBuilder) AllowEmpty(allow bool) *HostNameParamsBuilder {
	builder.hostNameParameters.noEmpty = !allow
	return builder
}

// SetPreferredVersion dictates the version to prefer when resolving host names.
func (builder *HostNameParamsBuilder) SetPreferredVersion(version IPVersion) *HostNameParamsBuilder {
	builder.hostNameParameters.preferredVersion = version
	return builder
}

// AllowBracketedIPv4 dictates whether to allow bracketed IPv4 addresses like "[1.2.3.4]".
func (builder *HostNameParamsBuilder) AllowBracketedIPv4(allow bool) *HostNameParamsBuilder {
	builder.hostNameParameters.noBracketedIPv4 = !allow
	return builder
}

// AllowBracketedIPv6 dictates whether to allow bracketed IPv6 addresses like "[1::2]".
func (builder *HostNameParamsBuilder) AllowBracketedIPv6(allow bool) *HostNameParamsBuilder {
	builder.hostNameParameters.noBracketedIPv6 = !allow
	return builder
}

// NormalizeToLowercase dictates whether to normalize the host name to lowercase characters when parsing.
func (builder *HostNameParamsBuilder) NormalizeToLowercase(allow bool) *HostNameParamsBuilder {
	builder.hostNameParameters.noNormalizeToLower = !allow
	return builder
}

// AllowIPAddress dictates whether to allow a host name to specify an IP address or subnet.
func (builder *HostNameParamsBuilder) AllowIPAddress(allow bool) *HostNameParamsBuilder {
	builder.hostNameParameters.noIPAddress = !allow
	return builder
}

// AllowPort dictates whether to allow a host name to specify a port.
func (builder *HostNameParamsBuilder) AllowPort(allow bool) *HostNameParamsBuilder {
	builder.hostNameParameters.noPort = !allow
	return builder
}

// AllowService dictates whether to allow a host name to specify a service, which typically maps to a port.
func (builder *HostNameParamsBuilder) AllowService(allow bool) *HostNameParamsBuilder {
	builder.hostNameParameters.noService = !allow
	return builder
}

// ExpectPort dictates whether a port should be inferred from a host like 1:2:3:4::80 that is ambiguous if a port might have been appended.
// The final segment would normally be considered part of the address, but can be interpreted as a port instead.
func (builder *HostNameParamsBuilder) ExpectPort(expect bool) *HostNameParamsBuilder {
	builder.hostNameParameters.expectPort = expect
	return builder
}
