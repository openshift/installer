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

/*
Package addrstr provides interfaces for specifying how to create specific strings from addresses and address sections,
as well as builder types to construct instances of those interfaces.

For example, StringOptionsBuilder produces instances implementing StringOptions for specifying generic strings.
More specific builders and corresponding interface types exist for more specific address versions and types.

Each instance produced by a builders is immutable.
*/
package addrstr

var (
	falseVal = false
	trueVal  = true
)

const (
	ipv6SegmentSeparator     = ':'
	ipv6ZoneSeparatorStr     = "%"
	ipv4SegmentSeparator     = '.'
	macColonSegmentSeparator = ':'
	rangeSeparatorStr        = "-"
	segmentWildcardStr       = "*"

	MinRadix = 2
	MaxRadix = 85

	invalidRadix = "invalid radix"
)

// Wildcards specifies the wildcards to use when constructing an address string.
// WildcardsBuilder can be used to build an instance of Wildcards.
type Wildcards interface {
	// GetRangeSeparator returns the wildcard used to separate the lower and upper boundary (inclusive) of a range of values.
	// If not set, then the default separator RangeSeparatorStr is used, which is the hyphen '-'.
	GetRangeSeparator() string

	// GetWildcard returns the wildcard used for representing any legitimate value, which is the asterisk '*' by default.
	GetWildcard() string

	// GetSingleWildcard returns the wildcard used for representing any single digit, which is the underscore '_' by default.
	GetSingleWildcard() string
}

type wildcards struct {
	rangeSeparator, wildcard, singleWildcard string //rangeSeparator cannot be empty, the other two can
}

// GetRangeSeparator returns the wildcard used to separate the lower and upper boundary (inclusive) of a range of values.
// If not set, then the default separator RangeSeparatorStr is used, which is the hyphen '-'
func (wildcards *wildcards) GetRangeSeparator() string {
	return wildcards.rangeSeparator
}

// GetWildcard returns the wildcard used for representing any legitimate value, which is the asterisk '*' by default.
func (wildcards *wildcards) GetWildcard() string {
	return wildcards.wildcard
}

// GetSingleWildcard returns the wildcard used for representing any single digit, which is the underscore '_' by default.
func (wildcards *wildcards) GetSingleWildcard() string {
	return wildcards.singleWildcard
}

// DefaultWildcards is the default Wildcards instance, using '-' and '*' as range separator and wildcard.
var DefaultWildcards Wildcards = &wildcards{rangeSeparator: rangeSeparatorStr, wildcard: segmentWildcardStr}

// WildcardsBuilder builds an instance of Wildcards
type WildcardsBuilder struct {
	wildcards
}

// SetRangeSeparator sets the wildcard used to separate the lower and upper boundary (inclusive) of a range of values.
// If not set, then the default separator RangeSeparatorStr is used, which is the hyphen '-'
func (wildcards *WildcardsBuilder) SetRangeSeparator(str string) *WildcardsBuilder {
	wildcards.rangeSeparator = str
	return wildcards
}

// SetWildcard sets the wildcard used for representing any legitimate value, which is the asterisk '*' by default.
func (wildcards *WildcardsBuilder) SetWildcard(str string) *WildcardsBuilder {
	wildcards.wildcard = str
	return wildcards
}

// SetSingleWildcard sets the wildcard used for representing any single digit, which is the underscore '_' by default.
func (wildcards *WildcardsBuilder) SetSingleWildcard(str string) *WildcardsBuilder {
	wildcards.singleWildcard = str
	return wildcards
}

// ToWildcards returns an immutable Wildcards instance built by this builder.
func (wildcards *WildcardsBuilder) ToWildcards() Wildcards {
	res := wildcards.wildcards
	if res.rangeSeparator == "" {
		//rangeSeparator cannot be empty
		res.rangeSeparator = rangeSeparatorStr
	}
	return &res
}

// StringOptions represents a clear way to create a specific type of string.
type StringOptions interface {
	// GetWildcards returns the wildcards specified for use in the string.
	GetWildcards() Wildcards

	// IsReverse indicates whether the string segments should be printed in reverse from the usual order, the usual order being most to least significant.
	IsReverse() bool

	// IsUppercase indicates whether to use uppercase for hexadecimal or other radices with alphabetic characters.
	IsUppercase() bool

	// IsExpandedSegments returns whether segments should be expanded to maximal width, typically by using leading zeros.
	IsExpandedSegments() bool

	// GetRadix returns the radix to be used.  The default is hexadecimal unless built using an IPv4 options builder in which case the default is decimal.
	GetRadix() int

	// GetSeparator returns the separator that separates the divisions of the address, typically ':' or '.'.  HasSeparator indicates if this method should be called.
	// the default is to have no separator, unless built using a MAC, IPv6 or IPv4 options builder in which case the separator is ':' for MAC and IPv6 and '.' for IPv4.
	GetSeparator() byte

	// HasSeparator indicates whether there is a separator.
	// The default is false, no separator, unless built using a MAC, IPv6 or IPv4 options builder in which case there is a default separator.
	HasSeparator() bool

	// GetAddressLabel returns a string to prepend to the entire address string, such as an octal, hex, or binary prefix.
	GetAddressLabel() string

	// GetSegmentStrPrefix returns the string prefix (if any) to prepend to each segment's values, such as an octal, hex or binary prefix.
	GetSegmentStrPrefix() string
}

