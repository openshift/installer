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

// CopyIPAddressStringParams produces an immutable copy of the original IPAddressStringParams.
// Copying an IPAddressStringParams created by an IPAddressStringParamsBuilder is unnecessary since it is already immutable.
func CopyIPAddressStringParams(orig IPAddressStringParams) IPAddressStringParams {
	if p, ok := orig.(*ipAddressStringParameters); ok {
		return p
	}
	return new(IPAddressStringParamsBuilder).Set(orig).ToParams()
}

// IPAddressStringParams provides parameters for parsing IP address strings,
// indicating what to allow, what to disallow, and other options.
//
// This allows you to control the validation performed by IPAddressString.
//
// IPAddressString uses a default permissive IPAddressStringParams instance when you do not specify one.
//
// If you wish to use parameters different from the default, then use this interface.  Immutable instances can be constructed with IPAddressStringParamsBuilder.
type IPAddressStringParams interface {
	AddressStringParams

	// AllowsPrefix indicates whether addresses with prefix length like 1.2.0.0/16 are allowed.
	AllowsPrefix() bool

	// EmptyStrParsedAs determines how a zero-length empty string is translated to an address.
	// If the option is ZeroAddressOption or LoopbackOption, then if defers to GetPreferredVersion() for the version.
	EmptyStrParsedAs() EmptyStrOption

	// AllStrParsedAs determines how the "all" string "*" is translated to addresses.
	// If the option is AllPreferredIPVersion, then it defers to GetPreferredVersion() for the version.
	AllStrParsedAs() AllStrOption

	// AllowsMask allows masks to follow valid addresses, such as 1.2.3.4/255.255.0.0 which has the mask 255.255.0.0
	// If the mask is the mask for a network prefix length, this is interpreted as the subnet for that network prefix length.
	// Otherwise the address is simply masked by the mask.
	// For instance, 1.2.3.4/255.0.255.0 is 1.0.3.0, while 1.2.3.4/255.255.0.0 is 1.2.0.0/16.
	AllowsMask() bool

	// GetPreferredVersion indicates the version to use for ambiguous addresses strings,
	// like prefix lengths less than 32 bits which are translated to masks,
	// the "all" address or the "empty" address.
	// The default is IPv6.
	//
	// If either of AllowsIPv4() or AllowsIPv6() returns false, then those settings take precedence over this setting.
	GetPreferredVersion() IPVersion

	// AllowsIPv4 allows IPv4 addresses and subnets.
	AllowsIPv4() bool

	// AllowsIPv6 allows IPv6 addresses and subnets.
	AllowsIPv6() bool

	// GetIPv4Params returns the parameters that apply specifically to IPv4 addresses and subnets.
	GetIPv4Params() IPv4AddressStringParams

	// GetIPv6Params returns the parameters that apply specifically to IPv6 addresses and subnets.
	GetIPv6Params() IPv6AddressStringParams
}

// EmptyStrOption is an option indicating how to translate an empty address string to an address.
type EmptyStrOption string

const (
	// NoAddressOption indicates that empty strings are not translated to addresses.
	NoAddressOption EmptyStrOption = "none"

	// ZeroAddressOption is the default, which means empty strings are translated to zero addresses.
	ZeroAddressOption EmptyStrOption = ""

	// LoopbackOption indicates empty strings are translated to loopback addresses.
	LoopbackOption EmptyStrOption = "loopback"
)

// AllStrOption is an option indicating how to translate an all address string, such as "*", to an address.
type AllStrOption string

const (
	// AllAddresses is the default, indicating the all address string refers to all addresses of all IP versions.
	AllAddresses AllStrOption = "" // the default

	// AllPreferredIPVersion indicates the all address string refers to all addresses of the preferred IP version.
	AllPreferredIPVersion AllStrOption = "preferred"
)

var _ IPAddressStringParams = &ipAddressStringParameters{}

// IPv4AddressStringParams provides parameters specific to IPv4 addresses and subnets
type IPv4AddressStringParams interface {
	IPAddressStringFormatParams

	// Allows_inet_aton_hex allows IPv4 inet_aton hexadecimal format "0xa.0xb.0xc.0cd".
	Allows_inet_aton_hex() bool

	// Allows_inet_aton_octal allows IPv4 inet_aton octal format, "04.05.06.07" being an example.
	// Can be overridden by allowLeadingZeros
	Allows_inet_aton_octal() bool

	// Allows_inet_aton_joinedSegments allows IPv4 joined segments like "1.2.3", "1.2", or just "1".
	//
	// For the case of just 1 segment, the behaviour is controlled by allowSingleSegment.
	Allows_inet_aton_joinedSegments() bool

	// Allows_inet_aton_single_segment_mask indicates whether you allow a mask that looks like a prefix length when you allow IPv4 joined segments: "1.2.3.5/255".
	Allows_inet_aton_single_segment_mask() bool

	// Allows single-segment inet_aton strings to have extraneous digits.
	// This allows up to 31 digits when parsing for both IPv4 and IPv6.
	// This allows an unlimited number of digits when parsing for just IPv4 (ie {@link IPAddressStringParameters#allowIPv6} is false).
	//
	// Digits that go beyond 32 bits are essentially ignored.
	// The number of digits before exceeding 32 bits depends on the radix.
	// The value of the most significant digit before exceeding 32 bits depends on the radix.
	//
	// The resulting address is the modulus of the address with the 32-bit unsigned int maximum value,
	// or equivalently the truncation of the address to 32 bits.
	Allows_inet_aton_extraneous_digits() bool

	// Allows_inet_aton_leading_zeros allows IPv4 inet_aton hexadecimal or octal to have leading zeros, such as in the first two segments of "0x0a.00b.c.d".
	// The first 0 is not considered a leading zero, it either denotes octal or hex depending on whether it is followed by an 'x'.
	// Zeros that appear afterwards are inet_aton leading zeros.
	Allows_inet_aton_leading_zeros() bool
}

var _ IPv4AddressStringParams = &ipv4AddressStringParameters{}

