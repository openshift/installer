//
// Copyright 2020-2023 Sean C Foley
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
	"fmt"
	"net"
	"net/netip"
	"strings"
	"unsafe"

	"github.com/seancfoley/ipaddress-go/ipaddr/addrerr"
	"github.com/seancfoley/ipaddress-go/ipaddr/addrstrparam"
)

const (
	PortSeparator    = ':'
	LabelSeparator   = '.'
	IPv6StartBracket = '['
	IPv6EndBracket   = ']'
)

func parseHostName(str string, params addrstrparam.HostNameParams) *HostName {
	str = strings.TrimSpace(str)
	res := &HostName{
		str:       str,
		hostCache: &hostCache{},
	}
	res.validate(params)
	return res
}

// NewHostName constructs a HostName that will parse the given string according to the default parameters.
func NewHostName(str string) *HostName {
	return parseHostName(str, defaultHostParameters)
}

// NewHostNameParams constructs a HostName that will parse the given string according to the given parameters.
func NewHostNameParams(str string, params addrstrparam.HostNameParams) *HostName {
	var prms addrstrparam.HostNameParams
	if params == nil {
		prms = defaultHostParameters
	} else {
		prms = addrstrparam.CopyHostNameParams(params)
	}
	return parseHostName(str, prms)
}

// NewHostNameFromAddrPort constructs a HostName from an IP address and a port.
func NewHostNameFromAddrPort(addr *IPAddress, port uint16) *HostName {
	portVal := PortInt(port)
	hostStr := toNormalizedAddrPortString(addr, portVal)
	parsedHost := parsedHost{
		originalStr:     hostStr,
		embeddedAddress: embeddedAddress{addressProvider: addr.getProvider()},
		labelsQualifier: parsedHostIdentifierStringQualifier{port: cachePorts(portVal)},
	}
	return &HostName{
		str:        hostStr,
		hostCache:  &hostCache{normalizedString: &hostStr},
		parsedHost: &parsedHost,
	}
}

// NewHostNameFromAddr constructs a HostName from an IP address.
func NewHostNameFromAddr(addr *IPAddress) *HostName {
	hostStr := addr.ToNormalizedString()
	return newHostNameFromAddr(hostStr, addr)
}

func newHostNameFromAddr(hostStr string, addr *IPAddress) *HostName { // same as HostName(String hostStr, ParsedHost parsed) {
	parsedHost := parsedHost{
		originalStr:     hostStr,
		embeddedAddress: embeddedAddress{addressProvider: addr.getProvider()},
	}
	return &HostName{
		str:        hostStr,
		hostCache:  &hostCache{normalizedString: &hostStr},
		parsedHost: &parsedHost,
	}
}

// NewHostNameFromNetTCPAddr constructs a HostName from a net.TCPAddr.
func NewHostNameFromNetTCPAddr(addr *net.TCPAddr) (*HostName, addrerr.AddressValueError) {
	return newHostNameFromSocketAddr(addr.IP, addr.Port, addr.Zone)
}

// NewHostNameFromNetUDPAddr constructs a HostName from a net.UDPAddr.
func NewHostNameFromNetUDPAddr(addr *net.UDPAddr) (*HostName, addrerr.AddressValueError) {
	return newHostNameFromSocketAddr(addr.IP, addr.Port, addr.Zone)
}

func newHostNameFromSocketAddr(ip net.IP, port int, zone string) (hostName *HostName, err addrerr.AddressValueError) {
	var ipAddr *IPAddress
	ipAddr, err = NewIPAddressFromNetIPAddr(&net.IPAddr{IP: ip, Zone: zone})
	if err != nil {
		return
	} else if ipAddr == nil {
		err = &addressValueError{addressError: addressError{key: "ipaddress.error.exceeds.size"}}
		return
	}
	portVal := PortInt(port)
	hostStr := toNormalizedAddrPortString(ipAddr, portVal)
	parsedHost := parsedHost{
		originalStr:     hostStr,
		embeddedAddress: embeddedAddress{addressProvider: ipAddr.getProvider()},
		labelsQualifier: parsedHostIdentifierStringQualifier{port: cachePorts(portVal)},
	}
	hostName = &HostName{
		str:        hostStr,
		hostCache:  &hostCache{normalizedString: &hostStr},
		parsedHost: &parsedHost,
	}
	return
}