type stringOptions struct {
	wildcards Wildcards

	base int // default is hex

	//the segment separator and in the case of split digits, the digit separator
	separator byte // default is ' ', but it's typically either '.' or ':'

	segmentStrPrefix,
	addrLabel string

	expandSegments,
	reverse,
	uppercase bool

	hasSeparator *bool // if not set, the default is false, no separator
}

// GetWildcards returns the wildcards specified for use in the string.
func (opts *stringOptions) GetWildcards() Wildcards {
	return opts.wildcards
}

// IsReverse indicates whether the string segments should be printed in reverse from the usual order, the usual order being most to least significant.
func (opts *stringOptions) IsReverse() bool {
	return opts.reverse
}

// IsUppercase indicates whether to use uppercase for hexadecimal or other radices with alphabetic characters.
func (opts *stringOptions) IsUppercase() bool {
	return opts.uppercase
}

// IsExpandedSegments returns whether segments should be expanded to maximal width, typically by using leading zeros.
func (opts *stringOptions) IsExpandedSegments() bool {
	return opts.expandSegments
}

// GetRadix returns the radix to be used.  The default is hexadecimal unless built using an IPv4 options builder in which case the default is decimal.
func (opts *stringOptions) GetRadix() int {
	return opts.base
}

// GetSeparator returns the separator that separates the divisions of the address, typically ':' or '.'.  HasSeparator indicates if this method should be called.
// the default is to have no separator, unless built using a MAC, IPv6 or IPv4 options builder in which case the separator is ':' for MAC and IPv6 and '.' for IPv4.
func (opts *stringOptions) GetSeparator() byte {
	return opts.separator
}

// HasSeparator indicates whether there is a separator.
// The default is false, no separator, unless built using a MAC, IPv6 or IPv4 options builder in which case there is a default separator.
func (opts *stringOptions) HasSeparator() bool {
	if opts.hasSeparator == nil {
		return false
	}
	return *opts.hasSeparator
}

// GetAddressLabel returns a string to prepend to the entire address string, such as an octal, hex, or binary prefix.
func (opts *stringOptions) GetAddressLabel() string {
	return opts.addrLabel
}

// GetSegmentStrPrefix returns the string prefix (if any) to prepend to each segment's values, such as an octal, hex or binary prefix.
func (opts *stringOptions) GetSegmentStrPrefix() string {
	return opts.segmentStrPrefix
}

var _ StringOptions = &stringOptions{}

func getDefaults(radix int, wildcards Wildcards, separator byte) (int, Wildcards, byte) {
	if radix == 0 {
		radix = 16
	}
	if wildcards == nil {
		wildcards = DefaultWildcards
	}
	if separator == 0 {
		separator = ' '
	}
	return radix, wildcards, separator
}

func getIPDefaults(zoneSeparator string) string {
	if len(zoneSeparator) == 0 {
		zoneSeparator = ipv6ZoneSeparatorStr
	}
	return zoneSeparator
}

func getIPv6Defaults(hasSeparator *bool, separator byte) (*bool, byte) {
	if hasSeparator == nil {
		hasSeparator = &trueVal
	}
	if separator == 0 {
		separator = ipv6SegmentSeparator
	}
	return hasSeparator, separator
}

func getIPv4Defaults(hasSeparator *bool, separator byte, radix int) (*bool, byte, int) {
	if hasSeparator == nil {
		hasSeparator = &trueVal
	}
	if radix == 0 {
		radix = 10
	}
	if separator == 0 {
		separator = ipv4SegmentSeparator
	}
	return hasSeparator, separator, radix
}

func getMACDefaults(hasSeparator *bool, separator byte) (*bool, byte) {
	if hasSeparator == nil {
		hasSeparator = &trueVal
	}
	if separator == 0 {
		separator = macColonSegmentSeparator
	}
	return hasSeparator, separator
}

// StringOptionsBuilder is used to build an immutable StringOptions instance.
type StringOptionsBuilder struct {
	stringOptions
}

// SetWildcards specifies the wildcards for use in the string.
func (builder *StringOptionsBuilder) SetWildcards(wildcards Wildcards) *StringOptionsBuilder {
	builder.wildcards = wildcards
	return builder
}

// SetReverse dictates whether the string segments should be printed in reverse from the usual order, the usual order being most to least significant.
func (builder *StringOptionsBuilder) SetReverse(reverse bool) *StringOptionsBuilder {
	builder.reverse = reverse
	return builder
}

// SetUppercase dictates whether to use uppercase for hexadecimal or other radices with alphabetic characters.
func (builder *StringOptionsBuilder) SetUppercase(uppercase bool) *StringOptionsBuilder {
	builder.uppercase = uppercase
	return builder
}

// SetExpandedSegments dictates whether segments should be expanded to maximal width, typically by using leading zeros.
func (builder *StringOptionsBuilder) SetExpandedSegments(expandSegments bool) *StringOptionsBuilder {
	builder.expandSegments = expandSegments
	return builder
}

// SetRadix sets the radix to be used.
// A radix less than MinRadix or greater than MaxRadix results in a panic.
func (builder *StringOptionsBuilder) SetRadix(base int) *StringOptionsBuilder {
	if base < MinRadix && base > MaxRadix {
		panic(invalidRadix)
	}
	builder.base = base
	return builder
}

