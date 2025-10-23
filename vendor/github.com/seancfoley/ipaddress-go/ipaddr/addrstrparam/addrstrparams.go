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

type AddressStringFormatParams interface {
	// AllowsWildcardedSeparator controls whether the wildcard '*' or '%' can replace the segment separators '.' and ':'.
	// If so, then you can write addresses like "*.*" or "*:*".
	AllowsWildcardedSeparator() bool

	// AllowsLeadingZeros indicates whether you allow addresses with segments that have leasing zeros like "001.2.3.004" or "1:000a::".
	// For IPV4, this option overrides inet_aton octal.
	//
	// Single segment addresses that must have the requisite length to be parsed are not affected by this flag.
	AllowsLeadingZeros() bool

	// AllowsUnlimitedLeadingZeros determines if you allow leading zeros that extend segments
	// beyond the usual segment length, which is 3 for IPv4 dotted-decimal and 4 for IPv6.
	// However, this only takes effect if leading zeros are allowed, which is when
	// AllowsLeadingZeros is true or the address is IPv4 and Allows_inet_aton_octal is true.
	//
	// For example, this determines whether you allow "0001.0002.0003.0004".
	AllowsUnlimitedLeadingZeros() bool

	// GetRangeParams returns the RangeParams describing whether ranges of values are allowed and what wildcards are allowed.
	GetRangeParams() RangeParams
}

type AddressStringParams interface {
	// AllowsEmpty indicates whether it allows zero-length address strings: ""
	AllowsEmpty() bool

	// AllowsSingleSegment allows an address to be specified as a single value, eg ffffffff, without the standard use of segments like "1.2.3.4" or "1:2:4:3:5:6:7:8"
	AllowsSingleSegment() bool

	// AllowsAll indicates if we allow the string of just the wildcard "*" to denote all addresses of all version.
	// If false, then for IP addresses we check the preferred version with GetPreferredVersion(), and then check AllowsWildcardedSeparator,
	// to determine if the string represents all addresses of that version.
	AllowsAll() bool
}

// RangeParams indicates what wildcards and ranges are allowed in the string.
type RangeParams interface {
	// AllowsWildcard indicates whether '*' is allowed to denote segments covering all possible segment values
	AllowsWildcard() bool

	// AllowsRangeSeparator indicates whether '-' (or the expected range separator for the address) is allowed to denote a range from lower to higher, like 1-10
	AllowsRangeSeparator() bool

	// AllowsSingleWildcard indicates whether to allow a segment terminating with '_' characters, which represent any digit
	AllowsSingleWildcard() bool

	// AllowsReverseRange indicates whether '-' (or the expected range separator for the address) is allowed to denote a range from higher to lower, like 10-1
	AllowsReverseRange() bool

	// AllowsInferredBoundary indicates whether a missing range value before or after a '-' is allowed to denote the mininum or maximum potential value
	AllowsInferredBoundary() bool
}

var _ AddressStringFormatParams = &addressStringFormatParameters{}
var _ AddressStringParams = &addressStringParameters{}
var _ RangeParams = &rangeParameters{}

type rangeParameters struct {
	noWildcard, noValueRange, noReverseRange, noSingleWildcard, noInferredBoundary bool
}

var (
	// NoRange - use no wildcards nor range separators
	NoRange RangeParams = &rangeParameters{
		noWildcard:         true,
		noValueRange:       true,
		noReverseRange:     true,
		noSingleWildcard:   true,
		noInferredBoundary: true,
	}

	// WildcardOnly - use this to support addresses like "1.*.3.4" or "1::*:3" or "1.2_.3.4" or "1::a__:3"
	WildcardOnly RangeParams = &rangeParameters{
		noValueRange:   true,
		noReverseRange: true,
	}

	// WildcardAndRange - use this to support addresses supported by the default wildcard options and also addresses like "1.2-3.3.4" or "1:0-ff::".
	WildcardAndRange RangeParams = &rangeParameters{}
)

// AllowsWildcard indicates whether '*' is allowed to denote segments covering all possible segment values.
func (builder *rangeParameters) AllowsWildcard() bool {
	return !builder.noWildcard
}

// AllowsRangeSeparator indicates whether '-' (or the expected range separator for the address) is allowed to denote a range from lower to higher, like 1-10.
func (builder *rangeParameters) AllowsRangeSeparator() bool {
	return !builder.noValueRange
}

// AllowsReverseRange indicates whether '-' (or the expected range separator for the address) is allowed to denote a range from higher to lower, like 10-1.
func (builder *rangeParameters) AllowsReverseRange() bool {
	return !builder.noReverseRange
}

// AllowsInferredBoundary indicates whether a missing range value before or after a '-' is allowed to denote the mininum or maximum potential value.
func (builder *rangeParameters) AllowsInferredBoundary() bool {
	return !builder.noInferredBoundary
}