// NewHostNameFromNetIP constructs a HostName from a net.IP.
func NewHostNameFromNetIP(bytes net.IP) (hostName *HostName, err addrerr.AddressValueError) {
	var addr *IPAddress
	addr, err = NewIPAddressFromNetIP(bytes)
	if err != nil {
		return
	} else if addr == nil {
		err = &addressValueError{addressError: addressError{key: "ipaddress.error.exceeds.size"}}
		return
	}
	hostName = NewHostNameFromAddr(addr)
	return
}

// NewHostNameFromPrefixedNetIP constructs a HostName from a net.IP paired with a prefix length.
func NewHostNameFromPrefixedNetIP(bytes net.IP, prefixLen PrefixLen) (hostName *HostName, err addrerr.AddressValueError) {
	var addr *IPAddress
	addr, err = NewIPAddressFromPrefixedNetIP(bytes, prefixLen)
	if err != nil {
		return
	} else if addr == nil {
		err = &addressValueError{addressError: addressError{key: "ipaddress.error.exceeds.size"}}
		return
	}

	hostName = NewHostNameFromAddr(addr)
	return
}

// NewHostNameFromNetIPAddr constructs a HostName from a net.IPAddr.
func NewHostNameFromNetIPAddr(addr *net.IPAddr) (hostName *HostName, err addrerr.AddressValueError) {
	var ipAddr *IPAddress
	ipAddr, err = NewIPAddressFromNetIPAddr(addr)
	if err != nil {
		return
	} else if ipAddr == nil {
		err = &addressValueError{addressError: addressError{key: "ipaddress.error.exceeds.size"}}
		return
	}
	hostName = NewHostNameFromAddr(ipAddr)
	return
}

// NewHostNameFromPrefixedNetIPAddr constructs a HostName from a net.IPAddr paired with a prefix length.
func NewHostNameFromPrefixedNetIPAddr(addr *net.IPAddr, prefixLen PrefixLen) (hostName *HostName, err addrerr.AddressValueError) {
	var ipAddr *IPAddress
	ipAddr, err = NewIPAddressFromPrefixedNetIPAddr(addr, prefixLen)
	if err != nil {
		return
	} else if ipAddr == nil {
		err = &addressValueError{addressError: addressError{key: "ipaddress.error.exceeds.size"}}
		return
	}
	hostName = NewHostNameFromAddr(ipAddr)
	return
}

// NewHostNameFromNetNetIPAddr constructs a host name from a netip.Addr.
func NewHostNameFromNetNetIPAddr(addr netip.Addr) *HostName {
	ipAddr := NewIPAddressFromNetNetIPAddr(addr)
	return NewHostNameFromAddr(ipAddr)
}

// NewHostNameFromNetNetIPPrefix constructs a host name from a netip.Prefix.
func NewHostNameFromNetNetIPPrefix(addr netip.Prefix) (hostName *HostName, err addrerr.AddressValueError) {
	var ipAddr *IPAddress
	ipAddr, err = NewIPAddressFromNetNetIPPrefix(addr)
	if err == nil {
		hostName = NewHostNameFromAddr(ipAddr)
	}
	return
}

// NewHostNameFromNetNetIPAddrPort constructs a host name from a netip.AddrPort.
func NewHostNameFromNetNetIPAddrPort(addrPort netip.AddrPort) *HostName {
	port := addrPort.Port()
	addr := addrPort.Addr()
	ipAddr := NewIPAddressFromNetNetIPAddr(addr)
	return NewHostNameFromAddrPort(ipAddr, port)
}

var defaultHostParameters = new(addrstrparam.HostNameParamsBuilder).ToParams()

var zeroHost = NewHostName("")

type resolveData struct {
	resolvedAddrs []*IPAddress
	err           error
}

type hostCache struct {
	resolveData      *resolveData
	normalizedString *string
}