// SetHasSeparator dictates whether there is a separator.
// The default is false, no separator, unless using a MAC, IPv6 or IPv4 options builder in which case there is a default separator.
func (builder *StringOptionsBuilder) SetHasSeparator(has bool) *StringOptionsBuilder {
	if has {
		builder.hasSeparator = &trueVal
	} else {
		builder.hasSeparator = &falseVal
	}
	return builder
}

// SetSeparator dictates the separator to separate the divisions of the address, typically ':' or '.'.
// HasSeparator indicates if this separator should be used or not.
func (builder *StringOptionsBuilder) SetSeparator(separator byte) *StringOptionsBuilder {
	builder.separator = separator
	builder.SetHasSeparator(true)
	return builder
}

// SetAddressLabel dictates a string to prepend to the entire address string, such as an octal, hex, or binary prefix.
func (builder *StringOptionsBuilder) SetAddressLabel(label string) *StringOptionsBuilder {
	builder.addrLabel = label
	return builder
}

// SetSegmentStrPrefix dictates a string prefix to prepend to each segment's values, such as an octal, hex or binary prefix.
func (builder *StringOptionsBuilder) SetSegmentStrPrefix(prefix string) *StringOptionsBuilder {
	builder.segmentStrPrefix = prefix
	return builder
}

// ToOptions returns an immutable StringOptions instance built by this builder.
func (builder *StringOptionsBuilder) ToOptions() StringOptions {
	res := builder.stringOptions
	res.base, res.wildcards, res.separator = getDefaults(res.base, res.wildcards, res.separator)
	return &res
}

// MACStringOptionsBuilder is used to build an immutable StringOptions instance for MAC address strings.
type MACStringOptionsBuilder struct {
	StringOptionsBuilder
}

// SetWildcards specifies the wildcards for use in the string.
func (builder *MACStringOptionsBuilder) SetWildcards(wildcards Wildcards) *MACStringOptionsBuilder {
	builder.StringOptionsBuilder.SetWildcards(wildcards)
	return builder
}

// SetReverse dictates whether the string segments should be printed in reverse from the usual order, the usual order being most to least significant.
func (builder *MACStringOptionsBuilder) SetReverse(reverse bool) *MACStringOptionsBuilder {
	builder.StringOptionsBuilder.SetReverse(reverse)
	return builder
}

// SetUppercase dictates whether to use uppercase for hexadecimal or other radices with alphabetic characters.
func (builder *MACStringOptionsBuilder) SetUppercase(uppercase bool) *MACStringOptionsBuilder {
	builder.StringOptionsBuilder.SetUppercase(uppercase)
	return builder
}

// SetExpandedSegments dictates whether segments should be expanded to maximal width, typically by using leading zeros.
func (builder *MACStringOptionsBuilder) SetExpandedSegments(expandSegments bool) *MACStringOptionsBuilder {
	builder.StringOptionsBuilder.SetExpandedSegments(expandSegments)
	return builder
}

// SetRadix sets the radix to be used.
func (builder *MACStringOptionsBuilder) SetRadix(base int) *MACStringOptionsBuilder {
	builder.StringOptionsBuilder.SetRadix(base)
	return builder
}

// SetHasSeparator dictates whether there is a separator.
// The default for MAC is true.
func (builder *MACStringOptionsBuilder) SetHasSeparator(has bool) *MACStringOptionsBuilder {
	builder.StringOptionsBuilder.SetHasSeparator(has)
	return builder
}

// SetSeparator dictates the separator to separate the divisions of the address, for MAC the default is ':'.
// HasSeparator indicates if this separator should be used or not.
func (builder *MACStringOptionsBuilder) SetSeparator(separator byte) *MACStringOptionsBuilder {
	builder.StringOptionsBuilder.SetSeparator(separator)
	return builder
}

// SetAddressLabel dictates a string to prepend to the entire address string, such as an octal, hex, or binary prefix.
func (builder *MACStringOptionsBuilder) SetAddressLabel(label string) *MACStringOptionsBuilder {
	builder.StringOptionsBuilder.SetAddressLabel(label)
	return builder
}

// SetSegmentStrPrefix dictates a string prefix to prepend to each segment's values, such as an octal, hex or binary prefix.
func (builder *MACStringOptionsBuilder) SetSegmentStrPrefix(prefix string) *MACStringOptionsBuilder {
	builder.StringOptionsBuilder.SetSegmentStrPrefix(prefix)
	return builder
}

// ToOptions returns an immutable StringOptions instance built by this builder.
func (builder *MACStringOptionsBuilder) ToOptions() StringOptions {
	b := &builder.StringOptionsBuilder
	b.hasSeparator, b.separator = getMACDefaults(b.hasSeparator, b.separator)
	return builder.StringOptionsBuilder.ToOptions()
}

// WildcardOption indicates options indicating when and where to use wildcards.
type WildcardOption string

const (

	// WildcardsNetworkOnly prints wildcards that are part of the network portion (only possible with subnet address notation, otherwise this option is ignored).
	WildcardsNetworkOnly WildcardOption = ""

	// WildcardsAll prints wildcards for any visible (non-compressed) segments.
	WildcardsAll WildcardOption = "allType"
)

// WildcardOptions indicates options indicating when and where to use wildcards, and what wildcards to use.
type WildcardOptions interface {

	// GetWildcardOption returns the WildcardOption to use.
	GetWildcardOption() WildcardOption

	// GetWildcards returns the wildcards to use.
	GetWildcards() Wildcards
}