// AllowsSingleWildcard indicates whether to allow a segment terminating with '_' characters, which represent any digit.
func (builder *rangeParameters) AllowsSingleWildcard() bool {
	return !builder.noSingleWildcard
}

// RangeParamsBuilder is used to build an immutable RangeParams for parsing address strings.
type RangeParamsBuilder struct {
	rangeParameters
	parent any
}

// ToParams returns an immutable RangeParams instance built by this builder.
func (builder *RangeParamsBuilder) ToParams() RangeParams {
	return &builder.rangeParameters
}

// Set initializes this builder with the values from the given RangeParams.
func (builder *RangeParamsBuilder) Set(rangeParams RangeParams) *RangeParamsBuilder {
	if rp, ok := rangeParams.(*rangeParameters); ok {
		builder.rangeParameters = *rp
	} else {
		builder.rangeParameters = rangeParameters{
			noWildcard:         !rangeParams.AllowsWildcard(),
			noValueRange:       !rangeParams.AllowsRangeSeparator(),
			noReverseRange:     !rangeParams.AllowsReverseRange(),
			noSingleWildcard:   !rangeParams.AllowsSingleWildcard(),
			noInferredBoundary: !rangeParams.AllowsInferredBoundary(),
		}
	}
	return builder
}

// GetIPv4ParentBuilder returns the IPv4AddressStringParamsBuilder if this builder was obtained by a call to IPv4AddressStringParamsBuilder.GetRangeParamsBuilder.
func (builder *RangeParamsBuilder) GetIPv4ParentBuilder() *IPv4AddressStringParamsBuilder {
	parent := builder.parent
	if p, ok := parent.(*IPv4AddressStringParamsBuilder); ok {
		return p
	}
	return nil
}

// GetIPv6ParentBuilder returns the IPv6AddressStringParamsBuilder if this builder was obtained by a call to IPv6AddressStringParamsBuilder.GetRangeParamsBuilder.
func (builder *RangeParamsBuilder) GetIPv6ParentBuilder() *IPv6AddressStringParamsBuilder {
	parent := builder.parent
	if p, ok := parent.(*IPv6AddressStringParamsBuilder); ok {
		return p
	}
	return nil
}

// GetMACParentBuilder returns the IPv6AddressStringParamsBuilder if this builder was obtained by a call to IPv6AddressStringParamsBuilder.GetRangeParamsBuilder.
func (builder *RangeParamsBuilder) GetMACParentBuilder() *MACAddressStringFormatParamsBuilder {
	parent := builder.parent
	if p, ok := parent.(*MACAddressStringFormatParamsBuilder); ok {
		return p
	}
	return nil
}

// AllowWildcard dictates whether '*' is allowed to denote segments covering all possible segment values.
func (builder *RangeParamsBuilder) AllowWildcard(allow bool) *RangeParamsBuilder {
	builder.noWildcard = !allow
	return builder
}

// AllowRangeSeparator dictates whether '-' (or the expected range separator for the address) is allowed to denote a range from lower to higher, like 1-10.
func (builder *RangeParamsBuilder) AllowRangeSeparator(allow bool) *RangeParamsBuilder {
	builder.noValueRange = !allow
	return builder
}

// AllowReverseRange dictates whether '-' (or the expected range separator for the address) is allowed to denote a range from higher to lower, like 10-1.
func (builder *RangeParamsBuilder) AllowReverseRange(allow bool) *RangeParamsBuilder {
	builder.noReverseRange = !allow
	return builder
}

// AllowInferredBoundary dictates whether a missing range value before or after a '-' is allowed to denote the mininum or maximum potential value.
func (builder *RangeParamsBuilder) AllowInferredBoundary(allow bool) *RangeParamsBuilder {
	builder.noInferredBoundary = !allow
	return builder
}

// AllowSingleWildcard dictates whether to allow a segment terminating with '_' characters, which represent any digit.
func (builder *RangeParamsBuilder) AllowSingleWildcard(allow bool) *RangeParamsBuilder {
	builder.noSingleWildcard = !allow
	return builder
}

type addressStringParameters struct {
	noEmpty, noAll, noSingleSegment bool
}

// AllowsEmpty indicates whether it allows zero-length address strings: "".
func (params *addressStringParameters) AllowsEmpty() bool {
	return !params.noEmpty
}

// AllowsSingleSegment allows an address to be specified as a single value, eg ffffffff, without the standard use of segments like "1.2.3.4" or "1:2:4:3:5:6:7:8".
func (params *addressStringParameters) AllowsSingleSegment() bool {
	return !params.noSingleSegment
}

// AllowsAll indicates if we allow the string of just the wildcard "*" to denote all addresses of all version.
// If false, then for IP addresses we check the preferred version with GetPreferredVersion(), and then check AllowsWildcardedSeparator(),
// to determine if the string represents all addresses of that version.
func (params *addressStringParameters) AllowsAll() bool {
	return !params.noAll
}