// HostName represents an internet host name.  Can be a fully qualified domain name, a simple host name, or an ip address string.
// It can also include a port number or service name (which maps to a port).
// It can include a prefix length or mask for either an ipaddress or host name string.  An IPv6 address can have an IPv6 zone.
//
// # Supported Formats
//
// You can use all host or address formats supported by nmap and all address formats supported by IPAddressString.
// All manners of domain names are supported. When adding a prefix length or mask to a host name string, it is to denote the subnet of the resolved address.
//
// Validation is done separately from DNS resolution to avoid unnecessary DNS lookups.
//
// See RFC 3513, RFC 2181, RFC 952, RFC 1035, RFC 1034, RFC 1123, RFC 5890 or the list of rfcs for IPAddress.  For IPv6 addresses in host, see RFC 2732 specifying "[]" notation
// and RFC 3986 and RFC 4038 (combining IPv6 "[]" notation with prefix or zone) and SMTP RFC 2821 for alternative uses of "[]" notation for both IPv4 and IPv6.
type HostName struct {
	str           string
	parsedHost    *parsedHost
	validateError addrerr.HostNameError
	*hostCache
}

func (host *HostName) init() *HostName {
	if host.parsedHost == nil && host.validateError == nil { // the only way params can be nil is when str == "" as well
		return zeroHost
	}
	return host
}

// GetValidationOptions returns the validation options supplied when constructing the HostName, or the default validation options if none were supplied.
// It returns nil if no options were used to construct.
func (host *HostName) GetValidationOptions() addrstrparam.HostNameParams {
	return host.init().parsedHost.params
}

func (host *HostName) validate(validationOptions addrstrparam.HostNameParams) {
	parsed, validateError := validator.validateHostName(host, validationOptions)
	if validateError != nil && parsed == nil {
		parsed = &parsedHost{originalStr: host.str, params: validationOptions}
	}
	host.parsedHost, host.validateError = parsed, validateError
}

// Validate validates that this string is a valid address, and if not, returns an error with a descriptive message indicating why it is not.
func (host *HostName) Validate() addrerr.HostNameError {
	return host.init().validateError
}

// String implements the [fmt.Stringer] interface,
// returning the original string used to create this HostName (altered by strings.TrimSpace if a host name and not an address),
// or "<nil>" if the receiver is a nil pointer.
func (host *HostName) String() string {
	if host == nil {
		return nilString()
	}
	return host.str
}

// Format implements the [fmt.Formatter] interface.
// It accepts the verbs hat are applicable to strings,
// namely the verbs %s, %q, %x and %X.
func (addrStr HostName) Format(state fmt.State, verb rune) {
	s := flagsFromState(state, verb)
	_, _ = state.Write([]byte(fmt.Sprintf(s, addrStr.str)))
}

// IsAddressString returns whether this host name is a string representing an IP address or subnet.
func (host *HostName) IsAddressString() bool {
	host = host.init()
	return host.IsValid() && host.parsedHost.isAddressString()
}

// IsAddress returns whether this host name is a string representing a valid specific IP address or subnet.
func (host *HostName) IsAddress() bool {
	if host.IsAddressString() {
		addr, _ := host.init().parsedHost.asAddress()
		return addr != nil
	}
	return false
}

// AsAddress returns the address if this host name represents an ip address.  Otherwise, this returns nil.
// Note that the translation includes prefix lengths and IPv6 zones.
//
// This does not resolve addresses or return resolved addresses.
// Call ToAddress or GetAddress to get the resolved address.
//
// In cases such as IPv6 literals and reverse-DNS hosts, you can check the relevant methods isIpv6Literal or isReverseDNS,
// in which case this method should return the associated address.
func (host *HostName) AsAddress() *IPAddress {
	if host.IsAddress() {
		addr, _ := host.parsedHost.asAddress()
		return addr
	}
	return nil
}

// IsAllAddresses returns whether this is an IP address that represents the set all all valid IP addresses (as opposed to an empty string, a specific address, or an invalid format).
func (host *HostName) IsAllAddresses() bool {
	host = host.init()
	return host.IsValid() && host.parsedHost.getAddressProvider().isProvidingAllAddresses()
}

// IsEmpty returns true if the host name is empty (zero-length).
func (host *HostName) IsEmpty() bool {
	host = host.init()
	return host.IsValid() &&
		((host.IsAddressString() &&
			host.parsedHost.getAddressProvider().isProvidingEmpty()) || len(host.GetNormalizedLabels()) == 0)
}

// GetAddress attempts to convert this host name to an IP address.
// If this represents an ip address, returns that address.
// If this represents a host, returns the resolved ip address of that host.
// Otherwise, returns nil.
// GetAddress is similar to ToAddress but does not return any errors.
//
// If you wish to get the represented address while avoiding DNS resolution, use AsAddress or AsAddressString.
func (host *HostName) GetAddress() *IPAddress {
	addr, _ := host.ToAddress()
	return addr
}