type wildcardOptions struct {
	wildcardOption WildcardOption
	wildcards      Wildcards
}

// GetWildcardOption returns the WildcardOption to use.
func (opts *wildcardOptions) GetWildcardOption() WildcardOption {
	return opts.wildcardOption
}

// GetWildcards returns the wildcards to use.
func (opts *wildcardOptions) GetWildcards() Wildcards {
	return opts.wildcards
}

var _ WildcardOptions = &wildcardOptions{}

// WildcardOptionsBuilder is used to build an immutable WildcardOptions instance for address strings.
type WildcardOptionsBuilder struct {
	wildcardOptions
}

// SetWildcardOptions dictates the WildcardOption to use.
func (builder *WildcardOptionsBuilder) SetWildcardOptions(wildcardOption WildcardOption) *WildcardOptionsBuilder {
	builder.wildcardOption = wildcardOption
	return builder
}

// SetWildcards dictates the wildcards to use.
func (builder *WildcardOptionsBuilder) SetWildcards(wildcards Wildcards) *WildcardOptionsBuilder {
	builder.wildcards = wildcards
	return builder
}

// ToOptions returns an immutable WildcardOptions instance built by this builder.
func (builder *WildcardOptionsBuilder) ToOptions() WildcardOptions {
	cpy := builder.wildcardOptions
	if builder.wildcards == nil {
		builder.wildcards = DefaultWildcards
	}
	return &cpy
}

// IPStringOptions represents a clear way to create a specific type of IP address or subnet string.
type IPStringOptions interface {
	StringOptions

	// GetAddressSuffix returns a suffix to be appended to the string.
	// .in-addr.arpa, .ip6.arpa, .ipv6-literal.net are examples of suffixes tacked onto the end of address strings.
	GetAddressSuffix() string

	// GetWildcardOption returns the WildcardOption to use.
	GetWildcardOption() WildcardOption

	// GetZoneSeparator indicates the delimiter that separates the zone from the address, the default being '%'.
	GetZoneSeparator() string
}

type ipStringOptions struct {
	stringOptions

	addrSuffix     string
	wildcardOption WildcardOption // default is WildcardsNetworkOnly
	zoneSeparator  string         // default is IPv6ZoneSeparator
}

// GetAddressSuffix returns a suffix to be appended to the string.
// .in-addr.arpa, .ip6.arpa, .ipv6-literal.net are examples of suffixes tacked onto the end of address strings.
func (opts *ipStringOptions) GetAddressSuffix() string {
	return opts.addrSuffix
}

// GetWildcardOptions returns the WildcardOptions to use.
func (opts *ipStringOptions) GetWildcardOptions() WildcardOptions {
	options := &wildcardOptions{
		opts.wildcardOption,
		opts.GetWildcards(),
	}
	return options
}

// GetWildcardOption returns the WildcardOption to use.
func (opts *ipStringOptions) GetWildcardOption() WildcardOption {
	return opts.wildcardOption

}

// GetZoneSeparator returns the delimiter that separates the address from the zone, the default being '%'.
func (opts *ipStringOptions) GetZoneSeparator() string {
	return opts.zoneSeparator
}

var _ IPStringOptions = &ipStringOptions{}

// IPStringOptionsBuilder is used to build an immutable IPStringOptions instance for IP address strings.
type IPStringOptionsBuilder struct {
	StringOptionsBuilder
	ipStringOptions ipStringOptions
}

// SetAddressSuffix dictates a suffix to be appended to the string.
// .in-addr.arpa, .ip6.arpa, .ipv6-literal.net are examples of suffixes tacked onto the end of address strings.
func (builder *IPStringOptionsBuilder) SetAddressSuffix(suffix string) *IPStringOptionsBuilder {
	builder.ipStringOptions.addrSuffix = suffix
	return builder
}

// SetWildcardOptions is a convenience method for setting both the WildcardOption and the Wildcards at the same time.
// It overrides previous calls to SetWildcardOption and SetWildcards,
// and is overridden by subsequent calls to those methods.
func (builder *IPStringOptionsBuilder) SetWildcardOptions(wildcardOptions WildcardOptions) *IPStringOptionsBuilder {
	builder.SetWildcards(wildcardOptions.GetWildcards())
	return builder.SetWildcardOption(wildcardOptions.GetWildcardOption())
}

// SetWildcardOption specifies the WildcardOption for use in the string.
func (builder *IPStringOptionsBuilder) SetWildcardOption(wildcardOption WildcardOption) *IPStringOptionsBuilder {
	builder.ipStringOptions.wildcardOption = wildcardOption
	return builder
}

// SetWildcards specifies the wildcards for use in the string.
func (builder *IPStringOptionsBuilder) SetWildcards(wildcards Wildcards) *IPStringOptionsBuilder {
	builder.StringOptionsBuilder.SetWildcards(wildcards)
	return builder
}

// SetZoneSeparator dictates the separator to separate the zone from the address, the default being '%'
// Zones apply to IPv6 addresses only, not IPv4.
func (builder *IPStringOptionsBuilder) SetZoneSeparator(separator string) *IPStringOptionsBuilder {
	builder.ipStringOptions.zoneSeparator = separator
	return builder
}

// SetReverse dictates whether the string segments should be printed in reverse from the usual order, the usual order being most to least significant.
func (builder *IPStringOptionsBuilder) SetReverse(reverse bool) *IPStringOptionsBuilder {
	builder.StringOptionsBuilder.SetReverse(reverse)
	return builder
}