// IPv6AddressStringParams provides parameters specific to IPv6 addresses and subnets.
type IPv6AddressStringParams interface {
	IPAddressStringFormatParams

	// AllowsMixed allows mixed-in embedded IPv4 like "a:b:c:d:e:f:1.2.3.4".
	AllowsMixed() bool

	// AllowsZone allows zones like "a:b:c:d:e:f:a:b%zone".
	AllowsZone() bool

	// AllowsEmptyZone allows the zone character % with no following zone.
	AllowsEmptyZone() bool

	// AllowsBase85 allows IPv6 single-segment base 85 addresses.
	AllowsBase85() bool

	// GetMixedParams provides the IP parameters that for parsing the embedded IPv4 section of a mixed IPv6/v4 address, if AllowsMixed is true.
	GetMixedParams() IPAddressStringParams

	// GetEmbeddedIPv4AddressParams returns the IPv4 parameters for parsing the embedded IPv4 section of a mixed IPv6/v4 address.
	GetEmbeddedIPv4AddressParams() IPv4AddressStringParams
}

var _ IPv6AddressStringParams = &ipv6AddressStringParameters{}

// IPAddressStringFormatParams provides format parameters that apply to all IP addresses, but can be different for IPv4 or IPv6,
// allowing for cases where you may wish to allow something for one version but not the same for the other version.
type IPAddressStringFormatParams interface {
	AddressStringFormatParams

	// AllowsPrefixesBeyondAddressSize allows prefix length values greater than 32 for IPv4 or greater than 128 for IPv6.
	AllowsPrefixesBeyondAddressSize() bool

	// AllowsPrefixLenLeadingZeros allows leading zeros in the prefix length like "1.2.3.4/016".
	AllowsPrefixLenLeadingZeros() bool

	// AllowsBinary allows binary addresses like "11111111.0.1.0" or "1111111111111111::".
	AllowsBinary() bool
}

func init() {
	defaultEmbeddedBuilder.
		AllowEmpty(false).
		AllowPrefix(false).
		AllowMask(false).
		AllowAll(false).
		AllowIPv6(false).
		GetIPv6AddressParamsBuilder().
		AllowZone(true).
		AllowEmptyZone(true)
	defaultEmbeddedParams =
		defaultEmbeddedBuilder.
			ToParams().(*ipAddressStringParameters)
}

var defaultEmbeddedBuilder IPAddressStringParamsBuilder
var defaultEmbeddedParams *ipAddressStringParameters

// ipAddressStringParameters has parameters for parsing IP address strings.
// They are immutable and can be constructed using an IPAddressStringParamsBuilder.
type ipAddressStringParameters struct {
	addressStringParameters
	ipv4Params ipv4AddressStringParameters
	ipv6Params ipv6AddressStringParameters

	emptyStringOption EmptyStrOption
	allStringOption   AllStrOption
	preferredVersion  IPVersion

	//emptyIsNotLoopback,
	//noPrefixOnly,
	noPrefix, noMask, noIPv6, noIPv4 bool
}

// AllowsPrefix indicates whether addresses with prefix length like 1.2.0.0/16 are allowed.
func (params *ipAddressStringParameters) AllowsPrefix() bool {
	return !params.noPrefix
}

// EmptyStrParsedAs determines how a zero-length empty string is translated to an address.
// If the option is ZeroAddressOption or LoopbackOption, then if defers to GetPreferredVersion() for the version.
func (params *ipAddressStringParameters) EmptyStrParsedAs() EmptyStrOption {
	return params.emptyStringOption
}

// AllStrParsedAs determines how the "all" string "*" is translated to addresses.
// If the option is AllPreferredIPVersion, then it defers to GetPreferredVersion() for the version.
func (params *ipAddressStringParameters) AllStrParsedAs() AllStrOption {
	return params.allStringOption
}

// GetPreferredVersion indicates the version to use for ambiguous addresses strings,
// like prefix lengths less than 32 bits which are translated to masks,
// the "all" address or the "empty" address.
// The default is IPv6.
//
// If either of AllowsIPv4() or AllowsIPv6() returns false, then those settings take precedence over this setting.
func (params *ipAddressStringParameters) GetPreferredVersion() IPVersion {
	return params.preferredVersion
}

// AllowsMask allows masks to follow valid addresses, such as 1.2.3.4/255.255.0.0 which has the mask 255.255.0.0
// If the mask is the mask for a network prefix length, this is interpreted as the subnet for that network prefix length.
// Otherwise the address is simply masked by the mask.
// For instance, 1.2.3.4/255.0.255.0 is 1.0.3.0, while 1.2.3.4/255.255.0.0 is 1.2.0.0/16.
func (params *ipAddressStringParameters) AllowsMask() bool {
	return !params.noMask
}

// AllowsIPv4 allows IPv4 addresses and subnets.
func (params *ipAddressStringParameters) AllowsIPv4() bool {
	return !params.noIPv4
}

// AllowsIPv6 allows IPv6 addresses and subnets.
func (params *ipAddressStringParameters) AllowsIPv6() bool {
	return !params.noIPv6
}

// GetIPv4Params returns the parameters that apply specifically to IPv4 addresses and subnets.
func (params *ipAddressStringParameters) GetIPv4Params() IPv4AddressStringParams {
	return &params.ipv4Params
}

// GetIPv6Params returns the parameters that apply specifically to IPv6 addresses and subnets.
func (params *ipAddressStringParameters) GetIPv6Params() IPv6AddressStringParams {
	return &params.ipv6Params
}

// IPAddressStringParamsBuilder builds an immutable IPAddressStringParameters for controlling parsing of IP address strings.
type IPAddressStringParamsBuilder struct {
	params ipAddressStringParameters
	AddressStringParamsBuilder
	ipv4Builder IPv4AddressStringParamsBuilder
	ipv6Builder IPv6AddressStringParamsBuilder

	parent *HostNameParamsBuilder
}

// GetParentBuilder returns the original HostNameParamsBuilder builder that this was obtained from, if this builder was obtained from a HostNameParamsBuilder.
func (builder *IPAddressStringParamsBuilder) GetParentBuilder() *HostNameParamsBuilder {
	return builder.parent
}

// ToParams returns an immutable IPAddressStringParams instance built by this builder.
func (builder *IPAddressStringParamsBuilder) ToParams() IPAddressStringParams {
	// We do not return a pointer to builder.params because that would make it possible to change a ipAddressStringParameters
	// by continuing to use the same builder,
	// and we want immutable objects for concurrency-safety,
	// so we cannot allow it
	result := builder.params
	result.addressStringParameters = *builder.AddressStringParamsBuilder.ToParams().(*addressStringParameters)
	result.ipv4Params = *builder.ipv4Builder.ToParams().(*ipv4AddressStringParameters)
	result.ipv6Params = *builder.ipv6Builder.ToParams().(*ipv6AddressStringParameters)
	return &result
}