// ToAddress resolves to an address.
// This method can potentially return a list of resolved addresses and an error as well, if some resolved addresses were invalid.
func (host *HostName) ToAddress() (addr *IPAddress, err addrerr.AddressError) {
	addresses, err := host.ToAddresses()
	if len(addresses) > 0 {
		addr = addresses[0]
	}
	return
}

// ToAddresses resolves to one or more addresses.
// The error can be addrerr.AddressStringError,addrerr.IncompatibleAddressError, or addrerr.HostNameError.
// This method can potentially return a list of resolved addresses and an error as well if some resolved addresses were invalid.
func (host *HostName) ToAddresses() (addrs []*IPAddress, err addrerr.AddressError) {
	host = host.init()
	data := (*resolveData)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&host.resolveData))))
	if data == nil {
		//note that validation handles empty address resolution
		err = host.Validate() //addrerr.HostNameError
		if err != nil {
			return
		}
		// http://networkbit.ch/golang-dns-lookup/
		parsedHost := host.parsedHost
		if parsedHost.isAddressString() {
			addr, addrErr := parsedHost.asAddress() //addrerr.IncompatibleAddressError
			addrs, err = []*IPAddress{addr}, addrErr
			//note there is no need to apply prefix or mask here, it would have been applied to the address already
		} else {
			strHost := parsedHost.getHost()
			validationOptions := host.GetValidationOptions()
			if len(strHost) == 0 {
				addrs = []*IPAddress{}
			} else {
				var ips []net.IP
				ips, lookupErr := net.LookupIP(strHost)
				if lookupErr != nil {
					//Note we do not set resolveData, so we will attempt to resolve again
					err = &hostNameNestedError{nested: lookupErr,
						hostNameError: hostNameError{addressError{str: strHost, key: "ipaddress.host.error.host.resolve"}}}
					return
				}
				count := len(ips)
				addrs = make([]*IPAddress, 0, count)
				var errs []addrerr.AddressError
				for j := 0; j < count; j++ {
					ip := ips[j]
					if ipv4 := ip.To4(); ipv4 != nil {
						ip = ipv4
					}
					networkPrefixLength := parsedHost.getNetworkPrefixLen()
					byteLen := len(ip)
					if networkPrefixLength == nil {
						mask := parsedHost.getMask()
						if mask != nil {
							maskBytes := mask.Bytes()
							if len(maskBytes) == byteLen {
								for i := 0; i < byteLen; i++ {
									ip[i] &= maskBytes[i]
								}
								networkPrefixLength = mask.GetBlockMaskPrefixLen(true)
							}
						}
					}
					ipAddr, addrErr := NewIPAddressFromPrefixedNetIP(ip, networkPrefixLength)
					if addrErr != nil {
						errs = append(errs, addrErr)
					} else {
						cache := ipAddr.cache
						if cache != nil {
							cache.identifierStr = &identifierStr{host}
						}
						addrs = append(addrs, ipAddr)
					}
				}
				if len(errs) > 0 {
					err = &mergedError{AddressError: &hostNameError{addressError{str: strHost, key: "ipaddress.host.error.host.resolve"}}, merged: errs}
				}
				count = len(addrs)
				if count > 0 {
					// sort by preferred version
					preferredVersion := IPVersion(validationOptions.GetPreferredVersion())
					if !preferredVersion.IsIndeterminate() {
						preferredAddrType := preferredVersion.toType()
						boundaryCase := 8 // we sort differently based on list size
						if count > boundaryCase {
							c := 0
							newAddrs := make([]*IPAddress, count)
							for _, val := range addrs {
								if val.getAddrType() == preferredAddrType {
									newAddrs[c] = val
									c++
								}
							}
							for i := 0; c < count; i++ {
								val := addrs[i]
								if val.getAddrType() != preferredAddrType {
									newAddrs[c] = val
									c++
								}
							}
							addrs = newAddrs
						} else {
							preferredIndex := 0
						top:
							for i := 0; i < count; i++ {
								val := addrs[i]
								if val.getAddrType() != preferredAddrType {
									var j int
									if preferredIndex == 0 {
										j = i + 1
									} else {
										j = preferredIndex
									}
									for ; j < len(addrs); j++ {
										if addrs[j].getAddrType() == preferredAddrType {
											// move the preferred into the non-preferred's spot
											addrs[i] = addrs[j]
											// don't swap so the non-preferred order is preserved,
											// instead shift each upwards by one spot
											k := i + 1
											for ; k < j; k++ {
												addrs[k], val = val, addrs[k]
											}
											addrs[k] = val
											preferredIndex = j + 1
											continue top
										}
									}
									// no more preferred so nothing more to do
									break
								}
							}
						}
					}
				}
			}
		}
		data = &resolveData{addrs, err}
		dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&host.resolveData))
		atomicStorePointer(dataLoc, unsafe.Pointer(data))
	}
	return data.resolvedAddrs, nil
}