// SetUppercase dictates whether to use uppercase for hexadecimal or other radices with alphabetic characters.
func (builder *IPStringOptionsBuilder) SetUppercase(uppercase bool) *IPStringOptionsBuilder {
	builder.StringOptionsBuilder.SetUppercase(uppercase)
	return builder
}

// SetExpandedSegments dictates whether segments should be expanded to maximal width, typically by using leading zeros.
func (builder *IPStringOptionsBuilder) SetExpandedSegments(expandSegments bool) *IPStringOptionsBuilder {
	builder.StringOptionsBuilder.SetExpandedSegments(expandSegments)
	return builder
}

// SetRadix sets the radix to be used.
func (builder *IPStringOptionsBuilder) SetRadix(base int) *IPStringOptionsBuilder {
	builder.StringOptionsBuilder.SetRadix(base)
	return builder
}

// SetHasSeparator dictates whether there is a separator.
// The default for IPStringOptionsBuilder is false.
func (builder *IPStringOptionsBuilder) SetHasSeparator(has bool) *IPStringOptionsBuilder {
	builder.StringOptionsBuilder.SetHasSeparator(has)
	return builder
}

// SetSeparator dictates the separator to separate the divisions of the address.
// HasSeparator indicates if this separator should be used or not.
func (builder *IPStringOptionsBuilder) SetSeparator(separator byte) *IPStringOptionsBuilder {
	builder.StringOptionsBuilder.SetSeparator(separator)
	return builder
}

// SetAddressLabel dictates a string to prepend to the entire address string, such as an octal, hex, or binary prefix.
func (builder *IPStringOptionsBuilder) SetAddressLabel(label string) *IPStringOptionsBuilder {
	builder.StringOptionsBuilder.SetAddressLabel(label)
	return builder
}

// SetSegmentStrPrefix dictates a string prefix to prepend to each segment's values, such as an octal, hex or binary prefix.
func (builder *IPStringOptionsBuilder) SetSegmentStrPrefix(prefix string) *IPStringOptionsBuilder {
	builder.StringOptionsBuilder.SetSegmentStrPrefix(prefix)
	return builder
}

// ToOptions returns an immutable IPStringOptions instance built by this builder.
func (builder *IPStringOptionsBuilder) ToOptions() IPStringOptions {
	builder.ipStringOptions.zoneSeparator = getIPDefaults(builder.ipStringOptions.zoneSeparator)
	res := builder.ipStringOptions
	res.stringOptions = *builder.StringOptionsBuilder.ToOptions().(*stringOptions)
	return &res
}

// IPv4StringOptionsBuilder is used to build an immutable IPStringOptions instance for IPv4 address strings.
type IPv4StringOptionsBuilder struct {
	IPStringOptionsBuilder
}

// SetAddressSuffix dictates a suffix to be appended to the string.
// .in-addr.arpa, .ip6.arpa, .ipv6-literal.net are examples of suffixes tacked onto the end of address strings.
func (builder *IPv4StringOptionsBuilder) SetAddressSuffix(suffix string) *IPv4StringOptionsBuilder {
	builder.IPStringOptionsBuilder.SetAddressSuffix(suffix)
	return builder
}

// SetWildcardOptions is a convenience method for setting both the WildcardOption and the Wildcards at the same time.
// It overrides previous calls to SetWildcardOption and SetWildcards,
// and is overridden by subsequent calls to those methods.
func (builder *IPv4StringOptionsBuilder) SetWildcardOptions(wildcardOptions WildcardOptions) *IPv4StringOptionsBuilder {
	builder.IPStringOptionsBuilder.SetWildcardOptions(wildcardOptions)
	return builder.SetWildcardOption(wildcardOptions.GetWildcardOption())
}

// SetWildcardOption specifies the WildcardOption for use in the string.
func (builder *IPv4StringOptionsBuilder) SetWildcardOption(wildcardOption WildcardOption) *IPv4StringOptionsBuilder {
	builder.IPStringOptionsBuilder.SetWildcardOption(wildcardOption)
	return builder
}

// SetWildcards specifies the wildcards for use in the string.
func (builder *IPv4StringOptionsBuilder) SetWildcards(wildcards Wildcards) *IPv4StringOptionsBuilder {
	builder.IPStringOptionsBuilder.SetWildcards(wildcards)
	return builder
}

// SetReverse dictates whether the string segments should be printed in reverse from the usual order, the usual order being most to least significant.
func (builder *IPv4StringOptionsBuilder) SetReverse(reverse bool) *IPv4StringOptionsBuilder {
	builder.IPStringOptionsBuilder.SetReverse(reverse)
	return builder
}

// SetUppercase dictates whether to use uppercase for hexadecimal or other radices with alphabetic characters.
func (builder *IPv4StringOptionsBuilder) SetUppercase(uppercase bool) *IPv4StringOptionsBuilder {
	builder.IPStringOptionsBuilder.SetUppercase(uppercase)
	return builder
}

// SetExpandedSegments dictates whether segments should be expanded to maximal width, typically by using leading zeros.
func (builder *IPv4StringOptionsBuilder) SetExpandedSegments(expandSegments bool) *IPv4StringOptionsBuilder {
	builder.IPStringOptionsBuilder.SetExpandedSegments(expandSegments)
	return builder
}