// GetIPv6AddressParamsBuilder returns a builder that builds the IPv6AddressStringParams for the IPAddressStringParams being built by this builder.
func (builder *IPAddressStringParamsBuilder) GetIPv6AddressParamsBuilder() (result *IPv6AddressStringParamsBuilder) {
	result = &builder.ipv6Builder
	result.parent = builder
	return
}

// GetIPv4AddressParamsBuilder returns a builder that builds the IPv4AddressStringParams for the IPAddressStringParams being built by this builder.
func (builder *IPAddressStringParamsBuilder) GetIPv4AddressParamsBuilder() (result *IPv4AddressStringParamsBuilder) {
	result = &builder.ipv4Builder
	result.parent = builder
	return
}

// Set populates this builder with the values from the given IPAddressStringParams.
func (builder *IPAddressStringParamsBuilder) Set(params IPAddressStringParams) *IPAddressStringParamsBuilder {
	return builder.set(params, false)
}

func (builder *IPAddressStringParamsBuilder) set(params IPAddressStringParams, isMixed bool) *IPAddressStringParamsBuilder {
	if p, ok := params.(*ipAddressStringParameters); ok {
		builder.params = *p
	} else {
		builder.params = ipAddressStringParameters{
			preferredVersion:  params.GetPreferredVersion(),
			emptyStringOption: params.EmptyStrParsedAs(),
			allStringOption:   params.AllStrParsedAs(),
			noPrefix:          !params.AllowsPrefix(),
			noMask:            !params.AllowsMask(),
			noIPv6:            !params.AllowsIPv6(),
			noIPv4:            !params.AllowsIPv4(),
		}
	}
	builder.AddressStringParamsBuilder.set(params)
	builder.ipv4Builder.Set(params.GetIPv4Params())
	builder.ipv6Builder.set(params.GetIPv6Params(), isMixed)
	return builder
}

// AllowEmpty dictates whether to allow empty zero-length address strings.
func (builder *IPAddressStringParamsBuilder) AllowEmpty(allow bool) *IPAddressStringParamsBuilder {
	builder.allowEmpty(allow)
	return builder
}

// AllowSingleSegment dictates whether to allow an address to be specified as a single value, eg "ffffffff", without the standard use of segments like "1.2.3.4" or "1:2:4:3:5:6:7:8".
func (builder *IPAddressStringParamsBuilder) AllowSingleSegment(allow bool) *IPAddressStringParamsBuilder {
	builder.allowSingleSegment(allow)
	return builder
}

// AllowAll dictates whether to alloww the string of just the wildcard "*" to denote all addresses of all version.
// If false, then for IP addresses we check the preferred version with GetPreferredVersion, and then check AllowsWildcardedSeparator,
// to determine if the string represents all addresses of that version.
func (builder *IPAddressStringParamsBuilder) AllowAll(allow bool) *IPAddressStringParamsBuilder {
	builder.allowAll(allow)
	return builder
}

// ParseEmptyStrAs dictates how a zero-length empty string is translated to an address.
// If the option is ZeroAddressOption or LoopbackOption, then if defers to GetPreferredVersion for the version.
func (builder *IPAddressStringParamsBuilder) ParseEmptyStrAs(option EmptyStrOption) *IPAddressStringParamsBuilder {
	builder.params.emptyStringOption = option
	builder.AllowEmpty(true)
	return builder
}

// ParseAllStrAs dictates how the "all" string "*" is translated to addresses.
// If the option is AllPreferredIPVersion, then it defers to GetPreferredVersion for the version.
func (builder *IPAddressStringParamsBuilder) ParseAllStrAs(option AllStrOption) *IPAddressStringParamsBuilder {
	builder.params.allStringOption = option
	return builder
}

// SetPreferredVersion dictates the version to use for ambiguous addresses strings,
// like prefix lengths less than 32 bits which are translated to masks,
// the "all" address or the "empty" address.
// The default is IPv6.
//
// If either of AllowsIPv4 or AllowsIPv6 returns false, then those settings take precedence over this setting.
func (builder *IPAddressStringParamsBuilder) SetPreferredVersion(version IPVersion) *IPAddressStringParamsBuilder {
	builder.params.preferredVersion = version
	return builder
}

// AllowPrefix dictates whether to allow addresses with prefix length like "1.2.0.0/16" are allowed.
func (builder *IPAddressStringParamsBuilder) AllowPrefix(allow bool) *IPAddressStringParamsBuilder {
	builder.params.noPrefix = !allow
	return builder
}

// AllowMask dictates whether to allow masks to follow valid addresses, such as "1.2.3.4/255.255.0.0" which has the mask "255.255.0.0".
// If the mask is the mask for a network prefix length, this is interpreted as the subnet for that network prefix length.
// Otherwise the address is simply masked by the mask.
// For instance, "1.2.3.4/255.0.255.0" is "1.0.3.0", while "1.2.3.4/255.255.0.0" is "1.2.0.0/16".
func (builder *IPAddressStringParamsBuilder) AllowMask(allow bool) *IPAddressStringParamsBuilder {
	builder.params.noMask = !allow
	return builder
}

// AllowIPv4 dictates whether to allow IPv4 addresses and subnets
func (builder *IPAddressStringParamsBuilder) AllowIPv4(allow bool) *IPAddressStringParamsBuilder {
	builder.params.noIPv4 = !allow
	return builder
}

// AllowIPv6 dictates whether to allow IPv6 addresses and subnets
func (builder *IPAddressStringParamsBuilder) AllowIPv6(allow bool) *IPAddressStringParamsBuilder {
	builder.params.noIPv6 = !allow
	return builder
}

// AllowWildcardedSeparator dictates whether the wildcard '*' or '%' can replace the segment separators '.' and ':'.
// If so, then you can write addresses like *.* or *:*
func (builder *IPAddressStringParamsBuilder) AllowWildcardedSeparator(allow bool) *IPAddressStringParamsBuilder {
	builder.GetIPv4AddressParamsBuilder().AllowWildcardedSeparator(allow)
	builder.GetIPv6AddressParamsBuilder().AllowWildcardedSeparator(allow)
	return builder
}

// SetRangeParams populates this builder with the values from the given RangeParams.
func (builder *IPAddressStringParamsBuilder) SetRangeParams(rangeParams RangeParams) *IPAddressStringParamsBuilder {
	builder.GetIPv4AddressParamsBuilder().SetRangeParams(rangeParams)
	builder.GetIPv6AddressParamsBuilder().SetRangeParams(rangeParams)
	return builder
}