// IsValid returns whether this represents a valid host name or IP address format.
func (host *HostName) IsValid() bool {
	return host.init().Validate() == nil
}

// AsAddressString returns the address string if this host name represents an ip address or an ip address string.  Otherwise, this returns nil.
// Note that translation includes prefix lengths and IPv6 zones.
// This does not resolve host names.  Call ToAddress or GetAddress to get the resolved address.
func (host *HostName) AsAddressString() *IPAddressString {
	host = host.init()
	if host.IsAddressString() {
		return host.parsedHost.asGenericAddressString()
	}
	return nil
}

// GetPort returns the port if a port was supplied, otherwise it returns nil.
func (host *HostName) GetPort() Port {
	host = host.init()
	if host.IsValid() {
		return host.parsedHost.getPort().copy()
	}
	return nil
}

// GetService returns the service name if a service name was supplied (which is typically mapped to a port), otherwise it returns an empty string.
func (host *HostName) GetService() string {
	host = host.init()
	if host.IsValid() {
		return host.parsedHost.getService()
	}
	return ""
}

// ToNormalizedString provides a normalized string which is lowercase for host strings, and which is the normalized string for addresses.
func (host *HostName) ToNormalizedString() string {
	if str := host.normalizedString; str != nil {
		return *str
	}
	return host.toNormalizedString(false, false)
}

// ToNormalizedWildcardString provides a normalized string which is lowercase for host strings, and which is a normalized string for addresses.
func (host *HostName) ToNormalizedWildcardString() string {
	return host.toNormalizedString(false, false)
}

// ToQualifiedString provides a normalized string which is lowercase for host strings, and which is a normalized string for addresses.
func (host *HostName) ToQualifiedString() string {
	return host.toNormalizedString(false, true)
}

func (host *HostName) toNormalizedString(wildcard, addTrailingDot bool) string {
	if host.IsValid() {
		var builder strings.Builder
		if host.IsAddress() {
			toNormalizedHostString(host.AsAddress(), wildcard, &builder)
		} else if host.IsAddressString() {
			builder.WriteString(host.AsAddressString().ToNormalizedString())
		} else {
			builder.WriteString(host.parsedHost.getHost())
			if addTrailingDot {
				builder.WriteByte(LabelSeparator)
			}
			/*
			 * If prefix or mask is supplied and there is an address, it is applied directly to the address provider, so
			 * we need only check for those things here
			 *
			 * Also note that ports and prefix/mask cannot appear at the same time, so this does not interfere with the port code below.
			 */
			networkPrefixLength := host.parsedHost.getEquivalentPrefixLen()
			if networkPrefixLength != nil {
				builder.WriteByte(PrefixLenSeparator)
				toUnsignedString(uint64(networkPrefixLength.bitCount()), 10, &builder)
			} else {
				mask := host.parsedHost.getMask()
				if mask != nil {
					builder.WriteByte(PrefixLenSeparator)
					builder.WriteString(mask.ToNormalizedString())
				}
			}
		}
		port := host.parsedHost.getPort()
		if port != nil {
			toNormalizedPortString(port.portNum(), &builder)
		} else {
			service := host.parsedHost.getService()
			if service != "" {
				builder.WriteByte(PortSeparator)
				builder.WriteString(service)
			}
		}
		return builder.String()
	}
	return host.str
}

func toNormalizedPortString(port PortInt, builder *strings.Builder) {
	builder.WriteByte(PortSeparator)
	toUnsignedString(uint64(port), 10, builder)
}