// SetRadix sets the radix to be used.
func (builder *IPv4StringOptionsBuilder) SetRadix(base int) *IPv4StringOptionsBuilder {
	builder.IPStringOptionsBuilder.SetRadix(base)
	return builder
}

// SetHasSeparator dictates whether there is a separator.
// The default for IPv4 is true.
func (builder *IPv4StringOptionsBuilder) SetHasSeparator(has bool) *IPv4StringOptionsBuilder {
	builder.IPStringOptionsBuilder.SetHasSeparator(has)
	return builder
}

// SetSeparator dictates the separator to separate the divisions of the address, for IPv4 the default is '.'.
// HasSeparator indicates if this separator should be used or not.
func (builder *IPv4StringOptionsBuilder) SetSeparator(separator byte) *IPv4StringOptionsBuilder {
	builder.IPStringOptionsBuilder.SetSeparator(separator)
	return builder
}

// SetAddressLabel dictates a string to prepend to the entire address string, such as an octal, hex, or binary prefix.
func (builder *IPv4StringOptionsBuilder) SetAddressLabel(label string) *IPv4StringOptionsBuilder {
	builder.IPStringOptionsBuilder.SetAddressLabel(label)
	return builder
}

// SetSegmentStrPrefix dictates a string prefix to prepend to each segment's values, such as an octal, hex or binary prefix.
func (builder *IPv4StringOptionsBuilder) SetSegmentStrPrefix(prefix string) *IPv4StringOptionsBuilder {
	builder.IPStringOptionsBuilder.SetSegmentStrPrefix(prefix)
	return builder
}

// ToOptions returns an immutable IPStringOptions instance built by this builder.
func (builder *IPv4StringOptionsBuilder) ToOptions() IPStringOptions {
	b := &builder.StringOptionsBuilder
	b.hasSeparator, b.separator, b.base = getIPv4Defaults(b.hasSeparator, b.separator, b.base)
	return builder.IPStringOptionsBuilder.ToOptions()
}

// IPv6StringOptions provides a clear way to create a specific type of IPv6 address string.
type IPv6StringOptions interface {
	IPStringOptions

	// GetIPv4Opts returns the options used for creating the embedded IPv4 address string in a mixed IPv6 address,
	// which comes from the last 32 bits of the IPv6 address.
	// For example: "a:b:c:d:e:f:1.2.3.4"
	GetIPv4Opts() IPStringOptions

	// GetCompressOptions returns the CompressOptions which specify how to compress zero-segments in the IPv6 address or subnet string.
	GetCompressOptions() CompressOptions

	// IsSplitDigits indicates whether every digit is separated from every other by separators.  If mixed, this option is ignored.
	IsSplitDigits() bool // can produce addrerr.IncompatibleAddressError for ranged series

	// IsMixed specifies that the last two segments of the IPv6 address should be printed as an IPv4 address, resulting in a mixed IPv6/v4 string.
	IsMixed() bool // can produce addrerr.IncompatibleAddressError for ranges in the IPv4 part of the series
}

type ipv6StringOptions struct {
	ipStringOptions
	ipv4Opts IPStringOptions

	//can be nil, which means no compression
	compressOptions CompressOptions

	splitDigits bool
}

// IsSplitDigits indicates whether every digit is separated from every other by separators.  If mixed, this option is ignored.
func (opts *ipv6StringOptions) IsSplitDigits() bool {
	return opts.splitDigits
}

// GetIPv4Opts returns the IPv4 string options to be used on the IPv4 address section in a mixed IPv6/v4 string.
func (opts *ipv6StringOptions) GetIPv4Opts() IPStringOptions {
	return opts.ipv4Opts
}

// GetCompressOptions returns the CompressOptions which specify how to compress zero-segments in the IPv6 address or subnet string.
func (opts *ipv6StringOptions) GetCompressOptions() CompressOptions {
	return opts.compressOptions
}

// IsMixed specifies that the last two segments of the IPv6 address should be printed as an IPv4 address, resulting in a mixed IPv6/v4 string.
func (opts *ipv6StringOptions) IsMixed() bool {
	return opts.ipv4Opts != nil
}

var _ IPv6StringOptions = &ipv6StringOptions{}

// IPv6StringOptionsBuilder is used to build an immutable IPv6StringOptions instance for IPv6 address strings.
type IPv6StringOptionsBuilder struct {
	opts ipv6StringOptions

	IPStringOptionsBuilder

	makeMixed bool
}

// IsMixed specifies whether the last two segments of the IPv6 address should be printed as an IPv4 address, resulting in a mixed IPv6/v4 string.
func (builder *IPv6StringOptionsBuilder) IsMixed() bool {
	return builder.makeMixed
}

// GetIPv4Opts returns the IPv4 string options to be used on the IPv4 address section in a mixed IPv6/v4 string.
func (builder *IPv6StringOptionsBuilder) GetIPv4Opts() IPStringOptions {
	return builder.opts.ipv4Opts
}

// GetCompressOptions returns the CompressOptions which specify how to compress zero-segments in the IPv6 address or subnet string.
func (builder *IPv6StringOptionsBuilder) GetCompressOptions() CompressOptions {
	return builder.opts.compressOptions
}

// SetSplitDigits dictates whether every digit is separated from every other by separators.  If mixed, this option is ignored.
func (builder *IPv6StringOptionsBuilder) SetSplitDigits(splitDigits bool) *IPv6StringOptionsBuilder {
	builder.opts.splitDigits = splitDigits
	return builder
}