// Allow_inet_aton dictates whether to allow any IPv4 inet_aton format, whether hex, octal, or joined segments.
func (builder *IPAddressStringParamsBuilder) Allow_inet_aton(allow bool) *IPAddressStringParamsBuilder {
	builder.GetIPv4AddressParamsBuilder().Allow_inet_aton(allow)
	builder.GetIPv6AddressParamsBuilder().Allow_mixed_inet_aton(allow)
	return builder
}

type ipAddressStringFormatParameters struct {
	addressStringFormatParameters

	allowPrefixesBeyondAddrSize,
	noPrefixLengthLeadingZeros,
	noBinary bool
}

// AllowsPrefixesBeyondAddressSize allows prefix length values greater than 32 for IPv4 or greater than 128 for IPv6.
func (params *ipAddressStringFormatParameters) AllowsPrefixesBeyondAddressSize() bool {
	return params.allowPrefixesBeyondAddrSize
}

// AllowsPrefixLenLeadingZeros allows leading zeros in the prefix length like "1.2.3.4/016".
func (params *ipAddressStringFormatParameters) AllowsPrefixLenLeadingZeros() bool {
	return !params.noPrefixLengthLeadingZeros
}

// AllowsBinary allows binary addresses like "11111111.0.1.0" or "1111111111111111::".
func (params *ipAddressStringFormatParameters) AllowsBinary() bool {
	return !params.noBinary
}

// IPAddressStringFormatParamsBuilder builds an immutable IPAddressStringFormatParams for controlling parsing of IP address strings.
type IPAddressStringFormatParamsBuilder struct {
	AddressStringFormatParamsBuilder

	ipParams ipAddressStringFormatParameters

	parent *IPAddressStringParamsBuilder
}

// GetParentBuilder returns the original IPAddressStringParamsBuilder builder that this was obtained from, if this builder was obtained from a IPAddressStringParamsBuilder.
func (builder *IPAddressStringFormatParamsBuilder) GetParentBuilder() *IPAddressStringParamsBuilder {
	return builder.parent
}

// ToParams returns an immutable IPAddressStringFormatParams instance built by this builder
func (builder *IPAddressStringFormatParamsBuilder) ToParams() IPAddressStringFormatParams {
	result := &builder.ipParams
	result.addressStringFormatParameters = *builder.AddressStringFormatParamsBuilder.ToParams().(*addressStringFormatParameters)
	return result
}

func (builder *IPAddressStringFormatParamsBuilder) set(params IPAddressStringFormatParams) {
	if p, ok := params.(*ipAddressStringFormatParameters); ok {
		builder.ipParams = *p
	} else {
		builder.ipParams = ipAddressStringFormatParameters{
			allowPrefixesBeyondAddrSize: params.AllowsPrefixesBeyondAddressSize(),
			noPrefixLengthLeadingZeros:  !params.AllowsPrefixLenLeadingZeros(),
			noBinary:                    !params.AllowsBinary(),
		}
	}
	builder.AddressStringFormatParamsBuilder.set(params)
}

// AllowsPrefixesBeyondAddressSize allows prefix length values greater than 32 for IPv4 or greater than 128 for IPv6.
func (builder *IPAddressStringFormatParamsBuilder) AllowsPrefixesBeyondAddressSize() bool {
	return builder.ipParams.AllowsPrefixesBeyondAddressSize()
}

// AllowsPrefixLenLeadingZeros allows leading zeros in the prefix length like "1.2.3.4/016".
func (builder *IPAddressStringFormatParamsBuilder) AllowsPrefixLenLeadingZeros() bool {
	return builder.ipParams.AllowsPrefixLenLeadingZeros()
}

// AllowsBinary allows binary addresses like 11111111.0.1.0 or 1111111111111111::
func (builder *IPAddressStringFormatParamsBuilder) AllowsBinary() bool {
	return builder.ipParams.AllowsBinary()
}

func (builder *IPAddressStringFormatParamsBuilder) allowBinary(allow bool) {
	builder.ipParams.noBinary = !allow
}

func (builder *IPAddressStringFormatParamsBuilder) allowPrefixesBeyondAddressSize(allow bool) {
	builder.ipParams.allowPrefixesBeyondAddrSize = allow
}

func (builder *IPAddressStringFormatParamsBuilder) allowPrefixLengthLeadingZeros(allow bool) {
	builder.ipParams.noPrefixLengthLeadingZeros = !allow
}

type ipv6AddressStringParameters struct {
	ipAddressStringFormatParameters

	noMixed, noZone, noBase85, noEmptyZone bool

	embeddedParams *ipAddressStringParameters
}

// AllowsMixed allows mixed-in embedded IPv4 like "a:b:c:d:e:f:1.2.3.4".
func (params *ipv6AddressStringParameters) AllowsMixed() bool {
	return !params.noMixed
}

// AllowsZone allows zones like "a:b:c:d:e:f:a:b%zone".
func (params *ipv6AddressStringParameters) AllowsZone() bool {
	return !params.noZone
}

// AllowsEmptyZone allows the zone character % with no following zone'
func (params *ipv6AddressStringParameters) AllowsEmptyZone() bool {
	return !params.noEmptyZone
}

// AllowsBase85 allows IPv6 single-segment base 85 addresses'
func (params *ipv6AddressStringParameters) AllowsBase85() bool {
	return !params.noBase85
}

// GetMixedParams provides the parameters that for parsing the embedded IPv4 section of a mixed IPv6/v4 address, if AllowsMixed is true'
func (params *ipv6AddressStringParameters) GetMixedParams() IPAddressStringParams {
	result := params.embeddedParams
	if result == nil {
		result = defaultEmbeddedParams
	}
	return result
}

// GetEmbeddedIPv4AddressParams returns the IPv4 parameters for parsing the embedded IPv4 section of a mixed IPv6/v4 address'
func (params *ipv6AddressStringParameters) GetEmbeddedIPv4AddressParams() IPv4AddressStringParams {
	return params.embeddedParams.GetIPv4Params()
}

// IPv6AddressStringParamsBuilder builds an immutable IPv6AddressStringParams for controlling parsing of IPv6 address strings'
type IPv6AddressStringParamsBuilder struct {
	// This is not anonymous since it clashes with IPAddressStringFormatParamsBuilder,
	// both have ipAddressStringFormatParameters and AddressStringFormatParams
	// and thee builder IPAddressStringFormatParamsBuilder takes precedence
	params ipv6AddressStringParameters

	embeddedBuilder *IPAddressStringParamsBuilder

	IPAddressStringFormatParamsBuilder
}