func toNormalizedHostString(addr *IPAddress, wildcard bool, builder *strings.Builder) {
	if addr.isIPv6() {
		if !wildcard && addr.IsPrefixed() { // prefix needs to be outside the brackets
			normalized := addr.ToNormalizedString()
			index := strings.IndexByte(normalized, PrefixLenSeparator)
			builder.WriteByte(IPv6StartBracket)
			translateReserved(addr.ToIPv6(), normalized[:index], builder)
			builder.WriteByte(IPv6EndBracket)
			builder.WriteString(normalized[index:])
		} else {
			normalized := addr.ToNormalizedWildcardString()
			builder.WriteByte(IPv6StartBracket)
			translateReserved(addr.ToIPv6(), normalized, builder)
			builder.WriteByte(IPv6EndBracket)
		}
	} else {
		if wildcard {
			builder.WriteString(addr.ToNormalizedWildcardString())
		} else {
			builder.WriteString(addr.ToNormalizedString())
		}
	}
}

func toNormalizedAddrPortString(addr *IPAddress, port PortInt) string {
	builder := strings.Builder{}
	toNormalizedHostString(addr, false, &builder)
	toNormalizedPortString(port, &builder)
	return builder.String()
}

// Equal returns true if the given host name matches this one.
// For hosts to match, they must represent the same addresses or have the same host names.
// Hosts are not resolved when matching.  Also, hosts must have the same port or service.
// They must have the same masks if they are host names.
// Even if two hosts are invalid, they match if they have the same invalid string.
func (host *HostName) Equal(other *HostName) bool {
	if host == nil {
		return other == nil
	} else if other == nil {
		return false
	}
	host = host.init()
	other = other.init()
	if host == other {
		return true
	}
	if host.IsValid() {
		if other.IsValid() {
			parsedHost := host.parsedHost
			otherParsedHost := other.parsedHost
			if parsedHost.isAddressString() {
				return otherParsedHost.isAddressString() &&
					parsedHost.asGenericAddressString().Equal(otherParsedHost.asGenericAddressString()) &&
					parsedHost.getPort().Equal(otherParsedHost.getPort()) &&
					parsedHost.getService() == otherParsedHost.getService()
			}
			if otherParsedHost.isAddressString() {
				return false
			}
			thisHost := parsedHost.getHost()
			otherHost := otherParsedHost.getHost()
			if thisHost != otherHost {
				return false
			}
			return parsedHost.getEquivalentPrefixLen().Equal(otherParsedHost.getEquivalentPrefixLen()) &&
				parsedHost.getMask().Equal(otherParsedHost.getMask()) &&
				parsedHost.getPort().Equal(otherParsedHost.getPort()) &&
				parsedHost.getService() == otherParsedHost.getService()
		}
		return false
	}
	return !other.IsValid() && host.String() == other.String()
}

// GetNormalizedLabels returns an array of normalized strings for this host name instance.
//
// If this represents an IP address, the address segments are separated into the returned array.
// If this represents a host name string, the domain name segments are separated into the returned array,
// with the top-level domain name (right-most segment) as the last array element.
//
// The individual segment strings are normalized in the same way as ToNormalizedString.
//
// Ports, service name strings, prefix lengths, and masks are all omitted from the returned array.
func (host *HostName) GetNormalizedLabels() []string {
	host = host.init()
	if host.IsValid() {
		return host.parsedHost.getNormalizedLabels()
	} else {
		str := host.str
		if len(str) == 0 {
			return []string{}
		}
		return []string{str}
	}
}

// GetHost returns the host string normalized but without port, service, prefix or mask.
//
// If an address, returns the address string normalized, but without port, service, prefix, mask, or brackets for IPv6.
//
// To get a normalized string encompassing all details, use ToNormalizedString.
//
// If not a valid host, returns the zero string.
func (host *HostName) GetHost() string {
	host = host.init()
	if host.IsValid() {
		return host.parsedHost.getHost()
	}
	return ""
}

// IsUncIPv6Literal returns whether this host name is an Uniform Naming Convention IPv6 literal host name.
func (host *HostName) IsUncIPv6Literal() bool {
	host = host.init()
	return host.IsValid() && host.parsedHost.isUNCIPv6Literal()
}

// IsReverseDNS returns whether this host name is a reverse-DNS string host name.
func (host *HostName) IsReverseDNS() bool {
	host = host.init()
	return host.IsValid() && host.parsedHost.isReverseDNS()
}