// SetCompressOptions sets the CompressOptions which specify how to compress zero-segments in the IPv6 address or subnet string.
func (builder *IPv6StringOptionsBuilder) SetCompressOptions(compressOptions CompressOptions) *IPv6StringOptionsBuilder {
	builder.opts.compressOptions = compressOptions
	return builder
}

// SetMixed dictates whether the string should be a mixed IPv6/v4 string, in which the last two segments of the IPv6 address should be printed as an IPv4 address.
func (builder *IPv6StringOptionsBuilder) SetMixed(makeMixed bool) *IPv6StringOptionsBuilder {
	builder.makeMixed = makeMixed
	return builder
}

// SetMixedOptions supplies the IPv4 options to be used on the IPv4 section of a mixed string.  Calling this method sets the string to be a mixed IPv6/v4 string.
func (builder *IPv6StringOptionsBuilder) SetMixedOptions(ipv4Options IPStringOptions) *IPv6StringOptionsBuilder {
	builder.makeMixed = true
	builder.opts.ipv4Opts = ipv4Options
	return builder
}

// SetWildcardOptions is a convenience method for setting both the WildcardOption and the Wildcards at the same time
// It overrides previous calls to SetWildcardOption and SetWildcards,
// and is overridden by subsequent calls to those methods.
func (builder *IPv6StringOptionsBuilder) SetWildcardOptions(wildcardOptions WildcardOptions) *IPv6StringOptionsBuilder {
	builder.IPStringOptionsBuilder.SetWildcardOptions(wildcardOptions)
	return builder
}

// SetWildcardOption specifies the WildcardOption for use in the string.
func (builder *IPv6StringOptionsBuilder) SetWildcardOption(wildcardOption WildcardOption) *IPv6StringOptionsBuilder {
	builder.IPStringOptionsBuilder.SetWildcardOption(wildcardOption)
	return builder
}

// SetWildcards specifies the wildcards for use in the string.
func (builder *IPv6StringOptionsBuilder) SetWildcards(wildcards Wildcards) *IPv6StringOptionsBuilder {
	builder.IPStringOptionsBuilder.SetWildcards(wildcards)
	return builder
}

// SetExpandedSegments dictates whether segments should be expanded to maximal width, typically by using leading zeros.
func (builder *IPv6StringOptionsBuilder) SetExpandedSegments(expandSegments bool) *IPv6StringOptionsBuilder {
	builder.IPStringOptionsBuilder.SetExpandedSegments(expandSegments)
	return builder
}

// SetRadix sets the radix to be used.
func (builder *IPv6StringOptionsBuilder) SetRadix(base int) *IPv6StringOptionsBuilder {
	builder.IPStringOptionsBuilder.SetRadix(base)
	return builder
}

// SetHasSeparator dictates whether there is a separator.
// The default for IPv6 is true.
func (builder *IPv6StringOptionsBuilder) SetHasSeparator(has bool) *IPv6StringOptionsBuilder {
	builder.IPStringOptionsBuilder.SetHasSeparator(has)
	return builder
}

// SetSeparator dictates the separator to separate the divisions of the address, for IPv6 the default is ':'.
// HasSeparator indicates if this separator should be used or not.
func (builder *IPv6StringOptionsBuilder) SetSeparator(separator byte) *IPv6StringOptionsBuilder {
	builder.IPStringOptionsBuilder.SetSeparator(separator)
	return builder
}

// SetZoneSeparator dictates the separator to separate the zone from the address, the default being '%'.
func (builder *IPv6StringOptionsBuilder) SetZoneSeparator(separator string) *IPv6StringOptionsBuilder {
	builder.IPStringOptionsBuilder.SetZoneSeparator(separator)
	return builder
}

// SetAddressSuffix dictates a suffix to be appended to the string.
// .in-addr.arpa, .ip6.arpa, .ipv6-literal.net are examples of suffixes tacked onto the end of address strings.
func (builder *IPv6StringOptionsBuilder) SetAddressSuffix(suffix string) *IPv6StringOptionsBuilder {
	builder.IPStringOptionsBuilder.SetAddressSuffix(suffix)
	return builder
}

// SetSegmentStrPrefix dictates a string prefix to prepend to each segment's values, such as an octal, hex or binary prefix.
func (builder *IPv6StringOptionsBuilder) SetSegmentStrPrefix(prefix string) *IPv6StringOptionsBuilder {
	builder.IPStringOptionsBuilder.SetSegmentStrPrefix(prefix)
	return builder
}

// SetReverse dictates whether the string segments should be printed in reverse from the usual order, the usual order being most to least significant.
func (builder *IPv6StringOptionsBuilder) SetReverse(reverse bool) *IPv6StringOptionsBuilder {
	builder.IPStringOptionsBuilder.SetReverse(reverse)
	return builder
}

// SetUppercase dictates whether to use uppercase for hexadecimal or other radices with alphabetic characters.
func (builder *IPv6StringOptionsBuilder) SetUppercase(upper bool) *IPv6StringOptionsBuilder {
	builder.IPStringOptionsBuilder.SetUppercase(upper)
	return builder
}