// ToParams returns an immutable IPv6AddressStringParams instance built by this builder'
func (builder *IPv6AddressStringParamsBuilder) ToParams() IPv6AddressStringParams {
	result := &builder.params
	result.ipAddressStringFormatParameters = *builder.IPAddressStringFormatParamsBuilder.ToParams().(*ipAddressStringFormatParameters)
	if emb := builder.embeddedBuilder; emb == nil {
		result.embeddedParams = defaultEmbeddedParams
	} else {
		result.embeddedParams = emb.ToParams().(*ipAddressStringParameters)
	}
	return result
}

// GetRangeParamsBuilder returns a builder that builds the range parameters for these IPv6 address string parameters.
func (builder *IPv6AddressStringParamsBuilder) GetRangeParamsBuilder() *RangeParamsBuilder {
	result := &builder.rangeParamsBuilder
	result.parent = builder
	return result
}

// AllowsMixed allows mixed-in embedded IPv4 like "a:b:c:d:e:f:1.2.3.4".
func (builder *IPv6AddressStringParamsBuilder) AllowsMixed() bool {
	return builder.params.AllowsMixed()
}

// AllowsZone allows zones like "a:b:c:d:e:f:a:b%zone".
func (builder *IPv6AddressStringParamsBuilder) AllowsZone() bool {
	return builder.params.AllowsZone()
}

// AllowsEmptyZone allows the zone character % with no following zone.
func (builder *IPv6AddressStringParamsBuilder) AllowsEmptyZone() bool {
	return builder.params.AllowsEmptyZone()
}

// AllowsBase85 allows IPv6 single-segment base 85 addresses.
func (builder *IPv6AddressStringParamsBuilder) AllowsBase85() bool {
	return builder.params.AllowsBase85()
}

// Set populates this builder with the values from the given IPv6AddressStringParams.
func (builder *IPv6AddressStringParamsBuilder) Set(params IPv6AddressStringParams) *IPv6AddressStringParamsBuilder {
	return builder.set(params, false)
}

func (builder *IPv6AddressStringParamsBuilder) set(params IPv6AddressStringParams, isMixed bool) *IPv6AddressStringParamsBuilder {
	if p, ok := params.(*ipv6AddressStringParameters); ok {
		builder.params = *p
	} else {
		builder.params = ipv6AddressStringParameters{
			noMixed:     !params.AllowsMixed(),
			noZone:      !params.AllowsZone(),
			noEmptyZone: !params.AllowsEmptyZone(),
			noBase85:    !params.AllowsBase85(),
		}
	}
	builder.IPAddressStringFormatParamsBuilder.set(params)
	if !isMixed {
		builder.getEmbeddedIPv4ParametersBuilder().ipv4Builder.Set(params.GetEmbeddedIPv4AddressParams())
	}
	return builder
}

// AllowBase85 dictates whether to allow IPv6 single-segment base 85 addresses.
func (builder *IPv6AddressStringParamsBuilder) AllowBase85(allow bool) *IPv6AddressStringParamsBuilder {
	builder.params.noBase85 = !allow
	return builder
}

// AllowZone dictates whether to allow zones like "a:b:c:d:e:f:a:b%zone".
func (builder *IPv6AddressStringParamsBuilder) AllowZone(allow bool) *IPv6AddressStringParamsBuilder {
	builder.params.noZone = !allow

	//we must decide whether to treat the % character as a zone when parsing the mixed part
	//if considered zone, then the zone character is actually part of the encompassing IPv6 address
	//otherwise, the zone character is an sql wildcard that is part of the mixed address
	//So whether we consider the % character a zone must match the same setting for the encompassing address

	// ipv4Builder can be nil when builder == &defaultEmbeddedBuilder.ipv6Builder, see getEmbeddedIPv4ParametersBuilder()
	if ipv4Builder := builder.getEmbeddedIPv4ParametersBuilder(); ipv4Builder != nil {
		ipv4Builder.GetIPv6AddressParamsBuilder().params.noZone = !allow
	}
	return builder
}

// AllowEmptyZone dictates whether to allow the zone character % with no following zone
func (builder *IPv6AddressStringParamsBuilder) AllowEmptyZone(allow bool) *IPv6AddressStringParamsBuilder {
	builder.params.noEmptyZone = !allow
	if ipv4Builder := builder.getEmbeddedIPv4ParametersBuilder(); ipv4Builder != nil {
		ipv4Builder.GetIPv6AddressParamsBuilder().params.noEmptyZone = !allow
	}
	return builder
}

// AllowMixed dictates whether to allow mixed-in embedded IPv4 like "a:b:c:d:e:f:1.2.3.4".
func (builder *IPv6AddressStringParamsBuilder) AllowMixed(allow bool) *IPv6AddressStringParamsBuilder {
	builder.params.noMixed = !allow
	return builder
}

func (builder *IPv6AddressStringParamsBuilder) getEmbeddedIPv4ParametersBuilder() (result *IPAddressStringParamsBuilder) {
	if builder == &defaultEmbeddedBuilder.ipv6Builder {
		return nil
	}
	if result = builder.embeddedBuilder; result == nil {
		result = &IPAddressStringParamsBuilder{}
		// copy in proper default values for embedded IPv4 addresses, which differ from defaults for typical ipv4AddrType addresses
		*result = defaultEmbeddedBuilder
		builder.embeddedBuilder = result
	}
	result.GetIPv4AddressParamsBuilder().mixedParent = builder
	return
}

// GetEmbeddedIPv4AddressParamsBuilder returns a builder to build the IPv4 parameters that controls parsing of the embedded IPv4 section of a mixed IPv6/v4 address.
func (builder *IPv6AddressStringParamsBuilder) GetEmbeddedIPv4AddressParamsBuilder() (result *IPv4AddressStringParamsBuilder) {
	return builder.getEmbeddedIPv4ParametersBuilder().GetIPv4AddressParamsBuilder()
}

// Allow_mixed_inet_aton dictates whether to allow inet_aton style formats, whether hex, octal, or joined segments, in the embedded IPv4 section of a mixed IPv6/v4 address.
func (builder *IPv6AddressStringParamsBuilder) Allow_mixed_inet_aton(allow bool) *IPv6AddressStringParamsBuilder {
	builder.getEmbeddedIPv4ParametersBuilder().GetIPv4AddressParamsBuilder().Allow_inet_aton(allow)
	if allow { // if we allow inet_aton in the mixed part, then of course that insinuates that we allow the mixed part
		builder.AllowMixed(allow)
	}
	return builder
}