// GetNetworkPrefixLen returns the prefix length, if a prefix length was supplied,
// either as part of an address or as part of a domain (in which case the prefix applies to any resolved address).
// Otherwise, GetNetworkPrefixLen returns nil.
func (host *HostName) GetNetworkPrefixLen() PrefixLen {
	if host.IsAddress() {
		addr, err := host.parsedHost.asAddress()
		if err == nil {
			return addr.getNetworkPrefixLen().copy()
		}
	} else if host.IsAddressString() {
		return host.parsedHost.asGenericAddressString().getNetworkPrefixLen().copy()
	} else if host.IsValid() {
		return host.parsedHost.getEquivalentPrefixLen().copy()
	}
	return nil
}

// GetMask returns the resulting mask value if a mask was provided with this host name.
func (host *HostName) GetMask() *IPAddress {
	if host.IsValid() {
		if host.parsedHost.isAddressString() {
			return host.parsedHost.getAddressProvider().getProviderMask()
		}
		return host.parsedHost.getMask()
	}
	return nil
}

// ResolvesToSelf returns whether this represents, or resolves to,
// a host or address representing the same host.
func (host *HostName) ResolvesToSelf() bool {
	if host.IsSelf() {
		return true
	} else if host.GetAddress() != nil {
		host.resolveData.resolvedAddrs[0].IsLoopback()
	}
	return false
}

// IsSelf returns whether this represents a host or address representing the same host.
// Also see IsLocalHost and IsLoopback.
func (host *HostName) IsSelf() bool {
	return host.IsLocalHost() || host.IsLoopback()
}

// IsLocalHost returns whether this host is "localhost".
func (host *HostName) IsLocalHost() bool {
	return host.IsValid() && strings.EqualFold(host.str, "localhost")
}

// IsLoopback returns whether this host has the loopback address, such as "::1" or "127.0.0.1".
//
// Also see IsSelf.
func (host *HostName) IsLoopback() bool {
	return host.IsAddress() && host.AsAddress().IsLoopback()
}

// ToNetTCPAddrService returns the TCPAddr if this HostName both resolves to an address and has an associated service or port, otherwise returns nil.
func (host *HostName) ToNetTCPAddrService(serviceMapper func(string) Port) *net.TCPAddr {
	if host.IsValid() {
		port := host.GetPort()
		if port == nil && serviceMapper != nil {
			service := host.GetService()
			if service != "" {
				port = serviceMapper(service)
			}
		}
		if port != nil {
			if addr := host.GetAddress(); addr != nil {
				return &net.TCPAddr{
					IP:   addr.GetNetIP(),
					Port: port.portNum(),
					Zone: string(addr.zone),
				}
			}
		}
	}
	return nil
}

// ToNetTCPAddr returns the TCPAddr if this HostName both resolves to an address and has an associated port.
// Otherwise, it returns nil.
func (host *HostName) ToNetTCPAddr() *net.TCPAddr {
	return host.ToNetTCPAddrService(nil)
}

// ToNetUDPAddrService returns the UDPAddr if this HostName both resolves to an address and has an associated service or port.
func (host *HostName) ToNetUDPAddrService(serviceMapper func(string) Port) *net.UDPAddr {
	tcpAddr := host.ToNetTCPAddrService(serviceMapper)
	if tcpAddr != nil {
		return &net.UDPAddr{
			IP:   tcpAddr.IP,
			Port: tcpAddr.Port,
			Zone: tcpAddr.Zone,
		}
	}
	return nil
}

// ToNetUDPAddr returns the UDPAddr if this HostName both resolves to an address and has an associated port.
func (host *HostName) ToNetUDPAddr(serviceMapper func(string) Port) *net.UDPAddr {
	return host.ToNetUDPAddrService(serviceMapper)
}

// ToNetIP is similar to ToAddress but returns the resulting address as a net.IP.
func (host *HostName) ToNetIP() net.IP {
	if addr, err := host.ToAddress(); addr != nil && err == nil {
		return addr.GetNetIP()
	}
	return nil
}

// ToNetIPAddr is similar to ToAddress but returns the resulting address as a net.IPAddr.
func (host *HostName) ToNetIPAddr() *net.IPAddr {
	if addr, err := host.ToAddress(); addr != nil && err == nil {
		return &net.IPAddr{
			IP:   addr.GetNetIP(),
			Zone: string(addr.zone),
		}
	}
	return nil
}