// AddressStringParamsBuilder builds an AddressStringParams.
type AddressStringParamsBuilder struct {
	addressStringParameters
}

func (builder *AddressStringParamsBuilder) set(params AddressStringParams) {
	if p, ok := params.(*addressStringParameters); ok {
		builder.addressStringParameters = *p
	} else {
		builder.addressStringParameters = addressStringParameters{
			noEmpty:         !params.AllowsEmpty(),
			noAll:           !params.AllowsAll(),
			noSingleSegment: !params.AllowsSingleSegment(),
		}
	}
}

// ToParams returns an immutable AddressStringParams instance built by this builder.
func (builder *AddressStringParamsBuilder) ToParams() AddressStringParams {
	return &builder.addressStringParameters
}

func (builder *AddressStringParamsBuilder) allowEmpty(allow bool) {
	builder.noEmpty = !allow
}

func (builder *AddressStringParamsBuilder) allowAll(allow bool) {
	builder.noAll = !allow
}

func (builder *AddressStringParamsBuilder) allowSingleSegment(allow bool) {
	builder.noSingleSegment = !allow
}

// AddressStringFormatParams are parameters specific to a given address type or version that is supplied.
type addressStringFormatParameters struct {
	rangeParams rangeParameters

	noWildcardedSeparator, noLeadingZeros, noUnlimitedLeadingZeros bool
}

// AllowsWildcardedSeparator controls whether the wildcard '*' or '%' can replace the segment separators '.' and ':'.
// If so, then you can write addresses like *.* or *:*
func (params *addressStringFormatParameters) AllowsWildcardedSeparator() bool {
	return !params.noWildcardedSeparator
}

// AllowsLeadingZeros indicates whether you allow addresses with segments that have leasing zeros like "001.2.3.004" or "1:000a::".
// For IPV4, this option overrides inet_aton octal.
//
// Single segment addresses that must have the requisite length to be parsed are not affected by this flag.
func (params *addressStringFormatParameters) AllowsLeadingZeros() bool {
	return !params.noLeadingZeros
}

// AllowsUnlimitedLeadingZeros determines if you allow leading zeros that extend segments
// beyond the usual segment length, which is 3 for IPv4 dotted-decimal and 4 for IPv6.
// However, this only takes effect if leading zeros are allowed, which is when
// AllowsLeadingZeros is true or the address is IPv4 and Allows_inet_aton_octal is true.
//
// For example, this determines whether you allow "0001.0002.0003.0004".
func (params *addressStringFormatParameters) AllowsUnlimitedLeadingZeros() bool {
	return !params.noUnlimitedLeadingZeros
}

// GetRangeParams returns the RangeParams describing whether ranges of values are allowed and what wildcards are allowed.
func (params *addressStringFormatParameters) GetRangeParams() RangeParams {
	return &params.rangeParams
}

// AddressStringFormatParamsBuilder creates parameters for parsing a specific address type or address version.
type AddressStringFormatParamsBuilder struct {
	addressStringFormatParameters

	rangeParamsBuilder RangeParamsBuilder
}

// ToParams returns an immutable AddressStringFormatParams instance built by this builder.
func (builder *AddressStringFormatParamsBuilder) ToParams() AddressStringFormatParams {
	result := &builder.addressStringFormatParameters
	result.rangeParams = *builder.rangeParamsBuilder.ToParams().(*rangeParameters)
	return result
}

func (builder *AddressStringFormatParamsBuilder) set(parms AddressStringFormatParams) {
	if p, ok := parms.(*addressStringFormatParameters); ok {
		builder.addressStringFormatParameters = *p
	} else {
		builder.addressStringFormatParameters = addressStringFormatParameters{
			noWildcardedSeparator:   !parms.AllowsWildcardedSeparator(),
			noLeadingZeros:          !parms.AllowsLeadingZeros(),
			noUnlimitedLeadingZeros: !parms.AllowsUnlimitedLeadingZeros(),
		}
	}
	builder.rangeParamsBuilder.Set(parms.GetRangeParams())
}

func (builder *AddressStringFormatParamsBuilder) setRangeParameters(rangeParams RangeParams) {
	builder.rangeParamsBuilder.Set(rangeParams)
}

// GetRangeParamsBuilder returns a builder that builds the range parameters for these address string format parameters.
func (builder *AddressStringFormatParamsBuilder) GetRangeParamsBuilder() RangeParams {
	return &builder.rangeParamsBuilder
}

func (builder *AddressStringFormatParamsBuilder) allowWildcardedSeparator(allow bool) {
	builder.noWildcardedSeparator = !allow
}

func (builder *AddressStringFormatParamsBuilder) allowLeadingZeros(allow bool) {
	builder.noLeadingZeros = !allow
}

func (builder *AddressStringFormatParamsBuilder) allowUnlimitedLeadingZeros(allow bool) {
	builder.noUnlimitedLeadingZeros = !allow
}