// AllowBinary dictates whether to allow binary addresses like "11111111.0.1.0" or "1111111111111111::".
func (builder *IPv6AddressStringParamsBuilder) AllowBinary(allow bool) *IPv6AddressStringParamsBuilder {
	builder.GetEmbeddedIPv4AddressParamsBuilder().AllowBinary(allow)
	builder.allowBinary(allow)
	return builder
}

// AllowWildcardedSeparator dictates whether the wildcard '*' or '%' can replace the segment separators '.' and ':'.
// If so, then you can write addresses like *.* or *:*
func (builder *IPv6AddressStringParamsBuilder) AllowWildcardedSeparator(allow bool) *IPv6AddressStringParamsBuilder {
	builder.GetEmbeddedIPv4AddressParamsBuilder().AllowWildcardedSeparator(allow)
	builder.allowWildcardedSeparator(allow)
	return builder
}

// AllowLeadingZeros dictates whether to allow addresses with segments that have leasing zeros like "001.2.3.004" or "1:000a::".
// For IPV4, this option overrides inet_aton octal.
//
// Single segment addresses that must have the requisite length to be parsed are not affected by this flag.
func (builder *IPv6AddressStringParamsBuilder) AllowLeadingZeros(allow bool) *IPv6AddressStringParamsBuilder {
	builder.GetEmbeddedIPv4AddressParamsBuilder().allowLeadingZeros(allow)
	builder.allowLeadingZeros(allow)
	return builder
}

// AllowUnlimitedLeadingZeros dictates whether to allow leading zeros that extend segments
// beyond the usual segment length, which is 3 for IPv4 dotted-decimal and 4 for IPv6.
// However, this only takes effect if leading zeros are allowed, which is when
// AllowsLeadingZeros is true or the address is IPv4 and Allows_inet_aton_octal is true.
//
// For example, this determines whether you allow 0001.0002.0003.0004
func (builder *IPv6AddressStringParamsBuilder) AllowUnlimitedLeadingZeros(allow bool) *IPv6AddressStringParamsBuilder {
	builder.GetEmbeddedIPv4AddressParamsBuilder().AllowUnlimitedLeadingZeros(allow)
	builder.allowUnlimitedLeadingZeros(allow)
	return builder
}

// SetRangeParams populates this builder with the values from the given RangeParams.
func (builder *IPv6AddressStringParamsBuilder) SetRangeParams(rangeParams RangeParams) *IPv6AddressStringParamsBuilder {
	builder.GetEmbeddedIPv4AddressParamsBuilder().SetRangeParams(rangeParams)
	builder.setRangeParameters(rangeParams)
	return builder
}

// AllowPrefixesBeyondAddressSize dictates whether to allow prefix length values greater than 32 for IPv4 or greater than 128 for IPv6.
func (builder *IPv6AddressStringParamsBuilder) AllowPrefixesBeyondAddressSize(allow bool) *IPv6AddressStringParamsBuilder {
	builder.allowPrefixesBeyondAddressSize(allow)
	return builder
}

// AllowPrefixLenLeadingZeros dictates whether to allow leading zeros in the prefix length like "1.2.3.4/016".
func (builder *IPv6AddressStringParamsBuilder) AllowPrefixLenLeadingZeros(allow bool) *IPv6AddressStringParamsBuilder {
	builder.allowPrefixLengthLeadingZeros(allow)
	return builder
}

type ipv4AddressStringParameters struct {
	ipAddressStringFormatParameters

	no_inet_aton_hex,
	no_inet_aton_octal,
	no_inet_aton_joinedSegments,
	inet_aton_single_segment_mask,
	inet_aton_extraneous_digits,
	no_inet_aton_leading_zeros bool
}

// Allows_inet_aton_hex allows IPv4 inet_aton hexadecimal format "0xa.0xb.0xc.0cd".
func (params *ipv4AddressStringParameters) Allows_inet_aton_hex() bool {
	return !params.no_inet_aton_hex
}

// Allows_inet_aton_octal allows IPv4 inet_aton octal format, "04.05.06.07" being an example.
// Can be overridden by AllowLeadingZeros
func (params *ipv4AddressStringParameters) Allows_inet_aton_octal() bool {
	return !params.no_inet_aton_octal
}

// Allows_inet_aton_joinedSegments allows IPv4 joined segments like "1.2.3", "1.2', or just "1".
//
// For the case of just 1 segment, the behaviour is controlled by allowSingleSegment
func (params *ipv4AddressStringParameters) Allows_inet_aton_joinedSegments() bool {
	return !params.no_inet_aton_joinedSegments
}

// Allows_inet_aton_single_segment_mask indicates whether you allow a mask that looks like a prefix length when you allow IPv4 joined segments: "1.2.3.5/255".
func (params *ipv4AddressStringParameters) Allows_inet_aton_single_segment_mask() bool {
	return params.inet_aton_single_segment_mask
}

// Allows single-segment inet_aton strings to have extraneous digits.
// This allows up to 31 digits when parsing for both IPv4 and IPv6.
// This allows an unlimited number of digits when parsing for just IPv4 (ie {@link IPAddressStringParameters#allowIPv6} is false).
//
// Digits that go beyond 32 bits are essentially ignored.
// The number of digits before exceeding 32 bits depends on the radix.
// The value of the most significant digit before exceeding 32 bits depends on the radix.
//
// The resulting address is the modulus of the address with the 32-bit unsigned int maximum value,
// or equivalently the truncation of the address to 32 bits.
func (params *ipv4AddressStringParameters) Allows_inet_aton_extraneous_digits() bool {
	return params.inet_aton_extraneous_digits
}

// Allows_inet_aton_leading_zeros allows IPv4 inet_aton hexadecimal or octal to have leading zeros, such as in the first two segments of "0x0a.00b.c.d".
// The first 0 is not considered a leading zero, it either denotes octal or hex depending on whether it is followed by an 'x'.
// Zeros that appear afterwards are inet_aton leading zeros.
func (params *ipv4AddressStringParameters) Allows_inet_aton_leading_zeros() bool {
	return !params.no_inet_aton_leading_zeros
}