// Compare returns a negative integer, zero, or a positive integer if this host name is less than, equal, or greater than the given host name.
// Any address item is comparable to any other.
func (host *HostName) Compare(other *HostName) int {
	if host == other {
		return 0
	} else if host == nil {
		return -1
	} else if other == nil {
		return 1
	}
	if host.IsValid() {
		if other.IsValid() {
			parsedHost := host.parsedHost
			otherParsedHost := other.parsedHost
			if parsedHost.isAddressString() {
				if otherParsedHost.isAddressString() {
					result := parsedHost.asGenericAddressString().Compare(otherParsedHost.asGenericAddressString())
					if result != 0 {
						return result
					}
					//fall through to compare ports
				} else {
					return -1
				}
			} else if otherParsedHost.isAddressString() {
				return 1
			} else {
				//both are non-address hosts
				normalizedLabels := parsedHost.getNormalizedLabels()
				otherNormalizedLabels := otherParsedHost.getNormalizedLabels()
				oneLen := len(normalizedLabels)
				twoLen := len(otherNormalizedLabels)
				var minLen int
				if oneLen < twoLen {
					minLen = oneLen
				} else {
					minLen = twoLen
				}
				for i := 1; i <= minLen; i++ {
					one := normalizedLabels[oneLen-i]
					two := otherNormalizedLabels[twoLen-i]
					result := strings.Compare(one, two)
					if result != 0 {
						return result
					}
				}
				if oneLen != twoLen {
					return oneLen - twoLen
				}

				//keep in mind that hosts can has masks/prefixes or ports, but not both
				networkPrefixLength := parsedHost.getEquivalentPrefixLen()
				otherPrefixLength := otherParsedHost.getEquivalentPrefixLen()
				if networkPrefixLength != nil {
					if otherPrefixLength != nil {
						if *networkPrefixLength != *otherPrefixLength {
							return otherPrefixLength.bitCount() - networkPrefixLength.bitCount()
						}
						//fall through to compare ports
					} else {
						return 1
					}
				} else {
					if otherPrefixLength != nil {
						return -1
					}
					mask := parsedHost.getMask()
					otherMask := otherParsedHost.getMask()
					if mask != nil {
						if otherMask != nil {
							ret := mask.Compare(otherMask)
							if ret != 0 {
								return ret
							}
							//fall through to compare ports
						} else {
							return 1
						}
					} else {
						if otherMask != nil {
							return -1
						}
						//fall through to compare ports
					}
				} //end non-address host compare
			}

			//two equivalent address strings or two equivalent hosts, now check port and service names
			portOne := parsedHost.getPort()
			portTwo := otherParsedHost.getPort()
			portRet := portOne.Compare(portTwo)
			if portRet != 0 {
				return portRet
			}
			serviceOne := parsedHost.getService()
			serviceTwo := otherParsedHost.getService()
			if serviceOne != "" {
				if serviceTwo != "" {
					ret := strings.Compare(serviceOne, serviceTwo)
					if ret != 0 {
						return ret
					}
				} else {
					return 1
				}
			} else if serviceTwo != "" {
				return -1
			}
			return 0
		} else {
			return 1
		}
	} else if other.IsValid() {
		return -1
	}
	return strings.Compare(host.String(), other.String())
}

// Wrap wraps this host name, returning a WrappedHostName, an implementation of ExtendedIdentifierString,
// which can be used to write code that works with a host identifier string including [IPAddressString], [MACAddressString], and [HostName].
func (host *HostName) Wrap() ExtendedIdentifierString {
	return WrappedHostName{host}
}

func translateReserved(addr *IPv6Address, str string, builder *strings.Builder) {
	//This is particularly targeted towards the zone
	if !addr.HasZone() {
		builder.WriteString(str)
		return
	}
	index := strings.IndexByte(str, IPv6ZoneSeparator)
	var translated = builder
	translated.WriteString(str[0:index])
	translated.WriteString("%25")
	for i := index + 1; i < len(str); i++ {
		c := str[i]
		if isReserved(c) {
			translated.WriteByte('%')
			toUnsignedString(uint64(c), 16, translated)
		} else {
			translated.WriteByte(c)
		}
	}
}