// ToOptions returns an immutable IPv6StringOptions instance built by this builder.
func (builder *IPv6StringOptionsBuilder) ToOptions() IPv6StringOptions {
	if builder.makeMixed {
		if builder.opts.ipv4Opts == nil {
			builder.opts.ipv4Opts = new(IPv4StringOptionsBuilder).SetExpandedSegments(builder.expandSegments).
				SetWildcardOption(builder.ipStringOptions.wildcardOption).
				SetWildcards(builder.wildcards).ToOptions()
		}
	} else {
		builder.opts.ipv4Opts = nil
	}
	b := &builder.IPStringOptionsBuilder.StringOptionsBuilder
	b.hasSeparator, b.separator = getIPv6Defaults(b.hasSeparator, b.separator)
	res := builder.opts
	res.ipStringOptions = *builder.IPStringOptionsBuilder.ToOptions().(*ipStringOptions)
	return &res
}

// CompressionChoiceOptions specify which zero-segments should be compressed.
type CompressionChoiceOptions string

const (
	// HostPreferred - if there is a host section, compress the host along with any adjoining zero-segments, otherwise compress a range of zero-segments.
	HostPreferred CompressionChoiceOptions = "host preferred"

	// MixedPreferred - if there is a mixed section that is compressible according to the MixedCompressionOptions, compress the mixed section along with any adjoining zero-segments, otherwise compress a range of zero-segments.
	MixedPreferred CompressionChoiceOptions = "mixed preferred"

	// ZerosOrHost - compress the largest range of zero or host segments.
	ZerosOrHost CompressionChoiceOptions = ""

	// ZerosCompression - compress the largest range of zero-segments.
	ZerosCompression CompressionChoiceOptions = "zeros"
)

// CompressHost indicates if a host of a prefixed address should be compressed.
func (choice CompressionChoiceOptions) CompressHost() bool {
	return choice != ZerosCompression
}

// MixedCompressionOptions specify which zero-segments should be compressed in mixed IPv6/v4 strings.
type MixedCompressionOptions string

const (
	// NoMixedCompression - do not allow compression of an IPv4 section.
	NoMixedCompression MixedCompressionOptions = "no mixed compression"

	// MixedCompressionNoHost - allow compression of the IPv4 section when there is no host of the prefixed address.
	MixedCompressionNoHost MixedCompressionOptions = "no host"

	// MixedCompressionCoveredByHost - compress the IPv4 section if it is part of the host of the prefixed address.
	MixedCompressionCoveredByHost MixedCompressionOptions = "covered by host"

	// AllowMixedCompression - allow compression of a the IPv4 section.
	AllowMixedCompression MixedCompressionOptions = ""
)

// CompressOptions specifies how to compress the zero-segments in an address or subnet string.
type CompressOptions interface {
	// GetCompressionChoiceOptions provides the CompressionChoiceOptions which specify which zero-segments should be compressed.
	GetCompressionChoiceOptions() CompressionChoiceOptions

	// GetMixedCompressionOptions provides the MixedCompressionOptions which specify which zero-segments should be compressed in mixed IPv6/v4 strings.
	GetMixedCompressionOptions() MixedCompressionOptions

	// CompressSingle indicates if a single zero-segment should be compressed on its own when there are no other segments to compress.
	CompressSingle() bool
}

type compressOptions struct {
	compressSingle bool

	rangeSelection CompressionChoiceOptions

	//options for addresses with an ipv4 section
	compressMixedOptions MixedCompressionOptions
}

// GetCompressionChoiceOptions provides the CompressionChoiceOptions which specify which zero-segments should be compressed.
func (opts *compressOptions) GetCompressionChoiceOptions() CompressionChoiceOptions {
	return opts.rangeSelection
}

// GetMixedCompressionOptions provides the MixedCompressionOptions which specify which zero-segments should be compressed in mixed IPv6/v4 strings.
func (opts *compressOptions) GetMixedCompressionOptions() MixedCompressionOptions {
	return opts.compressMixedOptions
}

// CompressSingle indicates if a single zero-segment should be compressed on its own when there are no other segments to compress.
func (opts *compressOptions) CompressSingle() bool {
	return opts.compressSingle
}

var _ CompressOptions = &compressOptions{}

// CompressOptionsBuilder is used to build an immutable CompressOptions instance for IPv6 address strings.
type CompressOptionsBuilder struct {
	compressOptions
}

// SetCompressSingle dictates whether a single zero-segment should be compressed on its own when there are no other segments to compress
func (builder *CompressOptionsBuilder) SetCompressSingle(compressSingle bool) *CompressOptionsBuilder {
	builder.compressSingle = compressSingle
	return builder
}

// SetCompressionChoiceOptions sets the CompressionChoiceOptions which specify which zero-segments should be compressed
func (builder *CompressOptionsBuilder) SetCompressionChoiceOptions(rangeSelection CompressionChoiceOptions) *CompressOptionsBuilder {
	builder.rangeSelection = rangeSelection
	return builder
}

// SetMixedCompressionOptions sets the MixedCompressionOptions which specify which zero-segments should be compressed in mixed IPv6/v4 strings
func (builder *CompressOptionsBuilder) SetMixedCompressionOptions(compressMixedOptions MixedCompressionOptions) *CompressOptionsBuilder {
	builder.compressMixedOptions = compressMixedOptions
	return builder
}

// ToOptions returns an immutable CompressOptions instance built by this builder
func (builder *CompressOptionsBuilder) ToOptions() CompressOptions {
	res := builder.compressOptions
	return &res
}