// IPv4AddressStringParamsBuilder builds an immutable IPv4AddressStringParams for controlling parsing of IPv4 address strings.
type IPv4AddressStringParamsBuilder struct {
	// This is not anonymous since it clashes with IPAddressStringFormatParamsBuilder,
	// both have ipAddressStringFormatParameters and AddressStringFormatParams
	// IPAddressStringFormatParamsBuilder takes precedence
	params ipv4AddressStringParameters

	IPAddressStringFormatParamsBuilder

	mixedParent *IPv6AddressStringParamsBuilder
}

// ToParams returns an immutable IPv4AddressStringParams instance built by this builder.
func (builder *IPv4AddressStringParamsBuilder) ToParams() IPv4AddressStringParams {
	result := &builder.params
	result.ipAddressStringFormatParameters = *builder.IPAddressStringFormatParamsBuilder.ToParams().(*ipAddressStringFormatParameters)
	return result
}

// GetEmbeddedIPv4AddressParentBuilder the parent IPv6AddressStringParamsBuilder,
// if this builder was obtained by a call to getEmbeddedIPv4ParamsBuilder() from IPv6AddressStringParamsBuilder.
func (builder *IPv4AddressStringParamsBuilder) GetEmbeddedIPv4AddressParentBuilder() *IPv6AddressStringParamsBuilder {
	return builder.mixedParent
}

// GetRangeParamsBuilder returns a builder that builds the range parameters for these IPv4 address string parameters.
func (builder *IPv4AddressStringParamsBuilder) GetRangeParamsBuilder() *RangeParamsBuilder {
	result := &builder.rangeParamsBuilder
	result.parent = builder
	return result
}

// Allows_inet_aton_hex allows IPv4 inet_aton hexadecimal format "0xa.0xb.0xc.0cd".
func (builder *IPv4AddressStringParamsBuilder) Allows_inet_aton_hex() bool {
	return !builder.params.no_inet_aton_hex
}

// Allows_inet_aton_octal allows IPv4 inet_aton octal format, "04.05.06.07" being an example.
// Can be overridden by AllowLeadingZeros
func (builder *IPv4AddressStringParamsBuilder) Allows_inet_aton_octal() bool {
	return !builder.params.no_inet_aton_octal
}

// Allows_inet_aton_joinedSegments allows IPv4 joined segments like "1.2.3", "1.2', or just "1".
//
// For the case of just 1 segment, the behaviour is controlled by allowSingleSegment
func (builder *IPv4AddressStringParamsBuilder) Allows_inet_aton_joinedSegments() bool {
	return !builder.params.no_inet_aton_joinedSegments
}

// Allows_inet_aton_single_segment_mask indicates whether you allow a mask that looks like a prefix length when you allow IPv4 joined segments: "1.2.3.5/255".
func (builder *IPv4AddressStringParamsBuilder) Allows_inet_aton_single_segment_mask() bool {
	return builder.params.inet_aton_single_segment_mask
}

// Allows single-segment inet_aton strings to have extraneous digits.
// This allows up to 31 digits when parsing for both IPv4 and IPv6.
// This allows an unlimited number of digits when parsing for just IPv4 (ie {@link IPAddressStringParameters#allowIPv6} is false).
//
// Digits that go beyond 32 bits are essentially ignored.
// The number of digits before exceeding 32 bits depends on the radix.
// The value of the most significant digit before exceeding 32 bits depends on the radix.
//
// The resulting address is the modulus of the address with the 32-bit unsigned int maximum value,
// or equivalently the truncation of the address to 32 bits.
func (builder *IPv4AddressStringParamsBuilder) Allows_inet_aton_extraneous_digits() bool {
	return builder.params.inet_aton_extraneous_digits
}

// Allows_inet_aton_leading_zeros allows IPv4 inet_aton hexadecimal or octal to have leading zeros, such as in the first two segments of "0x0a.00b.c.d".
// The first 0 is not considered a leading zero, it either denotes octal or hex depending on whether it is followed by an 'x'.
// Zeros that appear afterwards are inet_aton leading zeros.
func (builder *IPv4AddressStringParamsBuilder) Allows_inet_aton_leading_zeros() bool {
	return !builder.params.no_inet_aton_leading_zeros
}

// Set populates this builder with the values from the given IPv4AddressStringParams.
func (builder *IPv4AddressStringParamsBuilder) Set(params IPv4AddressStringParams) *IPv4AddressStringParamsBuilder {
	if p, ok := params.(*ipv4AddressStringParameters); ok {
		builder.params = *p
	} else {
		builder.params = ipv4AddressStringParameters{
			no_inet_aton_hex:              !params.Allows_inet_aton_hex(),
			no_inet_aton_octal:            !params.Allows_inet_aton_octal(),
			no_inet_aton_joinedSegments:   !params.Allows_inet_aton_joinedSegments(),
			inet_aton_single_segment_mask: params.Allows_inet_aton_single_segment_mask(),
			inet_aton_extraneous_digits:   params.Allows_inet_aton_extraneous_digits(),
			no_inet_aton_leading_zeros:    !params.Allows_inet_aton_leading_zeros(),
		}
	}
	builder.IPAddressStringFormatParamsBuilder.set(params)
	return builder
}

// Allow_inet_aton dictates whether to allow any IPv4 inet_aton format, whether hex, octal, or joined segments.
// Allows joined segments, resulting in just 2, 3 or 4 segments.  Allows octal or hex segments.
// Allows an unlimited number of leading zeros.
// To allow just a single segment, use AllowSingleSegment
// This does not affect whether extraneous digits are allowed, which can be allowed with Allow_inet_aton_extraneous_digits
func (builder *IPv4AddressStringParamsBuilder) Allow_inet_aton(allow bool) *IPv4AddressStringParamsBuilder {
	builder.params.no_inet_aton_joinedSegments = !allow
	builder.params.no_inet_aton_octal = !allow
	builder.params.no_inet_aton_hex = !allow
	builder.allowUnlimitedLeadingZeros(allow)
	return builder
}

// Allow_inet_aton_hex dictates whether to allow IPv4 inet_aton hexadecimal format "0xa.0xb.0xc.0cd".
func (builder *IPv4AddressStringParamsBuilder) Allow_inet_aton_hex(allow bool) *IPv4AddressStringParamsBuilder {
	builder.params.no_inet_aton_hex = !allow
	return builder
}

// Allow_inet_aton_octal dictates whether to allow IPv4 inet_aton octal format, "04.05.06.07" being an example.
func (builder *IPv4AddressStringParamsBuilder) Allow_inet_aton_octal(allow bool) *IPv4AddressStringParamsBuilder {
	builder.params.no_inet_aton_octal = !allow
	return builder
}

// Allow_inet_aton_leading_zeros dictates whether to allow IPv4 inet_aton hexadecimal or octal to have leading zeros, such as in the first two segments of "0x0a.00b.c.d".
//
//	The first 0 is not considered a leading zero, it either denotes octal or hex depending on whether it is followed by an 'x'.
//	Zeros that appear afterwards are inet_aton leading zeros.
func (builder *IPv4AddressStringParamsBuilder) Allow_inet_aton_leading_zeros(allow bool) *IPv4AddressStringParamsBuilder {
	builder.params.no_inet_aton_leading_zeros = !allow
	return builder
}

// Allow_inet_aton_joinedSegments dictates whether to allow IPv4 joined segments like "1.2.3", "1.2", or just "1".
//
// For the case of just 1 segment, the behaviour is controlled by AllowSingleSegment.
func (builder *IPv4AddressStringParamsBuilder) Allow_inet_aton_joinedSegments(allow bool) *IPv4AddressStringParamsBuilder {
	builder.params.no_inet_aton_joinedSegments = !allow
	return builder
}

// Allow_inet_aton_single_segment_mask dictates whether to allow a mask that looks like a prefix length when you allow IPv4 joined segments: "1.2.3.5/255".
func (builder *IPv4AddressStringParamsBuilder) Allow_inet_aton_single_segment_mask(allow bool) *IPv4AddressStringParamsBuilder {
	builder.params.inet_aton_single_segment_mask = allow
	return builder
}

// Allows single-segment inet_aton strings to have extraneous digits.
// This allows up to 31 digits when parsing for both IPv4 and IPv6.
// This allows an unlimited number of digits when parsing for just IPv4 (ie {@link IPAddressStringParameters#allowIPv6} is false).
//
// Digits that go beyond 32 bits are essentially ignored.
// The number of digits before exceeding 32 bits depends on the radix.
// The value of the most significant digit before exceeding 32 bits depends on the radix.
//
// The resulting address is the modulus of the address with the 32-bit unsigned int maximum value,
// or equivalently the truncation of the address to 32 bits.
func (builder *IPv4AddressStringParamsBuilder) Allow_inet_aton_extraneous_digits(allow bool) *IPv4AddressStringParamsBuilder {
	builder.params.inet_aton_extraneous_digits = allow
	return builder
}

// AllowWildcardedSeparator dictates whether the wildcard '*' or '%' can replace the segment separators '.' and ':'.
// If so, then you can write addresses like *.* or *:*
func (builder *IPv4AddressStringParamsBuilder) AllowWildcardedSeparator(allow bool) *IPv4AddressStringParamsBuilder {
	builder.allowWildcardedSeparator(allow)
	return builder
}

// AllowLeadingZeros dictates whether to allow addresses with segments that have leasing zeros like "001.2.3.004" or "1:000a::".
// For IPV4, this option overrides inet_aton octal.
//
// Single segment addresses that must have the requisite length to be parsed are not affected by this flag.
func (builder *IPv4AddressStringParamsBuilder) AllowLeadingZeros(allow bool) *IPv4AddressStringParamsBuilder {
	builder.allowLeadingZeros(allow)
	return builder
}

// AllowUnlimitedLeadingZeros dictates whether to allow leading zeros that extend segments
// beyond the usual segment length, which is 3 for IPv4 dotted-decimal and 4 for IPv6.
// However, this only takes effect if leading zeros are allowed, which is when
// AllowsLeadingZeros is true or the address is IPv4 and Allows_inet_aton_octal is true.
//
// For example, this determines whether you allow "0001.0002.0003.0004">
func (builder *IPv4AddressStringParamsBuilder) AllowUnlimitedLeadingZeros(allow bool) *IPv4AddressStringParamsBuilder {
	builder.allowUnlimitedLeadingZeros(allow)
	return builder
}

// SetRangeParams populates this builder with the values from the given RangeParams.
func (builder *IPv4AddressStringParamsBuilder) SetRangeParams(rangeParams RangeParams) *IPv4AddressStringParamsBuilder {
	builder.setRangeParameters(rangeParams)
	return builder
}

// AllowPrefixesBeyondAddressSize dictates whether to allow prefix length values greater than 32 for IPv4 or greater than 128 for IPv6.
func (builder *IPv4AddressStringParamsBuilder) AllowPrefixesBeyondAddressSize(allow bool) *IPv4AddressStringParamsBuilder {
	builder.allowPrefixesBeyondAddressSize(allow)
	return builder
}

// AllowPrefixLenLeadingZeros dictates whether to allow leading zeros in the prefix length like "1.2.3.4/016".
func (builder *IPv4AddressStringParamsBuilder) AllowPrefixLenLeadingZeros(allow bool) *IPv4AddressStringParamsBuilder {
	builder.allowPrefixLengthLeadingZeros(allow)
	return builder
}

// AllowBinary dictates whether to allow binary addresses like "11111111.0.1.0" or "1111111111111111::".
func (builder *IPv4AddressStringParamsBuilder) AllowBinary(allow bool) *IPv4AddressStringParamsBuilder {
	builder.allowBinary(allow)
	return builder
}

// IPVersion is the version type used by IP string parameters.
// It is interchangeable with the ipaddr.Version, the more generic version type used by the library as a whole.
type IPVersion int

const (
	// IndeterminateIPVersion represents an unspecified IP address version.
	IndeterminateIPVersion IPVersion = 0

	// IPv4 represents Internet Protocol version 4.
	IPv4 IPVersion = 4

	// IPv6 represents Internet Protocol version 6.
	IPv6 IPVersion = 6
)

// IsIPv6 returns true if this represents version 6
func (version IPVersion) IsIPv6() bool {
	return version == IPv6
}

// IsIPv4 returns true if this represents version 4
func (version IPVersion) IsIPv4() bool {
	return version == IPv4
}

// IsIndeterminate returns true if this represents an unspecified IP address version
func (version IPVersion) IsIndeterminate() bool {
	return version != IPv4 && version != IPv6
}

// String returns "IPv4", "IPv6", or the zero-value "" representing an indeterminate version
func (version IPVersion) String() string {
	switch version {
	case 4:
		return "IPv4"
	case 6:
		return "IPv6"
	}
	return ""
}
